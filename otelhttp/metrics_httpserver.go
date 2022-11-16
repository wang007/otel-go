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

	rs := c.RewriteStatus
	rpm := c.RewritePassiveMethod
	ras := c.RewriteActiveService

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
			if rs == nil && rpm == nil && ras == nil {
				return reporter
			}
			return &rewriteHttpServerReporter{
				HttpServerReporter: reporter,
				rs:                 rs,
				rpm:                rpm,
				ras:                ras,
				r:                  request,
				resp:               resp,
			}
		})
	})
}

var _ metrics.HttpServerReporter = (*rewriteHttpServerReporter)(nil)

type rewriteHttpServerReporter struct {
	metrics.HttpServerReporter
	rs  RewriteStatus
	rpm RewritePassiveMethod
	ras RewriteActiveService

	r    *http.Request
	resp metrics.StatusCodeResponseWriter
}

func (r *rewriteHttpServerReporter) Status() string {
	status := r.HttpServerReporter.Status()
	if r.rs != nil {
		rr := ResponseResult{
			StatusCode: r.resp.StatusCode(),
			Header:     r.resp.Header(),
			Err:        nil,
		}
		return r.rs(r.r, rr, status)
	}
	return status
}

func (r *rewriteHttpServerReporter) Mapping() string {
	m := r.HttpServerReporter.Mapping()
	if r.rpm != nil {
		return r.rpm(r.r, m)
	}
	return m
}

func (r *rewriteHttpServerReporter) ActiveService() string {
	s := r.HttpServerReporter.ActiveService()
	if r.ras != nil {
		return r.ras(r.r, s)
	}
	return s
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
