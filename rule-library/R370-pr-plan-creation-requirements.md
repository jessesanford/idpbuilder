# 🚨🚨🚨 BLOCKING RULE R370 - PR Plan Creation Requirements

**Criticality:** BLOCKING - Humans need clear PR instructions
**Enforcement:** MANDATORY - Final PR-Ready deliverable
**Created:** 2025-01-21

## PURPOSE
Create comprehensive MASTER-PR-PLAN.md that provides humans with exact, step-by-step instructions for creating and merging Pull Requests in the correct order with zero ambiguity.

## PLAN REQUIREMENTS

### Mandatory Sections
1. **Executive Summary** - What's being delivered
2. **Prerequisites** - What humans need before starting
3. **Verification Checklist** - Pre-flight checks
4. **PR Merge Order** - EXACT sequence with dependencies
5. **Individual PR Templates** - Complete, ready-to-paste
6. **Commands** - Copy-paste ready CLI/GUI instructions
7. **Rollback Plan** - If something goes wrong
8. **Progress Tracking** - Checklist for humans

### Information Per PR
- Branch name (exact)
- PR title (conventional format)
- PR body (complete template)
- Dependencies (which PRs must merge first)
- Verification steps
- Merge method (squash/rebase/merge)

## PLAN CREATION PROTOCOL

### Step 1: Gather Effort Information
```bash
echo "📊 Gathering effort information..."

# Collect all effort branches
EFFORTS_DATA=""
for effort_dir in efforts/*/; do
    if [ -d "$effort_dir/.git" ]; then
        cd "$effort_dir"
        BRANCH=$(git branch --show-current)
        COMMITS=$(git rev-list --count main..HEAD 2>/dev/null || echo 0)
        FILES=$(git diff --name-only main..HEAD 2>/dev/null | wc -l)
        ADDITIONS=$(git diff --stat main..HEAD 2>/dev/null | tail -1 | grep -o '[0-9]* insertion' | cut -d' ' -f1)
        DELETIONS=$(git diff --stat main..HEAD 2>/dev/null | tail -1 | grep -o '[0-9]* deletion' | cut -d' ' -f1)

        EFFORTS_DATA="$EFFORTS_DATA|$BRANCH:$COMMITS:$FILES:$ADDITIONS:$DELETIONS"
    fi
done

echo "Collected data for $(echo "$EFFORTS_DATA" | tr '|' '\n' | wc -l) efforts"
```

### Step 2: Determine Dependency Order
```bash
# Load dependency information
if [ -f ".software-factory/effort-dependencies.json" ]; then
    MERGE_ORDER=$(jq -r '.merge_order[]' .software-factory/effort-dependencies.json)
else
    # Fallback: alphabetical order with warning
    MERGE_ORDER=$(echo "$EFFORTS_DATA" | tr '|' '\n' | cut -d: -f1 | sort)
    echo "⚠️ WARNING: No dependency file, using alphabetical order"
fi
```

