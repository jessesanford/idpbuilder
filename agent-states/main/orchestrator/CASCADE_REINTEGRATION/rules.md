# Orchestrator - CASCADE_REINTEGRATION State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR CASCADE_REINTEGRATION STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322 Part A-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R360** - Cascade Reintegration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R360-cascade-reintegration-protocol.md`
   - Criticality: SUPREME LAW - Cascade management requirements
   - Summary: Handle cascading changes between dependent efforts

## 🔴🔴🔴 SUPREME DIRECTIVE: CASCADE IS PERSISTENT COORDINATOR 🔴🔴🔴

**THIS STATE ENFORCES R327/R352/R353 CASCADE LAWS - SUPPORTS OVERLAPPING CASCADES WITH ABSOLUTE FOCUS!**

CASCADE_REINTEGRATION is a PERSISTENT COORDINATOR that:
1. Maintains control after EVERY operation
2. Supports MULTIPLE overlapping cascade chains
3. Checks for new fixes after each operation
4. Returns here from ALL integration states when cascade_mode=true
5. **ENFORCES R353 CASCADE FOCUS - NO DIVERSIONS ALLOWED!**

**YOU CANNOT LEAVE THIS STATE UNTIL:**
- All cascade chains are complete or merged
- No pending fixes exist without chains
- All integrations are fresh
- No new fixes detected in final check

## 🔴🔴🔴 R353 CASCADE FOCUS PROTOCOL - SUPREME LAW 🔴🔴🔴

**DURING CASCADE OPERATIONS, NO DIVERSIONS ARE ALLOWED!**
- ❌ NO size checks or split evaluations
- ❌ NO transitions to split-related states
- ❌ NO quality assessments beyond "does it build"
- ✅ ONLY validate rebases and check for conflicts
- ✅ ALL spawned agents MUST receive cascade_mode=true context

## State Overview

In CASCADE_REINTEGRATION, you are the PERSISTENT COORDINATOR for ALL cascade operations. This state:
- Maintains control throughout ALL cascade operations (R352)
- Supports multiple overlapping cascade chains running simultaneously
- Returns here after EVERY integration state when cascade_mode=true
- Checks for new fixes after each operation
- Only releases control when ALL fixes reach project level

## Required Actions

### 1. 🔴🔴🔴 DETECT AND MANAGE CASCADE CHAINS (R352) 🔴🔴🔴

