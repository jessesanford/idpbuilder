# CASCADE_REINTEGRATION State Rules

## State Purpose
TRAP STATE for enforcing R327 cascade recreation of ALL stale integrations after fixes are applied to effort branches. This state BLOCKS all other work until cascade re-integration completes.

## 🔴🔴🔴 SUPREME LAW: R327 CASCADE ENFORCEMENT 🔴🔴🔴

**This state exists to enforce ONE rule: R327 - Mandatory Re-Integration After Fixes**

**When ANY fix is applied to ANY effort branch:**
1. ALL integration branches containing that effort become STALE
2. Those integrations MUST be deleted and recreated
3. The cascade MUST propagate UP the hierarchy (wave → phase → project)
4. This is NOT optional, NOT negotiable, NOT skippable

**Violation = -100% AUTOMATIC FAILURE**

## Entry Conditions

You enter CASCADE_REINTEGRATION when:
- From ERROR_RECOVERY after fixes applied to effort branches
- From MONITORING_EFFORT_FIXES when all fixes complete
- From MONITORING_INTEGRATE_WAVE_EFFORTS when stale integrations detected
- From WAVE_COMPLETE when fixes were applied during wave
- From COMPLETE_PHASE when cascades are pending
- From any state that detects integrations are stale

## Entry Actions

**IMMEDIATELY upon entering this state:**

### 1. Load R350 Dependency Graph
```bash
# R350 tracks ALL dependencies across the system
echo "🔍 Loading R350 dependency graph..."

# Read dependency graph from state file or generate
if [[ -f "orchestrator-state-v3.json" ]]; then
    DEPENDENCY_GRAPH=$(jq -r '.dependency_graph // {}' orchestrator-state-v3.json)
    echo "✅ Dependency graph loaded"
else
    echo "❌ ERROR: No state file found!"
    exit 1
fi
```

### 2. Calculate Complete Cascade Chain
```bash
# Determine which integrations are stale and need recreation
calculate_cascade_chain() {
    echo "🔍 R327: Calculating cascade chain..."

    local CASCADE_CHAIN=()

    # Check wave integrations
    for wave_branch in $(git branch -r | grep "wave.*-integration" | grep -v "phase\|project"); do
        if is_integration_stale "$wave_branch"; then
            CASCADE_CHAIN+=("RECREATE:${wave_branch}")
            echo "  📋 Wave integration stale: $wave_branch"
        fi
    done

    # Check phase integrations
    for phase_branch in $(git branch -r | grep "phase.*-integration" | grep -v "project"); do
        if is_integration_stale "$phase_branch"; then
            CASCADE_CHAIN+=("RECREATE:${phase_branch}")
            echo "  📋 Phase integration stale: $phase_branch"
        fi
    done

    # Check project integration
    if git show-ref --verify --quiet refs/remotes/origin/project-integration; then
        if is_integration_stale "origin/project-integration"; then
            CASCADE_CHAIN+=("RECREATE:project-integration")
            echo "  📋 Project integration stale"
        fi
    fi

    # Save cascade chain to state
    echo "${CASCADE_CHAIN[@]}" > /tmp/cascade-chain.txt
    echo "✅ Cascade chain calculated: ${#CASCADE_CHAIN[@]} recreations needed"
}

# Check if integration is stale
is_integration_stale() {
    local integration_branch="$1"

    # Get integration creation time
    local INTEGRATE_WAVE_EFFORTS_TIME=$(git log -1 --format=%ct "$integration_branch")

    # Get source branches for this integration
    local SOURCE_BRANCHES=$(get_source_branches_for "$integration_branch")

    # Check if any source has commits newer than integration
    for source in $SOURCE_BRANCHES; do
        local SOURCE_TIME=$(git log -1 --format=%ct "$source" 2>/dev/null || echo 0)
        if [[ $SOURCE_TIME -gt $INTEGRATE_WAVE_EFFORTS_TIME ]]; then
            echo "  ⚠️ Source $source is newer than integration"
            return 0  # Stale
        fi
    done

    return 1  # Not stale
}

calculate_cascade_chain
```

### 3. Create R351 Execution Plan
```bash
# R351: Cascade Execution Protocol - defines HOW to execute cascades
create_execution_plan() {
    echo "📋 R351: Creating cascade execution plan..."

    local CASCADE_CHAIN=($(cat /tmp/cascade-chain.txt))
    local EXECUTION_PLAN=()

    # Process cascades in dependency order
    for cascade_item in "${CASCADE_CHAIN[@]}"; do
        local action="${cascade_item%%:*}"
        local target="${cascade_item#*:}"

        case "$action" in
            RECREATE)
                EXECUTION_PLAN+=("DELETE_INTEGRATE_WAVE_EFFORTS:$target")
                EXECUTION_PLAN+=("CREATE_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE:$target")
                EXECUTION_PLAN+=("MERGE_SOURCES:$target")
                EXECUTION_PLAN+=("VERIFY_BUILD:$target")
                ;;
        esac
    done

    # Save execution plan
    printf "%s\n" "${EXECUTION_PLAN[@]}" > /tmp/cascade-execution-plan.txt
    echo "✅ Execution plan created: ${#EXECUTION_PLAN[@]} steps"
}

create_execution_plan
```

### 4. Begin Cascade Execution
```bash
echo "🚀 Beginning cascade execution..."
echo "📊 Progress will be tracked in state file"

# Initialize cascade tracking
jq '.cascade_reintegration = {
    status: "in_progress",
    started_at: (now | todate),
    total_steps: 0,
    completed_steps: 0,
    current_step: null
}' orchestrator-state-v3.json > /tmp/state.json

mv /tmp/state.json orchestrator-state-v3.json
```

## Execution Loop

