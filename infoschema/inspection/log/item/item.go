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

package item

import (
	"time"
)

type ItemType int16

type LevelType int16

const (
	LevelInvalid LevelType = iota
	LevelFATAL
	LevelERROR
	LevelWARN
	LevelINFO
	LevelDEBUG
)

const (
	TypeInvalid ItemType = iota
	TypeTiDB
	TypeTiKV
	TypePD
	TypeTiDBSlowQuery
)

// Item represent a log entity
type Item interface {
	// The host which produce the log
	GetHost() string

	// The service produced the log listening on
	GetPort() string

	// The component produced the log. eg. pd, tidb, tikv.
	GetComponent() string

	// The log file name
	GetFileName() string

	// The time when the log produced
	GetTime() time.Time

	// The log level
	GetLevel() LevelType

	// The whole log content
	GetContent() []byte

	// Append log content
	AppendContent([]byte) error
}
