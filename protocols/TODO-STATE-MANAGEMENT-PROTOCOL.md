# TODO STATE MANAGEMENT PROTOCOL

## ⚠️ PREREQUISITE: Requires .claude/settings.json
**This protocol DEPENDS on the PreCompact hooks in `.claude/settings.json`**
- Without settings.json, compaction recovery WILL NOT WORK
- See `CRITICAL-SETTINGS-JSON.md` for configuration details
- The PreCompact hooks automatically preserve TODO state during memory compaction

## CRITICAL: Prevent Lost TODOs During Mode Transitions

### Purpose
Ensure no tasks are lost when agents switch modes or experience context loss.

## 🚨 MANDATORY RULES FOR ALL AGENTS 🚨

### 1. BEFORE Mode Transition - SAVE TODOs

When switching from one mode to another in the state machine:

```bash
# Step 1: Create timestamped TODO file
TODO_FILE="/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/todos/${AGENT_NAME}-${CURRENT_MODE}-$(date '+%Y%m%d-%H%M%S').todo"

# Step 2: Write current TODOs to file
echo "# TODOs saved from ${AGENT_NAME} in ${CURRENT_MODE} mode" > $TODO_FILE
echo "# Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')" >> $TODO_FILE
echo "# Next Mode: ${NEXT_MODE}" >> $TODO_FILE
echo "" >> $TODO_FILE
# [Write your actual TODOs here in a structured format]

# Step 3: Verify file was created
if [ -f "$TODO_FILE" ]; then
    echo "✅ TODOs saved to: $TODO_FILE"
else
    echo "❌ ERROR: Failed to save TODOs!"
    # STOP - Do not proceed without saving TODOs
fi
```

### 2. AFTER Mode Transition - LOAD TODOs

When entering a new mode:

```bash
# Step 1: Check for existing TODO files
TODO_DIR="/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/todos"
MY_TODOS=$(ls -t $TODO_DIR/${AGENT_NAME}-*.todo 2>/dev/null | head -5)

# Step 2: If TODO files exist, load them
if [ ! -z "$MY_TODOS" ]; then
    echo "📋 Found TODO files to process:"
    for todo_file in $MY_TODOS; do
        echo "  - Loading: $todo_file"
        # READ the file and merge with current TODOs
        # De-duplicate against current task list
    done
fi

# Step 3: CRITICAL - Delete processed TODO files
for todo_file in $MY_TODOS; do
    rm -f $todo_file
    echo "  ✓ Deleted processed file: $todo_file"
done
```

## 📊 TODO File Format

### Standard TODO File Structure
```yaml
# TODOs saved from orchestrator in WAVE_COMPLETE mode
# Timestamp: 2025-08-21 16:30:00 UTC
# Next Mode: INTEGRATION_REVIEW

pending:
  - Create phase1/wave1-integration branch
  - Merge all Wave 1 splits into integration branch
  - Spawn architect for Wave 1 review
  
in_progress:
  - Fix Wave 2 webhook framework issue
  
completed:
  - E1.1.1 implementation
  - E1.1.2 implementation
  
blocked:
  - Wave 4 start (waiting on integration)
  
context:
  current_phase: 1
  current_wave: 3
  blocking_issue: "CHANGES_REQUIRED from architect"
```

## 🎯 Agent-Specific Mode Transitions

### ORCHESTRATOR (@agent-orchestrator-prompt-engineer-task-master)

File pattern: `orchestrator-{MODE}-{TIMESTAMP}.todo`

Key transitions:
- `WAVE_COMPLETE` → `INTEGRATION_REVIEW`: MUST save integration TODOs
- `INTEGRATION_REVIEW` → `ARCHITECT_REVIEW`: MUST save fixes needed
- `ARCHITECT_REVIEW` → `WAVE_START`: MUST save next wave tasks
- `CHANGES_REQUIRED` → `FIX_IMPLEMENTATION`: MUST save fix list

### SOFTWARE ENGINEER (@agent-kcp-go-lang-sr-sw-eng)

File pattern: `sw-eng-{MODE}-{TIMESTAMP}.todo`

Key transitions:
- `IMPLEMENTATION` → `MEASURE_SIZE`: MUST save progress
- `MEASURE_SIZE` → `SPLIT_REQUIRED`: MUST save completed work
- `FIX_REVIEW_ISSUES` → `IMPLEMENTATION`: MUST save fixes applied

### CODE REVIEWER (@agent-kcp-kubernetes-code-reviewer)