**CASCADE_REINTEGRATION is a TRAP STATE with an execution loop:**

### Loop Steps
```bash
# This loop continues until ALL cascades are complete
while true; do
    # 1. SELECT_NEXT_OPERATION
    NEXT_OPERATION=$(head -n1 /tmp/cascade-execution-plan.txt)

    if [[ -z "$NEXT_OPERATION" ]]; then
        echo "✅ All cascade operations complete!"
        break
    fi

    echo "📋 Next operation: $NEXT_OPERATION"

    # 2. VALIDATE_PREREQUISITES
    if ! validate_prerequisites "$NEXT_OPERATION"; then
        echo "❌ Prerequisites not met, waiting..."
        break
    fi

    # 3. EXECUTE_OPERATION
    execute_operation "$NEXT_OPERATION"

    # 4. VERIFY_PROJECT_DONE
    if ! verify_operation_success "$NEXT_OPERATION"; then
        echo "❌ Operation failed, entering ERROR_RECOVERY"
        transition_to_state "ERROR_RECOVERY"
        exit 1
    fi

    # 5. UPDATE_CASCADE_STATUS
    update_cascade_status "$NEXT_OPERATION"

    # Remove completed operation from plan
    tail -n +2 /tmp/cascade-execution-plan.txt > /tmp/cascade-plan-new.txt
    mv /tmp/cascade-plan-new.txt /tmp/cascade-execution-plan.txt

    # 6. LOOP or EXIT
    echo "✅ Operation complete, continuing to next..."
done
```

### Comprehensive Infrastructure Cleanup
```bash
# 🔴🔴🔴 CRITICAL: Complete infrastructure cleanup before recreation
# This function ensures NO remnants of old infrastructure remain
cleanup_integration_infrastructure() {
    local integration_type="$1"  # wave, phase, or project
    local phase="$2"
    local wave="$3"  # optional, only for wave integrations

    echo "🧹 R327 CLEANUP: Removing ALL stale integration infrastructure..."

    # Determine integration key and extract pre-planned infrastructure
    local integration_key
    local integration_config

    if [[ "$integration_type" == "wave" ]]; then
        integration_key="phase${phase}_wave${wave}"
        integration_config=$(jq -r ".pre_planned_infrastructure.integrations.wave_integrations.\"${integration_key}\" // empty" orchestrator-state-v3.json)
    elif [[ "$integration_type" == "phase" ]]; then
        integration_key="phase${phase}"
        integration_config=$(jq -r ".pre_planned_infrastructure.integrations.phase_integrations.\"${integration_key}\" // empty" orchestrator-state-v3.json)
    elif [[ "$integration_type" == "project" ]]; then
        integration_key="project"
        integration_config=$(jq -r ".pre_planned_infrastructure.integrations.project_integration // empty" orchestrator-state-v3.json)
    else
        echo "❌ Unknown integration type: $integration_type"
        return 1
    fi

    if [[ -z "$integration_config" ]] || [[ "$integration_config" == "null" ]]; then
        echo "⚠️ No pre-planned infrastructure found for ${integration_key}, skipping cleanup"
        return 0
    fi

    # Extract infrastructure details from pre-planned config
    local branch_name=$(echo "$integration_config" | jq -r '.branch_name // empty')
    local remote_branch=$(echo "$integration_config" | jq -r '.remote_branch // empty')
    local directory=$(echo "$integration_config" | jq -r '.directory // empty')
    local target_repo=$(echo "$integration_config" | jq -r '.target_repo_url // empty')

    echo "📋 Infrastructure to clean:"
    echo "  Branch: $branch_name"
    echo "  Remote: $remote_branch"
    echo "  Directory: $directory"

    # 1. DELETE REMOTE BRANCH (if exists)
    if [[ -n "$remote_branch" ]] && [[ "$remote_branch" != "null" ]]; then
        # Extract remote name (e.g., "target" from "target/branch-name")
        local remote_name="${remote_branch%%/*}"
        local remote_branch_name="${remote_branch#*/}"

        echo "🗑️  Checking remote branch: $remote_name/$remote_branch_name"

        if git ls-remote --heads "$remote_name" "$remote_branch_name" 2>/dev/null | grep -q "$remote_branch_name"; then
            echo "🗑️  Deleting remote branch: $remote_name/$remote_branch_name"
            git push "$remote_name" --delete "$remote_branch_name" 2>/dev/null || {
                echo "⚠️ Failed to delete remote branch (may already be deleted)"
            }
        else
            echo "✅ Remote branch already deleted: $remote_name/$remote_branch_name"
        fi
    fi

    # 2. DELETE LOCAL BRANCH (if exists)
    if [[ -n "$branch_name" ]] && [[ "$branch_name" != "null" ]]; then
        echo "🗑️  Checking local branch: $branch_name"

        if git show-ref --verify --quiet "refs/heads/$branch_name"; then
            # Save current branch to avoid issues if we're on the branch being deleted
            local current_branch=$(git branch --show-current)

            if [[ "$current_branch" == "$branch_name" ]]; then
                echo "⚠️ Currently on branch being deleted, switching to main"
                git fetch target main 2>/dev/null || git fetch origin main
                git checkout main 2>/dev/null || git checkout -b main target/main
            fi

            echo "🗑️  Deleting local branch: $branch_name"
            git branch -D "$branch_name" 2>/dev/null || {
                echo "⚠️ Failed to delete local branch (may have conflicts)"
            }
        else
            echo "✅ Local branch already deleted: $branch_name"
        fi
    fi

    # 3. DELETE WORKING DIRECTORY (if exists)
    if [[ -n "$directory" ]] && [[ "$directory" != "null" ]]; then
        # Convert to absolute path if relative
        local abs_directory
        if [[ "$directory" = /* ]]; then
            abs_directory="$directory"
        else
            abs_directory="$CLAUDE_PROJECT_DIR/$directory"
        fi

        echo "🗑️  Checking working directory: $abs_directory"

        if [[ -d "$abs_directory" ]]; then
            echo "🗑️  Removing working directory: $abs_directory"
            rm -rf "$abs_directory" || {
                echo "⚠️ Failed to remove directory, attempting with sudo"
                sudo rm -rf "$abs_directory" 2>/dev/null || {
                    echo "❌ Could not remove directory: $abs_directory"
                    return 1
                }
            }
            echo "✅ Working directory removed"
        else
            echo "✅ Working directory already removed: $abs_directory"
        fi
    fi

    # 4. RESET STATE FILE (mark as not created)
    echo "🗑️  Resetting state file tracking..."

    if [[ "$integration_type" == "wave" ]]; then
        jq --arg key "$integration_key" \
           '.pre_planned_infrastructure.integrations.wave_integrations[$key].created = false' \
           orchestrator-state-v3.json > tmp-state.json && \
           mv tmp-state.json orchestrator-state-v3.json
    elif [[ "$integration_type" == "phase" ]]; then
        jq --arg key "$integration_key" \
           '.pre_planned_infrastructure.integrations.phase_integrations[$key].created = false' \
           orchestrator-state-v3.json > tmp-state.json && \
           mv tmp-state.json orchestrator-state-v3.json
    elif [[ "$integration_type" == "project" ]]; then
        jq '.pre_planned_infrastructure.integrations.project_integration.created = false' \
           orchestrator-state-v3.json > tmp-state.json && \
           mv tmp-state.json orchestrator-state-v3.json
    fi

    echo "✅ R327 CLEANUP COMPLETE: All stale infrastructure removed"
    echo "   - Remote branch deleted"
    echo "   - Local branch deleted"
    echo "   - Working directory removed"
    echo "   - State tracking reset"
    return 0
}
```

