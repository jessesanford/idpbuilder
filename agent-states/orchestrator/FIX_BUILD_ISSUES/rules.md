# Orchestrator - FIX_BUILD_ISSUES State Rules

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

**YOU HAVE ENTERED FIX_BUILD_ISSUES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_FIX_BUILD_ISSUES
echo "$(date +%s) - Rules read and acknowledged for FIX_BUILD_ISSUES" > .state_rules_read_orchestrator_FIX_BUILD_ISSUES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY FIXING WORK UNTIL RULES ARE READ:
- ❌ Start fixing code
- ❌ Spawn engineers to fix
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

**⚠️ R006 WARNING FOR FIX_BUILD_ISSUES STATE:**
- DO NOT fix code issues yourself!
- DO NOT edit any .go, .py, .js, .ts, or other code files!
- DO NOT use cherry-pick, apply patches, or make direct fixes!
- Spawn SW Engineers for ALL code modifications
- You are a COORDINATOR ONLY - never a developer

### 🚨🚨🚨 RULE R151 - Parallel Agent Timestamp Requirement [BLOCKING]
**Orchestrator must ensure parallel agents start within 5 seconds**

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`

The orchestrator is a coordinator only and must spawn SW engineers for all code fixes.

### Build Issue Fix Requirements

**Every fix MUST be tracked for backporting**

ALL fixes must be documented with their original effort branch for backporting.

## 🎯 STATE OBJECTIVES

In the FIX_BUILD_ISSUES state, you are responsible for:

1. **Analyzing Build Failures**
   - Read BUILD-VALIDATION-REPORT.md
   - Read BUILD-ERRORS.txt
   - Identify root causes
   - Categorize by effort/component

2. **Creating Fix Plans**
   - Map each error to its source effort
   - Create specific fix instructions
   - Identify dependencies between fixes
   - Determine fix order

3. **Spawning Engineers for Fixes**
   - **CRITICAL**: Orchestrator NEVER fixes code directly
   - Spawn SW Engineers with specific fix instructions
   - One engineer per effort needing fixes
   - Provide clear fix requirements

4. **Tracking for Backports**
   - Document EVERY fix with its effort
   - Create backport manifest
   - Track which branches need updates
   - Prepare for BACKPORT_FIXES state

## 📝 REQUIRED ACTIONS

### Step 1: Analyze Build Failures
```bash
cd /efforts/integration-testing

# Read build validation report
if [ -f "BUILD-VALIDATION-REPORT.md" ]; then
    echo "📋 Reading build validation report..."
    cat BUILD-VALIDATION-REPORT.md
else
    echo "❌ No build validation report found!"
    exit 1
fi

# Read build errors
if [ -f "BUILD-ERRORS.txt" ]; then
    echo "📋 Analyzing build errors..."
    cat BUILD-ERRORS.txt
else
    echo "⚠️ No specific error file, checking build output..."
    grep -i "error" BUILD-OUTPUT.log | head -20
fi
```

### Step 2: Map Errors to Efforts
```bash
# Create error-to-effort mapping
cat > ERROR-TO-EFFORT-MAP.md << 'EOF'
# Error to Effort Mapping
Date: $(date)

## Compilation Errors

### Error 1: [Error Description]
- File: [path/to/file]
- Original Effort: [effort-name]
- Branch: [effort-branch-name]
- Fix Required: [Description of fix]

### Error 2: [Error Description]
- File: [path/to/file]
- Original Effort: [effort-name]
- Branch: [effort-branch-name]
- Fix Required: [Description of fix]

## Link Errors

### Error 1: [Error Description]
- Component: [component name]
- Original Effort: [effort-name]
- Branch: [effort-branch-name]
- Fix Required: [Description of fix]

## Missing Dependencies

### Dependency: [name]
- Required By: [component]
- Original Effort: [effort-name]
- Branch: [effort-branch-name]
- Fix Required: Add dependency to [file]

