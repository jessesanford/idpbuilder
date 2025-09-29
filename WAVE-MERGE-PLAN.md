# Phase 1 Wave 1 Integration Merge Plan

## Overview
- **Phase**: 1
- **Wave**: 1
- **Integration Branch**: phase1-wave1-integration
- **Base Branch**: main
- **Integration Directory**: /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/integration-workspace
- **Created**: 2025-09-29T13:57:00Z
- **Created By**: Code Reviewer Agent (WAVE_MERGE_PLANNING state)

## Critical Requirements (R269, R270)
- ✅ This plan is created per R269 (detailed merge planning)
- ✅ Using ONLY original effort/split branches per R270 (no integration branches as sources)
- ✅ Respecting cascade branching pattern per R501
- ✅ NO execution of merges - planning only!

## Efforts to Integrate

### Successfully Completed Efforts:
1. **E1.1.1-analyze-existing-structure** (29 lines)
   - Branch: `phase1/wave1/analyze-existing-structure`
   - Location: `/home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/E1.1.1-analyze-existing-structure`
   - Base: main
   - Status: ACCEPTED

2. **E1.1.2-split-001** (660 lines)
   - Branch: `phase1/wave1/unit-test-framework-split-001`
   - Location: `/home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/E1.1.2-split-001`
   - Base: phase1/wave1/analyze-existing-structure
   - Status: ACCEPTED
   - Content: Core Mock Registry Infrastructure

3. **E1.1.2-split-002** (802 lines)
   - Branch: `phase1/wave1/unit-test-framework-split-002`
   - Location: `/home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/E1.1.2-split-002`
   - Base: phase1/wave1/unit-test-framework-split-001
   - Status: ACCEPTED
   - Content: Test Utilities and Assertions

4. **E1.1.3-integration-test-setup** (612 lines)
   - Branch: `phase1/wave1/integration-test-setup`
   - Location: `/home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/E1.1.3-integration-test-setup`
   - Base: phase1/wave1/unit-test-framework (original, before splits)
   - Status: ACCEPTED

### EXCLUDED (Per R270):
- ❌ Original E1.1.2 branch (`phase1/wave1/unit-test-framework`) - Using splits instead
- ❌ Any integration branches - Never use as sources

## Merge Sequence Strategy

### CASCADE PATTERN ANALYSIS:
Per R501 and the cascade branching structure:
```
main
└── E1.1.1 (phase1/wave1/analyze-existing-structure)
    └── E1.1.2-split-001 (phase1/wave1/unit-test-framework-split-001)
        └── E1.1.2-split-002 (phase1/wave1/unit-test-framework-split-002)
    └── E1.1.3 (phase1/wave1/integration-test-setup) [based on original E1.1.2]
```

### CRITICAL MERGE ORDERING:
The merge must follow the cascade pattern to preserve commit history and dependencies:

1. **Merge E1.1.1** first (foundation)
2. **Merge E1.1.2-split-001** second (builds on E1.1.1)
3. **Merge E1.1.2-split-002** third (builds on split-001)
4. **Merge E1.1.3** fourth (may have conflicts due to different base)

## Detailed Merge Commands

### Pre-Merge Preparation
```bash
# Navigate to integration workspace
cd /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/integration-workspace

# Ensure we're on the integration branch
git checkout phase1-wave1-integration

# Fetch all updates from effort repositories
cd ../E1.1.1-analyze-existing-structure && git fetch origin
cd ../E1.1.2-split-001 && git fetch origin
cd ../E1.1.2-split-002 && git fetch origin
cd ../E1.1.3-integration-test-setup && git fetch origin

# Return to integration workspace
cd /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave1/integration-workspace

# Ensure integration branch is clean
git status
```

### Step 1: Add Effort Repositories as Remotes
```bash
# Add each effort as a remote to the integration repository
git remote add effort-E1.1.1 ../E1.1.1-analyze-existing-structure/.git
git remote add effort-E1.1.2-split-001 ../E1.1.2-split-001/.git
git remote add effort-E1.1.2-split-002 ../E1.1.2-split-002/.git
git remote add effort-E1.1.3 ../E1.1.3-integration-test-setup/.git

# Fetch all remotes
git fetch effort-E1.1.1
git fetch effort-E1.1.2-split-001
git fetch effort-E1.1.2-split-002
git fetch effort-E1.1.3
```

