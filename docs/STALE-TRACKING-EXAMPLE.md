# STALE INTEGRATION TRACKING - PRACTICAL EXAMPLE

## Scenario
Project "idpbuilder-oci-mgmt" with the following structure:
- Phase 1, Wave 2 completed with 3 efforts
- Integrations created at 03:24 and 17:53
- Fixes applied to efforts at 21:00
- Need to track staleness and cascade

## Initial State

```json
{
  "current_phase": 1,
  "current_wave": 2,
  "current_state": "MONITOR_FIXES",
  
  "current_wave_integration": {
    "workspace": "/efforts/phase1/wave2/wave-integration",
    "branch": "phase1-wave2-integration",
    "created_at": "2025-08-25T17:53:00Z",
    "is_stale": false,
    "merged_efforts": ["effort1", "effort2", "effort3"]
  },
  
  "stale_integration_tracking": {
    "stale_integrations": [],
    "staleness_cascade": [],
    "fix_tracking": {
      "fixes_applied": [],
      "fixes_pending_integration": []
    }
  }
}
```

## Event 1: Fix Applied to Effort

SW Engineer applies fix to effort3 at 21:00:

```json
{
  "stale_integration_tracking": {
    "fix_tracking": {
      "fixes_applied": [
        {
          "fix_id": "FIX-001",
          "commit": "abc123def",
          "branch": "phase1/wave2/effort3",
          "effort_name": "stack-manager-interface",
          "applied_at": "2025-08-25T21:00:00Z",
          "type": "build_fix",
          "description": "Fixed missing import causing build failure",
          "integrated_into": {
            "wave": false,
            "phase": false,
            "project": false
          },
          "made_stale": []
        }
      ]
    }
  }
}
```

## Event 2: Staleness Detection

Orchestrator runs freshness check and detects staleness:

```json
{
  "current_wave_integration": {
    "created_at": "2025-08-25T17:53:00Z",
    "is_stale": true,
    "staleness_reason": "Fix FIX-001 applied after integration creation",
    "stale_since": "2025-08-25T21:00:00Z",
    "stale_due_to_fixes": ["FIX-001"],
    "last_freshness_check": "2025-08-25T21:05:00Z"
  },
  
  "stale_integration_tracking": {
    "stale_integrations": [
      {
        "integration_type": "wave",
        "integration_id": "phase1-wave2-integration",
        "became_stale_at": "2025-08-25T21:00:00Z",
        "stale_reason": "Fix FIX-001 applied to effort3 after integration",
        "triggering_fixes": [
          {
            "fix_id": "FIX-001",
            "commit": "abc123def",
            "branch": "phase1/wave2/effort3",
            "applied_at": "2025-08-25T21:00:00Z",
            "description": "Fixed missing import causing build failure"
          }
        ],
        "affected_downstream": ["phase1-integration", "project-integration"],
        "recreation_required": true,
        "recreation_completed": false
      }
    ]
  }
}
```

## Event 3: Cascade Recording

System records the required cascade:

```json
{
  "stale_integration_tracking": {
    "staleness_cascade": [
      {
        "trigger": {
          "type": "effort_fix",
          "branch": "phase1/wave2/effort3",
          "commit": "abc123def",
          "timestamp": "2025-08-25T21:00:00Z"
        },
        "cascade_sequence": [
          {
            "level": "wave",
            "integration": "phase1-wave2-integration",
            "became_stale": "2025-08-25T21:00:00Z",
            "must_recreate": true,
            "recreation_status": "pending"
          },
          {
            "level": "phase",
            "integration": "phase1-integration",
            "became_stale": "2025-08-25T21:00:00Z",
            "must_recreate": true,
            "recreation_status": "pending"
          },
          {
            "level": "project",
            "integration": "project-integration",
            "became_stale": "2025-08-25T21:00:00Z",
            "must_recreate": true,
            "recreation_status": "pending"
          }
        ],
        "cascade_status": "pending",
        "cascade_started_at": null,
        "cascade_completed_at": null
      }
    ]
  }
}
```

## Event 4: Wave Integration Recreation

Orchestrator recreates wave integration:

```json
{
  "current_wave_integration": {
    "workspace": "/efforts/phase1/wave2/wave-integration-fresh",
    "branch": "phase1-wave2-integration",
    "created_at": "2025-08-25T21:30:00Z",
    "is_stale": false,
    "staleness_reason": null,
    "stale_since": null,
    "stale_due_to_fixes": [],
    "merged_efforts": ["effort1", "effort2", "effort3-with-fix"]
  },
  
  "stale_integration_tracking": {
    "stale_integrations": [
      {
        "integration_id": "phase1-wave2-integration",
        "recreation_completed": true,
        "recreation_at": "2025-08-25T21:30:00Z"
      }
    ],
    "staleness_cascade": [
      {
        "cascade_status": "in_progress",
        "cascade_started_at": "2025-08-25T21:30:00Z",
        "cascade_sequence": [
          {
            "level": "wave",
            "recreation_status": "completed",
            "recreated_at": "2025-08-25T21:30:00Z"
          },
          {
            "level": "phase",
            "recreation_status": "pending"
          }
        ]
      }
    ]
  }
}
```

