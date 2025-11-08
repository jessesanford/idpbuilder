# ⚠️⚠️⚠️ RULE R611 - ACTIVE AGENTS CLEANUP PROTOCOL (WARNING)

**Criticality:** WARNING - Performance degradation if violated
**Grading Impact:** -20% to -40% for violations
**Enforcement:** Query performance monitoring, state file size checks

---

## WARNING STATEMENT

**THE active_agents ARRAY MUST CONTAIN ONLY TRULY ACTIVE AGENTS. COMPLETED, ARCHIVED, OR STALE AGENTS MUST BE REMOVED. LEAVING INACTIVE AGENTS IN active_agents DEGRADES PERFORMANCE AND VIOLATES THE INTENT OF THE STATE FILE.**

---

## ⚠️⚠️⚠️ DEFINITION OF "ACTIVE" ⚠️⚠️⚠️

### Active Agent Criteria

**An agent is ACTIVE if ALL of the following are true:**

```yaml
active_agent_definition:
  state_check:
    - state != "COMPLETE"
    - state != "COMPLETED"
    - state != "ARCHIVED"
    - state != "FAILED" (if terminal failure)

  temporal_check:
    - spawned_at within reasonable timeframe
    - last_heartbeat_at recent (if heartbeat tracking)
    - not stuck in same state for >24 hours

  work_check:
    - assigned work not yet complete
    - not waiting indefinitely
    - making forward progress
```

### Inactive Agent Indicators

**An agent is INACTIVE if ANY of the following are true:**

```yaml
inactive_agent_indicators:
  terminal_states:
    - state == "COMPLETE"
    - state == "COMPLETED"
    - state == "ARCHIVED"
    - state == "FAILED"

  stuck_indicators:
    - same state for >24 hours with no progress
    - workspace deleted or inaccessible
    - branch merged and deleted
    - effort marked as "abandoned"

  stale_indicators:
    - spawned_at > 7 days ago
    - last_heartbeat_at > 2 hours ago (if heartbeats used)
    - effort_id not in current wave/phase
```

---

## 🔴🔴🔴 PART 1: CLEANUP RESPONSIBILITIES 🔴🔴🔴

### States That Perform Cleanup

**PRIMARY CLEANUP STATES** (automatic, per R610):

```yaml
MONITORING_SWE_PROGRESS:
  responsibility: Clean up completed SW-Engineer agents
  frequency: Every monitoring cycle
  trigger: Agent state == "COMPLETE"
  action: Move to agents_history within 60 seconds

MONITORING_EFFORT_REVIEWS:
  responsibility: Clean up completed Code-Reviewer agents
  frequency: Every monitoring cycle
  trigger: Agent state == "COMPLETE"
  action: Move to agents_history within 60 seconds

MONITORING_EFFORT_FIXES:
  responsibility: Clean up completed fix agents
  frequency: Every monitoring cycle
  trigger: Agent state == "COMPLETE"
  action: Move to agents_history within 60 seconds
```

**VALIDATION STATES** (safety net):

```yaml
COMPLETE_WAVE:
  responsibility: Validate no stale agents before wave completion
  frequency: Once per wave
  trigger: Wave completion checkpoint
  action: Run cleanup --validate, cleanup if needed

COMPLETE_PHASE:
  responsibility: Validate no stale agents before phase completion
  frequency: Once per phase
  trigger: Phase completion checkpoint
  action: Run cleanup --validate, cleanup if needed

COMPLETE_PROJECT:
  responsibility: Final validation before project done
  frequency: Once per project
  trigger: Project completion checkpoint
  action: Ensure active_agents empty or justify remaining
```

### Non-Cleanup States

**States that DO NOT perform cleanup:**

```yaml
states_without_cleanup:
  - SPAWN_* states (just created agents)
  - WAITING_FOR_* states (passive waiting)
  - CREATE_* states (infrastructure creation)
  - ANALYZE_* states (analysis only)

reasoning:
  - Too early: Agent just spawned
  - No visibility: Not monitoring agent progress
  - Wrong phase: Not appropriate time for cleanup
  - Responsibility: Other states handle cleanup
```

