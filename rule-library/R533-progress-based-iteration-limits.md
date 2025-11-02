# 🔴🔴🔴 RULE R533 - PROGRESS-BASED ITERATION LIMITS (SUPREME LAW)

**Criticality:** SUPREME LAW - PREVENTS INFINITE LOOPS
**Grading Impact:** -100% for violating iteration limits
**Enforcement:** CONTINUOUS - Every iteration decision point

---

## SUPREME LAW STATEMENT

**ITERATION LIMITS MUST BE PROGRESS-BASED, NOT BLIND COUNTERS. THE SYSTEM MUST TRACK WHETHER ACTUAL PROGRESS IS BEING MADE (BUGS FIXED) OR JUST CHURNING (SAME BUGS REPEATED). TWO-TIERED LIMITS PREVENT BOTH STALLS AND INFINITE SLOW PROGRESS.**

---

## 🚨🚨🚨 THE PROGRESS MANDATE 🚨🚨🚨

### The Problem with Blind Iteration Counters

**OLD APPROACH (WRONG)**:
```bash
iteration++
if [ $iteration -ge 5 ]; then
    ERROR_RECOVERY
fi
```

**Problems:**
- Treats all iterations equally
- Penalizes progress and non-progress the same
- Can't distinguish "stuck" from "slow but steady"
- May give up too early if making progress
- May continue too long if no progress

**NEW APPROACH (RIGHT)**:
```bash
# Track ACTUAL progress
if making_progress; then
    stall_counter=0  # Reset stall counter
else
    stall_counter++  # Increment stall counter
fi

# Two-tiered limits
if [ $stall_counter -ge 5 ]; then
    ERROR_RECOVERY  # No progress for 5 iterations
elif [ $iteration -ge 10 ]; then
    ERROR_RECOVERY  # Too many iterations even with progress
fi
```

---

## 🔴🔴🔴 TWO-TIERED ITERATION LIMITS 🔴🔴🔴

### Tier 1: No-Progress Limit (STRICT)

**Limit:** 5 consecutive iterations with NO progress
**Trigger:** stall_counter >= 5
**Condition:** bugs_closed == 0 in iteration
**Meaning:** System is stuck, not making progress
**Action:** Escalate to ERROR_RECOVERY immediately

**Example Scenario:**
```
Iteration 1: 5 bugs found, 0 fixed → stall_counter=1
Iteration 2: 5 bugs found, 0 fixed → stall_counter=2
Iteration 3: 6 bugs found, 0 fixed → stall_counter=3
Iteration 4: 6 bugs found, 0 fixed → stall_counter=4
Iteration 5: 6 bugs found, 0 fixed → stall_counter=5 → ERROR_RECOVERY
```

**Rationale:**
- If we can't fix ANY bugs in 5 tries, something is fundamentally wrong
- Continuing wastes resources
- Need different approach (architect intervention, replanning)

### Tier 2: Some-Progress Limit (GENEROUS)

**Limit:** 10 total iterations (regardless of progress)
**Trigger:** iteration >= 10
**Condition:** Overall iteration count reaches 10
**Meaning:** Even with progress, taking too long
**Action:** Escalate to ERROR_RECOVERY for replanning

**Example Scenario:**
```
Iteration 1: 5 bugs, 2 fixed → stall_counter=0 (progress)
Iteration 2: 4 bugs, 1 fixed → stall_counter=0 (progress)
Iteration 3: 5 bugs, 1 fixed → stall_counter=0 (progress)
...
Iteration 10: 3 bugs, 1 fixed → iteration>=10 → ERROR_RECOVERY
```

**Rationale:**
- Making progress, but taking too many iterations
- May indicate approach is inefficient
- Need to reassess strategy before continuing

---

## 🔴🔴🔴 PROGRESS DETECTION ALGORITHM 🔴🔴🔴

### Progress Categories

