# SPLIT-PLAN-002.md

## Split 002 of 3: Runtime Management and Testing
**Planner**: Code Reviewer code-reviewer (same for ALL splits)
**Parent Effort**: buildah-integration

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️ CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: Split 001 of phase2/wave1/buildah-integration
  - Path: efforts/phase2/wave1/buildah-integration/split-001/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-001
  - Summary: Implemented core storage and configuration types (747 lines)
  
- **This Split**: Split 002 of phase2/wave1/buildah-integration
  - Path: efforts/phase2/wave1/buildah-integration/split-002/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-002
  
- **Next Split**: Split 003 of phase2/wave1/buildah-integration
  - Path: efforts/phase2/wave1/buildah-integration/split-003/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-003

- **File Boundaries**:
  - This Split Start: Line 748 / File: pkg/oci/build/runtime.go
  - This Split End: Line 1478 / File: pkg/oci/build/config_test.go (line 329)
  - Next Split Start: Line 1479 / File: pkg/oci/build/extended_config.go

⚠️ NEVER reference splits from different efforts!
❌ WRONG: "Previous Split: Split 001 (oci-types)" when you're in buildah-integration
✅ RIGHT: "Previous Split: Split 001 of phase2/wave1/buildah-integration"

## Files in This Split (EXCLUSIVE - no overlap with other splits)

1. **pkg/oci/build/runtime.go** (401 lines)
   - RuntimeManager implementation
   - Rootless operation setup
   - Namespace configuration
   - Security capabilities management
   - Runtime validation and initialization

2. **pkg/oci/build/config_test.go** (329 lines)
   - Unit tests for ConfigManager
   - Validation test cases
   - Storage configuration tests
   - Integration test helpers

**Total Lines**: 730 (COMPLIANT - under 800 limit)

## Functionality

This split provides runtime management and testing capabilities:

- **Runtime Management**:
  - Initialize runtime environment for Buildah
  - Configure rootless UID/GID mappings
  - Setup security capabilities
  - Manage container runtime paths
  - Handle namespace isolation

- **Test Coverage**:
  - Configuration validation tests
  - Storage setup tests
  - Error handling verification
  - Integration test utilities

## Dependencies

- **From Split 001**:
  - BuildConfig struct
  - ConfigManager type
  - Storage configuration types
  
- **External Imports**:
  - github.com/opencontainers/runtime-spec/specs-go
  - github.com/containers/buildah
  - Standard testing packages

## Implementation Instructions

### Step 1: Create new branch
```bash
git checkout -b idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-002 origin/software-factory-2.0
```

### Step 2: Create directory structure
```bash
mkdir -p efforts/phase2/wave1/buildah-integration/split-002/pkg/oci/build
cd efforts/phase2/wave1/buildah-integration/split-002
```

### Step 3: Copy required files
```bash
# Copy runtime.go from main implementation
cp ../pkg/oci/build/runtime.go pkg/oci/build/

# Copy config_test.go from split-001 directory
cp ../split-001/pkg/oci/build/config_test.go pkg/oci/build/
```

### Step 4: Add dependency imports
Since this split depends on types from Split 001, you'll need to:
- Import the config types from Split 001's module
- Or copy just the type definitions needed (BuildConfig struct)

### Step 5: Create go.mod
```go
module github.com/idpbuilder/idpbuilder-oci-mgmt/split-002

go 1.21

require (
    github.com/containers/buildah v1.30.0
    github.com/opencontainers/runtime-spec v1.1.0
    // Add other dependencies
)
```

### Step 6: Compile and test
```bash
go mod tidy
go build ./pkg/oci/build/...
go test ./pkg/oci/build/...
```

### Step 7: Measure compliance
```bash
PROJECT_ROOT="/home/vscode/workspaces/idpbuilder-oci-mgmt"
$PROJECT_ROOT/tools/line-counter.sh -c idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-002 -b origin/software-factory-2.0
# Must show ≤800 lines
```

## Split Branch Strategy

- **Branch Name**: `idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-002`
- **Base Branch**: `origin/software-factory-2.0`
- **Merge Target**: `phase2/wave1/buildah-integration` (after all splits complete)

## Quality Checklist

- [ ] Only runtime.go and config_test.go present
- [ ] Dependencies on Split 001 types resolved
- [ ] Total lines ≤800 (target: 730)
- [ ] Tests pass successfully
- [ ] All imports resolved
- [ ] No functionality from other splits

## Integration Requirements

### Importing from Split 001:
The RuntimeManager in runtime.go needs the BuildConfig type from Split 001. Options:

1. **Module dependency** (preferred):
   ```go
   import "github.com/idpbuilder/idpbuilder-oci-mgmt/split-001/pkg/oci/build"
   ```

2. **Type duplication** (if module deps problematic):
   - Copy only the BuildConfig struct definition
   - Ensure it matches exactly with Split 001

### Test Dependencies:
The config_test.go file tests ConfigManager which is in Split 001. You'll need to:
- Import ConfigManager from Split 001
- Or adapt tests to only test runtime functionality

## Known Challenges

1. **Cross-split dependencies** - Runtime depends on config types from Split 001
2. **Test file location** - config_test.go tests Split 001 code but lives in Split 002
3. **Import resolution** - Need proper module setup for cross-split imports

## Success Criteria

✅ Line count ≤800 (target: 730)
✅ Runtime functionality complete
✅ Tests execute successfully
✅ Dependencies on Split 001 resolved
✅ No duplication with other splits
✅ Compiles and tests independently