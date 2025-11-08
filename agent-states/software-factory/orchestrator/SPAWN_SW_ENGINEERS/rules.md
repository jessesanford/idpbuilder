# SPAWN_SW_ENGINEERS State Rules

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
**State**: SPAWN_SW_ENGINEERS
**Agent**: orchestrator
**Type**: spawn
**Iteration Level**: wave
**Checkpoint**: No (but R313 requires STOP after spawning)

---

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
save_todos "SPAWNED_SW_ENGINEER"

# 4. R322: Stop conversation (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# 5. R405: CONTINUATION FLAG - MUST BE TRUE CHECKPOINT=R322 FOR SPAWNING!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"

# 6. Exit to end conversation
exit 0
```

**Enhanced Format**: The `CHECKPOINT=R322` context tells the test framework this is a normal R322 checkpoint, enabling automatic continuation.

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
6. [ ] **CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322 emitted** ← Critical!
7. [ ] Exited with `exit 0`

**Missing step 6 or wrong flag = Test failures = -100% grade**
**Using FALSE at R322 checkpoint = Defeats automation = -100% grade**

---


## State Description

Spawn SW Engineer agents to implement features according to effort implementation plans created by Code Reviewers.

This is the PRIMARY IMPLEMENTATION STATE in Software Factory 3.0 - where actual feature development begins.

---

## Entry Conditions

Before entering this state, verify:

- [ ] Implementation plans exist for all efforts in current wave
- [ ] Code Reviewers have completed effort planning (MONITORING_EFFORT_PLANNING complete)
- [ ] Effort infrastructure validated and ready (branches created, workspaces set up)
- [ ] orchestrator-state-v3.json shows current wave and efforts

---

## Responsibilities

When in this state, orchestrator MUST:

### 1. Load Implementation Plans

```bash
# For each effort in wave
for effort in wave_efforts:
    effort_plan = load("${effort}/implementation-plan.md")
    verify effort_plan exists and is valid
```

### 2. Determine Parallelization Strategy

**Per R356 Single-Effort Optimization**:
- If wave has 1 effort → Spawn 1 SW Engineer for entire wave
- If wave has >1 effort → Spawn 1 SW Engineer per effort (parallel)

**Per R151 Parallelization Requirements**:
- All parallel spawns MUST occur within 5 seconds
- Emit timestamp before first spawn
- Emit timestamp after last spawn
- Verify delta <5s

### 3. Pull Latest from Base Branch (R614 - SUPREME LAW)

**CRITICAL**: Before spawning ANY agent, MUST pull latest from base branch to ensure fresh code.

```bash
# For EACH effort directory, pull latest from base branch
for effort_dir in wave_effort_directories:
    echo "🔄 R614: Pulling latest for $effort_dir"

    # CD to effort directory (R208)
    cd "$effort_dir" || {
        echo "🚨 FATAL: Cannot CD to $effort_dir"
        exit 208
    }

    # Get base branch from pre_planned_infrastructure
    BASE_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$(basename $effort_dir)\".base_branch" \
                  orchestrator-state-v3.json)

    if [ -z "$BASE_BRANCH" ] || [ "$BASE_BRANCH" = "null" ]; then
        echo "🚨 FATAL: No base_branch found for $(basename $effort_dir)"
        exit 614
    fi

    echo "   Base branch: $BASE_BRANCH"

    # Fetch latest from origin
    git fetch origin "$BASE_BRANCH" || {
        echo "⚠️ WARNING: Cannot fetch $BASE_BRANCH"
        # Retry logic here if needed
        exit 614
    }

    # Pull latest from base branch (CRITICAL FOR CASCADE INTEGRITY!)
    git pull origin "$BASE_BRANCH" || {
        echo "🚨 FATAL: Cannot pull from $BASE_BRANCH"
        echo "This breaks cascade integrity!"
        exit 614
    }

    # Verify fresh base
    REMOTE_HEAD=$(git rev-parse origin/$BASE_BRANCH)
    MERGE_BASE=$(git merge-base HEAD origin/$BASE_BRANCH)

    if [[ "$MERGE_BASE" == "$REMOTE_HEAD" ]]; then
        echo "   ✅ On latest commit from $BASE_BRANCH"
        echo "      Commit: $(git rev-parse --short $REMOTE_HEAD)"
    else
        echo "   ⚠️ WARNING: May not have all latest commits"
        echo "      Remote: $(git rev-parse --short $REMOTE_HEAD)"
        echo "      Base:   $(git rev-parse --short $MERGE_BASE)"
    fi

    # Return to orchestrator directory
    cd - > /dev/null
