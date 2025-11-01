# ERROR_RECOVERY State Rules

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
Handle system errors, analyze failures, and determine recovery strategy while preserving existing work.

## Entry Conditions
- ANY state can transition to ERROR_RECOVERY when:
  - Critical validation fails
  - Required resources missing
  - State machine violations detected
  - Unrecoverable errors in agents
  - Manual intervention required

## Mandatory Actions

### 1. Capture Error Context
```bash
# Save error information
ERROR_FROM_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)
ERROR_TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "================================================================"
echo "ERROR RECOVERY INITIATED"
echo "================================================================"
echo "From State: ${ERROR_FROM_STATE}"
echo "Timestamp: ${ERROR_TIMESTAMP}"
echo "Error: ${ERROR_DESCRIPTION:-Unknown error}"
echo "================================================================"

# Create error record
jq --arg from "$ERROR_FROM_STATE" \
   --arg desc "${ERROR_DESCRIPTION:-Unknown error}" \
   --arg ts "$ERROR_TIMESTAMP" \
   '.error_recovery = {
       from_state: $from,
       description: $desc,
       timestamp: $ts,
       recovery_attempts: ((.error_recovery.recovery_attempts // 0) + 1)
   }' orchestrator-state-v3.json > /tmp/state.json

mv /tmp/state.json orchestrator-state-v3.json
```

### 2. Analyze Error Type
```bash
# Determine error category
determine_error_type() {
    local error_desc="$1"

    if [[ "$error_desc" == *"state machine"* ]]; then
        echo "STATE_MACHINE_VIOLATION"
    elif [[ "$error_desc" == *"missing"* ]] || [[ "$error_desc" == *"not found"* ]]; then
        echo "MISSING_RESOURCE"
    elif [[ "$error_desc" == *"validation"* ]] || [[ "$error_desc" == *"invalid"* ]]; then
        echo "VALIDATION_FAILURE"
    elif [[ "$error_desc" == *"size"* ]] || [[ "$error_desc" == *"limit"* ]]; then
        echo "LIMIT_EXCEEDED"
    elif [[ "$error_desc" == *"conflict"* ]] || [[ "$error_desc" == *"merge"* ]]; then
        echo "MERGE_CONFLICT"
    else
        echo "UNKNOWN_ERROR"
    fi
}

ERROR_TYPE=$(determine_error_type "${ERROR_DESCRIPTION}")
echo "Error Type: ${ERROR_TYPE}"
```

### 3. Assess System State
```bash
# Check what work has been completed
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)
EFFORTS_COMPLETE=$(jq -r '.efforts_completed | length' orchestrator-state-v3.json)
EFFORTS_IN_PROGRESS=$(jq -r '.efforts_in_progress | length' orchestrator-state-v3.json)

echo "System Status:"
echo "  Current Phase: ${PHASE}"
echo "  Current Wave: ${WAVE}"
echo "  Efforts Completed: ${EFFORTS_COMPLETE}"
echo "  Efforts In Progress: ${EFFORTS_IN_PROGRESS}"

# Check for work that needs preservation
if [[ "$EFFORTS_IN_PROGRESS" -gt 0 ]]; then
    echo "WARNING: ${EFFORTS_IN_PROGRESS} efforts in progress - need careful recovery"
    NEEDS_PRESERVATION=true
else
    NEEDS_PRESERVATION=false
fi
```

## Recovery Strategies

### 1. Automated Recovery Options
```bash
case "$ERROR_TYPE" in
    "MISSING_RESOURCE")
        echo "Recovery: Create missing resources and retry"
        RECOVERY_STRATEGY="CREATE_AND_RETRY"
        ;;

    "VALIDATION_FAILURE")
        echo "Recovery: Fix validation issues and retry"
        RECOVERY_STRATEGY="FIX_AND_RETRY"
        ;;

    "LIMIT_EXCEEDED")
        echo "Recovery: Split work and continue"
        RECOVERY_STRATEGY="SPLIT_AND_CONTINUE"
        ;;

    "MERGE_CONFLICT")
        echo "Recovery: Resolve conflicts manually"
        RECOVERY_STRATEGY="MANUAL_RESOLUTION"
        ;;

    "STATE_MACHINE_VIOLATION")
        echo "Recovery: Reset to valid state"
        RECOVERY_STRATEGY="RESET_TO_VALID_STATE"
        ;;

    *)
        echo "Recovery: Manual intervention required"
        RECOVERY_STRATEGY="MANUAL_INTERVENTION"
        ;;
esac
```