---

## 🔴🔴🔴 PART 2: PERFORMANCE REQUIREMENTS 🔴🔴🔴

### Query Performance Targets

**active_agents Array Query Performance:**

```yaml
performance_requirements:
  query_time:
    - Small array (1-5 agents): <10ms
    - Medium array (5-20 agents): <50ms
    - Large array (20-50 agents): <100ms
    - Oversized array (>50 agents): VIOLATION

  array_size_targets:
    - Typical: 1-5 agents (optimal)
    - Acceptable: 5-20 agents (monitoring multiple parallel)
    - Warning: 20-50 agents (wave with many efforts)
    - Critical: >50 agents (CLEANUP REQUIRED)

  state_file_size:
    - active_agents target: <25KB
    - active_agents acceptable: <50KB
    - active_agents warning: 50-100KB
    - active_agents critical: >100KB
```

### Performance Measurement

```bash
# Measure active_agents query performance
measure_active_agents_performance() {
    local state_file="orchestrator-state-v3.json"

    # Count agents
    local agent_count=$(jq '.active_agents | length' "$state_file")

    # Measure query time
    local start_time=$(date +%s%N)
    jq '.active_agents[] | select(.state == "COMPLETE")' "$state_file" > /dev/null
    local end_time=$(date +%s%N)

    local query_time_ms=$(( (end_time - start_time) / 1000000 ))

    echo "📊 Active Agents Performance:"
    echo "  Agent count: $agent_count"
    echo "  Query time: ${query_time_ms}ms"

    # Evaluate performance
    if [ "$agent_count" -gt 50 ]; then
        echo "  ❌ CRITICAL: Too many agents ($agent_count > 50)"
        return 1
    elif [ "$query_time_ms" -gt 100 ]; then
        echo "  ⚠️  WARNING: Slow query (${query_time_ms}ms > 100ms)"
        return 1
    else
        echo "  ✅ Performance within targets"
        return 0
    fi
}
```

### Size Calculation

```bash
# Calculate active_agents array size
calculate_active_agents_size() {
    local state_file="orchestrator-state-v3.json"

    # Extract active_agents array
    local size_bytes=$(jq '.active_agents' "$state_file" | wc -c)
    local size_kb=$(( size_bytes / 1024 ))

    echo "📊 Active Agents Size: ${size_kb}KB"

    if [ "$size_kb" -gt 100 ]; then
        echo "❌ CRITICAL: active_agents too large (${size_kb}KB > 100KB)"
        return 1
    elif [ "$size_kb" -gt 50 ]; then
        echo "⚠️  WARNING: active_agents large (${size_kb}KB > 50KB)"
        return 1
    else
        echo "✅ Size within target (<25KB)"
        return 0
    fi
}
```

---

## 🔴🔴🔴 PART 3: CLEANUP PROTOCOLS 🔴🔴🔴

### Automatic Cleanup (Primary Method)

**Performed by monitoring states per R610:**

```bash
# In MONITORING_SWE_PROGRESS state
automatic_cleanup_in_monitoring() {
    echo "🧹 R611: Checking for agents to cleanup..."

    # Find completed agents
    local completed=$(jq -r '
        .active_agents[] |
        select(.state == "COMPLETE" or .state == "COMPLETED") |
        .agent_id
    ' orchestrator-state-v3.json)

    if [ -z "$completed" ]; then
        echo "✅ R611: No completed agents found"
        return 0
    fi

    # Count
    local count=$(echo "$completed" | wc -l)
    echo "Found $count completed agents - initiating cleanup..."

    # Execute cleanup (R610)
    bash tools/cleanup-completed-agents.sh

    # Verify
    local remaining=$(jq '
        [.active_agents[] | select(.state == "COMPLETE" or .state == "COMPLETED")] | length
    ' orchestrator-state-v3.json)

    if [ "$remaining" -eq 0 ]; then
        echo "✅ R611: All completed agents cleaned up successfully"
        return 0
    else
        echo "❌ R611 VIOLATION: $remaining agents still present after cleanup"
        return 1
    fi
}
```

