# Orchestrator - INTEGRATION_TESTING State Rules

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

**YOU HAVE ENTERED INTEGRATION_TESTING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INTEGRATION_TESTING
echo "$(date +%s) - Rules read and acknowledged for INTEGRATION_TESTING" > .state_rules_read_orchestrator_INTEGRATION_TESTING
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INTEGRATION TESTING WORK UNTIL RULES ARE READ:
- ❌ Start merging efforts
- ❌ Start running tests
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES

### 🚨🚨🚨 RULE R006 - Orchestrator NEVER Writes Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality:** BLOCKING - Any code operation = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

### 🚨🚨🚨 RULE R329 - Orchestrator NEVER Performs Git Merges [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`
**Criticality:** BLOCKING - Any merge operation = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`

**⚠️ R006 + R329 WARNING FOR INTEGRATION_TESTING STATE:**
- DO NOT execute git merge commands yourself! (R329)
- DO NOT resolve merge conflicts yourself! (R006 + R329)
- DO NOT edit code to fix integration issues! (R006)
- DO NOT apply patches or fixes directly! (R006 + R329)
- MUST spawn Integration Agent for ALL merges (R329)
- Document all issues for appropriate agents to resolve
- You only coordinate - NEVER execute merges or modify code

### 🚨🚨🚨 RULE R271 - Mandatory Production Ready Validation [BLOCKING]
**MUST validate production readiness** | Source: rule-library/R271-mandatory-production-ready-validation.md

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`

This rule defines the mandatory production ready validation process required before any code can be considered complete.

### 🚨🚨🚨 RULE R273 - Runtime Specific Validation [BLOCKING]
**MUST validate runtime specific requirements** | Source: rule-library/R273-runtime-specific-validation.md

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R273-runtime-specific-validation.md`

This rule defines runtime-specific validation requirements based on the technology stack being used.

### 🚨🚨🚨 RULE R280 - Main Branch Protection [BLOCKING]
**SOFTWARE FACTORY NEVER MERGES TO MAIN** | Source: rule-library/R280-main-branch-protection.md

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection.md`

Software Factory creates MASTER-PR-PLAN.md for humans to execute PRs. We NEVER push to main ourselves.

### 🚨🚨🚨 RULE R328 - Integration Freshness Validation [BLOCKING]
**MUST verify integration branch freshness before merging** | Source: rule-library/R328-integration-freshness-validation.md

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R328-integration-freshness-validation.md`

Integration branches become stale when fixes are applied after creation. ALWAYS check timestamps!

## 🎯 STATE OBJECTIVES

In the INTEGRATION_TESTING state, you are responsible for:

1. **Coordinating Integration via Integration Agent (R329 MANDATORY)**
   - Identify all effort directories and branches
   - Spawn Code Reviewer to create merge plan
   - **SPAWN INTEGRATION AGENT TO EXECUTE ALL MERGES**
   - Monitor Integration Agent progress via reports
   - NEVER execute merges yourself (R329 VIOLATION)

2. **Verifying Integration Success (via Agent Reports)**
   - Review Integration Agent's INTEGRATION-REPORT.md
   - Check for conflicts documented by Integration Agent
   - Spawn Code Reviewer for build/test validation if needed
   - Document any issues for next states

3. **Creating Orchestration Report**
   - Document which agents were spawned
   - Summarize Integration Agent's findings
   - Note any issues requiring attention
   - Track state transitions

## 📝 REQUIRED ACTIONS

### 🔴🔴🔴 CRITICAL: UNDERSTANDING REPOSITORY CONTEXTS 🔴🔴🔴

**BEFORE YOU START, YOU MUST UNDERSTAND WHERE THINGS ARE:**

1. **Software Factory Repository** (`/home/vscode/software-factory-template/`)
   - Contains: SF code, rules, state files, agent configs
   - Branches: main, software-factory-2.0
   - **NEVER contains effort branches or integration branches**

2. **Target Repository Clones** (`efforts/*/*/`)
   - Contains: Actual project implementation code
   - Branches: effort branches (e.g., `phase1/wave1/effort-name`)
   - Location: Each effort has its own clone of target repo
   - **THIS IS WHERE EFFORT CODE LIVES**

3. **Integration Workspaces** (`efforts/*/integration-workspace/`)
   - Contains: Integration branches and merge operations
   - Branches: integration branches (e.g., `wave1-integration`, `phase1-integration`)
   - Location: Separate clones for merging work
   - **THIS IS WHERE INTEGRATION HAPPENS**

### Step 1: Identify All Efforts and Their Locations
```bash
# CRITICAL: Navigate from SF instance directory
SF_INSTANCE_DIR=$(pwd)
echo "📁 SF Instance: $SF_INSTANCE_DIR"

# Check state file for effort locations
echo "📊 Reading effort locations from state file..."
jq -r '.efforts_completed[] | "\(.name): \(.workspace // "NO_WORKSPACE_TRACKED")"' orchestrator-state.json
jq -r '.efforts_in_progress[] | "\(.name): \(.workspace // "NO_WORKSPACE_TRACKED")"' orchestrator-state.json

# List all effort directories (these are TARGET REPO CLONES, not SF branches!)
echo "📂 Effort directories (target repo clones):"
ls -la efforts/

# For each effort, identify the branch IN THE TARGET REPO CLONE
for effort_dir in efforts/*/; do
    if [[ -d "$effort_dir/.git" ]]; then
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "📁 Effort: $(basename $effort_dir)"
        echo "📍 Location: $effort_dir"
        cd "$effort_dir"
        
        # Verify this is a target repo clone, NOT software-factory
        REMOTE_URL=$(git remote get-url origin)
        if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
            echo "❌ ERROR: This is a Software Factory clone, not target repo!"
            continue
        fi
        
        echo "🔗 Repository: $REMOTE_URL"
        echo "🌿 Current branch: $(git branch --show-current)"
        echo "📝 Latest commit: $(git log -1 --oneline)"
        cd "$SF_INSTANCE_DIR"
    fi
