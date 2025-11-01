# Code Reviewer - CREATE_SPLIT_INVENTORY State Rules

## State Context
You are creating a split inventory and individual split plans within SF 3.0 for an implementation that exceeds the 800-line limit. The splits must be logical, buildable, and testable.

## SF 3.0 Split Planning Context

In this state, the Code Reviewer creates split plans:
- Reads oversized effort details from orchestrator-state-v3.json to understand what needs splitting
- Creates split inventory and individual split plans with proper metadata per R340/R383
- Reports split plan locations for orchestrator to record in orchestrator-state-v3.json `metadata_locations.split_plans`
- Ensures split strategy allows independent review and merging per R501
- All split planning artifacts stored with timestamps and atomic state updates per R288

## 🔴🔴🔴 CRITICAL: VERIFY MEASUREMENT WAS CORRECT! 🔴🔴🔴

**BEFORE CREATING SPLITS, CONFIRM THE SIZE VIOLATION IS REAL:**
1. **Was line-counter.sh used?** If not, STOP and remeasure!
2. **Did it auto-detect base?** Look for "🎯 Detected base:" in output
3. **Is the count realistic?** 11,876 lines for a ~500 line effort = WRONG!
4. **If measurement was manual** (wc -l, git diff, etc.) = INVALID, REMEASURE!

**COMMON MISTAKE TO AVOID:**
- Manual counting against main: Shows ALL code (11,876 lines)
- Correct tool usage: Shows ONLY changes (~500 lines)
- **If you see inflated numbers, STOP and use line-counter.sh properly!**

---

## 🔴🔴🔴 CRITICAL: Split Plan Creation Requirements 🔴🔴🔴

### Split Inventory Structure
You MUST create these files in the too-large branch:

1. **SPLIT-INVENTORY.md** - Master inventory of all splits
2. **SPLIT-PLAN-001.md** - Detailed plan for first split
3. **SPLIT-PLAN-002.md** - Detailed plan for second split
4. **SPLIT-PLAN-XXX.md** - Continue for each split needed

### Size Requirements for Splits
- Each split MUST be ≤700 lines (soft limit)
- HARD LIMIT: 800 lines (absolute maximum)
- Aim for 500-600 lines per split for safety margin
- Account for incremental building (each split adds to previous)

---

## Split Planning Strategy

### 1. Analyze Current Implementation
```python
def analyze_implementation_for_splits(effort_dir, total_lines):
    """Analyze implementation to determine split strategy"""
    
    analysis = {
        'total_lines': total_lines,
        'target_splits': math.ceil(total_lines / 600),  # Aim for 600 lines per split
        'components': identify_logical_components(effort_dir),
        'dependencies': analyze_dependencies(effort_dir),
        'test_boundaries': identify_test_boundaries(effort_dir)
    }
    
    # Identify natural split points
    split_points = []
    
    # Option 1: Component-based splits
    if has_clear_components(analysis['components']):
        split_points = plan_component_splits(analysis['components'])
    
    # Option 2: Layer-based splits (data, business, API)
    elif has_layer_architecture(effort_dir):
        split_points = plan_layer_splits(effort_dir)
    
    # Option 3: Feature-based splits
    else:
        split_points = plan_feature_splits(effort_dir)
    
    return {
        'strategy': determine_best_strategy(split_points),
        'estimated_splits': len(split_points),
        'split_boundaries': split_points
    }
```

### 2. Create Split Inventory

SPLIT-INVENTORY.md template:
```markdown
# Split Inventory
Date: [timestamp]
Effort: [effort-name]
Total Lines: [total-lines]
Required Splits: [number]

## Split Strategy
[Component-based / Layer-based / Feature-based]

## Split Overview
| Split | Description | Estimated Lines | Dependencies | Status |
|-------|-------------|-----------------|--------------|--------|
| 001   | [Core infrastructure and interfaces] | ~500 | None | PENDING |
| 002   | [Business logic implementation] | ~600 | Split-001 | PENDING |
| 003   | [API endpoints and handlers] | ~550 | Split-001, Split-002 | PENDING |

## Dependency Graph
```
Split-001 (Foundation)
    ↓
