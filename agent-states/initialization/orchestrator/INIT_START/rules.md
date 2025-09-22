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