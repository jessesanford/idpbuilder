---
name: cascade-orchestration
description: Execute full project cascade after Phase 1 Wave 1 rebuild through ERROR_RECOVERY protocol
---

# /cascade-orchestration

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║                   CASCADE ORCHESTRATION COMMAND                               ║
║                                                                               ║
║ Protocol: ERROR_RECOVERY CASCADE EXECUTION for Full Project Rebase           ║
║ Rules: R351 + R354 + R327 + R321 + STATE-MACHINE + GRADING                  ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🎯 AGENT IDENTITY ASSIGNMENT

**You are @agent-orchestrator in ERROR_RECOVERY cascade mode**

By invoking this command, you are the orchestrator executing a FULL PROJECT CASCADE through ERROR_RECOVERY state iterations. You must:
- Follow all orchestrator rules and protocols
- Execute cascade operations sequentially through ERROR_RECOVERY cycles
- Enforce R354 post-rebase reviews for EVERY rebase
- Track all operations in orchestrator-state.json
- **🔴 R322: MANDATORY STOP BEFORE STATE TRANSITIONS**
- **🔴 R354: EVERY REBASE REQUIRES CODE REVIEW**
- **🔴 R351: CASCADE EXECUTION IN EXACT ORDER**

## 🔄 CASCADE RESUMPTION CAPABILITY

This command is FULLY RESUMABLE and tracks meticulous state in orchestrator-state.json:
- Each operation tracks: status, started_at, completed_at, processed_efforts, version
- Possible statuses: pending, in_progress, rebased_pending_review, reviewing, fixing, completed
- On restart, automatically detects in-progress operations and resumes from correct point
- Tracks individual effort processing within each operation for granular resumption
- **Handles fix cascades**: Dynamically injects re-rebase operations when fixes discovered
- **Multi-cascade aware**: Tracks main cascade + multiple fix cascades without confusion

## 🚨 CASCADE CONTEXT - CRITICAL UNDERSTANDING 🚨

### Why This CASCADE is Required:
```
Phase 1 Wave 1 was REBUILT from scratch to fix critical issues.
ALL subsequent phases/waves were built on the OLD P1W1.
This cascade rebases EVERYTHING onto the new clean P1W1 foundation.
```

### CASCADE EXECUTION PATTERN:
```
For each operation in cascade_coordination.detailed_operations:
1. ERROR_RECOVERY → Read operation details
2. Execute rebase/recreation based on operation type
3. If rebase: SPAWN_CODE_REVIEWERS_FOR_REVIEW (R354)
4. MONITOR_REVIEWS → Track review completion
5. If fixes needed: SPAWN_ENGINEERS_FOR_FIXES → MONITOR_FIXES
6. Update cascade tracking
7. Return to ERROR_RECOVERY for next operation
```

## 🚨 MANDATORY PRE-FLIGHT CHECKS 🚨

### 1. Agent Identity Verification
```bash
WHO_AM_I="orchestrator"
CASCADE_MODE="FULL_PROJECT_CASCADE"
echo "✓ Confirming identity: $WHO_AM_I in $CASCADE_MODE mode"
```

### 2. CASCADE Rule Acknowledgment (MANDATORY)
```bash
echo "================================"
echo "CASCADE ORCHESTRATION RULE ACKNOWLEDGMENT"
echo "I am orchestrator in ERROR_RECOVERY cascade mode"
echo "I acknowledge these CASCADE-CRITICAL rules:"
echo "--------------------------------"
echo "🔴 R351: CASCADE EXECUTION PROTOCOL - Execute in exact"
echo "   dependency order with zero tolerance for shortcuts"
echo "--------------------------------"
echo "🔴 R354: POST-REBASE REVIEW REQUIREMENT - EVERY rebase"
echo "   MUST be followed by Code Reviewer validation"
echo "--------------------------------"
echo "🔴 R327: MANDATORY RE-INTEGRATION AFTER FIXES - All"
echo "   integrations must be recreated after source changes"
echo "--------------------------------"
echo "🔴 R321: IMMEDIATE BACKPORT DURING INTEGRATION - Fixes"
echo "   go to source branches, never to integrations"
echo "--------------------------------"
echo "🚨 R151: PARALLEL SPAWNING - When spawning multiple"
echo "   reviewers, use <5 second delta in SINGLE message"
echo "--------------------------------"
echo "🔴 R322: MANDATORY STOP BEFORE STATE TRANSITIONS"
echo "🔴 R232: TODOWRITE PENDING ITEMS ARE COMMANDS"
echo "🔴 R290: STATE RULE READING AND VERIFICATION"
echo "================================"
```

