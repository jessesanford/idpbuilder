---
name: R181 - Orchestrator MUST Set Up Workspaces
criticality: BLOCKING
agent: orchestrator
state: SETUP_EFFORT_INFRASTRUCTURE
---

# 🚨🚨🚨 RULE R181 - Orchestrator MUST Set Up Workspaces 🚨🚨🚨

## Rule Statement
The orchestrator agent MUST set up proper working copies/workspaces for all agents before spawning them. Each agent requires an isolated directory structure to prevent conflicts and ensure clean implementation.

## Rationale
Without proper workspace isolation, agents will:
- Conflict with each other's work
- Corrupt the main working copy
- Violate the one-agent-per-effort principle
- Create merge conflicts and integration issues

## Implementation Requirements

### 1. Workspace Creation Protocol
```bash
# For each effort, create isolated working copy
effort_dir="idpbuidler-oci-build-push-phase${PHASE}-wave${WAVE}-effort${EFFORT}"
mkdir -p "/tmp/workspaces/$effort_dir"
cd "/tmp/workspaces/$effort_dir"

# Clone or set up as needed
git clone --single-branch --branch "$BASE_BRANCH" "$REPO_URL" .
```

### 2. Directory Structure
```
/tmp/workspaces/
├── project-phase1-wave1-effort1/
├── project-phase1-wave1-effort2/  
├── project-phase1-wave1-effort3/
└── project-phase1-wave2-effort1/
```

### 3. Verification Steps
- Verify directory doesn't already exist
- Ensure proper permissions
- Confirm Git repository initialized
- Check branch tracking configured

## Enforcement
- **Trigger**: Before spawning any implementation agent
- **Validation**: Directory exists and is properly configured
- **Failure**: Stop orchestration if workspace setup fails

## Related Rules
- R209: Effort Directory Isolation Protocol
- R212: Phase Directory Isolation Protocol
- R191: Target Repository Configuration

---
**Status**: STUB - This rule file needs complete implementation details
**Created**: By software-factory-manager during rule library audit