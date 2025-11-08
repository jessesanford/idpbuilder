# R603 - Sequential Effort Infrastructure Timing

## Rule ID: R603
## Criticality: HIGH
## State Scope: CREATE_NEXT_INFRASTRUCTURE
## Agents: orchestrator
## Related: R360 (Just-In-Time Execution), R213 (Effort Metadata), R501 (Cascade Branching)

## Overview

When a wave contains SEQUENTIAL efforts (efforts with dependencies on each other), infrastructure must be created incrementally "just-in-time" as each effort completes. This ensures later efforts build on earlier efforts' commits, maintaining the cascade pattern.

**CRITICAL**: This rule addresses the timing of WHEN to create infrastructure for sequential efforts. R360 handles the general "just-in-time" principle, while R603 specifically handles dependency checking and sequential vs parallel determination.

## Problem Statement

### What Happens Without R603

**Scenario**: Wave 2.2 has 2 efforts:
- Effort 2.2.1: `"parallelizable": false`, `"depends_on": ["integration:wave2.1"]`
- Effort 2.2.2: `"parallelizable": false`, `"depends_on": ["effort:2.2.1"]`

**WITHOUT R603** (broken behavior):
```
1. CREATE_NEXT_INFRASTRUCTURE runs
2. Finds: Both 2.2.1 and 2.2.2 have created=false
3. Creates: 2.2.1 infrastructure (base: wave2.1-integration) ✅
4. Creates: 2.2.2 infrastructure (base: wave2.1-integration) ❌ WRONG!
5. Result: 2.2.2 missing config.go from 2.2.1
6. R509 validation fails: base branch doesn't have expected files
7. ERROR_RECOVERY triggered
```

**WITH R603** (correct behavior):
```
1. CREATE_NEXT_INFRASTRUCTURE runs
2. Finds: Both 2.2.1 and 2.2.2 have created=false
3. Checks dependencies:
   - 2.2.1 depends on wave2.1-integration (complete) ✓
   - 2.2.2 depends on effort:2.2.1 (NOT complete) ✗
4. Creates: ONLY 2.2.1 infrastructure
5. Spawns: SW Engineer for 2.2.1
6. 2.2.1 completes and is approved
7. CREATE_NEXT_INFRASTRUCTURE runs again
8. Checks dependencies:
   - 2.2.2 depends on effort:2.2.1 (NOW complete) ✓
9. Creates: 2.2.2 infrastructure (base: effort-2.2.1) ✅
10. Result: 2.2.2 HAS config.go from 2.2.1
11. R509 validation passes
```

## Requirements

### 1. R213 Metadata as Source of Truth

Infrastructure creation timing MUST read from wave implementation plan's R213 metadata:

```json
{
  "effort_id": "2.2.2",
  "depends_on": ["effort:2.2.1"],
  "parallelizable": false
}
```

**Fields to check**:
- `depends_on`: Array of dependencies (efforts or integrations)
- `parallelizable`: Boolean indicating if can run in parallel
- `can_parallelize`: Alternate field name (check both)

### 2. Dependency Checking Algorithm

Before creating infrastructure for any effort, verify ALL dependencies are complete:

```bash
# For each potential effort:
check_dependencies_satisfied() {
    local effort_id="$1"

    # Extract R213 metadata from wave plan
    WAVE_PLAN=$(get_wave_plan_path)
    depends_on=$(extract_r213_field "$effort_id" "depends_on" "$WAVE_PLAN")

    # Check each dependency
    for dep in $depends_on; do
        dep_type=$(echo "$dep" | cut -d: -f1)  # "effort" or "integration"
        dep_name=$(echo "$dep" | cut -d: -f2)

        if [ "$dep_type" = "effort" ]; then
            # Check if dependent effort is complete and approved
            dep_status=$(jq -r ".effort_status.\"$dep_name\".status" orchestrator-state-v3.json)
            if [ "$dep_status" != "approved" ]; then
                echo "Dependency $dep not satisfied (status: $dep_status)"
                return 1
            fi
        elif [ "$dep_type" = "integration" ]; then
            # Check if integration branch exists
            if ! git ls-remote --heads origin "$dep_name" | grep -q "$dep_name"; then
                echo "Dependency $dep not satisfied (branch missing)"
                return 1
            fi
        fi
    done

    echo "All dependencies satisfied for $effort_id"
    return 0
}
```

### 3. Parallel vs Sequential Detection

Determine if efforts can be created in batch (parallel) or must be incremental (sequential):

