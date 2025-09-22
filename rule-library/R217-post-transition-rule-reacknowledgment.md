# Rule R217: Post-Transition Rule Re-Acknowledgment Protocol

## 🔴🔴🔴 CRITICAL ENFORCEMENT WARNING 🔴🔴🔴

### YOUR READ TOOL CALLS ARE BEING MONITORED AND RECORDED!

**YOU WILL BE GRADED ON MAKING ALL REQUIRED READ CALLS**

- Every Read tool invocation is tracked and logged
- Missing Read calls = **-100% AUTOMATIC FAILURE**
- We are watching EVERY tool call you make
- We are monitoring compliance in real-time
- Your Read tool usage is being recorded for grading

**THIS IS NOT A SUGGESTION - THIS IS MANDATORY!**

The grading system is actively tracking whether you:
1. Use the Read tool for orchestrator.md after transitions
2. Use the Read tool for SOFTWARE-FACTORY-STATE-MACHINE.md
3. Use the Read tool for state-specific rules files
4. Use the Read tool for each critical rule file
5. Make ALL required Read calls, not just some

**FAILURE TO USE THE READ TOOL = AUTOMATIC FAILURE OF ENTIRE TASK**

## Rule Statement
The Orchestrator MUST re-read its agent configuration (orchestrator.md) AND re-acknowledge its CRITICAL and BLOCKING rules after EVERY state transition to prevent context drift and rule forgetting. This acknowledgment must occur IMMEDIATELY after updating current_state and BEFORE proceeding with new state activities.

## Criticality Level
**MANDATORY** - Forgetting critical rules leads to grading failures and system violations

## Enforcement Mechanism
- **Technical**: Inject acknowledgment sequence after every state update
- **Behavioral**: Cannot proceed to new state work until acknowledgment complete
- **Grading**: -30% for missing acknowledgments (Protocol violation)
- **MONITORING**: Read tool calls are being tracked and recorded
- **FAILURE**: Missing Read calls = -100% AUTOMATIC FAILURE

## Core Principle

```
State Transition Flow:
1. Validate transition (R206)
2. Update state file
3. RE-READ AGENT CONFIG (orchestrator.md) ← CRITICAL!
4. RE-ACKNOWLEDGE CRITICAL RULES (R217) 
5. Load new state rules (R203) - PRIMARY DIRECTIVES section
6. IMMEDIATELY PROCEED with state work (R021 - NEVER STOP!)

🚨🚨🚨 TRANSITIONS ARE CONTINUATIONS, NOT STOPS! 🚨🚨🚨
After re-acknowledging rules, IMMEDIATELY continue working!
NO "awaiting instructions" - NO "what next" - JUST CONTINUE!
```

## Detailed Requirements

### 🚨 CRITICAL: THIS REQUIRES ACTUAL TOOL USE! 🚨

**THE ORCHESTRATOR MUST USE THE READ TOOL TO ACTUALLY READ THESE FILES!**

## 🔴🔴🔴 WE ARE MONITORING YOUR COMPLIANCE 🔴🔴🔴

**ATTENTION: YOUR READ TOOL USAGE IS BEING TRACKED!**

The grading system is actively monitoring and recording:
- Every Read tool invocation you make
- The exact files you read
- The order in which you read them
- Whether you skip any required reads
- Your compliance percentage

**You will be graded on making ALL required Read calls:**
- Missing even ONE Read call = SEVERE PENALTY
- Skipping Read tool usage = -100% AUTOMATIC FAILURE
- We are watching and tracking EVERYTHING

This is NOT just echoing rule names or printing acknowledgments. The orchestrator MUST:
1. Use the Read tool to read orchestrator.md (ENTIRE CONFIG - NEW!) **[MONITORED]**
2. Use the Read tool to read SOFTWARE-FACTORY-STATE-MACHINE.md **[MONITORED]**
3. Use the Read tool to read the state-specific rules file (PRIMARY DIRECTIVES) **[MONITORED]**
4. Use the Read tool to read each rule listed in PRIMARY DIRECTIVES **[MONITORED]**
5. ONLY THEN acknowledge understanding

**WHY THIS MATTERS:**
- Context can be lost during long conversations
- State transitions are critical moments where rules matter most
- Simply echoing rule names doesn't refresh memory
- Actual file reading ensures rules are in active context
- **YOUR COMPLIANCE IS BEING MONITORED AND WILL AFFECT YOUR GRADE**

