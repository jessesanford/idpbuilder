# 🚨🚨🚨 RULE R192 - Repository Separation Enforcement

**Criticality:** BLOCKING - Mixing repos = Catastrophic confusion  
**Grading Impact:** -75% for repository boundary violations  
**Enforcement:** CONTINUOUS - Every file operation checked

## Rule Statement

Agents MUST maintain ABSOLUTE SEPARATION between Software Factory instance repository and target project repository. NO CODE in SF instance. NO CONFIGS in target clones.

## The Two Repository Rule

```
┌─────────────────────────────────────┐     ┌─────────────────────────────────────┐
│  SOFTWARE FACTORY INSTANCE REPO     │     │  TARGET PROJECT REPO (CLONES)       │
├─────────────────────────────────────┤     ├─────────────────────────────────────┤
│  PURPOSE: Planning & Orchestration  │     │  PURPOSE: Actual Code Development   │
│                                     │     │                                     │
│  ✅ ALLOWED:                        │     │  ✅ ALLOWED:                        │
│  - Configuration files (.yaml)      │     │  - Source code (.go, .py, .js...)  │
│  - Planning documents (.md)         │     │  - Tests (*_test.go, test_*.py)    │
│  - State files (state/*.yaml)       │     │  - Build files (Makefile, go.mod)  │
│  - TODO files (todos/*.todo)        │     │  - API definitions (*.proto, *.yaml)│
│  - Rules (rule-library/*.md)        │     │  - Documentation (pkg/*/README.md) │
│                                     │     │                                     │
│  ❌ FORBIDDEN:                      │     │  ❌ FORBIDDEN:                      │
│  - Source code                      │     │  - SF configuration files           │
│  - Application tests                │     │  - Planning documents               │
│  - Build artifacts                  │     │  - State files                      │
│  - Compiled binaries                │     │  - TODO files                       │
│  - Vendor dependencies              │     │  - Agent configurations             │
└─────────────────────────────────────┘     └─────────────────────────────────────┘
```

## Enforcement Mechanisms

### 1. Directory Guards
```bash
# Function every agent MUST use before writing
check_write_permission() {
    local target_path="$1"
    local file_type="$2"
    
    # Get absolute path
    local abs_path=$(realpath "$target_path")
    
    # Determine which repo we're in
    if [[ "$abs_path" =~ ^"$SF_ROOT" ]]; then
        # In SF instance repo
        if [[ "$file_type" =~ \.(go|py|js|java|cpp|c|rs|rb|php|ts|tsx|jsx)$ ]]; then
            echo "❌ VIOLATION R192: Cannot write source code in SF instance!"
            echo "Target path: $abs_path"
            echo "Detected source file type: $file_type"
            exit 1
        fi
        
        # Check if path is in writable areas
        if [[ ! "$abs_path" =~ (state/|todos/|planning/.*\.md$) ]]; then
            echo "⚠️ WARNING: Writing to restricted area in SF instance"
            echo "Only state/, todos/, and planning/ are writable"
        fi
    else
        # In target repo clone
        if [[ "$file_type" =~ (target-repo-config|rule-library|\.claude) ]]; then
            echo "❌ VIOLATION R192: Cannot write SF configs in target repo!"
            exit 1
        fi
    fi
}
```

### 2. Git Operation Guards
```bash
# Before any git operation
check_git_repo_context() {
    local repo_root=$(git rev-parse --show-toplevel 2>/dev/null)
    
    if [ -z "$repo_root" ]; then
        echo "❌ Not in a git repository!"
        return 1
    fi
    
    # Check for SF instance markers
    if [ -f "$repo_root/target-repo-config.yaml" ]; then
        echo "⚠️ In Software Factory instance repository"
        echo "Code changes should be in effort workspaces only!"
        
        # Check what's being committed
        local changed_files=$(git diff --cached --name-only)
        for file in $changed_files; do
            if [[ "$file" =~ \.(go|py|js|java|cpp|c|rs)$ ]]; then
                echo "❌ VIOLATION R192: Attempting to commit source code in SF instance!"
                echo "File: $file"
                return 1
            fi
        done
    fi
}
```

