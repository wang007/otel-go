package otelgrpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	g := c.GrpcCallCollector
	if g == nil {
		g = DefaultGrpcCallCollector
	}

	return grpc_middleware.ChainUnaryClient(TracesUnaryClientInterceptor(c.TracesOptions...), g.UnaryClientInterceptor(c.MetricsOptions...))
}

func StreamClientInterceptor(opts ...Option) grpc.StreamClientInterceptor {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	g := c.GrpcCallCollector
	if g == nil {
		g = DefaultGrpcCallCollector
	}
	return grpc_middleware.ChainStreamClient(TracesStreamClientInterceptor(c.TracesOptions...), g.StreamClientInterceptor(c.MetricsOptions...))
}

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	g := c.GrpcCallCollector
	if g == nil {
		g = DefaultGrpcCallCollector
	}
	return grpc_middleware.ChainUnaryServer(TracesUnaryServerInterceptor(c.TracesOptions...), g.UnaryServerInterceptor(c.MetricsOptions...))
}

func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	g := c.GrpcCallCollector
	if g == nil {
		g = DefaultGrpcCallCollector
	}
	return grpc_middleware.ChainStreamServer(TracesStreamServerInterceptor(c.TracesOptions...), g.StreamServerInterceptor(c.MetricsOptions...))
}
