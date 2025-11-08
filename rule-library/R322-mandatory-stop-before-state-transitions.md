# 🛑🛑🛑 RULE R322: MANDATORY ORCHESTRATOR CHECKPOINTS (SUPREME LAW)

**CRITICALITY**: 🔴🔴🔴 SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE

# 🔴🔴🔴 CRITICAL: R322 + R405 INTERACTION - READ THIS FIRST 🔴🔴🔴

## THE #1 MISTAKE CAUSING TEST FAILURES

**R322 requires "STOP" → Many developers incorrectly use CONTINUE-SOFTWARE-FACTORY=FALSE**

**THIS IS WRONG AND DEFEATS SOFTWARE FACTORY AUTOMATION!**

### THE CORRECT PATTERN FOR SPAWN STATES (MANDATORY):

```bash
# After spawning agents in any SPAWN_* state:

# 1. Complete state work
echo "✅ Spawned 3 SW Engineers for implementation"

# 2. Update state file to NEXT_STATE (R324)
update_state "MONITORING_SWE_PROGRESS"
commit_state()

# 3. R322: Stop conversation (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# 4. R405: CONTINUATION FLAG WITH CHECKPOINT CONTEXT (MANDATORY!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"

# 5. Exit to end conversation turn
exit 0
```

**CRITICAL**: The `CHECKPOINT=R322` context is now **MANDATORY** for all R322 stops!
- Tells automation this is a NORMAL checkpoint, not an error
- Enables test framework to auto-continue
- Makes intent explicit in logs

**R322 "STOP" means:**
- End this conversation turn (`exit 0`)
- Preserve context for next turn
- **NOT "stop the entire automation"!**

**R405 TRUE means:**
- Automation CAN restart the conversation
- System knows next state (from state file)
- Normal operation is proceeding

**These are INDEPENDENT decisions that both happen in normal operations!**

### 🚨 TEST 2 FAILURE ROOT CAUSE

**What happened:**
1. Orchestrator spawned agents ✅ (correct)
2. Orchestrator stopped per R322 ✅ (correct)
3. Orchestrator **did NOT emit `CONTINUE-SOFTWARE-FACTORY=TRUE`** ❌ (WRONG!)
4. Test framework saw no flag → stopped automation
5. Test 2 failed at iteration 8

**Fix:** ALWAYS emit `CONTINUE-SOFTWARE-FACTORY=TRUE` for normal spawning operations!

---

# 🚨 CRITICAL WARNING: R322 Stop ≠ CONTINUE-SOFTWARE-FACTORY=FALSE 🚨

**MOST COMMON MISCONCEPTION CAUSING -20% TO -100% PENALTIES:**

Developers REPEATEDLY misunderstand R322 "mandatory stop" as requiring
`CONTINUE-SOFTWARE-FACTORY=FALSE`. This is **COMPLETELY INCORRECT** and
**DEFEATS THE PURPOSE OF SOFTWARE FACTORY AUTOMATION!**

**THIS EXACT SCENARIO IS CAUSING VIOLATIONS:**

> "I'm in WAITING_FOR_PHASE_MERGE_PLAN state"
> "Merge plan exists and is validated ✅"
> "R322 says I must stop before SPAWN_INTEGRATION_AGENT_PHASE transition"
> "Therefore I set CONTINUE-SOFTWARE-FACTORY=FALSE"

**EVERYTHING ABOVE IS WRONG EXCEPT THE MERGE PLAN VALIDATION!**

**READ THIS ENTIRE SECTION BEFORE PROCEEDING WITH ANY R322 STOP!**

## 🔴🔴🔴 CRITICAL DISTINCTION: STOP INFERENCE vs CONTINUATION FLAG 🔴🔴🔴

**TWO COMPLETELY SEPARATE CONCEPTS - DO NOT CONFUSE THEM:**

### 1. STOPPING INFERENCE (Context Preservation)
- **Purpose**: Preserve context between states to prevent overflow
- **Required**: After spawning multiple agents or at state boundaries
- **This is NORMAL**: Stopping is part of healthy state transitions
- **Action**: `exit 0` to end conversation turn
- **User Experience**: User sees "/continue-orchestrating" as next step

