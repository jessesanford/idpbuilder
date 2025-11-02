# 📋📋📋 RULE R612 - AGENT HISTORY MANAGEMENT (STANDARD)

**Criticality:** STANDARD - Best practice for historical tracking
**Grading Impact:** -10% to -30% for violations
**Enforcement:** Schema validation, size monitoring, query performance

---

## STANDARD STATEMENT

**THE agents_history ARRAY PROVIDES A PERMANENT HISTORICAL RECORD OF ALL COMPLETED AGENTS. METADATA MUST BE MINIMAL TO PREVENT BLOAT, STRUCTURED FOR QUERYABILITY, AND VALIDATED AGAINST SCHEMA. AGENTS_HISTORY ENABLES AUDIT TRAILS, METRICS, AND DEBUGGING WITHOUT DEGRADING PERFORMANCE.**

---

## 📋📋📋 THE agents_history ARRAY 📋📋📋

### Purpose and Design

```yaml
agents_history_purpose:
  audit_trail:
    - Complete record of all agents ever spawned
    - Who did what, when, and with what outcome
    - Historical tracking for post-project analysis

  metrics_source:
    - Aggregate statistics (total lines changed, bugs found)
    - Performance metrics (agent completion times)
    - Success/failure rates by agent type

  debugging_aid:
    - Reference past agent behavior
    - Identify patterns in failures
    - Track which efforts required rework

  compliance:
    - Demonstrate all work reviewed
    - Prove agent lifecycle followed
    - Audit orchestrator decisions
```

### Design Principles

```yaml
design_principles:
  1_minimal_metadata:
    - Keep only essential fields
    - No verbose logs or detailed progress
    - Target: ~500 bytes per agent

  2_structured_data:
    - Consistent schema for all agent types
    - Queryable with jq/JSON tools
    - Machine-readable and human-readable

  3_bounded_growth:
    - Self-limiting through minimal metadata
    - Archive old projects if needed
    - Target: <100KB for typical project

  4_immutable_record:
    - Never modify historical entries
    - Only append new entries
    - Preserve original timestamps and outcomes
```

---

## 🔴🔴🔴 PART 1: SCHEMA SPECIFICATION 🔴🔴🔴

### Core Schema Fields

**REQUIRED for all agent types:**

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "agent_id",
    "agent_type",
    "final_state",
    "completed_at",
    "work_summary"
  ],
  "properties": {
    "agent_id": {
      "type": "string",
      "description": "Unique agent identifier",
      "pattern": "^(swe|reviewer|architect|integration)-.*$"
    },
    "agent_type": {
      "type": "string",
      "enum": ["sw-engineer", "code-reviewer", "architect", "integration"],
      "description": "Type of agent"
    },
    "final_state": {
      "type": "string",
      "description": "Terminal state when agent completed",
      "examples": ["COMPLETE", "COMPLETED", "FAILED"]
    },
    "completed_at": {
      "type": "string",
      "format": "date-time",
      "description": "ISO 8601 timestamp of completion"
    },
    "work_summary": {
      "type": "object",
      "description": "Agent-type-specific work summary"
    },
    "metrics": {
      "type": "object",
      "description": "Quantitative metrics from agent work"
    }
  }
}
```

### Agent-Type-Specific Schemas

**SW-Engineer Agent History:**

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
    "outcome": "approved",
    "splits": 0
  },
  "metrics": {
    "lines_added": 450,
    "lines_removed": 30,
    "files_modified": 8,
    "commits": 12,
    "duration_hours": 2.5
  }
}
```

**Code-Reviewer Agent History:**

```json
{
  "agent_id": "reviewer-1.2.1-20251102-043000",
  "agent_type": "code-reviewer",
  "final_state": "COMPLETE",
  "completed_at": "2025-11-02T04:45:00Z",
  "work_summary": {
    "reviewed_effort": "1.2.1",
    "review_type": "effort_review",
    "bugs_found": 3,
    "review_outcome": "approved_with_fixes"
  },
  "metrics": {
    "review_duration_minutes": 15,
    "files_reviewed": 8,
    "critical_issues": 0,
    "blocking_issues": 1,
    "warnings": 2
  }
}
```

**Architect Agent History:**

```json
{
  "agent_id": "architect-wave2-review-20251102-050000",
  "agent_type": "architect",
  "final_state": "COMPLETE",
  "completed_at": "2025-11-02T05:15:00Z",
  "work_summary": {
    "review_scope": "wave",
    "wave_id": "1.2",
    "architecture_outcome": "approved",
    "concerns_raised": 0
  },
  "metrics": {
    "efforts_reviewed": 3,
    "review_duration_minutes": 30,
    "recommendations_made": 2
  }
}
```