### Step 3: Generate MASTER-PR-PLAN.md
```bash
cat > MASTER-PR-PLAN.md << 'EOF'
# 🚀 MASTER PR PLAN - SOFTWARE FACTORY DELIVERY

**Generated**: $(date)
**Total PRs**: $(echo "$MERGE_ORDER" | wc -w)
**Project**: $(basename $(pwd))

## 📋 EXECUTIVE SUMMARY

Software Factory has completed all implementation, testing, and validation for this project. This document contains step-by-step instructions for creating Pull Requests to merge all completed work into the main branch.

**CRITICAL**: PRs must be created and merged in the EXACT order specified to avoid conflicts.

## ⚠️ PREREQUISITES

Before starting, ensure you have:
- [ ] Write access to the repository
- [ ] GitHub CLI (`gh`) installed OR web browser access
- [ ] All branches listed below exist in origin
- [ ] No uncommitted local changes
- [ ] Understanding that merge order is MANDATORY

## ✅ PRE-FLIGHT VERIFICATION

Run these checks before creating any PRs:

\`\`\`bash
# Verify all branches exist
for branch in $MERGE_ORDER; do
    git ls-remote --heads origin "$branch" > /dev/null 2>&1 && echo "✅ $branch" || echo "❌ $branch MISSING"
done

# Check integration test results
git log integration-testing --oneline -1

# Verify no uncommitted changes
cd efforts/ && find . -name ".git" -type d -exec sh -c 'cd "$(dirname {})" && pwd && git status --porcelain' \;
\`\`\`

## 🔄 MANDATORY PR MERGE ORDER

**⚠️ WARNING**: Deviating from this order WILL cause merge conflicts!

EOF

# Add each PR in order
PR_NUMBER=1
for BRANCH in $MERGE_ORDER; do
    cat >> MASTER-PR-PLAN.md << EOF

### PR #${PR_NUMBER}: $BRANCH
**Branch**: \`$BRANCH\`
**Depends On**: $([ $PR_NUMBER -eq 1 ] && echo "None (merge first)" || echo "PR #$((PR_NUMBER-1))")
**Merge Method**: Squash and merge

#### Quick Stats
- Commits: $(echo "$EFFORTS_DATA" | grep "$BRANCH" | cut -d: -f2)
- Files Changed: $(echo "$EFFORTS_DATA" | grep "$BRANCH" | cut -d: -f3)
- Additions: +$(echo "$EFFORTS_DATA" | grep "$BRANCH" | cut -d: -f4)
- Deletions: -$(echo "$EFFORTS_DATA" | grep "$BRANCH" | cut -d: -f5)

#### PR Title
\`\`\`
feat: $(echo $BRANCH | sed 's/effort-//' | sed 's/-/ /g')
\`\`\`

#### Merge Instructions
1. $([ $PR_NUMBER -gt 1 ] && echo "⚠️ WAIT for PR #$((PR_NUMBER-1)) to be fully merged")
2. Create PR from \`$BRANCH\` to \`main\`
3. Copy title and body from templates
4. Wait for CI/CD checks
5. Request reviews if required
6. Merge using "Squash and merge"
7. ✅ Delete branch after merge
8. Proceed to PR #$((PR_NUMBER+1))

---

EOF
    PR_NUMBER=$((PR_NUMBER + 1))
done

# Add PR body templates
cat >> MASTER-PR-PLAN.md << 'EOF'

## 📝 PR BODY TEMPLATES

Copy and paste these exactly into each PR:

EOF

PR_NUMBER=1
for BRANCH in $MERGE_ORDER; do
    cat >> MASTER-PR-PLAN.md << EOF

### PR Body for PR #${PR_NUMBER} ($BRANCH)

\`\`\`markdown
## Summary
Implementation of $(echo $BRANCH | sed 's/effort-//' | sed 's/-/ /g') functionality.

## Changes
- [Key change 1 from effort]
- [Key change 2 from effort]
- [Key change 3 from effort]

## Testing
- ✅ Unit tests: All passing
- ✅ Integration tests: Verified in integration-testing branch
- ✅ Build: Successful
- ✅ Manual testing: Completed

## Dependencies
$([ $PR_NUMBER -eq 1 ] && echo "None - this is the base implementation" || echo "Requires PR #$((PR_NUMBER-1)) to be merged first")

## Breaking Changes
None

## Documentation
- README: Updated
- API Docs: Current
- Comments: Added where needed

## Validation
This PR has been validated by Software Factory's PR-Ready validation suite.
- No SF artifacts present
- No stubs or TODOs
- All tests passing
- Production ready

## Review Notes
This is part of a sequence of ${#MERGE_ORDER[@]} PRs that must be merged in order.
Integration testing has verified all PRs work together correctly.

---
*Generated by Software Factory 2.0*
\`\`\`

EOF
    PR_NUMBER=$((PR_NUMBER + 1))
done

# Add CLI commands section
cat >> MASTER-PR-PLAN.md << 'EOF'

## 💻 CLI COMMANDS

### Using GitHub CLI (Recommended)

\`\`\`bash
# Set these variables
REPO_OWNER="your-org"
REPO_NAME="your-repo"

EOF

PR_NUMBER=1
for BRANCH in $MERGE_ORDER; do
    cat >> MASTER-PR-PLAN.md << EOF
# PR #${PR_NUMBER}: $BRANCH
gh pr create \\
  --repo \$REPO_OWNER/\$REPO_NAME \\
  --base main \\
  --head $BRANCH \\
  --title "feat: $(echo $BRANCH | sed 's/effort-//' | sed 's/-/ /g')" \\
  --body-file pr-body-${PR_NUMBER}.md

EOF
    PR_NUMBER=$((PR_NUMBER + 1))
done

cat >> MASTER-PR-PLAN.md << 'EOF'
\`\`\`

### Using GitHub Web Interface

1. Go to: https://github.com/[owner]/[repo]/pulls
2. Click "New pull request"
3. Set base: `main`
4. Set compare: `[effort-branch]`
5. Copy title and body from templates above
6. Click "Create pull request"
7. Repeat for each PR IN ORDER

## 🔄 ROLLBACK PROCEDURES

### If PR Fails CI/CD
\`\`\`bash
# Don't merge, investigate failure
gh pr checks [PR-NUMBER]
# Fix in branch and push
# CI will re-run automatically
\`\`\`

### If PR Causes Issues After Merge
\`\`\`bash
# Revert the problematic PR
gh pr revert [PR-NUMBER]
# This creates a new PR that undoes changes
# Merge the revert PR immediately
# Fix issues in new branch
\`\`\`

### If Multiple PRs Need Rollback
\`\`\`bash
# Revert in REVERSE order (last merged first)
for pr in [PR-N] [PR-N-1] [PR-N-2]; do
    gh pr revert $pr
done
\`\`\`

## 📊 PROGRESS TRACKING

Use this checklist to track your progress:

### PR Creation
EOF

for BRANCH in $MERGE_ORDER; do
    echo "- [ ] PR for $BRANCH created" >> MASTER-PR-PLAN.md
done

cat >> MASTER-PR-PLAN.md << 'EOF'

### PR Merging
EOF

for BRANCH in $MERGE_ORDER; do
    echo "- [ ] PR for $BRANCH merged" >> MASTER-PR-PLAN.md
done

cat >> MASTER-PR-PLAN.md << 'EOF'

### Final Verification
- [ ] All PRs merged successfully
- [ ] Main branch builds successfully
- [ ] All tests pass on main
- [ ] Deployment ready (if applicable)

## 🎯 COMPLETION CRITERIA

The project is complete when:
1. ✅ All PRs merged in order
2. ✅ Main branch builds successfully
3. ✅ All tests pass
4. ✅ No merge conflicts occurred
5. ✅ Branches deleted

## 📞 SUPPORT

If you encounter issues:
1. Check the error messages carefully
2. Verify you followed the exact order
3. Consult the rollback procedures
4. Contact the development team if blocked

## 🏁 FINAL NOTES

- Total implementation effort: [X] days
- Total lines of code: [Y]
- Test coverage: [Z]%
- All work validated by Software Factory 2.0

**Thank you for using Software Factory!**

---

*This document is the final deliverable from Software Factory 2.0*
*All implementation is complete and ready for production*
EOF

echo "✅ MASTER-PR-PLAN.md created successfully"
```

