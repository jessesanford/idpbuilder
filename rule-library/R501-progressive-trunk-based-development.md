# 🔴🔴🔴 SUPREME RULE R501: Progressive Trunk-Based Development with Final Merge Plan

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
Software Factory 2.0 uses PROGRESSIVE trunk-based development where branches CASCADE sequentially. Each effort branches from the PREVIOUS effort (not from main), creating a chain where each branch contains all previous changes. This enables clean, sequential merging back to main without conflicts. The orchestrator-state-v3.json MUST track the final_merge_plan showing the exact sequence for merging all branches.

## 🔴🔴🔴 THE PROGRESSIVE CASCADE LAW 🔴🔴🔴

**NO PARALLEL BRANCHING FROM MAIN - ONLY CASCADE BRANCHING!**

### The Progressive Principle:
```
Traditional (WRONG):              Progressive (CORRECT):
main                              main
 ├── effort-1                      └── effort-1
 ├── effort-2  ❌                      └── effort-2
 ├── effort-3  ❌                          └── effort-3
 └── effort-4  ❌                              └── effort-4
                                                   └── effort-5
```

## 🔴 MANDATORY FINAL_MERGE_PLAN TRACKING 🔴

### Required State File Structure:
Every orchestrator-state-v3.json MUST contain:

```json
{
  "final_merge_plan": {
    "description": "Sequential merge order for progressive trunk-based development",
    "last_updated": "2025-01-25T10:00:00Z",
    "merge_sequence": [
      {
        "order": 1,
        "branch": "project/phase1/wave1/command-tests",
        "base_branch": "main",
        "contains_efforts": ["command-tests"],
        "cumulative_changes": ["command-tests"],
        "merge_status": "pending",
        "merged_at": null,
        "merge_commit": null
      },
      {
        "order": 2,
        "branch": "project/phase1/wave1/command-skeleton",
        "base_branch": "project/phase1/wave1/command-tests",
        "contains_efforts": ["command-skeleton"],
        "cumulative_changes": ["command-tests", "command-skeleton"],
        "merge_status": "pending",
        "merged_at": null,
        "merge_commit": null
      },
      {
        "order": 3,
        "branch": "project/phase1/wave1/integration-tests",
        "base_branch": "project/phase1/wave1/command-skeleton",
        "contains_efforts": ["integration-tests"],
        "cumulative_changes": ["command-tests", "command-skeleton", "integration-tests"],
        "merge_status": "pending",
        "merged_at": null,
        "merge_commit": null
      }
    ],
    "merge_rules": {
      "sequential_only": true,
      "skip_integration_branches": true,
      "require_all_tests_pass": true,
      "allow_force_push": false
    },
    "next_merge_index": 1,
    "total_branches": 3,
    "branches_merged": 0,
    "estimated_completion": null
  }
}
```

## 🔴 CASCADE BRANCHING REQUIREMENTS 🔴

### 1. Branch Creation MUST Follow Cascade:
```bash
# ORCHESTRATOR creates branches in CASCADE
create_cascaded_branch() {
    local EFFORT=$1
    local INDEX=$2  # Position in cascade (1, 2, 3...)

    if [[ $INDEX -eq 1 && $PHASE -eq 1 && $WAVE -eq 1 ]]; then
        # ONLY the first effort of P1W1 branches from main
        BASE="main"
    else
        # ALL other efforts branch from previous effort
        BASE=$(get_previous_effort_branch $INDEX)
    fi

    # Create branch from cascade base
    git checkout -b "$EFFORT" "$BASE"

    # Update final_merge_plan
    add_to_final_merge_plan "$EFFORT" "$BASE" "$INDEX"
}
```

