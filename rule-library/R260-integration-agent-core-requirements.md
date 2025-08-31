# 🚨🚨🚨 RULE R260: Integration Agent Core Requirements 🚨🚨🚨

## Rule Definition
**Criticality:** BLOCKING
**Category:** Agent-Specific
**Applies To:** integration-agent

## Requirements

### 1. Git Operations Expertise
The integration agent MUST be an expert in:
- ALL git operations (merge, rebase, log, diff, branch)
- Trunk-based development patterns
- Feature branch management
- Git lineage and commit history analysis
- Branch dependency tracking

### 2. Branch Relationship Understanding
The agent MUST understand:
- How branches are BASED on other branches
- Parent-child branch relationships  
- Commit lineage and ancestry
- Branch divergence points
- Common ancestor identification

### 3. Branch Splitting Protocol Recognition
The agent MUST recognize:
- "TOO LARGE" branches that exceed 800 lines
- Split branches that subsume original branches
- Split branch naming conventions (prefix preservation)
- Split branch hierarchy and ordering

### 4. Commit History Preservation
The agent MUST:
- Preserve "commit-to-commit trails of intent"
- Maintain linear story where possible
- Never lose commit messages or metadata
- Preserve author information and timestamps
- Document commit relationships in work-log

## Enforcement

```bash
# Verify git expertise
verify_git_expertise() {
    local agent_log="$1"
    
    # Check for proper git command usage
    grep -q "git log --graph" "$agent_log" || echo "❌ Missing graph analysis"
    grep -q "git merge-base" "$agent_log" || echo "❌ Missing common ancestor check"
    grep -q "git branch --contains" "$agent_log" || echo "❌ Missing branch containment check"
}

# Verify branch relationship understanding
verify_branch_relationships() {
    local work_log="$1"
    
    # Must document branch relationships
    grep -q "Based on:" "$work_log" || echo "❌ Missing base branch documentation"
    grep -q "Parent branch:" "$work_log" || echo "❌ Missing parent documentation"
    grep -q "Divergence point:" "$work_log" || echo "❌ Missing divergence analysis"
}
```

## Grading Impact
- 10% - Git operations expertise demonstrated
- 10% - Branch relationship documentation
- 10% - Split protocol recognition
- 10% - Commit history preservation

## Related Rules
- R261 - Integration Planning Requirements
- R262 - Merge Operation Protocols
- R263 - Integration Documentation Requirements