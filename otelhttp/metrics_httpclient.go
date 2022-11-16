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
	serviceName := collector.Collector.ServerInfo().ServiceName

	rs := c.RewriteStatus
	rpm := c.RewritePassiveMethod
	rps := c.RewritePassiveService

	return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if ok, _ := req.Context().Value(metrics.DisableSendActiveServiceKey).(bool); !ok {
			req.Header.Set(metrics.ActiveServiceHeader, serviceName)
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
			return &httpClientReporter{
				req:            req,
				resp:           resp,
				passiveMethod:  passiveMethod,
				rpm:            rpm,
				passiveService: passiveService,
				rps:            rps,
				rs:             rs,
				err:            err,
			}
		})
		return resp, err
	})
}

var _ metrics.HttpClientReporter = (*httpClientReporter)(nil)

type httpClientReporter struct {
	req            *http.Request
	resp           *http.Response
	passiveMethod  string
	rpm            RewritePassiveMethod
	passiveService string
	rps            RewritePassiveService
	rs             RewriteStatus
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

	if h.rs != nil {
		rr := ResponseResult{
			StatusCode: h.resp.StatusCode,
			Header:     h.resp.Header,
			Err:        h.err,
		}
		return h.rs(h.req, rr, status)
	}
	return status
}

func (h *httpClientReporter) Mapping() string {
	m := h.passiveMethod
	if h.rpm != nil {
		m = h.rpm(h.req, m)
	}
	return m
}

func (h *httpClientReporter) PassiveService() string {
	s := h.passiveService
	if h.rps != nil {
		s = h.rps(h.req, s)
	}
	return s
}
