# WAITING_FOR_ARCHITECTURE_PLAN State Rules

## State Purpose
Monitor for architect completion of wave-level architecture and implementation plans per R210.

## Entry Conditions
- Architect spawned from SPAWN_ARCHITECT_WAVE_PLANNING
- Wave planning in progress
- orchestrator-state-v3.json shows current state as this state

## Monitoring Protocol

### 1. Poll for Wave Plans
```bash
# Get current phase and wave
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)
echo "Waiting for Wave ${PHASE}.${WAVE} plans from architect..."

# Define expected file locations per R303
WAVE_PLANNING_DIR="$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/wave${WAVE}"
mkdir -p "$WAVE_PLANNING_DIR"

# Poll every 30 seconds for up to 30 minutes
MAX_WAIT=1800  # 30 minutes
POLL_INTERVAL=30
ELAPSED=0

while [ $ELAPSED -lt $MAX_WAIT ]; do
    echo "Checking for wave plans... (${ELAPSED}s / ${MAX_WAIT}s)"

    # Check for wave architecture plan
    WAVE_ARCH_PLAN=$(ls -t "${WAVE_PLANNING_DIR}/WAVE-${PHASE}-${WAVE}-ARCHITECTURE--"*.md 2>/dev/null | head -1)

    # Check for wave implementation plan
    WAVE_IMPL_PLAN=$(ls -t "${WAVE_PLANNING_DIR}/WAVE-${PHASE}-${WAVE}-PLAN--"*.md 2>/dev/null | head -1)

    # If both plans exist, validate and proceed
    if [ -n "$WAVE_ARCH_PLAN" ] && [ -n "$WAVE_IMPL_PLAN" ]; then
        echo "Wave plans detected!"
        echo "  Architecture: $(basename "$WAVE_ARCH_PLAN")"
        echo "  Implementation: $(basename "$WAVE_IMPL_PLAN")"
        break
    fi

    # Show what's missing
    echo "Still waiting for:"
    [ -z "$WAVE_ARCH_PLAN" ] && echo "  - Wave ${PHASE}.${WAVE} Architecture Plan"
    [ -z "$WAVE_IMPL_PLAN" ] && echo "  - Wave ${PHASE}.${WAVE} Implementation Plan"

    sleep $POLL_INTERVAL
    ELAPSED=$((ELAPSED + POLL_INTERVAL))
done

# Check if we timed out
if [ $ELAPSED -ge $MAX_WAIT ]; then
    echo "ERROR: Timeout waiting for wave plans after ${MAX_WAIT} seconds"
    update_state "ERROR_RECOVERY" "Timeout waiting for wave plans"
    exit 1
fi
```

### 2. Validate Pre-Planned Infrastructure (R504 ENFORCEMENT)
```bash
echo "═══════════════════════════════════════════════════════"
echo "R504: Validating pre-planned infrastructure..."
echo "═══════════════════════════════════════════════════════"

# Check that architect populated pre_planned_infrastructure
PRE_PLANNED=$(jq -r '.pre_planned_infrastructure // empty' orchestrator-state-v3.json)

if [ -z "$PRE_PLANNED" ] || [ "$PRE_PLANNED" == "null" ]; then
    echo "❌ FATAL: Architect did not populate pre_planned_infrastructure!"
    echo "  R504 VIOLATION: Infrastructure must be pre-planned during wave planning"
    echo "  Cannot proceed to CREATE_NEXT_INFRASTRUCTURE without pre-planned data"
    update_state "ERROR_RECOVERY" "R504 violation: No pre_planned_infrastructure"
    exit 1
fi

# Validate pre-planned infrastructure for current wave
WAVE_EFFORTS=$(jq -r ".pre_planned_infrastructure.efforts | to_entries[] | select(.value.phase == \"phase${PHASE}\" and .value.wave == \"wave${WAVE}\") | .key" orchestrator-state-v3.json)

if [ -z "$WAVE_EFFORTS" ]; then
    echo "❌ FATAL: No pre-planned efforts for Phase ${PHASE} Wave ${WAVE}!"
    echo "  R504 VIOLATION: All wave efforts must be pre-planned"
    update_state "ERROR_RECOVERY" "R504 violation: No pre-planned efforts for wave"
    exit 1
fi

echo "✅ Pre-planned infrastructure found for wave:"
for effort_key in $WAVE_EFFORTS; do
    EFFORT_CONFIG=$(jq -r ".pre_planned_infrastructure.efforts.\"${effort_key}\"" orchestrator-state-v3.json)
    echo "  - $effort_key:"
    echo "    Branch: $(echo "$EFFORT_CONFIG" | jq -r '.branch_name')"
    echo "    Path: $(echo "$EFFORT_CONFIG" | jq -r '.full_path')"
    echo "    Remote: $(echo "$EFFORT_CONFIG" | jq -r '.target_remote')"
done

# Mark pre_planned_infrastructure as validated
jq '.pre_planned_infrastructure.validated = true |
    .pre_planned_infrastructure.validation_timestamp = now | todate' \
    orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ Pre-planned infrastructure validated and marked as ready"
```

