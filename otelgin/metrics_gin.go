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
	rewrite := c.rewriteServerReporter

	return metricsMiddlewareNext(collector, rewrite, func(g *gin.Context) { g.Next() })
}

func metricsMiddlewareNext(collector *metrics.HttpCallCollector, rewrite RewriteServerReporter, next func(*gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		collector.RecordPassiveHandleAndNext(func() metrics.HttpServerReporter {
			next(c)
			reporter := &ginReporter{
				c: c,
			}
			if rewrite != nil {
				return rewrite(c, reporter)
			}
			return reporter
		})
	}
}

var _ metrics.HttpServerReporter = (*ginReporter)(nil)

type ginReporter struct {
	c *gin.Context
}

func (r *ginReporter) Status() string {
	s := strconv.Itoa(r.c.Writer.Status())
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
	return s
}
