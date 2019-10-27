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
	"sort"
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
	Abs  string
}

type PromData struct {
	UnixTime int64  `json:"-"`
	Time     string `json:"time"`
	Value    string `json:"value"`
}

func getPromData(name string, s *model.SampleStream, tp string) PromDatas {
	data := PromDatas{Name: name, Data: []PromData{}}
	if s == nil {
		return data
	}

	var values []float64
	var total float64
	var avg, min, p80, p90, p95, max float64
	for _, value := range s.Values {
		unixTime := int64(value.Timestamp) / 1000
		time := time.Unix(unixTime, 0).Format(timeFormat)
		val := ""
		if tp == "duration" {
			val = fmt.Sprintf("%.2fms", 1000*value.Value)
		} else {
			val = fmt.Sprintf("%.2f", value.Value)
		}
		data.Data = append(data.Data, PromData{unixTime, time, val})
		if v := float64(value.Value); !math.IsNaN(v) {
			values = append(values, v)
			total += v
		}
	}
	sort.Sort(sort.Float64Slice(values))
	if count := float64(len(values)); count > 0 {
		avg = total / count
		min = values[0]
		p80 = values[int(count*0.8)]
		p90 = values[int(count*0.9)]
		p95 = values[int(count*0.95)]
		max = values[len(values)-1]
	}
	if tp == "duration" {
		data.Abs = fmt.Sprintf("avg: %.2fms, min: %.2fms, max: %.2fms, p80: %.2fms, p90: %.2fms, p95: %.2fms",
			avg*1000, min*1000, max*1000, p80*1000, p90*1000, p95*1000)
	} else {
		data.Abs = fmt.Sprintf("avg: %.2f, min: %.2f, max: %.2f, p80: %.2f, p90: %.2f, p95: %.2f",
			avg, min, max, p80, p90, p95)
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

	results = append(results, getPromData("parse_duration", data, "duration"))

	// get compile duration.
	query = `histogram_quantile(1, sum(rate(tidb_session_compile_duration_seconds_bucket{sql_type="general"}[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("compile_duration", data, "duration"))

	// get region read count.
	query = `histogram_quantile(1, sum(rate(tidb_tikvclient_txn_regions_num_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("read_region_count", data, ""))

	// get instance read duration.
	query = `histogram_quantile(1, sum(rate(tidb_tikvclient_request_seconds_bucket{type!="GC"}[1m])) by (le, instance))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("instance_read_duration", data, "duration"))

	// get backoff duration.
	query = `histogram_quantile(1, sum(rate(tidb_tikvclient_backoff_seconds_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("backoff_duration", data, "duration"))

	// get kv backoff duration.
	query = `histogram_quantile(1, sum(rate(tidb_tikvclient_backoff_seconds_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("kv_backoff_duration", data, "duration"))

	// get DistSQL duration
	query = `histogram_quantile(1, sum(rate(tidb_distsql_handle_query_duration_seconds_bucket[1m])) by (le, type))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("dist_sql_duration", data, "duration"))

	// get scan keys count
	query = `histogram_quantile(1, sum(rate(tidb_distsql_scan_keys_num_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("scan_keys_count", data, ""))

	// get coprocessor count.
	query = `sum(rate(tikv_grpc_msg_duration_seconds_count{type="coprocessor"}[1m])) by (type)`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("coprocessor_count", data, ""))

	// get coprocessor duration.
	query = `histogram_quantile(1, sum(rate(tikv_grpc_msg_duration_seconds_bucket{ type="coprocessor"}[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("coprocessor_duration", data, "duration"))

	// get resolve lock count
	query = `sum(rate(tikv_grpc_msg_duration_seconds_count{ type="kv_resolve_lock"}[1m]))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("resolve_lock_count", data, ""))

	// get resolve lock duration
	query = `histogram_quantile(1, sum(rate(tikv_grpc_msg_duration_seconds_bucket{ type="kv_resolve_lock"}[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("resolve_lock_duration", data, "duration"))

	// get coprocessor wait duration
	query = `histogram_quantile(1, sum(rate(tikv_coprocessor_request_wait_seconds_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("coprocessor_wait_duration", data, "duration"))

	// get coprocessor handle duration
	query = `histogram_quantile(1, sum(rate(tikv_coprocessor_request_handle_seconds_bucket[1m])) by (le))`
	data, err = getPromRange(api, ctx, query, start, end, promRangeStep)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results = append(results, getPromData("coprocessor_handle_duration", data, "duration"))

	return results, nil
}
