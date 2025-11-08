# 📊📊📊 RULE R613 - STATE FILE GROWTH MANAGEMENT (STANDARD)

**Criticality:** STANDARD - Performance and maintainability
**Grading Impact:** -15% to -40% for violations
**Enforcement:** Size monitoring, automated alerts, cleanup triggers

---

## STANDARD STATEMENT

**THE orchestrator-state-v3.json FILE MUST REMAIN MANAGEABLE IN SIZE (<500KB TARGET, <1MB HARD LIMIT). UNBOUNDED GROWTH DEGRADES PERFORMANCE, SLOWS GIT OPERATIONS, AND INDICATES MISSING CLEANUP. MONITORING AND AUTOMATED CLEANUP PREVENT STATE FILE BLOAT.**

---

## 📊📊📊 STATE FILE GROWTH FUNDAMENTALS 📊📊📊

### Why Size Matters

```yaml
performance_impact:
  small_file:
    size: "<100KB"
    jq_query_time: "10-50ms"
    git_operations: "fast"
    context_loading: "instant"

  medium_file:
    size: "100-500KB"
    jq_query_time: "50-200ms"
    git_operations: "acceptable"
    context_loading: "noticeable"

  large_file:
    size: "500KB-1MB"
    jq_query_time: "200-500ms"
    git_operations: "slow"
    context_loading: "delayed"

  critical_file:
    size: ">1MB"
    jq_query_time: ">500ms (unacceptable)"
    git_operations: "very slow"
    context_loading: "problematic"
    action_required: "IMMEDIATE CLEANUP"
```

### Growth Sources

```yaml
state_file_components:
  active_agents:
    typical_size: "5-25KB (1-5 agents)"
    bloat_risk: HIGH
    cleanup_mechanism: R610 (automatic)

  agents_history:
    typical_size: "25-75KB (50-150 agents)"
    bloat_risk: MEDIUM
    cleanup_mechanism: R612 (archival)

  phases:
    typical_size: "50-150KB (3-5 phases)"
    bloat_risk: LOW
    cleanup_mechanism: "None (project data)"

  waves:
    typical_size: "100-300KB (10-20 waves)"
    bloat_risk: LOW
    cleanup_mechanism: "None (project data)"

  efforts:
    typical_size: "150-400KB (30-80 efforts)"
    bloat_risk: LOW
    cleanup_mechanism: "None (project data)"

  metadata:
    typical_size: "10-50KB"
    bloat_risk: MEDIUM
    cleanup_mechanism: "Prune verbose logs"
```

---

## 🔴🔴🔴 PART 1: SIZE TARGETS AND THRESHOLDS 🔴🔴🔴

### Size Targets by Component

```yaml
component_size_targets:
  active_agents:
    target: "<25KB"
    acceptable: "25-50KB"
    warning: "50-100KB"
    critical: ">100KB - R610 cleanup required"

  agents_history:
    target: "<50KB"
    acceptable: "50-100KB"
    warning: "100-150KB"
    critical: ">150KB - R612 archival required"

  total_state_file:
    target: "<250KB (optimal)"
    acceptable: "250-500KB (good)"
    warning: "500KB-1MB (needs attention)"
    critical: ">1MB (IMMEDIATE ACTION REQUIRED)"
```

### Threshold Actions

```yaml
threshold_actions:
  at_warning_threshold:
    - Log warning in state transition
    - Monitor growth rate
    - Prepare cleanup plan
    - No blocking action

  at_critical_threshold:
    - BLOCKING action required
    - Run R610 cleanup
    - Run R612 archival if needed
    - Cannot proceed until resolved

  ongoing_monitoring:
    - Check size at each COMPLETE_* state
    - Alert if growth rate excessive
    - Automatic cleanup if thresholds exceeded
```

---

## 🔴🔴🔴 PART 2: SIZE MONITORING 🔴🔴🔴

### Monitoring Functions

**Measure Total Size:**

```bash
measure_state_file_size() {
    local state_file="${1:-orchestrator-state-v3.json}"

    if [ ! -f "$state_file" ]; then
        echo "❌ R613: State file not found: $state_file"
        return 1
    fi

    # Get file size
    local size_bytes=$(wc -c < "$state_file")
    local size_kb=$(( size_bytes / 1024 ))
    local size_mb=$(awk "BEGIN {printf \"%.2f\", $size_kb/1024}")

    echo "📊 R613: State File Size Report"
    echo "  File: $state_file"
    echo "  Size: ${size_kb}KB (${size_bytes} bytes, ${size_mb}MB)"

    # Evaluate against thresholds
    if [ "$size_bytes" -gt 1048576 ]; then
        # >1MB - CRITICAL
        echo "  ❌ CRITICAL: State file exceeds 1MB limit"
        echo "  ACTION REQUIRED: Immediate cleanup necessary"
        return 2
    elif [ "$size_bytes" -gt 512000 ]; then
        # >500KB - WARNING
        echo "  ⚠️  WARNING: State file approaching limit (>500KB)"
        echo "  RECOMMENDED: Review and cleanup"
        return 1
    elif [ "$size_bytes" -gt 256000 ]; then
        # >250KB - ACCEPTABLE
        echo "  ✅ Size acceptable (>250KB but <500KB)"
        return 0
    else
        # <250KB - TARGET
        echo "  ✅ Size optimal (<250KB)"
        return 0
    fi
}
```

