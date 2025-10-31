# WAITING_FOR_PHASE_PLANS State Rules


## 🚨 State Manager Bookend Pattern (MANDATORY)

**BEFORE this state**:
- State Manager validated transition via STARTUP_CONSULTATION
- You are here because State Manager directed you here
- orchestrator-state-v3.json shows validated_by: "state-manager"

**DURING this state**:
- Perform state-specific work
- NEVER call update_state directly
- Prepare results for State Manager
- Propose next state (don't decide!)

**AFTER this state**:
- Spawn State Manager SHUTDOWN_CONSULTATION
- Provide results and proposed next state
- State Manager validates and decides actual next state
- Transition to State Manager's required_next_state

**CRITICAL**: The orchestrator PROPOSES, the State Manager DECIDES!

---

## State Purpose
Actively monitor for architect completion of phase-level architecture and implementation plans per R210 and R233.

## Core Mandatory Rules

1. **🔴🔴🔴 R233** - All States Require Immediate Action (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: SUPREME - Violation = AUTOMATIC FAILURE
   - Summary: WAITING states must ACTIVELY poll/monitor - not passively wait
   - **CRITICAL**: "Waiting" means ACTIVELY CHECKING (poll every few seconds, report what you're checking, set timeouts)

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions
   - **CRITICAL**: Commit and push within 60 seconds of saving

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition
   - **CRITICAL**: Commit and push state changes immediately

## Entry Conditions
- Architect spawned from SPAWN_ARCHITECT_PHASE_PLANNING
- Phase planning in progress
- orchestrator-state-v3.json shows phase status as "planning"

## Monitoring Protocol

### 1. Poll for Phase Plans
```bash
# Get current phase
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
echo "Waiting for Phase ${PHASE} plans from architect..."

# Define expected file locations
PLANNING_DIR="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}"

# Poll every 30 seconds for up to 30 minutes
MAX_WAIT=1800  # 30 minutes
POLL_INTERVAL=30
ELAPSED=0

while [ $ELAPSED -lt $MAX_WAIT ]; do
    echo "Checking for phase plans... (${ELAPSED}s / ${MAX_WAIT}s)"

    # Check for architecture plan
    ARCH_PLAN=$(ls -t "${PLANNING_DIR}/PHASE-${PHASE}-ARCHITECTURE-PLAN--"*.md 2>/dev/null | head -1)

    # Check for implementation plan
    IMPL_PLAN=$(ls -t "${PLANNING_DIR}/PHASE-${PHASE}-PLAN--"*.md 2>/dev/null | head -1)

    # If both plans exist, validate and proceed
    if [ -n "$ARCH_PLAN" ] && [ -n "$IMPL_PLAN" ]; then
        echo "Phase plans detected!"
        echo "  Architecture: $(basename "$ARCH_PLAN")"
        echo "  Implementation: $(basename "$IMPL_PLAN")"
        break
    fi

    # Show what's missing
    echo "Still waiting for:"
    [ -z "$ARCH_PLAN" ] && echo "  - Phase ${PHASE} Architecture Plan"
    [ -z "$IMPL_PLAN" ] && echo "  - Phase ${PHASE} Implementation Plan"

    sleep $POLL_INTERVAL
    ELAPSED=$((ELAPSED + POLL_INTERVAL))
done

# Check if we timed out
if [ $ELAPSED -ge $MAX_WAIT ]; then
    echo "ERROR: Timeout waiting for phase plans after ${MAX_WAIT} seconds"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Timeout waiting for phase plans after ${MAX_WAIT}s"
    ERROR_OCCURRED=true
    # NOTE: Will transition via State Manager in completion checklist
fi
```

### 2. Validate Phase Plans (R502)
```bash
echo "═══════════════════════════════════════════════════════"
echo "R502: Validating phase plans..."
echo "═══════════════════════════════════════════════════════"

# Validate architecture plan structure
validate_architecture_plan() {
    local plan_file="$1"

    echo "Validating architecture plan: $(basename "$plan_file")"

    # Check required sections
    local required_sections=(
        "Phase Vision Alignment"
        "Analysis of Previous Phases"
        "Phase.*Architecture"
        "Core Architectural Decisions"
        "APIs and Contracts"
        "Abstractions and Interfaces"
    )

    local missing_sections=()
    for section in "${required_sections[@]}"; do
        if ! grep -qE "^#.*${section}" "$plan_file"; then
            missing_sections+=("$section")
        fi
    done

    if [ ${#missing_sections[@]} -gt 0 ]; then
        echo "ERROR: Architecture plan missing required sections:"
        printf '  - %s\n' "${missing_sections[@]}"
        return 1
    fi

    echo "✅ Architecture plan validation passed"
    return 0
}

# Validate implementation plan structure
validate_implementation_plan() {
    local plan_file="$1"

    echo "Validating implementation plan: $(basename "$plan_file")"

    # Check required sections
    local required_sections=(
        "Phase Overview"
        "Wave Breakdown"
        "Implementation Strategy"
        "Success Criteria"
    )

    local missing_sections=()
    for section in "${required_sections[@]}"; do
        if ! grep -qE "^#.*${section}" "$plan_file"; then
            missing_sections+=("$section")
        fi
    done

    if [ ${#missing_sections[@]} -gt 0 ]; then
        echo "ERROR: Implementation plan missing required sections:"
        printf '  - %s\n' "${missing_sections[@]}"
        return 1
    fi

    # Check that waves are defined
    if ! grep -qE "Wave [0-9]+" "$plan_file"; then
        echo "ERROR: Implementation plan must define waves"
        return 1
    fi

    echo "✅ Implementation plan validation passed"
    return 0
}

# Run validations
if ! validate_architecture_plan "$ARCH_PLAN"; then
    echo "ERROR: Architecture plan validation failed"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Invalid phase architecture plan"
    ERROR_OCCURRED=true
    # NOTE: Will transition via State Manager in completion checklist
fi

if ! validate_implementation_plan "$IMPL_PLAN"; then
    echo "ERROR: Implementation plan validation failed"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Invalid phase implementation plan"
    ERROR_OCCURRED=true
    # NOTE: Will transition via State Manager in completion checklist
fi
```

### 3. Update State File
```bash
echo "Updating orchestrator-state-v3.json with phase plan locations..."

# Update phase metadata with plan locations
jq --arg arch "$ARCH_PLAN" \
   --arg impl "$IMPL_PLAN" \
   --arg phase "$PHASE" \
   '.project_progression.phases["phase_" + $phase].architecture_plan = $arch |
    .project_progression.phases["phase_" + $phase].implementation_plan = $impl |
    .project_progression.phases["phase_" + $phase].status = "ready" |
    .project_progression.phases["phase_" + $phase].planning_complete = (now | todate)' \
   orchestrator-state-v3.json > orchestrator-state.tmp && \
   mv orchestrator-state.tmp orchestrator-state-v3.json

echo "Phase plans recorded in state file"
```

### 4. Extract Wave Information
```bash
echo "Extracting wave information from implementation plan..."

# Parse waves from implementation plan
WAVE_COUNT=$(grep -c "^##.*Wave [0-9]" "$IMPL_PLAN" || echo "0")

if [ "$WAVE_COUNT" -eq 0 ]; then
    echo "WARNING: No waves explicitly defined in implementation plan"
    echo "Defaulting to single wave"
    WAVE_COUNT=1
fi

echo "Phase ${PHASE} contains ${WAVE_COUNT} waves"

# Update phase metadata with wave count
jq --arg count "$WAVE_COUNT" \
   --arg phase "$PHASE" \
   '.project_progression.phases["phase_" + $phase].total_waves = ($count | tonumber)' \
   orchestrator-state-v3.json > orchestrator-state.tmp && \
   mv orchestrator-state.tmp orchestrator-state-v3.json
```

### 5. R505: Synchronize Pre-Planned Infrastructure
```bash
echo "═══════════════════════════════════════════════════════"
echo "R505: Synchronizing pre_planned_infrastructure from phase plans"
echo "═══════════════════════════════════════════════════════"

# Parse implementation plan for all efforts across all waves
parse_and_sync_phase_infrastructure() {
    local PHASE=$1
    local IMPL_PLAN=$2

    echo "Parsing implementation plan for effort infrastructure..."

    # Extract all efforts mentioned in the plan
    # Looking for patterns like "Effort 1:", "**Effort 1:**", "### Effort 1:", etc.
    local EFFORTS=$(grep -oE "Effort [0-9]+[a-z]*:" "$IMPL_PLAN" | sed 's/Effort \([^:]*\):.*/\1/' | sort -u)

    if [ -z "$EFFORTS" ]; then
        echo "WARNING: No efforts found in implementation plan"
        return 0
    fi

    echo "Found efforts to track: $EFFORTS"

    # Get project prefix and target repository
    local PROJECT_PREFIX=$(yq '.project_prefix // "project"' orchestrator-state-v3.json)
    local TARGET_REPO=$(yq '.repository' target-repo-config.yaml)

    # For each effort, calculate infrastructure per R308
    for EFFORT_NUM in $EFFORTS; do
        # Try to determine which wave this effort belongs to
        local WAVE=$(grep -B5 "Effort ${EFFORT_NUM}:" "$IMPL_PLAN" | grep -oE "Wave [0-9]+" | tail -1 | sed 's/Wave //')

        if [ -z "$WAVE" ]; then
            echo "WARNING: Could not determine wave for Effort ${EFFORT_NUM}, skipping..."
            continue
        fi

        # Sanitize effort name (convert to lowercase, replace spaces with hyphens)
        local EFFORT_NAME="effort-${EFFORT_NUM,,}"
        local EFFORT_ID="phase${PHASE}_wave${WAVE}_${EFFORT_NAME}"

        # Calculate base branch per R308 cascade algorithm
        local BASE_BRANCH=""
        local EFFORT_INDEX=$(echo "$EFFORT_NUM" | sed 's/[a-z]*//')  # Remove any suffix letters

        if [[ $PHASE -eq 1 && $WAVE -eq 1 && $EFFORT_INDEX -eq 1 ]]; then
            BASE_BRANCH="main"
        elif [[ $EFFORT_INDEX -eq 1 ]]; then
            # First effort of wave - need to determine previous wave's last effort
            if [[ $WAVE -eq 1 && $PHASE -gt 1 ]]; then
                # First wave of new phase - would cascade from previous phase's last effort
                BASE_BRANCH="<to-be-determined-from-previous-phase>"
            elif [[ $WAVE -gt 1 ]]; then
                # Subsequent waves - cascade from previous wave's last effort
                BASE_BRANCH="<to-be-determined-from-wave-$((WAVE-1))>"
            fi
        else
            # Subsequent efforts cascade from previous effort in same wave
            local PREV_EFFORT="effort-$((EFFORT_INDEX - 1))"
            BASE_BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${PREV_EFFORT}"
        fi

        # Create infrastructure entry
        echo "Adding infrastructure for: $EFFORT_ID"

        # Update orchestrator-state-v3.json with pre-planned infrastructure
        jq --arg id "$EFFORT_ID" \
           --arg path "$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/" \
           --arg branch "${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}" \
           --arg remote "origin/${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}" \
           --arg base "$BASE_BRANCH" \
           --arg repo "$TARGET_REPO" \
           --arg timestamp "$(date -Iseconds)" \
           '.pre_planned_infrastructure.efforts[$id] = {
              "full_path": $path,
              "branch_name": $branch,
              "remote_branch": $remote,
              "base_branch": $base,
              "target_remote": "target",
              "target_repository": $repo,
              "planning_remote": "planning",
              "split_pattern": ($branch + "--split-NNN"),
              "created": false,
              "validated": false,
              "added_at": $timestamp,
              "source": "R505-phase-sync"
           }' orchestrator-state-v3.json > orchestrator-state.tmp && \
           mv orchestrator-state.tmp orchestrator-state-v3.json
    done

    # Mark synchronization complete
    jq --arg timestamp "$(date -Iseconds)" \
       '.pre_planned_infrastructure.last_phase_sync = $timestamp |
        .pre_planned_infrastructure.validated = false |
        .pre_planned_infrastructure.needs_validation = true' \
       orchestrator-state-v3.json > orchestrator-state.tmp && \
       mv orchestrator-state.tmp orchestrator-state-v3.json

    echo "✅ Pre-planned infrastructure synchronized for Phase ${PHASE}"
}

# Run the synchronization
parse_and_sync_phase_infrastructure "$PHASE" "$IMPL_PLAN"

# Validate the synchronization
echo "Validating synchronized infrastructure..."
INFRA_COUNT=$(jq '.pre_planned_infrastructure.efforts | length' orchestrator-state-v3.json)
echo "Total infrastructure entries: $INFRA_COUNT"

if [ "$INFRA_COUNT" -eq 0 ]; then
    echo "WARNING: No infrastructure entries created - architecture plan may be incomplete"
    echo "System will transition to ERROR_RECOVERY if this is a problem"
fi
```

## Exit Conditions
- Both phase plans exist and are validated
- orchestrator-state-v3.json updated with plan locations
- Phase marked as "ready" in state file
- Wave count extracted and stored

## State Transitions
- **WAVE_START**: When both plans validated successfully
- **ERROR_RECOVERY**: When timeout occurs or validation fails

## Error Handling
```bash
handle_waiting_error() {
    local error_type="$1"
    local error_details="$2"

    case "$error_type" in
        "timeout")
            echo "ERROR: Architect did not complete phase planning in time"
            echo "Details: $error_details"
            # Consider spawning a monitor to check on architect
            ;;
        "invalid_plan")
            echo "ERROR: Phase plans do not meet R502 validation requirements"
            echo "Details: $error_details"
            # May need to re-spawn architect with clearer instructions
            ;;
        "partial_completion")
            echo "ERROR: Only one plan was created"
            echo "Details: $error_details"
            # Need both plans or neither
            ;;
    esac

    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="$error_type: $error_details"
    ERROR_OCCURRED=true
    # NOTE: Will transition via State Manager in completion checklist
}
```

## Success Transition Preparation
```bash
echo "═══════════════════════════════════════════════════════"
echo "✅ Phase ${PHASE} planning complete!"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "Phase Plans:"
echo "  Architecture: $(basename "$ARCH_PLAN")"
echo "  Implementation: $(basename "$IMPL_PLAN")"
echo ""
echo "Preparing transition to WAVE_START to begin Wave 1 of Phase ${PHASE}"

# Update current wave to 1 (in preparation for transition)
jq '.project_progression.current_wave.wave_number = 1' orchestrator-state-v3.json > orchestrator-state.tmp && \
    mv orchestrator-state.tmp orchestrator-state-v3.json

# Prepare transition proposal
PROPOSED_NEXT_STATE="WAVE_START"
TRANSITION_REASON="Phase ${PHASE} planning complete, starting Wave 1"
echo "📍 Will propose: $PROPOSED_NEXT_STATE to State Manager"
# NOTE: Actual transition via State Manager in completion checklist
```

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Waiting for Architect to complete phase planning
- ✅ Phase plan detected and validated
- ✅ Wave structure extracted successfully
- ✅ Ready to transition to WAVE_START
- ✅ Following designed workflow

**THIS IS NORMAL WORKFLOW.** Waiting for phase plans and transitioning is EXPECTED.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot access expected plan location
- ❌ Plan validation failed critically
- ❌ Cannot extract wave structure
- ❌ State machine corruption

**DO NOT set FALSE because:**
- ❌ Still waiting (NORMAL!)
- ❌ Plan complete, transitioning (EXPECTED!)
- ❌ R322 requires stop (stop ≠ FALSE!)

**Correct pattern:** `exit 0` + `CONTINUE-SOFTWARE-FACTORY=TRUE`

## Automation Flag
```bash
# After successful validation and transition:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_PHASE_PLANS:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - state complete [R287]"; then
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

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

**CRITICAL**: Steps 2-5 enforce State Manager bookend pattern - orchestrator PROPOSES, State Manager DECIDES!

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
