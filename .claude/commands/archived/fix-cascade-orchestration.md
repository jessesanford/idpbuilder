---
name: fix-cascade-orchestration
description: [DEPRECATED - Use /fix-cascade] Legacy command for GiteaClient fix cascade
---

# /fix-cascade-orchestration [DEPRECATED]

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        ⚠️ DEPRECATED COMMAND ⚠️                              ║
║                                                                               ║
║ This command is specific to the GiteaClient fix and has been replaced        ║
║ by the generic /fix-cascade command that follows R375.                       ║
║                                                                               ║
║ Use: /fix-cascade                                                            ║
║                                                                               ║
║ The new command will:                                                        ║
║ - Auto-detect fix plans (*-FIX-PLAN.md)                                     ║
║ - Use R375 dual state tracking                                              ║
║ - Support multiple concurrent fixes                                          ║
║ - Archive completed fix states                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🔄 MIGRATION TO NEW COMMAND

If you have an existing GiteaClient fix in progress:
1. Your state file `orchestrator-gitea-fix-state.json` will work with new command
2. Rename to `orchestrator-gitea-client-state.json` for consistency
3. Run `/fix-cascade` to continue

---

## ORIGINAL GITEA-SPECIFIC COMMAND (PRESERVED FOR REFERENCE)

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║            GITEA CLIENT FIX CASCADE (LEGACY - SPECIFIC)                      ║
║                                                                               ║
║ Protocol: FIX CASCADE EXECUTION for Interface Alignment                      ║
║ Rules: R351 + R354 + R327 + R321 + STATE-MACHINE + FIX-TRACKING             ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🎯 AGENT IDENTITY ASSIGNMENT

**You are @agent-orchestrator in FIX_CASCADE mode**

By invoking this command, you are the orchestrator executing a FIX CASCADE for the GiteaClient interface alignment issue. You must:
- Follow all orchestrator rules and fix tracking protocols
- Execute fix operations sequentially through ERROR_RECOVERY cycles
- Track all operations in both orchestrator-state.json AND fix-specific state
- **🔴 R322: MANDATORY STOP BEFORE STATE TRANSITIONS**
- **🔴 R354: EVERY FIXED BRANCH REQUIRES CODE REVIEW**
- **🔴 R351: CASCADE EXECUTION IN EXACT ORDER**

## 🔄 FIX CASCADE RESUMPTION CAPABILITY

This command is FULLY RESUMABLE and tracks dual state:
1. **Main state**: orchestrator-state.json (standard SF state)
2. **Fix state**: orchestrator-gitea-fix-state.json (fix-specific tracking)

On restart:
- Reads both state files to determine position
- Resumes from exact operation/phase
- Handles nested fix cascades if new issues discovered

## 🚨 FIX CASCADE CONTEXT - CRITICAL UNDERSTANDING 🚨

### Why This FIX CASCADE is Required:
```
GiteaClient implements Client interface instead of Registry interface.
This causes CLI to call stub implementations instead of real push logic.
Fix must be:
1. Backported to 4 original effort branches
2. Forward-ported through 12 pr-ready branches
3. Validated with comprehensive testing
```

### FIX CASCADE EXECUTION PATTERN:
```
For each phase in fix tracking plan:
1. ERROR_RECOVERY → Read fix operation details
2. Execute fix based on phase (ANALYSIS/BACKPORT/FORWARD_PORT/VALIDATION)
3. If code changes: SPAWN_CODE_REVIEWERS_FOR_REVIEW (R354)
4. MONITOR_REVIEWS → Track review completion
5. If fixes needed: SPAWN_ENGINEERS_FOR_FIXES → MONITOR_FIXES
6. Update both state files
7. Return to ERROR_RECOVERY for next phase
```

## 🚨 MANDATORY PRE-FLIGHT CHECKS 🚨

### 1. Agent Identity Verification
```bash
WHO_AM_I="orchestrator"
FIX_MODE="GITEA_CLIENT_INTERFACE_FIX"
echo "✓ Confirming identity: $WHO_AM_I in $FIX_MODE mode"
```

