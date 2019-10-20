// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package executor

import (
	gojson "encoding/json"
	"strconv"

	"github.com/cznic/mathutil"
	"github.com/pingcap/errors"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/types/json"
	"github.com/pingcap/tidb/util/chunk"
	"github.com/pingcap/tidb/util/mock"

	"golang.org/x/net/context"
)

// make sure `StreamReaderExecutor` implements `Executor`.
var _ Executor = &InspectionReaderExecutor{}

// StreamReaderExecutor reads data from a stream.
type InspectionReaderExecutor struct {
	baseExecutor
	Table   *model.TableInfo
	Columns []*model.ColumnInfo

	result *chunk.Chunk
	cursor int
	pos    int

	tp           string
	topic        string
	variableName string
}


// Open initialzes necessary variables for using this executor.
func (e *InspectionReaderExecutor) Open(ctx context.Context) error {
	 // tp, ok := e.Table.StreamProperties["type"]
	 // if !ok {
	 // 	return errors.New("Cannot find stream table type")
	 // }

	 // e.tp = tp
	 // e.setVariableName(strings.ToLower(tp))

	 // e.topic, ok = e.Table.StreamProperties["topic"]
	 // if !ok {
	 // 	return errors.New("Cannot find stream table topic")
	 // }

	//  var err error
	//  value, err := e.ctx.GetSessionVars().GlobalVarsAccessor.GetGlobalSysVar(e.variableName)
	//  if err != nil {
	//  	return errors.Trace(err)
	//  }
	//  if value != "" {
	//  	e.pos, err = strconv.Atoi(value)
	//  	if err != nil {
	//  		return errors.Trace(err)
	//  	}
	//  } else {
	//  	err = e.ctx.GetSessionVars().GlobalVarsAccessor.SetGlobalSysVar(e.variableName, "0")
	//  	if err != nil {
	//  		return errors.Trace(err)
	// 	}
	// }

	return nil
}

// Next fills data into the chunk passed by its caller.
func (e *InspectionReaderExecutor) Next(ctx context.Context, chk *chunk.Chunk) error {
	 // value, err := e.ctx.GetSessionVars().GlobalVarsAccessor.GetGlobalSysVar(e.variableName)
	 // if err != nil {
	 // 	return errors.Trace(err)
	 // }
	 // e.pos, err = strconv.Atoi(value)
	 // if err != nil {
	 // 	return errors.Trace(err)
	 // }
	 var err error
	 e.pos = 0

	 pos := 0
	 chk.GrowAndReset(e.maxChunkSize)
	 if e.result == nil {
	 	e.result = newFirstChunk(e)
	 	pos, err = e.fetchAll(e.pos)
	 	if err != nil {
	 		return errors.Trace(err)
	 	}
	 	iter := chunk.NewIterator4Chunk(e.result)
	 	for colIdx := 0; colIdx < e.Schema().Len(); colIdx++ {
	 		retType := e.Schema().Columns[colIdx].RetType
	 		if !types.IsTypeVarchar(retType.Tp) {
	 			continue
	 		}
	 		for row := iter.Begin(); row != iter.End(); row = iter.Next() {
	 			if valLen := len(row.GetString(colIdx)); retType.Flen < valLen {
	 				retType.Flen = valLen
	 			}
	 		}
	 	}
	 }
	 if e.cursor >= e.result.NumRows() {
	 	return nil
	 }
	 numCurBatch := mathutil.Min(chk.Capacity(), e.result.NumRows()-e.cursor)
	 chk.Append(e.result, e.cursor, e.cursor+numCurBatch)
	 e.cursor += numCurBatch

	 e.pos = pos
	 err = e.ctx.GetSessionVars().GlobalVarsAccessor.SetGlobalSysVar(e.variableName, strconv.Itoa(e.pos))
	 if err != nil {
	 	return errors.Trace(err)
	}

	return nil
}

// Close implements the Executor Close interface.
func (e *InspectionReaderExecutor) Close() error {
	return nil
}

func (e *InspectionReaderExecutor) fetchAll(cursor int) (int, error) {
	var pos int
	var err error
	pos, err = e.fetchMockData(cursor)
	if err != nil {
			return 0, errors.Trace(err)
	}

	return pos, nil
}

func (e *InspectionReaderExecutor) fetchMockData(cursor int) (int, error) {
	 var pos int
	 for i := cursor; i < maxFetchCnt && i < cursor+batchFetchCnt; {
	 	data, err := e.getData(mock.MockStreamJsonData[i])
	 	if err != nil {
	 		return 0, errors.Trace(err)
	 	}

	 	row := chunk.MutRowFromDatums(data).ToRow()
	 	e.result.AppendRow(row)

	 	i++
		pos = i
	 }

	return pos, nil
	//return 0, nil
}

func (e *InspectionReaderExecutor) getData(data string) ([]types.Datum, error) {
	m := make(map[string]interface{})
	err := gojson.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, errors.Trace(err)
	}

	row := make([]types.Datum, 0, len(e.Columns))
	for _, col := range e.Columns {
		name := col.Name.L
		if value, ok := m[name]; ok {
			data := types.NewDatum(value)
			val, err := data.ConvertTo(e.ctx.GetSessionVars().StmtCtx, &col.FieldType)
			if err != nil {
				return nil, errors.Trace(err)
			}
			row = append(row, val)
		} else {
			data := types.NewDatum(nil)
			row = append(row, data)
		}
	}

	return row, nil
}

func (e *InspectionReaderExecutor) appendRow(chk *chunk.Chunk, row []interface{}) {
	for i, col := range row {
		if col == nil {
			chk.AppendNull(i)
			continue
		}
		switch x := col.(type) {
		case nil:
			chk.AppendNull(i)
		case int:
			chk.AppendInt64(i, int64(x))
		case int64:
			chk.AppendInt64(i, x)
		case uint64:
			chk.AppendUint64(i, x)
		case float64:
			chk.AppendFloat64(i, x)
		case float32:
			chk.AppendFloat32(i, x)
		case string:
			chk.AppendString(i, x)
		case []byte:
			chk.AppendBytes(i, x)
		case types.BinaryLiteral:
			chk.AppendBytes(i, x)
		case *types.MyDecimal:
			chk.AppendMyDecimal(i, x)
		case types.Time:
			chk.AppendTime(i, x)
		case json.BinaryJSON:
			chk.AppendJSON(i, x)
		case types.Duration:
			chk.AppendDuration(i, x)
		case types.Enum:
			chk.AppendEnum(i, x)
		case types.Set:
			chk.AppendSet(i, x)
		default:
			chk.AppendNull(i)
		}
	}
}
