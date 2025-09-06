# Rule R202: Single SW Engineer Per Split - Sequential Execution Only

## Rule Statement
When an effort is split due to exceeding 800 lines, ALL splits from that effort MUST be implemented by a SINGLE SW engineering agent working sequentially. NEVER spawn multiple agents for splits from the same effort. The agent completes split 1, then split 2, then split 3, etc. in strict sequence.

## Criticality Level
**BLOCKING** - Multiple agents on splits cause conflicts and integration failures

## Enforcement Mechanism
- **Technical**: Single agent instance handles all splits sequentially
- **Behavioral**: Agent refuses if another agent working on sibling splits
- **Grading**: -40% for parallel split execution (Major architectural failure)

## Core Principle

```
1 Over-Limit Effort = 1 Code Reviewer (plans all splits) + 1 SW Engineer (implements all splits)
Splits MUST be done IN SEQUENCE, never in parallel
Split 1 → Complete → Split 2 → Complete → Split 3 → Complete
```

## Detailed Requirements

### ORCHESTRATOR: Sequential Split Execution

```bash
# ❌❌❌ WRONG - Multiple agents for splits (CAUSES CONFLICTS!)
handle_splits_wrong() {
    local effort="api-types"
    local splits_needed=3
    
    # WRONG! Spawning different agents for each split
    for split in 1 2 3; do
        Task: sw-engineer  # NEW agent each time = WRONG!
        Working directory: efforts/phase1/wave1/${effort}-split-${split}
        Implement split ${split}
    done
    # Result: Conflicts, integration issues, duplicated work!
}

# ✅✅✅ CORRECT - ONE agent for ALL splits, done sequentially
handle_splits_correct() {
    local effort="api-types"
    local splits_needed=3
    
    echo "═══════════════════════════════════════════════════════"
    echo "SPAWNING SINGLE AGENT FOR ALL SPLITS"
    echo "Effort: $effort"
    echo "Splits to implement: $splits_needed"
    echo "Execution: SEQUENTIAL (split-by-split)"
    echo "═══════════════════════════════════════════════════════"
    
    # Spawn ONE agent to handle ALL splits sequentially
    Task: sw-engineer
    Working directory: efforts/phase1/wave1/${effort}
    
    CRITICAL SPLIT IMPLEMENTATION INSTRUCTIONS:
    
    You are the ONLY agent implementing ALL splits for $effort.
    You MUST implement them SEQUENTIALLY, not in parallel.
    
    Your workflow:
    1. Read SPLIT-INVENTORY.md for complete split overview
    2. Start with SPLIT-PLAN-001.md
    3. Implement split 001 completely
    4. Commit and push split 001
    5. Verify split 001 is complete
    6. ONLY THEN move to SPLIT-PLAN-002.md
    7. Implement split 002 completely
    8. Continue sequentially until all $splits_needed splits done
    
    NEVER work on multiple splits simultaneously.
    Complete each split fully before starting the next.
    You have COMPLETE responsibility for ALL splits.
}
```

### SW ENGINEER: Sequential Split Implementation

```bash
# When assigned to implement splits
handle_split_implementation() {
    echo "═══════════════════════════════════════════════════════"
    echo "SPLIT IMPLEMENTATION ASSIGNMENT"
    echo "═══════════════════════════════════════════════════════"
    
    # Verify I'm the ONLY engineer for these splits
    if [ -f ".split-engineer-lock" ]; then
        EXISTING_ENGINEER=$(cat .split-engineer-lock)
        if [ "$EXISTING_ENGINEER" != "$MY_INSTANCE_ID" ]; then
            echo "❌ FATAL: Another engineer already implementing splits!"
            echo "Existing: $EXISTING_ENGINEER"
            exit 1
        fi
    else
        echo "$MY_INSTANCE_ID" > .split-engineer-lock
        echo "✅ I am the sole implementer for all splits"
    fi
    
    # Read split inventory
    TOTAL_SPLITS=$(grep -c "SPLIT-PLAN-" SPLIT-INVENTORY.md)
    echo "Total splits to implement: $TOTAL_SPLITS"
    echo "Execution mode: SEQUENTIAL"
    
    # Implement each split IN SEQUENCE
    for split_num in $(seq 1 $TOTAL_SPLITS); do
        implement_single_split $split_num
        verify_split_complete $split_num
        echo "✅ Split $split_num complete, moving to next..."
    done
}

# Implement one split at a time
implement_single_split() {
    local split_num=$1
    local split_plan="SPLIT-PLAN-$(printf "%03d" $split_num).md"
    
    echo "════════════════════════════════════════════════"
    echo "Starting Split $split_num of $TOTAL_SPLITS"
    echo "Plan: $split_plan"
    echo "════════════════════════════════════════════════"
    
    # Read the specific split plan
    cat $split_plan
    
    # Create split branch
    SPLIT_BRANCH="${CURRENT_BRANCH}-split-$(printf "%03d" $split_num)"
    git checkout -b $SPLIT_BRANCH
    
    # Implement ONLY files assigned to this split
    # ... implementation ...
    
    # Measure to ensure <800 lines
    $PROJECT_ROOT/tools/line-counter.sh
    
    # Commit and push this split
    git add -A
    git commit -m "feat: implement split $split_num of $TOTAL_SPLITS"
    git push -u origin $SPLIT_BRANCH
    
    echo "✅ Split $split_num implementation complete"
}
```

