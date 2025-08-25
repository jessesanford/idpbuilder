# SPLIT-PLAN-001.md
## Split 001 of 2: OCI Foundation Types
**Planner**: Code Reviewer code-reviewer-1756082519 (same for ALL splits)
**Parent Effort**: oci-types

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase1/wave1/oci-types
  - Path: efforts/phase1/wave1/oci-types/split-001/
  - Branch: phase1/wave1/oci-types-split-001
- **Next Split**: Split 002 of phase1/wave1/oci-types
  - Path: efforts/phase1/wave1/oci-types/split-002/
  - Branch: phase1/wave1/oci-types-split-002
- **File Boundaries**:
  - This Split Start: Line 1 / File: pkg/doc.go
  - This Split End: Line 191 / File: pkg/oci/manifest_test.go
  - Next Split Start: Line 1 / File: pkg/stack/constants.go

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- pkg/doc.go (39 lines) - Package documentation
- pkg/oci/constants.go (56 lines) - OCI constants and defaults
- pkg/oci/types.go (121 lines) - Core OCI type definitions
- pkg/oci/types_test.go (130 lines) - OCI types unit tests
- pkg/oci/manifest.go (124 lines) - OCI manifest handling
- pkg/oci/manifest_test.go (191 lines) - Manifest unit tests

**Total Lines**: 661 lines (COMPLIANT)

## Functionality
- Package documentation defining overall purpose
- Core OCI type definitions (OCIReference, OCIImage, OCIManifest, OCIDescriptor)
- OCI constants for media types and registry defaults
- Manifest parsing and validation logic
- Comprehensive unit tests for all OCI types and manifest operations

## Dependencies
- None (foundational split - no dependencies on stack package)
- External: Standard library, testing framework
- This split is self-contained and compilable

## Implementation Instructions
1. Create split working directory:
   ```bash
   mkdir -p efforts/phase1/wave1/oci-types/split-001
   cd efforts/phase1/wave1/oci-types/split-001
   ```

2. Create branch for this split:
   ```bash
   git checkout -b phase1/wave1/oci-types-split-001
   ```

3. Copy ONLY the files assigned to this split:
   ```bash
   mkdir -p pkg/oci
   cp ../pkg/doc.go pkg/
   cp ../pkg/oci/constants.go pkg/oci/
   cp ../pkg/oci/types.go pkg/oci/
   cp ../pkg/oci/types_test.go pkg/oci/
   cp ../pkg/oci/manifest.go pkg/oci/
   cp ../pkg/oci/manifest_test.go pkg/oci/
   ```

4. Implement complete functionality for OCI types
5. Ensure compilation: `go build ./pkg/oci/...`
6. Run unit tests: `go test ./pkg/oci/...`
7. Measure with line counter tool to verify compliance

## Testing Requirements
- All OCI type methods must have tests
- Manifest parsing edge cases must be covered
- Test coverage target: >80%
- Tests must pass independently without stack package

## Split Branch Strategy
- Branch: `phase1/wave1/oci-types-split-001`
- Must merge to: `phase1/wave1/oci-types` after review
- No dependencies on Split 002 (stack types)