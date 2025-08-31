# 🚨🚨🚨 RULE R258 - Mandatory Wave Review Report [BLOCKING]

## Rule Statement
**ARCHITECT MUST CREATE A PERMANENT WAVE REVIEW REPORT FILE BEFORE WAVE REVIEW COMPLETION**

The architect CANNOT signal wave review complete without:
1. Creating a standardized wave review report file
2. In the specified location with the exact naming convention
3. With all mandatory sections filled out
4. Verified by orchestrator before proceeding to next wave or phase assessment

## Why This Rule Exists

Previously, architects could provide wave review decisions through agent responses without creating any permanent record. This violated core principles:
- **No Audit Trail**: Wave decisions were not permanently documented
- **No Accountability**: No signed record of who approved what and when
- **No Verification**: Orchestrator couldn't verify review actually occurred
- **Lost Context**: Review details lost after conversation ended
- **No Metrics**: Integration quality and compliance data not preserved
- **No Decision Record**: Next steps unclear without documented decision

## Enforcement Mechanism

### 🚨 MANDATORY FILE REQUIREMENTS

**File Name**: `PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md`
- Example: `PHASE-3-WAVE-2-REVIEW-REPORT.md`
- Example: `PHASE-1-WAVE-3-REVIEW-REPORT.md`

**Location**: `wave-reviews/phase{N}/wave{W}/`
- Full path: `wave-reviews/phase3/wave2/PHASE-3-WAVE-2-REVIEW-REPORT.md`
- Directory must be created if it doesn't exist

**Creation Timing**: BEFORE signaling wave review complete
- Report MUST exist before architect returns decision
- Orchestrator MUST verify file exists and is valid
- No state transition without verified report

### 📋 MANDATORY REPORT STRUCTURE

```markdown
# Phase {N} Wave {W} Review Report

## Review Metadata
**Report ID**: WAVE-{N}-{W}-REVIEW-{YYYYMMDD}-{HHMMSS}
**Review Date**: {ISO-8601 timestamp}
**Reviewed By**: architect-reviewer
**Wave Branch**: {integration branch name}
**Efforts Reviewed**: {number of efforts in wave}

## Review Decision
**DECISION**: [PROCEED_NEXT_WAVE|PROCEED_PHASE_ASSESSMENT|CHANGES_REQUIRED|WAVE_FAILED]
**Confidence Level**: [HIGH|MEDIUM|LOW]
**Risk Assessment**: [LOW|MEDIUM|HIGH|CRITICAL]

## Integration Assessment
| Criterion | Status | Details |
|-----------|--------|---------|
| Branch Merges | [CLEAN|CONFLICTS_RESOLVED|ISSUES] | {details} |
| All Efforts Integrated | [YES|NO] | {count integrated}/{total} |
| Size Compliance | [ALL_UNDER_800|VIOLATIONS] | {violating efforts if any} |
| Test Results | [PASSING|FAILURES] | {pass rate %} |
| Build Status | [SUCCESS|FAILURE] | {details} |

## Architectural Review Scoring
| Category | Weight | Score (0-100) | Weighted Score |
|----------|--------|---------------|----------------|
| Pattern Compliance | 25% | {score} | {calculated} |
| API Consistency | 25% | {score} | {calculated} |
| Integration Quality | 20% | {score} | {calculated} |
| Code Quality | 15% | {score} | {calculated} |
| Test Coverage | 10% | {score} | {calculated} |
| Documentation | 5% | {score} | {calculated} |
| **OVERALL SCORE** | 100% | - | **{FINAL}** |

## Size Compliance Verification (Per R076)
| Effort | Measured Size | Status | Action Required |
|--------|---------------|--------|-----------------|
| {effort-name} | {lines} | [PASS|FAIL] | {split if >800} |
| {effort-name} | {lines} | [PASS|FAIL] | {split if >800} |

**Measurement Method Used**: `$PROJECT_ROOT/tools/line-counter.sh` (from effort directory)
**All Efforts Compliant**: [YES|NO]

## Wave Deliverables Checklist
- [ ] All planned efforts completed
- [ ] All code reviews passed
- [ ] All efforts ≤800 lines (verified)
- [ ] Integration branch created
- [ ] Integration branch tested
- [ ] All tests passing (>95% success)
- [ ] Documentation updated (work logs complete)
- [ ] No security vulnerabilities introduced
- [ ] Performance within acceptable limits

## Issues Identified
{Only if CHANGES_REQUIRED or WAVE_FAILED}

### Critical Issues (Must Fix)
1. **[CRITICAL]**: {Issue description}
   - Impact: {description}
   - Required Fix: {specific action}
   - Assigned To: {sw-engineer|code-reviewer}

### High Priority Issues (Should Fix)
1. **[HIGH]**: {Issue description}
   - Impact: {description}
   - Required Fix: {specific action}

### Medium Priority Issues (Consider Fixing)
1. **[MEDIUM]**: {Issue description}
   - Impact: {description}
   - Suggested Fix: {action}

## Required Actions
{Based on DECISION field}

### If PROCEED_NEXT_WAVE:
- Next Wave Number: Wave {W+1}
- Estimated Efforts: {count}
- Dependencies: {any prerequisites}

### If PROCEED_PHASE_ASSESSMENT:
- Phase ready for assessment
- All waves integrated
- Trigger phase assessment protocol

### If CHANGES_REQUIRED:
- Fix issues listed above
- Re-run wave review after fixes
- Do not proceed to next wave

### If WAVE_FAILED:
- Major architectural redesign required
- Return to planning state
- Consult with orchestrator for recovery

## Recommendations
- {Improvements for next wave}
- {Architecture suggestions}
- {Process improvements}

## Sign-Off
**Wave {W} Review Decision**: This wave [IS|IS NOT] ready to proceed.

**Reason**: {Brief explanation of decision}

**Architect Sign-Off**: {ISO-8601 timestamp}
**Report Hash**: {SHA256 of report content for integrity}
```

