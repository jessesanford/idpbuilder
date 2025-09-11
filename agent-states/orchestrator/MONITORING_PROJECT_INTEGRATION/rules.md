# Orchestrator - MONITORING_PROJECT_INTEGRATION State Rules

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

## State Context

**Purpose:**
Monitor the Integration Agent as it merges all phase branches into the project integration branch per R283.

## Primary Actions

### 🔴🔴🔴 CRITICAL: KNOW WHERE TO LOOK! 🔴🔴🔴

**Integration reports are in the integration workspace, NOT in SF directory!**

```bash
# Get integration workspace location from state file
locate_integration_workspace() {
    SF_INSTANCE_DIR=$(pwd)
    echo "📁 SF Instance: $SF_INSTANCE_DIR"
    
    # Read workspace location from state
    PROJECT_WORKSPACE=$(jq -r '.project_integration.workspace // .integration_workspaces.project.workspace // "NOT_SET"' orchestrator-state.json)
    MERGE_PLAN_PATH=$(jq -r '.project_integration.merge_plan // .integration_workspaces.project.merge_plan // "NOT_SET"' orchestrator-state.json)
    
    echo "📍 Project integration workspace: $PROJECT_WORKSPACE"
    echo "📄 Merge plan expected at: $MERGE_PLAN_PATH"
    
    if [[ -d "$PROJECT_WORKSPACE" ]]; then
        echo "✅ Integration workspace exists"
        return 0
    else
        echo "❌ Integration workspace not found!"
        echo "📂 Looking for integration workspaces..."
        find efforts -type d -name "integration-workspace" 2>/dev/null
        return 1
    fi
}
```

1. **Monitor Integration Progress**:
   ```bash
   # Navigate to integration workspace to find reports
   if locate_integration_workspace; then
       cd "$PROJECT_WORKSPACE"
       
       # Look for integration report
       if [[ -f "PROJECT-INTEGRATION-REPORT.md" ]]; then
           echo "✅ Found PROJECT-INTEGRATION-REPORT.md"
           cat PROJECT-INTEGRATION-REPORT.md
       elif [[ -f "INTEGRATION-REPORT.md" ]]; then
           echo "✅ Found INTEGRATION-REPORT.md"
           cat INTEGRATION-REPORT.md
       else
           echo "⚠️ No integration report found yet"
           echo "📂 Available files:"
           ls -la *.md
       fi
       
       # Check git log for merge activity
       echo "📊 Recent merge activity:"
       git log --oneline --graph -10
       
       cd "$SF_INSTANCE_DIR"
   fi
   ```
   - Track merge completion for each phase
   - Monitor for conflicts or failures

2. **Check for Documented Bugs (R266)**:
   ```bash
   # In integration workspace, check for bugs
   cd "$PROJECT_WORKSPACE"
   if grep -q "UPSTREAM BUGS IDENTIFIED" *.md 2>/dev/null; then
       echo "🐛 BUGS FOUND! Must fix before proceeding (R266)"
       grep -A 10 "UPSTREAM BUGS IDENTIFIED" *.md
   fi
   cd "$SF_INSTANCE_DIR"
   ```
   - If bugs found, will be detected during code review
   - Transition to PROJECT_INTEGRATION_CODE_REVIEW for comprehensive review

3. **Validate Integration Success**:
   ```bash
   # Check integration status from state file
   INTEGRATION_STATUS=$(jq -r '.project_integration.status' orchestrator-state.json)
   echo "📊 Integration status: $INTEGRATION_STATUS"
   
   # Check if all phases were merged
   PHASES_TO_MERGE=$(jq -r '.project_integration.phase_branches_to_merge[]' orchestrator-state.json)
   echo "📋 Phases that should be merged:"
   echo "$PHASES_TO_MERGE"
   ```
   - All phases merged successfully
   - No unresolved conflicts
   - Build/tests pass in integrated state
   - NO BUGS documented (if bugs exist, must fix first)

4. **Handle Failures per R321**:
   - If integration fails, must fix in source branches
   - Trigger IMMEDIATE_BACKPORT_REQUIRED if needed

## Valid State Transitions

- **SUCCESS** → PROJECT_INTEGRATION_CODE_REVIEW (all phases merged, need code review)
- **FAILURE** → ERROR_RECOVERY (integration failed catastrophically)
- **CONFLICTS** → IMMEDIATE_BACKPORT_REQUIRED (R321 enforcement)

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
