# 🔴🔴🔴 RULE R504 - Pre-Infrastructure Planning [SUPREME LAW]
**VIOLATION = SYSTEM CORRUPTION** | Source: [R504](RULE-REGISTRY.md#R504)

## PURPOSE
Eliminate ALL runtime naming/pathing decisions by pre-calculating and validating ALL infrastructure during planning phases. Infrastructure creation becomes pure mechanical execution of pre-planned structure. This prevents naming conflicts, path collisions, and ensures consistency across the entire system.

## WHEN THIS RULE APPLIES
- **DURING**: PROJECT_PLANNING, PHASE_PLANNING, WAVE_PLANNING states
- **BEFORE**: ANY infrastructure creation state (CREATE_NEXT_INFRASTRUCTURE, etc.)
- **ENFORCED IN**: All infrastructure creation and validation states
- **TRIGGER**: ANY attempt to create directories, branches, or remotes

## MANDATORY REQUIREMENTS

### 1. PARSE AND EXTRACT ALL EFFORTS
```bash
# Extract from IMPLEMENTATION-PLAN.md:
- ALL phases (Phase 1, Phase 2, etc.)
- ALL waves within each phase
- ALL efforts within each wave
- ALL dependencies between efforts
- ALL estimated sizes/complexities
```

### 2. PRE-CALCULATE COMPLETE INFRASTRUCTURE

**A. EFFORT INFRASTRUCTURE:**
```json
{
  "effort_id": {
    "full_path": "$CLAUDE_PROJECT_DIR/efforts/phase1/wave1/effort-name/",
    "branch_name": "project-prefix/phase1/wave1/effort-name",
    "remote_branch": "origin/project-prefix/phase1/wave1/effort-name",
    "target_remote": "target",
    "planning_remote": "planning",
    "split_pattern": "effort-name--split-NNN",
    "integration_branch": "project-prefix/phase1/wave1/integration",
    "target_repo_url": "https://github.com/user/repo.git"
  }
}
```

**B. INTEGRATE_WAVE_EFFORTS INFRASTRUCTURE (CRITICAL - NEW REQUIREMENT):**
```json
{
  "wave_integrations": {
    "phase1_wave1": {
      "branch_name": "project-prefix/phase1-wave1-integration",
      "base_branch": "main or previous-wave-integration",
      "target_repo_url": "MUST MATCH TARGET REPO",
      "directory": "/absolute/path/for/integration/work"
    }
  },
  "phase_integrations": {
    "phase1": {
      "branch_name": "project-prefix/phase1-integration",
      "base_branch": "last-wave-integration-of-phase",
      "target_repo_url": "MUST MATCH TARGET REPO",
      "directory": "/absolute/path/for/phase/integration"
    }
  }
}
```

### 3. VALIDATION REQUIREMENTS
**MUST validate BEFORE creating any infrastructure:**
- ✅ All paths follow R313 (effort directory structure)
- ✅ All branch names follow R327 (cascade branching)
- ✅ No path conflicts exist
- ✅ No branch name conflicts
- ✅ Remote configurations are valid
- ✅ Integration dependencies are correct
- ✅ Split patterns are pre-defined
- ✅ Each effort has target_repo_url matching parent
- ✅ Target repository URL is valid and accessible
- ✅ **ALL integration branches have target_repo_url**
- ✅ **Integration target_repo_url matches target-repo-config.yaml**
- ✅ **Wave integrations pre-planned during wave planning**
- ✅ **Phase integrations pre-planned during phase planning**
- ✅ **NO integration infrastructure created without pre-planning**

### 4. STATE FILE POPULATION
**orchestrator-state-v3.json MUST contain:**
```json
{
  "pre_planned_infrastructure": {
    "validated": true,
    "validation_timestamp": "ISO-8601",
    "project_prefix": "determined-from-plan",
    "target_repo_url": "https://github.com/user/repo.git",
    "efforts": {
      "phase1_wave1_effort1": {
        "full_path": "/absolute/path/to/effort/",
        "branch_name": "fully-qualified-branch-name",
        "remote_branch": "origin/fully-qualified-branch-name",
        "target_repo_url": "https://github.com/user/repo.git",
        "created": false,
        "validated": true
      }
    },
    "integrations": {
      "wave_integrations": {
        "phase1_wave1": {
          "phase": "phase1",
          "wave": "wave1",
          "branch_name": "project-prefix/phase1-wave1-integration",
          "remote_branch": "project-prefix/phase1-wave1-integration",
          "base_branch": "main",
          "target_repo_url": "https://github.com/user/repo.git",
          "directory": "/absolute/path/efforts/phase1/wave1/integration-workspace",
          "component_efforts": ["E1.1.1", "E1.1.2"],
          "created": false,
          "validated": false
        }
      },
      "phase_integrations": {
        "phase1": {
          "phase": "phase1",
          "branch_name": "project-prefix/phase1-integration",
          "remote_branch": "project-prefix/phase1-integration",
          "base_branch": "project-prefix/phase1-wave2-integration",
          "target_repo_url": "https://github.com/user/repo.git",
          "directory": "/absolute/path/efforts/phase1/integration",
          "component_waves": ["phase1_wave1", "phase1_wave2"],
          "created": false,
          "validated": false
        }
      }
    }
  }
}
```

### 5. INFRASTRUCTURE CREATION BECOMES MECHANICAL
**CREATE_NEXT_INFRASTRUCTURE state MUST:**
- ❌ NEVER make naming decisions
- ❌ NEVER calculate paths at runtime
- ❌ NEVER determine branch names dynamically
- ✅ ONLY read from pre_planned_infrastructure
- ✅ ONLY execute pre-validated structure
- ✅ ONLY mark items as "created": true

## VALIDATION CHECKPOINTS

### During Planning Phase:
```bash
# 1. Parse IMPLEMENTATION-PLAN.md
parse_implementation_plan() {
  local plan_file="$CLAUDE_PROJECT_DIR/PROJECT-IMPLEMENTATION-PLAN.md"
  # Extract all phases/waves/efforts
  # Return structured data
}

# 2. Pre-calculate all infrastructure
calculate_infrastructure() {
  local project_prefix="$1"
  # Generate all paths/branches
  # Apply naming rules
  # Return complete structure
}

# 3. Validate against rules
validate_infrastructure() {
  local infrastructure="$1"
  # Check R313 compliance (paths)
  # Check R327 compliance (branches)
  # Check for conflicts
  # Return validation result
}
```

### During Execution:
```bash
# Just read and execute
create_effort_infrastructure() {
  local effort_id="$1"
  local config=$(yq ".pre_planned_infrastructure.efforts.$effort_id" orchestrator-state-v3.json)

  # No decisions, just execution:
  mkdir -p "$(echo $config | yq '.full_path')"
  cd "$(echo $config | yq '.full_path')"
  git init
  git checkout -b "$(echo $config | yq '.branch_name')"

  # Mark as created
  yq -i ".pre_planned_infrastructure.efforts.$effort_id.created = true" orchestrator-state-v3.json
}
```

## FAILURE CONDITIONS (AUTOMATIC FAIL + SYSTEM HALT)
- 🔴 Creating ANY infrastructure without pre_planned_infrastructure = IMMEDIATE FAIL
- 🔴 Runtime path/branch calculation in CREATE states = SYSTEM CORRUPTION
- 🔴 Missing pre_planned_infrastructure.validated=true = CANNOT PROCEED
- 🔴 ANY deviation from pre-planned structure = CATASTROPHIC FAILURE
- 🔴 Incomplete or missing pre_planned_infrastructure = HALT FACTORY
- 🔴 Just-in-time naming decisions = -100% GRADE

## PROJECT_DONE CRITERIA
- ✅ ALL infrastructure pre-calculated during planning
- ✅ ALL paths/branches validated before ANY creation
- ✅ orchestrator-state-v3.json is SOLE SOURCE OF TRUTH
- ✅ Infrastructure creation is 100% mechanical
- ✅ Zero runtime naming decisions

## RELATED RULES
- R313: Effort directory structure
- R327: Cascade branching conventions
- R206: State machine validation
- R233: State action requirements

## ENFORCEMENT (MANDATORY AT ALL LEVELS)

### Planning States MUST:
1. Parse all phases/waves/efforts from PROJECT-IMPLEMENTATION-PLAN.md
2. Pre-calculate ALL paths, branches, remotes for entire project
3. Validate against target-repo-config.yaml
4. Populate pre_planned_infrastructure completely
5. Set validated=true only after full validation

### Infrastructure Creation States MUST:
1. **FIRST CHECK**: Verify pre_planned_infrastructure exists and validated=true
2. **FAIL IMMEDIATELY** if pre_planned_infrastructure missing or invalid
3. Read ALL data from pre_planned_infrastructure (NO calculations)
4. Execute mechanically from pre-planned specifications
5. Mark items as created=true after successful creation
6. **NEVER** make naming/pathing decisions at runtime

### Validation States MUST:
1. Compare actual infrastructure against pre_planned_infrastructure
2. Verify branch tracking matches pre-planned remote configurations
3. Check directory paths match pre-planned full_path values
4. Validate against target-repo-config.yaml for correct repository
5. FAIL if ANY deviation from pre-planned structure

**Penalty for violation:** -100% grade, immediate SYSTEM HALT, transition to ERROR_RECOVERY
## State Manager Coordination (SF 3.0)

State Manager enforces planning-before-infrastructure through transition guards:
- **PLANNING state**: Load plan into `orchestrator-state-v3.json`
- **Shutdown consultation**: Validate plan exists and is complete
- **SETUP_WAVE_INFRASTRUCTURE guard**: Requires completed plan
- **Transition rejection**: Cannot create infrastructure without plan

This prevents infrastructure creation before planning is done.

See: State machine SETUP_WAVE_INFRASTRUCTURE.requires.conditions[], R502
