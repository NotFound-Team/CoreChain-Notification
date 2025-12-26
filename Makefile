.PHONY: help build run test clean docker-up docker-down docker-logs deps

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

deps: ## Download Go dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/notification-service ./cmd/server

run: ## Run the application locally
	go run ./cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	docker build -t notification-service:latest -f deployments/docker/Dockerfile .

docker-up: ## Start all services with Docker Compose
	cd deployments/docker && docker-compose up -d

docker-down: ## Stop all Docker services
	cd deployments/docker && docker-compose down

docker-logs: ## View Docker logs
	cd deployments/docker && docker-compose logs -f

docker-restart: docker-down docker-up ## Restart Docker services

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...
