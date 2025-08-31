# Integration Agent - MERGING State Rules

## State Definition
The MERGING state is where actual branch integration occurs following the MERGE PLAN created by Code Reviewer.

## Required Actions

### 1. Verify Integration Branch Already Created (R014 Compliance)
```bash
# The orchestrator already created the integration branch
# We should already be on it from INIT state
CURRENT_BRANCH=$(git branch --show-current)
if [[ ! "$CURRENT_BRANCH" =~ integration ]]; then
    echo "❌ ERROR: Not on integration branch!"
    echo "Current branch: $CURRENT_BRANCH"
    exit 1
fi

# Verify branch follows R014 naming convention with project prefix
PROJECT_PREFIX=$(yq '.branch_naming.project_prefix' "$SF_INSTANCE_DIR/target-repo-config.yaml" 2>/dev/null || echo "")
if [ -n "$PROJECT_PREFIX" ] && [ "$PROJECT_PREFIX" != "null" ]; then
    if [[ ! "$CURRENT_BRANCH" =~ ^"$PROJECT_PREFIX"/ ]]; then
        echo "⚠️ WARNING: Integration branch missing project prefix: $PROJECT_PREFIX"
        echo "Branch should start with: $PROJECT_PREFIX/"
    fi
fi

echo "✅ On integration branch: $CURRENT_BRANCH"
```

### 2. Sequential Merging Following MERGE PLAN
```bash
# Read merge plan to extract branches
if [ -f "WAVE-MERGE-PLAN.md" ]; then
    MERGE_PLAN="WAVE-MERGE-PLAN.md"
elif [ -f "PHASE-MERGE-PLAN.md" ]; then
    MERGE_PLAN="PHASE-MERGE-PLAN.md"
else
    echo "❌ NO MERGE PLAN FOUND!"
    exit 1
fi

echo "📋 Following merge plan: $MERGE_PLAN"
```

**CRITICAL: Execute merges in EXACT order from plan**
- Extract merge commands from the plan
- Execute each merge with --no-ff
- Resolve conflicts as encountered
- Document EVERY operation in work-log.md
- Run tests after each merge (if specified in plan)

### 3. Conflict Resolution
- Resolve ALL conflicts completely
- Document resolution decisions
- Never leave conflict markers
- Test compilation after resolution

## SUPREME LAWS IN EFFECT
- R262 - Merge Operation Protocols
  - NEVER modify original branches
  - NEVER use cherry-pick
  - Create synthesis branches if needed

## Work Log Requirements
Every operation MUST be documented:
```markdown
## Operation N: [timestamp]
Command: [exact command]
Result: [output]
Conflicts: [if any]
Resolution: [how resolved]
```

## Transition Rules
- Can transition to: TESTING (after all merges complete)
- Cannot transition if: Conflicts unresolved
- Must have clean working tree

## Success Criteria
- All planned branches merged
- All conflicts resolved
- work-log.md contains every operation
- No original branches modified