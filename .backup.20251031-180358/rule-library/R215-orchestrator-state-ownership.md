# R215: Orchestrator State File Ownership Protocol

**Category:** Critical Rules  
**Agents:** Orchestrator (OWNER), All Others (READ-ONLY)  
**Criticality:** MISSION CRITICAL - State corruption causes system failure  
**Priority:** HIGHEST - Prevents conflicting state updates

## 🚨 ONLY THE ORCHESTRATOR CAN UPDATE orchestrator-state-v3.json 🚨

The orchestrator is the SOLE OWNER of the state file. All other agents may READ but MUST NEVER WRITE to orchestrator-state-v3.json.

## The Problem This Solves

Without exclusive orchestrator ownership:
- Multiple agents could update state simultaneously
- State conflicts and race conditions occur
- State machine transitions become inconsistent
- Recovery becomes impossible after context loss
- Grading and tracking fail

## The Solution: Orchestrator-Only Updates

### Part 1: Orchestrator OWNS the State File

```bash
# ONLY THE ORCHESTRATOR can do this:
update_orchestrator_state() {
    local NEW_STATE="$1"
    local REASON="$2"
    
    # Only orchestrator can update
    if [[ "$AGENT_TYPE" != "orchestrator" ]]; then 
        echo "❌ FATAL: Only orchestrator can update state!"; 
        exit 1; 
    fi
    
    cat >> orchestrator-state-v3.json << EOF
current_state: $NEW_STATE
previous_state: $CURRENT_STATE
transition_time: "$(date -Iseconds)"
transition_reason: "$REASON"
EOF
}
```

### Part 2: Other Agents Report Status, Don't Update State

When other agents need to communicate status to orchestrator:

```bash
# SW ENGINEER/CODE REVIEWER/ARCHITECT - Report, don't update
report_agent_status() {
    local STATUS="$1"
    local DETAILS="$2"
    
    # Write to agent-specific status file, NOT orchestrator state
    cat > agent-status.yaml << EOF
agent: $AGENT_TYPE
status: $STATUS
details: "$DETAILS"
timestamp: "$(date -Iseconds)"
EOF
    
    # Signal orchestrator to check status
    echo "STATUS_READY: Orchestrator should check agent-status.yaml"
}
```

### Part 3: Orchestrator Reads Agent Status

```bash
# ORCHESTRATOR - Reads agent status and updates its own state
check_agent_status() {
    local AGENT="$1"
    local STATUS_FILE="efforts/phase${PHASE}/wave${WAVE}/${EFFORT}/agent-status.yaml"
    
    if [ -f "$STATUS_FILE" ]; then 
        AGENT_STATUS=$(grep "status:" "$STATUS_FILE" | cut -d: -f2- | xargs); 
        
        # Orchestrator updates its own state based on agent status
        case "$AGENT_STATUS" in 
            "COMPLETED") 
                update_orchestrator_state "WAVE_COMPLETE" "All agents completed"; 
                ;; 
            "BLOCKED") 
                update_orchestrator_state "ERROR_RECOVERY" "Agent blocked: $AGENT"; 
                ;; 
            "NEEDS_SPLIT") 
                update_orchestrator_state "SPAWN_CODE_REVIEWER_WAVE_IMPL" "Split required"; 
                ;; 
        esac; 
    fi
}
```

## Integration with R206 (State Validation)

R206 requires state validation, but agents validate their OWN states, not the orchestrator's:

```bash
# CORRECT: Agent validates its own state transition
validate_agent_state_transition() {
    local AGENT_CURRENT="$1"
    local AGENT_TARGET="$2"
    
    # Validate against software-factory-3.0-state-machine.json
    validate_state_transition "$AGENT_CURRENT" "$AGENT_TARGET" "$AGENT_TYPE"
    
    # Update AGENT's internal state (NOT orchestrator-state-v3.json)
    echo "agent_state: $AGENT_TARGET" > agent-internal-state.yaml
}
```

## What Each Agent CAN and CANNOT Do

### Orchestrator
✅ **CAN**:
- Read orchestrator-state-v3.json
- Write/update orchestrator-state-v3.json
- Add new sections to state file
- Update any field in state file
- Create state backups

❌ **CANNOT**:
- Let other agents update its state

### SW Engineer / Code Reviewer / Architect
✅ **CAN**:
- Read orchestrator-state-v3.json (for context)
- Write to agent-status.yaml
- Track their own internal state
- Signal status to orchestrator

❌ **CANNOT**:
- Write to orchestrator-state-v3.json
- Update orchestrator state sections
- Modify state transitions
- Add fields to orchestrator state

## Enforcement Examples

### ✅ Correct: Orchestrator Updates After Agent Signal
```bash
# SW Engineer completes work
echo "status: COMPLETED" > agent-status.yaml

# Orchestrator checks and updates
AGENT_STATUS=$(cat agent-status.yaml)
if [[ "$AGENT_STATUS" == *"COMPLETED"* ]]; then
    # Only orchestrator updates state
    update_orchestrator_state "PROCESS_REVIEW" "SW Eng completed"
fi
```

### ❌ Wrong: Agent Updates Orchestrator State Directly
```bash
# SW ENGINEER - WRONG!
echo "current_state: COMPLETED" >> orchestrator-state-v3.json  # VIOLATION!

# CODE REVIEWER - WRONG!
./utilities/update-orchestrator-state-section.sh ...  # VIOLATION!
```

## Migration Required

Current code that violates this rule:
1. Architect updating phase_architecture_plans in orchestrator-state-v3.json
2. Code Reviewer updating implementation_plans sections
3. Helper script allowing non-orchestrator updates

These must be changed to:
1. Agents write to their own status files
2. Orchestrator reads status and updates its state
3. Helper scripts check agent type before allowing updates

## Validation Script

```bash
#!/bin/bash
# validate-r215-state-ownership.sh

check_state_file_ownership() {
    echo "🔍 R215: Checking orchestrator state file ownership..."
    
    # Check last modification
    if [ -f "orchestrator-state-v3.json" ]; then 
        # Check git log for who modified it
        git log -1 --format="%an: %s" orchestrator-state-v3.json; 
        
        # Look for violations in agent code
        echo "Checking for violation patterns..."; 
        
        # SW Engineers shouldn't update state
        if grep -r "orchestrator-state-v3.json" efforts/*/sw-engineer.sh 2>/dev/null; then 
            echo "❌ SW Engineer attempting to update orchestrator state!"; 
        fi; 
        
        # Code Reviewers shouldn't update state  
        if grep -r "update-orchestrator-state" efforts/*/code-reviewer.sh 2>/dev/null; then 
            echo "❌ Code Reviewer attempting to update orchestrator state!"; 
        fi; 
        
        echo "✅ Only orchestrator should update orchestrator-state-v3.json"; 
    fi
}
```

## Summary

- **Orchestrator** is the SOLE OWNER of orchestrator-state-v3.json
- **Other agents** can only READ, never WRITE
- **Agent communication** happens through agent-status.yaml files
- **State transitions** are orchestrator's responsibility
- **Prevents** race conditions and state corruption
- **Ensures** single source of truth for system state