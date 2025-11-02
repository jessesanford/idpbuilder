# R360 - Just-In-Time Infrastructure Execution

## Rule ID: R360
## Criticality: HIGH
## State Scope: CREATE_NEXT_INFRASTRUCTURE, WAVE_START, EFFORT_COMPLETE
## Agents: orchestrator
## Related: R504 (Pre-Infrastructure Planning)

## Overview

Infrastructure for efforts and splits must be EXECUTED "just in time" - only when the effort is ready to be worked on. This ensures efforts can build on completed dependencies within the same wave.

**IMPORTANT**: R504 handles ALL naming/pathing decisions upfront. R360 only controls WHEN to execute pre-planned infrastructure.

**RELATED**: R603 (Sequential Effort Infrastructure Timing) handles dependency checking and determines WHICH efforts are ready based on R213 metadata. R360 provides the general principle, R603 provides the implementation logic.

## Requirements

### 1. Infrastructure Creation Timing

**MANDATORY**: Create infrastructure ONLY when:
- All dependencies for the effort are complete
- The effort is next in the execution sequence
- The effort is about to be worked on

**PROHIBITED**: Do NOT create infrastructure:
- For all efforts at wave start
- Before dependencies are complete
- For efforts that may not be executed

### 2. Dependency Resolution

When creating infrastructure for an effort:

```bash
# Determine base branch based on dependencies
if effort.depends_on exists:
    base_branch = completed_effort(effort.depends_on).branch
else:
    base_branch = previous_wave_integration_branch
```

### 3. Effort Dependency Tracking

Track in orchestrator-state-v3.json:
```json
{
  "infrastructure_strategy": "just_in_time",
  "effort_dependencies": {
    "effort1": {
      "depends_on": null,
      "base_branch": "phase1/wave0/integration",
      "infrastructure_created": false,
      "status": "pending"
    },
    "effort2": {
      "depends_on": "effort1",
      "base_branch": null,  // Determined when created
      "infrastructure_created": false,
      "status": "waiting_for_dependency"
    }
  }
}
```

### 4. State Machine Flow

```
WAVE_START
    ↓
CREATE_NEXT_INFRASTRUCTURE (for first ready effort)
    ↓
SPAWN_SW_ENGINEERS (implement effort)
    ↓
EFFORT_COMPLETE
    ↓
CREATE_NEXT_INFRASTRUCTURE (for next ready effort)
    ↓
[Loop until all efforts complete]
```

### 5. Infrastructure Creation Process

For each effort needing infrastructure:

1. **Check Dependencies**:
   - All dependencies completed? ✓
   - Infrastructure not yet created? ✓

2. **Determine Base**:
   - If depends_on: use that effort's branch
   - Else: use previous wave's integration

3. **Create Infrastructure**:
   ```bash
   # Clone repository
   git clone "$TARGET_REPO" "$EFFORT_DIR"
   cd "$EFFORT_DIR"

   # Create branch from correct base
   git checkout -b "$EFFORT_BRANCH" "$BASE_BRANCH"
   git push -u origin "$EFFORT_BRANCH"

   # Lock git config (R312)
   chmod 444 .git/config
   ```

4. **Update Tracking**:
   - Set infrastructure_created = true
   - Record actual base_branch used
   - Update status to "ready"

### 6. Benefits

- **Proper Dependencies**: Later efforts get earlier efforts' commits
- **Resource Efficiency**: Only create what's needed when needed
- **Flexibility**: Can adjust execution based on results
- **Consistency**: Same approach for efforts and splits

### 7. Migration Path

**Phase 1**: Parallel support (both strategies available)
**Phase 2**: Default to JIT for new waves
**Phase 3**: Fully deprecate batch infrastructure creation

## Implementation

### For Orchestrator

```bash
# In CREATE_NEXT_INFRASTRUCTURE state
next_effort=$(jq -r '.effort_dependencies | to_entries[] |
  select(.value.infrastructure_created == false and
         (.value.depends_on == null or
          .effort_dependencies[.value.depends_on].status == "complete")) |
  .key' orchestrator-state-v3.json | head -1)

if [ -n "$next_effort" ]; then
    create_effort_infrastructure "$next_effort"
fi
```

### Validation

```bash
# Ensure no effort has infrastructure before dependencies complete
for effort in $(jq -r '.effort_dependencies | keys[]' orchestrator-state-v3.json); do
    depends_on=$(jq -r ".effort_dependencies.$effort.depends_on" orchestrator-state-v3.json)
    if [ "$depends_on" != "null" ]; then
        dep_status=$(jq -r ".effort_dependencies.$depends_on.status" orchestrator-state-v3.json)
        our_infra=$(jq -r ".effort_dependencies.$effort.infrastructure_created" orchestrator-state-v3.json)

        if [ "$dep_status" != "complete" ] && [ "$our_infra" = "true" ]; then
            echo "ERROR: $effort has infrastructure but dependency $depends_on not complete!"
            exit 360
        fi
    fi
done
```

## Enforcement

- Exit code 360 for violations
- Grading penalty: -30% for creating infrastructure before dependencies complete
- Orchestrator must validate dependency chain before infrastructure creation

## Related Rules

- R603: Sequential effort infrastructure timing (dependency checking implementation)
- R213: Effort metadata requirements (source of dependency information)
- R504: Pre-infrastructure planning protocol (naming/pathing decisions)
- R196: Base branch selection
- R308: Incremental branching strategy
- R312: Git config immutability
- R176: Workspace isolation

## Examples

### Correct: JIT with Dependencies
```
Wave 2 has: effort1, effort2 (depends on effort1), effort3 (independent)

1. CREATE_NEXT_INFRASTRUCTURE → effort1 (base: wave1/integration)
2. Implement effort1
3. CREATE_NEXT_INFRASTRUCTURE → effort3 (base: wave1/integration) [parallel]
4. CREATE_NEXT_INFRASTRUCTURE → effort2 (base: effort1) [after effort1 done]
5. Implement effort2 and effort3
```

### Wrong: Batch Creation
```
Wave 2 has: effort1, effort2 (depends on effort1)

1. CREATE_NEXT_INFRASTRUCTURE → Create both efforts based on wave1/integration
2. effort2 doesn't have effort1's commits! ❌
```