# Orchestrator - SPAWN_CODE_REVIEWER_BUILD_VALIDATION State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
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

**YOU HAVE ENTERED SPAWN_CODE_REVIEWER_BUILD_VALIDATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R313 ENFORCEMENT - MANDATORY STOP AFTER SPAWN 🔴🔴🔴

**THIS IS A SPAWN STATE - YOU MUST STOP IMMEDIATELY AFTER SPAWNING!**

Per R313, after spawning the Code Reviewer:
1. Record spawn details in state file
2. Save TODOs and commit state changes  
3. EXIT with clear continuation instructions
4. DO NOT wait for or process agent responses

**Waiting for responses causes context overflow and rule forgetting!**

## 📋 PRIMARY DIRECTIVES

### 🚨🚨🚨 RULE R006 - Orchestrator NEVER Writes Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality:** BLOCKING - Any code operation = -100% IMMEDIATE FAILURE

### 🚨🚨🚨 RULE R319 - Orchestrator NEVER Measures or Assesses Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality:** BLOCKING - Any technical assessment = -100% IMMEDIATE FAILURE

### 🚨🚨🚨 RULE R313 - Mandatory Stop After Spawning Agents [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R313-mandatory-stop-after-spawn.md`
**Criticality:** BLOCKING - Must stop after spawning to preserve context

**KEY POINTS:**
- Orchestrator CANNOT run builds (R006/R319)
- Orchestrator CANNOT validate artifacts (R319)
- Orchestrator MUST spawn Code Reviewer for ALL build validation
- Orchestrator MUST stop immediately after spawning (R313)

## 🎯 STATE OBJECTIVES

In SPAWN_CODE_REVIEWER_BUILD_VALIDATION state, you must:

1. **Prepare Build Validation Context**
   - Identify integration workspace location
   - Determine what needs validation
   - Prepare clear instructions for Code Reviewer

2. **Spawn Code Reviewer Agent**
   - Spawn with BUILD_VALIDATION state
   - Provide integration workspace path
   - Give clear validation requirements

3. **Update State Tracking**
   - Record spawn in orchestrator-state.yaml
   - Note purpose of validation
   - Track expected report location

4. **Stop Immediately**
   - Save TODOs
   - Commit state
   - Exit with continuation instructions

## 📝 REQUIRED ACTIONS

### Step 1: Prepare Spawn Context
```bash
# Get current phase and wave
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)

# Identify integration workspace
INTEGRATION_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"

# Verify workspace exists
if [ ! -d "$INTEGRATION_DIR" ]; then
    echo "❌ Integration workspace not found at: $INTEGRATION_DIR"
    echo "Must complete INTEGRATION state first!"
    exit 1
fi

echo "📦 Build validation needed for: $INTEGRATION_DIR"
```

### Step 2: Update State for Spawn
```bash
# Record spawn details BEFORE spawning
yq -i '.spawn_in_progress.agent = "code-reviewer"' orchestrator-state.yaml
yq -i '.spawn_in_progress.purpose = "build_validation"' orchestrator-state.yaml
yq -i '.spawn_in_progress.workspace = "'$INTEGRATION_DIR'"' orchestrator-state.yaml
yq -i '.spawn_in_progress.expected_report = "'$INTEGRATION_DIR'/BUILD-VALIDATION-REPORT.md"' orchestrator-state.yaml
yq -i '.spawn_in_progress.spawned_at = "'$(date -Iseconds)'"' orchestrator-state.yaml

echo "✅ State updated with spawn details"
```

### Step 3: Spawn Code Reviewer
```bash
# Create spawn task
cat > /tmp/build-validation-task.md << EOF
Validate that integrated code builds successfully.

CRITICAL REQUIREMENTS:
1. You are in BUILD_VALIDATION state
2. Working directory: $INTEGRATION_DIR
3. Run full build process
4. Verify all artifacts are generated
5. Document any warnings or errors
6. Create BUILD-VALIDATION-REPORT.md with findings
7. If build fails, identify which effort branches need fixes

DO NOT:
- Fix code directly (R321 - fixes go to source branches)
- Skip any build steps
- Proceed if build fails

Your report will determine if we proceed to PR creation or need fixes.
EOF

# Spawn the Code Reviewer
echo "🚀 Spawning Code Reviewer for build validation..."

Task: subagent_type="code-reviewer" \
      state="BUILD_VALIDATION" \
      workspace="$INTEGRATION_DIR" \
      prompt="$(cat /tmp/build-validation-task.md)" \
      description="Build validation for Phase $PHASE Wave $WAVE integration"

echo "✅ Code Reviewer spawned for build validation"
```