```bash
# Check if ANY effort in wave has parallelizable: false
has_sequential_efforts() {
    local wave_plan="$1"

    # Extract all parallelizable values from R213 metadata
    grep -A 20 '"effort_id"' "$wave_plan" | grep '"parallelizable"' | grep -q 'false'
    return $?
}

# Determine creation strategy
if has_sequential_efforts "$WAVE_PLAN"; then
    echo "Wave has sequential efforts - incremental infrastructure creation"
    CREATION_MODE="incremental"
else
    echo "All efforts parallelizable - batch infrastructure creation"
    CREATION_MODE="batch"
fi
```

### 4. Infrastructure Creation Strategy

**Batch Mode** (all efforts parallelizable):
```bash
# Create ALL effort infrastructure upfront
for effort_id in $(get_all_wave_efforts); do
    if check_dependencies_satisfied "$effort_id"; then
        create_effort_infrastructure "$effort_id"
    fi
done

# Then spawn ALL SW Engineers in parallel (R151 timing)
```

**Incremental Mode** (any effort sequential):
```bash
# Create infrastructure for FIRST ready effort only
next_effort_id=$(find_next_ready_effort)

if [ -n "$next_effort_id" ]; then
    create_effort_infrastructure "$next_effort_id"
    # After this effort completes, CREATE_NEXT_INFRASTRUCTURE runs again
    # to create infrastructure for next effort
fi
```

### 5. State Machine Flow for Sequential Efforts

```
SETUP_WAVE_INFRASTRUCTURE
    ↓
CREATE_NEXT_INFRASTRUCTURE
    ↓ (creates Effort 1 only)
SPAWN_SW_ENGINEERS (single SWE for Effort 1)
    ↓
MONITORING_SWE_PROGRESS
    ↓
SPAWN_CODE_REVIEWERS_EFFORT_REVIEW
    ↓
MONITORING_EFFORT_REVIEWS
    ↓ (Effort 1 approved)
CREATE_NEXT_INFRASTRUCTURE
    ↓ (creates Effort 2 based on Effort 1)
SPAWN_SW_ENGINEERS (single SWE for Effort 2)
    ↓
...
    ↓
COMPLETE_WAVE
```

## Implementation

### For CREATE_NEXT_INFRASTRUCTURE State

**Location**: `agent-states/software-factory/orchestrator/CREATE_NEXT_INFRASTRUCTURE/rules.md`

**Section to modify**: "Determine What Needs Infrastructure" (around line 318-342)

**NEW LOGIC**:

