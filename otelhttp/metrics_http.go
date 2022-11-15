package otelhttp

import (
	"github.com/wang007/otel-go/metrics"
	"net/http"
)

func NewMetricsHttpHandler(handler http.Handler, opts ...Option) http.Handler {
	return newMetricsHttpHandler(handler, false, opts...)
}

func NewMetrics404HttpHandler(handler http.Handler, opts ...Option) http.Handler {
	return newMetricsHttpHandler(handler, true, opts...)
}

func newMetricsHttpHandler(handler http.Handler, is404 bool, opts ...Option) http.Handler {
	c := config{}
	for _, opt := range opts {
		opt.apply(&c)
	}
	collector := c.collector
	if collector == nil {
		collector = metrics.DefaultCallCollector
	}

	hsc := metrics.HttpServerCallCollector{
		Collector: collector,
	}
	rws := c.rewriteStatus

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		resp := metrics.NewStatusCodeResponseWriter(writer)
		var reporter metrics.HttpServerReporter
		if is404 {
			reporter = metrics.New404HttpServerReporter(request, resp)
		} else {
			reporter = metrics.NewHttpServerReporter(request, resp)
		}

		if rws == nil {
			hsc.RecordAndNext(reporter, func() { handler.ServeHTTP(resp, request) })
		} else {
			rewriteReporter := &rewriteStatusReporter{
				HttpServerReporter: reporter,
				rewriteStatus:      rws,
				r:                  request,
				resp:               resp,
			}
			hsc.RecordAndNext(rewriteReporter, func() { handler.ServeHTTP(resp, request) })
		}
	})
}

var _ metrics.HttpServerReporter = (*rewriteStatusReporter)(nil)

type rewriteStatusReporter struct {
	metrics.HttpServerReporter
	rewriteStatus RewriteStatus

	r    *http.Request
	resp metrics.StatusCodeResponseWriter
}

func (r *rewriteStatusReporter) Status() string {
	s := r.HttpServerReporter.Status()
	return r.rewriteStatus(r.r, r.resp, s)
}