### MANDATORY Post-Transition Acknowledgment Sequence

```bash
# 🚨🚨 R217: POST-TRANSITION RULE RE-ACKNOWLEDGMENT 🚨🚨
reacknowledge_critical_rules() {
    local NEW_STATE="$1"
    local AGENT_TYPE="orchestrator"
    
    echo "════════════════════════════════════════════════════════════════"
    echo "🔄 R217: POST-TRANSITION RULE RE-ACKNOWLEDGMENT"
    echo "New State: $NEW_STATE"
    echo "Time: $(date '+%Y-%m-%d %H:%M:%S %Z')"
    echo "════════════════════════════════════════════════════════════════"
    
    # CRITICAL: YOU MUST USE THE READ TOOL TO ACTUALLY READ THESE FILES!
    echo "🚨 R217 REQUIRES ACTUAL FILE READING - NOT JUST ECHOING!"
    echo ""
    
    # Get project root (Software Factory 2.0 root)
    local SF_ROOT="${SF_ROOT:-/workspaces/software-factory-2.0-template}"
    
    # STEP 1: RE-READ SUPREME LAW
    echo "🔴🔴🔴 RE-READING SUPREME LAW 🔴🔴🔴"
    echo "MUST READ WITH READ TOOL: ${SF_ROOT}/SOFTWARE-FACTORY-STATE-MACHINE.md"
    # 🚨 USE READ TOOL: ${SF_ROOT}/SOFTWARE-FACTORY-STATE-MACHINE.md
    echo "After reading, I acknowledge: State Machine is ABSOLUTE"
    
    # STEP 2: RE-READ CRITICAL RULE FILES
    echo ""
    echo "🚨🚨🚨 RE-READING MISSION CRITICAL RULES 🚨🚨🚨"
    echo "MUST READ THESE RULE FILES WITH READ TOOL:"
    echo "1. ${SF_ROOT}/rule-library/R208-orchestrator-spawn-cd-protocol.md"
    echo "2. ${SF_ROOT}/rule-library/R209-effort-directory-isolation.md"
    echo "3. ${SF_ROOT}/rule-library/R151-parallel-agent-spawning-timing.md"
    echo "4. ${SF_ROOT}/rule-library/R176-workspace-isolation.md"
    echo "5. ${SF_ROOT}/rule-library/R007-size-limit-enforcement.md"
    echo "6. ${SF_ROOT}/rule-library/R218-orchestrator-parallel-code-reviewer-spawning.md"
    # 🚨 USE READ TOOL FOR EACH FILE PATH ABOVE!
    
    # STEP 3: RE-READ STATE-SPECIFIC RULES
    echo ""
    echo "📋 RE-READING STATE-SPECIFIC RULES FOR: $NEW_STATE"
    local STATE_RULES_FILE="${SF_ROOT}/agent-states/orchestrator/$NEW_STATE/rules.md"
    if [ -f "$STATE_RULES_FILE" ]; then
        echo "MUST READ WITH READ TOOL: $STATE_RULES_FILE"
        # 🚨 USE READ TOOL: $STATE_RULES_FILE
    else
        echo "⚠️ No state-specific rules file found at: $STATE_RULES_FILE"
    fi
    
    
    # STEP 4: CONFIRM ACTUAL READING
    echo ""
    echo "📚 ACKNOWLEDGMENT AFTER READING:"
    echo "✓ I have RE-READ the State Machine (not just echoed it)"
    echo "✓ I have RE-READ the critical rule files (not just listed them)"
    echo "✓ I have RE-READ the state-specific rules (not just mentioned them)"
    echo "✓ My memory is ACTUALLY refreshed with rule content"
    
    echo ""
    echo "════════════════════════════════════════════════════════════════"
    echo "✅ R217: CRITICAL RULES RE-READ AND RE-ACKNOWLEDGED!"
    echo "════════════════════════════════════════════════════════════════"
    
    return 0
}
```

### Integration with State Transition Function

