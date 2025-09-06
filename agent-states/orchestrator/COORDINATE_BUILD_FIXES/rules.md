# Orchestrator - COORDINATE_BUILD_FIXES State Rules

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

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED COORDINATE_BUILD_FIXES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_COORDINATE_BUILD_FIXES
echo "$(date +%s) - Rules read and acknowledged for COORDINATE_BUILD_FIXES" > .state_rules_read_orchestrator_COORDINATE_BUILD_FIXES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR COORDINATE_BUILD_FIXES STATE

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
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-frequency.md`
**Criticality**: BLOCKING - Save every 15 minutes/10 messages
**Summary**: Mandatory TODO saves during coordination

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml on all state changes

### State-Specific Rules

### 🚨🚨🚨 R151 - Parallel Agent Timestamp Requirement (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
**Criticality**: BLOCKING - Parallel agents must start within 5 seconds
**Summary**: When spawning multiple SW Engineers, ensure synchronized start

### 🔴🔴🔴 R321 - Immediate Backport During Integration (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: SUPREME LAW - Integration fixes must be backported immediately
**Summary**: Track all fixes for immediate backporting to source branches

## 🔴🔴🔴 CRITICAL: COORDINATE_BUILD_FIXES IS A VERB - START COORDINATING NOW! 🔴🔴🔴

**COORDINATE_BUILD_FIXES MEANS ACTIVELY COORDINATING FIX EFFORTS RIGHT NOW!**
- ❌ NOT "I'm in coordinate build fixes state"  
- ❌ NOT "Ready to coordinate fixes"
- ✅ YES "I'm distributing fix plans to engineers NOW"
- ✅ YES "I'm spawning SW Engineers for fixes NOW"
- ✅ YES "I'm tracking fix assignments NOW"

## State Context
COORDINATE_BUILD_FIXES = You ARE ACTIVELY coordinating the distribution and execution of build fixes THIS INSTANT. You have fix plans from Code Reviewer and must spawn SW Engineers to implement them.

## 🎯 STATE OBJECTIVES

In the COORDINATE_BUILD_FIXES state, you are responsible for:

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
State: COORDINATE_BUILD_FIXES

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
Assigned by: orchestrator/COORDINATE_BUILD_FIXES

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
State: COORDINATE_BUILD_FIXES
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
yq -i '.build_fix_coordination.spawned_engineers += [
    {"id": "SWE-1", "effort": "effort-1", "spawned_at": "'$(date +%s)'"},
    {"id": "SWE-2", "effort": "effort-2", "spawned_at": "'$(date +%s)'"}
]' orchestrator-state.yaml

# SEQUENTIAL GROUP - Will spawn after dependencies met
echo "📝 SWE-3 will be spawned after SWE-1 completes"
```

### Step 6: Set Up Progress Tracking
```bash
# Create monitoring dashboard
cat > BUILD-FIX-PROGRESS.md << 'EOF'
# Build Fix Progress Tracker
Date: $(date)
State: COORDINATE_BUILD_FIXES

## Active SW Engineers
| Engineer | Effort | Status | Started | Progress |
|----------|--------|--------|---------|----------|
| SWE-1 | effort-1 | IN_PROGRESS | [time] | 0% |
| SWE-2 | effort-2 | IN_PROGRESS | [time] | 0% |
| SWE-3 | effort-3 | PENDING | - | - |

## Fix Completion Checklist
### Compilation Fixes
- [ ] effort-1: Type errors resolved
- [ ] effort-1: Import errors fixed
- [ ] effort-2: Syntax errors corrected

### Dependency Fixes
- [ ] Missing packages added
- [ ] Version conflicts resolved
- [ ] Lock files updated

### Linking Fixes
- [ ] Undefined references resolved
- [ ] Symbol exports corrected
- [ ] Build order fixed

## Verification Status
- [ ] Individual builds tested
- [ ] Integration build attempted
- [ ] All tests passing
- [ ] Ready for backporting

## Next State
When all fixes complete: MONITOR_FIXES or BUILD_VALIDATION
EOF
```

## ⚠️ CRITICAL REQUIREMENTS

### 🔴 R006 VIOLATION DETECTION
If you find yourself about to:
- Edit any code file directly
- Write fix patches
- Make "quick fixes" yourself
- Use Edit/Write on source files

**STOP IMMEDIATELY - You are violating R006!**
Always spawn SW Engineers for ALL code changes!

### 🔴 R151 PARALLEL SPAWN TIMING
When spawning multiple SW Engineers:
- All parallel spawns MUST be within 5 seconds
- Each engineer must emit timestamp on startup
- Document spawn times in state file
- **Violation = -25% per late spawn**

### 🔴 R321 BACKPORT REQUIREMENT
**During integration context:**
- Track EVERY fix for backporting
- Document source branches
- Create backport manifest
- Immediate backport after fix success
- **Violation = -100% FAILURE**

## State Transitions

From COORDINATE_BUILD_FIXES state:
- **ENGINEERS_SPAWNED** → MONITOR_FIXES (Monitor fix progress)
- **NO_FIXES_NEEDED** → BUILD_VALIDATION (Rerun validation)
- **SPAWN_FAILED** → ERROR_RECOVERY (Cannot spawn engineers)

## Success Criteria
- ✅ All fix plans loaded and understood
- ✅ Fix work properly distributed
- ✅ SW Engineers spawned with clear instructions
- ✅ Parallel spawns within 5-second window (R151)
- ✅ Backport manifest created (R321)
- ✅ Progress tracking established

## Failure Triggers
- ❌ Auto-transition without stop = R322 VIOLATION
- ❌ Attempting to fix code yourself = R006 VIOLATION (-100%)
- ❌ Parallel spawn >5 seconds = R151 VIOLATION (-25%)
- ❌ Missing backport tracking = R321 VIOLATION (-100%)
- ❌ No fix instructions = -40% penalty

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**