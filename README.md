# GoKit Metrics

<p align="center">
  <a href="https://github.com/goxkit/metrics/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  </a>
  <a href="https://pkg.go.dev/github.com/goxkit/metrics">
    <img src="https://godoc.org/github.com/goxkit/metrics?status.svg" alt="Go Doc">
  </a>
  <a href="https://goreportcard.com/report/github.com/goxkit/metrics">
    <img src="https://goreportcard.com/badge/github.com/goxkit/metrics" alt="Go Report Card">
  </a>
  <a href="https://github.com/goxkit/metrics/actions">
    <img src="https://github.com/goxkit/metrics/actions/workflows/action.yml/badge.svg?branch=main" alt="Build Status">
  </a>
</p>

The metrics package provides tools for collecting, exporting, and monitoring metrics in Go applications with support for OpenTelemetry Protocol (OTLP) and more.

## Overview

The `metrics` package is a core component of the GoKit framework, offering a comprehensive solution for application monitoring through metrics collection. It features:

- OpenTelemetry Protocol (OTLP) metrics exporter with gRPC transport
- No-operation mode for development and testing
- HTTP middleware for collecting request metrics
- System metrics collectors for Go runtime statistics (memory, goroutines, GC, etc.)
- Easy integration with your application configuration

## Installation

```bash
go get github.com/goxkit/metrics
```

## Package Structure

```
metrics/
├── metrics.go             # Main package entry point
├── noop/                  # No-operation implementation
│   └── noop.go
├── otlp/                  # OpenTelemetry Protocol implementation
│   └── otlp.go
├── stdout/                # Standard output implementation
│   └── stdout.go
└── custom/                # Custom metrics implementations
    ├── http/              # HTTP metrics middleware
    │   └── http.go
    └── system/            # System metrics collectors
        ├── system.go
        ├── gouges_mem.go
        ├── gouges_sys.go
        └── type.go
```

## Usage

### Basic Setup

To initialize metrics in your application:

```go
import (
    "github.com/goxkit/configs"
    "github.com/goxkit/metrics"
)

func main() {
    // Create or load your application configs
    cfgs := configs.NewConfigs()

    // Install metrics provider
    provider, err := metrics.Install(cfgs)
    if err != nil {
        // Handle error
    }

    // Use the provider to create meters
    meter := provider.Meter("my-app-metrics")

    // Create instruments and collect metrics...
}
```

### HTTP Metrics Middleware

Collect metrics for HTTP requests in your application:

```go
import (
    "net/http"

    httpMetrics "github.com/goxkit/metrics/custom/http"
)

func setupHttpServer() {
    // Create the metrics middleware
    middleware, err := httpMetrics.NewHTTPMetricsMiddleware()
    if err != nil {
        // Handle error
    }

    // Create your HTTP handler
    myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    // Wrap your handler with the metrics middleware
    http.Handle("/", middleware.Handler(myHandler))
    http.ListenAndServe(":8080", nil)
}
```

### System Metrics Collection

Collect Go runtime metrics in your application:

```go
import (
    "github.com/goxkit/metrics/custom/system"
    "go.uber.org/zap"
)

func setupSystemMetrics(logger *zap.SugaredLogger) error {
    // Initialize basic metrics collectors (memory and system)
    return system.BasicMetricsCollector(logger)
}
```

## Core Components

### Main Package (`metrics.go`)

The main entry point that determines which metrics implementation to use based on configuration.

```go
provider, err := metrics.Install(configs)
```

### OTLP Implementation (`otlp/otlp.go`)

Configures the OpenTelemetry Protocol exporter for sending metrics to a collector.

### No-op Implementation (`noop/noop.go`)

A no-operation implementation useful for development and testing.

### HTTP Metrics (`custom/http/http.go`)

Middleware for collecting HTTP request metrics:
- Request counters with method, URI, and status code attributes
- Request duration histograms

### System Metrics (`custom/system/*`)

Collectors for Go runtime metrics:
- Memory usage stats (heap, GC, allocations)
- System stats (threads, goroutines, CGO calls)

## Configuration Integration

The metrics package integrates with the GoKit configs package:

```go
// Enable OTLP metrics exporter
configs.OTLPConfigs.Enabled = true
configs.OTLPConfigs.Endpoint = "localhost:4317"

// Install metrics with this configuration
provider, err := metrics.Install(configs)
```

## Best Practices

1. **Early Initialization**: Set up metrics early in your application lifecycle
2. **Proper Naming**: Use consistent naming conventions for your metrics
3. **Limited Cardinality**: Be cautious with high-cardinality labels/attributes
4. **Context Propagation**: Pass context to your metrics operations
5. **Integration**: Combine with tracing and logging for complete observability

## Integration with Other GoKit Packages

The metrics package is designed to work seamlessly with other GoKit components:
- `configs` - For configuration management
- `tracing` - For distributed tracing
- `logging` - For structured logging

## License

MIT License - See the LICENSE file for details.

## Documentation of Package Components

### metrics.go

The main entry point for the metrics package, responsible for installing the appropriate metrics provider based on configuration.

```go
func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error)
```

### noop/noop.go

Provides a no-operation implementation of the metrics provider for use in development or when metrics collection is disabled.

```go
func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error)
```

### otlp/otlp.go

Configures and installs the OpenTelemetry Protocol (OTLP) exporter for metrics collection.

```go
func Install(cfgs *configs.Configs) (*sdkmetric.MeterProvider, error)
```

### custom/http/http.go

Provides HTTP middleware for collecting request metrics, including request counts and durations.

```go
type HTTPMetricsMiddleware interface {
    Handler(next http.Handler) http.Handler
}

func NewHTTPMetricsMiddleware() (HTTPMetricsMiddleware, error)
```

### custom/system/system.go

Entry point for collecting system metrics, including memory usage and Go runtime statistics.

```go
func BasicMetricsCollector(logger *zap.SugaredLogger) error
```

### custom/system/gouges_mem.go

Collector for memory-related metrics from the Go runtime.

```go
func NewMemGauges(meter metric.Meter) (BasicGauges, error)
```

### custom/system/gouges_sys.go

Collector for system-related metrics including threads, CGO calls, and goroutines.

```go
func NewSysGauge(meter metric.Meter) (BasicGauges, error)
```

### custom/system/type.go

Defines interfaces and types for system metrics collection.

```go
type BasicGauges interface {
    Collect(meter metric.Meter)
}
```
