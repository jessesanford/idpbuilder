# WAITING_FOR_MASTER_ARCHITECTURE State Rules


## 🔴🔴🔴 STATE MANAGER BOOKEND REMINDER 🔴🔴🔴

**CRITICAL:** This state's completion checklist shows jq/yq commands to update state files.
**THESE ARE PROHIBITED IN SF 3.0!**

**MANDATORY PATTERN:**
```bash
# WRONG (SF 2.0 - DO NOT USE):
# Direct jq manipulation of orchestrator-state-v3.json is PROHIBITED
# Example: jq update to .state_machine.current_state field

# CORRECT (SF 3.0 - ALWAYS USE):
# Let State Manager handle it via variables:
NEXT_STATE="SPAWN_ARCHITECT_PHASE_PLANNING"
TRANSITION_REASON="Master architecture validated"
# Then use /continue-software-factory which calls State Manager
```

**YOU MUST:**
1. Set NEXT_STATE and TRANSITION_REASON variables
2. Complete state work
3. Set CONTINUE-SOFTWARE-FACTORY=TRUE
4. Exit and let State Manager update state file

**See:** $CLAUDE_PROJECT_DIR/rule-library/R600-state-manager-bookend-protocol.md

---

## State Purpose
Monitor for architect completion of master architecture plan per R210. Active monitoring checkpoint (R233) with validation (R340) before proceeding.

## Entry Conditions
- Architect spawned from SPAWN_ARCHITECT_MASTER_PLANNING
- Master architecture planning in progress
- orchestrator-state-v3.json shows state_machine.current_state as this state

## Monitoring Protocol

### 1. Poll for Master Architecture Plan (R550 Compliant)
```bash
echo "Waiting for Master Architecture from architect..."

# R550: Get expected file location from planning_files (NOT legacy planning_artifacts)
ARCH_PLAN_FILE=$(jq -r '.planning_files.project.architecture_plan // "planning/project/PROJECT-ARCHITECTURE-PLAN.md"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
ARCH_PLAN_PATH="$CLAUDE_PROJECT_DIR/$ARCH_PLAN_FILE"

echo "Expected R550 location: $ARCH_PLAN_PATH"
echo "  (from orchestrator-state-v3.json planning_files.project.architecture_plan)"

# Ensure planning directory exists (R550 standard structure)
mkdir -p "$(dirname "$ARCH_PLAN_PATH")"

# Poll every 30 seconds for up to 30 minutes
MAX_WAIT=1800  # 30 minutes
POLL_INTERVAL=30
ELAPSED=0

while [ $ELAPSED -lt $MAX_WAIT ]; do
    echo "Checking for project architecture plan at R550 location... (${ELAPSED}s / ${MAX_WAIT}s)"

    # Check if file exists at the R550 standardized location
    if [ -f "$ARCH_PLAN_PATH" ]; then
        echo "✅ Project architecture plan detected!"
        echo "  Location: $ARCH_PLAN_PATH"
        ARCH_PLAN="$ARCH_PLAN_PATH"
        break
    fi

    echo "Still waiting for architecture plan at: $ARCH_PLAN_PATH"

    sleep $POLL_INTERVAL
    ELAPSED=$((ELAPSED + POLL_INTERVAL))
done

# Check if we timed out
if [ $ELAPSED -ge $MAX_WAIT ]; then
    echo "❌ ERROR: Timeout waiting for project architecture plan after ${MAX_WAIT} seconds"
    echo "  Expected R550 location: $ARCH_PLAN_PATH"
    NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Timeout waiting for project architecture plan"
    exit 1
fi
```

### 2. Validate Master Architecture Plan (R340)
```bash
echo "═══════════════════════════════════════════════════════"
echo "R340: Validating master architecture plan..."
echo "═══════════════════════════════════════════════════════"

# Validate master architecture plan structure
validate_master_architecture_plan() {
    local plan_file="$1"

    echo "Validating master architecture plan: $(basename "$plan_file")"

    # Check required sections
    local required_sections=(
        "Project Vision"
        "System Architecture"
        "Core Components"
        "Technology Stack"
        "Phase Breakdown"
        "Success Criteria"
    )

    local missing_sections=()
    for section in "${required_sections[@]}"; do
        if ! grep -qE "^#.*${section}" "$plan_file"; then
            missing_sections+=("$section")
        fi
    done

    if [ ${#missing_sections[@]} -gt 0 ]; then
        echo "ERROR: Master architecture plan missing required sections:"
        printf '  - %s\n' "${missing_sections[@]}"
        return 1
    fi

    echo "✅ Master architecture plan validation passed"
    return 0
}

# Run validation
if ! validate_master_architecture_plan "$ARCH_PLAN"; then
    echo "❌ ERROR: Project architecture plan validation failed"
    NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Invalid project architecture plan"
    exit 1
fi
```