## Implementation Requirements

### For Architect (WAVE_REVIEW State)

```bash
# MANDATORY: Create wave review report
create_wave_review_report() {
    local PHASE=$1
    local WAVE=$2
    local DECISION=$3
    local REPORT_DIR="wave-reviews/phase${PHASE}/wave${WAVE}"
    local REPORT_FILE="${REPORT_DIR}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"
    
    # Create directory structure if needed
    mkdir -p "$REPORT_DIR"
    
    # Generate report with all mandatory sections
    cat > "$REPORT_FILE" << 'EOF'
# Phase ${PHASE} Wave ${WAVE} Review Report
[... complete report content with all sections ...]
EOF
    
    # Calculate and add hash for integrity
    local HASH=$(sha256sum "$REPORT_FILE" | cut -d' ' -f1)
    echo "**Report Hash**: $HASH" >> "$REPORT_FILE"
    
    # Commit and push immediately (per R288)
    git add "$REPORT_FILE"
    git commit -m "review: Phase ${PHASE} Wave ${WAVE} review report - ${DECISION}"
    git push
    
    echo "✅ Wave review report created: $REPORT_FILE"
    return 0
}

# MUST call this before signaling complete
create_wave_review_report "$CURRENT_PHASE" "$CURRENT_WAVE" "$DECISION"
```

### For Orchestrator (WAITING_FOR_WAVE_REVIEW State)