### 3. Environment and State Verification
```bash
# Verify we're in ERROR_RECOVERY state
CURRENT_STATE=$(jq -r '.current_state' orchestrator-state.json)
if [ "$CURRENT_STATE" != "ERROR_RECOVERY" ]; then
    echo "❌ ERROR: Not in ERROR_RECOVERY state (current: $CURRENT_STATE)"
    echo "Transition to ERROR_RECOVERY first!"
    exit 1
fi

# Verify cascade plan exists
CASCADE_OPS=$(jq '.cascade_coordination.detailed_operations | length' orchestrator-state.json)
if [ "$CASCADE_OPS" -eq 0 ]; then
    echo "❌ ERROR: No cascade operations found in state file!"
    exit 1
fi

echo "✅ State verified: ERROR_RECOVERY"
echo "✅ Cascade operations found: $CASCADE_OPS operations"

# Check for previous checkpoint (resumption case)
LAST_CHECKPOINT=$(jq -r '.cascade_coordination.last_checkpoint // empty' orchestrator-state.json)
if [ -n "$LAST_CHECKPOINT" ]; then
    echo ""
    echo "📍 CASCADE RESUMPTION DETECTED"
    echo "================================"
    CHECKPOINT_TIME=$(echo "$LAST_CHECKPOINT" | jq -r '.timestamp')
    CHECKPOINT_OP=$(echo "$LAST_CHECKPOINT" | jq -r '.operation')
    CHECKPOINT_STATUS=$(echo "$LAST_CHECKPOINT" | jq -r '.status')

    echo "Last checkpoint: $CHECKPOINT_TIME"
    echo "Last operation: #$CHECKPOINT_OP"
    echo "Last status: $CHECKPOINT_STATUS"

    # Show progress
    COMPLETED=$(jq '[.cascade_coordination.detailed_operations[] | select(.status == "completed")] | length' orchestrator-state.json)
    echo "Progress: $COMPLETED/$CASCADE_OPS operations completed"
    echo "================================"
fi
```