### 2. FIX CASCADE Rule Acknowledgment (MANDATORY)
```bash
echo "================================"
echo "FIX CASCADE ORCHESTRATION RULE ACKNOWLEDGMENT"
echo "I am orchestrator in FIX_CASCADE mode"
echo "I acknowledge these FIX-CRITICAL rules:"
echo "--------------------------------"
echo "🔴 R351: CASCADE EXECUTION PROTOCOL - Execute fixes in"
echo "   exact dependency order: backport → forward-port → validate"
echo "--------------------------------"
echo "🔴 R354: POST-FIX REVIEW REQUIREMENT - EVERY fixed branch"
echo "   MUST be reviewed by Code Reviewer for validation"
echo "--------------------------------"
echo "🔴 FIX-001: BACKPORT ORDERING - Must fix in order:"
echo "   split-001 → split-002 → split-003 → main branch"
echo "--------------------------------"
echo "🔴 FIX-002: CHANGE VALIDATION - Every fix must compile,"
echo "   pass tests, maintain <800 lines, keep compatibility"
echo "--------------------------------"
echo "🔴 FIX-003: FORWARD PORT PROTOCOL - Create from fixed"
echo "   efforts, preserve changes, rebase dependents"
echo "--------------------------------"
echo "🔴 R322: MANDATORY STOP BEFORE STATE TRANSITIONS"
echo "🔴 R232: TODOWRITE PENDING ITEMS ARE COMMANDS"
echo "================================"
```

### 3. Fix State Verification
```bash
# Check main orchestrator state
MAIN_STATE=$(jq -r '.current_state' orchestrator-state.json)
echo "Main orchestrator state: $MAIN_STATE"

# Check fix-specific state
FIX_STATE_FILE="orchestrator-gitea-fix-state.json"
if [ ! -f "$FIX_STATE_FILE" ]; then
    echo "❌ ERROR: Fix state file not found!"
    echo "Looking for fix documentation..."

    if [ -f "ORCHESTRATOR-GITEA-FIX-PROMPT.md" ]; then
        echo "✅ Found fix prompt - initializing fix state"

        # Initialize fix state from template
        cat > $FIX_STATE_FILE << 'EOF'
{
  "fix_id": "gitea-client-interface-alignment",
  "current_state": "INIT",
  "current_phase": "phase1",
  "timestamps": {
    "started": null,
    "last_updated": null,
    "completed": null
  },
  "backport_status": {
    "gitea-client-split-001": "PENDING",
    "gitea-client-split-002": "PENDING",
    "gitea-client-split-003": "PENDING",
    "gitea-client": "PENDING"
  },
  "forward_port_status": {
    "branches_completed": [],
    "branches_pending": [
      "registry-types-pr-ready",
      "cert-validation-pr-ready",
      "credential-management-pr-ready",
      "kubernetes-helpers-pr-ready",
      "registry-helpers-pr-ready",
      "gitea-client-split-001-pr-ready",
      "gitea-client-split-002-pr-ready",
      "gitea-client-split-003-pr-ready",
      "cert-chain-validator-pr-ready",
      "image-loader-pr-ready",
      "progress-tracker-pr-ready",
      "test-utils-pr-ready"
    ],
    "current_branch": null,
    "conflicts_resolved": 0
  },
  "validation_status": {
    "merge_test": "NOT_STARTED",
    "build_test": "NOT_STARTED",
    "push_test": "NOT_STARTED"
  },
  "error_log": [],
  "recovery_points": {},
  "working_directories": {
    "backport": "/tmp/gitea-fix-backport",
    "forward_port": "/tmp/gitea-fix-forward",
    "validation": "/tmp/gitea-fix-validation"
  }
}
EOF
        echo "✅ Fix state initialized"
    else
        echo "❌ No fix documentation found. Cannot proceed."
        echo "Required files:"
        echo "  - ORCHESTRATOR-GITEA-FIX-PROMPT.md"
        echo "  - GITEACLIENT-INTERFACE-FIX-PLAN.md"
        echo "  - SOFTWARE-FACTORY-FIX-TRACKING-PLAN.md"
        exit 1
    fi
fi

# Read fix state
FIX_CURRENT_STATE=$(jq -r '.current_state' $FIX_STATE_FILE)
FIX_PHASE=$(jq -r '.current_phase' $FIX_STATE_FILE)

echo "✅ Fix state loaded:"
echo "   - Fix State: $FIX_CURRENT_STATE"
echo "   - Fix Phase: $FIX_PHASE"

# Verify we can proceed
if [ "$MAIN_STATE" != "ERROR_RECOVERY" ] && [ "$MAIN_STATE" != "WAVE_COMPLETE" ]; then
    echo "⚠️ Main orchestrator not in ERROR_RECOVERY or WAVE_COMPLETE state"
    echo "Transitioning to ERROR_RECOVERY for fix cascade..."

    # Update main state
    jq '.previous_state = .current_state |
        .current_state = "ERROR_RECOVERY" |
        .transition_time = now |
        .transition_reason = "Entering FIX CASCADE for GiteaClient interface alignment"' \
        orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
fi
```

