.PHONY: all build test run clean download-opm generate-index help

# Configuration
INDEX_VERSION := v4.20
OPM_VERSION := 4.18.9
OPM_URL := https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/$(OPM_VERSION)/opm-linux-$(OPM_VERSION).tar.gz
OPM_TARBALL := opm-linux-$(OPM_VERSION).tar.gz
OPM_BINARY := ./bin/opm
INDEX_IMAGE := registry.redhat.io/redhat/redhat-operator-index:$(INDEX_VERSION)
INDEX_FILE := redhat-operator-index.json
BINARY_NAME := index-arch-sorter
BUILD_DIR := ./bin

# Default target
all: build

# Help target
help:
	@echo "Available targets:"
	@echo "  all             - Build the application (default)"
	@echo "  download-opm    - Download and extract the opm CLI tool"
	@echo "  generate-index  - Generate the operator index JSON file using opm"
	@echo "  build           - Build the Go application"
	@echo "  test            - Run Go tests"
	@echo "  run             - Run the application with the generated index"
	@echo "  run-json        - Run the application and output JSON format"
	@echo "  run-csv         - Run the application and output CSV format"
	@echo "  run-arch        - Filter by architecture (Usage: make run-arch ARCH=arm64)"
	@echo "  run-operator    - Filter by operator name (Usage: make run-operator OPERATOR=amq)"
	@echo "  run-filter      - Filter by operator and/or arch (Usage: make run-filter OPERATOR=amq ARCH=arm64)"
	@echo "  clean           - Remove generated files and binaries"
	@echo "  fmt             - Format Go code"
	@echo "  lint            - Run Go linter"
	@echo "  full            - Full workflow: download opm, generate index, build, and run"
	@echo ""
	@echo "Configuration:"
	@echo "  INDEX_VERSION   = $(INDEX_VERSION)"
	@echo "  OPM_VERSION     = $(OPM_VERSION)"
	@echo "  INDEX_IMAGE     = $(INDEX_IMAGE)"
	@echo "  INDEX_FILE      = $(INDEX_FILE)"
	@echo ""
	@echo "Examples:"
	@echo "  make run-arch ARCH=arm64"
	@echo "  make run-operator OPERATOR=cluster"
	@echo "  make run-filter OPERATOR=amq ARCH=amd64"

# Download and extract opm CLI
download-opm: $(OPM_BINARY)

$(OPM_BINARY):
	@echo "Downloading opm CLI version $(OPM_VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@if [ ! -f $(BUILD_DIR)/$(OPM_TARBALL) ]; then \
		curl -L $(OPM_URL) -o $(BUILD_DIR)/$(OPM_TARBALL); \
	fi
	@echo "Extracting opm CLI..."
	@tar -xzf $(BUILD_DIR)/$(OPM_TARBALL) -C $(BUILD_DIR)
	@# Rename extracted file to opm regardless of its original name (e.g., opm-rhel8, opm-linux, etc.)
	@if [ -f $(BUILD_DIR)/opm-rhel8 ]; then \
		mv $(BUILD_DIR)/opm-rhel8 $(OPM_BINARY); \
	elif [ -f $(BUILD_DIR)/opm-linux ]; then \
		mv $(BUILD_DIR)/opm-linux $(OPM_BINARY); \
	fi
	@chmod +x $(OPM_BINARY)
	@echo "opm CLI ready at $(OPM_BINARY)"

# Generate operator index JSON file
generate-index: $(OPM_BINARY)
	@echo "Generating operator index from $(INDEX_IMAGE)..."
	@echo "This may take several minutes..."
	@$(OPM_BINARY) render $(INDEX_IMAGE) -o json > $(INDEX_FILE)
	@echo "Index file generated: $(INDEX_FILE)"
	@wc -l $(INDEX_FILE)

# Build the Go application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sorter
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the application
run: build
	@if [ ! -f $(INDEX_FILE) ]; then \
		echo "Error: $(INDEX_FILE) not found. Run 'make generate-index' first."; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME) -file $(INDEX_FILE)

# Run with JSON output
run-json: build
	@if [ ! -f $(INDEX_FILE) ]; then \
		echo "Error: $(INDEX_FILE) not found. Run 'make generate-index' first."; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME) with JSON output..."
	@$(BUILD_DIR)/$(BINARY_NAME) -file $(INDEX_FILE) -format json

# Run with CSV output
run-csv: build
	@if [ ! -f $(INDEX_FILE) ]; then \
		echo "Error: $(INDEX_FILE) not found. Run 'make generate-index' first."; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME) with CSV output..."
	@$(BUILD_DIR)/$(BINARY_NAME) -file $(INDEX_FILE) -format csv

# Run with architecture filter (example: make run-arch ARCH=arm64)
run-arch: build
	@if [ ! -f $(INDEX_FILE) ]; then \
		echo "Error: $(INDEX_FILE) not found. Run 'make generate-index' first."; \
		exit 1; \
	fi
	@if [ -z "$(ARCH)" ]; then \
		echo "Error: ARCH variable not set. Usage: make run-arch ARCH=amd64|arm64|ppc64le|s390x"; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME) filtering by $(ARCH)..."
	@$(BUILD_DIR)/$(BINARY_NAME) -file $(INDEX_FILE) -arch $(ARCH)

# Run with operator filter (example: make run-operator OPERATOR=amq)
run-operator: build
	@if [ ! -f $(INDEX_FILE) ]; then \
		echo "Error: $(INDEX_FILE) not found. Run 'make generate-index' first."; \
		exit 1; \
	fi
	@if [ -z "$(OPERATOR)" ]; then \
		echo "Error: OPERATOR variable not set. Usage: make run-operator OPERATOR=<name>"; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME) filtering by operator: $(OPERATOR)..."
	@$(BUILD_DIR)/$(BINARY_NAME) -file $(INDEX_FILE) -operator $(OPERATOR)

# Run with both operator and architecture filters
run-filter: build
	@if [ ! -f $(INDEX_FILE) ]; then \
		echo "Error: $(INDEX_FILE) not found. Run 'make generate-index' first."; \
		exit 1; \
	fi
	@if [ -z "$(OPERATOR)" ] && [ -z "$(ARCH)" ]; then \
		echo "Error: Neither OPERATOR nor ARCH variable set. Usage: make run-filter OPERATOR=<name> ARCH=<arch>"; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME) with filters..."
	@CMD="$(BUILD_DIR)/$(BINARY_NAME) -file $(INDEX_FILE)"; \
	if [ -n "$(OPERATOR)" ]; then CMD="$$CMD -operator $(OPERATOR)"; fi; \
	if [ -n "$(ARCH)" ]; then CMD="$$CMD -arch $(ARCH)"; fi; \
	$$CMD

# Clean up generated files
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(INDEX_FILE)
	@echo "Clean complete"

# Full workflow: download opm, generate index, build, and run
full: download-opm generate-index build run

# Format Go code
fmt:
	@echo "Formatting Go code..."
	@go fmt ./...

# Run Go linter
lint:
	@echo "Running Go linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

# Show version information
version:
	@echo "Index Arch Sorter"
	@echo "Index Version: $(INDEX_VERSION)"
	@echo "OPM Version: $(OPM_VERSION)"
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "Binary: $(BUILD_DIR)/$(BINARY_NAME)"; \
	else \
		echo "Binary: Not built"; \
	fi
	@if [ -f $(OPM_BINARY) ]; then \
		echo "OPM: $(OPM_BINARY)"; \
		$(OPM_BINARY) version 2>/dev/null || echo "  (installed)"; \
	else \
		echo "OPM: Not downloaded"; \
	fi