### 4. CASCADE PROGRESS CHECK
```bash
# Determine where we are in the cascade
echo "🔍 Checking cascade progress..."

# Check for pending fix cascades FIRST
PENDING_FIXES=$(jq -r '.cascade_coordination.re_rebase_queue | length // 0' orchestrator-state.json)
if [ "$PENDING_FIXES" -gt 0 ]; then
    echo "⚠️ FIX CASCADE DETECTED - $PENDING_FIXES re-rebase operations queued"

    # Get highest priority re-rebase
    PRIORITY_FIX=$(jq -r '.cascade_coordination.re_rebase_queue |
        sort_by(.priority) | .[0]' orchestrator-state.json)

    echo "Processing priority fix: $(echo "$PRIORITY_FIX" | jq -r '.branch')"
    echo "Reason: $(echo "$PRIORITY_FIX" | jq -r '.reason')"

    # This becomes our next operation
    INJECT_FIX_CASCADE=true
fi

# Check for in-progress operations (resumption case)
IN_PROGRESS_OP=$(jq -r '.cascade_coordination.detailed_operations[] |
    select(.status == "in_progress" or .status == "rebased_pending_review" or
           .status == "reviewing" or .status == "fixing") |
    .order' orchestrator-state.json | head -1)

if [ -n "$IN_PROGRESS_OP" ]; then
    echo "⚠️ RESUMING CASCADE - Found in-progress operation: #$IN_PROGRESS_OP"

    # Get operation details and status
    OP_STATUS=$(jq -r --arg order "$IN_PROGRESS_OP" '
        .cascade_coordination.detailed_operations[] |
        select(.order == ($order | tonumber)) | .status' orchestrator-state.json)

    echo "Current status: $OP_STATUS"

    # Determine where to resume based on status
    case "$OP_STATUS" in
        "in_progress")
            echo "Resuming rebase operation..."
            NEXT_OP="$IN_PROGRESS_OP"
            ;;
        "rebased_pending_review")
            echo "Resuming at review stage..."
            NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_REVIEW"
            ;;
        "reviewing")
            echo "Resuming review monitoring..."
            NEXT_STATE="MONITOR_REVIEWS"
            ;;
        "fixing")
            echo "Resuming fix monitoring..."
            NEXT_STATE="MONITORING_FIX_PROGRESS"
            ;;
    esac
else
    # Find next operation to execute
    NEXT_OP=$(jq -r '.cascade_coordination.detailed_operations[] |
        select(.status != "completed" and .status != "skipped") |
        .order' orchestrator-state.json | head -1)

    if [ -z "$NEXT_OP" ]; then
        echo "🎉 CASCADE COMPLETE! All operations finished."
        echo "Decision: TRANSITION_TO_WAVE_COMPLETE"
    else
        echo "📊 Next operation: #$NEXT_OP"
    fi
fi

# Get operation details
if [ -n "$NEXT_OP" ]; then
    OP_DETAILS=$(jq --arg order "$NEXT_OP" '
        .cascade_coordination.detailed_operations[] |
        select(.order == ($order | tonumber))
    ' orchestrator-state.json)

    echo "Operation type: $(echo "$OP_DETAILS" | jq -r '.operation')"
    echo "Phase: $(echo "$OP_DETAILS" | jq -r '.phase // "N/A"')"
    echo "Wave: $(echo "$OP_DETAILS" | jq -r '.wave // "N/A"')"
    echo "Decision: EXECUTE_CASCADE_OPERATION"
fi
```

## 🔄 CASCADE OPERATION EXECUTION