done
```

### Step 2: Locate and Read Merge Order Plan
```bash
# CRITICAL: Merge plans are in integration workspaces, NOT in SF directory!
SF_INSTANCE_DIR=$(pwd)

# Check state file for integration workspace location
INTEGRATION_WORKSPACE=$(jq -r '.current_wave_integration.workspace // .current_phase_integration.workspace // "NOT_SET"' orchestrator-state.json)
MERGE_PLAN_PATH=$(jq -r '.current_wave_integration.merge_plan // .current_phase_integration.merge_plan // "NOT_SET"' orchestrator-state.json)

echo "📍 Integration workspace: $INTEGRATION_WORKSPACE"
echo "📄 Merge plan path: $MERGE_PLAN_PATH"

# Navigate to integration workspace if it exists
if [[ -d "$INTEGRATION_WORKSPACE" ]]; then
    cd "$INTEGRATION_WORKSPACE"
    echo "✅ Found integration workspace"
    
    # Look for merge plan
    if [[ -f "WAVE-MERGE-PLAN.md" ]]; then
        echo "📋 Found WAVE-MERGE-PLAN.md"
        cat WAVE-MERGE-PLAN.md
    elif [[ -f "PHASE-MERGE-PLAN.md" ]]; then
        echo "📋 Found PHASE-MERGE-PLAN.md"
        cat PHASE-MERGE-PLAN.md
    elif [[ -f "PROJECT-MERGE-PLAN.md" ]]; then
        echo "📋 Found PROJECT-MERGE-PLAN.md"
        cat PROJECT-MERGE-PLAN.md
    else
        echo "⚠️ No merge plan found - Code Reviewer should have created one!"
        ls -la *.md
    fi
    
    cd "$SF_INSTANCE_DIR"
else
    echo "❌ Integration workspace not found at: $INTEGRATION_WORKSPACE"
    echo "📂 Available workspaces:"
    find efforts -type d -name "integration-workspace" 2>/dev/null
fi
```

### Step 3: Verify Integration Freshness and Spawn Integration Agent
```bash
# 🔴🔴🔴 CRITICAL: CHECK INTEGRATION FRESHNESS FIRST! 🔴🔴🔴
verify_integration_freshness() {
    echo "🔍 Checking if integration branch is stale..."
    
    # Get integration branch creation time from state
    INTEGRATION_CREATED=$(jq -r '.current_wave_integration.created_at // .current_phase_integration.created_at' orchestrator-state.json)
    
    # Check each effort for newer commits
    STALE_INTEGRATION=false
    for effort in $(jq -r '.efforts_completed[].name, .efforts_in_progress[].name' orchestrator-state.json); do
        EFFORT_UPDATED=$(jq -r --arg name "$effort" '.efforts_completed[], .efforts_in_progress[] | select(.name == $name) | .last_updated_at // .completion_time' orchestrator-state.json)
        
        if [[ "$EFFORT_UPDATED" > "$INTEGRATION_CREATED" ]]; then
            echo "⚠️ Effort '$effort' has newer commits than integration branch!"
            echo "   Effort updated: $EFFORT_UPDATED"
            echo "   Integration created: $INTEGRATION_CREATED"
            STALE_INTEGRATION=true
        fi
    done
    
    if [[ "$STALE_INTEGRATION" == "true" ]]; then
        echo "🔴 CRITICAL: Integration branch is STALE!"
        echo "📝 Must recreate integration from fresh effort branches"
        echo "   The merge plan may reference outdated branches"
        return 1
    else
        echo "✅ Integration branch is fresh"
        return 0
    fi
}

