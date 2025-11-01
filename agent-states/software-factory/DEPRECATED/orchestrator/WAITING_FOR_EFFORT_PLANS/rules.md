# WAITING_FOR_EFFORT_PLANS State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
**Exit To**: ANALYZE_IMPLEMENTATION_PARALLELIZATION

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**WAIT for all Code Reviewer agents to complete their implementation plans, then validate plan quality.**

This is a MONITORING_SWE_PROGRESS and VALIDATION state - the orchestrator checks that all plans are created and meet quality standards.

## Required Inputs

### 1. List of Expected Plans
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Get list of efforts we're waiting for
EFFORTS=$(jq -r '.effort_planning_in_progress[]' orchestrator-state-v3.json 2>/dev/null)

if [ -z "$EFFORTS" ]; then
    # Fallback: read from infrastructure
    EFFORTS=$(jq -r '.infrastructure_created.efforts | keys[]' orchestrator-state-v3.json)
fi

if [ -z "$EFFORTS" ]; then
    echo "❌ FATAL: No efforts found to wait for"
    exit 1
fi

echo "⏳ Waiting for implementation plans for:"
echo "$EFFORTS"
```

## 🔴🔴🔴 PLAN CHECKING PROTOCOL 🔴🔴🔴

### Step 1: Check If All Plans Exist
```bash
echo "🔍 Checking for implementation plans..."

ALL_PLANS_EXIST=true
MISSING_PLANS=""

for effort in $EFFORTS; do
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"

    # Look for implementation plan with timestamp
    IMPL_PLAN=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)

    if [ -z "$IMPL_PLAN" ] || [ ! -f "$IMPL_PLAN" ]; then
        echo "❌ Missing plan for: $effort"
        ALL_PLANS_EXIST=false
        MISSING_PLANS="$MISSING_PLANS $effort"
    else
        echo "✅ Plan found: $effort"
    fi
done

if [ "$ALL_PLANS_EXIST" = false ]; then
    echo "⏳ Still waiting for plans: $MISSING_PLANS"
    echo "Orchestrator should check back in a few minutes or wait for agent completion signals"
    exit 0  # Not an error, just not ready yet
fi

echo "✅ All implementation plans found!"
```

### Step 2: Validate Plan Quality (R502)

For each plan, check:

```bash
echo "🔍 Validating plan quality per R502..."

VALIDATION_FAILURES=""

for effort in $EFFORTS; do
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    IMPL_PLAN=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)

    echo "Validating: $effort"

    # Check 1: File follows R383/R343 naming
    if ! echo "$IMPL_PLAN" | grep -q "IMPLEMENTATION-PLAN--[0-9]\{8\}-[0-9]\{6\}\.md"; then
        echo "⚠️ $effort: Plan filename doesn't follow R383 timestamp format"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Bad filename format"
    fi

    # Check 2: Contains required metadata
    if ! grep -q "^\\*\\*PHASE\\*\\*:" "$IMPL_PLAN"; then
        echo "⚠️ $effort: Missing PHASE metadata"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Missing PHASE"
    fi

    if ! grep -q "^\\*\\*WAVE\\*\\*:" "$IMPL_PLAN"; then
        echo "⚠️ $effort: Missing WAVE metadata"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Missing WAVE"
    fi

    if ! grep -q "^\\*\\*EFFORT\\*\\*:" "$IMPL_PLAN"; then
        echo "⚠️ $effort: Missing EFFORT metadata"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Missing EFFORT"
    fi

    if ! grep -q "^\\*\\*BRANCH\\*\\*:" "$IMPL_PLAN"; then
        echo "⚠️ $effort: Missing BRANCH metadata"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Missing BRANCH"
    fi

    # Check 3: Contains size estimate
    if ! grep -qi "estimate.*line\|line.*estimate\|size.*estimate" "$IMPL_PLAN"; then
        echo "⚠️ $effort: Missing size estimate"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Missing size estimate"
    fi

    # Check 4: Contains implementation steps
    if ! grep -qi "implementation.*step\|step.*implementation\|task\|action.*item" "$IMPL_PLAN"; then
        echo "⚠️ $effort: Missing implementation steps"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Missing implementation steps"
    fi

    # Check 5: Contains test requirements
    if ! grep -qi "test\|testing\|coverage" "$IMPL_PLAN"; then
        echo "⚠️ $effort: Missing test requirements"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Missing test requirements"
    fi

    # Check 6: Plan is not suspiciously short
    LINE_COUNT=$(wc -l < "$IMPL_PLAN")
    if [ "$LINE_COUNT" -lt 50 ]; then
        echo "⚠️ $effort: Plan seems too short ($LINE_COUNT lines)"
        VALIDATION_FAILURES="$VALIDATION_FAILURES\n- $effort: Suspiciously short plan ($LINE_COUNT lines)"
    fi
done

if [ -n "$VALIDATION_FAILURES" ]; then
    echo "❌ PLAN VALIDATION FAILURES:"
    echo -e "$VALIDATION_FAILURES"
    echo ""
    echo "🔄 Orchestrator should:"
    echo "1. Report validation failures"
    echo "2. Decide whether to:"
    echo "   a) Request Code Reviewers to fix issues"
    echo "   b) Proceed with warnings if issues are minor"
    echo "   c) Halt if issues are critical"
    exit 1
fi

