# Integration Agent - REPORTING State Rules

## State Definition
The REPORTING state completes all documentation and prepares final deliverables.

## Required Actions

### 1. Complete INTEGRATION-REPORT.md
Must include ALL sections:
- Overview (branches integrated, statistics)
- Errors and Issues Found
- Compensating/Remediation Recommendations  
- Build and Test Results
- Upstream Bugs Identified
- Integration Verification Checklist
- Final State

### 2. Verify Work Log Completeness
```bash
# Ensure work-log is replayable
grep "^Command:" work-log.md > replay.sh
bash -n replay.sh  # Verify syntax

# Count operations
OPERATION_COUNT=$(grep -c "^## Operation" work-log.md)
echo "Total operations documented: $OPERATION_COUNT"
```

### 3. Commit Documentation
```bash
# Add all documentation to integration branch
git add INTEGRATION-PLAN.md work-log.md INTEGRATION-REPORT.md
git commit -m "docs: complete integration documentation for [branch list]"
git push origin "$INTEGRATION_BRANCH"
```

## Documentation Quality Rules
- R263 - Integration Documentation Requirements
- R264 - Work Log Tracking Requirements

## Final Checklist
Before transition to COMPLETED:
- [ ] INTEGRATION-PLAN.md exists and was followed
- [ ] work-log.md is complete and replayable
- [ ] INTEGRATION-REPORT.md has all sections
- [ ] No original branches were modified
- [ ] No cherry-picks were used
- [ ] All documentation committed and pushed

## Transition Rules
- Can transition to: COMPLETED
- Cannot transition if: Documentation incomplete
- Must have pushed integration branch

## Success Criteria
- All three documents complete
- Documentation committed to integration branch
- Integration branch pushed to remote
- Ready for external review

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**
