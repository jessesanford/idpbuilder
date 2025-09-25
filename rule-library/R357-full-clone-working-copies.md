# ⚠️⚠️⚠️ RULE R357: Full Clone Working Copies (NO WORKTREES)

## Rule Definition
**Criticality:** WARNING
**Category:** Infrastructure
**Applies To:** orchestrator, all agents, all infrastructure setup
**Cross-References:** R271, R193, R176, R209

## ⚠️⚠️⚠️ CRITICAL CLARIFICATION: NO GIT WORKTREES ⚠️⚠️⚠️

**The Software Factory 2.0 uses FULL CLONE working copies, NOT git worktrees.**

## HISTORICAL CONTEXT

### OLD SYSTEM (Deprecated - DO NOT USE):
- Used `git worktree add` for isolation
- Pattern: `worktrees/phase{N}-integration`
- Shared .git directory between worktrees
- **STATUS: COMPLETELY DEPRECATED**

### CURRENT SYSTEM (MANDATORY):
- Uses `git clone` for EVERY working copy
- Each effort/integration gets its own FULL clone
- Pattern: `efforts/phase{N}/wave{M}/effort-name/`
- Pattern: `phase{N}/integration/`
- Each clone has its own .git directory
- **STATUS: REQUIRED FOR ALL INFRASTRUCTURE**

## CORRECT INFRASTRUCTURE CREATION

### ✅ CORRECT - Full Clone Pattern:
```bash
# For efforts:
EFFORT_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
mkdir -p "$(dirname "$EFFORT_DIR")"
git clone --single-branch --branch "$BASE_BRANCH" "$TARGET_REPO" "$EFFORT_DIR"
cd "$EFFORT_DIR"
git checkout -b "$EFFORT_BRANCH"

# For phase integration:
INTEGRATION_DIR="$CLAUDE_PROJECT_DIR/phase${PHASE}/integration"
mkdir -p "$(dirname "$INTEGRATION_DIR")"
git clone "$TARGET_REPO" "$INTEGRATION_DIR"
cd "$INTEGRATION_DIR"
git checkout -b "phase${PHASE}-integration" "$BASE_BRANCH"

# For project integration:
INTEGRATION_DIR="$CLAUDE_PROJECT_DIR/project/integration"
mkdir -p "$(dirname "$INTEGRATION_DIR")"
git clone "$TARGET_REPO" "$INTEGRATION_DIR"
cd "$INTEGRATION_DIR"
git checkout -b "project-integration" "$BASE_BRANCH"
```

### ❌ WRONG - Worktree Pattern (NEVER USE):
```bash
# NEVER DO THIS:
git worktree add -b "$BRANCH" "worktrees/phase2-integration" "$BASE"
git worktree add "worktrees/effort-name" "$BRANCH"
git worktree remove "worktrees/integration"
```

## DIRECTORY STRUCTURE

### Correct Structure (Full Clones):
```
$CLAUDE_PROJECT_DIR/
├── efforts/
│   └── phase1/
│       └── wave1/
│           ├── effort-1/        # Full clone with own .git/
│           └── effort-2/        # Full clone with own .git/
├── phase1/
│   └── integration/            # Full clone with own .git/
├── phase2/
│   └── integration/            # Full clone with own .git/
└── project/
    └── integration/            # Full clone with own .git/
```

### Wrong Structure (Worktrees - DEPRECATED):
```
# NEVER CREATE THIS STRUCTURE:
worktrees/
├── phase1-integration/         # ❌ Worktree
├── phase2-integration/         # ❌ Worktree
└── project-integration/        # ❌ Worktree
```

## WHY FULL CLONES?

1. **Complete Isolation**: Each working copy is completely independent
2. **No Shared State**: No risk of corruption from shared .git directory
3. **Parallel Safety**: Multiple agents can work without interference
4. **Clean Recovery**: Easy to delete and recreate if issues occur
5. **Standard Git**: Works with all git operations without special handling

## ENFORCEMENT

### State Machine Validation:
- SETUP_EFFORT_INFRASTRUCTURE: Must use git clone
- SETUP_PHASE_INTEGRATION_INFRASTRUCTURE: Must use git clone
- SETUP_PROJECT_INTEGRATION_INFRASTRUCTURE: Must use git clone
- CREATE_NEXT_SPLIT_INFRASTRUCTURE: Must use git clone

### Orchestrator Responsibilities:
1. Create full clones for ALL infrastructure
2. Never use `git worktree` commands
3. Ensure each clone has correct branch and remote
4. Verify isolation before spawning agents

## MIGRATION NOTES

If you find any references to:
- `git worktree add`
- `git worktree remove`
- `worktrees/` directories
- "worktree" in documentation

**These are OUTDATED and must be updated to use full clones.**

## Related Rules
- R271: Single-Branch Full Checkout Protocol (defines clone strategy)
- R193: Effort Clone Protocol (superseded, historical reference)
- R176: Workspace Isolation (requires separate working copies)
- R209: Effort Directory Isolation Protocol (directory structure)

## Grading Impact
- Using worktrees instead of full clones: -20% penalty
- Shared .git directory corruption: -50% penalty
- Infrastructure creation failures: -30% penalty