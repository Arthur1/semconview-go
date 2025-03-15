package main

import (
	"go.opentelemetry.io/otel/attribute"
	semconv1_20_0 "go.opentelemetry.io/otel/semconv/v1.20.0"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

func main() {
	_ = semconv1_20_0.HTTPStatusCodeKey.Int(200)
	_ = semconv1_20_0.HTTPMethodKey.String("GET")
	_ = attribute.String(string(semconv1_20_0.HTTPSchemeKey), "http")
	_ = semconv.UserAgentOriginal("Mozilla/5.0")
	_ = semconv.HTTPClientRequestDurationName
}
