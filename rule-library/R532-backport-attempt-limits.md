# 🔴🔴🔴 RULE R532: Backport Attempt Limits (SUPREME LAW)

## Classification
- **Category**: Integration / Loop Prevention / Safety Limits
- **Criticality Level**: 🔴🔴🔴 SUPREME LAW
- **Enforcement**: MANDATORY - Prevents infinite loops within same iteration
- **Penalty**: -100% for missing backport attempt tracking, allows runaway automation
- **Exit Code**: 532
- **Related Rules**: R531 (Integration Iteration Protocol), R321 (Immediate Backport Protocol)

## The Rule

**ALL backport fix operations within the SAME integration iteration MUST be limited to prevent infinite loops. If backport attempts exceed max_backport_attempts_per_iteration, the system MUST escalate to ERROR_RECOVERY.**

## The Problem This Rule Solves

### Infinite Loop Vulnerability

R531 tracks **iteration counters** for full integration cycles (integrate→build→review→fix). However, R531 intentionally does NOT increment the iteration counter when retrying the SAME attempt after build failures or backport fixes.

This creates a vulnerability:

```
START_WAVE_ITERATION (iteration=1)
  ↓
INTEGRATE_WAVE_EFFORTS (build fails)
  ↓
IMMEDIATE_BACKPORT_REQUIRED (fix upstream)
  ↓
START_WAVE_ITERATION (iteration=1 STILL - no increment per R531)
  ↓
INTEGRATE_WAVE_EFFORTS (build fails AGAIN)
  ↓
IMMEDIATE_BACKPORT_REQUIRED (fix upstream)
  ↓
[INFINITE LOOP - iteration counter stays at 1 forever!]
```

**The Problem:**
- R531's max_iterations check never triggers (iteration stays constant)
- System can loop forever spawning SW Engineers
- No escalation to ERROR_RECOVERY
- Runaway automation costs

**The Solution (R532):**
Track a SEPARATE counter for backport attempts within the same iteration.

## Dual Counter System

R532 works alongside R531 to provide complete loop protection:

```
┌──────────────────────────────────────────────┐
│  ITERATION 1 (R531 iteration counter)        │
├──────────────────────────────────────────────┤
│                                              │
│  START_WAVE_ITERATION                        │
│    backport_attempts_this_iteration = 0      │
│    ↓                                         │
│  INTEGRATE_WAVE_EFFORTS (build fails)        │
│    ↓                                         │
│  IMMEDIATE_BACKPORT_REQUIRED                 │
│    backport_attempts_this_iteration = 1      │
│    ↓                                         │
│  START_WAVE_ITERATION (RETRY SAME ITER)      │
│    Check: backport_attempts (1/3) ✅         │
│    ↓                                         │
│  INTEGRATE_WAVE_EFFORTS (build fails)        │
│    ↓                                         │
│  IMMEDIATE_BACKPORT_REQUIRED                 │
│    backport_attempts_this_iteration = 2      │
│    ↓                                         │
│  START_WAVE_ITERATION (RETRY SAME ITER)      │
│    Check: backport_attempts (2/3) ✅         │
│    ↓                                         │
│  INTEGRATE_WAVE_EFFORTS (build fails)        │
│    ↓                                         │
│  IMMEDIATE_BACKPORT_REQUIRED                 │
│    backport_attempts_this_iteration = 3      │
│    ↓                                         │
│  START_WAVE_ITERATION (RETRY SAME ITER)      │
│    Check: backport_attempts (3/3) ⚠️         │
│    ↓                                         │
│  INTEGRATE_WAVE_EFFORTS (build fails)        │
│    ↓                                         │
│  IMMEDIATE_BACKPORT_REQUIRED                 │
│    backport_attempts_this_iteration = 4      │
│    Check: EXCEEDED (4/3) ❌                  │
│    → ERROR_RECOVERY (ESCALATE!)              │
│                                              │
└──────────────────────────────────────────────┘
```

## Required State File Fields

### orchestrator-state-v3.json Schema