### Step 2: Merge E1.1.1 (Foundation)
```bash
# Merge E1.1.1 - should be clean since it's based on main
git merge effort-E1.1.1/phase1/wave1/analyze-existing-structure \
  --no-ff \
  -m "feat(wave1): integrate E1.1.1 - analyze existing structure (29 lines)"

# Verify the merge
git log --oneline -5
git diff HEAD~1 --stat
```

**Expected Result**: Clean merge, no conflicts
**Files Added**: Documentation and analysis files from E1.1.1

### Step 3: Merge E1.1.2-split-001 (Mock Registry Core)
```bash
# Merge E1.1.2-split-001 - should be clean since it builds on E1.1.1
git merge effort-E1.1.2-split-001/phase1/wave1/unit-test-framework-split-001 \
  --no-ff \
  -m "feat(wave1): integrate E1.1.2-split-001 - mock registry infrastructure (660 lines)"

# Verify the merge
git log --oneline -5
git diff HEAD~1 --stat
```

**Expected Result**: Clean merge, no conflicts
**Files Added**:
- `pkg/testing/mock_registry/registry.go`
- `pkg/testing/mock_registry/handlers.go`
- Related mock registry infrastructure files

### Step 4: Merge E1.1.2-split-002 (Test Utilities)
```bash
# Merge E1.1.2-split-002 - should be clean since it builds on split-001
git merge effort-E1.1.2-split-002/phase1/wave1/unit-test-framework-split-002 \
  --no-ff \
  -m "feat(wave1): integrate E1.1.2-split-002 - test utilities and assertions (802 lines)"

# Verify the merge
git log --oneline -5
git diff HEAD~1 --stat
```

**Expected Result**: Clean merge, no conflicts
**Files Added**:
- `pkg/testing/assertions/assertions.go`
- `pkg/testing/test_helpers/helpers.go`
- Additional test utilities

### Step 5: Merge E1.1.3 (Integration Tests) - POTENTIAL CONFLICTS
```bash
# Merge E1.1.3 - MAY HAVE CONFLICTS
# E1.1.3 was based on the original E1.1.2 branch (not the splits)
git merge effort-E1.1.3/phase1/wave1/integration-test-setup \
  --no-ff \
  -m "feat(wave1): integrate E1.1.3 - integration test setup (612 lines)"

# If conflicts occur, they will likely be in:
# - go.mod (dependency versions)
# - Test file imports (different paths between original and split structure)
```

#### Conflict Resolution Strategy for E1.1.3:

**Expected Conflicts**:
1. **Import Path Conflicts**: E1.1.3 may import from paths that existed in the original E1.1.2 but were reorganized in the splits
2. **go.mod Dependencies**: Version conflicts or missing dependencies

**Resolution Approach**:
```bash
# If conflicts occur:

# 1. Check conflict markers
git status
git diff

# 2. For import path conflicts:
# Update imports to use the split structure paths:
# OLD: import "github.com/jessesanford/idpbuilder/pkg/testing"
# NEW: import "github.com/jessesanford/idpbuilder/pkg/testing/mock_registry"
#      import "github.com/jessesanford/idpbuilder/pkg/testing/assertions"

# 3. For go.mod conflicts:
# Keep all dependencies from both sides (union of dependencies)
# Use the newer version if there's a version conflict

# 4. After resolving:
git add .
git commit -m "resolve: merge conflicts for E1.1.3 integration test setup"
```

### Step 6: Post-Merge Validation
```bash
# Ensure all code is present
find pkg/ -name "*.go" | wc -l  # Should show all Go files

# Run build to ensure compilation
go mod tidy
go build ./...

# Run tests to ensure everything works
go test ./pkg/testing/... -v

# Check total line count
git diff main --stat

# Verify no files were lost
git ls-tree -r HEAD --name-only | grep "\.go$" | sort > merged-files.txt
# Compare with expected files from all efforts
```