### For REBASE Operations:
```bash
if [[ "$(echo "$OP_DETAILS" | jq -r '.operation')" == "rebase_efforts" ]]; then
    echo "═══════════════════════════════════════════════════════"
    echo "EXECUTING REBASE OPERATION #$NEXT_OP"
    echo "═══════════════════════════════════════════════════════"

    PHASE=$(echo "$OP_DETAILS" | jq -r '.phase')
    WAVE=$(echo "$OP_DETAILS" | jq -r '.wave')
    EFFORTS=$(echo "$OP_DETAILS" | jq -r '.efforts[]')
    REBASE_ONTO=$(echo "$OP_DETAILS" | jq -r '.rebase_onto')

    echo "📋 Rebasing Phase $PHASE Wave $WAVE efforts onto $REBASE_ONTO"

    # Mark operation as in_progress
    jq --arg order "$NEXT_OP" '
        .cascade_coordination.detailed_operations |=
        map(if .order == ($order | tonumber)
            then . + {"status": "in_progress", "started_at": now | todate}
            else . end)
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Track which efforts have been processed
    PROCESSED_EFFORTS=$(jq -r --arg order "$NEXT_OP" '
        .cascade_coordination.detailed_operations[] |
        select(.order == ($order | tonumber)) |
        .processed_efforts // []' orchestrator-state.json)

    # For each effort, execute rebase
    for effort in $EFFORTS; do
        # Skip if already processed (for resumption)
        if echo "$PROCESSED_EFFORTS" | grep -q "$effort"; then
            echo "✅ Already processed: $effort"
            continue
        fi

        echo "🔄 Rebasing $effort..."

        # Acquire operation lock for this branch
        LOCK_ACQUIRED=$(jq --arg branch "$effort" '
            if .cascade_coordination.operation_locks[$branch] then
                false
            else
                .cascade_coordination.operation_locks[$branch] = {
                    "cascade_id": "main",
                    "operation": "'$NEXT_OP'",
                    "locked_at": now | todate
                } | true
            end' orchestrator-state.json)

        if [ "$LOCK_ACQUIRED" == "false" ]; then
            echo "⚠️ Branch locked by another operation, skipping..."
            continue
        fi

        # Clone effort workspace if not exists
        EFFORT_DIR="efforts/phase$PHASE/wave$WAVE/$effort"
        if [ ! -d "$EFFORT_DIR" ]; then
            echo "Creating effort workspace: $EFFORT_DIR"
            # Clone and setup workspace
            setup_effort_workspace "$PHASE" "$WAVE" "$effort"
        fi

        # Execute rebase
        cd "$EFFORT_DIR"
        git fetch origin
        git checkout "idpbuilder-oci-build-push/phase$PHASE/wave$WAVE/$effort"

        # Rebase onto new base
        echo "Rebasing onto $REBASE_ONTO..."
        if git rebase "origin/$REBASE_ONTO"; then
            echo "✅ Rebase successful"
            git push --force-with-lease origin

            # Mark effort as processed
            jq --arg order "$NEXT_OP" --arg effort "$effort" '
                .cascade_coordination.detailed_operations |=
                map(if .order == ($order | tonumber)
                    then . + {"processed_efforts": ((.processed_efforts // []) + [$effort])}
                    else . end)' $CLAUDE_PROJECT_DIR/orchestrator-state.json > tmp.json && \
                mv tmp.json $CLAUDE_PROJECT_DIR/orchestrator-state.json

            # Mark for R354 review
            jq --arg effort "$effort" '
                .cascade_coordination.pending_reviews += [{
                    "effort": $effort,
                    "type": "post_rebase",
                    "review_required": true,
                    "r354_mandated": true
                }]' $CLAUDE_PROJECT_DIR/orchestrator-state.json > tmp.json && \
                mv tmp.json $CLAUDE_PROJECT_DIR/orchestrator-state.json
        else
            echo "❌ Rebase conflict! Need manual resolution"
            # Handle conflict...
        fi

        cd $CLAUDE_PROJECT_DIR
    done

    # Update operation status
    jq --arg order "$NEXT_OP" '
        .cascade_coordination.detailed_operations |=
        map(if .order == ($order | tonumber)
            then . + {"status": "rebased_pending_review"}
            else . end)
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # 🔴 R354: MANDATORY POST-REBASE REVIEWS
    echo "🔴 R354: Spawning Code Reviewers for post-rebase validation..."

    # Transition to spawn reviewers
    NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_REVIEW"
fi
```

### For INTEGRATION RECREATION Operations:
```bash
if [[ "$(echo "$OP_DETAILS" | jq -r '.operation')" =~ recreate ]]; then
    echo "═══════════════════════════════════════════════════════"
    echo "EXECUTING INTEGRATION RECREATION #$NEXT_OP"
    echo "═══════════════════════════════════════════════════════"

    INTEGRATION_BRANCH=$(echo "$OP_DETAILS" | jq -r '.integration_branch')
    BASE_BRANCH=$(echo "$OP_DETAILS" | jq -r '.base')

    echo "🗑️ Deleting old integration: $INTEGRATION_BRANCH"
    git push origin --delete "$INTEGRATION_BRANCH" 2>/dev/null || true

    echo "📋 Creating integration plan..."

    # Transition to integration states
    if [[ "$INTEGRATION_BRANCH" =~ wave.*integration ]]; then
        NEXT_STATE="SETUP_INTEGRATION_INFRASTRUCTURE"
    elif [[ "$INTEGRATION_BRANCH" =~ phase.*integration ]]; then
        NEXT_STATE="SETUP_PHASE_INTEGRATION_INFRASTRUCTURE"
    elif [[ "$INTEGRATION_BRANCH" =~ project.*integration ]]; then
        NEXT_STATE="SETUP_PROJECT_INTEGRATION_INFRASTRUCTURE"
    fi

    # Mark operation for tracking
    jq --arg order "$NEXT_OP" --arg branch "$INTEGRATION_BRANCH" '
        .cascade_coordination.current_integration = {
            "branch": $branch,
            "operation_order": $order,
            "status": "pending_recreation"
        }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
fi
```

