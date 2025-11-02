# R104: Integration Branch Creation

## Rule Definition
**ID**: R104  
**Category**: Integration  
**Criticality**: BLOCKING  
**Created**: 2025-01-29  
**Updated**: 2025-01-29  

## Description
This rule defines where and how integration branches MUST be created for waves, phases, and project completion. Integration branches are ALWAYS created in the TARGET repository being developed, NOT in the software-factory-template repository.

## Requirements

### 1. Target Repository Identification
Before ANY integration work:
```yaml
# Read target repository configuration
TARGET_REPO_PATH: $CLAUDE_PROJECT_DIR/target-repo-config.yaml
```

The configuration MUST specify:
- `repository_path`: Absolute path to the target repository
- `repository_name`: Name for directory creation
- `default_branch`: Base branch for integration (usually main/master)

### 2. Integration Workspace Creation

**CRITICAL per R250/R271**: Integration workspaces clone the target repository DIRECTLY into the integration-workspace directory. NO subdirectories.

#### Wave Integration
```bash
# Wave integration workspace structure (R250-compliant)
$CLAUDE_PROJECT_DIR/efforts/phase{N}/wave{W}/integration-workspace/
├── .git/                    # Git repository at root level
├── [project files]          # Project files at root level
├── .software-factory/       # SF metadata (R343)
└── wave-{N}-{W}-integration # Integration branch created HERE

# WRONG (deprecated):
# $CLAUDE_PROJECT_DIR/efforts/wave{W}/integration-workspace/[target-repo-name]/
```

#### Phase Integration
```bash
# Phase integration workspace structure (R250-compliant)
$CLAUDE_PROJECT_DIR/efforts/phase{N}/integration-workspace/
├── .git/                     # Git repository at root level
├── [project files]           # Project files at root level
├── .software-factory/        # SF metadata (R343)
└── phase-{N}-integration     # Integration branch created HERE

# WRONG (deprecated):
# $CLAUDE_PROJECT_DIR/efforts/phase{N}/integration-workspace/[target-repo-name]/
```

#### Project Integration
```bash
# Project integration workspace structure (R250-compliant)
$CLAUDE_PROJECT_DIR/efforts/project/integration-workspace/
├── .git/                     # Git repository at root level
├── [project files]           # Project files at root level
├── .software-factory/        # SF metadata (R343)
└── project-integration       # Integration branch created HERE

# WRONG (deprecated):
# $CLAUDE_PROJECT_DIR/efforts/project/integration-workspace/[target-repo-name]/
```

### 3. Repository Verification
Before creating ANY integration branch:
```bash
# Verify we're in the TARGET repository
if [[ ! -d ".git" ]] || [[ "$(git remote get-url origin)" != *"target-repo"* ]]; then
    echo "ERROR: Not in target repository!"
    exit 1
fi
```

### 4. Branch Creation Protocol (R250/R271-Compliant)
```bash
# CRITICAL: Use direct cloning per R250 - NO subdirectories!
INTEGRATION_DIR="$CLAUDE_PROJECT_DIR/efforts/phase{N}/wave{W}/integration-workspace"  # or phase/project
BASE_BRANCH=$(yq '.default_branch' target-repo-config.yaml)  # Usually "main"
TARGET_REPO_URL=$(yq '.repository_path' target-repo-config.yaml)

# Clone DIRECTLY into integration-workspace (R250)
git clone --single-branch --branch "$BASE_BRANCH" "$TARGET_REPO_URL" "$INTEGRATION_DIR"

# Now IN the repository (integration-workspace IS the git repo)
cd "$INTEGRATION_DIR"

# Create integration branch from base
git checkout -b "phase{N}/wave{W}/integration" "$BASE_BRANCH"  # or phase{N}/integration or project/integration
git push -u origin "phase{N}/wave{W}/integration"

# WRONG (old pattern):
# cd integration-workspace/[target-repo-name]  # Extra subdirectory NO LONGER EXISTS
```

## Common Violations
1. Creating integration branches in software-factory-template repo
2. Creating branches in SWE working copies
3. Not cloning target repository first
4. Using wrong base branch

## Validation (R250-Compliant)
```bash
# Correct integration branch location per R250
pwd  # Should show: .../integration-workspace (NO subdirectory!)
git branch --show-current  # Should show: phase{N}/wave{W}/integration (or phase/project)
git remote -v  # Should show: TARGET repository, not software-factory
ls -la .git  # Should exist (we ARE in the git repository)

# WRONG (old pattern):
# pwd  # showing: .../integration-workspace/[target-repo-name]
```

## Dependencies
- R250: Wave Integration Protocol
- R034: Phase Integration Protocol
- R206: State Machine Validation
- target-repo-config.yaml

## Enforcement
- Orchestrator MUST verify correct repository before integration
- Integration branches in wrong repository = IMMEDIATE FAILURE
- All integration work happens in TARGET repository only