# PROJECT INTEGRATION FIX PLAN

## Bug Summary
- **Total Bugs Found**: 1
- **Critical**: 0
- **High**: 1 (blocks compilation)
- **Medium**: 0
- **Low**: 0

## Bug Details

### Bug #1: Incorrect Import Paths in Registry Package
- **Severity**: HIGH
- **File**: pkg/registry/gitea.go
- **Lines**: 14-16
- **Impact**: Registry package cannot compile, blocking registry functionality
- **Source Branch**: idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
- **Root Cause**: Import statements use wrong repository path (github.com/jessesanford/idpbuilder instead of github.com/cnoe-io/idpbuilder)

## Fix Strategy

### Fix Group 1: Import Path Corrections
This is a single, independent fix that can be completed immediately.

#### Bug #1: Incorrect Import Paths in Registry Package
- **Source Branch**: idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
- **Fix Location**: pkg/registry/gitea.go:14-16
- **Fix Instructions**:
  ```go
  // Current (INCORRECT) imports at lines 14-16:
  import (
      "github.com/jessesanford/idpbuilder/pkg/certs"
      "github.com/jessesanford/idpbuilder/pkg/certvalidation"
      "github.com/jessesanford/idpbuilder/pkg/fallback"
  )
  
  // Replace with (CORRECT):
  import (
      "github.com/cnoe-io/idpbuilder/pkg/certs"
      "github.com/cnoe-io/idpbuilder/pkg/certvalidation"
      "github.com/cnoe-io/idpbuilder/pkg/fallback"
  )
  ```
- **Testing Required**:
  1. After fix, run `go build ./pkg/registry/...` to verify compilation
  2. Run `go test ./pkg/registry/...` to ensure tests pass
  3. Build entire project with `go build ./...`
  4. Run integration tests if available
- **Assigned To**: sw-engineer-1

## Spawn Instructions

### Single Engineer Spawn:
Since there is only one bug to fix, we will spawn a single SW Engineer:

- **SW Engineer 1**:
  - Navigate to effort directory for gitea-client-split-002
  - Switch to branch: idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
  - Fix import paths in pkg/registry/gitea.go lines 14-16
  - Run tests to verify fix
  - Commit with message: "fix: correct import paths in registry package (R266 upstream bug fix)"
  - Push changes

### Work Directory:
The engineer should work in: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-002`

## Verification Steps

After the fix is applied:
1. Verify the registry package compiles: `go build ./pkg/registry/...`
2. Run registry tests: `go test ./pkg/registry/...`
3. Build full project: `go build ./...`
4. Run full test suite: `go test ./...`
5. Document fix completion in work-log.md

## Expected Outcome

After applying this fix:
- The registry package should compile successfully
- All tests should pass
- The project-integration branch should be ready for final integration testing
- The MVP functionality should be fully operational

## Timeline

- **Estimated Fix Time**: 15-30 minutes
- **Testing Time**: 15-30 minutes
- **Total Duration**: 30-60 minutes

---
**Plan Created**: 2025-09-09
**Orchestrator**: PROJECT_FIX_PLANNING state
**Next State**: SPAWN_SW_ENGINEER_PROJECT_FIXES