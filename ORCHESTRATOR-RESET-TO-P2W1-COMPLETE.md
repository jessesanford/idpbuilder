# 🔴🔴🔴 ORCHESTRATOR RESET PROMPT - PHASE 2 WAVE 1 COMPLETE 🔴🔴🔴

## Purpose
This prompt will reset the orchestrator to the END of Phase 2 Wave 1, with all efforts completed and ready for wave integration. This allows testing of the mandatory wave integration flow per R009 (formerly R336).

## Current Situation
- Phase 2 Wave 1: All efforts COMPLETE and REVIEWED
- Phase 2 Wave 2: NEVER HAPPENED (will be deleted)
- Next Step: MANDATORY wave integration before ANY new work

## Reset Instructions

### 1. FIRST: Save Current State (Optional)
```bash
# If you want to preserve current state for reference
cp orchestrator-state.json orchestrator-state-backup-$(date +%Y%m%d-%H%M%S).json
```

### 2. DELETE ALL PHASE 2 WAVE 2 WORK

#### A. Remove Effort Directories
```bash
# Delete Wave 2 effort directories
rm -rf efforts/phase2/wave2/
rm -rf /efforts/phase2/wave2/

# Delete any Wave 2 integration attempts
rm -rf efforts/phase2/wave2-integration/
rm -rf /efforts/phase2/wave2-integration/
```

#### B. Delete Local Branches
```bash
# Delete all Phase 2 Wave 2 local branches
for branch in $(git branch | grep -E "phase2.*wave2|P2W2|p2w2"); do
    git branch -D "$branch" 2>/dev/null || true
done

# Also check for any E2.2.* branches (Effort 2.2.*)
for branch in $(git branch | grep -E "E2\.2\.|effort.*2\.2"); do
    git branch -D "$branch" 2>/dev/null || true
done
```

#### C. Delete Remote Branches
```bash
# Delete all Phase 2 Wave 2 remote branches
for branch in $(git branch -r | grep -E "phase2.*wave2|P2W2|p2w2" | sed 's/origin\///'); do
    git push origin --delete "$branch" 2>/dev/null || true
done

# Also delete any E2.2.* remote branches
for branch in $(git branch -r | grep -E "E2\.2\.|effort.*2\.2" | sed 's/origin\///'); do
    git push origin --delete "$branch" 2>/dev/null || true
done
```

#### D. Remove Planning/Review Artifacts
```bash
# Remove Wave 2 planning documents
find . -name "*phase2*wave2*.md" -type f -delete
find . -name "*P2W2*.md" -type f -delete
find . -name "*E2.2*.md" -type f -delete

# Remove from phase-plans directory
rm -f phase-plans/*phase2-wave2*
rm -f phase-plans/*P2W2*
```

### 3. RESET ORCHESTRATOR STATE FILE

Create this exact state file:

```json
{
  "current_state": "WAVE_COMPLETE",
  "current_phase": 2,
  "current_wave": 1,
  "project_name": "Your Project Name",
  "base_branch": "main",
  
  "phases": [
    {
      "phase_number": 1,
      "status": "COMPLETE",
      "total_waves": 3,
      "integration_branch": "phase1-integration"
    },
    {
      "phase_number": 2,
      "status": "IN_PROGRESS",
      "total_waves": 2,
      "waves": {
        "wave1": {
          "status": "COMPLETE",
          "efforts": [
            {
              "id": "E2.1.1",
              "status": "COMPLETE",
              "branch": "phase2-wave1-effort1",
              "review_status": "APPROVED"
            },
            {
              "id": "E2.1.2", 
              "status": "COMPLETE",
              "branch": "phase2-wave1-effort2",
              "review_status": "APPROVED"
            },
            {
              "id": "E2.1.3",
              "status": "COMPLETE",
              "branch": "phase2-wave1-effort3",
              "review_status": "APPROVED"
            }
          ]
        }
      }
    }
  ],
  
  "efforts_completed": [
    "E1.1.1", "E1.1.2", "E1.1.3",
    "E1.2.1", "E1.2.2", "E1.2.3",
    "E1.3.1", "E1.3.2",
    "E2.1.1", "E2.1.2", "E2.1.3"
  ],
  
  "efforts_in_progress": [],
  
  "integrations": {
    "phase1": {
      "wave1": "phase1-wave1-integration",
      "wave2": "phase1-wave2-integration",
      "wave3": "phase1-wave3-integration",
      "phase": "phase1-integration"
    },
    "phase2": {
      "wave1": null
    }
  },
  
  "last_checkpoint": "P2W1_COMPLETE",
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
```

