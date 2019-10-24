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

package perfschema

import "github.com/pingcap/tidb/util"

// Performance Schema Name.
const (
	Name      = util.PerformanceSchemaName
	LowerName = util.PerformanceSchemaLowerName
)

// perfSchemaTables is a shortcut to involve all table names.
var perfSchemaTables = []string{
	tableGlobalStatus,
	tableSessionStatus,
	tableSetupActors,
	tableSetupObjects,
	tableSetupInstruments,
	tableSetupConsumers,
	tableStmtsCurrent,
	tableStmtsHistory,
	tableStmtsHistoryLong,
	tablePreparedStmtsInstances,
	tableTransCurrent,
	tableTransHistory,
	tableTransHistoryLong,
	tableStagesCurrent,
	tableStagesHistory,
	tableStagesHistoryLong,
	tableEventsStatementsSummaryByDigest,
	tableCpuProfile,
	tableMemoryProfile,
	tableMutexProfile,
	tableAllocsProfile,
	tableBlockProfile,
	tableGoroutine,
	tableTiKVCpuProfile,
}

// tableGlobalStatus contains the column name definitions for table global_status, same as MySQL.
const tableGlobalStatus = "CREATE TABLE performance_schema.global_status(" +
	"VARIABLE_NAME VARCHAR(64) not null," +
	"VARIABLE_VALUE VARCHAR(1024));"

// tableSessionStatus contains the column name definitions for table session_status, same as MySQL.
const tableSessionStatus = "CREATE TABLE performance_schema.session_status(" +
	"VARIABLE_NAME VARCHAR(64) not null," +
	"VARIABLE_VALUE VARCHAR(1024));"

// tableSetupActors contains the column name definitions for table setup_actors, same as MySQL.
const tableSetupActors = "CREATE TABLE if not exists performance_schema.setup_actors (" +
	"HOST			CHAR(60) NOT NULL  DEFAULT '%'," +
	"USER			CHAR(32) NOT NULL  DEFAULT '%'," +
	"ROLE			CHAR(16) NOT NULL  DEFAULT '%'," +
	"ENABLED		ENUM('YES','NO') NOT NULL  DEFAULT 'YES'," +
	"HISTORY		ENUM('YES','NO') NOT NULL  DEFAULT 'YES');"

// tableSetupObjects contains the column name definitions for table setup_objects, same as MySQL.
const tableSetupObjects = "CREATE TABLE if not exists performance_schema.setup_objects (" +
	"OBJECT_TYPE		ENUM('EVENT','FUNCTION','TABLE') NOT NULL  DEFAULT 'TABLE'," +
	"OBJECT_SCHEMA		VARCHAR(64)  DEFAULT '%'," +
	"OBJECT_NAME		VARCHAR(64) NOT NULL  DEFAULT '%'," +
	"ENABLED		ENUM('YES','NO') NOT NULL  DEFAULT 'YES'," +
	"TIMED			ENUM('YES','NO') NOT NULL  DEFAULT 'YES');"

// tableSetupInstruments contains the column name definitions for table setup_instruments, same as MySQL.
const tableSetupInstruments = "CREATE TABLE if not exists performance_schema.setup_instruments (" +
	"NAME			VARCHAR(128) NOT NULL," +
	"ENABLED		ENUM('YES','NO') NOT NULL," +
	"TIMED			ENUM('YES','NO') NOT NULL);"

// tableSetupConsumers contains the column name definitions for table setup_consumers, same as MySQL.
const tableSetupConsumers = "CREATE TABLE if not exists performance_schema.setup_consumers (" +
	"NAME			VARCHAR(64) NOT NULL," +
	"ENABLED			ENUM('YES','NO') NOT NULL);"

