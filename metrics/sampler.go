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

func (*NeverSampler) ShouldSample(_ float64, _ error) bool {
	return false
}
