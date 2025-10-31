# SPAWN_SW_ENGINEERS State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: MONITORING_EFFORT_REVIEWS, SPAWN_CODE_REVIEWER_FIX_PLAN
**Exit To**: MONITORING_EFFORT_FIXES

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**SPAWN SW Engineer agents to implement fixes identified during code review.**

This is a **STANDARD WORKFLOW STATE** - code review finding issues and requiring fixes is **NORMAL SOFTWARE DEVELOPMENT PRACTICE**. The review → fix → re-review cycle is designed to handle this automatically.

## Required Inputs

### 1. List of Efforts Needing Fixes
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Get efforts that need fixes
EFFORTS_NEEDING_FIXES=$(jq -r '.efforts_needing_fixes[]?' orchestrator-state-v3.json 2>/dev/null)

if [ -z "$EFFORTS_NEEDING_FIXES" ]; then
    echo "❌ ERROR: No efforts marked as needing fixes!"
    echo "State corruption - should not be in SPAWN_SW_ENGINEERS without efforts list"
    exit 1
fi

echo "🔧 Spawning engineers to fix:"
echo "$EFFORTS_NEEDING_FIXES"
```

### 2. Verify Fix Instructions Exist
```bash
for effort in $EFFORTS_NEEDING_FIXES; do
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"

    # Find the latest fix instructions
    FIX_INSTRUCTIONS=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-INSTRUCTIONS--"*.md 2>/dev/null | head -1)

    if [ -z "$FIX_INSTRUCTIONS" ] || [ ! -f "$FIX_INSTRUCTIONS" ]; then
        echo "❌ CRITICAL ERROR: Fix instructions missing for $effort"
        echo "Expected at: $EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-INSTRUCTIONS--*.md"
        exit 1
    fi

    echo "✅ Fix instructions found for $effort: $FIX_INSTRUCTIONS"
done
```

## 🔴🔴🔴 ENGINEER SPAWNING PROTOCOL 🔴🔴🔴

### Sequential Fix Processing (NOT Parallel)

**CRITICAL**: Fixes must be processed **ONE AT A TIME** to ensure proper review cycles and avoid conflicts.

```bash
echo "🚀 SPAWNING ENGINEERS FOR FIXES - $(date -Iseconds)"
echo "═══════════════════════════════════════════════════════"

# Process efforts sequentially
for effort in $EFFORTS_NEEDING_FIXES; do
    echo ""
    echo "🔧 Spawning engineer for: $effort"

    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    FIX_INSTRUCTIONS=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-INSTRUCTIONS--"*.md 2>/dev/null | head -1)

    # Spawn engineer with fix instructions
    echo "📝 Fix instructions: $FIX_INSTRUCTIONS"
    echo "📂 Working directory: $EFFORT_DIR"

    # Engineer will:
    # 1. Read fix instructions
    # 2. Implement fixes
    # 3. Run tests
    # 4. Commit changes
    # 5. Mark as REQUEST_REVIEW

    echo "Agent: @agent-sw-engineer"
    echo "Task: Implement fixes for $effort as specified in fix instructions"
    echo "Context: Fix instructions path: $FIX_INSTRUCTIONS"
done

