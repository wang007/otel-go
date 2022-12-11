package otelgrpc

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type TracesOption = otelgrpc.Option
type Filter = otelgrpc.Filter

type RewriteClientReporter func(ctx context.Context, cc *grpc.ClientConn, defaultReporter ClientReporter) ClientReporter
type RewriteClientSentServiceMark func(ctx context.Context, defaultServiceName string) context.Context

type RewriteServerReporter func(ctx context.Context, defaultReporter ServerReporter) ServerReporter
type RewriteServerRecvServiceMark func(ctx context.Context, defaultServiceName string) string

func WithInterceptorOptions(opts ...MetricsOption) Option {
	return optionFunc(func(c *config) {
		c.MetricsOptions = append(c.MetricsOptions, opts...)
	})
}

func WithGrpcCallCollector(g *GrpcCallCollector) Option {
	return optionFunc(func(c *config) {
		c.GrpcCallCollector = g
	})
}

func WithInterceptorFilter(f Filter) Option {
	return optionFunc(func(c *config) {
		c.TracesOptions = append(c.TracesOptions, otelgrpc.WithInterceptorFilter(f))
	})
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		c.TracesOptions = append(c.TracesOptions, otelgrpc.WithTracerProvider(tp))
	})
}

func WithPropagators(p propagation.TextMapPropagator) Option {
	return optionFunc(func(c *config) {
		c.TracesOptions = append(c.TracesOptions, otelgrpc.WithPropagators(p))
	})
}

type config struct {
	GrpcCallCollector *GrpcCallCollector
	MetricsOptions    []MetricsOption
	TracesOptions     []TracesOption
}

type Option interface {
	apply(c *config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

type metricsConfig struct {
	RewriteClientReporter        RewriteClientReporter
	RewriteClientSentServiceMark RewriteClientSentServiceMark
	RewriteServerReporter        RewriteServerReporter
	RewriteServerRecvServiceMark RewriteServerRecvServiceMark
}

type MetricsOption interface {
	apply(c *metricsConfig)
}

type interceptorOptionFunc func(*metricsConfig)

func (o interceptorOptionFunc) apply(c *metricsConfig) {
	o(c)
}

func WithRewriteClientReporter(r RewriteClientReporter) MetricsOption {
	return interceptorOptionFunc(func(i *metricsConfig) {
		i.RewriteClientReporter = r
	})
}

func WithRewriteClientSentServiceMark(r RewriteClientSentServiceMark) MetricsOption {
	return interceptorOptionFunc(func(i *metricsConfig) {
		i.RewriteClientSentServiceMark = r
	})
}

func WithRewriteServerReporter(r RewriteServerReporter) MetricsOption {
	return interceptorOptionFunc(func(i *metricsConfig) {
		i.RewriteServerReporter = r
	})
}

func WithRewriteServerRecvServiceMark(r RewriteServerRecvServiceMark) MetricsOption {
	return interceptorOptionFunc(func(i *metricsConfig) {
		i.RewriteServerRecvServiceMark = r
	})
}
