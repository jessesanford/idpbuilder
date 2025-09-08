# Orchestrator - BUILD_VALIDATION State Rules

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

**YOU HAVE ENTERED BUILD_VALIDATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_BUILD_VALIDATION
echo "$(date +%s) - Rules read and acknowledged for BUILD_VALIDATION" > .state_rules_read_orchestrator_BUILD_VALIDATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY BUILD WORK UNTIL RULES ARE READ:
- ❌ Start compilation
- ❌ Run build commands
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

**⚠️ R006 WARNING FOR BUILD_VALIDATION STATE:**
- DO NOT fix compilation errors yourself!
- DO NOT edit code to make it compile!
- DO NOT modify any source files!
- If build fails, document issues for SW Engineers
- You only run builds and report - NEVER fix code

### 🚨🚨🚨 RULE R319 - Orchestrator NEVER Measures or Assesses Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality:** BLOCKING - Any technical assessment = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`

**⚠️ R319 WARNING FOR BUILD_VALIDATION STATE:**
- DO NOT run builds yourself!
- DO NOT execute test commands!
- DO NOT validate artifacts!
- DO NOT assess technical compliance!
- You MUST spawn Code Reviewer to handle ALL build validation
- You only coordinate and read reports - NEVER validate directly

### 🚨🚨🚨 RULE R323 - Mandatory Final Artifact Build [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R323-mandatory-final-artifact-build.md`
**Criticality:** BLOCKING - No artifact = -50% to -100% FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R323-mandatory-final-artifact-build.md`

**⚠️ R323 ENFORCEMENT FOR BUILD_VALIDATION STATE:**
- Code Reviewer MUST build final deliverable artifact
- Code Reviewer MUST verify artifact exists
- Code Reviewer MUST document artifact path, size, type
- CANNOT transition to SUCCESS without artifact confirmation
- NO project can complete without built binary/package

### Build Validation Requirements

**This state coordinates build validation through the Code Reviewer agent.**

Key requirements:
- Code Reviewer must validate builds succeed
- Code Reviewer must verify artifacts are generated
- **🔴 R323: Code Reviewer MUST build final deliverable artifact**
- **🔴 R323: Cannot proceed without artifact confirmation**
- Code Reviewer must document warnings
- Any fixes must be tracked for backporting via SW Engineers

## 🎯 STATE OBJECTIVES

In the BUILD_VALIDATION state, you are responsible for:

1. **Spawning Code Reviewer for Build Validation**
   - Spawn Code Reviewer agent to validate builds
   - Provide clear instructions for what to validate
   - **🚨 R323: Instruct Code Reviewer to BUILD FINAL ARTIFACT**
   - Specify integration workspace location
   - Request comprehensive build report with artifact details

2. **Monitoring Validation Progress**
   - Check if Code Reviewer has completed validation
   - Read build validation reports
   - **🚨 R323: Verify artifact was built and documented**
   - Track any issues identified
   - Coordinate fixes if needed

3. **Coordinating Issue Resolution**
   - If build fails, spawn SW Engineers to fix
   - Track which efforts need fixes
   - Ensure fixes are backported
   - Re-validate after fixes
   - **🚨 R323: If no artifact built, BLOCK progress**

4. **Managing State Transitions**
   - Determine next state based on validation results
   - **🚨 R323: CANNOT transition to SUCCESS without artifact**
   - Update state file with build status AND artifact details
   - Track validation completion

## 📝 REQUIRED ACTIONS

### Step 1: Spawn Code Reviewer for Build Validation
```bash
# ✅ CORRECT: Delegate build validation to Code Reviewer
echo "🏗️ Build validation needed for integrated code"
echo "🚀 Spawning Code Reviewer to validate builds..."

# Update state to show spawning Code Reviewer
yq -i '.current_state = "SPAWN_CODE_REVIEWER_BUILD_VALIDATION"' orchestrator-state.yaml
yq -i '.spawn_in_progress.agent = "code-reviewer"' orchestrator-state.yaml
yq -i '.spawn_in_progress.purpose = "build_validation"' orchestrator-state.yaml
yq -i '.spawn_in_progress.workspace = "/efforts/integration-testing"' orchestrator-state.yaml