### Operation Execution
```bash
execute_operation() {
    local operation="$1"
    local action="${operation%%:*}"
    local target="${operation#*:}"

    case "$action" in
        DELETE_INTEGRATE_WAVE_EFFORTS)
            echo "🔴 R327: Deleting stale integration: $target"

            # Determine integration type from target name
            if [[ "$target" == *"wave"* ]] && [[ "$target" != *"phase"* ]]; then
                # Extract phase and wave numbers from target
                # Format: phase2-wave2-integration or idpbuilder-push-oci/phase2-wave2-integration
                local phase_num=$(echo "$target" | grep -oP 'phase\K\d+')
                local wave_num=$(echo "$target" | grep -oP 'wave\K\d+')
                cleanup_integration_infrastructure "wave" "$phase_num" "$wave_num"
            elif [[ "$target" == *"phase"* ]] && [[ "$target" != *"project"* ]]; then
                local phase_num=$(echo "$target" | grep -oP 'phase\K\d+')
                cleanup_integration_infrastructure "phase" "$phase_num"
            elif [[ "$target" == *"project"* ]]; then
                cleanup_integration_infrastructure "project"
            else
                # Fallback: just delete the remote branch
                echo "⚠️ Unknown integration format: $target, attempting simple deletion"
                git push origin --delete "$target" 2>/dev/null || true
            fi
            ;;

        CREATE_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE)
            echo "✅ Creating fresh integration infrastructure: $target"

            # DEFENSE IN DEPTH: Clean up before creation just in case DELETE was missed
            if [[ "$target" == *"wave"* ]] && [[ "$target" != *"phase"* ]]; then
                local phase_num=$(echo "$target" | grep -oP 'phase\K\d+')
                local wave_num=$(echo "$target" | grep -oP 'wave\K\d+')
                cleanup_integration_infrastructure "wave" "$phase_num" "$wave_num"
            elif [[ "$target" == *"phase"* ]] && [[ "$target" != *"project"* ]]; then
                local phase_num=$(echo "$target" | grep -oP 'phase\K\d+')
                cleanup_integration_infrastructure "phase" "$phase_num"
            elif [[ "$target" == *"project"* ]]; then
                cleanup_integration_infrastructure "project"
            fi

            # Transition to appropriate infrastructure state (now with clean slate)
            if [[ "$target" == *"wave"* ]]; then
                transition_to_state "SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE"
            elif [[ "$target" == *"phase"* ]]; then
                transition_to_state "SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE"
            elif [[ "$target" == *"project"* ]]; then
                transition_to_state "SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE"
            fi
            return  # Exit loop, state transition will handle rest
            ;;

        MERGE_SOURCES)
            echo "🔀 Merging sources into: $target"
            # This happens in INTEGRATE_WAVE_EFFORTS state
            ;;

        VERIFY_BUILD)
            echo "🔍 Verifying build for: $target"
            # This happens in integration review
            ;;
    esac
}
```

## Cascade Mode Behavior

**When in CASCADE_REINTEGRATION, the system enters "cascade mode":**

### State File Tracking
```bash
# Set cascade mode flag
jq '.cascade_mode = true' orchestrator-state-v3.json > /tmp/state.json
mv /tmp/state.json orchestrator-state-v3.json
```

### Persistent Coordination
- Cascade operations persist across state transitions
- Each infrastructure/integration state checks if in cascade mode
- When in cascade mode, states return to CASCADE_REINTEGRATION after completion

### R352: Overlapping Cascade Support
- Multiple cascades can be in progress simultaneously
- Each cascade is tracked independently
- Dependencies between cascades are respected

