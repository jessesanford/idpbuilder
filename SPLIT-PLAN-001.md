# SPLIT-PLAN-001.md
## Split 001 of 2: Core Error Reporting Types
**Planner**: Code Reviewer code-reviewer-1756085350 (same for ALL splits)
**Parent Effort**: error-reporting-types

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase1/wave1/effort-error-reporting-types
  - Path: efforts/phase1/wave1/effort-error-reporting-types/split-001/
  - Branch: phase1/wave1/error-reporting-types-split-001
- **Next Split**: Split 002 of phase1/wave1/effort-error-reporting-types
  - Path: efforts/phase1/wave1/effort-error-reporting-types/split-002/
  - Branch: phase1/wave1/error-reporting-types-split-002

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- `pkg/errors/types.go` (191 lines) - Core error interface and implementation
- `pkg/errors/codes.go` (148 lines) - Error code definitions and helpers
- `pkg/errors/constants.go` (65 lines) - Error category and code constants

**Total Lines**: 404 lines (well under 800 limit)

## Functionality
This split implements the core error reporting system:
- **OCIError Interface**: Standard interface for all OCI management errors
- **BaseError Type**: Concrete implementation with wrapping and context support
- **ErrorContext Type**: Structured context information for errors
- **ErrorStack Type**: Chain of related errors for tracking propagation
- **Error Codes**: Comprehensive error code system with categories
- **Error Categories**: Transient, permanent, configuration, validation categories

## Dependencies
- None (foundational split)
- Standard library only (fmt, time)

## Implementation Instructions
1. Create split branch from parent effort branch
2. Set up isolated workspace in split-001 directory
3. Create pkg/errors directory structure
4. Implement the three core files:
   - Start with constants.go (defines categories)
   - Then codes.go (uses categories from constants)
   - Finally types.go (uses both constants and codes)
5. Ensure each file compiles independently
6. Verify with: `go build ./pkg/errors`
7. Measure final size with: `${PROJECT_ROOT}/tools/line-counter.sh`
8. Commit with clear message about split 001 contents

## Testing Strategy
- Unit tests will be implemented in Split 002
- This split focuses on type definitions and compilation
- Manual verification that interfaces are properly defined

## Integration Points
- These types will be used throughout the OCI management system
- Error codes provide programmatic error handling
- Error wrapping enables detailed error chains
- Context support allows rich error information

## Split Branch Strategy
- Branch: `phase1/wave1/error-reporting-types-split-001`
- Base: `phase1/wave1/error-reporting-types`
- Must merge back to parent branch after review
- No direct merge to main branch
## 🚨 SPLIT INFRASTRUCTURE METADATA (Added by Orchestrator)
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase1/wave1/error-reporting-types--split-001
**BRANCH**: phase1/wave1/error-reporting-types--split-001
**REMOTE**: origin/phase1/wave1/error-reporting-types--split-001
**BASE_BRANCH**: main
**SPLIT_NUMBER**: 001
**TOTAL_SPLITS**: 2

### SW Engineer Instructions (R205)
1. READ this metadata FIRST
2. cd to WORKING_DIRECTORY above
3. Verify branch matches BRANCH above
4. ONLY THEN proceed with preflight checks
