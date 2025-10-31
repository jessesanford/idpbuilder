# SPAWN_CODE_REVIEWER_FIX_PLAN State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: MONITORING_EFFORT_REVIEWS, MONITORING_EFFORT_FIXES (if re-review finds more issues)
**Exit To**: WAITING_FOR_FIX_PLANS

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**SPAWN Code Reviewer agent to create comprehensive fix plans for issues discovered during code review.**

This is a **STANDARD WORKFLOW STATE** - spawning a code reviewer to create detailed fix instructions is **NORMAL PRACTICE** when code review identifies issues that need structured fix planning.

## Required Inputs

### 1. Identify Efforts Needing Fix Plans
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Efforts that need fix planning (from review or integration failures)
EFFORTS_NEEDING_FIX_PLANS=$(jq -r '.efforts_needing_fix_plans[]?' orchestrator-state-v3.json 2>/dev/null)

if [ -z "$EFFORTS_NEEDING_FIX_PLANS" ]; then
    echo "❌ ERROR: No efforts marked as needing fix plans!"
    echo "State corruption - should not be in SPAWN_CODE_REVIEWER_FIX_PLAN without efforts list"
    exit 1
fi

echo "📋 Creating fix plans for:"
echo "$EFFORTS_NEEDING_FIX_PLANS"
```

### 2. Verify Review Reports Exist
```bash
for effort in $EFFORTS_NEEDING_FIX_PLANS; do
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"

    # Find the latest review report
    REVIEW_REPORT=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-REPORT--"*.md 2>/dev/null | head -1)

    if [ -z "$REVIEW_REPORT" ] || [ ! -f "$REVIEW_REPORT" ]; then
        echo "❌ CRITICAL ERROR: Review report missing for $effort"
        echo "Expected at: $EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-REPORT--*.md"
        exit 1
    fi

    echo "✅ Review report found for $effort: $REVIEW_REPORT"
done
```

## 🔴🔴🔴 CODE REVIEWER SPAWNING PROTOCOL 🔴🔴🔴

### Spawn Reviewer to Create Fix Plans

```bash
echo "🚀 SPAWNING CODE REVIEWER FOR FIX PLANNING - $(date -Iseconds)"
echo "═══════════════════════════════════════════════════════"

# Spawn code reviewer to analyze issues and create fix plans
echo "Agent: @agent-code-reviewer"
echo "State: CREATE_FIX_PLAN"
echo ""
echo "Task: Analyze code review reports and create comprehensive fix instructions"
echo "Context:"
echo "  - Efforts needing fix plans: $EFFORTS_NEEDING_FIX_PLANS"
echo "  - Review reports contain issues that need structured fixes"
echo "  - Create CODE-REVIEW-INSTRUCTIONS files with detailed fix guidance"
echo ""

# Reviewer will:
# 1. Read review reports for each effort
# 2. Analyze issues and categorize by priority
# 3. Create detailed fix instructions (CODE-REVIEW-INSTRUCTIONS--*.md)
# 4. Include specific guidance for each issue
# 5. Mark as COMPLETED when done

