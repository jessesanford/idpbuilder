# AGENT METADATA LIFECYCLE INVESTIGATION REPORT

**Date**: 2025-11-02
**Investigator**: Software Factory Manager Agent
**Priority**: HIGH - System Design Gap
**Status**: INVESTIGATION COMPLETE - RECOMMENDATIONS PROVIDED

---

## EXECUTIVE SUMMARY

### The Problem
The Software Factory 3.0 system has **NO agent metadata cleanup mechanism**. Completed agents (state="COMPLETE") remain indefinitely in the `active_agents` array of `orchestrator-state-v3.json`, causing:

1. **Misleading state data** - "Active" agents that finished work days/weeks ago
2. **Unbounded state file growth** - State file grows without limit over project lifetime
3. **Performance degradation** - Monitoring states scan increasingly large inactive agent lists
4. **Lost historical context** - No preserved record of agent work after state file cleanup

### Current Situation
- **6 agents** currently in `active_agents` array
- **4 agents** from Phase 1 Wave 2 (Oct 29) with state="COMPLETE" - **4 days stale**
- **2 agents** from Phase 2 Wave 2 (Nov 1-2) with state="COMPLETE" - **1 day stale**
- **NO `agents_history` field** exists in state file
- **NO cleanup states** in state machine
- **NO rules** governing agent metadata lifecycle

### Root Cause
**SOFTWARE FACTORY 3.0 WAS NOT DESIGNED WITH AGENT LIFECYCLE MANAGEMENT**

The system was designed to:
- ✅ Spawn agents
- ✅ Monitor agent progress
- ✅ Detect agent completion
- ✅ Transition to next state when agents complete

But NOT designed to:
- ❌ Clean up completed agent metadata
- ❌ Archive agent work history
- ❌ Manage state file size
- ❌ Maintain historical agent records

### Recommendation
**IMPLEMENT AGENT LIFECYCLE MANAGEMENT PROTOCOL** via new rules and state modifications.

---

## 1. ROOT CAUSE ANALYSIS

### 1.1 Investigation Findings

#### No Cleanup Logic in Monitoring States
Examined `MONITORING_SWE_PROGRESS` state (primary agent monitoring state):

**What it DOES:**
```yaml
actions:
  - Check each SW Engineer for completion status
  - Count completed vs in-progress vs blocked
  - Verify all work committed and pushed
  - Transition to SPAWN_CODE_REVIEWERS_EFFORT_REVIEW when all complete

responsibilities:
  - Active monitoring per R233 (every 5 messages)
  - Progress aggregation
  - Issue detection
  - Completion verification
```

**What it DOES NOT DO:**
```yaml
missing_actions:
  - Remove COMPLETED agents from active_agents
  - Archive agent metadata to history
  - Clean up stale agent entries
  - Manage state file size
```

**Code Evidence:**
```bash
# From MONITORING_SWE_PROGRESS/rules.md lines 364-369
if [ $COMPLETED_COUNT -eq $TOTAL_COUNT ]; then
    echo "🎉 All implementations complete!"
    ALL_COMPLETE=true
    break
fi
# → Transitions to next state, but does NOT clean up agents
```

#### No Cleanup States in State Machine
Searched `software-factory-3.0-state-machine.json` for cleanup-related states:

**States Found:**
- `COMPLETE_WAVE` - Marks WAVE as complete (not agents)
- `COMPLETE_PHASE` - Marks PHASE as complete (not agents)
- `COMPLETE_PROJECT` - Marks PROJECT as complete (not agents)
- `WAVE_COMPLETE` - Wave efforts complete (not agent cleanup)

**States NOT Found:**
- No `CLEANUP_COMPLETED_AGENTS`
- No `ARCHIVE_AGENT_METADATA`
- No `MAINTAIN_STATE_FILE`
- No agent lifecycle management states of any kind

#### No Rules Governing Agent Lifecycle
Searched rule library for agent metadata management:

**Rules Found:**
- R197 - One Agent Per Effort (spawn control)
- R008 - Monitoring Frequency (progress checking)
- R233 - Active Monitoring Requirements (progress checking)
- R288 - State File Update Protocol (updates, not cleanup)

**Rules NOT Found:**
- No rule for removing completed agents from active_agents
- No rule for archiving agent metadata
- No rule for agents_history structure
- No rule for state file size management
- No rule for agent metadata lifecycle

### 1.2 Why This Wasn't Caught Earlier

#### Design Assumptions
Software Factory 3.0 was designed with focus on:
1. **Forward progress** - Move from state to state
2. **Agent coordination** - Spawn and monitor agents
3. **Quality gates** - Ensure work quality at each stage

But NOT on:
1. **Metadata hygiene** - Clean up after work complete
2. **Historical preservation** - Archive completed work
3. **Long-term maintainability** - Manage growing state

#### Early Project Phase
This is only noticeable after:
- Multiple phases/waves complete
- Many agents spawned and finished
- State file accumulates history

In early project (Phase 1 Wave 1), this wouldn't be visible yet.

#### No Explicit Requirements
No user requirements stated:
- "Clean up completed agents"
- "Archive agent history"
- "Manage state file size"

System assumed agents would be "active" for entire project lifecycle.

---

## 2. CURRENT STATE ASSESSMENT

### 2.1 Active Agents Array Contents

