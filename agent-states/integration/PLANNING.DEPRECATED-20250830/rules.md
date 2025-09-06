# Integration Agent - PLANNING State Rules

## State Definition
The PLANNING state is where the integration agent analyzes branches and creates the integration plan.

## Required Actions

### 1. Branch Analysis
```bash
# Analyze branch relationships
for branch in "${TARGET_BRANCHES[@]}"; do
    echo "Analyzing $branch..."
    git log --oneline "$branch" -10
    git merge-base main "$branch"
done
```

### 2. Create INTEGRATION-PLAN.md
Must include:
- Target branch identification
- Branch integration order (by lineage)
- Conflict prediction
- Merge strategy
- Expected outcomes

### 3. Document Branch Dependencies
- Identify parent-child relationships
- Document common ancestors
- Plan overlay order

## Rules in Effect
- R261 - Integration Planning Requirements (BLOCKING)
- R260 - Integration Agent Core Requirements

## Transition Rules
- Can transition to: MERGING (only after plan complete)
- Cannot transition if: INTEGRATION-PLAN.md missing
- Must document all target branches

## Success Criteria
- INTEGRATION-PLAN.md created
- All branches analyzed for lineage
- Merge order determined
- Conflict risks assessed

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
