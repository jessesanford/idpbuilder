# Rule R199: Single Code Reviewer for All Split Planning

## Rule Statement
When an effort exceeds 800 lines and requires splitting, the orchestrator MUST spawn EXACTLY ONE code reviewer agent to plan ALL splits for that effort. This single reviewer maintains complete context and prevents duplication across split boundaries.

## Criticality Level
**BLOCKING** - Multiple reviewers cause split conflicts and duplication

## Enforcement Mechanism
- **Technical**: Single agent instance for all split planning
- **Behavioral**: Reviewer refuses if another reviewer already planning splits
- **Grading**: -35% for multiple reviewers on same split effort

## Core Principle

```
1 Over-Limit Effort = 1 Code Reviewer = All Split Plans
```

## Detailed Requirements

### ORCHESTRATOR: Single Reviewer Protocol

```bash
# ❌❌❌ WRONG - Multiple reviewers for splits
handle_overlimit_wrong() {
    # Effort is 2400 lines, needs 3 splits
    
    # WRONG! Spawning multiple reviewers
    Task: code-reviewer
    Plan split 1 for api-types
    
    Task: code-reviewer  # WRONG! Different instance
    Plan split 2 for api-types
    
    Task: code-reviewer  # WRONG! Yet another instance
    Plan split 3 for api-types
}

# ✅✅✅ CORRECT - Single reviewer for ALL splits
handle_overlimit_correct() {
    local effort_name="api-types"
    local effort_size=2400
    local splits_needed=$((effort_size / 700 + 1))  # 4 splits
    
    echo "═══════════════════════════════════════════════════════"
    echo "EFFORT OVER LIMIT: $effort_name"
    echo "Size: $effort_size lines (limit: 800)"
    echo "Splits needed: $splits_needed"
    echo "═══════════════════════════════════════════════════════"
    
    # Spawn SINGLE reviewer for ALL splits
    Task: code-reviewer
    Working directory: efforts/phase1/wave1/$effort_name
    
    CRITICAL: You are the ONLY reviewer for this split effort.
    You must plan ALL $splits_needed splits to avoid duplication.
    
    Your task:
    1. Analyze the COMPLETE $effort_size line effort
    2. Create $splits_needed split plans (each <800 lines)
    3. Ensure NO duplication across splits
    4. Ensure NO gaps in coverage
    5. Create logical groupings that compile independently
    6. 🔴 PARAMOUNT (R307): Each split MUST be independently mergeable
       - Split 1 must work alone (even if others never merge)
       - Split 2 must work with just Split 1 merged
       - Use feature flags for incomplete features across splits
       - NO split can break existing functionality
    
    Output (with timestamps):
    - SPLIT-PLAN-001-YYYYMMDD-HHMMSS.md (lines 1-700)
    - SPLIT-PLAN-002-YYYYMMDD-HHMMSS.md (lines 701-1400) 
    - SPLIT-PLAN-003-YYYYMMDD-HHMMSS.md (lines 1401-2100)
    - SPLIT-PLAN-004-YYYYMMDD-HHMMSS.md (lines 2101-2400)
    
    You have COMPLETE context. Plan ALL splits now.
}
```

### CODE REVIEWER: Complete Split Context

```bash
# When assigned to plan splits
handle_split_planning() {
    echo "═══════════════════════════════════════════════════════"
    echo "SPLIT PLANNING ASSIGNMENT"
    echo "═══════════════════════════════════════════════════════"
    
    # Verify I'm the ONLY reviewer for this effort
    if [ -f ".split-reviewer-lock" ]; then
        EXISTING_REVIEWER=$(cat .split-reviewer-lock)
        if [ "$EXISTING_REVIEWER" != "$MY_INSTANCE_ID" ]; then
            echo "❌ FATAL: Another reviewer already planning splits!"
            echo "Existing: $EXISTING_REVIEWER"
            exit 1
        fi
    else
        echo "$MY_INSTANCE_ID" > .split-reviewer-lock
        echo "✅ I am the sole split planner for this effort"
    fi
    
    # Analyze ENTIRE effort
    echo "Analyzing complete effort for optimal splitting..."
    TOTAL_SIZE=$(./tools/line-counter.sh | grep "Total" | awk '{print $NF}')
    SPLITS_NEEDED=$((TOTAL_SIZE / 700 + 1))
    
    echo "Total size: $TOTAL_SIZE lines"
    echo "Splits needed: $SPLITS_NEEDED"
    echo "I will plan ALL $SPLITS_NEEDED splits with complete context"
    
    # Create comprehensive split plan
    create_all_split_plans
}

# Create ALL splits with deduplication awareness
create_all_split_plans() {
    local TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    local split_inventory="SPLIT-INVENTORY-${TIMESTAMP}.md"
    
    cat > "$split_inventory" << EOF
# Complete Split Plan for $EFFORT_NAME
Total Size: $TOTAL_SIZE lines
Splits: $SPLITS_NEEDED
Planner: $MY_INSTANCE_ID (sole reviewer)
Created: $TIMESTAMP

## Split Boundaries (NO OVERLAPS)
EOF
    
    # Plan each split with awareness of others
    for i in $(seq 1 $SPLITS_NEEDED); do
        echo "Planning split $i of $SPLITS_NEEDED..."
        
        # Create timestamped split plan with explicit boundaries
        cat > "SPLIT-PLAN-$(printf "%03d" $i)-${TIMESTAMP}.md" << EOF
# Split $i of $SPLITS_NEEDED
Reviewer: $MY_INSTANCE_ID (same for all splits)

## Boundaries
- Previous split ends at: [specific file:line]
- This split starts at: [specific file:line]
- This split ends at: [specific file:line]
- Next split starts at: [specific file:line]

## Files in This Split (NO duplication with other splits)
EOF
        
        # Add files ensuring no duplication
        plan_split_files $i $SPLITS_NEEDED
        
        # Update inventory
        echo "- Split $i: [describe scope] (lines X-Y)" >> "$split_inventory"
    done
    
    # Final verification
    echo "" >> "$split_inventory"
    echo "## Verification Checklist" >> "$split_inventory"
    echo "- [ ] No file appears in multiple splits" >> "$split_inventory"
    echo "- [ ] All files are covered" >> "$split_inventory"
    echo "- [ ] Each split <800 lines" >> "$split_inventory"
    echo "- [ ] Logical groupings maintained" >> "$split_inventory"
    echo "- [ ] Dependencies respected" >> "$split_inventory"
}
```

