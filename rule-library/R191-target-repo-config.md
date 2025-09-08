# 🚨🚨🚨 RULE R191 - Target Repository Configuration

**Criticality:** BLOCKING - Wrong repo = Total failure  
**Grading Impact:** -100% for working on wrong repository  
**Enforcement:** IMMEDIATE - Check before ANY code operation
**See Also:** R309 - NEVER Create Efforts in SF Repo

## 🔴🔴🔴 CRITICAL: TWO SEPARATE REPOSITORIES! 🔴🔴🔴

### THE FUNDAMENTAL SEPARATION:
1. **SOFTWARE FACTORY REPO** (Where this file lives)
   - Has: `.claude/`, `rule-library/`, `orchestrator-state.json`
   - Purpose: PLANNING AND ORCHESTRATION ONLY
   - **NEVER IMPLEMENT CODE HERE!**

2. **TARGET REPOSITORY** (What you clone for efforts)
   - Defined in: `target-repo-config.yaml`
   - Cloned to: `efforts/phaseX/waveY/effort-name/`
   - Purpose: ACTUAL PROJECT IMPLEMENTATION
   - **ALL CODE GOES HERE!**

### CRITICAL VALIDATION:
```bash
# The target URL must NEVER be the SF repo itself!
if [[ "$TARGET_REPO_URL" == *"software-factory"* ]]; then
    echo "🔴🔴🔴 FATAL: Target is Software Factory!"
    echo "This creates recursive pollution!"
    exit 191
fi
```

## Rule Statement

ALL agents MUST read and respect `target-repo-config.yaml` which defines the ACTUAL PROJECT repository to work on, NOT the Software Factory instance repository.

## Repository Architecture

```
/workspaces/
├── my-project-sw-factory/          # SOFTWARE FACTORY INSTANCE
│   ├── target-repo-config.yaml     # TARGET REPO CONFIGURATION
│   ├── rule-library/               # Rules (READ ONLY for agents)
│   ├── utilities/                  # Utilities (READ ONLY)
│   ├── .claude/                    # Agent configs (READ ONLY)
│   ├── planning/                   # Plans and docs
│   ├── state/                      # State files (WRITABLE)
│   ├── todos/                      # TODO files (WRITABLE)
│   └── efforts/                    # EFFORT WORKSPACES
│       ├── phase1/
│       │   ├── wave1/
│       │   │   ├── api-types/      # CLONED TARGET REPO
│       │   │   └── controllers/    # CLONED TARGET REPO
│       │   └── wave2/
│       │       └── webhooks/        # CLONED TARGET REPO
│       └── phase2/
│           └── wave1/
│               └── integration/     # CLONED TARGET REPO
```

## Configuration Loading

### Orchestrator MUST Load Config
```bash
# At startup, orchestrator loads target config
load_target_config() {
    local config_file="${SF_ROOT}/target-repo-config.yaml"
    
    if [ ! -f "$config_file" ]; then 
        echo "❌ CRITICAL: target-repo-config.yaml not found!"; 
        echo "Cannot proceed without target repository configuration"; 
        exit 1; 
    fi
    
    # Parse configuration
    # Use text_editor tool with view command to read target-repo-config.yaml:
    # Find target_repository.url, target_repository.base_branch, workspace.efforts_root
    TARGET_REPO_URL="<value from target_repository.url>"
    BASE_BRANCH="<value from target_repository.base_branch>"
    EFFORTS_ROOT="<value from workspace.efforts_root>"
    
    if [ -z "$TARGET_REPO_URL" ]; then 
        echo "❌ CRITICAL: No target repository URL configured!"; 
        exit 1; 
    fi
    
    echo "✅ Target repository: $TARGET_REPO_URL"
    echo "✅ Base branch: $BASE_BRANCH"
    echo "✅ Efforts root: $EFFORTS_ROOT"
    
    export TARGET_REPO_URL
    export BASE_BRANCH
    export EFFORTS_ROOT
}
```

### All Agents MUST Verify Config
```bash
# Every agent checks they have target config
verify_target_config() {
    if [ -z "$TARGET_REPO_URL" ]; then 
        echo "❌ CRITICAL: TARGET_REPO_URL not set!"; 
        echo "Orchestrator failed to provide target repository configuration"; 
        exit 1; 
    fi
    
    echo "✅ Working on target: $TARGET_REPO_URL"
}
```

## Branch Naming from Config

### Read Branch Format
```bash
get_effort_branch_name() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    
    # Read format and prefix from config
    # Use text_editor tool with view command to read target-repo-config.yaml:
    # Find branch_naming.effort_format and branch_naming.project_prefix
    local format="<value from branch_naming.effort_format>"
    local project_prefix="<value from branch_naming.project_prefix>"
    
    # Add prefix with slash if it exists
    local prefix=""
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then 
        prefix="${project_prefix}/"; 
    fi
    
    # Replace variables
    local branch_name="${format}"
    branch_name="${branch_name//\{prefix\}/$prefix}"
    branch_name="${branch_name//\{phase\}/$phase}"
    branch_name="${branch_name//\{wave\}/$wave}"
    branch_name="${branch_name//\{effort_name\}/$effort_name}"
    
    echo "$branch_name"
}

# Example usage
BRANCH=$(get_effort_branch_name 1 2 "api-types")
# Result without prefix: phase1/wave2/api-types
# Result with prefix "idpbuilder-oci-mgmt": idpbuilder-oci-mgmt/phase1/wave2/api-types
```