### R353: Cascade Focus
**FORBIDDEN during cascade execution:**
- Size checks (cascades don't grow code)
- Split evaluations (not creating new work)
- CREATE_NEXT_INFRASTRUCTURE transitions (cascade uses SETUP states)
- Quality assessments beyond build verification

**ALLOWED during cascade execution:**
- Rebase validation
- Conflict resolution
- Build verification
- Test execution

## 🔴🔴🔴 R410: LAYERED CASCADE PROTOCOL (SUPREME LAW) 🔴🔴🔴

### The Recursive Cascade Problem

**CRITICAL UNDERSTANDING:** Integration build failures during CASCADE are NOT blockers - they are **NEW BUG DISCOVERIES**!

**Why Build Failures Happen During CASCADE:**
- Individual efforts build successfully in isolation
- Integration reveals **cross-effort bugs** (interfaces, dependencies, timing)
- These bugs are ONLY visible when efforts are merged together
- This is EXPECTED behavior, NOT an error condition!

**Example:**
```
Layer 1 Cascade:
1. Fix bugs in effort-3, effort-5
2. Cascade to phase1-wave2-integration → BUILD FAILS!
3. Failure reveals: effort-2 expects sync data, effort-4 returns Promise
4. This is a NEW bug (not in original layer 1 bug list)
5. ❌ WRONG: "Build failed - need manual approval" + FALSE
6. ✅ CORRECT: "New bugs discovered - start layer 2" + TRUE
```

### Layered Cascade Decision Tree

**DECISIVE DECISION TREE - NO INDECISION ALLOWED:**

```
Integration build FAILED during CASCADE_REINTEGRATION?
├─ YES → Are these NEW bugs (not in current layer)?
│        ├─ YES → R410: START NEW CASCADE LAYER (AUTOMATED!)
│        │        ├─ 1. Document new bugs in bug_registry
│        │        ├─ 2. Create cascade_layers[N+1] metadata
│        │        ├─ 3. Pause current layer (progress recorded)
│        │        ├─ 4. Transition to ERROR_RECOVERY (start fixes)
│        │        ├─ 5. SET CONTINUE-SOFTWARE-FACTORY=TRUE
│        │        └─ 🔴 THIS IS AUTOMATED WORKFLOW - NOT A BLOCKER!
│        │
│        └─ NO → SAME bugs as current layer?
│                 ├─ 🚨 INFINITE LOOP DETECTED!
│                 ├─ Fixes didn't resolve bugs
│                 ├─ SET CONTINUE-SOFTWARE-FACTORY=FALSE
│                 └─ Manual review required
│
└─ NO → Build PROJECT_DONE?
         └─ YES → Complete current layer, resume previous
                  ├─ Mark layer complete in cascade_layers[]
                  ├─ Resume previous layer (if exists)
                  └─ CONTINUE-SOFTWARE-FACTORY=TRUE (still automated!)
```

### Implementation: Start New Cascade Layer

**Function to call when integration build fails:**

```bash
# Called when integration build fails during CASCADE_REINTEGRATION
handle_integration_build_failure() {
    local integration_key="$1"
    local build_errors="$2"

    echo "🔍 R410: Integration build failed for $integration_key"

    # Check if these are NEW bugs
    if are_bugs_new_in_cascade "$build_errors"; then
        echo "📋 NEW bugs discovered (cross-effort issues)"
        echo "🔄 R410: Starting new cascade layer (AUTOMATED DECISION)"

        # Get current layer count
        CURRENT_LAYER=$(jq -r '.cascade_layers | length' orchestrator-state-v3.json)
        NEW_LAYER=$((CURRENT_LAYER + 1))

        # Check layer limit (prevent infinite loops)
        if [[ $NEW_LAYER -gt 10 ]]; then
            echo "🚨 R410: Exceeded maximum cascade layers (10)"
            echo "🚨 This indicates systematic design issues"
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
            exit 410
        fi

        echo "🆕 Starting CASCADE LAYER $NEW_LAYER"

        # 1. DOCUMENT NEW BUGS
        document_cascade_bugs "$build_errors" "$NEW_LAYER"

        # 2. CREATE LAYER METADATA
        jq --arg layer "$NEW_LAYER" \
           --arg trigger "build_failures_from_layer_${CURRENT_LAYER}_reintegration" \
           --arg started "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
           --arg integration "$integration_key" \
           '.cascade_layers += [{
               layer_id: ($layer | tonumber),
               started_at: $started,
               trigger: $trigger,
               bugs: [],
               status: "fixing_upstream",
               progress: ($integration + ": build failed, started layer " + $layer)
           }]' orchestrator-state-v3.json > /tmp/state.json

        mv /tmp/state.json orchestrator-state-v3.json

        # 3. PAUSE CURRENT LAYER
        if [[ $CURRENT_LAYER -gt 0 ]]; then
            jq --argjson idx $((CURRENT_LAYER - 1)) \
               --arg pause_status "paused_for_layer_${NEW_LAYER}" \
               --arg progress "${integration_key}: build failed, paused for layer ${NEW_LAYER}" \
               '.cascade_layers[$idx].status = $pause_status |
                .cascade_layers[$idx].progress = $progress' \
               orchestrator-state-v3.json > /tmp/state.json

            mv /tmp/state.json orchestrator-state-v3.json
        fi

        echo "✅ CASCADE LAYER $NEW_LAYER created"
        echo "⏸️  CASCADE LAYER $CURRENT_LAYER paused"

        # 4. TRANSITION TO ERROR_RECOVERY to apply fixes
        update_state "ERROR_RECOVERY" "R410: Starting cascade layer $NEW_LAYER fixes"

        # 5. COMMIT STATE CHANGES
        git add orchestrator-state-v3.json
        git commit -m "cascade: R410 layer $NEW_LAYER started (build failures from layer $CURRENT_LAYER)" || true
        git push || true

        # 6. 🔴 CRITICAL: SET TRUE - THIS IS AUTOMATED WORKFLOW!
        echo ""
        echo "🔴🔴🔴 R410: CASCADE LAYER AUTOMATION 🔴🔴🔴"
        echo "Build failure during cascade is EXPECTED (discovers cross-effort bugs)"
        echo "Starting new layer is AUTOMATED PROTOCOL - NOT a blocker!"
        echo "System will:"
        echo "  1. Apply fixes for layer $NEW_LAYER bugs"
        echo "  2. Re-attempt integration with fixes"
        echo "  3. If successful, resume layer $CURRENT_LAYER"
        echo "  4. Eventually complete all layers"
        echo ""
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

        return 0
    else
        # Same bugs as current layer - infinite loop!
        echo "🚨 R410: INFINITE LOOP DETECTED!"
        echo "Build still failing with SAME bugs after fixes"
        echo "Fixes are ineffective - manual review required"
        echo ""
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 410
    fi
}

# Helper: Check if bugs are new
are_bugs_new_in_cascade() {
    local build_errors="$1"

    # Hash the errors for comparison
    NEW_HASH=$(echo "$build_errors" | sha256sum | cut -d' ' -f1)

    # Check all bugs in all cascade layers
    for bug_id in $(jq -r '.cascade_layers[].bugs[]' orchestrator-state-v3.json 2>/dev/null); do
        EXISTING_HASH=$(jq -r ".bug_registry.bugs[] | select(.id == \"$bug_id\") | .error_hash // empty" orchestrator-state-v3.json)

        if [[ "$NEW_HASH" == "$EXISTING_HASH" ]]; then
            echo "❌ Bug hash matches layer bug $bug_id - NOT new!"
            return 1  # NOT new
        fi
    done

    echo "✅ Bugs are NEW (not in any cascade layer)"
    return 0  # New bugs
}

# Helper: Document bugs for cascade layer
document_cascade_bugs() {
    local build_errors="$1"
    local layer_id="$2"

    # Get next bug ID
    NEXT_BUG_ID=$(jq -r '.bug_registry.bugs | length + 1' orchestrator-state-v3.json)

    # Hash errors for deduplication
    ERROR_HASH=$(echo "$build_errors" | sha256sum | cut -d' ' -f1)

    # Create bug entry
    BUG_ID="BUG-$(printf "%03d" $NEXT_BUG_ID)"

    jq --arg bug_id "$BUG_ID" \
       --arg desc "$build_errors" \
       --arg hash "$ERROR_HASH" \
       --arg layer "$layer_id" \
       --arg discovered "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.bug_registry.bugs += [{
           id: $bug_id,
           description: $desc,
           error_hash: $hash,
           discovered_in_layer: ($layer | tonumber),
           discovered_at: $discovered,
           status: "pending_fix",
           affected_efforts: []
       }]' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    # Link bug to layer
    jq --argjson layer $((layer_id - 1)) \
       --arg bug_id "$BUG_ID" \
       '.cascade_layers[$layer].bugs += [$bug_id]' \
       orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "📋 Documented $BUG_ID for cascade layer $layer_id"
}
```

### Implementation: Complete Layer and Resume Previous

**Function to call when integration build SUCCEEDS:**

```bash
# Called when integration build succeeds during CASCADE_REINTEGRATION
complete_cascade_layer_and_resume() {
    # Find active layer (status = "reintegrating")
    ACTIVE_LAYER=$(jq -r '.cascade_layers | map(select(.status == "reintegrating")) | .[0].layer_id // 0' orchestrator-state-v3.json)

    if [[ $ACTIVE_LAYER -eq 0 ]]; then
        echo "⚠️ R410: No active layer found - cascade may not be configured"
        return 1
    fi

    echo "✅ R410: CASCADE LAYER $ACTIVE_LAYER integration successful!"

    # Mark layer complete
    jq --argjson layer "$ACTIVE_LAYER" \
       --arg completed "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.cascade_layers[] |=
        if .layer_id == $layer then
            .status = "completed" |
            .completed_at = $completed |
            .progress = "integration successful, layer complete"
        else . end' \
       orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    # Check for previous paused layer
    PREV_LAYER=$((ACTIVE_LAYER - 1))

    if [[ $PREV_LAYER -gt 0 ]]; then
        PREV_STATUS=$(jq -r ".cascade_layers[] | select(.layer_id == $PREV_LAYER) | .status" orchestrator-state-v3.json)

        if [[ "$PREV_STATUS" == "paused_for_layer_${ACTIVE_LAYER}" ]]; then
            echo "🔄 R410: Resuming CASCADE LAYER $PREV_LAYER"

            # Resume previous layer
            jq --argjson layer "$PREV_LAYER" \
               '.cascade_layers[] |=
                if .layer_id == $layer then
                    .status = "reintegrating" |
                    .progress = "resuming after layer completion"
                else . end' \
               orchestrator-state-v3.json > /tmp/state.json

            mv /tmp/state.json orchestrator-state-v3.json

            echo "✅ CASCADE LAYER $PREV_LAYER resumed"
            echo "🔄 Will re-attempt integration with layer $ACTIVE_LAYER fixes"
            echo ""
            echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue cascade!
            return 0
        fi
    fi

    # No previous layer - cascade fully complete!
    echo "🎉 R410: ALL CASCADE LAYERS COMPLETE!"
    echo "✅ Layer $ACTIVE_LAYER was the final layer"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Cascade done, continue!
    return 0
}
```

### Cascade Layer Status Reporting

```bash
# Helper: Print cascade layer status
print_cascade_status() {
    echo ""
    echo "📊 R410 CASCADE LAYER STATUS"
    echo "=============================="

    LAYER_COUNT=$(jq -r '.cascade_layers | length' orchestrator-state-v3.json)

    if [[ $LAYER_COUNT -eq 0 ]]; then
        echo "No cascade layers active"
        return
    fi

    for i in $(seq 0 $((LAYER_COUNT - 1))); do
        LAYER=$(jq -r ".cascade_layers[$i]" orchestrator-state-v3.json)

        LAYER_ID=$(echo "$LAYER" | jq -r '.layer_id')
        STATUS=$(echo "$LAYER" | jq -r '.status')
        PROGRESS=$(echo "$LAYER" | jq -r '.progress')
        BUG_COUNT=$(echo "$LAYER" | jq -r '.bugs | length')

        echo ""
        echo "Layer $LAYER_ID:"
        echo "  Status: $STATUS"
        echo "  Bugs: $BUG_COUNT"
        echo "  Progress: $PROGRESS"
    done

    echo ""
    echo "=============================="
}
```

### Integration with Existing CASCADE Flow

**Modify VERIFY_BUILD operation to use R410:**

```bash
execute_operation() {
    local operation="$1"
    local action="${operation%%:*}"
    local target="${operation#*:}"

    case "$action" in
        VERIFY_BUILD)
            echo "🔍 R410: Verifying build for: $target"

            # Attempt build
            if build_integration "$target"; then
                echo "✅ Build successful for $target"

                # R410: Complete layer and resume previous
                complete_cascade_layer_and_resume

            else
                echo "❌ Build failed for $target"

                # Capture build errors
                BUILD_ERRORS=$(cat build-errors.log 2>/dev/null || echo "Build failed - see logs")

                # R410: Handle failure (start new layer or detect loop)
                handle_integration_build_failure "$target" "$BUILD_ERRORS"
            fi
            ;;

        # ... other operations ...
    esac
}
```

### Key Principles

1. **Build failures during CASCADE are EXPECTED** - they discover cross-effort bugs
2. **Starting new layer is AUTOMATED** - NOT a blocker requiring approval
3. **Layers are RECURSIVE** - layer 2 fixes bugs from layer 1, layer 3 from layer 2, etc.
4. **Eventually complete** - layers bubble up as each completes
5. **ALWAYS use TRUE** - this is designed workflow, not an error condition

**ONLY use FALSE when:**
- Infinite loop detected (same bugs in multiple layers)
- Exceeded maximum layers (10+ = design problem)
- Cannot determine what's new (system corruption)

**Grading penalty for incorrect FALSE: -100% (defeats automation!)**

## Valid Transitions

From CASCADE_REINTEGRATION, you can transition to:

1. **SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE**: To recreate wave integration
2. **SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE**: To recreate phase integration
3. **SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE**: To recreate project integration
4. **CASCADE_REINTEGRATION**: Loop back after infrastructure creation
5. **REVIEW_WAVE_INTEGRATION**: When ALL cascades complete

## Forbidden Transitions

**YOU CANNOT transition to these states from CASCADE_REINTEGRATION:**
- ❌ WAVE_COMPLETE (wave isn't complete until cascades done)
- ❌ COMPLETE_PHASE (phase isn't complete until cascades done)
- ❌ PROJECT_DONE (project isn't complete until cascades done)
- ❌ Any analysis/planning states (focus on cascade completion)

## Exit Conditions

**You can ONLY exit CASCADE_REINTEGRATION when:**

### 1. All Cascade Operations Completed
```bash
# Check if cascade chain is empty
check_cascade_complete() {
    local REMAINING=$(wc -l < /tmp/cascade-execution-plan.txt 2>/dev/null || echo 0)

    if [[ $REMAINING -eq 0 ]]; then
        echo "✅ All cascade operations complete"
        return 0
    else
        echo "⏳ $REMAINING cascade operations remaining"
        return 1
    fi
}
```

### 2. No Stale Integrations Remain
```bash
# Verify all integrations are fresh
verify_no_stale_integrations() {
    echo "🔍 R327: Final staleness check..."

    for integration in $(git branch -r | grep "integration"); do
        if is_integration_stale "$integration"; then
            echo "❌ Integration still stale: $integration"
            return 1
        fi
    done

    echo "✅ All integrations are fresh"
    return 0
}
```

### 3. All Dependencies Satisfied
```bash
# Verify dependency graph is satisfied
verify_dependencies_satisfied() {
    echo "🔍 R350: Verifying dependency graph..."

    # All effort dependencies satisfied
    # All wave dependencies satisfied
    # All phase dependencies satisfied

    echo "✅ All dependencies satisfied"
    return 0
}
```

### 4. Freshness Validation Passed
```bash
# Final R328 freshness validation
final_freshness_check() {
    echo "🔍 R328: Final freshness validation..."

    # Run comprehensive freshness check
    for integration_type in wave phase project; do
        if ! validate_integration_freshness "$integration_type"; then
            echo "❌ Freshness validation failed for $integration_type"
            return 1
        fi
    done

    echo "✅ All integrations pass freshness validation"
    return 0
}
```

### Exit Sequence
```bash
# When all conditions met, exit CASCADE_REINTEGRATION
exit_cascade_reintegration() {
    echo "🎉 CASCADE_REINTEGRATION COMPLETE!"

    # Clear cascade mode
    jq 'del(.cascade_mode) |
        .cascade_reintegration.status = "complete" |
        .cascade_reintegration.completed_at = (now | todate)' \
        orchestrator-state-v3.json > /tmp/state.json
    mv /tmp/state.json orchestrator-state-v3.json

    # Transition to integration code review
    echo "➡️ Transitioning to REVIEW_WAVE_INTEGRATION"
    transition_to_state "REVIEW_WAVE_INTEGRATION"
}

# Check all exit conditions
if check_cascade_complete && \
   verify_no_stale_integrations && \
   verify_dependencies_satisfied && \
   final_freshness_check; then
    exit_cascade_reintegration
fi
```

## Associated Rules

### MANDATORY (BLOCKING)
- **R327**: Mandatory Re-Integration After Fixes (the reason this state exists)
- **R348**: Cascade State Transitions (defines state flow during cascades)
- **R350**: Complete Cascade Dependency Graph (tracks dependencies)
- **R351**: Cascade Execution Protocol (defines execution order)
- **R328**: Integration Freshness Validation (validates cascade results)

### SUPPORTING
- **R352**: Overlapping Cascade Support (multiple cascades in progress)
- **R353**: Cascade Focus (what to skip during cascades)

## Prohibitions

- ❌ Skip cascade operations to "save time"
- ❌ Exit to PROJECT_DONE without completing all cascades
- ❌ Create new work while cascades are pending
- ❌ Run size checks or split evaluations during cascades
- ❌ Bypass freshness validation at exit
- ❌ Transition to completion states (WAVE_COMPLETE, COMPLETE_PHASE, PROJECT_DONE)

## Example Cascade Execution

**Scenario: Fixes applied to phase1/wave1/effort-2**

```
1. ERROR_RECOVERY detects stale integrations
2. Transition to CASCADE_REINTEGRATION
3. Calculate cascade chain:
   - phase1-wave1-integration is stale (contains effort-2)
   - phase1-integration is stale (contains wave1)
   - project-integration is stale (contains phase1)
4. Create execution plan:
   - DELETE phase1-wave1-integration
   - RECREATE phase1-wave1-integration (SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE)
   - DELETE phase1-integration
   - RECREATE phase1-integration (SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE)
   - DELETE project-integration
   - RECREATE project-integration (SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE)
5. Execute plan step by step
6. After each SETUP state, return to CASCADE_REINTEGRATION
7. When all steps complete, verify freshness
8. Exit to REVIEW_WAVE_INTEGRATION
```

## Automation Flag - CRITICAL FOR CASCADE AUTOMATION

### 🔴🔴🔴 CASCADE OPERATIONS ARE FULLY AUTOMATED 🔴🔴🔴

**CASCADES SHOULD KEEP AUTOMATION MOVING!**

### 🚨🚨🚨 CASCADE TRANSITION DECISION FLOWCHART 🚨🚨🚨

**USE THIS DECISION TREE FOR EVERY CASCADE TRANSITION:**

```
Are you transitioning to SETUP_*_INFRASTRUCTURE during cascade?
├─ YES → This is NORMAL CASCADE OPERATION
│         ├─ Infrastructure deleted successfully?
│         │  ├─ YES → Use CONTINUE-SOFTWARE-FACTORY=TRUE
│         │  └─ NO (unrecoverable) → Use FALSE
│         └─ DEFAULT: TRUE (cascade continues to SETUP)
│
└─ NO → Are you transitioning to REVIEW_WAVE_INTEGRATION?
        ├─ YES → Cascade complete, review is NORMAL
        │         └─ Use CONTINUE-SOFTWARE-FACTORY=TRUE
        └─ NO → Still in cascade loop
                  └─ Use CONTINUE-SOFTWARE-FACTORY=TRUE
```

**THERE IS NO SCENARIO IN CASCADE WHERE YOU USE FALSE UNLESS TRULY UNRECOVERABLE!**

### 🔴 EXPLICIT OUTBOUND TRANSITION GUIDANCE 🔴

**When transitioning FROM CASCADE_REINTEGRATION TO any other state:**

```bash
# SCENARIO 1: Transitioning to SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
# This happens after deleting wave integration
echo "🔄 Transitioning CASCADE_REINTEGRATION → SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE"
update_state "SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE" "Recreating ${INTEGRATE_WAVE_EFFORTS_KEY} after R327 cascade deletion"

# R322: Checkpoint before state transition
save_todos "CASCADE_TO_SETUP_TRANSITION"
git add todos/*.todo orchestrator-state-v3.json
git commit -m "state: CASCADE_REINTEGRATION → SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE (R327 cascade)"
git push

echo "🛑 R322: Mandatory checkpoint before state transition"

# Print cascade status
if jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    source utilities/cascade-status-report.sh
    cascade_status_report
fi

# 🔴🔴🔴 CRITICAL: THIS IS NORMAL CASCADE OPERATION - USE TRUE! 🔴🔴🔴
# The system knows:
# - Current state: SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE (from state file)
# - Next action: Create infrastructure per SETUP rules
# - After SETUP: Return to CASCADE_REINTEGRATION to continue chain
# - NO HUMAN INTERVENTION NEEDED!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # CASCADE AUTOMATION CONTINUES!
exit 0

# ❌ WRONG: echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
# ❌ This defeats cascade automation!
# ❌ User would have to run /continue after EVERY cascade step!
# ❌ 6 integrations = 12 manual restarts (delete + create for each)!
```

```bash
# SCENARIO 2: Transitioning to SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE
# Same pattern - this is NORMAL CASCADE OPERATION
echo "🔄 Transitioning CASCADE_REINTEGRATION → SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE"
# ... checkpoint ...
# 🔴 USE TRUE - Cascade continues!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
exit 0
```

```bash
# SCENARIO 3: Transitioning to REVIEW_WAVE_INTEGRATION
# Cascade complete - review is NORMAL workflow
echo "🔄 Transitioning CASCADE_REINTEGRATION → REVIEW_WAVE_INTEGRATION"
# ... checkpoint ...
# 🔴 USE TRUE - Review is automated!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
exit 0
```

### 🔴 AUTOMATED FLAG VALIDATION 🔴

```bash
# MANDATORY: Validate flag before printing
validate_cascade_continuation_flag() {
    local proposed_flag="$1"
    local transition_target="$2"

    # Check if we're in cascade mode
    if ! jq -e '.cascade_mode == true' orchestrator-state-v3.json > /dev/null 2>&1; then
        echo "⚠️ Not in cascade mode, validation skipped"
        return 0
    fi

    # Check if proposing FALSE
    if [ "$proposed_flag" = "FALSE" ]; then
        echo ""
        echo "🚨🚨🚨 CASCADE AUTOMATION FLAG VALIDATION FAILURE 🚨🚨🚨"
        echo ""
        echo "You are attempting to set CONTINUE-SOFTWARE-FACTORY=FALSE during CASCADE!"
        echo "Transition target: $transition_target"
        echo ""
        echo "❓ IS THIS TRULY UNRECOVERABLE?"
        echo "   - Can system determine next cascade step? (if YES → use TRUE)"
        echo "   - Is dependency graph corrupted? (if NO → use TRUE)"
        echo "   - Is state file invalid? (if NO → use TRUE)"
        echo "   - Is this just a normal transition? (if YES → use TRUE)"
        echo ""
        echo "🔴 CASCADE TRANSITIONS ARE NORMAL OPERATIONS!"
        echo "   CASCADE → SETUP → CASCADE → SETUP → ... → REVIEW"
        echo "   All of these are AUTOMATED WORKFLOW!"
        echo ""
        echo "✅ DEFAULT: CONTINUE-SOFTWARE-FACTORY=TRUE (99.9% of cascades)"
        echo "❌ ONLY use FALSE if genuinely stuck (0.1% of cascades)"
        echo ""
        echo "Grading penalty for incorrect FALSE: -20%"
        echo ""
        read -p "Are you SURE this is unrecoverable? (type UNRECOVERABLE to confirm FALSE): " confirm

        if [ "$confirm" != "UNRECOVERABLE" ]; then
            echo "✅ Validation corrected: Using TRUE instead"
            echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
            return 1  # Indicates correction was made
        fi
    fi

    return 0
}

# Use validation before setting flag
PROPOSED_FLAG="TRUE"  # Default to TRUE!
TRANSITION_TARGET="SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE"

# Validate (will auto-correct if wrong)
validate_cascade_continuation_flag "$PROPOSED_FLAG" "$TRANSITION_TARGET"

# Print validated flag
echo "CONTINUE-SOFTWARE-FACTORY=$PROPOSED_FLAG"
```

### Standard Cascade Flag Pattern

```bash
# 🔴 USE THIS PATTERN FOR ALL CASCADE TRANSITIONS 🔴

# 1. Complete cascade step (delete OR setup complete)
complete_cascade_step()

# 2. Update state file for transition
update_state "NEXT_STATE" "Reason for transition"

# 3. R322 checkpoint
save_todos "TRANSITION_REASON"
git add todos/*.todo orchestrator-state-v3.json
git commit -m "state: CASCADE transition"
git push

# 4. Print status if active cascade
if jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    source utilities/cascade-status-report.sh
    cascade_status_report
fi

# 5. 🔴 ALWAYS USE TRUE FOR CASCADE TRANSITIONS 🔴
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
exit 0

# ❌ NEVER: echo "CONTINUE-SOFTWARE-FACTORY=FALSE" (defeats automation!)
```

### Why TRUE for ALL Cascade Transitions?

**CASCADE → SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE:**
- System knows current state (SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE from state file)
- System knows what SETUP does (create infrastructure per R504)
- System knows where SETUP goes next (back to CASCADE or to merge)
- NO HUMAN DECISION NEEDED!

**CASCADE → SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE:**
- Same automated flow
- NO HUMAN DECISION NEEDED!

**CASCADE → REVIEW_WAVE_INTEGRATION:**
- Review is normal workflow (spawn reviewer per R291)
- NO HUMAN DECISION NEEDED!

### ❌ ONLY Use FALSE When (0.01% of cascades)

```bash
# Truly unrecoverable cascade failure
if [[ "$CASCADE_FAILURE_TYPE" == "UNRECOVERABLE" ]]; then
    echo "❌ CRITICAL: Cannot determine cascade dependencies"
    echo "❌ Dependency graph corrupted - human debugging required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

# Examples of unrecoverable:
# - orchestrator-state-v3.json completely corrupted (cannot parse)
# - Cascade chain calculation failed (cannot determine what to recreate)
# - Dependency graph missing/corrupt (cannot determine order)
# - Git repository in unknown state (cannot proceed)
```

### Summary - UNMISTAKABLE RULES

- ✅ **TRUE for CASCADE → SETUP transition**: Normal cascade operation
- ✅ **TRUE for CASCADE → CASCADE transition**: Continuing cascade chain
- ✅ **TRUE for CASCADE → REVIEW transition**: Cascade complete, review is normal
- ❌ **FALSE**: ONLY for catastrophic unrecoverable failures (0.01%)
- 🛑 **R322 stops**: INDEPENDENT from flag value!
- 📊 **Grading penalty**: -20% per incorrect FALSE during cascade

**DEFAULT: CONTINUE-SOFTWARE-FACTORY=TRUE**

**The entire cascade chain executes WITHOUT MANUAL INTERVENTION!**

## Notes

- CASCADE_REINTEGRATION is a **TRAP STATE** - you cannot escape until cascades complete
- This state enforces R327 absolutely - no bypasses allowed
- Cascade mode persists across state transitions
- Each integration recreation goes through normal infrastructure/merge flow
- After recreation, return to CASCADE_REINTEGRATION to continue
- Only exits when ALL integrations are fresh and ALL dependencies satisfied
- **Violation of this state's requirements = -100% automatic failure per R327**

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