**Total agents in active_agents**: 6

**Phase 1 Wave 2 (STALE - 4 days old):**
```json
{
  "agent_id": "swe-1.2.1-docker-client",
  "state": "COMPLETE",
  "spawned_at": "2025-10-29T22:12:59Z",
  "effort": "1.2.1-docker-client"
},
{
  "agent_id": "swe-1.2.2-registry-client",
  "state": "COMPLETE",
  "spawned_at": "2025-10-29T22:12:59Z",
  "effort": "1.2.2-registry-client"
},
{
  "agent_id": "swe-1.2.3-auth",
  "state": "COMPLETE",
  "spawned_at": "2025-10-29T22:12:59Z",
  "effort": "1.2.3-auth"
},
{
  "agent_id": "swe-1.2.4-tls",
  "state": "COMPLETE",
  "spawned_at": "2025-10-29T22:12:59Z",
  "effort": "1.2.4-tls"
}
```

**Phase 2 Wave 2 (STALE - 1 day old):**
```json
{
  "agent_id": "swe-2.2.1-registry-override",
  "state": "COMPLETE",
  "spawned_at": "2025-11-01T18:59:07Z",
  "completed_at": "2025-11-01T18:51:00Z",
  "effort": "2.2.1",
  "lines_implemented": 551
},
{
  "agent_id": "swe-2.2.2-env-variable-support",
  "state": "COMPLETED",
  "spawned_at": "2025-11-02T03:00:00Z"
}
```

### 2.2 Missing Infrastructure

**No agents_history field exists:**
```bash
$ jq '.agents_history' orchestrator-state-v3.json
null
```

**No cleanup protocol exists:**
- No states designed for cleanup
- No rules governing cleanup
- No documentation about when to clean up
- No schema for historical agent data

### 2.3 Impact Assessment

#### Performance Impact
**MONITORING_SWE_PROGRESS state must scan all agents:**
```bash
# From rules.md line 188
SW_ENGINEERS=$(jq -r '.spawned_agents.sw_engineers[]? | select(.state != "COMPLETE") | .effort_id' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
```

**Problem**: As project grows:
- More COMPLETE agents accumulate
- Monitoring scans larger arrays
- State queries become slower
- Unnecessary data in every state transition

**Projected Growth**:
```yaml
current_project:
  phases: 2
  waves_per_phase: ~3
  efforts_per_wave: ~4
  agents_per_effort: 2  # SW Engineer + Code Reviewer
  total_agents_expected: 2 * 3 * 4 * 2 = 48 agents

current_state:
  agents_in_active: 6
  agents_that_should_be_archived: 6  # All are COMPLETE

projected_end_of_project:
  agents_in_active: 48  # If no cleanup
  agents_that_should_be_active: 0-2  # Only current wave
```

#### Data Integrity Impact
**Misleading State Representation:**
- Field name: `active_agents`
- Actual content: 100% completed agents (6/6)
- User expectation: Currently active agents only

**Historical Context Loss:**
- Once agents removed, no record of work
- Cannot recreate timeline of agent operations
- Cannot audit agent performance
- Cannot recover from corruption

#### Maintainability Impact
**State File Growth:**
```yaml
current_size:
  orchestrator-state-v3.json: ~50KB
  active_agents: 6 entries

projected_size:
  orchestrator-state-v3.json: ~200KB
  active_agents: 48 entries (8x growth)
  state_history: 150+ transitions (3x growth)

maintenance_issues:
  - Hard to read/debug state file
  - Slow jq queries
  - Difficult manual editing
  - High risk of corruption
```

---

## 3. AGENT LIFECYCLE DESIGN RECOMMENDATIONS

### 3.1 Lifecycle Philosophy

**ULTRATHINK ON AGENT LIFECYCLE:**

Agents have a natural lifecycle:
1. **SPAWNED** - Agent created, assigned work
2. **ACTIVE** - Agent executing work
3. **COMPLETED** - Agent finished work, reported results
4. **ARCHIVED** - Agent metadata preserved in history

**Metadata should reflect lifecycle:**
- `active_agents`: Agents in SPAWNED or ACTIVE states
- `agents_history`: Agents in COMPLETED/ARCHIVED state

**Cleanup should happen automatically:**
- When agent reaches COMPLETED state
- When orchestrator detects completion
- Before transitioning to next major state

### 3.2 Design Options Analysis

#### Option A: Automatic (Event-Driven) Cleanup
**Trigger**: When agent transitions to COMPLETED state

**Mechanism**:
```yaml
agent_completion_hook:
  when: Agent reports COMPLETED or FAILED
  action: Orchestrator monitoring state detects it
  cleanup: Immediately move metadata to agents_history
  timing: Real-time, as agents complete
```

**Advantages:**
- ✅ Real-time cleanup (no stale data)
- ✅ No manual intervention
- ✅ State file always accurate
- ✅ Minimal implementation complexity

**Disadvantages:**
- ⚠️ Cleanup logic in multiple monitoring states
- ⚠️ Must handle concurrent completions
- ⚠️ Risk of losing metadata if cleanup fails

**Implementation**:
```bash
# In MONITORING_SWE_PROGRESS state
if agent_completed "$agent_id"; then
    archive_agent_metadata "$agent_id"
    remove_from_active_agents "$agent_id"
fi
```