### 2. Final Merge Plan Updates:
```bash
# Add effort to merge plan when created
add_to_final_merge_plan() {
    local BRANCH=$1
    local BASE=$2
    local ORDER=$3

    # Calculate cumulative changes
    CUMULATIVE=$(jq -r --arg base "$BASE" '
        .final_merge_plan.merge_sequence[] |
        select(.branch == $base) |
        .cumulative_changes | @json
    ' orchestrator-state-v3.json)

    # Add this effort to cumulative
    jq --arg branch "$BRANCH" \
       --arg base "$BASE" \
       --arg order "$ORDER" \
       --arg cumulative "$CUMULATIVE" \
       '.final_merge_plan.merge_sequence += [{
          "order": ($order | tonumber),
          "branch": $branch,
          "base_branch": $base,
          "contains_efforts": [$branch],
          "cumulative_changes": ($cumulative | fromjson) + [$branch],
          "merge_status": "pending"
       }] |
       .final_merge_plan.total_branches += 1 |
       .final_merge_plan.last_updated = now | todate' \
    orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
}
```

### 3. Split Branches Also Cascade:
```bash
# Splits cascade from each other, NOT from the oversized branch
create_split_cascade() {
    local OVERSIZED=$1
    local SPLIT_NUM=$2

    if [[ $SPLIT_NUM -eq 1 ]]; then
        # First split uses same base as oversized effort
        BASE=$(get_effort_original_base "$OVERSIZED")
    else
        # Subsequent splits cascade from previous split
        PREV_SPLIT=$((SPLIT_NUM - 1))
        BASE="${OVERSIZED}-split-$(printf "%03d" $PREV_SPLIT)"
    fi

    SPLIT_BRANCH="${OVERSIZED}-split-$(printf "%03d" $SPLIT_NUM)"
    git checkout -b "$SPLIT_BRANCH" "$BASE"
}
```

## 🔴 VERIFICATION REQUIREMENTS 🔴

### 1. Pre-Spawn Verification:
```bash
# BEFORE spawning any agent
verify_cascade_integrity() {
    local EFFORT=$1
    local EXPECTED_BASE=$2

    # Check final_merge_plan has correct base
    RECORDED_BASE=$(jq -r --arg e "$EFFORT" '
        .final_merge_plan.merge_sequence[] |
        select(.branch | contains($e)) |
        .base_branch
    ' orchestrator-state-v3.json)

    if [ "$RECORDED_BASE" != "$EXPECTED_BASE" ]; then
        echo "❌ FATAL: Cascade violation! $EFFORT not based on $EXPECTED_BASE"
        echo "Recorded base: $RECORDED_BASE"
        exit 1
    fi
}
```

### 2. Agent Startup Verification:
```bash
# EVERY agent MUST verify cascade on startup
verify_im_in_cascade() {
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

    # Get my base from merge plan
    MY_BASE=$(jq -r --arg b "$CURRENT_BRANCH" '
        .final_merge_plan.merge_sequence[] |
        select(.branch == $b) |
        .base_branch
    ' /workspaces/software-factory-2.0/orchestrator-state-v3.json)

    # Verify I'm actually based on it
    if ! git merge-base --is-ancestor "$MY_BASE" HEAD; then
        echo "❌ FATAL: Not in cascade! Branch $CURRENT_BRANCH not based on $MY_BASE"
        exit 1
    fi
}
```

### 3. Merge Sequence Validation:
```bash
# Validate entire merge sequence is cascaded
validate_merge_sequence() {
    jq -r '.final_merge_plan.merge_sequence[] |
           "\(.order):\(.branch):\(.base_branch)"' \
    orchestrator-state-v3.json | while IFS=: read -r order branch base; do
        if [[ $order -eq 1 ]]; then
            # First must be from main
            if [ "$base" != "main" ]; then
                echo "❌ First branch not from main!"
                exit 1
            fi
        else
            # Others must cascade from previous
            PREV_ORDER=$((order - 1))
            PREV_BRANCH=$(jq -r --arg o "$PREV_ORDER" '
                .final_merge_plan.merge_sequence[] |
                select(.order == ($o | tonumber)) |
                .branch
            ' orchestrator-state-v3.json)

            if [ "$base" != "$PREV_BRANCH" ]; then
                echo "❌ Branch $branch not cascaded from previous ($PREV_BRANCH)!"
                exit 1
            fi
        fi
    done
}
```

## 🔴 MERGE EXECUTION PROTOCOL 🔴