**1. PURE_PROGRESS** ✅✅✅
```yaml
condition:
  bugs_closed: > 0
  net_change: <= 0  # Total bugs decreasing or stable
example: "Fixed 3 bugs, found 1 new → Net -2"
stall_action: Reset stall_counter to 0
interpretation: Excellent! Closing bugs faster than finding them
```

**2. DISCOVERY_PROGRESS** ✅✅
```yaml
condition:
  bugs_closed: > 0
  net_change: > 0 and <= 3  # Few new bugs found
example: "Fixed 2 bugs, found 4 new → Net +2"
stall_action: Reset stall_counter to 0
interpretation: Good! Fixing bugs while discovering hidden issues
```

**3. SLOW_PROGRESS** ✅
```yaml
condition:
  bugs_closed: > 0
  net_change: > 3  # Many new bugs found
example: "Fixed 1 bug, found 6 new → Net +5"
stall_action: Don't increment stall_counter (but track trend)
interpretation: Making progress, but finding many issues
```

**4. STALL** ⚠️
```yaml
condition:
  bugs_closed: == 0
  net_change: == 0  # No movement
example: "0 bugs fixed, 0 new bugs"
stall_action: Increment stall_counter
interpretation: No progress - stuck on same bugs
```

**5. REGRESSION** ❌
```yaml
condition:
  bugs_closed: == 0
  net_change: > 0  # More bugs, none fixed
example: "0 bugs fixed, 3 new bugs found"
stall_action: Increment stall_counter
interpretation: Getting worse - no fixes, more bugs
```

### Progress Calculation Algorithm

```bash
calculate_progress() {
    local previous_iteration="$1"
    local current_iteration="$2"

    # Count bugs closed since last iteration
    bugs_closed=$(jq -r "
        .bugs[] |
        select(.closed_iteration == $current_iteration) |
        .bug_id
    " bug-tracking.json | wc -l)

    # Count total open bugs in each iteration
    prev_open=$(jq -r "
        .iteration_history[] |
        select(.iteration == $previous_iteration) |
        .bugs_open
    " orchestrator-state-v3.json)

    curr_open=$(jq -r "
        .iteration_history[] |
        select(.iteration == $current_iteration) |
        .bugs_open
    " orchestrator-state-v3.json)

    # Calculate net change
    net_change=$((curr_open - prev_open))

    # Determine progress category
    if [ "$bugs_closed" -eq 0 ]; then
        if [ "$net_change" -eq 0 ]; then
            echo "STALL"
        else
            echo "REGRESSION"
        fi
    else
        if [ "$net_change" -le 0 ]; then
            echo "PURE_PROGRESS"
        elif [ "$net_change" -le 3 ]; then
            echo "DISCOVERY_PROGRESS"
        else
            echo "SLOW_PROGRESS"
        fi
    fi

    # Return data
    jq -n \
        --arg progress "$progress_type" \
        --argjson closed "$bugs_closed" \
        --argjson net "$net_change" \
        '{
            progress_category: $progress,
            bugs_closed: $closed,
            net_change: $net
        }'
}
```

---

## 🔴🔴🔴 ITERATION DECISION LOGIC 🔴🔴🔴

### Complete Decision Algorithm

