# Software Factory 2.0 - Integration Directory Structure

## 🔴 CRITICAL: Repository Separation

**THE SOFTWARE FACTORY INSTANCE IS NOT THE TARGET REPOSITORY!**

- **Software Factory Instance**: Planning, configuration, state management
- **Target Repository**: Actual source code being developed
- **Integration Locations**: ALWAYS under `/efforts/` directory

## Directory Structure Overview

```
/                                    # System root
├── workspaces/                      # Development workspaces
│   └── software-factory-2.0/        # SF INSTANCE (planning only)
│       ├── orchestrator-state.json  # State management
│       ├── target-repo-config.yaml  # Target repo configuration
│       ├── agent-states/            # Agent-specific rules
│       └── utilities/               # Helper scripts
│
└── efforts/                         # ALL CODE DEVELOPMENT HERE
    ├── phase1/
    │   ├── wave1/
    │   │   ├── core-api-types/     # Effort workspace (target repo clone)
    │   │   ├── webhook-framework/  # Effort workspace (target repo clone)
    │   │   ├── controller-base/    # Effort workspace (target repo clone)
    │   │   └── integration-workspace/ # Wave integration (target repo clone)
    │   │
    │   ├── wave2/
    │   │   ├── rbac-controller/    # Effort workspace (target repo clone)
    │   │   ├── validation-hooks/   # Effort workspace (target repo clone)
    │   │   └── integration-workspace/ # Wave integration (target repo clone)
    │   │
    │   └── phase-integration-workspace/ # Phase integration (target repo clone)
    │
    └── phase2/
        ├── wave1/
        │   ├── cluster-controller/  # Effort workspace (target repo clone)
        │   ├── api-extensions/      # Effort workspace (target repo clone)
        │   └── integration-workspace/ # Wave integration (target repo clone)
        │
        └── phase-integration-workspace/ # Phase integration (target repo clone)
```

## Integration Types and Locations

### 1. Wave Integration
**Purpose**: Integrate all efforts from a single wave
**Location**: `/efforts/phase{X}/wave{Y}/integration-workspace/`
**When**: After all wave efforts are complete and reviewed
**Branch**: `{project-prefix}/phase{X}/wave{Y}-integration`

```bash
# Example: Phase 1, Wave 1 Integration
cd /efforts/phase1/wave1/integration-workspace
git clone "$TARGET_REPO_URL" .
git checkout -b "tmc-workspace/phase1/wave1-integration"
# Merge all wave1 efforts
```

### 2. Phase Integration
**Purpose**: Integrate all waves from a phase
**Location**: `/efforts/phase{X}/phase-integration-workspace/`
**When**: After all phase waves are complete and integrated
**Branch**: `{project-prefix}/phase{X}-integration`

```bash
# Example: Phase 1 Integration
cd /efforts/phase1/phase-integration-workspace
git clone "$TARGET_REPO_URL" .
git checkout -b "tmc-workspace/phase1-integration"
# Merge all phase1 wave integration branches
```

### 3. Final Integration
**Purpose**: Integrate all phases into main branch
**Location**: `/efforts/final-integration-workspace/`
**When**: After all phases complete
**Branch**: `{project-prefix}/main-integration-{timestamp}`

```bash
# Example: Final Integration
cd /efforts/final-integration-workspace
git clone "$TARGET_REPO_URL" .
git checkout -b "tmc-workspace/main-integration-20250825"
# Merge all phase integration branches
```

## Integration Workflow

### Wave Integration Process

1. **Create Integration Workspace**
   ```bash
   WAVE_DIR="/efforts/phase${PHASE}/wave${WAVE}"
   INTEGRATION_DIR="${WAVE_DIR}/integration-workspace"
   mkdir -p "$INTEGRATION_DIR"
   cd "$INTEGRATION_DIR"
   ```

2. **Clone Target Repository**
   ```bash
   TARGET_REPO_URL=$(yq '.target_repository.url' /workspaces/software-factory-2.0/target-repo-config.yaml)
   git clone "$TARGET_REPO_URL" .
   ```

