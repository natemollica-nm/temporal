#!/bin/bash

# Metrics demo script for Temporal learning repo

case "$1" in
    "prometheus")
        echo "🔧 Starting Temporal worker with Prometheus metrics..."
        echo "📊 Metrics will be available at: http://localhost:9090/metrics"
        echo "🌐 Web UI will be available at: http://localhost:4000"
        echo ""
        echo "In another terminal, run: ./scripts/metrics-demo.sh server"
        echo ""
        go run cmd/worker/main.go
        ;;
    "dogstatsd")
        echo "🔧 Starting Temporal worker with DogStatsD metrics..."
        echo "📊 Metrics will be sent to: 127.0.0.1:8125"
        echo "🌐 Web UI will be available at: http://localhost:4000"
        echo ""
        echo "Make sure DataDog agent or StatsD server is running!"
        echo "In another terminal, run: ./scripts/metrics-demo.sh server"
        echo ""
        # Create temporary config with DogStatsD
        cat > /tmp/temporal-config.yaml << EOF
metrics:
  provider: "dogstatsd"
  dogstatsd:
    hostPort: "127.0.0.1:8125"
    flushInterval: 1s
    flushBytes: 1432
EOF
        CONFIG_FILE=/tmp/temporal-config.yaml go run cmd/worker/main.go
        ;;
    "server")
        echo "🌐 Starting web server..."
        echo "📱 Visit: http://localhost:4000"
        go run cmd/server/main.go
        ;;
    *)
        echo "Temporal Metrics Demo"
        echo ""
        echo "Usage: $0 {prometheus|dogstatsd|server}"
        echo ""
        echo "  prometheus - Start worker with Prometheus metrics"
        echo "  dogstatsd  - Start worker with DogStatsD metrics"  
        echo "  server     - Start web server"
        echo ""
        echo "Example workflow:"
        echo "  1. Terminal 1: $0 prometheus"
        echo "  2. Terminal 2: $0 server"
        echo "  3. Visit http://localhost:4000 and submit a request"
        echo "  4. Check metrics at http://localhost:9090/metrics"
        exit 1
        ;;
esac
