# PR_SPAWN_CLEANUP_AGENTS State Rules

## 🔴🔴🔴 STATE PURPOSE: Spawn Cleanup Agents 🔴🔴🔴

### MANDATORY ACTIONS (R233 + R313 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Spawn SW Engineer(s) for artifact removal**
   ```bash
   # For parallel cleanup (if branches independent)
   /spawn agent-sw-engineer PR_ARTIFACT_REMOVAL \
     --branches "branch1,branch2,branch3" \
     --cleanup-manifest "PR-CLEANUP-MANIFEST.json" \
     --parallel true

   # OR for sequential cleanup (if dependencies exist)
   /spawn agent-sw-engineer PR_ARTIFACT_REMOVAL \
     --branch "branch1" \
     --cleanup-list "PR-CLEANUP-TASKS-branch1.md"
   ```

2. **Spawn Integration Agent if main needs reset**
   ```bash
   # Only if main is contaminated
   /spawn agent-integration PR_MAIN_RESET \
     --upstream "upstream/main" \
     --preserve-commits "[legitimate-sha-list]"
   ```

3. **Record cleanup assignments**
   ```json
   {
     "cleanup_agents": {
       "sw_engineer_1": {
         "id": "swe-cleanup-<timestamp>",
         "branches": ["branch1", "branch2"],
         "task": "artifact_removal"
       },
       "integration": {
         "id": "int-main-<timestamp>",
         "task": "main_reset",
         "status": "optional"
       }
     }
   }
   ```

### CRITICAL R313 REQUIREMENT
🚨 **MUST STOP IMMEDIATELY AFTER SPAWNING**
- Update state file with agent IDs
- Save cleanup assignments
- Commit and push state
- EXIT to preserve context

### EXIT CRITERIA
✅ Cleanup agents spawned
✅ Assignments recorded
✅ State saved
✅ IMMEDIATE STOP per R313

### TRANSITIONS
- After stop → PR_MONITOR_CLEANUP
- Spawn failure → PR_ERROR_DETECTED

### PROHIBITED ACTIONS
❌ Do NOT continue after spawning
❌ Do NOT perform cleanup yourself
❌ Do NOT monitor in same session
❌ Do NOT violate R313

### PARALLELIZATION RULES (R151)
- All parallel spawns within 5 seconds
- Sequential spawns if branch dependencies
- Document parallelization decision
- Monitor load on system

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