### 4. Fix Documentation Check
```bash
# Ensure fix documentation is available
echo "🔍 Checking for fix documentation..."

REQUIRED_DOCS=(
    "GITEACLIENT-INTERFACE-FIX-PLAN.md"
    "SOFTWARE-FACTORY-FIX-TRACKING-PLAN.md"
)

for doc in "${REQUIRED_DOCS[@]}"; do
    if [ -f "$doc" ]; then
        echo "✅ Found: $doc"
    else
        echo "❌ Missing: $doc"
        echo "Cannot proceed without fix documentation"
        exit 1
    fi
done
```

## 🔄 FIX CASCADE OPERATION EXECUTION

### Determine Current Fix Phase:
```bash
case "$FIX_CURRENT_STATE" in
    "INIT")
        echo "═══════════════════════════════════════════════════════"
        echo "PHASE 1: ANALYSIS - Creating Fix Implementation"
        echo "═══════════════════════════════════════════════════════"

        # Read fix plan
        echo "📋 Reading fix plan from GITEACLIENT-INTERFACE-FIX-PLAN.md"

        # Spawn code reviewer to create exact fix
        echo "Spawning Code Reviewer to create fix implementation..."

        Task: code-reviewer
        Working directory: $CLAUDE_PROJECT_DIR
        State: FIX_ANALYSIS
        Instructions:
        Create the exact code changes needed to align GiteaClient with Registry interface.

        Read: GITEACLIENT-INTERFACE-FIX-PLAN.md

        The fix should:
        1. Change GiteaClient.Push signature to match Registry interface
        2. Add missing Registry methods (List, Exists, Delete)
        3. Handle io.Reader to v1.Image conversion
        4. Maintain backward compatibility where possible

        Output exact code changes to: /tmp/gitea-fix-implementation.patch

        # Update state
        jq '.current_state = "ANALYSIS" |
            .timestamps.started = now' $FIX_STATE_FILE > tmp.json && \
            mv tmp.json $FIX_STATE_FILE

        NEXT_STATE="BACKPORT"
        ;;

    "ANALYSIS")
        echo "Verifying fix implementation exists..."
        if [ -f "/tmp/gitea-fix-implementation.patch" ]; then
            echo "✅ Fix implementation ready"
            jq '.current_state = "BACKPORT"' $FIX_STATE_FILE > tmp.json && \
                mv tmp.json $FIX_STATE_FILE
            FIX_CURRENT_STATE="BACKPORT"
        else
            echo "⚠️ Waiting for fix implementation..."
        fi
        ;;

    "BACKPORT")
        echo "═══════════════════════════════════════════════════════"
        echo "PHASE 2: BACKPORT - Applying Fix to Original Efforts"
        echo "═══════════════════════════════════════════════════════"

        # Get branches needing backport
        PENDING_BACKPORTS=$(jq -r '.backport_status | to_entries[] |
            select(.value == "PENDING") | .key' $FIX_STATE_FILE | head -1)

        if [ -z "$PENDING_BACKPORTS" ]; then
            echo "✅ All backports complete"
            jq '.current_state = "FORWARD_PORT"' $FIX_STATE_FILE > tmp.json && \
                mv tmp.json $FIX_STATE_FILE
            NEXT_STATE="FORWARD_PORT"
        else
            echo "📋 Backporting to: $PENDING_BACKPORTS"

            # Mark as in progress
            jq --arg branch "$PENDING_BACKPORTS" \
                '.backport_status[$branch] = "IN_PROGRESS"' \
                $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

            # Spawn software engineer for backport
            echo "Spawning Software Engineer for backport..."

            WORK_DIR="/tmp/gitea-fix-backport/$PENDING_BACKPORTS"

            Task: software-engineer
            Working directory: $WORK_DIR
            State: FIX_BACKPORT
            Instructions:
            Apply the GiteaClient interface fix to branch: $PENDING_BACKPORTS

            Steps:
            1. Clone idpbuilder if not exists
            2. Checkout the ORIGINAL effort branch (not pr-ready)
            3. Apply fix from /tmp/gitea-fix-implementation.patch
            4. Ensure it compiles (go build)
            5. Commit: "fix(registry): align GiteaClient with Registry interface [GITEA-FIX]"
            6. Create branch: ${PENDING_BACKPORTS}-fixed
            7. Push the fixed branch

            Report success/failure and any issues.

            NEXT_STATE="MONITOR_BACKPORT"
        fi
        ;;

    "FORWARD_PORT")
        echo "═══════════════════════════════════════════════════════"
        echo "PHASE 3: FORWARD PORT - Rebasing PR-Ready Branches"
        echo "═══════════════════════════════════════════════════════"

        # Get next branch to forward port
        NEXT_FORWARD=$(jq -r '.forward_port_status.branches_pending[0] // empty' \
            $FIX_STATE_FILE)

        if [ -z "$NEXT_FORWARD" ]; then
            echo "✅ All forward ports complete"
            jq '.current_state = "VALIDATION"' $FIX_STATE_FILE > tmp.json && \
                mv tmp.json $FIX_STATE_FILE
            NEXT_STATE="VALIDATION"
        else
            echo "📋 Forward porting: $NEXT_FORWARD"

            # Move from pending to current
            jq --arg branch "$NEXT_FORWARD" \
                '.forward_port_status.current_branch = $branch |
                 .forward_port_status.branches_pending -= [$branch]' \
                $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

            # Determine if this is a gitea-client branch needing recreation
            if [[ "$NEXT_FORWARD" =~ gitea-client-split ]]; then
                echo "Creating new pr-ready from fixed effort branch..."
                ACTION="recreate_from_fixed"
            else
                echo "Rebasing existing pr-ready branch..."
                ACTION="rebase_on_previous"
            fi

            # Spawn software engineer
            Task: software-engineer
            Working directory: /tmp/gitea-fix-forward/$NEXT_FORWARD
            State: FIX_FORWARD_PORT
            Instructions:
            Forward port fix to: $NEXT_FORWARD
            Action: $ACTION

            Dependencies: [list previous branches in order]

            Steps:
            1. If gitea-client branch: Create new pr-ready from fixed effort
            2. Otherwise: Rebase existing pr-ready onto previous branch
            3. Resolve conflicts maintaining fix changes
            4. Verify it builds
            5. Push updated branch as ${NEXT_FORWARD}-fixed

            NEXT_STATE="MONITOR_FORWARD_PORT"
        fi
        ;;

    "VALIDATION")
        echo "═══════════════════════════════════════════════════════"
        echo "PHASE 4: VALIDATION - Testing Fixed Implementation"
        echo "═══════════════════════════════════════════════════════"

        # Check validation status
        MERGE_STATUS=$(jq -r '.validation_status.merge_test' $FIX_STATE_FILE)
        BUILD_STATUS=$(jq -r '.validation_status.build_test' $FIX_STATE_FILE)
        PUSH_STATUS=$(jq -r '.validation_status.push_test' $FIX_STATE_FILE)

        if [ "$MERGE_STATUS" = "NOT_STARTED" ]; then
            echo "📋 Starting merge validation..."

            jq '.validation_status.merge_test = "IN_PROGRESS"' \
                $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

            Task: software-engineer
            Working directory: /tmp/gitea-fix-validation
            State: FIX_VALIDATION_MERGE
            Instructions:
            Perform comprehensive merge test of all fixed branches.

            Branches to merge (in order): [list all 12 fixed pr-ready branches]

            Steps:
            1. Create fresh clone
            2. Merge all branches sequentially
            3. Build: go build -o idpbuilder .
            4. Report success/failure

            NEXT_STATE="MONITOR_VALIDATION"

        elif [ "$BUILD_STATUS" = "NOT_STARTED" ] && [ "$MERGE_STATUS" = "PASSED" ]; then
            echo "📋 Starting build validation..."

            jq '.validation_status.build_test = "IN_PROGRESS"' \
                $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

            # Build test happens in merge validation
            jq '.validation_status.build_test = "PASSED"' \
                $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

        elif [ "$PUSH_STATUS" = "NOT_STARTED" ] && [ "$BUILD_STATUS" = "PASSED" ]; then
            echo "📋 Starting push functionality test..."

            jq '.validation_status.push_test = "IN_PROGRESS"' \
                $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

            Task: software-engineer
            Working directory: /tmp/gitea-fix-validation
            State: FIX_VALIDATION_PUSH
            Instructions:
            Test push functionality with fixed implementation.

            Credentials:
            - Registry: gitea.cnoe.localtest.me:8443
            - Username: giteaadmin
            - Token: dd547a556201ac9571d2328c3677c7ddcd52d3e4

            Steps:
            1. Build test image with idpbuilder build
            2. Push with idpbuilder push
            3. Verify NO "not yet implemented" error
            4. Confirm real implementation called

            NEXT_STATE="MONITOR_VALIDATION"

        elif [ "$PUSH_STATUS" = "PASSED" ]; then
            echo "✅ ALL VALIDATION COMPLETE!"

            jq '.current_state = "INTEGRATION" |
                .timestamps.completed = now' \
                $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

            NEXT_STATE="INTEGRATION"
        fi
        ;;

    "INTEGRATION")
        echo "═══════════════════════════════════════════════════════"
        echo "PHASE 5: INTEGRATION - Finalizing Fixed Branches"
        echo "═══════════════════════════════════════════════════════"

        echo "Creating final integration report..."

        # Create integration report
        cat > GITEA-FIX-INTEGRATION-REPORT.md << 'EOF'
# GiteaClient Interface Fix - Integration Report

## Fix Summary
Successfully aligned GiteaClient with Registry interface across all branches.

## Branches Updated
### Backported (4 branches):
- gitea-client-split-001-fixed
- gitea-client-split-002-fixed
- gitea-client-split-003-fixed
- gitea-client-fixed

### Forward Ported (12 branches):
[List all pr-ready-fixed branches]

## Validation Results
- Merge Test: PASSED
- Build Test: PASSED
- Push Test: PASSED - Real implementation now called

## Next Steps
1. Replace original branches with fixed versions
2. Archive pre-fix branches
3. Create pull requests from fixed branches
EOF

        echo "✅ FIX CASCADE COMPLETE!"

        jq '.current_state = "COMPLETE"' $FIX_STATE_FILE > tmp.json && \
            mv tmp.json $FIX_STATE_FILE

        # Update main orchestrator state
        jq '.current_state = "WAVE_COMPLETE" |
            .fix_cascade_completed = true |
            .fix_cascade_id = "gitea-client-interface-alignment"' \
            orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

        echo "🎉 Fix cascade completed successfully!"
        ;;

    *)
        echo "Unknown fix state: $FIX_CURRENT_STATE"
        ;;
esac
```

