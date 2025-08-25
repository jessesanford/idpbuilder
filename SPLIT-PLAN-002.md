# SPLIT-PLAN-002.md
## Split 002 of 2: Stack Types & Validation Logic
**Planner**: Code Reviewer code-reviewer-1756148460 (same for ALL splits)
**Parent Effort**: oci-stack-types

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 001 of phase1/wave1/oci-stack-types
  - Path: efforts/phase1/wave1/oci-stack-types/split-001/
  - Branch: idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types-split-001
  - Summary: Implemented core interfaces and base types (460 lines)
- **This Split**: Split 002 of phase1/wave1/oci-stack-types
  - Path: efforts/phase1/wave1/oci-stack-types/split-002/
  - Branch: idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types-split-002
- **Next Split**: None (final split)
  - Path: N/A
  - Branch: N/A

## Files in This Split (EXCLUSIVE - no overlap with Split 001)

### New File: pkg/oci/api/stack_types.go
1. **pkg/oci/api/stack_types.go** (~141 lines - extracted from types.go)
   - StackOCIConfig struct (from types.go lines 91-140)
   - StackImageInfo struct (from types.go lines 381-393)
   - StackHistoryEntry struct (from types.go lines 395-423)
   - ProgressEvent struct (from types.go lines 425-453)

### Complete File: validation.go
2. **pkg/oci/api/validation.go** (314 lines - complete)
   - All validation logic
   - Custom validators
   - Business logic validation
   - Helper functions

### Test File
3. **pkg/oci/api/validation_test.go** (~100 lines - NEW)
   - Validation function tests
   - Custom validator tests
   - Business logic validation tests
   - Edge cases and error conditions

## Functionality
This split provides stack-specific types and all validation:
- **Stack Configuration**: StackOCIConfig with versioning and metadata
- **Stack Information**: StackImageInfo and StackHistoryEntry
- **Progress Tracking**: ProgressEvent for build/push operations
- **Validation Logic**: Complete validation for all types
- **Custom Validators**: Image tag, semver, platform validators
- **Business Rules**: Rootless mode, authentication, timeout validations

## Dependencies
- **Split 001 Types**: Imports core types from Split 001
  - BuildConfig, RegistryConfig (for validation functions)
  - ImageInfo (embedded in StackImageInfo)
  - LayerInfo (referenced in validation)
- External: github.com/go-playground/validator/v10

## Implementation Instructions

### Step 1: Ensure Split 001 is Complete
```bash
# Verify Split 001 has been merged or is available
cd /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase1/wave1/oci-stack-types
git log --oneline | grep "split-001"
```

### Step 2: Create Split Working Directory
```bash
cd /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase1/wave1/oci-stack-types
git checkout -b idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types-split-002
mkdir -p split-002/pkg/oci/api
cd split-002
```

### Step 3: Copy Required Files from Split 001
```bash
# Copy the interfaces and base types from Split 001 for compilation
cp ../split-001/pkg/oci/api/interfaces.go pkg/oci/api/
cp ../split-001/pkg/oci/api/types.go pkg/oci/api/
```

### Step 4: Create stack_types.go
Extract from original types.go:
```go
package api

import "time"

// StackOCIConfig defines stack-specific OCI configuration.
// (lines 91-140 from original types.go)
type StackOCIConfig struct {
    // ... complete struct definition
}

// StackImageInfo extends ImageInfo with stack-specific information.
// (lines 381-393 from original types.go)
type StackImageInfo struct {
    *ImageInfo
    // ... additional fields
}

// StackHistoryEntry represents a single entry in stack build/push history.
// (lines 395-423 from original types.go)
type StackHistoryEntry struct {
    // ... complete struct definition
}

// ProgressEvent represents a progress update during build or push operations.
// (lines 425-453 from original types.go)
type ProgressEvent struct {
    // ... complete struct definition
}
```

### Step 5: Copy Complete validation.go
```bash
cp ../pkg/oci/api/validation.go pkg/oci/api/validation.go
```

### Step 6: Create validation_test.go
Create comprehensive tests for:
- ValidateBuildConfig function
- ValidateRegistryConfig function
- ValidateStackConfig function
- Custom validators (image_tag, semver, platform)
- Business logic validation (rootless mode, authentication)
- Helper functions (isValidDNSLabel, isValidRepository)

### Step 7: Verify Compilation
```bash
cd split-002
go mod init github.com/cnoe-io/idpbuilder-oci-mgmt
go get github.com/go-playground/validator/v10
go build ./pkg/oci/api/...
go test ./pkg/oci/api/...
```

### Step 8: Measure Size
```bash
# Count only the NEW files for this split
wc -l pkg/oci/api/stack_types.go pkg/oci/api/validation.go pkg/oci/api/validation_test.go
# Target: ~455 lines total
```

## Integration with Split 001
After both splits are complete:
1. Merge Split 001 to parent branch
2. Merge Split 002 to parent branch
3. Reorganize files:
   - Combine types.go and stack_types.go back into single types.go
   - Keep validation.go separate
   - Combine test files if needed

## Split Branch Strategy
- Branch: `idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types-split-002`
- Prerequisite: Split 001 must be complete for imports
- Must merge to: `idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types` after review

## Quality Requirements
- Must compile with Split 001 types available
- Must pass all validation tests
- Must be under 500 lines total
- Must maintain 80% test coverage for validation logic
- Must follow Go best practices

## Verification Checklist
- [ ] Stack types properly extracted
- [ ] All validation logic included
- [ ] Imports from Split 001 work correctly
- [ ] Validation tests comprehensive
- [ ] Business logic properly validated
- [ ] Size under 500 lines
- [ ] Ready for review

## Notes for Implementation
1. This split DEPENDS on Split 001 - ensure those types are available
2. The validation.go file references types from Split 001 (BuildConfig, etc.)
3. Focus on comprehensive validation testing
4. Ensure custom validators are properly registered
5. The split maintains logical cohesion around stack-specific functionality