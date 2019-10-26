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

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/config"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var timeFormat = "2006-01-02 15:04:05"

var promRangeStep = 15 * time.Second

type PromDatas struct {
	Name string
	Data []PromData
}

type PromData struct {
	UnixTime int64  `json:"ts"`
	Time     string `json:"time"`
	Value    string `json:"value"`
}

func getPromData(s *model.SampleStream, tp string) []PromData {
	data := []PromData{}
	if s == nil {
		return data
	}

	for _, value := range s.Values {
		unixTime := int64(value.Timestamp) / 1000
		time := time.Unix(unixTime, 0).Format(timeFormat)

		val := ""
		if tp == "duration" {
			val = fmt.Sprintf("%.2fms", 1000*value.Value)
		} else {
			val = fmt.Sprintf("%.2f", value.Value)
		}
		data = append(data, PromData{unixTime, time, val})
	}

	return data
}

func getPromRange(api v1.API, ctx context.Context, query string, start, end time.Time, step time.Duration) (*model.SampleStream, error) {
	if !start.Before(end) {
		return nil, errors.New("start time is not before end")
	}

	v, err := api.QueryRange(ctx, query, v1.Range{start, end, step})
	if err != nil {
		return nil, errors.Trace(err)
	}

	mat, ok := v.(model.Matrix)
	if !ok {
		return nil, errors.New("query prometheus: result type mismatch")
	}
	if len(mat) == 0 {
		return nil, nil
	}

	return mat[0], nil
}

func getValue(vec model.Vector, instance string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["instance"] == model.LabelValue(instance) {
			return val.Value
		}
	}

	return model.SampleValue(math.NaN())
}

func getStatementCount(vec model.Vector, tp string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["type"] == model.LabelValue(tp) {
			return val.Value
		}
	}

	return model.SampleValue(math.NaN())
}

func getKVCount(vec model.Vector, instance string, tp string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["instance"] == model.LabelValue(instance) &&
			val.Metric["type"] == model.LabelValue(tp) {
			return val.Value
		}
	}

	return model.SampleValue(math.NaN())
}

func getKVDuration(vec model.Vector, tp string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["type"] == model.LabelValue(tp) {
			return val.Value
		}
	}

	return model.SampleValue(math.NaN())
}

func getQPSCount(vec model.Vector, instance string, result string, tp string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["instance"] == model.LabelValue(instance) &&
			val.Metric["result"] == model.LabelValue(result) &&
			val.Metric["type"] == model.LabelValue(tp) {
			return val.Value
		}
	}

	return model.SampleValue(math.NaN())
}

func getTotalQPSCount(vec model.Vector, result string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["result"] == model.LabelValue(result) {
			return val.Value
		}
	}

	return model.SampleValue(math.NaN())
}

func GetSlowQueryMetrics(client api.Client, start, end time.Time) ([]PromDatas, error) {
	promAddr := config.GetGlobalConfig().PrometheusAddr
	if promAddr == "" {
		return nil, errors.New("Invalid Prometheus Address")
	}
	client, err := api.NewClient(api.Config{
		Address: fmt.Sprintf("http://%s", promAddr),
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), promReadTimeout)
	defer cancel()

	// get parse duration.
	query := `histogram_quantile(1, sum(rate(tidb_session_parse_duration_seconds_bucket{sql_type="general"}[1m])) by (le))`
	data, err := getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results := []PromDatas{}

	result := getPromData(data, "duration")
	results = append(results, PromDatas{"parse_duration", result})

	// get compile duration.
	query = `histogram_quantile(1, sum(rate(tidb_session_compile_duration_seconds_bucket{sql_type="general"}[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"compile_duration", result})

	// get region read count.
	query = `histogram_quantile(1, sum(rate(tidb_tikvclient_txn_regions_num_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "")
	results = append(results, PromDatas{"read_region_count", result})

	// get instance read duration.
	query = `histogram_quantile(1, sum(rate(tidb_tikvclient_request_seconds_bucket{type!="GC"}[1m])) by (le, instance))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"instance_read_duration", result})

	// get backoff duration.
	query = `histogram_quantile(1, sum(rate(tidb_tikvclient_backoff_seconds_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"backoff_duration", result})

	// get kv backoff duration.
	query = `histogram_quantile(1, sum(rate(tidb_tikvclient_backoff_seconds_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"kv_backoff_duration", result})

	// get DistSQL duration
	query = `histogram_quantile(1, sum(rate(tidb_distsql_handle_query_duration_seconds_bucket[1m])) by (le, type))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"dist_sql_duration", result})

	// get scan keys count
	query = `histogram_quantile(1, sum(rate(tidb_distsql_scan_keys_num_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "")
	results = append(results, PromDatas{"scan_keys_count", result})

	// get coprocessor count.
	query = `sum(rate(tikv_grpc_msg_duration_seconds_count{type="coprocessor"}[1m])) by (type)`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "")
	results = append(results, PromDatas{"coprocessor_count", result})

	// get coprocessor duration.
	query = `histogram_quantile(1, sum(rate(tikv_grpc_msg_duration_seconds_bucket{ type="coprocessor"}[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"coprocessor_duration", result})

	// get resolve lock count
	query = `sum(rate(tikv_grpc_msg_duration_seconds_count{ type="kv_resolve_lock"}[1m]))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "")
	results = append(results, PromDatas{"resolve_lock_count", result})

	// get resolve lock duration
	query = `histogram_quantile(1, sum(rate(tikv_grpc_msg_duration_seconds_bucket{ type="kv_resolve_lock"}[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"resolve_lock_duration", result})

	// get coprocessor wait duration
	query = `histogram_quantile(1, sum(rate(tikv_coprocessor_request_wait_seconds_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"coprocessor_wait_duration", result})

	// get coprocessor handle duration
	query = `histogram_quantile(1, sum(rate(tikv_coprocessor_request_handle_seconds_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result = getPromData(data, "duration")
	results = append(results, PromDatas{"coprocessor_handle_duration", result})

	return results, nil
}