## 🔴 R354 ENFORCEMENT - POST-FIX REVIEWS

### After Each Fix Operation:
```bash
# Check if review needed
NEEDS_REVIEW=$(jq -r '.pending_review // false' $FIX_STATE_FILE)

if [ "$NEEDS_REVIEW" = "true" ]; then
    echo "🔴 R354: Post-fix review required!"

    REVIEW_BRANCH=$(jq -r '.current_review_branch' $FIX_STATE_FILE)

    echo "Spawning Code Reviewer for $REVIEW_BRANCH..."

    Task: code-reviewer
    Working directory: [branch directory]
    State: POST_FIX_REVIEW
    Instructions:
    - Review fixed branch for interface alignment
    - Verify: builds pass, tests pass, interface matches
    - Check: no regression, maintains compatibility
    - Return: FIX_VALID or FIXES_NEEDED

    NEXT_STATE="MONITOR_REVIEWS"
fi
```

## 📊 FIX TRACKING AND STATE UPDATES

### Progress Visualization:
```bash
show_fix_progress() {
    echo "═══════════════════════════════════════════════════════"
    echo "📊 FIX CASCADE PROGRESS"
    echo "═══════════════════════════════════════════════════════"

    # Backport progress
    BACKPORT_COMPLETE=$(jq '[.backport_status[] | select(. == "COMPLETED")] | length' \
        $FIX_STATE_FILE)
    BACKPORT_TOTAL=$(jq '.backport_status | length' $FIX_STATE_FILE)
    echo "🔧 Backports: $BACKPORT_COMPLETE/$BACKPORT_TOTAL completed"

    # Forward port progress
    FORWARD_COMPLETE=$(jq '.forward_port_status.branches_completed | length' \
        $FIX_STATE_FILE)
    FORWARD_TOTAL=$((FORWARD_COMPLETE + \
        $(jq '.forward_port_status.branches_pending | length' $FIX_STATE_FILE)))
    echo "📦 Forward Ports: $FORWARD_COMPLETE/$FORWARD_TOTAL completed"

    # Validation status
    echo "✅ Validation:"
    echo "   - Merge Test: $(jq -r '.validation_status.merge_test' $FIX_STATE_FILE)"
    echo "   - Build Test: $(jq -r '.validation_status.build_test' $FIX_STATE_FILE)"
    echo "   - Push Test: $(jq -r '.validation_status.push_test' $FIX_STATE_FILE)"

    echo "═══════════════════════════════════════════════════════"
}

show_fix_progress
```

