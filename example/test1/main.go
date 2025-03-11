package main

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	semconv1_20_0 "go.opentelemetry.io/otel/semconv/v1.20.0"
	semconv1_30_0 "go.opentelemetry.io/otel/semconv/v1.30.0"
)

func main() {
	statusCode := semconv1_20_0.HTTPStatusCodeKey.Int(200)
	fmt.Println("Status Code:", statusCode)

	method := semconv1_20_0.HTTPMethodKey.String("GET")
	fmt.Println("Method:", method)

	scheme := attribute.String(string(semconv1_20_0.HTTPSchemeKey), "http")
	fmt.Println("Scheme:", scheme)

	userAgent := semconv1_30_0.UserAgentOriginal("Mozilla/5.0")
	fmt.Println("User Agent:", userAgent)

	_ = semconv1_30_0.HTTPClientRequestDurationName
}
