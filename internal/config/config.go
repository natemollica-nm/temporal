package config

import (
	"os"
	"time"
)

type Config struct {
	Temporal TemporalConfig
	Server   ServerConfig
	Metrics  MetricsConfig
}

type TemporalConfig struct {
	HostPort  string
	Namespace string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type MetricsConfig struct {
	Provider   string                 `yaml:"provider"`   // "prometheus" or "dogstatsd"
	Prometheus PrometheusConfig       `yaml:"prometheus"`
	DogStatsD  DogStatsDConfig        `yaml:"dogstatsd"`
}

type PrometheusConfig struct {
	ListenAddress string `yaml:"listenAddress"`
}

type DogStatsDConfig struct {
	HostPort      string        `yaml:"hostPort"`
	FlushInterval time.Duration `yaml:"flushInterval"`
	FlushBytes    int           `yaml:"flushBytes"`
}

func Default() Config {
	// Read metrics provider from environment, default to prometheus
	provider := os.Getenv("METRICS_PROVIDER")
	if provider == "" {
		provider = "prometheus"
	}

	return Config{
		Temporal: TemporalConfig{
			HostPort:  "127.0.0.1:7233",
			Namespace: "default",
		},
		Server: ServerConfig{
			Port:         4000,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Metrics: MetricsConfig{
			Provider: provider,
			Prometheus: PrometheusConfig{
				ListenAddress: ":9090",
			},
			DogStatsD: DogStatsDConfig{
				HostPort:      "127.0.0.1:8125",
				FlushInterval: time.Second,
				FlushBytes:    1432,
			},
		},
	}
}
