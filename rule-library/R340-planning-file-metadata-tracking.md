# 🔴🔴🔴 RULE R340: Planning File Metadata Tracking (BLOCKING) 🔴🔴🔴

## Rule Identifier
**Rule ID**: R340
**Category**: State Management / Metadata Tracking
**Criticality**: BLOCKING - Planning file discovery failures cause integration delays
**Priority**: ABSOLUTE

## Description

All planning files created by any agent MUST be tracked with complete metadata in orchestrator-state.json to enable instant discovery and prevent agents from searching through directories.

## Rationale

Without centralized planning file tracking:
- Integration agents waste time searching for merge plans
- SW Engineers cannot find their implementation plans
- Orchestrator loses track of which plans exist
- Splits and fixes become untraceable
- Recovery after context loss is impossible

## Requirements

### 1. Planning File Types to Track

The following planning file types MUST be tracked with metadata:

```json
"planning_files": {
  "effort_plans": {},      // EFFORT-IMPLEMENTATION-PLAN-*.md files
  "split_plans": {},       // SPLIT-PLAN-*.md files
  "fix_plans": {},         // FIX-PLAN-*.md files
  "merge_plans": {
    "wave": {},            // WAVE-MERGE-PLAN.md files
    "phase": {},           // PHASE-MERGE-PLAN.md files
    "project": {}          // PROJECT-MERGE-PLAN.md files
  },
  "architecture_plans": {
    "wave": {},            // WAVE-ARCHITECTURE-PLAN.md files
    "phase": {},           // PHASE-ARCHITECTURE-PLAN.md files
  },
  "review_reports": {}     // CODE-REVIEW-REPORT-*.md files
}
```

### 2. Required Metadata Fields

Each planning file entry MUST include:

```json
{
  "file_path": "/absolute/path/to/plan.md",
  "created_at": "2025-01-20T10:30:00Z",
  "created_by": "agent-name",
  "target_branch": "branch-name",
  "effort_name": "effort-identifier",
  "phase": 1,
  "wave": 2,
  "status": "active|deprecated|completed",
  "replaced_by": null  // Path to replacement if deprecated
}
```

### 3. Agent Responsibilities

#### Code Reviewer Agent
After creating ANY planning file, MUST report:
```markdown
## 📋 PLANNING FILE CREATED

**Type**: effort_plan|split_plan|fix_plan|merge_plan
**Path**: /absolute/path/to/plan.md
**Effort**: effort-name
**Phase**: X
**Wave**: Y
**Target Branch**: branch-name
**Created At**: ISO-8601-timestamp

ORCHESTRATOR: Please update planning_files.{type}.{identifier} in state file
```

#### Orchestrator Agent
Upon receiving planning file report, MUST:
1. Parse the metadata from the report
2. Update orchestrator-state.json immediately
3. Confirm update:
```bash
# Update state file
yq eval '.planning_files.effort_plans["effort-name"] = {
  "file_path": "/path/to/plan.md",
  "created_at": "timestamp",
  "created_by": "code-reviewer",
  "target_branch": "branch",
  "phase": 1,
  "wave": 2,
  "status": "active"
}' -i orchestrator-state.json

# Commit immediately
git add orchestrator-state.json
git commit -m "state: track planning file for effort-name"
git push
```

#### SW Engineer Agent
To find implementation plan, MUST:
```bash
# Read plan location from state
PLAN_PATH=$(yq '.planning_files.effort_plans["effort-name"].file_path' orchestrator-state.json)

# Verify plan exists
if [ ! -f "$PLAN_PATH" ]; then
    echo "❌ CRITICAL: Plan not found at tracked location!"
    echo "Expected: $PLAN_PATH"
    exit 1
fi

# Read the plan
cat "$PLAN_PATH"
```

#### Integration Agent
To find merge plan, MUST:
```bash
# For wave merge plan
MERGE_PLAN=$(yq '.planning_files.merge_plans.wave["phase1_wave1"].file_path' orchestrator-state.json)

# For phase merge plan
MERGE_PLAN=$(yq '.planning_files.merge_plans.phase["phase1"].file_path' orchestrator-state.json)
```

### 4. Standard Identifiers

Planning files use these identifiers in the state file:

- **Effort Plans**: Use effort name (e.g., "buildah-builder-interface")
- **Split Plans**: Use "effort-name-split-001", "effort-name-split-002"
- **Fix Plans**: Use "effort-name-fix-001", "effort-name-fix-002"
- **Wave Merge Plans**: Use "phase1_wave1", "phase2_wave3"
- **Phase Merge Plans**: Use "phase1", "phase2"
- **Project Merge Plan**: Use "project"

