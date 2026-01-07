.PHONY: help build clean test worker server deps fmt lint check install-tools

.DEFAULT_GOAL := help

##@ Build & Clean
build: ## Build all binaries
	@mkdir -p bin
	go build -o bin/worker cmd/worker/main.go
	go build -o bin/server cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/ worker server

##@ Development
worker: ## Start Temporal worker
	go run cmd/worker/main.go

server: ## Start web server
	go run cmd/server/main.go

dev-setup: deps install-tools ## Setup development environment
	@echo "Development environment ready!"

##@ Metrics & Observability
metrics-prometheus: ## Run worker with Prometheus metrics (default)
	@echo "Starting with Prometheus metrics on :9090"
	@echo "Visit http://localhost:9090/metrics to see metrics"
	go run cmd/worker/main.go

metrics-dogstatsd: ## Run worker with DogStatsD metrics
	@echo "Starting with DogStatsD metrics to 127.0.0.1:8125"
	@echo "Make sure your DataDog agent or StatsD server is running"
	METRICS_PROVIDER=dogstatsd go run cmd/worker/main.go

##@ Code Quality
test: ## Run tests
	go test -v ./...

fmt: ## Format code
	@gofumpt -l -w .

lint: ## Run linter
	golangci-lint run

check: fmt lint test ## Run all quality checks

##@ Dependencies
deps: ## Download and tidy dependencies
	go mod download
	go mod tidy

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

##@ Temporal Server
temporal-start: ## Start Temporal development server
	temporal server start-dev

temporal-stop: ## Stop Temporal server
	pkill -f "temporal server"

##@ Docker
docker-build: ## Build Docker images
	docker build -t temporal-worker -f deployments/Dockerfile.worker .
	docker build -t temporal-server -f deployments/Dockerfile.server .

##@ Workflows
all: clean deps check build ## Full build pipeline

##@ Help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)