**Measure Component Sizes:**

```bash
measure_component_sizes() {
    local state_file="${1:-orchestrator-state-v3.json}"

    echo "📊 R613: Component Size Breakdown"

    # active_agents
    local active_agents_size=$(jq '.active_agents' "$state_file" | wc -c)
    local active_agents_kb=$(( active_agents_size / 1024 ))
    echo "  active_agents: ${active_agents_kb}KB ($active_agents_size bytes)"

    # agents_history
    local history_size=$(jq '.agents_history // []' "$state_file" | wc -c)
    local history_kb=$(( history_size / 1024 ))
    echo "  agents_history: ${history_kb}KB ($history_size bytes)"

    # phases
    local phases_size=$(jq '.phases // []' "$state_file" | wc -c)
    local phases_kb=$(( phases_size / 1024 ))
    echo "  phases: ${phases_kb}KB ($phases_size bytes)"

    # waves (all phases)
    local waves_size=$(jq '[.phases[]?.waves[]?] // []' "$state_file" | wc -c)
    local waves_kb=$(( waves_size / 1024 ))
    echo "  waves: ${waves_kb}KB ($waves_size bytes)"

    # efforts (all waves)
    local efforts_size=$(jq '[.phases[]?.waves[]?.efforts[]?] // []' "$state_file" | wc -c)
    local efforts_kb=$(( efforts_size / 1024 ))
    echo "  efforts: ${efforts_kb}KB ($efforts_size bytes)"

    # metadata
    local metadata_size=$(jq 'del(.active_agents, .agents_history, .phases)' "$state_file" | wc -c)
    local metadata_kb=$(( metadata_size / 1024 ))
    echo "  metadata: ${metadata_kb}KB ($metadata_size bytes)"

    # Identify bloat
    echo ""
    echo "📊 R613: Bloat Analysis"

    if [ "$active_agents_kb" -gt 100 ]; then
        echo "  ❌ active_agents bloated (${active_agents_kb}KB > 100KB) - Run R610 cleanup"
    fi

    if [ "$history_kb" -gt 150 ]; then
        echo "  ⚠️  agents_history large (${history_kb}KB > 150KB) - Consider R612 archival"
    fi

    if [ "$metadata_kb" -gt 50 ]; then
        echo "  ⚠️  metadata large (${metadata_kb}KB > 50KB) - Prune verbose logs"
    fi
}
```

**Calculate Growth Rate:**

```bash
calculate_growth_rate() {
    local state_file="${1:-orchestrator-state-v3.json}"

    # Get current size
    local current_size=$(wc -c < "$state_file")

    # Check for previous measurement
    local history_file=".state-file-size-history"
    if [ ! -f "$history_file" ]; then
        # First measurement - create history
        echo "$(date +%s) $current_size" > "$history_file"
        echo "📊 R613: First size measurement recorded"
        return 0
    fi

    # Get previous measurement
    local previous_line=$(tail -1 "$history_file")
    local previous_time=$(echo "$previous_line" | awk '{print $1}')
    local previous_size=$(echo "$previous_line" | awk '{print $2}')

    # Calculate time delta (seconds)
    local current_time=$(date +%s)
    local time_delta=$(( current_time - previous_time ))

    # Calculate size delta (bytes)
    local size_delta=$(( current_size - previous_size ))

    # Calculate growth rate (bytes per hour)
    local growth_rate=0
    if [ "$time_delta" -gt 0 ]; then
        growth_rate=$(( (size_delta * 3600) / time_delta ))
    fi

    # Record current measurement
    echo "$current_time $current_size" >> "$history_file"

    echo "📊 R613: Growth Rate Analysis"
    echo "  Previous: ${previous_size} bytes"
    echo "  Current: ${current_size} bytes"
    echo "  Delta: ${size_delta} bytes over ${time_delta}s"
    echo "  Rate: ${growth_rate} bytes/hour"

    # Evaluate growth rate
    local growth_kb_per_hour=$(( growth_rate / 1024 ))
    if [ "$growth_kb_per_hour" -gt 100 ]; then
        echo "  ❌ WARNING: Excessive growth rate (${growth_kb_per_hour}KB/hour)"
        return 1
    else
        echo "  ✅ Growth rate acceptable (${growth_kb_per_hour}KB/hour)"
        return 0
    fi
}
```