```bash
# MANDATORY: Verify wave review report exists
verify_wave_review_report() {
    local PHASE=$1
    local WAVE=$2
    local REPORT_FILE="wave-reviews/phase${PHASE}/wave${WAVE}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"
    
    # Check file exists
    if [ ! -f "$REPORT_FILE" ]; then
        echo "❌ CRITICAL: No wave review report found!"
        echo "❌ Expected: $REPORT_FILE"
        echo "❌ Cannot proceed without wave review report (R258 violation)"
        transition_to "ERROR_RECOVERY"
        exit 1
    fi
    
    # Extract and verify decision
    local DECISION=$(grep "^\*\*DECISION\*\*:" "$REPORT_FILE" | cut -d: -f2 | xargs)
    
    case "$DECISION" in
        PROCEED_NEXT_WAVE|PROCEED_PHASE_ASSESSMENT|CHANGES_REQUIRED|WAVE_FAILED)
            echo "✅ Valid decision found: $DECISION"
            ;;
        *)
            echo "❌ Invalid decision in report: $DECISION"
            echo "❌ Must be one of: PROCEED_NEXT_WAVE, PROCEED_PHASE_ASSESSMENT, CHANGES_REQUIRED, WAVE_FAILED"
            transition_to "ERROR_RECOVERY"
            exit 1
            ;;
    esac
    
    # Extract overall score
    local SCORE=$(grep "^\*\*OVERALL SCORE\*\*" "$REPORT_FILE" | grep -o "[0-9]\+" | tail -1)
    
    # Extract size compliance status
    local SIZE_COMPLIANT=$(grep "^\*\*All Efforts Compliant\*\*:" "$REPORT_FILE" | grep -o "\(YES\|NO\)")
    
    # Update state file with report details
    yq -i ".wave_review.report_file = \"$REPORT_FILE\"" orchestrator-state.yaml
    yq -i ".wave_review.decision = \"$DECISION\"" orchestrator-state.yaml
    yq -i ".wave_review.score = $SCORE" orchestrator-state.yaml
    yq -i ".wave_review.size_compliant = \"$SIZE_COMPLIANT\"" orchestrator-state.yaml
    yq -i ".wave_review.timestamp = \"$(date -Iseconds)\"" orchestrator-state.yaml
    
    echo "✅ Wave review report verified:"
    echo "  📄 Report: $REPORT_FILE"
    echo "  📊 Decision: $DECISION"
    echo "  📈 Score: $SCORE"
    echo "  📏 Size Compliance: $SIZE_COMPLIANT"
    
    return 0
}

# MUST call this before processing architect decision
verify_wave_review_report "$CURRENT_PHASE" "$CURRENT_WAVE"

# Process decision based on report
process_wave_review_decision() {
    local DECISION=$(yq '.wave_review.decision' orchestrator-state.yaml)
    
    case "$DECISION" in
        PROCEED_NEXT_WAVE)
            transition_to "WAVE_START"
            increment_wave_counter
            ;;
        PROCEED_PHASE_ASSESSMENT)
            transition_to "SPAWN_ARCHITECT_PHASE_ASSESSMENT"
            ;;
        CHANGES_REQUIRED)
            transition_to "PROCESS_WAVE_FIXES"
            ;;
        WAVE_FAILED)
            transition_to "WAVE_RECOVERY"
            ;;
    esac
}
```

## Decision Types Explained

### PROCEED_NEXT_WAVE
- Wave fully approved with no blocking issues
- All efforts integrated successfully
- Ready to begin next wave of implementation
- Orchestrator should increment wave counter and start planning

### PROCEED_PHASE_ASSESSMENT
- Current wave was the last wave in phase
- All waves in phase have been integrated
- Ready for comprehensive phase assessment
- Orchestrator should trigger phase assessment protocol (R256, R257)

### CHANGES_REQUIRED
- Wave has fixable issues that must be addressed
- Specific actions identified in report
- Must fix before proceeding to next wave
- Orchestrator should coordinate fixes with appropriate agents

### WAVE_FAILED
- Wave has fundamental architectural problems
- Cannot proceed without major redesign
- Requires return to planning/design phase
- Orchestrator should initiate recovery protocol

## Validation Rules

### Report Validation Criteria
1. **File exists** at exact expected path
2. **All mandatory sections** present (not empty)
3. **Valid DECISION** value from allowed set
4. **Scores provided** for all architectural categories
5. **Size compliance** section completed with measurements
6. **Deliverables checklist** evaluated
7. **Sign-off section** completed with timestamp
8. **Report hash** present for integrity

