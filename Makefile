# OpenRouter CLI Makefile

BINARY_NAME=openrouter-cli
INSTALL_DIR=/usr/local/bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build clean install fmt test install-user

# Default target
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

# Install binary to system (requires sudo)
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_DIR)/
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete"

# Install binary to user's local bin (no sudo required)
install-user: build
	@echo "Installing $(BINARY_NAME) to $$HOME/.local/bin..."
	@mkdir -p $$HOME/.local/bin
	@cp $(BINARY_NAME) $$HOME/.local/bin/
	@chmod +x $$HOME/.local/bin/$(BINARY_NAME)
	@echo "Installation complete. Add $$HOME/.local/bin to your PATH if not already present."

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
