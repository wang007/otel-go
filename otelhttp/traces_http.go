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

package otelhttp

import (
	"net/http"
)
import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

// redirect to go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp

var DefaultTracesClient = otelhttp.DefaultClient

func NewTracesTransport(base http.RoundTripper, opts ...Option) *otelhttp.Transport {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	return otelhttp.NewTransport(base, c.tracesOptions...)
}

func NewTracesHandler(handler http.Handler, opts ...Option) http.Handler {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	return otelhttp.NewHandler(handler, c.service, c.tracesOptions...)
}