### Why Sequential Execution is MANDATORY

```yaml
parallel_splits_problems:
  - conflict: Multiple agents modifying related code
  - duplication: Overlapping implementations
  - integration: Merge conflicts between splits
  - dependencies: Split 2 may depend on Split 1
  - testing: Can't test partial implementations
  - coordination: No single source of truth

sequential_splits_benefits:
  - consistency: One mind, one approach
  - dependencies: Natural ordering preserved
  - integration: Clean merges guaranteed
  - testing: Can test after each split
  - debugging: Clear progression path
  - state: Single agent maintains context
```

### State Tracking

```yaml
# Orchestrator tracks single agent for all splits
efforts_in_progress:
  phase1_wave1_api_types_splits:
    status: "implementing_splits"
    total_splits: 3
    current_split: 2
    sw_engineer_instance: "sw-eng-087"  # SAME agent for all
    splits_completed: ["split-001"]
    splits_remaining: ["split-002", "split-003"]
    execution_mode: "SEQUENTIAL"  # NEVER parallel
```

### Integration After Splits

```bash
# After ALL splits complete sequentially
integrate_splits() {
    echo "All splits implemented by single agent"
    echo "Integrating splits in sequence..."
    
    # Merge splits back in order
    for split_num in $(seq 1 $TOTAL_SPLITS); do
        SPLIT_BRANCH="${BASE_BRANCH}-split-$(printf "%03d" $split_num)"
        git merge $SPLIT_BRANCH --no-ff
        echo "✅ Integrated split $split_num"
    done
}
```

## Common Violations to Avoid

### ❌ Spawning Multiple Agents
```bash
# NEVER DO THIS for splits!
Task: sw-engineer for split-001
Task: sw-engineer for split-002  # Different agent = WRONG!
```

### ❌ Parallel Split Execution
```bash
# NEVER implement splits simultaneously
implement_split_001 &  # Background = WRONG!
implement_split_002 &  # Parallel = CONFLICTS!
```

### ❌ Different Agents Per Split
```bash
# NEVER hand off between agents
Agent_1: implements split-001
Agent_2: implements split-002  # WRONG! Same agent must do all!
```

## Correct Flow

```
1. Effort exceeds 800 lines
2. ONE Code Reviewer plans ALL splits
3. ONE SW Engineer implements ALL splits:
   - Split 001 → Complete
   - Split 002 → Complete  
   - Split 003 → Complete
4. Integration of all splits
5. Final review of integrated code
```

## Integration with Other Rules

- **R199**: Single reviewer plans all splits
- **R202**: Single engineer implements all splits (THIS RULE)
- **R197**: One agent per effort (splits are sub-efforts of same effort)
- **R007**: Each split must be <800 lines

## Grading Impact

- **Multiple agents on splits**: -40% (Architectural failure)
- **Parallel split execution**: -35% (Coordination failure)
- **Split conflicts from parallelism**: -30% (Integration failure)
- **Incomplete split sequence**: -25% (Execution failure)

## Summary

**Remember**: 
- ONE agent implements ALL splits from an effort
- Splits are done SEQUENTIALLY, never in parallel
- Complete each split before starting the next
- Same agent maintains context across all splits
- This prevents conflicts and ensures clean integration