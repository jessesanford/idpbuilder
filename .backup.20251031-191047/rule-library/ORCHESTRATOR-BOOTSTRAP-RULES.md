# ORCHESTRATOR BOOTSTRAP RULES - MINIMAL ESSENTIAL SET

## Purpose
This file contains the ABSOLUTE MINIMUM rules needed for orchestrator bootstrap.
All other rules are loaded contextually based on state per R203.

## 🔴🔴🔴 ESSENTIAL BOOTSTRAP RULES (5 TOTAL) 🔴🔴🔴

These 5 rules are the ONLY rules that MUST be loaded at orchestrator startup.
All other rules are loaded based on current state.

### 1. R203 - State-Aware Startup Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R203-state-aware-agent-startup.md`
**Purpose**: Defines HOW to determine state and load state-specific rules
**Why Bootstrap**: Without this, agent doesn't know how to load other rules

### 2. R006 - Orchestrator Never Writes Code (BLOCKING)  
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Purpose**: Core identity - orchestrator is coordinator, not developer
**Why Bootstrap**: Fundamental constraint that applies to ALL states

### 3. R319 - Orchestrator Never Measures Code (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Purpose**: Core identity - orchestrator delegates measurement to reviewers
**Why Bootstrap**: Fundamental constraint that applies to ALL states

### 4. R322 - Mandatory Stop Before State Transitions (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
**Purpose**: Checkpoint control - MUST stop and await continuation
**Why Bootstrap**: Applies to EVERY state transition regardless of state

### 5. R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Purpose**: Maintain state persistence across transitions
**Why Bootstrap**: Needed to save state before stopping per R322

## STATE-SPECIFIC RULE LOADING

After loading these 5 bootstrap rules, the orchestrator:

1. **Determines current state** using R203 protocol:
   - Check if orchestrator-state-v3.json exists
   - If yes: read current_state field
   - If no: default to INIT state

2. **Loads state-specific rules** from:
   ```
   $CLAUDE_PROJECT_DIR/agent-states/software-factory/orchestrator/{STATE}/rules.md
   ```

3. **State-specific rules include:**
   - Rules only needed in that state
   - State transition requirements
   - State-specific validation

## RULES MOVED TO STATE-SPECIFIC LOCATIONS

### INIT State Only:
- R281 - Initial State File Creation
- R191 - Target Repository Configuration  
- R192 - Repository Separation Protocol

### SPAWN States Only:
- R151 - Parallel Spawning Requirements
- R208 - CD Before Spawn Protocol
- R218 - Parallel Code Reviewer Spawning
- R268 - Integration Agent Spawn Protocol

### Transition Points Only:
- R234 - Mandatory State Traversal
- R290 - State Rule Reading Verification
- R217 - Post-Transition Re-acknowledgment

### Specific Work States:
- R307/R308 - Branching strategies (WAVE_COMPLETE)
- R304 - Line Counter Tool (never for orchestrator)
- R287 - TODO Persistence (loaded when TODOs used)

## STARTUP SEQUENCE WITH BOOTSTRAP

```bash
1. Load these 5 bootstrap rules
2. Determine current state (R203)
3. Load state-specific rules
4. Acknowledge all loaded rules
5. Execute state work
6. Stop before transition (R322)
7. Update state file (R288)
8. Wait for continuation
```

## CONTEXT SAVINGS

**Before**: 23 rules loaded at startup (~40-50% context)
**After**: 5 rules loaded at startup (~10% context)
**Savings**: 75% reduction in bootstrap overhead

## IMPORTANT NOTES

1. **Supreme Laws** are distributed:
   - Some in bootstrap (R322, R288)
   - Others in state files where relevant
   - All still enforced, just loaded contextually

2. **Safety preserved**:
   - Core identity rules always loaded
   - Stop protocol always active
   - State persistence guaranteed

3. **Extensibility**:
   - New states just need their own rules.md
   - No changes to bootstrap needed
   - Clean separation of concerns

---
*Bootstrap Rules v2.0 - Optimized for Context Efficiency*
*Last Updated: 2025-09-06*