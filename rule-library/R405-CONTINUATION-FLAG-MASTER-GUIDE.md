# 🔴🔴🔴 CONTINUATION FLAG MASTER GUIDE - READ THIS FIRST 🔴🔴🔴

## 🚨 CRITICAL FORMAT REQUIREMENT 🚨

**YOU MUST OUTPUT THE COMPLETE FLAG, NOT JUST THE VALUE!**

### ❌ WRONG (WILL BREAK AUTOMATION):
```bash
TRUE
FALSE
R405 Continuation Flag: TRUE
R405 Continuation Flag: FALSE
```

### ✅ CORRECT (REQUIRED FORMAT):
```bash
CONTINUE-SOFTWARE-FACTORY=TRUE
CONTINUE-SOFTWARE-FACTORY=FALSE
```

**The automation framework greps for `^CONTINUE-SOFTWARE-FACTORY=` - anything else WILL FAIL!**

---

## THE GOLDEN RULE

**Default to TRUE. ALWAYS. Unless something is GENUINELY BROKEN AND UNRECOVERABLE.**

## What TRUE Means

TRUE = "System can auto-restart this conversation"
- Normal operations continuing ✅
- Designed workflows proceeding ✅
- Recoverable issues being handled ✅
- Success outcomes ✅
- Failures with known recovery paths ✅

## What FALSE Means

FALSE = "System CANNOT auto-restart - MANUAL INTERVENTION REQUIRED"
- State corruption with no recovery path ❌
- Files missing/corrupt with no way to recreate ❌
- Infrastructure completely destroyed ❌
- State machine invalid ❌
- Unknown error with no handling ❌

## Common Scenarios - DEFINITIVE ANSWERS

### ✅ Use TRUE (99.9% of cases)

**Success scenarios:**
- Tests passing → TRUE
- Build successful → TRUE
- Integration complete → TRUE
- Review approved → TRUE
- Deployment successful → TRUE

**Normal workflow scenarios:**
- Spawning agents → TRUE
- Waiting for results → TRUE
- Monitoring progress → TRUE
- State transitions → TRUE
- Checkpoints → TRUE

**Recoverable issue scenarios:**
- Tests failing (can fix) → TRUE
- Review found issues (can fix) → TRUE
- Build errors (can debug) → TRUE
- Integration conflicts (can resolve) → TRUE
- Size violations (can split) → TRUE

**All of these have recovery paths in the system!**

### 🔴 Use FALSE (0.1% of cases - TRULY EXCEPTIONAL)

**System corruption scenarios:**
- State file corrupted beyond parsing → FALSE
- Entire infrastructure directory missing → FALSE
- Git repository corrupted → FALSE
- Required files deleted with no backup → FALSE

**Unknown/Unhandled scenarios:**
- Error with no defined recovery protocol → FALSE
- State machine in invalid/undefined state → FALSE
- Critical dependency completely unavailable → FALSE

**Ask yourself: "Can the system recover from this automatically with its existing protocols?"**
- YES → TRUE
- NO (truly stuck, needs human) → FALSE

## 🔴🔴🔴 R322 INTERACTION - CRITICAL UNDERSTANDING 🔴🔴🔴

### The #1 Confusion in Software Factory

**WRONG thinking:**
> "R322 says mandatory stop → must use FALSE"
> "State work is complete and waiting for /continue-orchestrating → must use FALSE"
> "This is a checkpoint before state transition → must use FALSE"

**CORRECT thinking:**
> "R322 says stop conversation (`exit 0`) for context. Flag is independent - use TRUE for normal ops."
> "Completing state work and waiting for /continue is NORMAL OPERATION → use TRUE"
> "R322 checkpoints are EXPECTED WORKFLOW, not error conditions → use TRUE"

**R322 + TRUE = Normal operation** (99.9% of checkpoints)
**R322 + FALSE = Unrecoverable error** (0.1% of checkpoints)

### 🚨 CRITICAL: R322 Checkpoints Are NORMAL OPERATION 🚨

**R322 requires stops at specific points to:**
1. Preserve context between states (prevents overflow)
2. Create clean transition boundaries
3. Allow state file commits before continuation

**THESE ARE ALL NORMAL, DESIGNED BEHAVIORS!**

When you complete a state and stop per R322:
- ✅ State work completed → NORMAL
- ✅ Waiting for /continue-orchestrating → NORMAL (designed UX!)
- ✅ System ready for next state → NORMAL
- ✅ Use TRUE because automation knows what to do next!

**ONLY use FALSE if:**
- ❌ System is corrupted and cannot continue at all
- ❌ Data integrity violated beyond recovery
- ❌ Critical files missing with no way to recreate
- ❌ Once-per-10-years catastrophic situations

### 🚨 SPECIAL CASE: WAITING STATES 🚨

**WAITING_FOR_* states are FREQUENT SOURCES of incorrect FALSE usage!**

#### Example: WAITING_FOR_MERGE_PLAN

**Scenario:** Code Reviewer created merge plan, it exists and is validated