## 🔴 R354 ENFORCEMENT - POST-REBASE REVIEWS

### Check for Pending Reviews:
```bash
PENDING_REVIEWS=$(jq '.cascade_coordination.pending_reviews | length' orchestrator-state.json)

if [ "$PENDING_REVIEWS" -gt 0 ]; then
    echo "🔴 R354: $PENDING_REVIEWS post-rebase reviews required!"

    # Must spawn reviewers
    echo "Spawning Code Reviewers with cascade_mode=true..."

    # For each pending review
    EFFORTS_TO_REVIEW=$(jq -r '.cascade_coordination.pending_reviews[].effort' orchestrator-state.json)

    # Spawn reviewers (R151 compliant - parallel)
    for effort in $EFFORTS_TO_REVIEW; do
        Task: code-reviewer
        Working directory: [effort directory]
        State: POST_REBASE_REVIEW
        Instructions:
        - This is a CASCADE MODE post-rebase review (R354)
        - Focus ONLY on integration validation
        - Check: builds pass, tests pass, no conflicts
        - SKIP: size checks, style checks, documentation
        - Return: REBASE_VALID or FIXES_NEEDED
    done

    NEXT_STATE="MONITOR_REVIEWS"
fi
```

## 📊 CASCADE TRACKING AND STATE UPDATES

### After Each Operation:
```bash
# Update cascade progress
update_cascade_progress() {
    local OP_ORDER=$1
    local STATUS=$2

    # Mark operation complete
    jq --arg order "$OP_ORDER" --arg status "$STATUS" '
        .cascade_coordination.detailed_operations |=
        map(if .order == ($order | tonumber)
            then . + {"status": $status, "completed_at": now | todate}
            else . end)
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Calculate progress
    TOTAL_OPS=$(jq '.cascade_coordination.detailed_operations | length' orchestrator-state.json)
    COMPLETED=$(jq '[.cascade_coordination.detailed_operations[] |
                     select(.status == "completed")] | length' orchestrator-state.json)

    echo "📊 CASCADE PROGRESS: $COMPLETED/$TOTAL_OPS operations complete"

    # Check if cascade complete
    if [ "$COMPLETED" -eq "$TOTAL_OPS" ]; then
        echo "🎉 CASCADE COMPLETE!"
        jq '.cascade_coordination.status = "completed"' orchestrator-state.json > tmp.json && \
            mv tmp.json orchestrator-state.json

        # Transition out of ERROR_RECOVERY
        NEXT_STATE="WAVE_COMPLETE"
    else
        # Continue with next operation
        echo "Returning to ERROR_RECOVERY for next operation..."
        NEXT_STATE="ERROR_RECOVERY"
    fi
}
```

## 🔄 STATE TRANSITION PATTERNS

### CASCADE State Flow:
```
ERROR_RECOVERY (read next op) →
  [If rebase]: Execute rebase → SPAWN_CODE_REVIEWERS_FOR_REVIEW →
               MONITOR_REVIEWS → [If fixes]: SPAWN_ENGINEERS_FOR_FIXES →
               MONITOR_FIXES → ERROR_RECOVERY (next op)

  [If recreate]: Mark for recreation → SETUP_*_INTEGRATION_INFRASTRUCTURE →
                 SPAWN_INTEGRATION_AGENT → MONITORING_INTEGRATION →
                 ERROR_RECOVERY (next op)

  [If complete]: WAVE_COMPLETE → Continue normal flow
```

## 🛑 CASCADE COMPLETION CRITERIA

The cascade is complete when:
- [ ] All rebase operations executed
- [ ] All post-rebase reviews passed (R354)
- [ ] All integration branches recreated
- [ ] cascade_coordination.status = "completed"
- [ ] No pending operations remain

## 📝 CHECKPOINT PROTOCOL

