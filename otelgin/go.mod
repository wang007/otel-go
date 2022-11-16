module github.com/wang007/otel-go/otelgin

go 1.17

replace github.com/wang007/otel-go => ../

require (
	github.com/wang007/otel-go v0.0.0-20221115185326-733040ace4ae
	github.com/gin-gonic/gin v1.8.1
	go.opentelemetry.io/otel v1.11.0
    go.opentelemetry.io/otel/trace v1.11.0
)