### 2. Determine Recovery State

**FIRST: Check for R410 Layered Cascade Mode**

```bash
# 🔴🔴🔴 R410: Check if we're in layered cascade mode 🔴🔴🔴
ACTIVE_LAYER=$(jq -r '.cascade_layers | map(select(.status == "fixing_upstream")) | .[0].layer_id // 0' orchestrator-state-v3.json 2>/dev/null)

if [[ $ACTIVE_LAYER -gt 0 ]]; then
    echo "🔧 R410: LAYERED CASCADE MODE DETECTED"
    echo "📋 Currently fixing CASCADE LAYER $ACTIVE_LAYER"

    # Get bugs for this layer
    LAYER_BUGS=$(jq -r ".cascade_layers[] | select(.layer_id == $ACTIVE_LAYER) | .bugs[]" orchestrator-state-v3.json)
    BUG_COUNT=$(echo "$LAYER_BUGS" | wc -w)

    echo "📊 Fixing $BUG_COUNT bugs from layer $ACTIVE_LAYER:"
    echo "$LAYER_BUGS" | while read bug; do
        BUG_DESC=$(jq -r ".bug_registry.bugs[] | select(.id == \"$bug\") | .description" orchestrator-state-v3.json | head -c 100)
        echo "  - $bug: $BUG_DESC..."
    done

    echo ""
    echo "🔄 R410 FIX PROTOCOL:"
    echo "  1. Apply fixes to upstream effort branches"
    echo "  2. Update layer status to 'reintegrating'"
    echo "  3. Return to CASCADE_REINTEGRATION"
    echo "  4. Re-attempt integration with fixes"
    echo ""

    # Apply fixes for this layer
    apply_layer_fixes "$ACTIVE_LAYER" "$LAYER_BUGS"

    # Update layer status to reintegrating
    jq --argjson layer "$ACTIVE_LAYER" \
       '.cascade_layers[] |=
        if .layer_id == $layer then
            .status = "reintegrating" |
            .progress = "fixes applied, returning to cascade"
        else . end' \
       orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "✅ R410: Layer $ACTIVE_LAYER fixes applied"
    echo "🔄 Transitioning back to CASCADE_REINTEGRATION"

    # Return to CASCADE to re-attempt integration
    NEXT_STATE="CASCADE_REINTEGRATION"

    # R410: This is AUTOMATED - use TRUE!
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
    exit 0
fi
```

**IF NOT in cascade mode, use standard recovery:**

```bash
# Based on strategy, determine next state
case "$RECOVERY_STRATEGY" in
    "CREATE_AND_RETRY")
        # Return to state that failed after creating resources
        NEXT_STATE="$ERROR_FROM_STATE"
        ;;

    "FIX_AND_RETRY")
        # Return to state that failed after fixes
        NEXT_STATE="$ERROR_FROM_STATE"
        ;;

    "SPLIT_AND_CONTINUE")
        # Go to split planning
        NEXT_STATE="SPAWN_CODE_REVIEWER_SPLIT_PLAN"
        ;;

    "RESET_TO_VALID_STATE")
        # Find last valid checkpoint
        if [[ "$WAVE" -gt 0 ]]; then
            NEXT_STATE="WAVE_START"
        elif [[ "$PHASE" -gt 0 ]]; then
            NEXT_STATE="START_PHASE_ITERATION"
        else
            NEXT_STATE="INIT"
        fi
        ;;

    "MANUAL_RESOLUTION"|"MANUAL_INTERVENTION")
        # Stay in ERROR_RECOVERY until manual action
        NEXT_STATE="ERROR_RECOVERY"
        echo "⚠️ MANUAL INTERVENTION REQUIRED"
        echo "Please resolve the issue and update state manually"
        ;;

    *)
        NEXT_STATE="ERROR_RECOVERY"
        ;;
esac
```

### R410 Helper: Apply Layer Fixes

