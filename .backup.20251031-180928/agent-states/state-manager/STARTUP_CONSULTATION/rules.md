# State Manager - STARTUP_CONSULTATION State Rules

**State**: STARTUP_CONSULTATION
**Agent**: State Manager
**Purpose**: Validate current state files and provide directive report to Orchestrator at session startup
**Criticality**: SUPREME - Orchestrator depends on accurate state directives

---

## State Overview

This state is entered when the Orchestrator agent begins a work session and requests startup consultation from the State Manager. The State Manager reads all state files, validates their integrity, and generates a comprehensive directive report that guides the Orchestrator's next actions.

**Entry Conditions**:
- Orchestrator spawned State Manager for startup consultation
- State files exist (at least orchestrator-state-v3.json)
- Orchestrator awaiting directive report

**Exit Conditions**:
- Directive report generated and delivered to Orchestrator
- State validation complete (pass or fail documented)
- Orchestrator ready to proceed with work

---

## Core Responsibilities

### 1. Read All State Files

**Required Files** (MUST exist):
- `orchestrator-state-v3.json` - Primary state machine and project progression
- `bug-tracking.json` - Bug tracking state
- `integration-containers.json` - Integration iteration tracking

**Optional Files** (check existence):
- `fix-cascade-state.json` - Only exists during active fix cascades

**Action**:
```bash
# Read all required files
READ: orchestrator-state-v3.json
READ: bug-tracking.json
READ: integration-containers.json

# Check for optional file
if [ -f fix-cascade-state.json ]; then
    READ: fix-cascade-state.json
    CASCADE_ACTIVE=true
else
    CASCADE_ACTIVE=false
fi
```

### 2. Validate Schema Compliance

**Tool**: `tools/validate-state-file.sh`

**Validation Sequence**:
```bash
# Validate each file against its schema
bash tools/validate-state-file.sh orchestrator-state-v3.json
ORCH_VALID=$?

bash tools/validate-state-file.sh bug-tracking.json
BUGS_VALID=$?

bash tools/validate-state-file.sh integration-containers.json
CONTAINERS_VALID=$?

if [ "$CASCADE_ACTIVE" = "true" ]; then
    bash tools/validate-state-file.sh fix-cascade-state.json
    CASCADE_VALID=$?
fi

# Determine overall validation status
if [ $ORCH_VALID -eq 0 ] && [ $BUGS_VALID -eq 0 ] && [ $CONTAINERS_VALID -eq 0 ]; then
    VALIDATION_STATUS="PASS"
else
    VALIDATION_STATUS="FAIL"
fi
```

**On Validation Failure**:
- Document which file(s) failed
- Include validation error messages
- Mark validation_status as "FAIL" in directive report
- Orchestrator MUST fix validation errors before proceeding

### 3. Check State Machine Consistency

**State Machine**: `state-machines/software-factory-3.0-state-machine.json`

**Checks**:
1. **Current State Valid**:
   ```bash
   CURRENT_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)
   jq -e ".states[\"$CURRENT_STATE\"]" state-machines/software-factory-3.0-state-machine.json
   ```

2. **Allowed Transitions**:
   ```bash
   ALLOWED_TRANSITIONS=$(jq -r ".states[\"$CURRENT_STATE\"].allowed_transitions[]" state-machines/software-factory-3.0-state-machine.json)
   ```

3. **Transition Guards**:
   ```bash
   # Check if current state has required conditions
   REQUIRES=$(jq -r ".states[\"$CURRENT_STATE\"].requires.conditions[]?" state-machines/software-factory-3.0-state-machine.json)
   ```

4. **No Orphaned References**:
   - All referenced effort_ids exist in project_progression
   - All referenced bug_ids exist in bug-tracking.json
   - All referenced integration_ids exist in integration-containers.json

### 4. Loop Detection and Prevention (R530)

**Rule**: R530 - State Transition Cycle Detection
**Purpose**: Detect and prevent infinite state transition loops that waste API costs
**Performance**: O(1) time and space complexity

**Required Before Every State Transition**:

