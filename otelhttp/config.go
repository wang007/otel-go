package otelhttp

import (
	"github.com/wang007/otel-go/metrics"
	"net/http"
)

type config struct {
	collector     metrics.CallCollector
	rewriteStatus RewriteStatus
}

type Option interface {
	apply(c *config)
}

// RewriteStatus for rewrite metrics.CallCollector status when http code is conformable.
// eg: When http code is 200 and response body = {"code": "ERROR", "msg": "handle failed", ...}, it is actually a failure.
// So return "500" or "ERROR" by RewriteStatus to indicate the result of failure
type RewriteStatus func(r *http.Request, resp metrics.StatusCodeResponseWriter, recommendStatus string) string

type httpServerOptionFunc func(*config)

func (o httpServerOptionFunc) apply(c *config) {
	o(c)
}

func WithCallCollector(collector metrics.CallCollector) Option {
	return httpServerOptionFunc(func(config *config) {
		config.collector = collector
	})
}

func WithRewriteStatus(f RewriteStatus) Option {
	return httpServerOptionFunc(func(config *config) {
		config.rewriteStatus = f
	})
}
