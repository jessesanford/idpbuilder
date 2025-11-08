# 🚨 RULE R535 - Code Reviewer Size Violation Threshold (900 Lines)

**Criticality:** CRITICAL - Controls bug creation for size violations
**Grading Impact:** -20% for false positive size violation bugs
**Enforcement:** Code Reviewer ONLY
**Scope:** Code Reviewer (EXCLUSIVE - SW Engineers still see 800 line limit)

## Rule Statement

**Code Reviewers should ONLY create size violation bugs when measured implementation lines exceed 900 lines, NOT at 800 lines.**

This creates a 100-line "grace buffer" between the stated limit (800) and actual enforcement (900), reducing false positives while maintaining developer discipline.

## 🔴🔴🔴 STRATEGIC ASYMMETRY 🔴🔴🔴

### Different Agents See Different Thresholds:

| Agent | Sees Limit | Actual Enforcement | Rationale |
|-------|-----------|-------------------|-----------|
| **SW Engineer** | 800 lines | N/A | Keeps developers conservative |
| **Code Reviewer** | 900 lines | Creates bugs >900 | Only escalates true violations |
| **Orchestrator** | 800 lines | Delegates to reviewer | Follows standard guidance |
| **Architect** | 800 lines | Reviews for compliance | Architectural sizing |

**PURPOSE**: SW Engineers naturally stay under 800 lines providing a safety buffer. Code Reviewers only escalate genuine violations (>900), reducing noise and unnecessary splits.

## Size Violation Thresholds

### For Code Reviewer Bug Creation:

| Lines | SW Engineer Sees | Code Reviewer Action | Bug Created? |
|-------|-----------------|---------------------|--------------|
| 0-700 | ✅ COMPLIANT | ✅ APPROVE | NO |
| 700-800 | ⚠️ WARNING | ✅ APPROVE | NO |
| 800-900 | 🚨 VIOLATION (to them) | ✅ APPROVE (grace buffer) | **NO** |
| >900 | 🚨🚨🚨 VIOLATION | 🚨 CREATE BUG | **YES** |

### Bug Creation Decision Logic:

```bash
create_size_violation_bug() {
    local measured_lines="$1"
    local effort_id="$2"

    # Code Reviewer uses 900-line threshold
    if [ "$measured_lines" -gt 900 ]; then
        echo "🚨 SIZE VIOLATION: $measured_lines > 900 lines"
        echo "Creating BUG-XXX-SIZE-VIOLATION..."

        create_bug \
            --severity CRITICAL \
            --category SIZE_VIOLATION \
            --title "Hard Size Limit Violation - ${measured_lines} Lines" \
            --description "Exceeds 900-line enforcement threshold (${measured_lines} lines, +$((measured_lines - 900)) over enforcement limit)" \
            --effort "$effort_id"
    else
        echo "✅ SIZE WITHIN ENFORCEMENT THRESHOLD: $measured_lines <= 900 lines"
        echo "Grace buffer active (800-900 lines)"
        # NO BUG CREATED
    fi
}
```

## Code Reviewer Responsibilities

### Size Measurement and Review:

```bash
# During EFFORT_PLAN_CREATION or CREATE_SPLIT_PLAN

# 1. Measure lines using official tool
line_count=$($CLAUDE_PROJECT_DIR/tools/line-counter.sh -c effort-branch | grep "Total:" | awk '{print $2}')

echo "📊 Size Measurement:"
echo "  Measured lines: $line_count"
echo "  SW Engineer limit (they see): 800 lines"
echo "  Code Reviewer enforcement (grace buffer): 900 lines"

# 2. Apply Code Reviewer threshold (900 lines)
if [ "$line_count" -le 900 ]; then
    echo "✅ WITHIN ENFORCEMENT THRESHOLD"
    echo "Decision: APPROVE (no bug)"

    if [ "$line_count" -gt 800 ]; then
        echo "ℹ️ Note: SW Engineer would see this as violation"
        echo "ℹ️ Grace buffer: 800-900 lines = approved by reviewer"
    fi

    review_status="APPROVED"
    size_violation_bug="NONE"

elif [ "$line_count" -gt 900 ]; then
    echo "🚨 EXCEEDS ENFORCEMENT THRESHOLD"
    echo "Decision: CREATE SIZE VIOLATION BUG"
    echo "Lines over enforcement limit: $((line_count - 900))"

    review_status="NEEDS_FIXES"
    size_violation_bug="BUG-XXX-SIZE-VIOLATION"

    # Create the bug
    create_size_violation_bug "$line_count" "$effort_id"
fi
```

### In Review Reports:

```markdown
## 📊 SIZE MEASUREMENT REPORT
**Implementation Lines:** 848
**Command:** tools/line-counter.sh effort-branch
**Timestamp:** 2025-11-02T14:30:00Z

## Size Analysis (Code Reviewer Enforcement Threshold)
- **Current Lines**: 848
- **SW Engineer Limit (stated)**: 800 lines
- **Code Reviewer Enforcement**: 900 lines
- **Status**: ❌ EXCEEDS ENFORCEMENT THRESHOLD (+48 lines)
- **Grace Buffer**: N/A (exceeded 900)
- **Requires Split**: YES
- **Bug Created**: BUG-XXX-SIZE-VIOLATION

**Explanation**: While this exceeds the 800-line target shown to SW Engineers, Code Reviewers enforce at 900 lines to reduce false positives. This effort exceeds the enforcement threshold.
```

