.PHONY: build test clean install build-all release-dry-run

# Build the Go binary
build:
	@echo "Building dotclaude-go..."
	@go build -o bin/dotclaude-go cmd/dotclaude/main.go
	@echo "✓ Built: bin/dotclaude-go"

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
	@rm -rf bin/dotclaude-go bin/dotclaude-* dist/
	@echo "✓ Cleaned"

# Install to ~/bin (for local testing)
install: build
	@mkdir -p ~/bin
	@cp bin/dotclaude-go ~/bin/
	@chmod +x ~/bin/dotclaude-go
	@echo "✓ Installed to ~/bin/dotclaude-go"

# Build and run
run: build
	@./bin/dotclaude-go

# Cross-compile for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p bin
	@echo "  linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/dotclaude-linux-amd64 cmd/dotclaude/main.go
	@echo "  linux/arm64..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/dotclaude-linux-arm64 cmd/dotclaude/main.go
	@echo "  darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/dotclaude-darwin-amd64 cmd/dotclaude/main.go
	@echo "  darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/dotclaude-darwin-arm64 cmd/dotclaude/main.go
	@echo "  windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/dotclaude-windows-amd64.exe cmd/dotclaude/main.go
	@echo "✓ Built all platforms in bin/"

# Dry-run release (test goreleaser config)
release-dry-run:
	@goreleaser release --snapshot --clean --skip=publish

# Show help
help:
	@echo "dotclaude Makefile targets:"
	@echo "  make build          - Build the Go binary (current platform)"
	@echo "  make build-all      - Cross-compile for all platforms"
	@echo "  make test           - Run all tests (Go + shell)"
	@echo "  make test-go        - Run Go tests only"
	@echo "  make test-shell     - Run shell tests only"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make install        - Install to ~/bin"
	@echo "  make run            - Build and run"
	@echo "  make release-dry-run - Test GoReleaser configuration"