// tableStmtsCurrent contains the column name definitions for table events_statements_current, same as MySQL.
const tableStmtsCurrent = "CREATE TABLE if not exists performance_schema.events_statements_current (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID	BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"LOCK_TIME		BIGINT(20) UNSIGNED NOT NULL," +
	"SQL_TEXT		LONGTEXT," +
	"DIGEST			VARCHAR(32)," +
	"DIGEST_TEXT		LONGTEXT," +
	"CURRENT_SCHEMA	VARCHAR(64)," +
	"OBJECT_TYPE		VARCHAR(64)," +
	"OBJECT_SCHEMA	VARCHAR(64)," +
	"OBJECT_NAME		VARCHAR(64)," +
	"OBJECT_INSTANCE_BEGIN	BIGINT(20) UNSIGNED," +
	"MYSQL_ERRNO		INT(11)," +
	"RETURNED_SQLSTATE	VARCHAR(5)," +
	"MESSAGE_TEXT	VARCHAR(128)," +
	"ERRORS			BIGINT(20) UNSIGNED NOT NULL," +
	"WARNINGS		BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_AFFECTED	BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_SENT		BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_EXAMINED	BIGINT(20) UNSIGNED NOT NULL," +
	"CREATED_TMP_DISK_TABLES	BIGINT(20) UNSIGNED NOT NULL," +
	"CREATED_TMP_TABLES		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_FULL_JOIN		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_FULL_RANGE_JOIN	BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_RANGE	BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_RANGE_CHECK		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_SCAN		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_MERGE_PASSES		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_RANGE		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_ROWS		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_SCAN		BIGINT(20) UNSIGNED NOT NULL," +
	"NO_INDEX_USED	BIGINT(20) UNSIGNED NOT NULL," +
	"NO_GOOD_INDEX_USED		BIGINT(20) UNSIGNED NOT NULL," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE')," +
	"NESTING_EVENT_LEVEL		INT(11));"

// tableStmtsHistory contains the column name definitions for table events_statements_history, same as MySQL.
const tableStmtsHistory = "CREATE TABLE if not exists performance_schema.events_statements_history (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID		BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"LOCK_TIME		BIGINT(20) UNSIGNED NOT NULL," +
	"SQL_TEXT		LONGTEXT," +
	"DIGEST			VARCHAR(32)," +
	"DIGEST_TEXT		LONGTEXT," +
	"CURRENT_SCHEMA 	VARCHAR(64)," +
	"OBJECT_TYPE		VARCHAR(64)," +
	"OBJECT_SCHEMA	VARCHAR(64)," +
	"OBJECT_NAME		VARCHAR(64)," +
	"OBJECT_INSTANCE_BEGIN	BIGINT(20) UNSIGNED," +
	"MYSQL_ERRNO		INT(11)," +
	"RETURNED_SQLSTATE		VARCHAR(5)," +
	"MESSAGE_TEXT	VARCHAR(128)," +
	"ERRORS			BIGINT(20) UNSIGNED NOT NULL," +
	"WARNINGS		BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_AFFECTED	BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_SENT		BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_EXAMINED	BIGINT(20) UNSIGNED NOT NULL," +
	"CREATED_TMP_DISK_TABLES	BIGINT(20) UNSIGNED NOT NULL," +
	"CREATED_TMP_TABLES		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_FULL_JOIN		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_FULL_RANGE_JOIN	BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_RANGE	BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_RANGE_CHECK		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_SCAN		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_MERGE_PASSES		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_RANGE		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_ROWS		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_SCAN		BIGINT(20) UNSIGNED NOT NULL," +
	"NO_INDEX_USED	BIGINT(20) UNSIGNED NOT NULL," +
	"NO_GOOD_INDEX_USED		BIGINT(20) UNSIGNED NOT NULL," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE')," +
	"NESTING_EVENT_LEVEL		INT(11));"

// tableStmtsHistoryLong contains the column name definitions for table events_statements_history_long, same as MySQL.
const tableStmtsHistoryLong = "CREATE TABLE if not exists performance_schema.events_statements_history_long (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID	BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"LOCK_TIME		BIGINT(20) UNSIGNED NOT NULL," +
	"SQL_TEXT		LONGTEXT," +
	"DIGEST			VARCHAR(32)," +
	"DIGEST_TEXT		LONGTEXT," +
	"CURRENT_SCHEMA	VARCHAR(64)," +
	"OBJECT_TYPE		VARCHAR(64)," +
	"OBJECT_SCHEMA	VARCHAR(64)," +
	"OBJECT_NAME		VARCHAR(64)," +
	"OBJECT_INSTANCE_BEGIN	BIGINT(20) UNSIGNED," +
	"MYSQL_ERRNO		INT(11)," +
	"RETURNED_SQLSTATE		VARCHAR(5)," +
	"MESSAGE_TEXT	VARCHAR(128)," +
	"ERRORS			BIGINT(20) UNSIGNED NOT NULL," +
	"WARNINGS		BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_AFFECTED	BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_SENT		BIGINT(20) UNSIGNED NOT NULL," +
	"ROWS_EXAMINED	BIGINT(20) UNSIGNED NOT NULL," +
	"CREATED_TMP_DISK_TABLES	BIGINT(20) UNSIGNED NOT NULL," +
	"CREATED_TMP_TABLES		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_FULL_JOIN		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_FULL_RANGE_JOIN	BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_RANGE	BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_RANGE_CHECK		BIGINT(20) UNSIGNED NOT NULL," +
	"SELECT_SCAN		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_MERGE_PASSES		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_RANGE		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_ROWS		BIGINT(20) UNSIGNED NOT NULL," +
	"SORT_SCAN		BIGINT(20) UNSIGNED NOT NULL," +
	"NO_INDEX_USED	BIGINT(20) UNSIGNED NOT NULL," +
	"NO_GOOD_INDEX_USED		BIGINT(20) UNSIGNED NOT NULL," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE')," +
	"NESTING_EVENT_LEVEL		INT(11));"

