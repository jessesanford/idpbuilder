# Phase 2 Wave 1 Rebase Merge Plan

## 🔴 CRITICAL CONTEXT
**Created**: 2025-09-13T21:35:00Z
**Reviewer**: Code Reviewer Agent (WAVE_MERGE_PLANNING state)
**Pattern**: R327 Mandatory Re-integration
**Reason**: Phase 2 Wave 1 efforts created BEFORE Phase 1 integration finalized

## 📊 REBASE SUMMARY

### Target Integration Branch
- **Branch**: `idpbuilder-oci-build-push/phase1/integration`
- **Status**: COMPLETE - Contains all Phase 1 efforts
- **Commit**: Latest from origin
- **Contents**:
  - E1.1.1 kind-cert-extraction (650 lines)
  - E1.1.2 registry-tls-trust (700 lines)
  - E1.1.3 registry-auth-types-split-001 (types/constants)
  - E1.1.3 registry-auth-types-split-002 (implementation)
  - E1.2.1 cert-validation-split-001 (foundations)
  - E1.2.1 cert-validation-split-002 (implementation)
  - E1.2.1 cert-validation-split-003 (completion)
  - E1.2.2 fallback-strategies (560 lines)

### Efforts Requiring Rebase (In Dependency Order)
1. **gitea-client** (main effort)
   - Current: `idpbuilder-oci-build-push/phase2/wave1/gitea-client`
   - Base: Outdated (pre-Phase 1 integration)
   - Commits ahead: 94
   - **Must rebase FIRST**

