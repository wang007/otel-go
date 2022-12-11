package otelgorilla

import (
	"github.com/gorilla/mux"
)
import "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

//TracesMiddleware redirect go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelgorilla Middleware
func TracesMiddleware(opts ...Option) mux.MiddlewareFunc {
	c := config{}
	for _, o := range opts {
		o.apply(&c)
	}
	return otelmux.Middleware(c.service, c.tracesOptions...)
}
