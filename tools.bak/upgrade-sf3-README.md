# Software Factory 3.0 Upgrade Script Documentation

## Overview

The `upgrade-sf3.sh` script is a specialized tool for upgrading existing Software Factory 3.0 projects with the latest template improvements, bug fixes, and rule updates. Unlike the general `upgrade.sh` script, this tool is specifically designed for SF 3.0 projects and includes validation for SF 3.0 specific components.

## Key Features

### 1. SF 3.0 Specific Validation
- Verifies the project is SF 3.0 (not SF 2.0)
- Checks for State Manager agent presence
- Validates integration hierarchy in state machine
- Ensures critical rules (R288, R322, R405, R506) are present

### 2. Comprehensive Backup
- Creates timestamped backups in `.upgrade-backup/`
- **Special effort handling**: Preserves complete git history
- Creates manifest files for disaster recovery
- Backs up all configuration before changes

### 3. Smart Update System
- Only updates files that have changed
- Preserves all project work (efforts, todos, state)
- Reports what was updated vs unchanged
- Dry-run mode for safety

### 4. Critical SF 3.0 Updates
- **State Machine**: Updates with integration hierarchy fixes
- **State Manager**: Ensures agent is properly configured
- **Rules**: Updates all rules including critical SF 3.0 rules
- **Agent States**: Syncs all state-specific rules

## Usage

### Basic Usage

From your project directory:
```bash
bash /path/to/template/tools/upgrade-sf3.sh
```

### Options

```bash
--dry-run     # Show what would be updated without making changes
--verbose     # Show detailed output during upgrade
--no-backup   # Skip backup creation (NOT RECOMMENDED)
--force       # Skip confirmation prompts
--help        # Show help message
```

### Examples

#### Dry Run First (Recommended)
```bash
# See what would be changed
bash /home/vscode/software-factory-template/tools/upgrade-sf3.sh --dry-run

# If looks good, run the actual upgrade
bash /home/vscode/software-factory-template/tools/upgrade-sf3.sh --verbose
```

#### Automated Upgrade
```bash
# Skip prompts for CI/CD
bash /path/to/template/tools/upgrade-sf3.sh --force
```

## What Gets Updated

### Always Updated
- `state-machines/software-factory-3.0-state-machine.json` - Critical fixes!
- `.claude/agents/state-manager.md` - State Manager agent
- `.claude/agents/orchestrator.md` - Orchestrator with State Manager integration
- `rule-library/*.md` - All rules including R288, R322, R405, R506
- `agent-states/` - Complete directory structure
- `utilities/*.sh` - Helper scripts
- `tools/*.sh` - Tools including line-counter.sh

### Never Modified (Preserved)
- `orchestrator-state-v3.json` - Your project state
- `bug-tracking.json` - Your bug tracking
- `integration-containers.json` - Your integration data
- `efforts/` - All your implementation work
- `todos/` - Your TODO lists
- `*.md` project files - Your documentation

## What Gets Fixed

### Integration Hierarchy Enforcement
The latest state machine fixes ensure:
- Proper validation of integration hierarchy levels
- State Manager consultation at BUILD_VALIDATION
- Correct parent state requirements

### State Manager Validation
Ensures State Manager agent:
- Has proper consultation protocols (R288)
- Validates at correct hierarchy levels
- Properly integrated with orchestrator

### Critical Rules
Updates these essential SF 3.0 rules:
- **R288**: State Manager consultation protocol
- **R322**: Integration hierarchy mandatory
- **R405**: Orchestrator state transitions only
- **R506**: Absolute prohibition on pre-commit bypass

## Validation Checks

The script performs these validations:

### Pre-Upgrade
1. Verifies orchestrator-state-v3.json exists
2. Checks state machine version
3. Looks for SF 3.0 companion files
4. Validates State Manager agent presence

### Post-Upgrade
1. Confirms state machine has integration_hierarchy
2. Verifies State Manager agent installed
3. Checks critical rules are present
4. Validates state file version

## Backup and Recovery

### Backup Location
```
.upgrade-backup/
└── YYYYMMDD-HHMMSS/
    ├── orchestrator-state-v3.json
    ├── bug-tracking.json
    ├── integration-containers.json
    ├── .claude/
    ├── agent-states/
    ├── rule-library/
    ├── state-machines/
    ├── tools/
    ├── utilities/
    └── efforts-backup/
        ├── manifest.txt
        └── [complete effort copies with .git]
```

### Recovery Process
If something goes wrong:
```bash
# 1. Identify backup timestamp
ls -la .upgrade-backup/

# 2. Restore specific files
cp .upgrade-backup/TIMESTAMP/orchestrator-state-v3.json .

# 3. Or restore everything
cp -r .upgrade-backup/TIMESTAMP/* .
```

## Differences from upgrade.sh

| Feature | upgrade.sh | upgrade-sf3.sh |
|---------|-----------|----------------|
| **Target** | SF 2.0 and 3.0 | SF 3.0 only |
| **Validation** | Basic | Comprehensive SF 3.0 |
| **State Manager** | Generic copy | Validated and highlighted |
| **Integration Hierarchy** | Not checked | Explicitly validated |
| **Reporting** | Generic | SF 3.0 specific |
| **Rules Focus** | All rules | Critical SF 3.0 rules highlighted |
| **Version Marker** | Generic | SF 3.0 with fix list |

## When to Use

### Use upgrade-sf3.sh When:
- You have an SF 3.0 project
- You need the latest integration hierarchy fixes
- You want State Manager validation
- You need SF 3.0 specific reporting

### Use upgrade.sh When:
- You have an SF 2.0 project
- You need backward compatibility
- You want more configuration options
- You're upgrading from very old versions

## Troubleshooting

### "Not an SF 3.0 project" Error
- Ensure orchestrator-state-v3.json exists
- Check you're in the project root directory

### Validation Failures
- State machine missing: Will be installed
- State Manager missing: Will be installed
- Version mismatch: Will be noted for manual update

### Permission Errors
- Ensure you have write permissions
- Check disk space for backups

## Version Information

After upgrade, check `.sf-version` file:
```bash
cat .sf-version
```

Shows:
- Version upgraded to
- Template commit used
- Critical fixes applied
- Upgrade timestamp

## Best Practices

1. **Always dry-run first**
   ```bash
   upgrade-sf3.sh --dry-run
   ```

2. **Keep backups until verified**
   Don't delete `.upgrade-backup/` until you've tested

3. **Review State Manager**
   After upgrade, review State Manager agent config

4. **Validate state file**
   ```bash
   tools/validate-state-file.sh orchestrator-state-v3.json
   ```

5. **Check integration hierarchy**
   Ensure your current state respects hierarchy

## Support

For issues or questions:
1. Check the upgrade report for specific issues
2. Review validation output for failures
3. Examine `.sf-version` for upgrade details
4. Check backups if recovery needed