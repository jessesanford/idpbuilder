# Phase 2 Wave 1 Rebase-Merge Plan

**Created**: 2025-09-13T20:40:28Z
**Code Reviewer Agent**: code-reviewer
**State**: EFFORT_PLAN_CREATION (Rebase Planning)
**Rule Compliance**: R327 (Mandatory Re-integration After Fixes), R307 (Independent Branch Mergeability)

## 🔴 CRITICAL CONTEXT

Phase 2 Wave 1 efforts were implemented BEFORE Phase 1 integration was finalized. Per R327, these efforts MUST be rebased onto the current Phase 1 integration branch to ensure they build on the latest integrated foundation.

## 📊 Current State Analysis

### Phase 1 Integration Status
- **Integration Branch**: `idpbuilder-oci-build-push/phase1/wave2-integration`
- **Latest Commit**: fba8c88 (initialization of wave integration infrastructure)
- **Key Additions Since Phase 2 Base**:
  - Wave integration infrastructure (R308 compliance)
  - Integration validation checklist
  - Complete Phase 1 Wave 1 integration documentation
  - Kind Certificate Extraction (E1.1.1) integration
  - Merge conflict resolutions from E1.1.2

### Phase 2 Wave 1 Efforts Requiring Rebase
1. **gitea-client** (main effort)
   - Current Branch: `idpbuilder-oci-build-push/phase2/wave1/gitea-client`
   - Current Base: 5344541 (before Phase 1 integration)
   - Status: Implementation complete with demos

2. **gitea-client-split-001** (sequential dependency)
   - Current Branch: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001`
   - Depends On: gitea-client main effort
   - Status: Implementation complete

3. **gitea-client-split-002** (sequential dependency)
   - Current Branch: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002`
   - Depends On: gitea-client-split-001
   - Status: Implementation complete

## 🔄 REBASE SEQUENCE STRATEGY

### Critical Ordering Requirements (R327 Cascade Pattern)
The splits have SEQUENTIAL dependencies, requiring careful ordering:
1. **FIRST**: Rebase main effort (gitea-client)
2. **SECOND**: Rebase split-001 onto rebased main
3. **THIRD**: Rebase split-002 onto rebased split-001

### Phase 1: Rebase Main Effort (gitea-client)

```bash
# Step 1: Navigate to main effort directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client

# Step 2: Ensure clean working state
git status
# CRITICAL: If uncommitted changes exist, stash or commit them first
git stash push -m "Pre-rebase checkpoint"

# Step 3: Fetch latest integration branch
git fetch origin idpbuilder-oci-build-push/phase1/wave2-integration:refs/remotes/origin/idpbuilder-oci-build-push/phase1/wave2-integration

# Step 4: Create backup branch before rebase
git branch backup-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Step 5: Perform interactive rebase
git rebase origin/idpbuilder-oci-build-push/phase1/wave2-integration

# Step 6: Handle conflicts if they arise (see Conflict Resolution section)

# Step 7: Verify successful rebase
git log --oneline -5
# Should show Phase 1 integration commits as base

# Step 8: Force push rebased branch
git push --force-with-lease origin idpbuilder-oci-build-push/phase2/wave1/gitea-client
```

### Phase 2: Rebase Split-001

```bash
# Step 1: Navigate to split-001 directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-001

# Step 2: Fetch rebased main effort
git fetch origin idpbuilder-oci-build-push/phase2/wave1/gitea-client:refs/remotes/origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client

# Step 3: Create backup
git branch backup-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Step 4: Rebase onto rebased main effort
git rebase origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client

# Step 5: Handle conflicts if needed

# Step 6: Push rebased split-001
git push --force-with-lease origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
```

### Phase 3: Rebase Split-002

```bash
# Step 1: Navigate to split-002 directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-002

# Step 2: Fetch rebased split-001
git fetch origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001:refs/remotes/origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

# Step 3: Create backup
git branch backup-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Step 4: Rebase onto rebased split-001
git rebase origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

# Step 5: Handle conflicts if needed

# Step 6: Push rebased split-002
git push --force-with-lease origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
```

## ⚠️ CONFLICT RESOLUTION STRATEGY

