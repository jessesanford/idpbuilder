# Orchestrator - PROJECT_INTEGRATION State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED PROJECT_INTEGRATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_PROJECT_INTEGRATION
echo "$(date +%s) - Rules read and acknowledged for PROJECT_INTEGRATION" > .state_rules_read_orchestrator_PROJECT_INTEGRATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR PROJECT_INTEGRATION STATE

### 🚨🚨🚨 R283 - Project Integration Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R283-project-integration-protocol.md`
**Criticality**: BLOCKING - Mandatory for project completion
**Summary**: Final project integration MUST occur in isolated workspace with all phases merged

### 🚨🚨🚨 R006 - Orchestrator NEVER Writes Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Any code operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator is a COORDINATOR ONLY - never writes, edits, or modifies code

### 🚨🚨🚨 R269 - Merge Plan Requirement [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R269-merge-plan-requirement.md`
**Criticality**: BLOCKING - Must have formal merge plan
**Summary**: Code Reviewer must create detailed merge plan before integration

### 🚨🚨🚨 R270 - Merge Order Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R270-merge-order-protocol.md`
**Criticality**: BLOCKING - Merge order is critical
**Summary**: Phases must be merged in dependency order (Phase 1 → Phase 2 → ...)

### 🚨🚨🚨 R009 - Mandatory Wave/Phase Integration Protocol [SUPREME LAW]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R009-integration-branch-creation.md`
**Criticality**: BLOCKING - Must use target repository
**Summary**: Integration happens in target repository, NOT software-factory

### 🚨🚨🚨 R321 - Immediate Backport During Integration Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: BLOCKING - Fixes found during integration must be backported immediately
**Summary**: When integration issues are found, must stop and backport fixes to original branches

### 🚨🚨🚨 R280 - Main Branch Protection Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection.md`
**Criticality**: BLOCKING - Direct commits to main/master are forbidden
**Summary**: All changes must go through PR process with proper reviews

### 🚨🚨🚨 R307 - Branch Mergeability Check [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R307-independent-branch-mergeability.md`
**Criticality**: BLOCKING - Must verify branches are mergeable before attempting
**Summary**: Check for conflicts and mergeability before integration operations

### 🚨🚨🚨 R328 - Integration Freshness Validation [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R328-integration-freshness-validation.md`
**Criticality**: BLOCKING - Stale integrations cause failed merges and lost fixes
**Summary**: MUST verify all phase branches are fresh before creating project integration

## 🚨 PROJECT_INTEGRATION IS A VERB - COORDINATE PROJECT INTEGRATION NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING PROJECT_INTEGRATION

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Check if project integration infrastructure exists NOW
2. If NO infrastructure: Transition to SETUP_PROJECT_INTEGRATION_INFRASTRUCTURE
3. If infrastructure EXISTS: Transition to SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN
4. Update state file with the appropriate next state
5. Stop per R322 for state transition

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in PROJECT_INTEGRATION" [stops]
- ❌ "Successfully entered PROJECT_INTEGRATION state" [waits]
- ❌ "Ready to set up project integration" [pauses]
- ❌ "I'm in PROJECT_INTEGRATION state" [does nothing]
- ❌ Creating infrastructure yourself (PROJECT_INTEGRATION only coordinates!)
- ❌ Merging branches yourself (R329 violation!)

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "PROJECT_INTEGRATION STATE: Checking for existing project integration infrastructure..."
- ✅ "No infrastructure found, transitioning to SETUP_PROJECT_INTEGRATION_INFRASTRUCTURE..."
- ✅ "Infrastructure exists, transitioning to SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN..."

## State Context

**Prerequisites:**
- PHASE_COMPLETE has been reached (last phase completed)
- All phases have been individually integrated
- Architect has approved the final phase
- Ready for project-wide integration

**Purpose:**
**THIS STATE IS FOR COORDINATION ONLY!**

The PROJECT_INTEGRATION state is a decision point that:
1. **CHECKS** if project integration infrastructure exists
2. **TRANSITIONS** to SETUP_PROJECT_INTEGRATION_INFRASTRUCTURE if no infrastructure
3. **TRANSITIONS** to SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN if infrastructure exists