echo ""
echo "═══════════════════════════════════════════════════════"
```

## State Transition

### Update State After Spawning
```bash
# Clear efforts_needing_fixes and transition to monitoring
jq --arg timestamp "$(date -Iseconds)" \
   --argjson fixes "$(echo $EFFORTS_NEEDING_FIXES | jq -R -s -c 'split(" ") | map(select(length > 0))')" \
   '.efforts_needing_fixes = null |
    .fixes_in_progress = $fixes |
    .state_machine.current_state = "MONITORING_EFFORT_FIXES" |
    .state_machine.previous_state = "SPAWN_SW_ENGINEERS" |
    .state_transition_log += [{
        "from": "SPAWN_SW_ENGINEERS",
        "to": "MONITORING_EFFORT_FIXES",
        "timestamp": $timestamp,
        "reason": "Engineers spawned to implement fixes"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ State updated: MONITORING_EFFORT_FIXES"
```

## Integration with Rules

- **R313/R322**: Mandatory stop after spawning (context preservation)
- **R405**: Continuation flag (see below for critical guidance)
- **R287**: TODO Persistence before transition
- **Grading**: Workflow Compliance (25%), Sequential Fix Processing

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Fix instructions exist and are readable
- ✅ Effort directory exists and is accessible
- ✅ Engineer spawned successfully
- ✅ All normal operations completed
- ✅ State transition successful

**THIS IS THE STANDARD CASE.** Code review finding issues and requiring fixes is **NORMAL SOFTWARE DEVELOPMENT WORKFLOW**. The fix → re-review cycle is designed to handle this automatically.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Fix instructions file is missing or corrupt
- ❌ Cannot determine which effort needs fixes (state corruption)
- ❌ Effort infrastructure is broken/missing
- ❌ State machine corruption detected
- ❌ Unrecoverable error prevents spawning engineer
- ❌ Critical system violation detected

### ❌ DO NOT SET FALSE JUST BECAUSE:

**These are NOT reasons to use FALSE:**
- ❌ Code review found issues (this is **NORMAL**)
- ❌ Fixes are needed (this is **EXPECTED**)
- ❌ Transitioning to fix state (this is **STANDARD WORKFLOW**)
- ❌ "User might want to review" (only if truly exceptional)
- ❌ Multiple efforts need fixes (this is **DESIGNED FOR**)
- ❌ This is the second or third fix cycle (this is **EXPECTED**)

## R322 Clarification for This State

**R322 requires:**
1. ✅ Save state before transition (MANDATORY)
2. ✅ Emit continuation flag (MANDATORY)
3. ✅ Stop inference - context preservation (MANDATORY)

**R322 does NOT require:**
- ❌ Setting FALSE for normal operations
- ❌ Human review of standard workflow
- ❌ Stopping automation for expected processes

**The "stop" in R322 means:** Stop THIS conversation turn to preserve context.

**The flag determines:** Whether the system can auto-restart (TRUE) or requires human intervention (FALSE).

**For SPAWN_SW_ENGINEERS:** Almost always TRUE because spawning engineers for fixes is normal workflow!

## Grading Impact

**Correct flag usage:**
- ✅ Using TRUE for normal fix workflow: No penalty
- ✅ Using FALSE for genuine errors: No penalty

**Incorrect flag usage:**
- ❌ Using FALSE for normal operations: **-20%** (unnecessary human intervention)
- ❌ Pattern of excessive FALSE usage: **-50%** (defeats automation purpose)
- ❌ Using FALSE "just to be safe": **-20%** (misunderstanding automation)

## Common Scenarios

### Scenario 1: Code Review Found 3 Issues
**Status:** NORMAL
**Flag:** TRUE
**Reason:** Standard review-fix cycle

### Scenario 2: Fix Instructions Exist, Ready to Spawn
**Status:** NORMAL
**Flag:** TRUE
**Reason:** All prerequisites met

### Scenario 3: Second Fix Iteration (Still Issues After First Fix)
**Status:** NORMAL
**Flag:** TRUE
**Reason:** Iterative fixes are expected

### Scenario 4: Fix Instructions File Missing
**Status:** ERROR
**Flag:** FALSE
**Reason:** Cannot proceed without instructions (state corruption)

### Scenario 5: Effort Directory Not Found
**Status:** ERROR
**Flag:** FALSE
**Reason:** Infrastructure broken (unrecoverable)

## Exit Criteria

### Exit to MONITORING_EFFORT_FIXES when:
- ✅ Engineers spawned successfully
- ✅ Fix instructions provided to engineers
- ✅ State updated with fixes_in_progress

**BEFORE emitting continuation flag:**

```bash
# MANDATORY: Print cascade status if active (R406 auto-reporting)
if [ -f "orchestrator-state-v3.json" ] && \
   jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    echo ""
    echo "═══════════════════════════════════════════════════════════"
    echo "📊 R406 FIX CASCADE STATUS (automatic report)"
    echo "═══════════════════════════════════════════════════════════"
    source utilities/cascade-status-report.sh
    cascade_status_report
    echo "═══════════════════════════════════════════════════════════"
    echo ""
fi

# Emit continuation flag
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Standard case
```

### Exit to ERROR_RECOVERY when:
- 🚨 Critical infrastructure failure
- 🚨 State corruption detected
- 🚨 Cannot spawn engineers

**BEFORE emitting continuation flag:**

```bash
# MANDATORY: Print cascade status if active (R406 auto-reporting)
if [ -f "orchestrator-state-v3.json" ] && \
   jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    echo ""
    echo "═══════════════════════════════════════════════════════════"
    echo "📊 R406 FIX CASCADE STATUS (automatic report)"
    echo "═══════════════════════════════════════════════════════════"
    source utilities/cascade-status-report.sh
    cascade_status_report
    echo "═══════════════════════════════════════════════════════════"
    echo ""
fi

# Emit continuation flag
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # Exceptional case
```

---

**REMEMBER**: The software factory is designed to operate autonomously. Code review finding issues is the system **WORKING CORRECTLY**, not breaking! Use TRUE unless something is genuinely broken.

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