```bash
# UPDATED transition function with R217 and R230
perform_state_transition() {
    local CURRENT_STATE="$1"
    local NEW_STATE="$2"
    
    # STEP 1: Validate transition (R206)
    validate_state_transition "$CURRENT_STATE" "$NEW_STATE" "orchestrator"
    if [ $? -ne 0 ]; then
        echo "❌ Transition blocked by R206"
        exit 1
    fi
    
    # STEP 2: Update state file
    sed -i "s/current_state:.*/current_state: $NEW_STATE/" orchestrator-state.json
    echo "✓ State updated: $CURRENT_STATE → $NEW_STATE"
    
    # STEP 3: 🎯 VISUALIZE STATE POSITION (R230) 🎯
    display_state_machine_position "$NEW_STATE" "$CURRENT_PHASE" "$CURRENT_WAVE"
    
    # STEP 4: 🚨 RE-ACKNOWLEDGE CRITICAL RULES (R217) 🚨
    reacknowledge_critical_rules "$NEW_STATE"
    
    # STEP 5: Load new state rules (R203)
    if [ -f "agent-states/orchestrator/$NEW_STATE/rules.md" ]; then
        echo "Loading state-specific rules for $NEW_STATE..."
        # Read and process state rules
    fi
    
    # STEP 6: IMMEDIATELY proceed with new state work (R021 - NEVER STOP!)
    echo "🚀 IMMEDIATELY CONTINUING with $NEW_STATE activities!"
    echo "⚠️ R021 REMINDER: State transitions are waypoints, not destinations!"
    echo "🔴🔴🔴 VIOLATION WARNING: Stopping after transition = AUTOMATIC FAILURE!"
    echo "NOW EXECUTING $NEW_STATE WORK WITHOUT DELAY..."
}
```

### Acknowledgment Checklist for Common Transitions

```yaml
# Common state transitions and their critical rule focus
transition_acknowledgments:
  INIT_to_WAVE_START:
    - R206: State machine validation
    - R208: Spawn directory protocol
    - R151: Parallel spawn requirement
    
  WAVE_START_to_SPAWN_AGENTS:
    - R151: PARALLEL SPAWNING (CRITICAL!)
    - R208: Change directory before spawn
    - R209: Inject metadata into plans
    - R176: Workspace isolation
    
  SPAWN_AGENTS_to_MONITOR:
    - R007: Size limit monitoring
    - Progress tracking every 5 messages
    - Agent timeout detection
    
  MONITOR_to_WAVE_COMPLETE:
    - Integration branch requirements
    - Wave review preparation
    - R215: State file ownership
    
  WAVE_COMPLETE_to_INTEGRATION:
    - Branch creation protocol
    - Merge verification
    - Architect review preparation
```

## Common Violations to Avoid

### ❌ Forgetting Rules After Transition
```bash
# WRONG - Transition without re-acknowledgment
update_state "current_state" "SPAWN_AGENTS"
# Immediately starts spawning without refreshing rule memory
spawn_agents()  # Might forget R151 parallel requirement!
```

### ❌ Partial Acknowledgment
```bash
# WRONG - Only acknowledging some rules
echo "Transitioned to $NEW_STATE"
echo "Remember to validate states"  # Too vague!
```

### ✅ Correct Implementation with ACTUAL READING
```bash
# RIGHT - Full re-acknowledgment with TOOL USE
perform_state_transition "WAVE_START" "SPAWN_AGENTS"
# Inside reacknowledge_critical_rules():
# 1. READS ${SF_ROOT}/SOFTWARE-FACTORY-STATE-MACHINE.md (via Read tool)
# 2. READS ${SF_ROOT}/rule-library/R151-parallel-agent-spawning-timing.md (via Read tool)
# 3. READS ${SF_ROOT}/rule-library/R208-orchestrator-spawn-cd-protocol.md (via Read tool)
# 4. READS ${SF_ROOT}/rule-library/R218-orchestrator-parallel-code-reviewer-spawning.md (via Read tool)
# 5. READS ${SF_ROOT}/agent-states/orchestrator/SPAWN_AGENTS/rules.md (via Read tool)
# 6. Memory ACTUALLY refreshed with rule content!
```

### 📚 EXPLICIT FILE PATHS FOR READ TOOL

**When the orchestrator acknowledges rules, it MUST use these EXACT paths with the Read tool:**

