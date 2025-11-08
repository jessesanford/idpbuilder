# 🔴🔴🔴 RULE R536 - AUTOMATED FIX CASCADE PROTOCOL (SUPREME LAW)

**Criticality:** SUPREME LAW - ABSOLUTE AUTOMATION REQUIREMENT
**Grading Impact:** -100% for stopping automated fix cascades
**Enforcement:** Orchestrator, Integration Agent
**Scope:** ALL fix cascade operations (wave/phase/project levels)

---

## SUPREME LAW STATEMENT

**BUG FIXES AND FIX CASCADES ARE NORMAL OPERATION - NOT EXCEPTIONAL EVENTS. THE SYSTEM MUST AUTOMATICALLY LOOP THROUGH FIX CASCADES UNTIL CONVERGENCE (BUGS → 0) OR MAX ITERATIONS EXCEEDED. STOPPING FOR HUMAN APPROVAL OF BUG FIXES VIOLATES THE FUNDAMENTAL AUTOMATION PRINCIPLE.**

---

## 🔴🔴🔴 THE FUNDAMENTAL DESIGN PRINCIPLE 🔴🔴🔴

### Why This System Exists

The Software Factory is designed to **AUTOMATICALLY FIX BUGS FOUND DURING INTEGRATION**. This is not an exceptional workflow - it is THE CORE WORKFLOW:

```
NORMAL WORKFLOW (AUTOMATED):
Integration Review → Find Bugs → Create Fix Plan → Apply Fixes → Re-Integrate → Review Again → Repeat Until Clean
```

**This loop is DESIGNED INTO THE SYSTEM and should run AUTOMATICALLY.**

### What Was Wrong

The system was stopping at `CREATE_WAVE_FIX_PLAN` and asking for human approval. This is **COMPLETELY WRONG** because:

1. **Bug fixes are normal**: Finding bugs during integration is EXPECTED, not exceptional
2. **Fix cascades are designed**: The entire R300/R321 protocol exists FOR THIS PURPOSE
3. **Convergence is the goal**: The system should iterate until bugs → 0
4. **Human stops break flow**: Stopping disrupts the natural convergence process

### What Should Happen

**AUTOMATED FIX CASCADE LOOP:**

```
1. INTEGRATE_WAVE_EFFORTS → Review integration branch
2. Bugs found? → YES → CREATE_WAVE_FIX_PLAN (AUTOMATED)
3. Analyze bugs → Create fix assignments
4. FIX_WAVE_UPSTREAM_BUGS (AUTOMATED)
   - Apply fixes to effort branches
   - Push fixes to remotes
5. START_WAVE_ITERATION (AUTOMATED - RE-INTEGRATE)
   - Create fresh integration
   - Include fixed effort branches
6. INTEGRATE_WAVE_EFFORTS → Review integration branch (ITERATION N+1)
7. Bugs found?
   - YES + bugs decreasing → LOOP (convergence in progress)
   - YES + bugs increasing → ERROR_RECOVERY (divergence detected)
   - YES + iteration > MAX → ERROR_RECOVERY (iteration overflow)
   - NO → WAVE_COMPLETE (convergence achieved)
```

**KEY INSIGHT**: This entire loop should execute **WITHOUT HUMAN INTERVENTION**.

---

## 🚨🚨🚨 WHEN TO STOP (ONLY THESE SCENARIOS) 🚨🚨🚨

The system should **ONLY** stop and request human intervention for:

### STOP SCENARIO 1: Iteration Overflow

```yaml
condition:
  - Same bugs recurring after MAX_ITERATIONS (default: 10)
  - Bug count not decreasing over multiple iterations
  - Fix cascade not converging

action: ERROR_RECOVERY
reason: "Non-convergence detected - manual investigation required"
flag: CONTINUE-SOFTWARE-FACTORY=FALSE
```

**Example:**
```
Iteration 1: 5 bugs found → fixes applied
Iteration 2: 5 bugs found (SAME BUGS) → fixes applied
Iteration 3: 5 bugs found (STILL SAME) → fixes applied
...
Iteration 10: 5 bugs found → ERROR_RECOVERY (not converging)
```

