# Orchestrator - BACKPORT_FIXES State Rules

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


## 🔴🔴🔴 FULLY DEPRECATED STATE - DO NOT USE 🔴🔴🔴

**THIS STATE IS FULLY DEPRECATED AND MUST NOT BE USED!**

### WHY THIS STATE IS DEPRECATED:
1. **Violates R321**: Backporting must happen IMMEDIATELY, not deferred
2. **Mixed Responsibilities**: Combined planning and execution in one state
3. **Poor Separation**: Orchestrator was doing too much work

### USE THE NEW FLOW INSTEAD:
The backport process has been split into proper states with clear separation:

1. **SPAWN_CODE_REVIEWER_BACKPORT_PLAN** - Code Reviewer creates the plan
2. **WAITING_FOR_BACKPORT_PLAN** - Wait for plan completion
3. **SPAWN_SW_ENGINEER_BACKPORT_FIXES** - SW Engineers implement fixes
4. **MONITORING_BACKPORT_PROGRESS** - Monitor implementation progress

### IF YOU'RE IN THIS STATE:
**STOP IMMEDIATELY and transition to the new flow:**
```bash
# Transition to the new flow
cd $CLAUDE_PROJECT_DIR
cat > orchestrator-state.json << 'EOF'
current_state: SPAWN_CODE_REVIEWER_BACKPORT_PLAN
previous_state: BACKPORT_FIXES
migration_note: "Migrating from deprecated BACKPORT_FIXES to new flow"
EOF

git add orchestrator-state.json
git commit -m "migrate: from deprecated BACKPORT_FIXES to new backport flow"
git push
```

**DO NOT CONTINUE WITH THE OLD INSTRUCTIONS BELOW!**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED BACKPORT_FIXES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_BACKPORT_FIXES
echo "$(date +%s) - Rules read and acknowledged for BACKPORT_FIXES" > .state_rules_read_orchestrator_BACKPORT_FIXES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY BACKPORTING UNTIL RULES ARE READ:
- ❌ Start cherry-picking commits
- ❌ Switch to effort branches
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 CRITICAL WARNING - READ THIS FIRST 🔴🔴🔴

### THE ORCHESTRATOR IS BEING MONITORED FOR R006 VIOLATIONS!

**WE HAVE DETECTED ORCHESTRATORS TRYING TO:**
- Edit code directly instead of spawning SW Engineers
- Make excuses like "everything is already integrated"
- Claim backporting is "simpler" without engineers
- Apply fixes themselves using cherry-pick

**ALL OF THESE ARE R006 VIOLATIONS = IMMEDIATE -100% FAILURE**

The orchestrator is a MANAGER, not a DEVELOPER. You MUST spawn SW Engineers for ALL code changes, including backporting. NO EXCEPTIONS!

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 🔴🔴🔴 ABSOLUTE BLOCKING RULE - R006 ENFORCEMENT 🔴🔴🔴

### THE ORCHESTRATOR NEVER WRITES CODE - NO EXCEPTIONS!

**CRITICAL VIOLATIONS THAT CAUSE IMMEDIATE FAILURE:**
1. **Orchestrator editing ANY code file = IMMEDIATE -100% FAILURE**
2. **Orchestrator using Edit/Write tools on code = IMMEDIATE -100% FAILURE**
3. **Orchestrator applying fixes directly = IMMEDIATE -100% FAILURE**
4. **Orchestrator cherry-picking commits = IMMEDIATE -100% FAILURE**
5. **Orchestrator making excuses to avoid spawning engineers = IMMEDIATE -100% FAILURE**

### 🚨🚨🚨 NO EXCUSES ACCEPTED - THESE ARE ALL INVALID 🚨🚨🚨

**The following excuses are FORBIDDEN and indicate rule violation:**
- ❌ "Everything is already integrated, so I'll just apply fixes"
- ❌ "The backport process is simpler since it's all in one branch"
- ❌ "We only need to commit the fixes there"
- ❌ "I can handle this quickly myself"
- ❌ "It's just a small fix"
- ❌ "Spawning engineers would be overkill"

**REALITY CHECK:**
- Even if all efforts are integrated, backporting is MANDATORY
- Each effort branch MUST be independently mergeable to main
- Skipping backports = ALL PRs WILL FAIL
- The orchestrator is a MANAGER, not a DEVELOPER

## 📋 PRIMARY DIRECTIVES

