# STALE INTEGRATION TRACKING MECHANISM

## Overview
This document defines the comprehensive tracking mechanism for stale integrations in the Software Factory orchestrator-state.json file. This mechanism enforces R327 (Mandatory Re-Integration After Fixes) and R328 (Integration Freshness Validation).

## Core Tracking Structure

### 1. Main Tracking Object
Located at root level in orchestrator-state.json:

```json
"stale_integration_tracking": {
  "_comment": "Comprehensive tracking of stale integrations per R327 cascade requirements",
  "stale_integrations": [...],
  "staleness_cascade": [...],
  "fix_tracking": {...},
  "validation_history": [...]
}
```

### 2. Stale Integration Records

Each stale integration is tracked with:

```json
{
  "integration_type": "wave|phase|project",
  "integration_id": "phase1-wave2-integration",
  "became_stale_at": "2025-08-25T21:00:00Z",
  "stale_reason": "Fixes applied to source efforts after integration creation",
  "triggering_fixes": [
    {
      "fix_id": "FIX-001",
      "commit": "abc123def",
      "branch": "phase1/wave2/effort3",
      "applied_at": "2025-08-25T21:00:00Z",
      "description": "Fixed build error in authentication module"
    }
  ],
  "affected_downstream": ["phase1-integration", "project-integration"],
  "recreation_required": true,
  "recreation_completed": false,
  "recreation_at": null
}
```

### 3. Staleness Cascade Tracking

Tracks the cascade effect of stale integrations:

```json
"staleness_cascade": [
  {
    "trigger": {
      "type": "effort_fix",
      "branch": "phase1/wave2/effort3",
      "commit": "abc123",
      "timestamp": "2025-08-25T21:00:00Z"
    },
    "cascade_sequence": [
      {
        "level": "wave",
        "integration": "phase1-wave2-integration",
        "became_stale": "2025-08-25T21:00:00Z",
        "must_recreate": true
      },
      {
        "level": "phase",
        "integration": "phase1-integration",
        "became_stale": "2025-08-25T21:00:00Z",
        "must_recreate": true
      },
      {
        "level": "project",
        "integration": "project-integration",
        "became_stale": "2025-08-25T21:00:00Z",
        "must_recreate": true
      }
    ],
    "cascade_status": "pending|in_progress|completed",
    "cascade_started_at": null,
    "cascade_completed_at": null
  }
]
```

### 4. Fix Tracking

Comprehensive tracking of all fixes and their integration status:

```json
"fix_tracking": {
  "fixes_applied": [
    {
      "fix_id": "FIX-001",
      "commit": "abc123def",
      "branch": "phase1/wave2/effort3",
      "effort_name": "authentication-module",
      "applied_at": "2025-08-25T21:00:00Z",
      "type": "build_fix|test_fix|review_fix",
      "description": "Fixed missing import in auth handler",
      "integrated_into": {
        "wave": false,
        "phase": false,
        "project": false
      },
      "made_stale": [
        "phase1-wave2-integration",
        "phase1-integration",
        "project-integration"
      ]
    }
  ],
  "fixes_pending_integration": [
    {
      "fix_id": "FIX-002",
      "branch": "phase1/wave3/effort1",
      "pending_since": "2025-08-26T10:00:00Z",
      "blocks_integrations": ["phase1-wave3-integration"],
      "priority": "high"
    }
  ]
}
```

### 5. Integration Object Enhancement

Each integration object (wave/phase/project) gets enhanced tracking:

```json
"current_wave_integration": {
  "workspace": "/path/to/wave/integration-workspace",
  "branch": "wave-integration-branch-name",
  "created_at": "2025-08-25T10:00:00Z",
  "efforts_last_updated_at": "2025-08-25T09:30:00Z",
  "is_stale": false,
  "staleness_reason": null,
  "stale_since": null,
  "stale_due_to_fixes": [],
  "merged_efforts": ["effort1", "effort2", "effort3"],
  "last_freshness_check": "2025-08-25T10:05:00Z",
  "freshness_validation": {
    "last_check": "2025-08-25T10:05:00Z",
    "next_required": "2025-08-25T10:20:00Z",
    "check_frequency_minutes": 15,
    "validation_method": "timestamp_comparison|commit_hash|both"
  }
}
```

### 6. Validation History

Tracks all freshness validations performed:

```json
"validation_history": [
  {
    "timestamp": "2025-08-25T10:05:00Z",
    "integration_type": "wave",
    "integration_id": "phase1-wave2-integration",
    "result": "fresh|stale",
    "stale_components": [],
    "action_taken": "none|recreated|marked_stale",
    "validator": "orchestrator|manual_check"
  }
]
```

## Usage Patterns

### 1. Detecting Staleness

```bash
# Function to detect and record staleness
detect_staleness() {
    local integration_type=$1
    local integration_id=$2
    
    # Get integration creation time
    local created_at=$(jq -r ".current_${integration_type}_integration.created_at" orchestrator-state.json)
    
    # Check for newer fixes in source branches
    # ... detection logic ...
    
    if [[ "$is_stale" == "true" ]]; then
        # Record in stale_integration_tracking
        jq '.stale_integration_tracking.stale_integrations += [{
            "integration_type": "'$integration_type'",
            "integration_id": "'$integration_id'",
            "became_stale_at": "'$(date -Iseconds)'",
            "stale_reason": "Fixes applied after integration",
            "triggering_fixes": [...],
            "recreation_required": true
        }]' orchestrator-state.json > tmp.json
        mv tmp.json orchestrator-state.json
    fi
}
```