echo "✅ All plans validated successfully!"
```

### Step 3: Extract Metadata for State File
```bash
echo "💾 Extracting plan metadata..."

# Build metadata object for each effort
EFFORT_METADATA=$(jq -n '{}')

for effort in $EFFORTS; do
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    IMPL_PLAN=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)

    # Extract branch from plan
    BRANCH=$(grep "^\\*\\*BRANCH\\*\\*:" "$IMPL_PLAN" | cut -d: -f2- | xargs)

    # Extract estimated size
    SIZE_ESTIMATE=$(grep -i "estimate.*line" "$IMPL_PLAN" | grep -oE "[0-9]+" | head -1)
    if [ -z "$SIZE_ESTIMATE" ]; then
        SIZE_ESTIMATE="unknown"
    fi

    # Build metadata
    EFFORT_METADATA=$(echo "$EFFORT_METADATA" | jq \
        --arg effort "$effort" \
        --arg plan "$IMPL_PLAN" \
        --arg branch "$BRANCH" \
        --arg size "$SIZE_ESTIMATE" \
        '.[$effort] = {
            plan_file: $plan,
            branch: $branch,
            estimated_size: $size,
            status: "READY_FOR_IMPLEMENTATION"
        }')
done

echo "✅ Metadata extracted"
```

## State Transition

### Update orchestrator-state-v3.json
```bash
jq --argjson metadata "$EFFORT_METADATA" \
   --arg timestamp "$(date -Iseconds)" \
   '.effort_plans = $metadata |
    .effort_planning_in_progress = null |
    .state_machine.current_state = "ANALYZE_IMPLEMENTATION_PARALLELIZATION" |
    .state_machine.previous_state = "WAITING_FOR_EFFORT_PLANS" |
    .state_transition_log += [{
        "from": "WAITING_FOR_EFFORT_PLANS",
        "to": "ANALYZE_IMPLEMENTATION_PARALLELIZATION",
        "timestamp": $timestamp,
        "reason": "All effort plans completed and validated"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ State updated: Moving to parallelization analysis"
```

## Validation Requirements

### Pre-Check Validation
- ✅ Have list of expected efforts
- ✅ Know where to look for plans
- ✅ Infrastructure exists

### Post-Check Validation
- ✅ All plans exist
- ✅ All plans follow R383/R343 conventions
- ✅ All plans contain required metadata
- ✅ All plans have size estimates
- ✅ All plans have implementation steps
- ✅ Metadata extracted and saved to state

## Integration with Rules

- **R502**: Plan Validation Gates
- **R383**: Metadata File Timestamp Requirements
- **R343**: Metadata Directory Standardization
- **R219**: Dependency-Aware Effort Planning
- **R206**: State Machine Validation

## Exit Criteria

Before transitioning to ANALYZE_IMPLEMENTATION_PARALLELIZATION:
- ✅ ALL effort plans exist
- ✅ ALL plans validated per R502
- ✅ Metadata extracted and saved
- ✅ State file updated
- ✅ No critical validation failures

## Common Issues

### Issue: Plans Taking Too Long
**Detection**: Waiting >30 minutes with no progress
**Resolution**: Check agent status, may need to respawn Code Reviewers

### Issue: Plan Validation Failures
**Detection**: Plans exist but don't meet quality standards
**Resolution**: Decide severity:
  - Minor issues: Document and proceed with warnings
  - Major issues: Respawn Code Reviewers with feedback

### Issue: Missing Metadata
**Detection**: Plans don't contain required fields
**Resolution**: Request Code Reviewers to update plans

## Waiting Strategy

The orchestrator should:
1. Check every 2-3 minutes for plan completion
2. After 15 minutes, escalate monitoring frequency
3. After 30 minutes, investigate what's blocking agents
4. Provide progress updates to user

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Waiting for Code Reviewers to complete effort plans
- ✅ Some plans still in progress (normal waiting)
- ✅ All plans completed and validated
- ✅ Validation failures detected (system will handle recovery)
- ✅ Ready to transition to parallelization analysis
- ✅ Following designed workflow

**THIS IS NORMAL WORKFLOW.** Waiting for plans, checking periodically, and
transitioning when ready is the DESIGNED PROCESS. Even validation failures
are handled by the system - they don't require manual intervention.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot access plan locations or state file
- ❌ All Code Reviewers disappeared (no activity >1 hour)
- ❌ Critical infrastructure corruption
- ❌ State machine corruption

**DO NOT set FALSE because:**
- ❌ Still waiting for plans (NORMAL!)
- ❌ Validation failures (system handles!)
- ❌ Plans complete, transitioning (EXPECTED!)
- ❌ R322 requires stop (stop ≠ FALSE!)

**Correct pattern:** All waiting outcomes use TRUE

## Continuation Control

```bash
# Waiting states check periodically and continue or transition
if [ "$ALL_PLANS_EXIST" = false ]; then
    echo "⏳ Still waiting for plans to complete..."
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue waiting/checking
elif [ -n "$VALIDATION_FAILURES" ]; then
    echo "❌ Validation failures detected - transitioning to error recovery"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue per R405 - system will handle recovery
else
    echo "✅ All plans ready, proceeding to parallelization analysis"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"   # Continue to next state
fi
```

---

**REMEMBER**: This is a VALIDATION gate. Don't proceed to implementation until ALL plans are complete and validated!

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
