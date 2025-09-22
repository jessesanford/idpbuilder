# Orchestrator - INIT_REPO_DECISION State Rules

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