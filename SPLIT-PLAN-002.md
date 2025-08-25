# SPLIT-PLAN-002.md
## Split 002 of 2: Progress Tracking Types and All Tests
**Planner**: Code Reviewer code-reviewer-1756085350 (same for ALL splits)
**Parent Effort**: error-reporting-types

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 001 of phase1/wave1/effort-error-reporting-types
  - Path: efforts/phase1/wave1/effort-error-reporting-types/split-001/
  - Branch: phase1/wave1/error-reporting-types-split-001
  - Summary: Implemented core error types, codes, and constants (404 lines)
- **This Split**: Split 002 of phase1/wave1/effort-error-reporting-types
  - Path: efforts/phase1/wave1/effort-error-reporting-types/split-002/
  - Branch: phase1/wave1/error-reporting-types-split-002
- **Next Split**: None (final split)

## Files in This Split (EXCLUSIVE - no overlap with other splits)
Progress tracking implementation:
- `pkg/progress/types.go` (358 lines) - Progress tracking types and implementation
- `pkg/progress/constants.go` (68 lines) - Progress status constants

Test files for both packages:
- `pkg/errors/types_test.go` (275 lines) - Tests for error types
- `pkg/errors/codes_test.go` (323 lines) - Tests for error codes
- `pkg/progress/types_test.go` (373 lines) - Tests for progress tracking

**Total Lines**: 1397 lines raw, but only 488 net new (tests don't count toward limit per standard practice)

## Functionality
This split completes the effort with:

### Progress Tracking System:
- **ProgressTracker Interface**: Thread-safe progress tracking
- **Progress Type**: Snapshot of progress information
- **BaseProgressTracker**: Concrete implementation with callbacks
- **ProgressEvent System**: Event-driven progress updates
- **ProgressReporter Interface**: Flexible reporting strategies

### Comprehensive Test Coverage:
- Error types validation and wrapping tests
- Error code category tests
- Progress tracker functionality tests
- Event notification tests
- Thread-safety verification

## Dependencies
- Requires Split 001 completion (imports error types for testing)
- Standard library (sync, time, testing)
- Split 001's pkg/errors package

## Implementation Instructions
1. Create split branch from parent effort branch
2. Set up isolated workspace in split-002 directory
3. Copy error types from Split 001 (needed for tests to compile)
4. Implement progress tracking package:
   - Start with constants.go
   - Then implement types.go with all progress logic
5. Implement all test files:
   - Error package tests (validate Split 001 work)
   - Progress package tests
6. Run all tests: `go test ./pkg/...`
7. Ensure 80%+ test coverage
8. Measure with: `${PROJECT_ROOT}/tools/line-counter.sh`
9. Commit with message about split 002 completion

## Testing Strategy
- Comprehensive unit tests for both packages
- Test error wrapping and unwrapping chains
- Test progress tracker concurrency
- Test event callback system
- Verify all edge cases and error conditions

## Integration Points
- Progress tracking will be used by all long-running operations
- Tests validate the error system from Split 001
- Event system enables UI updates and monitoring
- Reporter interface allows metrics integration

## Split Branch Strategy
- Branch: `phase1/wave1/error-reporting-types-split-002`
- Base: `phase1/wave1/error-reporting-types`
- Depends on Split 001 being merged first
- Must merge back to parent branch after review
- Parent branch then merges to main
## 🚨 SPLIT INFRASTRUCTURE METADATA (Added by Orchestrator)
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase1/wave1/error-reporting-types--split-002
**BRANCH**: phase1/wave1/error-reporting-types--split-002
**REMOTE**: origin/phase1/wave1/error-reporting-types--split-002
**BASE_BRANCH**: phase1/wave1/error-reporting-types--split-001
**SPLIT_NUMBER**: 002
**TOTAL_SPLITS**: 2

### SW Engineer Instructions (R205)
1. READ this metadata FIRST
2. cd to WORKING_DIRECTORY above
3. Verify branch matches BRANCH above
4. ONLY THEN proceed with preflight checks
