# PR_ARTIFACT_REMOVAL State Rules

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
     git rm -rf phase-plans/ 2>/dev/null || true
     git rm -rf wave-plans/ 2>/dev/null || true
     git rm -rf protocols/ 2>/dev/null || true
     git rm -rf .claude/agents/ 2>/dev/null || true
     git rm -rf .claude/commands/ 2>/dev/null || true

     # Remove files
     git rm -f *-state.json 2>/dev/null || true
     git rm -f *.todo 2>/dev/null || true
     git rm -f CODE-REVIEW-REPORT.md 2>/dev/null || true
     git rm -f SPLIT-PLAN*.md 2>/dev/null || true
     git rm -f PROJECT-IMPLEMENTATION-PLAN.md 2>/dev/null || true
     git rm -f EFFORT-IMPLEMENTATION-PLAN.md 2>/dev/null || true
     git rm -f SOFTWARE-FACTORY-STATE-MACHINE.md 2>/dev/null || true
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