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

package parser

import (
	"time"

	"github.com/pingcap/tidb/infoschema/inspection/log/item"
)

// The log parser representation
type Parser interface {
	// Parse the first line of the log, and return ts and level,
	// if it's not the first line, nil will be returned.
	ParseHead(head []byte) (*time.Time, item.LevelType)
}

// All parsers this package has
func List() []Parser {
	return []Parser{
		//&PDLogV2Parser{},
		//&TiKVLogV2Parser{},
		&UnifiedLogParser{},
		//&SlowQueryParser{},
	}
}
