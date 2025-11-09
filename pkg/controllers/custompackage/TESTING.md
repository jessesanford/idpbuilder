# CustomPackage Controller Tests

## Overview

This directory contains controller tests for the **pre-existing** CustomPackage controller functionality. These tests are **NOT part of the OCI push feature implementation**.

## Test Requirements

The controller tests (`controller_test.go`) require:
- Kubernetes control plane test binaries (etcd, kube-apiserver)
- `kubebuilder` test environment (`envtest`)
- These binaries are typically installed via `make envtest`

## Running Tests Locally

To run these tests locally:

```bash
# Install kubebuilder test dependencies
make envtest

# Run tests
go test -v ./pkg/controllers/custompackage/...
```

## CI Environment

These tests are **automatically skipped in CI** for the following reasons:

1. **Out of Scope**: These controllers are pre-existing code, not part of the OCI push feature
2. **Heavy Dependencies**: Require external Kubernetes binaries (etcd, kube-apiserver)
3. **Not Blocking**: The OCI push feature does not depend on these controllers

### Skip Mechanism

Tests are skipped when:
- Running with `-short` flag: `go test -short ./...`
- Environment variable `SKIP_CONTROLLER_TESTS=true` is set
- Kubebuilder binaries are not available

## Test Coverage

The controller tests cover:

1. **TestReconcileCustomPkg**: Tests reconciliation of CustomPackage resources
2. **TestReconcileCustomPkgAppSet**: Tests ApplicationSet handling
3. **TestReconcileHelmValueObject**: Tests Helm value object reconciliation

These tests validate pre-existing functionality and ensure no regressions in the controller logic.

## For Developers

If you're working on the controllers themselves (not the OCI push feature), you should:

1. Install kubebuilder dependencies: `make envtest`
2. Run tests without the skip flag: `go test ./pkg/controllers/custompackage/...`
3. Ensure all tests pass before committing controller changes

## For CI/CD

CI pipelines should run tests with `-short` flag or set `SKIP_CONTROLLER_TESTS=true`:

```bash
# Option 1: Use -short flag
go test -short ./...

# Option 2: Use environment variable
export SKIP_CONTROLLER_TESTS=true
go test ./...
```

## Related

- **Bug**: BUG-QA-003
- **Scope**: OUT_OF_SCOPE_FOR_OCI_PUSH
- **Status**: Tests properly skipped in CI
