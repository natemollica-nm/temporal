#!/bin/bash

# Development helper script for Temporal learning repo

case "$1" in
    "worker")
        echo "Starting Temporal worker..."
        go run cmd/worker/main.go
        ;;
    "server")
        echo "Starting web server..."
        go run cmd/server/main.go
        ;;
    "build")
        echo "Building applications..."
        go build -o bin/worker cmd/worker/main.go
        go build -o bin/server cmd/server/main.go
        echo "Built binaries in bin/ directory"
        ;;
    "clean")
        echo "Cleaning build artifacts..."
        rm -rf bin/ worker server
        ;;
    *)
        echo "Usage: $0 {worker|server|build|clean}"
        echo "  worker - Start the Temporal worker"
        echo "  server - Start the web server"
        echo "  build  - Build both applications"
        echo "  clean  - Clean build artifacts"
        exit 1
        ;;
esac
