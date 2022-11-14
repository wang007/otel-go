package metrics

// CallCollector Record active client and passive server metrics,
type CallCollector interface {
	// RecordActiveRequest Record active metrics. eg: http client, grpc client, redis client, mysql client...
	RecordActiveRequest(passiveService, passiveMethod, methodType, status, protocol string, durationSec float64)

	//RecordPassiveHandle Record passive metrics. eg: http server, grpc server...
	RecordPassiveHandle(activeService, passiveMethod, methodType, status, protocol string, durationSec float64)

	// GetSampler return Sampler by key.  eg: key=redis, key=http_client,
	GetSampler(key string) Sampler

	ServerInfo() ServerInfo
}

var _ CallCollector = (*NoopCallCollector)(nil)

type NoopCallCollector struct {
}

func (n *NoopCallCollector) RecordActiveRequest(passiveService, passiveMethod, methodType, status, protocol string, durationSec float64) {
}

func (n *NoopCallCollector) RecordPassiveHandle(activeService, passiveMethod, methodType, status, protocol string, durationSec float64) {
}

func (n *NoopCallCollector) ServerInfo() ServerInfo {
	return ServerInfo{}
}

func (n *NoopCallCollector) GetSampler(key string) Sampler {
	return &NeverSampler{}
}
