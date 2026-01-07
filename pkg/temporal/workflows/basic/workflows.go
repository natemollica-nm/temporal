package basic

import (
	"fmt"
	"time"

	"github.com/natemollica-nm/temporal/pkg/temporal/activities/ip"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// GetAddressFromIP is the Temporal Workflow that retrieves the IP address and location info.
func GetAddressFromIP(ctx workflow.Context, name string) (string, error) {
	// Define the activity options, including the retry policy
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second, // amount of time that must elapse before the first retry occurs
			MaximumInterval:    time.Minute, // maximum interval between retries
			BackoffCoefficient: 2,           // how much the retry interval increases
			// MaximumAttempts: 5, // Uncomment this if you want to limit attempts
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var ipActivities *ip.IPActivities

	var ip string
	scheduledTimeNanos := workflow.Now(ctx).UnixNano()
	_ = workflow.Sleep(ctx, 500*time.Millisecond)
	err := workflow.ExecuteActivity(ctx, ipActivities.GetIP, scheduledTimeNanos).Get(ctx, &ip)
	if err != nil {
		return "", fmt.Errorf("failed to get IP: %s", err)
	}

	var location string
	err = workflow.ExecuteActivity(ctx, ipActivities.GetLocationInfo, ip, scheduledTimeNanos).Get(ctx, &location)
	if err != nil {
		return "", fmt.Errorf("failed to get location: %s", err)
	}
	var isp string
	err = workflow.ExecuteActivity(ctx, ipActivities.GetInternetServiceProvider, ip, scheduledTimeNanos).Get(ctx, &isp)
	return fmt.Sprintf("Hello, %s. Your IP is %s (%s) and your location is %s", name, ip, isp, location), nil
}