#### Step 0: Initialize Loop Detection Counters (If Missing)

**CRITICAL**: Always check if loop_detection exists before loading counters

```bash
# Check if loop_detection exists in orchestrator-state-v3.json
LOOP_DETECTION_EXISTS=$(jq -e '.state_machine.loop_detection' orchestrator-state-v3.json > /dev/null 2>&1 && echo "true" || echo "false")

if [ "$LOOP_DETECTION_EXISTS" == "false" ]; then
    echo "⚠️ Loop detection counters missing - initializing (R530)"

    # Initialize loop_detection with default values
    jq '.state_machine.loop_detection = {
        "last_two_states": [],
        "ping_pong_count": 0,
        "same_state_count": 0,
        "error_recovery_entries": 0,
        "last_progress_timestamp": (now | todate)
    }' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    echo "✅ Loop detection counters initialized"
fi
```

#### Step 1: Load Loop Detection Counters

```bash
# Load loop detection state from orchestrator-state-v3.json
LOOP_STATE=$(jq -r '.state_machine.loop_detection' orchestrator-state-v3.json)

# Extract individual counters
LAST_TWO_STATES=($(echo "$LOOP_STATE" | jq -r '.last_two_states[]' 2>/dev/null))
PING_PONG_COUNT=$(echo "$LOOP_STATE" | jq -r '.ping_pong_count // 0')
SAME_STATE_COUNT=$(echo "$LOOP_STATE" | jq -r '.same_state_count // 0')
ERROR_RECOVERY_ENTRIES=$(echo "$LOOP_STATE" | jq -r '.error_recovery_entries // 0')
LAST_PROGRESS_TS=$(echo "$LOOP_STATE" | jq -r '.last_progress_timestamp // null')

# Get current and target states
CURRENT_STATE=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)
TARGET_STATE="$REQUESTED_STATE"  # From Orchestrator's transition request
```

#### Step 2: Check for Ping-Pong Pattern

**Ping-Pong Loop**: A ↔ B ↔ A ↔ B (alternating between same 2 states)
**Threshold**: 5 cycles

```bash
# Check if target_state matches last_two_states[0] (ping-pong pattern)
if [[ "${#LAST_TWO_STATES[@]}" -eq 2 ]] && [[ "$TARGET_STATE" == "${LAST_TWO_STATES[0]}" ]]; then
    PING_PONG_COUNT=$((PING_PONG_COUNT + 1))
    echo "⚠️ Ping-pong pattern detected: $PING_PONG_COUNT cycles between ${LAST_TWO_STATES[@]}"

    if [[ $PING_PONG_COUNT -ge 5 ]]; then
        echo "🚨🚨🚨 LOOP DETECTED: Ping-pong threshold exceeded ($PING_PONG_COUNT cycles)"
        LOOP_DETECTED=true
        LOOP_REASON="Ping-pong loop: alternating between ${LAST_TWO_STATES[0]} and ${LAST_TWO_STATES[1]}"
        FORCE_STATE="ERROR_RECOVERY"
    fi
fi
```

#### Step 3: Check for Same-State Re-entry

**Same-State Loop**: A → A → A (re-entering same state repeatedly)
**Threshold**: 3 entries

```bash
# Check if target_state matches current_state (same-state re-entry)
if [[ "$TARGET_STATE" == "$CURRENT_STATE" ]]; then
    SAME_STATE_COUNT=$((SAME_STATE_COUNT + 1))
    echo "⚠️ Same-state re-entry detected: $SAME_STATE_COUNT entries to $CURRENT_STATE"

    if [[ $SAME_STATE_COUNT -ge 3 ]]; then
        echo "🚨🚨🚨 LOOP DETECTED: Same-state threshold exceeded ($SAME_STATE_COUNT re-entries)"
        LOOP_DETECTED=true
        LOOP_REASON="Same-state loop: stuck in $CURRENT_STATE"
        FORCE_STATE="ERROR_RECOVERY"
    fi
fi
```

#### Step 4: Check for ERROR_RECOVERY Escalation

