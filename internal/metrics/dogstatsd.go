package metrics

import (
	"fmt"
	"sort"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/natemollica-nm/temporal/internal/config"
	"github.com/uber-go/tally/v4"
	"go.temporal.io/sdk/log"
)

const (
	defaultFlushBytes    = 1432
	defaultFlushInterval = time.Second
)

type dogstatsdReporter struct {
	client *statsd.Client
	logger log.Logger
}

// NewDogStatsDReporter creates a new DogStatsD metrics reporter
func NewDogStatsDReporter(config config.DogStatsDConfig, logger log.Logger) (tally.StatsReporter, error) {
	hostPort := config.HostPort
	if hostPort == "" {
		hostPort = "127.0.0.1:8125"
	}

	flushInterval := config.FlushInterval
	if flushInterval == 0 {
		flushInterval = defaultFlushInterval
	}

	flushBytes := config.FlushBytes
	if flushBytes == 0 {
		flushBytes = defaultFlushBytes
	}

	client, err := statsd.New(
		hostPort,
		statsd.WithBufferFlushInterval(flushInterval),
		statsd.WithMaxBytesPerPayload(flushBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create DogStatsD client: %w", err)
	}

	return &dogstatsdReporter{
		client: client,
		logger: logger,
	}, nil
}

func (r *dogstatsdReporter) Capabilities() tally.Capabilities {
	return r
}

func (r *dogstatsdReporter) Reporting() bool {
	return true
}

func (r *dogstatsdReporter) Tagging() bool {
	return true
}

func (r *dogstatsdReporter) Flush() {
	if err := r.client.Flush(); err != nil {
		r.logger.Error("Failed to flush DogStatsD metrics", "error", err)
	}
}

func (r *dogstatsdReporter) ReportCounter(name string, tags map[string]string, value int64) {
	name = r.sanitizeMetricName(name)
	if err := r.client.Count(name, value, r.marshalTags(tags), 1); err != nil {
		r.logger.Error("Failed to report counter", "metric", name, "error", err)
	}
}

func (r *dogstatsdReporter) ReportGauge(name string, tags map[string]string, value float64) {
	name = r.sanitizeMetricName(name)
	if err := r.client.Gauge(name, value, r.marshalTags(tags), 1); err != nil {
		r.logger.Error("Failed to report gauge", "metric", name, "error", err)
	}
}

func (r *dogstatsdReporter) ReportTimer(name string, tags map[string]string, interval time.Duration) {
	name = r.sanitizeMetricName(name)
	if err := r.client.Timing(name, interval, r.marshalTags(tags), 1); err != nil {
		r.logger.Error("Failed to report timer", "metric", name, "error", err)
	}
}

func (r *dogstatsdReporter) ReportHistogramValueSamples(name string, tags map[string]string, buckets tally.Buckets, bucketLowerBound, bucketUpperBound float64, samples int64) {
	r.logger.Warn("Histogram value samples not supported in DogStatsD reporter", "metric", name)
}

func (r *dogstatsdReporter) ReportHistogramDurationSamples(name string, tags map[string]string, buckets tally.Buckets, bucketLowerBound, bucketUpperBound time.Duration, samples int64) {
	r.logger.Warn("Histogram duration samples not supported in DogStatsD reporter", "metric", name)
}

func (r *dogstatsdReporter) marshalTags(tags map[string]string) []string {
	if len(tags) == 0 {
		return nil
	}

	var keys []string
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dogTags []string
	for _, key := range keys {
		dogTags = append(dogTags, fmt.Sprintf("%s:%s", key, tags[key]))
	}
	return dogTags
}

func (r *dogstatsdReporter) sanitizeMetricName(name string) string {
	return fmt.Sprintf("temporal.%s", name)
}
