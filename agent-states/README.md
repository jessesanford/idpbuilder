# Agent States Directory Structure

This directory contains state-specific rules for Software Factory agents.

## Structure

```
agent-states/
├── software-factory/   # SF 3.0 main implementation states
│   ├── orchestrator/  (121 states)
│   ├── sw-engineer/   (9 states)
│   ├── code-reviewer/ (19 states)
│   ├── architect/     (10 states)
│   └── integration/   (7 states)
├── state-manager/      # State Manager consultation states
├── initialization/     # Initialization workflow states
├── integration/        # Integration agent states
├── pr-ready/           # PR-Ready transformation states
├── splitting/          # Code splitting workflow states
├── fix-cascade/        # Fix cascade management states
├── experimental/       # Experimental features
└── code-reviewer/      # LEGACY SF 2.0 (3 old states, deprecated)
```

## SF 3.0 Path Convention

**All SF 3.0 states use the full path:**
- `agent-states/software-factory/{agent}/{STATE}/rules.md`

**Examples:**
- `agent-states/software-factory/orchestrator/SETUP_WAVE_INFRASTRUCTURE/rules.md`
- `agent-states/software-factory/sw-engineer/IMPLEMENTATION/rules.md`
- `agent-states/software-factory/code-reviewer/CODE_REVIEW/rules.md`
- `agent-states/software-factory/architect/REVIEW_WAVE_ARCHITECTURE/rules.md`

## State Rule Files

Each state directory contains a `rules.md` file that defines:
- State-specific rules and requirements
- Entry conditions
- Exit conditions
- Required actions
- Prohibited actions

## Usage

Agents load state-specific rules according to R203:
1. Determine current state from context
2. Load rules from: `agent-states/{agent}/{STATE}/rules.md`
3. Acknowledge all loaded rules
4. Execute state-specific operations

## State Machines

Different state machines have their own subdirectories:
- `main/` - Primary Software Factory 2.0 implementation
- `integration/` - Integration flow states
- `pr-ready/` - PR readiness validation
- `initialization/` - Project initialization
- `splitting/` - Code splitting workflow
- `fix-cascade/` - Fix cascade management