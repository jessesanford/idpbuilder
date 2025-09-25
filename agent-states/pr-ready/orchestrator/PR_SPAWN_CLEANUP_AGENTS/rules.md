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