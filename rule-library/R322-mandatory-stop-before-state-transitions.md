# 🛑🛑🛑 RULE R322: MANDATORY ORCHESTRATOR CHECKPOINTS (SUPREME LAW)

**CRITICALITY**: 🔴🔴🔴 SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE

## Purpose
Enforce mandatory checkpoints at critical orchestrator transitions to:
1. Allow user review of plans before execution
2. Preserve context after spawning agents (prevents overflow)
3. Provide clear decision points for continuation
4. Ensure proper state persistence before moving forward

**THIS RULE CONSOLIDATES:**
- Former R313 (stop after spawning) - now Part A of this rule
- Checkpoint stops for user review - Part B of this rule
- Works with R324 (state update before stop) as companion rule

## 🔴🔴🔴 SUPREME LAW: STOP MEANS STOP 🔴🔴🔴

### MANDATORY STOP POINTS

The orchestrator MUST STOP and exit after these specific transitions:

#### PART A: SPAWN STATE STOPS (Context Preservation)
**After spawning ANY agents, MUST stop to prevent context overflow:**

- ALL `SPAWN_*` states → Next state
  - `SPAWN_AGENTS` → `MONITOR_IMPLEMENTATION`
  - `SPAWN_CODE_REVIEWERS_EFFORT_PLANNING` → `WAITING_FOR_EFFORT_PLANS`
  - `SPAWN_CODE_REVIEWERS_FOR_REVIEW` → `MONITOR_REVIEWS`
  - `SPAWN_ENGINEERS_FOR_FIXES` → `MONITOR_FIXES`
  - `SPAWN_INTEGRATION_AGENT` → `MONITORING_INTEGRATION`
  - `SPAWN_INTEGRATION_AGENT_PHASE` → `MONITORING_PHASE_INTEGRATION`
  - `SPAWN_INTEGRATION_AGENT_PROJECT` → `MONITORING_PROJECT_INTEGRATION`
  - `SPAWN_ARCHITECT_*` → Respective monitoring states
  - `SPAWN_CODE_REVIEWER_*` → Respective waiting states
  - `SPAWN_SW_ENGINEER_*` → Respective monitoring states
  
**Reason**: Agent outputs flood context window, causing rule amnesia

#### PART B: PLANNING → EXECUTION STOPS
**These transitions REQUIRE user review before execution:**

- `WAITING_FOR_MERGE_PLAN` → `SPAWN_INTEGRATION_AGENT`
  - User must review WAVE-MERGE-PLAN.md before execution
  - Stop allows verification of merge order and strategy
  
- `WAITING_FOR_PHASE_MERGE_PLAN` → `SPAWN_INTEGRATION_AGENT_PHASE`
  - User must review PHASE-MERGE-PLAN.md before execution
  - Stop allows verification of phase integration approach
  
- `WAITING_FOR_PROJECT_MERGE_PLAN` → `SPAWN_INTEGRATION_AGENT_PROJECT`
  - User must review PROJECT-MERGE-PLAN.md before execution
  - Stop allows verification of final project integration

- `WAITING_FOR_FIX_PLANS` → `DISTRIBUTE_FIX_PLANS`
  - User must review fix plans before distribution
  - Stop allows verification of fix approach

- `WAITING_FOR_BACKPORT_PLAN` → `SPAWN_SW_ENGINEER_BACKPORT_FIXES`
  - User must review backport plan before execution
  - Stop allows verification of backport strategy

#### PART C: ASSESSMENT → ACTION STOPS
**These transitions REQUIRE user decision before proceeding:**

- `WAITING_FOR_PHASE_ASSESSMENT` → `PHASE_COMPLETE`
  - User must acknowledge architect assessment
  - Stop allows review of phase achievements

- `WAITING_FOR_PROJECT_VALIDATION` → `CREATE_INTEGRATION_TESTING`
  - User must approve project validation results
  - Stop allows final review before testing

