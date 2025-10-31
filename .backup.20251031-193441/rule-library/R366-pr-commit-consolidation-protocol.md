# 🚨🚨🚨 BLOCKING RULE R366 - PR Commit Consolidation Protocol

**Criticality:** BLOCKING - Clean commit history for production
**Enforcement:** MANDATORY - All PR-Ready operations
**Created:** 2025-01-21

## PURPOSE
Consolidate multiple development commits into clean, atomic commits for PR submission while preserving all commit history and co-authorship information.

## CONSOLIDATION REQUIREMENTS

### When to Consolidate
- Branch has >10 commits for single feature
- Multiple "fix", "wip", or "debug" commits present
- Commits don't follow conventional commit standards
- Request for clean history from maintainers

### When NOT to Consolidate
- Commits already atomic and well-structured
- Multiple distinct features (need separate PRs)
- Commits from multiple authors need preservation
- Historical debugging value in commit sequence

## CONSOLIDATION PROTOCOL

### Step 1: Backup Original State
```bash
# CRITICAL: Always backup before consolidation
BRANCH=$(git branch --show-current)
BACKUP_BRANCH="${BRANCH}-backup-$(date +%Y%m%d-%H%M%S)"

# Create backup branch
git branch $BACKUP_BRANCH
echo "✅ Backup created: $BACKUP_BRANCH"

# Save commit history
git log --format="[%h] %an: %s" main..HEAD > commits-original-$BRANCH.txt
```

### Step 2: Analyze Commit Structure
```bash
# Count commits to consolidate
COMMIT_COUNT=$(git rev-list --count main..HEAD)
echo "Consolidating $COMMIT_COUNT commits"

# Extract unique authors for co-authorship
git log main..HEAD --format="%an <%ae>" | sort -u > authors-$BRANCH.txt

# Group commits by type
git log main..HEAD --format="%s" | grep -E "^(feat|fix|docs|test|refactor|chore):" | sort > commits-by-type.txt
```

### Step 3: Perform Soft Reset
```bash
# Reset to base while preserving changes
git reset --soft main

# Verify all changes preserved
git status --short
git diff --staged --stat

# Stage all changes
git add -A
```

### Step 4: Create Consolidated Commit
```bash
# Generate comprehensive commit message
cat > commit-message.txt << 'EOF'
feat: [Concise description of feature]

This commit consolidates $COMMIT_COUNT commits implementing [feature/fix].

## Changes
- [Primary change 1]
- [Primary change 2]
- [Primary change 3]

## Original Commits
$(cat commits-original-$BRANCH.txt | sed 's/^/- /')

## Testing
- Unit tests: ✅ Passing
- Integration tests: ✅ Verified
- Build: ✅ Successful

## Breaking Changes
None [or describe if any]

$(cat authors-$BRANCH.txt | sed 's/^/Co-authored-by: /')
EOF

# Create the consolidated commit
git commit -F commit-message.txt
```

### Step 5: Verify Consolidation
```bash
# Verify single commit created
COMMITS_AFTER=$(git rev-list --count main..HEAD)
if [ "$COMMITS_AFTER" -ne 1 ]; then
    echo "❌ ERROR: Expected 1 commit, got $COMMITS_AFTER"
    git reset --hard $BACKUP_BRANCH
    exit 1
fi

# Verify no changes lost
DIFF_BEFORE=$(git diff main $BACKUP_BRANCH --stat | tail -1)
DIFF_AFTER=$(git diff main HEAD --stat | tail -1)
if [ "$DIFF_BEFORE" != "$DIFF_AFTER" ]; then
    echo "❌ ERROR: Changes lost during consolidation!"
    git reset --hard $BACKUP_BRANCH
    exit 1
fi

echo "✅ Consolidation successful"
```

### Step 6: Force Push with Lease
```bash
# Push consolidated commit (safely)
git push origin $BRANCH --force-with-lease

# If lease fails, pull and retry
if [ $? -ne 0 ]; then
    echo "⚠️ Force-with-lease failed, checking remote..."
    git fetch origin $BRANCH

    # Compare with remote
    REMOTE_DIFFERS=$(git diff origin/$BRANCH --stat)
    if [ -n "$REMOTE_DIFFERS" ]; then
        echo "❌ Remote has unexpected changes, aborting"
        exit 1
    fi

    # Retry with regular force if safe
    git push origin $BRANCH --force
fi
```

## SAFETY REQUIREMENTS

### MANDATORY Safety Checks
1. **ALWAYS create backup branch** before consolidation
2. **ALWAYS verify no code lost** after reset
3. **ALWAYS preserve co-authorship** in commit message
4. **ALWAYS use --force-with-lease** for initial push
5. **ALWAYS save original commit history** to file

### PROHIBITED Actions
❌ **NEVER** consolidate without backup
❌ **NEVER** lose commit attribution
❌ **NEVER** drop changes during consolidation
❌ **NEVER** force push without verification
❌ **NEVER** consolidate commits from different features

## ERROR RECOVERY

### If Changes Lost
```bash
# Immediate recovery from backup
git reset --hard $BACKUP_BRANCH
git branch -D $BRANCH
git checkout -b $BRANCH
```

### If Push Conflicts
```bash
# Check remote state
git fetch origin $BRANCH
git log origin/$BRANCH..HEAD --oneline

# If remote ahead, rebase
git rebase origin/$BRANCH
```

### If Consolidation Fails
```bash
# Return to original state
git reset --hard $BACKUP_BRANCH
echo "Consolidation aborted, original state restored"
```

## OUTPUT ARTIFACTS

### Required Files
- `commits-original-${BRANCH}.txt` - Original commit history
- `authors-${BRANCH}.txt` - Co-author list
- `commit-message.txt` - Consolidated commit message
- `PR-CONSOLIDATION-REPORT.md` - Summary report

### Report Format
```markdown
# PR Commit Consolidation Report
Branch: ${BRANCH}
Date: $(date)

## Summary
- Original commits: ${COMMIT_COUNT}
- Consolidated to: 1
- Authors preserved: $(wc -l < authors-${BRANCH}.txt)
- Changes preserved: ✅

## Backup Information
- Backup branch: ${BACKUP_BRANCH}
- Can be deleted after PR merge

## Verification
- [ ] All changes preserved
- [ ] Co-authorship maintained
- [ ] Commit message comprehensive
- [ ] Successfully pushed
```

## PROJECT_DONE CRITERIA
✅ Single atomic commit created
✅ All changes preserved
✅ Co-authorship maintained
✅ Original history backed up
✅ Successfully pushed to origin

## GRADING PENALTIES
- No backup before consolidation: -40%
- Lost changes during consolidation: -100%
- Missing co-authorship: -30%
- No verification: -25%
- Force push without lease: -20%

## INTEGRATE_WAVE_EFFORTS WITH OTHER RULES
- Follows **R233** (Immediate Action)
- Prepares for **R370** (PR Plan Creation)
- Enables **R368** (Sequential Rebase)

---

*This rule ensures clean, professional commit history for production PRs.*