### STOP SCENARIO 2: Divergence Detection

```yaml
condition:
  - Bug count INCREASING over iterations
  - Fixes introducing more bugs than they solve
  - System becoming less stable

action: ERROR_RECOVERY
reason: "Divergence detected - fixes making things worse"
flag: CONTINUE-SOFTWARE-FACTORY=FALSE
```

**Example:**
```
Iteration 1: 3 bugs found → fixes applied
Iteration 2: 5 bugs found (MORE BUGS!) → fixes applied
Iteration 3: 8 bugs found (DIVERGING!) → ERROR_RECOVERY
```

### STOP SCENARIO 3: Catastrophic System Failure

```yaml
condition:
  - State file corruption
  - Git repository corruption
  - Missing critical files/branches
  - System integrity compromised

action: ERROR_RECOVERY
reason: "System corruption - manual recovery required"
flag: CONTINUE-SOFTWARE-FACTORY=FALSE
```

### STOP SCENARIO 4: Recursive Split Required

```yaml
condition:
  - Fix causes effort to exceed 900 lines
  - Split-of-a-split detected (R515 violation)
  - Size limit cannot be satisfied with normal fixes

action: ERROR_RECOVERY
reason: "Recursive split detected - architectural review required"
flag: CONTINUE-SOFTWARE-FACTORY=FALSE
```

---

## ✅ NORMAL OPERATIONS (NEVER STOP FOR THESE)

The system must **CONTINUE AUTOMATICALLY** through:

### ✅ NORMAL 1: Bugs Found During Integration

```bash
# Integration review finds bugs → CREATE FIX PLAN automatically
if [ "$bugs_found" -gt 0 ]; then
    echo "✅ NORMAL: Found $bugs_found bugs during integration review"
    echo "✅ AUTOMATED ACTION: Creating fix plan..."
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
    transition_to "CREATE_WAVE_FIX_PLAN"
fi
```

**NO STOP. NO HUMAN APPROVAL. AUTOMATIC CONTINUATION.**

### ✅ NORMAL 2: First Fix Cascade

```bash
# First time fixing bugs in this wave
echo "✅ NORMAL: First fix cascade for Wave 1"
echo "✅ AUTOMATED ACTION: Applying fixes to effort branches..."
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
transition_to "FIX_WAVE_UPSTREAM_BUGS"
```

**This is normal development. No stop required.**

### ✅ NORMAL 3: Second, Third, Fourth Fix Cascades

```bash
# Multiple fix iterations (convergence in progress)
echo "✅ NORMAL: Fix cascade iteration $ITERATION_COUNT"
echo "✅ CONVERGENCE: Bug count: 15 → 8 → 3"
echo "✅ AUTOMATED ACTION: Continuing fix cascade..."
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**Convergence is the GOAL. Continue until bugs → 0.**

### ✅ NORMAL 4: Integration Failures

```bash
# Integration attempt fails (merge conflicts, build failures)
echo "✅ NORMAL: Integration build failed"
echo "✅ AUTOMATED ACTION: Creating fix plan for build failures..."
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
transition_to "CREATE_WAVE_FIX_PLAN"
```

**Build failures ARE bugs. Fix them automatically.**

### ✅ NORMAL 5: Test Failures

```bash
# Tests fail during integration
echo "✅ NORMAL: 5 tests failing on integration branch"
echo "✅ AUTOMATED ACTION: Creating fix plan for test failures..."
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**Test failures ARE bugs. Fix them automatically.**

### ✅ NORMAL 6: Code Quality Issues

```bash
# Code review finds quality issues
echo "✅ NORMAL: Code quality issues detected"
echo "✅ AUTOMATED ACTION: Creating fix plan..."
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**Quality issues ARE bugs. Fix them automatically.**

### ✅ NORMAL 7: Size Violations

```bash
# Effort exceeds size limit after fixes
if [ "$line_count" -gt 900 ]; then
    echo "✅ NORMAL: Size violation detected ($line_count > 900)"
    echo "✅ AUTOMATED ACTION: Creating split plan..."
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
    transition_to "CREATE_SPLIT_PLAN"