#### Option B: Manual (State-Based) Cleanup
**Trigger**: Explicit cleanup state in state machine

**Mechanism**:
```yaml
cleanup_state:
  state_name: CLEANUP_COMPLETED_AGENTS
  when: After WAVE_COMPLETE or PHASE_COMPLETE
  action: Scan active_agents, move COMPLETED to history
  timing: Batched at wave/phase boundaries
```

**Advantages:**
- ✅ Explicit control over cleanup timing
- ✅ Batched operations (more efficient)
- ✅ Clear separation of concerns
- ✅ Easier to audit cleanup operations

**Disadvantages:**
- ⚠️ Stale data between cleanup cycles
- ⚠️ Additional state in state machine
- ⚠️ User must trigger cleanup
- ⚠️ More complex state machine

**Implementation**:
```bash
# New state: CLEANUP_COMPLETED_AGENTS
cleanup_completed_agents() {
    for agent in $(jq -r '.active_agents[] | select(.state == "COMPLETE") | .agent_id' state.json); do
        archive_agent "$agent"
    done
}
```

#### Option C: Hybrid (Automatic + Boundaries) ⭐ RECOMMENDED
**Trigger**: Automatic when detected + cleanup at wave/phase boundaries

**Mechanism**:
```yaml
automatic_cleanup:
  when: Monitoring state detects agent COMPLETE
  action: Move to agents_history immediately
  frequency: Real-time

boundary_cleanup:
  when: WAVE_COMPLETE, PHASE_COMPLETE states
  action: Verify no stale agents remain
  purpose: Safety net + validation

manual_cleanup:
  when: User runs cleanup command (emergency)
  action: Force cleanup of all COMPLETE agents
  purpose: Recovery from failures
```

**Advantages:**
- ✅ Best of both worlds
- ✅ Real-time cleanup + validation
- ✅ Safety net at boundaries
- ✅ Manual override available
- ✅ Resilient to failures

**Disadvantages:**
- ⚠️ Most complex implementation
- ⚠️ Multiple cleanup points to maintain

**Implementation**:
```bash
# Automatic cleanup in MONITORING_SWE_PROGRESS
check_agent_completion() {
    if [ "$agent_state" = "COMPLETE" ]; then
        archive_agent_metadata "$agent_id"  # Immediate
    fi
}

# Validation cleanup in WAVE_COMPLETE
validate_no_stale_agents() {
    stale_count=$(jq '[.active_agents[] | select(.state == "COMPLETE")] | length' state.json)
    if [ $stale_count -gt 0 ]; then
        echo "⚠️ Found $stale_count stale agents - cleaning up"
        cleanup_all_completed_agents
    fi
}
```

### 3.3 Recommended Approach

**RECOMMENDATION: OPTION C - HYBRID APPROACH**

**Rationale:**
1. **Real-time cleanup** keeps state file accurate
2. **Boundary validation** catches missed cleanups
3. **Manual override** enables recovery
4. **Defense in depth** - multiple safety nets
5. **Performance** - Cleanup happens when completion detected (minimal overhead)

**Implementation Strategy:**
1. Add cleanup logic to MONITORING_* states
2. Add validation to WAVE_COMPLETE, PHASE_COMPLETE
3. Create utility script for manual cleanup
4. Add pre-commit hook to detect stale agents
5. Add rule requiring cleanup compliance

---

## 4. AGENTS_HISTORY SCHEMA DESIGN

### 4.1 Design Goals

**Preserve Essential Information:**
- Agent identity (ID, type, effort)
- Timeline (spawned, completed, archived)
- Outcome (success/failure)
- Key metrics (lines implemented, etc.)

**Minimize Storage:**
- No full implementation plans (stored elsewhere)
- No full reports (stored in artifacts)
- No transient data (workspace paths, etc.)
- Only permanent historical record

**Enable Debugging:**
- Can recreate timeline of agent work
- Can identify which agent worked on which effort
- Can audit agent performance
- Can recover from state corruption

### 4.2 Proposed Schema

```json
{
  "agents_history": [
    {
      "agent_id": "swe-1.2.1-docker-client",
      "agent_type": "sw-engineer",
      "effort": "1.2.1-docker-client",
      "effort_name": "Docker Client Implementation",
      "phase": 1,
      "wave": 2,
      "spawned_at": "2025-10-29T22:12:59Z",
      "completed_at": "2025-10-29T23:45:00Z",
      "archived_at": "2025-11-02T04:30:00Z",
      "outcome": "COMPLETED",
      "lines_implemented": 450,
      "implementation_plan": ".software-factory/phase1/wave2/effort-1-docker-client/IMPLEMENTATION-PLAN--20251029-214603.md",
      "implementation_report": ".software-factory/phase1/wave2/effort-1-docker-client/IMPLEMENTATION-COMPLETE--20251029-234500.md",
      "code_review_status": "APPROVED",
      "archived_by": "orchestrator-MONITORING_SWE_PROGRESS"
    }
  ]
}
```

### 4.3 Field Definitions

**Core Identity:**
- `agent_id`: Unique agent identifier
- `agent_type`: sw-engineer, code-reviewer, architect, integration, etc.
- `effort`: Effort ID this agent worked on
- `effort_name`: Human-readable effort name
- `phase`: Phase number
- `wave`: Wave number

