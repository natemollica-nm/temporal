# Example 01: Basic Workflow

This example demonstrates the fundamental Temporal concepts:

## What it does
- Fetches your public IP address
- Gets geolocation information for that IP
- Retrieves ISP information
- Returns a formatted greeting with all the data

## Temporal concepts demonstrated
- **Sequential Activities**: Activities executed one after another
- **Retry Policies**: Automatic retries with exponential backoff
- **Activity Options**: Timeouts and retry configuration
- **Error Handling**: Proper error propagation in workflows

## Files involved
- `pkg/temporal/workflows/basic/workflows.go` - The main workflow
- `pkg/temporal/activities/ip/activities.go` - IP-related activities
- `pkg/temporal/shared/` - Shared constants and metrics

## Running the example
1. Start Temporal server: `temporal server start-dev`
2. Start worker: `go run cmd/worker/main.go`
3. Start web server: `go run cmd/server/main.go`
4. Visit http://localhost:4000 or use the API

## Key learning points
- Activities are where external API calls happen
- Workflows orchestrate the sequence and handle failures
- Retry policies make your workflows resilient to transient failures
