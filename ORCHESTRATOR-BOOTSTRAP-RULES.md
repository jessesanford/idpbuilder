# ORCHESTRATOR MINIMAL BOOTSTRAP RULES

## 🔴🔴🔴 CRITICAL: THIS IS THE MINIMAL SET FOR ORCHESTRATOR STARTUP 🔴🔴🔴

This file contains the ABSOLUTE MINIMUM rules the orchestrator needs to:
1. Determine its current state
2. Know to load state-specific rules
3. Validate state transitions
4. Persist state changes
5. Protect critical boundaries

## Purpose

The orchestrator has ZERO MEMORY between states. Each time it starts:
- It reads this minimal bootstrap
- Determines current state from orchestrator-state.json
- Loads state-specific rules for that state
- Has exactly what it needs, no more, no less

## 🔴🔴🔴 ESSENTIAL BOOTSTRAP RULES (9 Rules + 1 Reference) 🔴🔴🔴

### Core Functionality Rules (7 Essential)

1. **R283** - COMPLETE FILE READING (SUPREME LAW)
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R283-COMPLETE-FILE-READING-SUPREME-LAW.md`
   - Purpose: Read EVERY line of EVERY file - no partial reads
   - Why Bootstrap: Needed to properly read all other rules and files

2. **R290** - STATE RULE READING AND VERIFICATION (SUPREME LAW)
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R290-state-rule-reading-verification-supreme-law.md`
   - Purpose: Defines HOW to read state rules and create verification markers
   - Why Bootstrap: Without this, orchestrator won't know to read state rules

3. **R203** - STATE-AWARE AGENT STARTUP
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R203-state-aware-agent-startup.md`
   - Purpose: Defines the startup sequence and state determination
   - Why Bootstrap: Core identity and startup protocol

4. **R206** - STATE MACHINE TRANSITION VALIDATION
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R206-state-machine-transition-validation.md`
   - Purpose: Validate EVERY state transition against state machine
   - Why Bootstrap: Needed in all states to validate transitions

5. **R288** - STATE FILE UPDATE AND COMMIT PROTOCOL (SUPREME LAW)
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Purpose: Update orchestrator-state.json on EVERY transition
   - Why Bootstrap: Every state needs to update and persist state

6. **R287** - TODO PERSISTENCE COMPREHENSIVE
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Purpose: Save/commit/recover TODOs to survive compaction
   - Why Bootstrap: Needed for recovery in any state

7. **R322** - MANDATORY STOP BEFORE STATE TRANSITIONS (SUPREME LAW)
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Purpose: MUST STOP and summarize before EVERY state transition
   - Why Bootstrap: Critical for the stop-and-restart model

### Critical Protection Rules (2 Essential)

8. **R309** - NEVER CREATE EFFORTS IN SF REPO (SUPREME LAW)
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`
   - Purpose: NEVER create effort/wave/split branches in Software Factory repo
   - Why Bootstrap: Prevents catastrophic repository corruption

9. **R006** - ORCHESTRATOR NEVER WRITES CODE (BLOCKING)
   - Path: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Purpose: Orchestrator NEVER writes, measures, or reviews code - delegates ALL
   - Why Bootstrap: Core identity rule - orchestrator is a coordinator ONLY

### State Machine Reference (1 Reference)

10. **SOFTWARE-FACTORY-STATE-MACHINE.md** - STATE MACHINE AUTHORITY
    - Path: `$CLAUDE_PROJECT_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md`
    - Purpose: The ABSOLUTE authority on all valid states and transitions
    - Why Bootstrap: Referenced by R206 for all transition validation

## 🔴🔴🔴 STARTUP SEQUENCE WITH MINIMAL BOOTSTRAP 🔴🔴🔴

```yaml
orchestrator_startup:
  1_read_bootstrap:
    - Read this file (ORCHESTRATOR-BOOTSTRAP-RULES.md)
    - Read all 9 rule files listed above
    - Read SOFTWARE-FACTORY-STATE-MACHINE.md
    - Total: 10 files to read
    
  2_determine_state:
    - Read orchestrator-state.json (if exists)
    - Extract current_state field
    - If no file, state = INIT
    
  3_load_state_rules:
    - Per R290: Read agent-states/orchestrator/{STATE}/rules.md
    - Create verification marker: .state_rules_read_orchestrator_{STATE}
    - Acknowledge all state-specific rules
    
  4_begin_state_work:
    - Check marker exists (R290)
    - Execute state-specific actions
    - Never write code (R006)
    - Never create efforts in SF repo (R309)
    
  5_state_transition:
    - Validate with R206 and state machine
    - Update state file (R288)
    - Save TODOs (R287)
    - STOP per R322
    - Wait for /continue-orchestrating
```

## 🔴🔴🔴 WHAT'S NOT HERE (MOVED TO STATES) 🔴🔴🔴

These rules were moved to specific states where they're actually needed:

- **R234** (Mandatory State Traversal) → Only in sequence states
- **R208** (CD Before Spawn) → Only in SPAWN_* states
- **R151** (Parallel Spawning Timing) → Only in SPAWN_* states
- **R281** (State File Initialization) → Only in INIT state
- **R221** (Bash Directory Reset) → Only in states using bash
- **R235** (Pre-Flight Verification) → Only in work states
- **R280** (Main Branch Protection) → Only in git operation states
- **R307** (Branch Mergeability) → Only in integration states
- **R308** (Incremental Branching) → Only in SETUP_EFFORT_INFRASTRUCTURE
- **R319** (Never Measures Code) → Only in MONITOR_* states
- **R216** (Bash Syntax) → Only in states using bash
- **R321** (Immediate Backport) → Only in integration states

## 🔴🔴🔴 VERIFICATION CHECKLIST 🔴🔴🔴

Before using this bootstrap, verify:

□ All 9 rule files exist in rule-library/
□ SOFTWARE-FACTORY-STATE-MACHINE.md exists in project root
□ Each state directory has complete rules for that state
□ No state is missing critical rules it needs
□ orchestrator.md references this file correctly

## Migration Notes

**Date**: 2025-01-06
**Migrated From**: orchestrator.md (878 lines)
**Migrated To**: This file (9 rules) + state-specific rules
**Reduction**: ~77% fewer rules loaded per state
**Benefit**: Each state loads only what it needs