done

echo "✅ R614: All effort directories on fresh base branches"
```

**R614 ensures**:
- Sequential efforts have latest bug fixes from previous efforts
- No duplicate bug discovery
- No downstream rebases needed
- CASCADE pattern maintained

**Skip pull ONLY IF**:
- First effort in wave (base_branch = integration branch that doesn't exist yet)
- In that case, verify integration branch doesn't exist before skipping

### 4. Spawn SW Engineers

```bash
# Parallel spawn pattern (if multiple efforts)
SPAWN_START=$(date +%s)

for effort in wave_efforts (parallel):
    claude-code --agent sw-engineer \
        --workspace /efforts/phase${P}/wave${W}/${effort} \
        --branch ${effort}-implementation \
        --plan ${effort}/implementation-plan.md \
        --output json > ${effort}-spawn-result.json &

SPAWN_END=$(date +%s)
DELTA=$((SPAWN_END - SPAWN_START))

if [ $DELTA -gt 5 ]; then
    echo "🚨 VIOLATION: R151 parallelization timing exceeded (${DELTA}s > 5s)"
    exit 151
fi
```

**Single effort pattern (R356 optimization)**:
```bash
claude-code --agent sw-engineer \
    --workspace /efforts/phase${P}/wave${W} \
    --branch wave${W}-implementation \
    --plan wave${W}/implementation-plan.md \
    --output json > wave-spawn-result.json
```

### 5. Record Spawns

Update orchestrator-state-v3.json:
```json
{
  "active_agents": {
    "sw_engineers": [
      {
        "effort": "effort-123",
        "spawned_at": "2025-10-09T14:30:00Z",
        "workspace": "/efforts/phase1/wave1/effort-123",
        "plan_file": "effort-123/implementation-plan.md"
      }
    ]
  }
}
```

### 6. STOP per R313

**CRITICAL**: After spawning agents, orchestrator MUST:
- Update orchestrator-state-v3.json with current_state: "SPAWN_SW_ENGINEERS"
- Commit and push state file (R288 atomic update via State Manager)
- Exit with `CONTINUE-SOFTWARE-FACTORY=FALSE`
- **DO NOT monitor progress in same execution**

**R313 Rationale**: Prevents context/rule loss in spawned agents.

---

## State Transitions

### Valid Next States

After SPAWN_SW_ENGINEERS, the ONLY valid next state is:

1. **MONITORING_SWE_PROGRESS** (when continuation invoked)
   - Guard: SW Engineers spawned successfully
   - Condition: Spawn records exist in orchestrator-state-v3.json

2. **ERROR_RECOVERY** (if spawn fails)
   - Guard: Spawn errors detected
   - Condition: spawn_errors > 0

### Invalid Transitions

❌ **FORBIDDEN**:
- SPAWN_SW_ENGINEERS → MONITORING_SWE_PROGRESS (in same execution) - Violates R313
- SPAWN_SW_ENGINEERS → SPAWN_CODE_REVIEWERS_* - Wrong agent type
- SPAWN_SW_ENGINEERS → *_INTEGRATE_WAVE_EFFORTS - Skips implementation monitoring

---

## Exit Conditions

Can exit this state when:

- [ ] All SW Engineers spawned successfully
- [ ] Spawn records committed to orchestrator-state-v3.json
- [ ] State file pushed to remote
- [ ] R313 STOP executed

---

## State File Updates

### Updates to orchestrator-state-v3.json

**Shutdown Consultation with State Manager**:

```yaml
work_report:
  current_state: SPAWN_SW_ENGINEERS
  work_completed: "Spawned SW Engineers for wave implementation"
  results:
    sw_engineers_spawned: 3
    efforts: ["effort-auth", "effort-storage", "effort-api"]
    spawn_timing_delta: "2.1s"
    r151_compliant: true
  proposed_next_state: MONITORING_SWE_PROGRESS
  reasoning: "SW Engineers spawned, must monitor implementation progress per R313"
