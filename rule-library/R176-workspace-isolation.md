# 🚨🚨🚨 R176: Workspace Isolation Requirement 🚨🚨🚨

**Category:** Utility and Hook Management  
**Agents:** ALL (orchestrator, sw-engineer, code-reviewer, architect)  
**Criticality:** BLOCKING - Workspace contamination = AUTOMATIC FAILURE  
**Related:** R178 (Effort Directory Structure), R209 (Effort Directory Isolation)

## ABSOLUTE ISOLATION MANDATE

### 1. CORE PRINCIPLE

**EVERY EFFORT MUST BE COMPLETELY ISOLATED FROM:**
- Other efforts in the same wave
- Other waves in the same phase
- The main codebase
- The planning repository
- Any external codebases

### 2. ISOLATION STRUCTURE

```
target-repository/
├── efforts/                      # ISOLATION ROOT
│   ├── phase1/
│   │   ├── wave1/
│   │   │   ├── auth-base/       # ISOLATED EFFORT
│   │   │   │   ├── .git/        # Sparse checkout
│   │   │   │   ├── pkg/         # Effort's code
│   │   │   │   ├── go.mod       # If needed
│   │   │   │   └── README.md    # Effort docs
│   │   │   ├── user-service/    # ISOLATED EFFORT
│   │   │   │   ├── .git/
│   │   │   │   ├── pkg/
│   │   │   │   └── ...
│   │   │   └── api-gateway/     # ISOLATED EFFORT
│   │   └── wave2/
│   └── phase2/
└── pkg/                          # MAIN CODEBASE (DO NOT TOUCH!)
```

### 3. ISOLATION ENFORCEMENT

#### For SW-Engineer:
```bash
# MUST work ONLY in effort directory
EFFORT_DIR="/workspace/efforts/phase1/wave1/auth-base"

# FORBIDDEN - Working in main codebase
cd /workspace/pkg  # ❌ NEVER!

# REQUIRED - Work in isolated effort
cd $EFFORT_DIR/pkg  # ✅ CORRECT
```

#### For Code-Reviewer:
```bash
# MUST review ONLY effort's changes
REVIEW_DIR="/workspace/efforts/phase1/wave1/auth-base"

# Review ONLY this effort's diff
cd $REVIEW_DIR && git diff main...HEAD

# NOT the entire repository
cd /workspace && git diff  # ❌ WRONG!
```

#### For Orchestrator:
```bash
# MUST create isolated workspaces
create_effort_workspace() {
    local effort_name="$1"
    local effort_dir="efforts/phase$PHASE/wave$WAVE/$effort_name"
    
    # Create isolation
    mkdir -p "$effort_dir"
    
    # Sparse checkout for isolation
    cd "$effort_dir"
    git init
    git remote add origin $REPO_URL
    git sparse-checkout init
    git sparse-checkout set pkg/
    git checkout -b "$effort_name"
}
```

### 4. CONTAMINATION DETECTION

```bash
# Check for contamination
detect_contamination() {
    local effort_dir="$1"
    
    # Check file count
    local file_count=$(find "$effort_dir/pkg" -type f | wc -l)
    if [ $file_count -gt 1000 ]; then
        echo "❌ CONTAMINATION: $file_count files detected!"
        return 1
    fi
    
    # Check for foreign branches
    cd "$effort_dir"
    local branches=$(git branch -r | wc -l)
    if [ $branches -gt 5 ]; then
        echo "❌ CONTAMINATION: Multiple branches detected!"
        return 1
    fi
    
    echo "✅ Workspace is clean and isolated"
}
```

### 5. ISOLATION VIOLATIONS

**THESE ARE AUTOMATIC FAILURES:**
- ❌ Modifying files outside effort directory
- ❌ Importing code from other efforts directly
- ❌ Committing to main branch
- ❌ Working in planning repository
- ❌ Mixing changes from multiple efforts

### 6. PROPER ISOLATION WORKFLOW

```bash
# 1. Orchestrator creates isolation
cd /workspace
mkdir -p efforts/phase1/wave1/auth-base
cd efforts/phase1/wave1/auth-base
git init && git remote add origin $TARGET_REPO

# 2. SW-Engineer works in isolation
cd /workspace/efforts/phase1/wave1/auth-base
mkdir -p pkg/auth
# Create files ONLY here

# 3. Code-Reviewer reviews in isolation
cd /workspace/efforts/phase1/wave1/auth-base
git diff main...HEAD  # Only this effort's changes

# 4. Integration merges from isolation
cd /workspace/integration
git pull efforts/phase1/wave1/auth-base
```

### 7. CROSS-EFFORT DEPENDENCIES

When effort B depends on effort A:
```bash
# WRONG - Direct import
import "../auth-base/pkg/auth"  # ❌ NEVER!

# RIGHT - Through integration
# 1. Effort A completes and merges
# 2. Integration branch updated
# 3. Effort B rebases from integration
cd $EFFORT_B_DIR
git fetch origin integration
git rebase origin/integration
```

### 8. GRADING PENALTIES

```yaml
isolation_violations:
  working_in_main_codebase: -100%  # AUTOMATIC FAILURE
  cross_effort_contamination: -100%  # AUTOMATIC FAILURE
  modifying_planning_repo: -100%  # AUTOMATIC FAILURE
  no_isolation_structure: -50%
  importing_across_efforts: -40%
  branch_contamination: -30%
```

### 9. VALIDATION SCRIPT

```bash
#!/bin/bash
validate_workspace_isolation() {
    local effort_dir="$1"
    
    echo "🔍 Validating workspace isolation for $effort_dir"
    
    # Check we're in efforts directory
    if [[ ! "$effort_dir" == */efforts/* ]]; then
        echo "❌ Not in efforts directory!"
        return 176
    fi
    
    # Check for pkg directory
    if [ ! -d "$effort_dir/pkg" ]; then
        echo "❌ Missing pkg directory for isolation!"
        return 176
    fi
    
    # Check not in main workspace
    if [[ "$effort_dir" == */pkg ]] && [[ ! "$effort_dir" == */efforts/*/pkg ]]; then
        echo "❌ In main pkg directory!"
        return 176
    fi
    
    # Check git isolation
    cd "$effort_dir"
    local current_branch=$(git branch --show-current)
    if [[ "$current_branch" == "main" ]] || [[ "$current_branch" == "master" ]]; then
        echo "❌ On main branch - no isolation!"
        return 176
    fi
    
    echo "✅ Workspace properly isolated"
    return 0
}
```

### 10. RECOVERY FROM VIOLATIONS

If isolation is violated:
1. **STOP all work immediately**
2. **Identify contamination scope**
3. **Create new isolated workspace**
4. **Cherry-pick valid commits only**
5. **Verify isolation before continuing**

## INTEGRATION WITH OTHER RULES

- **R178**: Effort directory structure
- **R209**: Effort directory isolation protocol
- **R235**: Pre-flight verification checks isolation
- **R010**: Wrong location handling

## SUMMARY

**R176 Core Mandate: COMPLETE ISOLATION FOR EVERY EFFORT!**

- Each effort in its own directory
- No cross-contamination
- No main codebase modification
- Sparse checkouts for true isolation
- Dependencies only through integration

**Violation = Immediate failure and contamination cleanup required**

---
**Created**: Prevent catastrophic workspace contamination
**Purpose**: Ensure clean, isolated development environments
**Enforcement**: BLOCKING - Zero tolerance for violations