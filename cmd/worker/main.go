package main

import (
	"log"
	"net/http"

	"github.com/natemollica-nm/temporal/internal/config"
	"github.com/natemollica-nm/temporal/internal/metrics"
	"github.com/natemollica-nm/temporal/pkg/temporal/activities/ip"
	"github.com/natemollica-nm/temporal/pkg/temporal/shared"
	"github.com/natemollica-nm/temporal/pkg/temporal/workflows/basic"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	sdktally "go.temporal.io/sdk/contrib/tally"
	sdklog "go.temporal.io/sdk/log"
)

func main() {
	// Load configuration
	cfg := config.Default()
	
	// Create logger - use nil for default logger
	var logger sdklog.Logger
	
	// Initialize metrics
	if err := metrics.Initialize(cfg.Metrics, logger); err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}

	// Create the Temporal client
	c, err := client.Dial(client.Options{
		HostPort:       cfg.Temporal.HostPort,
		Namespace:      cfg.Temporal.Namespace,
		MetricsHandler: sdktally.NewMetricsHandler(metrics.GetScope()),
		Logger:         logger,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create the Temporal worker
	w := worker.New(c, shared.TaskQueueName, worker.Options{})

	// inject HTTP client into the Activities Struct
	activities := &ip.IPActivities{
		HTTPClient: http.DefaultClient,
	}

	// Register Workflow and Activities
	w.RegisterWorkflow(basic.GetAddressFromIP)
	w.RegisterActivity(activities)

	// Start the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
