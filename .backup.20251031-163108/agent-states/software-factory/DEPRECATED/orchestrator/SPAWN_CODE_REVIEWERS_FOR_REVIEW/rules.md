# SPAWN_CODE_REVIEWERS_EFFORT_REVIEW State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: MONITORING_SWE_PROGRESS
**Exit To**: MONITORING_EFFORT_REVIEWS

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**SPAWN Code Reviewer agents to perform code reviews for completed effort implementations.**

This state enforces the critical CODE REVIEW PROTOCOL - every implementation MUST be reviewed before being considered complete.

## Required Inputs

### 1. List of Efforts Ready for Review
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Get efforts ready for review
EFFORTS_TO_REVIEW=$(jq -r '.efforts_ready_for_review[]?' orchestrator-state-v3.json 2>/dev/null)

if [ -z "$EFFORTS_TO_REVIEW" ]; then
    echo "❌ FATAL: No efforts ready for review"
    echo "  This state should only be entered when efforts complete"
    exit 1
fi

echo "📋 Spawning reviewers for:"
echo "$EFFORTS_TO_REVIEW"
```

### 2. Effort Infrastructure Details
```bash
for effort in $EFFORTS_TO_REVIEW; do
    EFFORT_BRANCH=$(jq -r ".infrastructure_created.efforts.\"${effort}\".effort_branch" orchestrator-state-v3.json)
    BASE_BRANCH=$(jq -r ".infrastructure_created.efforts.\"${effort}\".base_branch" orchestrator-state-v3.json)
    EFFORT_DIR=$(jq -r ".infrastructure_created.efforts.\"${effort}\".directory" orchestrator-state-v3.json)

    echo "  - $effort: $EFFORT_BRANCH (base: $BASE_BRANCH)"
