package metrics

import (
	"fmt"
	"log"
	"time"

	"github.com/natemollica-nm/temporal/internal/config"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	sdktally "go.temporal.io/sdk/contrib/tally"
	sdklog "go.temporal.io/sdk/log"
)

// Factory creates metrics scopes based on configuration
type Factory struct {
	config config.MetricsConfig
	logger sdklog.Logger
}

// NewFactory creates a new metrics factory
func NewFactory(cfg config.MetricsConfig, logger sdklog.Logger) *Factory {
	return &Factory{
		config: cfg,
		logger: logger,
	}
}

// CreateScope creates a tally scope based on the configured provider
func (f *Factory) CreateScope() (tally.Scope, error) {
	switch f.config.Provider {
	case "dogstatsd":
		return f.createDogStatsDScope()
	case "prometheus":
		return f.createPrometheusScope()
	default:
		return f.createPrometheusScope() // Default to Prometheus
	}
}

func (f *Factory) createDogStatsDScope() (tally.Scope, error) {
	reporter, err := NewDogStatsDReporter(f.config.DogStatsD, f.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create DogStatsD reporter: %w", err)
	}

	scopeOpts := tally.ScopeOptions{
		Reporter:  reporter,
		Separator: ".",
		Prefix:    "temporal_samples",
	}

	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	return scope, nil
}

func (f *Factory) createPrometheusScope() (tally.Scope, error) {
	registry := prom.NewRegistry()
	reporter, err := prometheus.Configuration{
		ListenAddress: f.config.Prometheus.ListenAddress,
		TimerType:     "histogram",
	}.NewReporter(prometheus.ConfigurationOptions{
		Registry: registry,
		OnError: func(err error) {
			log.Println("prometheus reporter error:", err)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus reporter: %w", err)
	}

	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
		Prefix:          "temporal_samples",
	}

	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	return sdktally.NewPrometheusNamingScope(scope), nil
}
