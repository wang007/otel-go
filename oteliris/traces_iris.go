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
	"fmt"
	"github.com/kataras/iris/v12"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/wang007/otel-go/oteliris"
)

// TracesMiddleware sets up a handler to start tracing the incoming
// requests.  The service parameter should describe the name of the
// (virtual) server handling the request.
func TracesMiddleware(opts ...Option) iris.Handler {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(SemVersion()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	return func(ic iris.Context) {
		savedCtx := ic.Request().Context()
		defer func() {
			ic.ResetRequest(ic.Request().WithContext(savedCtx))
		}()
		ctx := cfg.Propagators.Extract(savedCtx, propagation.HeaderCarrier(ic.Request().Header))
		spanName := ""
		//TODO enhance O(n) -> O(1), iris v12.2.x version fix it
		if route := ic.GetCurrentRoute(); route != nil {
			spanName = route.Path()
		}
		routeStr := spanName
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", ic.Request().Method)
		}

		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", ic.Request())...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(ic.Request())...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(cfg.Service, routeStr, ic.Request())...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		ic.ResetRequest(ic.Request().WithContext(ctx))

		// serve the request to the next middleware
		ic.Next()

		status := ic.ResponseWriter().StatusCode()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(status, oteltrace.SpanKindServer)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
	}
}