### Backport Requirements - MUST BE DONE BY SW ENGINEERS

**ALL fixes made during integration testing MUST be backported to original effort branches BY SW ENGINEERS**

**MANDATORY PROCESS:**
1. **Orchestrator reads BACKPORT-REQUIRED-FIXES.md** (or similar manifest)
2. **Orchestrator spawns SW Engineer for EACH effort needing fixes**
3. **Each SW Engineer applies fixes to their assigned effort branch**
4. **Orchestrator monitors completion via state files**
5. **Orchestrator NEVER touches code directly**

Key requirements:
- Every fix must be tracked to its original effort branch
- Fixes must be applied BY SW ENGINEERS via cherry-pick or manual application
- Each backport must be verified BY SW ENGINEERS in the original branch
- Original branches must be updated and pushed BY SW ENGINEERS
- Orchestrator ONLY coordinates and monitors

After SW Engineers complete backporting, they verify the branch still builds and tests pass.

## 🎯 STATE OBJECTIVES - CRITICAL FOR PR SUCCESS

In the BACKPORT_FIXES state, the ORCHESTRATOR is responsible for:

1. **Reading Backport Manifest** (ORCHESTRATOR DOES THIS)
   - Identify all branches needing backports
   - List all fixes to be applied
   - Determine backport order

2. **Spawning SW Engineers for Backporting** (ORCHESTRATOR COORDINATES)
   - Spawn SW Engineer for each effort branch
   - Provide each engineer with their specific fixes
   - Monitor their progress via state files

3. **SW Engineers Apply Fixes** (SW ENGINEERS DO THIS)
   - SW Engineers check out their effort branches
   - SW Engineers apply the specific fixes
   - SW Engineers maintain branch integrity
   - SW Engineers verify builds and tests

4. **Creating Backport Completion Report** (ORCHESTRATOR DOES THIS)
   - Collect status from all SW Engineers
   - Document all backports completed
   - Verify all branches updated
   - Confirm PR readiness

## ⚠️⚠️⚠️ WHY THIS IS CRITICAL ⚠️⚠️⚠️

**WITHOUT PROPER BACKPORTING:**
- PRs will have merge conflicts with main
- Integration fixes won't be in PR branches
- Build will fail when PRs are tested
- Entire effort could be rejected

**This state ensures each PR branch has ALL fixes needed to merge cleanly!**

## 📝 REQUIRED ACTIONS - ORCHESTRATOR COORDINATES, SW ENGINEERS EXECUTE

### Step 1: ORCHESTRATOR - Load Backport Manifest
```bash
# ORCHESTRATOR DOES THIS - Read the backport manifest
cd /efforts/integration-testing

# Find and read the backport documentation
if [ -f "BACKPORT-REQUIRED-FIXES.md" ]; then
    echo "📋 Loading backport manifest..."
    cat BACKPORT-REQUIRED-FIXES.md
elif [ -f "BACKPORT-MANIFEST.md" ]; then
    echo "📋 Loading backport manifest..."
    cat BACKPORT-MANIFEST.md
else
    echo "❌ CRITICAL: No backport manifest found!"
    echo "Check for other backport documentation:"
    ls -la *BACKPORT* *backport* *FIX* 2>/dev/null
fi

# Identify which efforts need fixes
echo "📝 Identifying efforts needing backports..."
# Parse the manifest to get effort list
```

### Step 2: ORCHESTRATOR - Prepare Backport Instructions for Each SW Engineer
```bash
# ORCHESTRATOR CREATES INSTRUCTIONS FOR EACH SW ENGINEER
for effort in "${EFFORTS_NEEDING_BACKPORT[@]}"; do
    EFFORT_NAME="$effort"
    INSTRUCTION_FILE="/efforts/${EFFORT_NAME}/BACKPORT-INSTRUCTIONS.md"
    
    cat > "$INSTRUCTION_FILE" << 'EOF'
# Backport Instructions for SW Engineer

## Your Assignment
- Effort: EFFORT_NAME
- Branch: EFFORT_BRANCH
- Working Directory: /efforts/EFFORT_NAME

## Fixes to Apply
[List specific fixes from manifest for this effort]

## Process
1. Checkout your effort branch
2. Apply the listed fixes (cherry-pick or manual)
3. Verify build succeeds
4. Run tests
5. Push updated branch
6. Update sw-engineer-state.yaml when complete

## Success Criteria
- All fixes applied
- Build passes
- Tests pass
- Branch pushed
EOF
done
```