File pattern: `code-reviewer-{MODE}-{TIMESTAMP}.todo`

Key transitions:
- `PLANNING` → `REVIEW`: MUST save plan created
- `REVIEW` → `SPLIT_PLANNING`: MUST save issues found
- `SPLIT_EXECUTION` → `REVIEW`: MUST save split status

### ARCHITECT (@agent-kcp-architect-reviewer)

File pattern: `architect-{MODE}-{TIMESTAMP}.todo`

Key transitions:
- `WAVE_REVIEW` → `DECISION`: MUST save findings
- `PHASE_REVIEW` → `ASSESSMENT`: MUST save recommendations

## ⚠️ CRITICAL CLEANUP RULES

### MANDATORY: Delete After Processing
```bash
# After successfully loading and merging TODOs:
rm -f $TODO_FILE
echo "✅ Cleaned up processed TODO file"
```

### FORBIDDEN: Leave Stale TODOs
```bash
# Check for old TODO files (>1 hour old)
find $TODO_DIR -name "*.todo" -mmin +60 -exec rm {} \;
echo "🧹 Cleaned up stale TODO files"
```

### Weekly Cleanup (Orchestrator Only)
```bash
# Every Monday, clean ALL TODO files
if [ $(date +%u) -eq 1 ]; then
    rm -f $TODO_DIR/*.todo
    echo "🗑️ Weekly TODO cleanup completed"
fi
```

## 🔄 Example Mode Transition

### Orchestrator: WAVE_COMPLETE → INTEGRATION_REVIEW

```bash
# BEFORE leaving WAVE_COMPLETE mode
TODO_FILE="/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/todos/orchestrator-WAVE_COMPLETE-20250821-163000.todo"

cat > $TODO_FILE << 'EOF'
# TODOs saved from orchestrator in WAVE_COMPLETE mode
# Timestamp: 2025-08-21 16:30:00 UTC
# Next Mode: INTEGRATION_REVIEW

pending:
  - Create phase1/wave3-integration branch
  - Merge effort4-deepcopy-gen splits (3 branches)
  - Merge effort5-client-gen splits (3 branches)
  - Spawn architect for integration review
  - Update orchestrator-state.yaml integration_branches

in_progress:
  - None

completed:
  - E1.2.4 split into 3 compliant branches
  - E1.2.5 split into 3 compliant branches

context:
  current_phase: 1
  current_wave: 3
  total_splits_pending_merge: 6
EOF

echo "✅ Saved TODOs before mode transition"

# AFTER entering INTEGRATION_REVIEW mode
TODO_FILES=$(ls -t /workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/todos/orchestrator-*.todo 2>/dev/null)

if [ ! -z "$TODO_FILES" ]; then
    for file in $TODO_FILES; do
        echo "Loading TODOs from: $file"
        cat $file
        # Merge with current TODO list
        # De-duplicate
        rm -f $file
        echo "✓ Processed and deleted: $file"
    done
fi
```

## 🚨 Failure Scenarios

### If TODO Save Fails
```bash
# STOP IMMEDIATELY
echo "❌ CRITICAL: Failed to save TODOs before mode transition!"
echo "Current TODOs that would be lost:"
# [Print current TODOs]
echo "STOPPING - Cannot proceed without TODO persistence"
exit 1
```

### If TODO Load Fails
```bash
# Check for backup methods
if [ -f "/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/CURRENT-TODO-STATE.md" ]; then
    echo "⚠️ TODO files not found, checking backup..."
    cat /workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/CURRENT-TODO-STATE.md
fi
```

## 📝 Integration with TodoWrite Tool

When using the TodoWrite tool, also save a backup:
```bash
# After TodoWrite updates
TODO_BACKUP="/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/todos/backup-$(date '+%Y%m%d-%H%M%S').todo"
# Export current TODOs to backup
```

## ✅ Success Criteria

1. **No Lost TODOs**: Every task tracked across mode transitions
2. **No Stale TODOs**: Old files cleaned up immediately
3. **No Duplicate Work**: De-duplication prevents re-doing tasks
4. **Clear Audit Trail**: Timestamped files show transition history
5. **Fast Recovery**: TODOs restored quickly after context loss

## 🎯 Implementation Priority

**HIGHEST PRIORITY**: Orchestrator (manages entire workflow)
**HIGH PRIORITY**: Code Reviewer (manages splits)
**MEDIUM PRIORITY**: SW Engineer (implements code)
**LOW PRIORITY**: Architect (reviews only)

This protocol ensures zero task loss during the complex state machine transitions!