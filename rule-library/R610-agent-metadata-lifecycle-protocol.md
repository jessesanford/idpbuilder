# 🚨🚨🚨 RULE R610 - AGENT METADATA LIFECYCLE PROTOCOL (BLOCKING)

**Criticality:** BLOCKING - Agents must be cleaned up within 60 seconds of completion
**Grading Impact:** -30% to -75% for violations
**Enforcement:** Monitoring states, boundary states, automated validation

---

## BLOCKING STATEMENT

**ALL COMPLETED AGENTS MUST BE AUTOMATICALLY CLEANED UP FROM active_agents WITHIN 60 SECONDS OF COMPLETION DETECTION. COMPLETED AGENTS MOVED TO agents_history WITH MINIMAL METADATA. FAILURE TO CLEANUP = STATE FILE BLOAT AND DEGRADED PERFORMANCE.**

---

## 🚨🚨🚨 THE AGENT LIFECYCLE 🚨🚨🚨

### Agent States Through Lifecycle

```yaml
agent_lifecycle_states:
  SPAWNED:
    - Agent created and added to active_agents array
    - state: "INIT" or starting state
    - Full metadata present
    - Agent begins work

  ACTIVE:
    - Agent executing work in non-terminal states
    - state: varies (IMPLEMENTATION, CODE_REVIEW, etc.)
    - Metadata updated as work progresses
    - Agent in active_agents array

  COMPLETED:
    - Agent reached terminal state
    - state: "COMPLETE" or "COMPLETED"
    - Work finished (success or failure)
    - Ready for cleanup

  ARCHIVED:
    - Agent moved to agents_history
    - Minimal metadata retained
    - Removed from active_agents
    - Historical record preserved
```

---

## 🔴🔴🔴 PART 1: AUTOMATIC CLEANUP REQUIREMENTS 🔴🔴🔴

### Detection and Cleanup Timing

**WHEN**: Cleanup occurs automatically in MONITORING states

**MONITORING_SWE_PROGRESS:**
```yaml
detection:
  - Check all agents in active_agents
  - Find agents with state="COMPLETE"
  - Or agents with state="COMPLETED"

timing:
  - Detect completion in monitoring loop
  - Cleanup within 60 seconds of detection
  - Cannot transition to next state without cleanup

enforcement:
  - BLOCKING requirement in MONITORING states
  - Pre-transition validation in COMPLETE_* states
  - State file growth monitoring
```

**MONITORING_EFFORT_REVIEWS:**
```yaml
detection:
  - Check code-reviewer agents
  - state="COMPLETE" indicates review done
  - extraction_metrics.total_bugs extracted

cleanup_trigger:
  - All reviews in wave complete
  - Extract minimal metadata
  - Move to agents_history
  - Clear from active_agents
```

### Cleanup Execution

```bash
# Automatic cleanup in monitoring states
cleanup_completed_agents() {
    local state_file="orchestrator-state-v3.json"

    # Find completed agents
    local completed_agents=$(jq -r '
        .active_agents[] |
        select(.state == "COMPLETE" or .state == "COMPLETED") |
        .agent_id
    ' "$state_file")

    if [ -z "$completed_agents" ]; then
        echo "✅ No completed agents to clean up"
        return 0
    fi

    echo "🧹 R610: Cleaning up completed agents..."

    # Use cleanup utility
    bash "$CLAUDE_PROJECT_DIR/tools/cleanup-completed-agents.sh"

    # Verify cleanup succeeded
    local remaining=$(jq '
        [.active_agents[] |
         select(.state == "COMPLETE" or .state == "COMPLETED")] |
        length
    ' "$state_file")

    if [ "$remaining" -gt 0 ]; then
        echo "❌ R610 VIOLATION: $remaining completed agents not cleaned up"
        return 1
    fi

    echo "✅ R610: All completed agents cleaned up successfully"
    return 0
}
```

---

## 🔴🔴🔴 PART 2: AGENT HISTORY METADATA 🔴🔴🔴

### Minimal Metadata Retention

**REQUIRED FIELDS** (kept in agents_history):