echo "═══════════════════════════════════════════════════════"
```

## State Transition

### Update State After Spawning
```bash
# Transition to waiting for fix plans
jq --arg timestamp "$(date -Iseconds)" \
   --argjson efforts "$(echo $EFFORTS_NEEDING_FIX_PLANS | jq -R -s -c 'split(" ") | map(select(length > 0))')" \
   '.efforts_needing_fix_plans = null |
    .fix_planning_in_progress = $efforts |
    .state_machine.current_state = "WAITING_FOR_FIX_PLANS" |
    .state_machine.previous_state = "SPAWN_CODE_REVIEWER_FIX_PLAN" |
    .state_transition_log += [{
        "from": "SPAWN_CODE_REVIEWER_FIX_PLAN",
        "to": "WAITING_FOR_FIX_PLANS",
        "timestamp": $timestamp,
        "reason": "Code reviewer spawned to create fix plans"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ State updated: WAITING_FOR_FIX_PLANS"
```

## Integration with Rules

- **R313/R322**: Mandatory stop after spawning (context preservation)
- **R405**: Continuation flag (see below for critical guidance)
- **R287**: TODO Persistence before transition
- **R383/R343**: Metadata file standards for fix instructions
- **Grading**: Workflow Compliance (25%)

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Review reports exist and are readable
- ✅ Effort directories exist and are accessible
- ✅ Code reviewer spawned successfully
- ✅ All normal operations completed
- ✅ State transition successful

**THIS IS THE STANDARD CASE.** Spawning a code reviewer to create fix plans when issues are found is **NORMAL SOFTWARE DEVELOPMENT WORKFLOW**. The review → plan → fix → re-review cycle is designed to handle this automatically.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Review report files are missing or corrupt
- ❌ Cannot determine which efforts need fix plans (state corruption)
- ❌ Effort infrastructure is broken/missing
- ❌ State machine corruption detected
- ❌ Unrecoverable error prevents spawning reviewer
- ❌ Critical system violation detected

### ❌ DO NOT SET FALSE JUST BECAUSE:

**These are NOT reasons to use FALSE:**
- ❌ Code review found complex issues (this is **NORMAL**)
- ❌ Fix planning is needed (this is **EXPECTED**)
- ❌ Multiple efforts need fix plans (this is **DESIGNED FOR**)
- ❌ Issues are categorized as BLOCKING/HIGH (this is **NORMAL**)
- ❌ "User might want to review the plan" (only if truly exceptional)
- ❌ This is a complex fix requiring planning (this is **WHY THIS STATE EXISTS**)

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

**For SPAWN_CODE_REVIEWER_FIX_PLAN:** Almost always TRUE because spawning reviewers for fix planning is normal workflow!

## Grading Impact

**Correct flag usage:**
- ✅ Using TRUE for normal fix planning: No penalty
- ✅ Using FALSE for genuine errors: No penalty

**Incorrect flag usage:**
- ❌ Using FALSE for normal operations: **-20%** (unnecessary human intervention)
- ❌ Pattern of excessive FALSE usage: **-50%** (defeats automation purpose)
- ❌ Using FALSE "just to be safe": **-20%** (misunderstanding automation)

## Common Scenarios

### Scenario 1: Code Review Found Blocking Issues
**Status:** NORMAL
**Flag:** TRUE
**Reason:** Fix planning is the designed response

### Scenario 2: Multiple Efforts Need Detailed Fix Instructions
**Status:** NORMAL
**Flag:** TRUE
**Reason:** Reviewer can handle multiple efforts

### Scenario 3: Complex Issues Requiring Structured Planning
**Status:** NORMAL
**Flag:** TRUE
**Reason:** This is exactly what this state is for

### Scenario 4: Review Reports Missing
**Status:** ERROR
**Flag:** FALSE
**Reason:** Cannot create fix plans without review data (state corruption)

### Scenario 5: Effort Infrastructure Not Found
**Status:** ERROR
**Flag:** FALSE
**Reason:** Infrastructure broken (unrecoverable)

## Exit Criteria

### Exit to WAITING_FOR_FIX_PLANS when:
- ✅ Code reviewer spawned successfully
- ✅ Review reports provided to reviewer
- ✅ State updated with fix_planning_in_progress
- ✅ **Flag: CONTINUE-SOFTWARE-FACTORY=TRUE** (standard case)

### Exit to ERROR_RECOVERY when:
- 🚨 Critical infrastructure failure
- 🚨 State corruption detected
- 🚨 Cannot spawn reviewer
- 🚨 **Flag: CONTINUE-SOFTWARE-FACTORY=FALSE** (exceptional case)

---

**REMEMBER**: Creating fix plans is a designed part of the workflow. Issues requiring fix planning is the system **WORKING AS INTENDED**, not breaking! Use TRUE unless something is genuinely broken.

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
