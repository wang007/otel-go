package otelgorilla

import (
	"github.com/wang007/otel-go/metrics"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net/http"
)

type TracesOption = otelmux.Option

type RewriteServerReporter func(req *http.Request, resp metrics.StatusCodeResponseWriter, defaultReporter metrics.HttpServerReporter) metrics.HttpServerReporter

func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(c *config) {
		c.tracesOptions = append(c.tracesOptions, otelmux.WithPropagators(propagators))
	})
}

func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		c.tracesOptions = append(c.tracesOptions, otelmux.WithTracerProvider(provider))
	})
}

func WithService(service string) Option {
	return optionFunc(func(c *config) {
		c.service = service
	})
}

func WithHttpCallCollector(h *metrics.HttpCallCollector) Option {
	return optionFunc(func(c *config) {
		c.collector = h
	})
}

func WithRewriteServerReporter(f RewriteServerReporter) Option {
	return optionFunc(func(config *config) {
		config.rewriteServerReporter = f
	})
}

type config struct {
	collector             *metrics.HttpCallCollector
	rewriteServerReporter RewriteServerReporter

	service       string
	tracesOptions []TracesOption
}

type Option interface {
	apply(c *config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}
