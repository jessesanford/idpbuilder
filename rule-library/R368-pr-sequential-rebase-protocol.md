# 🔴🔴🔴 SUPREME LAW RULE R368 - PR Sequential Rebase Protocol

**Criticality:** SUPREME LAW - Incorrect rebase order = merge conflicts
**Enforcement:** MANDATORY - All PR-Ready operations
**Created:** 2025-01-21

## PURPOSE
Ensure all effort branches are rebased onto main in the correct dependency order to prevent merge conflicts and maintain clean, sequential mergeability as required by R363.

## REBASE PRINCIPLES

### Sequential Mergeability (R363)
- Every effort MUST merge directly to main
- Efforts MUST be mergeable in sequence
- NO integration branch intermediaries
- Dependencies determine sequence

### Dependency Order Rules
1. **Infrastructure First** - Core modules before features
2. **Dependencies Before Dependents** - Base before extensions
3. **Independent Parallel** - No deps = any order
4. **UI/Frontend Last** - After all backend ready

## REBASE PROTOCOL

### Step 1: Determine Dependency Order
```bash
# Extract dependency information from metadata
echo "📊 Analyzing effort dependencies..."

# Check for dependency documentation
if [ -f ".software-factory/effort-dependencies.json" ]; then
    DEPENDENCY_ORDER=$(jq -r '.rebase_order[]' .software-factory/effort-dependencies.json)
else
    # Fallback: analyze import statements
    echo "⚠️ No dependency file, analyzing code..."

    # Find efforts that don't depend on others (infrastructure)
    INFRASTRUCTURE=$(find efforts/ -name "*.md" | xargs grep -l "infrastructure\|core\|base")

    # Find efforts that depend on infrastructure
    DEPENDENT=$(find efforts/ -name "*.md" | xargs grep -l "depends on\|requires")

    # Build order
    DEPENDENCY_ORDER="$INFRASTRUCTURE $DEPENDENT"
fi

echo "Rebase order determined:"
echo "$DEPENDENCY_ORDER" | nl
```

### Step 2: Update Main Reference
```bash
# Ensure we have latest main
echo "📥 Updating main branch..."

git fetch origin main:main --force
MAIN_HEAD=$(git rev-parse main)
echo "Main is at: $MAIN_HEAD"
```

### Step 3: Sequential Rebase Process
```bash
# Process each effort in dependency order
REBASED_BRANCHES=""

for EFFORT in $DEPENDENCY_ORDER; do
    echo "🔄 Processing: $EFFORT"

    # Checkout effort branch
    git checkout $EFFORT

    # Record pre-rebase state
    PRE_REBASE=$(git rev-parse HEAD)

    # Perform rebase
    echo "Rebasing onto main..."
    git rebase main

    if [ $? -ne 0 ]; then
        echo "⚠️ Conflicts detected in $EFFORT"

        # Auto-resolve Software Factory conflicts
        # (SF files should not exist in main)
        git status --porcelain | grep -E "(orchestrator-state|\.todo|PLAN\.md)" | while read -r conflict; do
            FILE=$(echo $conflict | awk '{print $2}')
            echo "Auto-resolving SF artifact: $FILE"
            git rm "$FILE"
        done

        # Check if manual resolution needed
        if [ -n "$(git diff --name-only --diff-filter=U)" ]; then
            echo "❌ Manual conflict resolution required!"
            cat > REBASE-CONFLICT-$EFFORT.md << EOF
# Rebase Conflict Report
Effort: $EFFORT
Pre-rebase: $PRE_REBASE
Conflicts requiring manual resolution:
$(git diff --name-only --diff-filter=U)

## Resolution Instructions
1. Resolve conflicts manually
2. git add resolved files
3. git rebase --continue
4. Re-run sequential rebase
EOF
            git rebase --abort
            exit 1
        fi

        # Continue rebase if auto-resolved
        git rebase --continue
    fi

    # Verify rebase success
    POST_REBASE=$(git rev-parse HEAD)
    echo "✅ Rebased: $PRE_REBASE -> $POST_REBASE"

    # Update subsequent efforts to rebase on this
    REBASED_BRANCHES="$REBASED_BRANCHES $EFFORT"

    # Push rebased branch
    echo "Pushing rebased $EFFORT..."
    git push origin $EFFORT --force-with-lease
done
```

