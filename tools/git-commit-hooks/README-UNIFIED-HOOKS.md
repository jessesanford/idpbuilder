# Software Factory Unified Git Hook System

**Location**: `/tools/git-commit-hooks/`
**Purpose**: Centralized git hook management for all Software Factory repositories
**Version**: 3.0
**Last Updated**: 2025-11-01

---

## Overview

All git hooks for Software Factory are stored and managed from this single directory. The system automatically detects repository type (planning vs. effort) and Software Factory version (2.0 vs. 3.0) to apply appropriate validation.

---

## Directory Structure

```
tools/git-commit-hooks/
├── master-pre-commit.sh           # Master orchestrator hook (auto-detects repo type)
├── planning-pre-commit.sh         # Standalone planning repo hook
├── effort-pre-commit.sh           # Standalone effort repo hook
├── planning-hooks/                # Hooks specific to planning repositories
│   ├── branch-name-validation.hook
│   ├── efforts-protection.hook
│   └── README.md
├── effort-hooks/                  # Hooks specific to effort working copies
│   ├── branch-name-validation.hook
│   └── r251-main-branch-protection.hook
└── shared-hooks/                  # Hooks used by ALL repositories
    ├── metadata-placement-validation.hook
    └── orchestrator-state-validation.hook
```

---

## Hook Types

### 1. Master Pre-Commit Hook (`master-pre-commit.sh`)

**Purpose**: Universal pre-commit hook that auto-detects repository type and SF version

**Features**:
- Detects SF 2.0 vs SF 3.0 automatically
- Identifies repository type (template, planning, effort, general)
- Routes to appropriate validation hooks
- Validates state files (orchestrator-state-v3.json, bug-tracking.json, etc.)
- Runs R550 plan path consistency validation (SF 3.0)
- Enforces R506 (no --no-verify bypass)

**Usage**: This hook is installed in `.git/hooks/pre-commit` for ALL repositories

**Installation**:
```bash
# Automatic (via upgrade.sh):
bash tools/upgrade.sh /path/to/target/repo

# Manual:
cp tools/git-commit-hooks/master-pre-commit.sh /path/to/repo/.git/hooks/pre-commit
chmod +x /path/to/repo/.git/hooks/pre-commit
```

---

### 2. Planning Repository Hooks

**Location**: `tools/git-commit-hooks/planning-hooks/`

**Purpose**: Validation specific to planning repositories (where orchestrator runs)

#### Hooks:

1. **`branch-name-validation.hook`**
   - Validates branch naming conventions
   - Ensures feature branches follow patterns
   - Prevents commits to protected branches (if configured)

2. **`efforts-protection.hook`**
   - Prevents direct modification of `efforts/` directory
   - Enforces R309 (no efforts in SF repo)
   - Ensures effort infrastructure is created properly

#### When Applied:
- Repository has `orchestrator-state-v3.json` at root
- Repository matches planning patterns (PROJECT-IMPLEMENTATION-PLAN.md, etc.)
- Repository is NOT in `efforts/` subdirectory

---

### 3. Effort Repository Hooks

**Location**: `tools/git-commit-hooks/effort-hooks/`

**Purpose**: Validation specific to effort working copies

#### Hooks:

1. **`branch-name-validation.hook`**
   - Validates effort branch naming (phase/wave/effort patterns)
   - Ensures branches match planning structure
   - Enforces cascade branch naming

2. **`r251-main-branch-protection.hook`**
   - Prevents commits to `main` branch in effort repos
   - Enforces R251 (Repository Separation)
   - Directs developers to feature branches

#### When Applied:
- Repository is in `efforts/` directory structure
- Repository has `target-repo-config.yaml`
- Repository path matches `efforts/phase*/wave*` pattern

---

### 4. Shared Hooks

**Location**: `tools/git-commit-hooks/shared-hooks/`

**Purpose**: Validation applied to ALL repositories

#### Hooks:

1. **`metadata-placement-validation.hook`**
   - Enforces R383 (SUPREME LAW - Metadata File Organization)
   - Validates all metadata goes to `.software-factory/` directory
   - Prevents metadata pollution in repository root
   - Checks for R383 compliance markers in commits

2. **`orchestrator-state-validation.hook`**
   - Validates `orchestrator-state-v3.json` schema (SF 3.0)
   - Validates `orchestrator-state-v3.json` schema (SF 2.0)
   - Checks state machine compliance
   - Ensures state file integrity before commit

#### When Applied:
- ALL repositories (planning, effort, general, template)
- Always runs before repo-specific hooks

---

## Installation Methods

### Method 1: Automatic (via upgrade.sh)

The preferred method. Automatically detects repo type and installs correct hooks:

