# 🚨🚨🚨 BLOCKING RULE R503: Mandatory Planning Repository Tracking

## Rule Statement
ALL planning documents in the `$CLAUDE_PROJECT_DIR/planning/` directory MUST be tracked in orchestrator-state-v3.json under `planning_files.planning_repo_files`. The orchestrator MUST update tracking immediately when plans are validated or created. Missing tracking causes integration failures and lost plans.

## Criticality: 🚨🚨🚨 BLOCKING
Without tracking planning repository documents, agents cannot find required plans, integration gets blocked, and the entire project loses its architectural guidance.

## Tracking Structure

### Required State File Section
```json
{
  "planning_files": {
    "planning_repo_files": {
      "project_plans": {
        "PROJECT-IMPLEMENTATION-PLAN": {
          "file_path": "$CLAUDE_PROJECT_DIR/planning/project/PROJECT-IMPLEMENTATION-PLAN.md",
          "exists": true|false,
          "created_at": "ISO-8601",
          "created_by": "agent-name",
          "status": "required|optional|completed",
          "last_validated": "ISO-8601"
        }
      },
      "phase_plans": {
        "phase1": {
          "architecture": {
            "file_path": "$CLAUDE_PROJECT_DIR/planning/phase1/PHASE-1-ARCHITECTURE-PLAN.md",
            "exists": true|false,
            "created_at": "ISO-8601",
            "created_by": "architect|code-reviewer",
            "status": "required|completed",
            "last_validated": "ISO-8601"
          },
          "implementation": {
            "file_path": "$CLAUDE_PROJECT_DIR/planning/phase1/PHASE-1-IMPLEMENTATION-PLAN.md",
            "exists": true|false,
            "created_at": "ISO-8601",
            "created_by": "code-reviewer",
            "status": "required|completed",
            "last_validated": "ISO-8601"
          }
        }
      },
      "wave_plans": {
        "phase1_wave1": {
          "architecture": { ... },
          "implementation": { ... },
          "test": { ... }
        }
      }
    }
  }
}
```

## Mandatory Update Points

### 1. When Validating Plans (R502)
```bash
# After successfully validating plans
validate_and_track_plans() {
    local PHASE="$1"
    local WAVE="$2"

    # Validate phase plans
    if [ -f "$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN.md" ]; then
        yq eval ".planning_files.planning_repo_files.phase_plans[\"phase${PHASE}\"].architecture.exists = true" -i orchestrator-state-v3.json
        yq eval ".planning_files.planning_repo_files.phase_plans[\"phase${PHASE}\"].architecture.last_validated = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json
    fi

    # Validate wave plans
    if [ -f "$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md" ]; then
        yq eval ".planning_files.planning_repo_files.wave_plans[\"phase${PHASE}_wave${WAVE}\"].implementation.exists = true" -i orchestrator-state-v3.json
        yq eval ".planning_files.planning_repo_files.wave_plans[\"phase${PHASE}_wave${WAVE}\"].implementation.last_validated = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json
    fi

    git add orchestrator-state-v3.json
    git commit -m "state: update planning_repo_files tracking for phase${PHASE} wave${WAVE}"
    git push
}
```

### 2. When Creating Plans
```bash
# Architect creates phase architecture plan
create_phase_architecture() {
    local PHASE="$1"
    local PLAN_PATH="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"

    # Create the plan
    cat > "$PLAN_PATH" << 'EOF'
    # Phase Architecture Plan content
    EOF

    # MANDATORY: Update tracking
    yq eval ".planning_files.planning_repo_files.phase_plans[\"phase${PHASE}\"].architecture = {
        \"file_path\": \"$PLAN_PATH\",
        \"exists\": true,
        \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
        \"created_by\": \"architect\",
        \"status\": \"completed\",
        \"last_validated\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }" -i orchestrator-state-v3.json

    # Report to orchestrator
    echo "📋 PLANNING FILE CREATED"
    echo "Repository: planning"
    echo "Type: phase_architecture"
    echo "Path: $PLAN_PATH"
    echo "Phase: $PHASE"
    echo "ORCHESTRATOR: Please verify planning_repo_files.phase_plans[\"phase${PHASE}\"].architecture updated"
}
```

### 3. State Transitions Requiring Tracking

#### CREATE_NEXT_INFRASTRUCTURE
MUST update tracking after validating plans:
```bash
# R503 MANDATORY
update_planning_tracking "phase" "$CURRENT_PHASE"
update_planning_tracking "wave" "$CURRENT_PHASE" "$CURRENT_WAVE"
```

#### WAVE_START
MUST verify tracking exists:
```bash
# Check tracking is current
LAST_VALIDATED=$(yq ".planning_files.planning_repo_files.wave_plans[\"phase${PHASE}_wave${WAVE}\"].implementation.last_validated" orchestrator-state-v3.json)
if [ -z "$LAST_VALIDATED" ] || [ "$LAST_VALIDATED" = "null" ]; then
    echo "❌ CRITICAL: Wave plans not tracked in planning_repo_files!"
    exit 1
fi
```

