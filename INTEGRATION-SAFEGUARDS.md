# Integration Safeguards - Preventing Wrong Branch Selection

## Problem Analysis

The Phase 2 integration failure occurred because:
1. **Wrong Plan Type**: A Wave-level merge plan was used for Phase-level integration
2. **Timing Issue**: Plan was created before Wave 2 integration was complete
3. **Branch Confusion**: Plan specified individual effort branches instead of integration branches
4. **No Verification**: Integration agent didn't verify it was using the right type of plan

## Safeguards Implemented

### 1. Clear Naming Convention
- **Wave Merge Plans**: `WAVE-X-Y-MERGE-PLAN.md` (for merging efforts INTO a wave)
- **Phase Merge Plans**: `PHASE-X-INTEGRATION-MERGE-PLAN.md` (for merging waves INTO a phase)
- **Project Merge Plans**: `PROJECT-INTEGRATION-MERGE-PLAN.md` (for merging phases)

### 2. Plan Validation Headers
Every merge plan MUST start with:
```markdown
# [LEVEL] MERGE PLAN
**Plan Type**: [WAVE|PHASE|PROJECT]
**Merging**: [What type of branches - efforts/waves/phases]
```

### 3. Branch Type Verification
Plans must explicitly state branch types:
```markdown
## Branch Types Being Merged
- [ ] Individual Effort Branches (for Wave integration only)
- [ ] Wave Integration Branches (for Phase integration only)
- [ ] Phase Integration Branches (for Project integration only)
```

### 4. Timing Requirements
```markdown
## Prerequisites Verification
- [ ] All source branches are complete
- [ ] All source branches have been reviewed
- [ ] Integration branches exist (if merging integrations)
- [ ] This plan was created AFTER prerequisites completed
```

### 5. Integration Agent Checklist

Before executing ANY merge plan, the Integration Agent MUST:

```bash
# 1. Verify plan matches integration level
if [[ "$INTEGRATION_LEVEL" == "PHASE" ]]; then
    if ! grep -q "Plan Type: PHASE" "$MERGE_PLAN"; then
        echo "ERROR: Wrong plan type for phase integration!"
        exit 1
    fi
fi

# 2. Verify branches exist
for branch in $(grep "git merge" "$MERGE_PLAN" | awk '{print $3}'); do
    if ! git rev-parse "$branch" >/dev/null 2>&1; then
        echo "ERROR: Branch $branch does not exist!"
        exit 1
    fi
done

# 3. Verify no individual efforts in phase integration
if [[ "$INTEGRATION_LEVEL" == "PHASE" ]]; then
    if grep -q "phase2/wave2/cli-commands\|credential-management\|image-operations" "$MERGE_PLAN"; then
        echo "ERROR: Phase integration cannot merge individual efforts!"
        exit 1
    fi
fi

# 4. Verify integration branches for phase/project level
if [[ "$INTEGRATION_LEVEL" == "PHASE" ]] || [[ "$INTEGRATION_LEVEL" == "PROJECT" ]]; then
    if ! grep -q "integration-[0-9]" "$MERGE_PLAN"; then
        echo "ERROR: Phase/Project integration must use integration branches!"
        exit 1
    fi
fi
```

## Rules for Creating Merge Plans

### For Code Reviewer Creating Plans:

1. **Wave Merge Plans** merge individual effort branches:
   - Input: effort branches (e.g., `phase2/wave2/cli-commands`)
   - Output: wave integration branch

2. **Phase Merge Plans** merge wave integration branches:
   - Input: wave integration branches (e.g., `phase2/wave1/integration-20250915-125755`)
   - Output: phase integration branch

3. **Project Merge Plans** merge phase integration branches:
   - Input: phase integration branches (e.g., `phase1/integration`, `phase2/integration`)
   - Output: project integration branch

### For Orchestrator:

1. **Request Correct Plan Type**:
   ```
   For Wave Integration: "Create WAVE merge plan for efforts"
   For Phase Integration: "Create PHASE merge plan for wave integrations"
   For Project Integration: "Create PROJECT merge plan for phase integrations"
   ```

2. **Provide Branch Lists**:
   - List exact branch names to be merged
   - Specify if they are efforts, wave integrations, or phase integrations

3. **Archive Old Plans**:
   - Move superseded plans to `archived-merge-plans/`
   - Never leave multiple plans that could confuse the integration agent

## Verification Commands

### For Integration Agent:
```bash
# Before starting integration, verify you have the right plan
echo "Integration Level: $INTEGRATION_LEVEL"
echo "Merge Plan: $MERGE_PLAN_FILE"
echo "Plan Type: $(grep "Plan Type:" "$MERGE_PLAN_FILE")"

# They must match:
# WAVE -> WAVE plan
# PHASE -> PHASE plan
# PROJECT -> PROJECT plan
```

### For Orchestrator After Integration:
```bash
# Verify correct branches were merged
if [[ "$INTEGRATION_LEVEL" == "PHASE" ]]; then
    # Check for credential flags (Phase 2 specific)
    grep "pushUsername\|pushToken" pkg/cmd/push.go || echo "FAILURE: Credentials missing!"

    # Check git log shows integration branches
    git log --oneline | grep "integration-" || echo "FAILURE: No integration branches in history!"
fi
```

## Emergency Recovery

If wrong branches are merged:
1. **DO NOT PROCEED** - Stop immediately
2. **DO NOT PUSH** - Keep the mistake local
3. **ROLLBACK**:
   ```bash
   git reset --hard HEAD~  # Undo last merge
   # OR
   git checkout main
   git branch -D [bad-branch]
   ```
4. **GET CORRECT PLAN** - Request new plan from Code Reviewer
5. **VERIFY** - Use safeguards above before retrying

## Summary

The key to preventing future failures:
1. ✅ Use correct plan type for integration level
2. ✅ Create plans AFTER prerequisites are complete
3. ✅ Verify branch types match integration level
4. ✅ Archive old/incorrect plans immediately
5. ✅ Integration agent must validate before executing

These safeguards ensure that:
- Wave integrations merge efforts
- Phase integrations merge waves
- Project integrations merge phases
- Individual efforts are NEVER merged directly into phase/project level