```bash
echo "🔴🔴🔴 R327/R352 CASCADE COORDINATION ACTIVATED 🔴🔴🔴"
echo "Managing overlapping cascade chains..."

# CRITICAL: Set cascade_mode for persistent coordination (R352/R353)
echo "📝 Setting cascade_mode=true for persistent CASCADE_REINTEGRATION control"
echo "🔴 R353 CASCADE FOCUS PROTOCOL ACTIVE - NO DIVERSIONS ALLOWED"
jq '.cascade_coordination.cascade_mode = true |
    .cascade_coordination.persistent_coordination = true |
    .cascade_coordination.cascade_return_state = "CASCADE_REINTEGRATION" |
    .cascade_coordination.cascade_focus = {
        "allow_splits": false,
        "allow_size_checks": false,
        "only_validate_rebases": true,
        "protocol": "R353"
    }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# Initialize cascade coordination structure if needed
jq 'if .cascade_coordination.active_cascade_chains == null then
    .cascade_coordination.active_cascade_chains = []
else . end |
if .cascade_coordination.pending_fixes == null then
    .cascade_coordination.pending_fixes = {}
else . end' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# Check for new fixes that arrived (R352 requirement)
echo "🔍 Checking for new fixes that arrived during cascade operations..."
LAST_CHECK=$(jq -r '.cascade_coordination.last_fix_check // "1hour ago"' orchestrator-state.json)

for effort_dir in /efforts/phase*/wave*/effort-*; do
    if [[ -d "$effort_dir" ]]; then
        cd "$effort_dir"
        branch=$(git branch --show-current)
        NEW_FIXES=$(git log --since="$LAST_CHECK" --oneline --grep="fix:" 2>/dev/null | wc -l)
        
        if [[ "$NEW_FIXES" -gt 0 ]]; then
            echo "⚠️ Found $NEW_FIXES new fixes in $branch - creating new cascade chain"
            # Create new cascade chain for these fixes
            CHAIN_ID="cascade_$(date +%s)_${branch//\//_}"
            FIX_IDS=$(git log --since="$LAST_CHECK" --oneline --grep="fix:" --format=%H)
            
            jq --arg chain "$CHAIN_ID" --arg branch "$branch" --arg fixes "$FIX_IDS" '
                .cascade_coordination.active_cascade_chains += [{
                    "chain_id": $chain,
                    "trigger": {
                        "type": "fix_applied",
                        "location": $branch,
                        "timestamp": now | todate,
                        "fix_ids": ($fixes | split("\n") | map(select(. != "")))
                    },
                    "status": "pending",
                    "operations": [],
                    "started_at": now | todate
                }] |
                .cascade_coordination.pending_fixes[$branch] = {
                    "fix_ids": ($fixes | split("\n") | map(select(. != ""))),
                    "applied_at": now | todate,
                    "cascade_chain": $chain
                }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
        fi
    fi
done

# Update last check timestamp
jq --arg ts "$(date -Iseconds)" '.cascade_coordination.last_fix_check = $ts' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# Use stale-integration-manager to check all levels
./utilities/stale-integration-manager.sh check

# Check for active cascade chains (R352)
ACTIVE_CHAINS=$(jq -r '.cascade_coordination.active_cascade_chains[]? |
                       select(.status == "in_progress" or .status == "pending") |
                       .chain_id' orchestrator-state.json)

if [[ -z "$ACTIVE_CHAINS" ]]; then
    echo "🔍 No active cascade chains - checking if we can exit..."
    
    # R352: Check all exit conditions
    CAN_EXIT=true
    
    # Check for pending fixes
    PENDING_FIXES=$(jq '.cascade_coordination.pending_fixes | length' orchestrator-state.json)
    if [[ "$PENDING_FIXES" -gt 0 ]]; then
        echo "❌ Pending fixes exist without cascade chains"
        CAN_EXIT=false
    fi
    
    # Check for stale integrations
    STALE_COUNT=$(jq '[.stale_integration_tracking.stale_integrations[]? | 
                      select(.recreation_completed != true)] | length' orchestrator-state.json)
    if [[ "$STALE_COUNT" -gt 0 ]]; then
        echo "❌ $STALE_COUNT stale integrations remain"
        CAN_EXIT=false
    fi
    
    # Check all integration freshness
    for level in wave phase project; do
        IS_STALE=$(jq -r ".current_${level}_integration.is_stale // false" orchestrator-state.json)
        if [[ "$IS_STALE" == "true" ]]; then
            echo "❌ ${level} integration is still stale"
            CAN_EXIT=false
        fi
    done
    
    if [[ "$CAN_EXIT" == "true" ]]; then
        echo "✅ All cascade operations complete - can exit CASCADE_REINTEGRATION"
        # Clear cascade mode
        jq '.cascade_coordination.cascade_mode = false |
            .cascade_coordination.persistent_coordination = false' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
        NEXT_STATE="INTEGRATION_CODE_REVIEW"
    else
        echo "🔴 Cannot exit - cascade conditions not met"
        # Check legacy staleness cascade for operations
        CASCADE_PENDING=$(jq -r '
            .stale_integration_tracking.staleness_cascade[]? |
            select(.cascade_status != "completed") |
            "Cascade required from " + .trigger.branch
        ' orchestrator-state.json)
        
        if [[ -n "$CASCADE_PENDING" ]]; then
            echo "🔴 LEGACY CASCADE REQUIREMENTS DETECTED:"
            echo "$CASCADE_PENDING"
        fi
    fi
else
    echo "🔴 ACTIVE CASCADE CHAINS: $(echo "$ACTIVE_CHAINS" | wc -w)"
    for chain in $ACTIVE_CHAINS; do
        echo "  📋 $chain"
    done
    
    # Check for chain convergence
    echo "🔍 Checking for cascade chain convergence..."
    # Group chains by their current operation targets
    CONVERGENCE_CHECK=$(jq -r '
        [.cascade_coordination.active_cascade_chains[] |
         select(.status == "in_progress" or .status == "pending") |
         .operations[-1].target // "none"] |
        group_by(.) |
        map(select(length > 1)) |
        if length > 0 then "convergence detected" else empty end' orchestrator-state.json 2>/dev/null)
    
    if [[ -n "$CONVERGENCE_CHECK" ]]; then
        echo "🔀 Cascade chains are converging - will merge when they meet"
    fi
fi
```

