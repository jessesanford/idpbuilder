# 🔴🔴🔴 RULE R288 - MANDATORY STATE FILE UPDATE AND COMMIT PROTOCOL 🔴🔴🔴

## 🚨🚨🚨 BLOCKING: SUPREME LAW OF STATE PERSISTENCE 🚨🚨🚨

## THE ABSOLUTE DUAL REQUIREMENT:

**EVERY state transition MUST:**
1. **UPDATE** orchestrator-state.json IMMEDIATELY (within 30 seconds)
2. **COMMIT AND PUSH** EVERY SINGLE EDIT IMMEDIATELY (within 60 seconds)

**NO EXCEPTIONS. NO DEFERRALS. NO BATCHING. NO "LATER".**

## CRITICAL COMPANION RULE:
**R281 (SUPREME LAW #7)**: When creating the INITIAL state file, it MUST contain ALL phases, waves, and efforts from the implementation plan. See R281 for complete requirements.

## MANDATORY EXECUTION SEQUENCE:

```bash
# THE ONLY ACCEPTABLE PATTERN
perform_state_transition() {
    local OLD_STATE="$1"
    local NEW_STATE="$2"
    local REASON="$3"
    
    # Step 1: Validate transition is allowed
    validate_state_transition "$OLD_STATE" "$NEW_STATE"
    
    # Step 2: 🔴 UPDATE STATE FILE IMMEDIATELY 🔴
    # Use text_editor tool with str_replace to update orchestrator-state.json:
    # - Replace current_state value with NEW_STATE
    # - Replace previous_state value with OLD_STATE  
    # - Replace transition_time with current UTC timestamp
    # - Replace transition_reason with REASON
    # Agents should use multiple str_replace commands to update each field
    
    # Step 3: 🔴 COMMIT AND PUSH IMMEDIATELY 🔴
    git add orchestrator-state.json
    git commit -m "state: ${OLD_STATE} → ${NEW_STATE} - ${REASON} [R288]"
    git push
    
    # Step 4: Verify push succeeded
    if [ $? -ne 0 ]; then
        echo "🔴🔴🔴 CRITICAL: STATE PUSH FAILED! 🔴🔴🔴"
        git push --force-with-lease
        [ $? -ne 0 ] && exit 911  # Emergency exit
    fi
    
    # Step 5: Reload rules for new state (R217)
    reload_rules_for_state "$NEW_STATE"
    
    echo "✅ State transition complete and persisted: $OLD_STATE → $NEW_STATE"
}
```

## REQUIRED STATE FILE FIELDS:

### Core Fields (EVERY transition):
```yaml
state_machine:
  current_state: "NEW_STATE"
  previous_state: "OLD_STATE"
  transition_time: "2025-08-25T12:00:00Z"
  transition_reason: "Clear explanation"
  rules_reacknowledged: true  # After R217 compliance
```

### State-Specific Updates:
- **WAVE_COMPLETE**: Add waves_completed entry with metrics
- **INTEGRATION**: Add current_integration details
- **ERROR_RECOVERY**: Add error_context information
- **SUCCESS**: Add phase_completion summary
- **SPAWN_AGENTS**: Add agents_spawned records
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
- `state: PLANNING → SETUP_EFFORT_INFRASTRUCTURE - planning complete [R288]`
- `state: wave1 marked complete - all efforts reviewed [R288]`
- `state: effort1 status=BLOCKED - size limit exceeded [R288]`

## MANDATORY WRAPPER FUNCTION:

```bash
# ALL state updates MUST use this wrapper
update_and_commit_state() {
    local KEY="$1"
    local VALUE="$2"
    local REASON="${3:-update}"
    
    # 1. Make the edit
    # Use text_editor tool with str_replace to update orchestrator-state.json:
    # - Replace the KEY's current value with VALUE
    # Example: str_replace to change KEY: old_value to KEY: new_value
    
    # 2. IMMEDIATELY commit and push
    git add orchestrator-state.json
    git commit -m "state: ${KEY}=${VALUE} - ${REASON} [R288]"
    git push
    
    # 3. Verify success
    if [ $? -ne 0 ]; then
        echo "🔴 R288 VIOLATION: Failed to push state!"
        exit 288
    fi
}
```

## COMMON VIOLATIONS:

### ❌ FORBIDDEN PATTERNS:
```bash
# ❌ NO: Deferred commit
# Use text_editor tool with str_replace to update orchestrator-state.json:
# Replace current_state: "old_value" with current_state: "INTEGRATION"
do_other_work()  # VIOLATION! Must commit first!

# ❌ NO: Batch updates
# Use text_editor tool with str_replace to update orchestrator-state.json:
# - Replace current_state value with "WAVE_COMPLETE"
# - Replace wave1.status value with "COMPLETE"
git add orchestrator-state.json  # VIOLATION! Each edit needs commit!

# ❌ NO: Missing push
git add orchestrator-state.json
git commit -m "state: update"
# No push = VIOLATION!
```

### ✅ REQUIRED PATTERN:
```bash
# ✅ YES: Immediate update, commit, push
# Use text_editor tool with str_replace to update orchestrator-state.json:
# Replace current_state: "old_value" with current_state: "INTEGRATION"
git add orchestrator-state.json
git commit -m "state: transition to INTEGRATION [R288]"
git push

# ✅ YES: Each edit gets its own commit
update_and_commit_state "current_state" "WAVE_COMPLETE" "all efforts done"
update_and_commit_state "wave1.status" "COMPLETE" "wave finished"
```

## COMPLIANCE MONITORING:

```bash
check_r288_compliance() {
    # Check for uncommitted state changes
    if git status --porcelain | grep -q "orchestrator-state.json"; then
        echo "❌❌❌ R288 VIOLATION: Uncommitted state changes!"
        return 288
    fi
    
    # Check timestamp freshness
    # Use text_editor tool with view command to read orchestrator-state.json:
    # Find the state_machine.transition_time field
    local TIMESTAMP="<value from state_machine.transition_time>"
    local NOW=$(date +%s)
    local TRANS_TIME=$(date -d "$TIMESTAMP" +%s 2>/dev/null || echo 0)
    local AGE=$((NOW - TRANS_TIME))
    
    if [ $AGE -gt 60 ]; then
        echo "⚠️ R288 WARNING: State timestamp stale (${AGE}s old)"
    fi
    
    echo "✅ R288 Compliance: OK"
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
   # Use text_editor tool to increment r288_violations in orchestrator-state.json:
   # First view to get current value, then str_replace with incremented value
   git add orchestrator-state.json
   git commit -m "state: R288 violation logged [R288]"
   git push
   ```

## WHY THIS IS CRITICAL:

1. **State Loss Prevention**: Crashes don't lose progress
2. **Multi-Instance Safety**: All instances see current state
3. **Recovery Capability**: Can resume from exact point
4. **Complete Audit Trail**: Every transition tracked
5. **Debugging Support**: Full state history available

## THE SUPREME LAWS:

1. **No state transition is complete until it's in the file**
2. **No state update is safe until it's pushed to remote**
3. **Every edit deserves its own commit**
4. **The state file is the single source of truth**
5. **If it's not pushed, it didn't happen**

## FINAL WARNING:

**This rule consolidates and supersedes R252 and R253.**

**VIOLATION = FAILURE. NO EXCEPTIONS.**

**UPDATE → COMMIT → PUSH or FAIL!**