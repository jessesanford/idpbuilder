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
     "software-factory-3.0-state-machine.json"
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

