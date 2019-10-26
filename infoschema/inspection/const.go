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

package inspection

type inspectionVirtualTableItem struct {
	SQL string
	Attrs map[string]string
}

var inspectionVirtualTables = []inspectionVirtualTableItem{
	{tableLocalLog,
		map[string]string{
			"type": "log_tidb_local",
			"startTime": "2019-10-24T11:35:29",
			"endTime": "2019-10-24T11:35:47",
			"limit": "77",
		},
	},
	{tableRemoteLog,
		map[string]string{
			"type": "log_tidb_remote",
			"url": "http://127.0.0.1:10080/log",
			"startTime": "2019-10-24T11:35:29",
			"endTime": "2019-10-24T11:35:47",
			"limit": "5",
		},
	},
}

const tableLocalLog = "CREATE TABLE %s.local_log(" +
	"host varchar(256)," +
	"port varchar(256)," +
	"component varchar(256)," +
	"filename varchar(256)," +
	"time timestamp," +
	"level bigint," +
	"content text);"

const tableRemoteLog = "CREATE TABLE %s.remote_log(" +
	"host varchar(256)," +
	"port varchar(256)," +
	"component varchar(256)," +
	"filename varchar(256)," +
	"time timestamp," +
	"level bigint," +
	"content text);"

var inspectionPersistTables = []string{
	// tablePersistTest,
	tableTiDBClusterInfo,
	tableSystemInfo,
	tableTiDBClusterKeyMetrcisInfo,
	tableTiDBKeyMetrcisInfo,
	tableTiKVKeyMetrcisInfo,
	tableTiKVPerformanceInfo,
}

const tablePersistTest = "CREATE TABLE %s.test_persist(" +
	"click_id bigint(20)," +
	"user_id bigint(20));"

const tableTiDBClusterInfo = `CREATE TABLE %s.TIDB_CLUSTER_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  TYPE varchar(64) DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  ADDRESS varchar(64) DEFAULT NULL,
  STATUS_ADDRESS varchar(64) DEFAULT NULL,
  VERSION varchar(64) DEFAULT NULL,
  GIT_HASH varchar(64) DEFAULT NULL,
  CONFIG text DEFAULT NULL
)`

const tableSystemInfo = `CREATE TABLE %s.SYSTEM_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  TYPE varchar(64) DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  IP varchar(64) DEFAULT NULL,
  STATUS_ADDRESS varchar(64) DEFAULT NULL,
  CPU varchar(64) DEFAULT NULL,
  CPU_USAGE varchar(64) DEFAULT NULL,
  MEMORY varchar(64) DEFAULT NULL,
  MEMORY_USAGE varchar(64) DEFAULT NULL,
  LOAD1 varchar(64) DEFAULT NULL,
  LOAD5 varchar(64) DEFAULT NULL,
  LOAD15 varchar(64) DEFAULT NULL,
  KERNAL varchar(128) DEFAULT NULL
)`

const tableTiDBClusterKeyMetrcisInfo = `CREATE TABLE %s.TIDB_CLUSTER_KEY_METRICS_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  CONNECTION_COUNT varchar(64) DEFAULT NULL,
  QUERY_OK_COUNT varchar(64) DEFAULT NULL,
  QUERY_ERR_COUNT varchar(64) DEFAULT NULL,
  INSERT_COUNT varchar(64) DEFAULT NULL,
  UPDATE_COUNT varchar(64) DEFAULT NULL,
  DELETE_COUNT varchar(64) DEFAULT NULL,
  REPLACE_COUNT varchar(64) DEFAULT NULL,
  SELECT_COUNT varchar(64) DEFAULT NULL,
  80_QUERY_DURATION varchar(64) DEFAULT NULL,
  95_QUERY_DURATION varchar(64) DEFAULT NULL,
  99_QUERY_DURATION varchar(64) DEFAULT NULL,
  999_QUERY_DURATION varchar(64) DEFAULT NULL,
  AVAILABLE varchar(64) DEFAULT NULL,
  CAPACITY varchar(64) DEFAULT NULL
)`

const tableTiDBKeyMetrcisInfo = `CREATE TABLE %s.TIDB_KEY_METRICS_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  TYPE varchar(64) DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  IP varchar(64) DEFAULT NULL,
  STATUS_ADDRESS varchar(64) DEFAULT NULL,
  CONNECTION_COUNT varchar(64) DEFAULT NULL,
  QUERY_OK_COUNT varchar(64) DEFAULT NULL,
  QUERY_ERR_COUNT varchar(64) DEFAULT NULL,
  80_QUERY_DURATION varchar(64) DEFAULT NULL,
  95_QUERY_DURATION varchar(64) DEFAULT NULL,
  99_QUERY_DURATION varchar(64) DEFAULT NULL,
  999_QUERY_DURATION varchar(64) DEFAULT NULL,
  UPTIME varchar(64) DEFAULT NULL
)`

const tableTiKVKeyMetrcisInfo = `CREATE TABLE %s.TIKV_KEY_METRICS_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  TYPE varchar(64) DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  IP varchar(64) DEFAULT NULL,
  STATUS_ADDRESS varchar(64) DEFAULT NULL,
  AVAILABLE varchar(64) DEFAULT NULL,
  CAPACITY varchar(64) DEFAULT NULL,
  CPU varchar(64) DEFAULT NULL,
  MEMORY varchar(64) DEFAULT NULL,
  LEADER_COUNT varchar(64) DEFAULT NULL,
  REGION_COUNT varchar(64) DEFAULT NULL,
  KV_GET_COUNT varchar(64) DEFAULT NULL,
  KV_BATCH_GET_COUNT varchar(64) DEFAULT NULL,
  KV_SCAN_COUNT varchar(64) DEFAULT NULL,
  KV_PREWRITE_COUNT varchar(64) DEFAULT NULL,
  KV_COMMIT_COUNT varchar(64) DEFAULT NULL,
  KV_COPROCESSOR_COUNT varchar(64) DEFAULT NULL
)`

const tableTiKVPerformanceInfo = `CREATE TABLE %s.TIKV_PERFORMANCE_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  TYPE varchar(64) DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  IP varchar(64) DEFAULT NULL,
  STATUS_ADDRESS varchar(64) DEFAULT NULL,
  99_KV_GET_DURATION varchar(64) DEFAULT NULL,
  99_KV_BATCH_GET_DURATION varchar(64) DEFAULT NULL,
  99_KV_SCAN_DURATION varchar(64) DEFAULT NULL,
  99_KV_PREWRITE_DURATION varchar(64) DEFAULT NULL,
  99_KV_COMMIT_DURATION varchar(64) DEFAULT NULL,
  99_KV_COPROCESSOR_DURATION varchar(64) DEFAULT NULL,
  RAFT_STORE_CPU_USAGE varchar(64) DEFAULT NULL,
  ASYNC_APPLY_CPU_USAGE varchar(64) DEFAULT NULL,
  SCHEDULER_WORKER_CPU_USAGE varchar(64) DEFAULT NULL,
  COPROCESSOR_CPU_USAGE varchar(64) DEFAULT NULL,
  ROCKSDB_CPU_USAGE varchar(64) DEFAULT NULL
)`