## 📝 CHECKPOINT SAVING

### Before State Transitions:
```bash
save_fix_checkpoint() {
    local FIX_STATE=$1
    local DETAIL=$2

    echo "💾 Saving fix cascade checkpoint..."

    # Update fix state
    jq --arg detail "$DETAIL" \
        '.timestamps.last_updated = now |
         .last_checkpoint = {
            "state": "'$FIX_STATE'",
            "detail": $detail,
            "timestamp": now
         }' $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

    # Save TODOs
    TODO_FILE="todos/fix-cascade-$(date '+%Y%m%d-%H%M%S').todo"
    echo "[FIX CASCADE: GiteaClient Interface]" > "$TODO_FILE"
    echo "State: $FIX_STATE" >> "$TODO_FILE"
    echo "Detail: $DETAIL" >> "$TODO_FILE"

    # Commit both state files
    git add orchestrator-state.json $FIX_STATE_FILE todos/*.todo
    git commit -m "fix-cascade: $FIX_STATE - $DETAIL

FIX CASCADE CHECKPOINT:
- Fix: GiteaClient interface alignment
- State: $FIX_STATE
- Progress: $(show_fix_progress | grep completed)
- Can resume with: /fix-cascade-orchestration"

    git push

    echo "✅ Checkpoint saved - Resume anytime with /fix-cascade-orchestration"
}

# Always save before stopping
save_fix_checkpoint "$FIX_CURRENT_STATE" "State transition to $NEXT_STATE"
```