```

**State Manager performs atomic update**:
- orchestrator-state-v3.json: current_state → MONITORING_SWE_PROGRESS
- orchestrator-state-v3.json: active_agents.sw_engineers → populated
- bug-tracking.json: (no changes)
- integration-containers.json: (no changes)

### Updates to bug-tracking.json

No updates during SPAWN_SW_ENGINEERS (bugs discovered during review, not spawn).

### Updates to integration-containers.json

No updates during SPAWN_SW_ENGINEERS (iteration containers start after implementation complete).

---

## Common Issues and Resolutions

### Issue 1: Spawn Timing Violation (R151)

**Symptom**: Spawn timing delta >5s between first and last engineer spawn

**Cause**: Sequential spawning instead of parallel, or system slowness

**Resolution**:
1. Verify spawn commands use background processes (`&`)
2. Check system load during spawn
3. If systematic: Update R151 tolerance (requires rule change)
4. Document violation in state file

### Issue 2: Missing Implementation Plans

**Symptom**: Cannot find ${effort}/implementation-plan.md

**Cause**: Code Reviewer didn't create plan, or wrong directory

**Resolution**:
1. Verify MONITORING_EFFORT_PLANNING completed
2. Check effort directory structure
3. Search for plan files: `find . -name "*implementation-plan.md"`
4. If missing: Transition to ERROR_RECOVERY, re-run effort planning

### Issue 3: Workspace Conflicts

**Symptom**: SW Engineer spawn fails with "directory already exists"

**Cause**: Previous incomplete run, leftover workspaces

**Resolution**:
1. Check workspace directories: `ls /efforts/phase${P}/wave${W}/`
2. Verify no active git operations: `git status` in each workspace
3. If clean: Remove and recreate workspace
4. If dirty: Transition to ERROR_RECOVERY

### Issue 4: R313 Violation (Continued After Spawn)

**Symptom**: Orchestrator attempts to monitor SWE progress in same execution

**Cause**: Forgot R313 STOP requirement

**Resolution**:
- **NEVER allow this** - R313 is supreme law
- Code review orchestrator logic
- Ensure explicit exit after spawn with CONTINUE-SOFTWARE-FACTORY=FALSE

---

## Integration with SF 3.0 Architecture

### State Manager Bookend Pattern

**Startup Consultation**:
```
Orchestrator: "What should I do in SPAWN_SW_ENGINEERS?"

State Manager:
  - Load effort plans
  - Spawn SW Engineers (1 per effort OR 1 for wave per R356)
  - Ensure R151 timing compliance (<5s delta)
  - Update state file with spawn records
  - STOP per R313
  - Valid next state: MONITORING_SWE_PROGRESS
```

**Shutdown Consultation**:
```
Orchestrator: "I spawned 3 SW Engineers, transition to MONITORING_SWE_PROGRESS?"

State Manager:
  - Validates spawn records exist
  - Validates R151 timing compliance
  - Validates MONITORING_SWE_PROGRESS is in allowed_transitions
  - Performs atomic 4-file update (R288)
  - Returns: CONTINUE-SOFTWARE-FACTORY=FALSE (per R313 STOP)
```

### Iteration Container Context

**SPAWN_SW_ENGINEERS operates BEFORE iteration containers**:

```
WAVE_START
  ↓
CREATE_NEXT_INFRASTRUCTURE
  ↓
VALIDATE_EFFORT_INFRASTRUCTURE
  ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
  ↓
MONITORING_EFFORT_PLANNING
  ↓
SPAWN_SW_ENGINEERS  ← YOU ARE HERE
  ↓
MONITORING_SWE_PROGRESS
  ↓
SPAWN_CODE_REVIEWERS_EFFORT_REVIEW
  ↓
MONITORING_EFFORT_REVIEWS
  ↓
WAVE_COMPLETE
  ↓
