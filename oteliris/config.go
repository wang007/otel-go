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

package oteliris

import (
	"github.com/kataras/iris/v12"
	"github.com/wang007/otel-go/metrics"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type RewriteServerReporter func(c iris.Context, defaultReport metrics.HttpServerReporter) metrics.HttpServerReporter

// config is used to configure the iris middleware.
type config struct {
	Collector             *metrics.HttpCallCollector
	RewriteServerReporter RewriteServerReporter

	Service        string
	TracerProvider oteltrace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func WithRewriteServerReporter(f RewriteServerReporter) Option {
	return optionFunc(func(config *config) {
		config.RewriteServerReporter = f
	})
}

func WithHttpCallCollector(h *metrics.HttpCallCollector) Option {
	return optionFunc(func(c *config) {
		c.Collector = h
	})
}

func WithService(service string) Option {
	return optionFunc(func(c *config) {
		c.Service = service
	})
}

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(cfg *config) {
		if propagators != nil {
			cfg.Propagators = propagators
		}
	})
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.TracerProvider = provider
		}
	})
}