### 2. 🔴🔴🔴 ENFORCE POST-REBASE REVIEWS (R354) 🔴🔴🔴

**SUPREME LAW: EVERY REBASE REQUIRES IMMEDIATE REVIEW!**

```bash
# Function to enforce R354 after any rebase operation
enforce_post_rebase_review() {
    local effort="$1"
    local rebased_to="$2"
    local chain_id="$3"

    echo "🔴🔴🔴 R354 ENFORCEMENT: Post-rebase review MANDATORY 🔴🔴🔴"
    echo "Effort: $effort"
    echo "Rebased to: $rebased_to"

    # Add to pending reviews with R354 mandate
    jq --arg e "$effort" --arg base "$rebased_to" --arg chain "$chain_id" '
        .cascade_coordination.pending_reviews += [{
            "effort": $e,
            "rebased_to": $base,
            "review_type": "post_rebase",
            "r354_mandated": true,
            "cascade_mode": true,
            "chain_id": $chain,
            "timestamp": now | todate,
            "review_status": "pending"
        }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Set R354 enforcement flag
    jq --arg e "$effort" '
        .cascade_coordination.r354_enforcement = {
            "active": true,
            "blocking_effort": $e,
            "waiting_for": "post_rebase_review",
            "enforced_at": now | todate
        }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    echo "⏸️ CASCADE PAUSED per R354 - Awaiting post-rebase review"
    echo "MUST spawn Code Reviewer before continuing cascade!"

    # Return code 354 to signal review required
    return 354
}

# Check for pending post-rebase reviews (R354)
check_pending_rebase_reviews() {
    echo "🔍 R354 Check: Looking for pending post-rebase reviews..."

    PENDING_REVIEWS=$(jq -r '
        .cascade_coordination.pending_reviews[]? |
        select(.review_type == "post_rebase" and .review_status == "pending") |
        .effort' orchestrator-state.json)

    if [[ -n "$PENDING_REVIEWS" ]]; then
        echo "🔴 R354 BLOCKING: Post-rebase reviews required for:"
        echo "$PENDING_REVIEWS" | while read effort; do
            echo "  - $effort (MUST review before cascade continues)"
        done

        echo "🚨 MANDATORY ACTION: Transition to SPAWN_CODE_REVIEWERS_FOR_REVIEW"
        NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_REVIEW"

        # Set cascade context for reviewers
        jq '.cascade_coordination.reviewer_context = {
            "cascade_mode": true,
            "review_type": "post_rebase",
            "r354_enforcement": true,
            "skip_quality_checks": true,
            "focus": "integration_validation_only"
        }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

        return 354
    fi

    echo "✅ No pending post-rebase reviews"
    return 0
}

# FIRST THING: Check R354 compliance
if ! check_pending_rebase_reviews; then
    echo "🛑 R354 ENFORCEMENT: Cannot continue cascade until reviews complete!"
    # Will transition to SPAWN_CODE_REVIEWERS_FOR_REVIEW
fi
```

### 3. EXECUTE CASCADE WITH OVERLAPPING SUPPORT (R352 + R354)

**CRITICAL: Process operations from ANY active cascade chain WITH R354 ENFORCEMENT**

