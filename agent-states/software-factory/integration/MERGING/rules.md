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
PROJECT_PREFIX=$(jq '.branch_naming.project_prefix' "$SF_INSTANCE_DIR/target-repo-config.yaml" 2>/dev/null || echo "")
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
# R340: Read merge plan location from orchestrator state (NO SEARCHING!)
# Determine integration level from orchestrator state or context
PHASE=$(jq -r '.current_phase // empty' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave // empty' orchestrator-state-v3.json)

if [ -n "$WAVE" ] && [ "$WAVE" != "null" ]; then
    # Wave-level integration
    INTEGRATE_WAVE_EFFORTS_LEVEL="wave"
    WAVE_ID="phase${PHASE}_wave${WAVE}"

    # R340: Lookup merge plan from state
    MERGE_PLAN=$(jq -r ".planning_repo_files.merge_plans.wave[\"${WAVE_ID}\"].file_path // empty" orchestrator-state-v3.json)

    if [ -z "$MERGE_PLAN" ] || [ "$MERGE_PLAN" = "null" ]; then
        echo "❌ R340 VIOLATION: No wave merge plan tracked in state!"
        echo "   Expected at: planning_repo_files.merge_plans.wave[\"${WAVE_ID}\"]"
        exit 340
    fi

elif [ -n "$PHASE" ] && [ "$PHASE" != "null" ]; then
    # Phase-level integration
    INTEGRATE_WAVE_EFFORTS_LEVEL="phase"
    PHASE_ID="phase${PHASE}"

    # R340: Lookup merge plan from state
    MERGE_PLAN=$(jq -r ".planning_repo_files.merge_plans.phase[\"${PHASE_ID}\"].file_path // empty" orchestrator-state-v3.json)

    if [ -z "$MERGE_PLAN" ] || [ "$MERGE_PLAN" = "null" ]; then
        echo "❌ R340 VIOLATION: No phase merge plan tracked in state!"
        echo "   Expected at: planning_repo_files.merge_plans.phase[\"${PHASE_ID}\"]"
        exit 340
    fi

else
    # Project-level integration
    INTEGRATE_WAVE_EFFORTS_LEVEL="project"
    PROJECT_ID="project"

    # R340: Lookup merge plan from state
    MERGE_PLAN=$(jq -r ".planning_repo_files.merge_plans.project[\"${PROJECT_ID}\"].file_path // empty" orchestrator-state-v3.json)

    if [ -z "$MERGE_PLAN" ] || [ "$MERGE_PLAN" = "null" ]; then
        echo "❌ R340 VIOLATION: No project merge plan tracked in state!"
        echo "   Expected at: planning_repo_files.merge_plans.project[\"${PROJECT_ID}\"]"
        exit 340
    fi
fi

# Verify the merge plan file actually exists
if [ ! -f "$MERGE_PLAN" ]; then
    echo "❌ CRITICAL: Merge plan tracked in state but file missing!"
    echo "   Tracked path: $MERGE_PLAN"
    echo "   State was corrupted or file was deleted"
    exit 1
fi

echo "✅ R340: Merge plan loaded from orchestrator state"
echo "   Level: $INTEGRATE_WAVE_EFFORTS_LEVEL"
echo "   Path: $MERGE_PLAN"

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
                MERGED_COUNT=$(grep "${dep}-split-" .software-factory/work-log--${TIMESTAMP}.log 2>/dev/null | grep -c "MERGED" || echo 0)
                
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
            
            if ! grep -q "$PREV_SPLIT.*MERGED" .software-factory/work-log--${TIMESTAMP}.log 2>/dev/null; then
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
- Document EVERY operation in .software-factory/work-log--${TIMESTAMP}.log
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

### 3. Conflict Resolution (R361 Enforced)
```bash
# BEFORE resolving ANY conflict, check R361 compliance
validate_r361_compliance() {
    echo "🔴 R361 ENFORCEMENT: Checking integration limits"

    # Count current non-merge changes
    CHANGES=$(git diff --shortstat HEAD...origin/main --no-merges | grep -oE '[0-9]+ (insertion|deletion)' | awk '{sum+=$1} END {print sum}')

    if [ "$CHANGES" -gt 40 ]; then
        echo "⚠️ WARNING: Approaching 50-line limit!"
        echo "Current changes: $CHANGES lines"
        echo "Remaining budget: $((50 - CHANGES)) lines"
    fi

    # Check for new files (ABSOLUTELY FORBIDDEN)
    NEW_FILES=$(git status --porcelain | grep "^A" | wc -l)
    if [ "$NEW_FILES" -gt 0 ]; then
        echo "🔴🔴🔴 R361 VIOLATION: New files detected!"
        echo "Integration CANNOT create new files!"
        git status --porcelain | grep "^A"
        exit 361
    fi
}
```

**Conflict Resolution Rules:**
- Resolve ALL conflicts completely
- Document resolution decisions
- Never leave conflict markers
- Test compilation after resolution
- NEVER add new code to "fix" conflicts (R361)
- ONLY choose between existing versions

## SUPREME LAWS IN EFFECT
- R262 - Merge Operation Protocols
  - NEVER modify original branches
  - NEVER use cherry-pick
  - Create synthesis branches if needed
- R361 - Integration Conflict Resolution Only
  - NO new files or packages
  - NO adapter or wrapper code
  - Maximum 50 lines of changes total (excluding merge commits)
  - Integration = conflict resolution ONLY

## SF 3.0 Integration Container Context

This state operates within integration containers tracked in `integration-containers.json`:
- Reads container configuration including branches to merge, iteration number, and convergence settings
- Updates container status to "MERGING" during merge operations in `integration-containers.json`
- Records merge success/failure and conflict resolution details to container state
- Tracks convergence metrics and merge progress for the integration cycle
- If cross-container bugs are discovered during merging, the fix cascade mechanism creates `fix-cascade-state.json` per SF 3.0 architecture

All container state updates are atomic per R288. The orchestrator monitors `state_machine.current_state` in orchestrator-state-v3.json to track integration progress.

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
- .software-factory/work-log--${TIMESTAMP}.log contains every operation
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

