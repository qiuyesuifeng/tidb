![](docs/logo_with_text.png)

[![Build Status](https://travis-ci.org/pingcap/tidb.svg?branch=master)](https://travis-ci.org/pingcap/tidb)
[![Go Report Card](https://goreportcard.com/badge/github.com/pingcap/tidb)](https://goreportcard.com/report/github.com/pingcap/tidb)
![GitHub release](https://img.shields.io/github/release/pingcap/tidb.svg)
[![CircleCI Status](https://circleci.com/gh/pingcap/tidb.svg?style=shield)](https://circleci.com/gh/pingcap/tidb)
[![Coverage Status](https://coveralls.io/repos/github/pingcap/tidb/badge.svg?branch=master)](https://coveralls.io/github/pingcap/tidb?branch=master)

## What is TiDB and TBSSQL?

TiDB (The pronunciation is: /'taɪdiːbi:/ tai-D-B, etymology: titanium) is an open-source distributed scalable Hybrid Transactional and Analytical Processing (HTAP) database. It features infinite horizontal scalability, strong consistency, and high availability. TiDB is MySQL compatible and serves as a one-stop data warehouse for both OLTP (Online Transactional Processing) and OLAP (Online Analytical Processing) workloads.

**When TiDB meets steaming, what will happen?**

TSBSQL is short for TiDB Streaming and Batch SQL which aims to make TiDB not only works as a relational database management system (RDBMS) but also as a streaming system.

Maybe unify Batch and Streaming SQL in future:)

## TSBSQL Features
- [x] Create Stream Table
- [x] Drop Stream Table
- [x] Show Streams
- [x] Show Create Stream
- [x] Desc Stream Table
- [x] Create Table As Select Stream 
- [x] Select Stream Table
- [x] Select Aggregation Stream Table with Fixed Window
- [x] Select Sort Aggregation Stream Table with Fixed Window
- [x] Select Stream Table Join TiDB Table with Where Conditions
- [x] Global Session Variable for Stream Table Pos
- [x] Log Stream Data Source
- [x] Kafka Stream Data Source
- [ ] Pulsar Stream Data Source

## Quick start

It is much easiser, just like build TiDB.

In TiDB Repo, run `make server`, and then `./bin/tidb-server`

## Demo show

### Ads click analysis system

1. Create user and ads TiDB table
```
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '用户 ID',
  `brand` varchar(64) NOT NULL COMMENT '手机品牌',
  `model` varchar(64) NOT NULL COMMENT '手机型号',
  `name` varchar(32) DEFAULT NULL COMMENT '用户名称',
  `province` varchar(32) DEFAULT NULL COMMENT '用户所在省会',
  `city` varchar(32) DEFAULT NULL COMMENT '用户所在城市',
  `age` int(4) DEFAULT '0' COMMENT '用户年龄',
  `create_time` timestamp NULL DEFAULT NULL COMMENT '用户创建时间',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '用户最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_key` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `ads` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ads_id` bigint(20) NOT NULL COMMENT '广告 ID',
  `title` varchar(64) NOT NULL COMMENT '广告 Title',
  `content` varchar(256) NOT NULL COMMENT '广告内容',
  `status` varchar(16) NOT NULL COMMENT '用户状态',
  `price` int(10) DEFAULT '0' COMMENT '单次广告点击收益，单位分',
  `create_time` timestamp NULL DEFAULT NULL COMMENT '用户创建时间',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '用户最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ads_id_key` (`ads_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

2. Insert some test data into user and ads table.
```
Insert Into user values(1, 101, "Apple", "5", "Test User 101", "北京", "北京", 18, now(), now());
Insert Into user values(2, 102, "Apple", "6", "Test User 102", "北京", "北京", 19, now(), now());
Insert Into user values(3, 103, "Apple", "6S", "Test User 103", "北京", "北京", 20, now(), now());
Insert Into user values(4, 104, "Apple", "7", "Test User 104", "北京", "北京", 20, now(), now());
Insert Into user values(5, 105, "Apple", "8", "Test User 105", "北京", "北京", 30, now(), now());
Insert Into user values(6, 106, "Apple", "7", "Test User 106", "北京", "北京", 32, now(), now());
Insert Into user values(7, 107, "Apple", "8", "Test User 107", "北京", "北京", 20, now(), now());
Insert Into user values(8, 108, "Apple", "6S", "Test User 108", "北京", "北京", 35, now(), now());
Insert Into user values(9, 109, "Apple", "5S", "Test User 109", "北京", "北京", 40, now(), now());
Insert Into user values(10, 110, "Apple", "5", "Test User 110", "北京", "北京", 42, now(), now());
Insert Into user values(11, 111, "Apple", "5", "Test User 111", "北京", "北京", 18, now(), now());
Insert Into user values(12, 112, "Apple", "8", "Test User 112", "北京", "北京", 19, now(), now());
Insert Into user values(13, 113, "Apple", "6S", "Test User 113", "北京", "北京", 20, now(), now());
Insert Into user values(14, 114, "Apple", "6", "Test User 114", "北京", "北京", 20, now(), now());
Insert Into user values(15, 115, "Apple", "6S", "Test User 115", "北京", "北京", 30, now(), now());
Insert Into user values(16, 116, "Apple", "7", "Test User 116", "北京", "北京", 32, now(), now());
Insert Into user values(17, 117, "Apple", "8", "Test User 117", "北京", "北京", 20, now(), now());
Insert Into user values(18, 118, "Apple", "7", "Test User 118", "北京", "北京", 35, now(), now());
Insert Into user values(19, 119, "Apple", "5S", "Test User 119", "北京", "北京", 40, now(), now());
Insert Into user values(20, 120, "Apple", "6", "Test User 120", "北京", "北京", 42, now(), now());

Insert Into ads values(1, 1001, "Test ADS 1", "Hello TiDB Hackathon 1", "normal", 300, now(), now());
Insert Into ads values(2, 1002, "Test ADS 2", "Hello TiDB Hackathon 2", "normal", 500, now(), now());
Insert Into ads values(3, 1003, "Test ADS 3", "Hello TiDB Hackathon 3", "abnormal", 400, now(), now());
Insert Into ads values(4, 1004, "Test ADS 4", "Hello TiDB Hackathon 4", "normal", 600, now(), now());
Insert Into ads values(5, 1005, "Test ADS 5", "Hello TiDB Hackathon 5", "normal", 400, now(), now());
Insert Into ads values(6, 1006, "Test ADS 6", "Hello TiDB Hackathon 6", "abnormal", 200, now(), now());
Insert Into ads values(7, 1007, "Test ADS 7", "Hello TiDB Hackathon 7", "normal", 600, now(), now());
Insert Into ads values(8, 1008, "Test ADS 8", "Hello TiDB Hackathon 8", "normal", 500, now(), now());
Insert Into ads values(9, 1009, "Test ADS 9", "Hello TiDB Hackathon 9", "normal", 300, now(), now());
Insert Into ads values(10, 1010, "Test ADS 10", "Hello TiDB Hackathon 10", "abnormal", 400, now(), now());
```

3. Create Stream Table(Only demo for log data, if you want to create kafka stream table, you should start a [collector](https://github.com/qiuyesuifeng/collector))
```
create stream tidb_stream_table_demo(click_id bigint(20), user_id bigint(20), ads_id bigint(20), click_price int(10), create_time timestamp) with ('type' =  'demo', 'topic' = 'click');
```
4. Show Create Stream
```
mysql> show create stream tidb_stream_table_demo;
+------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Table                  | Create Stream                                                                                                                                                                                                                                                                            |
+------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| tidb_stream_table_demo | CREATE STREAM `tidb_stream_table_demo` (
  `click_id` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `ads_id` bigint(20) DEFAULT NULL,
  `click_price` int(10) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT NULL) WITH (
		'topic'='click'
		'type'='demo'
		); |
+------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)

mysql> desc tidb_stream_table_demo;
+-------------+------------+------+------+---------+-------+
| Field       | Type       | Null | Key  | Default | Extra |
+-------------+------------+------+------+---------+-------+
| click_id    | bigint(20) | YES  |      | NULL    |       |
| user_id     | bigint(20) | YES  |      | NULL    |       |
| ads_id      | bigint(20) | YES  |      | NULL    |       |
| click_price | int(10)    | YES  |      | NULL    |       |
| create_time | timestamp  | YES  |      | NULL    |       |
+-------------+------------+------+------+---------+-------+
5 rows in set (0.00 sec)
```

5. Mock some test data for tidb_stream_table_demo stream table, like
```
"{"click_id":1,"user_id":111,"ads_id":1006,"click_price":300,"create_time":"2018-12-02 03:56:21"}",
"{"click_id":2,"user_id":110,"ads_id":1016,"click_price":400,"create_time":"2018-12-02 03:56:22"}",
"{"click_id":3,"user_id":114,"ads_id":1012,"click_price":200,"create_time":"2018-12-02 03:56:24"}",
...
```

6. Select Stream Table
```
mysql> set @@global.tidb_stream_table_demo_pos = 0;
Query OK, 0 rows affected (0.00 sec)
mysql> select * from tidb_stream_table_demo;
+----------+---------+--------+-------------+---------------------+
| click_id | user_id | ads_id | click_price | create_time         |
+----------+---------+--------+-------------+---------------------+
|        1 |     103 |   1006 |         200 | 2018-12-02 10:43:36 |
|        2 |     118 |   1008 |         400 | 2018-12-02 10:43:37 |
|        3 |     119 |   1006 |         700 | 2018-12-02 10:43:38 |
|        4 |     107 |   1004 |         500 | 2018-12-02 10:43:39 |
|        5 |     105 |   1007 |         700 | 2018-12-02 10:43:40 |
|        6 |     113 |   1002 |        1000 | 2018-12-02 10:43:41 |
|        7 |     109 |   1005 |        1000 | 2018-12-02 10:43:42 |
|        8 |     116 |   1005 |        1000 | 2018-12-02 10:43:43 |
|        9 |     106 |   1004 |         400 | 2018-12-02 10:43:44 |
|       10 |     103 |   1010 |         800 | 2018-12-02 10:43:45 |
+----------+---------+--------+-------------+---------------------+
10 rows in set (0.01 sec)

mysql> select * from tidb_stream_table_demo;
+----------+---------+--------+-------------+---------------------+
| click_id | user_id | ads_id | click_price | create_time         |
+----------+---------+--------+-------------+---------------------+
|       11 |     105 |   1005 |         500 | 2018-12-02 10:43:46 |
|       12 |     114 |   1010 |         800 | 2018-12-02 10:43:47 |
|       13 |     111 |   1010 |         500 | 2018-12-02 10:43:48 |
|       14 |     112 |   1001 |         100 | 2018-12-02 10:43:49 |
|       15 |     115 |   1002 |         300 | 2018-12-02 10:43:50 |
|       16 |     103 |   1004 |         500 | 2018-12-02 10:43:51 |
|       17 |     118 |   1010 |         300 | 2018-12-02 10:43:52 |
|       18 |     106 |   1003 |         800 | 2018-12-02 10:43:53 |
|       19 |     114 |   1010 |         100 | 2018-12-02 10:43:54 |
|       20 |     113 |   1002 |         200 | 2018-12-02 10:43:55 |
+----------+---------+--------+-------------+---------------------+
10 rows in set (0.00 sec)

mysql> set @@global.tidb_stream_table_demo_pos = 2;
Query OK, 0 rows affected (0.00 sec)

mysql> select * from tidb_stream_table_demo;
+----------+---------+--------+-------------+---------------------+
| click_id | user_id | ads_id | click_price | create_time         |
+----------+---------+--------+-------------+---------------------+
|        3 |     119 |   1006 |         700 | 2018-12-02 10:43:38 |
|        4 |     107 |   1004 |         500 | 2018-12-02 10:43:39 |
|        5 |     105 |   1007 |         700 | 2018-12-02 10:43:40 |
|        6 |     113 |   1002 |        1000 | 2018-12-02 10:43:41 |
|        7 |     109 |   1005 |        1000 | 2018-12-02 10:43:42 |
|        8 |     116 |   1005 |        1000 | 2018-12-02 10:43:43 |
|        9 |     106 |   1004 |         400 | 2018-12-02 10:43:44 |
|       10 |     103 |   1010 |         800 | 2018-12-02 10:43:45 |
|       11 |     105 |   1005 |         500 | 2018-12-02 10:43:46 |
|       12 |     114 |   1010 |         800 | 2018-12-02 10:43:47 |
+----------+---------+--------+-------------+---------------------+
10 rows in set (0.00 sec)

mysql> set @@global.tidb_stream_table_demo_pos = 0;
Query OK, 0 rows affected (0.00 sec)

mysql> select * from tidb_stream_table_demo where click_id > 3 and click_id < 10;
+----------+---------+--------+-------------+---------------------+
| click_id | user_id | ads_id | click_price | create_time         |
+----------+---------+--------+-------------+---------------------+
|        4 |     107 |   1004 |         500 | 2018-12-02 10:43:39 |
|        5 |     105 |   1007 |         700 | 2018-12-02 10:43:40 |
|        6 |     113 |   1002 |        1000 | 2018-12-02 10:43:41 |
|        7 |     109 |   1005 |        1000 | 2018-12-02 10:43:42 |
|        8 |     116 |   1005 |        1000 | 2018-12-02 10:43:43 |
|        9 |     106 |   1004 |         400 | 2018-12-02 10:43:44 |
+----------+---------+--------+-------------+---------------------+
6 rows in set (0.00 sec)
```

7. Select Stream with Aggregation
```
mysql> set @@global.tidb_stream_table_demo_pos = 0;
Query OK, 0 rows affected (0.00 sec)

mysql> select user_id, count(*) as cnt, max(click_price) as max_price, min(click_price) as min_price, avg(click_price) as avg_price, sum(click_price) as total_cost from tidb_stream_table_demo group by user_id window tumbling ( size 60 second) order by total_cost desc, max_price desc;
+---------+------+-----------+-----------+-----------+------------+---------------------+---------------------+
| user_id | cnt  | max_price | min_price | avg_price | total_cost | window_start        | window_end          |
+---------+------+-----------+-----------+-----------+------------+---------------------+---------------------+
|     118 |    6 |      1000 |       100 |  600.0000 |       3600 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     109 |    4 |       800 |       600 |  725.0000 |       2900 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     111 |    4 |       900 |       400 |  625.0000 |       2500 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     120 |    3 |       900 |       700 |  800.0000 |       2400 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     103 |    4 |       900 |       300 |  600.0000 |       2400 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     104 |    4 |      1000 |       300 |  575.0000 |       2300 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     110 |    3 |       900 |       600 |  766.6667 |       2300 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     117 |    3 |       900 |       500 |  666.6667 |       2000 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     102 |    4 |      1000 |       200 |  475.0000 |       1900 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     114 |    4 |       700 |       200 |  475.0000 |       1900 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     116 |    3 |       900 |       400 |  600.0000 |       1800 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     113 |    3 |       900 |       400 |  566.6667 |       1700 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     108 |    2 |       800 |       400 |  600.0000 |       1200 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     112 |    4 |       800 |       100 |  275.0000 |       1100 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     119 |    1 |       800 |       800 |  800.0000 |        800 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     115 |    1 |       600 |       600 |  600.0000 |        600 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     106 |    3 |       300 |       100 |  200.0000 |        600 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     101 |    2 |       400 |       100 |  250.0000 |        500 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     107 |    2 |       300 |       200 |  250.0000 |        500 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     105 |    1 |       300 |       300 |  300.0000 |        300 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     109 |    3 |      1000 |       800 |  866.6667 |       2600 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     103 |    4 |       900 |       100 |  650.0000 |       2600 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     120 |    3 |      1000 |       400 |  766.6667 |       2300 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     119 |    3 |       900 |       400 |  733.3333 |       2200 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     114 |    4 |       800 |       100 |  550.0000 |       2200 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     111 |    4 |       800 |       100 |  475.0000 |       1900 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     108 |    3 |       900 |       300 |  566.6667 |       1700 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     105 |    2 |       900 |       500 |  700.0000 |       1400 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     104 |    3 |       900 |       200 |  466.6667 |       1400 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     116 |    3 |       900 |       100 |  366.6667 |       1100 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     107 |    2 |       700 |       400 |  550.0000 |       1100 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     117 |    1 |       900 |       900 |  900.0000 |        900 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     101 |    1 |       700 |       700 |  700.0000 |        700 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     112 |    1 |       200 |       200 |  200.0000 |        200 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     110 |    2 |       100 |       100 |  100.0000 |        200 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
+---------+------+-----------+-----------+-----------+------------+---------------------+---------------------+
35 rows in set (0.01 sec)
```

8. Select Stream Table Join Table Table
```
mysql> set @@global.tidb_stream_table_demo_pos = 0;
Query OK, 0 rows affected (0.00 sec)

mysql> select a.click_id, b.name, b.age, c.title, c.content, c.price from tidb_stream_table_demo a join user b on a.user_id = b.user_id join ads c on a.ads_id = c.ads_id;
+----------+---------------+------+-------------+-------------------------+-------+
| click_id | name          | age  | title       | content                 | price |
+----------+---------------+------+-------------+-------------------------+-------+
|        1 | Test User 103 |   20 | Test ADS 6  | Hello TiDB Hackathon 6  |   200 |
|        2 | Test User 118 |   35 | Test ADS 8  | Hello TiDB Hackathon 8  |   500 |
|        3 | Test User 119 |   40 | Test ADS 6  | Hello TiDB Hackathon 6  |   200 |
|        4 | Test User 107 |   20 | Test ADS 4  | Hello TiDB Hackathon 4  |   600 |
|        5 | Test User 105 |   30 | Test ADS 7  | Hello TiDB Hackathon 7  |   600 |
|        6 | Test User 113 |   20 | Test ADS 2  | Hello TiDB Hackathon 2  |   500 |
|        7 | Test User 109 |   40 | Test ADS 5  | Hello TiDB Hackathon 5  |   400 |
|        8 | Test User 116 |   32 | Test ADS 5  | Hello TiDB Hackathon 5  |   400 |
|        9 | Test User 106 |   32 | Test ADS 4  | Hello TiDB Hackathon 4  |   600 |
|       10 | Test User 103 |   20 | Test ADS 10 | Hello TiDB Hackathon 10 |   400 |
+----------+---------------+------+-------------+-------------------------+-------+
10 rows in set (0.00 sec)

mysql> set @@global.tidb_stream_table_demo_pos = 0;
Query OK, 0 rows affected (0.00 sec)

mysql> select a.click_id, b.name, b.age, c.title, c.content, c.price from tidb_stream_table_demo a join user b on a.user_id = b.user_id join ads c on a.ads_id = c.ads_id and a.click_id > 3 and a.click_id < 10;
+----------+---------------+------+------------+------------------------+-------+
| click_id | name          | age  | title      | content                | price |
+----------+---------------+------+------------+------------------------+-------+
|        4 | Test User 107 |   20 | Test ADS 4 | Hello TiDB Hackathon 4 |   600 |
|        5 | Test User 105 |   30 | Test ADS 7 | Hello TiDB Hackathon 7 |   600 |
|        6 | Test User 113 |   20 | Test ADS 2 | Hello TiDB Hackathon 2 |   500 |
|        7 | Test User 109 |   40 | Test ADS 5 | Hello TiDB Hackathon 5 |   400 |
|        8 | Test User 116 |   32 | Test ADS 5 | Hello TiDB Hackathon 5 |   400 |
|        9 | Test User 106 |   32 | Test ADS 4 | Hello TiDB Hackathon 4 |   600 |
+----------+---------------+------+------------+------------------------+-------+
6 rows in set (0.01 sec)
```

9. Create Table as Select Stream
```
mysql> create table click_tmp as select user_id, count(*) as cnt, max(click_price) as max_price, min(click_price) as min_price, avg(click_price) as avg_price, sum(click_price) as total_cost from tidb_stream_table_demo group by user_id window tumbling ( size 60 second) order by total_cost desc, max_price desc;
Query OK, 35 rows affected (0.02 sec)

mysql> select * from click_tmp;                                                                                                                                                   +---------+------+-----------+-----------+-----------+------------+---------------------+---------------------+
| user_id | cnt  | max_price | min_price | avg_price | total_cost | window_start        | window_end          |
+---------+------+-----------+-----------+-----------+------------+---------------------+---------------------+
|     118 |    6 |      1000 |       100 |  600.0000 |       3600 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     109 |    4 |       800 |       600 |  725.0000 |       2900 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     111 |    4 |       900 |       400 |  625.0000 |       2500 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     120 |    3 |       900 |       700 |  800.0000 |       2400 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     103 |    4 |       900 |       300 |  600.0000 |       2400 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     104 |    4 |      1000 |       300 |  575.0000 |       2300 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     110 |    3 |       900 |       600 |  766.6667 |       2300 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     117 |    3 |       900 |       500 |  666.6667 |       2000 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     102 |    4 |      1000 |       200 |  475.0000 |       1900 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     114 |    4 |       700 |       200 |  475.0000 |       1900 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     116 |    3 |       900 |       400 |  600.0000 |       1800 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     113 |    3 |       900 |       400 |  566.6667 |       1700 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     108 |    2 |       800 |       400 |  600.0000 |       1200 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     112 |    4 |       800 |       100 |  275.0000 |       1100 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     119 |    1 |       800 |       800 |  800.0000 |        800 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     115 |    1 |       600 |       600 |  600.0000 |        600 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     106 |    3 |       300 |       100 |  200.0000 |        600 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     101 |    2 |       400 |       100 |  250.0000 |        500 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     107 |    2 |       300 |       200 |  250.0000 |        500 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     105 |    1 |       300 |       300 |  300.0000 |        300 | 2018-12-02 10:52:04 | 2018-12-02 10:53:04 |
|     109 |    3 |      1000 |       800 |  866.6667 |       2600 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     103 |    4 |       900 |       100 |  650.0000 |       2600 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     120 |    3 |      1000 |       400 |  766.6667 |       2300 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     119 |    3 |       900 |       400 |  733.3333 |       2200 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     114 |    4 |       800 |       100 |  550.0000 |       2200 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     111 |    4 |       800 |       100 |  475.0000 |       1900 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     108 |    3 |       900 |       300 |  566.6667 |       1700 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     105 |    2 |       900 |       500 |  700.0000 |       1400 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     104 |    3 |       900 |       200 |  466.6667 |       1400 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     116 |    3 |       900 |       100 |  366.6667 |       1100 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     107 |    2 |       700 |       400 |  550.0000 |       1100 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     117 |    1 |       900 |       900 |  900.0000 |        900 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     101 |    1 |       700 |       700 |  700.0000 |        700 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     112 |    1 |       200 |       200 |  200.0000 |        200 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
|     110 |    2 |       100 |       100 |  100.0000 |        200 | 2018-12-02 10:53:05 | 2018-12-02 10:54:05 |
+---------+------+-----------+-----------+-----------+------------+---------------------+---------------------+
35 rows in set (0.00 sec)
```

## TIDB Architecture

![architecture](./docs/architecture.png)

## Contributing
Contributions are welcomed and greatly appreciated. See [CONTRIBUTING.md](CONTRIBUTING.md)
for details on submitting patches and the contribution workflow.

## Connect with us

- [**Contact PingCAP Team**](http://bit.ly/contact_us_via_github)
- **Twitter**: [@PingCAP](https://twitter.com/PingCAP)
- **Reddit**: https://www.reddit.com/r/TiDB/
- **Stack Overflow**: https://stackoverflow.com/questions/tagged/tidb
- **Mailing list**: [Google Group](https://groups.google.com/forum/#!forum/tidb-user)

## License
TiDB is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.

## Acknowledgments
- Thanks [cznic](https://github.com/cznic) for providing some great open source tools.
- Thanks [GolevelDB](https://github.com/syndtr/goleveldb), [BoltDB](https://github.com/boltdb/bolt), and [RocksDB](https://github.com/facebook/rocksdb) for their powerful storage engines.
