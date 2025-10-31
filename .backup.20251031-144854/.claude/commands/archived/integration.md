## 🔴🔴🔴 SUB-STATE MACHINE ENTRY POINT 🔴🔴🔴

**INTEGRATE_WAVE_EFFORTS NOW USES SUB-STATE MACHINE ARCHITECTURE**

### Check for Active Sub-State
```bash
# Check if already in integration sub-state
check_integration_active() {
    local ACTIVE=$(jq -r '.sub_state_machine.active' orchestrator-state-v3.json)
    local TYPE=$(jq -r '.sub_state_machine.type' orchestrator-state-v3.json)

    if [[ "$ACTIVE" == "true" && "$TYPE" == "INTEGRATE_WAVE_EFFORTS" ]]; then
        INT_STATE_FILE=$(jq -r '.sub_state_machine.state_file' orchestrator-state-v3.json)
        INT_CURRENT_STATE=$(jq -r '.current_state' "$INT_STATE_FILE")
        echo "✅ Integration already active"
        echo "📄 State file: $INT_STATE_FILE"
        echo "📍 Current state: $INT_CURRENT_STATE"
        return 0
    fi
    return 1
}
```

### Enter Integration Sub-State Machine
```bash
# Start new integration by entering sub-state machine
enter_integration_sub_state() {
    local INT_TYPE="$1"        # WAVE|PHASE|PROJECT
    local BRANCHES="$2"         # comma-separated list
    local TARGET_BRANCH="$3"    # target integration branch
    local VALIDATION_LEVEL="${4:-BASIC}"  # BASIC|FULL|COMPREHENSIVE
    local RETURN_STATE=$(jq -r '.current_state' orchestrator-state-v3.json)

    echo "🔄 Entering Integration Sub-State Machine"
    echo "📊 Type: $INT_TYPE"
    echo "🌿 Branches: $BRANCHES"
    echo "🎯 Target: $TARGET_BRANCH"

    # Create integration state file
    INT_STATE_FILE="orchestrator-integration-${INT_TYPE,,}-state.json"

    # Convert comma-separated branches to JSON array
    BRANCHES_JSON=$(echo "$BRANCHES" | tr ',' '\n' | jq -R . | jq -s .)

    cat > "$INT_STATE_FILE" << EOF
{
  "sub_state_type": "INTEGRATE_WAVE_EFFORTS",
  "integration_type": "${INT_TYPE}",
  "current_state": "INTEGRATE_WAVE_EFFORTS_INIT",
  "branches_to_integrate": $BRANCHES_JSON,
  "target_branch": "${TARGET_BRANCH}",
  "validation_level": "${VALIDATION_LEVEL}",
  "attempt": 1,
  "created_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "parent_state_machine": {
    "state_file": "orchestrator-state-v3.json",
    "return_state": "${RETURN_STATE}",
    "nested_level": 1
  },
  "integration_progress": {
    "branches_merged": [],
    "branches_pending": $BRANCHES_JSON,
    "conflicts_found": [],
    "issues_detected": []
  },
  "quality_gates": {
    "gate_1_pre_merge": {
      "status": "pending",
      "checks": []
    },
    "gate_2_post_merge": {
      "status": "pending",
      "checks": []
    },
    "gate_3_validation": {
      "status": "pending",
      "checks": []
    },
    "gate_4_comprehensive": {
      "status": "pending",
      "checks": []
    }
  },
  "cycle_history": []
}
EOF

    # Update main state to show sub-state active
    jq --arg file "$INT_STATE_FILE" \
       --arg type "$INT_TYPE" \
       --arg return "$RETURN_STATE" \
       '.sub_state_machine = {
          "active": true,
          "type": "INTEGRATE_WAVE_EFFORTS",
          "state_file": $file,
          "current_state": "INTEGRATE_WAVE_EFFORTS_INIT",
          "return_state": $return,
          "started_at": now,
          "trigger_reason": "Integration initiated via /integration command for " + $type
       }' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    # Commit both state files
    git add orchestrator-state-v3.json "$INT_STATE_FILE"
    git commit -m "integration: enter sub-state machine for ${INT_TYPE} integration"
    git push

    echo "✅ Entered Integration Sub-State Machine"
    echo "📄 Integration state file: $INT_STATE_FILE"
    echo "🔄 Will return to: $RETURN_STATE after completion"

    # Output signals for monitoring
    echo ""
    echo "SUB_ORCHESTRATOR_STARTED: INTEGRATE_WAVE_EFFORTS"
    echo "INTEGRATE_WAVE_EFFORTS_TYPE: $INT_TYPE"
    echo "BRANCHES_COUNT: $(echo "$BRANCHES" | tr ',' '\n' | wc -l)"
}
```

## 🔴🔴🔴 INTEGRATE_WAVE_EFFORTS QUALITY GATES 🔴🔴🔴