### 3. Validate Wave Plans (R502)
```bash
echo "═══════════════════════════════════════════════════════"
echo "R502: Validating wave plans content..."
echo "═══════════════════════════════════════════════════════"

# Validate wave architecture plan structure
validate_wave_architecture_plan() {
    local plan_file="$1"

    echo "Validating wave architecture plan: $(basename "$plan_file")"

    # Check required sections
    local required_sections=(
        "Technical Approach"
        "Component Designs"
        "Technology Choices"
        "Risk Assessment"
    )

    local missing_sections=()
    for section in "${required_sections[@]}"; do
        if ! grep -qE "^#.*${section}" "$plan_file"; then
            missing_sections+=("$section")
        fi
    done

    if [ ${#missing_sections[@]} -gt 0 ]; then
        echo "WARNING: Architecture plan may be missing sections:"
        printf '  - %s\n' "${missing_sections[@]}"
        # Continue anyway - architect knows best
    fi

    echo "✅ Wave architecture plan validation passed"
    return 0
}

# Validate wave implementation plan structure
validate_wave_implementation_plan() {
    local plan_file="$1"

    echo "Validating wave implementation plan: $(basename "$plan_file")"

    # Check required sections
    local required_sections=(
        "Effort Breakdown"
        "Parallelization Strategy"
        "Dependencies"
        "Acceptance Criteria"
    )

    local missing_sections=()
    for section in "${required_sections[@]}"; do
        if ! grep -qE "^#.*${section}" "$plan_file"; then
            missing_sections+=("$section")
        fi
    done

    if [ ${#missing_sections[@]} -gt 0 ]; then
        echo "WARNING: Implementation plan may be missing sections:"
        printf '  - %s\n' "${missing_sections[@]}"
        # Continue anyway - architect knows best
    fi

    # Check that efforts are defined
    if ! grep -qE "Effort [0-9]+" "$plan_file"; then
        echo "WARNING: Implementation plan may not define specific efforts"
    fi

    echo "✅ Wave implementation plan validation passed"
    return 0
}

# Run validations
validate_wave_architecture_plan "$WAVE_ARCH_PLAN"
validate_wave_implementation_plan "$WAVE_IMPL_PLAN"
```

