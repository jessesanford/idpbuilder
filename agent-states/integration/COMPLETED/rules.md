# Integration Agent - COMPLETED State Rules

## State Definition
The COMPLETED state indicates successful completion of all integration tasks.

## Final Validation

### 1. Self-Assessment
```bash
echo "=== INTEGRATION COMPLETE ==="
echo "Integration Branch: $INTEGRATION_BRANCH"
echo "Branches Integrated: ${#MERGED_BRANCHES[@]}"
echo "Conflicts Resolved: $CONFLICT_COUNT"
echo "Documentation Complete: YES"
```

### 2. Grading Self-Evaluation
Review against R267 criteria:
- Completeness of Integration (50%)
  - Branch merging: ___/20%
  - Conflict resolution: ___/15%
  - Branch integrity: ___/10%
  - Final validation: ___/5%
- Documentation Quality (50%)
  - Work log: ___/25%
  - Integration report: ___/25%

### 3. Final Verification
```bash
# Verify no violations
git log --grep="cherry picked" && echo "❌ VIOLATION" || echo "✅ No cherry-picks"

# Check documentation
for doc in INTEGRATION-PLAN.md work-log.md INTEGRATION-REPORT.md; do
    [[ -f "$doc" ]] && echo "✅ $doc exists" || echo "❌ Missing $doc"
done
```

## Success Indicators
- Integration branch exists and is pushed
- All documentation complete and committed
- No rule violations detected
- Ready for orchestrator review

## No Further Transitions
This is a terminal state.

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
