# MASTER_HANDLE_SUB_COMPLETION State Rules

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
Process completion of sub-orchestrators, integrate their results into master state, and determine next actions based on outcomes.

## Entry Criteria
- One or more sub-orchestrators have completed
- Output files available from completed subs
- Master state ready for updates

## Actions Required

### 1. Read Sub-Orchestrator Outputs
```bash
process_sub_output() {
  local SUB_ID="$1"
  local OUTPUT_FILE="/tmp/sub-orch-${SUB_ID}/output.json"

  if [[ ! -f "$OUTPUT_FILE" ]]; then
    handle_missing_output "$SUB_ID"
    return 1
  fi

  local OUTPUT=$(cat "$OUTPUT_FILE")
  local STATUS=$(echo "$OUTPUT" | jq -r '.status')
  local RESULTS=$(echo "$OUTPUT" | jq '.results')

  case "$STATUS" in
    "PROJECT_DONE")
      integrate_successful_results "$SUB_ID" "$RESULTS"
      ;;
    "FAILED")
      handle_failed_sub "$SUB_ID" "$OUTPUT"
      ;;
    "PARTIAL")
      handle_partial_results "$SUB_ID" "$RESULTS"
      ;;
  esac
}
```

### 2. Integrate Results by Type

#### FIX_CASCADE Results
```json
{
  "fix_cascade_complete": {
    "fix_id": "bug-123",
    "branches_fixed": ["branch1", "branch2"],
    "test_results": {
      "branch1": "PASSED",
      "branch2": "PASSED"
    },
    "all_successful": true
  }
}
```

#### INTEGRATE_WAVE_EFFORTS Results
```json
{
  "integration_complete": {
    "type": "WAVE",
    "integration_branch": "wave-1-integration",
    "merged_branches": ["effort-1", "effort-2"],
    "conflicts_resolved": [],
    "build_status": "PROJECT_DONE",
    "test_status": "PASSED"
  }
}
```

#### SPLIT_COORDINATION Results
```json
{
  "split_complete": {
    "original_effort": "E1.1",
    "splits_created": ["E1.1a", "E1.1b"],
    "all_within_limit": true,
    "line_counts": {
      "E1.1a": 650,
      "E1.1b": 680
    }
  }
}
```

### 3. Update Master State
```bash
update_master_state_with_results() {
  local SUB_TYPE="$1"
  local RESULTS="$2"

  # Remove from active list
  remove_from_active_subs "$SUB_ID"

  # Add to history
  add_to_sub_history "$SUB_ID" "$RESULTS"

  # Update type-specific state
  case "$SUB_TYPE" in
    "FIX_CASCADE")
      update_fix_tracking "$RESULTS"
      ;;
    "INTEGRATE_WAVE_EFFORTS")
      update_integration_status "$RESULTS"
      ;;
    "SPLIT_COORDINATION")
      update_split_tracking "$RESULTS"
      ;;
  esac
}
```

### 4. Determine Next Actions
```bash
determine_next_action() {
  local SUB_OUTPUT="$1"
  local NEXT_ACTION=$(echo "$SUB_OUTPUT" | jq -r '.next_action')

  case "$NEXT_ACTION" in
    "CONTINUE")
      # Proceed with workflow
      proceed_to_next_phase
      ;;
    "RETRY")
      # Schedule retry
      schedule_sub_retry "$SUB_ID"
      ;;
    "ESCALATE")
      # Escalate to error recovery
      trigger_escalation "$SUB_ID"
      ;;
    "SPAWN_NEXT")
      # Spawn follow-up sub-orchestrator
      prepare_next_sub_spawn
      ;;
  esac
}
```

### 5. Clean Up Sub-Orchestrator Artifacts
```bash
cleanup_completed_sub() {
  local SUB_ID="$1"

  # Archive important files
  archive_sub_artifacts "$SUB_ID"

  # Remove temporary files
  rm -f "/tmp/params-${SUB_ID}.json"
  rm -f "/tmp/sub-orch-${SUB_ID}/heartbeat.json"

  # Kill any lingering processes
  cleanup_sub_processes "$SUB_ID"

  # Update tracking
  mark_sub_cleaned "$SUB_ID"
}
```

## Exit Criteria
- All completed sub-orchestrators processed
- Master state updated with results
- Next actions determined
- Cleanup completed

## Success Transitions
- Work continues → Appropriate next state based on workflow
- More subs needed → `MASTER_SPAWN_SUB_ORCHESTRATOR`
- All complete → Continue main workflow

## Failure Transitions
- Critical failure → `MASTER_RECOVER_SUB_FAILURE`
- Retry needed → `MASTER_SPAWN_SUB_ORCHESTRATOR`
- Escalation → `ERROR_RECOVERY`

## Rules Applied
- R377: Communication Protocol (output contracts)
- R378: Lifecycle Management (cleanup)
- R206: State Machine Validation
- R287: TODO Persistence

## Result Integration Patterns

### Success Pattern
```json
{
  "sub_completion": {
    "id": "uuid-123",
    "type": "FIX_CASCADE",
    "status": "PROJECT_DONE",
    "duration_seconds": 1800,
    "results_integrated": true,
    "next_action": "CONTINUE"
  }
}
```

### Failure Pattern
```json
{
  "sub_completion": {
    "id": "uuid-456",
    "type": "INTEGRATE_WAVE_EFFORTS",
    "status": "FAILED",
    "error": "Merge conflicts unresolvable",
    "retry_count": 2,
    "next_action": "RETRY"
  }
}
```

### Partial Success Pattern
```json
{
  "sub_completion": {
    "id": "uuid-789",
    "type": "FIX_CASCADE",
    "status": "PARTIAL",
    "succeeded": ["branch1"],
    "failed": ["branch2"],
    "next_action": "SPAWN_NEXT"
  }
}
```

## State Updates Required
```json
{
  "sub_orchestrator_completions": [
    {
      "id": "uuid-123",
      "completed_at": "timestamp",
      "status": "PROJECT_DONE",
      "results_integrated": true
    }
  ],
  "pending_actions": [
    {
      "action": "SPAWN_FIX_FOR_BRANCH2",
      "reason": "partial_completion",
      "scheduled": true
    }
  ]
}
```

## Archive Requirements
- Keep output files for 7 days
- Compress logs > 1MB
- Store critical results in state history
- Maintain audit trail

## Common Issues
1. **Missing output file** - Sub crashed before writing
2. **Corrupted results** - JSON parsing failures
3. **State conflicts** - Multiple subs updating same data
4. **Cleanup failures** - Permissions or locks

## Notes
- Process completions in order received
- Handle partial successes gracefully
- Preserve failure information for debugging
- Consider cascading effects of failures

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

