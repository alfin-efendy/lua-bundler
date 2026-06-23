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

# ez-rbx-ui example smoke test (build every mode, serve one)
SMOKE_ENTRY     ?= testdata/ez-rbx-ui/example/main.lua
SMOKE_LIB_ENTRY ?= testdata/ez-rbx-ui/main.lua
SMOKE_LIB_OUT   ?= testdata/ez-rbx-ui/output/bundle.lua
SMOKE_DIR       ?= $(OUTPUT_DIR)/smoke
SMOKE_MODE      ?= release-o3
SMOKE_PORT      ?= 8080
SMOKE_SERVE     ?= 1

# Build information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.gitCommit=$(GIT_COMMIT)"

# Default target
.DEFAULT_GOAL := build

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build clean test run help install deps fmt vet lint check release example verify-ezui smoke-test

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
	@echo "  $(YELLOW)verify-ezui$(NC)  - Build+run ez-rbx-ui under mocked Roblox (plain/-O2/-O3)"
	@echo "  $(YELLOW)smoke-test$(NC)   - Build ez-rbx-ui example in all modes, then serve one"
	@echo ""
	@echo "$(GREEN)Configuration Variables:$(NC)"
	@echo "  $(YELLOW)ENTRY_FILE$(NC)   - Entry Lua file (default: $(ENTRY_FILE))"
	@echo "  $(YELLOW)OUTPUT_FILE$(NC)  - Output bundle file (default: $(OUTPUT_FILE))"
	@echo "  $(YELLOW)SMOKE_MODE$(NC)   - smoke-test mode to serve (default: $(SMOKE_MODE))"
	@echo "  $(YELLOW)SMOKE_PORT$(NC)   - smoke-test serve port (default: $(SMOKE_PORT))"
	@echo "  $(YELLOW)SMOKE_SERVE$(NC)  - set 0 to build all modes without serving (default: $(SMOKE_SERVE))"
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
	@if git submodule status testdata/ez-rbx-ui 2>/dev/null | grep -q '^-'; then \
		echo "$(GREEN)Initializing git submodule (ez-rbx-ui)...$(NC)"; \
		git submodule update --init --recursive; \
	else \
		echo "$(YELLOW)ez-rbx-ui submodule already checked out — leaving your version as-is.$(NC)"; \
		echo "$(YELLOW)  (run 'git submodule update --init --recursive' to reset it to the pinned version)$(NC)"; \
	fi

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

# Bundle the ez-rbx-ui submodule in release mode and run it under ez-rbx-ui's
# Roblox-faithful harness (loads + mocks Roblox + exercises CreateWindow). This
# is the runtime regression guard for the minifier and obfuscator.
verify-ezui: build
	@echo "$(GREEN)Verifying ez-rbx-ui --release bundle under mocked Roblox...$(NC)"
	@if [ ! -f testdata/ez-rbx-ui/main.lua ]; then \
		echo "$(YELLOW)Submodule missing — run: git submodule update --init --recursive$(NC)"; \
		exit 1; \
	fi
	@command -v lua5.1 >/dev/null 2>&1 || { echo "$(YELLOW)lua5.1 not found on PATH$(NC)"; exit 1; }
	@mkdir -p testdata/ez-rbx-ui/output
	@BIN="$(CURDIR)/$(BUILD_DIR)/$(BINARY_NAME)"; \
	cd testdata/ez-rbx-ui && \
	echo "$(GREEN)[1/3] Verifying plain --release bundle...$(NC)" && \
	"$$BIN" -e main.lua -o output/bundle.lua --release && \
	lua5.1 scripts/verify_bundle.lua && \
	echo "$(GREEN)[2/3] Verifying obfuscated --release -O 2 bundle...$(NC)" && \
	"$$BIN" -e main.lua -o output/bundle.lua --release -O 2 && \
	lua5.1 scripts/verify_bundle.lua && \
	echo "$(GREEN)[3/3] Verifying obfuscated --release -O 3 (string encryption) bundle...$(NC)" && \
	"$$BIN" -e main.lua -o output/bundle.lua --release -O 3 && \
	lua5.1 scripts/verify_bundle.lua
	@echo "$(GREEN)ez-rbx-ui plain, obfuscated (-O 2), and string-encrypted (-O 3) bundles all verified!$(NC)"

