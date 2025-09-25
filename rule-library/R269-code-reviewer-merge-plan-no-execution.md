# ⚠️⚠️⚠️ RULE R269: Code Reviewer Merge Plan No Execution ⚠️⚠️⚠️

## Rule Definition
**Criticality:** WARNING - Violation = Role confusion
**Category:** Agent-Specific
**Applies To:** code-reviewer
**Created:** 2025-08-27

## Core Requirement

**CODE REVIEWER CREATES MERGE PLANS ONLY - NEVER EXECUTES MERGES**

The Code Reviewer agent in WAVE_MERGE_PLANNING or PHASE_MERGE_PLANNING states:
- ✅ Creates detailed MERGE PLAN documents
- ✅ Analyzes branch dependencies and order
- ✅ Identifies potential conflicts
- ✅ Documents merge strategies
- ❌ NEVER executes git merge commands
- ❌ NEVER modifies any branches
- ❌ NEVER creates integration branches

## Role Separation

### Code Reviewer Responsibilities (Planning Only)
1. Analyze all branches to be integrated
2. Determine correct merge order based on git lineage
3. Identify split branches vs original branches
4. Create comprehensive merge plan document
5. Include validation steps in the plan
6. Document expected conflicts and resolutions

### Integration Agent Responsibilities (Execution Only)
1. Read the merge plan created by Code Reviewer
2. Execute merges in specified order
3. Resolve conflicts as they occur
4. Document all operations in work-log.md
5. Run tests and validation
6. Create integration report

## Enforcement

```bash
# Check Code Reviewer hasn't executed merges
verify_code_reviewer_compliance() {
    local agent_log="$1"
    
    # Code Reviewer should NOT have these
    if grep -q "git merge" "$agent_log" | grep -v "# "; then
        echo "❌ VIOLATION: Code Reviewer executed git merge!"
        return 1
    fi
    
    if grep -q "Resolving conflict" "$agent_log"; then
        echo "❌ VIOLATION: Code Reviewer resolved conflicts!"
        return 1
    fi
    
    # Code Reviewer SHOULD have these
    if ! grep -q "MERGE-PLAN.md" "$agent_log"; then
        echo "⚠️ WARNING: No merge plan document created"
    fi
    
    echo "✅ Code Reviewer compliance verified"
    return 0
}
```

## Example Compliant Behavior

### ✅ CORRECT: Code Reviewer Creates Plan
```markdown
# In WAVE-MERGE-PLAN.md

## Branches to Merge (IN ORDER)
1. phase3/wave1/effort1-api-types
   - Base: main at abc123
   - Merge Command:
   ```bash
   git merge origin/phase3/wave1/effort1-api-types --no-ff
   ```
```

### ❌ WRONG: Code Reviewer Executes Merge
```bash
# Code Reviewer should NEVER do this!
git merge origin/phase3/wave1/effort1-api-types --no-ff
```

## Rationale

This separation ensures:
1. **Clear accountability** - Each agent has distinct responsibilities
2. **Better documentation** - Planning and execution are separate
3. **Easier debugging** - Issues can be traced to planning or execution
4. **Replayability** - Integration Agent can re-execute from plan
5. **Review opportunity** - Plan can be reviewed before execution

## Related Rules
- R260 - Integration Agent Core Requirements
- R261 - Integration Planning Requirements (Integration Agent)
- R262 - Merge Operation Protocols
- R270 - No Integration Branches as Sources

## Grading Impact
- Code Reviewer loses 25% for executing merges
- Code Reviewer gains 10% for comprehensive merge plan
- Integration Agent loses 25% for not following plan