## ⚠️ ERROR HANDLING

### Fix Application Failures:
```bash
if [[ "$FIX_RESULT" == "FAILED" ]]; then
    echo "❌ Fix application failed!"

    ERROR_DETAIL=$(jq -r '.last_error' $FIX_STATE_FILE)

    # Log error
    jq --arg error "$ERROR_DETAIL" \
        '.error_log += [{
            "timestamp": now,
            "phase": "'$FIX_CURRENT_STATE'",
            "error": $error
        }]' $FIX_STATE_FILE > tmp.json && mv tmp.json $FIX_STATE_FILE

    # Determine recovery
    case "$FIX_CURRENT_STATE" in
        "BACKPORT")
            echo "Retry backport with adjusted patch..."
            ;;
        "FORWARD_PORT")
            echo "Manual conflict resolution needed..."
            ;;
        "VALIDATION")
            echo "Fix incomplete - return to analysis..."
            jq '.current_state = "ANALYSIS"' $FIX_STATE_FILE > tmp.json && \
                mv tmp.json $FIX_STATE_FILE
            ;;
    esac
fi
```

## 🛑 FIX CASCADE COMPLETION CRITERIA

The fix cascade is complete when:
- [ ] All 4 effort branches have fix backported
- [ ] All 12 pr-ready branches updated with fix
- [ ] Comprehensive merge succeeds
- [ ] Build produces working binary
- [ ] Push command uses real implementation (not stub)
- [ ] Integration report created
- [ ] State shows: COMPLETE

