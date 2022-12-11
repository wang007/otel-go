package otelgrpc

import (
	"github.com/wang007/otel-go/metrics"
	"time"
)

var (
	DefaultGrpcCallCollector = NewGrpcCallCollector(metrics.DefaultCallCollector, metrics.DefaultStreamCallCollector)
)

type GrpcType string

const (
	Unary        GrpcType = "UNARY"
	ClientStream GrpcType = "CLIENT_STREAMING"
	ServerStream GrpcType = "SERVER_STREAMING"
	BidiStream   GrpcType = "BIDI_STREAMING"
)

type ClientReporter interface {
	PassiveService() string
	Method() string
	MethodType() string
	Status() string
	Err() error
}

var _ ClientReporter = (*clientReporter)(nil)

type clientReporter struct {
	passiveService string
	method         string
	methodType     string
	status         string
	err            error
}

func (c *clientReporter) PassiveService() string {
	return c.passiveService
}

func (c *clientReporter) Method() string {
	return c.method
}

func (c *clientReporter) MethodType() string {
	return c.methodType
}

func (c *clientReporter) Status() string {
	return c.status
}

func (c *clientReporter) Err() error {
	return c.err
}

type ServerReporter interface {
	ActiveService() string
	Method() string
	MethodType() string
	Status() string
	Err() error
}

var _ ServerReporter = (*serverReporter)(nil)

type serverReporter struct {
	activeService string
	method        string
	methodType    string
	status        string
	err           error
}

func (c *serverReporter) ActiveService() string {
	return c.activeService
}

func (c *serverReporter) Method() string {
	return c.method
}

func (c *serverReporter) MethodType() string {
	return c.methodType
}

func (c *serverReporter) Status() string {
	return c.status
}

func (c *serverReporter) Err() error {
	return c.err
}

type GrpcCallCollector struct {
	collector           metrics.CallCollector
	streamCallCollector metrics.StreamCallCollector
	clientSampler       metrics.Sampler
	serverSampler       metrics.Sampler
}

func NewGrpcCallCollector(collector metrics.CallCollector, streamCollector metrics.StreamCallCollector) *GrpcCallCollector {
	return &GrpcCallCollector{
		collector:           collector,
		streamCallCollector: streamCollector,
		clientSampler:       collector.GetSampler("grpc_client"),
		serverSampler:       collector.GetSampler("grpc_server"),
	}
}

func (g *GrpcCallCollector) prepareRecordGrpcClientRequest() (end func(ClientReporter)) {
	start := time.Now()
	return func(r ClientReporter) {
		durationSec := time.Since(start).Seconds()
		if g.clientSampler.ShouldSample(durationSec, r.Err()) {
			g.collector.RecordActiveRequest(r.PassiveService(), r.Method(), r.MethodType(), r.Status(), "GRPC", durationSec)
		}
	}
}

func (g *GrpcCallCollector) prepareRecordGrpcServerHandled() (end func(ServerReporter)) {
	start := time.Now()
	return func(r ServerReporter) {
		durationSec := time.Since(start).Seconds()
		if g.serverSampler.ShouldSample(durationSec, r.Err()) {
			g.collector.RecordPassiveHandle(r.ActiveService(), r.Method(), r.MethodType(), r.Status(), "GRPC", durationSec)
		}
	}
}

func (g *GrpcCallCollector) activeStreamCallCollector(reporter ClientReporter) metrics.StreamCounter {
	return g.streamCallCollector.ActiveStreamCallCollector(
		reporter.PassiveService(),
		reporter.Method(),
		reporter.MethodType(),
		"GRPC",
		reporter.Status(),
	)
}

func (g *GrpcCallCollector) passiveStreamCallCollector(reporter ServerReporter) metrics.StreamCounter {
	return g.streamCallCollector.PassiveStreamCallCollector(
		reporter.ActiveService(),
		reporter.Method(),
		reporter.MethodType(),
		"GRPC",
		reporter.Status(),
	)
}
