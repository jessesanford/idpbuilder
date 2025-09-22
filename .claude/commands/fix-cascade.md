
## 🔴🔴🔴 SUB-STATE MACHINE ENTRY POINT 🔴🔴🔴

**FIX CASCADE NOW USES SUB-STATE MACHINE ARCHITECTURE**

### Check for Active Sub-State
```bash
# Check if already in fix cascade sub-state
check_fix_cascade_active() {
    local ACTIVE=$(jq -r '.sub_state_machine.active' orchestrator-state.json)
    local TYPE=$(jq -r '.sub_state_machine.type' orchestrator-state.json)

    if [[ "$ACTIVE" == "true" && "$TYPE" == "FIX_CASCADE" ]]; then
        FIX_STATE_FILE=$(jq -r '.sub_state_machine.state_file' orchestrator-state.json)
        FIX_CURRENT_STATE=$(jq -r '.current_state' "$FIX_STATE_FILE")
        echo "✅ Fix cascade already active"
        echo "📄 State file: $FIX_STATE_FILE"
        echo "📍 Current state: $FIX_CURRENT_STATE"
        return 0
    fi
    return 1
}
```

### Enter Fix Cascade Sub-State Machine
```bash
# Start new fix cascade by entering sub-state machine
enter_fix_cascade_sub_state() {
    local FIX_ID="$1"
    local RETURN_STATE=$(jq -r '.current_state' orchestrator-state.json)

    echo "🔄 Entering Fix Cascade Sub-State Machine"

    # Create fix state file per R375
    FIX_STATE_FILE="orchestrator-${FIX_ID}-state.json"
    cat > "$FIX_STATE_FILE" << EOF
{
  "sub_state_type": "FIX_CASCADE",
  "fix_identifier": "${FIX_ID}",
  "current_state": "FIX_CASCADE_INIT",
  "created_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "parent_state_machine": {
    "state_file": "orchestrator-state.json",
    "return_state": "${RETURN_STATE}",
    "nested_level": 1
  }
}
EOF

    # Update main state to show sub-state active
    jq --arg file "$FIX_STATE_FILE" \
       --arg return "$RETURN_STATE" \
       '.sub_state_machine = {
          "active": true,
          "type": "FIX_CASCADE",
          "state_file": $file,
          "current_state": "FIX_CASCADE_INIT",
          "return_state": $return,
          "started_at": now,
          "trigger_reason": "Fix cascade initiated via /fix-cascade command"
       }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Commit both state files
    git add orchestrator-state.json "$FIX_STATE_FILE"
    git commit -m "fix-cascade: enter sub-state machine for ${FIX_ID}"
    git push

    echo "✅ Entered Fix Cascade Sub-State Machine"
    echo "📄 Fix state file: $FIX_STATE_FILE"
    echo "🔄 Will return to: $RETURN_STATE after completion"
}
```

## 🔴🔴🔴 R376 ENFORCEMENT - FIX CASCADE QUALITY GATES 🔴🔴🔴

### GATE 1: POST-BACKPORT REVIEW
```bash
# MANDATORY after EVERY backport operation
after_backport_review() {
    local branch="$1"
    echo "🔴 GATE 1: Post-Backport Review (R376/R354)"

    # Mark review required in state
    jq --arg branch "$branch" '
        .quality_gates.gate_1_backport.status = "pending" |
        .quality_gates.gate_1_backport.current_branch = $branch |
        .quality_gates.current_gate = "gate_1_backport"
    ' "$FIX_STATE_FILE" > tmp.json && mv tmp.json "$FIX_STATE_FILE"

    # Spawn Code Reviewer with fix-cascade mode
    echo "🚨 BLOCKING: Cannot continue until review completes"
    echo ""
    echo "Task: code-reviewer"
    echo "State: FIX_CASCADE_REVIEW"
    echo "Review Type: post-backport"
    echo "Branch: $branch"
    echo "Focus: fix-correctness, build-success, test-pass"
    echo ""

    # Cannot proceed until review passes
    echo "⏸️ FIX CASCADE PAUSED - Awaiting mandatory review per R376"
}
```