### Get Integration Branch Name
```bash
get_integration_branch_name() {
    local phase="$1"
    local wave="$2"
    
    # Read format and prefix from config
    # Use text_editor tool with view command to read target-repo-config.yaml:
    # Find branch_naming.integration_format and branch_naming.project_prefix
    local format="<value from branch_naming.integration_format>"
    local project_prefix="<value from branch_naming.project_prefix>"
    
    # Add prefix with slash if it exists
    local prefix=""
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then 
        prefix="${project_prefix}/"; 
    fi
    
    # Replace variables
    local branch_name="${format}"
    branch_name="${branch_name//\{prefix\}/$prefix}"
    branch_name="${branch_name//\{phase\}/$phase}"
    branch_name="${branch_name//\{wave\}/$wave}"
    
    echo "$branch_name"
}

# Example usage
INTEGRATION_BRANCH=$(get_integration_branch_name 1 2)
# Result without prefix: phase1/wave2/integration
# Result with prefix "idpbuilder-oci-mgmt": idpbuilder-oci-mgmt/phase1/wave2/integration
```

### Get Phase Integration Branch Name
```bash
get_phase_integration_branch_name() {
    local phase="$1"
    
    # Read format and prefix from config
    # Use text_editor tool with view command to read target-repo-config.yaml:
    # Find branch_naming.phase_integration_format and branch_naming.project_prefix
    local format="<value from branch_naming.phase_integration_format>"
    local project_prefix="<value from branch_naming.project_prefix>"
    
    # Add prefix with slash if it exists
    local prefix=""
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then 
        prefix="${project_prefix}/"; 
    fi
    
    # Replace variables
    local branch_name="${format}"
    branch_name="${branch_name//\{prefix\}/$prefix}"
    branch_name="${branch_name//\{phase\}/$phase}"
    
    echo "$branch_name"
}

# Example usage
PHASE_BRANCH=$(get_phase_integration_branch_name 1)
# Result without prefix: phase1/integration
# Result with prefix "idpbuilder-oci-mgmt": idpbuilder-oci-mgmt/phase1/integration
```

## Workspace Path from Config

### Calculate Effort Path
```bash
get_effort_workspace() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    
    # Read patterns from config
    # Use text_editor tool with view command to read target-repo-config.yaml:
    # Find workspace.efforts_root and workspace.effort_path
    local efforts_root="<value from workspace.efforts_root>"
    local effort_path="<value from workspace.effort_path>"
    
    # Replace variables
    effort_path="${effort_path//\{phase\}/$phase}"
    effort_path="${effort_path//\{wave\}/$wave}"
    effort_path="${effort_path//\{effort_name\}/$effort_name}"
    
    echo "$SF_ROOT/$efforts_root/$effort_path"
}
```

## Validation Requirements

### Config File Validation
```bash
validate_target_config() {
    local config="$SF_ROOT/target-repo-config.yaml"
    
    # Required fields
    local required_fields=(
        ".target_repository.url"
        ".target_repository.base_branch"
        ".branch_naming.effort_format"
        ".workspace.efforts_root"
    )
    
    for field in "${required_fields[@]}"; do
        # Use text_editor tool with view command to read config file:
        # Find the specified field from the config
        local value="<value from field>"
        if [ -z "$value" ] || [ "$value" = "null" ]; then 
            echo "❌ Missing required field: $field"; 
            return 1; 
        fi
    done
    
    # Validate URL format
    if ! [[ "$TARGET_REPO_URL" =~ ^(https://|git@|ssh://) ]]; then 
        echo "❌ Invalid repository URL format"; 
        return 1; 
    fi
    
    echo "✅ Target configuration valid"
}
```

## Common Violations

### ❌ Working in SF Instance Repo
```bash
# BAD: Agent tries to create code in SF instance
cd /workspaces/my-project-sw-factory
mkdir src
echo "package main" > src/main.go
# VIOLATION: Writing code in SF instance repository!
```

### ✅ Working in Target Clone
```bash
# GOOD: Agent works in effort workspace
cd /workspaces/my-project-sw-factory/efforts/phase1/wave1/api-types
# This is a clone of the target repository
echo "package api" > pkg/api/types.go
# CORRECT: Writing code in target repository clone
```

## Grading Enforcement

### Critical Failures
- No target-repo-config.yaml: Cannot start (-100%)
- Invalid configuration: Cannot proceed (-100%)
- Working in wrong repository: Immediate failure (-100%)
- Missing TARGET_REPO_URL in environment: Stop work (-50%)

### Configuration Checks
- Config loaded at orchestrator startup: Required
- Config verified by all agents: Required
- Branch names follow config format: Required
- Workspace paths follow config structure: Required

## Integration with Other Rules

### Works with R192 (Separation)
- R191 defines WHAT to work on
- R192 enforces WHERE NOT to work

### Works with R193 (Clone Protocol)
- R191 provides the URL to clone
- R193 defines HOW to clone it

### Works with R194-R195 (Tracking/Push)
- R191 provides branch naming
- R194-R195 enforce proper git operations

## Example Configuration Usage

```yaml
# target-repo-config.yaml
target_repository:
  url: "https://github.com/cnoe-io/idpbuilder.git"
  base_branch: "main"

branch_naming:
  effort_format: "phase{phase}/wave{wave}/{effort_name}"

workspace:
  efforts_root: "efforts"
  effort_path: "phase{phase}/wave{wave}/{effort_name}"
```

Results in:
- Clone URL: `https://github.com/cnoe-io/idpbuilder.git`
- Branch: `phase1/wave2/api-types`
- Workspace: `/workspaces/idpbuilder-sw-factory/efforts/phase1/wave2/api-types/`

---
**Remember:** The Software Factory instance is for PLANNING. The target repository is for CODE!