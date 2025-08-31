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