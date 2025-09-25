# PR_CLEANUP_PLANNING State Rules

## 🔴🔴🔴 STATE PURPOSE: Plan SF Artifact Removal 🔴🔴🔴

### MANDATORY ACTIONS (R233 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Review discovered artifacts**
   ```bash
   # Load discovery results
   cat PR-DISCOVERY-REPORT-*.md
   cat PR-ARTIFACT-INVENTORY-*.json
   ```

2. **Create cleanup assignments**
   ```markdown
   ## SW Engineer Cleanup Tasks:

   ### Branch: phase1-wave1-branch
   - Remove: todos/, efforts/, agent-states/
   - Remove: orchestrator-state.json
   - Remove: *.todo files
   - Remove: SOFTWARE-FACTORY-STATE-MACHINE.md
   - Verify: No core files affected

   ### Branch: phase1-wave2-branch
   - Remove: rule-library/
   - Remove: templates/
   - Remove: CODE-REVIEW-REPORT.md
   - Remove: .claude/agents/
   ```

3. **Define cleanup validation criteria**
   ```bash
   # For each branch, must verify:
   - Zero SF artifacts remain
   - Core files untouched (main.*, Makefile, etc.)
   - No accidental deletions
   - Git history preserved
   ```

4. **Prioritize cleanup order**
   - Simple branches first (fewer artifacts)
   - Complex branches sequential
   - High-risk branches last
   - Parallel where possible

5. **Create rollback plan**
   ```bash
   # Before cleanup, save:
   - Branch SHA references
   - Artifact inventory
   - Backup strategy
   ```

### EXIT CRITERIA
✅ Cleanup tasks assigned
✅ Validation criteria defined
✅ Rollback plan created
✅ Agent instructions ready

### TRANSITIONS
- Success → PR_SPAWN_CLEANUP_AGENTS
- No artifacts found → PR_CONSOLIDATION_PLANNING
- Error → PR_ERROR_DETECTED

### PROHIBITED ACTIONS
❌ Do NOT perform cleanup yourself
❌ Do NOT delete any files
❌ Do NOT modify branches
❌ Do NOT skip planning

### SAFETY REQUIREMENTS
- Document every file to remove
- Never use wildcard deletions
- Preserve commit history
- Create cleanup manifest