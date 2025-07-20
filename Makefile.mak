# Makefile for BBS Sabacc Game

# Variables
BINARY_NAME=sabacc
BINARY_UNIX=$(BINARY_NAME)_unix
VERSION=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d_%T)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"

# Default target
.PHONY: all
all: clean test build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p bin
	@mkdir -p ansi
	@mkdir -p stats
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME) .
	@echo "Build complete: bin/$(BINARY_NAME)"

# Build for different architectures
.PHONY: build-all
build-all: build-linux build-windows build-darwin

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)_linux_amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build $(LDFLAGS) -o bin/$(BINARY_NAME)_linux_386 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)_linux_arm64 .

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)_windows_amd64.exe .
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build $(LDFLAGS) -o bin/$(BINARY_NAME)_windows_386.exe .

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)_darwin_amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)_darwin_arm64 .

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Test coverage report generated: coverage.html"

# Run tests with verbose output
.PHONY: test-verbose
test-verbose:
	go test -v -race -count=1 ./...

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Security check
.PHONY: security
security:
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed, skipping..."; \
		echo "Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Run the application (for testing)
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME) with test data..."
	@if [ ! -f door32.sys ]; then \
		echo "Creating test door32.sys file..."; \
		echo -e "2\n8\n38400\nBBS Sabacc Test\n1\nTest Player\nTestUser\n100\n90\n0\n1" > door32.sys; \
	fi
	./bin/$(BINARY_NAME) -path ./

# Create release package
.PHONY: package
package: build-all
	@echo "Creating release package..."
	@mkdir -p release/$(BINARY_NAME)-$(VERSION)
	@cp bin/* release/$(BINARY_NAME)-$(VERSION)/ 2>/dev/null || true
	@cp README.md release/$(BINARY_NAME)-$(VERSION)/
	@cp LICENSE release/$(BINARY_NAME)-$(VERSION)/ 2>/dev/null || echo "No LICENSE file found"
	@mkdir -p release/$(BINARY_NAME)-$(VERSION)/ansi
	@echo "# ANSI Art Directory" > release/$(BINARY_NAME)-$(VERSION)/ansi/README.md
	@echo "Place your ANSI art files here:" >> release/$(BINARY_NAME)-$(VERSION)/ansi/README.md
	@echo "- title.ans  (Title screen)" >> release/$(BINARY_NAME)-$(VERSION)/ansi/README.md
	@echo "- menu.ans   (Menu background)" >> release/$(BINARY_NAME)-$(VERSION)/ansi/README.md
	@echo "- game.ans   (Game screen)" >> release/$(BINARY_NAME)-$(VERSION)/ansi/README.md
	@cd release && tar -czf $(BINARY_NAME)-$(VERSION).tar.gz $(BINARY_NAME)-$(VERSION)
	@echo "Release package created: release/$(BINARY_NAME)-$(VERSION).tar.gz"

# Install to system (requires sudo)
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete"

# Uninstall from system (requires sudo)
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete"

# Development workflow
.PHONY: dev
dev: clean deps fmt lint test build
	@echo "Development build complete"

# Quick build for development
.PHONY: quick
quick:
	@echo "Quick build..."
	go build -o $(BINARY_NAME) .

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Run clean, test, and build"
	@echo "  build        - Build the application for Linux"
	@echo "  build-all    - Build for all supported platforms"
	@echo "  build-linux  - Build for Linux (multiple architectures)"
	@echo "  build-windows- Build for Windows"
	@echo "  build-darwin - Build for macOS"
	@echo "  test         - Run tests with coverage"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  bench        - Run benchmarks"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  security     - Run security checks"
	@echo "  run          - Build and run with test data"
	@echo "  package      - Create release package"
	@echo "  install      - Install to system (requires sudo)"
	@echo "  uninstall    - Uninstall from system (requires sudo)"
	@echo "  dev          - Complete development workflow"
	@echo "  quick        - Quick development build"
	@echo "  help         - Show this help"

# Create test drop file
.PHONY: create-test-data
create-test-data:
	@echo "Creating test drop file..."
	@echo -e "2\n8\n38400\nBBS Sabacc Test\n1\nTest Player\nTestUser\n100\n90\n0\n1" > door32.sys
	@echo "Test door32.sys created"

# Show project info
.PHONY: info
info:
	@echo "Project: BBS Sabacc"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(shell go version)"
	@echo "Git Commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Working Directory: $(shell pwd)"

# Validate configuration
.PHONY: validate
validate:
	@echo "Validating project structure..."
	@if [ ! -f go.mod ]; then echo "ERROR: go.mod not found"; exit 1; fi
	@if [ ! -f main.go ]; then echo "ERROR: main.go not found"; exit 1; fi
	@if [ ! -f cards.go ]; then echo "ERROR: cards.go not found"; exit 1; fi
	@if [ ! -f helpers.go ]; then echo "ERROR: helpers.go not found"; exit 1; fi
	@echo "Project structure valid"

# Default target info
.DEFAULT_GOAL := help