```bash
# Apply fixes for a specific cascade layer
apply_layer_fixes() {
    local layer_id="$1"
    local bug_list="$2"

    echo "🔧 Applying fixes for CASCADE LAYER $layer_id"

    # For each bug in this layer
    for bug_id in $bug_list; do
        echo "📋 Processing $bug_id"

        # Get affected efforts from bug registry
        AFFECTED_EFFORTS=$(jq -r ".bug_registry.bugs[] | select(.id == \"$bug_id\") | .affected_efforts[]" orchestrator-state-v3.json)

        if [[ -z "$AFFECTED_EFFORTS" ]]; then
            echo "⚠️ No affected efforts listed for $bug_id - analyzing bug..."
            # Determine affected efforts from bug description
            # (implementation-specific logic here)
        fi

        # Apply fix to each affected effort
        for effort in $AFFECTED_EFFORTS; do
            echo "🔧 Applying fix to $effort for $bug_id"

            # Get effort directory
            EFFORT_DIR=$(jq -r ".pre_planned_infrastructure.efforts[] | select(.effort_key == \"$effort\") | .directory" orchestrator-state-v3.json)

            if [[ -z "$EFFORT_DIR" ]] || [[ "$EFFORT_DIR" == "null" ]]; then
                echo "❌ Cannot find directory for $effort"
                continue
            fi

            # Apply fix (spawn SW Engineer with fix instructions)
            # (detailed fix application logic per R300)

            echo "✅ Fix applied to $effort"
        done

        # Update bug status
        jq --arg bug "$bug_id" \
           '.bug_registry.bugs[] |=
            if .id == $bug then
                .status = "fixed_in_effort" |
                .fixed_at = (now | todate)
            else . end' \
           orchestrator-state-v3.json > /tmp/state.json

        mv /tmp/state.json orchestrator-state-v3.json

        echo "✅ $bug_id marked as fixed"
    done

    echo "✅ All layer $layer_id fixes applied"
}
```

### 3. Create Recovery Report
```bash
cat <<EOF > "error-recovery-$(date +%Y%m%d-%H%M%S).md"
# Error Recovery Report

## Error Information
- **From State**: ${ERROR_FROM_STATE}
- **Error Type**: ${ERROR_TYPE}
- **Description**: ${ERROR_DESCRIPTION}
- **Timestamp**: ${ERROR_TIMESTAMP}

## System Status
- **Phase**: ${PHASE}
- **Wave**: ${WAVE}
- **Efforts Completed**: ${EFFORTS_COMPLETE}
- **Efforts In Progress**: ${EFFORTS_IN_PROGRESS}

## Recovery Strategy
- **Strategy**: ${RECOVERY_STRATEGY}
- **Next State**: ${NEXT_STATE}
- **Needs Preservation**: ${NEEDS_PRESERVATION}

## Recovery Actions
$(case "$RECOVERY_STRATEGY" in
    "CREATE_AND_RETRY")
        echo "1. Create missing resources"
        echo "2. Validate environment"
        echo "3. Retry from ${ERROR_FROM_STATE}"
        ;;
    "FIX_AND_RETRY")
        echo "1. Fix validation issues"
        echo "2. Re-run validations"
        echo "3. Retry from ${ERROR_FROM_STATE}"
        ;;
    "SPLIT_AND_CONTINUE")
        echo "1. Create split plan"
        echo "2. Execute splits sequentially"
        echo "3. Continue with reduced scope"
        ;;
    "MANUAL_RESOLUTION")
        echo "1. Review error details"
        echo "2. Manually resolve conflicts"
        echo "3. Update state file"
        echo "4. Continue with /continue-orchestrating"
        ;;
    *)
        echo "1. Review error details"
        echo "2. Determine manual fix"
        echo "3. Apply fix"
        echo "4. Update state and continue"
        ;;
esac)

## Notes
${RECOVERY_NOTES:-None}
EOF

echo "Recovery report created: error-recovery-$(date +%Y%m%d-%H%M%S).md"

## Exit Conditions - CRITICAL AUTOMATION GUIDANCE

### 🔴🔴🔴 DURING SEQUENTIAL FIX WORK (WITHIN ERROR_RECOVERY) 🔴🔴🔴

When working through multiple bugs sequentially:

```bash
# Fix Bug 1
spawn_engineer_for_bug_1()
wait_for_fix()
verify_fix()

# R322 checkpoint
echo "🛑 R322: Checkpoint after Bug 1 complete"
commit_state_per_r288()
exit 0

# Can automation continue to Bug 2? YES!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**Rationale:** Fixing bugs sequentially is NORMAL ERROR_RECOVERY workflow. The system knows:
- Current state: ERROR_RECOVERY
- Remaining work: Bugs 2, 3, 4, 5
- Next action: Spawn engineer for Bug 2