## Event 5: Multiple Fixes Across Waves

Multiple fixes in different waves:

```json
{
  "stale_integration_tracking": {
    "fix_tracking": {
      "fixes_applied": [
        {
          "fix_id": "FIX-001",
          "branch": "phase1/wave2/effort3",
          "integrated_into": {"wave": true, "phase": false, "project": false}
        },
        {
          "fix_id": "FIX-002",
          "commit": "def456ghi",
          "branch": "phase1/wave1/effort1",
          "applied_at": "2025-08-25T22:00:00Z",
          "type": "test_fix",
          "description": "Fixed failing unit test",
          "made_stale": ["phase1-wave1-integration", "phase1-integration"]
        },
        {
          "fix_id": "FIX-003",
          "commit": "ghi789jkl",
          "branch": "phase1/wave3/effort2",
          "applied_at": "2025-08-25T22:15:00Z",
          "type": "review_fix",
          "description": "Addressed review comments",
          "made_stale": ["phase1-wave3-integration", "phase1-integration"]
        }
      ]
    },
    "staleness_cascade": [
      {
        "trigger": {
          "type": "multiple_fixes",
          "branches": ["phase1/wave1/effort1", "phase1/wave3/effort2"],
          "timestamp": "2025-08-25T22:15:00Z"
        },
        "cascade_sequence": [
          {
            "level": "wave",
            "integration": "phase1-wave1-integration",
            "must_recreate": true
          },
          {
            "level": "wave",
            "integration": "phase1-wave3-integration",
            "must_recreate": true
          },
          {
            "level": "phase",
            "integration": "phase1-integration",
            "must_recreate": true,
            "depends_on": ["wave1", "wave2", "wave3"]
          }
        ]
      }
    ]
  }
}
```

## Queries and Reports

### Query: What made phase1-integration stale?

```bash
jq -r '.stale_integration_tracking.stale_integrations[] | 
       select(.integration_id == "phase1-integration") | 
       .triggering_fixes[] | 
       "Fix \(.fix_id): \(.description) in \(.branch)"' \
       orchestrator-state.json
```

Output:
```
Fix FIX-002: Fixed failing unit test in phase1/wave1/effort1
Fix FIX-003: Addressed review comments in phase1/wave3/effort2
```

### Query: What integrations need recreation?

```bash
jq -r '.stale_integration_tracking.stale_integrations[] | 
       select(.recreation_completed == false) | 
       "\(.integration_type): \(.integration_id)"' \
       orchestrator-state.json
```

Output:
```
phase: phase1-integration
project: project-integration
```

### Query: Which efforts received fixes?

```bash
jq -r '.stale_integration_tracking.fix_tracking.fixes_applied[] | 
       "\(.effort_name): \(.description) (Fix \(.fix_id))"' \
       orchestrator-state.json
```

Output:
```
stack-manager-interface: Fixed missing import causing build failure (Fix FIX-001)
oci-types: Fixed failing unit test (Fix FIX-002)
registry-auth-types: Addressed review comments (Fix FIX-003)
```

## Automated Staleness Check

```bash
#!/bin/bash
# check-staleness.sh

check_and_report_staleness() {
    echo "=== STALENESS CHECK ==="
    echo "Time: $(date -Iseconds)"
    echo ""
    
    # Check each integration level
    for level in wave phase project; do
        integration_created=$(jq -r ".current_${level}_integration.created_at // \"none\"" orchestrator-state.json)
        
        if [[ "$integration_created" != "none" ]]; then
            echo "Checking $level integration (created: $integration_created)..."
            
            # Find fixes after integration creation
            fixes_after=$(jq -r --arg created "$integration_created" '
                .stale_integration_tracking.fix_tracking.fixes_applied[] | 
                select(.applied_at > $created) | 
                .fix_id' orchestrator-state.json)
            
            if [[ -n "$fixes_after" ]]; then
                echo "  ⚠️ STALE! Fixes applied after creation: $fixes_after"
                
                # Mark as stale
                jq --arg level "$level" --arg fixes "$fixes_after" '
                    .["current_" + $level + "_integration"].is_stale = true |
                    .["current_" + $level + "_integration"].stale_due_to_fixes = ($fixes | split(" "))
                ' orchestrator-state.json > tmp.json
                mv tmp.json orchestrator-state.json
            else
                echo "  ✅ Fresh - no fixes after creation"
            fi
        fi
    done
    
    # Report cascade requirements
    echo ""
    echo "=== CASCADE REQUIREMENTS ==="
    jq -r '.stale_integration_tracking.staleness_cascade[] | 
           select(.cascade_status != "completed") | 
           "Pending cascade from \(.trigger.branch) - Status: \(.cascade_status)"' \
           orchestrator-state.json
}

check_and_report_staleness
```

## Benefits Demonstrated

1. **Clear Visibility**: Can see exactly which fix (FIX-001) made the integration stale
2. **Effort Tracking**: Know that "stack-manager-interface" received the fix
3. **Cascade Management**: Automatic tracking of phase and project staleness
4. **Temporal Tracking**: Timestamps show when staleness occurred
5. **Recovery Path**: Clear indication of what needs recreation
6. **Audit Trail**: Complete history of fixes and their impacts