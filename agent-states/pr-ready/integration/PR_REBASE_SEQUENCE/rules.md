# PR_REBASE_SEQUENCE State Rules

## 🔴🔴🔴 STATE PURPOSE: Execute Sequential Rebasing 🔴🔴🔴

### 🔴🔴🔴 RULE R368 - PR Sequential Rebase Protocol [SUPREME LAW]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R368-pr-sequential-rebase-protocol.md`
**Criticality:** SUPREME LAW - Incorrect rebase order = merge conflicts

This state implements R368 requirements for sequential rebasing to ensure R363 compliance.

### MANDATORY ACTIONS (R233 + R368 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Load rebase sequence**
   ```bash
   # Read assigned rebase order
   cat PR-REBASE-SEQUENCE.json
   BRANCHES=$(jq -r '.rebase_order[]' PR-REBASE-SEQUENCE.json)
   ```

2. **Execute sequential rebasing**
   ```bash
   # Start with first branch on main
   FIRST_BRANCH=$(echo $BRANCHES | cut -d' ' -f1)
   git checkout $FIRST_BRANCH
   git rebase main

   if [ $? -ne 0 ]; then
     # Handle conflicts
     echo "Conflicts in $FIRST_BRANCH"
     # Document conflict files
     git status --porcelain | grep "^UU"

     # Apply standard resolution
     git checkout --theirs .  # For add/add conflicts
     git add -A
     git rebase --continue
   fi

   git push origin $FIRST_BRANCH --force-with-lease

   # Each subsequent branch on previous
   PREV_BRANCH=$FIRST_BRANCH
   for BRANCH in $(echo $BRANCHES | cut -d' ' -f2-); do
     git checkout $BRANCH
     git rebase $PREV_BRANCH

     if [ $? -ne 0 ]; then
       # Document and resolve conflicts
       echo "Conflicts rebasing $BRANCH on $PREV_BRANCH"
       git status --porcelain | grep "^UU"

       # Standard resolution for known patterns
       git checkout --theirs pkg/certs/*  # Known conflict area
       git add -A
       git rebase --continue
     fi

     git push origin $BRANCH --force-with-lease
     PREV_BRANCH=$BRANCH
   done
   ```

3. **Create incremental branches if needed**
   ```bash
   # For truly incremental PRs
   for BRANCH in $BRANCHES; do
     if [ "$CREATE_INCREMENTAL" = "true" ]; then
       # Get only new changes
       git diff $PREV_BRANCH..$BRANCH > /tmp/$BRANCH.patch

       git checkout $PREV_BRANCH
       git checkout -b ${BRANCH}-incremental
       git apply /tmp/$BRANCH.patch
       git add -A
       git commit -m "feat: Incremental changes for $BRANCH"
       git push origin ${BRANCH}-incremental
     fi
   done
   ```

4. **Document conflict resolutions**
   ```markdown
   # Conflict Resolution Log

   ## Branch: phase1-wave1-effort1
   - Rebased on: main
   - Conflicts: None

   ## Branch: phase1-wave1-effort2
   - Rebased on: phase1-wave1-effort1
   - Conflicts in: pkg/certs/extractor.go
   - Resolution: Accepted incoming (--theirs)
   ```

### EXIT CRITERIA
✅ All branches rebased in sequence
✅ Conflicts resolved and documented
✅ Branches pushed to origin
✅ Dependency chain established

### OUTPUT FILES
- `PR-REBASE-REPORT.md`
- `PR-CONFLICT-RESOLUTIONS.md`
- `PR-INCREMENTAL-BRANCHES.txt` (if applicable)

### CONFLICT RESOLUTION STRATEGY
**Standard patterns:**
- `pkg/certs/*` → Use --theirs (accept incoming)
- Add/add conflicts → Use --theirs
- Modify/delete → Investigate case-by-case
- Whitespace conflicts → Auto-resolve

### PROHIBITED ACTIONS
❌ Do NOT skip conflict resolution
❌ Do NOT break dependency chain
❌ Do NOT merge branches (only rebase)
❌ Do NOT lose commits

### ERROR RECOVERY
- Rebase fails → Reset to original SHA
- Unresolvable conflict → Document and escalate
- Push rejected → Check branch protection
- Lost commits → Recover from reflog