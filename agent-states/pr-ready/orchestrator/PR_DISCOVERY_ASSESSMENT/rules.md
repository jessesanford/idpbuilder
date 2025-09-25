# PR_DISCOVERY_ASSESSMENT State Rules

## 🔴🔴🔴 STATE PURPOSE: Plan Discovery of SF Artifacts 🔴🔴🔴

### MANDATORY ACTIONS (R233 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Create discovery task assignments**
   ```markdown
   ## Integration Agent Tasks:
   - Enumerate all effort branches
   - Document branch dependencies
   - Check main branch status
   - Verify upstream connectivity

   ## SW Engineer Tasks:
   - Scan for SF artifacts in each branch
   - Count artifact files per branch
   - Identify artifact patterns
   - Check for .claude directories

   ## Code Reviewer Tasks:
   - Assess contamination severity
   - Identify high-risk branches
   - Verify no core files affected
   - Create contamination report
   ```

2. **Define artifact patterns to detect**
   ```bash
   ARTIFACTS=(
     "todos/" "efforts/" "agent-states/"
     "rule-library/" "templates/" "utilities/"
     "*-state.json" "*.todo"
     "CODE-REVIEW-REPORT.md" "SPLIT-PLAN*.md"
     "PROJECT-IMPLEMENTATION-PLAN.md"
     "SOFTWARE-FACTORY-STATE-MACHINE.md"
     "RECOVERY-*.md" "CURRENT-TODO-STATE.md"
     "phase-plans/" "wave-plans/" "protocols/"
     ".claude/agents/" ".claude/commands/"
   )
   ```

3. **Create discovery instructions document**
   - Save as `PR-DISCOVERY-INSTRUCTIONS.md`
   - Include artifact patterns
   - Specify reporting format
   - Set discovery timeouts

### EXIT CRITERIA
✅ Task assignments created
✅ Artifact patterns defined
✅ Instructions documented
✅ Agent spawn plan ready

### TRANSITIONS
- Success → PR_SPAWN_DISCOVERY_AGENTS
- No branches found → PR_READY_ABORT
- Error → PR_ERROR_DETECTED

### PROHIBITED ACTIONS
❌ Do NOT perform discovery yourself
❌ Do NOT modify any files
❌ Do NOT skip planning phase
❌ Do NOT spawn agents in this state

### TIMING REQUIREMENTS
- Complete within 120 seconds
- Document all patterns
- Save state before transition