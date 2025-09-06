# Integration Agent - MERGING State Rules

## State Definition
The MERGING state is where actual branch integration occurs following the MERGE PLAN created by Code Reviewer.

## Required Actions

### 1. Verify Integration Branch Already Created (R014 Compliance)
```bash
# The orchestrator already created the integration branch
# We should already be on it from INIT state
CURRENT_BRANCH=$(git branch --show-current)
if [[ ! "$CURRENT_BRANCH" =~ integration ]]; then
    echo "❌ ERROR: Not on integration branch!"
    echo "Current branch: $CURRENT_BRANCH"
    exit 1
fi

# Verify branch follows R014 naming convention with project prefix
PROJECT_PREFIX=$(yq '.branch_naming.project_prefix' "$SF_INSTANCE_DIR/target-repo-config.yaml" 2>/dev/null || echo "")
if [ -n "$PROJECT_PREFIX" ] && [ "$PROJECT_PREFIX" != "null" ]; then
    if [[ ! "$CURRENT_BRANCH" =~ ^"$PROJECT_PREFIX"/ ]]; then
        echo "⚠️ WARNING: Integration branch missing project prefix: $PROJECT_PREFIX"
        echo "Branch should start with: $PROJECT_PREFIX/"
    fi
fi

echo "✅ On integration branch: $CURRENT_BRANCH"
```

### 2. Sequential Merging Following MERGE PLAN

#### 🔴🔴🔴 CRITICAL: Split-Aware Merge Ordering 🔴🔴🔴

**SUPREME LAW: When efforts have splits, ALL splits must merge before dependent efforts!**

```bash
# Read merge plan to extract branches
if [ -f "WAVE-MERGE-PLAN.md" ]; then
    MERGE_PLAN="WAVE-MERGE-PLAN.md"
elif [ -f "PHASE-MERGE-PLAN.md" ]; then
    MERGE_PLAN="PHASE-MERGE-PLAN.md"
else
    echo "❌ NO MERGE PLAN FOUND!"
    exit 1
fi

echo "📋 Following merge plan: $MERGE_PLAN"

# Validate merge order respects split completeness
validate_merge_order() {
    local branch="$1"
    local effort=$(echo "$branch" | sed 's/-split-[0-9]*//')
    
    # If this is NOT a split, check all dependencies are complete
    if [[ ! "$branch" =~ -split- ]]; then
        # This is a regular effort, check its dependencies
        DEPENDENCIES=$(grep -A5 "branch: $branch" "$MERGE_PLAN" | grep "depends_on:" | sed 's/.*\[\(.*\)\].*/\1/')
        
        for dep in $DEPENDENCIES; do
            # Check if dependency has unmerged splits
            if grep -q "${dep}-split-" "$MERGE_PLAN"; then
                # Count how many splits exist
                SPLIT_COUNT=$(grep -c "${dep}-split-" "$MERGE_PLAN")
                # Count how many are already merged
                MERGED_COUNT=$(grep "${dep}-split-" work-log.md 2>/dev/null | grep -c "MERGED" || echo 0)
                
                if [ "$MERGED_COUNT" -lt "$SPLIT_COUNT" ]; then
                    echo "❌ ERROR: Cannot merge $branch yet!"
                    echo "   Dependency $dep has $SPLIT_COUNT splits but only $MERGED_COUNT merged"
                    echo "   All splits must be merged first!"
                    return 1
                fi
            fi
        done
    fi
    
    # If this is a split, ensure previous splits are merged
    if [[ "$branch" =~ -split-([0-9]+) ]]; then
        SPLIT_NUM="${BASH_REMATCH[1]}"
        if [ "$SPLIT_NUM" -gt 1 ]; then
            PREV_NUM=$((SPLIT_NUM - 1))
            PREV_SPLIT=$(printf "%s-split-%03d" "$effort" "$PREV_NUM")
            
            if ! grep -q "$PREV_SPLIT.*MERGED" work-log.md 2>/dev/null; then
                echo "❌ ERROR: Cannot merge $branch!"
                echo "   Previous split $PREV_SPLIT not merged yet"
                echo "   Splits must be merged sequentially!"
                return 1
            fi
        fi
    fi
    
    return 0
}
```

**CRITICAL: Execute merges in EXACT order from plan**
- Validate each merge against split/dependency rules
- Extract merge commands from the plan
- Execute each merge with --no-ff
- Resolve conflicts as encountered
- Document EVERY operation in work-log.md
- Run tests after each merge (if specified in plan)

**Example of Correct Split Handling:**
```
If E1 has 3 splits and E2 depends on E1:
1. Merge E1-split-001 ✓
2. Merge E1-split-002 ✓  
3. Merge E1-split-003 ✓
4. NOW merge E2 ✓ (E1 is complete)

NEVER:
1. Merge E1-split-001
2. Merge E2 ✗ (E1 incomplete!)
```

### 3. Conflict Resolution
- Resolve ALL conflicts completely
- Document resolution decisions
- Never leave conflict markers
- Test compilation after resolution

## SUPREME LAWS IN EFFECT
- R262 - Merge Operation Protocols
  - NEVER modify original branches
  - NEVER use cherry-pick
  - Create synthesis branches if needed

## Work Log Requirements
Every operation MUST be documented:
```markdown
## Operation N: [timestamp]
Command: [exact command]
Result: [output]
Conflicts: [if any]
Resolution: [how resolved]
```

## Transition Rules
- Can transition to: TESTING (after all merges complete)
- Cannot transition if: Conflicts unresolved
- Must have clean working tree

## Success Criteria
- All planned branches merged
- All conflicts resolved
- work-log.md contains every operation
- No original branches modified

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**