# Get integration workspace location
SF_INSTANCE_DIR=$(pwd)
INTEGRATION_WORKSPACE=$(jq -r '.current_wave_integration.workspace // .current_phase_integration.workspace // .project_integration.workspace' orchestrator-state.json)

if [[ ! -d "$INTEGRATION_WORKSPACE" ]]; then
    echo "❌ Integration workspace not found: $INTEGRATION_WORKSPACE"
    echo "📝 Need to create integration infrastructure first"
    exit 1
fi

# Verify freshness before proceeding
cd "$SF_INSTANCE_DIR"
if ! verify_integration_freshness; then
    echo "🔄 Integration is stale - need fresh integration"
    # Document the staleness
    cat > "$INTEGRATION_WORKSPACE/INTEGRATION-STALENESS-REPORT.md" << 'EOF'
# Integration Staleness Detected
Date: $(date)
Issue: Effort branches have been updated after integration branch creation
Resolution: Integration Agent must fetch fresh branches
Action: Integration Agent will handle fresh merges
EOF
fi

# 🔴🔴🔴 R329 ENFORCEMENT: SPAWN INTEGRATION AGENT 🔴🔴🔴
echo "📋 R329 MANDATORY: Spawning Integration Agent for ALL merges"
echo "🚫 Orchestrator MUST NOT execute merges directly"