```bash
echo "🔧 DETERMINING NEXT INFRASTRUCTURE TO CREATE (R603)..."

# Get current phase and wave
CURRENT_PHASE=$(yq '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
CURRENT_WAVE=$(yq '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)

# Get wave implementation plan path
WAVE_PLAN=$(jq -r ".planning_files.phases.phase${CURRENT_PHASE}.waves.wave${CURRENT_WAVE}.implementation_plan" \
           orchestrator-state-v3.json)

if [ ! -f "$WAVE_PLAN" ]; then
    echo "🚨 ERROR: Wave implementation plan not found: $WAVE_PLAN"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Missing wave implementation plan (R603 violation)"
    exit 603
fi

# Find ALL uncreated efforts in current wave
all_uncreated=$(yq '.pre_planned_infrastructure.efforts | to_entries[] |
  select(.value.phase == "phase'$CURRENT_PHASE'" and
         .value.wave == "wave'$CURRENT_WAVE'" and
         .value.created == false) |
  .key' orchestrator-state-v3.json)

if [ -z "$all_uncreated" ]; then
    echo "✅ All infrastructure for Phase $CURRENT_PHASE Wave $CURRENT_WAVE created"
    PROPOSED_NEXT_STATE="WAVE_COMPLETE"
    TRANSITION_REASON="No more infrastructure to create in wave"
    # Proceed to wave completion
    exit 0
fi

echo "📋 Uncreated efforts in wave: $(echo "$all_uncreated" | tr '\n' ' ')"

# R603: Check dependencies for each uncreated effort
next_effort_id=""
for effort_id in $all_uncreated; do
    echo "🔍 Checking dependencies for $effort_id..."

    # Extract R213 metadata from wave plan
    depends_on=$(grep -A 30 "\"effort_id\": \"${effort_id##*effort-}\"" "$WAVE_PLAN" | \
                 grep '"depends_on"' | \
                 sed 's/.*\[//' | sed 's/\].*//' | tr -d '"' | tr ',' ' ')

    echo "  Dependencies: $depends_on"

    # Check if all dependencies are satisfied
    dependencies_satisfied=true
    for dep in $depends_on; do
        dep_type=$(echo "$dep" | cut -d: -f1)
        dep_name=$(echo "$dep" | cut -d: -f2)

        if [ "$dep_type" = "effort" ]; then
            # Check if dependent effort is approved
            dep_status=$(jq -r ".effort_status.\"$dep_name\".status // \"not_found\"" \
                        orchestrator-state-v3.json)

            if [ "$dep_status" != "approved" ]; then
                echo "  ⏸️  Dependency $dep not satisfied (status: $dep_status)"
                dependencies_satisfied=false
                break
            fi
        elif [ "$dep_type" = "integration" ]; then
            # Check if integration branch exists
            TARGET_REPO=$(yq '.pre_planned_infrastructure.target_repo_url' orchestrator-state-v3.json)
            if ! git ls-remote --heads "$TARGET_REPO" | grep -q "refs/heads/$dep_name"; then
                echo "  ⏸️  Dependency $dep not satisfied (branch missing)"
                dependencies_satisfied=false
                break
            fi
        fi
    done

    if [ "$dependencies_satisfied" = "true" ]; then
        echo "  ✅ All dependencies satisfied"
        next_effort_id="$effort_id"
        break
    else
        echo "  ⏸️  Waiting for dependencies to complete"
    fi
done

if [ -z "$next_effort_id" ]; then
    echo "🚨 ERROR: No efforts ready for infrastructure creation"
    echo "All uncreated efforts have unsatisfied dependencies"
    echo "This indicates a dependency deadlock or planning error"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="No ready efforts (R603 dependency deadlock)"
    exit 603
fi

echo "🎯 Creating infrastructure for: $next_effort_id"
infrastructure_type="effort"
infrastructure_target="$next_effort_id"

# Check parallelizability for future optimization
parallelizable=$(grep -A 30 "\"effort_id\": \"${next_effort_id##*effort-}\"" "$WAVE_PLAN" | \
                grep '"parallelizable"' | grep -o 'true\|false' || echo "false")

if [ "$parallelizable" = "true" ]; then
    echo "📊 Note: Effort is parallelizable - could batch create with other parallel efforts"
    echo "   (Current implementation: incremental creation for safety)"
else
    echo "📊 Note: Effort is sequential - incremental creation required"
fi
```

## Validation

### Post-Implementation Checks

After implementing R603, verify:

```bash
# Test 1: Sequential efforts create incrementally
# Wave with: effort1 (parallelizable: false), effort2 (depends_on: effort1, parallelizable: false)
# Expected: Only effort1 infrastructure created initially

# Test 2: Parallel efforts can batch create
# Wave with: effort1 (parallelizable: true), effort2 (parallelizable: true, no dependencies)
# Expected: Both infrastructures can be created (though implementation may still do incremental)

# Test 3: Dependency checking
# Effort2 depends on effort1, but effort1 not approved
# Expected: CREATE_NEXT_INFRASTRUCTURE does NOT create effort2 infrastructure

# Test 4: Integration dependencies
# Effort1 depends on integration:phase1-wave2-integration
# Expected: Checks if integration branch exists before creating infrastructure
```

### Validation Script

```bash
#!/bin/bash
# Validate R603 compliance

validate_r603_sequential_infrastructure() {
    local wave_plan="$1"

    # Extract all efforts with parallelizable: false
    sequential_efforts=$(grep -B 5 '"parallelizable": false' "$wave_plan" | \
                        grep '"effort_id"' | cut -d'"' -f4)

    for effort in $sequential_efforts; do
        # Check if infrastructure was created before dependencies satisfied
        created_at=$(jq -r ".pre_planned_infrastructure.efforts.\"$effort\".created_at" \
                    orchestrator-state-v3.json)

        if [ "$created_at" != "null" ]; then
            # Infrastructure was created - verify dependencies were satisfied at that time
            depends_on=$(extract_r213_field "$effort" "depends_on" "$wave_plan")

            for dep in $depends_on; do
                # Check if dependency was complete before this effort's infrastructure created
                # (Implementation details depend on tracking)
                echo "Validating $effort was created after $dep completed"
            done
        fi
    done
}
```

## Enforcement

- Exit code 603 for violations
- Grading penalty: -50% for creating infrastructure before dependencies complete
- Orchestrator must validate R213 metadata before infrastructure creation
- Pre-commit hooks should validate dependency order in state file

## Benefits

### With R603 Compliance