**ERROR_RECOVERY Loop**: Multiple failures → ERROR_RECOVERY → failures again
**Threshold**: 3 ERROR_RECOVERY entries

```bash
# Count ERROR_RECOVERY entries
if [[ "$TARGET_STATE" == "ERROR_RECOVERY" ]]; then
    ERROR_RECOVERY_ENTRIES=$((ERROR_RECOVERY_ENTRIES + 1))
    echo "⚠️ ERROR_RECOVERY entry detected: $ERROR_RECOVERY_ENTRIES total entries"

    if [[ $ERROR_RECOVERY_ENTRIES -ge 3 ]]; then
        echo "🚨🚨🚨 LOOP DETECTED: ERROR_RECOVERY threshold exceeded ($ERROR_RECOVERY_ENTRIES entries)"
        LOOP_DETECTED=true
        LOOP_REASON="ERROR_RECOVERY escalation: recovery failed $ERROR_RECOVERY_ENTRIES times"
        FORCE_STATE="ERROR_RECOVERY"
    fi
fi
```

#### Step 5: Check for Genuine Progress (Reset Counters)

**Progress Detection**: Reset counters when genuine progress detected
**Triggers**:
- New effort created (`pre_planned_infrastructure.efforts[].created == true`)
- New effort validated (`pre_planned_infrastructure.efforts[].validated == true`)
- Integration completed (wave/phase/project integration branch created)
- New agent spawned (SWE, Code Reviewer, Architect spawned)
- State NOT in last_two_states (breaking the pattern)

```bash
PROGRESS_DETECTED=false

# Check 1: New effort created
if jq -e '.pre_planned_infrastructure.efforts[] | select(.created == true and .validated == false)' orchestrator-state-v3.json >/dev/null 2>&1; then
    echo "✅ Progress detected: New effort created"
    PROGRESS_DETECTED=true
fi

# Check 2: New effort validated
if jq -e '.pre_planned_infrastructure.efforts[] | select(.validated == true)' orchestrator-state-v3.json >/dev/null 2>&1; then
    echo "✅ Progress detected: Effort validated"
    PROGRESS_DETECTED=true
fi

# Check 3: Integration completed
if [[ "$TARGET_STATE" == "WAVE_COMPLETE" ]] && jq -e '.integration_branches.wave_integrations | length > 0' orchestrator-state-v3.json >/dev/null 2>&1; then
    echo "✅ Progress detected: Wave integration completed"
    PROGRESS_DETECTED=true
fi

# Check 4: New agent spawned (check spawn_history timestamp within last 60 seconds)
CURRENT_TS=$(date +%s)
LAST_SPAWN_TS=$(jq -r '.spawn_history[-1].timestamp // "1970-01-01T00:00:00Z"' orchestrator-state-v3.json | date -d - +%s 2>/dev/null || echo 0)
if [[ $((CURRENT_TS - LAST_SPAWN_TS)) -lt 60 ]]; then
    echo "✅ Progress detected: New agent spawned recently"
    PROGRESS_DETECTED=true
fi

# Check 5: State breaks the pattern (not in last_two_states)
if [[ "${#LAST_TWO_STATES[@]}" -gt 0 ]] && [[ "$TARGET_STATE" != "${LAST_TWO_STATES[0]}" ]] && [[ "$TARGET_STATE" != "${LAST_TWO_STATES[1]}" ]]; then
    echo "✅ Progress detected: New state breaks pattern (entering $TARGET_STATE)"
    PROGRESS_DETECTED=true
fi

# Reset counters on genuine progress
if [[ "$PROGRESS_DETECTED" == "true" ]]; then
    echo "✅ Genuine progress detected - resetting loop counters"
    PING_PONG_COUNT=0
    SAME_STATE_COUNT=0
    ERROR_RECOVERY_ENTRIES=0
    LAST_PROGRESS_TS="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
fi
```

#### Step 6: Update Loop Detection State

**Ring Buffer Update**: Maintain last 2 states for ping-pong detection

