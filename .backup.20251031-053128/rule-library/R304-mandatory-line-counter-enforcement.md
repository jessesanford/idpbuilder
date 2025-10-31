# Rule R304: Mandatory Line Counter Tool Enforcement

## Rule Statement
Code reviewers MUST use the line-counter.sh tool for all size measurements. The tool AUTOMATICALLY detects the correct base branch - eliminating human error. Manual counting results in AUTOMATIC -100% FAILURE.

## Criticality Level
**BLOCKING** - Violations result in immediate and automatic failure

## Enforcement Mechanism
- **Technical**: Tool validation and parameter checking
- **Behavioral**: Immediate failure on violation
- **Grading**: -100% for manual counting or wrong base branch

## 🔴🔴🔴 SUPREME ENFORCEMENT - NO EXCEPTIONS 🔴🔴🔴

### Mandatory Requirements

1. **MUST use line-counter.sh tool** - No manual counting allowed
2. **Tool automatically detects base branch** - No manual -b parameter needed
3. **Base branch auto-detection is ALWAYS correct** - Follows SF naming conventions
4. **Optional: specify branch to measure** - Defaults to current branch
5. **MUST document in review** - Include exact command and output

### Automatic -100% Failures

```bash
# ❌❌❌ THESE RESULT IN AUTOMATIC -100% GRADE:

# Manual counting = AUTOMATIC FAILURE
wc -l *.go
find . -name "*.go" | xargs wc -l
grep -c "^" *.go
cloc .
sloccount .

# Using old -b parameter (tool updated!) = CONFUSION/ERROR
# The tool no longer accepts -b parameter
# Base branch is AUTOMATICALLY determined

# ✅✅✅ CORRECT USAGE (NEW VERSION):
./tools/line-counter.sh                    # Measure current branch
./tools/line-counter.sh phase1/wave1/effort1  # Measure specific branch
./tools/line-counter.sh -d                 # With detailed output
./tools/line-counter.sh effort--split-003 -v  # With verbose output
```

## 🔴🔴🔴 CRITICAL DISTINCTION: EFFORTS vs SPLITS 🔴🔴🔴

### PROJECT PREFIX SUPPORT (ENHANCED!)
**The tool now reads project prefixes from target-repo-config.yaml:**

```bash
# Automatic prefix detection priority:
# 1. Reads from target-repo-config.yaml (100% accurate)
# 2. Searches in current dir, $CLAUDE_PROJECT_DIR, parent dirs
# 3. Falls back to pattern detection if no config found
# 4. Shows prefix source in output for transparency

# Example output with configured prefix:
# 🏷️  Project prefix: idpbuilder-oci-go-cr (from current directory)
# OR when using pattern detection:
# 🏷️  Project prefix: my-project (using pattern detection)
```

**The tool handles optional project prefixes automatically:**
- Standard: `phase1/wave1/effort-name`
- With prefix: `my-project/phase1/wave1/effort-name`
- Split with prefix: `my-project/phase1/wave1/effort--split-001`
- Integration with prefix: `my-project/phase1-wave1-integration`

### EFFORT MEASUREMENTS
**Base Branch**: Phase/wave integration branch
- Phase 1, Wave 1 efforts: Measure against `main`
- Phase 1, Wave 2 efforts: Measure against `phase1-wave1-integration`
- Phase 2, Wave 1 efforts: Measure against `phase1-integration`
- With prefix: Tool preserves prefix when looking for integration branches

### SPLIT MEASUREMENTS (SEQUENTIAL, NOT INCREMENTAL!)
**Base Branch**: IMMEDIATE PREDECESSOR split branch

Per R308, splits follow SEQUENTIAL branching:
- **split-001**: Measures against the ORIGINAL too-large branch
- **split-002**: Measures against split-001 (NOT original, NOT main!)
- **split-003**: Measures against split-002 (NOT split-001, NOT main!)
- **split-N**: Measures against split-(N-1)

