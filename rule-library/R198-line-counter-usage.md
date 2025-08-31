# Rule R198: Line Counter Tool Usage Protocol

## Rule Statement
Agents MUST use the line counter tool CORRECTLY by running it FROM WITHIN the effort directory with NO PARAMETERS. 

**CRITICAL FOR SPARSE CLONES**: The tool is in the orchestrator's project root, NOT in your sparse clone! You must find the project root first (where `orchestrator-state.yaml` lives), then use `${PROJECT_ROOT}/tools/line-counter.sh`.

## Criticality Level
**BLOCKING** - Incorrect measurement leads to size limit violations

## Enforcement Mechanism
- **Technical**: Tool execution validation
- **Behavioral**: Immediate correction on misuse
- **Grading**: -25% for incorrect line counting

## 🚨🚨🚨 CRITICAL: How to Use Line Counter

### ⚠️ SPARSE CLONE REALITY CHECK ⚠️
**Your effort directory is a SPARSE CLONE - it doesn't have the tools/ directory!**
The line counter lives in the orchestrator's project root, not in your sparse clone.

### ✅ CORRECT USAGE - Find Project Root FIRST

```bash
# STEP 1: Navigate to your effort directory
cd /path/to/efforts/phase1/wave1/api-types

# STEP 2: Find the orchestrator's project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Project root: $PROJECT_ROOT"

# STEP 3: Run the tool from PROJECT_ROOT with NO PARAMETERS
$PROJECT_ROOT/tools/line-counter.sh

# OR if you know the absolute path:
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh

# Output will show:
# Counting lines in phase1/wave1/api-types (excluding generated code)...
# Total non-generated lines: 245
```

### ❌❌❌ WRONG USAGE - NEVER DO THIS

```bash
# ❌ WRONG - ./tools doesn't exist in sparse clone!
./tools/line-counter.sh  # FAILS - no tools/ directory in sparse clone!

# ❌ WRONG - Don't pass branch names as parameters
$PROJECT_ROOT/tools/line-counter.sh -c phase1/wave1/api-types  # WRONG!

# ❌ WRONG - Don't use flags
$PROJECT_ROOT/tools/line-counter.sh --help  # WRONG!
$PROJECT_ROOT/tools/line-counter.sh -h      # WRONG!

# ❌ WRONG - Don't run from wrong directory
cd /home/vscode/workspaces/idpbuilder-oci-mgmt
./tools/line-counter.sh phase1/wave1/api-types  # WRONG!

# ❌ WRONG - Don't pass paths as arguments
$PROJECT_ROOT/tools/line-counter.sh efforts/phase1/wave1/api-types  # WRONG!
```

## Understanding the Tool

The line counter tool:
1. **Automatically detects** your current branch
2. **Compares against** the base branch (usually main)
3. **Excludes** generated code (zz_generated*, *.pb.go, etc.)
4. **Requires** you to be IN the git repository

### How It Works Internally
```bash
# The tool does this automatically:
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
git diff main..$CURRENT_BRANCH --numstat | 
  grep -v "zz_generated" | 
  grep -v ".pb.go" | 
  # ... counts lines
```

## Correct Workflow for Size Checking

### For Software Engineers (IN SPARSE CLONES)
```bash
# 1. Verify you're in the right directory
pwd
# Should output: /path/to/efforts/phase1/wave1/your-effort

# 2. Confirm you're on the right branch
git branch --show-current
# Should output: phase1/wave1/your-effort

# 3. Find the project root (where orchestrator lives)
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Found project root: $PROJECT_ROOT"

# 4. Run the line counter from PROJECT ROOT
$PROJECT_ROOT/tools/line-counter.sh
# Will output: Total non-generated lines: XXX

# 5. Check if under limit
SIZE=$($PROJECT_ROOT/tools/line-counter.sh | grep "Total" | awk '{print $NF}')
if [ "$SIZE" -gt 800 ]; then
    echo "❌ OVER LIMIT ($SIZE) - Need to split"
else
    echo "✅ Under limit ($SIZE) - can continue"
fi
```

### For Orchestrators (Verifying Efforts)
```bash
verify_effort_size() {
    local effort_dir=$1
    
    echo "Checking size of $effort_dir..."
    cd "$effort_dir"
    
    # Must be IN the directory
    if [ ! -f ./tools/line-counter.sh ]; then
        echo "❌ Line counter not found"
        return 1
    fi
    
    # Run with NO parameters
    SIZE=$(./tools/line-counter.sh | grep "Total" | awk '{print $NF}')
    
    if [ "$SIZE" -gt 800 ]; then
        echo "❌ Effort is $SIZE lines - OVER LIMIT"
        return 1
    else
        echo "✅ Effort is $SIZE lines - within limit"
        return 0
    fi
}

# Use it
verify_effort_size "/path/to/efforts/phase1/wave1/api-types"
```