fi
```

**Splits are normal operations. Execute automatically.**

---

## 🔴🔴🔴 ITERATION TRACKING AND CONVERGENCE DETECTION 🔴🔴🔴

### Tracking Iteration State

The system MUST track fix cascade iterations in `integration-containers.json`:

```json
{
  "current_iteration_container": {
    "level": "wave",
    "identifier": "phase1/wave1",
    "iteration_count": 3,
    "max_iterations": 10,
    "bug_history": [
      {"iteration": 1, "bugs_found": 15, "bugs_fixed": 0},
      {"iteration": 2, "bugs_found": 8, "bugs_fixed": 7},
      {"iteration": 3, "bugs_found": 3, "bugs_fixed": 5}
    ],
    "convergence_status": "CONVERGING"
  }
}
```

### Convergence Detection Logic

```bash
detect_convergence_status() {
    local current_bugs="$1"
    local previous_bugs="$2"
    local iteration="$3"
    local max_iterations=10

    # Check for convergence (bugs decreasing)
    if [ "$current_bugs" -eq 0 ]; then
        echo "CONVERGED"
        return 0
    elif [ "$current_bugs" -lt "$previous_bugs" ]; then
        echo "CONVERGING"
        return 0
    # Check for divergence (bugs increasing)
    elif [ "$current_bugs" -gt "$previous_bugs" ]; then
        echo "DIVERGING"
        return 1
    # Check for stagnation (same bugs)
    elif [ "$current_bugs" -eq "$previous_bugs" ]; then
        if [ "$iteration" -gt "$max_iterations" ]; then
            echo "ITERATION_OVERFLOW"
            return 1
        else
            echo "STAGNANT"
            return 0
        fi
    fi
}
```

### Decision Matrix

```bash
make_continuation_decision() {
    local status="$1"
    local iteration="$2"

    case "$status" in
        CONVERGED)
            echo "✅ CONVERGENCE ACHIEVED: bugs_found = 0"
            echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
            transition_to "WAVE_COMPLETE"
            ;;
        CONVERGING)
            echo "✅ CONVERGENCE IN PROGRESS: bugs decreasing"
            echo "✅ Iteration $iteration: Normal fix cascade"
            echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
            transition_to "CREATE_WAVE_FIX_PLAN"
            ;;
        STAGNANT)
            echo "⚠️ STAGNANT: Same bugs recurring"
            echo "✅ Iteration $iteration < MAX: Continue cascade"
            echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
            transition_to "CREATE_WAVE_FIX_PLAN"
            ;;
        DIVERGING)
            echo "❌ DIVERGENCE DETECTED: bugs increasing"
            echo "❌ Fixes making things worse"
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=DIVERGENCE"
            transition_to "ERROR_RECOVERY"
            ;;
        ITERATION_OVERFLOW)
            echo "❌ ITERATION OVERFLOW: $iteration > MAX_ITERATIONS"
            echo "❌ Same bugs after 10+ iterations"
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=NON_CONVERGENCE"
            transition_to "ERROR_RECOVERY"
            ;;
    esac
}
```

---

## 📊 IMPLEMENTATION CHECKLIST

### State Machine Updates (COMPLETED)

- [x] `CREATE_WAVE_FIX_PLAN`: Set `checkpoint: false`
- [x] `CREATE_PHASE_FIX_PLAN`: Set `checkpoint: false`
- [x] `CREATE_PROJECT_FIX_PLAN`: Set `checkpoint: false`
- [x] Update actions: Remove "Present plan to user for approval"
- [x] Update actions: Add "Execute fix plan automatically (R536)"

### R322 Updates (COMPLETED)

- [x] Confirm fix cascade states NOT in checkpoint list
- [x] Confirm "FIX PLANS ARE NORMAL" language present
- [x] Confirm automatic continuation guidance clear

### R405 Updates (REQUIRED)

- [ ] Add explicit examples for bug fix scenarios
- [ ] Clarify TRUE for bugs found during integration
- [ ] Clarify TRUE for fix plan creation
- [ ] Clarify TRUE for fix cascade iterations
- [ ] Clarify FALSE ONLY for overflow/divergence

### Orchestrator Config Updates (REQUIRED)

- [ ] Emphasize automated fix cascade loop
- [ ] Clarify when to stop (only overflow/divergence)
- [ ] Remove any suggestion of stopping for normal bugs
- [ ] Add iteration tracking requirements

### Code Reviewer Config Updates (REQUIRED)

- [ ] Clarify finding bugs is NORMAL, not exceptional
- [ ] Remove any suggestion to wait for approval
- [ ] Emphasize creating fix plans, not stopping
- [ ] Add bug tracking responsibilities

---

## 🔴 GRADING IMPACT

### Violations Result In:

- **-100% IMMEDIATE FAILURE** for stopping automated fix cascades unnecessarily
- **-75% PENALTY** for requesting human approval of normal bug fixes
- **-50% PENALTY** for treating bug fixes as exceptional events
- **-25% PENALTY** for missing iteration tracking

### Compliance Results In:

- **+20% BONUS** for perfect automated convergence
- **+10% BONUS** for proper iteration tracking
- **Clean fix cascade loop** executing dozens of iterations if needed
- **User confidence** that system handles bugs automatically

---

## 🎯 RELATIONSHIP TO OTHER RULES

### Works With:

- **R300**: Comprehensive Fix Management Protocol - Foundation for fix cascades
- **R321**: Immediate Backport During Integration - Fix at source immediately
- **R322**: Mandatory Checkpoints - Explicitly EXCLUDES fix cascades
- **R405**: Automation Continuation Flag - Controls automated continuation
- **R515**: Prohibition on Recursive Splits - Detection point for ERROR_RECOVERY
- **State Machine**: Convergence loop through fix cascade states

### Supersedes:

- Any guidance suggesting human approval of bug fixes
- Any pattern of stopping at fix plan creation
- Previous "checkpoint before fix execution" patterns
- Manual fix approval workflows

---

## 📋 EXAMPLES

### ✅ CORRECT: Automated Fix Cascade (Normal Operation)

```bash
# Wave integration review complete
echo "📊 Integration Review Results:"
echo "  - Bugs found: 12"
echo "  - Bug types: 8 logic errors, 4 integration issues"
echo "  - All bugs recorded in bug-tracking.json"