```json
{
  "project_progression": {
    "current_wave": {
      "wave_id": "P1W1",
      "iteration": 1,
      "max_iterations": 10,
      "backport_attempts_this_iteration": 2,
      "max_backport_attempts_per_iteration": 3,
      "status": "IN_PROGRESS"
    },
    "current_phase": {
      "phase_id": "P1",
      "iteration": 1,
      "max_iterations": 10,
      "backport_attempts_this_iteration": 0,
      "max_backport_attempts_per_iteration": 3,
      "status": "IN_PROGRESS"
    },
    "current_project": {
      "project_id": "MyProject",
      "iteration": 1,
      "max_iterations": 10,
      "backport_attempts_this_iteration": 1,
      "max_backport_attempts_per_iteration": 3,
      "status": "IN_PROGRESS"
    }
  }
}
```

**New Fields (per R532):**
- `backport_attempts_this_iteration`: Counter for backport cycles within current iteration
- `max_backport_attempts_per_iteration`: Safety limit (default: 3)

**Existing Fields (per R531):**
- `iteration`: Full integration cycle counter
- `max_iterations`: Safety limit for full cycles (default: 10)

## Backport Attempt Manager Tool

### Using tools/backport-attempt-manager.sh

```bash
# Increment backport attempt counter
NEW_COUNT=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
  increment_backport_attempts WAVE)

# Get current backport attempt count
CURRENT=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
  get_backport_attempt_count WAVE)

# Reset to 0 (when starting new iteration)
bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
  reset_backport_attempts WAVE

# Check if max exceeded
STATUS=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
  check_max_backport_attempts WAVE)

if [ "$STATUS" = "EXCEEDED" ]; then
    echo "❌ Max backport attempts exceeded - escalating to ERROR_RECOVERY"
    NEXT_STATE="ERROR_RECOVERY"
fi
```

## When to Increment Backport Attempt Counter

### Increment Points

Increment `backport_attempts_this_iteration` when:

1. **Entering IMMEDIATE_BACKPORT_REQUIRED**
   - Build/test failure detected
   - About to spawn SW Engineer for upstream fix
   - This is a new backport attempt

2. **Entering MONITORING_BACKPORT_PROGRESS** (optional - if tracked separately)
   - Monitoring SW Engineer fixing upstream bugs
   - Waiting for fixes to complete

### Reset Points

Reset `backport_attempts_this_iteration` to 0 when:

1. **Starting NEW iteration** (iteration counter increments)
   - Coming from SETUP_*_INFRASTRUCTURE
   - Coming from FIX_*_UPSTREAM_BUGS
   - Per R531, this is a new full integration attempt

### Check Points

Check `backport_attempts_this_iteration` against `max_backport_attempts_per_iteration`:

1. **In START_*_ITERATION states**
   - When coming from IMMEDIATE_BACKPORT_REQUIRED
   - When coming from MONITORING_BACKPORT_PROGRESS
   - When coming from ERROR_RECOVERY
   - BEFORE attempting integration again

2. **In IMMEDIATE_BACKPORT_REQUIRED**
   - After incrementing counter
   - Before spawning SW Engineers

## State Integration Points

### START_WAVE_ITERATION State

**BLOCKING CHECKLIST ITEM (add to existing R531 logic):**

```bash
# After R531 iteration counter logic
PREVIOUS_STATE=$(jq -r '.state_machine.previous_state // "SETUP_WAVE_INFRASTRUCTURE"' \
  "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

case "$PREVIOUS_STATE" in
    SETUP_WAVE_INFRASTRUCTURE|FIX_WAVE_UPSTREAM_BUGS)
        # R531: New iteration - increment iteration counter
        # R532: New iteration - reset backport counter
        bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" reset_backport_attempts WAVE
        echo "✅ CHECKLIST[1b]: Backport attempts counter reset to 0 (new iteration per R532)"
        ;;

    IMMEDIATE_BACKPORT_REQUIRED|MONITORING_BACKPORT_PROGRESS|ERROR_RECOVERY)
        # R531: Same iteration - do NOT increment iteration counter
        # R532: Check backport attempt limit
        BACKPORT_STATUS=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
            check_max_backport_attempts WAVE)

        if [ "$BACKPORT_STATUS" = "EXCEEDED" ]; then
            BACKPORT_COUNT=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
                get_backport_attempt_count WAVE)
            MAX_BACKPORT=$(jq -r '.project_progression.current_wave.max_backport_attempts_per_iteration // 3' \
                orchestrator-state-v3.json)

            echo "❌ R532 VIOLATION: Max backport attempts exceeded"
            echo "   Current iteration: $(get_iteration_count WAVE)"
            echo "   Backport attempts this iteration: ${BACKPORT_COUNT}/${MAX_BACKPORT}"
            echo "   This indicates ineffective fixes or systemic problems"
            echo "   REQUIRED ACTION: Escalate to ERROR_RECOVERY"

            PROPOSED_NEXT_STATE="ERROR_RECOVERY"
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MAX_BACKPORT_ATTEMPTS_EXCEEDED"
            exit 532
        fi

        BACKPORT_COUNT=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
            get_backport_attempt_count WAVE)
        MAX_BACKPORT=$(jq -r '.project_progression.current_wave.max_backport_attempts_per_iteration // 3' \
            orchestrator-state-v3.json)
        echo "✅ CHECKLIST[1b]: Backport attempts check passed (${BACKPORT_COUNT}/${MAX_BACKPORT} per R532)"
        ;;
esac
```