## Common Errors and Solutions

### Error: "Line counter not found at //tools/line-counter.sh"
```bash
# If you see:
❌ FATAL: Line counter not found at //tools/line-counter.sh

# CAUSE: The PROJECT_ROOT search returned "/" (root directory)
# This means orchestrator-state.yaml wasn't found
# SOLUTION: Find the actual orchestrator project directory
find /home -name "orchestrator-state.yaml" -type f 2>/dev/null | head -1
# Then use that directory's tools/line-counter.sh
```

### Error: "./tools/line-counter.sh: No such file or directory"
```bash
# CAUSE: You're in a sparse clone - tools/ isn't included!
# SOLUTION: Find project root first
PROJECT_ROOT=$(pwd); 
while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break; 
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT"); 
done
$PROJECT_ROOT/tools/line-counter.sh  # Use PROJECT_ROOT, not ./
```

### Error: "fatal: bad revision"
```bash
# If you see:
fatal: bad revision 'phase1/wave1/api-types..-c'

# CAUSE: You passed parameters to the tool
# SOLUTION: Run with NO parameters
$PROJECT_ROOT/tools/line-counter.sh  # Nothing after this!
```

### Error: "Not in a git repository"
```bash
# CAUSE: The sparse clone isn't set up properly
# SOLUTION: Verify you're in a git repo
git status  # Should work
git branch --show-current  # Should show your effort branch
```

## Multiple Effort Checking

If you need to check multiple efforts (orchestrator only):

```bash
# ✅ CORRECT - Check each from within its directory
check_all_efforts() {
    for effort_dir in efforts/phase1/wave1/*; do
        if [ -d "$effort_dir" ]; then
            echo "Checking $(basename $effort_dir)..."
            cd "$effort_dir"
            ./tools/line-counter.sh
            cd - > /dev/null
        fi
    done
}

# ❌ WRONG - Don't try to pass paths
for effort in efforts/phase1/wave1/*; do
    ./tools/line-counter.sh "$effort"  # WRONG!
done
```

## Integration with Size Limit Rule (R007)

This rule supports R007 (800 line limit) by ensuring accurate measurement:

```bash
# Continuous size monitoring
monitor_size_during_work() {
    while developing; do
        # Every ~100 lines of code
        SIZE=$(./tools/line-counter.sh | grep "Total" | awk '{print $NF}')
        
        if [ "$SIZE" -gt 700 ]; then
            echo "⚠️ WARNING: Approaching limit ($SIZE/800)"
        fi
        
        if [ "$SIZE" -gt 800 ]; then
            echo "❌ STOP: Over limit ($SIZE/800)"
            echo "Must split this effort"
            exit 1
        fi
    done
}
```

## Summary - Remember These Points

1. **ALWAYS cd to effort directory FIRST**
2. **Run tool with NO parameters**
3. **Tool auto-detects branch from git**
4. **Never pass branch names or paths**
5. **Never use flags like -c, --help, -h**
6. **Tool compares current branch vs main**
7. **Tool excludes generated code automatically**

## Quick Reference Card

```bash
# ✅ RIGHT WAY (for sparse clones):
cd efforts/phase1/wave1/my-effort
# Find project root first!
PROJECT_ROOT=$(pwd); while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break; 
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT"); 
done
$PROJECT_ROOT/tools/line-counter.sh  # NO PARAMETERS

# ❌ WRONG WAYS (all will fail):
./tools/line-counter.sh  # No tools/ in sparse clone!
$PROJECT_ROOT/tools/line-counter.sh -c branch-name  # No parameters!
$PROJECT_ROOT/tools/line-counter.sh --help  # No flags!
$PROJECT_ROOT/tools/line-counter.sh path/to/effort  # No paths!
$PROJECT_ROOT/tools/line-counter.sh -anything  # Nothing after command!
```

## Grading Impact

- **Using parameters with tool**: -25% (Incorrect usage)
- **Not checking size regularly**: -15% (Process violation)
- **Wrong directory execution**: -20% (Context error)
- **Exceeding 800 lines**: -40% (Size limit violation)