### Validation Cleanup (Safety Net)

**Performed by boundary states:**

```bash
# In COMPLETE_WAVE state
validation_cleanup_at_boundary() {
    echo "🔍 R611: Validating no stale agents before wave completion..."

    # Run validation mode
    if bash tools/cleanup-completed-agents.sh --validate; then
        echo "✅ R611: No stale agents found"
        return 0
    fi

    # Stale agents detected - cleanup
    echo "⚠️  R611: Stale agents detected - running cleanup..."
    bash tools/cleanup-completed-agents.sh

    # Re-validate
    if bash tools/cleanup-completed-agents.sh --validate; then
        echo "✅ R611: Cleanup successful"
        return 0
    else
        echo "❌ R611 VIOLATION: Cleanup failed"
        return 1
    fi
}
```

### Manual Cleanup

```bash
# Manual cleanup utility
bash tools/cleanup-completed-agents.sh

# Validation only (no cleanup)
bash tools/cleanup-completed-agents.sh --validate

# Cleanup specific agent
bash tools/cleanup-completed-agents.sh --agent-id "swe-1.2.1-docker-client"

# Dry run (show what would be cleaned)
bash tools/cleanup-completed-agents.sh --dry-run
```

---

## 🔴🔴🔴 PART 4: ACTIVE_AGENTS MAINTENANCE 🔴🔴🔴

### Array Structure Standards

**REQUIRED structure for each active agent:**

```json
{
  "agent_id": "swe-1.2.1-docker-client",
  "agent_type": "sw-engineer",
  "state": "IMPLEMENTATION",
  "spawned_at": "2025-11-02T04:00:00Z",
  "effort_id": "1.2.1",
  "effort_name": "docker-client",
  "workspace": "/efforts/idpbuilder-oci-mgmt/phase1/wave2/effort-1.2.1-docker-client",
  "branch_name": "idpbuilder-oci-mgmt/phase1/wave2/effort-1.2.1-docker-client",
  "status": "in_progress"
}
```

**MINIMAL structure (remove verbose fields):**

```yaml
remove_before_adding_to_active_agents:
  - Verbose error logs (keep summary only)
  - Full file listings (keep counts)
  - Intermediate checkpoints (keep latest)
  - Temporary metadata (keep essentials)
  - Debug information (remove entirely)
```

### Pruning Inactive Fields

```bash
# Prune verbose metadata before adding to active_agents
prune_agent_metadata() {
    local agent_data="$1"

    echo "$agent_data" | jq '{
        agent_id,
        agent_type,
        state,
        spawned_at,
        effort_id,
        effort_name,
        workspace,
        branch_name,
        status,
        last_update: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),

        # Keep essential metrics only
        line_count_tracking: {
            current_count: .line_count_tracking.current_count,
            limit: .line_count_tracking.limit
        },

        # Remove verbose fields
        # NO: full_progress_log
        # NO: detailed_error_stack
        # NO: intermediate_checkpoints
    }'
}
```

---

## 🔴 GRADING IMPACT

### Performance Grading

```yaml
r611_performance_grading:
  array_size:
    within_target: 30%         # <25KB
    acceptable: 20%            # 25-50KB
    warning: 10%               # 50-100KB
    critical: 0%               # >100KB

  query_performance:
    optimal: 25%               # <50ms
    acceptable: 15%            # 50-100ms
    slow: 5%                   # 100-200ms
    critical: 0%               # >200ms

  cleanup_compliance:
    automatic_cleanup: 25%     # Monitoring states cleanup
    validation_cleanup: 15%    # Boundary states validate
    no_stale_agents: 10%       # active_agents clean

total: 100%

violations:
  stale_agents_in_active: -20%        # Per stale agent
  oversized_array: -30%               # >100KB active_agents
  slow_queries: -15%                  # >200ms query time
  no_cleanup_in_monitoring: -40%      # Missing R610 integration
  cleanup_at_boundary_failed: -25%    # Validation cleanup failed
```

### Common Violations

