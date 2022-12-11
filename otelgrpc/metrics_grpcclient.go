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

const (
	GrpcHeaderFromServiceName = "active-service"
	unknown                   = "unknown"
)

func (g *GrpcCallCollector) UnaryClientInterceptor(opts ...MetricsOption) grpc.UnaryClientInterceptor {
	c := metricsConfig{}
	for _, opt := range opts {
		opt.apply(&c)
	}
	serviceName := g.collector.ServerInfo().ServiceName
	rewriteSent := c.RewriteClientSentServiceMark
	rewriteReporter := c.RewriteClientReporter

	interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if rewriteSent != nil {
			ctx = rewriteSent(ctx, serviceName)
		} else {
			ctx = metadata.AppendToOutgoingContext(ctx, GrpcHeaderFromServiceName, serviceName)
		}
		end := g.prepareRecordGrpcClientRequest()

		err := invoker(ctx, method, req, reply, cc, opts...)
		st, _ := status.FromError(err)
		s := strings.ToUpper(st.Code().String())

		var reporter ClientReporter = &clientReporter{
			passiveService: parsePassiveService(cc.Target()),
			method:         strings.TrimPrefix(method, "/"),
			methodType:     string(Unary),
			status:         s,
			err:            err,
		}
		if rewriteReporter != nil {
			reporter = rewriteReporter(ctx, cc, reporter)
		}
		end(reporter)
		return err
	}
	return interceptor
}

func (g *GrpcCallCollector) StreamClientInterceptor(opts ...MetricsOption) grpc.StreamClientInterceptor {
	c := metricsConfig{}
	for _, opt := range opts {
		opt.apply(&c)
	}
	serviceName := g.collector.ServerInfo().ServiceName
	rewriteSent := c.RewriteClientSentServiceMark
	rewriteReporter := c.RewriteClientReporter

	interceptor := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if rewriteSent != nil {
			ctx = rewriteSent(ctx, serviceName)
		} else {
			ctx = metadata.AppendToOutgoingContext(ctx, serviceName)
		}
		gt := clientStreamType(desc)

		end := g.prepareRecordGrpcClientRequest()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			st, _ := status.FromError(err)
			s := strings.ToUpper(st.Code().String())
			var reporter ClientReporter
			reporter = &clientReporter{
				passiveService: parsePassiveService(cc.Target()),
				method:         strings.TrimPrefix(method, "/"),
				methodType:     string(gt),
				status:         s,
				err:            err,
			}
			if rewriteReporter != nil {
				reporter = rewriteReporter(ctx, cc, reporter)
			}
			end(reporter)

			return clientStream, err
		}

		reporter := &clientReporter{
			passiveService: parsePassiveService(cc.Target()),
			method:         strings.TrimPrefix(method, "/"),
			methodType:     string(gt),
		}
		okStreamReporterCopy := *reporter
		okStreamReporterCopy.status = "OK"
		errStreamReporterCopy := *reporter
		errStreamReporterCopy.status = "ERROR"

		var okStreamReporter, errStreamReporter ClientReporter
		if rewriteReporter != nil {
			okStreamReporter = rewriteReporter(ctx, cc, &okStreamReporterCopy)
			errStreamReporter = rewriteReporter(ctx, cc, &errStreamReporterCopy)
		} else {
			okStreamReporter = &okStreamReporterCopy
			errStreamReporter = &errStreamReporterCopy
		}

		return &activeClientStream{
			ClientStream:       clientStream,
			clientReporter:     reporter,
			end:                end,
			streamOkCounter:    g.activeStreamCallCollector(okStreamReporter),
			streamErrorCounter: g.activeStreamCallCollector(errStreamReporter),
		}, nil

	}
	return interceptor
}

var _ grpc.ClientStream = (*activeClientStream)(nil)

type activeClientStream struct {
	grpc.ClientStream
	clientReporter     *clientReporter
	end                func(reporter ClientReporter)
	streamOkCounter    metrics.StreamCounter
	streamErrorCounter metrics.StreamCounter
}

func (a *activeClientStream) SendMsg(m interface{}) error {
	err := a.ClientStream.SendMsg(m)
	if err == nil {
		a.streamOkCounter.SentInc()
	} else {
		a.streamErrorCounter.SentInc()
	}
	return err
}

func (a *activeClientStream) RecvMsg(m interface{}) error {
	err := a.ClientStream.RecvMsg(m)
	if err == nil {
		a.streamOkCounter.ReceivedInc()
	} else {
		if err != io.EOF {
			a.streamErrorCounter.ReceivedInc()
		}
	}
	return err
}

func (a *activeClientStream) CloseSend() error {
	err := a.ClientStream.CloseSend()
	st, _ := status.FromError(err)
	s := strings.ToUpper(st.Code().String())

	a.clientReporter.status = s
	a.clientReporter.err = err
	a.end(a.clientReporter)

	return err
}

func parsePassiveService(target string) string {
	if strings.Contains(target, "://") {
		target = target[strings.Index(target, "://")+3:]
	}
	portIndex := strings.LastIndex(target, ":")
	if portIndex != -1 {
		if strings.Index(target, ":") != portIndex { // ipv6
			return unknown
		}
		target = target[:portIndex]
	}
	return target
}

func clientStreamType(desc *grpc.StreamDesc) GrpcType {
	if desc.ClientStreams && !desc.ServerStreams {
		return ClientStream
	} else if !desc.ClientStreams && desc.ServerStreams {
		return ServerStream
	}
	return BidiStream
}