### Step 4: Cross-Verify Sequential Mergeability
```bash
# Verify each branch can merge to main in sequence
echo "🔍 Verifying sequential mergeability..."

# Start from main
git checkout main

for EFFORT in $DEPENDENCY_ORDER; do
    echo "Testing merge of $EFFORT..."

    # Try merge (don't commit)
    git merge --no-commit --no-ff $EFFORT

    if [ $? -ne 0 ]; then
        echo "❌ $EFFORT cannot merge cleanly!"
        git merge --abort
        exit 1
    fi

    echo "✅ $EFFORT merges cleanly"

    # Don't abort - keep merged state for next test
    # This simulates sequential merging
done

# Clean up test merges
git reset --hard main
echo "✅ All efforts sequentially mergeable"
```

### Step 5: Create Rebase Report
```bash
cat > PR-REBASE-REPORT.md << 'EOF'
# PR Sequential Rebase Report
Date: $(date)
Main HEAD: $MAIN_HEAD

## Rebase Order (Dependency Sequence)
$(echo "$DEPENDENCY_ORDER" | nl)

## Rebase Results
$(for EFFORT in $DEPENDENCY_ORDER; do
    echo "- $EFFORT: ✅ Rebased and pushed"
done)

## Sequential Mergeability Test
✅ All efforts can merge to main in sequence without conflicts

## Next Steps
1. Efforts are ready for PR creation
2. PRs must be created in the above order
3. Each PR must be merged before the next is created

## Dependency Graph
\`\`\`
main
  ├── infrastructure-effort (no deps)
  ├── core-services-effort (depends on infrastructure)
  ├── feature-a-effort (depends on core)
  └── ui-effort (depends on all)
\`\`\`
EOF

echo "✅ Rebase report created"
```

## CONFLICT RESOLUTION

### Automatic Resolution
These conflicts are auto-resolved:
- Software Factory metadata files
- TODO files
- Planning documents
- Report files
- State tracking files

### Manual Resolution Required
These require manual intervention:
- Source code conflicts
- Test conflicts
- Configuration conflicts
- Documentation conflicts

### Resolution Protocol
```bash
# When manual resolution needed
1. Document conflicts in REBASE-CONFLICT-*.md
2. Escalate to orchestrator
3. Orchestrator spawns SW-Engineer for resolution
4. After resolution, re-run sequential rebase
```

## VERIFICATION REQUIREMENTS

### Pre-Rebase Checks
- [ ] Dependency order determined
- [ ] Main branch updated
- [ ] All branches pushed

### During Rebase
- [ ] Each branch rebases successfully
- [ ] Conflicts documented if any
- [ ] Auto-resolvable conflicts handled

### Post-Rebase Verification
- [ ] All branches rebased
- [ ] Sequential mergeability verified
- [ ] Report generated
- [ ] Branches pushed with lease

## ERROR RECOVERY

### Rebase Abort
```bash
# If rebase fails catastrophically
git rebase --abort
git checkout main
echo "Rebase aborted, original state preserved"
```

### Conflict Escalation
```bash
# If conflicts cannot be resolved
cat > REBASE-BLOCKED.md << EOF
# Rebase Blocked
Branch: $EFFORT
Reason: Unresolvable conflicts

Manual intervention required.
EOF
```

## PROJECT_DONE CRITERIA
✅ All efforts rebased onto main
✅ Correct dependency order maintained
✅ Sequential mergeability verified
✅ No unresolved conflicts
✅ Report generated

## GRADING PENALTIES
- Wrong rebase order: -50%
- Unresolved conflicts in PR: -40%
- No mergeability verification: -30%
- Breaking sequential merge: -100%

## INTEGRATE_WAVE_EFFORTS WITH OTHER RULES
- Enforces **R363** (Sequential Mergeability)
- Follows **R233** (Immediate Action)
- Prepares for **R370** (PR Plan Creation)
- Supports **R220** (Atomic PR Design)

---

*This rule ensures perfect sequential mergeability for production PRs.*