# 🚨🚨🚨 BLOCKING RULE R375: Fix State File Management Protocol

**Status**: ACTIVE
**Priority**: BLOCKING
**Impact**: State Management, Fix Tracking
**Actors**: Orchestrator, All Agents Executing Fixes

## ABSOLUTE REQUIREMENTS

### 1. SUB-STATE MACHINE ARCHITECTURE (MANDATORY)

**Main State File** (`orchestrator-state.json`):
- Tracks overall project/effort progress
- Current phase/wave/effort status
- Main state machine transitions only
- Permanent project lifecycle
- Sub-state machine references when active
- NEVER polluted with sub-state details

**Sub-State Machine Files** (`orchestrator-[identifier]-state.json`):
- Created for EACH sub-state machine entry (fix cascade, PR-ready, initialization)
- Tracks sub-state specific progress and validation
- Temporary lifecycle (archived after completion)
- Contains detailed sub-state execution
- References parent state for return

**Sub-State Machine Tracking in Main State:**
```json
{
  "sub_state_machine": {
    "active": true,
    "type": "FIX_CASCADE",
    "state_file": "orchestrator-gitea-fix-state.json",
    "current_state": "FIX_CASCADE_BACKPORT_IN_PROGRESS",
    "return_state": "MONITORING",
    "started_at": "2025-01-21T10:00:00Z"
  }
}
```

### 2. FIX STATE FILE CREATION REQUIREMENTS

**WHEN to create fix state file:**
- ANY hotfix cascade initiation
- ANY critical fix requiring multi-branch operations
- ANY error recovery requiring state tracking
- ANY fix effort spanning multiple repositories

**WHERE to store fix state files:**
- PRIMARY: Planning repository (ALWAYS)
- NEVER in target implementation repositories
- Must be in root directory alongside main state file

**NAMING CONVENTION (STRICT):**
```
orchestrator-[fix-identifier]-state.json

Examples:
- orchestrator-gitea-api-fix-state.json
- orchestrator-auth-critical-fix-state.json
- orchestrator-pr367-backport-state.json
```

### 3. SUB-STATE MACHINE FILE STRUCTURE

```json
{
  "sub_state_type": "FIX_CASCADE",
  "fix_identifier": "gitea-api-fix",
  "current_state": "FIX_CASCADE_BACKPORT_IN_PROGRESS",
  "fix_type": "HOTFIX|CRITICAL|BACKPORT|FORWARD_PORT",
  "created_at": "2025-01-21T07:00:00Z",
  "status": "IN_PROGRESS|VALIDATING|COMPLETED|FAILED",
  "parent_state_machine": {
    "state_file": "orchestrator-state.json",
    "return_state": "MONITORING",
    "nested_level": 1
  },
  "source_branch": "main",
  "target_branches": ["release-1.0", "release-1.1"],

  "backports": {
    "release-1.0": {
      "status": "COMPLETED",
      "pr_number": 123,
      "validation": "PASSED",
      "timestamp": "2025-01-21T07:30:00Z"
    }
  },

  "forward_ports": {
    "feature-branch": {
      "status": "IN_PROGRESS",
      "pr_number": null,
      "validation": "PENDING"
    }
  },

  "validation_results": {
    "unit_tests": "PASSED",
    "integration_tests": "PASSED",
    "smoke_tests": "RUNNING",
    "manual_verification": "PENDING"
  },

  "errors_encountered": [],
  "recovery_actions": [],
  "notes": []
}
```

### 4. SEPARATION OF CONCERNS (CRITICAL)

**WHAT GOES IN MAIN STATE:**
- Current effort/wave/phase
- Overall project progress
- Agent deployment status
- High-level error flags
- State machine position

**WHAT GOES IN FIX STATE:**
- Detailed fix progress
- Branch-specific status
- Validation test results
- Fix-specific errors
- Backport/forward-port tracking
- Recovery action history

### 5. LIFECYCLE MANAGEMENT

**FIX STATE LIFECYCLE:**
1. **CREATE**: At fix initiation
2. **UPDATE**: After each operation
3. **COMMIT**: Within 30 seconds of changes
4. **VALIDATE**: Before marking complete
5. **ARCHIVE**: Upon successful completion
6. **RETAIN**: In archived-fixes/ directory

**ARCHIVAL PROCESS (MANDATORY):**
```bash
# Upon fix completion
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
FIX_NAME="gitea-api-fix"

# Create archive directory if needed
mkdir -p archived-fixes/$(date +%Y)/$(date +%m)

# Archive with timestamp
mv orchestrator-${FIX_NAME}-state.json \
   archived-fixes/$(date +%Y)/$(date +%m)/${FIX_NAME}-${TIMESTAMP}.json

# Commit archival
git add archived-fixes/
git commit -m "archive: ${FIX_NAME} state completed at ${TIMESTAMP}"
git push
```