**THIS STATE NEVER:**
- ❌ Creates project integration workspace itself
- ❌ Sets up branches or directories itself
- ❌ Performs any actual integration work
- ❌ Executes git merges (R329 violation!)

You are the COORDINATOR of project integration flow.

## 🔴🔴🔴 CRITICAL: VERIFY PHASE BRANCH FRESHNESS FIRST! 🔴🔴🔴

**Before creating project integration, MUST verify all phase branches are fresh:**

```bash
verify_phase_branch_freshness() {
    echo "🔍 Checking phase branch freshness..."
    
    # Check if any effort branches have been updated after phase integrations
    STALE_PHASES=""
    
    for phase_num in $(seq 1 $(jq '.phases_planned' orchestrator-state.json)); do
        PHASE_INTEGRATION_TIME=$(jq -r --arg p "$phase_num" '.integration_branches[] | select(.phase == ($p | tonumber)) | .created_at' orchestrator-state.json)
        
        if [[ -z "$PHASE_INTEGRATION_TIME" ]]; then
            echo "⚠️ Phase $phase_num has no integration branch yet"
            continue
        fi
        
        # Check if any efforts in this phase were updated after integration
        EFFORTS_UPDATED=$(jq -r --arg p "$phase_num" --arg time "$PHASE_INTEGRATION_TIME" '
            .efforts_completed[] | 
            select(.phase == ($p | tonumber)) | 
            select((.last_updated_at // .completion_time) > $time) | 
            .name' orchestrator-state.json)
        
        if [[ -n "$EFFORTS_UPDATED" ]]; then
            echo "🔴 Phase $phase_num integration is STALE!"
            echo "   Updated efforts: $EFFORTS_UPDATED"
            STALE_PHASES="$STALE_PHASES $phase_num"
        fi
    done
    
    if [[ -n "$STALE_PHASES" ]]; then
        echo "❌ CRITICAL: Phase integrations are stale: $STALE_PHASES"
        echo "📝 These phases received fixes/updates after integration"
        echo "🔄 Must recreate phase integrations from fresh effort branches"
        return 1
    else
        echo "✅ All phase integrations are fresh"
        return 0
    fi
}
```

## 🔴🔴🔴 CRITICAL: R283 COMPLIANCE REQUIRED 🔴🔴🔴

**Project Integration Infrastructure per R283:**