```bash
# Function to recreate integration at specified level
recreate_integration() {
    local level="$1"  # wave, phase, or project
    local integration_id="$2"
    local chain_id="$3"
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🔴 RECREATING ${level^^} INTEGRATION: $integration_id"
    echo "📋 Chain: $chain_id"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Delete old integration branch
    echo "📝 Deleting stale integration branch..."
    git push origin --delete "$integration_id" 2>/dev/null || {
        echo "⚠️ Branch may not exist on remote, continuing..."
    }
    
    # Remove local workspace
    case "$level" in
        wave)
            WORKSPACE="/efforts/phase${PHASE}/wave${WAVE}/wave-integration"
            ;;
        phase)
            WORKSPACE="/efforts/phase${PHASE}/phase-integration"
            ;;
        project)
            WORKSPACE="/efforts/project-integration"
            ;;
    esac
    
    echo "📝 Removing old workspace: $WORKSPACE"
    rm -rf "$WORKSPACE"
    
    # Mark as pending recreation in state
    jq --arg id "$integration_id" '
        (.stale_integration_tracking.stale_integrations[]? | 
         select(.integration_id == $id)) |= 
         . + {"recreation_in_progress": true}
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    # Add operation to cascade chain
    jq --arg chain "$chain_id" --arg level "$level" --arg target "$integration_id" '
        (.cascade_coordination.active_cascade_chains[] | 
         select(.chain_id == $chain) | .operations) += [{
            "type": "recreate",
            "target": $target,
            "level": $level,
            "status": "in_progress",
            "started_at": now | todate
        }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    echo "✅ Old integration removed"
    echo "🔄 Transitioning to appropriate INTEGRATION state for recreation..."
    
    return 0
}

# Process cascade chains (R352 - multiple chains supported)
echo "🔍 Processing cascade chains..."

# Get next operation from any active chain
NEXT_CHAIN_OP=$(jq -r '
    [.cascade_coordination.active_cascade_chains[]? |
     select(.status == "in_progress" or .status == "pending")] |
    if length > 0 then .[0] else empty end' orchestrator-state.json)

if [[ -n "$NEXT_CHAIN_OP" ]]; then
    CHAIN_ID=$(echo "$NEXT_CHAIN_OP" | jq -r '.chain_id')
    echo "🎯 Processing chain: $CHAIN_ID"
    
    # Determine what needs to be done for this chain
    # This is simplified - in reality would check dependencies
    echo "📋 Analyzing cascade requirements for chain $CHAIN_ID"
fi

# Also check legacy staleness cascade
CASCADES=$(jq -r '
    .stale_integration_tracking.staleness_cascade[]? |
    select(.cascade_status != "completed") |
    .cascade_sequence[] |
    select(.recreation_status != "completed") |
    .level + ":" + .integration
' orchestrator-state.json | head -1)  # Process one at a time

if [[ -n "$CASCADES" ]]; then
    IFS=':' read -r LEVEL INTEGRATION_ID <<< "$CASCADES"
    
    echo "🎯 Next cascade target: $LEVEL - $INTEGRATION_ID"
    
    # Use the first active chain or create a default one
    ACTIVE_CHAIN=$(jq -r '.cascade_coordination.active_cascade_chains[0].chain_id // "default_cascade"' orchestrator-state.json)
    
    # Determine next state based on level
    case "$LEVEL" in
        wave)
            recreate_integration "wave" "$INTEGRATION_ID" "$ACTIVE_CHAIN"
            NEXT_STATE="INTEGRATION"  # Go to wave integration
            # Track current cascade operation
            jq --arg op "wave_reintegration" --arg target "$INTEGRATION_ID" '
                .cascade_coordination.current_cascade_operation = $op |
                .cascade_coordination.cascade_progress.current_operation_type = "recreate" |
                .cascade_coordination.cascade_progress.current_operation_target = $target
            ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
            ;;
        phase)
            recreate_integration "phase" "$INTEGRATION_ID" "$ACTIVE_CHAIN"
            NEXT_STATE="PHASE_INTEGRATION"  # Go to phase integration
            # Track current cascade operation
            jq --arg op "phase_reintegration" --arg target "$INTEGRATION_ID" '
                .cascade_coordination.current_cascade_operation = $op |
                .cascade_coordination.cascade_progress.current_operation_type = "recreate" |
                .cascade_coordination.cascade_progress.current_operation_target = $target
            ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
            ;;
        project)
            recreate_integration "project" "$INTEGRATION_ID" "$ACTIVE_CHAIN"
            NEXT_STATE="PROJECT_INTEGRATION"  # Go to project integration
            # Track current cascade operation
            jq --arg op "project_reintegration" --arg target "$INTEGRATION_ID" '
                .cascade_coordination.current_cascade_operation = $op |
                .cascade_coordination.cascade_progress.current_operation_type = "recreate" |
                .cascade_coordination.cascade_progress.current_operation_target = $target
            ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
            ;;
        *)
            echo "❌ Unknown integration level: $LEVEL"
            NEXT_STATE="ERROR_RECOVERY"
            ;;
    esac
    
    # Update cascade status
    jq --arg level "$LEVEL" --arg time "$(date -Iseconds)" '
        (.stale_integration_tracking.staleness_cascade[]? |
         .cascade_sequence[] |
         select(.level == $level)) |= 
         . + {"recreation_status": "in_progress", "started_at": $time}
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
elif [[ "$NEXT_STATE" != "INTEGRATION_CODE_REVIEW" ]]; then
    echo "✅ All current cascade operations complete!"
    # Check one more time for new fixes
    echo "🔍 Final check for new fixes..."
    # Would re-run the new fix detection here
    NEXT_STATE="INTEGRATION_CODE_REVIEW"
fi

# 🔴🔴🔴 R354 ENFORCEMENT: After ANY rebase operation 🔴🔴🔴
# Example: After rebasing an effort onto new integration
REBASE_JUST_COMPLETED="${REBASE_JUST_COMPLETED:-false}"
if [[ "$REBASE_JUST_COMPLETED" == "true" ]]; then
    echo "🔴 R354 TRIGGERED: Post-rebase review required!"

    # Call the enforce function
    enforce_post_rebase_review "$REBASED_EFFORT" "$NEW_BASE" "$ACTIVE_CHAIN" || {
        echo "🛑 R354 BLOCKING CASCADE - Must review rebased code"
        NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_REVIEW"
    }
fi
```