2. **gitea-client-split-001**
   - Current: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001`
   - Depends on: gitea-client (must wait for #1)
   - **Must rebase SECOND**

3. **gitea-client-split-002**
   - Current: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002`
   - Depends on: gitea-client-split-001 (must wait for #2)
   - **Must rebase THIRD**

## 🔧 REBASE EXECUTION STRATEGY

### PHASE 1: Rebase Main Effort (gitea-client)

```bash
# Step 1: Navigate to effort directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client

# Step 2: Ensure clean working state
git status
# MUST show: "nothing to commit, working tree clean"
# If not clean: git add -A && git commit -m "chore: save work before rebase"

# Step 3: Fetch latest integration branch
git fetch origin idpbuilder-oci-build-push/phase1/integration:refs/remotes/origin/idpbuilder-oci-build-push/phase1/integration

# Step 4: Create backup branch (safety)
git branch backup-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Step 5: Start interactive rebase
git rebase origin/idpbuilder-oci-build-push/phase1/integration

# Step 6: Handle conflicts if any
# For each conflict:
#   - Review the conflict markers
#   - Keep Phase 2 changes (ours) unless they conflict with Phase 1 APIs
#   - Ensure Phase 1 dependencies are properly imported
#   - git add <resolved-files>
#   - git rebase --continue

# Step 7: Verify compilation
go build ./... || make build

# Step 8: Run tests
go test ./... || make test

# Step 9: Push rebased branch
git push --force-with-lease origin idpbuilder-oci-build-push/phase2/wave1/gitea-client
```

### PHASE 2: Rebase Split-001 (AFTER Main Effort Complete)

```bash
# Step 1: Navigate to split-001 directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-001

# Step 2: Ensure clean state
git status

# Step 3: Fetch updated main effort branch
git fetch origin idpbuilder-oci-build-push/phase2/wave1/gitea-client:refs/remotes/origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client

# Step 4: Create backup
git branch backup-split-001-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Step 5: Rebase onto updated main effort
git rebase origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client

# Step 6: Handle conflicts
# Priority: Maintain split boundaries
# Do NOT duplicate code from main effort

# Step 7: Verify build and tests
go build ./... && go test ./...

# Step 8: Push rebased split
git push --force-with-lease origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
```

### PHASE 3: Rebase Split-002 (AFTER Split-001 Complete)

```bash
# Step 1: Navigate to split-002 directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client-split-002

# Step 2: Ensure clean state
git status

# Step 3: Fetch updated split-001 branch
git fetch origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001:refs/remotes/origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

# Step 4: Create backup
git branch backup-split-002-pre-rebase-$(date +%Y%m%d-%H%M%S)

# Step 5: Rebase onto updated split-001
git rebase origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

# Step 6: Handle conflicts
# Priority: Maintain sequential dependency chain
# Ensure all split-001 changes are preserved

# Step 7: Verify build and tests
go build ./... && go test ./...

# Step 8: Push rebased split
git push --force-with-lease origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
```

## ⚠️ CONFLICT RESOLUTION GUIDELINES

### Expected Conflict Areas
1. **go.mod/go.sum**: Phase 1 added new dependencies
   - Resolution: Keep BOTH sets of dependencies
   - Use `go mod tidy` after resolution

2. **Import statements**: Phase 1 created new packages
   - Resolution: Add Phase 1 imports where needed
   - Ensure proper import paths

3. **API changes**: Phase 1 may have modified interfaces
   - Resolution: Adapt Phase 2 code to use Phase 1 APIs
   - Maintain backward compatibility

### Conflict Resolution Priority
1. **Preserve Phase 1 integration**: All Phase 1 changes MUST be retained
2. **Adapt Phase 2 code**: Modify to work with Phase 1 changes
3. **Maintain split boundaries**: Don't mix code between splits
4. **Keep dependency order**: Splits must remain sequential

## ✅ POST-REBASE VERIFICATION CHECKLIST

### For EACH Rebased Effort:

#### 1. Git State Verification
- [ ] Branch rebased onto correct base
- [ ] No uncommitted changes
- [ ] Push successful with --force-with-lease
- [ ] Backup branch retained

#### 2. Code Compilation
- [ ] `go build ./...` succeeds
- [ ] No compilation errors
- [ ] All imports resolve

#### 3. Test Execution
- [ ] Unit tests pass
- [ ] Integration tests pass (if any)
- [ ] Demo runs successfully (R291 compliance)

#### 4. Dependency Verification
- [ ] Phase 1 packages accessible
- [ ] go.mod includes Phase 1 dependencies
- [ ] No duplicate dependencies

#### 5. Split Integrity (for splits only)
- [ ] Split maintains its boundaries
- [ ] No code duplication with other splits
- [ ] Dependencies on previous splits intact

#### 6. Size Compliance
- [ ] Run line counter to verify size
- [ ] Confirm still within limits (<800 lines)

### Final Integration Verification
```bash
# After all rebases complete, verify chain integrity
cd /home/vscode/workspaces/idpbuilder-oci-build-push

# Check main effort
cd efforts/phase2/wave1/gitea-client
git log --oneline -5
git merge-base HEAD origin/idpbuilder-oci-build-push/phase1/integration
# Should show Phase 1 integration as base

# Check split-001
cd ../gitea-client-split-001
git log --oneline -5
git merge-base HEAD origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client
# Should show rebased gitea-client as base

# Check split-002
cd ../gitea-client-split-002
git log --oneline -5
git merge-base HEAD origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
# Should show rebased split-001 as base
```

## 🚨 CRITICAL WARNINGS

1. **NEVER skip dependency order**: Splits MUST be rebased sequentially
2. **ALWAYS create backups**: Use timestamped backup branches
3. **FORCE PUSH carefully**: Use --force-with-lease, not --force
4. **TEST after each rebase**: Don't proceed if tests fail
5. **PRESERVE Phase 1 work**: All Phase 1 changes must be retained

## 📋 EXECUTION TRACKING

SW Engineers executing this plan should track progress:

| Effort | Start Time | Backup Branch | Rebase Status | Tests Pass | Push Complete | Verified |
|--------|------------|---------------|---------------|------------|---------------|----------|
| gitea-client | | | | | | |
| gitea-client-split-001 | | | | | | |
| gitea-client-split-002 | | | | | | |

## 🎯 SUCCESS CRITERIA

This rebase is complete when:
1. ✅ All three efforts rebased in correct order
2. ✅ All efforts base on Phase 1 integration (directly or transitively)
3. ✅ All code compiles successfully
4. ✅ All tests pass
5. ✅ All branches pushed to origin
6. ✅ Verification checklist complete for each effort
7. ✅ No Phase 1 functionality broken

## 📝 NOTES FOR SW ENGINEERS

- This is a MANDATORY rebase per R327
- Must be completed BEFORE any new implementation
- If major conflicts arise, consult with Code Reviewer
- Document any significant conflict resolutions
- Update orchestrator-state.json after completion

---
**Plan Created By**: Code Reviewer Agent
**State**: WAVE_MERGE_PLANNING
**Follows**: R327 Mandatory Re-integration Pattern
**Next Step**: SW Engineers execute this plan in order