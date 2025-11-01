# Orchestrator - INIT_REPO_DECISION State Rules

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
## Purpose
Analyze gathered requirements and determine the appropriate repository setup path.

## Entry Criteria
- Requirements gathering completed by architect
- Requirements document available in state file
- Answers include repository type information

## Required Actions

### 1. Read Requirements
Load init-state-${PROJECT_PREFIX}.json and examine:
- has_existing_codebase
- fork_url (if provided)
- wants_new_repository
- project_type

### 2. Apply Decision Logic
```
IF has_existing_codebase == true AND fork_url != null:
    next_state = "INIT_SETUP_UPSTREAM_FORK"
    repo_type = "upstream_fork"

ELIF wants_new_repository == true:
    next_state = "INIT_CREATE_NEW_REPO"
    repo_type = "new_repository"

ELSE:
    next_state = "INIT_LIBRARY_PROJECT"
    repo_type = "library_only"
```

### 3. Update State File
- current_state: chosen next state
- repo_type: determined type
- repo_decision_timestamp: current time
- repo_decision_rationale: brief explanation

### 4. Prepare for Repository Setup
Document what the SW Engineer agent needs to do:
- For upstream fork: clone URL, branch setup
- For new repo: directory, initial commit
- For library: structure only

## Exit Criteria
- Repository type determined
- Next state selected
- State file updated with decision
- Ready to spawn SW Engineer for setup

## Transitions
**MANDATORY ONE OF**:
- → INIT_SETUP_UPSTREAM_FORK (if has fork)
- → INIT_CREATE_NEW_REPO (if new project)
- → INIT_LIBRARY_PROJECT (if library/docs)

## Error Handling
- Ambiguous requirements → Return to INIT_REQUIREMENTS_GATHERING
- Invalid fork URL → Validate and re-ask
- Missing information → Identify gaps and gather

## R313 Compliance
After determining path and before spawning SW Engineer:
- Record decision in state file
- STOP after state transition
- Clear instructions for continuation

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

