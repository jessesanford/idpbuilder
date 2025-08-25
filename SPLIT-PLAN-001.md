# SPLIT-PLAN-001.md
## Split 001 of 2: Core Interfaces & Base Types
**Planner**: Code Reviewer code-reviewer-1756148460 (same for ALL splits)
**Parent Effort**: oci-stack-types

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase1/wave1/oci-stack-types
  - Path: efforts/phase1/wave1/oci-stack-types/split-001/
  - Branch: idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types-split-001
- **Next Split**: Split 002 of phase1/wave1/oci-stack-types
  - Path: efforts/phase1/wave1/oci-stack-types/split-002/
  - Branch: idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types-split-002
- **File Boundaries**:
  - This Split Start: pkg/oci/api/interfaces.go (line 1)
  - This Split End: pkg/oci/api/types.go (line 311)
  - Next Split Start: pkg/oci/api/types.go (line 312)

## Files in This Split (EXCLUSIVE - no overlap with other splits)

### Complete Files
1. **pkg/oci/api/interfaces.go** (149 lines)
   - All 5 interface definitions
   - OCIBuildService, OCIRegistryService, StackOCIManager
   - ProgressReporter, LayerProcessor

### Partial File: pkg/oci/api/types.go (Lines 1-311)
2. **pkg/oci/api/types.go** (311 lines from original 452)
   - Package declaration and imports (lines 1-9)
   - BuildConfig struct (lines 10-55)
   - RegistryConfig struct (lines 57-89)
   - BuildRequest struct (lines 142-179)
   - BuildResult struct (lines 181-209)
   - BuildStatus struct (lines 211-236)
   - BuildPhase enum and constants (lines 238-262)
   - BuildOptions struct (lines 264-295)
   - PushOptions struct (lines 297-310)
   - PullOptions struct (lines 312-325)
   - LayerInfo struct (lines 327-349)
   - ImageInfo struct (lines 351-379)
   
   **STOP at line 311** - Do NOT include:
   - StackOCIConfig (starts line 91)
   - StackImageInfo (starts line 381)
   - StackHistoryEntry (starts line 395)
   - ProgressEvent (starts line 425)

### Test File
3. **pkg/oci/api/types_test.go** (~100 lines - NEW)
   - Tests for BuildConfig validation
   - Tests for RegistryConfig validation
   - Tests for BuildRequest validation
   - Tests for enum values

## Functionality
This split provides the foundational contracts and types:
- **Service Contracts**: Complete interface definitions for all OCI services
- **Core Configuration**: Build and registry configuration structures
- **Request/Response Types**: All basic operation types
- **Options Types**: Build, push, and pull options
- **Information Types**: Layer and image information structures

## Dependencies
- None (foundational split)
- External: github.com/go-playground/validator/v10 (for validation tags only)

## Implementation Instructions

### Step 1: Create Split Working Directory
```bash
cd /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase1/wave1/oci-stack-types
git checkout -b idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types-split-001
mkdir -p split-001/pkg/oci/api
cd split-001
```

### Step 2: Copy Complete File - interfaces.go
```bash
cp ../pkg/oci/api/interfaces.go pkg/oci/api/interfaces.go
```

### Step 3: Extract Partial types.go (Lines 1-311 + selected structs)
Create `pkg/oci/api/types.go` with:
- Package header and imports
- Core configuration types (BuildConfig, RegistryConfig)
- Core request/response types (BuildRequest, BuildResult, BuildStatus)
- Enum types (BuildPhase and constants)
- Option types (BuildOptions, PushOptions, PullOptions)
- Info types (LayerInfo, ImageInfo)

**Critical**: The file must compile independently! Include these specific types:
```go
// Include lines 1-9 (package and imports)
// Include lines 10-55 (BuildConfig)
// Include lines 57-89 (RegistryConfig)
// SKIP lines 91-140 (StackOCIConfig)
// Include lines 142-179 (BuildRequest)
// Include lines 181-209 (BuildResult)
// Include lines 211-262 (BuildStatus and BuildPhase)
// Include lines 264-325 (Options types)
// Include lines 327-379 (LayerInfo, ImageInfo)
// STOP - do not include Stack types or ProgressEvent
```

### Step 4: Create Test File
Create `pkg/oci/api/types_test.go` with tests for:
- Basic struct creation
- JSON marshaling/unmarshaling
- Field validation using struct tags
- Enum constant values

### Step 5: Verify Compilation
```bash
cd split-001
go mod init github.com/cnoe-io/idpbuilder-oci-mgmt
go get github.com/go-playground/validator/v10
go build ./pkg/oci/api/...
go test ./pkg/oci/api/...
```

### Step 6: Measure Size
```bash
find pkg -name "*.go" -type f | xargs wc -l
# Target: ~460 lines total
```

## Split Branch Strategy
- Branch: `idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types-split-001`
- Must merge to: `idpbuidler-oci-mgmt/phase1/wave1/oci-stack-types` after review
- Integration: This split must be complete before Split 002 begins

## Quality Requirements
- Must compile independently
- Must pass all tests
- Must be under 500 lines total
- Must include complete documentation
- Must follow Go best practices

## Verification Checklist
- [ ] All interfaces included and complete
- [ ] Core types properly extracted
- [ ] No stack-specific types included
- [ ] File compiles independently
- [ ] Tests provide adequate coverage
- [ ] Size under 500 lines
- [ ] Ready for review