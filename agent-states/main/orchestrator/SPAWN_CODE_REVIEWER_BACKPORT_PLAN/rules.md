# Orchestrator - SPAWN_CODE_REVIEWER_BACKPORT_PLAN State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R256** - Fix Planning Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R256-fix-planning-protocol.md`

4. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

5. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`

6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md`

7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

8. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-state-transition-validation.md`


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

**YOU HAVE ENTERED SPAWN_CODE_REVIEWER_BACKPORT_PLAN STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_BACKPORT_PLAN
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWER_BACKPORT_PLAN" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_BACKPORT_PLAN
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

## 🎯 STATE OBJECTIVES - SPAWN CODE REVIEWER FOR BACKPORT PLANNING

In the SPAWN_CODE_REVIEWER_BACKPORT_PLAN state, the ORCHESTRATOR is responsible for:

1. **Preparing Fix Documentation for Code Reviewer**
   - Gather all integration failure reports
   - Collect all fixes that were made during integration
   - Identify which effort branches are affected
   - Document what fixes need to be backported where

2. **Creating Code Reviewer Assignment**
   - Prepare clear instructions for Code Reviewer
   - Specify they must create BACKPORT-PLAN.md
   - Include all fix information and affected branches
   - Set clear expectations for plan structure

3. **Spawning Code Reviewer Agent**
   - Spawn Code Reviewer with BACKPORT_PLAN_CREATION state
   - Provide path to fix documentation
   - Monitor Code Reviewer state file creation
   - Track Code Reviewer progress

4. **Stopping After Spawn (R322 Part A)**
   - Update orchestrator-state.json to WAITING_FOR_BACKPORT_PLAN
   - Commit and push state changes
   - STOP and wait for user continuation
   - Do NOT continue to next state automatically

## 📝 REQUIRED ACTIONS

### Step 1: Prepare Fix Documentation
```bash
# Create comprehensive fix documentation for Code Reviewer
cd /efforts/integration-testing

# Create fix manifest for Code Reviewer
cat > FIX-MANIFEST-FOR-BACKPORT.md << 'EOF'
# Fix Manifest for Backport Planning

## Integration Fixes Applied
[Document all fixes made during integration]

## Affected Effort Branches
[List each effort and what fixes it needs]

## Source Files Modified
[List all files that were changed]

## Test Failures Resolved
[Document what tests were failing and how they were fixed]

## Build Issues Fixed
[Document build problems and resolutions]

## Required Backports
Code Reviewer must analyze these fixes and create a detailed
BACKPORT-PLAN.md that specifies:
1. Which fixes go to which effort branches
2. Order of application
3. Verification requirements
4. Success criteria
EOF

echo "✅ Fix documentation prepared for Code Reviewer"
```

### Step 2: Create Code Reviewer Assignment
```bash
# Create clear assignment for Code Reviewer
cat > CODE-REVIEWER-BACKPORT-ASSIGNMENT.md << 'EOF'
# Code Reviewer Assignment: Create Backport Plan

## Your Task
You are being spawned to create a comprehensive BACKPORT-PLAN.md that will guide
SW Engineers in applying integration fixes back to source effort branches.

## Input Documentation
- FIX-MANIFEST-FOR-BACKPORT.md - Contains all fixes made during integration
- Integration test results and failure reports
- Current state of effort branches

## Required Output: BACKPORT-PLAN.md
Create a detailed plan that includes:

### For Each Effort Branch:
1. **Branch Name**: The exact branch to update
2. **Working Directory**: /efforts/[effort-name]
3. **Fixes Required**: Specific fixes from integration
4. **Files to Modify**: Exact file paths and changes
5. **Verification Steps**: How to verify fixes work
6. **Dependencies**: Any order requirements

### Plan Structure:
- Group fixes by effort branch
- Specify exact commands or cherry-picks needed
- Include validation criteria
- Document any risks or special considerations

## Success Criteria
- Every integration fix is mapped to source branches
- Clear instructions for SW Engineers
- No ambiguity about what goes where
- Verification steps for each backport

## Working Directory
/efforts/integration-testing

## State File
Create: /efforts/integration-testing/code-reviewer-state.yaml
Update current_state to BACKPORT_PLAN_COMPLETE when done
EOF

echo "✅ Code Reviewer assignment created"
```

### Step 3: Spawn Code Reviewer Agent
```bash
echo "🚀 Spawning Code Reviewer for backport planning..."

# Log the spawn action
echo "$(date): Spawning Code Reviewer for BACKPORT_PLAN_CREATION" >> SPAWN-LOG.md

# The actual spawn would be done through the Claude interface
cat > /tmp/spawn-code-reviewer-command.md << 'EOF'
@agent-code-reviewer

## BACKPORT PLAN CREATION ASSIGNMENT