```bash
# Update last_two_states ring buffer
if [[ "${#LAST_TWO_STATES[@]}" -eq 0 ]]; then
    LAST_TWO_STATES=("$CURRENT_STATE")
elif [[ "${#LAST_TWO_STATES[@]}" -eq 1 ]]; then
    LAST_TWO_STATES=("${LAST_TWO_STATES[0]}" "$TARGET_STATE")
else
    LAST_TWO_STATES=("${LAST_TWO_STATES[1]}" "$TARGET_STATE")
fi

# Persist loop detection state
jq --argjson last_two "$(printf '%s\n' "${LAST_TWO_STATES[@]}" | jq -R . | jq -s .)" \
   --arg ping_pong "$PING_PONG_COUNT" \
   --arg same_state "$SAME_STATE_COUNT" \
   --arg error_recovery "$ERROR_RECOVERY_ENTRIES" \
   --arg progress_ts "$LAST_PROGRESS_TS" \
   '.state_machine.loop_detection = {
       "last_two_states": $last_two,
       "ping_pong_count": ($ping_pong | tonumber),
       "same_state_count": ($same_state | tonumber),
       "error_recovery_entries": ($error_recovery | tonumber),
       "last_progress_timestamp": $progress_ts
   }' orchestrator-state-v3.json > orchestrator-state-v3.json.tmp && \
   mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json
```

#### Step 7: Enforce ERROR_RECOVERY if Loop Detected

**Hard Stop Enforcement**: Force ERROR_RECOVERY state when loop detected

```bash
if [[ "$LOOP_DETECTED" == "true" ]]; then
    echo "🚨🚨🚨 ERROR_RECOVERY ENFORCED: $LOOP_REASON"
    echo "📊 Loop Detection State:"
    echo "  - Ping-pong count: $PING_PONG_COUNT"
    echo "  - Same-state count: $SAME_STATE_COUNT"
    echo "  - ERROR_RECOVERY entries: $ERROR_RECOVERY_ENTRIES"
    echo "  - Last two states: ${LAST_TWO_STATES[@]}"

    # Force ERROR_RECOVERY state
    TARGET_STATE="ERROR_RECOVERY"

    # Add loop metadata to state file
    jq --arg reason "$LOOP_REASON" \
       --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.state_machine.loop_detection.hard_stop_reason = $reason |
        .state_machine.loop_detection.hard_stop_timestamp = $timestamp' \
       orchestrator-state-v3.json > orchestrator-state-v3.json.tmp && \
       mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json

    # Update directive report to reflect ERROR_RECOVERY
    echo ""
    echo "⚠️⚠️⚠️ DIRECTIVE OVERRIDE: Loop detected - forcing ERROR_RECOVERY"
    echo ""
    echo "Manual intervention required:"
    echo "1. Investigate root cause of loop: $LOOP_REASON"
    echo "2. Fix underlying issue in state files or validation logic"
    echo "3. Reset loop counters manually (see R530 recovery procedure)"
    echo "4. Manually set appropriate state to resume from"
    echo "5. Re-run factory continuation"
fi
```

**Loop Detection Summary** (include in directive report):
```markdown
## Loop Detection Status (R530)

**Loop Detection**: [✅ NO_LOOP / 🚨 LOOP_DETECTED]
**Ping-pong count**: $PING_PONG_COUNT / 5
**Same-state count**: $SAME_STATE_COUNT / 3
**ERROR_RECOVERY entries**: $ERROR_RECOVERY_ENTRIES / 3
**Last progress**: $LAST_PROGRESS_TS
**Last two states**: ${LAST_TWO_STATES[@]}

[If loop detected:]
**🚨 ERROR_RECOVERY ENFORCED**: $LOOP_REASON
**Manual intervention required** - See R530 recovery procedure
```

### 5. Generate Directive Report

**Format**: `directive_report` markdown structure

**Required Sections**:

