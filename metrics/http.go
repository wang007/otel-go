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
	"net/http"
	"time"
)

const (
	ActiveServiceHeader = "Active-Service"

	DisableSendActiveServiceKey = "__d_s_a_s_k"
	AllowFromURLKey             = "__a_f_u_k"
	PassiveMethodKey            = "__p_m_k"
	PassiveServiceKey           = "__p_s_k"
)

var (
	DefaultHttpCallCollector = NewHttpCallCollector(DefaultCallCollector)
)

type StatusCodeResponseWriter interface {
	http.ResponseWriter
	StatusCode() int
}

type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusCodeResponseWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func (s *statusCodeResponseWriter) StatusCode() int {
	return s.statusCode
}

func NewStatusCodeResponseWriter(response http.ResponseWriter) StatusCodeResponseWriter {
	return &statusCodeResponseWriter{
		ResponseWriter: response,
		statusCode:     http.StatusOK,
	}
}

type HttpClientReporter interface {
	Method() string
	Status() string
	Mapping() string
	PassiveService() string
	Err() error
}

type HttpServerReporter interface {
	Method() string
	Status() string
	Mapping() string
	ActiveService() string
	Err() error
}

func NewHttpCallCollector(collector CallCollector) *HttpCallCollector {
	return &HttpCallCollector{
		Collector:         collector,
		httpClientSampler: collector.GetSampler("http_client"),
		httpServerSampler: collector.GetSampler("http_server"),
	}
}

type HttpCallCollector struct {
	Collector         CallCollector
	httpClientSampler Sampler
	httpServerSampler Sampler
}

func (h *HttpCallCollector) RecordActiveRequestAndNext(next func() HttpClientReporter) {
	start := time.Now()
	reporter := next()
	durationSec := time.Since(start).Seconds()
	if h.httpClientSampler.ShouldSample(durationSec, reporter.Err()) {
		methodType := reporter.Method()
		status := reporter.Status()
		h.Collector.RecordActiveRequest(reporter.PassiveService(), reporter.Mapping(), methodType, status, "http", durationSec)
	}
}

func (h *HttpCallCollector) RecordPassiveHandleAndNext(next func() HttpServerReporter) {
	start := time.Now()
	reporter := next()
	durationSec := time.Since(start).Seconds()
	if h.httpServerSampler.ShouldSample(durationSec, reporter.Err()) {
		methodType := reporter.Method()
		status := reporter.Status()
		h.Collector.RecordPassiveHandle(reporter.ActiveService(), reporter.Mapping(), methodType, status, "http", durationSec)
	}
}
