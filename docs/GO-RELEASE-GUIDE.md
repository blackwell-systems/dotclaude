# Go Release Best Practices

This document covers release best practices for Go projects, with specific details on how dotclaude implements them.

## Table of Contents

- [Versioning](#versioning)
- [Release Workflow](#release-workflow)
- [Cross-Platform Builds](#cross-platform-builds)
- [GoReleaser](#goreleaser)
- [CI/CD Pipeline](#cicd-pipeline)
- [pkg.go.dev Integration](#pkggodev-integration)
- [Checksums and Signing](#checksums-and-signing)
- [Distribution Channels](#distribution-channels)

## Versioning

### Semantic Versioning (SemVer)

Go projects should follow [Semantic Versioning](https://semver.org/):

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
```

- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes (backwards compatible)
- **PRERELEASE**: Pre-release identifiers (alpha, beta, rc)

### Go Module Versioning

For v2+, the import path must include the major version:

```go
// v0.x.x and v1.x.x
import "github.com/user/project"

// v2.x.x and higher
import "github.com/user/project/v2"
```

### Git Tags

Always use the `v` prefix for release tags:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Release Workflow

### 1. Pre-Release Checklist

```bash
# Ensure tests pass
go test ./...

# Check for race conditions
go test -race ./...

# Verify module dependencies are tidy
go mod tidy
go mod verify

# Run linters
golangci-lint run

# Check for vulnerabilities
govulncheck ./...
```

### 2. Version Bumping

Update version in code (if applicable):

```go
// internal/version/version.go
const Version = "1.0.0"
```

Or inject at build time via ldflags (preferred):

```bash
go build -ldflags "-X main.version=1.0.0" ./cmd/myapp
```

### 3. Create and Push Tag

```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push tag to trigger release
git push origin v1.0.0
```

## Cross-Platform Builds

### Target Platforms

Common targets for CLI tools:

| GOOS    | GOARCH | Description              |
|---------|--------|--------------------------|
| linux   | amd64  | Linux x86_64             |
| linux   | arm64  | Linux ARM64 (Raspberry Pi 4+, cloud instances) |
| darwin  | amd64  | macOS Intel              |
| darwin  | arm64  | macOS Apple Silicon      |
| windows | amd64  | Windows x86_64           |

### CGO Considerations

For maximum portability, disable CGO:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp-linux-amd64 ./cmd/myapp
```

### Build Flags

Recommended build flags for releases:

```bash
go build -ldflags "-s -w -X main.version=${VERSION}" ./cmd/myapp
```

- `-s`: Omit symbol table
- `-w`: Omit DWARF debugging info
- `-X`: Set version variables

## GoReleaser

[GoReleaser](https://goreleaser.com/) is the industry standard for Go releases.

### Installation

```bash
# macOS
brew install goreleaser

# Linux
go install github.com/goreleaser/goreleaser@latest

# Or use in CI (see workflow below)
```

### Configuration (`.goreleaser.yml`)

```yaml
version: 2

builds:
  - main: ./cmd/myapp
    binary: myapp
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

release:
  github:
    owner: myorg
    name: myrepo
```

### Testing Locally

```bash
# Dry run (no publish)
goreleaser release --snapshot --clean

# Verify config
goreleaser check
```

## CI/CD Pipeline

### GitHub Actions Release Workflow

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Run tests
        run: go test -v ./...

      - uses: goreleaser/goreleaser-action@v6
        with:
          version: '~> v2'
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Branch Protection

Recommended branch protection rules for `main`:

- Require pull request reviews
- Require status checks to pass
- Require conversation resolution
- Do not allow bypassing the above settings

## pkg.go.dev Integration

### Automatic Indexing

pkg.go.dev automatically indexes public Go modules. To trigger indexing:

1. **Push a tag**: pkg.go.dev watches for new versions
2. **Request manually**: Visit `https://pkg.go.dev/github.com/user/project` and click "Request"
3. **Proxy fetch**: `GOPROXY=https://proxy.golang.org go get github.com/user/project@latest`

### Package Documentation

Create a `doc.go` file in your package root:

```go
// Package myproject provides description of what it does.
//
// # Features
//
//   - Feature 1
//   - Feature 2
//
// # Installation
//
//   go install github.com/user/myproject@latest
//
// # Usage
//
// Basic usage example:
//
//   myproject command
//
// See https://github.com/user/myproject for documentation.
package myproject
```

### README Badge

Add the pkg.go.dev badge to your README:

```markdown
[![Go Reference](https://pkg.go.dev/badge/github.com/user/project.svg)](https://pkg.go.dev/github.com/user/project)
```

### Go Report Card

Add code quality badge:

```markdown
[![Go Report Card](https://goreportcard.com/badge/github.com/user/project)](https://goreportcard.com/report/github.com/user/project)
```

## Checksums and Signing

### SHA256 Checksums

GoReleaser automatically generates checksums:

```yaml
checksum:
  name_template: 'checksums.txt'
  algorithm: sha256
```

Users verify downloads:

```bash
sha256sum -c checksums.txt
```

### GPG Signing (Optional)

For additional security, sign releases:

```yaml
signs:
  - artifacts: checksum
    args:
      - "--batch"
      - "--local-user"
      - "{{ .Env.GPG_FINGERPRINT }}"
      - "--output"
      - "${signature}"
      - "--detach-sign"
      - "${artifact}"
```

## Distribution Channels

### 1. GitHub Releases (Primary)

Releases are automatically published via GoReleaser.

### 2. go install

Users can install directly:

```bash
go install github.com/user/project/cmd/myapp@latest
```

### 3. Homebrew (Optional)

Add to `.goreleaser.yml`:

```yaml
brews:
  - name: myapp
    repository:
      owner: user
      name: homebrew-tap
    homepage: https://github.com/user/project
    description: My awesome tool
    install: |
      bin.install "myapp"
```

### 4. Docker (Optional)

Add to `.goreleaser.yml`:

```yaml
dockers:
  - image_templates:
      - "ghcr.io/user/myapp:{{ .Tag }}"
      - "ghcr.io/user/myapp:latest"
    dockerfile: Dockerfile.release
```

## dotclaude Implementation

dotclaude follows these best practices:

| Practice | Implementation |
|----------|----------------|
| Versioning | SemVer with `v` prefix tags |
| Build tool | GoReleaser v2 |
| Platforms | linux/darwin/windows on amd64/arm64 |
| CI/CD | GitHub Actions |
| Checksums | SHA256 via GoReleaser |
| Distribution | GitHub Releases, go install |
| Documentation | pkg.go.dev via doc.go |

### Release Process

1. Update version in code (if needed)
2. Update CHANGELOG.md
3. Create PR and merge to main
4. Create tag: `git tag v1.0.0 && git push origin v1.0.0`
5. GitHub Actions builds and publishes release
6. Verify at pkg.go.dev

### Local Testing

```bash
# Build all platforms
make build-all

# Test GoReleaser config
make release-dry-run
```

## References

- [Semantic Versioning](https://semver.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [pkg.go.dev About](https://pkg.go.dev/about)
- [Go Report Card](https://goreportcard.com/)
- [Go Module Reference](https://go.dev/ref/mod)
- [GitHub Actions for Go](https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go)

---

**Back to:** [README.md](../README.md) | [ARCHITECTURE.md](ARCHITECTURE.md) | [CONTRIBUTING.md](CONTRIBUTING.md)
