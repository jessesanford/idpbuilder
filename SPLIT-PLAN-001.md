# SPLIT-PLAN-001.md

## Split 001 of 3: Core Storage and Configuration
**Planner**: Code Reviewer code-reviewer (same for ALL splits)
**Parent Effort**: buildah-integration

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️ CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
  
- **This Split**: Split 001 of phase2/wave1/buildah-integration
  - Path: efforts/phase2/wave1/buildah-integration/split-001/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-001
  
- **Next Split**: Split 002 of phase2/wave1/buildah-integration
  - Path: efforts/phase2/wave1/buildah-integration/split-002/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-002

- **File Boundaries**:
  - This Split Start: Line 1 / File: pkg/oci/build/store.go
  - This Split End: Line 747 / File: pkg/oci/build/config.go (line 442)
  - Next Split Start: Line 748 / File: pkg/oci/build/runtime.go

## Files in This Split (EXCLUSIVE - no overlap with other splits)

1. **pkg/oci/build/store.go** (305 lines)
   - Storage backend initialization
   - Buildah store management
   - Rootless storage configuration
   - Storage lifecycle operations

2. **pkg/oci/build/config.go** (442 lines)
   - BuildConfig structure definition
   - Configuration validation
   - ConfigManager implementation
   - Storage and runtime options setup

**Total Lines**: 747 (COMPLIANT - under 800 limit)

## Functionality

This split provides the foundational components for Buildah integration:

- **Storage Management**:
  - Initialize Buildah storage backend
  - Configure storage drivers (vfs, overlay, btrfs, zfs)
  - Handle rootless storage setup
  - Manage storage lifecycle and cleanup

- **Configuration Management**:
  - Define core BuildConfig structure
  - Validate configuration parameters
  - Prepare storage options
  - Setup system context for operations

## Dependencies

- **External**: None (foundational split)
- **Internal**: None (first split)
- **Imports Required**:
  - github.com/containers/buildah
  - github.com/containers/storage
  - github.com/containers/image/v5/types

## Implementation Instructions

### Step 1: Clean up existing branch
```bash
# The branch already exists but has excess files
git checkout idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-001
rm -rf split-001/  # Remove duplicate directory structure
```

### Step 2: Create proper structure
```bash
# Ensure only these files exist in pkg/oci/build/
pkg/oci/build/store.go    # 305 lines
pkg/oci/build/config.go   # 442 lines (use the main one, not split-001 version)
```

### Step 3: Verify no duplicates
- Ensure NO files from split-001/pkg/ directory
- Use only the main pkg/oci/build/ versions
- Remove extended_config functionality (goes to Split 003)

### Step 4: Update imports
- Ensure all imports are properly resolved
- Add go.mod and go.sum for this split

### Step 5: Compile and test
```bash
go mod init github.com/idpbuilder/idpbuilder-oci-mgmt
go mod tidy
go build ./pkg/oci/build/...
```

### Step 6: Measure compliance
```bash
PROJECT_ROOT="/home/vscode/workspaces/idpbuilder-oci-mgmt"
$PROJECT_ROOT/tools/line-counter.sh -c idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-001 -b origin/software-factory-2.0
# Must show ≤800 lines
```

## Split Branch Strategy

- **Branch Name**: `idpbuilder-oci-mgmt/phase2/wave1/buildah-integration-split-001`
- **Base Branch**: `origin/software-factory-2.0`
- **Merge Target**: `phase2/wave1/buildah-integration` (after all splits complete)

## Quality Checklist

- [ ] Only store.go and config.go present
- [ ] No duplicate files from split-001/ directory
- [ ] Total lines ≤800 (target: 747)
- [ ] Code compiles independently
- [ ] All imports resolved
- [ ] No functionality from other splits

## Known Issues to Address

1. **Existing branch has 1870 lines** - Must be cleaned to only include the 747 lines for this split
2. **Duplicate config.go files** - Use only the 442-line version from pkg/oci/build/
3. **Remove split-001/ directory** - This shouldn't exist in the split branch

## Integration Notes

- This split provides core types that Split 002 and 003 will import
- BuildConfig struct will be needed by runtime.go in Split 002
- StoreManager will be referenced by extended config in Split 003

## Success Criteria

✅ Line count ≤800 (target: 747)
✅ Only designated files present
✅ No duplication with other splits
✅ Compiles independently
✅ Provides exported types for other splits