```bash
cd /home/vscode/software-factory-template
bash tools/upgrade.sh /path/to/target/repository
```

**What it does**:
1. Syncs entire `tools/git-commit-hooks/` directory to target
2. Installs `master-pre-commit.sh` as `.git/hooks/pre-commit`
3. Makes hooks executable
4. Installs effort hooks in all effort working copies

---

### Method 2: Manual Installation (Planning Repo)

For planning repositories:

```bash
TARGET_REPO="/path/to/planning/repo"

# Copy hooks directory
cp -r tools/git-commit-hooks "$TARGET_REPO/tools/"

# Install master hook
cp tools/git-commit-hooks/master-pre-commit.sh "$TARGET_REPO/.git/hooks/pre-commit"
chmod +x "$TARGET_REPO/.git/hooks/pre-commit"
```

---

### Method 3: Manual Installation (Effort Repo)

For effort working copies:

```bash
EFFORT_REPO="/path/to/efforts/phase1/wave1/effort-name"

# Copy hooks directory (if not already present from planning repo sync)
cp -r tools/git-commit-hooks "$EFFORT_REPO/tools/" 2>/dev/null || true

# Install master hook
cp tools/git-commit-hooks/master-pre-commit.sh "$EFFORT_REPO/.git/hooks/pre-commit"
chmod +x "$EFFORT_REPO/.git/hooks/pre-commit"
```

---

## Hook Execution Flow

### Planning Repository

```
Commit attempt
    ↓
.git/hooks/pre-commit (master-pre-commit.sh)
    ↓
Detect: repo_type=planning, sf_version=3.0
    ↓
Run shared-hooks/metadata-placement-validation.hook
    ↓
Run shared-hooks/orchestrator-state-validation.hook
    ↓
Run planning-hooks/branch-name-validation.hook
    ↓
Run planning-hooks/efforts-protection.hook
    ↓
Pass/Fail → Commit allowed/blocked
```

### Effort Repository

```
Commit attempt
    ↓
.git/hooks/pre-commit (master-pre-commit.sh)
    ↓
Detect: repo_type=effort, sf_version=3.0
    ↓
Run shared-hooks/metadata-placement-validation.hook
    ↓
Run shared-hooks/orchestrator-state-validation.hook
    ↓
Run effort-hooks/branch-name-validation.hook
    ↓
Run effort-hooks/r251-main-branch-protection.hook
    ↓
Pass/Fail → Commit allowed/blocked
```

---

## Repository Type Detection

The `master-pre-commit.sh` detects repository type using this logic:

```bash
1. Template Check:
   - Has .template-repository file?
   → template (allows multi-branch development)

2. Planning Check:
   - Has orchestrator-state-v3.json at root?
   - Has orchestrator-state-v3.json at root? (SF 2.0)
   - Matches planning name patterns?
   - Has PROJECT-IMPLEMENTATION-PLAN.md?
   → planning

3. Effort Check:
   - Path contains "efforts/"?
   - Has target-repo-config.yaml?
   - Path matches "efforts/phase*/wave*"?
   → effort

4. Default:
   → general (minimal validation)
```

---

## Software Factory Version Detection

```bash
1. SF 3.0 Check:
   - Has orchestrator-state-v3.json?
   - Has bug-tracking.json?
   - Has integration-containers.json?
   → SF 3.0 (use new validation)

2. SF 2.0 Check:
   - Has orchestrator-state-v3.json?
   → SF 2.0 (use legacy validation)

3. Default:
   → unknown (minimal validation)
```

---

## Rule Enforcement

### Critical Rules Enforced by Hooks:

- **R383**: Metadata File Organization (SUPREME LAW)
  - All metadata → `.software-factory/` directory
  - Enforced by: `shared-hooks/metadata-placement-validation.hook`

- **R251**: Repository Separation
  - No commits to `main` in effort repos
  - Enforced by: `effort-hooks/r251-main-branch-protection.hook`

- **R309**: No Efforts in SF Repo
  - Prevents `efforts/` directory in planning repo
  - Enforced by: `planning-hooks/efforts-protection.hook`

- **R506**: Absolute Prohibition on Pre-Commit Bypass
  - NEVER use `--no-verify`
  - All hooks warn about R506 on failure

- **R517**: Universal State Manager Consultation Law
  - State transitions must consult state manager
  - Enforced in SF 3.0 state validation

- **R550**: Plan Path Consistency and Discovery
  - All planning files tracked in state
  - Enforced by: SF 3.0 validation in `master-pre-commit.sh`

---

## Troubleshooting

### Hook Not Running