**No human intervention needed!**

### 🔴🔴🔴 BEFORE STATE TRANSITION (EXITING ERROR_RECOVERY) 🔴🔴🔴

When all fixes complete and ready to CASCADE:

```bash
# All bugs fixed, ready to cascade
mark_all_fixes_complete()
prepare_cascade()

# R322 stop before state transition
echo "🛑 R322: Checkpoint before CASCADE_REINTEGRATION transition"
update_state("CASCADE_REINTEGRATION")
commit_per_r288()
exit 0

# Can automation continue with cascade? YES!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**Rationale:** State transitions are state machine operations. The system knows:
- Current: ERROR_RECOVERY
- Next: CASCADE_REINTEGRATION
- What to do: Execute cascade per R327

**No human intervention needed!**

### ❌ ONLY Use FALSE When

```bash
# Truly unrecoverable situation
echo "❌ CRITICAL: Cannot determine which bug to fix"
echo "❌ Bug registry corrupted - human debugging required"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
exit 1
```

### Summary

- ✅ **TRUE**: 99.9% of ERROR_RECOVERY operations
- ❌ **FALSE**: 0.1% catastrophic failures only
- 🛑 **R322 stops**: Totally independent from flag value!

## Automation Flag - COMPREHENSIVE GUIDANCE

```bash
# Determine automation flag based on recovery status

if [[ "$RECOVERY_STRATEGY" == "MANUAL_INTERVENTION" ]] || \
   [[ "$RECOVERY_STRATEGY" == "MANUAL_RESOLUTION" ]]; then
    # Only these strategies require FALSE
    echo "⚠️ ERROR RECOVERY MODE"
    echo "Manual intervention required to resolve issues."
    echo "Review error details and determine recovery strategy."
    echo ""
    echo "Options:"
    echo "1. Fix issues and resume with: /continue-orchestrating"
    echo "2. Reset to specific state with manual state update"
    echo "3. Abort orchestration if unrecoverable"
    echo ""

    # MANDATORY: Print cascade status if active (R406 auto-reporting)
    if [ -f "orchestrator-state-v3.json" ] && \
       jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
        echo ""
        echo "═══════════════════════════════════════════════════════════"
        echo "📊 R406 FIX CASCADE STATUS (automatic report)"
        echo "═══════════════════════════════════════════════════════════"
        source utilities/cascade-status-report.sh
        cascade_status_report
        echo "═══════════════════════════════════════════════════════════"
    fi

    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # True system failure
else
    # All other strategies are automated
    echo "✅ ERROR RECOVERY COMPLETE"
    echo "Recovery strategy: ${RECOVERY_STRATEGY}"
    echo "Next state: ${NEXT_STATE}"
    echo ""
    echo "Automated recovery will proceed..."
    echo ""

    # MANDATORY: Print cascade status if active (R406 auto-reporting)
    if [ -f "orchestrator-state-v3.json" ] && \
       jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
        echo ""
        echo "═══════════════════════════════════════════════════════════"
        echo "📊 R406 FIX CASCADE STATUS (automatic report)"
        echo "═══════════════════════════════════════════════════════════"
        source utilities/cascade-status-report.sh
        cascade_status_report
        echo "═══════════════════════════════════════════════════════════"
        echo ""
    fi

    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Automation handles it
fi
```
```

## 🔴🔴🔴 CRITICAL: R327 CASCADE ENFORCEMENT AFTER FIXES 🔴🔴🔴

### SUPREME LAW: ALL FIXES MUST BE CASCADED

**When fixes are applied to effort branches in ERROR_RECOVERY, ALL integration branches containing those efforts become STALE and MUST be recreated!**

This is **R327 - Mandatory Re-Integration After Fixes** - a SUPREME LAW with -100% penalty for violation.