### 6. CONCURRENT FIX SUPPORT

**Multiple fixes CAN run simultaneously:**
- Each gets separate state file
- No interference between fixes
- Independent lifecycle management
- Clear tracking and audit trail

### 7. VALIDATION REQUIREMENTS

**BEFORE creating fix state:**
- Verify main state is healthy
- Ensure fix identifier is unique
- Check for existing fix state files

**DURING fix execution:**
- Update fix state after EVERY operation
- Commit and push within 30 seconds
- Never modify main state with fix details

**AFTER fix completion:**
- Validate all targets completed
- Archive state file (DO NOT DELETE)
- Update main state if needed (high-level only)
- Create completion report

## ENFORCEMENT CHECKLIST

```markdown
✅ Fix state file created for cascade
✅ Named according to convention
✅ Stored in planning repository
✅ Updated after each operation
✅ Committed within 30 seconds
✅ Main state remains unpolluted
✅ Archived upon completion
✅ Never deleted, only archived
```

## PENALTIES FOR VIOLATION

- **Creating fix details in main state**: -30% (state pollution)
- **Not creating fix state file**: -25% (tracking failure)
- **Wrong naming convention**: -15% (consistency violation)
- **Deleting instead of archiving**: -20% (audit trail loss)
- **Not committing updates**: -20% (state loss risk)

## EXAMPLES OF CORRECT USAGE

### Example 1: Hotfix Cascade
```bash
# Create fix state
cat > orchestrator-auth-hotfix-state.json << 'EOF'
{
  "fix_identifier": "auth-hotfix",
  "fix_type": "HOTFIX",
  "created_at": "2025-01-21T08:00:00Z",
  "status": "IN_PROGRESS",
  "source_branch": "main",
  "target_branches": ["release-2.0", "release-1.9"]
}
EOF

# Commit immediately
git add orchestrator-auth-hotfix-state.json
git commit -m "fix-state: initiate auth-hotfix cascade"
git push
```

### Example 2: Archival
```bash
# Complete and archive
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
mkdir -p archived-fixes/2025/01
mv orchestrator-auth-hotfix-state.json \
   archived-fixes/2025/01/auth-hotfix-${TIMESTAMP}.json
git add archived-fixes/
git commit -m "archive: auth-hotfix completed successfully"
git push
```

## SUB-STATE MACHINE TRANSITIONS

### Entering Sub-State Machine
```bash
# Main state machine enters sub-state
enter_sub_state_machine() {
    local SUB_TYPE="$1"  # FIX_CASCADE, PR_READY, INITIALIZATION
    local STATE_FILE="$2"
    local RETURN_STATE="$3"

    # Update main state to show sub-state active
    jq --arg type "$SUB_TYPE" \
       --arg file "$STATE_FILE" \
       --arg return "$RETURN_STATE" \
       --arg state "$4" \
       '.sub_state_machine = {
          "active": true,
          "type": $type,
          "state_file": $file,
          "current_state": $state,
          "return_state": $return,
          "started_at": now
       }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}
```

### Exiting Sub-State Machine
```bash
# Sub-state machine completes and returns to main
exit_sub_state_machine() {
    local RESULT="$1"  # SUCCESS, FAILURE, ABORTED

    # Get return state before clearing
    RETURN_STATE=$(jq -r '.sub_state_machine.return_state' orchestrator-state.json)

    # Clear sub-state tracking
    jq --arg state "$RETURN_STATE" \
       --arg result "$RESULT" \
       '.current_state = $state |
        .sub_state_machine.active = false |
        .sub_state_history += [{
          "type": .sub_state_machine.type,
          "completed_at": now,
          "result": $result
        }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}
```

## INTEGRATION WITH OTHER RULES

- **R206**: State machine transitions update appropriate state file
- **R287**: TODO persistence includes fix state updates
- **R341**: Error recovery creates fix state files
- **R367**: Hotfix execution uses fix state tracking
- **R376**: Quality gates use sub-state transitions

## RATIONALE

This pattern provides:
1. **Clean separation** of project vs fix concerns
2. **Concurrent fix support** without interference
3. **Complete audit trail** via archival
4. **Protection** of main state from pollution
5. **Scalability** for complex multi-branch operations

---

**Remember**: The main state file is sacred - it tracks the project's journey. Fix state files are temporal - they track specific recovery operations. Never mix the two!