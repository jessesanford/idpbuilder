# R356: Single Effort Parallelization Optimization

## Rule Type
PERFORMANCE OPTIMIZATION

## Criticality
⚠️⚠️⚠️ WARNING - Improves efficiency and reduces unnecessary state transitions

## Description
When a wave contains only a single effort, parallelization analysis is unnecessary and should be skipped. The orchestrator MUST detect single-effort waves and bypass parallelization analysis states, transitioning directly to spawning states.

## Rationale
Parallelization analysis exists to determine which efforts can run concurrently and which must run sequentially. When there is only ONE effort in a wave:
- Parallelization is impossible by definition
- Analysis wastes time and adds unnecessary state transitions
- The conclusion is always "single effort, no parallelization possible"
- Direct spawning is more efficient and clearer

## Requirements

### 1. Single Effort Detection
States that normally transition to parallelization analysis MUST check effort count:
```bash
EFFORT_COUNT=$(jq '.efforts_in_progress | length' orchestrator-state-v3.json)
if [ "$EFFORT_COUNT" -eq 1 ]; then
    echo "Single effort detected - skipping parallelization analysis (R356)"
    NEXT_STATE="[SPAWN_STATE]"  # Skip to spawn directly
else
    echo "Multiple efforts - analyzing parallelization"
    NEXT_STATE="[ANALYZE_STATE]"  # Normal parallelization analysis
fi
```

### 2. State Transition Optimization

#### For Code Reviewer Effort Planning:
- **VALIDATE_INFRASTRUCTURE** with 1 effort → **SPAWN_CODE_REVIEWERS_EFFORT_PLANNING** (skip ANALYZE_CODE_REVIEWER_PARALLELIZATION)
- **VALIDATE_INFRASTRUCTURE** with 2+ efforts → **ANALYZE_CODE_REVIEWER_PARALLELIZATION** (normal flow)

**CRITICAL**: Infrastructure validation (VALIDATE_INFRASTRUCTURE) is MANDATORY and can NEVER be skipped!

#### For SW Engineer Implementation:
- **WAITING_FOR_EFFORT_PLANS** with 1 effort → **SPAWN_SW_ENGINEERS** (skip ANALYZE_IMPLEMENTATION_PARALLELIZATION)
- **WAITING_FOR_EFFORT_PLANS** with 2+ efforts → **ANALYZE_IMPLEMENTATION_PARALLELIZATION** (normal flow)

### 3. Early Exit in Analysis States
If a parallelization analysis state is entered with single effort, it MUST immediately exit:
```bash
# In ANALYZE_*_PARALLELIZATION states
EFFORT_COUNT=$(jq '.efforts_in_progress | length' orchestrator-state-v3.json)
if [ "$EFFORT_COUNT" -eq 1 ]; then
    echo "🎯 R356: Single effort detected - no parallelization needed"
    echo "Strategy: Spawn single agent for the sole effort"
    # Immediately transition to spawn state
    safe_state_transition "SPAWN_[TYPE]" "Single effort - R356 optimization"
    exit 0
fi
```

### 4. Documentation Requirements
When skipping parallelization analysis, the orchestrator MUST:
1. Log the R356 optimization in the state file
2. Document the single-effort detection in transition_reason
3. Note in work logs that parallelization was skipped

## Affected States

### States That Should Check Effort Count:
1. **VALIDATE_INFRASTRUCTURE** - After validation passes, before transitioning to ANALYZE_CODE_REVIEWER_PARALLELIZATION
2. **WAITING_FOR_EFFORT_PLANS** - Before transitioning to ANALYZE_IMPLEMENTATION_PARALLELIZATION
3. **WAITING_FOR_REVIEW_RESULTS** - Before transitioning to parallelization analysis for fixes

**REMEMBER**: CREATE_NEXT_INFRASTRUCTURE must ALWAYS transition to VALIDATE_INFRASTRUCTURE first!

### States That Should Have Early Exit:
1. **ANALYZE_CODE_REVIEWER_PARALLELIZATION** - Exit immediately if single effort
2. **ANALYZE_IMPLEMENTATION_PARALLELIZATION** - Exit immediately if single effort
3. **Any future parallelization analysis states**

## Benefits
1. **Fewer State Transitions**: Eliminates 2 unnecessary states per wave for single efforts
2. **Faster Execution**: No time wasted analyzing obvious scenarios
3. **Clearer Logs**: No confusing "analyzing parallelization... conclusion: cannot parallelize single effort"
4. **Reduced Context**: Less state machine complexity for simple waves

## Edge Cases

### Multiple Efforts But Only One Ready
If multiple efforts exist but only one is ready for processing:
- Still skip parallelization (only 1 actionable effort)
- Document that other efforts are blocked/pending

### Single Effort That Will Be Split
If a single effort is known to require splitting:
- Still skip initial parallelization analysis
- Parallelization analysis happens AFTER split planning (multiple splits)

### Wave With Zero Efforts
If no efforts are found:
- This is an error condition (separate from R356)
- Should not reach parallelization analysis anyway

## Example Implementation

### In VALIDATE_INFRASTRUCTURE:
```bash
# After all validation checks pass
EFFORT_COUNT=$(jq '.efforts_in_progress | length' orchestrator-state-v3.json)

if [ "$EFFORT_COUNT" -eq 1 ]; then
    EFFORT_NAME=$(jq -r '.efforts_in_progress[0].name' orchestrator-state-v3.json)
    echo "═══════════════════════════════════════════════════════════════"
    echo "🎯 R356 OPTIMIZATION: Single Effort Detected"
    echo "═══════════════════════════════════════════════════════════════"
    echo "Effort: $EFFORT_NAME"
    echo "Parallelization Analysis: NOT NEEDED (only 1 effort)"
    echo "Next State: SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (direct)"
    echo "═══════════════════════════════════════════════════════════════"

    safe_state_transition "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" "Single effort - R356 optimization applied"
else
    echo "Multiple efforts detected: $EFFORT_COUNT"
    echo "Proceeding to parallelization analysis..."
    safe_state_transition "ANALYZE_CODE_REVIEWER_PARALLELIZATION" "Multiple efforts require analysis"
fi
```

## Validation
To verify R356 compliance:
1. Check state transitions for single-effort waves
2. Verify parallelization states are skipped when appropriate
3. Confirm spawn states handle single efforts correctly
4. Review logs for R356 optimization messages

## Related Rules
- **R234**: Mandatory State Traversal (R356 provides exception ONLY for skipping parallelization analysis, NOT validation)
- **R507**: Infrastructure Validation Required (MANDATORY - never skipped)
- **R508**: Repository Validation (MANDATORY - never skipped)
- **R151**: Parallel Agent Spawning (not applicable for single efforts)
- **R218**: Code Reviewer Parallelization (skipped for single efforts)
- **R219**: Implementation Parallelization (skipped for single efforts)

## Grading Impact
- **Correct Application**: +5% efficiency bonus
- **Unnecessary Analysis**: -5% for wasting time on single efforts
- **Incorrect Skip**: -20% if parallelization skipped for multiple efforts

## Version History
- **v1.0** (2025-01-15): Initial rule creation for single-effort optimization