```json
{
  "agent_id": "swe-1.2.1-docker-client",
  "agent_type": "sw-engineer",
  "final_state": "COMPLETE",
  "completed_at": "2025-11-02T04:30:00Z",
  "work_summary": {
    "effort_id": "1.2.1",
    "effort_name": "docker-client",
    "branch": "idpbuilder-oci-mgmt/phase1/wave2/effort-1.2.1-docker-client",
    "outcome": "approved"
  },
  "metrics": {
    "lines_added": 450,
    "files_modified": 8,
    "commits": 12
  }
}
```

**Estimated Size**: ~490 bytes per agent

**Fields REMOVED** (not in agents_history):
- Full workspace paths
- Detailed progress tracking
- Intermediate state history
- Verbose error logs
- Temporary metadata

### agents_history Array Structure

```json
{
  "agents_history": [
    {
      "agent_id": "swe-1.1.1-init-schemas",
      "agent_type": "sw-engineer",
      "final_state": "COMPLETE",
      "completed_at": "2025-11-01T10:00:00Z",
      "work_summary": {
        "effort_id": "1.1.1",
        "effort_name": "init-schemas",
        "branch": "idpbuilder-oci-mgmt/phase1/wave1/effort-1.1.1",
        "outcome": "approved"
      },
      "metrics": {
        "lines_added": 250,
        "files_modified": 5,
        "commits": 8
      }
    },
    {
      "agent_id": "reviewer-1.1.1-20251101-100500",
      "agent_type": "code-reviewer",
      "final_state": "COMPLETE",
      "completed_at": "2025-11-01T10:15:00Z",
      "work_summary": {
        "reviewed_effort": "1.1.1",
        "bugs_found": 3,
        "review_outcome": "approved_with_fixes"
      },
      "metrics": {
        "review_duration_minutes": 15,
        "files_reviewed": 5
      }
    }
  ]
}
```

**Growth Rate**: ~0.5KB per agent completed
**Typical Project**: 50-100 agents = 25-50KB total
**Maximum**: Keep <25KB (self-limiting through minimal metadata)

---

## 🔴🔴🔴 PART 3: STATE FILE INTEGRITY 🔴🔴🔴

### Atomic Cleanup Operations

**BACKUP BEFORE MODIFICATION:**

```bash
# Create backup before cleanup
backup_state_file() {
    local state_file="$1"
    local backup_file="${state_file}.backup-agent-cleanup-$(date +%Y%m%d-%H%M%S)"

    cp "$state_file" "$backup_file"
    echo "✅ Backup created: $backup_file"
}
```

**ATOMIC UPDATE PATTERN:**

```bash
# Atomic cleanup operation
atomic_cleanup() {
    local state_file="$1"
    local agent_id="$2"

    # 1. Backup
    backup_state_file "$state_file"

    # 2. Extract metadata
    local metadata=$(extract_agent_metadata "$state_file" "$agent_id")

    # 3. Update state file (atomic)
    local temp_file=$(mktemp)
    jq --arg agent_id "$agent_id" \
       --argjson metadata "$metadata" '
        # Add to history
        .agents_history += [$metadata] |
        # Remove from active
        .active_agents = [.active_agents[] | select(.agent_id != $agent_id)]
    ' "$state_file" > "$temp_file"

    # 4. Validate JSON
    if ! jq empty "$temp_file" 2>/dev/null; then
        echo "❌ JSON validation failed - restoring backup"
        rm "$temp_file"
        return 1
    fi

    # 5. Replace original
    mv "$temp_file" "$state_file"

    echo "✅ Agent $agent_id cleaned up atomically"
    return 0
}
```

### Validation Requirements

**POST-CLEANUP VALIDATION:**

```bash
validate_cleanup() {
    local state_file="$1"

    # Check JSON valid
    if ! jq empty "$state_file" 2>/dev/null; then
        echo "❌ State file corrupted after cleanup"
        return 1
    fi

    # Check no COMPLETE agents in active_agents
    local completed_count=$(jq '
        [.active_agents[] | select(.state == "COMPLETE" or .state == "COMPLETED")] | length
    ' "$state_file")

    if [ "$completed_count" -gt 0 ]; then
        echo "❌ $completed_count completed agents still in active_agents"
        return 1
    fi

    # Check agents_history array exists
    if ! jq -e '.agents_history' "$state_file" > /dev/null 2>&1; then
        echo "❌ agents_history array missing"
        return 1
    fi

    echo "✅ Cleanup validation passed"
    return 0
}
```