## 🎯 USAGE INSTRUCTIONS

This command should be invoked repeatedly throughout the fix cascade:
```bash
/fix-cascade-orchestration  # Start/continue fix cascade
# [Orchestrator executes current phase]
# [Orchestrator stops at state boundary]

/fix-cascade-orchestration  # Continue to next phase
# [Repeat until fix complete]
```

## 📊 GRADING COMPLIANCE

You will be graded on:
- ✅ Fix applied in correct order (backport → forward-port)
- ✅ R354: Reviews for every fixed branch
- ✅ R151: Parallel spawning where applicable
- ✅ R322: Stopping at EVERY state boundary
- ✅ Complete state tracking in both files
- ✅ Successful validation of fix

## 🔍 QUICK STATUS CHECK

```bash
# Check current position in fix cascade
echo "Main State: $(jq -r '.current_state' orchestrator-state.json)"
echo "Fix State: $(jq -r '.current_state' $FIX_STATE_FILE)"
echo "Fix Phase: $(jq -r '.current_phase' $FIX_STATE_FILE)"

# Show what's next
case "$(jq -r '.current_state' $FIX_STATE_FILE)" in
    "INIT") echo "Next: Create fix implementation" ;;
    "ANALYSIS") echo "Next: Start backporting to efforts" ;;
    "BACKPORT") echo "Next: Continue backporting or start forward-port" ;;
    "FORWARD_PORT") echo "Next: Continue forward-porting branches" ;;
    "VALIDATION") echo "Next: Run validation tests" ;;
    "INTEGRATION") echo "Next: Finalize and complete" ;;
    "COMPLETE") echo "Fix cascade complete!" ;;
esac
```

---

**Remember**:
- Track state in BOTH orchestrator-state.json AND orchestrator-gitea-fix-state.json
- EVERY fix needs review (R354)
- Stop at EVERY state boundary (R322)
- Execute in exact order (R351)
- The cascade continues until ALL operations complete

**"Fix in order, review every change, track everything, stop at boundaries"**