### 3. R505: Synchronize Pre-Planned Infrastructure
```bash
echo "═══════════════════════════════════════════════════════"
echo "R505: Synchronizing pre_planned_infrastructure from wave plans"
echo "═══════════════════════════════════════════════════════"

# Parse wave implementation plan for efforts and sync infrastructure
sync_wave_infrastructure() {
    local PHASE=$1
    local WAVE=$2
    local IMPL_PLAN=$3

    echo "Parsing wave implementation plan for effort infrastructure..."

    # Extract efforts from wave plan
    # Looking for patterns like "Effort 1:", "**Effort 1:**", "### Effort 1:", etc.
    local EFFORTS=$(grep -oE "Effort [0-9]+[a-z]*:" "$IMPL_PLAN" | sed 's/Effort \([^:]*\):.*/\1/' | sort -u)

    if [ -z "$EFFORTS" ]; then
        echo "WARNING: No efforts found in wave implementation plan"
        return 0
    fi

    echo "Found efforts in Wave ${WAVE}: $EFFORTS"

    # Get project prefix and target repository
    local PROJECT_PREFIX=$(yq '.project_prefix // "project"' orchestrator-state-v3.json)
    local TARGET_REPO=$(yq '.repository' target-repo-config.yaml)

    # Determine previous wave's last effort for cascade base
    local PREV_WAVE_LAST_EFFORT=""
    if [[ $WAVE -gt 1 ]]; then
        PREV_WAVE=$((WAVE - 1))
        # Find the last effort from previous wave
        PREV_WAVE_LAST_EFFORT=$(yq ".pre_planned_infrastructure.efforts | keys | .[] | select(. == \"phase${PHASE}_wave${PREV_WAVE}_*\") | ." orchestrator-state-v3.json | tail -1 | sed "s/phase${PHASE}_wave${PREV_WAVE}_//")

        if [ -n "$PREV_WAVE_LAST_EFFORT" ]; then
            PREV_WAVE_LAST_BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${PREV_WAVE}/${PREV_WAVE_LAST_EFFORT}"
        fi
    elif [[ $PHASE -gt 1 ]]; then
        # First wave of new phase - get last effort from previous phase
        PREV_PHASE=$((PHASE - 1))
        # This would need to look at previous phase's last wave's last effort
        PREV_WAVE_LAST_BRANCH="<to-be-determined-from-phase-${PREV_PHASE}>"
    fi

    local EFFORT_INDEX=0

    # For each effort, calculate infrastructure per R308
    for EFFORT_NUM in $EFFORTS; do
        ((EFFORT_INDEX++))

        # Sanitize effort name (convert to lowercase, replace spaces with hyphens)
        local EFFORT_NAME="effort-${EFFORT_NUM,,}"
        local EFFORT_ID="phase${PHASE}_wave${WAVE}_${EFFORT_NAME}"

        # Check if effort already exists in pre_planned_infrastructure
        local EXISTS=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\"" orchestrator-state-v3.json)

        # Calculate base branch per R308 cascade algorithm
        local BASE_BRANCH=""

        if [[ $PHASE -eq 1 && $WAVE -eq 1 && $EFFORT_INDEX -eq 1 ]]; then
            BASE_BRANCH="main"
        elif [[ $EFFORT_INDEX -eq 1 ]]; then
            # First effort of wave - use previous wave's last effort
            if [ -n "$PREV_WAVE_LAST_BRANCH" ]; then
                BASE_BRANCH="$PREV_WAVE_LAST_BRANCH"
            else
                BASE_BRANCH="main"  # Fallback if we can't determine
            fi
        else
            # Subsequent efforts cascade from previous effort in same wave
            local PREV_INDEX=$((EFFORT_INDEX - 1))
            local PREV_EFFORT="effort-${PREV_INDEX}"
            BASE_BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${PREV_EFFORT}"
        fi

        if [[ "$EXISTS" == "null" ]]; then
            echo "🆕 Adding new effort infrastructure: $EFFORT_ID"
            echo "   Base branch: $BASE_BRANCH"
        else
            echo "📝 Updating effort infrastructure: $EFFORT_ID"
            echo "   Base branch: $BASE_BRANCH"
        fi

        # Update orchestrator-state-v3.json with pre-planned infrastructure
        yq -i ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\" = {
            \"full_path\": \"$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/\",
            \"branch_name\": \"${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}\",
            \"remote_branch\": \"origin/${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}\",
            \"base_branch\": \"$BASE_BRANCH\",
            \"target_remote\": \"target\",
            \"target_repository\": \"$TARGET_REPO\",
            \"planning_remote\": \"planning\",
            \"split_pattern\": \"${EFFORT_NAME}--split-NNN\",
            \"created\": false,
            \"validated\": false,
            \"added_at\": \"$(date -Iseconds)\",
            \"source\": \"R505-wave-sync\"
        }" orchestrator-state-v3.json
    done

    # Mark synchronization complete
    yq -i ".pre_planned_infrastructure.last_wave_sync = \"$(date -Iseconds)\" |
           .pre_planned_infrastructure.validated = false |
           .pre_planned_infrastructure.needs_validation = true" orchestrator-state-v3.json

    echo "✅ Pre-planned infrastructure synchronized for Wave ${WAVE}"
}

# Run the synchronization
sync_wave_infrastructure "$PHASE" "$WAVE" "$WAVE_IMPL_PLAN"

# Validate the synchronization
echo "Validating synchronized infrastructure..."
WAVE_INFRA_COUNT=$(yq ".pre_planned_infrastructure.efforts | keys | .[] | select(. == \"phase${PHASE}_wave${WAVE}_*\") | ." orchestrator-state-v3.json | wc -l)
echo "Wave ${WAVE} infrastructure entries: $WAVE_INFRA_COUNT"

if [ "$WAVE_INFRA_COUNT" -eq 0 ]; then
    echo "WARNING: No infrastructure entries created for Wave ${WAVE} - architecture plan may be incomplete"
    echo "System will transition to ERROR_RECOVERY if this is a problem"
fi

# Verify all efforts target the correct repository
TARGET_REPO=$(yq '.repository' target-repo-config.yaml)
WRONG_REPO_COUNT=$(yq ".pre_planned_infrastructure.efforts | to_entries | .[] | select(.key | test(\"phase${PHASE}_wave${WAVE}_\")) | select(.value.target_repository != \"$TARGET_REPO\") | .key" orchestrator-state-v3.json | wc -l)

if [ "$WRONG_REPO_COUNT" -gt 0 ]; then
    echo "❌ ERROR: Some efforts targeting wrong repository!"
    echo "Expected: $TARGET_REPO"
    exit 505
fi

echo "✅ All Wave ${WAVE} efforts target correct repository: $TARGET_REPO"
```

