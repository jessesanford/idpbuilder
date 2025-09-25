# Wave Merge Plan - Phase 2 Wave 2

**Generated**: 2025-09-24 22:08:00 UTC
**Reviewer**: Code Reviewer Agent (state: WAVE_MERGE_PLANNING)
**Target Branch**: `idpbuilderpush/phase2/wave2/integration`
**Base Branch**: `idpbuilderpush/phase2/wave1/integration`

## R269/R270 Compliance Notice
Per Supreme Laws R269 and R270:
- This is a MERGE PLAN ONLY - NO merges will be executed by this agent
- Using ONLY original effort branches - NO integration branches as sources
- All merges will be performed by the Orchestrator or designated merge agent

## Current Integration State

The integration branch (`idpbuilderpush/phase2/wave2/integration`) needs to merge:
1. ❌ **Effort 2.2.2 (auth-flow)** - Branch: `idpbuilderpush/phase2/wave2/auth-flow` (now in correct repo)
2. ❌ **Effort 2.2.3 (push-command)** - Branch: `idpbuilderpush/phase2/wave2/push-command`

## Branches to Merge

### 1. Auth Flow Implementation
- **Branch**: `idpbuilderpush/phase2/wave2/auth-flow`
- **Size**: 151 lines
- **Files Modified**:
  - `pkg/oci/flow.go`
  - `pkg/oci/types.go`
- **Status**: Ready to merge (now in correct repo)

### 2. Push Command Implementation
- **Branch**: `idpbuilderpush/phase2/wave2/push-command`
- **Size**: ~790 lines (per git diff --stat)
- **Files Modified**:
  - `pkg/cmd/push/push.go` (346 lines)
  - `pkg/cmd/push/push_test.go` (383 lines)
  - `pkg/cmd/root.go` (2 lines)
  - `FIX_COMPLETE.marker` (59 lines)
- **Status**: Ready to merge

## Merge Order and Instructions

Both auth-flow and push-command need to be merged:

### Step 1: Prepare Integration Branch
```bash
# Navigate to integration workspace
cd /home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave2/integration-workspace

# Ensure we're on the integration branch
git checkout idpbuilderpush/phase2/wave2/integration

# Verify current state
git status
git log --oneline -5

# Pull latest changes
git pull origin idpbuilderpush/phase2/wave2/integration
```

### Step 2: Merge Auth Flow Implementation
```bash
# Fetch the branch to ensure we have latest
git fetch origin idpbuilderpush/phase2/wave2/auth-flow:idpbuilderpush/phase2/wave2/auth-flow

# Perform the merge
git merge idpbuilderpush/phase2/wave2/auth-flow \
  -m "integrate: auth-flow from effort 2.2.2 (Phase 2 Wave 2)"

# If conflicts occur, follow conflict resolution below
```

### Step 3: Merge Push Command Implementation
```bash
# Fetch the branch to ensure we have latest
git fetch origin idpbuilderpush/phase2/wave2/push-command:idpbuilderpush/phase2/wave2/push-command

# Perform the merge
git merge idpbuilderpush/phase2/wave2/push-command \
  -m "integrate: push-command from effort 2.2.3 (Phase 2 Wave 2)"

# If conflicts occur, follow conflict resolution below
```

## Expected Conflicts

Based on analysis of the changes:

### Potential Conflict in `pkg/cmd/root.go`
- **Reason**: Both auth-flow and push-command may modify the root command
- **Resolution Strategy**:
  ```go
  // Accept both changes - ensure both commands are registered
  // The file should include both:
  // - Authentication-related commands from auth-flow
  // - Push command from push-command effort
  ```

### No other conflicts expected
- The push command implementation is in new files (`pkg/cmd/push/`)
- Test files are also new and shouldn't conflict

## Conflict Resolution Process

If conflicts occur:

1. **Identify conflicts**:
   ```bash
   git status
   git diff --name-only --diff-filter=U
   ```

2. **For each conflicted file**:
   ```bash
   # Open the file and look for conflict markers
   # <<<<<<< HEAD
   # (current integration branch content)
   # =======
   # (incoming push-command content)
   # >>>>>>> idpbuilderpush/phase2/wave2/push-command

   # Resolution principle: Preserve ALL functionality
   # - Keep all command registrations
   # - Keep all imports
   # - Ensure no duplicate function definitions
   ```

3. **After resolution**:
   ```bash
   git add <resolved-file>
   git commit --no-edit  # Use the merge commit message
   ```

## Post-Merge Validation

After the merge is complete, validate:

### 1. Build Verification
```bash
# Ensure the project builds
go build ./...

# Run tests
go test ./...
```

### 2. Feature Verification
```bash
# Verify push command is available
./idpbuilder push --help

# Verify auth functionality still works
# (Check auth-related commands from previous efforts)
```

### 3. Size Verification
```bash
# Navigate to project root
cd /home/vscode/workspaces/idpbuilder-push

# Run line counter on integration branch
./tools/line-counter.sh

# Verify total is reasonable (sum of all efforts)
# Expected: ~600 lines total (200 + 151 + 250)
```

## Integration Summary

| Effort | Branch | Status | Lines | Notes |
|--------|--------|--------|-------|-------|
| 2.2.2 Auth Flow | idpbuilderpush/phase2/wave2/auth-flow | 🔄 To Merge | 151 | Now in correct repo |
| 2.2.3 Push Command | idpbuilderpush/phase2/wave2/push-command | 🔄 To Merge | 250 | Ready for integration |

**Total Expected Lines**: ~401 lines (within 800 line limit per R220)

## Notes for Orchestrator

1. **Two Merges Required**: Both auth-flow and push-command branches need to be merged
2. **Low Conflict Risk**: New functionality with minimal overlap
3. **Branch Availability**: The push-command branch has been verified to exist and is ready
4. **Size Compliance**: Total size remains well within limits

## Rollback Plan

If issues occur during merge:

```bash
# Abort the merge if conflicts are complex
git merge --abort

# Or reset to pre-merge state if merge was completed
git reset --hard HEAD~1

# Investigate issues and retry with adjusted strategy
```

---

**END OF MERGE PLAN**

*This plan created per R269 - NO merges executed, only planned*
*Orchestrator must execute the actual merge operations*