// Copyright 2016 PingCAP, Inc.
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

var inspectionVirtualTables = []string{
	tableTest,
}

const tableTest = "CREATE TABLE %s.test_virtual(" +
	"click_id bigint(20)," +
	"user_id bigint(20));"

var inspectionPersistTables = []string{
	tablePersistTest,
	tableTiDBClusterInfo,
	tableSystemInfo,
	tableTiDBKeyMetrcisInfo,
	tableTiKVKeyMetrcisInfo,
	tableTiDBPerformanceInfo,
	tableTiKVPerformanceInfo,
	tableTiDBErrorInfo,
	tableTiKVErrorInfo,
}

const tablePersistTest = "CREATE TABLE %s.test_persist(" +
	"click_id bigint(20)," +
	"user_id bigint(20));"

const tableTiDBClusterInfo = `CREATE TABLE %s.TIDB_CLUSTER_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  TYPE varchar(64) DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  ADDRESS varchar(64) DEFAULT NULL,
  VERSION varchar(64) DEFAULT NULL,
  GIT_HASH varchar(64) DEFAULT NULL,
  CONFIG text DEFAULT NULL
)`

const tableSystemInfo = `CREATE TABLE %s.SYSTEM_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  IP varchar(64) DEFAULT NULL,
  CPU varchar(64) DEFAULT NULL,
  CPU_USAGE double DEFAULT NULL,
  MEMORY varchar(64) DEFAULT NULL,
  MEMORY_USAGE double DEFAULT NULL,
  VERSION varchar(64) DEFAULT NULL,
  OS_VERSION varchar(128) DEFAULT NULL,
  KERNAL_VERSION varchar(128) DEFAULT NULL
)`

const tableTiDBKeyMetrcisInfo = `CREATE TABLE %s.TIDB_KEY_METRICS_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  STATE varchar(64) DEFAULT NULL,
  CONNECTION_COUNT bigint(21) unsigned DEFAULT NULL,
  UPTIME varchar(64) DEFAULT NULL
)`

const tableTiKVKeyMetrcisInfo = `CREATE TABLE %s.TIKV_KEY_METRICS_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  STATE varchar(64) DEFAULT NULL,
  VERSION varchar(64) DEFAULT NULL,
  CAPACITY varchar(64) DEFAULT NULL,
  AVAILABLE varchar(64) DEFAULT NULL,
  LEADER_COUNT bigint(21) unsigned DEFAULT NULL,
  LEADER_WEIGHT double DEFAULT NULL,
  LEADER_SCORE double DEFAULT NULL,
  LEADER_SIZE bigint(21) unsigned DEFAULT NULL,
  REGION_COUNT bigint(21) unsigned DEFAULT NULL,
  REGION_WEIGHT double DEFAULT NULL,
  REGION_SCORE double DEFAULT NULL,
  REGION_SIZE bigint(21) unsigned DEFAULT NULL,
  START_TS datetime DEFAULT NULL,
  LAST_HEARTBEAT_TS datetime DEFAULT NULL,
  UPTIME varchar(64) DEFAULT NULL
)`

const tableTiDBPerformanceInfo = `CREATE TABLE %s.TIDB_PERFORMANCE_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  QPS bigint(21) unsigned DEFAULT NULL,
  QUERY_DURATION datetime DEFAULT NULL,
  TPS bigint(21) unsigned DEFAULT NULL,
  TRANSACTION_DURATION datetime DEFAULT NULL,
  SLOW_QUERY_DURATION datetime DEFAULT NULL,
  EXPENSIVE_QUERY_COUNT bigint(21) unsigned DEFAULT NULL
)`

const tableTiKVPerformanceInfo = `CREATE TABLE %s.TIKV_PERFORMANCE_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  KV_GET_COUNT bigint(21) unsigned DEFAULT NULL,
  99_KV_GET_DURATION datetime DEFAULT NULL,
  KV_BATCHGET_COUNT bigint(21) unsigned DEFAULT NULL,
  99_KV_BATCH_GET_DURATION datetime DEFAULT NULL,
  KV_SCAN_COUNT bigint(21) unsigned DEFAULT NULL,
  99_KV_SCAN_DURATION datetime DEFAULT NULL,
  KV_PREWRITE_COUNT bigint(21) unsigned DEFAULT NULL,
  99_KV_PREWRITE_DURATION datetime DEFAULT NULL,
  KV_COMMIT_COUNT bigint(21) unsigned DEFAULT NULL,
  99_KV_COMMIT_DURATION datetime DEFAULT NULL,
  KV_COPROCESSOR_COUNT bigint(21) unsigned DEFAULT NULL,
  99_KV_COPROCESSOR_DURATION datetime DEFAULT NULL,
  RAFT_STORE_CPU_USAGE double DEFAULT NULL,
  ASYNC_APPLY_CPU_USAGE double DEFAULT NULL,
  SCHEDULER_WORKER_CPU_USAGE double DEFAULT NULL
)`

const tableTiDBErrorInfo = `CREATE TABLE %s.TIDB_ERROR_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  ERR_CRITICAL_COUNT bigint(21) unsigned DEFAULT NULL,
  FAILED_QPS_COUNT bigint(21) unsigned DEFAULT NULL
)`

const tableTiKVErrorInfo = `CREATE TABLE %s.TIKV_ERROR_INFO (
  ID bigint(21) unsigned DEFAULT NULL,
  NAME varchar(64) DEFAULT NULL,
  ERR_CRITICAL_COUNT bigint(21) unsigned DEFAULT NULL,
  ERR_SERVER_IS_BUSY_COUNT bigint(21) unsigned DEFAULT NULL,
  ERR_RAFTSTORE_COUNT bigint(21) unsigned DEFAULT NULL,
  ERR_COPROCESSOR_COUNT bigint(21) unsigned DEFAULT NULL,
  ERR_GRPC_COUNT bigint(21) unsigned DEFAULT NULL,
  ERR_CHANNEL_IS_FULL_COUNT bigint(21) unsigned DEFAULT NULL
)`