3. **Create Integration Branch**
   ```bash
   source /workspaces/software-factory-2.0/utilities/branch-naming-helpers.sh
   INTEGRATION_BRANCH=$(get_wave_integration_branch_name "$PHASE" "$WAVE")
   git checkout -b "$INTEGRATION_BRANCH"
   ```

4. **Merge Effort Branches**
   ```bash
   for effort_dir in /efforts/phase${PHASE}/wave${WAVE}/*/; do
       [[ "$effort_dir" == *"integration-workspace"* ]] && continue
       effort=$(basename "$effort_dir")
       EFFORT_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$effort")
       git fetch origin "$EFFORT_BRANCH"
       git merge "origin/$EFFORT_BRANCH" --no-ff -m "integrate: $effort into wave${WAVE}"
   done
   ```

5. **Validate and Push**
   ```bash
   # Run tests
   make test
   
   # Check size compliance
   $PROJECT_ROOT/tools/line-counter.sh -c "$INTEGRATION_BRANCH"
   
   # Push to remote
   git push -u origin "$INTEGRATION_BRANCH"
   ```

## Directory Creation Timing

### During SETUP_EFFORT_INFRASTRUCTURE State
- Create effort directories: `/efforts/phase{X}/wave{Y}/{effort-name}/`
- Each effort gets its own target repo clone
- Branches created with project prefix

### During INTEGRATION State
- Create wave integration: `/efforts/phase{X}/wave{Y}/integration-workspace/`
- Fresh clone of target repo for clean integration
- Never reuse effort workspaces for integration

### During Phase Completion
- Create phase integration: `/efforts/phase{X}/phase-integration-workspace/`
- Merge all wave integration branches
- Final phase-level validation

## Critical Rules

### ❌ NEVER DO THIS:
- Never integrate in the Software Factory instance directory
- Never merge code in `/workspaces/software-factory-2.0/`
- Never create integration branches in planning repositories
- Never mix planning files with source code

### ✅ ALWAYS DO THIS:
- Always use `/efforts/` for all code work
- Always create fresh clones for integration
- Always maintain separation between SF instance and target repo
- Always use project prefix in branch names

## Validation Checklist

Before any integration:
- [ ] Working in `/efforts/` directory, not SF instance
- [ ] Fresh clone of target repository created
- [ ] Integration workspace in correct location
- [ ] Branch name includes project prefix
- [ ] All effort branches fetched from remote
- [ ] Tests passing after integration
- [ ] Size compliance verified

## Environment Variables

```bash
# Set by orchestrator
export SF_INSTANCE_DIR="/workspaces/software-factory-2.0"
export EFFORTS_ROOT="/efforts"
export TARGET_REPO_URL="https://github.com/example/target-repo.git"
export PROJECT_PREFIX="tmc-workspace"

# Integration paths
export WAVE_INTEGRATION="${EFFORTS_ROOT}/phase${PHASE}/wave${WAVE}/integration-workspace"
export PHASE_INTEGRATION="${EFFORTS_ROOT}/phase${PHASE}/phase-integration-workspace"
```

## Common Errors and Solutions

### Error: "fatal: not a git repository"
**Cause**: You're in the SF instance directory
**Solution**: cd to `/efforts/phase{X}/wave{Y}/integration-workspace/`

### Error: "orchestrator-state.json" in git status
**Cause**: You're in the wrong repository
**Solution**: You're in SF instance, not target repo

### Error: "Branch already exists"
**Cause**: Integration workspace wasn't fresh
**Solution**: Remove and recreate integration-workspace directory

## Grading Impact

**AUTOMATIC FAILURE** if:
- Integration attempted in SF instance directory
- Code merged in planning repository
- Integration workspace not under `/efforts/`
- Effort code mixed with SF configuration

Remember: The Software Factory orchestrates work ON the target repository, it never CONTAINS target repository code!