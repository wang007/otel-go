// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func (n *NoopCallCollector) RecordActiveRequest(_, _, _, _, _ string, _ float64) {
}

func (n *NoopCallCollector) RecordPassiveHandle(_, _, _, _, _ string, _ float64) {
}

func (n *NoopCallCollector) ServerInfo() ServerInfo {
	return ServerInfo{}
}

func (n *NoopCallCollector) GetSampler(_ string) Sampler {
	return &NeverSampler{}
}
