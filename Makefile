.PHONY: help build run test clean docker-up docker-down migrate

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make docker-up   - Start Docker containers"
	@echo "  make docker-down - Stop Docker containers"
	@echo "  make migrate     - Run database migrations"

build:
	@echo "Building application..."
	@go build -o bin/go-chat ./cmd/server

run:
	@echo "Running application..."
	@go run ./cmd/server/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose -f docker/docker-compose.yml up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose -f docker/docker-compose.yml down

docker-logs:
	@docker-compose -f docker/docker-compose.yml logs -f

migrate:
	@echo "Running migrations..."
	@bash scripts/migrate.sh

install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Development
dev:
	@echo "Starting development mode with air (hot reload)..."
	@air

# Docker build and push for GCP
docker-build-gcp:
	@echo "Building Docker image for GCP..."
	@docker build -t gcr.io/YOUR_PROJECT_ID/chat-app:latest -f docker/Dockerfile .

docker-push-gcp:
	@echo "Pushing Docker image to GCP..."
	@docker push gcr.io/YOUR_PROJECT_ID/chat-app:latest