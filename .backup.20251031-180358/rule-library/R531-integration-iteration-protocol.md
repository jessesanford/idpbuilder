# 🔴🔴🔴 RULE R531: Integration Iteration Protocol (SUPREME LAW)

## Classification
- **Category**: Integration / State Management / Convergence Tracking
- **Criticality Level**: 🔴🔴🔴 SUPREME LAW
- **Enforcement**: MANDATORY - System integrity depends on iteration tracking
- **Penalty**: -100% for missing iteration tracking, -75% for exceeding max iterations without escalation
- **Exit Code**: 531

## The Rule

**ALL integration operations (wave, phase, project) MUST track iteration counts and enforce convergence limits. Integration containers must iterate until clean or escalate to ERROR_RECOVERY when max iterations exceeded.**

## Core Principles

### 1. Integration is Iterative by Nature

Integration operations are NOT one-shot processes:

```
┌──────────────────────────────────────────┐
│  WAVE INTEGRATION ITERATION CONTAINER    │
├──────────────────────────────────────────┤
│  Iteration 1: Merge → Build → Review     │
│               Found 3 bugs              │
│                                          │
│  Iteration 2: Fix bugs → Merge → Review │
│               Found 1 new bug           │
│                                          │
│  Iteration 3: Fix bug → Merge → Review  │
│               CLEAN! ✅                  │
│                                          │
│  Integration complete - proceed to next  │
└──────────────────────────────────────────┘
```

### 2. Iteration Counters Provide Convergence Metrics

Each integration level tracks:
- **Current iteration**: How many times have we attempted this integration?
- **Max iterations**: Safety limit to prevent infinite loops
- **Last progress timestamp**: When did we last make forward progress?

### 3. Three Integration Levels

```
┌─────────────────────────────────────────┐
│  PROJECT INTEGRATION (Level 3)          │
│  - Integrates all phase branches        │
│  - Tracked in: current_project.iteration│
│  - Max: 10 iterations                   │
├─────────────────────────────────────────┤
│  PHASE INTEGRATION (Level 2)            │
│  - Integrates all wave branches         │
│  - Tracked in: current_phase.iteration  │
│  - Max: 10 iterations                   │
├─────────────────────────────────────────┤
│  WAVE INTEGRATION (Level 1)             │
│  - Integrates all effort branches       │
│  - Tracked in: current_wave.iteration   │
│  - Max: 10 iterations                   │
└─────────────────────────────────────────┘
```

## Mandatory State File Tracking

### Required Fields in orchestrator-state-v3.json

```json
{
  "project_progression": {
    "current_wave": {
      "wave_id": "P1W1",
      "iteration": 0,
      "max_iterations": 10,
      "status": "IN_PROGRESS"
    },
    "current_phase": {
      "phase_id": "P1",
      "iteration": 0,
      "max_iterations": 10,
      "status": "IN_PROGRESS"
    },
    "current_project": {
      "project_id": "MyProject",
      "iteration": 0,
      "max_iterations": 10,
      "status": "IN_PROGRESS"
    },
    "iteration_tracking": {
      "last_progress_timestamp": "2025-10-31T12:00:00Z",
      "total_iterations": 0
    }
  }
}
```

## Iteration Management Tool

### Using tools/iteration-manager.sh

```bash
# Increment iteration counter
NEW_ITERATION=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" \
  increment_iteration WAVE)

# Check if max iterations exceeded
STATUS=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" \
  check_max_iterations WAVE)

if [ "$STATUS" = "EXCEEDED" ]; then
    echo "❌ Max iterations exceeded - escalating to ERROR_RECOVERY"
    NEXT_STATE="ERROR_RECOVERY"
fi

# Get current iteration count
CURRENT=$(bash "$CLAUDE_PROJECT_DIR/tools/iteration-manager.sh" \
  get_iteration_count WAVE)
```

## State Integration Points

### START_WAVE_ITERATION State

**BLOCKING CHECKLIST ITEMS:**

1. **Increment wave iteration counter**
   ```bash
   NEW_ITERATION=$(bash tools/iteration-manager.sh increment_iteration WAVE)
   echo "✅ CHECKLIST[1]: Wave iteration incremented to ${NEW_ITERATION}"
   ```

2. **Check max iterations not exceeded**
   ```bash
   ITERATION_STATUS=$(bash tools/iteration-manager.sh check_max_iterations WAVE)

   if [ "$ITERATION_STATUS" = "EXCEEDED" ]; then
       echo "❌ Max iterations exceeded for wave integration"
       PROPOSED_NEXT_STATE="ERROR_RECOVERY"
       echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
       exit 1
   fi

   CURRENT=$(bash tools/iteration-manager.sh get_iteration_count WAVE)
   MAX=$(jq -r '.project_progression.current_wave.max_iterations // 10' orchestrator-state-v3.json)
   echo "✅ CHECKLIST[2]: Max iterations check passed (${CURRENT}/${MAX})"
   ```

3. **Determine next state**
   ```bash
   PROPOSED_NEXT_STATE="INTEGRATE_WAVE_EFFORTS"
   echo "✅ CHECKLIST[3]: Next state determined: ${PROPOSED_NEXT_STATE}"
   ```