### Before EVERY State Transition:
```bash
# Save cascade progress
save_cascade_checkpoint() {
    local OP_ORDER=$1
    local STATUS=$2
    local DETAIL=$3

    echo "💾 Saving cascade checkpoint..."

    # Update cascade session tracking
    jq --arg session "$(date '+%Y%m%d-%H%M%S')" '
        .cascade_coordination.last_checkpoint = {
            "timestamp": now | todate,
            "session_id": $session,
            "operation": '$OP_ORDER',
            "status": "'$STATUS'",
            "detail": "'$DETAIL'"
        }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Save TODOs per R287
    TODO_FILE="todos/cascade-$(date '+%Y%m%d-%H%M%S').todo"
    echo "[CASCADE Operation #$OP_ORDER]" > "$TODO_FILE"
    echo "Status: $STATUS" >> "$TODO_FILE"
    echo "Detail: $DETAIL" >> "$TODO_FILE"
    echo "Next State: $NEXT_STATE" >> "$TODO_FILE"

    # Commit state
    git add orchestrator-state.json todos/*.todo
    git commit -m "cascade: Op #$OP_ORDER - $STATUS - $DETAIL

CASCADE CHECKPOINT:
- Operation: #$OP_ORDER
- Status: $STATUS
- Progress: $(jq '[.cascade_coordination.detailed_operations[] | select(.status == "completed")] | length' orchestrator-state.json)/$(jq '.cascade_coordination.detailed_operations | length' orchestrator-state.json)
- Next State: $NEXT_STATE
- Can resume with: /cascade-orchestration"

    git push

    echo "✅ Checkpoint saved - Can resume anytime with /cascade-orchestration"
}

# R322: STOP before transition
echo "🛑 STOPPING - State transition required"
echo "Next state: $NEXT_STATE"
echo "Use /cascade-orchestration to continue"
save_cascade_checkpoint "$NEXT_OP" "$STATUS" "State transition to $NEXT_STATE"
```

## ⚠️ ERROR HANDLING

### Rebase Conflicts:
```bash
if [[ "$REBASE_RESULT" == "CONFLICT" ]]; then
    echo "❌ Rebase conflict detected!"

    # Create fix plan
    jq '.cascade_coordination.blocked_operations += [{
        "operation_order": '$NEXT_OP',
        "reason": "rebase_conflict",
        "effort": "'$effort'",
        "action_required": "manual_resolution"
    }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Spawn SW Engineer to resolve
    NEXT_STATE="SPAWN_ENGINEERS_FOR_FIXES"
fi
```

### Review Failures:
```bash
if [[ "$REVIEW_RESULT" == "FIXES_NEEDED" ]]; then
    echo "❌ Post-rebase review found issues!"

    # Record fix cascade requirement
    FIXED_BRANCH=$(jq -r '.cascade_coordination.current_review.branch' orchestrator-state.json)
    PHASE=$(echo "$FIXED_BRANCH" | sed 's/.*phase\([0-9]\).*/\1/')
    WAVE=$(echo "$FIXED_BRANCH" | sed 's/.*wave\([0-9]\).*/\1/')

    # Calculate downstream impact
    echo "📊 Calculating cascade impact of fix..."

    # Find all operations that depend on this branch
    DOWNSTREAM_OPS=$(jq --arg phase "$PHASE" --arg wave "$WAVE" '
        .cascade_coordination.detailed_operations[] |
        select((.phase > ($phase | tonumber)) or
               (.phase == ($phase | tonumber) and .wave > ($wave | tonumber))) |
        select(.status != "completed")' orchestrator-state.json)

    if [ -n "$DOWNSTREAM_OPS" ]; then
        echo "⚠️ Fix will require re-cascading downstream operations"

        # Add to re_rebase_queue
        jq --arg branch "$FIXED_BRANCH" '
            .cascade_coordination.re_rebase_queue +=
            (.cascade_coordination.detailed_operations[] |
             select(.status != "completed") |
             select(.phase > '$PHASE' or (.phase == '$PHASE' and .wave > '$WAVE')) |
             .efforts[] | {
                "branch": .,
                "original_base": .rebase_onto,
                "new_base": ($branch + "-fixed"),
                "reason": "upstream fix in " + $branch,
                "priority": 2,
                "queued_at": now | todate
             })' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

        # Mark downstream operations for re-execution
        jq --arg phase "$PHASE" --arg wave "$WAVE" '
            .cascade_coordination.detailed_operations |=
            map(if ((.phase > ($phase | tonumber)) or
                   (.phase == ($phase | tonumber) and .wave > ($wave | tonumber)))
                then . + {"needs_re_rebase": true, "version": (.version // 1) + 1}
                else . end)' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    fi

    # Issues become part of cascade (R354)
    echo "Adding fixes to cascade chain..."

    # Spawn engineers for fixes
    NEXT_STATE="SPAWN_ENGINEERS_FOR_FIXES"
fi
```

