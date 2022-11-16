package otelhttp

import (
	"github.com/wang007/otel-go/metrics"
	"net/http"
)

type config struct {
	Collector *metrics.HttpCallCollector

	RewriteStatus         RewriteStatus
	RewritePassiveMethod  RewritePassiveMethod
	RewritePassiveService RewritePassiveService
	RewriteActiveService  RewriteActiveService
}

type Option interface {
	apply(c *config)
}

type ResponseResult struct {
	StatusCode int
	Header     http.Header
	Err        error
}

// RewritePassiveMethod for metrics http client and server
type RewritePassiveMethod func(r *http.Request, recommendPassiveMethod string) string

// RewritePassiveService for metrics http client
type RewritePassiveService func(r *http.Request, recommendPassiveService string) string

// RewriteActiveService for metrics http server
type RewriteActiveService func(r *http.Request, recommendActiveService string) string

// RewriteStatus for rewrite metrics.CallCollector status when http code is conformable.
// eg: When http code is 200 and response body = {"code": "ERROR", "msg": "handle failed", ...}, it is actually a failure.
// So return "500" or "ERROR" by RewriteStatus to indicate the result of failure
// for metrics http client and server
type RewriteStatus func(r *http.Request, resp ResponseResult, recommendStatus string) string

type httpServerOptionFunc func(*config)

func (o httpServerOptionFunc) apply(c *config) {
	o(c)
}

func WithHttpCallCollector(collector *metrics.HttpCallCollector) Option {
	return httpServerOptionFunc(func(config *config) {
		config.Collector = collector
	})
}

func WithRewriteStatus(f RewriteStatus) Option {
	return httpServerOptionFunc(func(config *config) {
		config.RewriteStatus = f
	})
}

func WithRewritePassiveMethod(f RewritePassiveMethod) Option {
	return httpServerOptionFunc(func(config *config) {
		config.RewritePassiveMethod = f
	})
}

func WithRewritePassiveService(f RewritePassiveService) Option {
	return httpServerOptionFunc(func(config *config) {
		config.RewritePassiveService = f
	})
}

func WithRewriteActiveService(f RewriteActiveService) Option {
	return httpServerOptionFunc(func(config *config) {
		config.RewriteActiveService = f
	})
}
