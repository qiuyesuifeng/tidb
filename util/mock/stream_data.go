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

package mock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/sessionctx/stmtctx"
	"github.com/pingcap/tidb/types"
	log "github.com/sirupsen/logrus"
)

// Mock data for stream reader demo.
var (
	MockStreamData       []Event
	MockStreamClickData  []Click
	MockKafkaStreamData  []KafkaEvent
	MockPulsarStreamData []PulsarEvent
	MockLogStreamData    []LogEvent
	MockStreamJsonData   []string
)

type Click struct {
	ClickID    int64  `json:"click_id"`
	UserID     int64  `json:"user_id"`
	AdsID      int64  `json:"ads_id"`
	ClickPrice int64  `json:"click_price"`
	CreateTime string `json:"create_time"`
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randInt64(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func getDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

type Event struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

type KafkaEvent struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

type PulsarEvent struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

type LogEvent struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

func genData(i int, t types.Time) Event {
	tt := types.Time{
		Time: types.FromDate(t.Time.Year(), t.Time.Month(), t.Time.Day(), t.Time.Hour(), t.Time.Minute(), t.Time.Second(), t.Time.Microsecond()),
		Type: mysql.TypeTimestamp,
		Fsp:  types.DefaultFsp,
	}
	content := fmt.Sprintf("TiDB Stream Demo Data %d", i)
	evt := Event{int64(i), content, tt}

	return evt
}

func genClickData(i int) Click {
	clickID := int64(i + 1)
	userID := randInt64(101, 120)
	adsID := randInt64(1001, 1010)
	clickPrice := randInt64(1, 10) * 100
	click := Click{clickID, userID, adsID, clickPrice, getDateTime(time.Now().Add(time.Second * time.Duration(i)))}

	return click
}

func genKafkaData(i int, t types.Time) KafkaEvent {
	tt := types.Time{
		Time: types.FromDate(t.Time.Year(), t.Time.Month(), t.Time.Day(), t.Time.Hour(), t.Time.Minute(), t.Time.Second(), t.Time.Microsecond()),
		Type: mysql.TypeTimestamp,
		Fsp:  types.DefaultFsp,
	}
	content := fmt.Sprintf("TiDB Stream Demo Data %d", i)
	evt := KafkaEvent{int64(i), content, tt}

	return evt
}

func genPulsarData(i int, t types.Time) PulsarEvent {
	evt := PulsarEvent{}
	return evt
}

func genLogData(i int, t types.Time) LogEvent {
	evt := LogEvent{}
	return evt
}

func init() {
	rand.Seed(time.Now().UnixNano())

	sc := &stmtctx.StatementContext{
		TimeZone: time.UTC,
	}
	secondDur, err := types.ParseDuration(sc, "00:00:01", types.MaxFsp)
	if err != nil {
		log.Fatalf("[parse duration failed]%v", err)
	}

	t := types.CurrentTime(mysql.TypeTimestamp)
	for i := 0; i < 1000; i++ {
		t, err = t.Add(sc, secondDur)
		if err != nil {
			log.Fatalf("[add time duration failed]%v", err)
		}

		MockStreamData = append(MockStreamData, genData(i, t))
		MockStreamClickData = append(MockStreamClickData, genClickData(i))
		MockKafkaStreamData = append(MockKafkaStreamData, genKafkaData(i, t))
		MockPulsarStreamData = append(MockPulsarStreamData, genPulsarData(i, t))
		MockLogStreamData = append(MockLogStreamData, genLogData(i, t))
	}

	for i := 0; i < 1000; i++ {
		data, err := json.Marshal(MockStreamClickData[i])
		if err != nil {
			log.Fatalf("[mock stream data marshal failed]%v", err)
		}

		MockStreamJsonData = append(MockStreamJsonData, string(data))
	}
}
