.PHONY: build test clean install build-all release-dry-run

# Binary name
BINARY = dotclaude

# Build the Go binary
build:
	@echo "Building $(BINARY)..."
	@mkdir -p bin
	@go build -o bin/$(BINARY) cmd/dotclaude/main.go
	@echo "✓ Built: bin/$(BINARY)"

# Run tests
test:
	@echo "Running Go tests..."
	@go test -v ./...
	@echo "Running shell tests..."
	@bats tests/*.bats

# Run Go tests only
test-go:
	@go test -v ./...

# Run shell tests only
test-shell:
	@bats tests/*.bats

# Clean build artifacts
clean:
	@rm -rf bin/$(BINARY) bin/dotclaude-* dist/
	@echo "✓ Cleaned"

# Install to ~/.local/bin
install: build
	@mkdir -p ~/.local/bin
	@cp bin/$(BINARY) ~/.local/bin/$(BINARY)
	@chmod +x ~/.local/bin/$(BINARY)
	@echo "✓ Installed to ~/.local/bin/$(BINARY)"

# Build and run
run: build
	@./bin/$(BINARY)

# Cross-compile for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p bin
	@echo "  linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(BINARY)-linux-amd64 cmd/dotclaude/main.go
	@echo "  linux/arm64..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/$(BINARY)-linux-arm64 cmd/dotclaude/main.go
	@echo "  darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(BINARY)-darwin-amd64 cmd/dotclaude/main.go
	@echo "  darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/$(BINARY)-darwin-arm64 cmd/dotclaude/main.go
	@echo "  windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(BINARY)-windows-amd64.exe cmd/dotclaude/main.go
	@echo "✓ Built all platforms in bin/"

# Dry-run release (test goreleaser config)
release-dry-run:
	@goreleaser release --snapshot --clean --skip=publish

# Show help
help:
	@echo "dotclaude Makefile targets:"
	@echo "  make build          - Build the Go binary (bin/$(BINARY))"
	@echo "  make build-all      - Cross-compile for all platforms"
	@echo "  make test           - Run all tests (Go + shell)"
	@echo "  make test-go        - Run Go tests only"
	@echo "  make test-shell     - Run shell tests only"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make install        - Install to ~/.local/bin"
	@echo "  make run            - Build and run"
	@echo "  make release-dry-run - Test GoReleaser configuration"
