# semconview-go

semconview-go is a tool to analyze Go codes and output information about the dependencies on OpenTelemetry Semantic Conventions.

## Example

```go
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
}
```

```console
$ ls
main.go
$ semconview-go list
Type       Name                 Version
attribute  http.method          v1.20.0
attribute  http.scheme          v1.20.0
attribute  http.status_code     v1.20.0
attribute  user_agent.original  v1.30.0
```

## Install

Currently we support macOS and Linux. You may be able to run otlc on Windows by using docker or go install.

### Homebrew

```sh
brew install Arthur1/tap/semconview-go
```

### Docker

- [`ghcr.io/arthur1/semconview-go`](https://github.com/Arthur1/semconview-go/pkgs/container/semconview-go)

### go install

```sh
go install github.com/Arthur1/semconview-go/cmd/semconview-go@latest
```
