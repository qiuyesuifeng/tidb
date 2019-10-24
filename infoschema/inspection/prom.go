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
	"github.com/prometheus/common/model"
)

func getValue(vec model.Vector, instance string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["instance"] == model.LabelValue(instance) {
			return val.Value
		}
	}

	return model.SampleValue(-1.0)
}

func getStatementCount(vec model.Vector, tp string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["type"] == model.LabelValue(tp) {
			return val.Value
		}
	}

	return model.SampleValue(-1.0)
}

func getQPSCount(vec model.Vector, instance string, result string, tp string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["instance"] == model.LabelValue(instance) &&
			val.Metric["result"] == model.LabelValue(result) &&
			val.Metric["type"] == model.LabelValue(tp) {
			return val.Value
		}
	}

	return model.SampleValue(-1.0)
}

func getTotalQPSCount(vec model.Vector, result string) model.SampleValue {
	for _, val := range vec {
		if val.Metric["result"] == model.LabelValue(result) {
			return val.Value
		}
	}

	return model.SampleValue(-1.0)
}