---

## 🔴🔴🔴 PART 3: CLEANUP STRATEGIES 🔴🔴🔴

### Automatic Cleanup Triggers

```bash
trigger_cleanup_if_needed() {
    local state_file="${1:-orchestrator-state-v3.json}"

    echo "🔍 R613: Checking if cleanup needed..."

    # Measure current size
    local size_bytes=$(wc -c < "$state_file")
    local size_kb=$(( size_bytes / 1024 ))

    # Check thresholds
    local cleanup_needed=false
    local cleanup_reason=""

    if [ "$size_bytes" -gt 1048576 ]; then
        cleanup_needed=true
        cleanup_reason="CRITICAL: >1MB"
    elif [ "$size_bytes" -gt 512000 ]; then
        cleanup_needed=true
        cleanup_reason="WARNING: >500KB"
    fi

    if [ "$cleanup_needed" = true ]; then
        echo "🧹 R613: Cleanup triggered - $cleanup_reason"

        # Run R610 cleanup (completed agents)
        echo "Running R610 cleanup (completed agents)..."
        bash tools/cleanup-completed-agents.sh

        # Check if R612 archival needed
        local history_kb=$(jq '.agents_history // [] | length' "$state_file")
        if [ "$history_kb" -gt 150 ]; then
            echo "Running R612 archival (agents_history)..."
            bash tools/archive-agents-history.sh 100
        fi

        # Measure after cleanup
        local new_size=$(wc -c < "$state_file")
        local new_kb=$(( new_size / 1024 ))
        local saved_kb=$(( size_kb - new_kb ))

        echo "✅ R613: Cleanup complete"
        echo "  Before: ${size_kb}KB"
        echo "  After: ${new_kb}KB"
        echo "  Saved: ${saved_kb}KB"

        return 0
    else
        echo "✅ R613: No cleanup needed (${size_kb}KB within limits)"
        return 0
    fi
}
```

### Manual Cleanup Options

```bash
# Full cleanup (R610 + R612)
cleanup_state_file_full() {
    echo "🧹 R613: Full state file cleanup..."

    # R610: Cleanup completed agents
    bash tools/cleanup-completed-agents.sh

    # R612: Archive old history if needed
    bash tools/archive-agents-history.sh 100

    # Measure result
    measure_state_file_size
}

# Prune verbose metadata
prune_verbose_metadata() {
    local state_file="${1:-orchestrator-state-v3.json}"

    echo "🧹 R613: Pruning verbose metadata..."

    # Backup
    cp "$state_file" "${state_file}.backup-prune-$(date +%Y%m%d-%H%M%S)"

    # Remove verbose fields
    local temp_file=$(mktemp)
    jq '
        # Remove verbose logs from active_agents
        .active_agents[] |= (
            del(.detailed_progress_log) |
            del(.verbose_error_stack) |
            del(.intermediate_checkpoints)
        ) |

        # Remove verbose metadata
        del(.debug_info) |
        del(.temporary_metadata)
    ' "$state_file" > "$temp_file"

    # Validate and replace
    if jq empty "$temp_file" 2>/dev/null; then
        mv "$temp_file" "$state_file"
        echo "✅ R613: Verbose metadata pruned"
    else
        echo "❌ R613: Pruning failed"
        rm "$temp_file"
        return 1
    fi
}
```

---

## 🔴🔴🔴 PART 4: MONITORING INTEGRATION 🔴🔴🔴

### Monitoring at Boundary States

**COMPLETE_WAVE:**

```bash
# In COMPLETE_WAVE state
echo "📊 R613: Monitoring state file size at wave completion..."

measure_state_file_size

# Trigger cleanup if needed
trigger_cleanup_if_needed

echo "✅ R613: State file size within limits"
```

**COMPLETE_PHASE:**

```bash
# In COMPLETE_PHASE state
echo "📊 R613: Monitoring state file size at phase completion..."

measure_state_file_size
measure_component_sizes
calculate_growth_rate

# Cleanup if critical
if measure_state_file_size | grep -q "CRITICAL"; then
    echo "🚨 R613: CRITICAL size - forcing cleanup..."
    cleanup_state_file_full
fi

echo "✅ R613: State file health validated"
```

**COMPLETE_PROJECT:**