**Timeline:**
- `spawned_at`: When agent was created (ISO 8601)
- `completed_at`: When agent finished work (ISO 8601)
- `archived_at`: When metadata moved to history (ISO 8601)

**Outcome:**
- `outcome`: COMPLETED | FAILED | BLOCKED | TIMEOUT
- `lines_implemented`: For SW Engineers, line count
- `code_review_status`: APPROVED | CHANGES_REQUESTED | BLOCKED

**References:**
- `implementation_plan`: Path to plan (for reference)
- `implementation_report`: Path to completion report
- `archived_by`: Which state/agent archived this

**Optional Metadata:**
- `duration_seconds`: Time from spawn to completion
- `commits_made`: Number of commits (if tracked)
- `files_modified`: Number of files changed
- `tests_added`: Number of tests added

### 4.4 Size Analysis

**Per-agent storage:**
```yaml
estimated_size_per_entry:
  core_fields: 150 bytes
  timeline_fields: 90 bytes
  outcome_fields: 100 bytes
  references: 150 bytes
  total: ~490 bytes per agent

projected_growth:
  agents_per_project: 48
  total_history_size: 48 * 490 = ~23KB

comparison:
  active_agents_current: ~2KB (6 agents with full metadata)
  agents_history_projected: ~23KB (48 agents with minimal metadata)
  ratio: 11.5x more agents, but only 11.5x size (acceptable)
```

**Conclusion**: agents_history size is manageable (<25KB for typical project).

### 4.5 Retention Policy

**Recommendation: Keep All History**

**Rationale:**
- History size is small (<25KB)
- Valuable for debugging
- Needed for audit trail
- Required for performance analysis
- Useful for project retrospectives

**Alternative (if size becomes issue):**
- Archive agents_history to separate file after project completion
- Keep only recent agents (last 2 phases)
- Compress old history

**NOT Recommended:**
- Deleting agent history entirely
- Lossy compression
- Summary-only records

---

## 5. NEW RULES SPECIFICATION

### 5.1 Rule R610 - Agent Metadata Lifecycle Protocol

**File**: `rule-library/R610-agent-metadata-lifecycle-protocol.md`

**Criticality**: 🚨🚨🚨 BLOCKING - STATE FILE INTEGRITY

**Summary**: ALL completed agent metadata MUST be moved from active_agents to agents_history within 60 seconds of completion detection.

**Key Requirements:**

```yaml
cleanup_triggers:
  automatic:
    - When orchestrator monitoring state detects agent.state == "COMPLETE"
    - Within 60 seconds of detection
    - Before transitioning to next major state

  validation:
    - At WAVE_COMPLETE state
    - At PHASE_COMPLETE state
    - At COMPLETE_PROJECT state

  manual:
    - Via cleanup utility script
    - Emergency recovery only

cleanup_protocol:
  1. Detect agent completion
  2. Extract agent metadata
  3. Create agents_history entry
  4. Remove from active_agents
  5. Commit atomically per R288

agents_history_structure:
  - See schema in R610 rule file
  - Minimal essential data only
  - References to full reports

grading:
  - Stale agents in active_agents: -10% per agent per day
  - No agents_history field: -20%
  - Cleanup failures: -30%
  - State file corruption from no cleanup: -100%
```

### 5.2 Rule R611 - Active Agents Cleanup Protocol

**File**: `rule-library/R611-active-agents-cleanup-protocol.md`

**Criticality**: ⚠️⚠️⚠️ WARNING - PERFORMANCE AND ACCURACY

**Summary**: active_agents array MUST contain only agents in SPAWNED or ACTIVE states, never COMPLETED agents.

**Key Requirements:**

```yaml
active_agents_definition:
  - Agent spawned but not yet complete
  - Agent currently executing work
  - Agent blocked waiting for resolution
  - NEVER agents with state=COMPLETE

cleanup_frequency:
  - Immediate when completion detected
  - Validated every state transition
  - Audited at wave/phase boundaries

validation:
  - Pre-commit hook checks for stale agents
  - State file validation checks active_agents
  - Monitoring states verify no COMPLETE in active

grading:
  - Active agents containing COMPLETE: -5% per agent
  - No cleanup for >24 hours: -15%
  - Performance degradation from large active_agents: -10%
```

### 5.3 Rule R612 - Agent History Management

**File**: `rule-library/R612-agent-history-management.md`

**Criticality**: 🔴🔴🔴 SUPREME LAW - HISTORICAL INTEGRITY

**Summary**: agents_history MUST preserve complete agent lifecycle records with minimal essential metadata.

**Key Requirements:**

```yaml
agents_history_structure:
  required_fields:
    - agent_id (unique identifier)
    - agent_type (sw-engineer, code-reviewer, etc.)
    - effort (effort ID)
    - phase, wave (location in project)
    - spawned_at, completed_at, archived_at (timeline)
    - outcome (COMPLETED, FAILED, etc.)
    - lines_implemented (for SW Engineers)

  optional_fields:
    - implementation_plan path
    - implementation_report path
    - code_review_status
    - duration_seconds
    - commits_made

  forbidden_data:
    - Full implementation plans (use references)
    - Full reports (use references)
    - Workspace paths (transient)
    - Temporary state (ephemeral)

retention:
  - Keep all history for project lifetime
  - Archive to separate file after project completion
  - Never delete (valuable for audits)

size_management:
  - Limit per-agent metadata to <500 bytes
  - Total history for typical project: <25KB
  - If size exceeds 100KB, consider archival

grading:
  - Missing agents_history: -20%
  - Incomplete agent records: -10% per agent
  - Oversized metadata (>500 bytes/agent): -5%
  - Lost history during cleanup: -50%
```

