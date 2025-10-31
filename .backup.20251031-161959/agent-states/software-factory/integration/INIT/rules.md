# Integration Agent - INIT State Rules

## 🚨 MANDATORY STARTUP ACKNOWLEDGMENT 🚨

### IMMEDIATE ACTIONS UPON STARTUP

**YOU MUST IMMEDIATELY:**
1. Acknowledge your INTEGRATE_WAVE_EFFORTS_DIR from the prompt
2. Set INTEGRATE_WAVE_EFFORTS_DIR environment variable
3. Verify you're in the correct directory
4. Read the MERGE PLAN (.software-factory/WAVE-MERGE-PLAN--${TIMESTAMP}.md or .software-factory/PHASE-MERGE-PLAN--${TIMESTAMP}.md)

### STARTUP ACKNOWLEDGMENT (REQUIRED)
```bash
echo "════════════════════════════════════"
echo "🔧 INTEGRATE_WAVE_EFFORTS AGENT STARTUP"
echo "════════════════════════════════════"
echo "INTEGRATE_WAVE_EFFORTS_DIR: ${INTEGRATE_WAVE_EFFORTS_DIR}"
echo "Current Directory: $(pwd)"
echo "Git Branch: $(git branch --show-current)"

# VERIFY CORRECT LOCATION
if [[ "$(pwd)" != *"$INTEGRATE_WAVE_EFFORTS_DIR"* ]]; then
    echo "❌ WRONG DIRECTORY!"
    exit 1
fi

# SET ENVIRONMENT
export INTEGRATE_WAVE_EFFORTS_DIR="${INTEGRATE_WAVE_EFFORTS_DIR}"
echo "✅ INTEGRATE_WAVE_EFFORTS_DIR acknowledged and set"

# READ MERGE PLAN
if [ -f ".software-factory/WAVE-MERGE-PLAN--${TIMESTAMP}.md" ]; then
    echo "✅ Found .software-factory/WAVE-MERGE-PLAN--${TIMESTAMP}.md"
    MERGE_PLAN=".software-factory/WAVE-MERGE-PLAN--${TIMESTAMP}.md"
elif [ -f ".software-factory/PHASE-MERGE-PLAN--${TIMESTAMP}.md" ]; then
    echo "✅ Found .software-factory/PHASE-MERGE-PLAN--${TIMESTAMP}.md"
    MERGE_PLAN=".software-factory/PHASE-MERGE-PLAN--${TIMESTAMP}.md"
else
    echo "❌ NO MERGE PLAN FOUND!"
    exit 1
fi

echo "📋 Merge plan: $MERGE_PLAN"
```

## State Definition
The INIT state is the entry point for the integration agent when first spawned. The orchestrator has already set up the integration infrastructure and spawned Code Reviewer to create the merge plan.

## SF 3.0 Architecture Context

On initialization, the Integration Agent:
- Reads `state_machine.current_state` from orchestrator-state-v3.json to understand orchestrator's state
- Loads integration container configuration from `integration-containers.json` for current iteration
- Verifies merge plan location tracked in orchestrator-state-v3.json per R340
- Understands integration scope (wave/phase/project level) from state file context

State transitions are managed through orchestrator-state-v3.json updates per R288. The Integration Agent operates within the integration container framework defined in SF 3.0 architecture Part 5.

## Required Actions

### 1. INTEGRATE_WAVE_EFFORTS_DIR Acknowledgment (CRITICAL)
```bash
# Extract INTEGRATE_WAVE_EFFORTS_DIR from prompt or environment
# This is passed by orchestrator in spawn command
echo "🎯 Acknowledging INTEGRATE_WAVE_EFFORTS_DIR: ${INTEGRATE_WAVE_EFFORTS_DIR}"

# Verify we're in the right place
if [[ "$(pwd)" != "$INTEGRATE_WAVE_EFFORTS_DIR" ]]; then
    cd "$INTEGRATE_WAVE_EFFORTS_DIR" || {
        echo "❌ Cannot access INTEGRATE_WAVE_EFFORTS_DIR: $INTEGRATE_WAVE_EFFORTS_DIR"
        exit 1
    }
fi

# Confirm integration branch is checked out
CURRENT_BRANCH=$(git branch --show-current)
if [[ ! "$CURRENT_BRANCH" =~ integration ]]; then
    echo "❌ Not on integration branch! Current: $CURRENT_BRANCH"
    exit 1
fi
```

### 2. Merge Plan Verification
```bash
# Verify merge plan exists and is valid
validate_merge_plan() {
    local plan="$1"
    
    echo "📖 Reading merge plan: $plan"
    
    # Check required sections
    for section in "Target Integration Branch" "Branches to Merge" "Validation Steps"; do
        if ! grep -q "## $section" "$plan"; then
            echo "❌ Missing section in merge plan: $section"
            return 1
        fi
    done
    
    # Count merge operations
    MERGE_COUNT=$(grep -c "git merge origin/" "$plan")
    echo "📊 Total merges to execute: $MERGE_COUNT"
    
    if [[ $MERGE_COUNT -eq 0 ]]; then
        echo "❌ No merge commands found in plan!"
        return 1
    fi
    
    echo "✅ Merge plan validated successfully"
    return 0
}

validate_merge_plan "$MERGE_PLAN"
```

### 3. Rule Acknowledgment
The agent MUST acknowledge:
- R260 - Integration Agent Core Requirements (Git expertise)
- R261 - Integration Planning Requirements (Follow merge plan)
- R262 - Merge Operation Protocols (SUPREME - preserve history)
- R263 - Integration Documentation Requirements (.software-factory/work-log--${TIMESTAMP}.log)
- R267 - Integration Agent Grading Criteria (50/50 split)

### 4. Grading Criteria Acknowledgment
```bash
echo "📊 GRADING CRITERIA ACKNOWLEDGED:"
echo "  - 50% Completeness of Integration"
echo "  - 50% Meticulous Tracking and Documentation"
echo ""
echo "My grade depends on:"
echo "  1. Successfully executing ALL merges in the plan"
echo "  2. Creating comprehensive .software-factory/work-log--${TIMESTAMP}.log"
echo "  3. Generating detailed INTEGRATE_WAVE_EFFORTS-REPORT.md"
echo "  4. Preserving complete git history"
echo "  5. Documenting all conflicts and resolutions"
```

## Transition Rules
- Can transition to: MERGING (to execute merge plan)
- Cannot skip directly to: TESTING, REPORTING
- Must complete INTEGRATE_WAVE_EFFORTS_DIR acknowledgment before transition
- Must verify merge plan exists before transition

## Success Criteria
- ✅ INTEGRATE_WAVE_EFFORTS_DIR acknowledged and verified
- ✅ Correct directory confirmed (pwd matches INTEGRATE_WAVE_EFFORTS_DIR)
- ✅ Integration branch checked out
- ✅ Merge plan found and validated
- ✅ All core rules acknowledged
- ✅ Grading criteria understood
- ✅ Ready to begin planning phase

---
### 🚨🚨🚨 RULE R260 - Integration Agent INTEGRATE_WAVE_EFFORTS_DIR Acknowledgment
**Source:** rule-library/R260-integration-agent-core-requirements.md
**Criticality:** BLOCKING - Must acknowledge directory

Integration Agent MUST acknowledge INTEGRATE_WAVE_EFFORTS_DIR, verify location, and set environment variable.
---

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