# Spawn Code Reviewer with build validation task
Task: subagent_type="code-reviewer" \
      state="BUILD_VALIDATION" \
      prompt="Validate that integrated code builds successfully. Run full build process, verify all artifacts are generated, document any warnings or issues. Create BUILD-VALIDATION-REPORT.md with findings." \
      workspace="/efforts/integration-testing" \
      description="Build validation for integrated code"

echo "⏳ Code Reviewer spawned for build validation"
echo "📋 Waiting for BUILD-VALIDATION-REPORT.md"
```

### Step 2: Monitor Validation Progress
```bash
# Check if Code Reviewer has completed
VALIDATION_REPORT="/efforts/integration-testing/BUILD-VALIDATION-REPORT.md"

if [ -f "$VALIDATION_REPORT" ]; then
    echo "✅ Build validation report received"
    
    # Read the report
    BUILD_STATUS=$(grep "Status:" "$VALIDATION_REPORT" | head -1)
    echo "Build validation result: $BUILD_STATUS"
    
    # Update state file
    yq -i '.build_validation.completed = true' orchestrator-state.yaml
    yq -i ".build_validation.status = \"$(echo $BUILD_STATUS | cut -d: -f2 | xargs)\"" orchestrator-state.yaml
else
    echo "⏳ Waiting for Code Reviewer to complete validation"
    echo "📍 Expected report: $VALIDATION_REPORT"
fi
```

### Step 3: Handle Validation Results
```bash
# Based on Code Reviewer's report, determine next action
if grep -q "Status: SUCCESS" "$VALIDATION_REPORT"; then
    echo "✅ Build validation passed - ready for PR creation"
    NEXT_STATE="PR_PLAN_CREATION"
elif grep -q "Status: FAILED" "$VALIDATION_REPORT"; then
    echo "❌ Build validation failed - fixes required"
    
    # Check if fixes need backporting
    if grep -q "Backport Required: Yes" "$VALIDATION_REPORT"; then
        echo "📝 Backporting required to effort branches"
        NEXT_STATE="IMMEDIATE_BACKPORT_REQUIRED"
    else
        echo "🔧 Build fixes needed"
        NEXT_STATE="FIX_BUILD_ISSUES"
    fi
else
    echo "⚠️ Unable to determine build status"
    NEXT_STATE="ERROR_RECOVERY"
fi

echo "📊 Next state: $NEXT_STATE"
yq -i ".current_state = \"$NEXT_STATE\"" orchestrator-state.yaml

```

### Step 4: Coordinate Fix Resolution (If Needed)
```bash
# If Code Reviewer reports build failures, spawn SW Engineers to fix
if grep -q "Status: FAILED" "$VALIDATION_REPORT"; then
    echo "🔧 Build failures detected - need SW Engineers to fix"
    
    # Extract which efforts need fixes
    EFFORTS_NEEDING_FIXES=$(grep "Effort:" "$VALIDATION_REPORT" | cut -d: -f2)
    
    for effort in $EFFORTS_NEEDING_FIXES; do
        echo "🚀 Spawning SW Engineer to fix build issues in $effort"
        
        Task: subagent_type="software-engineer" \
              state="FIX_BUILD_ISSUES" \
              prompt="Fix build issues identified in BUILD-VALIDATION-REPORT.md. Apply fixes and ensure build succeeds." \
              workspace="/efforts/$effort" \
              description="Fix build issues in $effort"
    done
    
    echo "⏳ SW Engineers spawned for build fixes"
fi
```

### Step 5: Verify Backport Requirements
```bash
# Check if Code Reviewer identified fixes that need backporting
if grep -q "Backport Required: Yes" "$VALIDATION_REPORT"; then
    echo "📝 Backporting required to original effort branches"
    
    # Extract backport requirements
    BACKPORT_LIST=$(grep -A10 "Backport Requirements:" "$VALIDATION_REPORT")
    echo "$BACKPORT_LIST"
    
    # Update state for backporting
    yq -i '.backport_required = true' orchestrator-state.yaml
    yq -i '.current_state = "IMMEDIATE_BACKPORT_REQUIRED"' orchestrator-state.yaml
    
    echo "🔄 Transitioning to IMMEDIATE_BACKPORT_REQUIRED state"
