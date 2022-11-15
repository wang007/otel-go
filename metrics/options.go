// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"os"
	"strconv"
	"strings"
)

type CollectorType string

const (
	Noop                 = "noop"
	PrometheusHistogram  = "prometheus_histogram"
	PrometheusCounter    = "prometheus_counter"
	OtelMetricsHistogram = "otel_metrics_histogram" //TODO
	OtelMetricsCounter   = "otel_metrics_counter"   //TODO
)

var DefaultOptions = NewOptionsFromEnv()

type Options struct {
	Sampler       SamplerOptions
	ServerInfo    ServerInfo
	CollectorType CollectorType
}

type SamplerOptions struct {
	ThresholdSec   float64
	OnErrorSampled bool
	RatioMap       map[string]int
}

func NewOptionsFromEnv() Options {

	collectorType := os.Getenv("METRICS_COLLECTOR_TYPE")
	if collectorType == "" {
		collectorType = "prometheus_histogram"
	}

	thresholdSecStr := os.Getenv("METRICS_SAMPLER_THRESHOLD_SEC")
	if thresholdSecStr == "" {
		thresholdSecStr = "1"
	}
	thresholdSec, err := strconv.ParseFloat(thresholdSecStr, 10)
	if err != nil {
		panic("METRICS_SAMPLER_THRESHOLD_SEC must be integer. err: " + err.Error())
	}
	if thresholdSec < 0 {
		thresholdSec = 0 // always sampler
	}

	onErrorSampledStr := os.Getenv("METRICS_SAMPLER_ONERROR_SAMPLED")
	if onErrorSampledStr == "" {
		onErrorSampledStr = "true"
	}
	onErrorSampled, err := strconv.ParseBool(onErrorSampledStr)
	if err != nil {
		panic("METRICS_SAMPLER_ONERROR_SAMPLED must be bool. err: " + err.Error())
	}

	ratioMapStr := os.Getenv("METRICS_SAMPLER_RATIO_MAP")
	ratioMap := map[string]int{}
	if ratioMapStr != "" {
		split := strings.Split(ratioMapStr, ",")
		for _, s := range split {
			if s == "" {
				continue
			}
			kv := strings.Split(s, "=")
			if len(kv) != 2 {
				panic("METRICS_SAMPLER_RATIO_MAP must be key=value,key1=value1... format. eg: sql=10,http_client=20,db=50")
			}
			key := strings.Trim(kv[0], " ")
			v, err := strconv.Atoi(strings.Trim(kv[1], " "))
			if err != nil {
				panic("ratio must be integer and value=[0-100]")
			}
			ratioMap[key] = v
		}
	}

	serviceName := os.Getenv("METRICS_SERVICE_NAME")
	serviceInstance := os.Getenv("METRICS_SERVICE_INSTANCE")

	serverInfo := NewServerInfoInK8sCluster()
	if serviceName != "" {
		serverInfo.ServiceName = serviceName
	}
	if serviceInstance != "" {
		serverInfo.ServiceInstance = serviceInstance
	}

	return Options{
		Sampler: SamplerOptions{
			ThresholdSec:   thresholdSec,
			OnErrorSampled: onErrorSampled,
			RatioMap:       ratioMap,
		},
		ServerInfo:    serverInfo,
		CollectorType: CollectorType(collectorType),
	}
}