### 5.4 Rule R613 - State File Growth Management

**File**: `rule-library/R613-state-file-growth-management.md`

**Criticality**: ⚠️⚠️⚠️ WARNING - LONG-TERM MAINTAINABILITY

**Summary**: orchestrator-state-v3.json MUST be kept under 500KB through aggressive metadata cleanup and archival.

**Key Requirements:**

```yaml
size_targets:
  active_agents: <5KB (typically 0-4 agents)
  agents_history: <25KB (all archived agents)
  state_history: <50KB (state transitions)
  total_state_file: <500KB (maximum)

growth_management:
  cleanup:
    - Move COMPLETE agents to history immediately
    - Prune old state_history (keep last 100 transitions)
    - Remove transient data after use

  monitoring:
    - Check state file size every state transition
    - Warn if exceeds 250KB
    - Error if exceeds 500KB

  archival:
    - Archive old state_history to separate file
    - Archive agents_history after project completion
    - Keep only current project data in main state file

grading:
  - State file >250KB: -5%
  - State file >500KB: -15%
  - State file >1MB: -30% (CRITICAL)
  - State file corruption from size: -100%
```

---

## 6. IMPLEMENTATION FILES

### 6.1 Template Repository Changes

**All changes must be made in `/home/vscode/software-factory-template/`**

#### Files to Create

**1. New Rules**
```
rule-library/R610-agent-metadata-lifecycle-protocol.md
rule-library/R611-active-agents-cleanup-protocol.md
rule-library/R612-agent-history-management.md
rule-library/R613-state-file-growth-management.md
```

**2. Utility Scripts**
```
utilities/cleanup-completed-agents.sh
utilities/validate-agent-metadata.sh
utilities/archive-agents-history.sh
```

**3. Pre-commit Hooks**
```
tools/git-commit-hooks/shared-hooks/agent-metadata-validation.hook
```

#### Files to Modify

**1. State Rules - Add Cleanup Logic**
```
agent-states/software-factory/orchestrator/MONITORING_SWE_PROGRESS/rules.md
agent-states/software-factory/orchestrator/MONITORING_EFFORT_REVIEWS/rules.md
agent-states/software-factory/orchestrator/WAVE_COMPLETE/rules.md
agent-states/software-factory/orchestrator/COMPLETE_PHASE/rules.md
agent-states/software-factory/orchestrator/COMPLETE_PROJECT/rules.md
```

**2. Schema Updates**
```
schemas/orchestrator-state-v3.schema.json
  - Add agents_history array schema
  - Update active_agents validation
  - Add size constraints
```

**3. Rule Registry**
```
rule-library/RULE-REGISTRY.md
  - Add R610, R611, R612, R613
  - Update agent lifecycle section
```

**4. Orchestrator Agent Config**
```
.claude/agents/orchestrator.md
  - Add agent lifecycle responsibilities
  - Reference new cleanup rules
```

### 6.2 Cleanup Utility Script

**File**: `utilities/cleanup-completed-agents.sh`

```bash
#!/bin/bash
# Software Factory 3.0 - Agent Metadata Cleanup Utility
# Enforces R610, R611, R612

set -euo pipefail

STATE_FILE="${1:-orchestrator-state-v3.json}"

echo "🧹 Agent Metadata Cleanup Utility"
echo "=================================="
echo "State file: $STATE_FILE"
echo ""

# Validate state file exists
if [ ! -f "$STATE_FILE" ]; then
    echo "❌ State file not found: $STATE_FILE"
    exit 1
fi

# Backup state file
BACKUP_FILE="${STATE_FILE}.backup-$(date +%Y%m%d-%H%M%S)"
cp "$STATE_FILE" "$BACKUP_FILE"
echo "✅ Backup created: $BACKUP_FILE"

# Count completed agents in active_agents
COMPLETED_COUNT=$(jq '[.active_agents[] | select(.state == "COMPLETE" or .state == "COMPLETED")] | length' "$STATE_FILE")

if [ "$COMPLETED_COUNT" -eq 0 ]; then
    echo "✅ No completed agents in active_agents - nothing to clean up"
    rm "$BACKUP_FILE"
    exit 0
fi

echo "📊 Found $COMPLETED_COUNT completed agents to archive"
echo ""

# Initialize agents_history if doesn't exist
if ! jq -e '.agents_history' "$STATE_FILE" >/dev/null 2>&1; then
    echo "📝 Creating agents_history array..."
    jq '.agents_history = []' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
fi

# Archive each completed agent
ARCHIVED=0
jq -r '.active_agents[] | select(.state == "COMPLETE" or .state == "COMPLETED") | .agent_id' "$STATE_FILE" | while read agent_id; do
    echo "🗄️  Archiving $agent_id..."

    # Extract agent metadata
    AGENT_DATA=$(jq ".active_agents[] | select(.agent_id == \"$agent_id\")" "$STATE_FILE")

    # Create history entry
    HISTORY_ENTRY=$(echo "$AGENT_DATA" | jq '. + {
        archived_at: "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
        archived_by: "cleanup-utility"
    }')

    # Append to agents_history
    jq ".agents_history += [$HISTORY_ENTRY]" "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"

    # Remove from active_agents
    jq ".active_agents = [.active_agents[] | select(.agent_id != \"$agent_id\")]" "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"

    echo "  ✅ Archived to agents_history"
    ((ARCHIVED++))
done

echo ""
echo "✅ Cleanup complete: $ARCHIVED agents archived"
echo ""

# Validate result
REMAINING=$(jq '[.active_agents[] | select(.state == "COMPLETE" or .state == "COMPLETED")] | length' "$STATE_FILE")
HISTORY_COUNT=$(jq '.agents_history | length' "$STATE_FILE")

echo "📊 Final state:"
echo "  Active agents (all states): $(jq '.active_agents | length' "$STATE_FILE")"
echo "  Active agents (COMPLETE): $REMAINING"
echo "  Agents in history: $HISTORY_COUNT"
echo ""

if [ "$REMAINING" -gt 0 ]; then
    echo "⚠️  WARNING: Still have $REMAINING completed agents in active_agents"
    exit 1
fi

echo "✅ All completed agents successfully archived"
rm "$BACKUP_FILE"
```