**WRONG interpretation:**
> "R322 mandates stop before SPAWN_INTEGRATION_AGENT transition"
> "State work is complete (plan validated)"
> "User needs to invoke /continue-orchestrating to proceed"
> "Therefore I must set CONTINUE-SOFTWARE-FACTORY=FALSE"

**CORRECT interpretation:**
> "R322 checkpoint is NORMAL procedure for context preservation"
> "State work completed successfully = NORMAL outcome"
> "Waiting for /continue-orchestrating is DESIGNED user experience"
> "System KNOWS next step: transition to SPAWN_INTEGRATION_AGENT"
> "NO manual intervention required, just normal continuation"
> "Therefore set CONTINUE-SOFTWARE-FACTORY=TRUE"

**The key distinction:**
- **Stopping inference** (`exit 0`) = Context management (ALWAYS at R322 points)
- **Continuation flag** = Can automation proceed? (TRUE unless catastrophic failure)

**These are INDEPENDENT decisions!**

You STOP for R322 checkpoint AND set flag=TRUE for normal continuation.

#### All WAITING States Follow This Pattern:

- `WAITING_FOR_ASSESSMENT` → Assessment complete → TRUE (continue to next state)
- `WAITING_FOR_REVIEW` → Review complete → TRUE (continue based on result)
- `WAITING_FOR_EFFORT_PLANS` → Plans ready → TRUE (continue to parallelization)
- `WAITING_FOR_MERGE_PLAN` → Plan validated → TRUE (continue to spawn integration)
- `WAITING_FOR_IMPLEMENTATION` → Work done → TRUE (continue to reviews)
- `WAITING_FOR_REVIEW_RESULTS` → Reviews complete → TRUE (continue based on results)

**PATTERN:** WAITING states check for completion, validate results, then continue → ALL use TRUE!

**ONLY use FALSE in WAITING states if:**
- ❌ The thing we're waiting for completely disappeared (agents crashed with no recovery)
- ❌ Results arrived but are completely corrupted/unreadable
- ❌ State file corruption prevents determining what to wait for
- ❌ System deadlock with no automated resolution

### 🚨 SPECIAL CASE: CASCADE OPERATIONS (R327) 🚨

**CASCADE transitions are THE #1 SOURCE of incorrect FALSE usage!**

**CASCADE_REINTEGRATION → SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE:**
```bash
# After deleting stale integration (correct)
# Transitioning to SETUP to recreate (correct)
# R322 checkpoint (correct)
# Flag? → MUST BE TRUE!

# Why? The system knows:
# - Current state: SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE (from state file)
# - What SETUP does: Create infrastructure per R504
# - Where SETUP goes: Back to CASCADE or to merge
# - NO HUMAN DECISION NEEDED!

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # CASCADE CONTINUES!
# ❌ NOT FALSE! Defeats cascade automation!
```

**CASCADE_REINTEGRATION → CASCADE_REINTEGRATION (loop):**
```bash
# After SETUP creates infrastructure (correct)
# Returning to CASCADE for next deletion (correct)
# R322 checkpoint (correct)
# Flag? → MUST BE TRUE!

# Why? The system knows:
# - Current state: CASCADE_REINTEGRATION (from state file)
# - Remaining work: Delete next stale integration
# - NO HUMAN DECISION NEEDED!

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # CASCADE CHAIN CONTINUES!
```

**See CASCADE_REINTEGRATION/rules.md for complete decision flowchart.**

### The Distinction

**R322 Stop (Conversation Checkpoint):**
- **Purpose:** Preserve context, prevent memory loss
- **When:** Before state transitions, after significant work
- **How:** `exit 0` to end conversation
- **Result:** User runs `/continue` command to resume

**CONTINUE-SOFTWARE-FACTORY Flag:**
- **Purpose:** Can the system AUTO-RESUME without manual intervention?
- **When:** As last output before exit
- **How:** Echo the flag value
- **Result:** Determines if automation can pick up where it left off

### They Are Independent!

You can (and usually should) do BOTH:
1. Stop per R322 (conversation checkpoint)
2. Set TRUE (automation can continue)

### Examples

#### ✅ CORRECT: Sequential Fix Work
```bash
# In ERROR_RECOVERY, fixing bugs one by one
fix_bug_1()
commit_per_r288()

# R322: Stop after significant work
echo "🛑 R322: Checkpoint after Bug 1 fixed"
exit 0

# Automation can continue with Bug 2
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

#### ✅ CORRECT: State Transition
```bash
# Completed ERROR_RECOVERY, transitioning to CASCADE_REINTEGRATION
update_state("CASCADE_REINTEGRATION")
commit_per_r288()

# R322: Stop before state transition
echo "🛑 R322: Checkpoint before state transition"
exit 0

# System can auto-transition
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

#### ❌ WRONG: Unnecessary FALSE
```bash
# Fixed a bug successfully
fix_bug()
commit_state()

echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ❌ Why? Everything is fine!
```

### Decision Tree

```
Are you stopping per R322? ──┐
                            YES → Is automation blocked?
                                  ├─ NO → TRUE (99% of cases)
                                  └─ YES → Is it unrecoverable?
                                           ├─ NO → TRUE (system can retry)
                                           └─ YES → FALSE (needs human)
```

