# 🚨🚨🚨 RULE R505 - Phase/Wave Pre-Planned Infrastructure Synchronization [BLOCKING]
**FAIL = INFRASTRUCTURE MISMATCH** | Source: [R505](RULE-REGISTRY.md#R505)

## PURPOSE
Ensure pre_planned_infrastructure field in orchestrator-state-v3.json remains synchronized and up-to-date throughout the project lifecycle as new phase/wave plans are created. This is an ONGOING synchronization requirement, not just initial population.

## WHEN THIS RULE APPLIES
- **DURING**: SPAWN_ARCHITECT_PHASE_PLANNING state
- **DURING**: SPAWN_ARCHITECT_WAVE_PLANNING state
- **AFTER**: Phase/Wave plans are created
- **BEFORE**: CREATE_NEXT_INFRASTRUCTURE state
- **TRIGGER**: Any time new efforts are added to plans

## MANDATORY SYNCHRONIZATION POINTS

### 1. SPAWN_ARCHITECT_PHASE_PLANNING
**When spawning architect for phase planning:**
```bash
# After architect creates phase plan
PHASE_PLAN="$CLAUDE_PROJECT_DIR/phase-plans/phase-${PHASE}-plan.md"

# Parse and sync infrastructure
sync_phase_infrastructure() {
    local PHASE=$1
    local PLAN_FILE=$2

    # Extract all efforts from phase plan
    local EFFORTS=$(grep -E "^\*\*Effort [0-9]+" "$PLAN_FILE" | sed 's/.*Effort \([^:]*\):.*/\1/')

    # For each effort, calculate and update infrastructure
    for EFFORT in $EFFORTS; do
        calculate_effort_infrastructure "$PHASE" "$EFFORT"
        update_pre_planned_infrastructure "$PHASE" "$EFFORT"
    done

    # Mark as synchronized
    yq -i ".pre_planned_infrastructure.last_phase_sync = \"$(date -Iseconds)\"" orchestrator-state-v3.json
}
```

### 2. SPAWN_ARCHITECT_WAVE_PLANNING
**When spawning architect for wave planning:**
```bash
# After architect creates wave plan
WAVE_PLAN="$CLAUDE_PROJECT_DIR/phase-plans/phase-${PHASE}-wave-${WAVE}-plan.md"

# Parse and sync infrastructure
sync_wave_infrastructure() {
    local PHASE=$1
    local WAVE=$2
    local PLAN_FILE=$3

    # Extract efforts from wave plan
    local EFFORTS=$(parse_wave_efforts "$PLAN_FILE")

    # Check for NEW efforts not in pre_planned_infrastructure
    for EFFORT in $EFFORTS; do
        local EFFORT_ID="phase${PHASE}_wave${WAVE}_${EFFORT}"

        # Check if effort already exists
        local EXISTS=$(yq ".pre_planned_infrastructure.efforts.$EFFORT_ID" orchestrator-state-v3.json)

        if [[ "$EXISTS" == "null" ]]; then
            echo "🆕 New effort detected: $EFFORT_ID"
            add_new_effort_infrastructure "$PHASE" "$WAVE" "$EFFORT"
        else
            echo "✅ Effort already tracked: $EFFORT_ID"
        fi
    done

    # Update wave synchronization timestamp
    yq -i ".pre_planned_infrastructure.last_wave_sync = \"$(date -Iseconds)\"" orchestrator-state-v3.json
}
```

### 3. WAITING_FOR_PHASE_PLANS / WAITING_FOR_ARCHITECTURE_PLAN
**After plans are received:**
```bash
# Orchestrator must sync before transitioning
validate_and_sync_infrastructure() {
    local STATE=$1

    case $STATE in
        WAITING_FOR_PHASE_PLANS)
            # Read the created phase plan
            local PHASE=$(yq '.current_phase' orchestrator-state-v3.json)
            local PLAN_FILE="$CLAUDE_PROJECT_DIR/phase-plans/phase-${PHASE}-plan.md"

            if [[ -f "$PLAN_FILE" ]]; then
                sync_phase_infrastructure "$PHASE" "$PLAN_FILE"
                echo "✅ Phase $PHASE infrastructure synchronized"
            else
                echo "❌ FATAL: Phase plan not found!"
                exit 505
            fi
            ;;

        WAITING_FOR_ARCHITECTURE_PLAN)
            # Read the created wave plan
            local PHASE=$(yq '.current_phase' orchestrator-state-v3.json)
            local WAVE=$(yq '.current_wave' orchestrator-state-v3.json)
            local PLAN_FILE="$CLAUDE_PROJECT_DIR/phase-plans/phase-${PHASE}-wave-${WAVE}-plan.md"

            if [[ -f "$PLAN_FILE" ]]; then
                sync_wave_infrastructure "$PHASE" "$WAVE" "$PLAN_FILE"
                echo "✅ Wave $WAVE infrastructure synchronized"
            else
                echo "❌ FATAL: Wave plan not found!"
                exit 505
            fi
            ;;
    esac
}
```

## INFRASTRUCTURE CALCULATION (PER R308)

### Calculate Complete Infrastructure Entry
```bash
calculate_effort_infrastructure() {
    local PHASE=$1
    local WAVE=$2
    local EFFORT=$3
    local EFFORT_INDEX=$4

    # Determine base branch per R308 cascade algorithm
    local BASE_BRANCH=""

    # First effort in P1W1 starts from main
    if [[ $PHASE -eq 1 && $WAVE -eq 1 && $EFFORT_INDEX -eq 1 ]]; then
        BASE_BRANCH="main"
    # First effort in new wave uses previous wave's last effort
    elif [[ $EFFORT_INDEX -eq 1 ]]; then
        BASE_BRANCH=$(get_previous_wave_last_effort $PHASE $WAVE)
    # Subsequent efforts cascade from previous effort
    else
        BASE_BRANCH=$(get_previous_effort_branch $PHASE $WAVE $((EFFORT_INDEX - 1)))
    fi

    # Get target repository from config
    local TARGET_REPO=$(yq '.repository' target-repo-config.yaml)
    local PROJECT_PREFIX=$(yq '.project_prefix' orchestrator-state-v3.json)

    # Build infrastructure entry
    cat <<EOF
{
  "full_path": "$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/wave${WAVE}/${EFFORT}/",
  "branch_name": "${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT}",
  "remote_branch": "origin/${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT}",
  "base_branch": "${BASE_BRANCH}",
  "target_remote": "target",
  "target_repo_url": "${TARGET_REPO}",
  "planning_remote": "planning",
  "split_pattern": "${EFFORT}--split-NNN",
  "created": false,
  "validated": false,
  "added_at": "$(date -Iseconds)",
  "source": "R505-sync"
}
EOF
}
```

### Update Pre-Planned Infrastructure
```bash
update_pre_planned_infrastructure() {
    local PHASE=$1
    local WAVE=$2
    local EFFORT=$3
    local INFRASTRUCTURE=$4

    local EFFORT_ID="phase${PHASE}_wave${WAVE}_${EFFORT}"

    # Add or update effort entry
    echo "$INFRASTRUCTURE" | yq -P '.' > /tmp/effort_infra.yaml

    yq -i ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\" = load(\"/tmp/effort_infra.yaml\")" orchestrator-state-v3.json

    # Update validation status
    yq -i ".pre_planned_infrastructure.validated = false" orchestrator-state-v3.json
    yq -i ".pre_planned_infrastructure.needs_validation = true" orchestrator-state-v3.json

    echo "✅ Updated infrastructure for $EFFORT_ID"
}
```

## VALIDATION REQUIREMENTS

### Pre-Infrastructure Creation Validation
```bash
validate_infrastructure_sync() {
    local PHASE=$1
    local WAVE=$2

    echo "🔍 Validating pre_planned_infrastructure synchronization..."

    # Check all efforts in current wave have entries
    local WAVE_EFFORTS=$(yq ".efforts_in_wave[]" orchestrator-state-v3.json)
    local MISSING_COUNT=0

    for EFFORT in $WAVE_EFFORTS; do
        local EFFORT_ID="phase${PHASE}_wave${WAVE}_${EFFORT}"
        local EXISTS=$(yq ".pre_planned_infrastructure.efforts.$EFFORT_ID" orchestrator-state-v3.json)

        if [[ "$EXISTS" == "null" ]]; then
            echo "❌ Missing infrastructure for: $EFFORT_ID"
            ((MISSING_COUNT++))
        fi
    done

    if [[ $MISSING_COUNT -gt 0 ]]; then
        echo "❌ FATAL: $MISSING_COUNT efforts missing infrastructure!"
        echo "Run synchronization before CREATE_NEXT_INFRASTRUCTURE"
        exit 505
    fi

    # Verify target repository is correct
    local TARGET_REPO=$(yq '.repository' target-repo-config.yaml)
    local INFRA_COUNT=$(yq '.pre_planned_infrastructure.efforts | length' orchestrator-state-v3.json)
    local CORRECT_REPO_COUNT=$(yq ".pre_planned_infrastructure.efforts.* | select(.target_repo_url == \"$TARGET_REPO\") | length" orchestrator-state-v3.json)

    if [[ "$INFRA_COUNT" != "$CORRECT_REPO_COUNT" ]]; then
        echo "❌ FATAL: Some efforts targeting wrong repository!"
        echo "Expected: $TARGET_REPO"
        exit 505
    fi

    echo "✅ Infrastructure synchronized and valid"
}
```

## CHANGE DETECTION

### Detect Plan Changes
```bash
detect_plan_changes() {
    local PLAN_FILE=$1
    local LAST_SYNC=$2

    # Check if plan modified after last sync
    local PLAN_MODIFIED=$(stat -c %Y "$PLAN_FILE" 2>/dev/null || stat -f %m "$PLAN_FILE" 2>/dev/null)
    local SYNC_TIME=$(date -d "$LAST_SYNC" +%s 2>/dev/null || date -j -f "%Y-%m-%dT%H:%M:%S" "$LAST_SYNC" +%s 2>/dev/null)

    if [[ $PLAN_MODIFIED -gt $SYNC_TIME ]]; then
        echo "⚠️ Plan modified after last sync - resynchronization required"
        return 0
    else
        echo "✅ Plan unchanged since last sync"
        return 1
    fi
}
```

## FAILURE CONDITIONS
- 🚨 Missing pre_planned_infrastructure entries for efforts = FAIL
- 🚨 Wrong target repository in any effort = FAIL
- 🚨 Creating infrastructure without synchronized data = FAIL
- 🚨 Plan changes not synchronized = FAIL
- 🚨 Base branch calculation differs from R308 = FAIL

## PROJECT_DONE CRITERIA
- ✅ ALL efforts have pre_planned_infrastructure entries
- ✅ Synchronization runs after every plan creation/update
- ✅ Target repository matches target-repo-config.yaml
- ✅ Base branches follow R308 cascade algorithm
- ✅ Validation passes before infrastructure creation

## STATE-SPECIFIC ENFORCEMENT

### SPAWN_ARCHITECT_PHASE_PLANNING
```bash
# In orchestrator state rules
echo "📋 Spawning architect for phase planning..."
spawn_architect_phase_planning

# Wait for completion
wait_for_phase_plans

# MANDATORY: Synchronize infrastructure
sync_phase_infrastructure $CURRENT_PHASE "$PHASE_PLAN_FILE"

# Validate synchronization
validate_infrastructure_sync $CURRENT_PHASE 1
```

### SPAWN_ARCHITECT_WAVE_PLANNING
```bash
# In orchestrator state rules
echo "📋 Spawning architect for wave planning..."
spawn_architect_wave_planning

# Wait for completion
wait_for_architecture_plan

# MANDATORY: Synchronize infrastructure
sync_wave_infrastructure $CURRENT_PHASE $CURRENT_WAVE "$WAVE_PLAN_FILE"

# Validate synchronization
validate_infrastructure_sync $CURRENT_PHASE $CURRENT_WAVE
```

### CREATE_NEXT_INFRASTRUCTURE
```bash
# MUST validate synchronization first
validate_infrastructure_sync $CURRENT_PHASE $CURRENT_WAVE

# Only then create infrastructure from pre_planned_infrastructure
create_from_pre_planned_infrastructure
```

## INTEGRATE_WAVE_EFFORTS WITH OTHER RULES

- **R504**: Initial pre-infrastructure planning
- **R308**: CASCADE algorithm for base branch determination
- **R507**: Infrastructure validation after creation
- **R508**: Target repository enforcement
- **R509**: Base branch validation
- **R510**: Infrastructure creation protocol

## MONITORING_SWE_PROGRESS AND REPORTING

```bash
# Generate sync status report
report_infrastructure_sync_status() {
    echo "📊 PRE-PLANNED INFRASTRUCTURE SYNC STATUS"
    echo "========================================="

    local TOTAL_EFFORTS=$(yq '.pre_planned_infrastructure.efforts | length' orchestrator-state-v3.json)
    local VALIDATED=$(yq '.pre_planned_infrastructure.validated' orchestrator-state-v3.json)
    local LAST_PHASE_SYNC=$(yq '.pre_planned_infrastructure.last_phase_sync' orchestrator-state-v3.json)
    local LAST_WAVE_SYNC=$(yq '.pre_planned_infrastructure.last_wave_sync' orchestrator-state-v3.json)

    echo "Total Efforts Tracked: $TOTAL_EFFORTS"
    echo "Validation Status: $VALIDATED"
    echo "Last Phase Sync: $LAST_PHASE_SYNC"
    echo "Last Wave Sync: $LAST_WAVE_SYNC"

    # Check for out-of-sync efforts
    local CURRENT_PHASE=$(yq '.current_phase' orchestrator-state-v3.json)
    local CURRENT_WAVE=$(yq '.current_wave' orchestrator-state-v3.json)

    echo ""
    echo "Current Phase: $CURRENT_PHASE, Wave: $CURRENT_WAVE"
    echo ""

    # List any missing efforts
    echo "Checking for missing efforts..."
    validate_infrastructure_sync $CURRENT_PHASE $CURRENT_WAVE
}
```

## ENFORCEMENT
**Orchestrator MUST:**
1. Run synchronization after EVERY phase/wave plan creation
2. Validate synchronization before ANY infrastructure creation
3. Detect and sync when plans are modified
4. Ensure ALL efforts target the CORRECT repository
5. Calculate base branches per R308 CASCADE algorithm
6. Update pre_planned_infrastructure as the SOLE SOURCE OF TRUTH

**Penalty for violation:** -50% grade for missing synchronization, -100% for wrong repository
## State Manager Coordination (SF 3.0)

State Manager maintains infrastructure sync through state file validation:
- **orchestrator-state-v3.json** tracks current phase/wave in `.state_machine.iteration_container`
- **infrastructure creation** updates both filesystem (directories) and state file (JSON)
- **Atomic commits** ensure filesystem and state file never desync
- **Rollback on failure** prevents partial infrastructure

If directory creation succeeds but state update fails → rollback deletes directories.

See: R281 (initial creation), `tools/atomic-state-update.sh`
