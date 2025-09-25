# SOFTWARE FACTORY SUB-STATE MACHINE ARCHITECTURE

## Overview
The Software Factory uses a hierarchical state machine architecture that supports sub-state machines for complex workflows. This allows the main orchestration flow to divert to specialized sub-flows (fix cascades, PR preparation, initialization) while maintaining clean state separation and return capability.

## Architecture Principles

### 1. State Isolation
- **Main State**: Tracks high-level project progress only
- **Sub-States**: Handle specific workflow details in isolation
- **No Pollution**: Sub-state details NEVER contaminate main state

### 2. Clean Transitions
- Entry points clearly defined
- Return states recorded before entry
- Exit always returns to recorded state
- Sub-state completion archived

### 3. Nested Support
- Sub-states can spawn nested sub-states
- Maximum nesting depth: 3 levels
- Each level maintains its own state file

## Available Sub-State Machines

### 1. FIX-CASCADE
- **File**: SOFTWARE-FACTORY-FIX-CASCADE-STATE-MACHINE.md
- **Purpose**: Handle hotfixes, backports, forward-ports
- **Entry**: From ERROR_RECOVERY or MONITORING
- **State File**: orchestrator-[fix-id]-state.json

### 2. PR-READY
- **File**: SOFTWARE-FACTORY-PR-READY-STATE-MACHINE.md
- **Purpose**: Prepare branches for production PRs
- **Entry**: From PROJECT_INTEGRATION or manual trigger
- **State File**: pr-ready-state.json

### 3. INITIALIZATION
- **File**: SOFTWARE-FACTORY-INITIALIZATION-STATE-MACHINE.md
- **Purpose**: Initialize new projects
- **Entry**: From INIT state
- **State File**: initialization-state.json

### 4. SPLITTING
- **File**: SOFTWARE-FACTORY-SPLITTING-STATE-MACHINE.md
- **Purpose**: Split oversized efforts into compliant chunks
- **Entry**: From MEASURE_SIZE, CODE_REVIEW, MONITOR_IMPLEMENTATION
- **State File**: splitting-[effort-name]-state.json

### 5. INTEGRATION
- **File**: SOFTWARE-FACTORY-INTEGRATION-STATE-MACHINE.md
- **Purpose**: Handle complex integration workflows with fix cycles
- **Entry**: From WAVE_COMPLETE, PHASE_COMPLETE, PROJECT_INTEGRATION
- **State File**: integration-[type]-[identifier]-state.json

## State File Structure

### Main State File
```json
{
  "current_state": "MONITORING",
  "sub_state_machine": {
    "active": true,
    "type": "FIX_CASCADE",
    "state_file": "orchestrator-gitea-fix-state.json",
    "current_state": "FIX_CASCADE_BACKPORT_IN_PROGRESS",
    "return_state": "MONITORING",
    "started_at": "2025-01-21T10:00:00Z",
    "trigger_reason": "Critical API fix needed",
    "nested_level": 0,
    "parent_sub_state": null
  },
  "sub_state_history": [
    {
      "type": "FIX_CASCADE",
      "state_file": "orchestrator-auth-fix-state.json",
      "started_at": "2025-01-20T09:00:00Z",
      "completed_at": "2025-01-20T10:30:00Z",
      "result": "SUCCESS",
      "archived_to": "archived-fixes/2025/01/auth-fix-20250120.json"
    }
  ]
}
```

### Sub-State File
```json
{
  "sub_state_type": "FIX_CASCADE",
  "current_state": "FIX_CASCADE_VALIDATION",
  "parent_state_machine": {
    "state_file": "orchestrator-state.json",
    "return_state": "MONITORING",
    "nested_level": 1
  },
  "sub_state_details": {
    // Specific to sub-state type
  }
}
```

## Entry Protocol

### 1. Decision to Enter Sub-State
```python
def should_enter_sub_state(context):
    if context.fix_cascade_required:
        return ("FIX_CASCADE", "orchestrator-fix-state.json")
    elif context.pr_preparation_needed:
        return ("PR_READY", "pr-ready-state.json")
    elif context.initialization_required:
        return ("INITIALIZATION", "initialization-state.json")
    elif context.size_violation_detected:
        return ("SPLITTING", f"splitting-{context.effort_name}-state.json")
    return (None, None)
```