## Fix Order
1. [First fix - usually dependencies]
2. [Second fix - usually compilation]
3. [Third fix - usually linking]
EOF
```

### Step 3: Create Fix Instructions
```bash
# For each effort needing fixes, create instructions
for effort in $(grep "Original Effort:" ERROR-TO-EFFORT-MAP.md | cut -d: -f2 | sort -u); do
    effort_trim=$(echo $effort | xargs)
    
    cat > "FIX-INSTRUCTIONS-${effort_trim}.md" << EOF
# Fix Instructions for $effort_trim
Date: $(date)

## Issues to Fix

### Issue 1: [Specific issue]
**File**: [path/to/file]
**Line**: [line number if known]
**Current Code**:
\`\`\`
[problematic code]
\`\`\`
**Required Fix**:
\`\`\`
[corrected code]
\`\`\`

### Issue 2: [Specific issue]
**File**: [path/to/file]
**Required Action**: [Add import/dependency/etc]

## Testing Requirements
After fixes:
1. Run build command: [specific command]
2. Verify no errors for this component
3. Run tests if applicable

## Backport Note
**CRITICAL**: After fixes are verified, these changes MUST be backported to branch: [branch-name]
EOF
done
```

### Step 4: Create Backport Manifest
```bash
# CRITICAL: Track all fixes for backporting
cat > BACKPORT-MANIFEST.md << 'EOF'
# Backport Manifest
Date: $(date)
State: FIX_BUILD_ISSUES

## Fixes Requiring Backport

### Effort: [effort-name-1]
- Branch: [original-branch]
- Files Modified:
  - [file1.ext] - [type of change]
  - [file2.ext] - [type of change]
- Commits to Cherry-pick: [will be added after fixes]

### Effort: [effort-name-2]
- Branch: [original-branch]
- Files Modified:
  - [file1.ext] - [type of change]
- Commits to Cherry-pick: [will be added after fixes]

## Backport Execution Plan
1. After all fixes verified in integration-testing
2. For each effort branch:
   a. Checkout effort branch
   b. Cherry-pick or apply fixes
   c. Test build in isolation
   d. Commit with reference to integration fix
   e. Push to remote

## Verification
- [ ] All fixes documented
- [ ] All branches identified
- [ ] Fix commits recorded
- [ ] Ready for BACKPORT_FIXES state
EOF
```

### Step 5: Spawn Engineers for Fixes
```bash
# For each effort needing fixes, spawn an engineer
for fix_file in FIX-INSTRUCTIONS-*.md; do
    if [ -f "$fix_file" ]; then
        effort_name=$(basename "$fix_file" .md | sed 's/FIX-INSTRUCTIONS-//')
        
        echo "🔧 Spawning SW Engineer to fix: $effort_name"
        echo "Instructions file: $fix_file"
        
        # Record spawn in state file
        cat >> spawned-engineers.txt << EOF
$(date): Spawned engineer for $effort_name
- Instruction File: $fix_file
- Working Directory: /efforts/integration-testing
- Task: Fix build issues for $effort_name component
EOF
        
        # ACTUAL SPAWN COMMAND WOULD GO HERE
        # Example: spawn_sw_engineer --fix-mode --instructions "$fix_file"
    fi
done
```

### Step 6: Monitor Fix Progress
```bash
# Create monitoring checklist
cat > FIX-PROGRESS-TRACKER.md << 'EOF'
# Fix Progress Tracker
Date: $(date)

## Engineers Spawned
- [ ] Engineer 1: [effort-name] - Status: [IN_PROGRESS/COMPLETE]
- [ ] Engineer 2: [effort-name] - Status: [IN_PROGRESS/COMPLETE]

## Fixes Applied
- [ ] Compilation errors resolved
- [ ] Link errors resolved
- [ ] Dependencies added
- [ ] Build re-tested

## Verification Status
- [ ] Build completes without errors
- [ ] All tests pass
- [ ] Backport manifest updated with commits

## Next State
Ready for: [BACKPORT_FIXES/BUILD_VALIDATION]
EOF
```