### 5. Deprecation Handling

When a plan is replaced:
```json
{
  "old-plan": {
    "file_path": "/path/to/old-plan.md",
    "status": "deprecated",
    "replaced_by": "/path/to/new-plan.md",
    "deprecated_at": "2025-01-20T11:00:00Z"
  },
  "new-plan": {
    "file_path": "/path/to/new-plan.md",
    "status": "active",
    "replaces": "/path/to/old-plan.md"
  }
}
```

### 6. Discovery Functions

Agents MUST use these patterns to discover plans:

```bash
# Function to find any plan
find_plan() {
    local plan_type="$1"  # effort_plans, split_plans, etc.
    local identifier="$2"
    
    local path=$(yq ".planning_files.${plan_type}[\"${identifier}\"].file_path" orchestrator-state.json)
    
    if [ "$path" = "null" ] || [ -z "$path" ]; then
        echo "❌ No plan tracked for ${plan_type}/${identifier}"
        return 1
    fi
    
    if [ ! -f "$path" ]; then
        echo "❌ Tracked plan not found at: $path"
        return 1
    fi
    
    echo "$path"
}

# Usage examples
EFFORT_PLAN=$(find_plan "effort_plans" "buildah-builder-interface")
SPLIT_PLAN=$(find_plan "split_plans" "oci-types-split-001")
WAVE_MERGE=$(find_plan "merge_plans.wave" "phase1_wave1")
```

## Validation

### Pre-Integration Validation
```bash
# Verify all efforts have tracked plans
for effort in $(yq '.efforts_in_progress[].name' orchestrator-state.json); do
    plan_path=$(yq ".planning_files.effort_plans[\"$effort\"].file_path" orchestrator-state.json)
    if [ "$plan_path" = "null" ]; then
        echo "❌ VIOLATION: No plan tracked for effort: $effort"
        exit 1
    fi
done
```

### State File Validation
```bash
# Check for orphaned plan references
for plan_path in $(yq '.planning_files.**.file_path' orchestrator-state.json); do
    if [ ! -f "$plan_path" ]; then
        echo "⚠️ WARNING: Tracked plan missing: $plan_path"
    fi
done
```

## Common Violations

1. ❌ **Creating plans without reporting metadata**
   - Code Reviewer creates plan but doesn't report to Orchestrator
   - Fix: Add reporting block after every plan creation

2. ❌ **Searching for plans instead of reading from state**
   - Using `find` or `ls` to locate plans
   - Fix: Always read location from orchestrator-state.json

3. ❌ **Not updating state when plans move/change**
   - Plan relocated but state not updated
   - Fix: Update state immediately after any plan modification

4. ❌ **Missing required metadata fields**
   - Incomplete metadata in state file
   - Fix: Use complete metadata structure

## Enforcement

- **Grading Impact**: -20% for each untracked planning file
- **Integration Block**: Cannot proceed if merge plan location unknown
- **Recovery Block**: Cannot resume work without plan location

## Related Rules

- R251: Repository Separation (plans in .software-factory/)
- R324: State file update requirements
- R337: Base branch single source of truth

## Implementation Notes

1. Planning files are stored in `.software-factory/` within effort directories
2. Paths in state file must be absolute paths
3. State updates must be committed immediately (R288)
4. All agents must validate plan existence before use

## Example State File Section

```json
{
  "planning_files": {
    "effort_plans": {
      "buildah-builder-interface": {
        "file_path": "/efforts/phase1/wave2/buildah-builder-interface/.software-factory/phase1/wave2/buildah-builder-interface/IMPLEMENTATION-PLAN-20250120-100000.md",
        "created_at": "2025-01-20T10:00:00Z",
        "created_by": "code-reviewer",
        "target_branch": "phase1/wave2/buildah-builder-interface",
        "phase": 1,
        "wave": 2,
        "status": "active"
      }
    },
    "merge_plans": {
      "wave": {
        "phase1_wave1": {
          "file_path": "/efforts/phase1/wave1/integration-workspace/WAVE-MERGE-PLAN.md",
          "created_at": "2025-01-20T12:00:00Z",
          "created_by": "code-reviewer",
          "target_branch": "phase1-wave1-integration",
          "phase": 1,
          "wave": 1,
          "status": "completed"
        }
      }
    }
  }
}
```

---

**Remember**: Every planning file MUST be tracked. No exceptions. Agents MUST read locations from state, never search for plans.