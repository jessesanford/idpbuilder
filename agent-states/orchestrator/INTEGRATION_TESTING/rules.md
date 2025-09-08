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

**⚠️ R006 WARNING FOR INTEGRATION_TESTING STATE:**
- DO NOT resolve merge conflicts yourself!
- DO NOT edit code to fix integration issues!
- DO NOT apply patches or fixes directly!
- Document all issues for SW Engineers to resolve
- You only coordinate merging - NEVER modify code

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

## 🎯 STATE OBJECTIVES

In the INTEGRATION_TESTING state, you are responsible for:

1. **Merging All Efforts Sequentially**
   - Read all effort directories to identify branches
   - Determine dependency order based on implementation plan
   - Merge each effort one at a time
   - Handle conflicts systematically

2. **Verifying Integration Success**
   - After each merge, verify no conflicts remain
   - Run basic compilation/build checks if applicable
   - Document any issues encountered

3. **Creating Integration Report**
   - Document merge order executed
   - List any conflicts resolved
   - Note any issues for BUILD_VALIDATION state

## 📝 REQUIRED ACTIONS

### Step 1: Identify All Efforts
```bash
# List all effort directories
ls -la efforts/

# For each effort, identify the branch
for effort in efforts/*/; do
    echo "Effort: $(basename $effort)"
    cd "$effort"
    git branch --show-current
done
```

### Step 2: Read Merge Order Plan
Check for existing merge plan from previous states:
```bash
# Check for merge plan
if [ -f "INTEGRATION-MERGE-PLAN.md" ]; then
    cat INTEGRATION-MERGE-PLAN.md
else
    echo "⚠️ No merge plan found - will use alphabetical order"
fi
```

### Step 3: Execute Sequential Merges
```bash
# Change to integration-testing workspace
cd efforts/integration-testing

# Ensure we're on integration-testing branch
git checkout integration-testing

# For each effort in order
for effort_branch in "${ORDERED_EFFORTS[@]}"; do
    echo "📝 Merging: $effort_branch"
    
    # Fetch the effort branch
    git fetch origin "$effort_branch"
    
    # Attempt merge
    if git merge "origin/$effort_branch" --no-ff -m "integrate: merge $effort_branch into integration-testing"; then
        echo "✅ Successfully merged $effort_branch"
    else
        echo "⚠️ Conflicts detected in $effort_branch"
        # Document conflicts
        git status --short > "CONFLICTS-$effort_branch.txt"
        
        # Attempt automatic resolution for simple conflicts
        # If complex, transition to FIX_BUILD_ISSUES state
    fi
done
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
INTEGRATION_TESTING → FIX_BUILD_ISSUES
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
