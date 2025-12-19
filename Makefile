# OpenRouter CLI Makefile

BINARY_NAME=openrouter-cli
INSTALL_DIR=/usr/local/bin

.PHONY: build clean install fmt

# Default target
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_DIR)/
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
