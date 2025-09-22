# 🚨🚨🚨 BLOCKING RULE R377: Master-Sub Orchestrator Communication Protocol 🚨🚨🚨

## Rule Definition
**ID**: R377
**Category**: ORCHESTRATION
**Criticality**: BLOCKING
**Created**: 2025-01-21
**Dependencies**: R206 (State Machine Validation), R313 (Stop After Spawn)

## Description
Defines the mandatory communication protocol between master orchestrators and sub-state orchestrators. All communication MUST occur through well-defined file-based contracts with explicit input/output schemas.

## 🔴 MANDATORY COMMUNICATION PATTERN

### 1. SPAWNING SUB-ORCHESTRATORS

```bash
# Master MUST use existing commands in background mode
# Each command already handles its own sub-state machine

# For FIX_CASCADE:
claude -p "$CLAUDE_PROJECT_DIR" \
  --command "/fix-cascade" \
  --params "fix_id=[fix-id] branches=[branch-list]" \
  > /tmp/fix-cascade-[id].log 2>&1 &

# For SPLITTING:
claude -p "$CLAUDE_PROJECT_DIR" \
  --command "/splitting" \
  --params "effort=[effort-name] size=[measured-size]" \
  > /tmp/splitting-[id].log 2>&1 &

SUB_PID=$!
```

### 2. PARAMETER PASSING

```bash
# Parameters are passed directly to commands
# Each command handles its own state management

# FIX_CASCADE parameters:
fix_id="fix-bug-123"  # Unique identifier for the fix
branches="branch1,branch2,branch3"  # Comma-separated branch list

# SPLITTING parameters:
effort="effort2-controller"  # Effort name that exceeded size
size="1250"  # Measured line count

# Commands create their own sub-state files:
# - orchestrator-${fix_id}-state.json (for fix-cascade)
# - splitting-${effort}-state.json (for splitting)
```

### 3. HEARTBEAT CONTRACT

Sub-orchestrators MUST write heartbeat every 30 seconds:

```json
{
  "pid": 12345,
  "status": "RUNNING|COMPLETING|FAILED",
  "progress_percentage": 75,
  "current_state": "FIX_CASCADE_BACKPORT",
  "last_update": "2025-01-21T10:30:00Z",
  "estimated_completion": "2025-01-21T10:45:00Z",
  "messages": [],
  "checkpoint_data": {}
}
```

### 4. OUTPUT CONTRACT

Upon completion, sub-orchestrator MUST write:

```json
{
  "status": "SUCCESS|FAILED|PARTIAL",
  "completion_timestamp": "2025-01-21T10:45:00Z",
  "duration_seconds": 900,
  "results": {
    // Type-specific results
  },
  "errors": [],
  "next_action": "CONTINUE|RETRY|ESCALATE",
  "state_transitions_executed": [],
  "artifacts_created": []
}
```

## 🛑 SUB-ORCHESTRATOR TYPE CONTRACTS

### FIX_CASCADE CONTRACT

**Input:**
```json
{
  "fix_id": "fix-bug-123",
  "branches_to_fix": ["branch1", "branch2"],
  "issue_description": "Detailed issue description",
  "validation_requirements": {
    "tests_must_pass": true,
    "linting_required": true,
    "review_required": false
  },
  "source_fix": {
    "branch": "fix/original-branch",
    "commit": "abc123"
  }
}
```

**Output:**
```json
{
  "status": "SUCCESS",
  "branches_fixed": ["branch1", "branch2"],
  "test_results": {
    "branch1": "PASSED",
    "branch2": "PASSED"
  },
  "commits_created": {
    "branch1": "def456",
    "branch2": "ghi789"
  },
  "validation_passed": true
}
```

### INTEGRATION CONTRACT

**Input:**
```json
{
  "integration_type": "WAVE|PHASE|PROJECT",
  "branches_to_merge": ["effort-1", "effort-2"],
  "target_branch": "wave-1-integration",
  "validation_level": "BASIC|FULL|COMPREHENSIVE",
  "conflict_resolution": "MANUAL|AUTO_THEIRS|AUTO_OURS"
}
```

