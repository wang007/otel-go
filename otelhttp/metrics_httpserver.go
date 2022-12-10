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

package otelhttp

import (
	"github.com/wang007/otel-go/metrics"
	"net/http"
	"strconv"
)

func NewMetricsHttpHandler(handler http.Handler, opts ...Option) http.Handler {
	return newMetricsHttpHandler(handler, false, opts...)
}

func NewMetrics404HttpHandler(handler http.Handler, opts ...Option) http.Handler {
	return newMetricsHttpHandler(handler, true, opts...)
}

func newMetricsHttpHandler(handler http.Handler, is4xx bool, opts ...Option) http.Handler {
	c := config{}
	for _, opt := range opts {
		opt.apply(&c)
	}
	collector := c.Collector
	if collector == nil {
		collector = metrics.DefaultHttpCallCollector
	}
	rewrite := c.RewriteServerReporter

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		resp := metrics.NewStatusCodeResponseWriter(writer)

		collector.RecordPassiveHandleAndNext(func() metrics.HttpServerReporter {
			handler.ServeHTTP(resp, request)

			var reporter metrics.HttpServerReporter
			if is4xx {
				reporter = New4xxHttpServerReporter(request, resp)
			}
			if resp.StatusCode() >= http.StatusBadRequest && resp.StatusCode() < http.StatusInternalServerError {
				reporter = New4xxHttpServerReporter(request, resp)
			} else {
				reporter = NewHttpServerReporter(request, resp)
			}

			if rewrite != nil {
				reporter = rewrite(request, resp, reporter)
			}
			return reporter
		})
	})
}

var _ metrics.HttpServerReporter = (*httpServerReporter)(nil)

type httpServerReporter struct {
	r     *http.Request
	resp  metrics.StatusCodeResponseWriter
	is404 bool
}

func (h *httpServerReporter) Err() error {
	return nil
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
	return h.r.Header.Get(metrics.ActiveServiceHeader)
}

func NewHttpServerReporter(r *http.Request, resp metrics.StatusCodeResponseWriter) metrics.HttpServerReporter {
	return &httpServerReporter{
		r:     r,
		resp:  resp,
		is404: false,
	}
}

func New4xxHttpServerReporter(r *http.Request, resp metrics.StatusCodeResponseWriter) metrics.HttpServerReporter {
	return &httpServerReporter{
		r:     r,
		resp:  resp,
		is404: true,
	}
}
