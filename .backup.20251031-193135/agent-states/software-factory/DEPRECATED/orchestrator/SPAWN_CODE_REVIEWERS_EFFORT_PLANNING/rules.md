# SPAWN_CODE_REVIEWERS_EFFORT_PLANNING State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: CREATE_NEXT_INFRASTRUCTURE
**Exit To**: WAITING_FOR_EFFORT_PLANS

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**SPAWN Code Reviewer agents to create implementation plans for EACH effort in the wave.**

Code Reviewers create implementation plans because they understand:
- How to break work into PR-sized chunks
- What needs review and testing
- Size constraints and split strategies

## Required Inputs

### 1. Infrastructure State
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Read infrastructure that was just created
INFRA=$(jq -r '.infrastructure_created' orchestrator-state-v3.json)

if [ "$INFRA" = "null" ] || [ -z "$INFRA" ]; then
    echo "❌ FATAL: No infrastructure found in state file"
    echo "  CREATE_NEXT_INFRASTRUCTURE must complete before spawning reviewers"
    exit 1
fi

echo "✅ Infrastructure found for phase $PHASE wave $WAVE"
```

### 2. Effort List
```bash
# Extract all efforts from infrastructure
EFFORTS=$(jq -r '.infrastructure_created.efforts | keys[]' orchestrator-state-v3.json | sort)

if [ -z "$EFFORTS" ]; then
    echo "❌ FATAL: No efforts found in infrastructure"
    exit 1
fi

echo "📋 Efforts to create plans for:"
echo "$EFFORTS"
```

### 3. Wave Plan
```bash
WAVE_PLAN=$(ls -t "$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-PLAN--"*.md 2>/dev/null | head -1)

if [ -z "$WAVE_PLAN" ] || [ ! -f "$WAVE_PLAN" ]; then
    echo "❌ FATAL: Wave plan not found"
    exit 1
fi

echo "✅ Wave plan: $WAVE_PLAN"
```

## 🔴🔴🔴 REVIEWER SPAWNING PROTOCOL 🔴🔴🔴

### For EACH Effort, Spawn Code Reviewer

```bash
for effort in $EFFORTS; do
    # Get effort details from infrastructure
    EFFORT_BRANCH=$(jq -r ".infrastructure_created.efforts.\"${effort}\".effort_branch" orchestrator-state-v3.json)
    EFFORT_DIR=$(jq -r ".infrastructure_created.efforts.\"${effort}\".directory" orchestrator-state-v3.json)
    BASE_BRANCH=$(jq -r ".infrastructure_created.efforts.\"${effort}\".base_branch" orchestrator-state-v3.json)

    echo "🚀 Spawning Code Reviewer for: $effort"
    echo "  Branch: $EFFORT_BRANCH"
    echo "  Directory: $EFFORT_DIR"
    echo "  Base: $BASE_BRANCH"

    # Build task prompt for Code Reviewer
    TASK_PROMPT=$(cat << EOF
📋 CODE REVIEWER TASK: Create Implementation Plan for ${effort}

You are a Code Reviewer agent spawned to create an implementation plan.

**Context:**
- Project: $PROJECT_PREFIX
- Phase: $PHASE
- Wave: $WAVE
- Effort: $effort
- Wave Plan: $WAVE_PLAN

**Your Working Environment:**
- TARGET_DIRECTORY: $CLAUDE_PROJECT_DIR/$EFFORT_DIR
- BRANCH: $EFFORT_BRANCH
- BASE_BRANCH: $BASE_BRANCH

**Your Deliverable:**

Create an IMPLEMENTATION PLAN at:
\$CLAUDE_PROJECT_DIR/$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-PLAN--$(date +%Y%m%d-%H%M%S).md

**Plan Must Include:**
1. **Effort Overview**
   - What functionality this effort implements
   - How it relates to the wave plan
   - Dependencies on other efforts (if any)

2. **Implementation Steps**
   - Ordered list of development tasks
   - Expected file changes per step
   - Estimated line count per step

3. **Size Constraints**
   - Target: 300-600 lines total
   - Hard limit: 800 lines maximum
   - Split strategy if effort may exceed limits

4. **Test Requirements**
   - Unit tests needed
   - Integration tests needed
   - Test coverage expectations

5. **Review Checkpoints**
   - When to measure size
   - When to run code review
   - Split triggers

6. **Technical Metadata**
   - **PHASE**: $PHASE
   - **WAVE**: $WAVE
   - **EFFORT**: $effort
   - **BRANCH**: $EFFORT_BRANCH
   - **BASE**: $BASE_BRANCH

**Critical Requirements:**
- Follow R383 (timestamp in filename)
- Follow R343 (metadata directory structure)
- Size estimate MUST be realistic
- Plan MUST be achievable in one PR (<800 lines)
- If effort is too large, note that splitting will be required

**State Transition:**
After creating plan, you should transition to state indicating plan is complete.

EOF
)

    # Spawn Code Reviewer agent using Task tool
    # Subagent type: code-reviewer
    # Task description: "Create implementation plan for $effort"
    # Prompt: $TASK_PROMPT
    # Working directory: $CLAUDE_PROJECT_DIR/$EFFORT_DIR

    echo "✅ Code Reviewer spawned for: $effort"