---

## 🔴🔴🔴 PART 4: INTEGRATION WITH R611 🔴🔴🔴

### R611: Active Agents Definition

**R610 + R611 Working Together:**

```yaml
r610_responsibility:
  - Automatic cleanup timing (60 seconds)
  - Lifecycle state management
  - Metadata extraction and archival
  - State file integrity during cleanup

r611_responsibility:
  - Definition of "active" vs "completed"
  - Performance requirements for active_agents queries
  - Which states perform cleanup
  - active_agents array maintenance

integration:
  - R610 defines WHEN and HOW to cleanup
  - R611 defines WHAT cleanup means
  - Both ensure active_agents contains only active agents
  - Both prevent state file bloat
```

---

## 🔴🔴🔴 PART 5: ENFORCEMENT MECHANISMS 🔴🔴🔴

### Automatic Enforcement Points

**1. MONITORING_SWE_PROGRESS State:**
```yaml
enforcement:
  - BLOCKING checklist item: cleanup completed agents
  - Cannot transition without cleanup
  - Must acknowledge cleanup in checklist

acknowledgment_format:
  "✅ CHECKLIST[n]: Cleaned up completed agents per R610 [count] agents archived"
```

**2. MONITORING_EFFORT_REVIEWS State:**
```yaml
enforcement:
  - BLOCKING checklist item: cleanup completed reviewers
  - Extract bug metrics before cleanup
  - Archive reviewer metadata

acknowledgment_format:
  "✅ CHECKLIST[n]: Archived completed code reviewers per R610 [count] reviewers"
```

**3. COMPLETE_WAVE/COMPLETE_PHASE States (Validation):**
```yaml
enforcement:
  - STANDARD task: validate no stale agents
  - Run cleanup if stale agents found
  - Safety net for missed cleanups

validation_command:
  "bash tools/cleanup-completed-agents.sh --validate"
```

### Manual Validation

```bash
# Check for stale completed agents
validate_no_stale_agents() {
    local state_file="orchestrator-state-v3.json"

    local stale_count=$(jq '
        [.active_agents[] |
         select(.state == "COMPLETE" or .state == "COMPLETED")] |
        length
    ' "$state_file")

    if [ "$stale_count" -gt 0 ]; then
        echo "⚠️ Found $stale_count stale completed agents"
        echo "Running cleanup..."
        bash tools/cleanup-completed-agents.sh
        return $?
    fi

    echo "✅ No stale agents found"
    return 0
}
```

---

## 🔴 GRADING IMPACT

### Compliance Grading

```yaml
r610_compliance:
  automatic_cleanup:
    within_60_seconds: 25%
    correct_metadata_extraction: 20%
    atomic_operations: 20%

  state_integrity:
    json_validation: 15%
    backup_creation: 10%
    no_data_loss: 10%

total: 100%

violations:
  no_cleanup: -50%                    # Completed agents left in active_agents
  cleanup_delayed_over_60s: -30%      # Cleanup too slow
  corrupted_state_file: -75%          # JSON corruption during cleanup
  missing_history_metadata: -20%      # Incomplete agents_history entries
  no_backup: -15%                     # No backup before cleanup
```

### Performance Impact

```yaml
state_file_bloat:
  without_r610:
    - active_agents grows unbounded
    - 100 agents × 5KB each = 500KB
    - State file queries slow (>1s)
    - Git operations bloated

  with_r610:
    - active_agents: typically 1-5 agents (5-25KB)
    - agents_history: 100 agents × 0.5KB = 50KB
    - Total: ~75KB vs 500KB (85% reduction)
    - Query performance: <100ms (10× faster)
```

---

## 📊 METADATA EXTRACTION EXAMPLES

### Software Engineer Agent

```bash
extract_swe_metadata() {
    local agent_id="$1"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: "sw-engineer",
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                effort_id: .effort_id,
                effort_name: .effort_name,
                branch: .branch_name,
                outcome: .status
            },
            metrics: {
                lines_added: .line_count_tracking.final_count,
                files_modified: (.commits_made // 0),
                commits: (.commits_made // 0)
            }
        }
    ' orchestrator-state-v3.json
}
```

