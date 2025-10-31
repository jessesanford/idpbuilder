# Software Factory Commands - Quick Reference

## 🎯 Primary Commands (Use These)

### Development & Coordination
```bash
/continue-orchestrating   # Main orchestrator - coordinates all work
/fix-cascade             # Fix any identified issues across branches
/pr-ready-transform      # Transform effort branches to PR-ready state
/integration             # Integration sub-state for merging branches
/splitting               # Splitting sub-state for oversized efforts
```

### Agent Continuations
```bash
/continue-implementing   # Continue as software engineer
/continue-reviewing      # Continue as code reviewer
/continue-architecting   # Continue as architect
```

### Utilities
```bash
/check-status           # Check project status and health
/reset-state            # Reset when state is corrupted
```

## 📁 Active Commands Location
All active commands are in: `.claude/commands/`

## 📦 Archived Commands
Rarely used or deprecated commands moved to: `.claude/commands/archived/`
- Legacy cascade commands
- Phase-specific spawn commands
- Manual integration commands

## 🔄 Typical Workflows

### Normal Development
```bash
/continue-orchestrating   # Start/resume orchestrator
# Orchestrator spawns other agents as needed
```

### Fixing Issues
```bash
# 1. Create fix plan
vim BUGNAME-FIX-PLAN.md

# 2. Run fix cascade
/fix-cascade
```

### Preparing PRs
```bash
/pr-ready-transform      # Clean up branches for PRs
```

### Recovery
```bash
/check-status            # Diagnose issues
/reset-state             # Reset if needed
/continue-orchestrating  # Resume work
```