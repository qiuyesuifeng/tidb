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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	log2 "github.com/pingcap/tidb/infoschema/inspection/log"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/chunk"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// make sure `LocalLogReaderExecutor` implements `Executor`.
var _ Executor = &RemoteLogReaderExecutor{}

type RemoteLogReaderExecutor struct {
	baseExecutor
	Table   *model.TableInfo
	Columns []*model.ColumnInfo

	startTimeStr string
	endTimeStr   string
	LimitStr     string
	result       *chunk.Chunk
	cnt          int
	limit        int

	pattern  string
	level    string
	filename string
	nodes    string
	fds      []*resultBuffer
	address  string
}

type resultBuffer struct {
	tp string
	//buffer      *log2.TiDBLogBatch
	buffer      []*log2.LogItem
	idxInBuffer int
	valid       bool
	statusAddr  string
	fd          string
}

func (e *RemoteLogReaderExecutor) Open(ctx context.Context) error {
	local, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		return err
	}

	st, err := time.ParseInLocation(types.TimeStampLayout, e.startTimeStr, local)
	if err != nil {
		return err
	}
	et, err := time.ParseInLocation(types.TimeStampLayout, e.endTimeStr, local)
	if err != nil {
		return err
	}
	l, err := strconv.ParseInt(e.LimitStr, 10, 64)
	if err != nil {
		return err
	}
	e.limit = int(l)
	e.cnt = 0
	nodes := strings.Split(e.nodes, ";")
	for _, node := range nodes {
		if node == "" {
			continue
		}
		segments := strings.Split(node, "@")
		tp := segments[0]
		statusAddr := segments[1]
		if e.address != "" && statusAddr != e.address {
			continue
		}

		b := resultBuffer{
			tp:          tp,
			idxInBuffer: 0,
			valid:       true,
			statusAddr:  statusAddr,
		}
		var url string
		if tp == "tidb" {
			url = fmt.Sprintf("http://%s/log/open?start_time=%s&end_time=%s", statusAddr, e.startTimeStr, e.endTimeStr)
		} else if tp == "tikv" {
			url = fmt.Sprintf("http://%s/log/open?start_time=%d&end_time=%d", statusAddr, st.Unix()*1000, et.Unix()*1000)
		}
		if e.pattern != "" {
			url = fmt.Sprintf("%s&pattern=%s", url, e.pattern)
		}
		if e.filename != "" {
			url = fmt.Sprintf("%s&filename=%s", url, e.filename)
		}
		if e.level != "" {
			url = fmt.Sprintf("%s&level=%s", url, e.level)
		}
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		v := make(map[string]string)
		if err := json.Unmarshal(body, &v); err != nil {
			return err
		}
		fd, exists := v["fd"]
		if !exists {
			return errors.New("illegal http return. missing field: `fd`")
		}
		b.fd = fd
		e.fds = append(e.fds, &b)
	}
	return nil
}

