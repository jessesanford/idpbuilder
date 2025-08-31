# Integration Agent - INIT State Rules

## 🚨 MANDATORY STARTUP ACKNOWLEDGMENT 🚨

### IMMEDIATE ACTIONS UPON STARTUP

**YOU MUST IMMEDIATELY:**
1. Acknowledge your INTEGRATION_DIR from the prompt
2. Set INTEGRATION_DIR environment variable
3. Verify you're in the correct directory
4. Read the MERGE PLAN (WAVE-MERGE-PLAN.md or PHASE-MERGE-PLAN.md)

### STARTUP ACKNOWLEDGMENT (REQUIRED)
```bash
echo "════════════════════════════════════"
echo "🔧 INTEGRATION AGENT STARTUP"
echo "════════════════════════════════════"
echo "INTEGRATION_DIR: ${INTEGRATION_DIR}"
echo "Current Directory: $(pwd)"
echo "Git Branch: $(git branch --show-current)"

# VERIFY CORRECT LOCATION
if [[ "$(pwd)" != *"$INTEGRATION_DIR"* ]]; then
    echo "❌ WRONG DIRECTORY!"
    exit 1
fi

# SET ENVIRONMENT
export INTEGRATION_DIR="${INTEGRATION_DIR}"
echo "✅ INTEGRATION_DIR acknowledged and set"

# READ MERGE PLAN
if [ -f "WAVE-MERGE-PLAN.md" ]; then
    echo "✅ Found WAVE-MERGE-PLAN.md"
    MERGE_PLAN="WAVE-MERGE-PLAN.md"
elif [ -f "PHASE-MERGE-PLAN.md" ]; then
    echo "✅ Found PHASE-MERGE-PLAN.md"
    MERGE_PLAN="PHASE-MERGE-PLAN.md"
else
    echo "❌ NO MERGE PLAN FOUND!"
    exit 1
fi

echo "📋 Merge plan: $MERGE_PLAN"
```

## State Definition
The INIT state is the entry point for the integration agent when first spawned. The orchestrator has already set up the integration infrastructure and spawned Code Reviewer to create the merge plan.

## Required Actions

### 1. INTEGRATION_DIR Acknowledgment (CRITICAL)
```bash
# Extract INTEGRATION_DIR from prompt or environment
# This is passed by orchestrator in spawn command
echo "🎯 Acknowledging INTEGRATION_DIR: ${INTEGRATION_DIR}"

# Verify we're in the right place
if [[ "$(pwd)" != "$INTEGRATION_DIR" ]]; then
    cd "$INTEGRATION_DIR" || {
        echo "❌ Cannot access INTEGRATION_DIR: $INTEGRATION_DIR"
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
- R263 - Integration Documentation Requirements (work-log.md)
- R267 - Integration Agent Grading Criteria (50/50 split)

### 4. Grading Criteria Acknowledgment
```bash
echo "📊 GRADING CRITERIA ACKNOWLEDGED:"
echo "  - 50% Completeness of Integration"
echo "  - 50% Meticulous Tracking and Documentation"
echo ""
echo "My grade depends on:"
echo "  1. Successfully executing ALL merges in the plan"
echo "  2. Creating comprehensive work-log.md"
echo "  3. Generating detailed INTEGRATION-REPORT.md"
echo "  4. Preserving complete git history"
echo "  5. Documenting all conflicts and resolutions"
```

## Transition Rules
- Can transition to: PLANNING (to analyze merge plan)
- Cannot skip directly to: MERGING, TESTING, REPORTING
- Must complete INTEGRATION_DIR acknowledgment before transition
- Must verify merge plan exists before transition

## Success Criteria
- ✅ INTEGRATION_DIR acknowledged and verified
- ✅ Correct directory confirmed (pwd matches INTEGRATION_DIR)
- ✅ Integration branch checked out
- ✅ Merge plan found and validated
- ✅ All core rules acknowledged
- ✅ Grading criteria understood
- ✅ Ready to begin planning phase

---
### 🚨🚨🚨 RULE R260 - Integration Agent INTEGRATION_DIR Acknowledgment
**Source:** rule-library/R260-integration-agent-integration-dir.md
**Criticality:** BLOCKING - Must acknowledge directory

Integration Agent MUST acknowledge INTEGRATION_DIR, verify location, and set environment variable.
---