### 2. Entry Execution
```bash
enter_sub_state_machine() {
    local SUB_TYPE="$1"
    local STATE_FILE="$2"
    local RETURN_STATE=$(jq -r '.current_state' orchestrator-state.json)

    # Create sub-state file
    create_sub_state_file "$STATE_FILE" "$SUB_TYPE"

    # Update main state
    jq --arg type "$SUB_TYPE" \
       --arg file "$STATE_FILE" \
       --arg return "$RETURN_STATE" \
       '.sub_state_machine = {
          "active": true,
          "type": $type,
          "state_file": $file,
          "return_state": $return,
          "started_at": now
       }' orchestrator-state.json > tmp.json && \
       mv tmp.json orchestrator-state.json

    echo "✅ Entered $SUB_TYPE sub-state machine"
}
```

## Exit Protocol

### 1. Completion Handling
```bash
complete_sub_state_machine() {
    local RESULT="$1"  # SUCCESS, FAILURE, ABORTED
    local SUB_STATE_FILE=$(jq -r '.sub_state_machine.state_file' orchestrator-state.json)

    # Archive sub-state file
    archive_sub_state_file "$SUB_STATE_FILE" "$RESULT"

    # Return to main state
    RETURN_STATE=$(jq -r '.sub_state_machine.return_state' orchestrator-state.json)

    # Update main state
    jq --arg state "$RETURN_STATE" \
       --arg result "$RESULT" \
       '.current_state = $state |
        .sub_state_machine.active = false |
        .sub_state_history += [{
          "type": .sub_state_machine.type,
          "state_file": .sub_state_machine.state_file,
          "completed_at": now,
          "result": $result
        }]' orchestrator-state.json > tmp.json && \
        mv tmp.json orchestrator-state.json

    echo "✅ Returned to main state: $RETURN_STATE"
}
```

### 2. Archival Process
```bash
archive_sub_state_file() {
    local STATE_FILE="$1"
    local RESULT="$2"
    local TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    local TYPE=$(jq -r '.sub_state_type' "$STATE_FILE")

    # Create archive directory
    mkdir -p "archived-states/$(date +%Y)/$(date +%m)"

    # Archive with result and timestamp
    ARCHIVE_NAME="${TYPE}-${RESULT}-${TIMESTAMP}.json"
    mv "$STATE_FILE" "archived-states/$(date +%Y)/$(date +%m)/$ARCHIVE_NAME"

    # Commit archival
    git add archived-states/
    git commit -m "archive: $TYPE completed with $RESULT at $TIMESTAMP"
    git push
}
```

## Nested Sub-States

### Support for Nested Execution
```json
{
  "sub_state_machine": {
    "active": true,
    "type": "FIX_CASCADE",
    "state_file": "orchestrator-main-fix-state.json",
    "nested_level": 1,
    "parent_sub_state": null,
    "nested_sub_state": {
      "active": true,
      "type": "FIX_CASCADE",
      "state_file": "orchestrator-nested-fix-state.json",
      "nested_level": 2,
      "parent_sub_state": "orchestrator-main-fix-state.json"
    }
  }
}
```

### Maximum Nesting Rules
- Level 0: Main state machine
- Level 1: First sub-state machine
- Level 2: Nested sub-state machine
- Level 3: Maximum depth (no further nesting)

## Command Routing

### /continue-orchestrating
```bash
# Check for active sub-state and route appropriately
if [[ $(jq -r '.sub_state_machine.active' orchestrator-state.json) == "true" ]]; then
    SUB_TYPE=$(jq -r '.sub_state_machine.type' orchestrator-state.json)
    case "$SUB_TYPE" in
        FIX_CASCADE)
            exec_fix_cascade_continuation
            ;;
        PR_READY)
            exec_pr_ready_continuation
            ;;
        INITIALIZATION)
            exec_initialization_continuation
            ;;
        SPLITTING)
            exec_splitting_continuation
            ;;
        INTEGRATION)
            exec_integration_continuation
            ;;
    esac
else
    exec_main_orchestration_continuation
fi
```

### Sub-State Specific Commands
- `/fix-cascade` - Enter or continue fix cascade
- `/pr-ready` - Enter or continue PR preparation
- `/initialize` - Start project initialization
- `/splitting` - Enter or continue effort splitting
- `/integrate-wave` - Start wave integration
- `/integrate-phase` - Start phase integration
- `/integrate-project` - Start project integration

## State Machine Discovery

