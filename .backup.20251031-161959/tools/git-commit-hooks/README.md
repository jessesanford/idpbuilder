# Git Commit Hooks Organization

## Overview

The Software Factory 2.0 template uses a structured git hook system that automatically validates commits based on repository type. This ensures compliance with Software Factory rules and prevents common mistakes.

## Directory Structure

```
tools/git-commit-hooks/
├── master-pre-commit.sh         # Main hook that orchestrates all validations
├── effort-hooks/                # Hooks specific to effort repositories
│   ├── branch-name-validation.hook     # Validates branch naming conventions
│   └── r251-main-branch-protection.hook # Prevents commits to main branch
├── planning-hooks/              # Hooks specific to planning repository
│   └── efforts-protection.hook  # Prevents committing efforts/ directory
└── shared-hooks/                # Hooks that apply to all repositories
    └── orchestrator-state-validation.hook # Validates state and state machine files
```

## Hook Types and Application

### Planning Repository Hooks

Applied to the main Software Factory template/planning repository:

1. **efforts-protection.hook** - Prevents any files in `efforts/` directory from being committed
   - Ensures implementation code goes to target repository
   - Protects repository separation (R251)

2. **orchestrator-state-validation.hook** - Validates orchestrator-state-v3.json
   - Ensures state file is always valid
   - Validates against state machine schema

### Effort Repository Hooks

Applied to individual effort working copies:

1. **branch-name-validation.hook** - Ensures branches follow naming conventions
   - Pattern: `{prefix}/phase{X}/wave{Y}/{effort-name}`
   - Allows splits: `{effort-name}--split-{NNN}`
   - Allows fixes: `{effort-name}-fix`

2. **r251-main-branch-protection.hook** - Prevents commits to main/master branch
   - Enforces feature branch development
   - Complies with R251 repository separation rule

3. **orchestrator-state-validation.hook** - Same as planning repository

### Shared Hooks

Applied to all repository types:

1. **orchestrator-state-validation.hook** - Common validation logic
   - Validates orchestrator-state-v3.json if present
   - Validates state machine files if present
   - Automatically installs jsonschema if needed

## Installation

### Automatic Installation via upgrade.sh

The `upgrade.sh` script automatically installs hooks when upgrading a project:

```bash
./upgrade.sh /path/to/project
```

### Manual Installation

Use the unified installation script:

```bash
# Install in current directory
./utilities/install-hooks.sh

# Install in specific repository
./utilities/install-hooks.sh /path/to/repo

# Install in all effort repositories
./utilities/install-hooks.sh /path/to/project effort

# Install everywhere (planning + all efforts)
./utilities/install-hooks.sh /path/to/project all
```

### Legacy Installation Scripts

For backward compatibility, these scripts still work:

- `utilities/install-branch-validation-hook.sh` - Installs branch validation only
- `utilities/install-efforts-prevention-hook.sh` - Installs efforts protection only

## How It Works

### Master Hook Orchestration

The `master-pre-commit.sh` script:

1. **Detects Repository Type**
   - Checks remote URL for software-factory-template
   - Looks for target-repo-config.yaml
   - Examines directory structure

2. **Loads Appropriate Hooks**
   - Planning repos: efforts-protection + state validation
   - Effort repos: branch validation + R251 protection + state validation
   - General repos: state validation only

3. **Executes Validations**
   - Runs each applicable hook in sequence
   - Stops on first failure
   - Provides clear error messages

### Repository Type Detection

The system automatically detects repository type:

- **Planning Repository**: Contains software-factory-template in remote URL
- **Effort Repository**: Located in efforts/ directory or has target-repo-config.yaml
- **General Repository**: Any other git repository

## Configuration

### Environment Variables

Control hook behavior with environment variables:

```bash
# Skip specific validations (emergency use only)
SKIP_BRANCH_VALIDATION=true git commit -m "message"
SKIP_EFFORTS_PROTECTION=true git commit -m "message"
SKIP_STATE_VALIDATION=true git commit -m "message"
SKIP_MAIN_PROTECTION=true git commit -m "message"
```