**Integration Agent History:**

```json
{
  "agent_id": "integration-wave2-20251102-051500",
  "agent_type": "integration",
  "final_state": "COMPLETE",
  "completed_at": "2025-11-02T05:30:00Z",
  "work_summary": {
    "integration_type": "wave",
    "wave_id": "1.2",
    "efforts_integrated": 3,
    "outcome": "success"
  },
  "metrics": {
    "merge_conflicts": 0,
    "test_failures": 0,
    "integration_duration_minutes": 15
  }
}
```

---

## 🔴🔴🔴 PART 2: METADATA EXTRACTION 🔴🔴🔴

### Extraction Functions by Agent Type

**Extract SW-Engineer Metadata:**

```bash
extract_swe_agent_history() {
    local agent_id="$1"
    local state_file="${2:-orchestrator-state-v3.json}"

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
                outcome: .status,
                splits: (.split_count // 0)
            },
            metrics: {
                lines_added: (.line_count_tracking.final_count // 0),
                lines_removed: (.line_count_tracking.lines_removed // 0),
                files_modified: (.files_modified_count // 0),
                commits: (.commits_made // 0),
                duration_hours: (
                    if .completed_at and .spawned_at then
                        (((.completed_at | fromdate) - (.spawned_at | fromdate)) / 3600)
                    else 0 end
                )
            }
        }
    ' "$state_file"
}
```

**Extract Code-Reviewer Metadata:**

```bash
extract_reviewer_agent_history() {
    local agent_id="$1"
    local state_file="${2:-orchestrator-state-v3.json}"

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
                review_type: (.review_type // "effort_review"),
                bugs_found: (.extraction_metrics.total_bugs // 0),
                review_outcome: (.review_result // "completed")
            },
            metrics: {
                review_duration_minutes: (.duration_minutes // 0),
                files_reviewed: (.files_reviewed_count // 0),
                critical_issues: (.extraction_metrics.critical_count // 0),
                blocking_issues: (.extraction_metrics.blocking_count // 0),
                warnings: (.extraction_metrics.warning_count // 0)
            }
        }
    ' "$state_file"
}
```

**Extract Architect Metadata:**

```bash
extract_architect_agent_history() {
    local agent_id="$1"
    local state_file="${2:-orchestrator-state-v3.json}"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: "architect",
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                review_scope: (.review_scope // "wave"),
                wave_id: (.wave_id // .focus),
                architecture_outcome: (.outcome // "approved"),
                concerns_raised: (.concerns_count // 0)
            },
            metrics: {
                efforts_reviewed: (.efforts_reviewed_count // 0),
                review_duration_minutes: (.duration_minutes // 0),
                recommendations_made: (.recommendations_count // 0)
            }
        }
    ' "$state_file"
}
```

**Extract Integration Agent Metadata:**

```bash
extract_integration_agent_history() {
    local agent_id="$1"
    local state_file="${2:-orchestrator-state-v3.json}"

    jq --arg agent_id "$agent_id" '
        .active_agents[] |
        select(.agent_id == $agent_id) |
        {
            agent_id: .agent_id,
            agent_type: "integration",
            final_state: .state,
            completed_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            work_summary: {
                integration_type: (.integration_type // "wave"),
                wave_id: (.wave_id // .focus),
                efforts_integrated: (.efforts_merged_count // 0),
                outcome: (.integration_outcome // "success")
            },
            metrics: {
                merge_conflicts: (.merge_conflicts_count // 0),
                test_failures: (.test_failures_count // 0),
                integration_duration_minutes: (.duration_minutes // 0)
            }
        }
    ' "$state_file"
}
```

---

## 🔴🔴🔴 PART 3: APPENDING TO HISTORY 🔴🔴🔴

### Atomic Append Operation

```bash
append_to_agents_history() {
    local state_file="$1"
    local agent_metadata="$2"

    # Validate metadata is valid JSON
    if ! echo "$agent_metadata" | jq empty 2>/dev/null; then
        echo "❌ R612: Invalid JSON metadata"
        return 1
    fi

    # Validate against schema
    if ! validate_agent_history_entry "$agent_metadata"; then
        echo "❌ R612: Metadata doesn't match schema"
        return 1
    fi

    # Backup state file
    local backup_file="${state_file}.backup-history-$(date +%Y%m%d-%H%M%S)"
    cp "$state_file" "$backup_file"

    # Atomic append
    local temp_file=$(mktemp)
    jq --argjson metadata "$agent_metadata" '
        # Initialize agents_history if missing
        if .agents_history == null then
            .agents_history = []
        else . end |

        # Append new entry
        .agents_history += [$metadata]
    ' "$state_file" > "$temp_file"

    # Validate result
    if ! jq empty "$temp_file" 2>/dev/null; then
        echo "❌ R612: State file corrupted during append"
        rm "$temp_file"
        return 1
    fi

    # Replace original
    mv "$temp_file" "$state_file"

    echo "✅ R612: Agent history entry appended successfully"
    return 0
}
```

