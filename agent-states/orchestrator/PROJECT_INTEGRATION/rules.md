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

### 🚨🚨🚨 R104 - Target Repository Integration [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R104-target-repository-integration.md`
**Criticality**: BLOCKING - Must use target repository
**Summary**: Integration happens in target repository, NOT software-factory

### 🚨🚨🚨 R321 - Immediate Backport During Integration Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: BLOCKING - Fixes found during integration must be backported immediately
**Summary**: When integration issues are found, must stop and backport fixes to original branches

### 🚨🚨🚨 R280 - Main Branch Protection Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection-protocol.md`
**Criticality**: BLOCKING - Direct commits to main/master are forbidden
**Summary**: All changes must go through PR process with proper reviews

### 🚨🚨🚨 R307 - Branch Mergeability Check [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R307-branch-mergeability-check.md`
**Criticality**: BLOCKING - Must verify branches are mergeable before attempting
**Summary**: Check for conflicts and mergeability before integration operations

## 🚨 PROJECT_INTEGRATION IS A VERB - SET UP INFRASTRUCTURE NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING PROJECT_INTEGRATION

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Create project integration workspace per R283
2. Set up isolated integration environment
3. Prepare for Code Reviewer merge plan creation
4. Document integration infrastructure
5. Transition to SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in PROJECT_INTEGRATION" [stops]
- ❌ "Successfully entered PROJECT_INTEGRATION state" [waits]
- ❌ "Ready to set up project integration" [pauses]
- ❌ "I'm in PROJECT_INTEGRATION state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering PROJECT_INTEGRATION, creating isolated workspace NOW..."
- ✅ "START PROJECT INTEGRATION SETUP per R283..."
- ✅ "PROJECT_INTEGRATION: Setting up integration infrastructure..."

## State Context

**Prerequisites:**
- PHASE_COMPLETE has been reached (last phase completed)
- All phases have been individually integrated
- Architect has approved the final phase
- Ready for project-wide integration

**Purpose:**
PROJECT_INTEGRATION creates the infrastructure for merging ALL phases:
1. Create isolated project integration workspace (R283)
2. Clone target repository (NOT software-factory)
3. Create project-integration branch
4. Prepare for merging all phase integration branches
5. Document project integration plan

## 🔴🔴🔴 CRITICAL: R283 COMPLIANCE REQUIRED 🔴🔴🔴

**Project Integration Infrastructure per R283:**

```bash
create_project_integration_infrastructure() {
    echo "🏭 PROJECT_INTEGRATION: Starting infrastructure setup per R283..."
    
    # Save SF instance directory
    SF_INSTANCE_DIR=$(pwd)
    echo "📁 SF Instance: $SF_INSTANCE_DIR"
    
    # Read target repository configuration (R104)
    TARGET_CONFIG="$SF_INSTANCE_DIR/target-repo-config.yaml"
    TARGET_REPO_URL=$(yq '.repository_path' "$TARGET_CONFIG")
    TARGET_REPO_NAME=$(yq '.repository_name' "$TARGET_CONFIG")
    DEFAULT_BRANCH=$(yq '.default_branch' "$TARGET_CONFIG")
    
    # Create isolated project integration workspace (R283)
    PROJECT_INTEGRATION_DIR="$SF_INSTANCE_DIR/efforts/project/integration-workspace"
    rm -rf "$PROJECT_INTEGRATION_DIR"
    mkdir -p "$PROJECT_INTEGRATION_DIR"
    cd "$PROJECT_INTEGRATION_DIR"
    
    # Clone TARGET repository (NOT software-factory!) per R104
    echo "🔄 Cloning target repository (R104)..."
    git clone "$TARGET_REPO_URL" "$TARGET_REPO_NAME"
    cd "$TARGET_REPO_NAME"
    
    # CRITICAL SAFETY CHECK - Verify correct repository (R104)
    REMOTE_URL=$(git remote get-url origin)
    if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
        echo "❌ CRITICAL: Cloned orchestrator repository instead of target!"
        echo "Expected: Target project repository"
        echo "Got: $REMOTE_URL"
        exit 104  # R104 violation
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
    
    # Update orchestrator state
    cd "$SF_INSTANCE_DIR"
    jq '.project_integration = {
        "workspace": "'$PROJECT_INTEGRATION_DIR'",
        "branch": "'$INTEGRATION_BRANCH'",
        "status": "INFRASTRUCTURE_READY"
    }' orchestrator-state.json
    
    echo "✅ Project integration infrastructure ready"
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