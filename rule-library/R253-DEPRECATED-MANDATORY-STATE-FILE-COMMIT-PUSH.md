# ⚠️ DEPRECATED - Subsumed by R288
This rule has been consolidated into R288-state-file-update-and-commit-protocol.md
Please refer to R288 for current state file update and commit requirements.

---

# 🔴🔴🔴 RULE R253 - MANDATORY STATE FILE COMMIT AND PUSH 🔴🔴🔴

## ⚠️⚠️⚠️ SUPREME LAW - COMMIT AND PUSH ON EVERY EDIT ⚠️⚠️⚠️

## THE ABSOLUTE LAW:

**EVERY SINGLE EDIT TO orchestrator-state.yaml MUST BE IMMEDIATELY COMMITTED AND PUSHED!**

**NO EXCEPTIONS! NO DEFERRALS! NO "BATCH LATER"! NO "BUT FIRST"!**

## WORKS WITH R252:
**R252** defines WHAT must be updated in the state file during transitions. **R253** enforces that EVERY edit must be immediately committed and pushed. Together they ensure complete, persistent state tracking.

## WHY THIS IS CRITICAL:

1. **State Loss Prevention**: If the orchestrator crashes, state is preserved
2. **Multi-Instance Safety**: Other instances can see current state
3. **Recovery Capability**: Can resume from exact last state
4. **Audit Trail**: Complete history of all state transitions
5. **Debugging**: Can trace exact sequence of state changes

## MANDATORY PROTOCOL:

```bash
# EVERY TIME YOU EDIT orchestrator-state.yaml, YOU MUST:
commit_and_push_state() {
    local CHANGE_DESCRIPTION="$1"
    
    # 1. IMMEDIATELY stage the file
    git add orchestrator-state.yaml
    
    # 2. IMMEDIATELY commit with descriptive message
    git commit -m "state: $CHANGE_DESCRIPTION [R253]"
    
    # 3. IMMEDIATELY push to remote
    git push
    
    # 4. VERIFY push succeeded
    if [ $? -ne 0 ]; then
        echo "🔴🔴🔴 CRITICAL: STATE PUSH FAILED! 🔴🔴🔴"
        echo "RETRY IMMEDIATELY OR FAIL!"
        git push --force-with-lease  # Try with lease
        if [ $? -ne 0 ]; then
            echo "❌❌❌ FATAL: Cannot push state! STOP ALL WORK! ❌❌❌"
            exit 911  # Emergency exit
        fi
    fi
    
    echo "✅ State committed and pushed: $CHANGE_DESCRIPTION"
}
```

## ENFORCEMENT EXAMPLES:

### ✅ CORRECT - Immediate Commit/Push:
```bash
# Update state
yq -i '.current_state = "WAVE_COMPLETE"' orchestrator-state.yaml
git add orchestrator-state.yaml
git commit -m "state: transition to WAVE_COMPLETE [R253]"
git push

# Update wave completion
yq -i '.waves_completed.phase1.wave1.status = "COMPLETE"' orchestrator-state.yaml
git add orchestrator-state.yaml
git commit -m "state: mark phase1/wave1 complete [R253]"
git push

# Update effort tracking
yq -i '.efforts_in_progress.effort1.status = "REVIEWING"' orchestrator-state.yaml
git add orchestrator-state.yaml
git commit -m "state: effort1 now in review [R253]"
git push
```

### ❌ WRONG - Deferred or Batch Commits:
```bash
# ❌❌❌ NEVER DO THIS:
yq -i '.current_state = "WAVE_COMPLETE"' orchestrator-state.yaml
yq -i '.waves_completed.phase1.wave1.status = "COMPLETE"' orchestrator-state.yaml
# ... do other work ...
git add orchestrator-state.yaml  # TOO LATE! VIOLATION!
git commit -m "state: various updates"  # BATCH = FAIL!

# ❌❌❌ NEVER DO THIS:
yq -i '.current_state = "INTEGRATION"' orchestrator-state.yaml
echo "I'll commit this after I finish integration"  # NO! VIOLATION!

# ❌❌❌ NEVER DO THIS:
for i in 1 2 3; do
    yq -i ".effort$i.status = 'complete'" orchestrator-state.yaml
done
git add orchestrator-state.yaml  # BATCHING = VIOLATION!
```

## INTEGRATION WITH OTHER FUNCTIONS:

### Update ALL State Functions to Include Commit/Push:

