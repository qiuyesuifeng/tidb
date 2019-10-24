package inspection

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/domain"
	"github.com/pingcap/tidb/domain/infosync"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/store/helper"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/pingcap/tidb/util"
	"github.com/pingcap/tidb/util/sqlexec"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	pmodel "github.com/prometheus/common/model"
)

const promReadTimeout = 10 * time.Second

func NewInspectionHelper(ctx sessionctx.Context) *InspectionHelper {
	return &InspectionHelper{
		ctx:           ctx,
		p:             parser.New(),
		dbName:        fmt.Sprintf("%s_%s", "TIDB_INSPECTION", time.Now().Format("20060102150405")),
		tableNames:    []string{},
		items:         []ClusterItem{},
		nodeExporters: make(map[string]string),
	}
}

type ClusterItem struct {
	ID      int64
	Type    string
	Name    string
	IP      string
	Address string
}

type InspectionHelper struct {
	ctx        sessionctx.Context
	p          *parser.Parser
	dbName     string
	tableNames []string

	items         []ClusterItem
	isInit        bool
	nodeExporters map[string]string
	promClient    api.Client
}

func getIPfromAdress(address string) string {
	return strings.Split(address, ":")[0]
}

func (i *InspectionHelper) GetDBName() string {
	return i.dbName
}

func (i *InspectionHelper) GetTableNames() []string {
	return i.tableNames
}

func (i *InspectionHelper) CreateInspectionDB() error {
	err := domain.GetDomain(i.ctx).DDL().CreateSchema(i.ctx, model.NewCIStr(i.dbName), nil)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (i *InspectionHelper) CreateInspectionTables() error {
	// Create inspection tables
	for _, template := range inspectionVirtualTables {
		sql := fmt.Sprintf(template, i.dbName)
		stmt, err := i.p.ParseOneStmt(sql, "", "")
		if err != nil {
			return errors.Trace(err)
		}

		s, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			return errors.New(fmt.Sprintf("Fail to create inspection table. Maybe create table statment is illegal: %s", sql))
		}

		s.Table.TableInfo = &model.TableInfo{IsInspection: true, InspectionInfo: make(map[string]string)}
		if err := domain.GetDomain(i.ctx).DDL().CreateTable(i.ctx, s); err != nil {
			return errors.Trace(err)
		}

		i.tableNames = append(i.tableNames, s.Table.Name.O)
	}

	for _, template := range inspectionPersistTables {
		sql := fmt.Sprintf(template, i.dbName)
		stmt, err := i.p.ParseOneStmt(sql, "", "")
		if err != nil {
			return errors.Trace(err)
		}

		s, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			return errors.New(fmt.Sprintf("Fail to create inspection table. Maybe create table statment is illegal: %s", sql))
		}
		s.Table.TableInfo = &model.TableInfo{IsInspection: true, InspectionInfo: make(map[string]string)}
		if err := domain.GetDomain(i.ctx).DDL().CreateTable(i.ctx, s); err != nil {
			return errors.Trace(err)
		}

		i.tableNames = append(i.tableNames, s.Table.Name.O)
	}

	return nil
}