### Finding Active State
```bash
get_active_state_machine() {
    # Check main state for sub-state
    if [[ $(jq -r '.sub_state_machine.active' orchestrator-state.json) == "true" ]]; then
        STATE_FILE=$(jq -r '.sub_state_machine.state_file' orchestrator-state.json)
        STATE_TYPE=$(jq -r '.sub_state_machine.type' orchestrator-state.json)
        CURRENT_STATE=$(jq -r '.current_state' "$STATE_FILE")
        echo "Active: $STATE_TYPE in state $CURRENT_STATE"
    else
        CURRENT_STATE=$(jq -r '.current_state' orchestrator-state.json)
        echo "Active: Main orchestrator in state $CURRENT_STATE"
    fi
}
```

### State File Locations
```
/
├── orchestrator-state.json                 # Main state
├── orchestrator-*-state.json              # Fix cascade states
├── pr-ready-state.json                    # PR preparation state
├── initialization-state.json              # Project init state
├── integration-*-*-state.json             # Integration states (type-identifier)
├── splitting-*-state.json                 # Splitting states (effort-name)
└── archived-states/                       # Archived sub-states
    └── YYYY/MM/
        └── TYPE-RESULT-TIMESTAMP.json
```

## Rules Integration

### R375 - Fix State File Management
- Enforces dual state file pattern
- Mandates archival on completion
- Prevents state pollution

### R376 - Fix Cascade Quality Gates
- Applied within fix cascade sub-state
- Gates tracked in sub-state file
- No progression without gate pass

### R206 - State Machine Validation
- Validates transitions in correct state file
- Checks sub-state machine transitions
- Ensures return state validity

## Benefits

### 1. Clean Separation
- Main flow uncluttered
- Sub-flows isolated
- Easy to understand

### 2. Reusability
- Sub-state machines can be reused
- Common patterns abstracted
- Consistent behavior

### 3. Maintainability
- Changes isolated to sub-state
- Main flow protected
- Clear boundaries

### 4. Debuggability
- State files show exact location
- History tracked
- Clear audit trail

### 5. Scalability
- New sub-states easily added
- No impact on main flow
- Supports complex workflows

## Implementation Checklist

When implementing a new sub-state machine:

- [ ] Create state machine definition (SOFTWARE-FACTORY-[NAME]-STATE-MACHINE.md)
- [ ] Create state directories in agent-states/[name]/
- [ ] Create state rules for each state
- [ ] Add entry logic to parent states
- [ ] Implement state file structure
- [ ] Add to command routing
- [ ] Update this documentation
- [ ] Test transitions in and out
- [ ] Verify archival works
- [ ] Document in relevant rules

## Common Patterns

### Pattern 1: Error Recovery to Fix
```
Main: ERROR_RECOVERY
  → Sub: FIX_CASCADE
    → Complete
  → Main: Return to ERROR_RECOVERY
  → Main: Continue to next state
```

### Pattern 2: Integration to PR
```
Main: PROJECT_INTEGRATION
  → Sub: PR_READY
    → Complete
  → Main: Return to PROJECT_INTEGRATION
  → Main: SUCCESS
```

### Pattern 3: Nested Fix Discovery
```
Main: MONITORING
  → Sub: FIX_CASCADE (level 1)
    → Nested: FIX_CASCADE (level 2)
      → Complete
    → Sub: Continue
    → Complete
  → Main: MONITORING
```

### Pattern 4: Size Violation Splitting
```
Main: MONITOR_IMPLEMENTATION
  → Sub: SPLITTING
    → Split Analysis
    → Create Split Plan
    → Execute Splits Sequentially
    → Complete
  → Main: Return to MONITOR_IMPLEMENTATION
  → Main: Continue with split branches
```

### Pattern 5: Integration-Fix Cycles
```
Main: WAVE_COMPLETE
  → Sub: INTEGRATION
    → Merge → Build fails
    → Exit with FIX_REQUIRED
  → Main: Trigger FIX_CASCADE
    → Fix issues in source branches
    → Complete fixes
  → Main: Re-enter INTEGRATION
  → Sub: INTEGRATION (attempt 2)
    → Delete stale branch
    → Re-merge with fixes
    → Build OK → Tests fail
    → Exit with FIX_REQUIRED
  → Main: Trigger FIX_CASCADE again
  → Sub: INTEGRATION (attempt 3)
    → All validations pass
    → Complete
  → Main: Return to WAVE_COMPLETE
  → Main: Continue to next wave
```

## Error Handling

### Sub-State Failure
- Archive partial state
- Return to parent with FAILURE
- Parent decides next action

### Lost State Recovery
- Check sub_state_machine.active
- Locate state file
- Resume or abort

### Corruption Recovery
- Use archived states
- Rebuild from history
- Manual intervention if needed