## ⚠️ CRITICAL REQUIREMENTS

### 🔴 R006 VIOLATION DETECTION 
If you find yourself about to:
- Use Edit/Write tools on any code file
- Run `git cherry-pick` or `git apply`
- Fix compilation errors directly
- Modify any .go, .py, .js, .ts files
- Apply patches yourself
- Make excuses like "it's simpler to just fix it"

**STOP IMMEDIATELY - You are violating R006!**
Spawn SW Engineers instead - NO EXCEPTIONS!

### Orchestrator NEVER Writes Code
- **VIOLATION = AUTOMATIC FAILURE**
- All code fixes MUST be done by SW Engineers
- Orchestrator only coordinates and tracks

### Complete Backport Documentation
**Every single fix MUST have:**
1. Original effort/branch identified
2. Files modified documented
3. Fix commit hash recorded
4. Backport plan created

### Fix Verification Before Backport
- Fixes must work in integration-testing first
- Build must succeed completely
- Only then proceed to backporting

## 🚫 FORBIDDEN ACTIONS

1. **NEVER have orchestrator write code directly**
2. **NEVER skip backport documentation**
3. **NEVER proceed without verifying fixes work**
4. **NEVER merge fixes without tracking**
5. **NEVER modify effort branches directly**

## ✅ SUCCESS CRITERIA

Before transitioning to next state:
- [ ] All build errors analyzed
- [ ] Fix instructions created for each issue
- [ ] Engineers spawned for all fixes
- [ ] Fixes verified to work
- [ ] Complete backport manifest created
- [ ] All changes committed and pushed

## 🔄 STATE TRANSITIONS

### To Backport State:
```
FIX_BUILD_ISSUES → BACKPORT_FIXES
```
- All fixes complete and verified
- Backport manifest ready
- Need to apply fixes to original branches

### Re-validation Path:
```
FIX_BUILD_ISSUES → BUILD_VALIDATION
```
- Fixes applied in integration branch
- Need to re-run build validation
- Backporting will happen after validation

### Complex Issues Path:
```
FIX_BUILD_ISSUES → ERROR_RECOVERY
```
- Issues too complex to fix
- Structural problems found
- Need architectural intervention

## 📊 VERIFICATION CHECKLIST

Before leaving this state:
```bash
# Verify fix instructions created
ls -la FIX-INSTRUCTIONS-*.md

# Verify backport manifest exists
ls -la BACKPORT-MANIFEST.md

# Verify engineers spawned
cat spawned-engineers.txt

# Check fix progress
cat FIX-PROGRESS-TRACKER.md

# Verify build now works
cd /efforts/integration-testing
[run build command]

# Commit all documentation
git add -A
git commit -m "fix: build issues resolved - ready for backporting"
git push
```

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Proper Delegation** (30%)
   - Engineers spawned for fixes
   - Orchestrator didn't write code
   
2. **Fix Quality** (25%)
   - Clear instructions provided
   - All issues addressed
   
3. **Backport Documentation** (25%)
   - Every fix tracked
   - Branches identified
   
4. **Verification** (20%)
   - Fixes tested
   - Build succeeds

## 💡 TIPS FOR SUCCESS

1. **Analyze Thoroughly**: Understand root causes
2. **Document Everything**: Backporting depends on it
3. **Test Each Fix**: Verify before proceeding
4. **Think Dependencies**: Fix in correct order

## 🚨 COMMON PITFALLS TO AVOID

1. **Orchestrator Coding**: Instant failure
2. **Missing Backport Info**: Can't complete PRs
3. **Partial Fixes**: Build still fails
4. **Wrong Fix Order**: Creates more errors

Remember: Every fix here must make it back to the original branches!
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
