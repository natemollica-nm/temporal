package shared

import (
	"time"

	"go.temporal.io/sdk/client"
)

const (
	activityLatency        = "activity_latency"
	scheduleToStartLatency = "schedule_to_start_latency"

	activityStartedCount = "activity_started"
	activityFailedCount  = "activity_failed"
	activitySuccessCount = "activity_succeeded"
)

func RecordActivityStart(handler client.MetricsHandler, activityType string, timeStart int64) client.MetricsHandler {
	handler = handler.WithTags(map[string]string{
		"operation": activityType,
	})
	scheduleToStartDuration := time.Duration(time.Now().UnixNano() - timeStart)
	handler.Timer(scheduleToStartLatency).Record(scheduleToStartDuration)
	handler.Counter(activityStartedCount).Inc(1)
	return handler
}

// RecordActivityEnd emits metrics at the end of an activity function
func RecordActivityEnd(handler client.MetricsHandler, startTime time.Time, err error) {
	handler.Timer(activityLatency).Record(time.Since(startTime))
	if err != nil {
		handler.Counter(activityFailedCount).Inc(1)
		return
	}
	handler.Counter(activitySuccessCount).Inc(1)
}