### Deduplication Strategies

```bash
# Strategy 1: File-based splitting (no file in multiple splits)
split_by_files() {
    # List all files
    find . -name "*.go" -o -name "*.ts" | sort > all_files.txt
    
    # Assign files to splits
    current_split=1
    current_size=0
    
    while read -r file; do
        file_size=$(wc -l < "$file")
        
        if [ $((current_size + file_size)) -gt 700 ]; then
            # Start new split
            current_split=$((current_split + 1))
            current_size=0
        fi
        
        echo "$file" >> "split-$current_split-files.txt"
        current_size=$((current_size + file_size))
    done < all_files.txt
}

# Strategy 2: Module-based splitting (keep related code together)
split_by_modules() {
    # Group by package/module
    for module in api controllers webhooks middleware; do
        module_size=$(find . -path "*/$module/*" | xargs wc -l | tail -1 | awk '{print $1}')
        echo "$module: $module_size lines" >> module_sizes.txt
    done
    
    # Assign modules to splits ensuring no overlap
    # Each module goes to exactly ONE split
}

# Strategy 3: Dependency-based splitting
split_by_dependencies() {
    # Ensure dependencies are in same or earlier splits
    # Types → Split 1
    # Interfaces using types → Split 1 or 2
    # Implementations → Split 2 or 3
    # Tests → Same split as implementation
}
```

## State Management

### Tracking Single Reviewer Assignment
```yaml
# orchestrator-state-v3.json
split_planning:
  phase1_wave1_api_types:
    status: "planning"
    reviewer_instance: "code-reviewer-042"  # ONLY ONE
    total_size: 2400
    splits_needed: 4
    started: "2024-01-20T10:00:00Z"
    
  # If another effort needs splitting
  phase1_wave1_controllers:
    status: "planning"
    reviewer_instance: "code-reviewer-043"  # DIFFERENT effort, different reviewer
    total_size: 1800
    splits_needed: 3
```

### Split Plan Outputs
```
efforts/phase1/wave1/api-types/
├── SPLIT-INVENTORY.md          # Master list of all splits
├── SPLIT-PLAN-001.md           # Split 1 details
├── SPLIT-PLAN-002.md           # Split 2 details
├── SPLIT-PLAN-003.md           # Split 3 details
├── SPLIT-PLAN-004.md           # Split 4 details
└── .split-reviewer-lock        # Ensures single reviewer
```

## Benefits of Single Reviewer

1. **Complete Context**: Sees entire effort, makes optimal decisions
2. **No Duplication**: Single mind planning prevents overlaps
3. **Consistent Strategy**: One approach across all splits
4. **Dependency Awareness**: Can properly sequence dependencies
5. **Efficient Planning**: One pass through code, not multiple

## Common Violations to Avoid

### ❌ Multiple Reviewers
```bash
# WRONG - Each split gets different reviewer
for split in 1 2 3; do
    Task: code-reviewer  # NEW instance each time!
    Plan split $split
done
```

### ❌ Sequential Additive Planning
```bash
# WRONG - Planning splits one at a time without full context
Task: code-reviewer
Plan split 1

# Later, after split 1 done...
Task: code-reviewer  # Different instance!
Plan split 2 based on what's left
```

### ❌ Overlapping Assignments
```bash
# WRONG - Multiple reviewers working simultaneously
# Reviewer A planning splits 1-2
# Reviewer B planning splits 3-4
# Result: Conflicts and duplication!
```

## Correct Implementation Flow

```
1. SW Engineer implements effort
2. Effort exceeds 800 lines
3. Orchestrator detects overlimit
4. Orchestrator spawns SINGLE code reviewer
5. Code reviewer analyzes ENTIRE effort
6. Code reviewer creates ALL split plans
7. Code reviewer ensures no duplication
8. Code reviewer completes and terminates
9. Orchestrator spawns SW engineers for each split
```

## Integration with Other Rules

- **R197**: One agent per effort (splits are sub-efforts)
- **R056**: Split plan creation (single reviewer creates all)
- **R007**: Size limits (each split must be <800)
- **R153**: Review effectiveness (complete context improves quality)

## Grading Impact

- **Multiple reviewers for same split effort**: -35% (Major violation)
- **Duplication across splits**: -25% (Planning failure)
- **Missing coverage in splits**: -20% (Incomplete planning)
- **Overlapping split boundaries**: -30% (Coordination failure)

## Summary

**Remember**: When splitting is needed, spawn ONE code reviewer who will:
- Analyze the COMPLETE over-limit effort
- Plan ALL splits with full context
- Ensure NO duplication between splits
- Create a comprehensive split inventory
- Then terminate (per R197)