### Grading Penalty

Incorrectly using FALSE when automation can continue: **-20%**
This defeats the purpose of Software Factory automation!

## Decision Tree

```
Is something broken?
├─ NO → TRUE (normal operation)
└─ YES → Can system recover automatically?
         ├─ YES → TRUE (use recovery protocols)
         └─ NO (truly stuck) → FALSE
```

## Examples from Real Scenarios

### ❌ WRONG - Past Violations

```bash
# VIOLATION 1: Success treated as needing review
Integration_Status="PROJECT_DONE"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # WRONG! Success is normal!

# VIOLATION 2: Code review finding issues
Issues_Found=true
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # WRONG! This is normal workflow!

# VIOLATION 3: Spawning agents
Spawning_Engineer=true
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # WRONG! This is normal operation!

# VIOLATION 4: Integration complete
Integration_Complete=true
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # WRONG! This is success!
```

### ✅ CORRECT - Proper Usage

```bash
# CORRECT 1: Success continues automation
Integration_Status="PROJECT_DONE"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # RIGHT! Continue automation!

# CORRECT 2: Recoverable issues continue
Issues_Found=true  # Will enter review-fix cycle (designed process)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # RIGHT! System handles this!

# CORRECT 3: Spawning continues
Spawning_Engineer=true  # Normal workflow
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # RIGHT! Let it continue!

# CORRECT 4: Only true errors stop
State_File_Corrupt=true  # Cannot parse, cannot recover
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # RIGHT! Need human intervention!
```

## Grading Impact

- **-20%** per incorrect FALSE for normal operations
- **-50%** for pattern of excessive FALSE usage (3+)
- **-100%** for complete automation defeat
- **+0%** for correct usage (no penalty, as expected)

## References

- R405: Automation Flag Continuation Principle (base rule)
- R322: Mandatory Stop Before State Transitions (context preservation)
- R287: TODO Persistence (state management)
- R288: State File Updates (state management)

## When in Doubt

**If you're not sure whether to use TRUE or FALSE:**

1. Ask: "Is the system completely stuck with no way forward?"
   - NO → Use TRUE
   - YES → Maybe FALSE, but double-check recovery protocols first

2. Check: Does a recovery protocol exist for this situation?
   - YES → Use TRUE (let protocol handle it)
   - NO → Maybe FALSE, but check state machine for error handling first

3. Verify: Is this truly exceptional or just normal workflow?
   - Normal workflow → Use TRUE
   - Truly exceptional → Maybe FALSE

**Default assumption: Use TRUE unless you have a VERY good reason for FALSE.**

The Software Factory is designed to be resilient and self-recovering.
Give it a chance to work before stopping automation!

## Code Review Scenarios - DEFINITIVE ANSWERS

**CRITICAL: This scenario causes REPEATED violations. Read carefully!**

### ✅ Review approved → TRUE
```bash
REVIEW_STATUS="APPROVED"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Success, proceed
```

### ✅ Review finds issues → TRUE
```bash
REVIEW_STATUS="NEEDS_FIXES"
# CRITICAL: This is NORMAL and EXPECTED!
# System has automatic fix protocol (review → fix → re-review)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Let fix protocol work!
```

### ✅ Review requires split → TRUE
```bash
REVIEW_STATUS="NEEDS_SPLIT"
# System has automatic split protocol
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Let split protocol work!
```

### 🔴 Review report corrupt → FALSE
```bash
if [ ! -f "REVIEW-REPORT.md" ] || [ ! -s "REVIEW-REPORT.md" ]; then
    # Truly broken - cannot determine review outcome
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # Need human intervention
fi
```

**KEY INSIGHT:** The entire PURPOSE of code review is to find issues so they
can be fixed. Finding issues is PROJECT_DONE, not failure!

**USER'S WORDS:** "WE HAVE A FIX PROTOCOL THAT IS AUTOMATIC FOR THIS REASON!"

## Systemic Violations Being Fixed

**ALERT: The following violations have been identified and corrected:**

1. **SPAWN_SW_ENGINEERS** - Fixed: Setting FALSE when spawning is normal workflow
2. **MONITORING_EFFORT_REVIEWS** - Fixed: Setting FALSE when review complete is normal
3. **SPAWN_CODE_REVIEWERS_EFFORT_REVIEW** - Fixed: Setting FALSE when spawning reviewers is normal
4. **MONITORING_PROJECT_INTEGRATE_WAVE_EFFORTS** - Fixed: Setting FALSE after PROJECT_DONE is completely wrong
5. **PROJECT_REVIEW_WAVE_INTEGRATION** - Fixed: Setting FALSE when spawning reviewer is normal
6. **WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION** - Fixed: Setting FALSE when review finds issues

**Root Cause:** Misunderstanding that:
- R322 "stop" does NOT mean use FALSE
- PROJECT_DONE outcomes should ALWAYS use TRUE
- Recoverable issues should use TRUE (they have fix protocols)
- Normal state transitions should use TRUE
- **Review finding issues is NORMAL and EXPECTED (not an error!)**

**This guide exists to prevent future violations. Read it before setting the flag in ANY state.**