### Orchestrator Validation Steps
```bash
# Comprehensive validation function
validate_wave_review_report() {
    local REPORT="$1"
    
    # Check mandatory sections exist
    local REQUIRED_SECTIONS=(
        "Review Decision"
        "Integration Assessment"
        "Architectural Review Scoring"
        "Size Compliance Verification"
        "Wave Deliverables Checklist"
        "Sign-Off"
    )
    
    for section in "${REQUIRED_SECTIONS[@]}"; do
        grep -q "^## $section" "$REPORT" || {
            echo "❌ Missing mandatory section: $section"
            return 1
        }
    done
    
    # Verify decision is valid
    grep -q "^\*\*DECISION\*\*: \(PROCEED_NEXT_WAVE\|PROCEED_PHASE_ASSESSMENT\|CHANGES_REQUIRED\|WAVE_FAILED\)" "$REPORT" || {
        echo "❌ Invalid or missing DECISION"
        return 1
    }
    
    # Verify size compliance is documented
    grep -q "^\*\*All Efforts Compliant\*\*: \(YES\|NO\)" "$REPORT" || {
        echo "❌ Missing size compliance status"
        return 1
    }
    
    # Verify sign-off exists
    grep -q "^\*\*Architect Sign-Off\*\*:" "$REPORT" || {
        echo "❌ Missing architect sign-off"
        return 1
    }
    
    # Verify report hash exists
    grep -q "^\*\*Report Hash\*\*:" "$REPORT" || {
        echo "❌ Missing report hash for integrity"
        return 1
    }
    
    echo "✅ Wave review report validation passed"
    return 0
}
```

## State Machine Integration

### State Transitions Based on Report
```yaml
WAVE_REVIEW:
  transitions:
    - trigger: report_decision_PROCEED_NEXT_WAVE
      target: WAVE_START
      action: increment_wave_counter
    
    - trigger: report_decision_PROCEED_PHASE_ASSESSMENT
      target: SPAWN_ARCHITECT_PHASE_ASSESSMENT
      action: prepare_phase_assessment
    
    - trigger: report_decision_CHANGES_REQUIRED
      target: PROCESS_WAVE_FIXES
      action: extract_required_fixes
    
    - trigger: report_decision_WAVE_FAILED
      target: WAVE_RECOVERY
      action: initiate_recovery_protocol
```

## Violations

**Severity**: BLOCKING - Wave cannot complete without report

**Detection**:
- Architect signals wave review complete without creating report file
- Report file missing mandatory sections
- Invalid decision value in report
- Orchestrator proceeding without verifying report
- Report not committed to git repository
- Size compliance section incomplete or missing

**Consequences**:
- Immediate transition to ERROR_RECOVERY state
- Wave marked incomplete
- Cannot proceed to next wave or phase
- Architect must create proper report
- Orchestrator blocked from state progression
- Grading: -50 points (critical violation)

## Related Rules
- R074 - Wave Completion Architectural Review (requires review)
- R076 - Effort Size Compliance Verification (≤800 lines)
- R022 - Architect Size Verification Protocol (measurement method)
- R288 - Mandatory State File Updates (state tracking)
- R288 - Mandatory State File Commit Push (persistence)
- R256 - Mandatory Phase Assessment Gate (phase progression)
- R257 - Mandatory Phase Assessment Report (phase reporting)

## Audit Commands

```bash
# Verify all completed waves have review reports
find wave-reviews -type d -name "wave*" | while read wave_dir; do
    phase=$(echo "$wave_dir" | grep -o "phase[0-9]*" | grep -o "[0-9]*")
    wave=$(echo "$wave_dir" | grep -o "wave[0-9]*" | grep -o "[0-9]*")
    report_file="${wave_dir}/PHASE-${phase}-WAVE-${wave}-REVIEW-REPORT.md"
    
    if [ ! -f "$report_file" ]; then
        echo "❌ VIOLATION: Phase $phase Wave $wave missing review report!"
    else
        decision=$(grep "^\*\*DECISION\*\*:" "$report_file" | cut -d: -f2 | xargs)
        echo "✅ Phase $phase Wave $wave has review report (Decision: $decision)"
    fi
done

# Check for waves marked complete without reports
yq '.efforts_completed[] | select(.wave != null)' orchestrator-state.yaml | while read effort; do
    wave=$(echo "$effort" | yq '.wave')
    phase=$(echo "$effort" | yq '.phase')
    report="wave-reviews/phase${phase}/wave${wave}/PHASE-${phase}-WAVE-${wave}-REVIEW-REPORT.md"
    
    if [ ! -f "$report" ]; then
        echo "⚠️ WARNING: Phase $phase Wave $wave completed but no review report"
    fi
done

# Verify report decisions match state transitions
for report in wave-reviews/*/*/PHASE-*-WAVE-*-REVIEW-REPORT.md; do
    decision=$(grep "^\*\*DECISION\*\*:" "$report" | cut -d: -f2 | xargs)
    phase=$(basename "$report" | grep -o "PHASE-[0-9]*" | grep -o "[0-9]*")
    wave=$(basename "$report" | grep -o "WAVE-[0-9]*" | grep -o "[0-9]*")
    
    echo "Phase $phase Wave $wave: $decision"
    
    # Verify appropriate follow-up based on decision
    case "$decision" in
        PROCEED_NEXT_WAVE)
            next_wave=$((wave + 1))
            if [ -d "wave-reviews/phase${phase}/wave${next_wave}" ]; then
                echo "  ✅ Next wave started as expected"
            fi
            ;;
        PROCEED_PHASE_ASSESSMENT)
            if [ -f "phase-assessments/phase${phase}/PHASE-${phase}-ASSESSMENT-REPORT.md" ]; then
                echo "  ✅ Phase assessment conducted as expected"
            fi
            ;;
    esac
done
```