```yaml
supreme_law:
  path: ${SF_ROOT}/SOFTWARE-FACTORY-STATE-MACHINE.md
  description: The ABSOLUTE authority on all states and transitions

critical_rules:
  - path: ${SF_ROOT}/rule-library/R151-parallel-agent-spawning-timing.md
    criticality: 50% of grade - parallel spawn timing
  
  - path: ${SF_ROOT}/rule-library/R208-orchestrator-spawn-cd-protocol.md
    criticality: MISSION CRITICAL - spawn directory protocol
  
  - path: ${SF_ROOT}/rule-library/R209-effort-directory-isolation.md
    criticality: MISSION CRITICAL - effort isolation
  
  - path: ${SF_ROOT}/rule-library/R218-orchestrator-parallel-code-reviewer-spawning.md
    criticality: MANDATORY - wave plan parallelization reading
  
  - path: ${SF_ROOT}/rule-library/R176-workspace-isolation.md
    criticality: BLOCKING - workspace isolation
  
  - path: ${SF_ROOT}/rule-library/R007-size-limit-enforcement.md
    criticality: BLOCKING - 800 line limit

state_specific:
  pattern: ${SF_ROOT}/agent-states/orchestrator/${STATE}/rules.md
  example: ${SF_ROOT}/agent-states/orchestrator/SPAWN_AGENTS/rules.md
```

**Where SF_ROOT is typically: /workspaces/software-factory-2.0-template**

## Grading Impact

### 🔴🔴🔴 READ TOOL COMPLIANCE IS MONITORED 🔴🔴🔴

**YOUR READ TOOL USAGE IS BEING TRACKED AND GRADED:**
- **Missing ANY required Read tool call**: -100% AUTOMATIC FAILURE
- **We are monitoring EVERY Read tool invocation**
- **The grading system tracks your compliance in real-time**
- **You cannot fake compliance - we see everything**

**Standard Penalties:**
- **No acknowledgment after transition**: -30% (Protocol violation)
- **Partial acknowledgment**: -15% (Incomplete compliance)
- **Forgetting R151 due to no refresh**: -50% (Critical failure)
- **Wrong directory spawn after transition**: -20% (Workspace violation)
- **FAILURE TO USE READ TOOL**: -100% (AUTOMATIC FAILURE)

## Integration with Other Rules

- **R206**: Validates transition is legal
- **R217**: Re-acknowledges rules AFTER transition (THIS RULE)
- **R203**: Loads new state-specific rules
- **R151**: Critical parallel spawn rule (must not forget!)
- **R208**: Directory protocol (must not forget!)

## Implementation Priority

**IMMEDIATE** - Add to orchestrator configuration NOW to prevent:
- Forgetting R151 parallel spawn requirement (50% grade loss!)
- Forgetting R208 directory protocol (mission critical!)
- Context drift during long orchestration sessions

## Summary

### 🔴🔴🔴 FINAL WARNING: READ TOOL MONITORING ACTIVE 🔴🔴🔴

**THE GRADING SYSTEM IS WATCHING YOUR READ TOOL USAGE!**

We are actively monitoring and will grade you on:
- Using the Read tool for EVERY required file
- Making ALL Read calls, not just some
- Actually invoking the Read tool, not just echoing file names
- Full compliance with R217 requirements

**Missing Read tool calls = -100% AUTOMATIC FAILURE**

**R217 ensures the orchestrator maintains perfect rule memory by:**
1. **ACTUALLY READING** (not echoing) the State Machine after every transition **[MONITORED]**
2. **ACTUALLY READING** critical rule files via the Read tool **[MONITORED]**
3. **ACTUALLY READING** state-specific rules for the new state **[MONITORED]**
4. Refreshing context with real rule content (not just names)
5. Creating an audit trail of actual file reads **[WE ARE TRACKING THIS]**
6. **IMMEDIATELY CONTINUING** work after reading (no stopping!)
7. **TREATING TRANSITIONS AS WAYPOINTS** not destinations

**REMEMBER:** This rule requires ACTUAL TOOL USE. The orchestrator MUST use the Read tool to read these files, not just print their names. This prevents the #1 cause of orchestrator failures: forgetting critical rules after state transitions due to context loss.

**YOUR COMPLIANCE IS BEING MONITORED. WE SEE EVERY TOOL CALL. YOU WILL BE GRADED ON THIS.**