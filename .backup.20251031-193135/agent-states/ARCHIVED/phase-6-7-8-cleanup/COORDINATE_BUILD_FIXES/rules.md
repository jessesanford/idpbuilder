# Orchestrator - FIX_WAVE_UPSTREAM_BUGS State Rules


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

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED FIX_WAVE_UPSTREAM_BUGS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_FIX_WAVE_UPSTREAM_BUGS
echo "$(date +%s) - Rules read and acknowledged for FIX_WAVE_UPSTREAM_BUGS" > .state_rules_read_orchestrator_FIX_WAVE_UPSTREAM_BUGS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR FIX_WAVE_UPSTREAM_BUGS STATE

### Core Mandatory Rules

### 🚨🚨🚨 R006 - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Automatic termination, 0% grade
**Summary**: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents
**CRITICAL**: Even simple fixes must be delegated to SW Engineers!

### 🚨🚨🚨 R319 - ORCHESTRATOR NEVER MEASURES CODE (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality**: BLOCKING - Orchestrator MUST NOT use line-counter.sh
**Summary**: Code Reviewers measure code size, NOT orchestrators

### ⚠️⚠️⚠️ R317 - Working Directory Restrictions (WARNING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R317-working-directory-restrictions.md`
**Criticality**: WARNING - -25% for violations
**Summary**: MUST NOT enter agent working directories - operate from root only

### 🚨🚨🚨 R318 - Agent Failure Escalation Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R318-agent-failure-escalation-protocol.md`
**Criticality**: BLOCKING - -40% for attempting forbidden fixes
**Summary**: NEVER fix agent failures directly - respawn with better instructions

### 🚨🚨🚨 R287 - TODO Save Frequency Requirements (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
**Criticality**: BLOCKING - Save every 15 minutes/10 messages
**Summary**: Mandatory TODO saves during coordination

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state-v3.json on all state changes

### State-Specific Rules

### 🚨🚨🚨 R151 - Parallel Agent Timestamp Requirement (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
**Criticality**: BLOCKING - Parallel agents must start within 5 seconds
**Summary**: When spawning multiple SW Engineers, ensure synchronized start

### 🔴🔴🔴 R321 - Immediate Backport During Integration (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: SUPREME LAW - Integration fixes must be backported immediately
**Summary**: Track all fixes for immediate backporting to source branches

## 🔴🔴🔴 CRITICAL: FIX_WAVE_UPSTREAM_BUGS IS A VERB - START COORDINATING NOW! 🔴🔴🔴

**FIX_WAVE_UPSTREAM_BUGS MEANS ACTIVELY COORDINATING FIX EFFORTS RIGHT NOW!**
- ❌ NOT "I'm in coordinate build fixes state"  
- ❌ NOT "Ready to coordinate fixes"
- ✅ YES "I'm distributing fix plans to engineers NOW"
- ✅ YES "I'm spawning SW Engineers for fixes NOW"
- ✅ YES "I'm tracking fix assignments NOW"

## State Context
FIX_WAVE_UPSTREAM_BUGS = You ARE ACTIVELY coordinating the distribution and execution of build fixes THIS INSTANT. You have fix plans from Code Reviewer and must spawn SW Engineers to implement them.

## 🎯 STATE OBJECTIVES

In the FIX_WAVE_UPSTREAM_BUGS state, you are responsible for:

1. **Reading Fix Plans**
   - Load FIX-PLAN-*.md files from Code Reviewer
   - Understand fix requirements for each effort
   - Identify dependencies between fixes
   - Determine parallelization opportunities

2. **Distributing Fix Work**
   - Create fix instructions for each SW Engineer
   - Assign specific fixes to specific engineers
   - Provide clear success criteria
   - Set up tracking mechanisms

3. **Spawning SW Engineers**
   - Spawn engineers with specific fix instructions
   - Ensure proper working directories
   - Coordinate parallel vs sequential work
   - Track R151 timestamp requirements

4. **Setting Up Monitoring**
   - Create fix progress tracking
   - Document backport requirements (R321)
   - Prepare for verification
   - Set up success metrics

## 📝 REQUIRED ACTIONS

### Step 1: Load Fix Plans from Code Reviewer
```bash
cd /efforts/integration-testing

# Find all fix plans created by Code Reviewer
echo "📋 Loading fix plans..."
for fix_plan in FIX-PLAN-*.md; do
    if [ -f "$fix_plan" ]; then
        echo "Found: $fix_plan"
        effort_name=$(echo "$fix_plan" | sed 's/FIX-PLAN-//;s/.md//')
        echo "  Effort: $effort_name"
        
        # Read key information
        grep -A5 "## Required Fixes" "$fix_plan"
    fi
done
```

