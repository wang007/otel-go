package oteliris

import (
	"github.com/kataras/iris/v12"
	"github.com/wang007/otel-go/metrics"
	"strconv"
)

func MetricsMiddleware(opts ...Option) iris.Handler {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	collector := c.Collector
	if collector == nil {
		collector = metrics.DefaultHttpCallCollector
	}
	rewrite := c.RewriteServerReporter

	return metricsMiddlewareNext(collector, rewrite, func(c iris.Context) {
		c.Next()
	})
}

func metricsMiddlewareNext(collector *metrics.HttpCallCollector, rewrite RewriteServerReporter, next func(c iris.Context)) iris.Handler {
	return func(c iris.Context) {
		collector.RecordPassiveHandleAndNext(func() metrics.HttpServerReporter {
			next(c)
			r := &irisReporter{c: c}
			if rewrite != nil {
				return rewrite(c, r)
			}
			return r
		})
	}
}

var _ metrics.HttpServerReporter = (*irisReporter)(nil)

type irisReporter struct {
	c iris.Context
}

func (r *irisReporter) Status() string {
	return strconv.Itoa(r.c.GetStatusCode())
}

func (r *irisReporter) Err() error {
	return nil
}

func (r *irisReporter) Method() string {
	return r.c.Method()
}

func (r *irisReporter) Mapping() string {
	if route := r.c.GetCurrentRoute(); route != nil {
		return route.Path()
	}
	return ""
}

func (r *irisReporter) ActiveService() string {
	return r.c.GetHeader(metrics.ActiveServiceHeader)
}
