# PR_ARTIFACT_REMOVAL State Rules

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
## 🔴🔴🔴 STATE PURPOSE: Remove All SF Artifacts 🔴🔴🔴

### 🚨🚨🚨 RULE R367 - PR Branch Cleanup Protocol [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R367-pr-branch-cleanup-protocol.md`
**Criticality:** BLOCKING - Clean branches for production PRs

This state implements R367 requirements for comprehensive branch cleanup.

### MANDATORY ACTIONS (R233 + R367 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Load cleanup manifest**
   ```bash
   # Read assigned branches and artifacts
   cat PR-CLEANUP-MANIFEST.json
   BRANCHES=$(jq -r '.branches[]' PR-CLEANUP-MANIFEST.json)
   ```

2. **Create backup reference**
   ```bash
   # Before any deletions, save current SHA
   for BRANCH in $BRANCHES; do
     git checkout $BRANCH
     BACKUP_SHA=$(git rev-parse HEAD)
     echo "$BRANCH:$BACKUP_SHA" >> PR-CLEANUP-BACKUPS.txt
   done
   ```

3. **Remove SF artifacts systematically**
   ```bash
   for BRANCH in $BRANCHES; do
     echo "Cleaning branch: $BRANCH"
     git checkout $BRANCH

     # Remove directories
     git rm -rf todos/ 2>/dev/null || true
     git rm -rf efforts/ 2>/dev/null || true
     git rm -rf agent-states/ 2>/dev/null || true
     git rm -rf rule-library/ 2>/dev/null || true
     git rm -rf templates/ 2>/dev/null || true
     git rm -rf utilities/ 2>/dev/null || true
     git rm -rf .software-factory/planning/ 2>/dev/null || true
     git rm -rf protocols/ 2>/dev/null || true
     git rm -rf .claude/agents/ 2>/dev/null || true
     git rm -rf .claude/commands/ 2>/dev/null || true

     # Remove files
     git rm -f *-state.json 2>/dev/null || true
     git rm -f *.todo 2>/dev/null || true
     git rm -f .software-factory/phase*/wave*/*/CODE-REVIEW-REPORT--*.md 2>/dev/null || true
     git rm -f .software-factory/phase*/wave*/*/SPLIT-PLAN--*.md 2>/dev/null || true
     git rm -f .software-factory/PROJECT-IMPLEMENTATION-PLAN--*.md 2>/dev/null || true
     git rm -f .software-factory/phase*/wave*/*/IMPLEMENTATION-PLAN--*.md 2>/dev/null || true
     git rm -f software-factory-3.0-state-machine.json 2>/dev/null || true
     git rm -f RECOVERY-*.md 2>/dev/null || true
     git rm -f CURRENT-TODO-STATE.md 2>/dev/null || true

     # Verify core files still exist
     ls main.* Makefile README.* LICENSE 2>/dev/null || {
       echo "ERROR: Core files missing!"
       git reset --hard $BACKUP_SHA
       exit 1
     }

     # Commit cleanup
     git commit -m "cleanup: Remove Software Factory artifacts

   Removed all SF-specific files and directories while preserving
   core application functionality." || echo "No artifacts to remove"
   done
   ```

4. **Verify cleanup completeness**
   ```bash
   # Scan for any remaining artifacts
   for BRANCH in $BRANCHES; do
     git checkout $BRANCH

     # Should find nothing
     REMAINING=$(find . -name "*.todo" -o -name "*-state.json" \
       -o -path "*/todos/*" -o -path "*/efforts/*" | wc -l)

     if [ $REMAINING -gt 0 ]; then
       echo "WARNING: $REMAINING artifacts remain in $BRANCH"
     fi
   done
   ```

5. **Push cleaned branches**
   ```bash
   for BRANCH in $BRANCHES; do
     git push origin $BRANCH --force-with-lease
   done
   ```

### EXIT CRITERIA
✅ All artifacts removed
✅ Core files preserved
✅ Cleanup committed
✅ Branches pushed

### OUTPUT FILES
- `PR-CLEANUP-REPORT.md`
- `PR-CLEANUP-BACKUPS.txt`
- `PR-CLEANUP-VERIFICATION.json`

### CRITICAL SAFETY RULES
🚨 **NEVER DELETE**:
- main.* (main.go, main.py, main.rs, etc.)
- Makefile
- README.md or README.*
- LICENSE
- package.json, go.mod, Cargo.toml
- src/ directory (unless explicitly SF)
- pkg/ directory (unless explicitly SF)
- Any application source code

### PROHIBITED ACTIONS
❌ Do NOT use `rm -rf *` or wildcards on root
❌ Do NOT delete without checking
❌ Do NOT skip verification
❌ Do NOT proceed if core files affected

### ERROR RECOVERY
- Core file deleted → Reset to backup SHA
- Cleanup incomplete → Re-run removal
- Push rejected → Check branch protection
- Artifacts remain → Manual inspection

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

