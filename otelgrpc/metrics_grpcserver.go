package otelgrpc

import (
	"context"
	"github.com/wang007/otel-go/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"strings"
)

func (g *GrpcCallCollector) UnaryServerInterceptor(opts ...MetricsOption) func(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	c := metricsConfig{}
	for _, opt := range opts {
		opt.apply(&c)
	}
	rewriteRecv := c.RewriteServerRecvServiceMark
	rewriteReporter := c.RewriteServerReporter

	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		var activeService string
		header, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := header.Get(GrpcHeaderFromServiceName)
			if len(values) > 0 {
				activeService = values[0]
			}
		}
		if rewriteRecv != nil {
			activeService = rewriteRecv(ctx, activeService)
		}

		end := g.prepareRecordGrpcServerHandled()

		resp, err := handler(ctx, req)
		st, _ := status.FromError(err)
		s := strings.ToUpper(st.Code().String())

		reporter := &serverReporter{
			activeService: activeService,
			method:        strings.TrimPrefix(info.FullMethod, "/"),
			methodType:    string(Unary),
			status:        s,
			err:           err,
		}

		if rewriteReporter != nil {
			end(rewriteReporter(ctx, reporter))
		} else {
			end(reporter)
		}
		return resp, err
	}
	return interceptor
}

func (g *GrpcCallCollector) StreamServerInterceptor(opts ...MetricsOption) func(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	c := metricsConfig{}
	for _, opt := range opts {
		opt.apply(&c)
	}
	rewriteRecv := c.RewriteServerRecvServiceMark
	rewriteReporter := c.RewriteServerReporter

	interceptor := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		ctx := ss.Context()
		var activeService string
		header, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := header.Get(GrpcHeaderFromServiceName)
			if len(values) > 0 {
				activeService = values[0]
			}
		}
		if rewriteRecv != nil {
			activeService = rewriteRecv(ctx, activeService)
		}

		end := g.prepareRecordGrpcServerHandled()

		reporter := &serverReporter{
			activeService: activeService,
			method:        strings.TrimPrefix(info.FullMethod, "/"),
			methodType:    string(serverStreamType(info)),
		}
		okStreamReporterCopy := *reporter
		okStreamReporterCopy.status = "OK"
		errStreamReporterCopy := *reporter
		errStreamReporterCopy.status = "ERROR"

		var okStreamReporter, errStreamReporter ServerReporter
		if rewriteReporter != nil {
			okStreamReporter = rewriteReporter(ctx, &okStreamReporterCopy)
			errStreamReporter = rewriteReporter(ctx, &errStreamReporterCopy)
		} else {
			okStreamReporter = &okStreamReporterCopy
			errStreamReporter = &errStreamReporterCopy
		}
		wrapped := &passiveServerStream{
			ServerStream:       ss,
			streamOkCounter:    g.passiveStreamCallCollector(okStreamReporter),
			streamErrorCounter: g.passiveStreamCallCollector(errStreamReporter),
		}

		err := handler(srv, wrapped)

		st, _ := status.FromError(err)
		s := strings.ToUpper(st.Code().String())
		reporter.status = s
		reporter.err = err
		if rewriteReporter != nil {
			end(rewriteReporter(ctx, reporter))
		} else {
			end(reporter)
		}

		return err
	}
	return interceptor
}

var _ grpc.ServerStream = (*passiveServerStream)(nil)

type passiveServerStream struct {
	grpc.ServerStream
	streamOkCounter    metrics.StreamCounter
	streamErrorCounter metrics.StreamCounter
}

func (p *passiveServerStream) SendMsg(m interface{}) error {
	err := p.ServerStream.SendMsg(m)
	if err == nil {
		p.streamOkCounter.SentInc()
	} else {
		p.streamErrorCounter.SentInc()
	}
	return err
}

func (p *passiveServerStream) RecvMsg(m interface{}) error {
	err := p.ServerStream.RecvMsg(m)
	if err == nil {
		p.streamOkCounter.ReceivedInc()
	} else {
		if err != io.EOF {
			p.streamErrorCounter.ReceivedInc()
		}
	}
	return err
}

func serverStreamType(info *grpc.StreamServerInfo) GrpcType {
	if info.IsClientStream && !info.IsServerStream {
		return ClientStream
	} else if !info.IsClientStream && info.IsServerStream {
		return ServerStream
	}
	return BidiStream
}
