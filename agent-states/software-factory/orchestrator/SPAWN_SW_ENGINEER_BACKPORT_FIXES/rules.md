# Orchestrator - SPAWN_SW_ENGINEER_BACKPORT_FIXES State Rules

## ✅ BACKPORTS ARE NORMAL - NO USER REVIEW REQUIRED

**IMPORTANT:** Backports are NORMAL operations. While R322 Part A requires stopping AFTER spawning (for context preservation), there is NO requirement to stop BEFORE spawning backport fixes. Backports should proceed automatically!

# PRIMARY DIRECTIVES

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R232** - Line Counting Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R232-enforcement-examples.md`

3. **R220** - Size Limit Compliance
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R220-atomic-pr-design-requirement.md`

4. **R256** - Fix Planning Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R256-mandatory-phase-assessment-gate.md`

5. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

6. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`

7. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`

8. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

9. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-mandatory-line-counter-auto-detection.md`


# 🔴🔴🔴 MANDATORY: R322 STOP + R405 CONTINUATION FLAG 🔴🔴🔴

**CRITICAL FOR SPAWN STATES - READ THIS FIRST OR FAIL TEST 2!**

## 🚨 THE PATTERN THAT FAILED TEST 2 🚨

**WHAT HAPPENED IN TEST 2:**
- Orchestrator spawned Code Reviewers ✅ (correct)
- Orchestrator stopped per R322 ✅ (correct)
- Orchestrator **DID NOT emit `CONTINUE-SOFTWARE-FACTORY=TRUE`** ❌ (WRONG!)
- Test framework saw no continuation flag → stopped automation
- Test 2 FAILED at iteration 8

**ROOT CAUSE:** Confusion between R322 "stop" and R405 continuation flag

## 🔴 CRITICAL DISTINCTION: TWO INDEPENDENT DECISIONS 🔴

### Decision 1: Should Agent Stop? (R322 - Context Preservation)
**YES - ALWAYS stop after spawning for context preservation**

- **Purpose**: Prevent context overflow between states
- **Action**: `exit 0` to end conversation
- **User Experience**: User sees "/continue-orchestrating" as next step
- **This is NORMAL!** Not an error!

### Decision 2: Should Factory Continue? (R405 - Automation Control)
**YES - ALWAYS emit TRUE for normal spawning operations**

- **Purpose**: Tell automation whether it CAN restart
- **Action**: `echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"` (LAST output before exit)
- **Automation**: Framework will auto-restart orchestrator
- **This is NORMAL!** Designed behavior!

## ✅ REQUIRED PATTERN FOR ALL SPAWN STATES

```bash
# 1. Complete spawning work
echo "✅ Spawned [agent type] for [purpose]"

# 2. Update state file per R324/R288
update_state "[NEXT_STATE]"
commit_state_files_per_r288()

# 3. Save TODOs per R287
save_todos "SPAWNED_[AGENT_TYPE]"

# 4. R322: Stop conversation (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# 5. R405: CONTINUATION FLAG - MUST BE TRUE FOR SPAWNING!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Exit to end conversation
exit 0
```

## ❌ WRONG PATTERN (CAUSES TEST FAILURES)

```bash
# ❌ THIS KILLS AUTOMATION - DO NOT DO THIS!
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
exit 0

# Result: Test framework stops, Test 2 fails at iteration 8
```

## 🎯 WHY TRUE IS CORRECT FOR SPAWNING

**Spawning is NORMAL operation:**
- ✅ System knows next state (from state machine)
- ✅ Automation can continue (designed workflow)
- ✅ No manual intervention needed
- ✅ Context preservation ≠ error condition

**The orchestrator stopping (`exit 0`) is for:**
- Preserving context between conversation turns
- Allowing state file commits
- Creating clean state boundaries

**The TRUE flag indicates:**
- Automation CAN restart the conversation
- System knows what to do next (check state file)
- Normal operation is proceeding

## 🔴 WHEN TO USE FALSE (NOT FOR SPAWNING!)

**FALSE should ONLY be used for catastrophic failures:**
- ❌ State file corrupted beyond parsing
- ❌ Critical infrastructure destroyed
- ❌ Unrecoverable system errors
- ❌ **NEVER for normal spawning operations!**

