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

import "github.com/prometheus/client_golang/prometheus"

var (
	DefaultStreamCallCollector = NewPrometheusCounterStreamCallCollector()
)

func init() {
	prometheus.DefaultRegisterer.MustRegister(DefaultStreamCallCollector)
}

type StreamCounter interface {
	SentInc()
	ReceivedInc()
}

type StreamCallCollector interface {
	ActiveStreamCallCollector(passiveService, passiveMethod, methodType, protocol, status string) StreamCounter
	PassiveStreamCallCollector(activeService, passiveMethod, methodType, protocol, status string) StreamCounter
}

var _ StreamCounter = (*prometheusCallCounter)(nil)

type prometheusCallCounter struct {
	sentCounter     prometheus.Counter
	receivedCounter prometheus.Counter
}

func (p *prometheusCallCounter) SentInc() {
	p.sentCounter.Inc()
}

func (p *prometheusCallCounter) ReceivedInc() {
	p.receivedCounter.Inc()
}

var _ StreamCallCollector = (*PrometheusCounterStreamCallCollector)(nil)
var _ prometheus.Collector = (*PrometheusCounterStreamCallCollector)(nil)

type PrometheusCounterStreamCallCollector struct {
	clientStreamMsgReceived *prometheus.CounterVec
	clientStreamMsgSent     *prometheus.CounterVec

	serverStreamMsgReceived *prometheus.CounterVec
	serverStreamMsgSent     *prometheus.CounterVec
}

func (p *PrometheusCounterStreamCallCollector) ActiveStreamCallCollector(passiveService, passiveMethod, methodType, protocol, status string) StreamCounter {
	sentCounter := p.clientStreamMsgSent.WithLabelValues(passiveService, passiveMethod, methodType, protocol, status)
	receivedCounter := p.clientStreamMsgReceived.WithLabelValues(passiveService, passiveMethod, methodType, protocol, status)
	return &prometheusCallCounter{
		sentCounter:     sentCounter,
		receivedCounter: receivedCounter,
	}
}

func (p *PrometheusCounterStreamCallCollector) PassiveStreamCallCollector(activeService, passiveMethod, methodType, protocol, status string) StreamCounter {
	sentCounter := p.serverStreamMsgSent.WithLabelValues(activeService, passiveMethod, methodType, protocol, status)
	receivedCounter := p.serverStreamMsgReceived.WithLabelValues(activeService, passiveMethod, methodType, protocol, status)
	return &prometheusCallCounter{
		sentCounter:     sentCounter,
		receivedCounter: receivedCounter,
	}
}

func (p *PrometheusCounterStreamCallCollector) Describe(descs chan<- *prometheus.Desc) {
	p.clientStreamMsgReceived.Describe(descs)
	p.clientStreamMsgSent.Describe(descs)
	p.serverStreamMsgReceived.Describe(descs)
	p.serverStreamMsgSent.Describe(descs)
}

func (p *PrometheusCounterStreamCallCollector) Collect(metrics chan<- prometheus.Metric) {
	p.clientStreamMsgReceived.Collect(metrics)
	p.clientStreamMsgSent.Collect(metrics)
	p.serverStreamMsgReceived.Collect(metrics)
	p.serverStreamMsgSent.Collect(metrics)
}

func NewPrometheusCounterStreamCallCollector() *PrometheusCounterStreamCallCollector {
	return &PrometheusCounterStreamCallCollector{
		clientStreamMsgReceived: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "active_stream_received_total",
				Help: "Total number of stream messages received by the client.",
			}, []string{"passive_service", "passive_method", "method_type", "protocol", "status"}),

		clientStreamMsgSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "active_stream_sent_total",
				Help: "Total number of stream messages sent by the client.",
			}, []string{"passive_service", "passive_method", "method_type", "protocol", "status"}),

		serverStreamMsgReceived: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "passive_stream_received_total",
				Help: "Total number of stream messages received on the server.",
			}, []string{"active_service", "passive_method", "method_type", "protocol", "status"}),

		serverStreamMsgSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "passive_stream_sent_total",
				Help: "Total number of stream messages sent by the server.",
			}, []string{"active_service", "passive_method", "method_type", "protocol", "status"}),
	}
}

var _ StreamCallCollector = NoopStreamCallCollector{}

type NoopStreamCallCollector struct{}

func (n NoopStreamCallCollector) ActiveStreamCallCollector(_, _, _, _, _ string) StreamCounter {
	return noopStreamCounter
}

func (n NoopStreamCallCollector) PassiveStreamCallCollector(_, _, _, _, _ string) StreamCounter {
	return noopStreamCounter
}

var noopStreamCounter StreamCounter = NoopStreamCounter{}

type NoopStreamCounter struct{}

func (n NoopStreamCounter) SentInc() {}

func (n NoopStreamCounter) ReceivedInc() {}
