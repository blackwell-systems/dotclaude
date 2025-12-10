.PHONY: build test clean install

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
	@rm -f bin/dotclaude-go
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

# Show help
help:
	@echo "dotclaude Makefile targets:"
	@echo "  make build    - Build the Go binary"
	@echo "  make test     - Run all tests (Go + shell)"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make install  - Install to ~/bin"
	@echo "  make run      - Build and run"
