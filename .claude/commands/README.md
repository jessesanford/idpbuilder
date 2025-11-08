# Software Factory 2.0 Commands

> **Quick Start**: See `COMMANDS-QUICK-REFERENCE.md` for the essential commands you need.

This directory contains slash commands for the Software Factory. Most users only need a few primary commands.

## Primary Commands (Frequently Used)

### 🎯 Core Operations

#### `/continue-orchestrating`
**Purpose**: Main orchestrator command for coordinating all Software Factory work
**When to Use**: Starting or resuming orchestration
**Note**: Never writes code, only coordinates other agents

#### `/fix-cascade`
**Purpose**: Generic fix cascade for any identified issue (R375 compliant)
**When to Use**: After creating a `*-FIX-PLAN.md` document
**Note**: Auto-detects fix plans and manages dual state files

#### `/pr-ready-transform`
**Purpose**: Transform Software Factory effort branches into PR-ready branches
**When to Use**: Preparing branches for upstream pull requests
**Note**: Removes all SF artifacts and consolidates commits

#### `/integration`
**Purpose**: Execute integration sub-state machine for merging branches
**When to Use**: Integrating multiple effort branches (WAVE/PHASE/PROJECT level)
**Parameters**: `type=WAVE|PHASE|PROJECT`, `branches=effort1,effort2`, `target=integration-branch`, `validation=BASIC|FULL|COMPREHENSIVE`
**Note**: Handles complex merge scenarios with automatic cycle management

#### `/splitting`
**Purpose**: Execute splitting sub-state machine for oversized efforts
**When to Use**: When effort exceeds size limits and needs to be split
**Note**: Creates and manages sequential split branches

### 🔧 Agent Continuations

#### `/continue-implementing`
**Purpose**: Continue as Software Engineer
**When to Use**: When manually implementing (usually spawned by orchestrator)

#### `/continue-reviewing`
**Purpose**: Continue as Code Reviewer
**When to Use**: When manually reviewing (usually spawned by orchestrator)

#### `/continue-architecting`
**Purpose**: Continue as Architect
**When to Use**: When manually architecting (usually spawned by orchestrator)

### 📊 Utilities

#### `/check-status`
**Purpose**: Check current system status and health
**When to Use**: Diagnosing issues or checking progress

#### `/reset-state`
**Purpose**: Reset corrupted state files
**When to Use**: Only when state is unrecoverable
**Levels**: 1 (TODOs), 2 (State), 3 (Full)

## Archived Commands

Rarely used or deprecated commands have been moved to `.claude/commands/archived/`:
- Legacy cascade orchestration commands
- Phase-specific spawn commands
- Manual integration commands
- Superseded fix commands

These are preserved for reference but should not be used in normal operations.

## Typical Usage Patterns

### Standard Development Flow
```bash
/continue-orchestrating   # Orchestrator coordinates everything
# Agents are spawned automatically as needed
```

### Fix Cascade Flow
```bash
# 1. Document the issue
vim BUGNAME-FIX-PLAN.md

# 2. Execute fix cascade
/fix-cascade              # Auto-detects plan and manages fix
```

### PR Preparation Flow
```bash
/pr-ready-transform       # Clean up branches for upstream PRs
```

### Recovery Flow
```bash
/check-status            # Diagnose the problem
/reset-state             # Reset if needed (choose level)
/continue-orchestrating  # Resume work
```

## Key Points

- **Most work** is done through `/continue-orchestrating`
- **Fixes** use the generic `/fix-cascade` command
- **PR prep** uses `/pr-ready-transform`
- **Other agent commands** are usually spawned automatically
- **Archived commands** are in `archived/` subdirectory

For detailed documentation on each command, see the individual command files.