### Schema Validation

```bash
validate_agent_history_entry() {
    local metadata="$1"
    local schema_file="$CLAUDE_PROJECT_DIR/schemas/agents-history-schema.json"

    # Check required fields
    local required_fields=("agent_id" "agent_type" "final_state" "completed_at" "work_summary")

    for field in "${required_fields[@]}"; do
        if ! echo "$metadata" | jq -e ".$field" >/dev/null 2>&1; then
            echo "❌ R612: Missing required field: $field"
            return 1
        fi
    done

    # Validate agent_type
    local agent_type=$(echo "$metadata" | jq -r '.agent_type')
    case "$agent_type" in
        "sw-engineer"|"code-reviewer"|"architect"|"integration")
            # Valid
            ;;
        *)
            echo "❌ R612: Invalid agent_type: $agent_type"
            return 1
            ;;
    esac

    # Validate timestamp format
    local completed_at=$(echo "$metadata" | jq -r '.completed_at')
    if ! date -d "$completed_at" >/dev/null 2>&1; then
        echo "❌ R612: Invalid timestamp format: $completed_at"
        return 1
    fi

    echo "✅ R612: Metadata validation passed"
    return 0
}
```

---

## 🔴🔴🔴 PART 4: QUERYING HISTORY 🔴🔴🔴

### Common Query Patterns

**Count Agents by Type:**

```bash
count_agents_by_type() {
    local state_file="${1:-orchestrator-state-v3.json}"

    echo "📊 Agent History Summary:"
    jq -r '
        .agents_history |
        group_by(.agent_type) |
        map({
            agent_type: .[0].agent_type,
            count: length
        }) |
        .[] |
        "  \(.agent_type): \(.count) agents"
    ' "$state_file"
}
```

**Total Lines Changed:**

```bash
total_lines_changed() {
    local state_file="${1:-orchestrator-state-v3.json}"

    jq '
        [.agents_history[] |
         select(.agent_type == "sw-engineer") |
         .metrics.lines_added // 0] |
        add // 0
    ' "$state_file"
}
```

**Average Review Duration:**

```bash
average_review_duration() {
    local state_file="${1:-orchestrator-state-v3.json}"

    jq '
        [.agents_history[] |
         select(.agent_type == "code-reviewer") |
         .metrics.review_duration_minutes // 0] |
        (add / length)
    ' "$state_file"
}
```

**Total Bugs Found:**

```bash
total_bugs_found() {
    local state_file="${1:-orchestrator-state-v3.json}"

    jq '
        [.agents_history[] |
         select(.agent_type == "code-reviewer") |
         .work_summary.bugs_found // 0] |
        add // 0
    ' "$state_file"
}
```

**List Failed Agents:**

```bash
list_failed_agents() {
    local state_file="${1:-orchestrator-state-v3.json}"

    jq -r '
        .agents_history[] |
        select(.final_state == "FAILED") |
        "\(.agent_id) (\(.agent_type)) - \(.completed_at)"
    ' "$state_file"
}
```

---

## 🔴🔴🔴 PART 5: SIZE MANAGEMENT 🔴🔴🔴

### Size Targets and Monitoring

```yaml
agents_history_size_targets:
  per_entry:
    target: 400-600 bytes
    acceptable: 600-1000 bytes
    warning: 1000-2000 bytes
    critical: ">2000 bytes (too verbose)"

  total_array:
    small_project: "<25KB (50 agents)"
    medium_project: "25-75KB (50-150 agents)"
    large_project: "75-150KB (150-300 agents)"
    critical: ">150KB (needs archival)"
```

### Size Calculation

