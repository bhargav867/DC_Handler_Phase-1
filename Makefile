.PHONY: help build run test clean setup dev prod docker-build docker-run

help:
	@echo "DC_Handler Phase 1 - Makefile Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup        - Setup project (download deps)"
	@echo "  make install      - Install dependencies"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Run in development mode"
	@echo "  make run          - Run the application"
	@echo "  make build        - Build binary"
	@echo "  make test         - Run tests"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make docker-dev   - Docker Compose development"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"

# Setup
setup: install
	@echo "✓ Setup complete"

install:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies installed"

# Development
env:
	@if [ ! -f .env ]; then cp .env.example .env && echo "✓ Created .env"; else echo "✓ .env exists"; fi

dev: env
	@echo "Starting development server..."
	go run main.go

run: build
	@echo "Running application..."
	./dc-handler

build:
	@echo "Building binary..."
	go build -o dc-handler .
	@echo "✓ Build complete: ./dc-handler"

test:
	@echo "Running tests..."
	go test -v ./...

# Code quality
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Formatting complete"

lint:
	@echo "Linting code..."
	golangci-lint run ./... || echo "Install: https://golangci-lint.run/usage/install/"

# Docker
docker-build:
	@echo "Building Docker image..."
	docker build -t dc-handler:latest .
	@echo "✓ Docker image built"

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env dc-handler:latest

docker-dev:
	@echo "Starting Docker Compose development..."
	docker-compose up

docker-dev-build:
	@echo "Building and starting Docker Compose..."
	docker-compose up --build

# Maintenance
clean:
	@echo "Cleaning build artifacts..."
	rm -f dc-handler
	go clean
	@echo "✓ Clean complete"

vendor:
	@echo "Creating vendor directory..."
	go mod vendor
	@echo "✓ Vendor directory created"

deps:
	@echo "Listing dependencies..."
	go list -m all

.DEFAULT_GOAL := help