// tablePreparedStmtsInstances contains the column name definitions for table prepared_statements_instances, same as MySQL.
const tablePreparedStmtsInstances = "CREATE TABLE if not exists performance_schema.prepared_statements_instances (" +
	"OBJECT_INSTANCE_BEGIN	BIGINT(20) UNSIGNED NOT NULL," +
	"STATEMENT_ID	BIGINT(20) UNSIGNED NOT NULL," +
	"STATEMENT_NAME	VARCHAR(64)," +
	"SQL_TEXT		LONGTEXT NOT NULL," +
	"OWNER_THREAD_ID	BIGINT(20) UNSIGNED NOT NULL," +
	"OWNER_EVENT_ID	BIGINT(20) UNSIGNED NOT NULL," +
	"OWNER_OBJECT_TYPE		ENUM('EVENT','FUNCTION','TABLE')," +
	"OWNER_OBJECT_SCHEMA		VARCHAR(64)," +
	"OWNER_OBJECT_NAME		VARCHAR(64)," +
	"TIMER_PREPARE	BIGINT(20) UNSIGNED NOT NULL," +
	"COUNT_REPREPARE	BIGINT(20) UNSIGNED NOT NULL," +
	"COUNT_EXECUTE	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_TIMER_EXECUTE		BIGINT(20) UNSIGNED NOT NULL," +
	"MIN_TIMER_EXECUTE		BIGINT(20) UNSIGNED NOT NULL," +
	"AVG_TIMER_EXECUTE		BIGINT(20) UNSIGNED NOT NULL," +
	"MAX_TIMER_EXECUTE		BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_LOCK_TIME	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_ERRORS		BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_WARNINGS	BIGINT(20) UNSIGNED NOT NULL," +
	"		SUM_ROWS_AFFECTED		BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_ROWS_SENT	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_ROWS_EXAMINED		BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_CREATED_TMP_DISK_TABLES	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_CREATED_TMP_TABLES	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SELECT_FULL_JOIN	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SELECT_FULL_RANGE_JOIN	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SELECT_RANGE		BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SELECT_RANGE_CHECK	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SELECT_SCAN	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SORT_MERGE_PASSES	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SORT_RANGE	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SORT_ROWS	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_SORT_SCAN	BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_NO_INDEX_USED		BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_NO_GOOD_INDEX_USED	BIGINT(20) UNSIGNED NOT NULL);"

// tableTransCurrent contains the column name definitions for table events_transactions_current, same as MySQL.
const tableTransCurrent = "CREATE TABLE if not exists performance_schema.events_transactions_current (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID	BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"STATE			ENUM('ACTIVE','COMMITTED','ROLLED BACK')," +
	"TRX_ID			BIGINT(20) UNSIGNED," +
	"GTID			VARCHAR(64)," +
	"XID_FORMAT_ID	INT(11)," +
	"XID_GTRID		VARCHAR(130)," +
	"XID_BQUAL		VARCHAR(130)," +
	"XA_STATE		VARCHAR(64)," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"ACCESS_MODE		ENUM('READ ONLY','READ WRITE')," +
	"ISOLATION_LEVEL	VARCHAR(64)," +
	"AUTOCOMMIT		ENUM('YES','NO') NOT NULL," +
	"NUMBER_OF_SAVEPOINTS	BIGINT(20) UNSIGNED," +
	"NUMBER_OF_ROLLBACK_TO_SAVEPOINT	BIGINT(20) UNSIGNED," +
	"NUMBER_OF_RELEASE_SAVEPOINT		BIGINT(20) UNSIGNED," +
	"OBJECT_INSTANCE_BEGIN	BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE'));"