### Step 4: Save State and Exit (R313 MANDATORY)
```bash
# Save TODOs
save_todos "SPAWN_CODE_REVIEWER_BUILD_VALIDATION"

# Commit state
cd $CLAUDE_PROJECT_DIR
git add orchestrator-state.yaml todos/
git commit -m "state: spawned Code Reviewer for build validation - Phase $PHASE Wave $WAVE"
git push

# Create continuation instructions
cat > CONTINUE.md << EOF
# Continuation Instructions

## Current Status
- State: SPAWN_CODE_REVIEWER_BUILD_VALIDATION
- Spawned: Code Reviewer for build validation
- Workspace: $INTEGRATION_DIR
- Expected Report: BUILD-VALIDATION-REPORT.md

## Next Steps
1. Wait for Code Reviewer to complete validation
2. Check for BUILD-VALIDATION-REPORT.md in integration workspace
3. Based on report:
   - If SUCCESS: Transition to PR_PLAN_CREATION
   - If FAILED: Transition to IMMEDIATE_BACKPORT_REQUIRED
   - If FIXES_NEEDED: Spawn SW Engineers for fixes

## To Continue
/spawn orchestrator WAITING_FOR_BUILD_VALIDATION
EOF

echo "📋 Continuation instructions saved to CONTINUE.md"

# EXIT IMMEDIATELY (R313)
echo "🛑 STOPPING per R313 - Context preserved"
echo "Continue with: /spawn orchestrator WAITING_FOR_BUILD_VALIDATION"
exit 0
```

## ⚠️ CRITICAL REQUIREMENTS

### R313 Compliance Is MANDATORY
**You MUST stop after spawning. No exceptions.**
- ❌ DO NOT wait for Code Reviewer response
- ❌ DO NOT read any output from Code Reviewer
- ❌ DO NOT continue to next state
- ✅ STOP immediately after spawn
- ✅ Save all state information
- ✅ Exit with clear instructions

### Build Validation Is Code Reviewer's Job
**Per R006/R319, you CANNOT:**
- ❌ Run build commands yourself
- ❌ Check if builds succeed
- ❌ Validate artifacts
- ❌ Assess build output
- ✅ ONLY spawn Code Reviewer to do this

### Clear Instructions Are Critical
The Code Reviewer needs:
- Exact workspace location
- What to validate
- Report format required
- Where to save report

## 🚫 FORBIDDEN ACTIONS

1. **NEVER run builds yourself** - Immediate failure
2. **NEVER wait for agent response** - R313 violation
3. **NEVER assess build results** - R319 violation  
4. **NEVER fix build issues** - R006 violation
5. **NEVER skip the spawn** - Must delegate

## ✅ SUCCESS CRITERIA

Before exiting:
- [ ] Code Reviewer spawned with BUILD_VALIDATION state
- [ ] Workspace path provided correctly
- [ ] State file updated with spawn details
- [ ] TODOs saved
- [ ] State committed and pushed
- [ ] Continuation instructions created
- [ ] Exited immediately per R313

## 🔄 STATE TRANSITIONS

After spawning (in next session):
```
SPAWN_CODE_REVIEWER_BUILD_VALIDATION 
    ↓ (spawn then stop)
WAITING_FOR_BUILD_VALIDATION
    ↓ (check report)
[Based on report:]
    → PR_PLAN_CREATION (if success)
    → IMMEDIATE_BACKPORT_REQUIRED (if fixes needed)
    → ERROR_RECOVERY (if issues)
```

## 📊 VERIFICATION

System will verify:
1. Code Reviewer was actually spawned
2. State file was updated
3. You stopped after spawning
4. No build commands were executed by orchestrator

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Proper Delegation** (40%)
   - Code Reviewer spawned correctly
   - No attempt to validate yourself
   
2. **R313 Compliance** (30%)
   - Stopped immediately after spawn
   - No waiting for response
   
3. **State Management** (30%)
   - State file updated
   - TODOs saved
   - Clear continuation

## 💡 REMEMBER

You are a COORDINATOR:
- You identify WHAT needs validation
- You spawn WHO can validate
- You track WHEN validation happens
- You NEVER validate yourself

The Code Reviewer is the technical expert who:
- Runs the builds
- Checks the output
- Validates artifacts
- Reports findings

Together, you maintain separation of duties!
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