### Step 3: ORCHESTRATOR - Spawn SW Engineers
```bash
# ORCHESTRATOR SPAWNS SW ENGINEERS - NEVER DOES THE WORK ITSELF

echo "🚀 Spawning SW Engineers for backporting..."

# For each effort needing backports
for effort in "${EFFORTS_NEEDING_BACKPORT[@]}"; do
    echo "📋 Spawning SW Engineer for: $effort"
    
    # Create spawn command
    cat > /tmp/spawn-sw-engineer-${effort}.md << 'EOF'
@agent-software-engineer

## BACKPORT_FIXES Assignment

You are being spawned to apply integration fixes to your effort branch.

1. Read your instructions at: /efforts/${effort}/BACKPORT-INSTRUCTIONS.md
2. Apply all listed fixes to your branch
3. Verify build and tests pass
4. Push your updated branch
5. Update your state file when complete

Your working directory: /efforts/${effort}
Your branch: [branch name from manifest]

DO NOT modify anything outside your assigned effort.
EOF
    
    # Log the spawn
    echo "$(date): Spawned SW Engineer for $effort" >> BACKPORT-SPAWN-LOG.md
done

echo "✅ All SW Engineers spawned. Monitoring their progress..."
```

### Step 4: ORCHESTRATOR - Monitor SW Engineer Progress
```bash
# ORCHESTRATOR MONITORS - DOES NOT DO THE WORK

echo "📊 Monitoring backport progress..."

# Check each SW Engineer's state file
while true; do
    COMPLETED_COUNT=0
    TOTAL_COUNT=${#EFFORTS_NEEDING_BACKPORT[@]}
    
    for effort in "${EFFORTS_NEEDING_BACKPORT[@]}"; do
        STATE_FILE="/efforts/${effort}/sw-engineer-state.yaml"
        
        if [ -f "$STATE_FILE" ]; then
            STATE=$(grep "current_state:" "$STATE_FILE" | awk '{print $2}')
            
            if [ "$STATE" = "BACKPORT_COMPLETE" ]; then
                ((COMPLETED_COUNT++))
                echo "✅ $effort: Backport complete"
            else
                echo "⏳ $effort: In progress (state: $STATE)"
            fi
        else
            echo "⏳ $effort: Waiting for SW Engineer to start"
        fi
    done
    
    echo "Progress: $COMPLETED_COUNT / $TOTAL_COUNT complete"
    
    if [ $COMPLETED_COUNT -eq $TOTAL_COUNT ]; then
        echo "✅ All backports complete!"
        break
    fi
    
    sleep 30  # Check every 30 seconds
done
```

### Step 5: ORCHESTRATOR - Collect Results and Create Report
```bash
# ORCHESTRATOR COLLECTS RESULTS - DOES NOT APPLY FIXES

### Step 6: Create Backport Completion Report
```bash
cd /efforts/integration-testing

cat > BACKPORT-COMPLETION-REPORT.md << 'EOF'
# Backport Completion Report
Date: $(date)
State: BACKPORT_FIXES

## Branches Backported

### ✅ Successfully Backported
1. **Branch**: effort-1-branch
   - Commits Applied: 3
   - Build Status: PASS
   - Test Status: PASS
   - Push Status: SUCCESS

2. **Branch**: effort-2-branch
   - Commits Applied: 2
   - Build Status: PASS
   - Test Status: PASS
   - Push Status: SUCCESS

### ⚠️ Backport Issues
[List any branches with problems]

## Verification Summary
- Total Branches Updated: X
- All Builds Passing: YES/NO
- All Tests Passing: YES/NO
- All Branches Pushed: YES/NO

## PR Readiness
✅ All effort branches now contain integration fixes
✅ Each branch builds independently
✅ Each branch passes tests
✅ Ready for PR creation

## Files Changed Per Branch
[List for PR descriptions]

### effort-1-branch:
- src/file1.go - Fixed import issue
- config/config.yaml - Added missing field

### effort-2-branch:
- pkg/module/module.go - Fixed type error

## Next Steps
Proceed to PR_PLAN_CREATION state
EOF
```

## ⚠️ CRITICAL REQUIREMENTS

### Every Branch Must Be Updated
- **NO EXCEPTIONS**
- Missing even one branch = PR failures
- Check BACKPORT-MANIFEST.md thoroughly

### Verification Is Mandatory
**For EACH branch after backporting:**
1. Build must succeed
2. Tests must pass
3. Branch must push successfully

### Commit Messages Must Reference Integration
```bash
git commit -m "backport: fixes from integration testing

