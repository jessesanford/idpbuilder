# State Machine Migration Guide

## ⚠️ CRITICAL: State Machine Authority

Per **Rule R206**, there is only ONE authoritative state machine:
- **SOFTWARE-FACTORY-STATE-MACHINE.md** - SINGLE SOURCE OF TRUTH

## Legacy System (DEPRECATED)

The `state-machines/` directory contains the OLD state machine system:
- More detailed states (e.g., SPAWN_CODE_REVIEWER_PLANNING vs SPAWN_AGENTS)
- Agent-specific files
- Different state names and transitions

## Current System (AUTHORITATIVE)

The `SOFTWARE-FACTORY-STATE-MACHINE.md` contains:
- Simplified, unified states
- All agents in one file
- R206 compliance required
- New architecture-driven states (R210, R211)

## Migration Rules

1. **ALWAYS** validate against SOFTWARE-FACTORY-STATE-MACHINE.md
2. **NEVER** use state-machines/*.md for validation
3. **Map** old detailed states to new simplified states
4. **Update** any references to use new state names

## State Mapping Examples

### Orchestrator
- OLD: `SPAWN_CODE_REVIEWER_PLANNING` → NEW: `SPAWN_AGENTS`
- OLD: `SPAWN_SW_ENG` → NEW: `SPAWN_AGENTS`
- OLD: `SPAWN_SW_ENG_FIX` → NEW: `SPAWN_AGENTS`
- OLD: `CREATE_SPLIT_PLAN` → NEW: `SPAWN_AGENTS` (for Code Reviewer)
- NEW: `SPAWN_ARCHITECT_PHASE_PLANNING` (R210, no old equivalent)
- NEW: `SPAWN_CODE_REVIEWER_WAVE_IMPL` (R211, no old equivalent)

### Code Reviewer
- OLD: `PLANNING` → NEW: `EFFORT_PLAN_CREATION`
- NEW: `PHASE_IMPLEMENTATION_PLANNING` (R211, no old equivalent)
- NEW: `WAVE_IMPLEMENTATION_PLANNING` (R211, no old equivalent)
- NEW: `WAVE_DIRECTORY_ACKNOWLEDGMENT` (R214, no old equivalent)

### Architect
- OLD: `REVIEW` → NEW: `WAVE_REVIEW` or `PHASE_ASSESSMENT`
- NEW: `PHASE_ARCHITECTURE_PLANNING` (R210, no old equivalent)
- NEW: `WAVE_ARCHITECTURE_PLANNING` (R210, no old equivalent)

## Action Required

All agents must:
1. Read SOFTWARE-FACTORY-STATE-MACHINE.md on startup
2. Validate state transitions against it
3. Ignore state-machines/*.md files
4. Use new state names in all transitions
