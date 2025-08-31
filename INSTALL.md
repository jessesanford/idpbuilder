# Software Factory 2.0 Installation Guide

## Quick Installation

### Option 1: Run from Template Directory
```bash
cd /workspaces/software-factory-2.0-template
./setup.sh                    # Create new project
./migrate-planning-only.sh    # Migrate planning from SF 1.0
./migrate-from-1.0.sh        # Full migration from SF 1.0
```

### Option 2: Add to PATH
```bash
# Add to your ~/.bashrc or ~/.zshrc:
export PATH="/workspaces/software-factory-2.0-template:$PATH"

# Then from anywhere:
sf2 new
sf2 migrate-planning
sf2 help
```

### Option 3: Create Alias
```bash
# Add to your ~/.bashrc or ~/.zshrc:
alias sf2="/workspaces/software-factory-2.0-template/sf2"

# Then from anywhere:
sf2 new
sf2 migrate-planning
```

## Known Issues Fixed

### Issue: Quotes in Project Description
**Problem:** Project descriptions containing single quotes (') or double quotes (") would break the setup script.
**Status:** ✅ FIXED - The script now properly escapes quotes in all contexts.

### Issue: Path Dependencies
**Problem:** Scripts had hardcoded paths to `/workspaces/software-factory-2.0-template/`
**Status:** ✅ FIXED - Scripts now use relative paths based on their actual location.

## Troubleshooting

### "setup.sh not found" Error
If you get this error, make sure you're running the scripts from the template directory or using the full path:
```bash
# Wrong:
./setup.sh                    # From wrong directory

# Right:
/workspaces/software-factory-2.0-template/setup.sh

# Or:
cd /workspaces/software-factory-2.0-template
./setup.sh
```

### Permission Denied
Make sure scripts are executable:
```bash
chmod +x /workspaces/software-factory-2.0-template/*.sh
chmod +x /workspaces/software-factory-2.0-template/sf2
```

### Scripts Not Finding Template Files
The scripts look for template files relative to their own location. Ensure you haven't moved scripts without the template files:
```bash
# Template structure required:
software-factory-2.0-template/
├── setup.sh
├── migrate-*.sh
├── 🚨-CRITICAL/
├── agent-states/
├── state-machines/
├── utilities/
├── expertise/
├── quick-reference/
└── .claude/
```

## Testing Your Installation

Test that everything works:
```bash
# Test 1: Check help
/workspaces/software-factory-2.0-template/sf2 help

# Test 2: Create test project (interactive)
cd /tmp
/workspaces/software-factory-2.0-template/setup.sh
# Answer: test-project, Test description, /tmp/test-project, etc.

# Test 3: Clean up
rm -rf /tmp/test-project
```

## For IDPBuilder Migration

To migrate your IDPBuilder project:
```bash
# From anywhere:
/workspaces/software-factory-2.0-template/migrate-planning-only.sh

# When prompted:
# Source: /home/vscode/workspaces/idpbuilder-1.0-version
# Project name: idpbuilder-oci-mgmt
# Answer the setup wizard questions

# The script will:
# 1. Run setup.sh to create SF 2.0 structure
# 2. Extract planning from your SF 1.0 project
# 3. Create MASTER-IMPLEMENTATION-PLAN.md
# 4. Set up for SF 2.0 regeneration of implementation
```

## Support

If you encounter issues:
1. Check this INSTALL.md for known issues
2. Verify template files are in place
3. Ensure scripts have execute permissions
4. Check that quotes in inputs are properly handled