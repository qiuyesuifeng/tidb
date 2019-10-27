// Copyright 2019 PingCAP, Inc.
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

package log

import (
	"path"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/infoschema/inspection/log/item"
	"github.com/pingcap/tidb/infoschema/inspection/log/parser"
)

var (
	tidbLogPath = "/Users/duchuan.dc/pingcap/hackathon2019/tidb_log/"
	fw          []*parser.FileWrapper
)

type LogItem struct {
	Address  string `json:"address"`
	FileName string `json:"filename"`
	Time     int64  `json:"time"`
	Level    string `json:"level"`
	Content  string `json:"content"`
}

type TiDBLogBatch struct {
	Logs []*LogItem `json:"logs"`
	Cnt  int        `json:"cnt"`
}

func ResetTiDBLogPath(p string) {
	tidbLogPath = p
}

func GetTiDBLogPath() string {
	dir, _ := path.Split(config.GetGlobalConfig().Log.File.Filename)
	return dir
}

func ParseLevelToStr(level item.LevelType) string {
	if level == item.LevelInvalid {
		return "invalid"
	}
	if level == item.LevelFATAL {
		return "fatal"
	}
	if level == item.LevelERROR {
		return "error"
	}
	if level == item.LevelWARN {
		return "warn"
	}
	if level == item.LevelINFO {
		return "info"
	}
	if level == item.LevelDEBUG {
		return "debug"
	}
	return "unknown"
}
