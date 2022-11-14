package metrics

var _ CallCollector = (*prometheusCallCollector)(nil)

type prometheusCallCollector struct {
}

func (p *prometheusCallCollector) ServerInfo() ServerInfo {
	//TODO implement me
	panic("implement me")
}

func (p *prometheusCallCollector) RecordActiveRequest(passiveService, passiveMethod, methodType, status, protocol string, durationSec float64) {
	//TODO implement me
	panic("implement me")
}

func (p *prometheusCallCollector) RecordPassiveHandle(activeService, passiveMethod, methodType, status, protocol string, durationSec float64) {
	//TODO implement me
	panic("implement me")
}

func (p *prometheusCallCollector) GetSampler(key string) Sampler {
	//TODO implement me
	panic("implement me")
}