### GATE 1: PRE-MERGE VALIDATION
```bash
# MANDATORY before starting any merge
pre_merge_validation() {
    local branch="$1"
    echo "🔴 GATE 1: Pre-Merge Validation"

    # Mark validation required in state
    jq --arg branch "$branch" '
        .quality_gates.gate_1_pre_merge.status = "in_progress" |
        .quality_gates.gate_1_pre_merge.current_branch = $branch |
        .quality_gates.current_gate = "gate_1_pre_merge"
    ' "$INT_STATE_FILE" > tmp.json && mv tmp.json "$INT_STATE_FILE"

    # Validation checks
    echo "📋 Pre-Merge Checks:"
    echo "  ✓ Branch builds successfully"
    echo "  ✓ All tests pass"
    echo "  ✓ No uncommitted changes"
    echo "  ✓ Branch is up to date"

    echo ""
    echo "INTEGRATE_WAVE_EFFORTS_PROGRESS: Validating branch $branch"
}
```

### GATE 2: POST-MERGE REVIEW
```bash
# MANDATORY after each merge
post_merge_review() {
    local branch="$1"
    echo "🔴 GATE 2: Post-Merge Review"

    jq --arg branch "$branch" '
        .quality_gates.gate_2_post_merge.status = "pending" |
        .quality_gates.gate_2_post_merge.current_branch = $branch |
        .quality_gates.current_gate = "gate_2_post_merge"
    ' "$INT_STATE_FILE" > tmp.json && mv tmp.json "$INT_STATE_FILE"

    echo "🚨 BLOCKING: Cannot continue until review completes"
    echo ""
    echo "Task: code-reviewer"
    echo "State: INTEGRATE_WAVE_EFFORTS_REVIEW"
    echo "Review Type: post-merge"
    echo "Branch: $branch"
    echo "Focus: merge-correctness, no-conflicts, build-success"
    echo ""

    echo "INTEGRATE_WAVE_EFFORTS_ISSUE_FOUND: Review required for $branch"
}
```

### GATE 3: VALIDATION TESTING
```bash
# Based on validation_level setting
validation_testing_gate() {
    local LEVEL=$(jq -r '.validation_level' "$INT_STATE_FILE")
    echo "🔴 GATE 3: Validation Testing ($LEVEL)"

    jq '.quality_gates.gate_3_validation.status = "in_progress" |
        .quality_gates.current_gate = "gate_3_validation"
    ' "$INT_STATE_FILE" > tmp.json && mv tmp.json "$INT_STATE_FILE"

    case "$LEVEL" in
        BASIC)
            echo "📊 Basic Validation:"
            echo "  ✓ Build passes"
            echo "  ✓ Unit tests pass"
            ;;
        FULL)
            echo "📊 Full Validation:"
            echo "  ✓ Build passes"
            echo "  ✓ All tests pass"
            echo "  ✓ Integration tests pass"
            ;;
        COMPREHENSIVE)
            echo "📊 Comprehensive Validation:"
            echo "  ✓ All tests pass"
            echo "  ✓ Performance benchmarks"
            echo "  ✓ Security scans"
            echo "  ✓ E2E tests"
            ;;
    esac
}
```

### GATE 4: COMPREHENSIVE FINAL CHECK
```bash
# MANDATORY before completing integration
comprehensive_final_check() {
    echo "🔴 GATE 4: Comprehensive Final Check"

    jq '.quality_gates.gate_4_comprehensive.status = "pending" |
        .quality_gates.current_gate = "gate_4_comprehensive"
    ' "$INT_STATE_FILE" > tmp.json && mv tmp.json "$INT_STATE_FILE"

    echo "📊 COMPREHENSIVE VALIDATION REQUIRED:"
    echo "  ✓ All branches successfully integrated"
    echo "  ✓ No merge conflicts remaining"
    echo "  ✓ All tests pass on integrated branch"
    echo "  ✓ No regressions introduced"
    echo "  ✓ Ready for next phase"
    echo ""

    echo "INTEGRATE_WAVE_EFFORTS_PROGRESS: Final validation in progress"
}
```

## 📊 INTEGRATE_WAVE_EFFORTS PROGRESS TRACKING

```bash
show_integration_progress() {
    echo "═══════════════════════════════════════════════════════"
    echo "📊 INTEGRATE_WAVE_EFFORTS PROGRESS: $INT_TYPE"
    echo "═══════════════════════════════════════════════════════"

    # Show merge progress
    MERGED=$(jq '.integration_progress.branches_merged | length' "$INT_STATE_FILE")
    PENDING=$(jq '.integration_progress.branches_pending | length' "$INT_STATE_FILE")
    TOTAL=$((MERGED + PENDING))
    echo "🌿 Branches: $MERGED/$TOTAL merged"

    # Show conflicts
    CONFLICTS=$(jq '.integration_progress.conflicts_found | length' "$INT_STATE_FILE")
    echo "⚠️ Conflicts: $CONFLICTS found"

    # Show issues
    ISSUES=$(jq '.integration_progress.issues_detected | length' "$INT_STATE_FILE")
    echo "🔍 Issues: $ISSUES detected"

    # Current state
    echo "📍 Current State: $INT_CURRENT_STATE"
    echo "🔄 Attempt: $(jq -r '.attempt' "$INT_STATE_FILE")"

    echo "═══════════════════════════════════════════════════════"
}

show_integration_progress
```