### GATE 2: POST-FORWARD-PORT REVIEW
```bash
# MANDATORY after EVERY forward-port operation
after_forward_port_review() {
    local branch="$1"
    echo "🔴 GATE 2: Post-Forward-Port Review (R376/R354)"

    # Mark review required in state
    jq --arg branch "$branch" '
        .quality_gates.gate_2_forward_port.status = "pending" |
        .quality_gates.gate_2_forward_port.current_branch = $branch |
        .quality_gates.current_gate = "gate_2_forward_port"
    ' "$FIX_STATE_FILE" > tmp.json && mv tmp.json "$FIX_STATE_FILE"

    # Spawn Code Reviewer
    echo "🚨 BLOCKING: Cannot continue until review completes"
    echo ""
    echo "Task: code-reviewer"
    echo "State: FIX_CASCADE_REVIEW"
    echo "Review Type: post-forward-port"
    echo "Branch: $branch"
    echo "Focus: integration-correctness, no-regressions"
    echo ""

    echo "⏸️ FIX CASCADE PAUSED - Awaiting mandatory review per R376"
}
```

### GATE 3: CONFLICT RESOLUTION REVIEW
```bash
# MANDATORY after resolving any conflicts
after_conflict_resolution_review() {
    local branch="$1"
    echo "🔴 GATE 3: Conflict Resolution Review (R376/R354)"

    jq --arg branch "$branch" '
        .quality_gates.gate_3_conflict_resolution.status = "pending" |
        .quality_gates.gate_3_conflict_resolution.current_branch = $branch |
        .quality_gates.current_gate = "gate_3_conflict_resolution"
    ' "$FIX_STATE_FILE" > tmp.json && mv tmp.json "$FIX_STATE_FILE"

    echo "🚨 BLOCKING: Conflict resolution must be reviewed"
    echo ""
    echo "Task: code-reviewer"
    echo "State: FIX_CASCADE_REVIEW"
    echo "Review Type: conflict-resolution"
    echo "Branch: $branch"
    echo "Focus: correct-resolution, no-lost-code, both-sides-merged"
    echo ""
}
```

### GATE 4: COMPREHENSIVE FINAL VALIDATION
```bash
# MANDATORY before completing fix cascade
comprehensive_validation_gate() {
    echo "🔴 GATE 4: Comprehensive Fix Validation (R376)"

    jq '.quality_gates.gate_4_comprehensive.status = "pending" |
        .quality_gates.current_gate = "gate_4_comprehensive"
    ' "$FIX_STATE_FILE" > tmp.json && mv tmp.json "$FIX_STATE_FILE"

    echo "📊 COMPREHENSIVE VALIDATION REQUIRED:"
    echo "  ✓ All branches must build successfully"
    echo "  ✓ All tests must pass on all branches"
    echo "  ✓ Fix must resolve original issue"
    echo "  ✓ No regressions introduced"
    echo "  ✓ Code quality maintained"
    echo ""
    echo "Task: code-reviewer"
    echo "State: FIX_CASCADE_COMPREHENSIVE_REVIEW"
    echo "Review Type: comprehensive-validation"
    echo "All Branches: true"
    echo "Checklist: build,test,functional,fix-verification,no-regressions"
    echo ""

    echo "⏸️ CANNOT COMPLETE until comprehensive validation passes"
}
```

### QUALITY GATE ENFORCEMENT
```bash
# Check if we can proceed based on quality gates
check_quality_gates() {
    local current_gate=$(jq -r '.quality_gates.current_gate' "$FIX_STATE_FILE")
    local gate_status=$(jq -r ".quality_gates.${current_gate}.status" "$FIX_STATE_FILE")

    if [[ "$gate_status" == "pending" || "$gate_status" == "in_progress" ]]; then
        echo "🛑 R376 VIOLATION: Quality gate ${current_gate} not passed!"
        echo "❌ CANNOT CONTINUE fix cascade until gate passes"
        return 376  # R376 violation code
    fi

    if [[ "$gate_status" == "failed" ]]; then
        echo "❌ Quality gate ${current_gate} FAILED - fixes required"
        echo "Must address issues before continuing"
        return 1
    fi

    echo "✅ Quality gate ${current_gate} passed"
    return 0
}

# Call quality gate checks at state transitions
check_quality_gates || exit 376
```