### 🚨 THE PROBLEM THIS SOLVES 🚨
```bash
# ❌❌❌ OLD TOOL - Human error selecting wrong base:
./tools/line-counter.sh -b main -c split-003  # OLD VERSION
# WRONG OUTPUT: 5,584 lines (includes split-001 + split-002 + split-003!)

# ✅✅✅ NEW TOOL - Automatic correct base detection:
./tools/line-counter.sh split-003  # NEW VERSION - NO -b NEEDED!
# Tool automatically detects: base = split-002
# CORRECT OUTPUT: 280 lines (ONLY split-003's incremental work)
```

**MEASURING AGAINST MAIN INCLUDES ALL PREVIOUS SPLITS' WORK!**
This leads to false violations and rejected valid splits.

## Correct Usage Pattern

### Step 1: Navigate to Effort Directory
```bash
cd /path/to/efforts/phase1/wave1/effort-name
```

### 🔴🔴🔴 CRITICAL: EFFORTS ARE SEPARATE GIT REPOSITORIES! 🔴🔴🔴

**UNDERSTANDING THE EFFORT REPOSITORY STRUCTURE:**
- Each effort is a SEPARATE git repository (not just a branch)
- Has its own .git directory and branch structure
- Line-counter.sh needs BRANCH NAMES from THIS repository
- NOT directory names, NOT paths, ONLY branch names!

**BEFORE MEASURING, VERIFY:**
```bash
# 1. Confirm you're in a git repository
pwd  # Should show: /path/to/efforts/phase1/wave1/effort-name
ls -la .git  # Should exist - this is the effort's git repo!

# 2. Check git status - code MUST be committed
git status  # Should show: "nothing to commit, working tree clean"
# If not clean, COMMIT FIRST:
git add -A
git commit -m "feat: implementation complete for measurement"

# 3. Push to remote (REQUIRED for valid measurement)
git push

# 4. Get the current branch name (NOT directory name!)
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"  # e.g., phase1-wave1-effort-name

# 5. Find the base/integration branch IN THIS REPO
git branch -a | grep -E "integration|main"
# Should show something like: phase1/integration or main
BASE_BRANCH="phase1/integration"  # Use what exists in THIS repo
```

### Step 2: Find Project Root
```bash
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Project root: $PROJECT_ROOT"
```

### Step 3: Run Tool - Automatic Base Detection!
```bash
# CRITICAL: Run FROM WITHIN the effort repository!
CURRENT_BRANCH=$(git branch --show-current)

# Tool might be at different relative paths:
if [[ -f "$PROJECT_ROOT/tools/line-counter.sh" ]]; then
    TOOL="$PROJECT_ROOT/tools/line-counter.sh"
elif [[ -f "../../tools/line-counter.sh" ]]; then
    TOOL="../../tools/line-counter.sh"
elif [[ -f "../../../tools/line-counter.sh" ]]; then
    TOOL="../../../tools/line-counter.sh"
else
    echo "ERROR: Cannot find line-counter.sh!"
    exit 1
fi

# NEW: Just run tool - it figures out the base automatically!
OUTPUT=$($TOOL)  # Auto-detects current branch and base
echo "$OUTPUT"  # Shows detected base branch

# Extract the line count
SIZE=$(echo "$OUTPUT" | grep "Total" | awk '{print $NF}')
echo "Measured size: $SIZE lines"

# The tool will show what base it detected:
# 🎯 Detected base: phase1-wave1-integration
# ✅ Total non-generated lines: 456
```

### ❌❌❌ COMMON FATAL ERRORS ❌❌❌