#### PART D: MONITORING → STATE CHANGE STOPS
**These transitions REQUIRE state persistence:**

- `MONITORING_INTEGRATION` → `WAVE_REVIEW`
  - Must save integration results
  - Stop ensures proper state capture

- `MONITORING_PHASE_INTEGRATION` → `SPAWN_ARCHITECT_PHASE_ASSESSMENT`
  - Must save phase integration results
  - Stop ensures proper state capture

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

2. **Update orchestrator-state.json with NEW state**
   ```bash
   # CRITICAL: Set current_state to NEXT state
   jq '.current_state = "NEXT_STATE"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
   jq '.previous_state = "CURRENT_STATE"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
   jq '.transition_time = "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
   ```

3. **Save TODO state per R287**
   ```bash
   # Save current TODO list
   # Include reason for stop
   save_todos "R322_CHECKPOINT_BEFORE_${NEXT_STATE}"
   ```

4. **Commit and push state changes**
   ```bash
   git add orchestrator-state.json todos/*.todo
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
        ("WAITING_FOR_FIX_PLANS", "DISTRIBUTE_FIX_PLANS"),
        ("WAITING_FOR_BACKPORT_PLAN", "SPAWN_SW_ENGINEER_BACKPORT_FIXES"),
        ("WAITING_FOR_PHASE_ASSESSMENT", "PHASE_COMPLETE"),
        ("WAITING_FOR_PROJECT_VALIDATION", "CREATE_INTEGRATION_TESTING"),
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

### Supersedes:
- **R313**: Now incorporated as Part A of this rule
- Any guidance suggesting automatic continuation at checkpoints
- Any pattern of immediate transition without checkpoint
- Previous "continue if successful" patterns at spawn points

## EXAMPLES

### ✅ CORRECT: Proper R322 Checkpoint
```bash
# In WAITING_FOR_MERGE_PLAN state
echo "✅ Merge plan created and validated"

# Update state BEFORE stopping
jq '.current_state = "SPAWN_INTEGRATION_AGENT"' orchestrator-state.json > tmp.json
mv tmp.json orchestrator-state.json

# Save TODOs
save_todos "R322_CHECKPOINT_BEFORE_SPAWN_INTEGRATION"

# Commit state
git add orchestrator-state.json todos/*.todo
git commit -m "state: R322 checkpoint - merge plan ready for review"
git push

# Display checkpoint
cat << 'EOF'
## 🛑 R322 STATE TRANSITION CHECKPOINT

### ✅ Merge Plan Created:
- Location: integration-workspace/WAVE-MERGE-PLAN.md
- Efforts: 5 branches identified
- Strategy: Sequential merge with conflict resolution

### 📊 Ready to Spawn Integration Agent
- Current State: WAITING_FOR_MERGE_PLAN ✅
- Next State: SPAWN_INTEGRATION_AGENT (ready)

### ⏸️ STOPPED - User Review Required
Please review WAVE-MERGE-PLAN.md before execution.
To proceed: /continue-orchestrating
EOF

exit 0  # STOP EXECUTION
```

### ❌ WRONG: Skipping Checkpoint
```bash
# In WAITING_FOR_MERGE_PLAN state
echo "✅ Merge plan created"

# VIOLATION: Continuing without stop!
echo "Spawning integration agent..."
/spawn integration-agent  # R322 VIOLATION = -100%
```

## AUDIT REQUIREMENTS

Every R322 checkpoint must create:
1. State file update with transition recorded
2. TODO save with checkpoint reason
3. Git commit with R322 reference
4. User-visible checkpoint message
5. Clean exit (not error)

## SUMMARY

**R322 is a SUPREME LAW that ensures:**
- Users can review plans before execution
- Context is preserved between major transitions
- State is properly saved at critical points
- System provides clear decision points
- No automatic execution of unreviewed plans

**REMEMBER:** When in doubt, STOP and checkpoint. It's better to have an extra stop than to violate R322 and fail immediately.