### START_PHASE_ITERATION State

Same pattern as WAVE, but use `PHASE` level:
```bash
bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" <command> PHASE
```

### START_PROJECT_ITERATION State

Same pattern as WAVE, but use `PROJECT` level:
```bash
bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" <command> PROJECT
```

### IMMEDIATE_BACKPORT_REQUIRED State

**BLOCKING CHECKLIST ITEM:**

```bash
# Before spawning SW Engineers for backport fixes
# Increment backport attempt counter per R532
bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" increment_backport_attempts WAVE

BACKPORT_COUNT=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
    get_backport_attempt_count WAVE)
MAX_BACKPORT=$(jq -r '.project_progression.current_wave.max_backport_attempts_per_iteration // 3' \
    orchestrator-state-v3.json)

echo "✅ CHECKLIST[X]: Backport attempt counter incremented to ${BACKPORT_COUNT}/${MAX_BACKPORT} (R532)"

# Note: The actual max check happens in START_WAVE_ITERATION
# We increment here but check there (before retrying integration)
```

### MONITORING_BACKPORT_PROGRESS State

**OPTIONAL CHECK:**

```bash
# When deciding next state after monitoring SW Engineers
BACKPORT_STATUS=$(bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" \
    check_max_backport_attempts WAVE)

if [ "$BACKPORT_STATUS" = "EXCEEDED" ]; then
    echo "❌ R532: Max backport attempts exceeded during monitoring"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MAX_BACKPORT_ATTEMPTS_EXCEEDED"
    exit 532
fi
```

## Default Values

### max_backport_attempts_per_iteration

**Default: 3 attempts per iteration**

**Rationale:**
- 3 backport cycles should be enough for most build/test fixes
- More than 3 indicates:
  - Fixes are not effective
  - Upstream branches have systemic problems
  - Architecture issues need review
- Human intervention needed after 3 failed attempts

**Configurable:**
Users can adjust in `orchestrator-state-v3.json`:
```json
{
  "current_wave": {
    "max_backport_attempts_per_iteration": 5
  }
}
```

## Relationship with R531

### R531: Iteration Counter (Full Cycles)

Tracks complete integration attempts:
```
Iteration 1: START → INTEGRATE → BUILD → REVIEW → FIX
Iteration 2: START → INTEGRATE → BUILD → REVIEW (clean!)
```

**Increments:** Only when full cycle completes (from FIX state)
**Max:** 10 iterations (full cycles)
**Purpose:** Track convergence over time

### R532: Backport Attempt Counter (Retries Within Same Cycle)

Tracks retry attempts within SAME iteration:
```
Iteration 1:
  Attempt 1: START → INTEGRATE (build fails) → BACKPORT
  Attempt 2: START → INTEGRATE (build fails) → BACKPORT
  Attempt 3: START → INTEGRATE (build succeeds) → REVIEW → FIX
```

**Increments:** Every backport fix cycle within same iteration
**Max:** 3 attempts per iteration
**Purpose:** Prevent infinite loops at same iteration number

### Working Together

