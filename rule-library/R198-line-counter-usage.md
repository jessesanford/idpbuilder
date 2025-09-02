# Rule R198: Line Counter Tool Usage Protocol

## Rule Statement
Agents MUST use the line counter tool CORRECTLY with the MANDATORY `-b` parameter to specify the correct base branch. Manual counting or using wrong base branches = AUTOMATIC -100% FAILURE.

**CRITICAL FOR SPARSE CLONES**: The tool is in the orchestrator's project root, NOT in your sparse clone! You must find the project root first (where `orchestrator-state.yaml` lives), then use `${PROJECT_ROOT}/tools/line-counter.sh -b [BASE_BRANCH] -c [CURRENT_BRANCH]`.

## Criticality Level
**BLOCKING** - Incorrect measurement leads to size limit violations

## Enforcement Mechanism
- **Technical**: Tool execution validation
- **Behavioral**: Immediate correction on misuse
- **Grading**: -100% for manual counting or wrong base branch

## 🚨🚨🚨 CRITICAL: How to Use Line Counter

### ⚠️ SPARSE CLONE REALITY CHECK ⚠️
**Your effort directory is a SPARSE CLONE - it doesn't have the tools/ directory!**
The line counter lives in the orchestrator's project root, not in your sparse clone.

### ✅ CORRECT USAGE - MUST SPECIFY BASE BRANCH

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

# STEP 3: Identify the CORRECT base branch
# For efforts: Use the phase integration branch, NOT "main"!
BASE_BRANCH="phase1/integration"  # From orchestrator-state.yaml
CURRENT_BRANCH=$(git branch --show-current)

# STEP 4: Run the tool WITH MANDATORY PARAMETERS
$PROJECT_ROOT/tools/line-counter.sh -b $BASE_BRANCH -c $CURRENT_BRANCH

# CORRECT examples:
$PROJECT_ROOT/tools/line-counter.sh -b phase1/integration -c phase1/wave1/api-types
$PROJECT_ROOT/tools/line-counter.sh -b phase2/integration -c phase2/wave1/effort1

# Output will show:
# Counting lines in phase1/wave1/api-types against phase1/integration...
# Total non-generated lines: 245
```

### ❌❌❌ WRONG USAGE - AUTOMATIC -100% FAILURE

```bash
# ❌❌❌ FATAL - Manual counting = AUTOMATIC FAILURE
wc -l *.go  # -100% FAILURE!
find . -name "*.go" | xargs wc -l  # -100% FAILURE!

# ❌❌❌ FATAL - Using "main" as base for efforts = AUTOMATIC FAILURE  
$PROJECT_ROOT/tools/line-counter.sh -b main -c phase1/wave1/api-types  # -100% FAILURE!

# ❌❌❌ FATAL - No parameters = WRONG (outdated usage)
$PROJECT_ROOT/tools/line-counter.sh  # WRONG - must specify base branch!

# ❌ WRONG - ./tools doesn't exist in sparse clone!
./tools/line-counter.sh  # FAILS - no tools/ directory in sparse clone!

# ❌ WRONG - Missing base branch parameter
$PROJECT_ROOT/tools/line-counter.sh -c phase1/wave1/api-types  # Missing -b!

# ❌ WRONG - Wrong directory
cd /home/vscode/workspaces/idpbuilder-oci-mgmt
./tools/line-counter.sh -b phase1/integration  # Wrong directory!
```

## Understanding the Tool

The line counter tool:
1. **REQUIRES** `-b` parameter to specify base branch
2. **REQUIRES** `-c` parameter to specify current branch (or auto-detects)
3. **Compares** current branch against specified base branch
4. **Excludes** generated code (zz_generated*, *.pb.go, etc.)
5. **Must be run** from within the effort directory

### How It Works With Parameters
```bash
# The tool uses your specified parameters:
BASE_BRANCH="$1"  # From -b parameter (MANDATORY)
CURRENT_BRANCH="$2"  # From -c parameter (or auto-detected)
git diff $BASE_BRANCH..$CURRENT_BRANCH --numstat | 
  grep -v "zz_generated" | 
  grep -v ".pb.go" | 
  # ... counts lines

# CRITICAL: For efforts, BASE must be phase integration branch!
# NOT "main" - that would count ALL phase code as your effort!
```

## Correct Workflow for Size Checking

### For Code Reviewers (MANDATORY PROCESS)
```bash
# 1. Verify you're in the right directory
pwd
# Should output: /path/to/efforts/phase1/wave1/your-effort

# 2. Confirm you're on the right branch
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"

# 3. Find the project root (where orchestrator lives)
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Found project root: $PROJECT_ROOT"

# 4. CRITICAL: Identify the CORRECT base branch
# Check orchestrator-state.yaml for phase integration branch
BASE_BRANCH=$(grep "current_phase_integration:" $PROJECT_ROOT/orchestrator-state.yaml -A 2 | grep "branch:" | awk '{print $2}')
echo "Base branch: $BASE_BRANCH"

# 5. Run the line counter WITH MANDATORY PARAMETERS
SIZE=$($PROJECT_ROOT/tools/line-counter.sh -b $BASE_BRANCH -c $CURRENT_BRANCH | grep "Total" | awk '{print $NF}')
echo "Measured size: $SIZE lines"

# 6. Check if under limit
if [ "$SIZE" -gt 800 ]; then
    echo "❌ OVER LIMIT ($SIZE) - MUST CREATE SPLIT PLAN"
else
    echo "✅ Under limit ($SIZE) - can approve"
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
# ✅ RIGHT WAY - MUST SPECIFY BASE BRANCH:
cd efforts/phase1/wave1/my-effort
# Find project root first!
PROJECT_ROOT=$(pwd); while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break; 
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT"); 
done
# Get base branch from orchestrator-state.yaml
BASE_BRANCH="phase1/integration"  # NOT "main"!
$PROJECT_ROOT/tools/line-counter.sh -b $BASE_BRANCH -c $(git branch --show-current)

# ❌❌❌ AUTOMATIC FAILURES (-100% GRADE):
wc -l *.go  # Manual counting = -100% FAILURE!
$PROJECT_ROOT/tools/line-counter.sh -b main  # Wrong base = -100% FAILURE!
$PROJECT_ROOT/tools/line-counter.sh  # No parameters = WRONG!

# ❌ OTHER WRONG WAYS:
./tools/line-counter.sh  # No tools/ in sparse clone!
$PROJECT_ROOT/tools/line-counter.sh -c branch-name  # Missing -b parameter!
```

## Grading Impact

- **Manual counting (wc -l, etc.)**: -100% (AUTOMATIC FAILURE)
- **Using "main" as base for efforts**: -100% (AUTOMATIC FAILURE)
- **Missing -b parameter**: -50% (Critical parameter missing)
- **Not checking size regularly**: -15% (Process violation)
- **Wrong directory execution**: -20% (Context error)
- **Exceeding 800 lines**: -40% (Size limit violation)