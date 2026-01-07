## Temporal Go Learning Repo

### Getting Started

**_Install Temporal GO SDK_**

```shell
go get go.temporal.io/sdk
```

**_Install Temporal CLI_**

```shell
brew install temporal
```

**_Start Temporal Server_**

```shell
temporal server start-dev
```

This command starts a local Temporal Service. It starts the Web UI, creates the default Namespace, and uses an in-memory database.

* The Temporal Service will be available on `localhost:7233`.
* The Temporal Web UI will be available at http://localhost:8233.

Leave the local Temporal Service running as you work through tutorials and other projects. 

The `temporal server start-dev` command uses an **_in-memory_** database, so stopping the server will erase all your Workflows and all your Task Queues. If you want to retain those between runs, start the server and specify a database filename using the `--db-filename` option, like this:

```shell
temporal server start-dev --db-filename your_temporal.db
```

---

# Temporal - Get Address from IP - Go

This application demonstrates using Temporal by calling two APIs in sequence.
It fetches the user's IP address and then uses that address to geolocate that user.

You can use the app in two ways:

- Through a web front-end
- Through a JSON POST request

In both cases, you provide a name that's included in the greeting.

## Using the app

The app requires the Temporal Service to be running.

### Quick Start

1. **Start Temporal Server:**
   ```bash
   make temporal-start
   # or: temporal server start-dev
   ```

2. **Start the Application:**
   ```bash
   # Terminal 1: Start the worker
   make worker
   
   # Terminal 2: Start the web server  
   make server
   ```

3. **Use the Application:**
   - **Web UI:** Visit http://localhost:4000 and enter your name
   - **API:** Send a POST request:
     ```bash
     curl -X POST http://localhost:4000/api \
       -H "Content-Type: application/json" \
       -d '{"name":"Your Name"}'
     ```

4. **Monitor Workflows:**
   - **Temporal UI:** http://localhost:8233
   - **Metrics:** http://localhost:9090/metrics (Prometheus)

### Alternative: Using Make Targets

```bash
# Full development setup
make dev-setup

# Start with different metrics providers
make metrics-prometheus  # Default Prometheus metrics
make metrics-dogstatsd   # DataDog StatsD metrics

# Interactive demo
./scripts/metrics-demo.sh prometheus
```

## Metrics & Observability

This repo supports multiple metrics exporters:

### Prometheus (Default)
```bash
make metrics-prometheus
# Metrics available at http://localhost:9090/metrics
```

### DogStatsD
```bash
make metrics-dogstatsd  
# Sends metrics to DataDog agent at 127.0.0.1:8125
```

### Interactive Demo
```bash
./scripts/metrics-demo.sh prometheus  # or dogstatsd
```

## Project Structure

```
temporal/
├── cmd/                          # Application entry points
│   ├── worker/                   # Temporal worker
│   └── server/                   # Web server
├── internal/                     # Private application code
│   ├── config/                   # Configuration management
│   ├── handlers/                 # HTTP handlers
│   └── metrics/                  # Metrics and observability
├── pkg/temporal/                 # Reusable Temporal components
│   ├── activities/ip/            # IP-related activities
│   ├── workflows/basic/          # Basic workflow patterns
│   └── shared/                   # Common types and utilities
├── examples/                     # Learning examples by pattern
│   └── 01-basic-workflow/        # Current IP geolocation example
└── web/static/                   # Web UI assets
```