### START_PHASE_ITERATION State

**BLOCKING CHECKLIST ITEMS:**

1. **Increment phase iteration counter**
   ```bash
   NEW_ITERATION=$(bash tools/iteration-manager.sh increment_iteration PHASE)
   echo "✅ CHECKLIST[1]: Phase iteration incremented to ${NEW_ITERATION}"
   ```

2. **Check max iterations not exceeded**
   ```bash
   ITERATION_STATUS=$(bash tools/iteration-manager.sh check_max_iterations PHASE)

   if [ "$ITERATION_STATUS" = "EXCEEDED" ]; then
       echo "❌ Max iterations exceeded for phase integration"
       PROPOSED_NEXT_STATE="ERROR_RECOVERY"
       echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
       exit 1
   fi

   echo "✅ CHECKLIST[2]: Max iterations check passed"
   ```

### START_PROJECT_ITERATION State

**BLOCKING CHECKLIST ITEMS:**

1. **Increment project iteration counter**
   ```bash
   NEW_ITERATION=$(bash tools/iteration-manager.sh increment_iteration PROJECT)
   echo "✅ CHECKLIST[1]: Project iteration incremented to ${NEW_ITERATION}"
   ```

2. **Check max iterations not exceeded**
   ```bash
   ITERATION_STATUS=$(bash tools/iteration-manager.sh check_max_iterations PROJECT)

   if [ "$ITERATION_STATUS" = "EXCEEDED" ]; then
       echo "❌ Max iterations exceeded for project integration"
       PROPOSED_NEXT_STATE="ERROR_RECOVERY"
       echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
       exit 1
   fi

   echo "✅ CHECKLIST[2]: Max iterations check passed"
   ```

## Convergence Detection

### What Indicates Progress?

**Integration iteration should converge** - each iteration should reduce bugs:

```
Iteration 1: 10 bugs found → Fix → Push
Iteration 2: 3 bugs found  → Fix → Push
Iteration 3: 0 bugs found  → CONVERGED! ✅
```

### What Indicates Non-Convergence?

**Escalate to ERROR_RECOVERY if:**

1. **Max iterations exceeded** (10 by default)
   - Indicates systemic problems
   - Cannot achieve clean integration
   - Needs human intervention

2. **No progress made** (same bugs reappearing)
   - Tracked via bug_registry changes
   - If iteration N+1 has same bugs as iteration N
   - Indicates ineffective fixes

3. **Bug count increasing**
   - Iteration 1: 3 bugs
   - Iteration 2: 5 bugs
   - Iteration 3: 8 bugs
   - Indicates cascading problems

## Integration with R336

**R336 mandates wave integration before next wave**

R531 defines HOW to track those iterations:

```
R336: "MUST integrate wave before next wave starts"
R531: "Track wave integration iterations, max 10, escalate if exceeded"
```

**Together they ensure:**
- Wave integration is mandatory (R336)
- Wave integration will converge or escalate (R531)
- No infinite loops (R531 max iterations)
- Clean state before next wave (R336 + R531)

## Integration with R520

**R520 tracks integration attempt metadata**

R531 focuses on iteration counters:

```
R520: Integration attempt tracking (attempt_count, last_attempt_result, ready_for_retry)
R531: Iteration counter management (increment, check max, enforce limits)
```

**Relationship:**
- R520 tracks WHAT happened in each attempt
- R531 tracks HOW MANY times we've tried
- Both work together for complete integration tracking

## Error Escalation Protocol

### When Max Iterations Exceeded

```bash
# In START_WAVE_ITERATION state
if [ "$ITERATION_STATUS" = "EXCEEDED" ]; then
    echo "🔴🔴🔴 R531 VIOLATION: Max iterations exceeded"
    echo "Current iteration: $(get_iteration_count WAVE)"
    echo "Maximum allowed: 10"
    echo ""
    echo "This indicates:"
    echo "- Systemic integration problems"
    echo "- Fixes are not effective"
    echo "- Bug count not converging"
    echo ""
    echo "REQUIRED ACTION: Escalate to ERROR_RECOVERY"
    echo "Human intervention needed to:"
    echo "- Analyze root cause of non-convergence"
    echo "- Revise integration strategy"
    echo "- Potentially restructure efforts"

    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MAX_ITERATIONS_EXCEEDED"
    exit 531
fi
```

### Recovery Procedures

**In ERROR_RECOVERY state:**

1. **Analyze iteration history**
   ```bash
   # Review what happened in each iteration
   jq '.integration_attempts[]' orchestrator-state-v3.json
   ```

2. **Identify convergence blockers**
   - Same bugs reappearing?
   - New bugs introduced by fixes?
   - Integration conflicts increasing?

3. **Reset iteration counter after fixes**
   ```bash
   # After resolving systemic issues
   jq '.project_progression.current_wave.iteration = 0' \
     orchestrator-state-v3.json > tmp.json
   mv tmp.json orchestrator-state-v3.json
   ```

## Verification Protocol

### Pre-Transition Verification

