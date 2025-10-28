# Rule R198: Line Counter Tool Usage Protocol

## Rule Statement
Agents MUST use the line counter tool for all size measurements. The tool AUTOMATICALLY detects the correct base branch - no manual specification needed! Manual counting = AUTOMATIC -100% FAILURE.

**CRITICAL**: The tool ONLY counts implementation code! Tests, demos, docs, configs are automatically excluded.

**CRITICAL FOR SPARSE CLONES**: The tool is in the orchestrator's project root, NOT in your sparse clone! You must find the project root first (where `orchestrator-state-v3.json` lives), then use `${PROJECT_ROOT}/tools/line-counter.sh [BRANCH_TO_MEASURE]`.

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

### ✅ CORRECT USAGE - AUTOMATIC BASE DETECTION

```bash
# STEP 1: Navigate to your effort directory
cd /path/to/efforts/phase1/wave1/api-types

# STEP 2: Find the orchestrator's project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Project root: $PROJECT_ROOT"

# STEP 3: Run the tool - it auto-detects the base!
$PROJECT_ROOT/tools/line-counter.sh  # Measures current branch
# OR specify a branch to measure
$PROJECT_ROOT/tools/line-counter.sh phase1/wave1/api-types

# The tool automatically determines:
# - For splits: Previous split or original effort
# - For efforts: Wave/phase integration branch
# - For integrations: Previous phase or main

# Output will show what base was detected:
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 📊 Line Counter - Software Factory 2.0
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 📌 Analyzing branch: phase1/wave1/api-types
# 🎯 Detected base:    phase1-wave1-integration
# 🏷️  Project prefix:  idpbuilder-oci-go-cr (from current directory)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# ✅ Total implementation lines: 245 (excludes tests/demos/docs)
```

### ❌❌❌ WRONG USAGE - AUTOMATIC -100% FAILURE

```bash
# ❌❌❌ FATAL - Manual counting = AUTOMATIC FAILURE
wc -l *.go  # -100% FAILURE!
find . -name "*.go" | xargs wc -l  # -100% FAILURE!

# ❌❌❌ FATAL - Using old -b parameter syntax (OUTDATED!)
$PROJECT_ROOT/tools/line-counter.sh -b main -c phase1/wave1/api-types  # OLD VERSION - NO LONGER NEEDED!

# ❌ WRONG - ./tools doesn't exist in sparse clone!
./tools/line-counter.sh  # FAILS - no tools/ directory in sparse clone!

# ❌ WRONG - Wrong directory
cd /home/vscode/workspaces/idpbuilder-oci-mgmt
./tools/line-counter.sh  # Wrong directory!

# ✅ CORRECT - Let the tool auto-detect!
$PROJECT_ROOT/tools/line-counter.sh  # Auto-detects base from branch name
$PROJECT_ROOT/tools/line-counter.sh phase1/wave1/api-types  # Measure specific branch
```

## Understanding the Tool

The line counter tool:
1. **AUTOMATICALLY** detects the correct base branch from naming conventions
2. **OPTIONALLY** accepts a branch name to measure (default: current branch)
3. **Compares** branch against auto-detected base
4. **EXCLUDES** all non-implementation files automatically:
   - Test files (*_test.go, test/*, tests/*, *.test.*)
   - Demo files (demos/*, demo-*, DEMO.md, example-*)
   - Documentation (*.md, docs/*, README*, LICENSE*)
   - Generated code (*.pb.go, *_generated.*, *.gen.go)
   - Configuration (*.json, *.yaml, *.yml, *.toml)
   - Dependencies (vendor/*, node_modules/*, .cache/*)
   - Build artifacts (bin/*, dist/*, build/*, *.o, *.so)
5. **Shows** what base branch was detected in output
6. **ONLY COUNTS** critical path implementation code

### How Auto-Detection Works
```bash
# The tool analyzes branch naming patterns:
# SUPPORTS OPTIONAL PROJECT PREFIXES (e.g., my-project/phase1/wave1/effort)

# PROJECT PREFIX DETECTION (NEW!):
#   1. Reads target-repo-config.yaml from:
#      - Current directory
#      - $CLAUDE_PROJECT_DIR (if set)
#      - Parent directories up to orchestrator root
#   2. Uses configured prefix if found (100% accurate)
#   3. Falls back to pattern detection if no config
#   4. Shows prefix source in output for transparency

# For splits (--split-NNN):
#   - First split: base is original effort
#   - Later splits: base is previous split
#   - Preserves project prefix in base detection

# For efforts (phase*/wave*/effort-name or project/phase*/wave*/effort-name):
#   - Checks for wave integration branch first
#   - Falls back to phase integration  
#   - Uses main/master as last resort
#   - Handles project-prefixed branches automatically

# For integration branches:
#   - Uses previous phase or main
#   - Supports project-prefixed integration branches

# Shows detected base and prefix source in output:
# 🎯 Detected base: phase1-wave1-integration
# 🏷️  Project prefix: idpbuilder-oci-go-cr (from current directory)
# OR when no config:
# 🎯 Detected base: my-project/phase1-wave1-integration
# 🏷️  Project prefix: (using pattern detection)
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
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Found project root: $PROJECT_ROOT"