### Sequential Merge Process:
```bash
# Execute merges in sequence per plan
execute_final_merge() {
    # Get next branch to merge
    NEXT_INDEX=$(jq '.final_merge_plan.next_merge_index' orchestrator-state-v3.json)

    MERGE_INFO=$(jq -r --arg idx "$NEXT_INDEX" '
        .final_merge_plan.merge_sequence[] |
        select(.order == ($idx | tonumber)) |
        "\(.branch):\(.base_branch)"
    ' orchestrator-state-v3.json)

    IFS=: read -r BRANCH BASE <<< "$MERGE_INFO"

    # Checkout main and merge
    git checkout main
    git merge --no-ff "$BRANCH" -m "merge: $BRANCH (order: $NEXT_INDEX)"

    # Update merge plan
    jq --arg idx "$NEXT_INDEX" \
       --arg commit "$(git rev-parse HEAD)" \
       '.final_merge_plan.merge_sequence |= map(
          if .order == ($idx | tonumber) then
            .merge_status = "completed" |
            .merged_at = now | todate |
            .merge_commit = $commit
          else . end
       ) |
       .final_merge_plan.next_merge_index += 1 |
       .final_merge_plan.branches_merged += 1' \
    orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
}
```

## 🚨 COMMON VIOLATIONS (AUTOMATIC FAILURE) 🚨

### ❌ VIOLATION 1: Multiple Branches from Main
```bash
# WRONG - All branching from main
git checkout -b effort-1 main
git checkout -b effort-2 main  # ❌ Should be from effort-1!
git checkout -b effort-3 main  # ❌ Should be from effort-2!
```

### ✅ CORRECT: Cascade Branching
```bash
# RIGHT - Progressive cascade
git checkout -b effort-1 main
git checkout -b effort-2 effort-1
git checkout -b effort-3 effort-2
```

### ❌ VIOLATION 2: Missing Final Merge Plan
```json
{
  "orchestrator-state": {
    "efforts_in_progress": [...],
    // ❌ No final_merge_plan section!
  }
}
```

### ❌ VIOLATION 3: Integration Branches in Merge Plan
```json
{
  "final_merge_plan": {
    "merge_sequence": [
      {"branch": "effort-1", ...},
      {"branch": "phase1-wave1-integration", ...}  // ❌ NO!
    ]
  }
}
```

## 🔴 INTEGRATE_WAVE_EFFORTS WITH OTHER RULES 🔴

### R308 (Incremental Branching)
- R308 defines the CASCADE strategy
- R501 enforces the merge plan tracking
- Both work together for progressive development

### R337 (State File Truth)
- R337 requires state file as truth
- R501 adds final_merge_plan to state
- Plan becomes part of single source

### R196 (Base Branch Selection)
- R196 reads bases from state
- R501 ensures bases follow cascade
- Both prevent parallel branching

## 🔴 GRADING IMPACT 🔴

- **Parallel branching from main**: -100% (CASCADE VIOLATION)
- **Missing final_merge_plan**: -100% (Tracking failure)
- **Wrong base in cascade**: -100% (Progressive violation)
- **Integration branches in plan**: -50% (Plan contamination)
- **Out-of-order merging**: -75% (Sequence violation)

## 🔴 WHY THIS MATTERS 🔴

### Without Progressive Cascade:
- **Merge Conflicts**: Every branch conflicts with others
- **Integration Hell**: Big-bang integration nightmares
- **Lost Changes**: Commits overwritten during merges
- **Broken Builds**: Incompatible changes merged together

### With Progressive Cascade:
- **No Conflicts**: Each merge assumes previous already merged
- **Smooth Integration**: Sequential, predictable merging
- **Preserved History**: All changes accumulate properly
- **Always Buildable**: Each merge maintains working state

## 🔴 THE FINAL TRUTH 🔴

**REMEMBER: Progressive trunk-based development means EVERY branch builds on the previous one, creating a cascade of changes that merge back cleanly and sequentially!**

- First effort: From main
- All others: From previous effort
- Splits: From previous split
- Merges: In exact sequence
- No skipping, no parallel, ONLY CASCADE!

**This is the FOUNDATION of Software Factory 2.0's development model!**