# Awesome-Shell Makefile
# Run `make` to see all targets.

.PHONY: build build-linux build-darwin clean test run help

BINARY_NAME := as
BIN_DIR     := bin
MAIN_PKG    := ./cmd/awesome-shell

# Default target
help:
	@echo "Available targets:"
	@echo "  build         - Build binary to $(BIN_DIR)/$(BINARY_NAME) (current OS)"
	@echo "  build-linux   - Build for Linux (GOOS=linux)"
	@echo "  build-darwin  - Build for macOS (GOOS=darwin)"
	@echo "  clean         - Remove built binary"
	@echo "  test          - Run tests"
	@echo "  run           - Build and run (e.g. make run -- password)"
	@echo "  deps          - Download Go modules"

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PKG)
	@echo "Built $(BIN_DIR)/$(BINARY_NAME)"

build-linux:
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PKG)
	@echo "Built $(BIN_DIR)/$(BINARY_NAME)-linux-amd64"

build-darwin:
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PKG)
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PKG)
	@echo "Built $(BIN_DIR)/$(BINARY_NAME)-darwin-*"

clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME) $(BIN_DIR)/$(BINARY_NAME)-* $(BIN_DIR)/awesome-shell $(BIN_DIR)/awesome-shell-*

test:
	go test ./...

run: build
	$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

deps:
	go mod download
	go mod tidy
