# Sub-Orchestrator Simplification Report

## Executive Summary

The sub-orchestrator pattern has been simplified to use existing commands in background mode rather than creating duplicate commands. This reduces complexity and maintenance burden while preserving all functionality.

## Changes Made

### 1. Archived Duplicate Commands
The following redundant commands have been archived:
- `/sub-orchestrate-fix` → Use `/fix-cascade` instead
- `/sub-orchestrate-split` → Use `/splitting` instead
- `/sub-orchestrate-integration` → No replacement (integration not implemented)

### 2. Updated Rules

#### R377 - Master-Sub-Orchestrator Communication
- **Before**: Referenced `/sub-orchestrate-[type]` commands with parameter files
- **After**: Uses existing commands directly with inline parameters
- **Example**:
  ```bash
  # OLD WAY
  claude -p . --command "/sub-orchestrate-fix" --params "file=/tmp/params.json"

  # NEW WAY
  claude -p . --command "/fix-cascade" --params "fix_id=bug-123 branches=b1,b2,b3"
  ```

#### R378 - Sub-Orchestrator Lifecycle Management
- **Before**: Referenced SUB_ORCHESTRATOR signals
- **After**: Uses SUB_PROCESS signals (more generic)
- **Rationale**: Commands are sub-processes, not separate orchestrators

#### R379 - Sub-Process Monitoring & Recovery
- **Before**: Used `/sub-orchestrate-recover` command
- **After**: Routes to appropriate existing command based on type
- **Mapping**:
  - FIX_CASCADE → `/fix-cascade`
  - SPLITTING → `/splitting`
  - Other → `/continue-orchestrating`

## Why This is Better

### 1. Single Source of Truth
Each workflow has ONE command that handles both direct invocation and background spawning:
- `/fix-cascade` - Handles all fix cascade operations
- `/splitting` - Handles all splitting operations

### 2. Reduced Maintenance
- No duplicate logic to keep synchronized
- No confusion about which command to use
- Fewer commands to document and test

### 3. Cleaner Architecture
- Commands are self-contained with their own sub-state machines
- No artificial distinction between "regular" and "sub" orchestration
- Simpler mental model for developers

### 4. Better Command Names
- `/fix-cascade` clearly indicates its purpose
- `/splitting` clearly indicates its purpose
- No generic `/sub-orchestrate-X` naming

## Implementation Details

### Background Spawning Pattern
The master orchestrator spawns commands in background using:
```python
result = Bash(
    command=f"claude -p . --command '/fix-cascade' --params 'fix_id={id} branches={branches}'",
    run_in_background=True,
    description="Spawn fix cascade sub-process"
)
shell_id = result['shell_id']  # Track for monitoring
```

### State Management
Each command manages its own sub-state file:
- **Fix Cascade**: `orchestrator-${fix_id}-state.json`
- **Splitting**: `splitting-${effort}-state.json`

### Signal Output
Commands output standard signals for monitoring:
```bash
echo "SUB_PROCESS_PID: $$"
echo "SUB_PROCESS_STARTED: $(date -Iseconds)"
echo "SUB_PROCESS_TYPE: FIX_CASCADE"
echo "SUB_PROCESS_PROGRESS: 25% - Processing branch 1 of 4"
echo "SUB_PROCESS_COMPLETE: Success"
```

## Migration Guide

### For Existing Code
If you have code referencing sub-orchestrate commands:

1. **Replace command names**:
   - `/sub-orchestrate-fix` → `/fix-cascade`
   - `/sub-orchestrate-split` → `/splitting`

2. **Update parameter passing**:
   - Remove parameter file creation
   - Pass parameters directly in command

3. **Update signal parsing**:
   - `SUB_ORCHESTRATOR_*` → `SUB_PROCESS_*`

### For New Development
Simply use the existing commands:
- Need a fix cascade? Use `/fix-cascade`
- Need splitting? Use `/splitting`
- Spawn in background with `run_in_background=True`

## Conclusion

This simplification removes unnecessary abstraction layers while preserving all required functionality. The system is now easier to understand, maintain, and extend.

**Key Principle**: Don't create new commands unless they add unique value. Existing commands can be spawned in background mode when needed for parallel execution.