Split-002 (Logic)
    ↓
Split-003 (API)
```

## Testing Strategy
- Split-001: Unit tests for core components
- Split-002: Integration tests for business logic
- Split-003: End-to-end API tests

## Build Verification
Each split must:
1. Build successfully on its own
2. Pass its own tests
3. Not break previous splits
```

### 3. Create Individual Split Plans

SPLIT-PLAN-XXX.md template:
```markdown
# Split Plan 001 - [Description]
Date: [timestamp]
Effort: [effort-name]
Split: 001 of [total]

## Scope
This split implements [specific functionality].

## Files to Include
### New Files
- `path/to/new/file1.go` - [description]
- `path/to/new/file2.go` - [description]

### Modified Files
- `path/to/existing/file.go` - Add [specific changes]

## Implementation Checklist
- [ ] Create base interfaces
- [ ] Implement core types
- [ ] Add configuration structures
- [ ] Write unit tests
- [ ] Update documentation

## Size Estimate
- New code: ~400 lines
- Tests: ~100 lines
- Total: ~500 lines

## Dependencies
- Prerequisites: None (first split)
- External deps: [list any new imports]

## Build Commands
```bash
# After implementing this split
go build ./...
go test ./...
```

## Success Criteria
- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] No impact on existing code
- [ ] Within 700-line limit

## Notes for SW Engineer
[Any special instructions or warnings]
```

### 4. Sequential Split Requirements

**CRITICAL**: Splits build on each other!
```python
def validate_split_sequence(splits):
    """Ensure splits can be implemented sequentially"""
    
    cumulative_lines = 0
    
    for i, split in enumerate(splits):
        # Each split adds to the previous total
        cumulative_lines += split['estimated_lines']
        
        if i == 0:
            # First split from base
            if split['estimated_lines'] > 700:
                return False, "First split too large"
        else:
            # Subsequent splits are incremental
            if split['estimated_lines'] > 700:
                return False, f"Split {i+1} too large"
        
        # Verify dependencies are met
        for dep in split['dependencies']:
            if not is_dependency_satisfied(dep, splits[:i]):
                return False, f"Dependency {dep} not met for split {i+1}"
    
    return True, "Split sequence valid"
```

## Common Split Patterns

### Component-Based Splits
```
Split-001: Core interfaces and types
Split-002: Data access layer
Split-003: Business logic
Split-004: API/Controller layer
```

### Feature-Based Splits
```
Split-001: Feature A (complete vertical slice)
Split-002: Feature B (complete vertical slice)
Split-003: Feature C (complete vertical slice)
```

### Layer-Based Splits
```
Split-001: Database models and migrations
Split-002: Repository layer
Split-003: Service layer
Split-004: HTTP handlers
```

## File Commit Requirements

After creating all split files:
```bash
# In the too-large branch
git add SPLIT-INVENTORY.md
git add .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT}/SPLIT-PLAN--*.md
git commit -m "feat: create split plans for oversized implementation

- Total lines: [number]
- Splits required: [number]
- Strategy: [component/layer/feature]-based"
git push
```

## State Transitions

From CREATE_SPLIT_INVENTORY state:
- **SPLITS_PLANNED** → COMPLETED (Plans created and committed)
- **PLANNING_ERROR** → ERROR_RECOVERY (Cannot create viable splits)

## Success Criteria
- ✅ SPLIT-INVENTORY.md created with overview
- ✅ Individual SPLIT-PLAN-XXX.md for each split
- ✅ Each split ≤700 lines (soft) / 800 lines (hard)
- ✅ Dependencies clearly mapped
- ✅ Build/test strategy defined
- ✅ All files committed and pushed

## Failure Triggers
- ❌ Missing split inventory = -40% penalty
- ❌ Splits exceed 800 lines = -100% FAILURE
- ❌ Unbuildable split sequence = -50% penalty
- ❌ Files not committed = -30% penalty

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