### 2. CONTINUATION FLAG (Automation Control)
- **Purpose**: Tell external automation whether to restart automatically
- **TRUE**: Normal operations - automation will restart Claude Code
- **FALSE**: Exceptional cases - human intervention REQUIRED (not just helpful)
- **Action**: `echo "CONTINUE-SOFTWARE-FACTORY=TRUE"` or `FALSE` (LAST output)

### 🎯 KEY INSIGHT: We stop inference but set flag=TRUE for normal operations!

**NORMAL R322 CHECKPOINT PATTERN:**
```bash
# State work completed successfully
echo "✅ Merge plan validated and ready"
echo "✅ Transitioning to SPAWN_INTEGRATION_AGENT_PHASE"

# Update state file per R324
update_state "SPAWN_INTEGRATION_AGENT_PHASE"
commit_state_per_r288()

# Stop for R322 checkpoint (context preservation)
echo "🛑 R322: Checkpoint before state transition"

# Set flag based on whether automation can continue
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # ← YES! This is NORMAL!

# Stop inference
exit 0
```

**WHAT HAPPENS NEXT:**
1. User runs `/continue-orchestrating` (designed UX)
2. Orchestrator restarts, reads state file
3. Sees current_state = "SPAWN_INTEGRATION_AGENT_PHASE"
4. Executes that state (spawns integration agent)
5. **THIS IS NORMAL OPERATION! NOT AN ERROR!**

### ❌ COMMON MISCONCEPTION CAUSING VIOLATIONS:

**WRONG INTERPRETATION:**
> "R322 requires mandatory stop before state transitions"
> "State work is complete and I'm waiting for /continue-orchestrating"
> "Therefore I should set CONTINUE-SOFTWARE-FACTORY=FALSE"

**CORRECT INTERPRETATION:**
> "R322 requires ending the conversation (`exit 0`) to preserve context."
> "The flag indicates whether automation CAN restart (TRUE) or CANNOT (FALSE)."
> "Completing state work and waiting for /continue is NORMAL, DESIGNED behavior."
> "System knows what to do next (it's in the state file), so use TRUE!"

**REMEMBER:**
- **Stop inference** = End this conversation turn (context preservation) [ALWAYS at R322 points]
- **Continuation flag** = Can automation proceed when restarted? [TRUE for normal, FALSE for catastrophic]
- **These are independent!** Normal operations: Stop conversation (`exit 0`) + Flag TRUE

**R322 CHECKPOINTS ARE NORMAL OPERATION, NOT ERRORS!**

## Purpose
Enforce mandatory checkpoints ONLY at EXCEPTIONAL orchestrator transitions to:
1. Allow user review of CRITICAL plans before execution (integration/merge plans ONLY)
2. Preserve context after spawning agents (prevents overflow)
3. Provide clear decision points for EXCEPTIONAL situations
4. Ensure proper state persistence at CRITICAL junctures

**IMPORTANT CLARIFICATION**: Normal state transitions do NOT require stops! Only the specific transitions listed below require checkpoints.

**THIS RULE CONSOLIDATES:**
- Former R313 (stop after spawning) - now Part A of this rule
- Checkpoint stops for user review - Part B of this rule
- Works with R324 (state update before stop) as companion rule

## 🔴🔴🔴 SUPREME LAW: STOP MEANS STOP 🔴🔴🔴

### MANDATORY STOP POINTS

The orchestrator MUST STOP and exit after these specific transitions:

#### PART A: SPAWN STATE STOPS (Context Preservation)
**🚨 REMEMBER**: Stop inference BUT set CONTINUE-SOFTWARE-FACTORY=TRUE for normal operations!

**WHEN TO STOP INFERENCE after spawning:**
- ✅ Spawning MULTIPLE implementation agents (context flood)
- ✅ Spawning agents that generate LARGE outputs (>1000 lines)
- ✅ Spawning for COMPLEX operations (integration work)

**WHAT FLAG TO SET:**
- For normal spawning: CONTINUE-SOFTWARE-FACTORY=**TRUE** (auto-restart allowed)
- Only set FALSE if system error or corruption detected