func (i *InspectionHelper) TestWriteTable() error {
	sql := fmt.Sprintf("insert into %s.test_persist values (1,1), (2,2);", i.dbName)
	_, _, err := i.ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (i *InspectionHelper) GetClusterInfo() error {
	// get tidb servers info.
	tidbItems, err := infosync.GetAllServerInfo(context.Background())
	if err != nil {
		return errors.Trace(err)
	}

	idx := 0
	for _, item := range tidbItems {
		tp := "tidb"
		name := fmt.Sprintf("tidb-%d", idx)
		tidbAddr := fmt.Sprintf("%s:%d", item.IP, item.Port)
		tidbStatusAddr := fmt.Sprintf("%s:%d", item.IP, item.StatusPort)
		tidbConfig := fmt.Sprintf("http://%s/config", tidbStatusAddr)

		sql := fmt.Sprintf(`insert into %s.TIDB_CLUSTER_INFO values (%d, "%s", "%s", "%s", "%s", "%s", "%s", "%s");`,
			i.dbName, idx, tp, name, tidbAddr, tidbStatusAddr, item.Version, item.GitHash, tidbConfig)

		_, _, err := i.ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
		if err != nil {
			return errors.Trace(err)
		}

		i.items = append(i.items, ClusterItem{int64(idx), tp, name, item.IP, tidbAddr})
		idx++
	}

	// get pd servers info.
	tikvStore, ok := i.ctx.GetStore().(tikv.Storage)
	if !ok {
		return errors.New("Information about TiKV store status can be gotten only when the storage is TiKV")
	}
	tikvHelper := &helper.Helper{
		Store:       tikvStore,
		RegionCache: tikvStore.GetRegionCache(),
	}

	pdHosts, err := tikvHelper.GetPDAddrs()
	if err != nil {
		return errors.Trace(err)
	}
	for ii, host := range pdHosts {
		host = strings.TrimSpace(host)

		// get pd config
		config := fmt.Sprintf("http://%s/pd/api/v1/config", host)

		// get pd version
		url := fmt.Sprintf("http://%s/pd/api/v1/config/cluster-version", host)
		d, err := util.Get(url).Bytes()
		if err != nil {
			return errors.Trace(err)
		}

		version := strings.Trim(strings.Trim(string(d), "\n"), "\"")

		// get pd git_hash
		url = fmt.Sprintf("http://%s/pd/api/v1/status", host)
		dd, err := util.Get(url).Bytes()
		if err != nil {
			return errors.Trace(err)
		}

		m := make(map[string]interface{})
		err = json.Unmarshal(dd, &m)
		if err != nil {
			return errors.Trace(err)
		}

		githash := m["git_hash"]

		tp := "pd"
		name := fmt.Sprintf("pd-%d", ii)
		sql := fmt.Sprintf(`insert into %s.TIDB_CLUSTER_INFO values (%d, "%s", "%s", "%s","%s", "%s", "%s","%s");`,
			i.dbName, idx, tp, name, host, host, version, githash, config)

		_, _, err = i.ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
		if err != nil {
			return errors.Trace(err)
		}

		i.items = append(i.items, ClusterItem{int64(idx), tp, name, getIPfromAdress(host), host})
		idx++
	}

	// get tikv servers info.
	storesStat, err := tikvHelper.GetStoresStat()
	if err != nil {
		return errors.Trace(err)
	}
	for ii, storeStat := range storesStat.Stores {
		tp := "tikv"
		name := fmt.Sprintf("tikv-%d", ii)
		tikvConfig := fmt.Sprintf("http://%s/config", storeStat.Store.StatusAddress)

		sql := fmt.Sprintf(`insert into %s.TIDB_CLUSTER_INFO values (%d, "%s", "%s", "%s", "%s", "%s", "%s", "%s");`,
			i.dbName, idx, tp, name, storeStat.Store.Address, storeStat.Store.StatusAddress, storeStat.Store.Version, storeStat.Store.GitHash, tikvConfig)

		_, _, err := i.ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
		if err != nil {
			return errors.Trace(err)
		}

		i.items = append(i.items, ClusterItem{int64(idx), tp, name, getIPfromAdress(storeStat.Store.Address), storeStat.Store.Address})
		idx++
	}

	i.isInit = true
	return nil
}

func (i *InspectionHelper) initProm() error {
	if !i.isInit {
		return errors.New("InspectionHelper is not init.")
	}

	if i.promClient != nil {
		return nil
	}

	promAddr := config.GetGlobalConfig().PrometheusAddr
	if promAddr == "" {
		return errors.New("Invalid Prometheus Address")
	}

	var err error
	i.promClient, err = api.NewClient(api.Config{
		Address: fmt.Sprintf("http://%s", promAddr),
	})
	if err != nil {
		return errors.Trace(err)
	}

	// get node exporter info.
	api := v1.NewAPI(i.promClient)
	ctx, cancel := context.WithTimeout(context.Background(), promReadTimeout)
	defer cancel()

	targets, err := api.Targets(ctx)
	if err != nil {
		return errors.Trace(err)
	}

	for _, target := range targets.Active {
		if target.Labels["group"] == "node_exporter" {
			neAddr := string(target.Labels["instance"])
			if neAddr != "" {
				i.nodeExporters[getIPfromAdress(neAddr)] = neAddr
			}
		}
	}

	return nil
}

func (i *InspectionHelper) getSystemInfo(item ClusterItem) error {
	api := v1.NewAPI(i.promClient)
	ctx, cancel := context.WithTimeout(context.Background(), promReadTimeout)
	defer cancel()

	neAddr, ok := i.nodeExporters[item.IP]
	if !ok {
		return errors.New("Can not find node exporter address")
	}

	// get cpu count
	cpuCountQuery := fmt.Sprintf(`count(node_cpu_seconds_total{instance="%s", mode="user"})`, neAddr)
	result, err := api.Query(ctx, cpuCountQuery, time.Now())
	if err != nil {
		return errors.Trace(err)
	}
	cpuCount := result.(pmodel.Vector)[0].Value

	// get cpu usage.
	cpuUsageQuery := fmt.Sprintf(`1 - (sum(rate(node_cpu_seconds_total{instance="%s", mode="idle"}[1m])) / count(node_cpu_seconds_total{instance="%s", mode="idle"}) or 
  sum(irate(node_cpu_seconds_total{instance="%s", mode="idle"}[30s])) / count(node_cpu_seconds_total{instance="%s", mode="idle"}))`,
		neAddr, neAddr, neAddr, neAddr)
	result, err = api.Query(ctx, cpuUsageQuery, time.Now())
	if err != nil {
		return errors.Trace(err)
	}
	cpuUsage := fmt.Sprintf("%.2f%%", 100*result.(pmodel.Vector)[0].Value)

	// get total memory.
	memoryQuery := fmt.Sprintf(`node_memory_MemTotal_bytes{instance="%s"}`, neAddr)
	result, err = api.Query(ctx, memoryQuery, time.Now())
	if err != nil {
		return errors.Trace(err)
	}
	memory := fmt.Sprintf("%.2fGiB", result.(pmodel.Vector)[0].Value/1024/1024/1024)

	// get memory usage.
	memoryUsageQuery := fmt.Sprintf(`1 - (node_memory_MemAvailable_bytes{instance="%s"} or 
    (node_memory_MemFree_bytes{instance="%s"} + node_memory_Buffers_bytes{instance="%s"} + node_memory_Cached_bytes{instance="%s"})) / node_memory_MemTotal_bytes{instance="%s"}`,
		neAddr, neAddr, neAddr, neAddr, neAddr)
	result, err = api.Query(ctx, memoryUsageQuery, time.Now())
	if err != nil {
		return errors.Trace(err)
	}
	memoryUsage := fmt.Sprintf("%.2f%%", 100*result.(pmodel.Vector)[0].Value)

	// get load1/load5/load15
	load1Query := fmt.Sprintf(`node_load1{instance="%s"}`, neAddr)
	result, err = api.Query(ctx, load1Query, time.Now())
	if err != nil {
		return errors.Trace(err)
	}
	load1 := fmt.Sprintf("%.2f", result.(pmodel.Vector)[0].Value)

	load5Query := fmt.Sprintf(`node_load5{instance="%s"}`, neAddr)
	result, err = api.Query(ctx, load5Query, time.Now())
	if err != nil {
		return errors.Trace(err)
	}
	load5 := fmt.Sprintf("%.2f", result.(pmodel.Vector)[0].Value)

	load15Query := fmt.Sprintf(`node_load15{instance="%s"}`, neAddr)
	result, err = api.Query(ctx, load15Query, time.Now())
	if err != nil {
		return errors.Trace(err)
	}
	load15 := fmt.Sprintf("%.2f", result.(pmodel.Vector)[0].Value)

	// get kernel version.
	kernelQuery := fmt.Sprintf(`node_uname_info{instance="%s"}`, neAddr)
	result, err = api.Query(ctx, kernelQuery, time.Now())
	if err != nil {
		return errors.Trace(err)
	}
	metric := result.(pmodel.Vector)[0].Metric
	os := metric["sysname"]
	machine := metric["machine"]
	kernelVersion := metric["release"]
	kernel := fmt.Sprintf("%s-%s-%s", os, machine, kernelVersion)

	sql := fmt.Sprintf(`insert into %s.SYSTEM_INFO values (%d, "%s", "%s", "%s", "%s", 
		"%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s");`,
		i.dbName, item.ID, item.Type, item.Name, item.IP, item.Address,
		cpuCount, cpuUsage, memory, memoryUsage, load1, load5, load15, kernel)

	_, _, err = i.ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (i *InspectionHelper) GetSystemInfo() error {
	err := i.initProm()
	if err != nil {
		return errors.Trace(err)
	}

	for _, item := range i.items {
		err = i.getSystemInfo(item)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}