You are being spawned to create a comprehensive backport plan for integration fixes.

### Your Immediate Tasks:
1. Read your assignment at: /efforts/integration-testing/CODE-REVIEWER-BACKPORT-ASSIGNMENT.md
2. Read the fix manifest at: /efforts/integration-testing/FIX-MANIFEST-FOR-BACKPORT.md
3. Analyze what fixes need to go to which effort branches
4. Create detailed BACKPORT-PLAN.md with clear instructions for SW Engineers
5. Update your state file when complete

### Working Directory: 
/efforts/integration-testing

### State Transition:
- Initial State: INIT
- Target State: BACKPORT_PLAN_CREATION
- Final State: BACKPORT_PLAN_COMPLETE

### Critical Requirements:
- Map EVERY integration fix to its source branches
- Provide EXACT instructions for SW Engineers
- Include verification steps
- No ambiguity allowed

Start immediately upon spawn.
EOF

echo "✅ Code Reviewer spawn command prepared"
```

### Step 4: Update State and STOP (R322 Part A Enforcement)
```bash
# Update orchestrator state to waiting
cd $CLAUDE_PROJECT_DIR

# Update state file
cat > orchestrator-state.json << 'EOF'
current_state: WAITING_FOR_BACKPORT_PLAN
previous_state: SPAWN_CODE_REVIEWER_BACKPORT_PLAN
agents_spawned:
  - agent: code-reviewer
    purpose: Create backport plan for integration fixes
    state: BACKPORT_PLAN_CREATION
    timestamp: $(date +%s)
backport_status: PLANNING
integration_branch: integration-testing
waiting_for:
  - Code Reviewer to complete BACKPORT-PLAN.md
  - Plan to map all fixes to source branches
EOF

# Commit state change
git add orchestrator-state.json
git commit -m "state: transition to WAITING_FOR_BACKPORT_PLAN after spawning Code Reviewer"
git push

echo "✅ State updated and committed"
echo "🛑 STOPPING per R322 Part A - Must stop after spawn state"
```

## ⚠️ CRITICAL REQUIREMENTS

### Clear Separation of Responsibilities
- **Orchestrator**: ONLY coordinates and spawns
- **Code Reviewer**: Creates the backport plan
- **SW Engineers**: Will implement the fixes (next state)

### No Direct Code Work
- Orchestrator MUST NOT analyze fixes directly
- Orchestrator MUST NOT create the plan itself
- Orchestrator ONLY prepares documentation and spawns

### R322 Part A Compliance
- MUST stop after spawning Code Reviewer
- MUST update state to WAITING_FOR_BACKPORT_PLAN
- MUST NOT continue automatically
- Wait for /continue-orchestrating

## 🚫 FORBIDDEN ACTIONS

1. **Creating the backport plan yourself** - Must be done by Code Reviewer
2. **Analyzing code changes directly** - Code Reviewer's responsibility
3. **Continuing past spawn without stopping** - R322 Part A violation
4. **Spawning SW Engineers in this state** - That's the next state
5. **Making any code edits** - R006 violation

## ✅ SUCCESS CRITERIA

Before transitioning to WAITING_FOR_BACKPORT_PLAN:
- [ ] Fix documentation prepared for Code Reviewer
- [ ] Clear assignment created with expectations
- [ ] Code Reviewer spawned with proper state
- [ ] Orchestrator state updated to WAITING_FOR_BACKPORT_PLAN
- [ ] State changes committed and pushed
- [ ] STOPPED per R322 Part A requirements

## 🔄 STATE TRANSITIONS

### Success Path:
```
SPAWN_CODE_REVIEWER_BACKPORT_PLAN → WAITING_FOR_BACKPORT_PLAN
```
- Code Reviewer spawned successfully
- Waiting for plan completion
- Will continue after user command

### Next States After Waiting:
```
WAITING_FOR_BACKPORT_PLAN → SPAWN_SW_ENGINEER_BACKPORT_FIXES
```
- Once plan is ready
- Spawn SW Engineers to implement

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Separation of Concerns** (40%)
   - Code Reviewer does planning
   - Orchestrator only coordinates
   
2. **R322 Part A Compliance** (30%)
   - Proper stop after spawn
   - State update before stop
   
3. **Documentation Quality** (20%)
   - Clear instructions for Code Reviewer
   - Complete fix information provided
   
4. **State Management** (10%)
   - Proper state transitions
   - State file updates

## 💡 TIPS FOR SUCCESS

1. **Let Code Reviewer analyze** - Don't try to understand fixes yourself
2. **Provide complete information** - Give Code Reviewer everything needed
3. **Stop means stop** - R322 Part A is absolute
4. **Clear expectations** - Tell Code Reviewer exactly what output you need

Remember: This state is about DELEGATION, not doing the work yourself!

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