```bash
# MANDATORY: Wrapper for ALL yq edits to state file
update_orchestrator_state() {
    local KEY="$1"
    local VALUE="$2"
    local REASON="${3:-update}"
    
    # 1. Make the edit
    yq -i ".${KEY} = \"${VALUE}\"" orchestrator-state.yaml
    
    # 2. IMMEDIATELY commit and push (R253)
    git add orchestrator-state.yaml
    git commit -m "state: ${KEY}=${VALUE} - ${REASON} [R253]"
    git push
    
    # 3. Verify success
    if [ $? -ne 0 ]; then
        echo "🔴 R253 VIOLATION: Failed to push state!"
        exit 253
    fi
}

# MANDATORY: Enhanced wave completion with auto-commit
mark_wave_complete() {
    local PHASE="$1"
    local WAVE="$2"
    
    # Multiple edits, multiple commits (R253)
    yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
    git add orchestrator-state.yaml && git commit -m "state: wave${WAVE} completion time [R253]" && git push
    
    yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.status = \"COMPLETE\"" orchestrator-state.yaml
    git add orchestrator-state.yaml && git commit -m "state: wave${WAVE} status COMPLETE [R253]" && git push
    
    yq -i ".current_wave_complete = true" orchestrator-state.yaml
    git add orchestrator-state.yaml && git commit -m "state: current wave marked complete [R253]" && git push
}
```

## COMMIT MESSAGE FORMAT:

All state commits MUST follow this format:
```
state: <what changed> - <why> [R253]
```

Examples:
- `state: current_state=INTEGRATION - wave complete, starting integration [R253]`
- `state: effort1 status=BLOCKED - size limit exceeded [R253]`
- `state: phase1/wave1 complete - all efforts reviewed [R253]`
- `state: spawn architect-reviewer-123 - wave review needed [R253]`

## FREQUENCY REQUIREMENT:

**EVERY SINGLE EDIT** means:
- After EVERY yq command that touches orchestrator-state.yaml
- After EVERY state transition
- After EVERY status update
- After EVERY effort tracking change
- After EVERY completion marker
- After EVERY error recording

**NOT** at the end of a sequence!
**NOT** after multiple changes!
**NOT** when convenient!
**IMMEDIATELY** after each edit!

## RECOVERY PROTOCOL:

If you realize you've violated R253:
1. **STOP** everything immediately
2. **COMMIT** the current state with message: `state: RECOVERY - missed commits detected [R253-VIOLATION]`
3. **PUSH** immediately
4. **LOG** the violation in the state file itself:
   ```bash
   yq -i ".r253_violations += 1" orchestrator-state.yaml
   git add orchestrator-state.yaml
   git commit -m "state: R253 violation logged [R253]"
   git push
   ```

## GRADING IMPACT:

**AUTOMATIC FAILURE CONDITIONS:**
- Any state edit without immediate commit/push
- Batch commits of multiple state changes
- State file with uncommitted changes for >30 seconds
- Missing [R253] tag in state commit messages
- Any "git status" showing modified orchestrator-state.yaml

**GRADING PENALTIES:**
- First violation: -20% on state management score
- Second violation: -50% on state management score
- Third violation: AUTOMATIC FAIL of entire orchestration

## MONITORING COMPLIANCE:

```bash
# Check for R253 compliance
check_r253_compliance() {
    # Check for uncommitted changes
    if git status --porcelain | grep -q "orchestrator-state.yaml"; then
        echo "❌❌❌ R253 VIOLATION: Uncommitted state changes detected!"
        echo "COMMIT AND PUSH IMMEDIATELY!"
        return 253
    fi
    
    # Check recent commits for R253 tag
    local RECENT=$(git log --oneline -n 10 | grep "orchestrator-state.yaml")
    if ! echo "$RECENT" | grep -q "\[R253\]"; then
        echo "⚠️ Warning: Recent state commits missing [R253] tag"
    fi
    
    echo "✅ R253 Compliance: OK"
}
```

## THE GOLDEN RULES:

1. **EDIT → COMMIT → PUSH** (always in this order, always immediate)
2. **One Edit = One Commit** (never batch)
3. **No State Left Behind** (every change tracked)
4. **Push or Perish** (local commits are not enough)
5. **[R253] Tag Required** (mark your compliance)

## VERIFICATION CHECKLIST:

Before ANY state operation:
- [ ] Git repository is clean (no uncommitted changes)
- [ ] Remote is accessible and pushable
- [ ] You understand R253 requirements

After EVERY state edit:
- [ ] File immediately staged with `git add`
- [ ] Commit created with descriptive message and [R253] tag
- [ ] Changes pushed to remote successfully
- [ ] No uncommitted changes remain

## FINAL WARNING:

**This is not a suggestion. This is not a best practice. This is MANDATORY.**

**Every violation is tracked. Every violation affects grading. Every violation risks state loss.**

**COMMIT AND PUSH ON EVERY EDIT OR FAIL!**