### Detection: Are Integrations Stale?
```bash
# MANDATORY CHECK after applying fixes:
detect_stale_integrations() {
    echo "🔍 R327 ENFORCEMENT: Checking for stale integrations"

    local STALE_DETECTED=false

    # Check each integration level
    for integration_type in wave phase project; do
        local INTEGRATE_WAVE_EFFORTS_BRANCH="${integration_type}-integration"

        # Does this integration exist?
        if git show-ref --verify --quiet "refs/remotes/origin/${INTEGRATE_WAVE_EFFORTS_BRANCH}"; then
            # Get integration creation timestamp
            local INTEGRATE_WAVE_EFFORTS_TIME=$(git log -1 --format=%ct "origin/${INTEGRATE_WAVE_EFFORTS_BRANCH}")

            # Find newest fix commit in ANY effort branch
            local NEWEST_FIX_TIME=0
            for effort_branch in $(git branch -r | grep "effort-" | grep -v "integration"); do
                local FIX_TIME=$(git log -1 --grep="fix:" --format=%ct "$effort_branch" 2>/dev/null || echo 0)
                [[ $FIX_TIME -gt $NEWEST_FIX_TIME ]] && NEWEST_FIX_TIME=$FIX_TIME
            done

            # Compare timestamps
            if [[ $NEWEST_FIX_TIME -gt $INTEGRATE_WAVE_EFFORTS_TIME ]]; then
                echo "❌ R327 VIOLATION DETECTED!"
                echo "   ${INTEGRATE_WAVE_EFFORTS_BRANCH} is STALE"
                echo "   Integration created: $(date -d "@$INTEGRATE_WAVE_EFFORTS_TIME" 2>/dev/null || echo "timestamp: $INTEGRATE_WAVE_EFFORTS_TIME")"
                echo "   Newest fix applied: $(date -d "@$NEWEST_FIX_TIME" 2>/dev/null || echo "timestamp: $NEWEST_FIX_TIME")"
                echo "   🔴 CASCADE RE-INTEGRATE_WAVE_EFFORTS MANDATORY!"
                STALE_DETECTED=true
            fi
        fi
    done

    if [[ "$STALE_DETECTED" == "true" ]]; then
        echo "🔴🔴🔴 STALE INTEGRATE_WAVE_EFFORTSS DETECTED 🔴🔴🔴"
        echo "MUST transition to CASCADE_REINTEGRATION"
        return 1  # Stale detected
    else
        echo "✅ All integrations are current (no fixes after creation)"
        return 0  # All current
    fi
}

# Run detection after fixes applied
if ! detect_stale_integrations; then
    echo "🔴 R327 CASCADE ENFORCEMENT ACTIVATED"
    NEXT_STATE="CASCADE_REINTEGRATION"
fi
```

### Enforcement: Mandatory Cascade Path

**AFTER EFFORT FIXES, YOU MUST:**
1. ✅ Detect stale integrations (above)
2. ✅ Transition to CASCADE_REINTEGRATION state
3. ✅ Execute complete cascade (wave → phase → project)
4. ✅ ONLY THEN proceed to validation/success

**FORBIDDEN TRANSITIONS AFTER EFFORT FIXES:**
- ❌ ERROR_RECOVERY → PROJECT_INTEGRATE_WAVE_EFFORTS (skips cascade!)
- ❌ ERROR_RECOVERY → MONITORING_PROJECT_INTEGRATE_WAVE_EFFORTS (stale!)
- ❌ ERROR_RECOVERY → PROJECT_DONE (broken integrations!)

**REQUIRED TRANSITION AFTER EFFORT FIXES:**
- ✅ ERROR_RECOVERY → CASCADE_REINTEGRATION → [cascade execution] → REVIEW_WAVE_INTEGRATION → PROJECT_DONE

### Why This Matters

**Without R327 Cascade (WHAT WENT WRONG):**
```
1. Project integration created with broken code
2. ERROR_RECOVERY applies fixes to effort branches
3. Orchestrator goes: ERROR_RECOVERY → PROJECT_INTEGRATE_WAVE_EFFORTS → PROJECT_DONE
4. Result: project-integration has OLD/BROKEN code
5. Result: effort branches have NEW/FIXED code
6. Result: Integration is UNUSABLE
```

**With R327 Cascade (CORRECT):**
```
1. Project integration created with broken code
2. ERROR_RECOVERY applies fixes to effort branches
3. Detect: effort branches have commits newer than integrations
4. Transition: ERROR_RECOVERY → CASCADE_REINTEGRATION
5. DELETE wave integration (stale per R327)
6. RECREATE wave integration from fixed efforts
7. DELETE phase integration (cascade up)
8. RECREATE phase integration from new wave
9. DELETE project integration (cascade up)
10. RECREATE project integration from new phase
11. Result: ALL integrations have ALL fixes
12. THEN proceed to validation/success
```