// Next fills data into the chunk passed by its caller.
func (e *RemoteLogReaderExecutor) Next(ctx context.Context, chk *chunk.Chunk) error {
	var err error

	chk.GrowAndReset(chk.Capacity())
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
func (e *RemoteLogReaderExecutor) Close() error {
	for _, b := range e.fds {
		url := fmt.Sprintf("http://%s/log/close?fd=%s", b.statusAddr, b.fd)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *RemoteLogReaderExecutor) fetchAll() error {
	err := e.fetchRemoteLog()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *RemoteLogReaderExecutor) fetchRemoteLog() error {
	l := e.result.Capacity()
	if e.limit-e.cnt < l {
		l = e.limit - e.cnt
	}
	if l > 10 {
		l = 10
	}
	for i := 0; i < l; i++ {
		var item *log2.LogItem = nil
		var idx = 0
		for j, n := range e.fds {
			if !n.valid {
				continue
			}
			if n.buffer == nil || n.idxInBuffer == len(n.buffer) {
				var err error
				if n.tp == "tidb" {
					err = e.fetchOneTiDB(n)
				} else if n.tp == "tikv" {
					err = e.fetchOneTiKV(n)
				}
				if err != nil {
					return err
				}
			}
			if !n.valid {
				continue
			}
			if item == nil {
				item = n.buffer[n.idxInBuffer]
				idx = j
			} else {
				if n.buffer[n.idxInBuffer].Time < item.Time {
					item = n.buffer[n.idxInBuffer]
					idx = j
				}
			}
		}
		if item != nil {
			e.fds[idx].idxInBuffer = e.fds[idx].idxInBuffer + 1
			data, err := e.getData(item, e.fds[idx].tp, e.fds[idx].statusAddr)
			if err != nil {
				return errors.Trace(err)
			}
			row := chunk.MutRowFromDatums(data).ToRow()
			e.result.AppendRow(row)
			e.cnt++
		} else {
			break
		}
	}
	return nil
}

func (e *RemoteLogReaderExecutor) getData(data *log2.LogItem, tp, host string) ([]types.Datum, error) {
	row := make([]types.Datum, 0, len(e.Columns))
	for _, col := range e.Columns {
		switch col.Name.L {
		case "address":
			row = append(row, types.NewStringDatum(host))
		case "type":
			row = append(row, types.NewStringDatum(tp))
		case "filename":
			row = append(row, types.NewStringDatum(data.FileName))
		case "time":
			tm := types.Time{
				Time: types.FromGoTime(time.Unix(data.Time/1000, 0)),
				Type: mysql.TypeDatetime,
				Fsp:  0,
			}
			row = append(row, types.NewTimeDatum(tm))
		case "level":
			row = append(row, types.NewStringDatum(data.Level))
		case "content":
			row = append(row, types.NewStringDatum(data.Content))
		default:
			data := types.NewDatum(nil)
			row = append(row, data)
		}
	}
	return row, nil
}

func (e *RemoteLogReaderExecutor) fetchOneTiDB(r *resultBuffer) error {
	l := e.result.Capacity()
	if e.limit-e.cnt < l {
		l = e.limit - e.cnt
	}
	url := fmt.Sprintf("http://%s/log/next?fd=%s&limit=%d", r.statusAddr, r.fd, l)
	if e.pattern != "" {
		url = fmt.Sprintf("%s&pattern=%s", url, e.pattern)
	}
	resp, err := http.Get(url)
	if err != nil {
		r.valid = false
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	logrus.Infof("data=%s|", string(body))
	if err != nil {
		//logrus.Info("ERR=%s", err)
		return err
	}
	logBatch := &log2.TiDBLogBatch{}
	if err := json.Unmarshal(body, logBatch); err != nil {
		return err
	}
	if logBatch.Cnt == 0 {
		r.valid = false
	}
	r.buffer = logBatch.Logs
	r.idxInBuffer = 0
	return nil
}

func (e *RemoteLogReaderExecutor) fetchOneTiKV(r *resultBuffer) error {
	l := e.result.Capacity()
	if e.limit-e.cnt < l {
		l = e.limit - e.cnt
	}
	url := fmt.Sprintf("http://%s/log/next?fd=%s&limit=%d", r.statusAddr, r.fd, l)
	if e.pattern != "" {
		url = fmt.Sprintf("%s&pattern=%s", url, e.pattern)
	}
	resp, err := http.Get(url)
	if err != nil {
		r.valid = false
		return err
	}
	defer resp.Body.Close()
	//logrus.Infof("URL=%s|", url)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//logrus.Info("ERR=%s", err)
		return err
	}
	//logrus.Info("RET=%s", string(body))
	logBatch := make([]*log2.LogItem, 0, 0)
	if err := json.Unmarshal(body, &logBatch); err != nil {
		//logrus.Info("KV_UMA_ERR=%s", err.Error())
		return err
	}
	if len(logBatch) == 0 {
		r.valid = false
	}
	r.buffer = logBatch
	//logrus.Infof("TiKV_BATCH=%+v|", logBatch)
	r.idxInBuffer = 0
	return nil
}

/*
func (e *RemoteLogReaderExecutor) xxx() error {
	l := e.result.Capacity()
	if e.limit - e.cnt < l {
		 l = e.limit - e.cnt
	}
	url := fmt.Sprintf("http://%s/log/next?fd=%s&limit=%d", e.address, e.fd, l)
	logrus.Infof("URL=%s|",url)
	if e.pattern != "" {
		url = fmt.Sprintf("%s&pattern=%s", url, e.pattern)
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logrus.Infof("FD=%s", e.fd)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("ERR=%s", err)
		return err
	}
	logrus.Info("RET=%s", string(body))
	logBatch := &log2.TiDBLogBatch{}
	if err := json.Unmarshal(body, logBatch); err != nil {
		return err
	}
	for _, item := range logBatch.Logs {
		if e.cnt >= e.limit {
			break
		}
		if item == nil {
			continue
		}
		data, err := e.getData(item)
		if err != nil {
			return errors.Trace(err)
		}
		row := chunk.MutRowFromDatums(data).ToRow()
		e.result.AppendRow(row)
		e.cnt++
	}
	return nil
}

*/

/*
func (e *RemoteLogReaderExecutor) xxx() error {
	l := e.result.Capacity()
	if e.limit - e.cnt < l {
		 l = e.limit - e.cnt
	}
	url := fmt.Sprintf("http://%s/log/next?fd=%s&limit=%d", e.address, e.fd, l)
	logrus.Infof("URL=%s|",url)
	if e.pattern != "" {
		url = fmt.Sprintf("%s&pattern=%s", url, e.pattern)
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logrus.Infof("FD=%s", e.fd)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("ERR=%s", err)
		return err
	}
	logrus.Info("RET=%s", string(body))
	logBatch := &log2.TiDBLogBatch{}
	if err := json.Unmarshal(body, logBatch); err != nil {
		return err
	}
	for _, item := range logBatch.Logs {
		if e.cnt >= e.limit {
			break
		}
		if item == nil {
			continue
		}
		data, err := e.getData(item)
		if err != nil {
			return errors.Trace(err)
		}
		row := chunk.MutRowFromDatums(data).ToRow()
		e.result.AppendRow(row)
		e.cnt++
	}
	return nil
}

*/
