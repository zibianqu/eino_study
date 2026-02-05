.PHONY: help build run test clean docker-up docker-down migrate

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building..."
	@go build -o build/bin/server ./cmd/server
	@go build -o build/bin/cli ./cmd/cli

run: ## Run the application
	@go run cmd/server/main.go

dev: ## Run in development mode with air (hot reload)
	@air

test: ## Run tests
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build files
	@echo "Cleaning..."
	@rm -rf build/bin
	@rm -f coverage.out coverage.html

docker-up: ## Start docker containers
	@docker-compose -f build/docker-compose.yaml up -d

docker-down: ## Stop docker containers
	@docker-compose -f build/docker-compose.yaml down

migrate: ## Run database migrations
	@psql -U postgres -d eino_study -f scripts/init_db.sql

lint: ## Run linter
	@golangci-lint run

tidy: ## Tidy go modules
	@go mod tidy

install-tools: ## Install development tools
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest