# Code Reviewer - WAVE_MERGE_PLANNING State Rules

## 🚨 CRITICAL: MERGE PLAN CREATION ONLY - NO EXECUTION! 🚨

### PRIMARY PURPOSE
Create a detailed MERGE PLAN document that the Integration Agent will follow.
DO NOT execute any merges - only plan them!

### 🔴🔴🔴 CRITICAL RULE: NO INTEGRATION BRANCHES IN MERGE PLAN! 🔴🔴🔴

**YOU MUST USE ONLY ORIGINAL EFFORT BRANCHES!**
- ✅ CORRECT: phase3/wave2/effort1-api-gateway
- ✅ CORRECT: phase3/wave2/effort2-controller-split1  
- ❌ WRONG: wave2-integration-20250827
- ❌ WRONG: phase3-integration-20250827

Integration branches are TARGETS, not SOURCES! Never merge from an integration branch.

## State Context
You are creating a comprehensive merge plan for wave integration. The orchestrator has already set up the integration infrastructure (directory and branch). Your role is to analyze what needs to be merged and create a detailed plan.

## 🔴🔴🔴 CRITICAL: FIRST ACTION - CD TO INTEGRATION DIRECTORY 🔴🔴🔴

**YOUR FIRST ACTION MUST BE:**
```bash
# The orchestrator will tell you the integration directory path
# You MUST CD to this directory before creating WAVE-MERGE-PLAN.md
cd [integration directory path from orchestrator]
pwd  # Verify you are in the integration workspace
```

**NEVER create WAVE-MERGE-PLAN.md in the root directory!**
**ALWAYS create it IN the integration workspace!**

## MERGE PLAN REQUIREMENTS

### 1. Analyze All Effort Branches
```bash
#!/bin/bash
# Find all effort branches for this wave
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)

echo "📊 Analyzing effort branches for Phase $PHASE Wave $WAVE..."

# List all effort directories
for effort_dir in /efforts/phase${PHASE}/wave${WAVE}/*/; do
    [[ "$effort_dir" == *"integration-workspace"* ]] && continue
    
    effort=$(basename "$effort_dir")
    cd "$effort_dir"
    
    # Get branch info
    current_branch=$(git branch --show-current)
    
    # Check if this is a split
    if [[ "$current_branch" == *"-split"* ]]; then
        echo "✅ INCLUDE: $current_branch (split branch)"
    elif git log --oneline | grep -q "too large"; then
        echo "❌ EXCLUDE: $current_branch (original, too large)"
    else
        echo "✅ INCLUDE: $current_branch (within size limit)"
    fi
    
    # Identify base branch
    base_commit=$(git merge-base HEAD origin/main)
    echo "  Base: main at $base_commit"
done
```

### 2. Determine Merge Order
```markdown
## Merge Order Analysis

The correct merge order is critical to avoid conflicts and preserve intent.

### Checking Branch Bases:
- Branch A based on: main at commit abc123  
- Branch B based on: main at commit def456
- Branch C based on: Branch A at commit ghi789

### Therefore Order Must Be:
1. Branch A (base: main)
2. Branch C (base: Branch A) - depends on A
3. Branch B (base: main) - independent

### Dependency Graph:
```
main
├── Branch A
│   └── Branch C (depends on A)
└── Branch B (independent)
```
```

### 3. Handle Splits Correctly

#### 🔴🔴🔴 CRITICAL: Split Merge Ordering with Dependencies 🔴🔴🔴

**SUPREME LAW: Dependencies are at EFFORT level, not SPLIT level!**

When creating merge plans with splits and dependencies:
1. ALL splits of an effort must be listed sequentially
2. ALL splits must complete before dependent efforts
3. NEVER interleave dependent efforts between splits

```markdown
## Split Handling

### Effort E3.1.3 - Contexts Implementation
Original branch exceeded 800 lines and was split:

**EXCLUDE (too large):**
- phase3/wave1/effort3-contexts (1,234 lines)

**INCLUDE IN THIS ORDER (properly sized splits):**  
1. phase3/wave1/effort3-contexts-split-001 (423 lines) - base: main
2. phase3/wave1/effort3-contexts-split-002 (389 lines) - base: split-001
3. phase3/wave1/effort3-contexts-split-003 (401 lines) - base: split-002

Total after splits: 1,213 lines (compliant)

### Dependency Handling Example

If E3.1.4 depends on E3.1.3:
✅ CORRECT ORDER:
1. effort3-contexts-split-001
2. effort3-contexts-split-002
3. effort3-contexts-split-003
4. effort4-dependent-feature  # NOW has complete E3.1.3

❌ WRONG ORDER:
1. effort3-contexts-split-001
2. effort4-dependent-feature  # ERROR: E3.1.3 incomplete!
3. effort3-contexts-split-002
4. effort3-contexts-split-003
```

### 4. Create WAVE-MERGE-PLAN.md

**CRITICAL LOCATION REQUIREMENT:**
- **YOU MUST BE IN THE INTEGRATION DIRECTORY** (cd there first!)
- **Create the file IN the current directory** (./WAVE-MERGE-PLAN.md)
- **DO NOT create it in the root directory or anywhere else!**

```bash
# Verify you are in the integration workspace
pwd  # Should show: /efforts/phase${PHASE}/wave${WAVE}/integration-workspace

# Create the merge plan HERE
cat > WAVE-MERGE-PLAN.md << 'EOF'
[merge plan content]
EOF
```

```markdown
# Wave ${WAVE} Merge Plan

**Generated:** $(date -u +%Y-%m-%dT%H:%M:%SZ)
**Code Reviewer:** code-reviewer
**State:** WAVE_MERGE_PLANNING

## Target Integration Branch
- **Branch Name:** phase${PHASE}-wave${WAVE}-integration-${timestamp}
- **Base:** main at ${commit}
- **Location:** /efforts/phase${PHASE}/wave${WAVE}/integration-workspace

## Branches to Merge (IN ORDER)

### 1. phase${PHASE}/wave${WAVE}/effort1-api-types
- **Type:** Original effort branch
- **Base:** main at abc123
- **Size:** 542 lines
- **Dependencies:** None
- **Conflicts Expected:** None
- **Merge Command:**
  ```bash
  git fetch origin phase${PHASE}/wave${WAVE}/effort1-api-types
  git merge origin/phase${PHASE}/wave${WAVE}/effort1-api-types --no-ff \
    -m "Integrate effort1-api-types into wave ${WAVE}"
  ```

### 2. phase${PHASE}/wave${WAVE}/effort2-controller-split1
- **Type:** Split branch (1 of 3)
- **Base:** main at abc123  
- **Size:** 398 lines
- **Dependencies:** None
- **Conflicts Expected:** Possible in controller.go
- **Merge Command:**
  ```bash
  git fetch origin phase${PHASE}/wave${WAVE}/effort2-controller-split1
  git merge origin/phase${PHASE}/wave${WAVE}/effort2-controller-split1 --no-ff \
    -m "Integrate effort2-controller-split1 into wave ${WAVE}"
  ```

### 3. phase${PHASE}/wave${WAVE}/effort2-controller-split2
- **Type:** Split branch (2 of 3)
- **Base:** effort2-controller-split1 at def456
- **Size:** 412 lines
- **Dependencies:** Must merge after split1
- **Conflicts Expected:** None (sequential splits)
- **Merge Command:**
  ```bash
  git fetch origin phase${PHASE}/wave${WAVE}/effort2-controller-split2
  git merge origin/phase${PHASE}/wave${WAVE}/effort2-controller-split2 --no-ff \
    -m "Integrate effort2-controller-split2 into wave ${WAVE}"
  ```

[Continue for all branches...]

## Excluded Branches (Too Large)
These original branches were split and should NOT be merged:
- phase${PHASE}/wave${WAVE}/effort2-controller (original, 1,456 lines)

## Merge Strategy
1. **Merge Type:** --no-ff (preserve branch history)
2. **Conflict Resolution:** Favor newer implementation when conflicts occur
3. **Testing:** Run unit tests after each merge
4. **Validation:** Check total size after all merges

## Expected Conflicts
Based on branch analysis, conflicts are likely in:
- `pkg/controller/controller.go` - Modified by both effort2 and effort3
- `api/v1/types.go` - Extended by multiple efforts

## Validation Steps
1. After each merge:
   ```bash
   go test ./...
   ```
2. After all merges:
   ```bash
   make test-integration
   $PROJECT_ROOT/tools/line-counter.sh -c $(git branch --show-current)
   ```
3. Final validation:
   - Verify all effort features are present
   - Confirm no effort was missed
   - Check combined size is reasonable

## Risk Assessment
- **Low Risk:** Sequential splits should merge cleanly
- **Medium Risk:** Controller modifications may conflict
- **Mitigation:** Test after each merge to catch issues early

## Integration Agent Instructions
1. CD to integration directory before starting
2. Execute merges in the EXACT order specified
3. Run tests after EACH merge
4. Document any conflicts encountered
5. Create work-log.md with all operations
6. Generate INTEGRATION-REPORT.md when complete
```

## Validation Before Completion

```bash
#!/bin/bash
# Validate the merge plan is complete

validate_merge_plan() {
    local plan_file="$1"
    
    # Check all required sections exist
    for section in "Target Integration Branch" "Branches to Merge" "Merge Strategy" "Validation Steps"; do
        if ! grep -q "## $section" "$plan_file"; then
            echo "❌ Missing required section: $section"
            return 1
        fi
    done
    
    # Verify no integration branches as sources
    if grep -E "origin/.*-integration" "$plan_file" | grep -v "^##"; then
        echo "❌ ERROR: Integration branches found as merge sources!"
        return 1
    fi
    
    # Check merge commands are present
    merge_count=$(grep -c "git merge origin/" "$plan_file")
    if [[ $merge_count -lt 1 ]]; then
        echo "❌ No merge commands found in plan"
        return 1
    fi
    
    echo "✅ Merge plan validation passed"
    echo "📊 Total merges planned: $merge_count"
    return 0
}

# Run validation
validate_merge_plan "WAVE-MERGE-PLAN.md"
```

## State Transitions

From WAVE_MERGE_PLANNING state:
- **PLAN_COMPLETE** → Return to orchestrator
- **VALIDATION_FAILED** → Fix and re-validate

## Critical Success Criteria

1. ✅ CD'd to integration directory FIRST (verified with pwd)
2. ✅ WAVE-MERGE-PLAN.md created IN the integration directory (not root!)
3. ✅ All effort branches analyzed and categorized
4. ✅ Merge order determined based on dependencies
5. ✅ NO integration branches used as sources
6. ✅ Split branches handled correctly
7. ✅ Clear instructions for Integration Agent
8. ✅ Validation steps included
9. ✅ Risk assessment documented

## Common Mistakes to Avoid

1. **Using integration branches as sources**
   - ❌ WRONG: Merge from wave1-integration
   - ✅ RIGHT: Merge from effort branches only

2. **Including "too large" original branches**
   - ❌ WRONG: Include both original and splits
   - ✅ RIGHT: Exclude original, include only splits

3. **Wrong merge order**
   - ❌ WRONG: Random order
   - ✅ RIGHT: Dependency-aware ordering

4. **Executing merges**
   - ❌ WRONG: Actually running git merge
   - ✅ RIGHT: Only documenting what to merge

5. **Missing validation steps**
   - ❌ WRONG: No test commands
   - ✅ RIGHT: Test after each merge

---
### ⚠️⚠️⚠️ RULE R261 - Code Reviewer Merge Plan No Execution
**Source:** rule-library/R261-code-reviewer-merge-plan-no-execution.md
**Criticality:** WARNING - Violation = Role confusion

Code Reviewer creates merge plans ONLY. NEVER executes merges. That's the Integration Agent's job.
---

---
### 🔴🔴🔴 RULE R262 - No Integration Branches as Sources
**Source:** rule-library/R262-no-integration-branches-as-sources.md  
**Criticality:** SUPREME - Violation = Recursive integration chaos

CRITICAL: Only original effort branches in merge plans. Integration branches are TARGETS not SOURCES. This prevents recursive integration issues.
---