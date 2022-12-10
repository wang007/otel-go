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
	gootelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// TracesMiddleware redirect go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin
func TracesMiddleware(opts ...Option) gin.HandlerFunc {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	return gootelgin.Middleware(c.service, c.tracesOptions...)
}