### Exit State Selection Logic
```bash
# Determine correct exit state from ERROR_RECOVERY
determine_exit_state() {
    local ERROR_TYPE="$1"
    local RECOVERY_STRATEGY="$2"

    # FIRST: Check for stale integrations (R327 enforcement)
    if ! detect_stale_integrations; then
        echo "🔴 R327 ENFORCEMENT: Stale integrations detected"
        echo "MANDATORY transition to CASCADE_REINTEGRATION"
        echo "CASCADE_REINTEGRATION"
        return
    fi

    # THEN: Check recovery strategy
    case "$RECOVERY_STRATEGY" in
        "CREATE_AND_RETRY")
            echo "$ERROR_FROM_STATE"
            ;;
        "FIX_AND_RETRY")
            # Fixes were applied - double-check cascade
            if ! detect_stale_integrations; then
                echo "CASCADE_REINTEGRATION"
            else
                echo "$ERROR_FROM_STATE"
            fi
            ;;
        "SPLIT_AND_CONTINUE")
            echo "SPAWN_CODE_REVIEWER_SPLIT_PLAN"
            ;;
        "RESET_TO_VALID_STATE")
            if [[ "$WAVE" -gt 0 ]]; then
                echo "WAVE_START"
            elif [[ "$PHASE" -gt 0 ]]; then
                echo "START_PHASE_ITERATION"
            else
                echo "INIT"
            fi
            ;;
        *)
            echo "ERROR_RECOVERY"
            ;;
    esac
}

NEXT_STATE=$(determine_exit_state "$ERROR_TYPE" "$RECOVERY_STRATEGY")
```

## Exit Conditions

### Success Criteria
- Error captured and documented
- Recovery strategy determined
- System state preserved
- Next state identified
- **R327 cascade requirements verified (CRITICAL)**

### State Transitions
- **CASCADE_REINTEGRATION**: If effort fixes applied and integrations are stale (R327 MANDATORY)
- **[PREVIOUS_STATE]**: If error resolved and can retry (only if NO stale integrations)
- **INIT**: If full reset required
- **START_PHASE_ITERATION**: If resetting to phase boundary
- **WAVE_START**: If resetting to wave boundary
- **ERROR_RECOVERY**: If manual intervention required

### State Update Requirements
```bash
# Update state based on recovery
update_state() {
    local next_state="$1"
    local notes="${2:-Error recovery complete}"

    if [[ "$next_state" == "ERROR_RECOVERY" ]]; then
        # Staying in error recovery - increment counter
        jq --arg notes "$notes" \
           --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
           '.error_recovery.manual_intervention_required = true |
            .error_recovery.last_checked = $timestamp |
            .error_recovery.notes = $notes' \
           orchestrator-state-v3.json > /tmp/state.json
    else
        # Moving to recovery state
        jq --arg state "$next_state" \
           --arg notes "$notes" \
           --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
           '.state_machine.current_state = $state |
            .error_recovery.resolved = true |
            .error_recovery.resolution = $notes |
            .last_transition = {
                from: "ERROR_RECOVERY",
                to: $state,
                timestamp: $timestamp,
                notes: $notes
            }' orchestrator-state-v3.json > /tmp/state.json
    fi

    mv /tmp/state.json orchestrator-state-v3.json
    echo "State updated: ${next_state}"
}

# Apply recovery decision
if [[ "$NEXT_STATE" != "ERROR_RECOVERY" ]]; then
    echo "Attempting automatic recovery to: ${NEXT_STATE}"
    update_state "${NEXT_STATE}" "Automated recovery: ${RECOVERY_STRATEGY}"
else
    echo "Manual intervention required - staying in ERROR_RECOVERY"
    update_state "ERROR_RECOVERY" "Manual intervention required: ${ERROR_DESCRIPTION}"
fi
```

## Associated Rules
- **R290**: State rule reading verification (SUPREME LAW)
- **R233**: Single operation per state (SUPREME LAW)
- **R285**: Error handling protocol
- **R288**: State file update requirements
- **R283**: Recovery checkpoint creation

## Prohibitions
- ❌ Ignore errors and continue
- ❌ Lose completed work during recovery
- ❌ Retry indefinitely without resolution
- ❌ Corrupt state file during recovery
- ❌ Skip error documentation

## Notes
- ERROR_RECOVERY is the universal error handler
- Preserves work when possible
- Documents all errors for analysis
- Supports both automated and manual recovery
- Creates checkpoints for rollback if needed
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
