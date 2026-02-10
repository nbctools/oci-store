# Makefile for OCI-Store
# Provides convenient commands for building, testing, and releasing

# Variables
BINARY_NAME=oci-store
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=build
DIST_DIR=dist
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"

# Go settings
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
CGO_ENABLED=0

# Build targets
.PHONY: help build test clean lint install release demo

# Default target
help: ## Show this help message
	@echo "OCI-Store Build System"
	@echo "====================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development targets
build: ## Build the binary for current platform
	@echo "Building $(BINARY_NAME) v$(VERSION) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build binaries for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)

	@echo "Building linux/amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	@tar czf $(DIST_DIR)/$(BINARY_NAME)-linux-amd64.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-linux-amd64

	@echo "Building linux/arm64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	@tar czf $(DIST_DIR)/$(BINARY_NAME)-linux-arm64.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-linux-arm64

	@echo "Building darwin/amd64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@tar czf $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-darwin-amd64

	@echo "Building darwin/arm64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@tar czf $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-darwin-arm64

	@echo "Building windows/amd64..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@cd $(DIST_DIR) && zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "All builds completed in $(DIST_DIR)/"

test: ## Run all tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-quick: ## Run tests without coverage
	@echo "Running quick tests..."
	go test -v ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	rm -f coverage.out coverage.html
	go clean -cache

install: build ## Install binary to system
	@echo "Installing $(BINARY_NAME)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

uninstall: ## Remove binary from system
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Demo and development targets
demo: build ## Run the multi-backend demo
	@echo "Starting S3 demo with localstack..."
	@if [ -d "demo" ]; then \
		cd demo && docker-compose up -d && \
		echo "Waiting for services..." && sleep 10 && \
		cd .. && \
		./$(BUILD_DIR)/$(BINARY_NAME) --help && \
		echo "Demo ready! Run './demo/s3-demo.sh' to test all backends"; \
	else \
		echo "Demo directory not found"; \
	fi

demo-stop: ## Stop demo services
	@echo "Stopping demo services..."
	@if [ -d "demo" ]; then \
		cd demo && docker-compose down; \
	fi

test-cli: build ## Test CLI structure
	@echo "Testing CLI structure..."
	./$(BUILD_DIR)/$(BINARY_NAME) --help
	./$(BUILD_DIR)/$(BINARY_NAME) s3 --help
	./$(BUILD_DIR)/$(BINARY_NAME) gcs --help
	./$(BUILD_DIR)/$(BINARY_NAME) azure --help

# Release targets
release-check: ## Perform pre-release checks
	@echo "Performing release checks..."
	@echo "Version: $(VERSION)"
	@git status --porcelain
	@echo "Running tests..."
	$(MAKE) test
	@echo "Running linter..."
	$(MAKE) lint
	@echo "Building all platforms..."
	$(MAKE) build-all
	@echo "Release checks completed!"

checksums: build-all ## Generate checksums for release
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && \
	sha256sum *.tar.gz *.zip 2>/dev/null | tee checksums.txt
	@echo "Checksums generated in $(DIST_DIR)/checksums.txt"

# Homebrew targets
brew-formula: ## Generate Homebrew formula
	@echo "Generating Homebrew formula..."
	@mkdir -p homebrew
	@sed 's/sha256_placeholder/$(shell cat $(DIST_DIR)/checksums.txt | grep darwin-amd64 | cut -d' ' -f1)/g' \
		oci-store.rb > homebrew/oci-store.rb
	@echo "Formula generated: homebrew/oci-store.rb"

# Development helpers
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Quick development cycle
dev: fmt vet test-quick ## Quick development cycle (format, vet, test)

# Full quality check
qc: fmt vet lint test ## Full quality check

# Show version
version: ## Show current version
	@echo "Version: $(VERSION)"	