// tableTransHistory contains the column name definitions for table events_transactions_history, same as MySQL.
//
const tableTransHistory = "CREATE TABLE if not exists performance_schema.events_transactions_history (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID	BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"STATE			ENUM('ACTIVE','COMMITTED','ROLLED BACK')," +
	"TRX_ID			BIGINT(20) UNSIGNED," +
	"GTID			VARCHAR(64)," +
	"XID_FORMAT_ID	INT(11)," +
	"XID_GTRID		VARCHAR(130)," +
	"XID_BQUAL		VARCHAR(130)," +
	"XA_STATE		VARCHAR(64)," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"ACCESS_MODE		ENUM('READ ONLY','READ WRITE')," +
	"ISOLATION_LEVEL	VARCHAR(64)," +
	"AUTOCOMMIT		ENUM('YES','NO') NOT NULL," +
	"NUMBER_OF_SAVEPOINTS	BIGINT(20) UNSIGNED," +
	"NUMBER_OF_ROLLBACK_TO_SAVEPOINT	BIGINT(20) UNSIGNED," +
	"NUMBER_OF_RELEASE_SAVEPOINT		BIGINT(20) UNSIGNED," +
	"OBJECT_INSTANCE_BEGIN	BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE'));"

// tableTransHistoryLong contains the column name definitions for table events_transactions_history_long, same as MySQL.
const tableTransHistoryLong = "CREATE TABLE if not exists performance_schema.events_transactions_history_long (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID	BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"STATE			ENUM('ACTIVE','COMMITTED','ROLLED BACK')," +
	"TRX_ID			BIGINT(20) UNSIGNED," +
	"GTID			VARCHAR(64)," +
	"XID_FORMAT_ID	INT(11)," +
	"XID_GTRID		VARCHAR(130)," +
	"XID_BQUAL		VARCHAR(130)," +
	"XA_STATE		VARCHAR(64)," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"ACCESS_MODE		ENUM('READ ONLY','READ WRITE')," +
	"ISOLATION_LEVEL	VARCHAR(64)," +
	"AUTOCOMMIT		ENUM('YES','NO') NOT NULL," +
	"NUMBER_OF_SAVEPOINTS	BIGINT(20) UNSIGNED," +
	"NUMBER_OF_ROLLBACK_TO_SAVEPOINT	BIGINT(20) UNSIGNED," +
	"NUMBER_OF_RELEASE_SAVEPOINT		BIGINT(20) UNSIGNED," +
	"OBJECT_INSTANCE_BEGIN	BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE'));"

// tableStagesCurrent contains the column name definitions for table events_stages_current, same as MySQL.
const tableStagesCurrent = "CREATE TABLE if not exists performance_schema.events_stages_current (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID	BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"WORK_COMPLETED	BIGINT(20) UNSIGNED," +
	"WORK_ESTIMATED	BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE'));"

// tableStagesHistory contains the column name definitions for table events_stages_history, same as MySQL.
const tableStagesHistory = "CREATE TABLE if not exists performance_schema.events_stages_history (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID	BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"WORK_COMPLETED	BIGINT(20) UNSIGNED," +
	"WORK_ESTIMATED	BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE'));"

// tableStagesHistoryLong contains the column name definitions for table events_stages_history_long, same as MySQL.
const tableStagesHistoryLong = "CREATE TABLE if not exists performance_schema.events_stages_history_long (" +
	"THREAD_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"EVENT_ID		BIGINT(20) UNSIGNED NOT NULL," +
	"END_EVENT_ID	BIGINT(20) UNSIGNED," +
	"EVENT_NAME		VARCHAR(128) NOT NULL," +
	"SOURCE			VARCHAR(64)," +
	"TIMER_START		BIGINT(20) UNSIGNED," +
	"TIMER_END		BIGINT(20) UNSIGNED," +
	"TIMER_WAIT		BIGINT(20) UNSIGNED," +
	"WORK_COMPLETED	BIGINT(20) UNSIGNED," +
	"WORK_ESTIMATED	BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_ID		BIGINT(20) UNSIGNED," +
	"NESTING_EVENT_TYPE		ENUM('TRANSACTION','STATEMENT','STAGE'));"