#### START_PHASE_ITERATION (implicit in first wave)
MUST ensure phase tracking:
```bash
# Verify phase plans are tracked
PHASE_ARCH_TRACKED=$(yq ".planning_files.planning_repo_files.phase_plans[\"phase${PHASE}\"].architecture.exists" orchestrator-state-v3.json)
if [ "$PHASE_ARCH_TRACKED" != "true" ]; then
    echo "❌ CRITICAL: Phase architecture plan not tracked!"
    exit 1
fi
```

## Agent Responsibilities

### Orchestrator
- **MUST** maintain planning_repo_files section
- **MUST** update tracking when validating plans
- **MUST** verify tracking before state transitions
- **MUST** commit tracking updates immediately

### Architect
- **MUST** report plan creation to orchestrator
- **MUST** include full path and metadata
- **MUST** request tracking update confirmation

### Code Reviewer
- **MUST** report implementation plan creation
- **MUST** distinguish between planning repo and effort repo files
- **MUST** verify orchestrator updated tracking

### SW Engineer
- **MUST** check planning_repo_files for wave plans
- **MUST** report if tracked plan doesn't exist
- **MUST** use tracked paths, not search for plans

## Verification Commands

### Check if tracking exists:
```bash
# Verify structure exists
yq '.planning_files.planning_repo_files' orchestrator-state-v3.json

# Check specific phase tracking
yq '.planning_files.planning_repo_files.phase_plans["phase1"]' orchestrator-state-v3.json

# Check wave tracking
yq '.planning_files.planning_repo_files.wave_plans["phase1_wave1"]' orchestrator-state-v3.json
```

### Validate tracking accuracy:
```bash
# Compare tracked files with actual files
for plan in $(yq '.planning_files.planning_repo_files.wave_plans[].*.file_path' orchestrator-state-v3.json); do
    if [ ! -f "$plan" ]; then
        echo "❌ Tracked plan doesn't exist: $plan"
    fi
done
```

## Common Violations and Fixes

### ❌ VIOLATION: Missing planning_repo_files section
```bash
# WRONG - Old structure only
{
  "planning_files": {
    "effort_plans": { ... }
  }
}

# CORRECT - Include planning_repo_files
{
  "planning_files": {
    "planning_repo_files": { ... },
    "effort_plans": { ... }
  }
}
```

### ❌ VIOLATION: Not updating tracking after validation
```bash
# WRONG - Validate but don't track
if [ -f "$PHASE_PLAN" ]; then
    echo "✅ Plan exists"
fi

# CORRECT - Validate AND track
if [ -f "$PHASE_PLAN" ]; then
    echo "✅ Plan exists"
    yq eval ".planning_files.planning_repo_files.phase_plans[\"phase${PHASE}\"].exists = true" -i orchestrator-state-v3.json
    yq eval ".planning_files.planning_repo_files.phase_plans[\"phase${PHASE}\"].last_validated = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state-v3.json
fi
```

### ❌ VIOLATION: Searching for plans instead of using tracking
```bash
# WRONG - Search filesystem
PLAN=$(find planning/ -name "WAVE-*-IMPLEMENTATION-PLAN.md" | head -1)

# CORRECT - Use tracked location
PLAN=$(yq ".planning_files.planning_repo_files.wave_plans[\"phase${PHASE}_wave${WAVE}\"].implementation.file_path" orchestrator-state-v3.json)
```

## Integration with Other Rules

### Works with:
- **R502**: Validates plans exist before tracking
- **R340**: Extends tracking to planning repository
- **R303**: Uses correct planning/ directory structure
- **R210/R211**: Tracks plans created by architect/code-reviewer

### Blocks:
- State transitions without tracking
- Infrastructure creation without tracked plans
- Agent spawn without plan tracking
- Integration without merge plan tracking

## Grading Impact

### Penalties for Violations:
- Missing planning_repo_files section: **-40%**
- Not updating tracking after validation: **-30%**
- Proceeding without verified tracking: **-35%**
- Agents not using tracked locations: **-25%**
- Stale tracking (>1 hour old): **-20%**

### Required Evidence:
- planning_repo_files section in state file
- Tracking updates in git history
- Validation checks before state transitions
- Agent logs showing tracked path usage

## Recovery Protocol

If tracking is missing or corrupted:
1. **STOP** all work immediately
2. **Scan** planning/ directory for all plans
3. **Rebuild** planning_repo_files section
4. **Validate** each plan exists
5. **Update** timestamps and metadata
6. **Commit** corrected state file
7. **Resume** from validation point

## Summary

This rule ensures that ALL planning documents in the planning/ directory are properly tracked in the state file. Without this tracking, the Software Factory loses its ability to find and use architectural guidance, causing cascading failures throughout the system.

The planning_repo_files section is NOT optional - it's a MANDATORY part of the orchestrator's state management responsibilities.

---
*Rule Type*: State Management / Tracking
*Agents*: All Agents (Primary: Orchestrator)
*Enforcement*: Blocking at validation and state transition points
*Created*: Response to missing planning document tracking in orchestrator-state-v3.json