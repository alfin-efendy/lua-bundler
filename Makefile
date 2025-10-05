# Makefile for Lua Bundler
# Go-based tool for bundling Lua scripts

# Variables
BINARY_NAME=lua-bundler
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
MAIN_FILE=main.go
BUILD_DIR=build
OUTPUT_DIR=output

# Entry and output configuration
ENTRY_FILE ?= example/myscript/main.lua
OUTPUT_FILE ?= $(OUTPUT_DIR)/example_bundle.lua

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Default target
.DEFAULT_GOAL := build

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build clean test run help install deps fmt vet lint check release example

# Show help
help:
	@echo "$(GREEN)Lua Bundler - Available commands:$(NC)"
	@echo "  $(YELLOW)build$(NC)        - Build the binary for current platform"
	@echo "  $(YELLOW)build-all$(NC)    - Build binaries for all platforms"
	@echo "  $(YELLOW)run$(NC)          - Run the program with example"
	@echo "  $(YELLOW)test$(NC)         - Run tests"
	@echo "  $(YELLOW)clean$(NC)        - Clean build artifacts"
	@echo "  $(YELLOW)install$(NC)      - Install binary to GOPATH/bin"
	@echo "  $(YELLOW)deps$(NC)         - Download dependencies"
	@echo "  $(YELLOW)fmt$(NC)          - Format code"
	@echo "  $(YELLOW)vet$(NC)          - Run go vet"
	@echo "  $(YELLOW)lint$(NC)         - Run golint (if available)"
	@echo "  $(YELLOW)check$(NC)        - Run fmt, vet, and lint"
	@echo "  $(YELLOW)example$(NC)      - Build and run example"
	@echo "  $(YELLOW)release$(NC)      - Create release build"
	@echo ""
	@echo "$(GREEN)Configuration Variables:$(NC)"
	@echo "  $(YELLOW)ENTRY_FILE$(NC)   - Entry Lua file (default: $(ENTRY_FILE))"
	@echo "  $(YELLOW)OUTPUT_FILE$(NC)  - Output bundle file (default: $(OUTPUT_FILE))"
	@echo ""
	@echo "$(GREEN)Usage Examples:$(NC)"
	@echo "  make run ENTRY_FILE=my_script.lua OUTPUT_FILE=my_bundle.lua"
	@echo "  make run-copy ENTRY_FILE=my_script.lua"
	@echo "  make example ENTRY_FILE=another_script.lua"

# Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...

# Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	go vet ./...

# Run golint if available
lint:
	@echo "$(GREEN)Running golint...$(NC)"
	@command -v golint >/dev/null 2>&1 && golint ./... || echo "$(YELLOW)golint not available, skipping...$(NC)"

# Run all checks
check: fmt vet lint
	@echo "$(GREEN)All checks completed!$(NC)"

# Build for current platform
build: deps
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Build for all platforms
build-all: deps
	@echo "$(GREEN)Building for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	@echo "$(YELLOW)Building for Linux...$(NC)"
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) $(MAIN_FILE)
	
	# Windows
	@echo "$(YELLOW)Building for Windows...$(NC)"
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_WINDOWS) $(MAIN_FILE)
	
	# macOS
	@echo "$(YELLOW)Building for macOS...$(NC)"
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin $(MAIN_FILE)
	
	@echo "$(GREEN)All builds completed!$(NC)"
	@ls -la $(BUILD_DIR)/

# Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v ./...

# Clean build artifacts
clean:
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -rf $(OUTPUT_DIR)
	rm -f bundle.lua
	rm -f *_bundle.lua
	go clean

# Install to GOPATH/bin
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME)...$(NC)"
	go install $(LDFLAGS) .
	@echo "$(GREEN)Installation completed!$(NC)"

# Run the program with example
run: build
	@echo "$(GREEN)Running with example script...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	@if [ -f "$(ENTRY_FILE)" ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) -e $(ENTRY_FILE) -o $(OUTPUT_FILE); \
		echo "$(GREEN)Copying output file to clipboard...$(NC)"; \
		if [ -f "$(OUTPUT_FILE)" ]; then \
			if command -v xclip >/dev/null 2>&1; then \
				cat "$(OUTPUT_FILE)" | xclip -selection clipboard; \
				echo "$(GREEN)✓ Content copied to clipboard using xclip!$(NC)"; \
			elif command -v xsel >/dev/null 2>&1; then \
				cat "$(OUTPUT_FILE)" | xsel --clipboard --input; \
				echo "$(GREEN)✓ Content copied to clipboard using xsel!$(NC)"; \
			elif command -v wl-copy >/dev/null 2>&1; then \
				cat "$(OUTPUT_FILE)" | wl-copy; \
				echo "$(GREEN)✓ Content copied to clipboard using wl-copy (Wayland)!$(NC)"; \
			else \
				echo "$(RED)No clipboard tool found! Please install xclip, xsel, or wl-copy$(NC)"; \
				echo "$(YELLOW)Install with: sudo apt-get install xclip$(NC)"; \
			fi; \
		else \
			echo "$(RED)Output file $(OUTPUT_FILE) not found!$(NC)"; \
		fi; \
	else \
		echo "$(RED)Entry file $(ENTRY_FILE) not found!$(NC)"; \
		exit 1; \
	fi

# Run example in release mode
example: build
	@echo "$(GREEN)Running example in release mode...$(NC)"
	@mkdir -p $(OUTPUT_DIR)
	@if [ -f "$(ENTRY_FILE)" ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) -e $(ENTRY_FILE) -o $(OUTPUT_FILE) --release --obfuscate 3; \
		echo "$(GREEN)Example bundle created: $(OUTPUT_FILE)$(NC)"; \
	else \
		echo "$(RED)Entry file $(ENTRY_FILE) not found!$(NC)"; \
		exit 1; \
	fi

# Create release build (optimized)
release: check
	@echo "$(GREEN)Creating release build...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -a -installsuffix cgo $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)Release build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Development workflow - watch for changes and rebuild
watch:
	@echo "$(GREEN)Watching for changes...$(NC)"
	@command -v entr >/dev/null 2>&1 || (echo "$(RED)entr is required for watch mode. Install with: apt-get install entr$(NC)" && exit 1)
	find . -name "*.go" | entr -r make build

# Show project info
info:
	@echo "$(GREEN)Project Information:$(NC)"
	@echo "  Binary name: $(BINARY_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Go version: $(shell go version)"
	@echo "  Build dir: $(BUILD_DIR)"
	@echo "  Output dir: $(OUTPUT_DIR)"
	@echo "  Main file: $(MAIN_FILE)"
	@echo ""
	@echo "$(GREEN)Current Configuration:$(NC)"
	@echo "  Entry file: $(ENTRY_FILE)"
	@echo "  Output file: $(OUTPUT_FILE)"

# Quick development test
dev-test: build
	@echo "$(GREEN)Running development test...$(NC)"
	@if [ -f "$(ENTRY_FILE)" ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) -e $(ENTRY_FILE) -o dev_test_bundle.lua; \
		echo "$(GREEN)Development test completed!$(NC)"; \
		rm -f dev_test_bundle.lua; \
	else \
		echo "$(YELLOW)No entry file $(ENTRY_FILE) found, building only...$(NC)"; \
	fi