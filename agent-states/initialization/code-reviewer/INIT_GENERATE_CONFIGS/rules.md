# Code Reviewer - INIT_GENERATE_CONFIGS State Rules

## Purpose
Generate all required configuration files based on gathered requirements.

## Entry Criteria
- Requirements gathering complete
- Repository setup done
- All required information available

## Required Actions

### 1. Generate setup-config.yaml

**Location**: `$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs/setup-config.yaml`

**Template**:
```yaml
# Software Factory 2.0 Setup Configuration
# Generated: [TIMESTAMP]

project:
  name: "[PROJECT_NAME]"
  prefix: "[PROJECT_PREFIX]"
  description: "[PROJECT_IDEA]"
  type: "[cli|service|library|application]"

technology:
  primary_language: "[LANGUAGE]"
  languages:
    - "[LANGUAGE_1]"
    - "[LANGUAGE_2]"

  frameworks:
    - name: "[FRAMEWORK_1]"
      purpose: "[PURPOSE]"
    - name: "[FRAMEWORK_2]"
      purpose: "[PURPOSE]"

  build_system: "[BUILD_SYSTEM]"
  package_manager: "[PACKAGE_MANAGER]"

  testing:
    framework: "[TEST_FRAMEWORK]"
    coverage_target: [NUMBER]
    type: "[unit|integration|e2e|all]"

architecture:
  pattern: "[monolith|microservices|serverless|library]"
  components:
    - name: "[COMPONENT_1]"
      type: "[api|cli|library|service]"
      description: "[DESCRIPTION]"

deployment:
  environment: "[local|cloud|kubernetes|docker]"
  cloud_provider: "[aws|gcp|azure|none]"
  containerization: "[docker|podman|none]"
  orchestration: "[kubernetes|docker-compose|none]"
  ci_cd: "[github-actions|jenkins|gitlab-ci|none]"

quality:
  code_style: "[STYLE_GUIDE]"
  linting: "[LINTER]"
  security_scanning: "[enabled|disabled]"
  performance_monitoring: "[enabled|disabled]"

development:
  ide: "[vscode|intellij|vim|emacs]"
  local_environment: "[docker|native|vagrant]"
  hot_reload: "[enabled|disabled]"
```

### 2. Generate target-repo-config.yaml

**Location**: `$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs/target-repo-config.yaml`

**For Upstream Fork**:
```yaml
# Target Repository Configuration
# Generated: [TIMESTAMP]

repository:
  type: "upstream_fork"

upstream:
  url: "[UPSTREAM_URL]"
  branch: "[MAIN_BRANCH]"
  organization: "[ORG_NAME]"
  name: "[REPO_NAME]"

fork:
  url: "[FORK_URL]"
  branch: "[DEV_BRANCH]"
  owner: "[USERNAME]"

paths:
  target_repo: "efforts/[PROJECT_PREFIX]/target-repo"
  project: "efforts/[PROJECT_PREFIX]/project"

directories:
  include:
    - "src"
    - "cmd"
    - "pkg"
  exclude:
    - "vendor"
    - "node_modules"
    - ".git"

integration:
  strategy: "rebase"  # or "merge"
  pr_target: "upstream"
  branch_prefix: "sf2-"
```

**For New Repository**:
```yaml
# Target Repository Configuration
# Generated: [TIMESTAMP]

repository:
  type: "new_project"

project:
  path: "efforts/[PROJECT_PREFIX]/project"
  remote_url: "none"  # Will be added later
  branch: "main"

directories:
  source: "src"
  tests: "tests"
  docs: "docs"
  configs: "configs"

version_control:
  initialized: true
  initial_commit: true
  branch_strategy: "git-flow"  # or "trunk-based"
```

### 3. Generate Initial .claude/CLAUDE.md

**Location**: `$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX]/.claude/CLAUDE.md`

**Template**:
```markdown
# Project-Specific Claude Configuration

## Project: [PROJECT_NAME]

### Overview
[PROJECT_IDEA expanded]

### Technology Context
- Primary Language: [LANGUAGE]
- Key Frameworks: [FRAMEWORKS]
- Build System: [BUILD_SYSTEM]
- Testing: [TEST_FRAMEWORK]

### Coding Standards
Follow [LANGUAGE] best practices:
- [Standard 1]
- [Standard 2]
- [Standard 3]

### Architecture Guidelines
Pattern: [ARCHITECTURE_PATTERN]
- [Guideline 1]
- [Guideline 2]

### Development Workflow
1. All changes in feature branches
2. Tests required for all new code
3. Code review before merge
4. Documentation updates with code

### Key Commands
- Build: `[BUILD_COMMAND]`
- Test: `[TEST_COMMAND]`
- Run: `[RUN_COMMAND]`
- Lint: `[LINT_COMMAND]`

### Important Paths
- Source: `[SOURCE_PATH]`
- Tests: `[TEST_PATH]`
- Configs: `[CONFIG_PATH]`

### References
- IMPLEMENTATION-PLAN.md - Development phases
- setup-config.yaml - Project configuration
- target-repo-config.yaml - Repository settings
```

### 4. Validate Generated Files

Check each file for:
- Valid YAML syntax
- No missing required fields
- No placeholder values remaining
- Paths are correct
- Values match requirements

### 5. Update State File
Record generation status:
```json
"configs_generated": [
  {
    "file": "setup-config.yaml",
    "path": "efforts/[PREFIX]/configs/setup-config.yaml",
    "valid": true
  },
  {
    "file": "target-repo-config.yaml",
    "path": "efforts/[PREFIX]/configs/target-repo-config.yaml",
    "valid": true
  },
  {
    "file": "CLAUDE.md",
    "path": "efforts/[PREFIX]/.claude/CLAUDE.md",
    "valid": true
  }
],
"config_generation_complete": true
```

## Field Mapping from Requirements

### From requirements → setup-config.yaml
- requirements.technology.languages → technology.primary_language
- requirements.technology.frameworks → technology.frameworks
- requirements.technology.build_system → technology.build_system
- requirements.technology.test_framework → technology.testing.framework
- requirements.architecture.type → architecture.pattern
- requirements.deployment.environment → deployment.environment

### From requirements → target-repo-config.yaml
- requirements.codebase.upstream_url → upstream.url
- requirements.codebase.fork_url → fork.url
- requirements.codebase.main_branch → upstream.branch

## Exit Criteria
- All config files generated
- No placeholders remaining
- YAML syntax valid
- All fields populated
- Files saved to correct locations

## Transition
**MANDATORY**: → INIT_SYNTHESIZE_PLAN

## Common Issues to Avoid
- Empty or null values
- Invalid YAML indentation
- Missing required fields
- Incorrect file paths
- Placeholder text left in