### 3. Agent Working Directory Enforcement
```bash
# SW Engineers MUST be in effort workspace
verify_sw_engineer_workspace() {
    local cwd=$(pwd)
    
    # Must be under efforts/
    if [[ ! "$cwd" =~ /efforts/phase[0-9]+/wave[0-9]+/ ]]; then
        echo "❌ VIOLATION R192: SW Engineer not in effort workspace!"
        echo "Current: $cwd"
        echo "Expected: */efforts/phase*/wave*/*"
        exit 1
    fi
    
    # Must be a git repo (clone of target)
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo "❌ VIOLATION R192: Not in a git repository!"
        echo "Effort workspace should be a clone of target repo"
        exit 1
    fi
    
    # Must NOT have SF markers
    if [ -f "target-repo-config.yaml" ] || [ -d ".claude" ]; then
        echo "❌ VIOLATION R192: This appears to be SF instance, not target clone!"
        exit 1
    fi
}
```

## Workspace Structure Enforcement

### Correct Structure Example
```
/workspaces/idpbuilder-sw-factory/           # SF INSTANCE
├── target-repo-config.yaml                  # ✅ Config in SF instance
├── rule-library/                            # ✅ Rules in SF instance
├── .claude/agents/                          # ✅ Agent configs in SF instance
├── planning/                                # ✅ Plans in SF instance
│   ├── phase1-plan.md
│   └── architecture.md
├── state/                                   # ✅ State in SF instance
│   └── orchestrator-state.yaml
├── todos/                                   # ✅ TODOs in SF instance
│   └── orchestrator-WAVE_START-*.todo
└── efforts/                                 # EFFORT WORKSPACES
    └── phase1/
        └── wave1/
            └── api-types/                   # ← THIS IS A TARGET REPO CLONE
                ├── .git/                    # ✅ Git repo (target)
                ├── pkg/                     # ✅ Source code here
                │   └── api/
                │       └── types.go
                ├── go.mod                   # ✅ Build files here
                └── Makefile                 # ✅ Build tools here
```

### Common Violations to Prevent

#### ❌ Source Code in SF Instance
```bash
# WRONG - Creating source in SF root
cd /workspaces/idpbuilder-sw-factory
mkdir pkg
echo "package main" > pkg/main.go
# VIOLATION: Source code in SF instance!
```

#### ❌ SF Configs in Target Clone
```bash
# WRONG - Copying configs to target
cd /workspaces/idpbuilder-sw-factory/efforts/phase1/wave1/api-types
cp ../../../../target-repo-config.yaml .
# VIOLATION: SF config in target clone!
```

#### ❌ Mixed State Files
```bash
# WRONG - State files in target
cd /workspaces/idpbuilder-sw-factory/efforts/phase1/wave1/api-types
echo "state: COMPLETE" > orchestrator-state.yaml
# VIOLATION: State belongs in SF instance!
```

## Grading Penalties

### Separation Violations
- Source code in SF instance: -75%
- SF configs in target clone: -50%
- State files in wrong repo: -40%
- TODOs in target clone: -30%
- Planning docs in target: -25%

### Pattern Violations
- Not using efforts/ directory: -50%
- Wrong directory structure: -30%
- Missing git repo in effort: -40%
- SF markers in target clone: -60%

## Recovery from Violations

### If Source Code in SF Instance
```bash
# Move to correct location
mv $SF_ROOT/pkg $SF_ROOT/efforts/phase1/wave1/api-types/
cd $SF_ROOT/efforts/phase1/wave1/api-types/
git add pkg/
git commit -m "fix: move source code to target clone"
```

### If Configs in Target Clone
```bash
# Remove from target
cd $EFFORT_WORKSPACE
rm target-repo-config.yaml
git rm target-repo-config.yaml
git commit -m "fix: remove SF config from target clone"
```

## Integration with Other Rules

### Depends on R191 (Config)
- R191 defines target repository
- R192 enforces separation from it

### Enables R193 (Clone Protocol)
- R192 requires separate clones
- R193 defines how to create them

### Supports R176-R180 (Workspace Isolation)
- Workspace isolation within efforts/
- Repository separation at higher level

## Validation Commands

```bash
# Orchestrator validates separation
validate_repository_separation() {
    echo "🔍 Validating repository separation..."
    
    # Check SF instance for source code
    if find "$SF_ROOT" -type f \( -name "*.go" -o -name "*.py" -o -name "*.js" \) 
        | grep -v "/efforts/" | head -1; then
        echo "❌ Found source code in SF instance!"
        return 1
    fi
    
    # Check efforts for SF configs
    if find "$SF_ROOT/efforts" -name "target-repo-config.yaml" | head -1; then
        echo "❌ Found SF configs in effort workspaces!"
        return 1
    fi
    
    echo "✅ Repository separation maintained"
}
```

---
**Remember:** Planning happens in SF instance. Coding happens in target clones. NEVER MIX!