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
	"github.com/pingcap/tidb/domain"
	"github.com/pingcap/tidb/domain/infosync"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/store/helper"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/pingcap/tidb/util"
	"github.com/pingcap/tidb/util/sqlexec"
)

func NewInspectionHelper(ctx sessionctx.Context) *InspectionHelper {
	return &InspectionHelper{
		ctx:        ctx,
		p:          parser.New(),
		dbName:     fmt.Sprintf("%s_%s", "TIDB_INSPECTION", time.Now().Format("20060102150405")),
		tableNames: []string{},
	}
}

type InspectionHelper struct {
	ctx        sessionctx.Context
	p          *parser.Parser
	dbName     string
	tableNames []string
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
		tidbStatusAddr := fmt.Sprintf("%s:%d", item.IP, item.StatusPort)
		tidbConfig := fmt.Sprintf("http://%s/config", tidbStatusAddr)
		sql := fmt.Sprintf(`insert into %s.TIDB_CLUSTER_INFO values (%d, "tidb", "tidb-%d", "%s:%d", "%s", "%s", "%s", "%s");`,
			i.dbName, idx, idx, item.IP, item.Port, tidbStatusAddr, item.Version, item.GitHash, tidbConfig)

		_, _, err := i.ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
		if err != nil {
			return errors.Trace(err)
		}

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
		sql := fmt.Sprintf(`insert into %s.TIDB_CLUSTER_INFO values (%d, "pd", "pd-%d", "%s","%s", "%s", "%s","%s");`,
			i.dbName, idx, ii, host, host, version, githash, config)

		_, _, err = i.ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
		if err != nil {
			return errors.Trace(err)
		}

		idx++
	}

	// get tikv servers info.
	storesStat, err := tikvHelper.GetStoresStat()
	if err != nil {
		return errors.Trace(err)
	}
	for ii, storeStat := range storesStat.Stores {
		tikvConfig := fmt.Sprintf("http://%s/config", storeStat.Store.StatusAddress)
		sql := fmt.Sprintf(`insert into %s.TIDB_CLUSTER_INFO values (%d, "tikv", "tikv-%d", "%s", "%s", "%s", "%s", "%s");`,
			i.dbName, idx, ii, storeStat.Store.Address, storeStat.Store.StatusAddress, storeStat.Store.Version, storeStat.Store.GitHash, tikvConfig)

		_, _, err := i.ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
		if err != nil {
			return errors.Trace(err)
		}

		idx++
	}

	return nil
}