```bash
# Used in START_WAVE_ITERATION, START_PHASE_ITERATION, START_PROJECT_ITERATION

decide_iteration_continuation() {
    local scope="$1"  # WAVE, PHASE, or PROJECT
    local current_iteration="$2"
    local previous_iteration=$((current_iteration - 1))

    echo "📊 R533: Analyzing iteration progress..."

    # Step 1: Calculate progress
    PROGRESS_ANALYSIS=$(bash tools/bug-progress-analyzer.sh \
        analyze_progress "$scope" "$current_iteration" "$previous_iteration")

    PROGRESS_CATEGORY=$(echo "$PROGRESS_ANALYSIS" | jq -r '.progress_category')
    BUGS_CLOSED=$(echo "$PROGRESS_ANALYSIS" | jq -r '.bugs_closed')
    BUGS_REOPENED=$(echo "$PROGRESS_ANALYSIS" | jq -r '.bugs_reopened')
    NET_CHANGE=$(echo "$PROGRESS_ANALYSIS" | jq -r '.net_change')

    echo "  Progress: $PROGRESS_CATEGORY"
    echo "  Bugs closed: $BUGS_CLOSED"
    echo "  Bugs reopened: $BUGS_REOPENED"
    echo "  Net change: $NET_CHANGE"

    # Step 2: Get current stall counter
    STALL_COUNT=$(jq -r ".project_progression.current_${scope,,}.bugs.progress_stalls" \
        orchestrator-state-v3.json)

    # Step 3: Update stall counter based on progress
    case "$PROGRESS_CATEGORY" in
        "PURE_PROGRESS"|"DISCOVERY_PROGRESS")
            # Making progress - reset stall counter
            NEW_STALL_COUNT=0
            echo "  ✅ Progress detected - resetting stall counter"
            ;;
        "SLOW_PROGRESS")
            # Making some progress - don't increment stall counter
            NEW_STALL_COUNT="$STALL_COUNT"
            echo "  ⚠️  Slow progress - maintaining stall counter at $STALL_COUNT"
            ;;
        "STALL"|"REGRESSION")
            # No progress - increment stall counter
            NEW_STALL_COUNT=$((STALL_COUNT + 1))
            echo "  ❌ No progress - incrementing stall counter to $NEW_STALL_COUNT"
            ;;
    esac

    # Step 4: Check for bug flapping (reopened bugs)
    if [ "$BUGS_REOPENED" -gt 0 ]; then
        echo "  🚨 WARNING: $BUGS_REOPENED bugs REOPENED - possible flapping!"
        # Treat flapping as no-progress
        NEW_STALL_COUNT=$((NEW_STALL_COUNT + 1))
    fi

    # Step 5: Apply two-tiered limits
    if [ "$NEW_STALL_COUNT" -ge 5 ]; then
        echo "❌❌❌ R533 NO-PROGRESS LIMIT EXCEEDED ❌❌❌"
        echo "Stall counter: $NEW_STALL_COUNT / 5"
        echo "No bugs fixed for 5 consecutive iterations"
        echo "ESCALATING TO ERROR_RECOVERY"

        # Update state file
        jq --argjson stall "$NEW_STALL_COUNT" \
            ".project_progression.current_${scope,,}.bugs.progress_stalls = \$stall |
             .state_machine.next_state = \"ERROR_RECOVERY\" |
             .state_machine.transition_reason = \"R533: No-progress limit exceeded ($NEW_STALL_COUNT stalls)\"" \
            orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

        return 1  # Signal to transition to ERROR_RECOVERY
    fi

    if [ "$current_iteration" -ge 10 ]; then
        echo "❌❌❌ R533 SOME-PROGRESS LIMIT EXCEEDED ❌❌❌"
        echo "Total iterations: $current_iteration / 10"
        echo "Maximum iterations reached even with progress"
        echo "ESCALATING TO ERROR_RECOVERY for replanning"

        # Update state file
        jq ".state_machine.next_state = \"ERROR_RECOVERY\" |
            .state_machine.transition_reason = \"R533: Some-progress limit exceeded ($current_iteration iterations)\"" \
            orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

        return 1  # Signal to transition to ERROR_RECOVERY
    fi

    # Step 6: Update state file with new stall counter
    jq --argjson stall "$NEW_STALL_COUNT" \
        ".project_progression.current_${scope,,}.bugs.progress_stalls = \$stall" \
        orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    echo "✅ R533: Iteration continuation approved"
    echo "  Stall counter: $NEW_STALL_COUNT / 5"
    echo "  Total iterations: $current_iteration / 10"

    return 0  # Signal to continue with next iteration
}
```

---

## 🔴🔴🔴 BUG FLAPPING DETECTION 🔴🔴🔴

### What is Flapping?

**Definition:** Bug is marked CLOSED, then reappears as OPEN in later iteration