```markdown
# State Manager Directive Report

**Generated**: [TIMESTAMP_UTC]
**Session**: [SESSION_ID]
**Orchestrator State**: [CURRENT_STATE]

## Validation Results

### Schema Validation
- orchestrator-state-v3.json: [✅ PASS / ❌ FAIL]
- bug-tracking.json: [✅ PASS / ❌ FAIL]
- integration-containers.json: [✅ PASS / ❌ FAIL]
- fix-cascade-state.json: [✅ PASS / ❌ FAIL / ⚪ NOT_PRESENT]

**Overall Status**: [✅ PASS / ❌ FAIL]

### State Machine Validation
- Current state valid: [✅ YES / ❌ NO]
- State exists in state machine: [✅ YES / ❌ NO]
- Allowed transitions: [COUNT] transitions available

### Reference Validation
- Orphaned effort references: [COUNT]
- Orphaned bug references: [COUNT]
- Orphaned integration references: [COUNT]

**Overall Status**: [✅ PASS / ⚠️ WARNINGS / ❌ FAIL]

## Current State Analysis

**Current State**: `[STATE_NAME]`
**State Type**: [ACTION / SPAWN / WAITING_FOR / MONITORING / INTEGRATE_WAVE_EFFORTS / ITERATION_CONTAINER / CONSULTATION / ERROR / TERMINAL]
**State Description**: [from state machine]

**Allowed Transitions** ([COUNT] options):
1. `[STATE_NAME_1]` - [description]
2. `[STATE_NAME_2]` - [description]
...

**Transition Guards** (must satisfy before transition):
- [GUARD_1]: [✅ SATISFIED / ❌ NOT_SATISFIED]
- [GUARD_2]: [✅ SATISFIED / ❌ NOT_SATISFIED]
...

## Project Progression Status

**Current Phase**: [PHASE_NAME]
**Current Wave**: [WAVE_NAME]
**Efforts in Progress**: [COUNT]
**Efforts Completed**: [COUNT]
**Total Efforts**: [COUNT]
**Progress**: [PERCENTAGE]%

**Active Integrations**: [COUNT]
**Active Fix Cascades**: [COUNT]

## Recommended Next Actions

### Priority 1: [ACTION_CATEGORY]
- [SPECIFIC_ACTION_1]
- [SPECIFIC_ACTION_2]

### Priority 2: [ACTION_CATEGORY]
- [SPECIFIC_ACTION_1]

### Blockers (if any)
- ⚠️ [BLOCKER_DESCRIPTION]
- ⚠️ [BLOCKER_DESCRIPTION]

## Detailed Findings

### Validation Errors (if any)
```
[VALIDATION_ERROR_OUTPUT]
```

### Warnings (if any)
- ⚠️ [WARNING_1]
- ⚠️ [WARNING_2]

### State Inconsistencies (if any)
- 🔴 [INCONSISTENCY_1]
- 🔴 [INCONSISTENCY_2]

## Directive Summary

**Can Proceed**: [✅ YES / ❌ NO - FIX ERRORS FIRST]

**Recommended State Transition**: [CURRENT_STATE] → [NEXT_STATE]

**Required Actions Before Transition**:
1. [ACTION_1]
2. [ACTION_2]

**Estimated Work**: [TIME_ESTIMATE]

---

**State Manager Sign-Off**: [TIMESTAMP_UTC]
```

### 5. Deliver Directive Report to Orchestrator

**Delivery Method**: Output the `directive_report` markdown as final message

**Orchestrator Response**:
- If `VALIDATION_STATUS == "PASS"`: Proceed with recommended actions
- If `VALIDATION_STATUS == "FAIL"`: Fix validation errors, re-request consultation
- If blockers present: Resolve blockers before proceeding

---

## Rules Enforcement

### R600: Checklist Execution Protocol
- Directive report must be generated before Orchestrator proceeds
- Validation failures block all work until resolved

### R288: Multi-File Atomic Update Protocol
- State Manager validates ALL 4 files atomically
- No partial validation allowed

### R506: Pre-Commit Hook Enforcement
- Validation uses same schema checks as pre-commit hooks
- Ensures consistency between runtime and commit-time validation