```
┌─────────────────────────────────────────────┐
│  R531 Iteration 1                           │
├─────────────────────────────────────────────┤
│  R532 Backport Attempt 1 (build fails)      │
│  R532 Backport Attempt 2 (build fails)      │
│  R532 Backport Attempt 3 (build succeeds)   │
│  → Integration completes                    │
│  → Review finds bugs                        │
│  → FIX state applies fixes                  │
└─────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────┐
│  R531 Iteration 2 (counter increments)      │
├─────────────────────────────────────────────┤
│  R532 Backport Attempts RESET to 0          │
│  Fresh start for new iteration              │
│  → Integration succeeds first try           │
│  → Review clean                             │
└─────────────────────────────────────────────┘
```

## Error Escalation Protocol

### When Max Backport Attempts Exceeded

```bash
if [ "$BACKPORT_STATUS" = "EXCEEDED" ]; then
    echo "🔴🔴🔴 R532 VIOLATION: Max backport attempts exceeded"
    echo "Current iteration: $(get_iteration_count WAVE)"
    echo "Backport attempts this iteration: $(get_backport_attempt_count WAVE)"
    echo "Maximum allowed: $(jq -r '.project_progression.current_wave.max_backport_attempts_per_iteration' orchestrator-state-v3.json)"
    echo ""
    echo "This indicates:"
    echo "- Backport fixes are not effective"
    echo "- Upstream branches have systemic problems"
    echo "- Same issues reoccurring despite fixes"
    echo "- Possible architecture issues"
    echo ""
    echo "REQUIRED ACTION: Escalate to ERROR_RECOVERY"
    echo "Human intervention needed to:"
    echo "- Analyze root cause of repeated failures"
    echo "- Review upstream branch quality"
    echo "- Consider architectural changes"
    echo "- Manually fix problematic code"

    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MAX_BACKPORT_ATTEMPTS_EXCEEDED"
    exit 532
fi
```

### In ERROR_RECOVERY State

1. **Review backport history**
   ```bash
   jq '.project_progression.current_wave.backport_attempts_this_iteration' orchestrator-state-v3.json
   ```

2. **Analyze repeated failures**
   - What keeps failing? (build? tests? specific module?)
   - Are fixes introducing new bugs?
   - Are upstream branches fundamentally broken?

3. **Manual intervention options**
   - Fix upstream branches manually
   - Restructure efforts differently
   - Adjust architecture
   - Reset backport counter after fixing root cause

4. **Reset counter after manual fixes**
   ```bash
   bash "$CLAUDE_PROJECT_DIR/tools/backport-attempt-manager.sh" reset_backport_attempts WAVE
   ```

## Verification Protocol

### Pre-Transition Verification

```bash
# Before starting integration iteration
echo "🔍 R532: Verifying backport attempt tracking..."

# Check field exists
BACKPORT_COUNT=$(jq -r '.project_progression.current_wave.backport_attempts_this_iteration' \
    orchestrator-state-v3.json)
if [ "$BACKPORT_COUNT" = "null" ]; then
    echo "⚠️ R532 WARNING: No backport attempt tracking, initializing to 0"
    jq '.project_progression.current_wave.backport_attempts_this_iteration = 0' \
      orchestrator-state-v3.json > tmp.json
    mv tmp.json orchestrator-state-v3.json
fi

# Check max configured
MAX_BACKPORT=$(jq -r '.project_progression.current_wave.max_backport_attempts_per_iteration' \
    orchestrator-state-v3.json)
if [ "$MAX_BACKPORT" = "null" ]; then
    echo "⚠️ R532 WARNING: No max_backport_attempts set, using default 3"
    jq '.project_progression.current_wave.max_backport_attempts_per_iteration = 3' \
      orchestrator-state-v3.json > tmp.json
    mv tmp.json orchestrator-state-v3.json
fi

echo "✅ R532: Backport attempt tracking verified"
```

## Common Violations

### VIOLATION 1: Not Incrementing Backport Counter

```bash
# ❌ WRONG: Entering IMMEDIATE_BACKPORT_REQUIRED without incrementing
current_wave.backport_attempts_this_iteration: 1  # Still 1 on 3rd backport!

# ✅ CORRECT: Increment in IMMEDIATE_BACKPORT_REQUIRED
bash tools/backport-attempt-manager.sh increment_backport_attempts WAVE
```

### VIOLATION 2: Not Checking Backport Limit