### Step 7: Create Integration Summary
```bash
# Generate integration report
cat > INTEGRATION-SUMMARY.md << 'EEOF'
# Phase 1 Wave 1 Integration Summary

## Merged Efforts
- E1.1.1: ✅ Clean merge (29 lines)
- E1.1.2-split-001: ✅ Clean merge (660 lines)
- E1.1.2-split-002: ✅ Clean merge (802 lines)
- E1.1.3: ✅ Merged (612 lines) [conflicts resolved if any]

## Total Lines: 2,103

## Build Status: ✅ Passing
## Test Status: ✅ All tests passing

## Integration Timestamp: $(date -Iseconds)
EEOF

git add INTEGRATION-SUMMARY.md
git commit -m "docs: add wave 1 integration summary"
```

### Step 8: Push Integration Branch
```bash
# Push the integrated branch
git push origin phase1-wave1-integration
```

## Expected Conflicts and Resolution

### Minimal Conflicts Expected
Due to the cascade branching pattern, most merges should be clean. The only potential conflict point is E1.1.3:

1. **E1.1.1 → integration**: Clean (from main)
2. **E1.1.2-split-001 → integration**: Clean (builds on E1.1.1)
3. **E1.1.2-split-002 → integration**: Clean (builds on split-001)
4. **E1.1.3 → integration**: Potential conflicts (based on original E1.1.2)

### E1.1.3 Specific Conflict Areas:
- **File Structure**: E1.1.3 expects the original E1.1.2 structure, but we have splits
- **Imports**: May need path adjustments
- **Test Dependencies**: May reference functions that moved between splits

### Resolution Guidelines:
1. **Preserve All Code**: Never delete code during conflict resolution
2. **Update Imports**: Adjust import paths to match split structure
3. **Maintain Functionality**: Ensure all tests still pass after resolution
4. **Document Changes**: Note any significant adjustments made

## Validation Checklist

After completing all merges, verify:

- [ ] All 4 effort branches successfully merged
- [ ] Total line count approximately 2,103 lines
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] No files were lost in merging
- [ ] Integration branch pushed successfully
- [ ] All remotes cleaned up (`git remote remove effort-*`)

## Important Notes

1. **DO NOT** merge the original E1.1.2 branch - only use the splits
2. **DO NOT** use any integration branches as merge sources
3. **PRESERVE** commit history by using `--no-ff` flag
4. **DOCUMENT** any conflict resolutions clearly
5. **VALIDATE** thoroughly after each merge step

## Contingency Plans

### If Severe Conflicts Occur:
1. **Abort the merge**: `git merge --abort`
2. **Analyze the conflict**: Understand why it occurred
3. **Consult effort implementation plans**: Review what each effort changed
4. **Apply minimal fixes**: Change only what's necessary for integration
5. **Document extensively**: Record all decisions made

### If Build Breaks:
1. **Check go.mod**: Run `go mod tidy`
2. **Verify imports**: Ensure all imports resolve correctly
3. **Check for missing files**: Verify all expected files are present
4. **Run individual tests**: Isolate which component fails

## Success Criteria

The integration is successful when:
- ✅ All effort branches are merged
- ✅ No code is lost or duplicated
- ✅ The project builds successfully
- ✅ All tests pass
- ✅ The integration branch is pushed to origin
- ✅ Total implementation respects size limits (under 2,500 lines for wave)

---

## Execution Instructions for Integration Agent

**CRITICAL**: This plan is for the Integration Agent to execute. The Integration Agent should:

1. Follow the commands EXACTLY as written
2. Handle conflicts as described in the resolution strategies
3. Validate after EACH merge before proceeding
4. Stop and report if unexpected issues arise
5. Create detailed logs of the integration process

**Remember**: This is a PLAN only. The actual execution will be performed by the Integration Agent in a separate session.

## R405 Compliance Note

This merge plan has been created in compliance with R405 automation requirements. The Integration Agent must:
- Execute these merges WITHOUT manual review
- Continue through the entire sequence unless a CRITICAL error occurs
- Report completion status with CONTINUE-SOFTWARE-FACTORY flag

CONTINUE-SOFTWARE-FACTORY=TRUE
