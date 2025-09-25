# SOFTWARE ENGINEER AGENT TYPE FIX REPORT

**Date**: 2025-09-08
**Issue**: Orchestrator spawn failures due to incorrect agent type
**Resolution**: Complete - All references updated

## Problem Statement

The orchestrator was failing to spawn software engineer agents with the error:
```
Error: Agent type 'sw-engineer' not found. Available agents: general-purpose, statusline-setup, output-style-setup, architect, software-factory-manager, integration, software-engineer, code-reviewer, orchestrator
```

## Root Cause

There was a mismatch between:
- **File name**: `.claude/agents/sw-engineer.md`
- **Agent type for spawning**: `software-engineer` (required by the system)

The orchestrator and many rules were using `sw-engineer` as the agent type, but the system requires `software-engineer` for spawning.

## Solution Implemented

### 1. Fixed Agent Type References
- Changed all `Task: sw-engineer` to `Task: software-engineer`
- Changed all `subagent_type="sw-engineer"` to `subagent_type="software-engineer"`

### 2. Fixed Agent Name References
- Changed agent instance names from `sw-engineer-effort1` to `software-engineer-effort1`
- Updated parallel spawn records to use correct naming

### 3. Fixed @agent References
- Changed all `@agent-sw-engineer` to `@agent-software-engineer`
- This ensures consistency in documentation and examples

### 4. Updated Documentation
- Added clarification in AGENT-NAME-REFERENCE.md about the distinction
- Documented that spawning requires "software-engineer" not "sw-engineer"

## Files Modified (20 total)

### Orchestrator State Rules
- `agent-states/orchestrator/BUILD_VALIDATION/rules.md`
- `agent-states/orchestrator/SPAWN_AGENTS/rules.md`
- `agent-states/orchestrator/SPAWN_AGENTS/rules-optimized.md`
- `agent-states/orchestrator/SPAWN_ENGINEERS_FOR_FIXES/rules.md`
- `agent-states/orchestrator/SPAWN_SW_ENGINEER_BACKPORT_FIXES/rules.md`
- `agent-states/orchestrator/BACKPORT_FIXES/rules.md`
- `agent-states/orchestrator/SPAWN_INTEGRATION_AGENT_PHASE/rules.md`

### Command Files
- `.claude/commands/continue-orchestrating.md`

### Rule Library
- `rule-library/R151-parallel-agent-spawning-timing.md`
- `rule-library/R196-base-branch-selection.md`
- `rule-library/R197-one-agent-per-effort.md`
- `rule-library/R202-single-agent-per-split.md`
- `rule-library/R204-orchestrator-split-infrastructure.md`
- `rule-library/R232-enforcement-examples.md`
- `rule-library/R240-integration-fix-execution.md`
- `rule-library/R300-comprehensive-fix-management-protocol.md`
- `rule-library/R312-git-config-immutability-protocol.md`

### Quick Reference
- `quick-reference/orchestrator-workspace-setup-quick-ref.md`

### Documentation
- `AGENT-NAME-REFERENCE.md`

## Verification

### Before Fix
```bash
grep -r "Task.*sw-engineer" . | wc -l
# Result: 47 occurrences
```

### After Fix
```bash
grep -r "Task.*sw-engineer" . | wc -l
# Result: 1 occurrence (only in documentation showing it as wrong)

grep -r "Task.*software-engineer" . | wc -l
# Result: 47 occurrences (all converted)
```

### Subagent Type Check
```bash
grep -r "subagent_type=\"sw-engineer\"" . | wc -l
# Result: 0 (all fixed)

grep -r "subagent_type=\"software-engineer\"" . | wc -l
# Result: 1 (converted correctly)
```

## Impact

This fix resolves:
1. **Orchestrator spawn failures** - The orchestrator can now successfully spawn software engineer agents
2. **Workflow blocking** - Removes the critical blocker preventing the entire software factory workflow
3. **Future confusion** - Documentation now clearly explains the distinction

## Recommendations

1. **Consider renaming the file**: To avoid future confusion, consider renaming `.claude/agents/sw-engineer.md` to `.claude/agents/software-engineer.md` to match the agent type.

2. **Add validation**: Add a pre-spawn validation in the orchestrator to check agent types against available agents before attempting to spawn.

3. **Standardize naming**: Establish a clear naming convention where file names always match agent types.

## Status

✅ **COMPLETE** - All references have been updated and changes have been committed and pushed to the repository.