# 🚨🚨🚨 RULE R264: Work Log Tracking Requirements 🚨🚨🚨

## Rule Definition
**Criticality:** BLOCKING
**Category:** Agent-Specific
**Applies To:** integration-agent

## METICULOUS TRACKING REQUIREMENTS

### 1. File Naming Requirements (R383 MANDATORY)

**🔴🔴🔴 SUPREME LAW: All work logs MUST have timestamps per R383**

```bash
# MANDATORY: Use sf_metadata_path helper function
source $CLAUDE_PROJECT_DIR/utilities/sf-metadata-path.sh

# Create work log with timestamp
WORK_LOG=$(sf_metadata_path "$PHASE" "$WAVE" "$EFFORT" "work-log" "log")
echo "# Integration Work Log" > "$WORK_LOG"

# ❌ FORBIDDEN: work-log.md (no timestamp)
# ✅ REQUIRED: work-log--20250121-143052.log
```

### 2. Replayability Standard
The work log MUST be **100% REPLAYABLE**:
- Anyone can execute commands in sequence
- Produces EXACT same integration result
- No missing steps or assumed knowledge
- All decisions documented

### 3. Required Elements for Each Operation

```markdown
## Operation [number]: [timestamp]
Purpose: [what this operation achieves]
Command: [exact command with all flags]
Working Directory: [pwd output]
Result: [command output or summary]
Conflicts: [if any]
Resolution: [how conflicts were resolved]
Verification: [how result was verified]
```

### 4. Command Categories to Track

#### MANDATORY - Always Document
- ALL git checkout commands
- ALL git merge commands
- ALL git branch operations
- ALL conflict resolutions
- ALL git commits
- ALL test/build commands
- ALL verification steps

#### Context Commands
- git status (before/after major operations)
- git log (to show state)
- git diff (to show changes)
- ls/tree (to show file structure)

### 5. Replay Script Generation
The work log MUST support replay script extraction:

```bash
# Extract all commands from work-log for replay
grep "^Command:" work-log.md | cut -d: -f2- > replay.sh

# The replay script should work
bash replay.sh  # Should reproduce the integration
```

## Work Log Template

```markdown
# Integration Work Log
Agent: integration-agent
Start Time: YYYY-MM-DD HH:MM:SS
Integration ID: [unique identifier]

## Initial State
Command: pwd
Working Directory: /workspaces/project
Result: /workspaces/project

Command: git status
Result: On branch main, clean working tree

Command: git log --oneline -5
Result: [last 5 commits]

## Operation 1: Create Integration Branch [HH:MM:SS]
Purpose: Create clean branch for integration work
Command: git checkout -b phase1-wave2-integration-20240315 main
Working Directory: /workspaces/project  
Result: Switched to a new branch 'phase1-wave2-integration-20240315'
Verification: git branch --show-current
> phase1-wave2-integration-20240315

## Operation 2: Merge First Feature [HH:MM:SS]
Purpose: Integrate feature-auth branch
Command: git merge origin/feature-auth --no-ff -m "integrate: feature-auth into wave2"
Working Directory: /workspaces/project
Result: 
  Merge made by 'recursive' strategy.
  5 files changed, 234 insertions(+), 12 deletions(-)
Conflicts: NONE
Verification: git log --oneline -1
> abc1234 integrate: feature-auth into wave2

## Operation 3: Merge Second Feature [HH:MM:SS]
Purpose: Integrate feature-api branch
Command: git merge origin/feature-api --no-ff -m "integrate: feature-api into wave2"
Working Directory: /workspaces/project
Result: 
  Auto-merging src/api/handler.go
  CONFLICT (content): Merge conflict in src/api/handler.go
Conflicts: YES - src/api/handler.go
Resolution:
  Command: vim src/api/handler.go
  Action: Kept feature-api version (newer implementation)
  Command: git add src/api/handler.go
  Command: git status
  Result: All conflicts fixed but you are still merging
  Command: git commit -m "resolve: conflict in handler.go - use feature-api version"
  Result: [commit sha]
Verification: grep -c "<<<<<<" src/api/handler.go
> 0 (no conflict markers remain)

## Final State [HH:MM:SS]
Command: git log --graph --oneline -10
Result: [integration graph showing all merges]

Command: git diff --stat main..HEAD
Result: [summary of all changes]

## Replay Verification
All commands extracted: 23 total
Replay tested: YES
Replay successful: YES
```

## Enforcement

```bash
# Verify work log replayability
test_replay() {
    local work_log="$1"
    
    # Extract commands
    grep "^Command:" "$work_log" | cut -d: -f2- > /tmp/replay.sh
    
    # Check if replay script is valid bash
    if ! bash -n /tmp/replay.sh; then
        echo "❌ Work log commands contain syntax errors"
        return 1
    fi
    
    # Count operations
    local op_count=$(grep -c "^## Operation" "$work_log")
    local cmd_count=$(grep -c "^Command:" "$work_log")
    
    echo "✅ Work log has $op_count operations with $cmd_count commands"
    
    # Verify critical elements
    grep -q "Working Directory:" "$work_log" || echo "⚠️ Missing working directory tracking"
    grep -q "Verification:" "$work_log" || echo "⚠️ Missing verification steps"
}
```

## Grading Impact
- 15% - Work log replayability
- 10% - Command completeness
- 10% - Decision documentation
- 15% - Verification steps included

## Related Rules
- R263 - Integration Documentation Requirements
- R265 - Integration Testing Requirements
- R267 - Integration Agent Grading Criteria