```bash
# In COMPLETE_PROJECT state (PROJECT_DONE)
echo "📊 R613: Final state file health check..."

measure_state_file_size
measure_component_sizes

# Generate final report
echo ""
echo "📊 R613: Project State File Report"
jq -r '
    "Total Phases: \(.phases | length)",
    "Total Waves: \([.phases[]?.waves[]?] | length)",
    "Total Efforts: \([.phases[]?.waves[]?.efforts[]?] | length)",
    "Total Agents (history): \(.agents_history // [] | length)",
    "Active Agents: \(.active_agents | length)"
' orchestrator-state-v3.json

echo "✅ R613: Project complete - state file final size recorded"
```

---

## 🔴 GRADING IMPACT

### Size Compliance Grading

```yaml
r613_compliance:
  size_targets:
    optimal_size: 40%           # <250KB
    acceptable_size: 30%        # 250-500KB
    warning_size: 10%           # 500KB-1MB
    critical_size: 0%           # >1MB

  monitoring:
    regular_monitoring: 25%     # Check at boundary states
    growth_rate_tracking: 15%   # Track growth over time
    component_analysis: 10%     # Identify bloat sources

  cleanup:
    automatic_cleanup: 10%      # Trigger cleanup when needed

total: 100%

violations:
  exceeds_1mb: -40%                     # CRITICAL violation
  exceeds_500kb_no_cleanup: -25%        # WARNING, no action
  no_monitoring: -20%                   # No size checks
  excessive_growth_rate: -15%           # >100KB/hour
  bloated_active_agents: -30%           # >100KB active_agents (R610 violation)
```

### Performance Impact

```yaml
performance_degradation:
  size_under_250kb:
    jq_queries: "fast (<50ms)"
    git_commits: "instant"
    state_loading: "negligible"

  size_250_500kb:
    jq_queries: "acceptable (50-150ms)"
    git_commits: "fast"
    state_loading: "noticeable"

  size_500kb_1mb:
    jq_queries: "slow (150-300ms)"
    git_commits: "slow"
    state_loading: "delayed"

  size_over_1mb:
    jq_queries: "very slow (>300ms)"
    git_commits: "very slow"
    state_loading: "problematic"
    action: "IMMEDIATE CLEANUP REQUIRED"
```

---

## 📊 INTEGRATION WITH OTHER RULES

### R613 + R610 (Agent Lifecycle)

```yaml
integration:
  - R613 monitors total state file size
  - R610 cleans up completed agents (primary bloat source)
  - R613 triggers R610 if active_agents bloated
  - Both prevent unbounded growth
```

### R613 + R612 (Agent History)

```yaml
integration:
  - R613 monitors agents_history size
  - R612 manages agents_history archival
  - R613 triggers R612 archival if history >150KB
  - Both keep historical data manageable
```

### R613 + R288 (State File Updates)

```yaml
integration:
  - R288 ensures atomic state file updates
  - R613 monitors result of updates
  - Both ensure state file integrity
  - R613 validates size after R288 operations
```

---

## 📋 MONITORING CHECKLIST

### At Each Boundary State

```yaml
monitoring_checklist:
  COMPLETE_WAVE:
    - ✅ Measure total state file size
    - ✅ Check active_agents size
    - ✅ Trigger cleanup if >500KB

  COMPLETE_PHASE:
    - ✅ Measure total state file size
    - ✅ Measure component sizes
    - ✅ Calculate growth rate
    - ✅ Trigger cleanup if critical

  COMPLETE_PROJECT:
    - ✅ Final state file size report
    - ✅ Component breakdown
    - ✅ Growth rate analysis
    - ✅ Archive if needed
```

---

## ✅ SUCCESS CRITERIA

```yaml
r613_success:
  - State file remains <500KB throughout project
  - Regular size monitoring at boundary states
  - Automatic cleanup triggered when needed
  - Growth rate tracked and reasonable
  - Component sizes within targets
  - Performance maintained (<200ms queries)

r613_failure:
  - State file exceeds 1MB
  - No size monitoring
  - Cleanup not triggered when needed
  - Excessive growth rate (>100KB/hour)
  - Bloated components
  - Performance degradation
```

---

**Remember:** R613 is the health monitor for the state file. It ensures the orchestrator-state-v3.json remains manageable through regular monitoring and automated cleanup. Combined with R610 and R612, it prevents state file bloat and maintains system performance.

**See Also:**
- R610: Agent Metadata Lifecycle Protocol (primary cleanup)
- R611: Active Agents Cleanup Protocol (active_agents maintenance)
- R612: Agent History Management (agents_history archival)
- R288: State File Update Requirements (atomic updates)
- tools/cleanup-completed-agents.sh (R610 cleanup)
- tools/measure-state-file-size.sh (size monitoring)