### Code Reviewer Agent

```bash
extract_reviewer_metadata() {
    local agent_id="$1"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: "code-reviewer",
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                reviewed_effort: .focus,
                bugs_found: (.extraction_metrics.total_bugs // 0),
                review_outcome: .review_result
            },
            metrics: {
                review_duration_minutes: .duration_minutes,
                files_reviewed: .files_reviewed_count
            }
        }
    ' orchestrator-state-v3.json
}
```

---

## 🛡️ ERROR RECOVERY

### Cleanup Failure Scenarios

**Scenario 1: JSON Corruption**
```bash
if atomic_cleanup fails:
    1. Restore from backup
    2. Log error details
    3. Transition to ERROR_RECOVERY
    4. Notify orchestrator
    5. Do NOT proceed with state transition
```

**Scenario 2: Partial Cleanup**
```bash
if some_agents_cleaned_but_not_all:
    1. Identify which agents failed
    2. Retry cleanup for failed agents
    3. If retry fails, ERROR_RECOVERY
    4. Never leave half-cleaned state
```

**Scenario 3: Missing Metadata**
```bash
if metadata_extraction_incomplete:
    1. Use best-effort metadata
    2. Mark as "incomplete_metadata": true
    3. Log warning
    4. Proceed with cleanup (don't block)
```

---

## 📋 INTEGRATION WITH OTHER RULES

### R610 + R288 (State File Updates)
- Cleanup is a state file update operation
- Must use R288 atomic update protocol
- Commit cleanup changes immediately

### R610 + R612 (Agent History Management)
- R610 defines cleanup timing and process
- R612 defines agents_history schema and integrity
- Both ensure historical data preserved

### R610 + R613 (State File Growth)
- R610 prevents growth by removing completed agents
- R613 monitors overall state file size
- Both work together to keep state file manageable

---

## ✅ CORRECT BEHAVIOR EXAMPLES

### Example 1: Monitoring State Cleanup

```bash
# In MONITORING_SWE_PROGRESS state
echo "📋 CHECKLIST[3]: Cleaning up completed agents per R610..."

# Check for completed agents
COMPLETED_COUNT=$(jq '[.active_agents[] | select(.state == "COMPLETE")] | length' \
    orchestrator-state-v3.json)

if [ "$COMPLETED_COUNT" -gt 0 ]; then
    echo "Found $COMPLETED_COUNT completed agents - cleaning up..."
    bash tools/cleanup-completed-agents.sh
else
    echo "No completed agents to clean up"
fi

echo "✅ CHECKLIST[3]: Cleaned up completed agents per R610 [$COMPLETED_COUNT agents archived]"
```

### Example 2: Wave Completion Validation

```bash
# In COMPLETE_WAVE state
echo "📋 Validating no stale agents per R610..."

bash tools/cleanup-completed-agents.sh --validate || {
    echo "⚠️ Found stale completed agents - running cleanup..."
    bash tools/cleanup-completed-agents.sh
}

echo "✅ No stale completed agents remaining"
```

---

## 🎯 SUCCESS CRITERIA

```yaml
r610_success:
  - All COMPLETE agents cleaned within 60s
  - agents_history contains all completed agents
  - active_agents contains only truly active agents
  - State file remains valid JSON
  - Backups created before each cleanup
  - No agent metadata lost
  - Performance maintained (<100ms queries)

r610_failure:
  - Completed agents left in active_agents
  - Cleanup delayed beyond 60 seconds
  - State file corrupted during cleanup
  - Agent history incomplete or missing
  - No backups created
  - Performance degradation from bloat
```

---

**Remember:** R610 is the automatic janitor of the Software Factory. It ensures that completed agents don't linger in active_agents, bloating the state file and degrading performance. Cleanup happens automatically in monitoring states, with validation at boundary states as a safety net.

**See Also:**
- R611: Active Agents Cleanup Protocol (definition of "active")
- R612: Agent History Management (agents_history schema)
- R613: State File Growth Management (size monitoring)
- R288: State File Update Requirements (atomic updates)
- tools/cleanup-completed-agents.sh (cleanup utility)
