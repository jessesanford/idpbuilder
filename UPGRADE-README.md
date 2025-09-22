# Software Factory 2.0 Upgrade Tool

## Purpose
The `upgrade.sh` script allows you to update an existing Software Factory 2.0 instance with the latest rules, agent configurations, and utilities from the template, while preserving all work in progress.

## What Gets Updated

### In Your Project Directory
- ✅ **Rule Library** (`rule-library/*.md`) - All rule definitions
- ✅ **Agent Configurations** (`.claude/agents/*.md`) - Agent behaviors and instructions
- ✅ **Command Configurations** (`.claude/commands/*.md`) - Command definitions
- ✅ **Utilities** (`utilities/*.sh`) - Helper scripts
- ✅ **Critical Files** - Grading rubric, acknowledgment rules, etc.
- ✅ **Claude Settings** (`.claude/settings.json`)

### In Global ~/.claude Directory
The script also updates your global Claude configurations:
- ✅ **~/.claude/agents/** - Updated agent configurations for global access
- ✅ **~/.claude/commands/** - Updated command configurations
- ✅ **~/.claude/utilities/** - Updated utility scripts
- ✅ **~/.claude/settings.json** - Updated settings (with prompt if exists)

## What Gets Preserved
- ✅ **efforts/** - All your implementation work
- ✅ **todos/** - All TODO state files  
- ✅ **orchestrator-state.json** - Current orchestrator state
- ✅ **project-config.yaml** - Your project configuration
- ✅ **target-repo-config.yaml** - Target repository settings
- ✅ **.git/** - Your git repository
- ✅ **checkpoints/** - Any saved checkpoints
- ✅ **snapshots/** - Any saved snapshots
- ✅ **work-logs/** - All work logs

## Usage

### Basic Upgrade
```bash
# From the SF template directory
./upgrade.sh /path/to/your/project
```

### With Configuration File
If you have your original `setup-config.yaml`, you can use it to populate variables:
```bash
./upgrade.sh /path/to/your/project --config setup-config.yaml
```

### Dry Run (Preview Changes)
See what would be updated without making changes:
```bash
./upgrade.sh /path/to/your/project --dry-run
```

### Skip Backup
By default, a backup is created. To skip:
```bash
./upgrade.sh /path/to/your/project --no-backup
```

### Force (Skip Confirmations)
```bash
./upgrade.sh /path/to/your/project --force
```

## Examples

### Typical Upgrade Workflow

1. **Check what would be updated:**
   ```bash
   ./upgrade.sh /workspaces/my-project --dry-run
   ```

2. **Run the upgrade with your config:**
   ```bash
   ./upgrade.sh /workspaces/my-project --config my-setup-config.yaml
   ```

3. **Verify the upgrade:**
   ```bash
   cd /workspaces/my-project
   ls -la rule-library/  # Check new rules
   cat .sf-version       # Check upgrade timestamp
   ```

### Upgrading Multiple Projects
```bash
for project in /workspaces/project1 /workspaces/project2; do
    ./upgrade.sh "$project" --force
done
```

## Configuration File Variables

If you provide a `--config` file, the following variables will be substituted in the updated files:

- `${PROJECT_NAME}` - Your project name
- `${TARGET_REPO_URL}` - Target repository URL
- `${TARGET_BASE_BRANCH}` - Target repository base branch
- `${GITHUB_URL}` - GitHub repository URL
- `${PROJECT_DESC}` - Project description

## Safety Features

### Automatic Backup
By default, the script creates a backup of your instance (excluding efforts/ to save space):
```
/workspaces/my-project.backup.20240821-143022/
```

### Preservation of Work
The script temporarily preserves your work files, updates the rules/configs, then restores your work files. This ensures no loss of:
- Implementation efforts
- TODO states
- Orchestrator state
- Git history

### Dry Run Mode
Use `--dry-run` to preview all changes without modifying anything.

## Troubleshooting

### "Target doesn't appear to be a Software Factory instance"
This warning appears if `project-config.yaml` is missing. You can continue anyway if you're sure it's an SF instance.

### "Permission denied"
Make sure the upgrade script is executable:
```bash
chmod +x upgrade.sh
```

### Restore from Backup
If something goes wrong, restore from the automatic backup:
```bash
# Find your backup
ls -la /workspaces/my-project.backup.*

# Restore it
mv /workspaces/my-project /workspaces/my-project.broken
mv /workspaces/my-project.backup.20240821-143022 /workspaces/my-project
```

## When to Upgrade

Run the upgrade when:
- New rules are added to the template (like R196, R197, R198)
- Agent behaviors are improved
- Bug fixes are made to utilities
- Command definitions are updated

## What Happens During Upgrade

1. **Backup Creation** (unless --no-backup)
   - Copies everything except efforts/ and .git/
   - Creates efforts-list.txt for reference

2. **File Preservation**
   - Temporarily saves work files to /tmp
   - Includes state files, todos, configs

3. **Update Process**
   - Copies new rules from template
   - Updates agent configurations in project
   - Updates utilities and commands
   - Performs variable substitution if config provided

4. **Global Updates**
   - Copies updated agents to ~/.claude/agents/
   - Copies updated commands to ~/.claude/commands/
   - Copies updated utilities to ~/.claude/utilities/
   - Updates ~/.claude/settings.json (with prompt)

5. **Restoration**
   - Restores all preserved work files
   - Maintains directory structure
   - Preserves git state

6. **Verification**
   - Creates .sf-version marker
   - Shows summary of changes including global updates

## Best Practices

1. **Always do a dry run first:**
   ```bash
   ./upgrade.sh /path/to/project --dry-run
   ```

2. **Keep your setup-config.yaml:**
   Save your original setup configuration for future upgrades.

3. **Commit before upgrading:**
   ```bash
   cd /workspaces/my-project
   git add -A
   git commit -m "checkpoint: Before SF upgrade"
   git push
   ```

4. **Test after upgrading:**
   Run a simple command to ensure everything works.

5. **Review new rules:**
   Check `rule-library/` for new rules that might affect your workflow.

## Version Tracking

After upgrade, check the version marker:
```bash
cat /workspaces/my-project/.sf-version
```

Shows:
- Upgrade timestamp
- Source template path

## Getting Help

If you encounter issues:
1. Check the backup was created successfully
2. Review what was preserved vs updated
3. Ensure you have write permissions
4. Try with --dry-run first
5. Check the error messages for specific issues

## Advanced Usage

### Custom Preservation
To preserve additional files, modify the `preserve_items` array in the script:
```bash
local preserve_items=(
    "orchestrator-state.json"
    "todos"
    "efforts"
    "my-custom-file.txt"  # Add your files
)
```

### Selective Updates
To update only specific components, comment out sections in `main_upgrade()`:
```bash
# Comment out what you don't want to update
# update_rule_library
update_agent_configs
# update_utilities
```

## Safety Notes

- The script NEVER deletes efforts/ directory
- Work in progress is always preserved
- Backups are created by default
- Dry run mode available for safety
- All operations are logged to console