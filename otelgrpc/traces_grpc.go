package otelgrpc

import (
	"google.golang.org/grpc"
)
import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

// TracesUnaryClientInterceptor redirect go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
func TracesUnaryClientInterceptor(opts ...TracesOption) grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(opts...)
}

// TracesStreamClientInterceptor redirect go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
func TracesStreamClientInterceptor(opts ...TracesOption) grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor(opts...)
}

// TracesUnaryServerInterceptor redirect go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
func TracesUnaryServerInterceptor(opts ...TracesOption) grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor(opts...)
}

// TracesStreamServerInterceptor redirect go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
func TracesStreamServerInterceptor(opts ...TracesOption) grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor(opts...)
}