```bash
create_project_integration_infrastructure() {
    # FIRST: Verify phase branches are fresh
    if ! verify_phase_branch_freshness; then
        echo "❌ Cannot create project integration with stale phase branches"
        echo "📝 Document the staleness issue"
        cat > PHASE-STALENESS-REPORT.md << 'EOF'
# Phase Integration Staleness Detected
Date: $(date)
Issue: Some phase integration branches are outdated
Action Required: Recreate phase integrations from latest effort branches
Stale Phases: $STALE_PHASES
EOF
        return 1
    fi
    echo "🏭 PROJECT_INTEGRATION: Starting infrastructure setup per R283..."
    
    # Save SF instance directory
    SF_INSTANCE_DIR=$(pwd)
    echo "📁 SF Instance: $SF_INSTANCE_DIR"
    
    # Read target repository configuration (R009 requirement)
    TARGET_CONFIG="$SF_INSTANCE_DIR/target-repo-config.yaml"
    TARGET_REPO_URL=$(yq '.repository_path' "$TARGET_CONFIG")
    TARGET_REPO_NAME=$(yq '.repository_name' "$TARGET_CONFIG")
    DEFAULT_BRANCH=$(yq '.default_branch' "$TARGET_CONFIG")
    
    # Create isolated project integration workspace (R283)
    PROJECT_INTEGRATION_DIR="$SF_INSTANCE_DIR/efforts/project/integration-workspace"
    rm -rf "$PROJECT_INTEGRATION_DIR"
    mkdir -p "$PROJECT_INTEGRATION_DIR"
    cd "$PROJECT_INTEGRATION_DIR"
    
    # Clone TARGET repository (NOT software-factory!) per R009
    echo "🔄 Cloning target repository (R009)..."
    git clone "$TARGET_REPO_URL" "$TARGET_REPO_NAME"
    cd "$TARGET_REPO_NAME"
    
    # CRITICAL SAFETY CHECK - Verify correct repository (R009)
    REMOTE_URL=$(git remote get-url origin)
    if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
        echo "❌ CRITICAL: Cloned orchestrator repository instead of target!"
        echo "Expected: Target project repository"
        echo "Got: $REMOTE_URL"
        exit 9  # R009 violation
    fi
    
    # Create project integration branch
    git checkout "$DEFAULT_BRANCH"
    git pull origin "$DEFAULT_BRANCH"
    INTEGRATION_BRANCH="project-integration"
    git checkout -b "$INTEGRATION_BRANCH"
    
    # Document infrastructure
    cat > PROJECT-INTEGRATION-INFO.md << EOF
# Project Integration Infrastructure

## Environment Details
- Created: $(date -u +%Y-%m-%dT%H:%M:%SZ)
- SF Instance: $SF_INSTANCE_DIR
- Target Repository: $TARGET_REPO_NAME
- Integration Branch: $INTEGRATION_BRANCH
- Base Branch: $DEFAULT_BRANCH

## Phase Integration Branches to Merge
$(jq '.phases[].integration_branch' "$SF_INSTANCE_DIR/orchestrator-state.json")

## R283 Compliance
- ✅ Isolated workspace created
- ✅ Target repository cloned (not SF)
- ✅ Project integration branch created
- ✅ Ready for merge plan creation
EOF
    
    # Update orchestrator state with comprehensive metadata
    cd "$SF_INSTANCE_DIR"
    
    # Collect phase integration branch information
    PHASE_BRANCHES=$(jq -r '.integration_branches[] | "\(.branch):\(.created_at)"' orchestrator-state.json)
    
    # Update state file with detailed tracking
    jq --arg workspace "$PROJECT_INTEGRATION_DIR" \
       --arg branch "$INTEGRATION_BRANCH" \
       --arg repo_name "$TARGET_REPO_NAME" \
       --arg created "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.project_integration = {
        "workspace": $workspace,
        "branch": $branch,
        "repository": $repo_name,
        "status": "INFRASTRUCTURE_READY",
        "created_at": $created,
        "phase_branches_to_merge": (.integration_branches | map(.branch)),
        "merge_plan": ($workspace + "/PROJECT-MERGE-PLAN.md"),
        "is_stale": false
    } | .integration_workspaces.project = {
        "workspace": $workspace,
        "merge_plan": ($workspace + "/PROJECT-MERGE-PLAN.md"),
        "created_at": $created
    } | .repository_contexts.current_operation = {
        "type": "PROJECT_INTEGRATION",
        "target_repo_workspace": $workspace,
        "software_factory_dir": "'$SF_INSTANCE_DIR'"
    }' orchestrator-state.json > orchestrator-state.json.tmp && mv orchestrator-state.json.tmp orchestrator-state.json
    
    echo "✅ Project integration infrastructure ready with metadata tracking"
    echo "📊 State file updated with:"
    echo "   - Integration workspace location"
    echo "   - Merge plan path"
    echo "   - Repository context tracking"
    echo "   - Phase branches to merge"
}
```

## Valid State Transitions

From PROJECT_INTEGRATION state:
- **SUCCESS** → SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN (infrastructure ready)
- **FAILURE** → ERROR_RECOVERY (setup failed)

## Success Criteria

Before transitioning to SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN:
1. ✅ Project integration workspace created at correct location
2. ✅ Target repository cloned (verified not SF repo)
3. ✅ Project integration branch created
4. ✅ All phase integration branches identified
5. ✅ Infrastructure documented
6. ✅ State file updated

## Common Pitfalls to Avoid

1. **Wrong Repository**: Cloning software-factory instead of target
2. **Wrong Location**: Not using isolated workspace per R283
3. **Missing Phases**: Not identifying all phase integration branches
4. **Merge Attempts**: Orchestrator trying to merge (violates R006)
5. **Code Writing**: Orchestrator creating merge scripts (violates R006)

## Summary

The PROJECT_INTEGRATION state is responsible for:
1. Creating isolated project integration infrastructure per R283
2. Preparing for Code Reviewer to create merge plan
3. Setting up for Integration Agent to execute merges
4. Ensuring all phases will be integrated together
5. Maintaining separation between SF and target repositories

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```