```bash
# ❌ WRONG: Using old -b syntax (tool updated!)
./tools/line-counter.sh -b main -c go-containerregistry-image-builder
# ERROR: Tool no longer accepts -b parameter!

# ❌ WRONG: Manual counting
wc -l *.go
find . -name "*.go" | xargs wc -l
# ERROR: Manual counting = -100% FAILURE!

# ❌ WRONG: Measuring uncommitted code
# Making changes...
./tools/line-counter.sh  # Without committing first
# ERROR: Uncommitted changes won't be measured by git diff!

# ✅ RIGHT: From effort repo with committed code
cd efforts/phase1/wave1/my-effort  # CD into effort repo
git add -A && git commit -m "feat: complete" && git push  # Commit & push
../../tools/line-counter.sh  # Tool auto-detects everything!
# Output shows detected base and accurate count
```

## Why This Matters

### The Problem We're Solving
Code reviewers were:
1. Doing manual counts with `wc -l`
2. Using "main" as base branch, counting entire phase as effort
3. **Using directory names instead of branch names** (CRITICAL ERROR!)
4. **Not understanding efforts are separate git repositories**
5. Creating incorrect split plans thinking 856-line efforts were 10,000+ lines
6. Causing massive integration failures and rework

### Real Example of the Problem (SOLVED!)
```bash
# BEFORE: Code reviewer manually specified wrong base:
cd efforts/phase2/wave1/go-containerregistry-image-builder
../../tools/line-counter.sh -b main -c phase2-wave1-gcr-image-builder
# Result: 11,180 lines (WRONG! Included all phase work)

# NOW: Tool auto-detects correct base:
cd efforts/phase2/wave1/go-containerregistry-image-builder
../../tools/line-counter.sh  # No parameters needed!
# Tool output:
# 🎯 Detected base: phase2-wave1-integration
# ✅ Total non-generated lines: 856
# Result: 856 lines (CORRECT! Only effort changes)

# WITH PROJECT PREFIX - Also works automatically:
../../tools/line-counter.sh idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder
# Tool output:
# 🎯 Detected base: idpbuilder-oci-go-cr/phase2-wave1-integration
# ✅ Total non-generated lines: 856
```

### The Solution
- Tool enforces consistent measurement
- Correct base branch ensures only effort changes are counted
- Automated exclusion of generated code
- Verifiable and reproducible measurements

## Integration with Other Rules

- **R198**: Line counter usage protocol
- **R007**: 800-line size limit enforcement
- **R200**: Measure only effort changeset
- **R221**: CD before every bash command

## Split Measurement Validation Protocol

### MANDATORY: Identify Split Context Before Measuring
```bash
identify_split_base() {
    local CURRENT_BRANCH="$1"
    
    # Check if this is a split branch
    if [[ "$CURRENT_BRANCH" =~ split-([0-9]{3}) ]]; then
        SPLIT_NUM="${BASH_REMATCH[1]}"
        SPLIT_NUM_INT=$((10#$SPLIT_NUM))  # Convert to int, removing leading zeros
        
        if [[ $SPLIT_NUM_INT -eq 1 ]]; then
            # split-001 measures against original branch
            echo "🔍 This is split-001, finding original branch..."
            # Original branch is typically the effort branch before splitting
            ORIGINAL=$(git branch -a | grep -v "split-" | grep -E "effort|feature" | head -1)
            echo "Base: $ORIGINAL"
            return 0
        else
            # split-N measures against split-(N-1)
            PREV_NUM=$((SPLIT_NUM_INT - 1))
            PREV_SPLIT=$(printf "split-%03d" $PREV_NUM)
            echo "🔍 This is split-$SPLIT_NUM, base is: $PREV_SPLIT"
            echo "Base: $PREV_SPLIT"
            return 0
        fi
    else
        # Not a split, use phase/wave integration
        echo "🔍 Regular effort, using phase/wave integration as base"
        return 1
    fi
}

# Example usage in code review:
CURRENT=$(git branch --show-current)
if identify_split_base "$CURRENT"; then
    echo "✅ Split detected, using sequential base"
else
    echo "✅ Regular effort, using integration base"
fi
```