**Example:**
```
Iteration 1: W2-BUG-001 OPEN (authentication fails)
Iteration 2: W2-BUG-001 CLOSED (fixed by commit abc123)
Iteration 3: W2-BUG-001 OPEN again (still fails!) → REOPENED/FLAPPING
```

**Why Flapping is Critical:**
- Indicates incomplete fix
- Suggests systemic issue (not understanding root cause)
- Wastes effort (false confidence in "fix")
- Should be treated as NO PROGRESS

### Flapping Detection

```bash
detect_flapping() {
    local current_iteration="$1"
    local previous_iteration="$2"

    # Find bugs that were CLOSED in previous iteration
    closed_bugs=$(jq -r "
        .bugs[] |
        select(.closed_iteration == $previous_iteration) |
        .bug_id
    " bug-tracking.json)

    # Find bugs that are OPEN in current iteration
    open_bugs=$(jq -r "
        .bugs[] |
        select(.first_found_iteration <= $current_iteration and
               (.closed_iteration == null or .closed_iteration < $current_iteration)) |
        .bug_id
    " bug-tracking.json)

    # Find intersection (reopened bugs)
    reopened=$(comm -12 \
        <(echo "$closed_bugs" | sort) \
        <(echo "$open_bugs" | sort))

    reopen_count=$(echo "$reopened" | grep -c '^' || true)

    # Mark as REOPENED in bug-tracking.json
    for bug_id in $reopened; do
        jq --arg id "$bug_id" \
           --argjson iter "$current_iteration" '
            .bugs |= map(
                if .bug_id == $id then
                    .status = "REOPENED" |
                    .reopened_iteration = $iter |
                    .reopen_count = (.reopen_count // 0) + 1
                else
                    .
                end
            )
        ' bug-tracking.json > tmp.json && mv tmp.json bug-tracking.json

        echo "🚨 FLAPPING DETECTED: $bug_id (reopen count: $(jq -r ".bugs[] | select(.bug_id==\"$bug_id\") | .reopen_count" bug-tracking.json))"
    done

    echo "$reopen_count"
}
```

### Flapping Treatment

```yaml
flapping_impact:
  progress_calculation:
    - Treat REOPENED bugs as NO PROGRESS
    - Increment stall_counter
    - Do NOT count as "bugs_closed" even if other bugs were fixed

  severity_levels:
    1_reopen: "Warning - fix may have been incomplete"
    2_reopens: "Concerning - likely not understanding root cause"
    3_reopens: "Critical - systemic issue, need architect review"

  actions:
    - Flag high-reopen-count bugs in ERROR_RECOVERY report
    - Escalate to architect if reopen_count >= 3
    - Require detailed root cause analysis
```

---

## 🔴🔴🔴 STATE FILE TRACKING 🔴🔴🔴

### Required Fields in orchestrator-state-v3.json

```json
{
  "project_progression": {
    "current_wave": {
      "bugs": {
        "total_found": 15,
        "total_closed": 8,
        "total_open": 7,
        "reopened_count": 2,
        "history": [
          {
            "iteration": 1,
            "bugs_found": 5,
            "bugs_closed": 0,
            "bugs_open": 5,
            "bugs_reopened": 0,
            "bug_ids": ["W2-BUG-001", "W2-BUG-002", "W2-BUG-003", "W2-BUG-004", "W2-BUG-005"],
            "progress_category": "STALL",
            "timestamp": "2025-11-02T12:00:00Z"
          },
          {
            "iteration": 2,
            "bugs_found": 3,
            "bugs_closed": 2,
            "bugs_open": 6,
            "bugs_reopened": 0,
            "bug_ids": ["W2-BUG-001", "W2-BUG-003", "W2-BUG-006", "W2-BUG-007", "W2-BUG-008", "W2-BUG-009"],
            "progress_category": "DISCOVERY_PROGRESS",
            "timestamp": "2025-11-02T14:30:00Z"
          }
        ],
        "progress_stalls": 1,
        "max_stalls_allowed": 5,
        "max_iterations": 10,
        "current_iteration": 2
      }
    },
    "current_phase": {
      "bugs": { /* Same structure */ }
    }
  },
  "project_level": {
    "bugs": { /* Same structure */ }
  }
}
```