### 3. CHAIN CONVERGENCE AND MERGING (R352)

```bash
# Check for converging cascade chains
check_chain_convergence() {
    local TARGET="$1"
    
    echo "🔍 Checking for cascade chain convergence at $TARGET"
    
    # Find all chains with operations targeting this integration
    CONVERGING_CHAINS=$(jq -r --arg target "$TARGET" '
        [.cascade_coordination.active_cascade_chains[] |
         select(.operations[].target == $target) |
         .chain_id] | unique | @json' orchestrator-state.json)
    
    CHAIN_COUNT=$(echo "$CONVERGING_CHAINS" | jq 'length')
    
    if [[ "$CHAIN_COUNT" -gt 1 ]]; then
        echo "🔀 $CHAIN_COUNT cascade chains converging at $TARGET"
        
        # Get chain IDs
        PRIMARY=$(echo "$CONVERGING_CHAINS" | jq -r '.[0]')
        SECONDARY=$(echo "$CONVERGING_CHAINS" | jq -r '.[1:][]')
        
        # Merge chains
        echo "📋 Merging chains: Primary=$PRIMARY, Secondary=$SECONDARY"
        
        # Collect all fixes from all chains
        jq --arg primary "$PRIMARY" --arg secondary "$SECONDARY" --arg target "$TARGET" '
            # Find primary chain
            (.cascade_coordination.active_cascade_chains[] | 
             select(.chain_id == $primary)) |= . + {
                merged_with: ((.merged_with // []) + ($secondary | split(" "))),
                merge_point: $target,
                merge_timestamp: now | todate
            } |
            # Mark secondary chains as merged
            (.cascade_coordination.active_cascade_chains[] |
             select(.chain_id | IN($secondary | split(" ")[]))) |= . + {
                status: "merged_into",
                merged_into: $primary,
                merge_timestamp: now | todate
            }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
        
        echo "✅ Chains merged successfully"
    fi
}

# Call after each operation completes
if [[ -n "$INTEGRATION_ID" ]]; then
    check_chain_convergence "$INTEGRATION_ID"
fi
```