## 📊 PROGRESS TRACKING

```bash
show_fix_progress() {
    echo "═══════════════════════════════════════════════════════"
    echo "📊 FIX CASCADE PROGRESS: $FIX_NAME"
    echo "═══════════════════════════════════════════════════════"
    
    # Show backport progress
    BACKPORT_DONE=$(jq '[.backport_status[] | select(. == "COMPLETED")] | length' \
        "$FIX_STATE_FILE")
    BACKPORT_TOTAL=$(jq '.backport_status | length' "$FIX_STATE_FILE")
    echo "🔧 Backports: $BACKPORT_DONE/$BACKPORT_TOTAL"
    
    # Show forward port progress
    FORWARD_DONE=$(jq '.forward_port_status.branches_completed | length' \
        "$FIX_STATE_FILE")
    FORWARD_PENDING=$(jq '.forward_port_status.branches_pending | length' \
        "$FIX_STATE_FILE")
    FORWARD_TOTAL=$((FORWARD_DONE + FORWARD_PENDING))
    echo "📦 Forward Ports: $FORWARD_DONE/$FORWARD_TOTAL"
    
    # Show validation status
    echo "✅ Validation:"
    jq '.validation_status' "$FIX_STATE_FILE"
    
    echo "═══════════════════════════════════════════════════════"
}

show_fix_progress
```

## 💾 CHECKPOINT MANAGEMENT (R375)

```bash
save_fix_checkpoint() {
    echo "💾 Saving fix cascade checkpoint (R375)..."
    
    # Update timestamp
    jq '.timestamps.last_updated = now' "$FIX_STATE_FILE" > tmp.json && \
        mv tmp.json "$FIX_STATE_FILE"
    
    # Commit to planning repo (not target repo)
    git add "$FIX_STATE_FILE" "${FIX_NAME}-"*.md
    git commit -m "fix-cascade: ${FIX_NAME} - checkpoint at $FIX_CURRENT_STATE"
    git push
    
    echo "✅ Checkpoint saved - Resume with: /fix-cascade"
}

save_fix_checkpoint
```

## ⚠️ ERROR RECOVERY

```bash
handle_fix_error() {
    local ERROR_MSG="$1"
    
    echo "❌ Fix error: $ERROR_MSG"
    
    # Log to fix state
    jq --arg error "$ERROR_MSG" \
       '.error_log += [{
           "timestamp": now,
           "state": "'$FIX_CURRENT_STATE'",
           "error": $error
       }]' "$FIX_STATE_FILE" > tmp.json && mv tmp.json "$FIX_STATE_FILE"
    
    # R327: Check for nested fix cascade
    if [[ "$ERROR_MSG" == *"new issue discovered"* ]]; then
        echo "🔄 R327: Nested fix cascade may be needed"
        echo "Create new fix plan for discovered issue"
    fi
}
```

## 🎯 USAGE

```bash
# Start a new fix cascade
1. Create fix plan: BUGNAME-FIX-PLAN.md
2. Run: /fix-cascade
3. Command auto-detects fix and initializes state

# Continue existing fix cascade
1. Run: /fix-cascade
2. Command resumes from last checkpoint

# Multiple concurrent fixes
- Each fix gets its own state file per R375
- Main orchestrator state remains clean
- Fixes tracked independently
```

## 📋 COMPLETION CRITERIA

Fix cascade is complete when:
- [ ] All backports applied and reviewed
- [ ] All forward ports completed
- [ ] Validation tests pass
- [ ] Integration report created
- [ ] Fix state archived per R375
- [ ] Main state updated

---

**Remember**: R375 mandates dual state tracking - main state stays clean, fix state tracks details.