# 🚨🚨🚨 RULE R328 - Integration Freshness Validation [BLOCKING]

**File**: `rule-library/R328-integration-freshness-validation.md`
**Criticality**: BLOCKING - Stale integrations cause failed merges and lost fixes
**Created**: 2025-09-09
**Purpose**: Ensure integration branches are created from the latest effort/phase branches

## The Problem This Solves

The orchestrator has been merging outdated integration branches that were created BEFORE fixes were applied to effort branches. This causes:
- Lost fixes and improvements
- Failed builds and tests
- Confusion about which branches contain the latest code
- Wasted time debugging "already fixed" issues

## The Solution

**ALWAYS verify integration branch freshness BEFORE merging:**

1. **Check Creation Timestamps**: Compare when integration branch was created vs when efforts were last updated
2. **Detect Staleness**: If any effort has commits newer than integration branch creation, it's STALE
3. **Recreate if Stale**: Don't use stale integration branches - recreate from fresh efforts
4. **Track in State File**: Maintain metadata about branch creation times and updates

## Implementation Requirements

### 1. State File Tracking

The orchestrator-state.json MUST track:
```json
{
  "efforts_completed": [{
    "last_commit": "abc123",
    "last_updated_at": "2025-08-25T09:00:00Z",
    "has_fixes_applied": true,
    "fixes_applied_at": "2025-08-25T09:30:00Z"
  }],
  "current_wave_integration": {
    "created_at": "2025-08-25T08:00:00Z",
    "efforts_last_updated_at": "2025-08-25T09:30:00Z",
    "is_stale": true
  }
}
```

### 2. Freshness Check Function

```bash
check_integration_freshness() {
    local integration_type="$1"  # wave, phase, or project
    
    # Get integration creation time
    local integration_created=$(jq -r ".current_${integration_type}_integration.created_at" orchestrator-state.json)
    
    # Check each relevant effort/component
    local is_stale=false
    local stale_components=""
    
    # Logic to check timestamps...
    
    if [[ "$is_stale" == "true" ]]; then
        echo "🔴 CRITICAL: ${integration_type} integration is STALE!"
        echo "Components with newer updates: $stale_components"
        return 1
    fi
    
    return 0
}
```

### 3. Enforcement Points

**MUST check freshness at these states:**
- INTEGRATION_TESTING (before merging efforts into wave)
- WAVE_INTEGRATION (before merging waves into phase)
- PHASE_INTEGRATION (before merging phases into project)
- PROJECT_INTEGRATION (before final project merge)

### 4. Recovery Protocol

When staleness is detected:
1. **STOP** - Do not proceed with stale branch
2. **DOCUMENT** - Create STALENESS-REPORT.md
3. **DECIDE**:
   - Option A: Recreate integration from fresh branches
   - Option B: Apply fixes to existing integration (if minor)
   - Option C: Spawn SW Engineers to handle complex updates
4. **UPDATE** - Mark old integration as deprecated, create new one

## Repository Context Clarification

### Software Factory Repository
- **Location**: `/home/vscode/software-factory-template/`
- **Contains**: SF code, rules, state files, agent configs
- **Branches**: main, software-factory-2.0
- **NEVER contains**: Effort branches, integration branches

### Target Repository Clones
- **Location Pattern**: `efforts/*/*/`
- **Contains**: Actual project implementation code
- **Branches**: Effort branches (e.g., `phase1/wave1/effort-name`)
- **Purpose**: Where SW Engineers write code

### Integration Workspaces
- **Location Pattern**: `efforts/*/integration-workspace/`
- **Contains**: Integration branches and merge operations
- **Branches**: Integration branches (e.g., `wave1-integration`)
- **Purpose**: Where merging happens

## Common Mistakes to Avoid

1. **Looking for effort branches in SF repo** - They don't exist there!
2. **Assuming integration branches are always fresh** - Check timestamps!
3. **Merging without checking merge plan** - Plans list correct branch locations
4. **Not tracking workspace locations** - State file must track where things are

## Grading Impact

- Using stale integration branch: -30% penalty
- Not checking freshness before merge: -20% penalty
- Losing fixes due to stale branch: -50% penalty
- Not documenting staleness issues: -15% penalty

## Related Rules

- R283: Project Integration Protocol
- R269: Merge Plan Requirement
- R321: Immediate Backport During Integration
- R104: Target Repository Integration

## Summary

**NEVER trust an integration branch without checking its freshness!**

Integration branches become stale when:
- Fixes are applied to effort branches after integration creation
- Backports are made to source branches
- Emergency patches are applied
- Split implementations receive updates

Always verify, always track, always use fresh branches!