```bash
calculate_agents_history_size() {
    local state_file="${1:-orchestrator-state-v3.json}"

    # Extract agents_history array
    local history_json=$(jq '.agents_history' "$state_file")
    local size_bytes=$(echo "$history_json" | wc -c)
    local size_kb=$(( size_bytes / 1024 ))

    # Count entries
    local entry_count=$(echo "$history_json" | jq 'length')

    # Calculate average entry size
    local avg_entry_size=0
    if [ "$entry_count" -gt 0 ]; then
        avg_entry_size=$(( size_bytes / entry_count ))
    fi

    echo "📊 R612: agents_history Size Report"
    echo "  Total entries: $entry_count"
    echo "  Total size: ${size_kb}KB ($size_bytes bytes)"
    echo "  Average entry: ${avg_entry_size} bytes"

    # Evaluate
    if [ "$size_kb" -gt 150 ]; then
        echo "  ❌ CRITICAL: History too large (${size_kb}KB > 150KB)"
        echo "  Recommendation: Archive old entries"
        return 1
    elif [ "$size_kb" -gt 75 ]; then
        echo "  ⚠️  WARNING: History growing large (${size_kb}KB > 75KB)"
        return 1
    elif [ "$avg_entry_size" -gt 1000 ]; then
        echo "  ⚠️  WARNING: Average entry too large (${avg_entry_size} > 1000 bytes)"
        echo "  Recommendation: Reduce metadata verbosity"
        return 1
    else
        echo "  ✅ Size within targets"
        return 0
    fi
}
```

### Archival Strategy

```bash
archive_old_agents_history() {
    local state_file="${1:-orchestrator-state-v3.json}"
    local keep_count="${2:-100}"  # Keep last 100 agents

    # Count current entries
    local total_count=$(jq '.agents_history | length' "$state_file")

    if [ "$total_count" -le "$keep_count" ]; then
        echo "✅ R612: No archival needed ($total_count <= $keep_count)"
        return 0
    fi

    # Archive old entries
    local archive_file="agents-history-archive-$(date +%Y%m%d-%H%M%S).json"
    local archive_count=$(( total_count - keep_count ))

    echo "📦 R612: Archiving oldest $archive_count entries to $archive_file..."

    # Extract old entries
    jq --arg count "$archive_count" '
        .agents_history[:($count | tonumber)]
    ' "$state_file" > "$archive_file"

    # Keep only recent entries
    local temp_file=$(mktemp)
    jq --arg count "$keep_count" '
        .agents_history = .agents_history[-($count | tonumber):]
    ' "$state_file" > "$temp_file"

    # Validate
    if jq empty "$temp_file" 2>/dev/null; then
        mv "$temp_file" "$state_file"
        echo "✅ R612: Archived $archive_count entries, keeping $keep_count most recent"
        return 0
    else
        echo "❌ R612: Archival failed - state file corrupted"
        rm "$temp_file"
        return 1
    fi
}
```

---

## 🔴 GRADING IMPACT

### Compliance Grading

```yaml
r612_compliance:
  schema_validation:
    all_entries_valid: 30%
    required_fields_present: 20%

  size_management:
    entries_minimal: 20%          # <1KB per entry
    total_size_reasonable: 15%    # <150KB total

  data_integrity:
    no_duplicates: 10%
    timestamps_valid: 5%

total: 100%

violations:
  invalid_schema: -20%              # Per invalid entry
  oversized_entries: -15%           # Entries >2KB
  oversized_history: -30%           # Total >150KB
  missing_required_fields: -25%     # Per entry missing fields
  duplicate_agent_ids: -10%         # Per duplicate
```

---

## 📊 INTEGRATION WITH OTHER RULES

### R612 + R610 (Lifecycle Protocol)

```yaml
integration:
  - R610 extracts metadata when cleaning up
  - R612 validates and stores metadata
  - R610 removes from active_agents
  - R612 preserves in agents_history
```

### R612 + R611 (Active Agents Cleanup)

```yaml
integration:
  - R611 identifies inactive agents
  - R612 archives their metadata
  - Both ensure complete tracking
```

### R612 + R613 (State File Growth)

```yaml
integration:
  - R612 manages agents_history size
  - R613 monitors total state file size
  - Both trigger archival if needed
```

---

## ✅ SUCCESS CRITERIA

```yaml
r612_success:
  - All completed agents in agents_history
  - All entries validate against schema
  - Average entry size <1KB
  - Total history size <150KB
  - No duplicate agent_ids
  - All timestamps valid ISO 8601
  - Queryable with standard jq commands

r612_failure:
  - Missing agent history entries
  - Invalid schema entries
  - Oversized entries (>2KB)
  - Oversized history (>150KB)
  - Duplicate agent_ids
  - Invalid timestamps
  - Unqueryable data
```

---

**Remember:** R612 ensures agents_history is a clean, minimal, queryable historical record. Every completed agent gets exactly one entry with just the essential metadata. This enables metrics, auditing, and debugging without bloating the state file.

**See Also:**
- R610: Agent Metadata Lifecycle Protocol (cleanup timing)
- R611: Active Agents Cleanup Protocol (active vs completed)
- R613: State File Growth Management (overall size)
- schemas/agents-history-schema.json (JSON schema)
- tools/cleanup-completed-agents.sh (uses R612 for archival)