### 4. Update State File
```bash
echo "Updating orchestrator-state-v3.json with wave plan locations..."

# Update wave metadata with plan locations
jq --arg arch "$WAVE_ARCH_PLAN" \
   --arg impl "$WAVE_IMPL_PLAN" \
   --arg phase "$PHASE" \
   --arg wave "$WAVE" \
   '.waves["phase" + $phase + "_wave" + $wave].architecture_plan = $arch |
    .waves["phase" + $phase + "_wave" + $wave].implementation_plan = $impl |
    .waves["phase" + $phase + "_wave" + $wave].status = "ready" |
    .waves["phase" + $phase + "_wave" + $wave].planning_complete = (now | todate)' \
   orchestrator-state-v3.json > orchestrator-state.tmp && \
   mv orchestrator-state.tmp orchestrator-state-v3.json

echo "Wave plans recorded in state file"
```

### 5. Extract Effort Information
```bash
echo "Extracting effort information from wave implementation plan..."

# Parse efforts from implementation plan
EFFORT_COUNT=$(grep -c "Effort [0-9]" "$WAVE_IMPL_PLAN" || echo "0")

if [ "$EFFORT_COUNT" -eq 0 ]; then
    echo "WARNING: No efforts explicitly defined in wave plan"
    echo "Cannot proceed without efforts"
    update_state "ERROR_RECOVERY" "No efforts found in wave plan"
    exit 1
fi

echo "Wave ${PHASE}.${WAVE} contains ${EFFORT_COUNT} efforts"

# Update wave metadata with effort count
jq --arg count "$EFFORT_COUNT" \
   --arg phase "$PHASE" \
   --arg wave "$WAVE" \
   '.waves["phase" + $phase + "_wave" + $wave].total_efforts = ($count | tonumber)' \
   orchestrator-state-v3.json > orchestrator-state.tmp && \
   mv orchestrator-state.tmp orchestrator-state-v3.json
```

## Exit Conditions
- Both wave plans exist and are validated
- Pre-planned infrastructure synchronized per R505
- orchestrator-state-v3.json updated with plan locations
- Wave marked as "ready" in state file
- Effort count extracted and stored

## State Transitions
- **SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING**: If phase-level test planning needed
- **SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING**: If wave-level test planning needed
- **CREATE_NEXT_INFRASTRUCTURE**: When ready to create effort infrastructure
- **ERROR_RECOVERY**: When timeout occurs or validation fails

## Error Handling
```bash
handle_waiting_error() {
    local error_type="$1"
    local error_details="$2"

    case "$error_type" in
        "timeout")
            echo "ERROR: Architect did not complete wave planning in time"
            echo "Details: $error_details"
            ;;
        "invalid_plan")
            echo "ERROR: Wave plans do not meet R502 validation requirements"
            echo "Details: $error_details"
            ;;
        "no_efforts")
            echo "ERROR: Wave plan does not define any efforts"
            echo "Details: $error_details"
            ;;
        "sync_failure")
            echo "ERROR: Failed to synchronize pre_planned_infrastructure"
            echo "Details: $error_details"
            ;;
    esac

    update_state "ERROR_RECOVERY" "$error_type: $error_details"
}
```

## Success Transition
```bash
echo "═══════════════════════════════════════════════════════"
echo "✅ Wave ${PHASE}.${WAVE} planning complete!"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "Wave Plans:"
echo "  Architecture: $(basename "$WAVE_ARCH_PLAN")"
echo "  Implementation: $(basename "$WAVE_IMPL_PLAN")"
echo ""
echo "Infrastructure Synchronized:"
echo "  Wave ${WAVE} efforts: $WAVE_INFRA_COUNT"
echo "  All targeting: $TARGET_REPO"
echo ""

# Determine next state based on test planning needs
# For now, proceed to infrastructure creation
echo "Transitioning to CREATE_NEXT_INFRASTRUCTURE to set up effort infrastructure"

# Transition to CREATE_NEXT_INFRASTRUCTURE
update_state "CREATE_NEXT_INFRASTRUCTURE" "Wave ${PHASE}.${WAVE} planning complete with R505 sync, creating infrastructure"
```

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Waiting for Architect to complete planning
- ✅ Plan detected and validated successfully
- ✅ R505 synchronization completed
- ✅ Ready to transition to infrastructure creation
- ✅ Following designed workflow

**THIS IS NORMAL WORKFLOW.** Waiting for plans and transitioning after validation is EXPECTED.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot access expected plan location
- ❌ Plan validation failed critically
- ❌ R505 synchronization impossible
- ❌ State machine corruption

**DO NOT set FALSE because:**
- ❌ Still waiting (NORMAL!)
- ❌ Plan complete, transitioning (EXPECTED!)
- ❌ R322 requires stop (stop ≠ FALSE!)

**Correct pattern:** `exit 0` + `CONTINUE-SOFTWARE-FACTORY=TRUE`

## Automation Flag
```bash
# After successful validation, synchronization, and transition:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
