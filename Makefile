.PHONY: build run test clean help install dev backend1 backend2

# Variables
BINARY_NAME=lb
GO=go
GOFLAGS=-v

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the load balancer binary
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

run: build ## Build and run the load balancer
	./$(BINARY_NAME)

test: ## Run all tests
	$(GO) test $(GOFLAGS) ./...

test-coverage: ## Run tests with coverage
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

bench: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./...

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	$(GO) clean

install: ## Install dependencies
	$(GO) mod download
	$(GO) mod tidy

fmt: ## Format code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

lint: fmt vet ## Run all linters

dev: ## Run in development mode with auto-reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && $(GO) install github.com/cosmtrek/air@latest)
	air

backend1: ## Start a simple backend server on port 8081
	@echo "Starting backend server on :8081..."
	@$(GO) run examples/backend.go -port 8081

backend2: ## Start a simple backend server on port 8082
	@echo "Starting backend server on :8082..."
	@$(GO) run examples/backend.go -port 8082

all: clean lint test build ## Run all checks and build
