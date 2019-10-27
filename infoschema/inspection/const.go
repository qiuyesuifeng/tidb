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

const tableClusterLog = "CREATE TABLE %s.CLUSTER_LOG(" +
	"ADDRESS varchar(64)," +
	"TYPE varchar(64)," +
	"FILENAME varchar(256)," +
	"TIME timestamp," +
	"LEVEL varchar(10)," +
	"CONTENT text);"

var inspectionPersistTables = []string{
	tableTiDBClusterInfo,
	tableSystemInfo,
	tableTiDBClusterKeyMetrcisInfo,
	tableTiDBKeyMetrcisInfo,
	tableTiKVKeyMetrcisInfo,
	tableTiKVPerformanceInfo,
	tableInspectionResult,
	tableTiDBCpuProfile,
	tableTiKVCpuProfile,
	tableSlowQueryDetail,
}

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
  SLOW_QUERY_COUNT varchar(64) DEFAULT NULL,
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

const tableInspectionResult = `CREATE TABLE %s.RESULT (
  ID bigint(21) unsigned DEFAULT NULL,
  METRICS TEXT DEFAULT NULL,
  RESULT TEXT DEFAULT NULL
)`

// tableTiDBCpuProfile contains the columns name definitions for table tidb_cpu_profile
const tableTiDBCpuProfile = "CREATE TABLE IF NOT EXISTS %s.TIDB_CPU_PROFILE (" +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"PERCENT_ABS VARCHAR(8) NOT NULL," +
	"PERCENT_REL VARCHAR(8) NOT NULL," +
	"ROOT_CHILD INT(8) NOT NULL," +
	"DEPTH INT(8) NOT NULL," +
	"FILE VARCHAR(512) NOT NULL);"

	// tableTiKVCpuProfile contains the columns name definitions for table tikv_cpu_profile
const tableTiKVCpuProfile = "CREATE TABLE IF NOT EXISTS %s.TIKV_CPU_PROFILE (" +
	"NAME VARCHAR(16) NOT NULL," +
	"ADDRESS VARCHAR(64) NOT NULL," +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"PERCENT_ABS VARCHAR(8) NOT NULL," +
	"PERCENT_REL VARCHAR(8) NOT NULL," +
	"ROOT_CHILD INT(8) NOT NULL," +
	"DEPTH INT(8) NOT NULL," +
	"FILE VARCHAR(512) NOT NULL);"

const tableSlowQueryDetail = "CREATE TABLE IF NOT EXISTS %s.SLOW_QUERY_DETAIL (" +
	"ID BIGINT(20) NOT NULL," +
	"TYPE VARCHAR(8) NOT NULL," +
	"NAME VARCHAR(128) NOT NULL," +
	"DATA TEXT NOT NULL);"