### 2. Cascade Tracking

```bash
# Track cascade when fix is applied
track_cascade() {
    local fix_branch=$1
    local fix_commit=$2
    
    # Determine cascade impact
    local affected_wave=$(determine_wave_from_branch "$fix_branch")
    local affected_phase=$(determine_phase_from_branch "$fix_branch")
    
    # Record cascade
    jq '.stale_integration_tracking.staleness_cascade += [{
        "trigger": {
            "type": "effort_fix",
            "branch": "'$fix_branch'",
            "commit": "'$fix_commit'",
            "timestamp": "'$(date -Iseconds)'"
        },
        "cascade_sequence": [
            {"level": "wave", "integration": "'$affected_wave'-integration", "must_recreate": true},
            {"level": "phase", "integration": "'$affected_phase'-integration", "must_recreate": true},
            {"level": "project", "integration": "project-integration", "must_recreate": true}
        ],
        "cascade_status": "pending"
    }]' orchestrator-state.json > tmp.json
    mv tmp.json orchestrator-state.json
}
```

### 3. Querying Stale Integrations

```bash
# Get all currently stale integrations
get_stale_integrations() {
    jq -r '.stale_integration_tracking.stale_integrations[] | 
           select(.recreation_completed == false) | 
           "\(.integration_type): \(.integration_id) - Stale since: \(.became_stale_at)"' \
           orchestrator-state.json
}

# Get fixes that made integrations stale
get_staleness_causes() {
    local integration_id=$1
    jq -r '.stale_integration_tracking.stale_integrations[] | 
           select(.integration_id == "'$integration_id'") | 
           .triggering_fixes[] | 
           "Fix: \(.commit) in \(.branch) - \(.description)"' \
           orchestrator-state.json
}
```

### 4. Recreation Tracking

```bash
# Mark integration as recreated
mark_integration_recreated() {
    local integration_id=$1
    
    # Update stale integration record
    jq '(.stale_integration_tracking.stale_integrations[] | 
         select(.integration_id == "'$integration_id'")) |= 
         . + {"recreation_completed": true, "recreation_at": "'$(date -Iseconds)'"}' \
         orchestrator-state.json > tmp.json
    mv tmp.json orchestrator-state.json
    
    # Update cascade status
    jq '(.stale_integration_tracking.staleness_cascade[] | 
         select(.cascade_sequence[].integration == "'$integration_id'")) |= 
         . + {"cascade_status": "in_progress"}' \
         orchestrator-state.json > tmp.json
    mv tmp.json orchestrator-state.json
}
```

## Integration with Rules

### R327 Enforcement
The tracking mechanism directly supports R327 by:
1. Recording when integrations become stale
2. Tracking the cascade requirement
3. Preventing use of stale integrations
4. Documenting recreation completion

### R328 Validation
Supports R328 by:
1. Recording freshness check timestamps
2. Tracking validation frequency
3. Maintaining validation history
4. Identifying stale components

## State Machine Integration

### Required State Transitions
When staleness is detected:
1. Current state saves staleness info
2. Transition to appropriate recreation state:
   - Wave stale → INTEGRATION
   - Phase stale → PHASE_INTEGRATION
   - Project stale → PROJECT_INTEGRATION
3. Track cascade through states
4. Validate freshness before proceeding

### State File Updates
Each state that detects staleness MUST:
1. Update `stale_integration_tracking`
2. Set `is_stale` flags
3. Record triggering fixes
4. Plan cascade recreations

## Reporting

### Staleness Report Generation

```bash
generate_staleness_report() {
    cat > STALENESS-REPORT.md << EOF
# STALE INTEGRATION REPORT
Generated: $(date -Iseconds)

## Currently Stale Integrations
$(jq -r '.stale_integration_tracking.stale_integrations[] | 
        select(.recreation_completed == false) | 
        "- \(.integration_type): \(.integration_id)\n  Stale since: \(.became_stale_at)\n  Reason: \(.stale_reason)"' \
        orchestrator-state.json)

## Pending Cascades
$(jq -r '.stale_integration_tracking.staleness_cascade[] | 
        select(.cascade_status != "completed") | 
        "- Trigger: \(.trigger.branch) at \(.trigger.timestamp)\n  Status: \(.cascade_status)"' \
        orchestrator-state.json)

## Fixes Pending Integration
$(jq -r '.stale_integration_tracking.fix_tracking.fixes_pending_integration[] | 
        "- Fix \(.fix_id) in \(.branch)\n  Pending since: \(.pending_since)\n  Priority: \(.priority)"' \
        orchestrator-state.json)
EOF
}
```

## Benefits

1. **Complete Traceability**: Know exactly which fixes made integrations stale
2. **Cascade Visibility**: See the full impact of any fix
3. **Prevention**: Detect staleness before attempting to use integration
4. **Recovery**: Clear path to recreation with tracking
5. **Audit Trail**: Full history of validations and recreations
6. **Compliance**: Automatic R327/R328 enforcement

## Implementation Checklist

- [ ] Add tracking structure to orchestrator-state.json
- [ ] Update state machine to check staleness
- [ ] Implement detection functions
- [ ] Add cascade tracking logic
- [ ] Create validation scripts
- [ ] Update orchestrator to use tracking
- [ ] Add reporting capabilities
- [ ] Document in state transition guides