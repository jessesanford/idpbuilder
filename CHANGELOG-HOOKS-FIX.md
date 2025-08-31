# Changelog: Hooks to Utilities Fix

## Problem
The setup.sh script was failing with:
```
chmod: cannot access '/home/vscode/workspaces/idpbuilder-oci-mgmt/hooks/*.sh': No such file or directory
```

## Root Cause
We renamed `hooks/` directory to `utilities/` to accurately reflect that these are manual scripts, not automatic hooks. However, setup.sh still had references to the old `hooks/` directory.

## Changes Made

### 1. Fixed setup.sh (Line 795-804)
**Before:**
```bash
# Make hooks executable and update paths
chmod +x "$TARGET_DIR/hooks"/*.sh

# Update paths in hook scripts to point to new project
echo -e "${CYAN}Updating hook script paths...${NC}"
for hook in "$TARGET_DIR/hooks"/*.sh; do
    if [ -f "$hook" ]; then
        sed -i "s|/workspaces/software-factory-2.0-template|$TARGET_DIR|g" "$hook"
    fi
done
echo -e "${GREEN}✓ Hook scripts configured${NC}"
```

**After:**
```bash
# Utilities are already made executable above in the utilities section
# No need to update paths as utilities use relative paths
```

### 2. Verified Correct Sections
- Lines 450-466: Correctly handles `utilities/` directory
- Lines 425-429: Correctly describes PreCompact hook only
- `.claude/settings.json`: Only contains real PreCompact hook

## Testing
```bash
# The setup should now work without errors:
cd /your/project
/path/to/software-factory-2.0-template/setup.sh

# Should see:
✓ Utility scripts installed and made executable
✓ Support directories created (todos/, checkpoints/, snapshots/)
Note: These are manual utility scripts, not automatic hooks!
```

## Key Points
1. **utilities/** directory contains manual scripts (not automatic hooks)
2. **PreCompact** is the only real Claude Code hook (in `.claude/settings.json`)
3. Scripts in utilities/ must be run manually
4. Setup.sh now correctly handles the utilities directory

## Verification
After setup, you should have:
```
your-project/
├── .claude/
│   └── settings.json       # Only PreCompact hook
├── utilities/              # Manual scripts
│   ├── pre-compact.sh
│   ├── post-compact.sh
│   ├── todo-preservation.sh
│   ├── state-snapshot.sh
│   └── recovery-assistant.sh
├── todos/                  # Created by setup
├── checkpoints/           # Created by setup
└── snapshots/             # Created by setup
```

The error is now fixed and the setup will complete successfully.