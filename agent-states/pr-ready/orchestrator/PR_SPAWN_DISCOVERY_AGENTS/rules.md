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