## 📊 CASCADE STATUS VISUALIZATION

```bash
# Display comprehensive cascade status
show_cascade_status() {
    echo "═══════════════════════════════════════════════════════"
    echo "📊 CASCADE STATUS OVERVIEW"
    echo "═══════════════════════════════════════════════════════"

    # Main cascade progress
    TOTAL_OPS=$(jq '.cascade_coordination.detailed_operations | length' orchestrator-state.json)
    COMPLETED=$(jq '[.cascade_coordination.detailed_operations[] | select(.status == "completed")] | length' orchestrator-state.json)
    IN_PROGRESS=$(jq '[.cascade_coordination.detailed_operations[] | select(.status == "in_progress")] | length' orchestrator-state.json)

    echo "🏗️ Main Cascade:"
    echo "  Progress: $COMPLETED/$TOTAL_OPS operations completed"
    echo "  In Progress: $IN_PROGRESS operations"

    # Fix cascades
    FIX_CASCADES=$(jq '.cascade_coordination.fix_cascades | length // 0' orchestrator-state.json)
    if [ "$FIX_CASCADES" -gt 0 ]; then
        echo ""
        echo "🔧 Active Fix Cascades: $FIX_CASCADES"
        jq -r '.cascade_coordination.fix_cascades[] |
            "  • \(.fix_id)\n    Branch: \(.trigger_branch)\n    Status: \(.status)"' \
            orchestrator-state.json 2>/dev/null
    fi

    # Re-rebase queue
    REBASE_QUEUE=$(jq '.cascade_coordination.re_rebase_queue | length // 0' orchestrator-state.json)
    if [ "$REBASE_QUEUE" -gt 0 ]; then
        echo ""
        echo "📋 Re-rebase Queue: $REBASE_QUEUE operations pending"
        jq -r '.cascade_coordination.re_rebase_queue |
            sort_by(.priority) | .[:3] |
            .[] | "  • \(.branch) (Priority: \(.priority))\n    Reason: \(.reason)"' \
            orchestrator-state.json 2>/dev/null
    fi

    # Operations needing re-rebase
    NEEDS_REREBASE=$(jq '[.cascade_coordination.detailed_operations[] | select(.needs_re_rebase == true)] | length' orchestrator-state.json)
    if [ "$NEEDS_REREBASE" -gt 0 ]; then
        echo ""
        echo "⚠️ Operations Requiring Re-execution: $NEEDS_REREBASE"
    fi

    echo "═══════════════════════════════════════════════════════"
}

# Call status display at key points
show_cascade_status
```

## 🎯 USAGE INSTRUCTIONS

This command should be invoked repeatedly throughout the cascade:
```bash
/cascade-orchestration  # Start/continue cascade
# [Orchestrator executes current operation]
# [Orchestrator stops at state boundary]

/cascade-orchestration  # Continue to next operation
# [Repeat until cascade complete]
```

## 📊 GRADING COMPLIANCE

You will be graded on:
- ✅ R351: Exact cascade execution order
- ✅ R354: Post-rebase reviews for EVERY rebase
- ✅ R151: Parallel spawning where applicable
- ✅ R322: Stopping at EVERY state boundary
- ✅ State file tracking completeness
- ✅ No operations skipped or shortcuts taken

---

**Remember**:
- This is a MULTI-DAY operation with many state transitions
- EVERY rebase needs review (R354)
- Track EVERYTHING in orchestrator-state.json
- Stop at EVERY state boundary (R322)
- The cascade continues until ALL operations complete

**"Execute in order, review every rebase, track everything, stop at boundaries"**