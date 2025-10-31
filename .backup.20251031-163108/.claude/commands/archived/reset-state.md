---
name: reset-state
description: Emergency reset of Software Factory state machine when corrupted
---

# /reset-state

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║                      STATE RESET COMMAND                                      ║
║                                                                               ║
║ Rules: EMERGENCY-PROTOCOLS + STATE-RECOVERY + BACKUP-VERIFICATION            ║
║ + CAREFUL-RESET + CONTEXT-PRESERVATION                                        ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## ⚠️ CRITICAL WARNING ⚠️

This command performs a controlled reset of the Software Factory state machine.
**USE WITH EXTREME CAUTION** - Only when state is corrupted or unrecoverable.

## 🚨 PRE-RESET SAFETY CHECKS 🚨

### MANDATORY: Verify Reset is Needed
```bash
# Before resetting, MUST verify state is actually corrupted:
echo "🔍 DIAGNOSING STATE CORRUPTION:"
echo "1. Is orchestrator-state-v3.json readable and valid?"
echo "2. Are TODO files recoverable?"
echo "3. Are current branches in valid state?"
echo "4. Is work in progress that can be salvaged?"

# Read current state if possible
if [[ -f "./agent-configs/[project]/orchestrator-state-v3.json" ]]; then
    echo "📄 CURRENT STATE FILE EXISTS - Review before resetting:"
    cat "./agent-configs/[project]/orchestrator-state-v3.json"
    echo ""
    echo "❓ Are you sure state is corrupted and cannot be repaired?"
    echo "❓ Consider using context recovery instead of full reset"
fi
```

### MANDATORY: Backup Current State
```bash
# Create backup before any reset operations
BACKUP_DIR="./agent-configs/[project]/backups/reset-$(date '+%Y%m%d-%H%M%S')"
mkdir -p "$BACKUP_DIR"

echo "💾 CREATING STATE BACKUP TO: $BACKUP_DIR"

# Backup critical files
cp -r "./agent-configs/[project]/todos/" "$BACKUP_DIR/todos-backup/" 2>/dev/null
cp "./agent-configs/[project]/orchestrator-state-v3.json" "$BACKUP_DIR/" 2>/dev/null
cp "./agent-configs/[project]/"*-STATE*.md "$BACKUP_DIR/" 2>/dev/null
cp "./agent-configs/[project]/CURRENT-"*.md "$BACKUP_DIR/" 2>/dev/null

# Backup git state
git branch -a > "$BACKUP_DIR/git-branches.txt"
git status > "$BACKUP_DIR/git-status.txt"
git log --oneline -20 > "$BACKUP_DIR/git-recent-commits.txt"

echo "✅ BACKUP COMPLETED"
```

## 🔄 RESET TYPE SELECTION

Choose the appropriate reset level based on the severity of state corruption:

### LEVEL 1: TODO-ONLY RESET (Safest)
Use when only TODO state is corrupted, but main state machine is intact.

```bash
echo "🔄 PERFORMING LEVEL 1 RESET: TODO-ONLY"

# Clear corrupt TODO files
if [[ -d "./agent-configs/[project]/todos/" ]]; then
    echo "🗑️ Clearing corrupt TODO files..."
    mv "./agent-configs/[project]/todos/" "./agent-configs/[project]/todos-corrupted-$(date '+%Y%m%d-%H%M%S')/"
    mkdir -p "./agent-configs/[project]/todos/"
fi

# Reset TODO state for each agent type
echo "📋 Initializing fresh TODO state for all agents..."

# The actual TODOs will be populated by each agent when they start
echo "✅ LEVEL 1 RESET COMPLETE - TODOs cleared"
echo "ℹ️  Agents will initialize fresh TODO lists on startup"
```

### LEVEL 2: STATE-MACHINE RESET (Moderate)
Use when state machine is corrupted but work branches are intact.