**WARNING**: Use skip flags only in emergencies. They bypass critical safety checks.

### Custom Branch Prefixes

For projects requiring custom branch prefixes, create `target-repo-config.yaml`:

```yaml
branch_naming:
  project_prefix: "myproject"
```

This will require branches like: `myproject/phase1/wave1/feature-name`

## Troubleshooting

### Hook Not Running

```bash
# Check if hook is executable
ls -la .git/hooks/pre-commit

# Make executable if needed
chmod +x .git/hooks/pre-commit
```

### Wrong Repository Type Detected

```bash
# Check detected type
bash .git/hooks/pre-commit
# Look for "Repository type: ..." in output

# Force reinstall if incorrect
./utilities/install-hooks.sh . single true
```

### Python/jsonschema Issues

The orchestrator-state validation requires Python and jsonschema:

```bash
# Manual installation if auto-install fails
pip3 install --user jsonschema

# Or with system packages flag
pip3 install --user --break-system-packages jsonschema
```

### Bypassing Hooks (Emergency Only)

```bash
# Bypass all hooks (DANGEROUS - use only in emergencies)
git commit --no-verify -m "emergency commit"

# Or set environment variable
GIT_SKIP_HOOKS=true git commit -m "message"
```

## Error Messages

### Efforts Directory Protection

```
❌ ERROR: WRONG! THIS IS THE PLANNING REPO. YOU ARE NOT ALLOWED TO COMMIT EFFORT WORK
TO THE PLANNING REPO. REMOVE THESE COMMITS AND INSTEAD PUSH THEM TO THE TARGET REPO
```

**Solution**: Navigate to effort directory and commit there instead

### Branch Name Validation

```
❌ ERROR: BRANCH NAME VALIDATION FAILED!
WRONG! YOUR BRANCH NAME SHOULD BE phase1/wave1/feature-name
```

**Solution**: Create correct branch and move changes there

### R251 Main Branch Protection

```
❌ ERROR: R251 VIOLATION: Cannot commit to main branch in effort repository!
```

**Solution**: Create feature branch for your changes

### State Validation

```
❌ ERROR: orchestrator-state-v3.json validation failed!
```

**Solution**: Fix JSON syntax or schema violations

## Hook Development

### Adding New Hooks

1. Create hook file in appropriate directory:
   - `effort-hooks/` for effort-specific validations
   - `planning-hooks/` for planning-specific validations
   - `shared-hooks/` for universal validations

2. Follow naming convention: `{feature}.hook`

3. Use standard structure:
```bash
#!/bin/bash
set -euo pipefail

# Color codes
RED='\033[0;31m'
# ... other colors

# Main validation logic
main() {
    # Validation code here
    exit 0  # Success
    exit 1  # Failure
}

main "$@"
```

4. Update `master-pre-commit.sh` to include new hook

### Testing Hooks

Test script for validating hooks:

```bash
# Create test commit
echo "test" > test.txt
git add test.txt
git commit -m "test"

# Check if hook ran
# Look for validation output
```

## Related Documentation

- [Efforts Directory Protection](../../docs/EFFORTS-DIRECTORY-PROTECTION.md)
- [Branch Naming Validation](../../docs/BRANCH-NAMING-VALIDATION-REPORT.md)
- [Rule R251 - Repository Separation](../../rule-library/R251-REPOSITORY-SEPARATION-LAW.md)
- [Rule R410 - Mandatory Branch Validation](../../rule-library/R410-mandatory-branch-validation-hook.md)

## Summary

The organized hook system provides:

- ✅ Automatic repository type detection
- ✅ Appropriate validations per repository type
- ✅ Clear separation of concerns
- ✅ Easy maintenance and updates
- ✅ Legacy compatibility
- ✅ Emergency bypass options

This ensures Software Factory compliance while maintaining flexibility for different repository types.