**SPECIFIC STOPS REQUIRED (with TRUE flag):**
- `SPAWN_SW_ENGINEERS` → `MONITORING_SWE_PROGRESS` (stop inference, flag=TRUE)
- `SPAWN_CODE_REVIEWERS_EFFORT_PLANNING` → `WAITING_FOR_EFFORT_PLANS` (stop, flag=TRUE)
- `SPAWN_CODE_REVIEWERS_EFFORT_REVIEW` → `MONITORING_EFFORT_REVIEWS` (stop, flag=TRUE)
- `SPAWN_SW_ENGINEERS` → `MONITORING_EFFORT_FIXES` (stop, flag=TRUE)
- `SPAWN_INTEGRATION_AGENT` → `MONITORING_INTEGRATE_WAVE_EFFORTS` (stop, flag=TRUE)
- `SPAWN_INTEGRATION_AGENT_PHASE` → `MONITORING_INTEGRATE_PHASE_WAVES` (stop, flag=TRUE)
- `SPAWN_INTEGRATION_AGENT_PROJECT` → `MONITORING_PROJECT_INTEGRATE_WAVE_EFFORTS` (stop, flag=TRUE)

**ALSO STOP WITH TRUE FLAG (single agent but important transition):**
- `SPAWN_ARCHITECT_PHASE_ASSESSMENT` → Stop inference, flag=TRUE (normal)
- `SPAWN_ARCHITECT_REVIEW_WAVE_ARCHITECTURE` → Stop inference, flag=TRUE (normal)
- `SPAWN_CODE_REVIEWER_FIX_PLAN` → Stop inference, flag=TRUE (normal)

**Reason**: Context preservation between states
**Flag**: TRUE because these are NORMAL operations that should auto-continue

#### PART B: CRITICAL INTEGRATE_WAVE_EFFORTS PLAN REVIEW STOPS
**ONLY these SPECIFIC integration transitions REQUIRE user review:**

- `WAITING_FOR_MERGE_PLAN` → `SPAWN_INTEGRATION_AGENT`
  - User must review WAVE-MERGE-PLAN.md before execution
  - Stop allows verification of merge order and strategy

- `WAITING_FOR_PHASE_MERGE_PLAN` → `SPAWN_INTEGRATION_AGENT_PHASE`
  - User must review PHASE-MERGE-PLAN.md before execution
  - Stop allows verification of phase integration approach

- `WAITING_FOR_PROJECT_MERGE_PLAN` → `SPAWN_INTEGRATION_AGENT_PROJECT`
  - User must review PROJECT-MERGE-PLAN.md before execution
  - Stop allows verification of final project integration

**IMPORTANT**: Normal implementation → review transitions do NOT require stops!

**CLARIFICATION ON NORMAL OPERATIONS (NO STOPS REQUIRED):**
- ✅ `WAITING_FOR_FIX_PLANS` → `CREATE_WAVE_FIX_PLAN` - FIX PLANS ARE NORMAL (continue automatically)
- ✅ `WAITING_FOR_BACKPORT_PLAN` → `SPAWN_SW_ENGINEER_BACKPORT_FIXES` - BACKPORTS ARE NORMAL (continue automatically)
- ✅ Any fix cascade transitions - FIXES ARE NORMAL OPERATIONS
- ✅ Any split implementation transitions - SPLITS ARE NORMAL OPERATIONS

**These are NORMAL SOFTWARE DEVELOPMENT ACTIVITIES that should proceed automatically!**

#### PART C: ASSESSMENT → ACTION STOPS
**These transitions REQUIRE user decision before proceeding:**

- ❌ REMOVED - `WAITING_FOR_PHASE_ASSESSMENT` → `COMPLETE_PHASE` IS NORMAL OPERATION
  - ✅ COMPLETE_PHASE is normal progression - CONTINUE AUTOMATICALLY
  - ✅ Architect already approved - no user review needed

- `WAITING_FOR_PROJECT_VALIDATION` → `CREATE_INTEGRATE_WAVE_EFFORTS_TESTING`
  - User must approve project validation results
  - Stop allows final review before testing

#### PART D: EXCEPTIONAL MONITORING_SWE_PROGRESS STOPS
**ONLY these specific monitoring transitions require stops:**

- `MONITORING_INTEGRATE_WAVE_EFFORTS` → `REVIEW_WAVE_ARCHITECTURE`
  - Must save integration results
  - Stop ensures proper state capture