- Applied build fixes from integration-testing branch
- Fixes compilation error in module X
- Adds missing dependency Y
- Ref: integration-testing commits abc123, def456"
```

### Track Everything
- Starting commit hash
- Ending commit hash
- Files modified
- Build/test results

## 🚫 FORBIDDEN ACTIONS - VIOLATIONS = IMMEDIATE FAILURE

### ORCHESTRATOR-SPECIFIC FORBIDDEN ACTIONS:
1. **NEVER edit any code file yourself** - R006 VIOLATION = -100%
2. **NEVER use git cherry-pick yourself** - R006 VIOLATION = -100%
3. **NEVER apply patches yourself** - R006 VIOLATION = -100%
4. **NEVER use Edit/Write tools on code** - R006 VIOLATION = -100%
5. **NEVER make excuses to avoid spawning SW Engineers** - R006 VIOLATION = -100%
6. **NEVER rationalize that "integration makes it simpler"** - STILL MUST SPAWN ENGINEERS

### GENERAL FORBIDDEN ACTIONS:
1. **NEVER skip a branch listed in manifest**
2. **NEVER force push (--force) to branches**
3. **NEVER merge integration branch into effort branches**
4. **NEVER proceed if builds fail after backport**
5. **NEVER allow SW Engineers to modify code beyond documented fixes**

### DETECTION TRIGGERS FOR VIOLATIONS:
- If orchestrator uses `git cherry-pick` = FAIL
- If orchestrator uses `Edit` on .go/.py/.js files = FAIL  
- If orchestrator says "I'll apply the fixes" = FAIL
- If orchestrator says "simpler to just commit" = FAIL
- If orchestrator doesn't spawn SW Engineers = FAIL

## ✅ SUCCESS CRITERIA

Before transitioning to PR_PLAN_CREATION:
- [ ] All branches in manifest updated
- [ ] Each branch builds successfully
- [ ] Each branch tests pass
- [ ] All branches pushed to remote
- [ ] Backport completion report created
- [ ] All changes documented

## 🔄 STATE TRANSITIONS

### Success Path:
```
BACKPORT_FIXES → PR_PLAN_CREATION
```
- All backports complete
- All verifications pass
- Ready for PR planning

### Failure Path:
```
BACKPORT_FIXES → ERROR_RECOVERY
```
- Backport conflicts unresolvable
- Branches won't build after backport
- Structural issues found

## 📊 VERIFICATION CHECKLIST

Before leaving this state:
```bash
# Verify all branches updated
for effort in /efforts/*/; do
    if [ -d "$effort/.git" ]; then
        echo "Checking: $(basename $effort)"
        cd "$effort"
        git log --oneline -5
        echo "---"
    fi
done

# Verify backport report complete
cat BACKPORT-COMPLETION-REPORT.md

# Final verification
echo "Backport Checklist:"
echo "✓ All branches listed in manifest: YES/NO"
echo "✓ All branches build: YES/NO"
echo "✓ All branches pushed: YES/NO"
echo "✓ Ready for PRs: YES/NO"
```

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Completeness** (35%)
   - Every branch backported
   - No branches missed
   
2. **Correctness** (35%)
   - Right fixes to right branches
   - Builds work after backport
   
3. **Verification** (20%)
   - Each branch tested
   - Results documented
   
4. **Documentation** (10%)
   - Clear completion report
   - Ready for PR creation

## 💡 TIPS FOR SUCCESS

1. **Work Systematically**: One branch at a time
2. **Verify Immediately**: Test right after backporting
3. **Document Issues**: If something fails, record it
4. **Use Git Log**: Track what was applied

## 🚨 COMMON PITFALLS TO AVOID

1. **Missing Branches**: Check manifest twice
2. **Wrong Fixes**: Ensure fixes go to right branch
3. **Skip Verification**: Leads to PR failures
4. **Force Pushing**: Destroys history

## 🎯 THE ULTIMATE GOAL

**After this state, each effort branch should:**
- Contain its original implementation
- PLUS all integration fixes
- Build successfully standalone
- Be ready for a clean PR to main

Remember: This state ensures PRs will merge without conflicts!
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
