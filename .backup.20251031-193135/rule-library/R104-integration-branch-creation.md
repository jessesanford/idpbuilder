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

#### Wave Integration
```bash
# Wave integration workspace structure
$CLAUDE_PROJECT_DIR/efforts/waveX/integration-workspace/
└── [target-repo-name]/  # Clone of TARGET repository
    ├── .git/
    ├── [project files]
    └── wave-X-integration  # Integration branch created HERE
```

#### Phase Integration  
```bash
# Phase integration workspace structure
$CLAUDE_PROJECT_DIR/efforts/phaseX/integration-workspace/
└── [target-repo-name]/  # Clone of TARGET repository
    ├── .git/
    ├── [project files]
    └── phase-X-integration  # Integration branch created HERE
```

#### Project Integration
```bash
# Project integration workspace structure
$CLAUDE_PROJECT_DIR/efforts/project/integration-workspace/
└── [target-repo-name]/  # Clone of TARGET repository
    ├── .git/
    ├── [project files]
    └── project-integration  # Integration branch created HERE
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

### 4. Branch Creation Protocol
```bash
# Always create from the target repo's default branch
cd $CLAUDE_PROJECT_DIR/efforts/[wave|phase|project]/integration-workspace/[target-repo-name]
git checkout main  # or master, from config
git pull origin main
git checkout -b [wave|phase|project]-X-integration
```

## Common Violations
1. Creating integration branches in software-factory-template repo
2. Creating branches in SWE working copies
3. Not cloning target repository first
4. Using wrong base branch

## Validation
```bash
# Correct integration branch location
pwd  # Should show: .../integration-workspace/[target-repo-name]
git branch --show-current  # Should show: [wave|phase|project]-X-integration
git remote -v  # Should show: TARGET repository, not software-factory
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