- `MONITORING_INTEGRATE_PHASE_WAVES` → `SPAWN_ARCHITECT_PHASE_ASSESSMENT`
  - Must save phase integration results
  - Stop ensures proper state capture

**CLARIFICATION**: Regular monitoring transitions like:
- `MONITORING_SWE_PROGRESS` → `SPAWN_CODE_REVIEWERS_EFFORT_REVIEW` - NO STOP REQUIRED
- `MONITORING_EFFORT_REVIEWS` → `CREATE_NEXT_INFRASTRUCTURE` - NO STOP REQUIRED
- `MONITORING_EFFORT_FIXES` → Next state - NO STOP REQUIRED

These are NORMAL OPERATIONS and should continue automatically!

## STOP PROTOCOL

### 🚨 CRITICAL COMPANION: R324 MUST BE FOLLOWED
**Before ANY stop, MUST update current_state per R324 to prevent infinite loops!**

### ✅ REQUIRED ACTIONS BEFORE STOPPING:

1. **Complete ALL work for current state**
   ```bash
   # Finish all TODOs for current state
   # Process any pending items
   # Clean up working directory
   ```

2. **Update orchestrator-state-v3.json with NEW state**
   ```bash
   # SF 3.0: State Manager handles this via shutdown consultation
   # State Manager validates and performs atomic 4-file update

   # SF 2.0 fallback: Manual update
   jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
   jq '.state_machine.previous_state = "CURRENT_STATE"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
   jq '.transition_time = "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
   ```

3. **Save TODO state per R287**
   ```bash
   # Save current TODO list
   # Include reason for stop
   save_todos "R322_CHECKPOINT_BEFORE_${NEXT_STATE}"
   ```

4. **Commit and push state changes**
   ```bash
   # SF 3.0: State Manager performs atomic commit
   # Invokes tools/atomic-state-update.sh for all 4 state files
   # Validates schema compliance before committing
   # Rolls back on any validation failure

   # SF 2.0 fallback: Manual commit
   git add orchestrator-state-v3.json todos/*.todo
   git commit -m "state: R322 checkpoint before ${NEXT_STATE}"
   git push
   ```

5. **Display checkpoint message**
   ```markdown
   ## 🛑 R322 STATE TRANSITION CHECKPOINT
   
   ### ✅ Current State Work Completed:
   - [List completed work]
   - [Include any created artifacts]
   
   ### 📊 Current Status:
   - Previous State: CURRENT_STATE ✅
   - Next State: NEXT_STATE (ready to enter)
   - Checkpoint Reason: [e.g., "Merge plan ready for review"]
   
   ### 📋 Artifacts for Review:
   - [List files created: plans, reports, etc.]
   
   ### ⏸️ STOPPED - User Review Required
   Please review the above artifacts before continuing.
   To proceed, use: /continue-orchestrating
   ```

### ❌ FORBIDDEN ACTIONS:

1. **DO NOT continue to next state automatically**
   ```bash
   # ❌ WRONG: Automatic transition
   echo "Moving to SPAWN_INTEGRATION_AGENT..."
   spawn_integration_agent()  # VIOLATION!
   ```

2. **DO NOT start work for new state**
   ```bash
   # ❌ WRONG: Starting new work
   echo "Now spawning integration agent..."
   /spawn integration-agent  # VIOLATION!
   ```

3. **DO NOT assume permission to continue**
   ```bash
   # ❌ WRONG: Assuming continuation
   echo "Plan looks good, proceeding with integration..."
   ```

## ENFORCEMENT MECHANISM