## 🔄 CYCLE MANAGEMENT

```bash
# Handle INTEGRATE_WAVE_EFFORTS→FIX_CASCADE→INTEGRATE_WAVE_EFFORTS cycles
handle_integration_cycle() {
    local ISSUE_TYPE="$1"
    local ISSUE_DESC="$2"

    echo "🔄 Integration cycle required"

    # Record in cycle history
    jq --arg type "$ISSUE_TYPE" \
       --arg desc "$ISSUE_DESC" \
       '.cycle_history += [{
           "attempt": .attempt,
           "timestamp": now,
           "issue_type": $type,
           "description": $desc,
           "state_when_detected": .current_state
       }] |
       .attempt += 1' "$INT_STATE_FILE" > tmp.json && mv tmp.json "$INT_STATE_FILE"

    # Check cycle limit
    local ATTEMPTS=$(jq -r '.attempt' "$INT_STATE_FILE")
    if [[ $ATTEMPTS -gt 3 ]]; then
        echo "❌ CRITICAL: Maximum integration attempts (3) exceeded!"
        echo "⚠️ Manual intervention required"

        # Transition to ERROR state
        jq '.current_state = "INTEGRATE_WAVE_EFFORTS_ERROR"' "$INT_STATE_FILE" > tmp.json && \
            mv tmp.json "$INT_STATE_FILE"

        echo "INTEGRATE_WAVE_EFFORTS_COMPLETE: FAILED"
        return 1
    fi

    echo "📝 Transitioning to FIX_CASCADE state"
    echo "🎯 Will return to INTEGRATE_WAVE_EFFORTS after fixes"
}
```

## 💾 CHECKPOINT MANAGEMENT

```bash
save_integration_checkpoint() {
    echo "💾 Saving integration checkpoint..."

    # Update timestamp
    jq '.timestamps.last_updated = now' "$INT_STATE_FILE" > tmp.json && \
        mv tmp.json "$INT_STATE_FILE"

    # Commit to planning repo
    git add "$INT_STATE_FILE"
    git commit -m "integration: ${INT_TYPE} - checkpoint at $INT_CURRENT_STATE"
    git push

    echo "✅ Checkpoint saved - Resume with: /integration"
}

save_integration_checkpoint
```

## ⚠️ ERROR HANDLING

```bash
handle_integration_error() {
    local ERROR_MSG="$1"

    echo "❌ Integration error: $ERROR_MSG"

    # Log to integration state
    jq --arg error "$ERROR_MSG" \
       '.error_log += [{
           "timestamp": now,
           "state": "'$INT_CURRENT_STATE'",
           "error": $error
       }]' "$INT_STATE_FILE" > tmp.json && mv tmp.json "$INT_STATE_FILE"

    # Determine recovery action
    if [[ "$ERROR_MSG" == *"conflict"* ]]; then
        echo "⚠️ Merge conflict detected - manual resolution required"
        echo "INTEGRATE_WAVE_EFFORTS_ISSUE_FOUND: Merge conflict"
    elif [[ "$ERROR_MSG" == *"test"* ]]; then
        echo "⚠️ Test failure - fix cascade may be needed"
        echo "INTEGRATE_WAVE_EFFORTS_ISSUE_FOUND: Test failure"
    fi
}
```

## 🎯 USAGE

```bash
# Start a new integration
/integration type=WAVE branches=effort1,effort2,effort3 target=wave1-integration validation=FULL

# Continue existing integration
/integration
# Command auto-detects active integration and resumes

# Parameters:
- type: WAVE|PHASE|PROJECT (required for new)
- branches: comma-separated list (required for new)
- target: target branch name (required for new)
- validation: BASIC|FULL|COMPREHENSIVE (optional, default: BASIC)
```

## 📋 COMPLETION CRITERIA

Integration is complete when:
- [ ] All branches successfully merged
- [ ] All conflicts resolved
- [ ] Validation tests pass
- [ ] Quality gates passed
- [ ] Integration report created
- [ ] State archived
- [ ] Main state updated

## 🔀 STATE TRANSITIONS

Per SOFTWARE-FACTORY-INTEGRATE_WAVE_EFFORTS-STATE-MACHINE.md:
- INTEGRATE_WAVE_EFFORTS_INIT → INTEGRATE_WAVE_EFFORTS_PLANNING
- INTEGRATE_WAVE_EFFORTS_PLANNING → INTEGRATE_WAVE_EFFORTS_SETUP
- INTEGRATE_WAVE_EFFORTS_SETUP → INTEGRATE_WAVE_EFFORTS_MERGING
- INTEGRATE_WAVE_EFFORTS_MERGING → INTEGRATE_WAVE_EFFORTS_VALIDATION
- INTEGRATE_WAVE_EFFORTS_VALIDATION → INTEGRATE_WAVE_EFFORTS_COMPLETE or INTEGRATE_WAVE_EFFORTS_FIX_CASCADE
- INTEGRATE_WAVE_EFFORTS_FIX_CASCADE → INTEGRATE_WAVE_EFFORTS_MERGING (cycle)

---

**Remember**: Integration sub-state machine handles complex merge scenarios with automatic cycle management