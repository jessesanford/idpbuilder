# PR_REBASE_SEQUENCE State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
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

