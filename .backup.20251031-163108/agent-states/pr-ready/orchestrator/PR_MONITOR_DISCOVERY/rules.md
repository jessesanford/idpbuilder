# PR_MONITOR_DISCOVERY State Rules

## 🔴🔴🔴 STATE PURPOSE: Monitor Discovery Progress 🔴🔴🔴

### MANDATORY ACTIONS (R233 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Check agent completion status**
   ```bash
   # Check for discovery reports with proper metadata paths
   ls -la .software-factory/pr-ready/PR-DISCOVERY-REPORT--*.md 2>/dev/null || echo "No discovery reports found"
   ls -la .software-factory/pr-ready/PR-ARTIFACT-INVENTORY--*.json 2>/dev/null || echo "No artifact inventory found"
   ls -la .software-factory/pr-ready/PR-CONTAMINATION-ASSESSMENT--*.md 2>/dev/null || echo "No assessment files found"
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
         "branch2": ["orchestrator-state-v3.json", ...]
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

### MONITORING_SWE_PROGRESS REQUIREMENTS
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

