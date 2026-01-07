package metrics

import (
	"github.com/natemollica-nm/temporal/internal/config"
	"github.com/uber-go/tally/v4"
	sdklog "go.temporal.io/sdk/log"
)

var globalScope tally.Scope

// Initialize sets up the global metrics scope based on configuration
func Initialize(cfg config.MetricsConfig, logger sdklog.Logger) error {
	factory := NewFactory(cfg, logger)
	scope, err := factory.CreateScope()
	if err != nil {
		return err
	}
	globalScope = scope
	return nil
}

// GetScope returns the global metrics scope
func GetScope() tally.Scope {
	if globalScope == nil {
		// Fallback to noop scope if not initialized
		return tally.NoopScope
	}
	return globalScope
}
