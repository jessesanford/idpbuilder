# UPGRADE INSTRUCTIONS - idpbuilder-oci-push-planning

## Current State Analysis

**Project**: idpbuilder-oci-push-planning
**Current State**: PR_PLAN_CREATION
**SF Version**: 3.0.0
**Status**: Ready for upgrade with critical bug fixes needed

## Why Upgrade Now?

The Software Factory template was just updated with **CRITICAL FIXES**:

1. **Integration Hierarchy Enforcement** - Fixed validation bugs
2. **State Manager Validation** - Now properly validates at BUILD_VALIDATION level
3. **Enhanced Tests** - Better coverage of state transitions
4. **Critical Rules** - R288, R322, R405, R506 updates

Your project is currently in `PR_PLAN_CREATION` state, which is a safe point for upgrade.

## Pre-Upgrade Checklist

- [ ] Current state is PR_PLAN_CREATION (confirmed)
- [ ] No uncommitted changes in efforts/ directories
- [ ] Template repository is at `/home/vscode/software-factory-template`
- [ ] You have at least 500MB free disk space for backups

## Upgrade Steps

### Step 1: Commit Current Work (If Any)

```bash
# Check for uncommitted changes
git status

# If you have changes, commit them
git add -A
git commit -m "chore: save work before SF 3.0 upgrade"
git push
```

### Step 2: Manual Backup (Recommended)

Even though the script creates backups, make a manual one:

```bash
# Create manual safety backup
cp orchestrator-state-v3.json orchestrator-state-v3.json.pre-upgrade
cp -r .claude .claude.pre-upgrade
```

### Step 3: Dry Run First

**ALWAYS** do a dry run to see what will change:

```bash
# Run from your project directory
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Dry run to preview changes
bash /home/vscode/software-factory-template/tools/upgrade-sf3.sh --dry-run
```

Review the output carefully. You should see:
- State machine will be updated
- State Manager agent will be updated/validated
- Rules will be updated
- Agent states will be synced

### Step 4: Run Actual Upgrade

If the dry run looks good:

```bash
# Run the actual upgrade with verbose output
bash /home/vscode/software-factory-template/tools/upgrade-sf3.sh --verbose
```

The script will:
1. Validate this is an SF 3.0 project ✓
2. Create comprehensive backup in `.upgrade-backup/`
3. Update state machine with fixes
4. Update State Manager agent
5. Update all rules including R288, R322, R405, R506
6. Sync agent states
7. Preserve all your work

### Step 5: Verify Critical Updates

After upgrade, verify the critical fixes are in place:

```bash
# 1. Check State Manager agent exists
ls -la .claude/agents/state-manager.md

# 2. Verify integration hierarchy in state machine
jq '.integration_hierarchy' state-machines/software-factory-3.0-state-machine.json

# 3. Check critical rules
ls -la rule-library/R288*.md rule-library/R322*.md

# 4. Validate state file
tools/validate-state-file.sh orchestrator-state-v3.json
```

### Step 6: Update State File Version (If Needed)

The upgrade script will note if version needs updating:

```bash
# Check current version
jq -r '.state_machine.version' orchestrator-state-v3.json

# If not 3.0.0, update it
jq '.state_machine.version = "3.0.0"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json
```

### Step 7: Commit Upgrade Changes

```bash
# Add all updated files
git add -A

# Commit with detailed message
git commit -m "upgrade: SF 3.0 template with critical fixes

- Integration hierarchy enforcement fixed
- State Manager validation at BUILD_VALIDATION
- Enhanced test coverage
- Rules R288, R322, R405, R506 updated
- All agent states synchronized

Template commit: $(cd /home/vscode/software-factory-template && git rev-parse --short HEAD)"

git push
```

## Post-Upgrade Validation

### 1. State Machine Validation

```bash
# Ensure state machine is valid
jq -e . state-machines/software-factory-3.0-state-machine.json > /dev/null && echo "✓ State machine valid"

# Check integration hierarchy exists
jq -e '.integration_hierarchy' state-machines/software-factory-3.0-state-machine.json > /dev/null && echo "✓ Integration hierarchy present"
```

### 2. Current State Validation

```bash
# Your current state should still be valid
jq -r '.current_state' orchestrator-state-v3.json
# Should show: PR_PLAN_CREATION

# Validate entire state file
tools/validate-state-file.sh orchestrator-state-v3.json
```

### 3. Test State Manager Integration

The State Manager should now be properly integrated. When you continue:

```bash
/continue-software-factory
```

The orchestrator should consult the State Manager for the next transition.

## Rollback Procedure (If Needed)

If something goes wrong, you have multiple recovery options:

### Option 1: Restore from Automatic Backup

```bash
# List available backups
ls -la .upgrade-backup/

# Restore from latest backup (example timestamp)
cp -r .upgrade-backup/20251030-HHMMSS/* .
```

### Option 2: Restore from Manual Backup

```bash
# Restore manual backups
cp orchestrator-state-v3.json.pre-upgrade orchestrator-state-v3.json
cp -r .claude.pre-upgrade .claude
```

### Option 3: Git Reset

```bash
# If you committed before upgrade
git reset --hard HEAD~1
```

## Expected Improvements After Upgrade

After successful upgrade, you should have:

1. **Fixed State Transitions**
   - Integration hierarchy properly enforced
   - State Manager consulted at correct levels
   - No more validation errors at BUILD_VALIDATION

2. **Better Error Handling**
   - Clear error messages for state violations
   - Proper validation at each hierarchy level
   - State Manager guidance for transitions

3. **Enhanced Rules**
   - R288: State Manager consultation required
   - R322: Integration hierarchy mandatory
   - R405: Automated state transitions
   - R506: Pre-commit bypass protection

## Troubleshooting

### "Not an SF 3.0 project" Error
- Ensure you're in `/home/vscode/workspaces/idpbuilder-oci-push-planning`
- Check orchestrator-state-v3.json exists

### Validation Failures After Upgrade
- Run: `tools/validate-state-file.sh orchestrator-state-v3.json`
- Check error messages for specific issues
- Ensure state file version is 3.0.0

### State Manager Not Working
- Verify: `ls -la .claude/agents/state-manager.md`
- Check orchestrator includes State Manager calls
- Review R288 for consultation requirements

## Next Steps After Upgrade

1. **Continue with current state**:
   ```bash
   /continue-software-factory
   ```

2. **The orchestrator should**:
   - Complete PR_PLAN_CREATION
   - Consult State Manager for next transition
   - Move to CREATE_MASTER_PR or next appropriate state

3. **Monitor for improvements**:
   - State transitions should be smoother
   - No integration hierarchy errors
   - Clear State Manager guidance

## Support

If you encounter issues:

1. Check `.upgrade-backup/` for recovery files
2. Review `.sf-version` for upgrade details
3. Validate state with `tools/validate-state-file.sh`
4. Check integration hierarchy compliance

The upgrade should take less than 1 minute and preserve all your work while fixing critical bugs.

## Summary

**Action Required**: Run the SF 3.0 upgrade to get critical bug fixes

**Risk Level**: LOW - Comprehensive backups, safe current state

**Time Required**: ~1 minute for upgrade, 2-3 minutes for validation

**Command**:
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
bash /home/vscode/software-factory-template/tools/upgrade-sf3.sh --verbose
```