# Check convergence status
ITERATION_COUNT=$(jq '.current_iteration_container.iteration_count' integration-containers.json)
echo "📊 Iteration: $ITERATION_COUNT/10"

# AUTOMATED DECISION: Create fix plan
echo "✅ NORMAL OPERATION: Bugs found during integration"
echo "✅ AUTOMATED ACTION: Creating fix plan for 12 bugs..."
echo "✅ No human approval required - this is normal development"

# Update state
update_state "CREATE_WAVE_FIX_PLAN"
commit_state()

# R405 flag: TRUE because this is normal operation
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# Exit for State Manager consultation
exit 0
```

**What happens next:**
1. System automatically runs `/continue-software-factory`
2. Orchestrator enters `CREATE_WAVE_FIX_PLAN` state
3. Creates fix assignments for 12 bugs
4. Transitions to `FIX_WAVE_UPSTREAM_BUGS` automatically
5. Applies fixes to upstream effort branches
6. Transitions to `START_WAVE_ITERATION` automatically
7. Creates fresh integration including fixed branches
8. Runs integration review again (Iteration 2)
9. **LOOPS until bugs → 0 or max iterations**

**NO HUMAN INTERVENTION NEEDED!**

### ✅ CORRECT: Convergence Achieved

```bash
# Iteration 5 review complete
echo "📊 Integration Review Results (Iteration 5):"
echo "  - Bugs found: 0"
echo "  - Previous iterations: 15 → 8 → 4 → 1 → 0"
echo "  - Convergence: ACHIEVED"

