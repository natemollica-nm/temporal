package metrics

import (
	"log"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	sdktally "go.temporal.io/sdk/contrib/tally"
)

var (
	globalScope tally.Scope
	registry    *prom.Registry
)

func init() {
	registry = prom.NewRegistry()
	reporter, err := prometheus.Configuration{
		ListenAddress: ":9090",
		TimerType:     "histogram",
	}.NewReporter(prometheus.ConfigurationOptions{
		Registry: registry,
		OnError: func(err error) {
			log.Println("prometheus reporter error:", err)
		},
	})
	if err != nil {
		log.Fatal("failed to create prometheus reporter:", err)
	}

	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
		Prefix:          "temporal_samples",
	}

	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	globalScope = sdktally.NewPrometheusNamingScope(scope)
}

func GetScope() tally.Scope {
	return globalScope
}


