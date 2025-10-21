.PHONY: help test test-integration test-unit build lint fmt vet clean install-tools

# Default target
help:
	@echo "Available targets:"
	@echo "  make test              - Run all tests"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-integration  - Run integration tests"
	@echo "  make build             - Build the project"
	@echo "  make lint              - Run golangci-lint"
	@echo "  make fmt               - Format code"
	@echo "  make vet               - Run go vet"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make install-tools     - Install development tools"

# Run all tests
test:
	@echo "Running all tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run unit tests only (short mode)
test-unit:
	@echo "Running unit tests..."
	go test -v -short -race ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@echo "Make sure you have set the following environment variables:"
	@echo "  GHL_CLIENT_ID"
	@echo "  GHL_CLIENT_SECRET"
	@echo "  GHL_ACCESS_TOKEN or (GHL_AUTH_CODE and GHL_REDIRECT_URI)"
	@echo "  GHL_LOCATION_ID"
	@echo ""
	go test -v -race ./...

# Build the project
build:
	@echo "Building..."
	go build -v ./...

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found. Run 'make install-tools' first."; exit 1; }
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f coverage.out coverage.html
	go clean -cache -testcache

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

# Run examples
run-example-simple:
	@echo "Running simple token example..."
	@echo "Make sure you have set: GHL_ACCESS_TOKEN, GHL_LOCATION_ID"
	go run examples/simple_token/main.go

run-example-refresh:
	@echo "Running refresh token example..."
	@echo "Make sure you have set: GHL_CLIENT_ID, GHL_CLIENT_SECRET, GHL_ACCESS_TOKEN, GHL_REFRESH_TOKEN, GHL_LOCATION_ID"
	go run examples/refresh_token/main.go

run-example-basic:
	@echo "Running basic usage example..."
	@echo "Make sure you have set the required environment variables"
	go run examples/basic/main.go
