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
	"github.com/prometheus/client_golang/prometheus"
)

var (
	defBuckets           = []float64{0.02, 0.050, 0.100, 0.200, 0.500, 1, 2, 5}
	DefaultCallCollector = NewCallCollector(DefaultOptions)
)

func NewCallCollector(opt Options) CallCollector {
	switch opt.CollectorType {
	case Noop:
		return &NoopCallCollector{}
	case PrometheusCounter:
		collector := NewPrometheusCounterCallCollector(opt)
		prometheus.DefaultRegisterer.MustRegister(collector)
		return collector
	case OtelMetricsHistogram:
		fallthrough
	case OtelMetricsCounter:
		fallthrough
	case PrometheusHistogram:
		fallthrough
	default:
		collector := NewPrometheusHistogramCallCollector(opt)
		prometheus.DefaultRegisterer.MustRegister(collector)
		return collector
	}
}

var _ CallCollector = (*PrometheusHistogramCallCollector)(nil)
var _ prometheus.Collector = (*PrometheusHistogramCallCollector)(nil)

type PrometheusHistogramCallCollector struct {
	active  *prometheus.HistogramVec
	passive *prometheus.HistogramVec

	serverInfo ServerInfo

	thresholdSec   float64
	onErrorSampled bool
	ratioMap       map[string]int

	defaultSampler AlwaysSampler
}

func (p *PrometheusHistogramCallCollector) Describe(descs chan<- *prometheus.Desc) {
	p.active.Describe(descs)
	p.passive.Describe(descs)
}

func (p *PrometheusHistogramCallCollector) Collect(metrics chan<- prometheus.Metric) {
	p.active.Collect(metrics)
	p.passive.Collect(metrics)
}

func (p *PrometheusHistogramCallCollector) ServerInfo() ServerInfo {
	return p.serverInfo
}

func (p *PrometheusHistogramCallCollector) RecordActiveRequest(passiveService, passiveMethod, methodType, status, protocol string, durationSec float64) {
	p.active.WithLabelValues(passiveService, passiveMethod, methodType, status, protocol).Observe(durationSec)
}

func (p *PrometheusHistogramCallCollector) RecordPassiveHandle(activeService, passiveMethod, methodType, status, protocol string, durationSec float64) {
	p.passive.WithLabelValues(activeService, passiveMethod, methodType, status, protocol).Observe(durationSec)
}

func (p *PrometheusHistogramCallCollector) GetSampler(key string) Sampler {
	v, ok := p.ratioMap[key]
	if !ok {
		return p.defaultSampler
	}
	return NewSampler(p.thresholdSec, v, p.onErrorSampled)
}

func NewPrometheusHistogramCallCollector(opt Options) *PrometheusHistogramCallCollector {
	p := &PrometheusHistogramCallCollector{
		active: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "active_requests_duration_seconds",
				Help:    "Histogram of request latency (seconds) of active.",
				Buckets: opt.HistogramBuckets,
			},
			[]string{"passive_service", "passive_method", "method_type", "status", "protocol"},
		),
		passive: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "passive_handled_duration_seconds",
				Help:    "Histogram of response latency (seconds) of passive.",
				Buckets: opt.HistogramBuckets,
			},
			[]string{"active_service", "passive_method", "method_type", "status", "protocol"},
		),

		serverInfo: opt.ServerInfo,

		thresholdSec:   opt.SamplerOptions.ThresholdSec,
		onErrorSampled: opt.SamplerOptions.OnErrorSampled,
		ratioMap:       opt.SamplerOptions.RatioMap,

		defaultSampler: AlwaysSampler{},
	}
	return p
}

var _ CallCollector = (*PrometheusCounterCallCollector)(nil)
var _ prometheus.Collector = (*PrometheusCounterCallCollector)(nil)

type PrometheusCounterCallCollector struct {
	active  *prometheus.CounterVec
	passive *prometheus.CounterVec

	serverInfo ServerInfo

	thresholdSec   float64
	onErrorSampled bool
	ratioMap       map[string]int

	defaultSampler AlwaysSampler
}

func (p *PrometheusCounterCallCollector) Describe(descs chan<- *prometheus.Desc) {
	p.active.Describe(descs)
	p.passive.Describe(descs)
}

func (p *PrometheusCounterCallCollector) Collect(metrics chan<- prometheus.Metric) {
	p.active.Collect(metrics)
	p.passive.Collect(metrics)
}

func (p *PrometheusCounterCallCollector) RecordActiveRequest(passiveService, passiveMethod, methodType, status, protocol string, _ float64) {
	p.active.WithLabelValues(passiveService, passiveMethod, methodType, status, protocol).Inc()

}

func (p *PrometheusCounterCallCollector) RecordPassiveHandle(activeService, passiveMethod, methodType, status, protocol string, _ float64) {
	p.passive.WithLabelValues(activeService, passiveMethod, methodType, status, protocol).Inc()
}

func (p *PrometheusCounterCallCollector) GetSampler(key string) Sampler {
	v, ok := p.ratioMap[key]
	if !ok {
		return p.defaultSampler
	}
	return NewSampler(p.thresholdSec, v, p.onErrorSampled)
}

func (p *PrometheusCounterCallCollector) ServerInfo() ServerInfo {
	return p.serverInfo
}

func NewPrometheusCounterCallCollector(opt Options) *PrometheusCounterCallCollector {
	p := &PrometheusCounterCallCollector{
		active: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "active_requests_total",
				Help: "Counter of request of active.",
			},
			[]string{"passive_service", "passive_method", "method_type", "status", "protocol"},
		),
		passive: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "passive_handled_total",
				Help: "Counter of response of passive.",
			},
			[]string{"active_service", "passive_method", "method_type", "status", "protocol"},
		),

		serverInfo: opt.ServerInfo,

		thresholdSec:   opt.SamplerOptions.ThresholdSec,
		onErrorSampled: opt.SamplerOptions.OnErrorSampled,
		ratioMap:       opt.SamplerOptions.RatioMap,

		defaultSampler: AlwaysSampler{},
	}
	return p
}
