package otelgin

import (
	"github.com/gin-gonic/gin"
	"github.com/wang007/otel-go/metrics"
	"strconv"
)

// MetricsMiddleware returns a Gin measuring middleware.
func MetricsMiddleware(opts ...Option) gin.HandlerFunc {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	collector := c.httpCallCollector
	if collector == nil {
		collector = metrics.DefaultHttpCallCollector
	}
	rs := c.RewriteStatus
	ras := c.RewriteActiveService

	return metricsMiddlewareNext(collector, rs, ras, func(g *gin.Context) { g.Next() })
}

func metricsMiddlewareNext(collector *metrics.HttpCallCollector, rs RewriteStatus, ras RewriteActiveService, next func(*gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		collector.RecordPassiveHandleAndNext(func() metrics.HttpServerReporter {
			next(c)
			return &ginReporter{
				c:   c,
				rs:  rs,
				ras: ras,
			}
		})
	}
}

var _ metrics.HttpServerReporter = (*ginReporter)(nil)

type ginReporter struct {
	c   *gin.Context
	rs  RewriteStatus
	ras RewriteActiveService
}

func (r *ginReporter) Status() string {
	s := strconv.Itoa(r.c.Writer.Status())
	if r.rs != nil {
		return r.rs(r.c, s)
	}
	return s
}

func (r *ginReporter) Err() error {
	return nil
}

func (r *ginReporter) Method() string {
	return r.c.Request.Method
}

func (r *ginReporter) Mapping() string {
	return r.c.FullPath()
}

func (r *ginReporter) ActiveService() string {
	s := r.c.GetHeader(metrics.ActiveServiceHeader)
	if r.ras != nil {
		return r.ras(r.c, s)
	}
	return s
}