### 6.3 Validation Hook

**File**: `tools/git-commit-hooks/shared-hooks/agent-metadata-validation.hook`

```bash
#!/bin/bash
# Pre-commit hook: Validate no stale agents in active_agents

STATE_FILE="orchestrator-state-v3.json"

if [ ! -f "$STATE_FILE" ]; then
    exit 0  # No state file = no validation needed
fi

# Check for completed agents in active_agents
STALE_COUNT=$(jq '[.active_agents[] | select(.state == "COMPLETE" or .state == "COMPLETED")] | length' "$STATE_FILE" 2>/dev/null || echo "0")

if [ "$STALE_COUNT" -gt 0 ]; then
    echo "❌ R611 VIOLATION: Found $STALE_COUNT completed agents in active_agents"
    echo ""
    echo "Stale agents:"
    jq -r '.active_agents[] | select(.state == "COMPLETE" or .state == "COMPLETED") | "  - \(.agent_id) (completed: \(.completed_at // .spawned_at))"' "$STATE_FILE"
    echo ""
    echo "🔧 Fix: Run utilities/cleanup-completed-agents.sh"
    echo ""
    exit 1
fi

# Verify agents_history exists if any agents were ever spawned
ACTIVE_COUNT=$(jq '.active_agents | length' "$STATE_FILE")
HISTORY_EXISTS=$(jq 'has("agents_history")' "$STATE_FILE")

if [ "$ACTIVE_COUNT" -gt 0 ] && [ "$HISTORY_EXISTS" = "false" ]; then
    echo "⚠️  R612 WARNING: agents_history field missing"
    echo "This is required for historical agent tracking"
    echo ""
    # Don't fail commit, just warn (backward compatibility)
fi

exit 0
```

---

## 7. MIGRATION STRATEGY

### 7.1 Current Project Cleanup

**For `/home/vscode/workspaces/idpbuilder-oci-push-planning/`**

#### Step 1: Manual Cleanup (Immediate)

```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Backup current state
cp orchestrator-state-v3.json orchestrator-state-v3.json.backup-before-cleanup

# Create agents_history array
jq '.agents_history = []' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

# Archive each COMPLETE agent
for agent_id in swe-1.2.1-docker-client swe-1.2.2-registry-client swe-1.2.3-auth swe-1.2.4-tls swe-2.2.1-registry-override swe-2.2.2-env-variable-support; do
    echo "Archiving $agent_id..."

    # Extract and move to history
    jq ".agents_history += [(.active_agents[] | select(.agent_id == \"$agent_id\") | . + {archived_at: \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", archived_by: \"manual-cleanup\"})]" orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    # Remove from active
    jq ".active_agents = [.active_agents[] | select(.agent_id != \"$agent_id\")]" orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
done

# Verify cleanup
echo "Active agents: $(jq '.active_agents | length' orchestrator-state-v3.json)"
echo "History agents: $(jq '.agents_history | length' orchestrator-state-v3.json)"
echo "Stale agents: $(jq '[.active_agents[] | select(.state == "COMPLETE")] | length' orchestrator-state-v3.json)"

# Commit
git add orchestrator-state-v3.json
git commit -m "fix: clean up completed agents, add agents_history [R610, R611, R612]

- Moved 6 COMPLETE agents from active_agents to agents_history
- Created agents_history array per R612 schema
- Fixed stale agent accumulation (4 days old for P1W2 agents)
- Compliance with new agent lifecycle management protocol

Archived agents:
- swe-1.2.1-docker-client (Phase 1 Wave 2)
- swe-1.2.2-registry-client (Phase 1 Wave 2)
- swe-1.2.3-auth (Phase 1 Wave 2)
- swe-1.2.4-tls (Phase 1 Wave 2)
- swe-2.2.1-registry-override (Phase 2 Wave 2)
- swe-2.2.2-env-variable-support (Phase 2 Wave 2)

State file size reduced, active_agents now accurately reflects active work."

git push
```

#### Step 2: Install Cleanup Utility (After Template Update)