## 📋 SPAWN STATE CHECKLIST

**Before exiting this spawn state, verify:**
1. [ ] All agents spawned successfully
2. [ ] State file updated to next state per R324
3. [ ] State files committed per R288
4. [ ] TODOs saved per R287
5. [ ] R322 stop message displayed
6. [ ] **CONTINUE-SOFTWARE-FACTORY=TRUE emitted** ← Critical!
7. [ ] Exited with `exit 0`

**Missing step 6 = Test 2 failure = -100% grade**

---


## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state-v3.json with new state
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

**YOU HAVE ENTERED SPAWN_SW_ENGINEER_BACKPORT_FIXES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_SPAWN_SW_ENGINEER_BACKPORT_FIXES-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for SPAWN_SW_ENGINEER_BACKPORT_FIXES" > "markers/state-verification/state_rules_read_orchestrator_SPAWN_SW_ENGINEER_BACKPORT_FIXES-$(date +%Y%m%d-%H%M%S)"
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

## 🎯 STATE OBJECTIVES - SPAWN SW ENGINEERS FOR BACKPORT IMPLEMENTATION

In the SPAWN_SW_ENGINEER_BACKPORT_FIXES state, the ORCHESTRATOR is responsible for:

1. **Reading the Backport Plan**
   - Load BACKPORT-PLAN.md created by Code Reviewer
   - Identify all effort branches needing fixes
   - Understand fix groupings and dependencies
   - Note any sequential requirements

2. **Preparing Individual SW Engineer Assignments**
   - Extract fixes for each effort from the plan
   - Create clear instructions for each SW Engineer
   - Specify exact working directories and branches
   - Include verification requirements

3. **Spawning SW Engineers (Parallel or Sequential)**
   - Analyze if fixes can be done in parallel
   - Spawn SW Engineers according to plan requirements
   - If parallel: spawn all within 5 seconds (R151)
   - If sequential: document why and spawn in order

4. **Stopping After Spawn (R322 Part A)**
   - Update orchestrator-state-v3.json to MONITORING_BACKPORT_PROGRESS
   - Document all spawned engineers
   - Commit and push state changes
   - STOP and wait for user continuation

## 📝 REQUIRED ACTIONS

### Step 1: Load and Parse Backport Plan
```bash
# Read the backport plan created by Code Reviewer
cd /efforts/integration-testing

# Load the plan
if [ -f "BACKPORT-PLAN.md" ]; then
    echo "📋 Loading backport plan from Code Reviewer..."
    cat BACKPORT-PLAN.md
    
    # Parse out effort branches needing fixes
    echo "📊 Identifying efforts requiring backports..."
    grep -E "^##.*Branch:|Effort:" BACKPORT-PLAN.md
else
    echo "❌ CRITICAL: No BACKPORT-PLAN.md found!"
    echo "Cannot proceed without Code Reviewer's plan"
    exit 1
fi
```

### Step 2: Create Individual SW Engineer Instructions
```bash
# For each effort in the backport plan, create specific instructions

# Example for each effort (this would be done for ALL efforts in plan)
EFFORTS_FROM_PLAN=("effort-1" "effort-2" "effort-3")  # Parsed from plan

for EFFORT_NAME in "${EFFORTS_FROM_PLAN[@]}"; do
    INSTRUCTION_FILE="/efforts/${EFFORT_NAME}/BACKPORT-INSTRUCTIONS.md"
    
    echo "📝 Creating instructions for ${EFFORT_NAME}..."
    
    cat > "$INSTRUCTION_FILE" << 'EOF'
# Backport Implementation Instructions

## Your Assignment
- Effort: ${EFFORT_NAME}
- Working Directory: /efforts/${EFFORT_NAME}
- Branch: [branch-name-from-plan]

## Fixes to Apply (From Code Reviewer's Plan)
[Extract specific fixes for this effort from BACKPORT-PLAN.md]

### Fix 1: [Description]
- Files: [List files]
- Changes: [Specific changes needed]
- Verification: [How to verify]

### Fix 2: [Description]
- Files: [List files]
- Changes: [Specific changes needed]
- Verification: [How to verify]

## Implementation Process
1. Ensure you're on your effort branch
2. Apply fixes as specified above
3. Run build to verify compilation
4. Run tests to ensure nothing breaks
5. Commit with clear message referencing integration fixes
6. Push your updated branch
7. Update sw-engineer-state.yaml when complete

## Success Criteria
- ✅ All listed fixes applied exactly as specified
- ✅ Build succeeds after fixes
- ✅ All tests pass
- ✅ Branch pushed to remote
- ✅ State file shows BACKPORT_COMPLETE

## Important Notes
- Do NOT modify anything beyond the specified fixes
- Do NOT merge from integration branch
- Apply fixes directly to your effort branch
- If you encounter issues, document them and update state to BLOCKED
EOF
    
    echo "✅ Instructions created for ${EFFORT_NAME}"
done
```

