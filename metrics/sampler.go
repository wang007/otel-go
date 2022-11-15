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

import (
	"math/rand"
	"time"
)

// Sampler determine whether to sample
type Sampler interface {
	ShouldSample(durationSec float64, err error) bool
}

type sampler struct {
	thresholdSec   float64
	onErrorSampled bool
	ratio          uint8
	rand           *rand.Rand
}

func (s *sampler) ShouldSample(durationSec float64, err error) bool {
	if err != nil && s.onErrorSampled {
		return true
	}
	if durationSec > s.thresholdSec {
		return true
	}
	r := s.ratio
	if r == 0 {
		return false
	}
	if r == 100 {
		return true
	}
	scope := rand.Intn(101)
	return scope <= int(r)
}

func NewSampler(thresholdSec float64, ratio int, onErrorSampled bool) Sampler {
	if ratio > 100 {
		ratio = 100
	}
	if ratio < 0 {
		ratio = 0
	}
	return &sampler{
		thresholdSec:   thresholdSec,
		onErrorSampled: onErrorSampled,
		ratio:          uint8(ratio),
		rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type NeverSampler struct{}

func (NeverSampler) ShouldSample(_ float64, _ error) bool {
	return false
}

type AlwaysSampler struct{}

func (AlwaysSampler) ShouldSample(_ float64, _ error) bool {
	return true
}
