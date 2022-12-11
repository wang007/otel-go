package otelgorilla

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Middleware(opts ...Option) mux.MiddlewareFunc {
	metricsMiddleware := MetricsMiddleware(opts...)
	tracesMiddleware := TracesMiddleware(opts...)
	return func(handler http.Handler) http.Handler {
		return tracesMiddleware(metricsMiddleware(handler))
	}
}