```bash
echo "🔄 PERFORMING LEVEL 2 RESET: STATE-MACHINE"

# Level 1 reset first
# [Include Level 1 steps above]

# Reset state machine files
echo "🎯 Resetting state machine to INIT..."

# Reset orchestrator state
cat > "./agent-configs/[project]/orchestrator-state-v3.json" << 'EOF'
# Software Factory 2.0 Orchestrator State
# Reset to INIT state on $(date)

project: "[project]"
state_machine_version: "2.0"
current_state: "INIT"
current_phase: 1
current_wave: 0
last_update: "$(date '+%Y-%m-%d %H:%M:%S')"

phases_planned: 0
waves_planned: 0
efforts_planned: 0

efforts_in_progress: {}
efforts_completed: {}
integration_branches: {}

last_architect_review: null
last_line_count_check: null

reset_history:
  - timestamp: "$(date '+%Y-%m-%d %H:%M:%S')"
    type: "STATE-MACHINE RESET"
    reason: "State corruption recovery"
    backup_location: "$BACKUP_DIR"
EOF

echo "✅ LEVEL 2 RESET COMPLETE - State machine reset to INIT"
echo "ℹ️  Orchestrator will need to re-plan phases and waves"
```

### LEVEL 3: FULL RESET (Dangerous)
Use only when everything is corrupted and must start completely fresh.

```bash
echo "🔄 PERFORMING LEVEL 3 RESET: FULL RESET"
echo "⚠️  THIS WILL RESET ALL PROGRESS - LAST CHANCE TO ABORT"
echo "   Press Ctrl+C within 10 seconds to abort..."
sleep 10

# Level 1 & 2 resets first
# [Include Level 1 & 2 steps above]

# Reset all configuration and planning files
echo "🗑️ Resetting all configuration files..."

# Clear current state files
rm -f "./agent-configs/[project]/CURRENT-"*.md
rm -f "./agent-configs/[project]/"*-CURRENT-*.md

# Reset work directories (preserve backups)
echo "🏗️ Resetting work directories..."
if [[ -d "./efforts/" ]]; then
    mv "./efforts/" "./efforts-reset-backup-$(date '+%Y%m%d-%H%M%S')/"
    mkdir -p "./efforts/"
fi

# Reset git to clean state (preserve work in backup branches)
echo "🌿 Creating backup branches for current work..."
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "main" && "$CURRENT_BRANCH" != "master" ]]; then
    git checkout -b "backup-before-reset-$(date '+%Y%m%d-%H%M%S')"
    git checkout main 2>/dev/null || git checkout master 2>/dev/null
fi

echo "✅ LEVEL 3 RESET COMPLETE - Full factory reset"
echo "ℹ️  All progress reset - start with fresh orchestration"
echo "ℹ️  Backup branches and files preserved for recovery"
```

## 📋 POST-RESET INITIALIZATION

### Agent State Verification
```bash
echo "🔍 POST-RESET VERIFICATION:"

# Verify file structure
echo "📁 Checking directory structure..."
ls -la "./agent-configs/[project]/" | head -10

# Verify state file
echo "📄 Checking state file..."
if [[ -f "./agent-configs/[project]/orchestrator-state-v3.json" ]]; then
    echo "✅ State file exists"
    head -10 "./agent-configs/[project]/orchestrator-state-v3.json"
else
    echo "❌ State file missing - may need manual creation"
fi

# Verify git state
echo "🌿 Checking git state..."
echo "Current branch: $(git branch --show-current)"
echo "Status: $(git status --porcelain | wc -l) modified files"
```

### Agent Restart Requirements
```bash
echo "🔄 AGENT RESTART REQUIREMENTS:"
echo ""
echo "After reset, ALL agents must:"
echo "1. ✅ Perform complete startup sequence"
echo "2. ✅ Read their core configuration files"
echo "3. ✅ Initialize fresh TODO lists"
echo "4. ✅ Verify environment and identity"
echo "5. ✅ Load appropriate expertise modules"
echo ""
echo "Orchestrator must:"
echo "1. ✅ Read reset state from orchestrator-state-v3.json"
echo "2. ✅ Begin fresh planning cycle"
echo "3. ✅ Re-establish phase and wave structure"
echo "4. ✅ Update grading checkpoints"
```

