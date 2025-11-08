# 🚨🚨🚨 RULE R222: Code Review Gate Requirement 🚨🚨🚨

**Category:** Gate Requirements  
**Agents:** orchestrator (enforcement), code-reviewer (execution)  
**Criticality:** BLOCKING - No wave completion without ALL reviews passing  
**States:** MONITOR → WAVE_COMPLETE transition

## 🔴🔴🔴 ABSOLUTE GATE REQUIREMENT 🔴🔴🔴

**NO WAVE CAN BE MARKED COMPLETE WITHOUT:**
1. ALL efforts have been reviewed
2. ALL reviews have PASSED
3. ALL size compliance checks PASSED (<800 lines per effort)
4. NO effort is in FIX_ISSUES state
5. NO effort has pending review issues

## 🔴🔴🔴 ORCHESTRATOR RESPONSIBILITIES - CRITICAL! 🔴🔴🔴

**THE ORCHESTRATOR MUST ACTIVELY SPAWN CODE REVIEWERS!**

Reviews don't happen automatically! The orchestrator MUST:
1. **DETECT** when implementations complete (in MONITOR state)
2. **SPAWN** Code Reviewers IMMEDIATELY for completed implementations
3. **TRACK** review_status for EVERY effort (NOT_STARTED → IN_PROGRESS → PASSED/FAILED)
4. **BLOCK** WAVE_COMPLETE until ALL reviews are PASSED
5. **COORDINATE** review-fix loops when reviews fail

**CRITICAL**: If the orchestrator doesn't spawn Code Reviewers, no reviews will ever happen!
This is a GRADING FAILURE worth -100%!

### When to Spawn Code Reviewers:
```bash
# In MONITOR state, for EACH effort:
if implementation_status == "COMPLETE" && review_status != "IN_PROGRESS" && review_status != "PASSED":
    MUST SPAWN CODE REVIEWER IMMEDIATELY!
```

## ENFORCEMENT PROTOCOL

### 1. MANDATORY PRE-TRANSITION VERIFICATION

Before ANY transition to WAVE_COMPLETE state:

```bash
# R222 ENFORCEMENT - MUST CHECK EVERY EFFORT!
echo "🔍 R222: Verifying ALL reviews passed..."
ALL_PASSED=true

# Use text_editor tool with view command to read orchestrator-state-v3.json:
# Find the efforts_in_progress array
for effort in <efforts_in_progress array>; do
    # Check review status
    # Use text_editor tool with view command to read orchestrator-state-v3.json:
    # Find effort in efforts_completed array and get its review_status
    REVIEW_STATUS="<review_status for this effort>"
    
    # Check size compliance
    # Use text_editor tool with view command to read orchestrator-state-v3.json:
    # Find effort in efforts_completed array and get its lines_changed
    SIZE_LINES="<lines_changed for this effort>"
    
    if [ "$REVIEW_STATUS" != "PASSED" ] && [ "$REVIEW_STATUS" != "passed" ]; then
        echo "❌ R222 VIOLATION: $effort review status: $REVIEW_STATUS"
        echo "🚫 CANNOT TRANSITION TO WAVE_COMPLETE!"
        echo "🔄 Must remain in MONITOR and fix issues"
        ALL_PASSED=false
    fi
    
    if [ "$SIZE_LINES" -gt 800 ]; then
        echo "❌ SIZE VIOLATION: $effort has $SIZE_LINES lines (>800 limit)!"
        echo "🚫 CANNOT TRANSITION TO WAVE_COMPLETE!"
        ALL_PASSED=false
    fi
done

if [ "$ALL_PASSED" = false ]; then
    echo "🔴🔴🔴 CRITICAL: INVALID STATE TRANSITION BLOCKED!"
    echo "Must execute review-fix loops until all pass"
    exit 222
fi

echo "✅ R222 VERIFIED: All reviews passed, proceeding with completion"
```

### 2. GATE VIOLATIONS

**Attempting to enter WAVE_COMPLETE with failed reviews results in:**
- Immediate state machine violation
- -100% grading penalty
- Mandatory return to MONITOR state
- Required review-fix loop execution

### 3. REVIEW STATUS TRACKING

The orchestrator MUST maintain accurate review status in state file:

```yaml
efforts_completed:
  - name: effort-001
    review_status: PASSED  # REQUIRED: PASSED, FAILED, PENDING, or NOT_STARTED
    lines_changed: 450     # REQUIRED: Actual line count from line-counter.sh
    review_report: path/to/review-report.md
    issues_found: 0
    issues_fixed: 0
```

### 4. REVIEW-FIX LOOP REQUIREMENTS

When ANY review fails:

```yaml
state_transitions:
  MONITOR:
    review_failed: FIX_ISSUES  # Mandatory transition
  FIX_ISSUES:
    fixes_complete: CODE_REVIEW  # Re-review required
  CODE_REVIEW:
    review_passed: MONITOR  # Return to monitoring
    review_failed: FIX_ISSUES  # Continue loop
```

## VERIFICATION CHECKPOINTS

### Entry to WAVE_COMPLETE Checklist:
- [ ] ALL efforts show review_status: PASSED
- [ ] ALL efforts have lines_changed <= 800
- [ ] NO efforts in efforts_in_progress array
- [ ] ALL efforts in efforts_completed array
- [ ] Review reports exist for ALL efforts
- [ ] State file updated with review outcomes

### Monitoring Responsibilities:
1. Track review status in real-time
2. Prevent premature transitions
3. Enforce review-fix loops
4. Validate size compliance
5. Document all review outcomes

## PENALTIES FOR VIOLATIONS

| Violation | Penalty | Recovery |
|-----------|---------|----------|
| Enter WAVE_COMPLETE with ANY failed review | -100% | Return to MONITOR, fix all issues |
| Enter WAVE_COMPLETE with pending reviews | -50% | Complete all reviews first |
| Missing review status tracking | -25% | Update state file completely |
| Incomplete review-fix loop | -40% | Complete full loop |
| Size violation (>800 lines) | -100% | Split and re-review |

## ACKNOWLEDGMENT REQUIREMENT

The orchestrator MUST acknowledge before WAVE_COMPLETE transition:

```markdown
"I acknowledge R222: I CANNOT transition to WAVE_COMPLETE unless:
- ALL reviews are PASSED
- ALL efforts are <800 lines
- NO pending issues exist
- ALL review reports are complete
Attempting transition with ANY failure = -100% penalty"
```

## IMPLEMENTATION NOTES

1. **Review Status Source of Truth**: orchestrator-state-v3.json
2. **Size Measurement Tool**: $CLAUDE_PROJECT_DIR/tools/line-counter.sh
3. **Review Reports Location**: review-reports/[wave]/[effort]/
4. **State Machine Authority**: software-factory-3.0-state-machine.json

## RELATED RULES

- R108: Code Review Protocol (review execution)
- R153: Review Turnaround (timing requirements)
- R206: State Machine Validation (transition rules)
- R215: Orchestrator State Ownership (state file authority)
- R234: Mandatory State Traversal (state progression)

---

**REMEMBER**: This gate is ABSOLUTE. No exceptions, no workarounds, no "almost passing". ALL reviews MUST pass before wave completion.