# Check if merge plan exists
if [[ -f "$INTEGRATION_WORKSPACE/WAVE-MERGE-PLAN.md" ]] || \
   [[ -f "$INTEGRATION_WORKSPACE/PHASE-MERGE-PLAN.md" ]] || \
   [[ -f "$INTEGRATION_WORKSPACE/PROJECT-MERGE-PLAN.md" ]]; then
    
    MERGE_PLAN=$(ls "$INTEGRATION_WORKSPACE"/*MERGE-PLAN.md | head -1)
    echo "✅ Found merge plan: $(basename $MERGE_PLAN)"
    
    # Spawn Integration Agent to execute merges
    echo "🚀 SPAWNING INTEGRATION AGENT (R329 COMPLIANCE)"
    
    cat > /tmp/integration-agent-task.md << EOF
# Integration Testing Execution Task

## Critical Requirements (R329)
- You are the Integration Agent, responsible for ALL merge operations
- The orchestrator has delegated this work per R329
- You MUST execute all merges according to the plan

## Working Directory
$INTEGRATION_WORKSPACE

## Merge Plan Location  
$(basename $MERGE_PLAN)

## Instructions
1. CD to integration workspace: $INTEGRATION_WORKSPACE
2. Read and follow $(basename $MERGE_PLAN) EXACTLY
3. Fetch latest branches from origin
4. Execute merges in specified order
5. Document any conflicts in INTEGRATION-REPORT.md
6. Handle merge conflicts if resolvable
7. Create comprehensive report of all operations

## Expected Outputs
- INTEGRATION-REPORT.md with full details
- work-log.md showing all commands executed
- Merged integration branch pushed to remote

## Staleness Check
$(if [[ "$STALE_INTEGRATION" == "true" ]]; then
    echo "⚠️ CRITICAL: Integration branches are stale"
    echo "MUST fetch fresh branches before merging"
else
    echo "✅ Integration branches are fresh"
fi)
EOF

    # Spawn the Integration Agent
    echo "📋 Spawning Integration Agent with task..."
    /spawn integration-agent EXECUTE_MERGES "$(cat /tmp/integration-agent-task.md)"
    
    # Update state to track spawning
    echo "📝 Updating state to MONITORING_INTEGRATION"
    # Transition to monitoring the Integration Agent
    
else
    echo "❌ No merge plan found in $INTEGRATION_WORKSPACE"
    echo "📝 Must spawn Code Reviewer first to create merge plan"
    
    # Spawn Code Reviewer to create merge plan
    echo "🚀 Spawning Code Reviewer to create merge plan..."
    /spawn code-reviewer CREATE_MERGE_PLAN "Create merge plan for integration testing in $INTEGRATION_WORKSPACE"
fi

cd "$SF_INSTANCE_DIR"

# 🔴 CRITICAL R329 REMINDER 🔴
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚨 R329 ENFORCEMENT COMPLETE"
echo "✅ Integration Agent spawned for merge execution"
echo "❌ Orchestrator did NOT execute any merges"
echo "📊 Waiting for Integration Agent reports..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
```

### Step 4: Create Integration Report
```bash
cat > INTEGRATION-TESTING-REPORT.md << 'EOF'
# Integration Testing Report
Date: $(date)
State: INTEGRATION_TESTING

## Efforts Merged
1. [effort-name] - branch-name - ✅ Success / ⚠️ Conflicts
2. ...

## Conflicts Encountered
- [If any, list here with resolution approach]

## Build Status
- Compilation: [PENDING - will check in BUILD_VALIDATION]
- Tests: [PENDING - will run in PRODUCTION_READY_VALIDATION]

## Next Steps
Transition to PRODUCTION_READY_VALIDATION for full validation
EOF
```

## ⚠️ CRITICAL REQUIREMENTS

### 🔴 R006 VIOLATION DETECTION 
If you find yourself about to:
- Edit code to resolve merge conflicts
- Modify source files during integration
- Apply fixes to make branches compatible
- Use git commands to edit code content
- Make excuses like "it's simpler to fix during merge"

**STOP IMMEDIATELY - You are violating R006!**
Spawn SW Engineers to handle ALL conflict resolution and fixes!

### Merge Order Matters
- Dependencies must be merged before dependents
- Core/shared libraries first
- Feature branches after infrastructure
- UI/frontend typically last

### Conflict Resolution Protocol
1. **ANY Conflicts**: STOP - Spawn SW Engineers to resolve
2. **NEVER resolve conflicts yourself** - R006 VIOLATION
3. **Document conflicts for SW Engineers to fix**
4. **Breaking Changes**: Require SW Engineer intervention

### Backport Tracking
**CRITICAL**: Track any fixes made during conflict resolution:
```bash
# If you fix conflicts during merge
echo "Fixed conflict in file X" >> BACKPORT-REQUIREMENTS.txt
```

These MUST be backported to original branches later!

## 🚫 FORBIDDEN ACTIONS

1. **NEVER edit any code files yourself** - R006 VIOLATION = -100%
2. **NEVER resolve merge conflicts yourself** - R006 VIOLATION = -100%
3. **NEVER apply patches or fixes directly** - R006 VIOLATION = -100%
4. **NEVER merge directly to main branch**
5. **NEVER skip efforts even if they seem independent**
6. **NEVER force merge with conflicts unresolved**
7. **NEVER modify effort code during merge** - ALL code changes require SW Engineers

## ✅ SUCCESS CRITERIA

Before transitioning to PRODUCTION_READY_VALIDATION:
- [ ] All identified efforts merged into integration-testing
- [ ] All merge conflicts documented
- [ ] INTEGRATION-TESTING-REPORT.md created
- [ ] No unresolved conflicts remain
- [ ] Branch pushed to remote

## 🔄 STATE TRANSITIONS

### Success Path:
```
INTEGRATION_TESTING → PRODUCTION_READY_VALIDATION
```
- All efforts merged successfully
- Integration report created
- Ready for validation

### Error Path:
```
INTEGRATION_TESTING → ANALYZE_BUILD_FAILURES
```
- Complex merge conflicts require fixes
- Structural incompatibilities found
- Need to create fix plan

## 📊 VERIFICATION CHECKLIST

Before leaving this state:
```bash
# Verify all efforts merged
git log --oneline --graph -20

# Verify no uncommitted changes
git status

# Verify integration report exists
ls -la INTEGRATION-TESTING-REPORT.md

# Verify branch pushed
git push origin integration-testing
```

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Systematic Approach** (25%)
   - Following merge order plan
   - Documenting each merge
   
2. **Conflict Handling** (25%)
   - Proper conflict identification
   - Appropriate resolution vs escalation
   
3. **Documentation** (25%)
   - Complete integration report
   - Backport requirements tracked
   
4. **Verification** (25%)
   - All efforts actually merged
   - No broken state left behind

## 💡 TIPS FOR SUCCESS

1. **Take Time to Plan**: Review all efforts before starting merges
2. **Test After Each Merge**: Run quick sanity checks
3. **Document Everything**: Future states need your reports
4. **Think Dependencies**: Merge order can prevent conflicts

## 🚨 COMMON PITFALLS TO AVOID

1. **Merging in Wrong Order**: Can create unnecessary conflicts
2. **Ignoring Small Conflicts**: They compound into bigger issues
3. **Not Tracking Fixes**: Backporting becomes impossible
4. **Skipping Verification**: Next states fail mysteriously

Remember: This state proves all efforts work together. Take it seriously!
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**


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