### Step 3: Determine Parallelization Strategy
```bash
# Analyze if fixes can be parallelized
echo "🔍 Analyzing parallelization requirements..."

# Check backport plan for dependencies
if grep -q "SEQUENTIAL\|DEPENDENCY\|ORDER" BACKPORT-PLAN.md; then
    echo "⚠️ Sequential execution required - dependencies detected"
    PARALLEL_EXECUTION=false
else
    echo "✅ Parallel execution possible - no dependencies"
    PARALLEL_EXECUTION=true
fi

# Document the strategy
cat > BACKPORT-SPAWN-STRATEGY.md << 'EOF'
# Backport SW Engineer Spawn Strategy

## Execution Mode: ${PARALLEL_EXECUTION}

## Efforts to Spawn:
1. effort-1 - 3 fixes - no dependencies
2. effort-2 - 2 fixes - no dependencies  
3. effort-3 - 4 fixes - no dependencies

## Spawn Order:
$(if [ "$PARALLEL_EXECUTION" = true ]; then
    echo "All engineers spawned simultaneously (R151 compliance)"
else
    echo "Sequential spawn order based on dependencies"
fi)

## Timing Requirements (R151):
- All parallel spawns must occur within 5 seconds
- First spawn timestamp: [will be recorded]
- Last spawn timestamp: [will be recorded]
- Delta must be < 5 seconds for parallel execution
EOF
```

### Step 4: Spawn SW Engineers
```bash
echo "🚀 Beginning SW Engineer spawn sequence..."

# Record start time for R151 compliance
START_TIME=$(date +%s)
echo "Spawn sequence started at: $(date '+%Y-%m-%d %H:%M:%S')"

# Spawn engineers based on strategy
if [ "$PARALLEL_EXECUTION" = true ]; then
    echo "📊 Spawning ALL SW Engineers in parallel (R151 compliance)..."
    
    for EFFORT_NAME in "${EFFORTS_FROM_PLAN[@]}"; do
        echo "🚀 Spawning SW Engineer for ${EFFORT_NAME}..."
        
        # Create spawn command
        cat > /tmp/spawn-sw-${EFFORT_NAME}.md << 'EOF'
@agent-software-engineer

## BACKPORT IMPLEMENTATION ASSIGNMENT

You are being spawned to implement backport fixes for your effort branch.

### Immediate Actions Required:
1. Read your instructions at: /efforts/${EFFORT_NAME}/BACKPORT-INSTRUCTIONS.md
2. Navigate to your working directory: /efforts/${EFFORT_NAME}
3. Verify you're on the correct branch
4. Apply ALL fixes specified in your instructions
5. Verify build and tests pass
6. Commit and push changes
7. Update your state file to BACKPORT_COMPLETE

### Critical Requirements:
- Apply fixes EXACTLY as specified
- Do NOT merge from other branches
- Do NOT modify beyond specified fixes
- Must complete ALL fixes or report BLOCKED

### State Progression:
INIT → BACKPORT_IMPLEMENTATION → BACKPORT_COMPLETE

Start immediately - this is time-critical per R321.
EOF
        
        # Log spawn with timestamp
        echo "$(date +%s): Spawned SW Engineer for ${EFFORT_NAME}" >> SPAWN-LOG.md
    done
    
else
    echo "📊 Spawning SW Engineers sequentially due to dependencies..."
    
    # Sequential spawn with order from plan
    for EFFORT_NAME in "${EFFORTS_FROM_PLAN[@]}"; do
        echo "🚀 Spawning SW Engineer for ${EFFORT_NAME} (sequential)..."
        # [Same spawn command as above]
        echo "⏸️ Waiting for ${EFFORT_NAME} to complete before next spawn..."
        break  # Would actually wait in real implementation
    done
fi

# Record end time and verify R151 compliance
END_TIME=$(date +%s)
SPAWN_DURATION=$((END_TIME - START_TIME))

echo "Spawn sequence completed at: $(date '+%Y-%m-%d %H:%M:%S')"
echo "Total spawn duration: ${SPAWN_DURATION} seconds"

if [ "$PARALLEL_EXECUTION" = true ] && [ $SPAWN_DURATION -gt 5 ]; then
    echo "⚠️ WARNING: R151 VIOLATION - Spawn duration exceeded 5 seconds!"
    echo "Duration was ${SPAWN_DURATION} seconds"
fi
```

