# SPLIT-PLAN-003.md

## Split 003 of 3: Extended Configuration and Integration
**Planner**: Code Reviewer code-reviewer (same for ALL splits)
**Parent Effort**: buildah-integration

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️ CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: Split 002 of phase2/wave1/buildah-integration
  - Path: efforts/phase2/wave1/buildah-integration/split-002/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-002
  - Summary: Implemented runtime management and tests (730 lines)
  
- **This Split**: Split 003 of phase2/wave1/buildah-integration (FINAL)
  - Path: efforts/phase2/wave1/buildah-integration/split-003/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-003
  
- **Next Split**: None (this is the final split)

- **File Boundaries**:
  - This Split Start: Line 1479 / File: pkg/oci/build/extended_config.go
  - This Split End: Line 1939 / File: pkg/oci/build/extended_config.go (line 460)
  - Next Split Start: N/A (final split)

## Files in This Split (EXCLUSIVE - no overlap with other splits)

1. **pkg/oci/build/extended_config.go** (460 lines)
   - Extended configuration features
   - Advanced validation logic
   - Integration helpers
   - Utility functions for buildah operations
   - Configuration merging and inheritance

**Total Lines**: 460 (HIGHLY COMPLIANT - well under 800 limit)

## Important Note About File Origin

The `extended_config.go` file is actually the content from:
- Original location: `split-001/pkg/oci/build/config.go` (460 lines)
- This appears to be an extended/enhanced version of the base config.go
- It should be renamed to `extended_config.go` to avoid confusion

## Functionality

This split provides advanced configuration and integration features:

- **Extended Configuration**:
  - Advanced configuration options beyond base config
  - Configuration inheritance and merging
  - Environment-specific overrides
  - Profile-based configuration

- **Integration Helpers**:
  - Builder pattern for complex configurations
  - Configuration validation utilities
  - Helper functions for common buildah operations
  - Configuration serialization/deserialization

- **Utility Functions**:
  - Path resolution helpers
  - Security policy utilities
  - Network configuration helpers
  - Cache management utilities

## Dependencies

- **From Split 001**:
  - BuildConfig base struct
  - ConfigManager interface
  - StoreManager types
  
- **From Split 002**:
  - RuntimeManager types
  - Runtime configuration structures

- **External Imports**:
  - github.com/containers/buildah
  - github.com/containers/storage
  - github.com/containers/image/v5/types

## Implementation Instructions

### Step 1: Create new branch
```bash
git checkout -b idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-003 origin/software-factory-2.0
```

### Step 2: Create directory structure
```bash
mkdir -p efforts/phase2/wave1/buildah-integration/split-003/pkg/oci/build
cd efforts/phase2/wave1/buildah-integration/split-003
```

### Step 3: Copy and rename the extended config file
```bash
# Copy the extended config from split-001 directory and rename it
cp ../split-001/pkg/oci/build/config.go pkg/oci/build/extended_config.go

# Update the package declaration if needed
# Ensure it's still: package build
```

### Step 4: Update imports and references
- Change any self-references to use the new filename
- Import types from Split 001 and Split 002 as needed
- Update any struct names to avoid conflicts (e.g., ExtendedConfigManager)

### Step 5: Add dependency resolution
Since this split depends on types from Split 001 and 002:
```go
// In extended_config.go, add imports:
import (
    // Types from Split 001
    baseconfig "github.com/idpbuilder/idpbuilder-oci-mgmt/split-001/pkg/oci/build"
    // Types from Split 002
    runtime "github.com/idpbuilder/idpbuilder-oci-mgmt/split-002/pkg/oci/build"
)
```

### Step 6: Create go.mod
```go
module github.com/idpbuilder/idpbuilder-oci-mgmt/split-003

go 1.21

require (
    github.com/containers/buildah v1.30.0
    github.com/containers/storage v1.48.0
    github.com/containers/image/v5 v5.26.1
    // Dependencies on other splits
    github.com/idpbuilder/idpbuilder-oci-mgmt/split-001 v0.0.0
    github.com/idpbuilder/idpbuilder-oci-mgmt/split-002 v0.0.0
)

replace (
    github.com/idpbuilder/idpbuilder-oci-mgmt/split-001 => ../split-001
    github.com/idpbuilder/idpbuilder-oci-mgmt/split-002 => ../split-002
)
```

### Step 7: Compile and verify
```bash
go mod tidy
go build ./pkg/oci/build/...
```

### Step 8: Measure compliance
```bash
PROJECT_ROOT="/home/vscode/workspaces/idpbuilder-oci-mgmt"
$PROJECT_ROOT/tools/line-counter.sh -c idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-003 -b origin/software-factory-2.0
# Must show ≤800 lines (should be ~460)
```

## Split Branch Strategy

- **Branch Name**: `idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-003`
- **Base Branch**: `origin/software-factory-2.0`
- **Merge Target**: `phase2/wave1/buildah-integration` (after all splits complete)

## Quality Checklist

- [ ] Only extended_config.go present (renamed from split-001/config.go)
- [ ] Dependencies on Split 001 & 002 resolved
- [ ] Total lines ≤800 (target: 460)
- [ ] No duplicate functionality with base config.go
- [ ] All imports resolved
- [ ] Compiles successfully

## Integration Strategy

This is the final split that brings together advanced features:

1. **Extends base configuration** from Split 001
2. **Leverages runtime capabilities** from Split 002
3. **Provides high-level integration** for the complete buildah solution

## File Rename Justification

Renaming `split-001/pkg/oci/build/config.go` to `extended_config.go`:
- Avoids confusion with the base `config.go` in Split 001
- Clearly indicates this contains extended/advanced features
- Prevents import conflicts when merging splits

## Final Integration Notes

After all three splits are complete and reviewed:
1. Merge Split 001 to parent branch
2. Merge Split 002 to parent branch
3. Merge Split 003 to parent branch
4. Ensure no duplicate code in final integration
5. Verify combined functionality works as expected

## Success Criteria

✅ Line count ≤800 (target: 460)
✅ Extended configuration features complete
✅ No duplication with base config
✅ Dependencies properly resolved
✅ Compiles independently
✅ Ready for final integration