## Utility Script

Create `utilities/verify-wave-review-report.sh`:
```bash
#!/bin/bash
# Verify wave review report compliance with R258

PHASE=$1
WAVE=$2

if [ -z "$PHASE" ] || [ -z "$WAVE" ]; then
    echo "Usage: $0 <phase> <wave>"
    exit 1
fi

REPORT_FILE="wave-reviews/phase${PHASE}/wave${WAVE}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"

# Function to check report
check_wave_review_report() {
    if [ ! -f "$REPORT_FILE" ]; then
        echo "❌ CRITICAL: Wave review report not found!"
        echo "Expected: $REPORT_FILE"
        echo "This violates R258 - Mandatory Wave Review Report"
        return 1
    fi
    
    # Extract key fields
    local DECISION=$(grep "^\*\*DECISION\*\*:" "$REPORT_FILE" | cut -d: -f2 | xargs)
    local SCORE=$(grep "^\*\*OVERALL SCORE\*\*" "$REPORT_FILE" | grep -o "[0-9]\+" | tail -1)
    local SIZE_OK=$(grep "^\*\*All Efforts Compliant\*\*:" "$REPORT_FILE" | grep -o "\(YES\|NO\)")
    local SIGNOFF=$(grep "^\*\*Architect Sign-Off\*\*:" "$REPORT_FILE")
    
    # Display report summary
    echo "📊 Wave Review Report Summary:"
    echo "  Phase: $PHASE, Wave: $WAVE"
    echo "  Decision: $DECISION"
    echo "  Score: $SCORE"
    echo "  Size Compliance: $SIZE_OK"
    
    # Validate decision
    case "$DECISION" in
        PROCEED_NEXT_WAVE|PROCEED_PHASE_ASSESSMENT|CHANGES_REQUIRED|WAVE_FAILED)
            echo "  ✅ Valid decision type"
            ;;
        *)
            echo "  ❌ Invalid decision: $DECISION"
            return 1
            ;;
    esac
    
    # Check critical fields
    if [ -z "$SIGNOFF" ]; then
        echo "  ❌ Missing architect sign-off"
        return 1
    fi
    
    if [ "$SIZE_OK" = "NO" ] && [ "$DECISION" = "PROCEED_NEXT_WAVE" ]; then
        echo "  ❌ Cannot proceed with size violations!"
        return 1
    fi
    
    echo "✅ Wave review report is compliant with R258"
    return 0
}

# Run check
check_wave_review_report
exit $?
```

## Summary

**The Wave Review Report is MANDATORY and BLOCKING**:
- Every wave review MUST produce a permanent report file
- The report MUST follow the exact structure and naming convention
- The report MUST contain one of four valid decisions
- Orchestrator MUST verify the report before any state transition
- Decision field determines the next state in the state machine
- No verbal reviews allowed - everything must be documented
- Creates permanent audit trail for all wave-level decisions

This ensures accountability, traceability, and verifiability of all wave review decisions. NO WAVE CAN PROCEED WITHOUT A PROPER REVIEW REPORT.