### Step 5: Update State and STOP (R322 Part A Enforcement)
```bash
# Update orchestrator state
cd $CLAUDE_PROJECT_DIR

# Create comprehensive state update
cat > orchestrator-state-v3.json << 'EOF'
current_state: MONITORING_BACKPORT_PROGRESS
previous_state: SPAWN_SW_ENGINEER_BACKPORT_FIXES
backport_status: IN_PROGRESS
agents_spawned:
  - agent: sw-engineer-1
    effort: effort-1
    purpose: Implement backport fixes
    timestamp: $(date +%s)
    state: BACKPORT_IMPLEMENTATION
  - agent: sw-engineer-2
    effort: effort-2
    purpose: Implement backport fixes
    timestamp: $(date +%s)
    state: BACKPORT_IMPLEMENTATION
  - agent: sw-engineer-3
    effort: effort-3
    purpose: Implement backport fixes
    timestamp: $(date +%s)
    state: BACKPORT_IMPLEMENTATION
parallelization:
  mode: parallel
  r151_compliant: true
  spawn_duration: ${SPAWN_DURATION}
backport_plan: /efforts/integration-testing/BACKPORT-PLAN.md
monitoring:
  total_efforts: 3
  completed: 0
  in_progress: 3
  blocked: 0
EOF

# Commit state change
git add orchestrator-state-v3.json
git commit -m "state: transition to MONITORING_BACKPORT_PROGRESS after spawning SW Engineers

- Spawned SW Engineers for all efforts requiring backports
- Using backport plan from Code Reviewer
- R151 compliant parallel spawn
- R322 Part A stopping after spawn"
git push

echo "✅ State updated and committed"
echo "🛑 STOPPING per R322 Part A - Must stop after spawn state"
```

## ⚠️ CRITICAL REQUIREMENTS

### R151 Compliance for Parallel Spawns
- If spawning multiple engineers: ALL within 5 seconds
- Document spawn timestamps
- Flag any R151 violations

### Clear Separation of Work
- **Code Reviewer**: Created the plan (already done)
- **Orchestrator**: Distributes work and spawns engineers
- **SW Engineers**: Actually implement the fixes

### No Direct Implementation
- Orchestrator MUST NOT apply any fixes
- Orchestrator MUST NOT use git commands on code
- Orchestrator ONLY coordinates

### R322 Part A Enforcement
- MUST stop after spawning engineers
- MUST update state before stopping
- MUST NOT continue automatically

## 🚫 FORBIDDEN ACTIONS

1. **Implementing any fixes directly** - R006 violation
2. **Using cherry-pick or applying patches** - SW Engineer's job
3. **Continuing without stopping** - R322 Part A violation
4. **Modifying the backport plan** - Use it as-is from Code Reviewer
5. **Spawning agents slowly in parallel mode** - R151 violation

## ✅ PROJECT_DONE CRITERIA

