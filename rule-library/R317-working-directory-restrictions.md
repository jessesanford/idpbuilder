# ⚠️⚠️⚠️ RULE R317 - Working Directory Restrictions

**Criticality:** WARNING  
**Grading Impact:** -25% for violations  
**Enforcement:** Directory operation monitoring
**Applies To:** Orchestrator during all operations

## Rule Statement

The Orchestrator MUST NOT enter or operate within agent working directories. The Orchestrator must remain in the project root or orchestration directories only. Agent workspaces are FORBIDDEN territory.

## 🔴🔴🔴 DIRECTORY BOUNDARIES 🔴🔴🔴

### ✅ ORCHESTRATOR ALLOWED DIRECTORIES:
- `/` (project root)
- `/todos/` (TODO persistence)
- `/orchestration/` (if exists)
- `/phase-plans/` (planning documents)
- `/templates/` (read-only)
- `/rule-library/` (read-only)

### ❌ ORCHESTRATOR FORBIDDEN DIRECTORIES:
- `/efforts/**/` (ALL effort directories)
- `/working-copies/**/` (ALL agent workspaces)
- `/agent-workspaces/**/` (ALL agent areas)
- `/implementation/**/` (ALL code directories)
- `/src/**/` (source code directories)
- `/test/**/` (test directories)
- Any directory containing code files

## Detection Mechanism

```bash
# Validate before ANY cd command
validate_directory_change() {
    local target_dir="$1"
    
    # Check forbidden patterns
    if [[ "$target_dir" =~ efforts/|working-copies/|agent-workspaces/|src/|test/|implementation/ ]]; then
        echo "🚨 R317 VIOLATION: Orchestrator cannot enter $target_dir"
        echo "❌ This is an agent working directory!"
        echo "✅ Correct: Spawn agent to work in this directory"
        return 1
    fi
    
    # Check if directory contains code
    if ls "$target_dir"/*.{go,py,js,ts,java,cpp,c,rs,rb,php} 2>/dev/null | grep -q .; then
        echo "🚨 R317 VIOLATION: Directory contains code files!"
        echo "❌ Orchestrator must not operate in code directories"
        return 1
    fi
    
    echo "✅ Directory change validated: $target_dir"
    return 0
}
```

## Common Violations

### ❌ VIOLATION: Entering Effort Directories
```bash
# Orchestrator tries to check split progress
cd efforts/phase1/wave1/split1  # VIOLATION
ls -la  # Operating in forbidden directory
```

### ❌ VIOLATION: Navigating to Implementation
```bash
# Orchestrator enters code directory
cd working-copies/sw-engineer-1/src  # VIOLATION
cat main.go  # Reading from forbidden location
```

### ❌ VIOLATION: Operating in Agent Workspace
```bash
# Orchestrator tries to help setup
cd agent-workspaces/reviewer  # VIOLATION
mkdir reports  # Creating in forbidden area
```

## Correct Patterns

### ✅ GOOD: Operate from Root
```bash
# Stay in project root
pwd  # Should be project root

# Create infrastructure from root
mkdir -p efforts/phase1/wave1/split2

# Check contents without entering
ls -la efforts/phase1/wave1/split1/

# Read files using full paths
cat efforts/phase1/wave1/PLAN.md
```

### ✅ GOOD: Use Absolute Paths
```bash
# Never cd into agent directories
# Always use full paths from root

# Instead of:
# cd efforts/phase1 && ls

# Do:
ls efforts/phase1/

# Instead of:
# cd working-copies/sw-engineer && git status  

# Do:
git -C working-copies/sw-engineer status
```

### ✅ GOOD: Delegate Directory Operations
```bash
# Need work done in effort directory?
echo "📋 Task: Verify split1 test coverage"
echo "🚀 Spawning SW Engineer..."

Task: Check test coverage in split1
Agent: sw-engineer
Working Directory: efforts/phase1/wave1/split1
```

## Boundary Enforcement

```bash
# Wrapper for cd command
safe_cd() {
    local target="${1:-.}"
    
    # Resolve to absolute path
    local abs_path=$(realpath "$target")
    
    # Validate against R317
    validate_directory_change "$abs_path" || {
        echo "❌ Directory change blocked by R317"
        return 1
    }
    
    # Proceed with change
    cd "$target"
    echo "📍 Now in: $(pwd)"
}

# Alias cd to safe version for orchestrator
alias cd='safe_cd'
```

## Working Directory Audit

```bash
# Orchestrator self-check
audit_working_directory() {
    local current=$(pwd)
    
    echo "🔍 R317 Working Directory Audit"
    echo "Current: $current"
    
    # Check if in forbidden area
    if [[ "$current" =~ efforts/|working-copies/|agent-workspaces/ ]]; then
        echo "🚨 VIOLATION: Currently in forbidden directory!"
        echo "⚠️ Must return to project root immediately"
        cd /
        return 1
    fi
    
    echo "✅ Working directory compliant"
    return 0
}
```

## Recovery Protocol

If orchestrator finds itself in wrong directory:
```bash
# 1. Immediate evacuation
cd /  # Return to project root

# 2. Log the violation
echo "⚠️ R317 Recovery: Evacuated from $(pwd)"

# 3. Verify compliance
audit_working_directory

# 4. Continue from safe location
```

## Integration with State Machine

```yaml
# orchestrator-state.json should track
current_directory: "/"  # Always root for orchestrator
allowed_paths:
  - "/"
  - "/todos"
  - "/orchestration"
forbidden_paths:
  - "/efforts"
  - "/working-copies"
  - "/agent-workspaces"
```

## Grading Impact

- Entering effort directory: -10%
- Operating in agent workspace: -15%  
- Modifying files from within: -25%
- Repeated violations: Progressive penalties

## Monitoring Commands

```bash
# Track orchestrator movement
history | grep "cd " | tail -20

# Verify current location
pwd

# Check for violations in history
history | grep -E "cd.*(efforts|working-copies|agent-workspaces)"
```

---
**REMEMBER:** Orchestrator supervises from above, never enters the workspace. Stay in root!