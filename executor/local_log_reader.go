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
	"fmt"
	"strconv"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/config"
	log2 "github.com/pingcap/tidb/infoschema/inspection/log"
	"github.com/pingcap/tidb/infoschema/inspection/log/item"
	"github.com/pingcap/tidb/infoschema/inspection/log/search"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/chunk"
	"golang.org/x/net/context"
)

// make sure `LocalLogReaderExecutor` implements `Executor`.
var _ Executor = &LocalLogReaderExecutor{}

type LocalLogReaderExecutor struct {
	baseExecutor
	Table   *model.TableInfo
	Columns []*model.ColumnInfo

	se *search.Sequence
	result *chunk.Chunk
	cnt int
	limit int
	startTimeStr string
	endTimeStr string
	LimitStr string
}

func (e *LocalLogReaderExecutor) Open(ctx context.Context) error {
	TimeStampLayout := "2006-01-02T15:04:05"
	local, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		return err
	}
	// startTime, _ := time.ParseInLocation(TimeStampLayout ,"1970-01-01T00:00:00", local)
	// endTime, _ := time.ParseInLocation(TimeStampLayout , "2030-01-01T00:00:00", local)
	startTime, err := time.ParseInLocation(TimeStampLayout ,e.startTimeStr, local)
	if err != nil {
		return err
	}
	endTime, err := time.ParseInLocation(TimeStampLayout , e.endTimeStr, local)
	if err != nil {
		return err
	}
	l, err := strconv.ParseInt(e.LimitStr, 10, 64)
	if err != nil {
		return err
	}
	e.limit = int(l)
	if e.se, err = search.NewSequence(log2.GetTiDBLogPath(), startTime, endTime); err != nil {
		return err
	}
	e.cnt = 0
	return nil
}

// Next fills data into the chunk passed by its caller.
func (e *LocalLogReaderExecutor) Next(ctx context.Context, chk *chunk.Chunk) error {
	var err error

	chk.GrowAndReset(e.maxChunkSize)
	if e.result == nil {
		e.result = newFirstChunk(e)
	}
	e.result.Reset()
	err = e.fetchAll()
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
	chk.Append(e.result, 0, e.result.NumRows())
	return nil
}

// Close implements the Executor Close interface.
func (e *LocalLogReaderExecutor) Close() error {
	return nil
}

func (e *LocalLogReaderExecutor) fetchAll() error {
	err := e.fetchData()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *LocalLogReaderExecutor) fetchData() error {
	for i := 0; i < 10; i++ {
		if e.cnt >= e.limit {
			break
		}
		item, err := e.se.Next()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return errors.Trace(err)
			}
		}
		data, err := e.parseData(item)
		if err != nil {
			return errors.Trace(err)
		}
		row := chunk.MutRowFromDatums(data).ToRow()
		e.result.AppendRow(row)
		e.cnt++
	}
	return nil
}

func (e *LocalLogReaderExecutor) parseData(data item.Item) ([]types.Datum, error) {
	row := make([]types.Datum, 0, len(e.Columns))
	for _, col := range e.Columns {
		switch col.Name.L {
		case "address":
			row = append(row, types.NewStringDatum(fmt.Sprintf("%s:%d", config.GetGlobalConfig().Status.StatusHost, config.GetGlobalConfig().Status.StatusPort)))
		case "component":
			row = append(row, types.NewStringDatum("tidb"))
		case "filename":
			row = append(row, types.NewStringDatum(data.GetFileName()))
		case "time":
			tm := types.Time{
				Time: types.FromGoTime(data.GetTime()),
				Type: mysql.TypeDatetime,
				Fsp:  0,
			}
			row = append(row, types.NewTimeDatum(tm))
		case "level":
			row = append(row, types.NewStringDatum(log2.ParseLevelToStr(data.GetLevel())))
		case "content":
			row = append(row, types.NewStringDatum(string(data.GetContent())))
		default:
			data := types.NewDatum(nil)
			row = append(row, data)
		}
	}
	return row, nil
}