```yaml
violation_examples:
  1_stale_completed_agents:
    description: "Agent state=COMPLETE left in active_agents"
    penalty: -20% per agent
    fix: "Run cleanup-completed-agents.sh"

  2_oversized_active_agents:
    description: "active_agents array >100KB"
    penalty: -30%
    fix: "Cleanup completed agents, prune metadata"

  3_no_monitoring_cleanup:
    description: "MONITORING_SWE_PROGRESS state doesn't cleanup"
    penalty: -40%
    fix: "Add R610 cleanup to monitoring states"

  4_slow_performance:
    description: "active_agents queries >200ms"
    penalty: -15%
    fix: "Reduce array size through cleanup"
```

---

## 📊 INTEGRATION WITH OTHER RULES

### R611 + R610 (Lifecycle Protocol)
```yaml
relationship:
  - R610 defines WHEN cleanup happens (60s after completion)
  - R611 defines WHAT active_agents should contain (only active)
  - Both prevent state file bloat
  - Both improve performance

workflow:
  1. R610 detects completed agent
  2. R610 triggers cleanup
  3. R611 validates result (only active agents remain)
  4. R611 monitors performance
```

### R611 + R612 (Agent History)
```yaml
relationship:
  - R611 removes from active_agents
  - R612 ensures data preserved in agents_history
  - Both maintain complete agent record

workflow:
  1. R611 identifies inactive agent
  2. R612 extracts metadata
  3. R612 adds to agents_history
  4. R611 removes from active_agents
```

### R611 + R613 (State File Growth)
```yaml
relationship:
  - R611 controls active_agents size
  - R613 monitors total state file size
  - Both prevent state file bloat

monitoring:
  - R611 checks active_agents size
  - R613 checks total orchestrator-state-v3.json size
  - Both trigger cleanup if thresholds exceeded
```

---

## ✅ CORRECT BEHAVIOR EXAMPLES

### Example 1: Monitoring State Compliance

```bash
# CORRECT: Cleanup in MONITORING_SWE_PROGRESS
echo "📋 CHECKLIST[3]: Cleanup completed agents per R610/R611..."

# Check for completed agents
COMPLETED_COUNT=$(jq '[.active_agents[] | select(.state == "COMPLETE")] | length' \
    orchestrator-state-v3.json)

if [ "$COMPLETED_COUNT" -gt 0 ]; then
    echo "R611: Found $COMPLETED_COUNT completed agents"
    bash tools/cleanup-completed-agents.sh
fi

# Measure performance
bash tools/measure-active-agents-performance.sh

echo "✅ CHECKLIST[3]: Active agents array clean and performant"
```

### Example 2: Boundary State Validation

```bash
# CORRECT: Validation in COMPLETE_WAVE
echo "🔍 R611: Validating active_agents before wave completion..."

# Validate no stale agents
if ! bash tools/cleanup-completed-agents.sh --validate; then
    echo "⚠️  R611: Found stale agents - cleaning up..."
    bash tools/cleanup-completed-agents.sh
fi

# Check performance
bash tools/measure-active-agents-performance.sh

echo "✅ R611: active_agents validated and clean"
```

---

## 🎯 SUCCESS CRITERIA

```yaml
r611_success:
  - active_agents contains only truly active agents
  - No COMPLETE/COMPLETED agents in active_agents
  - Query performance <100ms
  - Array size <50KB
  - Cleanup automatic in monitoring states
  - Validation successful at boundary states
  - Performance metrics tracked

r611_failure:
  - Stale agents lingering in active_agents
  - Oversized array (>100KB)
  - Slow queries (>200ms)
  - No cleanup in monitoring states
  - Validation failures at boundaries
  - Performance degradation
```

---

**Remember:** R611 ensures that active_agents truly contains only active agents. This is essential for performance - keeping the array small and queries fast. Cleanup happens automatically in monitoring states (R610), with validation at boundaries as a safety net.

**See Also:**
- R610: Agent Metadata Lifecycle Protocol (automatic cleanup)
- R612: Agent History Management (agents_history schema)
- R613: State File Growth Management (size monitoring)
- tools/cleanup-completed-agents.sh (cleanup utility)
- tools/measure-active-agents-performance.sh (performance monitoring)
