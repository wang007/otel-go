package oteliris

import (
	"github.com/kataras/iris/v12"
	"github.com/wang007/otel-go/metrics"
)

func Middleware(opts ...Option) iris.Handler {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	collector := c.Collector
	if collector == nil {
		collector = metrics.DefaultHttpCallCollector
	}
	rewrite := c.RewriteServerReporter

	tm := TracesMiddleware(opts...)
	return metricsMiddlewareNext(collector, rewrite, func(c iris.Context) {
		tm(c)
	})
}