# AUTOMATED DECISION: Wave complete
echo "✅ CONVERGENCE ACHIEVED: Zero bugs remaining"
echo "✅ AUTOMATED ACTION: Marking wave as complete..."
echo "✅ System iterated 5 times to achieve clean integration"

# Update state
update_state "WAVE_COMPLETE"
commit_state()

# R405 flag: TRUE because convergence succeeded
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

exit 0
```

**Result**: System automatically proceeds to next wave or phase integration.

### ❌ CORRECT: Iteration Overflow Detection (Stop Required)

```bash
# Iteration 11 review complete
echo "📊 Integration Review Results (Iteration 11):"
echo "  - Bugs found: 3"
echo "  - Bug history: 3 → 3 → 3 → 3 → 3 (5 iterations)"
echo "  - Convergence: FAILED (non-convergence detected)"
echo "  - Max iterations: 10 (EXCEEDED)"

# EXCEPTIONAL CONDITION: Iteration overflow
echo "❌ ITERATION OVERFLOW DETECTED"
echo "❌ Same 3 bugs recurring after 11 iterations"
echo "❌ System not converging - manual investigation required"

# Display bug details
cat << 'EOF'
## 🔴 FIX CASCADE NON-CONVERGENCE - MANUAL INTERVENTION REQUIRED

### ❌ Iteration Overflow:
- Iterations: 11/10 (exceeded maximum)
- Bug count stable at 3 bugs
- Fixes not resolving root cause

### 🐛 Recurring Bugs:
1. BUG-001: Database connection timeout (5 iterations)
2. BUG-002: Race condition in async handler (11 iterations)
3. BUG-003: Memory leak in cache (8 iterations)

### 🛠️ HUMAN ACTION REQUIRED:
1. Review bug patterns for systemic issues
2. Consider architectural changes
3. Investigate why fixes aren't working
4. May require phase/project redesign

### ⏸️ AUTOMATION STOPPED
EOF

# Update state
update_state "ERROR_RECOVERY"
commit_state()

# R405 flag: FALSE because human investigation required
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=ITERATION_OVERFLOW"

exit 1
```

**Result**: Automation stops. Human investigates why fixes aren't converging.

### ❌ WRONG: Stopping at First Bug (VIOLATION!)

```bash
# Integration review finds bugs
echo "📊 Integration Review Results:"
echo "  - Bugs found: 5"

# ❌ VIOLATION: Stopping for human approval
echo "🛑 R322: Stopping for user to review bug fixes"
echo "🛑 Please approve fix plan before proceeding"

# ❌ VIOLATION: Using FALSE for normal operation
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=BUGS_FOUND"

exit 0
```

**Why this is wrong:**
- Bugs found is NORMAL, not exceptional
- No approval needed for normal bug fixes
- Should use TRUE and continue automatically
- **Penalty**: -100% for violating R536

---

## 🔴 THE FIX CASCADE MANTRA

**Repeat before every integration review:**

1. **"Bugs are normal - fixes are automatic"**
2. **"Continue until convergence (bugs → 0)"**
3. **"Stop ONLY for overflow or divergence"**
4. **"Flag TRUE for normal operations"**
5. **"Iteration tracking prevents infinite loops"**

---

**Remember:** The entire point of the Software Factory is to automatically find and fix bugs through iterative integration. Stopping this process defeats the core automation principle. Let the system converge naturally through automated fix cascades.

**See Also:**
- R300: Comprehensive Fix Management Protocol (foundation)
- R321: Immediate Backport During Integration (fix at source)
- R322: Mandatory Checkpoints (excludes fix cascades)
- R405: Automation Continuation Flag (controls continuation)
- R515: Prohibition on Recursive Splits (error detection)
- State Machine: `integration-containers.json` iteration tracking
