# 🔴🔴🔴 RULE R288 - MANDATORY STATE FILE UPDATE AND COMMIT PROTOCOL 🔴🔴🔴

## 🚨🚨🚨 BLOCKING: SUPREME LAW OF STATE PERSISTENCE 🚨🚨🚨

## THE ABSOLUTE DUAL REQUIREMENT:

**EVERY state transition MUST:**
1. **UPDATE** all 4 state files ATOMICALLY using atomic-state-update.sh (within 30 seconds)
2. **COMMIT AND PUSH** all 4 files in a SINGLE ATOMIC COMMIT (within 60 seconds)

**NO EXCEPTIONS. NO DEFERRALS. NO BATCHING. NO "LATER".**

## SF 3.0 ARCHITECTURE - 4-FILE ATOMIC UPDATES:

Software Factory 3.0 uses **4 separate state files** that MUST be updated atomically:
1. **orchestrator-state-v3.json** - State machine, project progression, references
2. **bug-tracking.json** - All bugs discovered during reviews
3. **integration-containers.json** - Active iteration containers (wave/phase/project)
4. **fix-cascade-state.json** - Cross-container bug cascade tracking (created on demand)

**CRITICAL**: All 4 files MUST be updated and committed together using `tools/atomic-state-update.sh`

## CRITICAL COMPANION RULE:
**R281 (SUPREME LAW #7)**: When creating the INITIAL state file, it MUST contain ALL phases, waves, and efforts from the implementation plan. See R281 for complete requirements.

## MANDATORY EXECUTION SEQUENCE:

```bash
# THE ONLY ACCEPTABLE PATTERN (SF 3.0)
perform_state_transition() {
    local OLD_STATE="$1"
    local NEW_STATE="$2"
    local REASON="$3"

    # Step 1: Validate transition is allowed
    validate_state_transition "$OLD_STATE" "$NEW_STATE"

    # Step 2: 🔴 UPDATE ALL 4 STATE FILES ATOMICALLY 🔴
    # Use text_editor tool with str_replace to update orchestrator-state-v3.json:
    # - Replace current_state value with NEW_STATE
    # - Replace previous_state value with OLD_STATE
    # - Replace last_transition_timestamp with current UTC timestamp
    # - APPEND to state_history array (CRITICAL - see below)
    # Also update bug-tracking.json, integration-containers.json if needed

    # CRITICAL: Append to state_history array
    # Create transition entry for state_history
    local TRANSITION_ENTRY=$(cat <<EOF
{
  "from_state": "$OLD_STATE",
  "to_state": "$NEW_STATE",
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "validated_by": "state-manager",
  "reason": "$REASON"
}
EOF
)
    # Use jq to append to state_history array in orchestrator-state-v3.json
    jq ".state_machine.state_history += [$TRANSITION_ENTRY]" orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    # Step 3: 🔴🔴🔴 USE ATOMIC-STATE-UPDATE.SH 🔴🔴🔴
    # BLOCKING: All 4 files MUST be updated atomically
    bash "$CLAUDE_PROJECT_DIR/tools/atomic-state-update.sh" \
        --reason "${OLD_STATE} → ${NEW_STATE} - ${REASON}" || {
        echo "❌❌❌ R288 VIOLATION: Atomic state update failed!"
        echo "Rollback executed - files restored from backup"
        exit 288
    }

    # atomic-state-update.sh handles:
    # - Backup of all 4 files
    # - JSON validation of all 4 files
    # - Single atomic commit with all 4 files
    # - Push to remote
    # - Automatic rollback on any failure

    # Step 4: Reload rules for new state (R217)
    reload_rules_for_state "$NEW_STATE"

    echo "✅ State transition complete and persisted: $OLD_STATE → $NEW_STATE"
    echo "✅ State history updated with transition record"
    echo "✅ All 4 state files committed atomically"
}
```

## REQUIRED STATE FILE FIELDS:

### Core Fields (EVERY transition):
```yaml
state_machine:
  current_state: "NEW_STATE"
  previous_state: "OLD_STATE"
  last_transition_timestamp: "2025-08-25T12:00:00Z"
  state_history: [
    # CRITICAL: Append new transition entry to this array on EVERY state transition
    {
      "from_state": "OLD_STATE",
      "to_state": "NEW_STATE",
      "timestamp": "2025-08-25T12:00:00Z",
      "validated_by": "state-manager",  # MUST be "state-manager" (SF 3.0)
      "reason": "Clear explanation of why transition occurred"
    }
  ]
```

### SF 3.0 Proposal/Decision Metadata Fields (added by State Manager):

When orchestrator proposes a next state during SHUTDOWN_CONSULTATION, State Manager adds these fields:

```yaml
state_history:
  - from_state: "WAITING_FOR_PROJECT_TEST_PLAN"
    to_state: "CREATE_PROJECT_INTEGRATION_BRANCH_EARLY"
    timestamp: "2025-10-12T19:25:42Z"
    validated_by: "state-manager"

    # Orchestrator's proposal and State Manager's decision
    orchestrator_proposal: "INIT"  # What orchestrator proposed
    proposal_accepted: false       # State Manager rejected proposal
    proposal_rejected_reason: "Cannot return to INIT - must continue mandatory sequence project_initialization"

    # Mandatory sequence tracking (if applicable)
    mandatory_sequence: "project_initialization"
    sequence_position: "6/9"
```

**Field Definitions**:
- `orchestrator_proposal` (string, optional): State that orchestrator proposed during SHUTDOWN_CONSULTATION
- `proposal_accepted` (boolean, required if orchestrator_proposal present):
  - `true`: State Manager accepted orchestrator's proposal
  - `false`: State Manager rejected/overrode orchestrator's proposal
- `proposal_rejected_reason` (string, required if proposal_accepted=false): Explanation why State Manager overrode the proposal
- `mandatory_sequence` (string, optional): Name of mandatory sequence (e.g., "project_initialization")
- `sequence_position` (string, optional): Position in sequence (e.g., "6/9")

**IMPORTANT**: These fields are ONLY added by State Manager during SHUTDOWN_CONSULTATION. Orchestrator NEVER writes these fields.

**CRITICAL - STATE_HISTORY REQUIREMENT (SF 3.0)**:
- The `state_machine.state_history` array MUST be appended to on EVERY state transition
- This array provides live tracking of state progression during execution
- Runtime tests read this array to validate state machine compliance
- Without state_history updates, tests will FAIL even if orchestrator works correctly
- State Manager automatically appends to state_history during SHUTDOWN_CONSULTATION

### State-Specific Updates:
- **WAVE_COMPLETE**: Add waves_completed entry with metrics
- **INTEGRATE_WAVE_EFFORTS**: Add current_integration details
- **ERROR_RECOVERY**: Add error_context information
- **PROJECT_DONE**: Add phase_completion summary
- **SPAWN_SW_ENGINEERS**: Add agents_spawned records
- **MONITOR**: Update monitoring_status

## ENFORCEMENT PROTOCOL:

### The Golden Pattern:
**EDIT → COMMIT → PUSH** (ALWAYS in this order, ALWAYS immediate)

### Timing Requirements:
- **Update**: Within 30 seconds of state transition decision
- **Commit**: Within 60 seconds of update
- **Push**: Immediately after commit (no delay)

### Commit Message Format:
```
state: <what changed> - <why> [R288]
```

Examples:
- `state: PLANNING → CREATE_NEXT_INFRASTRUCTURE - planning complete [R288]`
- `state: wave1 marked complete - all efforts reviewed [R288]`
- `state: effort1 status=BLOCKED - size limit exceeded [R288]`

## MANDATORY WRAPPER FUNCTION:

```bash
# ALL state updates MUST use this wrapper (SF 3.0)
update_and_commit_state() {
    local KEY="$1"
    local VALUE="$2"
    local REASON="${3:-update}"

    # 1. Make the edit to appropriate state file(s)
    # Use text_editor tool with str_replace to update the relevant state file:
    # - orchestrator-state-v3.json: state machine, project progression
    # - bug-tracking.json: bugs array
    # - integration-containers.json: active integrations
    # - fix-cascade-state.json: cascade tracking
    # Replace the KEY's current value with VALUE

    # 2. 🔴🔴🔴 USE ATOMIC-STATE-UPDATE.SH 🔴🔴🔴
    # Validates all 4 files, creates single atomic commit, handles rollback
    bash "$CLAUDE_PROJECT_DIR/tools/atomic-state-update.sh" \
        --reason "${KEY}=${VALUE} - ${REASON}" || {
        echo "❌❌❌ R288 VIOLATION: Atomic state update failed!"
        echo "All changes rolled back automatically"
        exit 288
    }

    echo "✅ State updated atomically: ${KEY}=${VALUE}"
}
```

## COMMON VIOLATIONS:

### ❌ FORBIDDEN PATTERNS:
```bash
# ❌ NO: Deferred commit
# Use text_editor tool with str_replace to update orchestrator-state-v3.json:
# Replace current_state: "old_value" with current_state: "INTEGRATE_WAVE_EFFORTS"
do_other_work()  # VIOLATION! Must commit first!

# ❌ NO: Batch updates
# Use text_editor tool with str_replace to update orchestrator-state-v3.json:
# - Replace current_state value with "WAVE_COMPLETE"
# - Replace wave1.status value with "COMPLETE"
git add orchestrator-state-v3.json  # VIOLATION! Each edit needs commit!

# ❌ NO: Missing push
git add orchestrator-state-v3.json
git commit -m "state: update"
# No push = VIOLATION!
```

### ✅ REQUIRED PATTERN:
```bash
# ✅ YES: Immediate atomic update using atomic-state-update.sh
# Use text_editor tool with str_replace to update orchestrator-state-v3.json:
# Replace current_state: "old_value" with current_state: "INTEGRATE_WAVE_EFFORTS"
bash tools/atomic-state-update.sh --reason "transition to INTEGRATE_WAVE_EFFORTS"
# This commits ALL 4 files atomically and pushes

# ✅ YES: Each logical change gets its own atomic commit
update_and_commit_state "current_state" "WAVE_COMPLETE" "all efforts done"
update_and_commit_state "wave1.status" "COMPLETE" "wave finished"
```

## COMPLIANCE MONITORING_SWE_PROGRESS:

```bash
check_r288_compliance() {
    # Check for uncommitted state changes (ANY of the 4 files)
    local STATE_FILES="orchestrator-state-v3.json bug-tracking.json integration-containers.json fix-cascade-state.json"
    for file in $STATE_FILES; do
        if git status --porcelain | grep -q "$file"; then
            echo "❌❌❌ R288 VIOLATION: Uncommitted state changes in $file!"
            return 288
        fi
    done

    # Check timestamp freshness
    # Use text_editor tool with view command to read orchestrator-state-v3.json:
    # Find the transition_time field at root level
    local TIMESTAMP="<value from .state_machine.transition_time>"
    local NOW=$(date +%s)
    local TRANS_TIME=$(date -d "$TIMESTAMP" +%s 2>/dev/null || echo 0)
    local AGE=$((NOW - TRANS_TIME))

    if [ $AGE -gt 60 ]; then
        echo "⚠️ R288 WARNING: State timestamp stale (${AGE}s old)"
    fi

    echo "✅ R288 Compliance: OK - All 4 state files committed"
}
```

## GRADING PENALTIES:

### AUTOMATIC FAILURE CONDITIONS:
- State transition without immediate update: **FAIL**
- State update without immediate commit/push: **FAIL**
- Batch commits of multiple changes: **FAIL**
- Uncommitted state changes >30 seconds: **FAIL**
- Missing [R288] tag in commits: **-10%**
- Stale timestamp (>60s): **-20%**

### VIOLATION PENALTIES:
- First violation: **-20%** on state management
- Second violation: **-50%** on state management
- Third violation: **AUTOMATIC FAIL**
- Lost state due to non-persistence: **-100% IMMEDIATE FAIL**

## RECOVERY PROTOCOL:

If you detect a violation:
1. **STOP** all work immediately
2. **COMMIT** current state: `git commit -m "state: RECOVERY - R288 violation [R288-VIOLATION]"`
3. **PUSH** immediately
4. **LOG** violation in state file:
   ```bash
   # Use text_editor tool to increment r288_violations in orchestrator-state-v3.json:
   # First view to get current value, then str_replace with incremented value
   bash tools/atomic-state-update.sh --reason "R288 violation logged"
   # This atomically commits all 4 state files with violation log
   ```

## WHY THIS IS CRITICAL:

1. **State Loss Prevention**: Crashes don't lose progress
2. **Multi-Instance Safety**: All instances see current state
3. **Recovery Capability**: Can resume from exact point
4. **Complete Audit Trail**: Every transition tracked
5. **Debugging Support**: Full state history available

## THE SUPREME LAWS:

1. **No state transition is complete until all 4 files are updated atomically**
2. **No state update is safe until it's pushed to remote**
3. **Every logical change deserves its own atomic commit**
4. **The 4 state files together are the single source of truth**
5. **If it's not pushed, it didn't happen**
6. **All 4 files must be committed together or none at all (atomicity)**

## State Manager Coordination (SF 3.0)

**THIS IS THE STATE MANAGER'S PRIMARY RESPONSIBILITY**

State Manager implements R288 through its bookend pattern:

1. **Startup Consultation** (STARTUP_CONSULTATION state):
   - Read current state from all 4 files
   - Validate file consistency
   - Provide `directive_report` to Orchestrator with current state context

2. **Shutdown Consultation** (SHUTDOWN_CONSULTATION state):
   - Validate proposed state changes from Orchestrator
   - Execute atomic update via `tools/atomic-state-update.sh`
   - Validate all 4 files against schemas
   - Rollback on ANY validation failure
   - Commit all 4 files together with [R288] tag
   - Return `validation_result` to Orchestrator

The bookend pattern ensures NO Orchestrator work occurs without proper state management.

See: `.claude/agents/state-manager.md`, `agent-states/state-manager/STARTUP_CONSULTATION/rules.md`, `agent-states/state-manager/SHUTDOWN_CONSULTATION/rules.md`, `tools/atomic-state-update.sh`

## FINAL WARNING:

**This rule consolidates and supersedes R252 and R253.**

**VIOLATION = FAILURE. NO EXCEPTIONS.**

**UPDATE → COMMIT → PUSH or FAIL!**