- ✅ Sequential efforts always have access to previous effort's commits
- ✅ Base branch cascade maintained automatically
- ✅ R509 validation passes (expected files present)
- ✅ No manual rebasing or cherry-picking needed
- ✅ Clean merge sequence preserved

### Without R603 (Current Problem)

- ❌ Later efforts missing earlier efforts' commits
- ❌ R509 validation failures
- ❌ Cascade broken
- ❌ Implementation blocked
- ❌ Manual intervention required

## Related Rules

- **R360**: Just-In-Time Infrastructure Execution (general principle)
- **R213**: Effort Metadata Requirements (source of dependency information)
- **R501**: Progressive Trunk-Based Development (cascade pattern)
- **R509**: Mandatory Base Branch Validation (detects cascade violations)
- **R504**: Pre-Infrastructure Planning Protocol (pre-planning dependencies)
- **R514**: Infrastructure Creation Protocol (HOW to create infrastructure)

## Migration Path

**Phase 1**: Update CREATE_NEXT_INFRASTRUCTURE state rules with R603 logic
**Phase 2**: Add validation scripts to detect R603 violations
**Phase 3**: Add pre-commit hooks to prevent R603 violations
**Phase 4**: Document best practices for wave planning with sequential efforts

## Examples

### Example 1: Simple Sequential (Wave 2.2)

**Wave Plan**:
```json
{
  "effort_id": "2.2.1",
  "depends_on": ["integration:phase2-wave2.1"],
  "parallelizable": false
}
{
  "effort_id": "2.2.2",
  "depends_on": ["effort:2.2.1"],
  "parallelizable": false
}
```

**Correct Flow** (R603 compliant):
1. CREATE_NEXT_INFRASTRUCTURE → Creates 2.2.1 (depends on integration, satisfied)
2. SPAWN_SW_ENGINEERS → Implements 2.2.1
3. MONITORING_EFFORT_REVIEWS → Approves 2.2.1
4. CREATE_NEXT_INFRASTRUCTURE → Creates 2.2.2 (depends on 2.2.1, now satisfied)
5. SPAWN_SW_ENGINEERS → Implements 2.2.2

### Example 2: Mixed Parallel and Sequential

**Wave Plan**:
```json
{
  "effort_id": "3.1.1",
  "depends_on": ["integration:phase3-wave1.0"],
  "parallelizable": true
}
{
  "effort_id": "3.1.2",
  "depends_on": ["integration:phase3-wave1.0"],
  "parallelizable": true
}
{
  "effort_id": "3.1.3",
  "depends_on": ["effort:3.1.1", "effort:3.1.2"],
  "parallelizable": false
}
```

**Correct Flow** (R603 compliant):
1. CREATE_NEXT_INFRASTRUCTURE → Creates 3.1.1 (parallel, deps satisfied)
2. CREATE_NEXT_INFRASTRUCTURE → Creates 3.1.2 (parallel, deps satisfied)
3. SPAWN_SW_ENGINEERS → Implements 3.1.1 and 3.1.2 in parallel
4. MONITORING_EFFORT_REVIEWS → Approves both 3.1.1 and 3.1.2
5. CREATE_NEXT_INFRASTRUCTURE → Creates 3.1.3 (sequential, deps now satisfied)
6. SPAWN_SW_ENGINEERS → Implements 3.1.3

## Anti-Patterns

### ❌ WRONG: Batch Creating All Efforts

```bash
# This violates R603
for effort in $(get_all_wave_efforts); do
    create_effort_infrastructure "$effort"  # Creates ALL at once
done
```

### ❌ WRONG: Ignoring depends_on Field

```bash
# This violates R603
next_effort=$(find_first_uncreated_effort)  # Doesn't check dependencies
create_effort_infrastructure "$next_effort"
```

### ✅ RIGHT: Dependency-Aware Creation

```bash
# This complies with R603
for effort in $(get_all_wave_efforts); do
    if check_dependencies_satisfied "$effort"; then
        create_effort_infrastructure "$effort"
        break  # Create one at a time for sequential
    fi
done
```

## Summary

R603 ensures sequential efforts are implemented in the correct order by creating infrastructure incrementally, only when dependencies are satisfied. This maintains the cascade pattern and prevents R509 validation failures.

**Key Principle**: Infrastructure creation timing must respect R213 dependency metadata.

**Implementation**: CHECK dependencies → CREATE infrastructure → WAIT for completion → CREATE next

**Result**: Clean cascade, proper base branches, successful validation, smooth merges.
