# 🚨🚨🚨 BLOCKING RULE R296: Deprecated Branch Marking Protocol

## Rule Definition
When an effort branch exceeds size limits and is successfully split into smaller branches, the original oversized branch MUST be marked as deprecated to prevent accidental integration.

## Criticality: 🚨🚨🚨 BLOCKING
Integration of deprecated branches causes critical failures and corrupted merges.

## Requirements

### 1. Branch Renaming After Successful Splits
When ALL splits for an effort are completed and verified:
```bash
# Original branch that was split
OLD_BRANCH="myproject/phase1/wave1/effort1"

# Rename to mark as deprecated
NEW_BRANCH="${OLD_BRANCH}-deprecated-split"

# Execute renaming
git branch -m "$OLD_BRANCH" "$NEW_BRANCH"
git push origin ":$OLD_BRANCH"  # Delete old remote
git push origin "$NEW_BRANCH"   # Push renamed branch
```

### 2. State File Update Requirements
The orchestrator-state.json MUST be updated to track deprecated branches:
```yaml
efforts_completed:
  effort-001-feature:
    status: "SPLIT_DEPRECATED"
    deprecated_branch: "myproject/phase1/wave1/effort1-deprecated-split"
    replacement_splits:
      - "myproject/phase1/wave1/effort1-split1"
      - "myproject/phase1/wave1/effort1-split2"
      - "myproject/phase1/wave1/effort1-split3"
    do_not_integrate: true
    split_completed_at: "2025-01-20T14:30:00Z"
```

### 3. Integration Prevention Check
BEFORE any integration operation, validate no deprecated branches:
```bash
# Check for deprecated suffix
for branch in "${BRANCHES_TO_MERGE[@]}"; do
    if [[ "$branch" == *"-deprecated-split" ]]; then
        echo "❌ BLOCKED: Cannot integrate deprecated branch: $branch"
        echo "Use replacement splits instead!"
        exit 1
    fi
done

# Check state file for SPLIT_DEPRECATED status
for effort in "${EFFORTS[@]}"; do
    STATUS=$(jq ".efforts_completed.\"$effort\".status" orchestrator-state.json)
    if [[ "$STATUS" == "SPLIT_DEPRECATED" ]]; then
        echo "❌ BLOCKED: $effort was split and deprecated"
        SPLITS=$(jq ".efforts_completed.\"$effort\".replacement_splits[]" orchestrator-state.json)
        echo "Use these splits instead: $SPLITS"
        exit 1
    fi
done
```

## When This Rule Applies

### MUST Apply When:
1. SW Engineer completes all splits for an effort
2. Orchestrator marks effort splits as complete
3. Integration agent starts merging efforts
4. Phase integration begins
5. Final project integration occurs

### Triggered By States:
- SW Engineer: SPLIT_IMPLEMENTATION (completion)
- Orchestrator: MONITORING_SPLITS (when all complete)
- Orchestrator: INTEGRATION (pre-check)
- Orchestrator: PHASE_INTEGRATION (pre-check)
- Orchestrator: FINAL_INTEGRATION (pre-check)
- Integration Agent: All states (pre-check)

## Implementation Requirements

### For Orchestrator:
1. After ALL splits complete for an effort:
   - Rename original branch with `-deprecated-split` suffix
   - Update state file with SPLIT_DEPRECATED status
   - List all replacement splits
   - Set do_not_integrate flag

2. Before ANY integration:
   - Check all branches for deprecated suffix
   - Verify state file for SPLIT_DEPRECATED status
   - Block if deprecated branch detected
   - Guide to use replacement splits

### For SW Engineer:
1. After completing final split:
   - Report to orchestrator that splits are complete
   - Provide list of all split branches created
   - Confirm original branch can be deprecated

### For Integration Agent:
1. Before EVERY merge operation:
   - Run deprecated branch check
   - Verify against state file
   - Abort if deprecated branch detected
   - Use only replacement splits

## Error Messages

### When Deprecated Branch Detected:
```
❌ CRITICAL: Attempted to integrate deprecated branch!
Branch: myproject/phase1/wave1/effort1-deprecated-split
Status: SPLIT_DEPRECATED - This branch was split due to size violations

Replacement branches to use instead:
  1. myproject/phase1/wave1/effort1-split1 (250 lines)
  2. myproject/phase1/wave1/effort1-split2 (300 lines)
  3. myproject/phase1/wave1/effort1-split3 (150 lines)

Action Required: Update integration list to use replacement splits
```

## Verification Steps

### Check for Deprecated Branches:
```bash
# List all deprecated branches
git branch -r | grep "deprecated-split"

# Check state file for deprecated efforts
jq '.efforts_completed | to_entries | .[] | select(.value.status == "SPLIT_DEPRECATED") | .key' orchestrator-state.json
```

### Validate Integration List:
```bash
# Ensure no deprecated branches in integration list
for branch in "${INTEGRATION_BRANCHES[@]}"; do
    if [[ "$branch" == *"deprecated-split" ]]; then
        echo "ERROR: Deprecated branch in integration list: $branch"
        exit 1
    fi
done
```

## Related Rules
- R204: Orchestrator Split Infrastructure
- R034: Integration Requirements
- R260: Integration Agent Core Requirements
- R282: Phase Integration Protocol
- R283: Project Integration Protocol

## Penalties
- Integrating deprecated branch: -50% (Critical failure)
- Not marking split branches as deprecated: -20%
- Missing state file updates: -15%
- No pre-integration checks: -25%

---
*Rule Type*: Protocol
*Agents*: Orchestrator, SW Engineer, Integration Agent
*Enforcement*: Automated via pre-checks and state validation