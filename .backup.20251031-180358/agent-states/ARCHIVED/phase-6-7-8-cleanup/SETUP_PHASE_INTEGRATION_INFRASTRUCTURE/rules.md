# Orchestrator - SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R308** - Incremental Branching Strategy (EFFORT CASCADE ONLY)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`
   - Criticality: SUPREME LAW - Defines EFFORT cascade flow
   - Summary: R308 is for EFFORT branching, NOT integration base determination

6. **🔴🔴🔴 R009** - Mandatory Wave/Phase Integration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R009-integration-branch-creation.md`
   - Criticality: SUPREME LAW - Wave N+1 REQUIRES Wave N integration
   - Summary: Enforces trunk-based development

## 🚨🚨🚨 CRITICAL: INFRASTRUCTURE PATH REQUIREMENTS 🚨🚨🚨

**VIOLATION = -100% AUTOMATIC FAILURE**

### UNIVERSAL RULE: ALL INFRASTRUCTURE MUST BE UNDER efforts/

**CORRECT PATHS:**
- Phase integration: `$CLAUDE_PROJECT_DIR/efforts/phase{X}/integration/`
- Wave integration: `$CLAUDE_PROJECT_DIR/efforts/phase{X}/wave{Y}/integration/`
- Project integration: `$CLAUDE_PROJECT_DIR/efforts/project/integration/`
- Effort workspaces: `$CLAUDE_PROJECT_DIR/efforts/phase{X}/wave{Y}/{effort-name}/`

**NEVER USE:**
- `$CLAUDE_PROJECT_DIR/phase{X}/` ❌
- `$CLAUDE_PROJECT_DIR/wave{Y}/` ❌
- `$CLAUDE_PROJECT_DIR/integration/` ❌
- `$CLAUDE_PROJECT_DIR/project/` ❌

### MANDATORY PATH VALIDATION FUNCTION:
```bash
# Validate infrastructure path - MUST be called before creating infrastructure
validate_infrastructure_path() {
    local path="$1"

    # CRITICAL: ALL infrastructure MUST be under efforts/
    if [[ ! "$path" =~ ^.*efforts/.* ]]; then
        echo "🔴🔴🔴 CRITICAL PATH VIOLATION 🔴🔴🔴"
        echo "Infrastructure path MUST be under efforts/ directory!"
        echo "Got: $path"
        echo "Required pattern: \$CLAUDE_PROJECT_DIR/efforts/**/*"
        exit 1
    fi

    echo "✅ Path validation passed: $path"
}

# ALWAYS validate before creating infrastructure
INTEGRATE_WAVE_EFFORTS_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/integration"
validate_infrastructure_path "$INTEGRATE_WAVE_EFFORTS_DIR"
```

**Penalty for violation**: -100% (IMMEDIATE FAILURE)

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

## 🔴🔴🔴 CRITICAL: PHASE INTEGRATE_WAVE_EFFORTS BASE BRANCH DETERMINATION (R282/R512) 🔴🔴🔴

**VIOLATION = -100% AUTOMATIC FAILURE**

### PHASE INTEGRATE_WAVE_EFFORTS MUST FOLLOW R282/R512 SEQUENTIAL REBUILD MODEL:

**Phase Integration branches are created from the FIRST EFFORT of the phase:**
- Phase 1 Integration: from phase1/wave1/effort1 (first effort of phase 1)
- Phase 2 Integration: from phase2/wave1/effort1 (first effort of phase 2)
- Phase 3 Integration: from phase3/wave1/effort1 (first effort of phase 3)
- **NOT from wave integrations** (those are testing checkpoints per R364)

### 🔴 CRITICAL R282/R512 ENFORCEMENT FOR PHASE INTEGRATE_WAVE_EFFORTS:
```bash
# Example: Phase 2 Integration (after completing Phase 2 Wave 3)
# WRONG - AUTOMATIC FAILURE:
BASE_BRANCH="main"  # ❌ NEVER for phase integration!
BASE_BRANCH="phase1-integration"  # ❌ WRONG! Not from previous phase!
BASE_BRANCH="phase2-wave3-integration"  # ❌ WRONG! Not from wave integration (R512/R270 violation)!

# CORRECT (R282 Sequential Rebuild Model):
BASE_BRANCH="phase2/wave1/effort1"  # ✅ From FIRST EFFORT of THIS phase!
```

### PHASE INTEGRATE_WAVE_EFFORTS BASE DETERMINATION LOGIC (R282 SEQUENTIAL REBUILD):
```bash
determine_phase_integration_base() {
    local PHASE=$1

    echo "🔴 R282/R512: Phase $PHASE integration base determination"
    echo "📊 Using SEQUENTIAL REBUILD model (base = first effort of phase)"

    # Phase integration uses FIRST EFFORT of THIS phase (R282 Sequential Rebuild)
    FIRST_EFFORT=$(echo "✅ State file updated to: $NEXT_STATE"
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
git commit -m "state: SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE → $NEXT_STATE - SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE complete [R287]"
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

