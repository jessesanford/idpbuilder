# 🚨🚨🚨 RULE R261: Integration Planning Requirements 🚨🚨🚨

## Rule Definition
**Criticality:** BLOCKING
**Category:** Agent-Specific
**Applies To:** integration-agent

## Requirements

### 1. Mandatory Plan Creation
The integration agent MUST:
- **ALWAYS** create an INTEGRATION-PLAN.md BEFORE any merge operations
- Document target branches and their relationships
- Identify merge order based on git commit lineage
- Plan for conflict resolution strategy
- Document expected outcomes

### 2. Branch Ordering by Lineage
The plan MUST:
- Order branches by git commit lineage
- Respect base branch dependencies
- Group related branches together
- Maximize linear story preservation
- Document overlay order to reduce collisions

### 3. Collision Prediction
The plan MUST identify:
- Branches with divergent histories
- Potential merge conflicts
- Files touched by multiple branches
- Strategies to minimize conflicts
- Alternative merge paths if conflicts arise

### 4. Naming Strategy for Synthesis
When recommending synthesized branches:
- Preserve original branch PREFIX
- Use clear suffix indicating synthesis (e.g., -integrated, -merged)
- Document relationship to original branches
- Never modify naming of original branches

## Plan Template

```markdown
# INTEGRATION-PLAN.md

## Target Branch
- Branch: main (or specified target)
- Current HEAD: [commit SHA]

## Branches to Integrate (in order)
1. branch-name-1
   - Based on: [parent branch]
   - Divergence point: [commit SHA]
   - Files modified: [count]
   - Conflict risk: LOW/MEDIUM/HIGH
   
2. branch-name-2
   - [same structure]

## Merge Strategy
- Order rationale: [explain lineage-based ordering]
- Conflict mitigation: [strategies]
- Synthesis requirements: [if needed]

## Expected Outcome
- Final branch state
- Test requirements
- Build validation needed
```

## Enforcement

```bash
# Verify plan exists before merging
verify_integration_plan() {
    local effort_dir="$1"
    
    if [[ ! -f "$effort_dir/INTEGRATION-PLAN.md" ]]; then
        echo "❌ BLOCKING: No INTEGRATION-PLAN.md found!"
        return 1
    fi
    
    # Check plan completeness
    grep -q "Target Branch" "$effort_dir/INTEGRATION-PLAN.md" || echo "❌ Missing target"
    grep -q "Branches to Integrate" "$effort_dir/INTEGRATION-PLAN.md" || echo "❌ Missing branch list"
    grep -q "Merge Strategy" "$effort_dir/INTEGRATION-PLAN.md" || echo "❌ Missing strategy"
}
```

## Grading Impact
- 15% - Plan completeness and accuracy
- 10% - Correct branch ordering by lineage
- 5% - Collision prediction accuracy

## Related Rules
- R260 - Integration Agent Core Requirements
- R262 - Merge Operation Protocols
- R264 - Work Log Tracking Requirements