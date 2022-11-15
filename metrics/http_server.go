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
	"strconv"
	"time"
)

const (
	ActiveServiceHeader = "Active-Service"
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

type HttpServerReporter interface {
	Method() string
	Status() string
	Mapping() string
	ActiveService() string
}

var _ HttpServerReporter = (*httpServerReporter)(nil)

type httpServerReporter struct {
	r     *http.Request
	resp  StatusCodeResponseWriter
	is404 bool
}

func (h *httpServerReporter) Method() string {
	return h.r.Method
}

func (h *httpServerReporter) Status() string {
	return strconv.Itoa(h.resp.StatusCode())
}

func (h *httpServerReporter) Mapping() string {
	if h.is404 {
		return ""
	}
	return h.r.URL.Path
}

func (h *httpServerReporter) ActiveService() string {
	return h.r.Header.Get(ActiveServiceHeader)
}

func NewHttpServerReporter(r *http.Request, resp StatusCodeResponseWriter) HttpServerReporter {
	return &httpServerReporter{
		r:     r,
		resp:  resp,
		is404: false,
	}
}

func New404HttpServerReporter(r *http.Request, resp StatusCodeResponseWriter) HttpServerReporter {
	return &httpServerReporter{
		r:     r,
		resp:  resp,
		is404: true,
	}
}

type HttpServerCallCollector struct {
	Collector CallCollector
}

func (h *HttpServerCallCollector) RecordAndNext(reporter HttpServerReporter, next func()) {
	start := time.Now()
	defer func() {
		h.Collector.RecordPassiveHandle(
			reporter.ActiveService(),
			reporter.Mapping(),
			reporter.Method(),
			reporter.Status(),
			"http",
			time.Since(start).Seconds(),
		)
	}()
	next()
}