SETUP_WAVE_INFRASTRUCTURE  ← Iteration container STARTS here
```

**Iteration containers handle INTEGRATE_WAVE_EFFORTS bugs, not IMPLEMENTATION bugs**:
- Implementation bugs → Fixed in effort branches via CREATE_EFFORT_FIX_PLAN loop
- Integration bugs → Fixed via iteration container FIX_*_UPSTREAM_BUGS states

---

## Testing and Validation

### Pre-Spawn Validation Checklist

```bash
✅ Phase plan exists: /plans/phase${P}-plan.md
✅ Wave plan exists: /plans/phase${P}/wave${W}-plan.md
✅ Effort plans exist: /efforts/phase${P}/wave${W}/${effort}/implementation-plan.md
✅ Effort branches exist: git branch -r | grep ${effort}-implementation
✅ Workspaces clean: git status (all clean)
✅ orchestrator-state-v3.json current_state: SPAWN_SW_ENGINEERS
```

### Post-Spawn Validation Checklist

```bash
✅ Spawn records in orchestrator-state-v3.json
✅ R151 timing: spawn_timing_delta < 5s
✅ State file committed: git log -1 --grep "SPAWN_SW_ENGINEERS"
✅ State file pushed: git ls-remote origin HEAD
✅ CONTINUE-SOFTWARE-FACTORY=FALSE emitted
```

---

## Example Execution

```bash
# Orchestrator enters SPAWN_SW_ENGINEERS state
echo "📍 STATE: SPAWN_SW_ENGINEERS"
echo "   Phase: 1, Wave: 1, Efforts: 3"

# 1. Consult State Manager
state-manager --mode startup --state SPAWN_SW_ENGINEERS

# 2. Load implementation plans
for effort in effort-auth effort-storage effort-api; do
    plan=/efforts/phase1/wave1/${effort}/implementation-plan.md
    echo "✅ Loaded: $plan"
done

# 3. Spawn SW Engineers (parallel, R151 compliant)
SPAWN_START=$(date +%s)

claude-code --agent sw-engineer --workspace /efforts/phase1/wave1/effort-auth ... &
claude-code --agent sw-engineer --workspace /efforts/phase1/wave1/effort-storage ... &
claude-code --agent sw-engineer --workspace /efforts/phase1/wave1/effort-api ... &

wait

SPAWN_END=$(date +%s)
DELTA=$((SPAWN_END - SPAWN_START))
echo "✅ R151 Timing: ${DELTA}s (compliant: $( [ $DELTA -lt 5 ] && echo YES || echo NO ))"

# 4. Update state file (via State Manager)
state-manager --mode shutdown \
    --current-state SPAWN_SW_ENGINEERS \
    --next-state MONITORING_SWE_PROGRESS \
    --spawned-agents 3 \
    --timing-delta ${DELTA}

# 5. STOP per R313
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
exit 0
```

---

## Grading Criteria

This state is evaluated on:

1. **R151 Parallelization** (15%)
   - ✓ Parallel spawns within 5s timing delta
   - ✓ Timestamp emissions before/after spawn

2. **R313 STOP Enforcement** (25%)
   - ✓ Orchestrator STOPS after spawning
   - ✓ No monitoring in same execution
   - ✓ CONTINUE-SOFTWARE-FACTORY=FALSE emitted

3. **R356 Optimization** (10%)
   - ✓ Single-effort waves use single engineer
   - ✓ Multi-effort waves spawn in parallel

4. **State File Updates** (20%)
   - ✓ orchestrator-state-v3.json updated atomically
   - ✓ Spawn records accurate
   - ✓ Committed and pushed per R288

5. **Agent Spawning** (30%)
   - ✓ All SW Engineers spawned successfully
   - ✓ Correct workspaces assigned
   - ✓ Implementation plans provided

---

**CRITICAL**: This state enables Item #703 testing - orchestrator spawning SW Engineers via CLI with JSON output verification.

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_SW_ENGINEERS:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_SW_ENGINEERS complete - [describe what was accomplished]"
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
  "state_completed": "SPAWN_SW_ENGINEERS",
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
  --current-state "SPAWN_SW_ENGINEERS" \
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

if ! git commit -m "state: SPAWN_SW_ENGINEERS → $NEXT_STATE - SPAWN_SW_ENGINEERS complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: SPAWN_SW_ENGINEERS"
    echo "Attempted transition from: SPAWN_SW_ENGINEERS"
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
save_todos "SPAWN_SW_ENGINEERS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_SW_ENGINEERS complete [R287]"; then
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
# Enhanced format with CHECKPOINT=R322 for spawn states
# This tells automation this is a normal R322 checkpoint, enabling auto-continue

echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

**Enhanced Format**: The `CHECKPOINT=R322` context is now **MANDATORY** for R322 checkpoints.
- Tells framework this is normal operation (not error)
- Enables automatic continuation in tests
- Makes intent explicit in logs

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