// tableEventsStatementsSummaryByDigest contains the column name definitions for table
// events_statements_summary_by_digest, same as MySQL.
const tableEventsStatementsSummaryByDigest = "CREATE TABLE if not exists events_statements_summary_by_digest (" +
	"SCHEMA_NAME VARCHAR(64) DEFAULT NULL," +
	"DIGEST VARCHAR(64) DEFAULT NULL," +
	"DIGEST_TEXT LONGTEXT DEFAULT NULL," +
	"EXEC_COUNT BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_LATENCY BIGINT(20) UNSIGNED NOT NULL," +
	"MAX_LATENCY BIGINT(20) UNSIGNED NOT NULL," +
	"MIN_LATENCY BIGINT(20) UNSIGNED NOT NULL," +
	"AVG_LATENCY BIGINT(20) UNSIGNED NOT NULL," +
	"SUM_ROWS_AFFECTED BIGINT(20) UNSIGNED NOT NULL," +
	"FIRST_SEEN TIMESTAMP(6) NOT NULL," +
	"LAST_SEEN TIMESTAMP(6) NOT NULL," +
	"QUERY_SAMPLE_TEXT LONGTEXT DEFAULT NULL);"

// tableCpuProfile contains the columns name definitions for table events_cpu_profile_graph
const tableCpuProfile = "CREATE TABLE IF NOT EXISTS " + tableNameCpuProfile + " (" +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"PERCENT_ABS VARCHAR(8) NOT NULL," +
	"PERCENT_REL VARCHAR(8) NOT NULL," +
	"DEPTH INT(8) NOT NULL," +
	"FILE VARCHAR(512) NOT NULL);"

// tableMemoryProfile contains the columns name definitions for table events_memory_profile_graph
const tableMemoryProfile = "CREATE TABLE IF NOT EXISTS " + tableNameMemoryProfile + " (" +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"PERCENT_ABS VARCHAR(8) NOT NULL," +
	"PERCENT_REL VARCHAR(8) NOT NULL," +
	"DEPTH INT(8) NOT NULL," +
	"FILE VARCHAR(512) NOT NULL);"

// tableMutexProfile contains the columns name definitions for table events_mutex_profile_graph
const tableMutexProfile = "CREATE TABLE IF NOT EXISTS " + tableNameMutexProfile + " (" +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"PERCENT_ABS VARCHAR(8) NOT NULL," +
	"PERCENT_REL VARCHAR(8) NOT NULL," +
	"DEPTH INT(8) NOT NULL," +
	"FILE VARCHAR(512) NOT NULL);"

// tableAllocsProfile contains the columns name definitions for table events_allocs_profile_graph
const tableAllocsProfile = "CREATE TABLE IF NOT EXISTS " + tableNameAllocsProfile + " (" +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"PERCENT_ABS VARCHAR(8) NOT NULL," +
	"PERCENT_REL VARCHAR(8) NOT NULL," +
	"DEPTH INT(8) NOT NULL," +
	"FILE VARCHAR(512) NOT NULL);"

// tableBlockProfile contains the columns name definitions for table events_block_profile_graph
const tableBlockProfile = "CREATE TABLE IF NOT EXISTS " + tableNameBlockProfile + " (" +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"PERCENT_ABS VARCHAR(8) NOT NULL," +
	"PERCENT_REL VARCHAR(8) NOT NULL," +
	"DEPTH INT(8) NOT NULL," +
	"FILE VARCHAR(512) NOT NULL);"

// tableGoroutine contains the columns name definitions for table events_goroutine
const tableGoroutine = "CREATE TABLE IF NOT EXISTS " + tableNameGoroutines + " (" +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"ID INT(8) NOT NULL," +
	"STATE VARCHAR(16) NOT NULL," +
	"LOCATION VARCHAR(512));"

const tableTiKVCpuProfile = "CREATE TABLE IF NOT EXISTS " + tableNameTiKVCpuProfile + " (" +
	"NAME VARCHAR(16) NOT NULL," +
	"ADDRESS VARCHAR(64) NOT NULL," +
	"FUNCTION VARCHAR(512) NOT NULL," +
	"PERCENT_ABS VARCHAR(8) NOT NULL," +
	"PERCENT_REL VARCHAR(8) NOT NULL," +
	"DEPTH INT(8) NOT NULL," +
	"FILE VARCHAR(512) NOT NULL);"