## 🔄 RECOVERY OPTIONS

### Partial Recovery from Backup
```bash
# If you need to recover specific items after reset:
RECOVERY_FROM="$BACKUP_DIR"

echo "🔄 PARTIAL RECOVERY OPTIONS:"
echo "1. Recover TODO state: cp $RECOVERY_FROM/todos-backup/* ./agent-configs/[project]/todos/"
echo "2. Recover state file: cp $RECOVERY_FROM/orchestrator-state-v3.json ./agent-configs/[project]/"
echo "3. Recover current files: cp $RECOVERY_FROM/CURRENT-*.md ./agent-configs/[project]/"
echo "4. Check git branches: cat $RECOVERY_FROM/git-branches.txt"
echo "5. Review recent commits: cat $RECOVERY_FROM/git-recent-commits.txt"
```

### Gradual Re-initialization
```bash
echo "📋 GRADUAL RE-INITIALIZATION STEPS:"
echo ""
echo "Instead of full reset, consider:"
echo "1. Start with orchestrator context recovery"
echo "2. Load TODO state from most recent backup"
echo "3. Verify work branch integrity"
echo "4. Resume from last known good state"
echo "5. Re-validate current progress"
echo ""
echo "Use: /check-status to assess current state first"
```

## 🚨 SAFETY PROTOCOLS

### Validation Before Proceeding
```bash
# Before allowing any agent to proceed after reset:
echo "🔍 MANDATORY POST-RESET VALIDATION:"
echo ""
echo "✅ Verify agent identity and configuration"
echo "✅ Confirm state machine is in valid state"
echo "✅ Check all required files are present"
echo "✅ Validate git repository integrity"
echo "✅ Confirm backup was successful"
echo "✅ Test basic agent functionality"
```

### Reset Documentation
```bash
# Document the reset for future reference
RESET_LOG="./agent-configs/[project]/RESET-HISTORY.md"
cat >> "$RESET_LOG" << EOF

# Reset Event: $(date '+%Y-%m-%d %H:%M:%S')

## Reset Level: [1/2/3]
## Reason: [State corruption/Recovery/Testing]
## Backup Location: $BACKUP_DIR

## Pre-Reset State:
- Current phase: [if known]
- Current wave: [if known]
- Efforts in progress: [count]
- TODOs preserved: [yes/no]

## Post-Reset Actions Required:
- [ ] Orchestrator re-initialization
- [ ] Agent identity verification
- [ ] TODO state recreation
- [ ] Progress assessment
- [ ] Planning restart

## Recovery Notes:
[Any specific recovery steps or considerations]

---
EOF
```

## ⚠️ FINAL WARNINGS

### When NOT to Use Reset
```bash
echo "❌ DO NOT USE RESET IF:"
echo "- State can be recovered through normal context recovery"
echo "- Only minor TODO synchronization issues exist"
echo "- Work is in progress that hasn't been backed up"
echo "- Git branches contain unmerged critical work"
echo "- Reset won't actually solve the underlying problem"
```

### Alternative Solutions
```bash
echo "🔧 CONSIDER THESE ALTERNATIVES FIRST:"
echo "1. Use /check-status to diagnose specific issues"
echo "2. Use context recovery protocols from CLAUDE.md"
echo "3. Manually fix specific state file corruption"
echo "4. Recover TODOs from most recent backup"
echo "5. Use git operations to fix branch issues"
echo "6. Consult specific agent continuation commands"
```

This reset command provides controlled state recovery while preserving as much work as possible and ensuring proper safety protocols are followed.