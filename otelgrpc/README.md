# otelgrpc
> opentelemetry 对 grpc 的支持

## 使用方式

### 引入 otel-go 做 trace provider 初始化
先引入 otel-go 包并做好初始化，参考 [otel-go README](../README.md)

### 引入 otelgrpc
```go
go get e.codingcorp.net/devops/coding-infra/otel-go/otelgrpc latest
```

###  client interceptor
```go
conn, err := grpc.Dial(*addr,
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
    grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
    grpc.WithInsecure())
```

### server interceptor
```go
s := grpc.NewServer(
	    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()), 
)
```
