# PR_COMMIT_SQUASH State Rules

## 🔴🔴🔴 STATE PURPOSE: Consolidate Commits in Branches 🔴🔴🔴

### 🚨🚨🚨 RULE R366 - PR Commit Consolidation Protocol [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R366-pr-commit-consolidation-protocol.md`
**Criticality:** BLOCKING - Clean commit history for production

This state implements R366 requirements for commit consolidation.

### MANDATORY ACTIONS (R233 + R366 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Process each assigned branch**
   ```bash
   # For each branch in assignment
   BRANCH="<assigned-branch>"
   git checkout $BRANCH

   # Capture all commit messages for preservation
   git log --format="- %s" main..HEAD > /tmp/commits-$BRANCH.txt
   COMMIT_COUNT=$(git rev-list --count main..HEAD)
   ```

2. **Perform soft reset to consolidate**
   ```bash
   # Reset to main while keeping all changes
   git reset --soft main

   # Stage all changes
   git add -A
   ```

3. **Create consolidated commit**
   ```bash
   # Generate comprehensive commit message
   cat > /tmp/commit-msg-$BRANCH.txt << EOF
   feat: [Concise feature description for $BRANCH]

   This commit consolidates $COMMIT_COUNT commits implementing [feature].

   Original commits included:
   $(cat /tmp/commits-$BRANCH.txt)

   Changes:
   - [Major change 1]
   - [Major change 2]
   - [Major change 3]

   Testing: All tests passing
   Dependencies: [List any dependencies]

   Co-authored-by: Software Factory <sf@company.io>
   EOF

   # Create the consolidated commit
   git commit -F /tmp/commit-msg-$BRANCH.txt
   ```

4. **Verify consolidation success**
   ```bash
   # Should show single commit
   git log --oneline main..HEAD

   # Verify all changes preserved
   git diff main --stat

   # Ensure no files lost
   git status --porcelain
   ```

5. **Force push consolidated branch**
   ```bash
   # Push with lease for safety
   git push origin $BRANCH --force-with-lease
   ```

### EXIT CRITERIA
✅ All assigned branches consolidated
✅ Commit history preserved in messages
✅ No changes lost
✅ Branches pushed to origin

### OUTPUT FILES
- `PR-CONSOLIDATION-REPORT.md`
- `PR-COMMIT-MESSAGES/` (directory with preserved messages)

### CRITICAL REQUIREMENTS
- MUST preserve ALL commit messages
- MUST maintain co-authorship
- MUST verify no code lost
- MUST use force-with-lease

### PROHIBITED ACTIONS
❌ Do NOT lose commit history
❌ Do NOT drop changes
❌ Do NOT merge branches
❌ Do NOT skip verification

### ERROR RECOVERY
- If changes lost → Reset to original SHA
- If push fails → Check branch protection
- If conflicts → Document and escalate

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