fi
```

## ⚠️ CRITICAL REQUIREMENTS

### 🔴 R006 VIOLATION DETECTION 
If you find yourself about to:
- Edit code to fix build errors
- Modify imports or dependencies in code
- Fix syntax errors directly
- Apply compilation fixes yourself
- Make excuses like "just a quick fix to make it compile"

**STOP IMMEDIATELY - You are violating R006!**
Document the issues and spawn SW Engineers to fix them!

### Backport Tracking Is MANDATORY
**EVERY fix made during build MUST be:**
1. Documented in BACKPORT-REQUIRED-FIXES.md
2. Tagged with the original effort branch
3. Backported before PR creation
4. Verified in original branch

### Build Must Be Reproducible
- Same commands must produce same output
- No timestamps in artifacts (unless versioned)
- Dependencies must be locked/pinned
- Build environment must be documented

### Zero Tolerance for Critical Errors
- Compilation errors = MUST FIX
- Linking errors = MUST FIX
- Missing dependencies = MUST FIX
- Runtime errors in build = MUST FIX

## 🚫 FORBIDDEN ACTIONS

1. **NEVER ignore build failures**
2. **NEVER skip backport tracking**
3. **NEVER modify code without documenting for backport**
4. **NEVER proceed with partial builds**
5. **NEVER fake build success**

## ✅ SUCCESS CRITERIA

Before transitioning to next state:
- [ ] Build completes without errors
- [ ] All expected artifacts generated
- [ ] Build warnings documented
- [ ] Any fixes tracked for backporting
- [ ] Build validation report created
- [ ] All changes committed and pushed

## 🔄 STATE TRANSITIONS

### Success Path (No Fixes Needed):
```
BUILD_VALIDATION → PR_PLAN_CREATION
```
- Build successful
- No fixes required
- Ready for PR planning

### Fix Required Path:
```
BUILD_VALIDATION → FIX_BUILD_ISSUES
```
- Build failures found
- Fixes too complex for inline resolution
- Need dedicated fixing state

### Backport Required Path:
```
BUILD_VALIDATION → BACKPORT_FIXES
```
- Build succeeded after fixes
- Fixes must be backported first
- Then proceed to PR_PLAN_CREATION

## 📊 VERIFICATION CHECKLIST

Before leaving this state:
```bash
# Verify build output captured
ls -la BUILD-OUTPUT.log

# Verify report created
ls -la BUILD-VALIDATION-REPORT.md

# Check for backport requirements
if [ -f "BACKPORT-REQUIRED-FIXES.md" ]; then
    echo "⚠️ BACKPORT REQUIRED - Must transition to BACKPORT_FIXES"
    cat BACKPORT-REQUIRED-FIXES.md
fi

# Commit all changes
git add -A
git commit -m "build: validation complete - $(if [ "$BUILD_SUCCESS" = true ]; then echo "SUCCESS"; else echo "FIXES REQUIRED"; fi)"
git push
```

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Build Execution** (25%)
   - Correct build commands used
   - All modules built
   
2. **Issue Resolution** (25%)
   - Build errors identified
   - Appropriate fixes applied
   
3. **Backport Tracking** (25%)
   - ALL fixes documented
   - Original branches identified
   
4. **Documentation** (25%)
   - Complete build report
   - Clear next steps

## 💡 TIPS FOR SUCCESS

1. **Clean Builds**: Always start fresh
2. **Read Errors Carefully**: First error often causes cascade
3. **Track Everything**: Every change needs backporting
4. **Test Fixes**: Verify fixes actually work

## 🚨 COMMON PITFALLS TO AVOID

1. **Forgetting Backport Tracking**: Critical failure point
2. **Ignoring Warnings**: They matter in production
3. **Partial Fixes**: Fix completely or escalate
4. **Not Verifying Artifacts**: Ensure deployment ready

Remember: Build validation ensures deployability. Backport tracking ensures clean PRs!
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
