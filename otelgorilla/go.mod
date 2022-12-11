module github.com/wang007/otel-go/otelgorilla

go 1.17

replace github.com/wang007/otel-go => ../

require (
	github.com/gorilla/mux v1.8.0
	github.com/wang007/otel-go v0.0.0-20221115185326-733040ace4ae
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.36.3
	go.opentelemetry.io/otel v1.11.0
	go.opentelemetry.io/otel/trace v1.11.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	golang.org/x/sys v0.0.0-20221006211917-84dc82d7e875 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