# Build the ez-rbx-ui example playground in every mode
# ({normal,release} x {O0,O1,O2,O3}), syntax-check each, then serve one
# (SMOKE_MODE) so it can be loaded in Roblox via loadstring+HttpGet.
#   make smoke-test                                  # build all 8, serve $(SMOKE_MODE) on :$(SMOKE_PORT)
#   make smoke-test SMOKE_MODE=o2 SMOKE_PORT=8081    # serve a different mode/port
#   make smoke-test SMOKE_SERVE=0                    # build + syntax-check all 8, don't serve (CI)
smoke-test: build
	@if [ ! -f "$(SMOKE_ENTRY)" ]; then \
		echo "$(YELLOW)ez-rbx-ui submodule missing — run: git submodule update --init --recursive$(NC)"; \
		exit 1; \
	fi
	@mkdir -p "$(SMOKE_DIR)" "$(dir $(SMOKE_LIB_OUT))"
	@BIN="$(CURDIR)/$(BUILD_DIR)/$(BINARY_NAME)"; \
	echo "$(GREEN)Building ez-rbx-ui library bundle (example dependency)...$(NC)"; \
	"$$BIN" -e "$(SMOKE_LIB_ENTRY)" -o "$(SMOKE_LIB_OUT)" >/dev/null || { echo "$(RED)library bundle failed$(NC)"; exit 1; }; \
	echo "$(GREEN)Building example in all modes -> $(SMOKE_DIR)/$(NC)"; \
	fail=0; serveflags=""; servefound=0; \
	for spec in "normal:" "o1:-O 1" "o2:-O 2" "o3:-O 3" "release:--release" "release-o1:--release -O 1" "release-o2:--release -O 2" "release-o3:--release -O 3"; do \
		name="$${spec%%:*}"; flags="$${spec#*:}"; \
		out="$(SMOKE_DIR)/ezui-$$name.lua"; \
		printf "  %-12s " "$$name"; \
		if "$$BIN" -e "$(SMOKE_ENTRY)" -o "$$out" $$flags >/dev/null 2>&1 && [ -s "$$out" ]; then \
			if command -v luac5.1 >/dev/null 2>&1 && ! luac5.1 -p "$$out" >/dev/null 2>&1; then \
				echo "$(RED)PARSE FAILED$(NC)"; fail=1; \
			else \
				echo "$(GREEN)OK$(NC) ($$(wc -c < "$$out") bytes)"; \
			fi; \
		else \
			echo "$(RED)BUNDLE FAILED$(NC)"; fail=1; \
		fi; \
		if [ "$$name" = "$(SMOKE_MODE)" ]; then serveflags="$$flags"; servefound=1; fi; \
	done; \
	if [ "$$fail" -ne 0 ]; then echo "$(RED)smoke-test: one or more modes failed$(NC)"; exit 1; fi; \
	echo "$(GREEN)All 8 modes built and parse OK.$(NC)"; \
	if [ "$(SMOKE_SERVE)" != "1" ]; then echo "$(YELLOW)SMOKE_SERVE=0 — skipping serve.$(NC)"; exit 0; fi; \
	if [ "$$servefound" -ne 1 ]; then \
		echo "$(RED)Unknown SMOKE_MODE='$(SMOKE_MODE)' (valid: normal o1 o2 o3 release release-o1 release-o2 release-o3)$(NC)"; exit 1; \
	fi; \
	echo "$(GREEN)Serving '$(SMOKE_MODE)' on :$(SMOKE_PORT) — in Roblox: loadstring(game:HttpGet('http://localhost:$(SMOKE_PORT)/ezui-$(SMOKE_MODE).lua'))()$(NC)"; \
	"$$BIN" -e "$(SMOKE_ENTRY)" -o "$(SMOKE_DIR)/ezui-$(SMOKE_MODE).lua" $$serveflags --serve --port $(SMOKE_PORT)

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