### Potential Conflict Areas
Based on Phase 1 integration changes, conflicts may occur in:
1. **go.mod/go.sum**: Dependency version conflicts
2. **pkg/**: If Phase 1 added new packages or interfaces
3. **cmd/**: If command structure changed
4. **Integration points**: Where Phase 2 interfaces with Phase 1 code

### Resolution Approach
```bash
# When conflicts occur during rebase:

# 1. Check conflict status
git status

# 2. For each conflicted file:
# Open file and look for conflict markers
<<<<<<< HEAD
# Phase 1 integration version
=======
# Your Phase 2 changes
>>>>>>> commit-hash

# 3. Resolution guidelines:
# - Preserve Phase 1 integration changes (base)
# - Apply Phase 2 changes on top
# - Ensure no functionality is lost
# - Test compilation after each resolution

# 4. After resolving a file:
git add <resolved-file>

# 5. Continue rebase
git rebase --continue

# 6. If rebase becomes too complex:
git rebase --abort
# Consider alternative merge strategy
```

### Conflict Resolution Principles
1. **Preserve Phase 1 Integration**: Phase 1 changes take precedence
2. **Maintain Phase 2 Functionality**: Ensure all Phase 2 features remain
3. **Test Incrementally**: Compile and test after each conflict resolution
4. **Document Changes**: Note any significant resolution decisions

## ✅ POST-REBASE VALIDATION

### Build Verification
```bash
# For each rebased effort:

# 1. Clean build artifacts
make clean || go clean -cache

# 2. Verify dependencies
go mod tidy
go mod verify

# 3. Build the project
go build ./...

# 4. Run unit tests
go test ./...

# 5. Run integration tests (if applicable)
make test-integration || go test -tags=integration ./...

# 6. Verify demos still work (R291 compliance)
./run-demos.sh || make demos
```

### Integration Points Verification
```bash
# Verify Phase 1 integration points:

# 1. Check for Phase 1 interfaces implementation
grep -r "phase1" pkg/

# 2. Verify import statements
go list -m all | grep idpbuilder

# 3. Check for API compatibility
# Run any API validation tests

# 4. Verify feature flags (if any)
# Check configuration files
```

### Size Verification (R304 Compliance)
```bash
# CRITICAL: Verify size limits after rebase

# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# For each effort, run line counter
cd /path/to/effort
$PROJECT_ROOT/tools/line-counter.sh

# Ensure still under 800 line limit
```

## 📋 EFFORT PLANNING READINESS

### Prerequisites Confirmed
After successful rebase, each effort will have:
- ✅ Clean rebase onto Phase 1 integration branch
- ✅ All tests passing
- ✅ Build successful
- ✅ Size limits maintained
- ✅ Demos functional (R291)
- ✅ Ready for parallel implementation planning

### Branch Tracking Updates
```bash
# Update each effort's tracking to new base:

# For main effort:
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client
git branch --set-upstream-to=origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client

# For split-001:
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-001
git branch --set-upstream-to=origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

# For split-002:
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-002
git branch --set-upstream-to=origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
```

## 🎯 NEXT STEPS FOR ORCHESTRATOR

1. **Execute Rebase Sequence**: Follow the sequential rebase process above
2. **Validate Each Step**: Ensure builds and tests pass after each rebase
3. **Handle Conflicts**: Use the conflict resolution strategy if needed
4. **Verify Readiness**: Confirm all validation checks pass
5. **Spawn Parallel Code Reviewers**: Once rebased, spawn parallel Code Reviewers for effort planning:
   - One for gitea-client main effort
   - One for gitea-client-split-001
   - One for gitea-client-split-002

## 🚨 CRITICAL REMINDERS

1. **R327 Compliance**: This rebase is MANDATORY per cascade pattern
2. **Sequential Dependencies**: Splits MUST be rebased in order
3. **Backup Branches**: Always create backups before rebasing
4. **Force Push Safety**: Use `--force-with-lease` not `--force`
5. **Test Everything**: Validate after each rebase step

## 📊 Success Criteria

The rebase is successful when:
- ✅ All three efforts rebased onto Phase 1 integration
- ✅ Sequential dependencies maintained (splits build on each other)
- ✅ All builds passing
- ✅ All tests passing
- ✅ Demos functional
- ✅ Size limits maintained
- ✅ Ready for parallel implementation planning

---

**End of Rebase-Merge Plan**

*Created by Code Reviewer Agent in EFFORT_PLAN_CREATION state*
*Compliance: R327, R307, R308, R304*