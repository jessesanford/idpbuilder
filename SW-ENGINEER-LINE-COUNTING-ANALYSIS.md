# SW Engineer Line Counting Analysis Report

## Executive Summary

This report analyzes the current state of line counting guidance for Software Engineers in the Software Factory 2.0 system. The analysis reveals significant gaps in measurement protocols that could lead to size violations and grading failures.

## Current State Assessment

### 1. Rules Identified

#### IMPLEMENTATION State
- **Location**: `agent-states/sw-engineer/IMPLEMENTATION/rules.md`
- **Clarity Level**: PARTIALLY CLEAR
- **Key Findings**:
  - References line counter tool but with legacy path `$PROJECT_ROOT/tools/line-counter.sh`
  - Mentions thresholds (700/800 lines) but lacks clear measurement frequency
  - Shows measurement in examples but doesn't mandate continuous monitoring
  - No clear base branch determination guidance

#### MEASURE_SIZE State
- **Location**: `agent-states/sw-engineer/MEASURE_SIZE/rules.md`
- **Clarity Level**: CLEAR AND COMPREHENSIVE
- **Key Findings**:
  - Excellent "universal line counter finder" script
  - Clear instructions for both normal efforts and splits
  - Detailed measurement scripts with compliance analysis
  - Good decision matrix based on results
  - BUT: This is AFTER size issues, not DURING implementation

#### SPLIT_IMPLEMENTATION State
- **Location**: `agent-states/sw-engineer/SPLIT_IMPLEMENTATION/rules.md`
- **Clarity Level**: MODERATE
- **Key Findings**:
  - References line counter but with simple usage
  - Focuses on split workflow, not measurement details
  - No clear base branch guidance for splits

## Critical Gaps Identified

### GAP 1: NO CONTINUOUS MONITORING PROTOCOL
**Issue**: SW Engineers lack clear guidance on WHEN to measure during implementation
**Impact**: Engineers may only discover size violations at completion
**Risk Level**: HIGH - Could result in -20% to -100% grading penalty

### GAP 2: UNCLEAR BASE BRANCH DETERMINATION
**Issue**: No clear guidance on what base branch to use for efforts vs splits
**Impact**: Wrong measurements leading to incorrect size assessments
**Risk Level**: CRITICAL - Could cause R304 violations (-100% failure)

### GAP 3: WRONG TOOL PATHS
**Issue**: References to legacy tool paths that don't exist in template
**Impact**: Engineers can't find measurement tool
**Risk Level**: MEDIUM - Causes delays and confusion

### GAP 4: NO SELF-MONITORING REQUIREMENTS
**Issue**: No mandate for engineers to proactively monitor their size
**Impact**: Reactive rather than preventive size management
**Risk Level**: HIGH - Leads to emergency splits and rework

### GAP 5: INCONSISTENT EXPECTATIONS
**Issue**: Code Reviewer (R304) has strict requirements that SW Engineers aren't aware of
**Impact**: Misalignment between development and review
**Risk Level**: CRITICAL - Systematic failures

## Comparison: SW Engineer vs Code Reviewer Expectations

| Aspect | SW Engineer Rules | Code Reviewer Rules (R304) | Gap |
|--------|------------------|---------------------------|-----|
| Tool Path | Legacy path referenced | PROJECT_ROOT/tools/line-counter.sh | MISALIGNED |
| Base Branch | Not specified clearly | MUST use phase integration, NEVER main | CRITICAL GAP |
| Measurement Frequency | Vague ("milestones") | Every review | UNCLEAR |
| Parameters | Sometimes shown | MANDATORY -b and -c | INCOMPLETE |
| Failure Conditions | Not emphasized | -100% for violations | NOT COMMUNICATED |

## Risk Assessment

### Scenario Analysis

#### Scenario A: Fresh Effort Implementation
**Current State**: SW Engineer starts E1.1.1
**Problems**:
- Doesn't know to use phase1/integration as base
- May not measure until completion
- Could exceed 800 lines before first measurement

#### Scenario B: Split Implementation
**Current State**: Working on split of large effort
**Problems**:
- Unclear what base to use for each split
- No guidance on incremental measurement
- Risk of splits exceeding limits