```bash
# Check if hook is installed
ls -la /path/to/repo/.git/hooks/pre-commit

# Check if executable
stat /path/to/repo/.git/hooks/pre-commit

# Make executable if needed
chmod +x /path/to/repo/.git/hooks/pre-commit
```

### Hook Failing Incorrectly

```bash
# Check repository type detection
cd /path/to/repo
bash -x .git/hooks/pre-commit 2>&1 | grep "Repository type"

# Check SF version detection
bash -x .git/hooks/pre-commit 2>&1 | grep "Software Factory Version"
```

### R383 Validation Failing

```bash
# Check metadata file locations
find /path/to/repo -name "*.json" -o -name "*.yaml" -o -name "*.md" | grep -v ".git" | grep -v ".software-factory"

# Move files to correct location
mkdir -p /path/to/repo/.software-factory/
mv incorrect-location-file.json /path/to/repo/.software-factory/
```

### State File Validation Failing

```bash
# Run validator manually
bash /path/to/repo/tools/validate-state-file.sh /path/to/repo/orchestrator-state-v3.json

# Check schema
cat /path/to/repo/schemas/orchestrator-state-v3.schema.json | jq .
```

---

## Development Guidelines

### Adding a New Hook

1. **Determine Type**: Planning, effort, or shared?

2. **Create Hook File**:
   ```bash
   touch tools/git-commit-hooks/{type}-hooks/new-validation.hook
   chmod +x tools/git-commit-hooks/{type}-hooks/new-validation.hook
   ```

3. **Follow Template**:
   ```bash
   #!/bin/bash
   set -euo pipefail

   # Hook description and purpose
   # Rule enforcement: R###

   # Validation logic here

   exit 0  # Pass
   exit 1  # Fail
   ```

4. **Update Master Hook**:
   - Add `run_hook` call in appropriate section of `master-pre-commit.sh`

5. **Test**:
   ```bash
   # Test in isolation
   bash tools/git-commit-hooks/{type}-hooks/new-validation.hook

   # Test via master hook
   bash tools/git-commit-hooks/master-pre-commit.sh
   ```

6. **Document**: Update this README

### Modifying Existing Hook

1. **Test Current Behavior**: Document what currently works
2. **Make Changes**: Edit hook file
3. **Test Changes**: Run in test repository
4. **Update Documentation**: Update this README if behavior changes
5. **Sync to Projects**: Run upgrade.sh on affected projects

---

## Maintenance

### Syncing Hooks to Existing Projects

After modifying hooks in the template:

```bash
# Sync to a single project
bash tools/upgrade.sh /path/to/project

# Sync to all projects (example)
for project in ~/projects/*; do
    if [ -f "$project/orchestrator-state-v3.json" ]; then
        echo "Syncing hooks to: $project"
        bash tools/upgrade.sh "$project"
    fi
done
```

### Verifying Hook Installation

```bash
# Check hook is present and correct
TARGET="/path/to/repo"
if diff "$TARGET/.git/hooks/pre-commit" \
        "tools/git-commit-hooks/master-pre-commit.sh" >/dev/null 2>&1; then
    echo "✅ Hook is up to date"
else
    echo "⚠️  Hook differs from template"
fi
```

---

## Migration from Legacy Hooks

**Old Locations** (deprecated):
- `.githooks/pre-commit` - R517-only validation
- `docs/hooks/pre-commit` - Comprehensive but outdated

**New Location**:
- `tools/git-commit-hooks/master-pre-commit.sh` - Unified system

**Migration**:
1. Run `bash tools/upgrade.sh /path/to/repo`
2. Old hooks are automatically replaced
3. New unified system is installed

---

## Related Documentation

- [Rule Library](/home/vscode/software-factory-template/rule-library/README.md)
- [R506 - Absolute Prohibition on Pre-Commit Bypass](/home/vscode/software-factory-template/rule-library/R506-ABSOLUTE-PROHIBITION-PRE-COMMIT-BYPASS-SUPREME-LAW.md)
- [R383 - Metadata File Organization](/home/vscode/software-factory-template/rule-library/R383-software-factory-metadata-file-organization.md)
- [R517 - State Manager Consultation Law](/home/vscode/software-factory-template/rule-library/R517-universal-state-manager-consultation-law.md)
- [Upgrade Script](/home/vscode/software-factory-template/tools/upgrade.sh)

---

## Summary

**One Hook System, All Repositories**:
- **Storage**: `tools/git-commit-hooks/`
- **Installation**: `.git/hooks/pre-commit` (via master-pre-commit.sh)
- **Auto-Detection**: Repository type and SF version
- **Consistency**: Same storage, same installation, same validation

**Key Principle**: All hooks for ALL repository types are managed from this single directory. No exceptions, no special cases, no legacy locations.