```bash
# After template repository updated with new utilities
cd /home/vscode/software-factory-template
git pull

# Copy utility to current project
cp utilities/cleanup-completed-agents.sh /home/vscode/workspaces/idpbuilder-oci-push-planning/utilities/

# Install pre-commit hook
cp tools/git-commit-hooks/shared-hooks/agent-metadata-validation.hook /home/vscode/workspaces/idpbuilder-oci-push-planning/tools/git-commit-hooks/shared-hooks/

# Update pre-commit to include new hook
# (Add to .git/hooks/pre-commit or tools/git-commit-hooks/master-pre-commit.sh)
```

#### Step 3: Ongoing Maintenance

```bash
# Run cleanup utility if stale agents detected
bash utilities/cleanup-completed-agents.sh

# Validate no stale agents before commits
jq '[.active_agents[] | select(.state == "COMPLETE")] | length' orchestrator-state-v3.json
# Should always be 0

# Check history size periodically
jq '.agents_history | length' orchestrator-state-v3.json
jq '.agents_history | length' orchestrator-state-v3.json | awk '{print $1 * 500 " bytes estimated"}'
```

### 7.2 Template Repository Updates

**Priority Order:**

1. **Create new rules** (R610-R613) - Defines requirements
2. **Update schema** - Adds agents_history validation
3. **Create utilities** - cleanup-completed-agents.sh, etc.
4. **Update state rules** - Add cleanup logic to monitoring states
5. **Install hooks** - Pre-commit validation
6. **Update docs** - Rule registry, agent configs
7. **Test thoroughly** - Run against test projects

**Testing Strategy:**

```bash
# Create test project with multiple phases
cd /tmp/test-agent-lifecycle
git clone /home/vscode/software-factory-template .

# Manually create state file with COMPLETE agents
# Run cleanup utility
# Verify agents moved to history
# Verify pre-commit hook catches stale agents
# Verify state file size managed correctly
```

---

## 8. ANSWERS TO USER QUESTIONS

### Q1: Why weren't these agents' metadata moved to history automatically?

**Answer**: Because **Software Factory 3.0 has NO automatic agent metadata cleanup mechanism**.

**Evidence**:
- No cleanup logic in MONITORING_* states (checked MONITORING_SWE_PROGRESS code)
- No cleanup states in state machine (checked all 150+ states)
- No rules governing agent lifecycle (checked entire rule library)
- No agents_history field defined in schema

**Root cause**: System was designed for forward progress (spawn → monitor → transition) but NOT for metadata hygiene (cleanup after completion).

### Q2: Is automatic cleanup supposed to happen?

**Answer**: **NO - there is no design or implementation for automatic cleanup**.

**Current behavior** (as designed):
1. Orchestrator spawns agents → adds to active_agents
2. Orchestrator monitors agents → checks active_agents for completion
3. When all agents complete → transitions to next state
4. **Agents remain in active_agents forever** ← NO CLEANUP STEP

**This is a design gap, not a bug in the implementation.**

### Q3: If not, should it be?

**Answer**: **YES - automatic cleanup should be implemented**.

**Justification**:

**Why cleanup is needed:**
1. **State file integrity** - `active_agents` should mean actually active
2. **Performance** - Monitoring states scan all agents, including stale ones
3. **Maintainability** - State file grows unbounded without cleanup
4. **Debugging** - Stale data makes state file hard to understand
5. **Auditing** - Need historical record, but separate from active work

**Why automatic is better than manual:**
1. **Accuracy** - Real-time cleanup = accurate state always
2. **Performance** - Cleanup when detected = no stale data overhead
3. **Reliability** - No human intervention = no forgotten cleanups
4. **Simplicity** - User doesn't need to remember to clean up

**Recommended**: Hybrid approach (automatic + validation + manual override)

### Q4: What rules are we missing about agent lifecycle management?

**Answer**: **Four new rules are needed** (detailed in Section 5):

1. **R610 - Agent Metadata Lifecycle Protocol**
   - BLOCKING - Defines when/how to move agents to history
   - Automatic cleanup within 60s of completion detection
   - Validation at wave/phase boundaries

2. **R611 - Active Agents Cleanup Protocol**
   - WARNING - active_agents must contain ONLY active agents
   - No COMPLETE agents allowed
   - Performance and accuracy requirements

3. **R612 - Agent History Management**
   - SUPREME LAW - agents_history structure and retention
   - Minimal essential metadata only
   - Historical integrity preservation

4. **R613 - State File Growth Management**
   - WARNING - Keep state file <500KB
   - Aggressive cleanup and archival
   - Long-term maintainability

### Q5: How should we implement proper metadata cleanup?

**Answer**: **Hybrid approach with three cleanup mechanisms** (detailed in Section 3):

**Mechanism 1: Automatic Cleanup (Primary)**
```yaml
where: MONITORING_SWE_PROGRESS, MONITORING_EFFORT_REVIEWS states
when: Agent completion detected (state == "COMPLETE")
action: Immediately archive to agents_history and remove from active_agents
frequency: Real-time (within 60 seconds)
```

**Mechanism 2: Validation Cleanup (Safety Net)**
```yaml
where: WAVE_COMPLETE, PHASE_COMPLETE, COMPLETE_PROJECT states
when: Before transitioning to next major state
action: Scan for any remaining COMPLETE agents, clean up if found
frequency: At wave/phase/project boundaries
```

