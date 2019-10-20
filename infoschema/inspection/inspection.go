package inspection

import (
	"errors"
	"fmt"
	"time"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/domain"

	"github.com/pingcap/parser/model"
	"github.com/pingcap/tidb/sessionctx"
)

func NewInspectionHelper (ctx sessionctx.Context) *InspectionHelper {
	return &InspectionHelper{
		ctx: ctx,
		p: parser.New(),
		dbName: fmt.Sprintf("%s_%s", "tidb_inspection",time.Now().Format("20060102150405")),
	}
}

type InspectionHelper struct {
	ctx sessionctx.Context
	p *parser.Parser
	dbName string
}

func (i *InspectionHelper) GetDBName() string {
	return i.dbName
}

func (i *InspectionHelper) CreateInspectionDB() error {
	err := domain.GetDomain(i.ctx).DDL().CreateSchema(i.ctx, model.NewCIStr(i.dbName), nil)
	if err != nil {
		return err
	}
	return nil
}

func (i *InspectionHelper) CreateInspectionTables() error {
	// Create inspection tables
	for _, template := range inspectionTables {
		sql := fmt.Sprintf(template, i.dbName)
		stmt, err := i.p.ParseOneStmt(sql, "", "")
		if err != nil {
			return err
		}
		s, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			return errors.New(fmt.Sprintf("Fail to create inspection table. Maybe create table statment is illegal: %s", sql))
		}
		s.Table.TableInfo = &model.TableInfo{IsInspection:true, InspectionInfo: make(map[string]string)}
		if err := domain.GetDomain(i.ctx).DDL().CreateTable(i.ctx, s); err != nil {
			return err
		}
	}
	return nil
}

/* TODO: The Following are inspection tables. They should be memtable like information schemas
func tableFromMeta(alloc autoid.Allocator, meta *model.TableInfo) (table.Table, error) {
	return createInspectionTable(meta), nil
}

// createPerfSchemaTable creates all perfSchemaTables
func createInspectionTable(meta *model.TableInfo) *inspectTable {
	columns := make([]*table.Column, 0, len(meta.Columns))
	for _, colInfo := range meta.Columns {
		col := table.ToColumn(colInfo)
		columns = append(columns, col)
	}
	t := &inspectTable{
		meta: meta,
		cols: columns,
	}
	return t
}

// inspectTable stands for the fake table all its data is in the memory.
type inspectTable struct {
	infoschema.VirtualTable
	meta *model.TableInfo
	cols []*table.Column
}

// Cols implements table.Table Type interface.
func (vt *inspectTable) Cols() []*table.Column {
	return vt.cols
}

// WritableCols implements table.Table Type interface.
func (vt *inspectTable) WritableCols() []*table.Column {
	return vt.cols
}

// GetID implements table.Table GetID interface.
func (vt *inspectTable) GetPhysicalID() int64 {
	return vt.meta.ID
}

// Meta implements table.Table Type interface.
func (vt *inspectTable) Meta() *model.TableInfo {
	return vt.meta
}

func (vt *inspectTable) getRows(ctx sessionctx.Context, cols []*table.Column) (fullRows [][]types.Datum, err error) {
	// switch vt.meta.Name.O {
	// case tableNameEventsStatementsSummaryByDigest:
	// 	fullRows = stmtsummary.StmtSummaryByDigestMap.ToDatum()
	// case tableNameCpuProfile:
	// 	fullRows, err = cpuProfileGraph()
	// case tableNameMemoryProfile:
	// 	fullRows, err = profileGraph("heap")
	// case tableNameMutexProfile:
	// 	fullRows, err = profileGraph("mutex")
	// case tableNameAllocsProfile:
	// 	fullRows, err = profileGraph("allocs")
	// case tableNameBlockProfile:
	// 	fullRows, err = profileGraph("block")
	// case tableNameGoroutines:
	// 	fullRows, err = goroutinesList()
	// }
	// if err != nil {
	// 	return
	// }
	// if len(cols) == len(vt.cols) {
	// 	return
	// }
	rows := make([][]types.Datum, len(fullRows))
	for i, fullRow := range fullRows {
		row := make([]types.Datum, len(cols))
		for j, col := range cols {
			row[j] = fullRow[col.Offset]
		}
		rows[i] = row
	}
	return rows, nil
}

// IterRecords implements table.Table IterRecords interface.
func (vt *inspectTable) IterRecords(ctx sessionctx.Context, startKey kv.Key, cols []*table.Column,
	fn table.RecordIterFunc) error {
	if len(startKey) != 0 {
		return table.ErrUnsupportedOp
	}
	rows, err := vt.getRows(ctx, cols)
	if err != nil {
		return err
	}
	for i, row := range rows {
		more, err := fn(int64(i), row, cols)
		if err != nil {
			return err
		}
		if !more {
			break
		}
	}
	return nil
}

 */
