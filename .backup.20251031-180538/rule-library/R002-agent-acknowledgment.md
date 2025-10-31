# 🚨 RULE R002 - Agent Acknowledgment Protocol

**Criticality:** CRITICAL - Required for all agents  
**Grading Impact:** -20% for missing acknowledgment  
**Enforcement:** IMMEDIATE - Must occur on every startup

## Rule Statement

ON EVERY STARTUP, EACH AGENT MUST OUTPUT a formal acknowledgment of all applicable rules.

## Mandatory Acknowledgment Format

```bash
================================
RULE ACKNOWLEDGMENT
I am {agent-name} in state {CURRENT_STATE}
I acknowledge these CRITICAL rules:
--------------------------------
[List all CRITICAL and BLOCKING rules with format:]
R###: Rule description [CRITICALITY]
================================
```

## Required Elements

### 1. Agent Identification
- Agent name (orchestrator, sw-engineer, code-reviewer, architect)
- Current state from state machine
- Timestamp of acknowledgment

### 2. Rule Categories to Acknowledge
Each agent MUST acknowledge rules from:
- **SUPREME LAWS** (if any)
- **CRITICAL** rules (automatic failure if violated)
- **BLOCKING** rules (must complete before proceeding)
- **State-specific** rules for current state

### 3. Source Locations
Agents must read and acknowledge rules from:
- `.claude/agents/{agent-name}.md` - Core agent configuration
- `agent-states/{agent-name}/{STATE}/rules.md` - State-specific rules
- `rule-library/R*.md` - Referenced rule files
- `.claude/CLAUDE.md` - Project overrides

## Agent-Specific Requirements

### Orchestrator
Must acknowledge:
- R006: Never writes code
- R151: Parallel spawn timing
- R288: State file updates
- R287: TODO persistence
- State-specific rules

### SW Engineer
Must acknowledge:
- R152: Implementation speed metrics
- R200: Measure only changeset
- R287: TODO persistence
- State-specific rules

### Code Reviewer
Must acknowledge:
- R153: Review turnaround metrics
- R219: Dependency-aware planning
- R287: TODO persistence
- State-specific rules

### Architect
Must acknowledge:
- R158: Pattern compliance
- R210: Architecture planning
- State-specific rules

## Verification

```bash
# Check acknowledgment occurred
verify_acknowledgment() {
    local agent="$1"
    local log="$2"
    
    if ! grep -q "RULE ACKNOWLEDGMENT" "$log"; then
        echo "❌ VIOLATION: R002 - No acknowledgment from $agent"
        return 1
    fi
    
    if ! grep -q "I am $agent in state" "$log"; then
        echo "❌ VIOLATION: R002 - Invalid acknowledgment format"
        return 1
    fi
    
    echo "✅ R002: Valid acknowledgment from $agent"
}
```

## Grading Impact

- Missing acknowledgment: -20%
- Incomplete acknowledgment: -10%
- Wrong format: -5%
- Pattern of violations: -50%

## Example Compliance

### GOOD: Complete acknowledgment
```
================================
RULE ACKNOWLEDGMENT
I am orchestrator in state WAVE_START
I acknowledge these CRITICAL rules:
--------------------------------
R006: Never write code [CRITICAL]
R151: Parallel spawn <5s [CRITICAL]  
R288: Update state file on transitions [BLOCKING]
R287: Comprehensive TODO persistence [BLOCKING]
R203: State-aware startup [BLOCKING]
================================
```

### BAD: Missing acknowledgment
```
Starting orchestration work...
[No acknowledgment output]
```

---
**Remember:** Acknowledgment demonstrates you understand your responsibilities. No acknowledgment = No trust = No job.
## Software Factory 3.0 Integration

**State Tracking**: In SF 3.0, state transitions are tracked in `orchestrator-state-v3.json`:
```json
{
  "state_machine": {
    "current_state": "CURRENT_STATE_NAME",
    "previous_state": "PREVIOUS_STATE_NAME",
    "state_history": [...]
  }
}
```

**Compliance**: This rule applies to SF 3.0 state machine with appropriate state name mappings per R516 naming conventions.

**Reference**: See `docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md` Part 2 for state machine design.

