module github.com/wang007/otel-go/otelhttp

go 1.17

replace github.com/wang007/otel-go => ../

require (
	github.com/wang007/otel-go v0.0.0-20221114180426-f6b46cd09e80
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.3
    go.opentelemetry.io/otel v1.11.0
    go.opentelemetry.io/otel/metric v0.32.3
    go.opentelemetry.io/otel/trace v1.11.0
)
