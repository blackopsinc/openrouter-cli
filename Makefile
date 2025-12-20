# Makefile for cross-compiling openrouter-cli

BINARY_NAME=openrouter-cli
INSTALL_DIR=/usr/local/bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build for all platforms
.PHONY: all
all: linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-386

# Build for Linux AMD64
.PHONY: linux-amd64
linux-amd64: bin
	@echo "Building for Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 .
	@chmod +x bin/$(BINARY_NAME)-linux-amd64

# Build for Linux ARM64
.PHONY: linux-arm64
linux-arm64: bin
	@echo "Building for Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-linux-arm64 .
	@chmod +x bin/$(BINARY_NAME)-linux-arm64

# Build for macOS AMD64 (Intel)
.PHONY: darwin-amd64
darwin-amd64: bin
	@echo "Building for macOS AMD64 (Intel)..."
	@GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 .
	@chmod +x bin/$(BINARY_NAME)-darwin-amd64

# Build for macOS ARM64 (Apple Silicon)
.PHONY: darwin-arm64
darwin-arm64: bin
	@echo "Building for macOS ARM64 (Apple Silicon)..."
	@GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 .
	@chmod +x bin/$(BINARY_NAME)-darwin-arm64

# Build for Windows AMD64 (64-bit)
.PHONY: windows-amd64
windows-amd64: bin
	@echo "Building for Windows AMD64 (64-bit)..."
	@GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe .

# Build for Windows 386 (32-bit)
.PHONY: windows-386
windows-386: bin
	@echo "Building for Windows 386 (32-bit)..."
	@GOOS=windows GOARCH=386 go build -o bin/$(BINARY_NAME)-windows-386.exe .

# Build for current platform
.PHONY: build
build: bin
	@echo "Building for current platform..."
	@go build -o bin/$(BINARY_NAME) .
	@chmod +x bin/$(BINARY_NAME)
	@echo "Build complete: bin/$(BINARY_NAME)"

# Create bin directory if it doesn't exist
bin:
	@mkdir -p bin

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@go clean
	@rm -rf bin/
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

# Install binary to system (requires sudo)
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp bin/$(BINARY_NAME) $(INSTALL_DIR)/
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete"

# Install binary to user's local bin (no sudo required)
.PHONY: install-user
install-user: build
	@echo "Installing $(BINARY_NAME) to $$HOME/.local/bin..."
	@mkdir -p $$HOME/.local/bin
	@cp bin/$(BINARY_NAME) $$HOME/.local/bin/
	@chmod +x $$HOME/.local/bin/$(BINARY_NAME)
	@echo "Installation complete. Add $$HOME/.local/bin to your PATH if not already present."

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              - Build for all platforms (Linux, macOS, Windows)"
	@echo "  linux-amd64      - Build for Linux AMD64"
	@echo "  linux-arm64      - Build for Linux ARM64"
	@echo "  darwin-amd64     - Build for macOS AMD64 (Intel)"
	@echo "  darwin-arm64     - Build for macOS ARM64 (Apple Silicon)"
	@echo "  windows-amd64    - Build for Windows AMD64 (64-bit)"
	@echo "  windows-386      - Build for Windows 386 (32-bit)"
	@echo "  build            - Build for current platform"
	@echo "  install          - Build and install to system (requires sudo)"
	@echo "  install-user     - Build and install to user's local bin"
	@echo "  clean            - Remove build artifacts"
	@echo "  fmt              - Format code"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make all              # Build for all platforms"
	@echo "  make windows-amd64    # Build for Windows 64-bit"
	@echo "  make linux-amd64      # Build for Linux 64-bit"
	@echo "  make darwin-arm64     # Build for macOS Apple Silicon"
	@echo "  make build            # Build for current platform"
	@echo "  make install-user     # Install to ~/.local/bin"