### Step 2: Create Fix Assignment Matrix
```bash
# Document who fixes what
cat > FIX-ASSIGNMENT-MATRIX.md << 'EOF'
# Fix Assignment Matrix
Date: $(date)
State: FIX_WAVE_UPSTREAM_BUGS

## Assignment Overview
| SW Engineer | Effort | Fix Type | Priority | Dependencies |
|-------------|--------|----------|----------|--------------|
| SWE-1 | [effort-1] | Compilation | HIGH | None |
| SWE-2 | [effort-2] | Dependencies | HIGH | None |
| SWE-3 | [effort-3] | Linking | MEDIUM | SWE-1 |

## Parallelization Strategy
### Parallel Group 1 (Independent fixes)
- SWE-1: Fix compilation in effort-1
- SWE-2: Add dependencies in effort-2

### Sequential Group (Dependent fixes)
- After SWE-1 completes → SWE-3 starts

## Timing Requirements (R151)
For parallel spawns:
- Maximum spawn delay: 5 seconds
- All parallel agents must emit timestamps
- Verify synchronization

## Success Criteria per Engineer
### SWE-1
- [ ] All compilation errors in effort-1 resolved
- [ ] Build succeeds for component
- [ ] Tests pass

### SWE-2
- [ ] Missing dependencies added
- [ ] Package versions correct
- [ ] No version conflicts

### SWE-3
- [ ] Linking errors resolved
- [ ] Symbols properly exported
- [ ] Integration tests pass
EOF
```

### Step 3: Create Individual Fix Instructions
```bash
# For each SW Engineer, create specific instructions
for effort in effort-1 effort-2 effort-3; do
    cat > "FIX-INSTRUCTIONS-SWE-${effort}.md" << EOF
# Fix Instructions for SW Engineer - ${effort}
Date: $(date)
Assigned by: orchestrator/FIX_WAVE_UPSTREAM_BUGS

## Your Assignment
Fix build failures in ${effort} component.

## Working Directory
/efforts/integration-testing

## Fix Plan Reference
See: FIX-PLAN-${effort}.md for detailed changes

## Specific Tasks
1. [Task 1 from fix plan]
2. [Task 2 from fix plan]
3. [Task 3 from fix plan]

## Files to Modify
- path/to/file1.go - [specific change]
- path/to/file2.go - [specific change]

## Testing Requirements
After fixes:
1. Run: go build ./...
2. Run: go test ./...
3. Verify no compilation errors
4. Create FIX-COMPLETE.marker when done

## Backport Tracking (R321)
**CRITICAL**: Document all changes for backporting:
- Original branch: [branch-name]
- Files modified: [track all]
- Commits: [record hashes]

## Success Indicators
- Build completes without errors
- Tests pass
- No new issues introduced
EOF
done
```

### Step 4: Create Backport Manifest (R321)
```bash
# CRITICAL: Track all fixes for immediate backporting
cat > BACKPORT-MANIFEST.md << 'EOF'
# Backport Manifest
Date: $(date)
State: FIX_WAVE_UPSTREAM_BUGS
Integration Context: ACTIVE

## 🔴🔴🔴 R321 ENFORCEMENT 🔴🔴🔴
ALL fixes MUST be immediately backported to source branches!

## Fixes Requiring Backport

### Effort: [effort-1]
- Source Branch: [original-branch]
- SW Engineer: SWE-1
- Fix Status: IN_PROGRESS
- Files to Backport:
  - [ ] file1.go
  - [ ] file2.go
- Backport Status: PENDING

### Effort: [effort-2]
- Source Branch: [original-branch]
- SW Engineer: SWE-2
- Fix Status: IN_PROGRESS
- Files to Backport:
  - [ ] package.json
  - [ ] dependencies.go
- Backport Status: PENDING

## Backport Execution Plan
1. After each fix completes in integration
2. Immediately checkout source branch
3. Apply same fixes to source
4. Test in isolation
5. Commit with integration reference
6. Push to remote

## Tracking
- [ ] All fixes documented
- [ ] Source branches identified
- [ ] Fix commits recorded
- [ ] Ready for IMMEDIATE_BACKPORT_REQUIRED
EOF
```

### Step 5: Spawn SW Engineers
```bash
# Spawn engineers based on parallelization strategy

# PARALLEL GROUP 1 - Independent fixes (R151: must be within 5 seconds)
echo "🚀 Spawning parallel SW Engineers..."
START_TIME=$(date +%s)

# Spawn SWE-1
echo "$(date +%s): Spawning SWE-1 for effort-1 fixes"
# [Actual spawn command with FIX-INSTRUCTIONS-SWE-effort-1.md]

# Spawn SWE-2 (must be within 5 seconds of SWE-1)
echo "$(date +%s): Spawning SWE-2 for effort-2 fixes"
# [Actual spawn command with FIX-INSTRUCTIONS-SWE-effort-2.md]

END_TIME=$(date +%s)
SPAWN_DURATION=$((END_TIME - START_TIME))

if [ $SPAWN_DURATION -gt 5 ]; then
    echo "⚠️ WARNING: R151 VIOLATION - Spawn time ${SPAWN_DURATION}s > 5s"
fi

# Document spawns in state file
cat <<EOF | echo "✅ State file updated to: $NEXT_STATE"
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
git commit -m "state: FIX_WAVE_UPSTREAM_BUGS → $NEXT_STATE - FIX_WAVE_UPSTREAM_BUGS complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "FIX_WAVE_UPSTREAM_BUGS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - FIX_WAVE_UPSTREAM_BUGS complete [R287]"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

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
PROPOSED_NEXT_STATE="NEXT_STATE"
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
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

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
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

