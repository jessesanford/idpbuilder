# PR_MONITOR_DISCOVERY State Rules

## 🔴🔴🔴 STATE PURPOSE: Monitor Discovery Progress 🔴🔴🔴

### MANDATORY ACTIONS (R233 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Check agent completion status**
   ```bash
   # Check for discovery reports
   ls -la PR-DISCOVERY-REPORT-*.md
   ls -la PR-ARTIFACT-INVENTORY-*.json
   ls -la PR-CONTAMINATION-ASSESSMENT-*.md
   ```

2. **Aggregate discovery findings**
   ```json
   {
     "discovery_results": {
       "branches_found": ["branch1", "branch2", ...],
       "total_artifacts": {
         "todos": 45,
         "state_files": 12,
         "planning_docs": 23,
         "rule_library": 156
       },
       "contamination_by_branch": {
         "branch1": ["todos/", "efforts/", ...],
         "branch2": ["orchestrator-state.json", ...]
       },
       "high_risk_branches": [],
       "dependency_order": ["branch1", "branch2", ...]
     }
   }
   ```

3. **Validate discovery completeness**
   - All branches scanned
   - All artifact types checked
   - Dependency tree complete
   - Core file status verified

4. **Create cleanup priority list**
   - Order branches by contamination
   - Identify quick wins
   - Flag complex cases

### EXIT CRITERIA
✅ All agents completed
✅ Discovery aggregated
✅ Artifacts inventoried
✅ Cleanup plan ready

### TRANSITIONS
- All complete → PR_CLEANUP_PLANNING
- Agents still running → Continue monitoring
- Timeout (>10 min) → PR_ERROR_DETECTED
- Discovery failed → PR_ERROR_DETECTED

### MONITORING REQUIREMENTS
- Check every 30 seconds
- Log progress updates
- Track agent heartbeats
- Detect stalled agents

### PROHIBITED ACTIONS
❌ Do NOT perform discovery yourself
❌ Do NOT modify discovered data
❌ Do NOT proceed with incomplete discovery
❌ Do NOT skip validation

### TIMEOUT HANDLING
- 10 minute hard timeout
- Warning at 7 minutes
- Kill stalled agents at timeout
- Preserve partial results