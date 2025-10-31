# Rule R017: Checkpoint Creation Protocol

## Rule Statement
Agents MUST create checkpoints at critical junctures to enable recovery from failures. Checkpoints include saving TODOs, committing work, updating state files, and creating recovery markers before high-risk operations.

## Criticality Level
**WARNING** - Missing checkpoints can lead to significant work loss

## Enforcement Mechanism
- **Technical**: Automatic checkpoint creation at defined triggers
- **Behavioral**: Save state before transitions and spawns
- **Grading**: -20% for missing checkpoints, -50% for work loss due to missing checkpoints

## Core Principle

```
Checkpoint = Save TODOs → Commit Work → Update State → Create Marker
ALWAYS checkpoint before state transitions
ALWAYS checkpoint before spawning agents  
ALWAYS checkpoint before risky operations
```

## Detailed Requirements

### Mandatory Checkpoint Triggers

1. **Before State Transitions**
   ```bash
   # Save current work state
   save_todos "PRE_TRANSITION_${current_state}_TO_${next_state}"
   git add -A && git commit -m "checkpoint: before transition to ${next_state}"
   git push
   ```

2. **Before Agent Spawning**
   ```bash
   # Create spawn checkpoint
   save_todos "PRE_SPAWN_${agent_type}"
   update_state_file "spawning_agent" "${agent_type}"
   git commit -m "checkpoint: before spawning ${agent_type}"
   git push
   ```

3. **Before High-Risk Operations**
   - Integration merges
   - Large refactoring
   - Deletion operations
   - Configuration changes

### Checkpoint Components

1. **TODO State**: Current task list with statuses
2. **Git State**: All work committed and pushed
3. **State File**: Current state and metadata
4. **Recovery Marker**: Timestamp and operation identifier

### Recovery from Checkpoint

```bash
# Find latest checkpoint
latest_checkpoint=$(ls todos/*.todo | sort -r | head -1)

# Restore TODO state
load_todos "${latest_checkpoint}"

# Verify git state
git status
git log --oneline -5

# Resume from checkpoint
echo "Resumed from checkpoint: ${latest_checkpoint}"
```

## Relationship to Other Rules
- **R287**: TODO persistence comprehensive
- **R173**: State preservation
- **R174**: Recovery detection
- **R288**: State file update and commit

## Implementation Notes
- Checkpoints must be atomic operations
- Recovery must be possible within 2 minutes
- Checkpoint files must include timestamps
- Maximum 30 seconds between trigger and checkpoint creation