### 4. VERIFY RESET SUCCESS

```bash
# Check no Wave 2 branches exist
echo "=== Checking for Wave 2 branches ==="
git branch -a | grep -E "wave2|W2|2\.2" || echo "✅ No Wave 2 branches found"

# Check no Wave 2 directories exist
echo "=== Checking for Wave 2 directories ==="
find . -type d -name "*wave2*" 2>/dev/null || echo "✅ No Wave 2 directories found"

# Verify state file
echo "=== Verifying state file ==="
jq '.current_state' orchestrator-state.json  # Should be "WAVE_COMPLETE"
jq '.current_wave' orchestrator-state.json   # Should be 1
jq '.efforts_in_progress | length' orchestrator-state.json  # Should be 0

# Check Phase 2 Wave 1 integration doesn't exist yet
echo "=== Checking integration status ==="
if git ls-remote --heads origin "phase2-wave1-integration" > /dev/null 2>&1; then
    echo "⚠️ WARNING: Wave 1 integration already exists!"
    echo "You may want to delete it to test fresh integration:"
    echo "git push origin --delete phase2-wave1-integration"
else
    echo "✅ Wave 1 integration not created yet (correct for testing)"
fi
```

### 5. EXPECTED NEXT BEHAVIOR

When you run `/continue-orchestrating`, the orchestrator MUST:

1. **Recognize state**: WAVE_COMPLETE for Phase 2 Wave 1
2. **Next transition**: INTEGRATION (per R009 - MANDATORY!)
3. **NOT allowed**: 
   - ❌ Transition to WAVE_START (R009 violation!)
   - ❌ Transition to PLANNING (R009 violation!)
   - ❌ Start Wave 2 without integration (R009 violation!)
4. **Required sequence**:
   ```
   WAVE_COMPLETE 
   → INTEGRATION (setup infrastructure)
   → SPAWN_CODE_REVIEWER_MERGE_PLAN
   → SPAWN_INTEGRATION_AGENT
   → MONITORING_INTEGRATION
   → WAVE_REVIEW
   → [Only then] WAVE_START (for Wave 2)
   ```

### 6. TESTING THE INTEGRATION REQUIREMENT

After reset, test with:

```bash
# This MUST trigger integration, not skip to Wave 2
/continue-orchestrating
```

The orchestrator should:
- Create integration infrastructure for Phase 2 Wave 1
- Spawn code reviewer for merge plan
- Spawn integration agent to execute merges
- Create `phase2-wave1-integration` branch
- Only AFTER successful integration, allow Wave 2 to start

### 7. R009 ENFORCEMENT VERIFICATION

The following should cause IMMEDIATE FAILURE:

```bash
# If orchestrator tries this, it's a R009 violation:
current_state="WAVE_COMPLETE"
next_state="WAVE_START"  # ❌ WRONG! Must be INTEGRATION

# Or if Wave 2 efforts try to use wrong base:
base_branch="main"  # ❌ WRONG! Must be phase2-wave1-integration
```

## Summary

This reset creates a perfect test scenario for R009 (Mandatory Wave/Phase Integration Protocol):
- Phase 2 Wave 1 is complete but NOT integrated
- Wave 2 has been completely removed
- Orchestrator MUST create integration before proceeding
- Any attempt to skip integration = -100% FAILURE

## Commit After Reset

```bash
git add orchestrator-state.json
git commit -m "reset: orchestrator to P2W1 complete, ready for mandatory integration test"
git push
```

---
**Remember**: The purpose is to test that R009 is properly enforced. Integration is NOT optional!