done
```

## 🔴🔴🔴 CODE REVIEWER SPAWNING PROTOCOL 🔴🔴🔴

### For EACH Completed Effort, Spawn Code Reviewer

```bash
for effort in $EFFORTS_TO_REVIEW; do
    # Get effort details
    EFFORT_BRANCH=$(jq -r ".infrastructure_created.efforts.\"${effort}\".effort_branch" orchestrator-state-v3.json)
    EFFORT_DIR=$(jq -r ".infrastructure_created.efforts.\"${effort}\".directory" orchestrator-state-v3.json)
    BASE_BRANCH=$(jq -r ".infrastructure_created.efforts.\"${effort}\".base_branch" orchestrator-state-v3.json)

    # Get implementation plan
    IMPL_PLAN=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)

    echo "🚀 Spawning Code Reviewer for: $effort"
    echo "  Branch: $EFFORT_BRANCH"
    echo "  Directory: $EFFORT_DIR"
    echo "  Base: $BASE_BRANCH"
    echo "  Plan: $IMPL_PLAN"

    # Build task prompt for Code Reviewer
    TASK_PROMPT=$(cat << EOF
📋 CODE REVIEWER TASK: Review Implementation of ${effort}

You are a Code Reviewer agent spawned to review a completed implementation.

**Context:**
- Project: $PROJECT_PREFIX
- Phase: $PHASE
- Wave: $WAVE
- Effort: $effort
- Implementation Plan: $IMPL_PLAN

**Your Working Environment:**
- TARGET_DIRECTORY: $CLAUDE_PROJECT_DIR/$EFFORT_DIR
- BRANCH: $EFFORT_BRANCH
- BASE_BRANCH: $BASE_BRANCH

**Your Deliverable:**

Create a CODE REVIEW REPORT at:
\$CLAUDE_PROJECT_DIR/$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-REPORT--$(date +%Y%m%d-%H%M%S).md

**Review Must Cover:**

1. **Code Quality**
   - Follows project coding standards
   - Proper error handling
   - Appropriate comments and documentation
   - No obvious bugs or issues

2. **Size Compliance (R220/R304)**
   - **CRITICAL**: Use line-counter.sh to measure size
   - Command: \`$CLAUDE_PROJECT_DIR/tools/line-counter.sh -b $BASE_BRANCH -c $EFFORT_BRANCH\`
   - Verify total lines ≤ 800 (HARD LIMIT)
   - Document actual line count in report

3. **Implementation Completeness**
   - All items from implementation plan completed
   - Acceptance criteria met
   - No incomplete features

4. **Test Coverage**
   - Tests written per plan requirements
   - Tests passing
   - Coverage adequate

5. **Architecture Compliance**
   - Follows wave architecture plan
   - Consistent with project patterns
   - No architectural violations

6. **Git Hygiene**
   - All changes committed
   - Commit messages clear
   - No uncommitted or untracked files

**Findings Classification:**
- **BLOCKING**: Must fix before merge (bugs, size violations, missing tests)
- **HIGH**: Should fix before merge (quality issues, incomplete features)
- **MEDIUM**: Should fix soon (minor improvements, documentation)
- **LOW**: Nice to have (suggestions, optimizations)

**Required Actions:**
If you find issues:
1. Document ALL findings in CODE-REVIEW-REPORT
2. Create CODE-REVIEW-INSTRUCTIONS at:
   \$CLAUDE_PROJECT_DIR/$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-INSTRUCTIONS--$(date +%Y%m%d-%H%M%S).md
3. Instructions must be actionable for SW Engineer to fix

If no issues found:
1. Document clean review in CODE-REVIEW-REPORT
2. Mark implementation as APPROVED

**Critical Requirements:**
- Follow R383 (timestamp in filename)
- Follow R343 (metadata directory structure)
- Follow R304 (MUST use line-counter.sh, never wc -l)
- Size measurement is MANDATORY and BLOCKING

**Working Directory:** $CLAUDE_PROJECT_DIR/$EFFORT_DIR
**Branch:** $EFFORT_BRANCH

EOF
)

    # Spawn Code Reviewer agent using Task tool
    # Subagent type: code-reviewer
    # Task description: "Review implementation of $effort"
    # Prompt: $TASK_PROMPT
    # Working directory: $CLAUDE_PROJECT_DIR/$EFFORT_DIR

    echo "✅ Code Reviewer spawned for: $effort"
done
```

## Parallelization Strategy

**ALL Code Reviewers should be spawned IN PARALLEL per R151.**

Code reviews are independent:
- Each reviews different code in isolated directory
- No conflicts between reviewers
- All can run simultaneously

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
REVIEWER_LIST=$(echo "$EFFORTS_TO_REVIEW" | jq -R -s -c 'split("\n") | map(select(length > 0))')

jq --argjson reviewers "$REVIEWER_LIST" \
   --arg timestamp "$(date -Iseconds)" \
   '.reviews_in_progress = $reviewers |
    .efforts_ready_for_review = null |
    .state_machine.current_state = "MONITORING_EFFORT_REVIEWS" |
    .state_machine.previous_state = "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW" |
    .state_transition_log += [{
        "from": "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW",
        "to": "MONITORING_EFFORT_REVIEWS",
        "timestamp": $timestamp,
        "reason": "Code Reviewers spawned for \($reviewers | length) completed efforts"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ State updated: Monitoring reviews"
```

## Validation Requirements

### Pre-Spawn Validation
- ✅ Efforts completed implementation
- ✅ All code committed and pushed
- ✅ Implementation plans exist
- ✅ Effort directories accessible
- ✅ Git branches exist

### Post-Spawn Validation
- ✅ One Code Reviewer spawned per effort
- ✅ All spawns issued in parallel (single message)
- ✅ State file updated with review tracking
- ✅ Timestamps recorded

## Integration with Rules

- **R151**: Parallelization Timestamp Requirements (5s deviation max)
- **R304**: Mandatory Line Counter Tool Enforcement (BLOCKING)
- **R220/R221**: Size Limits and Continuous Delivery
- **R208/R209**: Directory Isolation
- **R383/R343**: Metadata File Standards
- **Grading**: Workflow Compliance (25%), Size Compliance (20%)

## Exit Criteria

Before transitioning to MONITORING_EFFORT_REVIEWS:
- ✅ Code Reviewer spawned for EVERY completed effort
- ✅ All spawns completed in parallel
- ✅ Each reviewer given complete context and tools
- ✅ State file updated with tracking info
- ✅ Review requirements clearly specified

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Fixes completed successfully
- ✅ Code committed and pushed
- ✅ Ready to spawn Code Reviewer for re-review
- ✅ Review instructions/metadata prepared
- ✅ All infrastructure working normally
- ✅ Code Reviewers spawned successfully
- ✅ Following designed workflow

**THIS IS NORMAL WORKFLOW!** After engineers fix issues, spawning a code reviewer
to verify the fixes is the DESIGNED PROCESS. This is automation working correctly.
Spawning agents and transitioning to monitoring states is EXPECTED BEHAVIOR.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot determine which effort needs re-review
- ❌ Fix commit failed or is corrupt
- ❌ Effort infrastructure broken/missing
- ❌ State machine corruption detected
- ❌ Cannot access fix results
- ❌ Required metadata files missing or corrupt
- ❌ Unrecoverable error prevents proceeding

**DO NOT set FALSE because:**
- ❌ Spawning reviewer (this is NORMAL workflow)
- ❌ Re-review needed after fixes (this is EXPECTED process)
- ❌ R322 requires stop (stop ≠ FALSE flag!)
- ❌ "User might want to see results" (only if exceptional)
- ❌ Transitioning to monitoring state (NORMAL operation)

### Critical Distinction: R322 Stop vs Continuation Flag

**R322 requires:**
1. Stop conversation (`exit 0`) ✅ - Context preservation
2. Save state ✅ - State persistence
3. Emit flag ✅ - Automation control

**R322 does NOT require:**
- Setting FALSE for normal operations ❌
- Human review of standard workflow ❌
- Halting automation ❌

**Correct pattern for this state:**
```bash
# State work completed successfully
exit 0  # R322 stop for context preservation
```

**Last line before exit:**
```bash
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Normal operation - allow automation
```

**Grading Impact:**
- Using FALSE for normal operations: -20% per violation
- Pattern of incorrect FALSE usage: -50%
- Breaking automation unnecessarily: -30%

## Common Issues

### Issue: Implementation Not Actually Complete
**Detection**: Code not committed, tests not written
**Resolution**: Don't spawn reviewer yet, continue monitoring implementation

### Issue: Multiple Review Rounds
**Detection**: First review finds issues, need second review after fixes
**Resolution**: Track review iteration count, spawn new reviewer after fixes

### Issue: Size Violation at Review Time
**Detection**: Code Reviewer discovers >800 lines
**Resolution**: BLOCK merge, transition to split planning state

## R313 Enforcement - MANDATORY STOP (Context Preservation)

```bash
# This is the ABSOLUTE LAST thing that happens in this state
echo ""
echo "🛑 R313 ENFORCEMENT: STOPPING INFERENCE NOW (to preserve context)"
echo "The orchestrator MUST stop inference to prevent context overflow."
echo "Code Reviewers have been spawned to review completed implementations."
echo ""
echo "⚠️ IMPORTANT: This is a NORMAL stop for context preservation, not an error!"
echo "Next state will be: MONITORING_EFFORT_REVIEWS"
echo "System will automatically continue when ready."
```

## Automation Flag

```bash
# After successful spawn and state transition:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # TRUE because code review is NORMAL operation!
# The system stops inference but sets TRUE to allow automatic restart
```

---

**REMEMBER**: Code review is MANDATORY per grading criteria. Every implementation must be reviewed before being considered complete!

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