### 4. VALIDATE CASCADE COMPLETION (R352)

```bash
# Comprehensive cascade completion check per R352
echo "🔍 Validating cascade completion per R352..."

validate_cascade_complete() {
    # 1. Check all cascade chains are complete or merged
    ACTIVE_CHAINS=$(jq -r '
        .cascade_coordination.active_cascade_chains[]? |
        select(.status != "completed" and .status != "merged_into") |
        .chain_id' orchestrator-state.json)
    
    if [[ -n "$ACTIVE_CHAINS" ]]; then
        echo "❌ Active cascade chains remain"
        return 1
    fi
    
    # 2. Check no pending fixes without cascade chains
    ORPHAN_FIXES=$(jq -r '
        .cascade_coordination.pending_fixes | 
        to_entries[] | 
        select(.value.cascade_chain == null) |
        .key' orchestrator-state.json)
    
    if [[ -n "$ORPHAN_FIXES" ]]; then
        echo "❌ Pending fixes without cascade chains: $ORPHAN_FIXES"
        return 1
    fi
    
    # 3. Check all integrations are fresh
    ALL_FRESH=true
    for level in wave phase project; do
        IS_STALE=$(jq -r ".current_${level}_integration.is_stale // false" orchestrator-state.json)
        if [[ "$IS_STALE" == "true" ]]; then
            echo "❌ ${level} integration is still stale!"
            ALL_FRESH=false
        else
            echo "✅ ${level} integration is fresh"
        fi
    done
    
    if [[ "$ALL_FRESH" != "true" ]]; then
        return 1
    fi
    
    # 4. Final check for any stale integrations
    STALE_COUNT=$(jq '[.stale_integration_tracking.stale_integrations[]? | 
                      select(.recreation_completed != true)] | length' orchestrator-state.json)
    
    if [[ "$STALE_COUNT" -gt 0 ]]; then
        echo "❌ $STALE_COUNT stale integrations remain"
        return 1
    fi
    
    echo "✅✅✅ CASCADE COMPLETE - All conditions met per R352!"
    
    # Clear cascade tracking
    jq '
        .cascade_coordination.cascade_mode = false |
        .cascade_coordination.persistent_coordination = false |
        .cascade_coordination.current_cascade_operation = null |
        .cascade_coordination.active_cascade_chains = 
            (.cascade_coordination.active_cascade_chains | 
             map(. + {status: "completed", completed_at: now|todate})) |
        .stale_integration_tracking.staleness_cascade = [] |
        .stale_integration_tracking.stale_integrations = 
            (.stale_integration_tracking.stale_integrations | 
             map(. + {recreation_completed: true, recreation_at: now|todate}))
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    return 0
}

if ! validate_cascade_complete; then
    echo "🔴 CASCADE INCOMPLETE - Cannot exit CASCADE_REINTEGRATION!"
    echo "Must continue cascade process..."
    # Stay in CASCADE_REINTEGRATION - it will loop back
fi
```

### 5. Update State File

```bash
# Update orchestrator state
jq ".current_state = \"$NEXT_STATE\"" orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq ".state_transition_history += [{
    \"from\": \"CASCADE_REINTEGRATION\", 
    \"to\": \"$NEXT_STATE\", 
    \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", 
    \"reason\": \"Cascade operation: ${LEVEL:-completion} (R327/R352)\"
}]" orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# Commit state change
git add orchestrator-state.json
git commit -m "state: CASCADE_REINTEGRATION → $NEXT_STATE - R327/R352 cascade coordination"
git push

# 🛑 STOP per R322 - State has been updated, now stop!
echo "🛑 Stopping before $NEXT_STATE state (per R322)"
echo "Next action: ${LEVEL:-Review} integration"

# CRITICAL R352: If cascade_mode is still true, we'll return here
if [[ "$(jq -r '.cascade_coordination.cascade_mode' orchestrator-state.json)" == "true" ]]; then
    echo "🔄 CASCADE MODE ACTIVE - Will return to CASCADE_REINTEGRATION after $NEXT_STATE"
fi

echo "When restarted with /continue-orchestrating, will continue from $NEXT_STATE"
```