**Output:**
```json
{
  "status": "SUCCESS",
  "integration_branch": "wave-1-integration",
  "merge_results": {
    "effort-1": "MERGED",
    "effort-2": "MERGED"
  },
  "conflicts_resolved": [],
  "test_status": "PASSED",
  "build_status": "SUCCESS"
}
```

### SPLIT_COORDINATION CONTRACT

**Input:**
```json
{
  "original_effort": "E1.1",
  "split_plan": {
    "splits": ["E1.1a", "E1.1b", "E1.1c"],
    "dependencies": [],
    "sequential": true
  },
  "size_limit": 700
}
```

**Output:**
```json
{
  "status": "SUCCESS",
  "splits_completed": ["E1.1a", "E1.1b", "E1.1c"],
  "line_counts": {
    "E1.1a": 650,
    "E1.1b": 680,
    "E1.1c": 450
  },
  "all_within_limit": true
}
```

## ⚠️ MONITORING REQUIREMENTS

### Master MUST:
1. Check heartbeat file every 15 seconds
2. Detect stale processes (no heartbeat > 2 minutes)
3. Maintain process tracking in master state
4. Handle sub-orchestrator failures gracefully
5. Clean up zombie processes

### Sub-Orchestrator MUST:
1. Write heartbeat every 30 seconds
2. Update progress percentage accurately
3. Write checkpoint data for recovery
4. Signal completion explicitly
5. Clean exit with proper status code

## 🔴 FAILURE HANDLING

### Stale Process Detection
```bash
# If heartbeat older than 2 minutes
if [[ $HEARTBEAT_AGE -gt 120 ]]; then
  # Process is stale, attempt recovery
  if [[ -f "$CHECKPOINT_FILE" ]]; then
    # Can recover from checkpoint
    spawn_recovery_sub_orchestrator
  else
    # Must restart from beginning
    mark_as_failed_and_retry
  fi
fi
```

### Maximum Retry Policy
- 3 retries for transient failures
- 1 retry for configuration errors
- 0 retries for invalid state errors

## 🚨 VIOLATIONS

### BLOCKING Violations (-100% penalty):
1. Direct variable passing instead of file-based communication
2. Missing heartbeat for > 5 minutes
3. No output file written after completion
4. Invalid status values in contracts
5. Missing mandatory contract fields

### WARNING Violations (-25% penalty):
1. Heartbeat interval > 45 seconds
2. Progress percentage not updated
3. Missing checkpoint data
4. No estimated completion time
5. Incomplete error reporting

## 📋 Implementation Checklist

- [ ] Parameter file created before spawn
- [ ] Heartbeat file initialized
- [ ] Output location prepared
- [ ] Process tracking added to master state
- [ ] Monitoring loop established
- [ ] Completion detection implemented
- [ ] Cleanup procedures in place
- [ ] Recovery mechanisms tested

## 🔗 Related Rules
- R378: Sub-Orchestrator Lifecycle Management
- R379: Sub-Process Monitoring and Recovery
- R313: Mandatory Stop After Spawn
- R206: State Machine Validation

## 📝 Examples

### Correct Master Spawn:
```bash
# Prepare parameters
cat > /tmp/fix-params-$UUID.json <<EOF
{
  "sub_orchestrator_type": "FIX_CASCADE",
  "unique_id": "$UUID",
  "input_parameters": {...}
}
EOF

# Spawn fix-cascade command in background
claude -p . --command "/fix-cascade" \
  --params "fix_id=$FIX_ID branches=$BRANCHES" &

# Track in state
update_master_state "sub_orchestrator_spawned" "$!"
```

### Correct Sub-Orchestrator Heartbeat:
```bash
write_heartbeat() {
  cat > "$HEARTBEAT_FILE" <<EOF
{
  "pid": $$,
  "status": "RUNNING",
  "progress_percentage": $PROGRESS,
  "current_state": "$CURRENT_STATE",
  "last_update": "$(date -Iseconds)"
}
EOF
}
```

---

**Remember**: All communication between master and sub-orchestrators MUST be file-based. No direct process communication, shared memory, or environment variables allowed!