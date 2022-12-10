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
	"errors"
	"github.com/wang007/otel-go/metrics"
	"io"
	"net"
	"net/http"
	"strconv"
	"syscall"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

func NewMetricsTransport(base http.RoundTripper, opts ...Option) http.RoundTripper {
	c := config{}
	for _, opt := range opts {
		opt.apply(&c)
	}
	collector := c.Collector
	if collector == nil {
		collector = metrics.DefaultHttpCallCollector
	}
	rewrite := c.RewriteClientReporter
	rewriteSent := c.RewriteClientSentServiceMark

	serviceName := collector.Collector.ServerInfo().ServiceName

	return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if ok, _ := req.Context().Value(metrics.DisableSendActiveServiceKey).(bool); !ok {
			if rewriteSent != nil {
				rewriteSent(req, serviceName)
			} else {
				req.Header.Set(metrics.ActiveServiceHeader, serviceName)
			}
		}

		allow, _ := req.Context().Value(metrics.AllowFromURLKey).(bool)

		passiveService, ok := req.Context().Value(metrics.PassiveServiceKey).(string)
		if !ok && allow {
			passiveService = req.URL.Hostname()
		}

		passiveMethod, ok := req.Context().Value(metrics.PassiveMethodKey).(string)
		if !ok && allow {
			passiveMethod = req.URL.Path
		}

		if base == nil {
			base = http.DefaultTransport
		}

		var resp *http.Response
		var err error
		collector.RecordActiveRequestAndNext(func() metrics.HttpClientReporter {
			resp, err = base.RoundTrip(req)
			reporter := &httpClientReporter{
				req:            req,
				resp:           resp,
				passiveMethod:  passiveMethod,
				passiveService: passiveService,
				err:            err,
			}
			if rewrite != nil {
				return rewrite(req, resp, reporter)
			}
			return reporter
		})
		return resp, err
	})
}

var _ metrics.HttpClientReporter = (*httpClientReporter)(nil)

type httpClientReporter struct {
	req            *http.Request
	resp           *http.Response
	passiveMethod  string
	passiveService string
	err            error
}

func (h *httpClientReporter) Err() error {
	return h.err
}

func (h *httpClientReporter) Method() string {
	return h.req.Method
}

func (h *httpClientReporter) Status() string {
	status := ""
	err := h.err
	if err != nil {
		status = "CONNECT_ERROR"
		if errors.Is(err, io.EOF) {
			status = "RST_ERROR"
		} else if errors.Is(err, syscall.ECONNRESET) {
			status = "RST_ERROR"
		} else if errors.Is(err, syscall.EPIPE) {
			status = "RST_ERROR"
		} else {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				status = "TIMEOUT"
			}
		}
	} else {
		status = strconv.Itoa(h.resp.StatusCode)
	}
	return status
}

func (h *httpClientReporter) Mapping() string {
	return h.passiveMethod
}

func (h *httpClientReporter) PassiveService() string {
	return h.passiveService
}
