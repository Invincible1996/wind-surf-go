.PHONY: run build clean dev test deps install-tools kill-port

# Go binary path
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin

# Port number
PORT=8083

# Default target
all: deps build

# Install project dependencies
deps:
	go mod download
	go mod tidy

# Install development tools
install-tools:
	go install github.com/air-verse/air@latest

# Kill process on port 8083 (more aggressive cleanup)
kill-port:
	@echo "Cleaning up processes..."
	@-pgrep -f "tmp/main" | xargs kill -9 2>/dev/null || true
	@-lsof -ti :$(PORT) | xargs kill -9 2>/dev/null || true
	@-pkill -f "air" 2>/dev/null || true
	@-rm -rf ./tmp
	@echo "Cleaned up processes and temporary files"

# Build the application
build:
	go build -o ./bin/app ./cmd/main.go

# Run the application
run:
	go run ./cmd/main.go

# Run with hot reload using air
dev: kill-port
	@mkdir -p ./tmp
	$(GOBIN)/air -c .air.toml

# Clean build artifacts
clean: kill-port
	rm -rf bin/
	rm -rf tmp/

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Help command
help:
	@echo "Available commands:"
	@echo "  make deps         - Download and tidy dependencies"
	@echo "  make install-tools- Install development tools (air for hot reload)"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make dev         - Run with hot reload using air"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make help        - Show this help message"
