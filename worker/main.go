package main

import (
	"log"
	"net/http"

	"github.com/natemollica-nm/temporal/iplocate"
	"github.com/natemollica-nm/temporal/metrics"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	sdktally "go.temporal.io/sdk/contrib/tally"
)

func main() {
	// Create the Temporal client
	c, err := client.Dial(client.Options{
		HostPort:       "127.0.0.1:7233",
		Namespace:      "default",
		MetricsHandler: sdktally.NewMetricsHandler(metrics.GetScope()),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create the Temporal worker
	w := worker.New(c, iplocate.TaskQueueName, worker.Options{})

	// inject HTTP client into the Activities Struct
	activities := &iplocate.IPActivities{
		HTTPClient: http.DefaultClient,
	}

	// Register Workflow and Activities
	w.RegisterWorkflow(iplocate.GetAddressFromIP)
	w.RegisterActivity(activities)

	// Start the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