```bash
# ❌ WRONG: Allowing unlimited backport retries
current_wave.backport_attempts_this_iteration: 7  # Way over limit!
current_state: START_WAVE_ITERATION  # Should be ERROR_RECOVERY!

# ✅ CORRECT: Check in START_WAVE_ITERATION
if [ "$BACKPORT_STATUS" = "EXCEEDED" ]; then
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
fi
```

### VIOLATION 3: Not Resetting on New Iteration

```bash
# ❌ WRONG: Backport counter carries over to new iteration
Iteration 1: backport_attempts = 3
Iteration 2: backport_attempts = 3  # Should be 0!

# ✅ CORRECT: Reset when iteration increments
if [ "$PREVIOUS_STATE" = "FIX_WAVE_UPSTREAM_BUGS" ]; then
    # New iteration starting
    bash tools/backport-attempt-manager.sh reset_backport_attempts WAVE
fi
```

### VIOLATION 4: Missing State File Fields

```json
// ❌ WRONG: No backport tracking
{
  "current_wave": {
    "wave_id": "P1W1",
    "iteration": 2
    // Missing: backport_attempts_this_iteration, max_backport_attempts_per_iteration
  }
}

// ✅ CORRECT: Complete tracking
{
  "current_wave": {
    "wave_id": "P1W1",
    "iteration": 2,
    "backport_attempts_this_iteration": 1,
    "max_backport_attempts_per_iteration": 3
  }
}
```

## Grading Impact

- **Missing backport attempt tracking fields**: -50% (allows infinite loops)
- **Not incrementing backport counter**: -40% (incorrect tracking)
- **Not checking max backport attempts**: -75% (allows runaway automation)
- **Exceeding max without ERROR_RECOVERY**: -100% (AUTOMATIC FAILURE - infinite loop!)
- **Not resetting on new iteration**: -30% (incorrect metrics)
- **Corrupting backport counters**: -60% (invalid state)

## Success Metrics

- ✅ All backport attempts limited to max_backport_attempts_per_iteration
- ✅ Exceeding limit always escalates to ERROR_RECOVERY
- ✅ Counter resets on new iteration (when R531 iteration increments)
- ✅ No infinite loops possible at any integration level
- ✅ Backward compatible (defaults to 3 if missing)
- ✅ Works for WAVE, PHASE, and PROJECT levels

## Integration with Other Rules

- **R531**: Integration Iteration Protocol (R532 complements with retry tracking)
- **R321**: Immediate Backport Protocol (R532 limits backport attempts)
- **R288**: State File Updates (R532 requires new fields)
- **R510**: Checklist Compliance (R532 defines new checklist items)
- **R336**: Mandatory Wave Integration (R532 prevents integration loops)

## Visual Summary

```
┌───────────────────────────────────────────────┐
│  R532: BACKPORT ATTEMPT LIMITS                │
├───────────────────────────────────────────────┤
│                                               │
│  Prevents infinite loops within SAME          │
│  iteration by limiting backport retries       │
│                                               │
│  Dual Counter System:                         │
│  ├─ R531 iteration counter (full cycles)      │
│  └─ R532 backport counter (retries)           │
│                                               │
│  Flow:                                        │
│  1. START_*_ITERATION                         │
│     - Check backport attempts < max           │
│     - Reset counter if new iteration          │
│  2. INTEGRATE (build fails)                   │
│  3. IMMEDIATE_BACKPORT_REQUIRED               │
│     - Increment backport counter              │
│  4. Back to START_*_ITERATION                 │
│     - Check backport attempts again           │
│     - Escalate if exceeded                    │
│                                               │
│  Safety:                                      │
│  - Max 3 backport attempts per iteration      │
│  - Escalate to ERROR_RECOVERY if exceeded     │
│  - Human intervention required                │
│                                               │
└───────────────────────────────────────────────┘
```

## The Golden Rule

**Backport attempts within the same iteration must be limited. Track them separately from full iteration cycles. Escalate when limits exceeded.**

R531 prevents infinite iterations. R532 prevents infinite retries WITHIN the same iteration. Together they ensure the system always converges or escalates.

---

**REMEMBER**: Without R532, a broken build can loop forever at the same iteration number, spawning agents endlessly without ever triggering max_iterations!