```bash
# Before transitioning to integration
echo "🔍 R531: Verifying iteration tracking..."

# Check iteration field exists
ITERATION=$(jq -r '.project_progression.current_wave.iteration' orchestrator-state-v3.json)
if [ "$ITERATION" = "null" ]; then
    echo "❌ R531 VIOLATION: No iteration tracking for wave"
    exit 531
fi

# Check max_iterations field exists
MAX=$(jq -r '.project_progression.current_wave.max_iterations' orchestrator-state-v3.json)
if [ "$MAX" = "null" ]; then
    echo "⚠️ R531 WARNING: No max_iterations set, using default 10"
    jq '.project_progression.current_wave.max_iterations = 10' \
      orchestrator-state-v3.json > tmp.json
    mv tmp.json orchestrator-state-v3.json
fi

echo "✅ R531: Iteration tracking verified (${ITERATION}/${MAX})"
```

### Post-Integration Verification

```bash
# After successful integration
echo "📊 R531: Recording final iteration count..."

FINAL_ITERATION=$(get_iteration_count WAVE)
jq --arg iter "$FINAL_ITERATION" \
  '.completed_integrations[-1].total_iterations = ($iter | tonumber)' \
  orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json

echo "✅ Wave integration completed in ${FINAL_ITERATION} iterations"
```

## Common Violations

### VIOLATION 1: Not Incrementing Iteration Counter

```bash
# ❌ WRONG: Starting integration without incrementing
current_state: INTEGRATE_WAVE_EFFORTS
current_wave.iteration: 0  # Still 0 on iteration 3!

# ✅ CORRECT: Increment in START_WAVE_ITERATION
NEW_ITERATION=$(bash tools/iteration-manager.sh increment_iteration WAVE)
```

### VIOLATION 2: Not Checking Max Iterations

```bash
# ❌ WRONG: Proceeding even when limit exceeded
current_wave.iteration: 11
current_wave.max_iterations: 10
current_state: INTEGRATE_WAVE_EFFORTS  # Should be ERROR_RECOVERY!

# ✅ CORRECT: Check and escalate
if [ "$ITERATION_STATUS" = "EXCEEDED" ]; then
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
fi
```

### VIOLATION 3: Missing Iteration Tracking Fields

```json
// ❌ WRONG: No iteration tracking
{
  "current_wave": {
    "wave_id": "P1W1"
    // Missing: iteration, max_iterations
  }
}

// ✅ CORRECT: Complete iteration tracking
{
  "current_wave": {
    "wave_id": "P1W1",
    "iteration": 2,
    "max_iterations": 10
  }
}
```

## Grading Impact

- **Missing iteration tracking fields**: -50% (cannot track convergence)
- **Not incrementing iteration counter**: -40% (incorrect metrics)
- **Not checking max iterations**: -75% (allows infinite loops)
- **Exceeding max without ERROR_RECOVERY**: -100% (AUTOMATIC FAILURE)
- **Corrupting iteration counters**: -60% (invalid state)

## Success Metrics

- ✅ All wave integrations converge within 10 iterations
- ✅ All phase integrations converge within 10 iterations
- ✅ All project integrations converge within 10 iterations
- ✅ Iteration counters accurately reflect integration attempts
- ✅ Max iterations exceeded always escalates to ERROR_RECOVERY
- ✅ Iteration tracking never corrupted or missing

## Integration with Other Rules

- **R336**: Mandatory wave integration (R531 tracks iterations)
- **R520**: Integration attempt tracking (R531 provides counters)
- **R288**: State file updates (R531 requires iteration fields)
- **R510**: Checklist compliance (R531 defines checklist items)
- **R406**: Fix cascade tracking (R531 iterations drive cascade)

## Visual Summary

```
┌────────────────────────────────────────────────┐
│  R531: INTEGRATION ITERATION PROTOCOL          │
├────────────────────────────────────────────────┤
│                                                │
│  START_*_ITERATION states:                     │
│    1. Increment iteration counter              │
│    2. Check max iterations                     │
│    3. Determine next state                     │
│                                                │
│  Integration states:                           │
│    - Execute integration                       │
│    - Track bugs found                          │
│    - Record progress                           │
│                                                │
│  Convergence tracking:                         │
│    - Iteration N: X bugs                       │
│    - Iteration N+1: Y bugs (Y < X = progress)  │
│    - Iteration N+2: 0 bugs (CONVERGED!)        │
│                                                │
│  Safety limits:                                │
│    - Max 10 iterations per integration         │
│    - Exceeded = escalate to ERROR_RECOVERY     │
│    - Prevents infinite loops                   │
│                                                │
└────────────────────────────────────────────────┘
```

## The Golden Rule

**Integration is iterative. Track iterations. Enforce limits. Converge or escalate.**

Every integration container MUST track how many times it has attempted to reach a clean state. If it cannot converge within max iterations, escalate to ERROR_RECOVERY for human intervention.

---

**REMEMBER**: Integration iterations provide the convergence metrics that prove your system is moving toward stability. Missing or corrupted iteration tracking means you have no visibility into convergence!
