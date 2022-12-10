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

var DefaultClient = &http.Client{Transport: NewTransport(http.DefaultTransport)}

func NewTransport(base http.RoundTripper, opts ...Option) http.RoundTripper {
	mt := NewMetricsTransport(base, opts...)
	return NewTracesTransport(mt, opts...)
}

func NewHandler(handler http.Handler, opts ...Option) http.Handler {
	mh := NewMetricsHttpHandler(handler, opts...)
	return NewTracesHandler(mh, opts...)
}

func New4xxHandler(handler http.Handler, opts ...Option) http.Handler {
	mh := NewMetrics404HttpHandler(handler, opts...)
	return NewTracesHandler(mh, opts...)
}