# 4. Run the line counter - it auto-detects the base!
OUTPUT=$($PROJECT_ROOT/tools/line-counter.sh)
echo "$OUTPUT"
SIZE=$(echo "$OUTPUT" | grep "Total" | awk '{print $NF}')
echo "Measured size: $SIZE lines"

# 5. Check if under limit
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
    local branch_name=$2  # Optional
    
    echo "Checking size of $effort_dir..."
    
    # Find project root
    PROJECT_ROOT=$(cd "$effort_dir" && pwd)
    while [ "$PROJECT_ROOT" != "/" ]; do
        if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then
            break
        fi
        PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
    done
    
    # Run line counter (auto-detects base)
    if [ -n "$branch_name" ]; then
        SIZE=$($PROJECT_ROOT/tools/line-counter.sh "$branch_name" | grep "Total" | awk '{print $NF}')
    else
        cd "$effort_dir"
        SIZE=$($PROJECT_ROOT/tools/line-counter.sh | grep "Total" | awk '{print $NF}')
    fi
    
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
# This means orchestrator-state-v3.json wasn't found
# SOLUTION: Find the actual orchestrator project directory
find /home -name "orchestrator-state-v3.json" -type f 2>/dev/null | head -1
# Then use that directory's tools/line-counter.sh
```

### Error: "./tools/line-counter.sh: No such file or directory"
```bash
# CAUSE: You're in a sparse clone - tools/ isn't included!
# SOLUTION: Find project root first
PROJECT_ROOT=$(pwd); 
while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break; 
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT"); 
done
$PROJECT_ROOT/tools/line-counter.sh  # Use PROJECT_ROOT, not ./
```

### Error: "Could not determine base branch"
```bash
# If you see:
Error: Could not determine base branch for 'my-branch'

# CAUSE: Branch doesn't follow naming conventions
# SOLUTION: Ensure branch follows Software Factory patterns:
#   - Efforts: phaseN/waveM/effort-name
#   - Efforts with prefix: project-name/phaseN/waveM/effort-name
#   - Splits: phaseN/waveM/effort--split-NNN
#   - Splits with prefix: project-name/phaseN/waveM/effort--split-NNN
#   - Integration: phaseN-waveM-integration
#   - Integration with prefix: project-name/phaseN-waveM-integration
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
# ✅ CORRECT - Check each effort by branch name
check_all_efforts() {
    # Find project root first
    PROJECT_ROOT=$(pwd)
    while [ "$PROJECT_ROOT" != "/" ]; do
        [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
        PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
    done
    
    # Check each effort branch
    for branch in phase1/wave1/*; do
        echo "Checking $branch..."
        $PROJECT_ROOT/tools/line-counter.sh "$branch"
    done
}

# OR check from within each directory
check_all_efforts_v2() {
    for effort_dir in efforts/phase1/wave1/*; do
        if [ -d "$effort_dir" ]; then
            echo "Checking $(basename $effort_dir)..."
            cd "$effort_dir"
            $PROJECT_ROOT/tools/line-counter.sh
            cd - > /dev/null
        fi
    done
}
```

## Integration with Size Limit Rule (R007)

This rule supports R007 (800 line limit) by ensuring accurate measurement:

```bash
# Continuous size monitoring
monitor_size_during_work() {
    # Find project root once
    PROJECT_ROOT=$(pwd)
    while [ "$PROJECT_ROOT" != "/" ]; do
        [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
        PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
    done
    
    while developing; do
        # Every ~100 lines of code - auto-detects base!
        SIZE=$($PROJECT_ROOT/tools/line-counter.sh | grep "Total" | awk '{print $NF}')
        
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

1. **Find PROJECT_ROOT first** (where orchestrator-state-v3.json lives)
2. **Run tool - it auto-detects the base**
3. **Tool shows detected base in output**
4. **Optional: pass branch name to measure**
5. **Never use old -b/-c parameters**
6. **Tool compares against correct base automatically**
7. **Tool excludes generated code automatically**

## Quick Reference Card

```bash
# ✅ RIGHT WAY - AUTO-DETECTION:
cd efforts/phase1/wave1/my-effort
# Find project root first!
PROJECT_ROOT=$(pwd); while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break; 
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT"); 
done
# Run tool - it auto-detects everything!
$PROJECT_ROOT/tools/line-counter.sh  # Measures current branch
$PROJECT_ROOT/tools/line-counter.sh phase1/wave1/api-types  # Measure specific branch
$PROJECT_ROOT/tools/line-counter.sh my-project/phase1/wave1/api-types  # With project prefix

# ❌❌❌ AUTOMATIC FAILURES (-100% GRADE):
wc -l *.go  # Manual counting = -100% FAILURE!
find . -name "*.go" | xargs wc -l  # Manual counting = -100% FAILURE!

# ❌ OUTDATED/WRONG WAYS:
./tools/line-counter.sh  # No tools/ in sparse clone!
$PROJECT_ROOT/tools/line-counter.sh -b main -c branch  # OLD SYNTAX - not needed!
```

## Grading Impact

- **Manual counting (wc -l, etc.)**: -100% (AUTOMATIC FAILURE)
- **Not using line-counter.sh tool**: -100% (AUTOMATIC FAILURE)
- **Not checking size regularly**: -15% (Process violation)
- **Wrong directory execution**: -20% (Context error)
- **Exceeding 800 lines**: -40% (Size limit violation)
- **Tool auto-detects correct base**: NO PENALTY (working as designed)