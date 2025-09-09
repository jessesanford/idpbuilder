# Orchestrator - WAITING_FOR_PROJECT_MERGE_PLAN State Rules

## 🛑🛑🛑 R322 MANDATORY CHECKPOINT BEFORE SPAWN_INTEGRATION_AGENT_PROJECT 🛑🛑🛑

**THIS IS A CRITICAL R322 CHECKPOINT STATE!**

### SUPREME LAW - PROJECT MERGE PLAN REQUIRES USER REVIEW

When transitioning from WAITING_FOR_PROJECT_MERGE_PLAN → SPAWN_INTEGRATION_AGENT_PROJECT:
- **MUST STOP** to allow user review of PROJECT-MERGE-PLAN.md
- **MUST UPDATE** state file to SPAWN_INTEGRATION_AGENT_PROJECT before stopping
- **MUST DISPLAY** checkpoint message listing ALL phases to be merged
- **MUST EXIT** cleanly to preserve context
- **VIOLATION = -100% IMMEDIATE FAILURE**

### CHECKPOINT PROTOCOL:
```markdown
## 🛑 R322 PROJECT INTEGRATION CHECKPOINT

### ✅ Project Merge Plan Created:
- Location: project-integration/PROJECT-MERGE-PLAN.md
- Phases to merge: [List all phase branches]
- Integration strategy: [Sequential/Parallel]

### 📊 Ready for Final Integration:
- Current State: WAITING_FOR_PROJECT_MERGE_PLAN ✅
- Next State: SPAWN_INTEGRATION_AGENT_PROJECT (pending approval)

### ⚠️ CRITICAL REVIEW REQUIRED
This is the FINAL integration merging ALL phases!
Please review the plan carefully before execution.

### ⏸️ STOPPED FOR USER REVIEW
To proceed after review: /continue-orchestrating
```

**STOP MEANS STOP - NO automatic continuation!**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

---

## State Context

**Purpose:**
Monitor Code Reviewer creating the project-level merge plan that will integrate all phase branches.

## Primary Actions

1. **Check for Merge Plan**:
   - Look for PROJECT-MERGE-PLAN.md in project integration workspace
   - Verify plan includes all phases in correct order (R270)
2. **Validate Plan Completeness**:
   - All phase branches listed
   - Merge order specified
   - Conflict resolution strategy documented
3. **Update State** when plan is ready

## Valid State Transitions

- **SUCCESS** → SPAWN_INTEGRATION_AGENT_PROJECT (plan ready)
- **TIMEOUT** → ERROR_RECOVERY (plan not created)
- **FAILURE** → ERROR_RECOVERY (invalid plan)