---

## 🔴🔴🔴 INTEGRATE_WAVE_EFFORTS WITH OTHER RULES 🔴🔴🔴

### R533 + R531 (Original Iteration Limits)

**R531 (OLD):** Simple counter-based limits
**R533 (NEW):** Progress-aware limits

**Relationship:**
- R533 SUPERSEDES R531 for bug-fix iterations
- R533 provides more nuanced decision-making
- R531 may still apply to other iteration types (not bug-fix)

### R533 + R534 (Bug Lifecycle Tracking)

**R534:** Defines how bugs are tracked, identified, and verified
**R533:** Uses R534 data to make iteration decisions

**Synergy:**
- R534 provides the bug data
- R533 analyzes the bug data for progress
- Together they prevent infinite loops while allowing legitimate progress

---

## 🔴🔴🔴 GRADING IMPACT 🔴🔴🔴

### Compliance Criteria

```yaml
r533_compliance:
  progress_tracking: 30%        # Track bugs_closed correctly
  two_tiered_limits: 30%        # Apply both 5 and 10 limits
  flapping_detection: 20%       # Detect and handle reopened bugs
  correct_escalation: 20%       # Escalate at right time

violations:
  ignoring_progress: -50%           # Using blind counters
  exceeding_no_progress_limit: -100%  # More than 5 stalls
  exceeding_some_progress_limit: -100%  # More than 10 iterations
  missing_flapping_detection: -30%    # Not detecting reopened bugs
  incorrect_progress_category: -20%   # Miscalculating progress
```

### Automatic Failure Conditions

```yaml
automatic_failure:
  - More than 5 iterations with NO bugs fixed
  - More than 10 total iterations (even with progress)
  - Flapping bugs not detected or handled
  - Progress calculation wrong (claiming progress when none)
```

---

## 🔴 SUMMARY: THE TWO-TIERED LIMIT PRINCIPLE

```
┌─────────────────────────────────────────────────────────────┐
│                TWO-TIERED ITERATION LIMITS                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  TIER 1: NO-PROGRESS LIMIT                                 │
│  ════════════════════════                                  │
│  Limit: 5 consecutive iterations with NO bugs fixed        │
│  Trigger: stall_counter >= 5                               │
│  Action: ERROR_RECOVERY (system is stuck)                  │
│                                                             │
│  TIER 2: SOME-PROGRESS LIMIT                               │
│  ══════════════════════════                                │
│  Limit: 10 total iterations (even with progress)           │
│  Trigger: iteration >= 10                                  │
│  Action: ERROR_RECOVERY (need replanning)                  │
│                                                             │
│  PROGRESS CATEGORIES:                                      │
│  ═══════════════════                                       │
│  ✅ PURE_PROGRESS      → Reset stall_counter               │
│  ✅ DISCOVERY_PROGRESS → Reset stall_counter               │
│  ⚠️  SLOW_PROGRESS      → Keep stall_counter               │
│  ❌ STALL              → Increment stall_counter           │
│  ❌ REGRESSION          → Increment stall_counter           │
│  🚨 FLAPPING (reopened) → Increment stall_counter          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Remember:** Progress is measured by BUGS FIXED, not bugs found. The system must close bugs to make progress. Finding new bugs while fixing others is acceptable (DISCOVERY_PROGRESS), but finding bugs without fixing any is STALL/REGRESSION.

**See Also:**
- R534: Bug Lifecycle Tracking Protocol (defines bug states and identification)
- R531: Iteration Container Limits (original, now superseded by R533 for bug-fix iterations)
- tools/bug-progress-analyzer.sh (implements progress calculation)
- bug-tracking.json schema (stores bug data)
- orchestrator-state-v3.json schema (stores iteration history)