### Step 4: Create Individual PR Body Files
```bash
# Create separate files for easy copying
PR_NUMBER=1
for BRANCH in $MERGE_ORDER; do
    cat > pr-body-${PR_NUMBER}.md << EOF
## Summary
[Specific summary for $BRANCH]

## Changes
[Detailed changes for this effort]

## Testing
- Unit tests: ✅ Passing
- Integration: ✅ Verified
- Build: ✅ Successful

## Dependencies
$([ $PR_NUMBER -eq 1 ] && echo "None" || echo "PR #$((PR_NUMBER-1))")

[Rest of template...]
EOF
    PR_NUMBER=$((PR_NUMBER + 1))
done
```

### Step 5: Validate Plan Completeness
```bash
# Ensure plan has all required sections
echo "🔍 Validating PR plan completeness..."

REQUIRED_SECTIONS=(
    "EXECUTIVE SUMMARY"
    "PREREQUISITES"
    "MANDATORY PR MERGE ORDER"
    "PR BODY TEMPLATES"
    "CLI COMMANDS"
    "ROLLBACK PROCEDURES"
    "PROGRESS TRACKING"
)

for section in "${REQUIRED_SECTIONS[@]}"; do
    if grep -q "$section" MASTER-PR-PLAN.md; then
        echo "✅ Section present: $section"
    else
        echo "❌ Missing section: $section"
        exit 1
    fi
done

echo "✅ PR plan validation passed"
```

## PROJECT_DONE CRITERIA
✅ MASTER-PR-PLAN.md created
✅ All efforts included
✅ Dependency order specified
✅ PR templates complete
✅ Commands ready to copy
✅ Rollback procedures documented

## OUTPUT ARTIFACTS
- `MASTER-PR-PLAN.md` - Complete PR instructions
- `pr-body-*.md` - Individual PR body files
- `pr-commands.sh` - Executable PR creation script

## GRADING PENALTIES
- Missing dependency order: -40%
- Incomplete PR templates: -30%
- No rollback plan: -25%
- Missing verification: -20%

## INTEGRATE_WAVE_EFFORTS WITH OTHER RULES
- Depends on **R368** (Sequential Rebase)
- Requires **R369** (Validation Complete)
- Implements **R279** (MASTER-PR-PLAN)
- Supports **R220** (Atomic PR Design)

---

*This rule ensures humans have perfect instructions for PR creation.*