#### Scenario C: Fixing After Review
**Current State**: Adding fixes to existing effort
**Problems**:
- No guidance on re-measurement requirements
- Base branch confusion after merges
- Cumulative size not tracked

## Recommendations

### PRIORITY 1: Create R305 - SW Engineer Self-Monitoring Protocol

### PRIORITY 2: Update IMPLEMENTATION State Rules
Add clear continuous monitoring requirements:
- Measure at startup (baseline)
- Re-measure every 100 lines or 1 hour
- Stop at 700 lines for assessment
- Document all measurements

### PRIORITY 3: Fix Tool Path References
Update all references from legacy paths to:
```bash
$CLAUDE_PROJECT_DIR/tools/line-counter.sh
```

### PRIORITY 4: Add Base Branch Determination Guide
Clear instructions for:
- Efforts: Use phase integration branch
- Splits: Use original effort branch
- Fixes: Use current integration branch

### PRIORITY 5: Align with Code Reviewer Expectations
Ensure SW Engineers know about R304 requirements BEFORE development

## Proposed R305 Rule

```markdown
# Rule R305: SW Engineer Self-Monitoring Protocol

## Rule Statement
Software Engineers MUST continuously self-monitor code size during implementation using the line-counter.sh tool with correct parameters to prevent size violations.

## Criticality Level
**BLOCKING** - Failure to self-monitor results in grading penalties and potential automatic failure

## Monitoring Requirements

### 1. Baseline Measurement (START)
Upon entering IMPLEMENTATION state:
- Run line-counter.sh to establish baseline
- Document starting point in work-log.md
- Set monitoring schedule

### 2. Continuous Monitoring (DURING)
MUST measure:
- Every 100 lines added (estimated)
- Every 60 minutes of active coding
- Before any commit
- After completing any feature/component

### 3. Threshold Actions
- At 600 lines: Yellow alert - increase monitoring frequency
- At 700 lines: Red alert - assess completion feasibility
- At 750 lines: Critical - prepare for immediate completion
- At 800 lines: STOP - Automatic failure if exceeded

### 4. Base Branch Determination
MANDATORY base branches:
- For efforts: phase[N]/integration (from orchestrator-state.yaml)
- For splits: original effort branch (before split)
- For fixes: current integration branch
- NEVER use "main" or "master" as base

### 5. Correct Tool Usage
# Find tool
LC=$(find /home -name "line-counter.sh" -path "*/tools/*" 2>/dev/null | head -1)

# Determine base
BASE="phase1/integration"  # From orchestrator-state.yaml

# Measure
$LC -b $BASE -c $(git branch --show-current)

### 6. Documentation Requirements
Each measurement MUST be logged:
- Timestamp
- Current size
- Utilization percentage
- Growth rate
- Decision (continue/optimize/stop)

## Grading Impact
- No self-monitoring: -20% per violation
- Wrong base branch: -100% (R304 violation)
- Exceeding 800 lines: -100% (automatic failure)
- Poor monitoring discipline: -15% per incident
```

## Implementation Priority

1. **IMMEDIATE**: Update IMPLEMENTATION/rules.md with correct tool paths
2. **HIGH**: Create and publish R305
3. **HIGH**: Add base branch determination to all states
4. **MEDIUM**: Update examples with correct usage
5. **MEDIUM**: Create monitoring automation scripts

## Success Metrics

After implementation:
- SW Engineers measure proactively every 100 lines
- Zero efforts exceed 800 lines unexpectedly
- 100% alignment with Code Reviewer expectations
- Splits stay within limits from the start

## Conclusion

The current SW Engineer rules have critical gaps that set engineers up for failure. They lack:
1. Clear continuous monitoring requirements
2. Correct base branch determination
3. Awareness of R304 strict requirements
4. Proactive size management protocols

These gaps create systematic failures where engineers only discover size violations after the fact, leading to emergency splits, rework, and grading penalties.

**Recommendation**: Implement all five priorities immediately to ensure SW Engineers can successfully self-monitor and prevent size violations.

---

*Report Generated: 2025-09-02*
*Factory Manager: software-factory-manager*
*Status: CRITICAL GAPS REQUIRING IMMEDIATE ACTION*