done
```

## Parallelization Strategy

**ALL Code Reviewers should be spawned IN PARALLEL per R151.**

This is effort PLANNING, not implementation - there are no code conflicts.
All reviewers:
- Read the same wave plan
- Work in isolated directories
- Create independent planning documents
- Can run simultaneously

### Spawn All Reviewers in Single Message
```bash
# When actually spawning, use Task tool multiple times in one message:
# Task 1: Code Reviewer for effort-1
# Task 2: Code Reviewer for effort-2
# Task 3: Code Reviewer for effort-3
# ... etc

# This ensures R151 compliance (timestamps within 5 seconds)
```

## State Update

### Update orchestrator-state-v3.json
```bash
# Build list of spawned reviewers
REVIEWER_LIST=$(echo "$EFFORTS" | jq -R -s -c 'split("\n") | map(select(length > 0))')

jq --argjson reviewers "$REVIEWER_LIST" \
   --arg timestamp "$(date -Iseconds)" \
   '.state_machine.current_state = "WAITING_FOR_EFFORT_PLANS" |
    .state_machine.previous_state = "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING" |
    .effort_planning_in_progress = $reviewers |
    .state_transition_log += [{
        "from": "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING",
        "to": "WAITING_FOR_EFFORT_PLANS",
        "timestamp": $timestamp,
        "reason": "Code Reviewers spawned for all \($reviewers | length) efforts"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ State updated: Waiting for effort plans"
```

## Validation Requirements

### Pre-Spawn Validation
- ✅ Infrastructure exists in state file
- ✅ All effort directories exist
- ✅ All effort branches exist and pushed
- ✅ Wave plan exists and readable
- ✅ No existing implementation plans

### Post-Spawn Validation
- ✅ One Code Reviewer spawned per effort
- ✅ All spawns issued in parallel (single message)
- ✅ State file updated with effort list
- ✅ Timestamps recorded

## Integration with Rules

- **R151**: Parallelization Timestamp Requirements (5s deviation max)
- **R208**: Directory Spawn Protocol
- **R209**: Effort Directory Isolation
- **R383**: Metadata File Timestamp Requirements
- **R343**: Metadata Directory Standardization
- **R219**: Dependency-Aware Effort Planning

## Exit Criteria

Before transitioning to WAITING_FOR_EFFORT_PLANS:
- ✅ Code Reviewer spawned for EVERY effort
- ✅ All spawns completed in parallel
- ✅ Each reviewer given complete context
- ✅ Target directories specified correctly
- ✅ State file updated with tracking info

## Common Issues

### Issue: Missing Infrastructure
**Detection**: infrastructure_created is null in state file
**Resolution**: Return to CREATE_NEXT_INFRASTRUCTURE state

### Issue: Effort Directory Missing
**Detection**: Effort directory doesn't exist
**Resolution**: Re-run infrastructure creation

### Issue: Sequential Spawning
**Detection**: Reviewers spawned one at a time (timestamp deviation >5s)
**Resolution**: Use Task tool multiple times in SINGLE message

## R313 Enforcement - MANDATORY STOP

```bash
# This is the ABSOLUTE LAST thing that happens in this state
echo ""
echo "🛑 R313 ENFORCEMENT: STOPPING NOW"
echo "The orchestrator MUST NOT continue past this point."
echo "Code Reviewers have been spawned for effort planning."
echo "Wait for all reviewers to complete their implementation plans."
echo ""
echo "Next state will be: WAITING_FOR_EFFORT_PLANS"
echo "Resume with: /continue-orchestrating"
```

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Code Reviewer agents spawned successfully for ALL efforts
- ✅ All spawns completed in parallel (R151 compliance)
- ✅ Each reviewer given wave architecture and effort context
- ✅ State file updated with spawn tracking
- ✅ Ready to transition to WAITING_FOR_EFFORT_PLANS
- ✅ Following designed workflow

**THIS IS NORMAL WORKFLOW.** Spawning Code Reviewers for effort planning is the DESIGNED PROCESS.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot identify efforts to plan
- ❌ Wave architecture missing or corrupt
- ❌ Cannot spawn Code Reviewer agents
- ❌ State machine corruption

**DO NOT set FALSE because:**
- ❌ Spawning multiple reviewers (parallelization is NORMAL!)
- ❌ R322 requires stop (stop ≠ FALSE!)
- ❌ Waiting for plans (EXPECTED!)

**Correct pattern:** `exit 0` + `CONTINUE-SOFTWARE-FACTORY=TRUE`

## Automation Flag

```bash
# After successful spawn and state transition:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue orchestration per R405 - system will handle monitoring
```

---

**REMEMBER**: Code Reviewers create implementation plans because they understand PR sizing and review requirements. Spawn them ALL IN PARALLEL for R151 compliance!

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
