# PR_CLEANUP_PLANNING State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
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
   - Remove: orchestrator-state-v3.json
   - Remove: *.todo files
   - Remove: software-factory-3.0-state-machine.json
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