### For 800-900 Line Range:

```markdown
## 📊 SIZE MEASUREMENT REPORT
**Implementation Lines:** 820
**Command:** tools/line-counter.sh effort-branch
**Timestamp:** 2025-11-02T14:30:00Z

## Size Analysis (Code Reviewer Enforcement Threshold)
- **Current Lines**: 820
- **SW Engineer Limit (stated)**: 800 lines
- **Code Reviewer Enforcement**: 900 lines
- **Status**: ✅ WITHIN ENFORCEMENT THRESHOLD
- **Grace Buffer**: ACTIVE (800-900 range)
- **Requires Split**: NO
- **Bug Created**: NONE

**Explanation**: While this exceeds the 800-line target shown to SW Engineers, Code Reviewers enforce at 900 lines to provide a grace buffer. This effort is within the enforcement threshold and requires no split.
```

## SW Engineer Perspective (UNCHANGED)

**SW Engineers continue to see 800-line limit in their documentation:**
- R007 still shows 800 lines
- SW Engineer agent config still shows 800 lines
- SW Engineer state rules still show 800 lines

**This strategic asymmetry is intentional:**
- Keeps developers conservative (aiming for <800)
- Provides natural safety buffer
- Reduces reviewer workload
- Eliminates false positive bugs

## Rationale for 900-Line Threshold

### Problems with 800-Line Enforcement:
1. **False Positives**: 810-line efforts create bugs unnecessarily
2. **Cascade Disruption**: Forced splits during integration for minor overages
3. **Reviewer Overhead**: Constant split planning for marginal violations
4. **Developer Frustration**: "I'm 10 lines over, really?"

### Benefits of 900-Line Threshold:
1. **Buffer Zone**: 100-line grace period catches minor overages
2. **Real Violations Only**: >900 lines is genuinely oversized
3. **Natural Discipline**: SW Engineers still aim for <800
4. **Reduced Noise**: Fewer unnecessary bugs in tracking

### Interaction with R339 Fix Grace Period:
- **R339**: Allows 900 lines during FIX_ISSUES state (first fix round)
- **R535**: Allows 900 lines for initial implementation (Code Reviewer enforcement)
- **Combined Effect**: Consistent 900-line enforcement threshold across contexts

## Implementation Checklist

### Code Reviewer Changes:
- ✅ Update `.claude/agents/code-reviewer.md` with 900-line threshold
- ✅ Update `agent-states/code-reviewer/EFFORT_PLAN_CREATION/rules.md`
- ✅ Update `agent-states/code-reviewer/CREATE_SPLIT_PLAN/rules.md`
- ✅ Update bug creation logic to use 900-line threshold
- ✅ Update review report templates

### SW Engineer (NO CHANGES):
- ❌ **DO NOT** update SW Engineer agent config
- ❌ **DO NOT** update SW Engineer state rules
- ❌ **DO NOT** mention 900-line threshold to SW Engineers

### Documentation:
- ✅ Add R535 to RULE-REGISTRY.md
- ✅ Cross-reference with R007, R339
- ✅ Note strategic asymmetry in architecture docs

## Related Rules

- **R007**: Size Limit Compliance (800 Line Maximum) - Base rule
- **R339**: Fix Grace Period Protocol - Similar 900-line threshold for fixes
- **R220**: Atomic PR Design Requirement - Foundation for size limits
- **R304**: Mandatory Line Counting Enforcement - Measurement requirements

## Common Questions

### Q: Won't SW Engineers notice the asymmetry?
**A:** SW Engineers don't see Code Reviewer thresholds. They see 800 lines, aim for <800, and rarely hit the grace buffer.

### Q: What if SW Engineer asks why 820 lines was approved?
**A:** Code Reviewer response: "Within acceptable variance for initial implementation. Focus on functionality, I'll flag true violations."

### Q: Does this violate transparency?
**A:** No. Different roles have different enforcement thresholds. SW Engineers have design constraints (800), Code Reviewers have enforcement thresholds (900). This is like "speed limit 65 mph" with enforcement at 75 mph.

### Q: Should we tell SW Engineers about the buffer?
**A:** **NO.** The buffer exists because developers aim for the stated limit. Announcing it would defeat the purpose.

## Enforcement

**ONLY Code Reviewers use 900-line threshold:**
- Bug creation: >900 lines
- Review approval: ≤900 lines
- Split requirements: >900 lines

**SW Engineers continue using 800-line limit:**
- Planning: Target <800 lines
- Self-monitoring: Warning at 700, stop at 800
- Checkpoint verification: Check against 800

---

**Remember:** This asymmetry is strategic. It maintains developer discipline while reducing false positive bugs. The 100-line grace buffer exists precisely because developers aim for 800.