### Validation Checklist for Reviewers
```bash
echo "📋 SPLIT MEASUREMENT VALIDATION CHECKLIST"
echo "========================================="
echo "[ ] 1. Is this a split branch? (contains 'split-XXX')"
echo "[ ] 2. If split-001: Base = original too-large branch"
echo "[ ] 3. If split-N: Base = split-(N-1)"
echo "[ ] 4. NEVER use 'main' as base for splits"
echo "[ ] 5. NEVER use phase integration as base for splits"
echo "[ ] 6. Document the base branch selection reasoning"
echo "[ ] 7. Verify measurement shows ONLY incremental work"
```

## Code Reviewer Implementation

```python
def measure_effort_size(effort_dir, project_root):
    """Measure effort size using mandatory line counter tool
    
    NEW: No base_branch parameter needed - tool auto-detects!
    """
    
    if not os.path.exists(f"{project_root}/tools/line-counter.sh"):
        raise FileNotFoundError("❌ FATAL: Line counter tool not found!")
    
    # Get current branch
    current_branch = subprocess.check_output(
        ["git", "branch", "--show-current"],
        cwd=effort_dir
    ).decode().strip()
    
    # NEW: Run tool without -b parameter - it auto-detects base!
    cmd = [
        f"{project_root}/tools/line-counter.sh",
        current_branch  # Just specify what to measure
    ]
    
    result = subprocess.run(
        cmd,
        cwd=effort_dir,
        capture_output=True,
        text=True,
        check=True
    )
    
    # Parse output to extract detected base
    detected_base = None
    for line in result.stdout.split('\n'):
        if 'Detected base:' in line:
            detected_base = line.split('Detected base:')[1].strip()
    
    # Parse size
    for line in result.stdout.split('\n'):
        if 'Total' in line:
            size = int(line.split()[-1])
            return {
                'size': size,
                'compliant': size <= 800,
                'command': ' '.join(cmd),
                'auto_detected_base': detected_base,  # NEW!
                'current_branch': current_branch,
                'output': result.stdout
            }
    
    raise ValueError("Could not parse line counter output")
```

## Review Report Documentation

Every code review MUST include:

```yaml
size_measurement:
  tool_used: "${PROJECT_ROOT}/tools/line-counter.sh"
  command_executed: "./tools/line-counter.sh phase1/wave1/effort1"  # NEW: No -b parameter!
  auto_detected_base: "phase1/integration"  # Tool figured this out automatically!
  current_branch: "phase1/wave1/effort1"
  measured_lines: 687
  limit: 800
  compliant: true
  timestamp: "2025-01-20T10:30:00Z"
  raw_output: |
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    📊 Line Counter - Software Factory 2.0
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    📌 Analyzing branch: phase1/wave1/effort1
    🎯 Detected base:    phase1/integration
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ✅ Total non-generated lines: 687
```

## Grading Impact

| Violation | Penalty | Recovery |
|-----------|---------|----------|
| Manual counting (wc -l, etc.) | -100% | None - Automatic failure |
| Not using line-counter.sh | -100% | None - Automatic failure |
| Modifying tool's base detection | -80% | Use tool as designed |
| No documentation in review | -30% | Must add measurement details |
| Incorrect tool path | -20% | Must find project root first |
| Using old -b parameter | -0% | Tool ignores it, shows warning |

## Summary

**THE TOOL NOW HANDLES EVERYTHING AUTOMATICALLY:**
1. Use the line-counter.sh tool (MANDATORY)
2. Tool auto-detects the correct base branch (no -b needed!)
3. Optionally specify branch to measure (defaults to current)
4. Document the exact command and output

**BENEFITS OF AUTO-DETECTION:**
- Eliminates human error in base selection
- Always measures against correct predecessor
- Handles splits correctly (split-N vs split-N-1)
- Handles efforts correctly (vs integration branches)
- Makes 5,584 vs 280 line mistakes IMPOSSIBLE

**ANY METHOD OTHER THAN line-counter.sh = AUTOMATIC FAILURE**

This is not negotiable. This is not optional. This is MANDATORY.