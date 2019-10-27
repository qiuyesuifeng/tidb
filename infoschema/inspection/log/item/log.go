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
	"fmt"
	"time"
)

const MAX_LOG_SIZE = 1024 * 1024 * 50

// The LogItem struct implements Item interface
type LogItem struct {
	File      string    // name of log file
	Time      time.Time // timestamp of a single line log
	Level     LevelType // log level
	Host      string    // host ip
	Port      string    // port of component
	Component string    // name of component
	Content   []byte    // content of a entire line log
	Type      ItemType  // type of log file
}

func (l *LogItem) GetHost() string {
	return l.Host
}

func (l *LogItem) GetPort() string {
	return l.Port
}

func (l *LogItem) GetComponent() string {
	return l.Component
}

func (l *LogItem) GetFileName() string {
	return l.File
}

func (l *LogItem) GetTime() time.Time {
	return l.Time
}

func (l *LogItem) GetLevel() LevelType {
	return l.Level
}

func (l *LogItem) GetContent() []byte {
	return l.Content
}

func (l *LogItem) AppendContent(content []byte) error {
	if len(l.Content) > MAX_LOG_SIZE {
		return fmt.Errorf("log size exceeds limit, log content:\n%s\n", string(l.Content))
	}
	l.Content = append(l.Content, byte('\n'))
	l.Content = append(l.Content, content...)
	return nil
}