Before transitioning to MONITORING_BACKPORT_PROGRESS:
- [ ] Backport plan loaded from Code Reviewer
- [ ] Instructions created for each SW Engineer
- [ ] Parallelization strategy determined
- [ ] All SW Engineers spawned appropriately
- [ ] R151 compliance verified if parallel
- [ ] State updated to MONITORING_BACKPORT_PROGRESS
- [ ] State changes committed and pushed
- [ ] STOPPED per R322 Part A

## 🔄 STATE TRANSITIONS

### Success Path:
```
SPAWN_SW_ENGINEER_BACKPORT_FIXES → MONITORING_BACKPORT_PROGRESS
```
- All engineers spawned
- Monitoring their progress
- Waiting for completion

### From Monitoring:
```
MONITORING_BACKPORT_PROGRESS → PR_PLAN_CREATION
```
- All backports complete
- Ready for PR planning

```
MONITORING_BACKPORT_PROGRESS → ERROR_RECOVERY
```
- Backports failed
- Need intervention

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **R151 Compliance** (35%)
   - Parallel spawns within 5 seconds
   - Proper timestamp documentation
   
2. **R322 Part A Compliance** (25%)
   - Stop after spawning
   - No automatic continuation
   
3. **Work Distribution** (25%)
   - Clear instructions for each engineer
   - Following Code Reviewer's plan exactly
   
4. **State Management** (15%)
   - Proper state transitions
   - Complete state documentation

## 💡 TIPS FOR PROJECT_DONE

1. **Trust the plan** - Code Reviewer analyzed the fixes, use their plan
2. **Spawn quickly** - For parallel execution, speed matters (R151)
3. **Clear instructions** - Each engineer needs to know exactly what to do
4. **Stop means stop** - R322 Part A is non-negotiable

Remember: You're the COORDINATOR, not the implementer!

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state-v3.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
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



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_SW_ENGINEER_BACKPORT_FIXES:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_SW_ENGINEER_BACKPORT_FIXES complete - [describe what was accomplished]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Update State File (R288 - SUPREME LAW)
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "SPAWN_SW_ENGINEER_BACKPORT_FIXES",
  "work_accomplished": [
    "[List accomplishments from state work]"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "SPAWN_SW_ENGINEER_BACKPORT_FIXES" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Commit and push state files
# 4. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Validate State File (R324)
```bash
# Validate state file before committing
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ State file validated"
```

---

### ✅ Step 5: Commit State File (R288)
```bash
# Commit and push state file immediately
git add orchestrator-state-v3.json

if ! git commit -m "state: SPAWN_SW_ENGINEER_BACKPORT_FIXES → $NEXT_STATE - SPAWN_SW_ENGINEER_BACKPORT_FIXES complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: SPAWN_SW_ENGINEER_BACKPORT_FIXES"
    echo "Attempted transition from: SPAWN_SW_ENGINEER_BACKPORT_FIXES"
    echo ""
    echo "Common causes:"
    echo "  - Schema validation failure (check pre-commit hook output above)"
    echo "  - Missing required fields in JSON files"
    echo "  - Invalid JSON syntax"
    echo ""
    echo "🛑 Cannot proceed - manual intervention required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
    exit 1
fi

git push || echo "⚠️ WARNING: Push failed - committed locally"
echo "✅ State file committed and pushed"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_SW_ENGINEER_BACKPORT_FIXES_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_SW_ENGINEER_BACKPORT_FIXES complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 8: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No state update = state machine broken (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: No commit = state lost on compaction (R288 violation, -100%)
- Missing Step 6: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 7: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 8: No exit = R322 violation (-100%)

**ALL 8 STEPS ARE MANDATORY - NO EXCEPTIONS**

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="MONITORING_BACKPORT_PROGRESS"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 SPAWN STATE PATTERN - R322 + R405 USAGE 🚨

**Spawning operations require R322 stop for context preservation:**
```bash
# After spawning agent(s)
echo "✅ Spawned agents for work"

# R322 checkpoint (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# Flag? → MUST BE TRUE (normal operation!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# Stop inference
exit 0
```

**Why TRUE is correct:**
- Spawning is NORMAL operation
- System knows next state
- Automation can continue
- **Context preservation ≠ manual intervention needed!**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

