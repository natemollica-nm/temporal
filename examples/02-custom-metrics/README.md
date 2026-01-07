# Example 02: Custom Metrics Exporters

This example demonstrates how to use custom metrics exporters with Temporal, specifically showing DogStatsD integration alongside the default Prometheus exporter.

## What it demonstrates

- **Custom Metrics Reporters**: How to implement and use custom metrics exporters
- **DogStatsD Integration**: Sending Temporal metrics to DataDog via StatsD protocol
- **Configuration-based Metrics**: Switching between different metrics providers via configuration
- **Metrics Factory Pattern**: Clean abstraction for different metrics backends

## Metrics Providers Supported

### Prometheus (Default)
- Exposes metrics on HTTP endpoint (default: `:9090`)
- Compatible with Prometheus scraping
- Includes histograms and detailed metrics

### DogStatsD
- Sends metrics via UDP to DataDog agent
- Uses DataDog's tag format (`key:value`)
- Configurable flush intervals and buffer sizes

## Configuration

Set the metrics provider in your configuration:

```yaml
metrics:
  provider: "dogstatsd"  # or "prometheus"
  dogstatsd:
    hostPort: "127.0.0.1:8125"
    flushInterval: 1s
    flushBytes: 1432
```

## Running with DogStatsD

1. Start DataDog agent (or StatsD-compatible service)
2. Update configuration to use `dogstatsd` provider
3. Start worker: `go run cmd/worker/main.go`
4. Start server: `go run cmd/server/main.go`
5. Metrics will be sent to your StatsD endpoint

## Key Implementation Details

- **Factory Pattern**: `internal/metrics/factory.go` creates appropriate reporters
- **DogStatsD Reporter**: `internal/metrics/dogstatsd.go` implements Tally interface
- **Tag Formatting**: Converts Temporal tags to DataDog format (`key:value`)
- **Metric Naming**: Prefixes all metrics with `temporal.` for namespacing

## Learning Points

- How to implement custom Tally reporters
- DataDog StatsD protocol and tag formatting
- Configuration-driven metrics selection
- Proper error handling in metrics reporters