### Detection Pattern
```python
def detect_r322_violation(current_state, next_state, action):
    """Detect R322 violations in state transitions"""

    CHECKPOINT_REQUIRED = [
        ("WAITING_FOR_MERGE_PLAN", "SPAWN_INTEGRATION_AGENT"),
        ("WAITING_FOR_PHASE_MERGE_PLAN", "SPAWN_INTEGRATION_AGENT_PHASE"),
        ("WAITING_FOR_PROJECT_MERGE_PLAN", "SPAWN_INTEGRATION_AGENT_PROJECT"),
        # REMOVED: Fix plans and backports are NORMAL operations
        # ("WAITING_FOR_FIX_PLANS", "CREATE_WAVE_FIX_PLAN"),  # NO STOP - NORMAL
        # ("WAITING_FOR_BACKPORT_PLAN", "SPAWN_SW_ENGINEER_BACKPORT_FIXES"),  # NO STOP - NORMAL
        # ("WAITING_FOR_PHASE_ASSESSMENT", "COMPLETE_PHASE"),  # NO STOP - NORMAL PROGRESSION
        ("WAITING_FOR_PROJECT_VALIDATION", "CREATE_INTEGRATE_WAVE_EFFORTS_TESTING"),
    ]
    
    if (current_state, next_state) in CHECKPOINT_REQUIRED:
        if action != "STOP":
            return {
                'violation': True,
                'severity': 'SUPREME_LAW',
                'penalty': '-100%',
                'reason': f'R322: Must stop between {current_state} and {next_state}'
            }
    
    return {'violation': False}
```

### Validation Script
```bash
#!/bin/bash
# Validate R322 compliance in state transition

validate_r322_checkpoint() {
    local CURRENT_STATE="$1"
    local NEXT_STATE="$2"
    
    # Check if this transition requires R322 stop
    case "${CURRENT_STATE}->${NEXT_STATE}" in
        "WAITING_FOR_MERGE_PLAN->SPAWN_INTEGRATION_AGENT"|\
        "WAITING_FOR_PHASE_MERGE_PLAN->SPAWN_INTEGRATION_AGENT_PHASE"|\
        "WAITING_FOR_PROJECT_MERGE_PLAN->SPAWN_INTEGRATION_AGENT_PROJECT")
            echo "🛑 R322 CHECKPOINT REQUIRED!"
            echo "User must review plan before execution"
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}
```

## GRADING IMPACT

### Violations Result In:
- **-100% IMMEDIATE FAILURE** for skipping required checkpoint
- **-50% PENALTY** for improper stop protocol
- **-25% PENALTY** for missing state persistence
- **AUTOMATIC REJECTION** of any work done after violation

### Compliance Results In:
- **+10% BONUS** for perfect checkpoint execution
- **Clean audit trail** for review and debugging
- **User confidence** in system operation
- **Context preservation** preventing rule loss

## RELATIONSHIP TO OTHER RULES

### Works With:
- **R324**: State file update before stop (CRITICAL COMPANION)
- **R287**: TODO persistence at checkpoints
- **R206**: State validation requirements
- **R231**: Continuous operation (EXCEPT at R322 checkpoints)
- **State Manager** (SF 3.0): Shutdown consultation validates and atomically commits all state changes

### Supersedes:
- **R313**: Now incorporated as Part A of this rule
- Any guidance suggesting automatic continuation at checkpoints
- Any pattern of immediate transition without checkpoint
- Previous "continue if successful" patterns at spawn points

## EXAMPLES

### ✅ CORRECT: R322 Checkpoint with TRUE CHECKPOINT=R322 (Spawn State)
```bash
# After spawning multiple SWE agents
echo "✅ Spawned 5 SWE agents for implementation"
echo "📊 Agents working on efforts E1.1, E1.2, E1.3, E1.4, E1.5"

# Update state and save TODOs
update_state "MONITORING_SWE_PROGRESS"
save_todos "SPAWNED_MULTIPLE_AGENTS"

# R322 checkpoint message
echo "🛑 R322: Stopping to preserve context (normal operation)"

# R405 flag with R322 checkpoint context (MANDATORY!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"

# Stop this conversation turn
exit 0
```

**What happens next:**
1. Test framework sees `TRUE CHECKPOINT=R322`
2. Framework recognizes this as normal R322 checkpoint
3. Framework **automatically** runs `/continue-software-factory`
4. Orchestrator restarts, reads state=MONITORING_SWE_PROGRESS
5. Orchestrator executes monitoring state
6. **No manual intervention needed!**