### R516: State Naming Convention Compliance
- Verify current_state matches R516 naming patterns
- Report any non-compliant state names

---

## Exit Criteria

**Must Complete**:
1. ✅ All state files read
2. ✅ Schema validation performed on all files
3. ✅ State machine consistency checked
4. ✅ Directive report generated with all required sections
5. ✅ Report delivered to Orchestrator

**Validation Success**:
- All schema validations pass
- Current state exists in state machine
- No critical inconsistencies found
- Recommended actions provided

**Validation Failure**:
- Document all failures in directive report
- Provide specific fix instructions
- Mark `Can Proceed: ❌ NO` in directive summary

---

## State Transition

**On Successful Consultation**:
- Transition to: [Exit this agent, return control to Orchestrator]
- Orchestrator proceeds with recommended actions

**On Validation Failure**:
- Transition to: [Exit this agent, return control to Orchestrator]
- Orchestrator MUST fix errors before work
- May re-spawn State Manager for re-validation

---

## Examples

### Example 1: Clean Startup (All Valid)

**Scenario**: Orchestrator starts session, all state files valid

**Directive Report Excerpt**:
```markdown
## Validation Results
- orchestrator-state-v3.json: ✅ PASS
- bug-tracking.json: ✅ PASS
- integration-containers.json: ✅ PASS
- fix-cascade-state.json: ⚪ NOT_PRESENT

**Overall Status**: ✅ PASS

## Recommended Next Actions
### Priority 1: Continue Wave Integration
- Proceed to INTEGRATE_WAVE_EFFORTS state
- Merge effort branches: effort-1, effort-2, effort-3
- Spawn Code Reviewer after integration

**Can Proceed**: ✅ YES
```

### Example 2: Schema Validation Failure

**Scenario**: bug-tracking.json has invalid schema

**Directive Report Excerpt**:
```markdown
## Validation Results
- orchestrator-state-v3.json: ✅ PASS
- bug-tracking.json: ❌ FAIL
- integration-containers.json: ✅ PASS

**Overall Status**: ❌ FAIL

### Detailed Findings
```
ERROR: bug-tracking.json validation failed
Line 45: Missing required field "timestamp" in bug entry
Line 67: Invalid severity value "CATASTROPHIC" (must be CRITICAL|HIGH|MEDIUM|LOW)
```

**Can Proceed**: ❌ NO - FIX ERRORS FIRST

**Required Actions Before Transition**:
1. Fix bug-tracking.json schema errors (2 issues)
2. Re-run validation: bash tools/validate-state-file.sh bug-tracking.json
3. Re-request State Manager consultation
```

### Example 3: State Machine Inconsistency

**Scenario**: Current state not in allowed transitions list

**Directive Report Excerpt**:
```markdown
## State Machine Validation
- Current state valid: ❌ NO
- Current state: "INVALID_STATE_NAME"
- State exists in state machine: ❌ NO

### State Inconsistencies
- 🔴 Current state "INVALID_STATE_NAME" does not exist in state machine
- 🔴 Possible corruption or manual state modification detected

**Can Proceed**: ❌ NO - FIX ERRORS FIRST

**Required Actions Before Transition**:
1. Restore orchestrator-state-v3.json from backup
2. OR manually correct current_state to valid state name
3. Re-run State Manager consultation
```

---

## Implementation Notes

### Performance
- Startup consultation should complete in < 5 seconds
- Read all files in parallel if possible
- Cache state machine JSON (single read per session)

### Error Handling
- If ANY file read fails: Report error, mark validation FAIL
- If schema validation crashes: Report crash, mark validation FAIL
- If state machine missing: CRITICAL ERROR, cannot proceed

### Logging
- Log all validation results to session log
- Include timestamps for all checks
- Preserve validation output for debugging

---

**Rule Compliance**: R203 (State-Aware Startup), R287 (TODO Persistence), R288 (Atomic Updates), R506 (Pre-Commit), R516 (State Naming), R530 (Loop Detection), R600 (Checklist Execution)

**Last Updated**: 2025-10-08
**Version**: 1.0