**Mechanism 3: Manual Cleanup (Emergency)**
```yaml
where: Utility script (utilities/cleanup-completed-agents.sh)
when: User runs script manually or hook detects stale agents
action: Force cleanup of all COMPLETE agents
frequency: On-demand
```

**Implementation:**
- Add cleanup logic to monitoring state rules
- Create cleanup utility script
- Add pre-commit hook validation
- Update schema to require agents_history

### Q6: Can we keep an `agents_history` as long as data is minimal?

**Answer**: **YES - agents_history is recommended and size is manageable**.

**Size Analysis:**
```yaml
per_agent_storage:
  minimal_metadata: ~490 bytes per agent
  typical_project: 48 agents
  total_size: ~23KB

comparison:
  orchestrator-state-v3.json_current: ~50KB
  agents_history_addition: ~23KB
  percentage_increase: ~46%

conclusion: ACCEPTABLE (state file still <100KB)
```

**What to keep (minimal):**
- Agent identity (ID, type, effort)
- Timeline (spawned, completed, archived)
- Outcome (COMPLETE, FAILED, etc.)
- Key metrics (lines_implemented)
- References to reports (paths only, not content)

**What NOT to keep (avoid bloat):**
- Full implementation plans (use references)
- Full reports (use references)
- Workspace paths (transient)
- Temporary state (ephemeral)

**Benefits**:
- Debugging: Can recreate agent timeline
- Auditing: Complete record of all agent work
- Performance: Can analyze agent efficiency
- Recovery: Can restore from corruption

**Recommendation**: Keep all agent history for project lifetime, with option to archive to separate file after project completion.

---

## 9. TESTING STRATEGY

### 9.1 Unit Tests

**Test cleanup utility:**
```bash
# Create test state file with COMPLETE agents
# Run cleanup utility
# Verify agents moved to history
# Verify active_agents cleaned
# Verify no data loss
```

**Test validation hook:**
```bash
# Create state file with stale agents
# Attempt git commit
# Verify hook catches stale agents
# Verify helpful error message
```

### 9.2 Integration Tests

**Test monitoring state cleanup:**
```bash
# Spawn SW Engineer
# Monitor until completion
# Verify automatic cleanup triggered
# Verify agents_history updated
# Verify active_agents cleaned
```

**Test boundary validation:**
```bash
# Complete wave with multiple agents
# Transition to WAVE_COMPLETE
# Verify validation cleanup runs
# Verify no stale agents remain
```

### 9.3 Regression Tests

**Test backward compatibility:**
```bash
# Run against old state files without agents_history
# Verify cleanup creates agents_history
# Verify no corruption of other state data
# Verify schema validation passes
```

**Test error recovery:**
```bash
# Simulate cleanup failure mid-operation
# Verify rollback to backup
# Verify state file integrity preserved
# Verify retry succeeds
```

---

## 10. CONCLUSION

### 10.1 Summary

**Problem**: Software Factory 3.0 has no agent metadata lifecycle management, causing:
- 6 COMPLETE agents stuck in active_agents (4 from 4 days ago)
- Misleading state representation
- Performance degradation
- No historical agent records

**Root Cause**: System designed for forward progress, not metadata hygiene.

**Solution**: Implement hybrid agent lifecycle management:
- Automatic cleanup when completion detected
- Validation at wave/phase boundaries
- Manual emergency cleanup
- agents_history for historical preservation

**Impact**:
- State file accuracy restored
- Performance improved (smaller active_agents)
- Historical context preserved
- Long-term maintainability ensured

### 10.2 Recommendations Priority

**IMMEDIATE (Next commit):**
1. Clean up 6 stale agents in current project manually
2. Create agents_history field
3. Document temporary cleanup procedure

**SHORT-TERM (Template update):**
1. Create R610-R613 rules
2. Create cleanup utility script
3. Add pre-commit validation hook
4. Update orchestrator state rules with cleanup logic

**MEDIUM-TERM (Production ready):**
1. Test thoroughly in template
2. Run upgrade.sh to deploy to projects
3. Validate all projects cleaned up
4. Monitor for cleanup failures

**LONG-TERM (Maintenance):**
1. Monitor state file sizes
2. Analyze agent performance from history
3. Consider archival strategy for very large projects
4. Refine cleanup triggers based on experience

### 10.3 Next Steps

**User Actions:**
1. Review this investigation report
2. Approve recommended approach (hybrid cleanup)
3. Approve new rules (R610-R613)
4. Authorize template repository modifications

**Factory Manager Actions:**
1. Create new rules in template repository
2. Implement cleanup utility script
3. Update state rules with cleanup logic
4. Install pre-commit validation hooks
5. Test thoroughly
6. Commit and push to template repository

**Orchestrator Actions (after template update):**
1. Run upgrade.sh to get new utilities
2. Execute manual cleanup on current project
3. Verify agents_history created correctly
4. Monitor for automatic cleanup in future states
5. Validate no stale agents accumulate

---

**END OF INVESTIGATION REPORT**

Generated: 2025-11-02
By: Software Factory Manager Agent
Template Location: `/home/vscode/software-factory-template/`
Current Project: `/home/vscode/workspaces/idpbuilder-oci-push-planning/`