### ✅ CORRECT: Catastrophic Error requiring Human Intervention (FALSE with REASON)
```bash
# EXCEPTIONAL CASE: State file corrupted, cannot proceed
echo "❌ ERROR: orchestrator-state-v3.json is corrupted"
echo "❌ JSON parsing failed - invalid schema"
echo "❌ Cannot determine current state - CANNOT PROCEED"

# Display error details
cat << 'EOF'
## 🔴 CATASTROPHIC ERROR - MANUAL INTERVENTION REQUIRED

### ❌ State File Corruption Detected:
- File: orchestrator-state-v3.json
- Error: Invalid JSON schema
- Fields missing: current_state, state_machine
- Cannot recover automatically

### 🛠️ HUMAN ACTION REQUIRED:
1. Investigate state file corruption cause
2. Restore from backup or reconstruct manually
3. Validate JSON schema before continuing

### ⏸️ AUTOMATION STOPPED
EOF

# FALSE flag with REASON (this is TRUE exception case)
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=STATE_CORRUPTION"
exit 1  # Exit with error code
```

**What happens next:**
1. Test framework sees `FALSE REASON=STATE_CORRUPTION`
2. Framework knows specific error type
3. Framework **STOPS** and alerts human
4. Human must fix state file before continuing
5. **This is EXCEPTIONAL - not a normal R322 checkpoint!**

**NOTE**: R322 checkpoints at spawn states use `TRUE CHECKPOINT=R322`, NOT FALSE!

### ❌ WRONG: Using FALSE at R322 Checkpoint (COMMON MISTAKE!)
```bash
# After spawning agents - WRONG FLAG!
echo "✅ Spawned 3 SW Engineers"
update_state "MONITORING_SWE_PROGRESS"

echo "🛑 R322: Stopping after spawn"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ❌ WRONG! Defeats automation!
exit 0
```

**Why this is wrong:**
- R322 checkpoint is NORMAL workflow, not an error
- Using FALSE stops all automation unnecessarily
- Tests fail because framework doesn't know to continue
- **Correct**: Use `TRUE CHECKPOINT=R322` for R322 stops!

### ❌ WRONG: Missing Continuation Flag at R322 Checkpoint
```bash
# After spawning agents - NO FLAG!
echo "✅ Spawned 3 SW Engineers"
update_state "MONITORING_SWE_PROGRESS"

echo "🛑 R322: Stopping after spawn"
exit 0  # ❌ WRONG! No continuation flag!
```

**Why this is wrong:**
- R405 requires continuation flag on EVERY state completion
- Framework doesn't know whether to continue or stop
- Tests may hang or fail unexpectedly
- **Correct**: ALWAYS emit `CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322`

### ❌ WRONG: Skipping R322 Checkpoint Entirely
```bash
# After spawning agents - CONTINUING WITHOUT STOP!
echo "✅ Spawned 3 SW Engineers"
update_state "MONITORING_SWE_PROGRESS"

# ❌ VIOLATION: Continuing to next state without R322 stop!
echo "Now monitoring implementation progress..."
monitor_implementation()  # Context overflow risk!
```

**Why this is wrong:**
- R322 requires mandatory stop after spawning
- Prevents context overflow from agent outputs
- **Correct**: Stop with `exit 0` and flag `TRUE CHECKPOINT=R322`

## AUDIT REQUIREMENTS

Every R322 checkpoint must create:
1. State file update with transition recorded
2. TODO save with checkpoint reason
3. Git commit with R322 reference
4. User-visible checkpoint message
5. Clean exit (not error)

## SUMMARY

**R322 is a SUPREME LAW that ensures:**
- Users can review CRITICAL INTEGRATE_WAVE_EFFORTS PLANS before execution (not normal operations)
- Context is preserved after SPAWNING agents (not during monitoring)
- State is properly saved at EXCEPTIONAL transitions (not routine ones)
- System provides clear decision points for EXCEPTIONAL situations
- No automatic execution of CRITICAL INTEGRATE_WAVE_EFFORTS plans without review

**CRITICAL CLARIFICATION**:
- Normal state transitions (monitoring → spawning reviewers, reviews → fixes, etc.) do NOT require stops
- Only the SPECIFIC transitions listed above require checkpoints
- When in doubt about NORMAL operations: DON'T STOP - CONTINUE (set CONTINUE-SOFTWARE-FACTORY=TRUE)
- When in doubt about EXCEPTIONAL situations: STOP (set CONTINUE-SOFTWARE-FACTORY=FALSE)

**DEFAULT BEHAVIOR**: Unless explicitly listed above, transitions should continue automatically for full automation.