### 3. Update State File Planning Files (R550 Pattern)
```bash
echo "Updating orchestrator-state-v3.json with project architecture plan completion status..."

# R550: Update planning_files.project.architecture_plan (NEW standard)
# NOTE: planning_artifacts is deprecated but can still be updated for backward compatibility if needed
jq --arg plan "$ARCH_PLAN_FILE" \
   '.planning_files.project.architecture_plan = $plan' \
   "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" > "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.tmp" && \
   mv "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.tmp" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

echo "✅ R550: Project architecture plan tracked in state file"
echo "  Location: $ARCH_PLAN_FILE"
echo "  Field: planning_files.project.architecture_plan"
```

### 4. Extract Phase Information
```bash
echo "Extracting phase information from project architecture plan..."

# Parse phases from project architecture plan
PHASE_COUNT=$(grep -c "^##.*Phase [0-9]" "$ARCH_PLAN" || echo "0")

if [ "$PHASE_COUNT" -eq 0 ]; then
    echo "⚠️ WARNING: No phases explicitly defined in architecture plan"
    echo "Defaulting to single phase"
    PHASE_COUNT=1
fi

echo "✅ Project architecture plan defines ${PHASE_COUNT} phases"

# Update metadata with phase count
jq --arg count "$PHASE_COUNT" \
   '.project_progression.total_phases = ($count | tonumber)' \
   orchestrator-state-v3.json > orchestrator-state-v3.tmp && \
   mv orchestrator-state-v3.tmp orchestrator-state-v3.json
```

## Exit Conditions
- Project architecture plan exists and is validated per R340 and R550
- orchestrator-state-v3.json updated with plan location in planning_files.project.architecture_plan
- Phase count extracted and stored

## State Transitions
- **SPAWN_ARCHITECT_PHASE_PLANNING**: When master architecture validated successfully
- **ERROR_RECOVERY**: When timeout occurs or validation fails

## Automation Flag
```bash
# After successful validation and transition:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_MASTER_ARCHITECTURE:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Determine Next State
```bash
# Based on state work results, determine next state
NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[BRIEF_REASON_FOR_TRANSITION]"
echo "Next state determined: $NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Commit Work Products (if any)
```bash
# Commit any deliverables BEFORE state transition
# (Plans, reports, configurations, etc.)
git add [work-products]
git commit -m "feat/doc: [description of work products]"
git push
echo "✅ Work products committed"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "WAITING_FOR_MASTER_ARCHITECTURE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_MASTER_ARCHITECTURE complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Verify State Transition Variables Set
```bash
# Verify required variables are set for State Manager
if [ -z "$NEXT_STATE" ] || [ -z "$TRANSITION_REASON" ]; then
    echo "❌ FATAL: State transition variables not set!"
    echo "NEXT_STATE='$NEXT_STATE'"
    echo "TRANSITION_REASON='$TRANSITION_REASON'"
    exit 1
fi
echo "✅ State transition variables verified"
echo "   NEXT_STATE=$NEXT_STATE"
echo "   TRANSITION_REASON=$TRANSITION_REASON"
```

---

### ✅ Step 6: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 7: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-software-factory to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**SF 3.0 uses 7-step exit (not 8-step) - State Manager handles state file updates!**

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever (variables not set)
- Missing Step 3: Work products not committed = lost deliverables
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: Variables not verified = silent failures
- Missing Step 6: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 7: No exit = R322 violation (-100%)

**ALL 7 STEPS ARE MANDATORY - NO EXCEPTIONS**

**ELIMINATED STEPS FROM SF 2.0:**
- ❌ Step 3 (old): Manual jq/yq state updates → State Manager handles this
- ❌ Step 4 (old): Manual state validation → State Manager validates
- ❌ Step 5 (old): Manual state commit → State Manager commits

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required
