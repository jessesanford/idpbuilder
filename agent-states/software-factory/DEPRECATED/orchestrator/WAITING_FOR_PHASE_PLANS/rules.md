# WAITING_FOR_PHASE_PLANS State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Purpose
Monitor for architect completion of phase-level architecture and implementation plans per R210.

## Entry Conditions
- Architect spawned from SPAWN_ARCHITECT_PHASE_PLANNING
- Phase planning in progress
- orchestrator-state-v3.json shows phase status as "planning"

## Monitoring Protocol

### 1. Poll for Phase Plans
```bash
# Get current phase
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
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
    update_state "ERROR_RECOVERY" "Timeout waiting for phase plans"
    exit 1
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
    update_state "ERROR_RECOVERY" "Invalid phase architecture plan"
    exit 1
fi

if ! validate_implementation_plan "$IMPL_PLAN"; then
    echo "ERROR: Implementation plan validation failed"
    update_state "ERROR_RECOVERY" "Invalid phase implementation plan"
    exit 1
fi
```

### 3. Update State File
```bash
echo "Updating orchestrator-state-v3.json with phase plan locations..."

# Update phase metadata with plan locations
jq --arg arch "$ARCH_PLAN" \
   --arg impl "$IMPL_PLAN" \
   --arg phase "$PHASE" \
   '.phases["phase_" + $phase].architecture_plan = $arch |
    .phases["phase_" + $phase].implementation_plan = $impl |
    .phases["phase_" + $phase].status = "ready" |
    .phases["phase_" + $phase].planning_complete = (now | todate)' \
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
   '.phases["phase_" + $phase].total_waves = ($count | tonumber)' \
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

    update_state "ERROR_RECOVERY" "$error_type: $error_details"
}
```

## Success Transition
```bash
echo "═══════════════════════════════════════════════════════"
echo "✅ Phase ${PHASE} planning complete!"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "Phase Plans:"
echo "  Architecture: $(basename "$ARCH_PLAN")"
echo "  Implementation: $(basename "$IMPL_PLAN")"
echo ""
echo "Transitioning to WAVE_START to begin Wave 1 of Phase ${PHASE}"

# Update current wave to 1
jq '.current_wave = 1' orchestrator-state-v3.json > orchestrator-state.tmp && \
    mv orchestrator-state.tmp orchestrator-state-v3.json

# Transition to WAVE_START
update_state "WAVE_START" "Phase ${PHASE} planning complete, starting Wave 1"
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
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
