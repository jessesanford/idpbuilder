# SPLIT-PLAN-002.md
## Split 002 of 2: Stack Configuration Types
**Planner**: Code Reviewer code-reviewer-1756082519 (same for ALL splits)
**Parent Effort**: oci-types

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 001 of phase1/wave1/oci-types
  - Path: efforts/phase1/wave1/oci-types/split-001/
  - Branch: phase1/wave1/oci-types-split-001
  - Summary: Implemented OCI foundation types, constants, and manifest handling
- **This Split**: Split 002 of phase1/wave1/oci-types
  - Path: efforts/phase1/wave1/oci-types/split-002/
  - Branch: phase1/wave1/oci-types-split-002
- **Next Split**: None (final split)
  - Path: N/A
  - Branch: N/A
- **File Boundaries**:
  - Previous Split End: Line 191 / File: pkg/oci/manifest_test.go
  - This Split Start: Line 1 / File: pkg/stack/constants.go
  - This Split End: Line 164 / File: pkg/stack/types_test.go (final file)

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- pkg/stack/constants.go (42 lines) - Stack-specific constants
- pkg/stack/types.go (107 lines) - Stack configuration types
- pkg/stack/types_test.go (164 lines) - Stack types unit tests

**Total Lines**: 313 lines (COMPLIANT)

## Functionality
- Stack configuration type definitions (StackConfiguration, StackComponent, StackDependency)
- Component type enumerations and status tracking
- Stack validation and business logic
- Comprehensive unit tests for stack types

## Dependencies
- **REQUIRES**: Split 001 completion (imports oci package types)
- The stack package has: `import "github.com/cnoe-io/idpbuilder/pkg/oci"`
- Must ensure Split 001's OCI types are available
- External: Standard library, testing framework

## Implementation Instructions
1. **PREREQUISITE**: Ensure Split 001 is complete and merged
   ```bash
   # Verify Split 001 is available
   git log --oneline | grep "split-001"
   ```

2. Create split working directory:
   ```bash
   mkdir -p efforts/phase1/wave1/oci-types/split-002
   cd efforts/phase1/wave1/oci-types/split-002
   ```

3. Create branch for this split:
   ```bash
   git checkout -b phase1/wave1/oci-types-split-002
   ```

4. Copy ONLY the files assigned to this split:
   ```bash
   mkdir -p pkg/stack
   # Note: pkg/doc.go is in Split 001, not needed here
   cp ../pkg/stack/constants.go pkg/stack/
   cp ../pkg/stack/types.go pkg/stack/
   cp ../pkg/stack/types_test.go pkg/stack/
   ```

5. Ensure OCI package dependency is available:
   ```bash
   # May need to copy or reference OCI types from Split 001
   # Or update import paths as needed
   ```

6. Implement complete stack functionality
7. Ensure compilation: `go build ./pkg/stack/...`
8. Run unit tests: `go test ./pkg/stack/...`
9. Measure with line counter tool to verify compliance

## Testing Requirements
- All stack type methods must have tests
- Component lifecycle states must be tested
- Dependency resolution logic must be covered
- Test coverage target: >80%
- Tests must verify integration with OCI types

## Split Branch Strategy
- Branch: `phase1/wave1/oci-types-split-002`
- Must merge to: `phase1/wave1/oci-types` after review
- Depends on Split 001 being completed first
- Final split - completes the oci-types effort

## Integration Notes
- After both splits are complete and reviewed:
  1. Merge split-001 to phase1/wave1/oci-types
  2. Merge split-002 to phase1/wave1/oci-types
  3. Verify complete functionality works together
  4. Total implementation: 974 lines (661 + 313)
## 🚨 SPLIT INFRASTRUCTURE METADATA (Added by Orchestrator)
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase1/wave1/oci-types--split-002
**BRANCH**: phase1/wave1/oci-types--split-002
**REMOTE**: origin/phase1/wave1/oci-types--split-002
**BASE_BRANCH**: phase1/wave1/oci-types--split-001
**SPLIT_NUMBER**: 002
**TOTAL_SPLITS**: 2

### SW Engineer Instructions (R205)
1. READ this metadata FIRST
2. cd to WORKING_DIRECTORY above
3. Verify branch matches BRANCH above
4. ONLY THEN proceed with preflight checks