## Valid Transitions

Based on R327/R352/R353 cascade requirements:

1. **To Integration States (with return)**: 
   - `CASCADE_REINTEGRATION` → `INTEGRATION` → `CASCADE_REINTEGRATION`
   - `CASCADE_REINTEGRATION` → `PHASE_INTEGRATION` → `CASCADE_REINTEGRATION`
   - `CASCADE_REINTEGRATION` → `PROJECT_INTEGRATION` → `CASCADE_REINTEGRATION`
   
   **CRITICAL (R353)**: cascade_mode MUST remain true during these transitions!

2. **Self-loop for more operations**:
   - `CASCADE_REINTEGRATION` → `CASCADE_REINTEGRATION` (processing chains)
   
3. **Completion (only when ALL conditions met)**:
   - `CASCADE_REINTEGRATION` → `INTEGRATION_CODE_REVIEW` (R352 conditions satisfied)

4. **Error**:
   - `CASCADE_REINTEGRATION` → `ERROR_RECOVERY` (cascade failed)

## 🔴🔴🔴 FORBIDDEN TRANSITIONS DURING CASCADE (R353) 🔴🔴🔴

**THESE TRANSITIONS ARE ABSOLUTELY FORBIDDEN WHILE cascade_mode=true:**
- ❌ CASCADE_REINTEGRATION → CREATE_NEXT_SPLIT_INFRASTRUCTURE
- ❌ CASCADE_REINTEGRATION → SPLIT_PLANNING  
- ❌ CASCADE_REINTEGRATION → ANALYZE_SPLIT_REQUIREMENTS
- ❌ Any state → Split-related state (when cascade_mode=true)

## Grading Criteria

- ✅ **+30%**: Support multiple overlapping cascade chains (R352)
- ✅ **+30%**: Check for new fixes after each operation (R352)
- ✅ **+20%**: Properly merge converging chains (R352)
- ✅ **+20%**: Validate ALL R352 exit conditions before leaving

## Common Violations

- ❌ **-100%**: Exiting with pending fixes not cascaded
- ❌ **-100%**: Not returning to CASCADE_REINTEGRATION when cascade_mode=true
- ❌ **-50%**: Losing track of cascade chains
- ❌ **-50%**: Not checking for new fixes during operations

## Related Rules

- R327: Mandatory Re-Integration After Fixes (base cascade requirement)
- R348: Cascade State Transitions (state enforcement)
- R350: Complete Cascade Dependency Graph (dependency calculation)
- R351: Cascade Execution Protocol (execution order)
- R352: Overlapping Cascade Protocol (SUPREME - multiple chains)
- R353: Cascade Focus Protocol (SUPREME - NO diversions during cascade)
- R354: Post-Rebase Review Requirement (SUPREME - every rebase needs review)
- R322: Mandatory Stop After State Transitions

## CASCADE COORDINATION MANTRA (R352)

```
Multiple fixes may arrive,
While cascades are in flight.
CASCADE_REINTEGRATION stays alive,
Until ALL reach project site.

Chains may merge, chains may split,
New fixes join the queue.
The coordinator won't quit,
Until EVERY fix is through!
```

## 🔴🔴🔴 REMEMBER: PERSISTENT COORDINATION (R352) 🔴🔴🔴

**CASCADE_REINTEGRATION is a PERSISTENT COORDINATOR:**
- Maintains control after EVERY operation
- Supports unlimited overlapping cascade chains
- Checks for new fixes continuously
- Only releases when ALL fixes reach project
- Integration states return here when cascade_mode=true

**Violation = -100% AUTOMATIC FAILURE**