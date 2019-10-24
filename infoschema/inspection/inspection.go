package inspection

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/ddl"
	"github.com/pingcap/tidb/infoschema"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/meta/autoid"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/table"
	"github.com/pingcap/tidb/types"
)

func InitInspectionDB(name string) {
	p := parser.New()
	tbls := make([]*model.TableInfo, 0)
	dbID := autoid.GenLocalSchemaID()

	for _, sql := range inspectionTables {
		stmt, err := p.ParseOneStmt(sql, "", "")
		if err != nil {
			panic(err)
		}
		meta, err := ddl.BuildTableInfoFromAST(stmt.(*ast.CreateTableStmt))
		if err != nil {
			panic(err)
		}
		tbls = append(tbls, meta)
		meta.ID = autoid.GenLocalSchemaID()
		for _, c := range meta.Columns {
			c.ID = autoid.GenLocalSchemaID()
		}
	}
	dbInfo := &model.DBInfo{
		ID:      dbID,
		Name:    model.NewCIStr(name),
		Charset: mysql.DefaultCharset,
		Collate: mysql.DefaultCollationName,
		Tables:  tbls,
	}
	infoschema.RegisterVirtualTable(dbInfo, tableFromMeta)
}

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

