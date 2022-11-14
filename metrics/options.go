package metrics

import (
	"os"
	"strconv"
	"strings"
)

var DefaultOptions = NewOptionsFromEnv()

type Options struct {
	Sampler         SamplerOptions
	ServiceName     string
	ServiceInstance string
}

type SamplerOptions struct {
	ThresholdSec   float64
	OnErrorSampled bool
	RatioMap       map[string]int
}

func NewOptionsFromEnv() Options {
	thresholdSecStr := os.Getenv("METRICS_SAMPLER_THRESHOLD_SEC")
	if thresholdSecStr == "" {
		thresholdSecStr = "1"
	}
	thresholdSec, err := strconv.ParseFloat(thresholdSecStr, 10)
	if err != nil {
		panic("METRICS_SAMPLER_THRESHOLD_SEC must be integer. err: " + err.Error())
	}
	if thresholdSec < 0 {
		thresholdSec = 0 // always sampler
	}

	onErrorSampledStr := os.Getenv("METRICS_SAMPLER_ONERROR_SAMPLED")
	if onErrorSampledStr == "" {
		onErrorSampledStr = "true"
	}
	onErrorSampled, err := strconv.ParseBool(onErrorSampledStr)
	if err != nil {
		panic("METRICS_SAMPLER_ONERROR_SAMPLED must be bool. err: " + err.Error())
	}

	ratioMapStr := os.Getenv("METRICS_SAMPLER_RATIO_MAP")
	ratioMap := map[string]int{}
	if ratioMapStr != "" {
		split := strings.Split(ratioMapStr, ",")
		for _, s := range split {
			if s == "" {
				continue
			}
			kv := strings.Split(s, "=")
			if len(kv) != 2 {
				panic("METRICS_SAMPLER_RATIO_MAP must be key=value,key1=value1... format. eg: sql=10,http_client=20,db=50")
			}
			key := strings.Trim(kv[0], " ")
			v, err := strconv.Atoi(strings.Trim(kv[1], " "))
			if err != nil {
				panic("ratio must be integer and value=[0-100]")
			}
			ratioMap[key] = v
		}
	}

	serviceName := os.Getenv("METRICS_SERVICE_NAME")
	serviceInstance := os.Getenv("METRICS_SERVICE_INSTANCE")

	return Options{
		Sampler: SamplerOptions{
			ThresholdSec:   thresholdSec,
			OnErrorSampled: onErrorSampled,
			RatioMap:       ratioMap,
		},
		ServiceName:     serviceName,
		ServiceInstance: serviceInstance,
	}
}
