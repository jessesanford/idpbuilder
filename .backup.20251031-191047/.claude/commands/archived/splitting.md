# /splitting Command
## Enter or Continue SPLITTING Sub-State Machine

### Purpose
Handle effort splitting when size limits are exceeded (>800 lines). This command manages the SPLITTING sub-state machine for breaking oversized implementations into compliant chunks.

### Usage
```
/splitting [effort-name] [size]
```

### Parameters
- `effort-name` (optional): The effort to split (e.g., "effort2-controller")
- `size` (optional): The measured size in lines

If no parameters provided, continues existing split operation.

### Entry Conditions

The SPLITTING sub-state can be entered from:
1. **MEASURE_SIZE** (SW Engineer) - Size violation detected
2. **CODE_REVIEW** (Code Reviewer) - Split planning complete
3. **MONITORING_SWE_PROGRESS** (Orchestrator) - Monitoring detects oversized
4. **MONITORING_EFFORT_REVIEWS** (Orchestrator) - Review requires split

### Command Behavior

#### New Split Operation
```bash
/splitting effort2-controller 1250
```

This will:
1. Check if already in SPLITTING sub-state
2. Create `splitting-effort2-controller-state.json`
3. Initialize split context
4. Enter SPLIT_INIT state
5. Begin split workflow

#### Continue Existing Split
```bash
/splitting
```

This will:
1. Check for active SPLITTING sub-state
2. Load current split state
3. Continue from last checkpoint
4. Route to appropriate agent/state

### State Machine Flow

```
Entry → INIT
      → CREATE_SPLIT_PLAN (Code Reviewer)
      → CREATE_SPLIT_INVENTORY (Code Reviewer)
      → CREATE_NEXT_INFRASTRUCTURE (Orchestrator)
      → SPLIT_IMPLEMENTATION (SW Engineer)
      → MEASURE_SIZE
      → SPLIT_REVIEW (Code Reviewer)
      → [Loop for each split]
      → SPLIT_VALIDATION
      → SPLIT_COMPLETE
      → Return to Main State
```

### Sub-State File Structure

```json
{
  "sub_state_type": "SPLITTING",
  "current_state": "SPLIT_IMPLEMENTATION",
  "effort_id": "effort2-controller",
  "original_effort": {
    "size": 1250,
    "branch": "phase1-wave2-effort2"
  },
  "split_plan": {
    "total_splits": 3,
    "completed_splits": 1,
    "current_split": "split-002"
  }
}
```

### Integration with Main State

When SPLITTING is active, the main orchestrator-state-v3.json shows:
```json
{
  "current_state": "MONITORING_SWE_PROGRESS",
  "sub_state_machine": {
    "active": true,
    "type": "SPLITTING",
    "state_file": "splitting-effort2-controller-state.json",
    "return_state": "MONITORING_SWE_PROGRESS"
  }
}
```

### Command Implementation

```python
def handle_splitting_command(effort_name=None, size=None):
    """Handle /splitting command"""

    # Check if already in SPLITTING
    main_state = load_json("orchestrator-state-v3.json")

    if main_state.get("sub_state_machine", {}).get("type") == "SPLITTING":
        # Continue existing split
        split_state_file = main_state["sub_state_machine"]["state_file"]
        split_state = load_json(split_state_file)

        print(f"Continuing SPLITTING operation")
        print(f"Current state: {split_state['current_state']}")
        print(f"Effort: {split_state['effort_id']}")
        print(f"Progress: {split_state['split_plan']['completed_splits']}/{split_state['split_plan']['total_splits']}")

        # Route to appropriate continuation
        route_splitting_continuation(split_state)

    elif effort_name and size:
        # Start new split operation
        if size <= 800:
            print(f"⚠️ Size {size} is within limit (≤800)")
            print("Split not required")
            return

        print(f"🔄 Starting SPLITTING operation for {effort_name}")
        print(f"📏 Size: {size} lines (limit: 800)")

        # Enter SPLITTING sub-state
        enter_splitting_sub_state(effort_name, size)

    else:
        print("❌ No active SPLITTING operation")
        print("Usage: /splitting [effort-name] [size]")
```

### Quality Gates

The SPLITTING sub-state enforces quality gates at:
1. **Split Size**: Each split MUST be ≤800 lines
2. **Code Review**: Each split requires full review
3. **Integration**: All splits must work together

### Exit Conditions

SPLITTING completes when:
- ✅ All splits implemented and reviewed
- ✅ Each split ≤800 lines
- ✅ Integration validation passed
- ✅ Original marked SPLIT_DEPRECATED (R296)

### Error Handling

If SPLITTING fails:
1. Partial state is archived
2. Returns to ERROR_RECOVERY in main state
3. Manual intervention may be required

### Related Commands

- `/continue-orchestrating` - Continues any active sub-state
- `/check-split-status` - Shows current split progress
- `/abort-splitting` - Emergency abort (use carefully)

### Rules Integration

- **R296**: Original effort marked SPLIT_DEPRECATED
- **R304**: MUST use line-counter.sh for measurements
- **R204**: Infrastructure created just-in-time
- **R308**: Sequential split dependencies

### Example Workflow

```bash
# 1. Size violation detected
SW Engineer: "Detected 1250 lines in effort2-controller"

# 2. Enter SPLITTING
/splitting effort2-controller 1250

# 3. Code Reviewer analyzes and plans
# 4. Orchestrator creates infrastructure
# 5. SW Engineer implements split-001
# 6. Code review of split-001
# 7. Continue with split-002...
# 8. Complete all splits
# 9. Return to main flow

# Main state continues with split branches instead of original
```

### Monitoring

Check split progress:
```bash
jq '.current_state, .split_plan.completed_splits, .split_plan.total_splits' \
    splitting-effort2-controller-state.json
```

### Troubleshooting

**Issue**: Split still too large after planning
**Solution**: SPLIT_REPLANNING state handles re-splitting

**Issue**: Lost split context
**Solution**: Check archived-states/ for recovery

**Issue**: Conflicts between splits
**Solution**: Sequential implementation prevents this

---

**Command Version**: 1.0.0
**Sub-State Machine**: SPLITTING
**Created**: 2025-01-21