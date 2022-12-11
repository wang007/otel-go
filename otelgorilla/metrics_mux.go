package otelgorilla

import (
	"github.com/gorilla/mux"
	"github.com/wang007/otel-go/metrics"
	"net/http"
	"strconv"
)

func MetricsMiddleware(opts ...Option) mux.MiddlewareFunc {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	collector := c.collector
	if collector == nil {
		collector = metrics.DefaultHttpCallCollector
	}
	rewrite := c.rewriteServerReporter

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := mux.CurrentRoute(r)
			writer := metrics.NewStatusCodeResponseWriter(w)

			collector.RecordPassiveHandleAndNext(func() metrics.HttpServerReporter {
				handler.ServeHTTP(writer, r)
				reporter := &muxReporter{
					w:     writer,
					route: route,
					r:     r,
				}
				if rewrite != nil {
					return rewrite(r, writer, reporter)
				}
				return reporter
			})
		})
	}
}

var _ metrics.HttpServerReporter = (*muxReporter)(nil)

type muxReporter struct {
	w     metrics.StatusCodeResponseWriter
	route *mux.Route
	r     *http.Request
}

func (r *muxReporter) Status() string {
	s := strconv.Itoa(r.w.StatusCode())
	return s
}

func (r *muxReporter) Err() error {
	return nil
}

func (r *muxReporter) Method() string {
	return r.r.Method
}

func (r *muxReporter) Mapping() string {
	if r.route == nil { //
		return ""
	}
	template, err := r.route.GetPathTemplate()
	if err != nil {
		return "unknown"
	}
	return template
}

func (r *muxReporter) ActiveService() string {
	return r.r.Header.Get(metrics.ActiveServiceHeader)
}
