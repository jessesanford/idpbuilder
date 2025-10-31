# PR_SPAWN_DISCOVERY_AGENTS State Rules

## 🔴🔴🔴 STATE PURPOSE: Spawn Agents for Discovery 🔴🔴🔴

### MANDATORY ACTIONS (R233 + R313 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Spawn Integration Agent**
   ```bash
   /spawn agent-integration PR_BRANCH_INVENTORY \
     --task "Inventory all effort branches and dependencies" \
     --instructions "PR-DISCOVERY-INSTRUCTIONS.md"
   ```

2. **Spawn SW Engineer Agent**
   ```bash
   /spawn agent-sw-engineer PR_ARTIFACT_SCAN \
     --task "Scan branches for SF artifacts" \
     --branches "[list from state file]" \
     --patterns "[artifact patterns]"
   ```

3. **Spawn Code Reviewer Agent**
   ```bash
   /spawn agent-code-reviewer PR_CONTAMINATION_ASSESSMENT \
     --task "Assess SF contamination levels" \
     --critical-files "[core file list]"
   ```

4. **Record spawned agents**
   ```json
   {
     "spawned_agents": {
       "integration": {
         "id": "int-<timestamp>",
         "state": "PR_BRANCH_INVENTORY",
         "task": "branch inventory"
       },
       "sw_engineer": {
         "id": "swe-<timestamp>",
         "state": "PR_ARTIFACT_SCAN",
         "task": "artifact scanning"
       },
       "code_reviewer": {
         "id": "rev-<timestamp>",
         "state": "PR_CONTAMINATION_ASSESSMENT",
         "task": "contamination assessment"
       }
     }
   }
   ```

### CRITICAL R313 REQUIREMENT
🚨 **MUST STOP IMMEDIATELY AFTER SPAWNING**
- Save state file
- Commit changes
- EXIT to preserve context
- Wait for continuation command

### EXIT CRITERIA
✅ All agents spawned
✅ Agent IDs recorded
✅ State file updated
✅ IMMEDIATE STOP per R313

### TRANSITIONS
- After stop → PR_MONITOR_DISCOVERY
- Spawn failure → PR_ERROR_DETECTED

### PROHIBITED ACTIONS
❌ Do NOT continue after spawning
❌ Do NOT monitor in same session
❌ Do NOT perform any discovery
❌ Do NOT violate R313 stop requirement

### TIMING REQUIREMENTS
- All spawns within 5 seconds (R151)
- Save state immediately
- Stop within 10 seconds of spawning

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

