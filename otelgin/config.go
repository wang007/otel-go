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

package otelgin

import (
	"github.com/gin-gonic/gin"
	"github.com/wang007/otel-go/metrics"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type RewriteServerReporter func(c *gin.Context, defaultReport metrics.HttpServerReporter) metrics.HttpServerReporter

type TracesOption = otelgin.Option

type config struct {
	httpCallCollector     *metrics.HttpCallCollector
	rewriteServerReporter RewriteServerReporter
	service               string
	tracesOptions         []TracesOption
}

type Option interface {
	apply(c *config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func WithRewriteServerReporter(f RewriteServerReporter) Option {
	return optionFunc(func(config *config) {
		config.rewriteServerReporter = f
	})
}

func WithHttpCallCollector(h *metrics.HttpCallCollector) Option {
	return optionFunc(func(c *config) {
		c.httpCallCollector = h
	})
}

func WithService(service string) Option {
	return optionFunc(func(c *config) {
		c.service = service
	})
}

func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(c *config) {
		c.tracesOptions = append(c.tracesOptions, otelgin.WithPropagators(propagators))
	})
}

func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		c.tracesOptions = append(c.tracesOptions, otelgin.WithTracerProvider(provider))
	})
}
