# Testing Guide

## Prerequisites

### For Controller Tests

Controller tests require kubebuilder test environment binaries. Install them with:

```bash
make envtest
./bin/setup-envtest use 1.29.1
```

This creates a symlink at `bin/k8s/1.29.1-<platform>/` pointing to the kubebuilder binaries.

### For Integration Tests

Integration tests use `IDPBUILDER_TEST_MODE=true` to suppress emoji output that could trigger false test failures. This is automatically set via `TestMain` in test files.

## Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/cmd/push/...
go test ./pkg/controllers/custompackage/...

# Run with verbose output
go test -v ./...

# Skip slow integration tests
go test -short ./...
```

## Bug Fixes

### BUG-QA-002: Test False Positive from Emoji
- **Issue**: ❌ emoji in error messages triggered go test failures
- **Fix**: Added `IDPBUILDER_TEST_MODE` environment variable to suppress emojis during tests
- **Implementation**: TestMain sets this globally, FormatError checks it before adding emojis

### BUG-QA-003: Missing Kubebuilder Dependencies  
- **Issue**: Controller tests failed with "no such file or directory" for k8s binaries
- **Fix**: Document requirement to run `make envtest` before controller tests
- **Note**: These tests validate pre-existing controller code, not OCI push feature
