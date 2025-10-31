# Orchestrator - INIT_START State Rules

## Purpose
Begin project initialization by capturing basic project information and setting up initialization infrastructure.

## Entry Criteria
- User has run `/init-software-factory` command
- Project name/prefix has been provided
- Initial idea/description has been captured

## Required Actions

### 1. Validate Inputs
```bash
# Verify project name is alphanumeric
if [[ ! "$PROJECT_PREFIX" =~ ^[a-zA-Z0-9]+$ ]]; then
    ERROR: Invalid project name
fi

# Verify idea is not empty
if [[ -z "$PROJECT_IDEA" ]]; then
    ERROR: No project idea provided
fi
```

### 2. Create Directory Structure
```bash
mkdir -p $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX
mkdir -p $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs
mkdir -p $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/plans
```

### 3. Initialize State File
Create or update init-state-${PROJECT_PREFIX}.json with:
- initialization_mode: true
- current_state: "INIT_START"
- project_prefix: provided value
- project_idea: provided description
- timestamp: current UTC time

### 4. Transition to INIT_LOAD_EXAMPLES
- Update state to "INIT_LOAD_EXAMPLES"
- Prepare to spawn architect agent
- Document what architect should do

## Exit Criteria
- Project directories created
- State file initialized
- Ready to spawn architect for template loading

## Transition
**MANDATORY**: → INIT_LOAD_EXAMPLES (spawn architect)

## Error Handling
- Invalid project name → Request valid name
- Missing directories → Create with proper permissions
- State file issues → Retry with error logging

## R313 Compliance
After transitioning to INIT_LOAD_EXAMPLES, if spawning architect:
- Record spawn in state file
- STOP immediately after spawn
- Provide continuation instructions

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

