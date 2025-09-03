# Rule R304: Mandatory Line Counter Tool Enforcement

## Rule Statement
Code reviewers MUST use the line-counter.sh tool with the CORRECT base branch parameter. Manual counting or using incorrect base branches results in AUTOMATIC -100% FAILURE.

## Criticality Level
**BLOCKING** - Violations result in immediate and automatic failure

## Enforcement Mechanism
- **Technical**: Tool validation and parameter checking
- **Behavioral**: Immediate failure on violation
- **Grading**: -100% for manual counting or wrong base branch

## 🔴🔴🔴 SUPREME ENFORCEMENT - NO EXCEPTIONS 🔴🔴🔴

### Mandatory Requirements

1. **MUST use line-counter.sh tool** - No manual counting allowed
2. **MUST specify -b parameter** - Base branch is mandatory
3. **MUST use correct base branch** - Phase integration branch, NOT "main"
4. **MUST specify -c parameter** - Current branch (or auto-detect)
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

# Wrong base branch = AUTOMATIC FAILURE
./tools/line-counter.sh -b main -c phase1/wave1/effort1
./tools/line-counter.sh -b master -c phase1/wave1/effort1

# Missing parameters = MAJOR FAILURE
./tools/line-counter.sh  # No parameters
./tools/line-counter.sh -c phase1/wave1/effort1  # Missing -b
```

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
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Project root: $PROJECT_ROOT"
```

### Step 3: Identify Correct Base Branch
```bash
# From orchestrator-state.yaml
BASE_BRANCH=$(grep "current_phase_integration:" $PROJECT_ROOT/orchestrator-state.yaml -A 2 | \
              grep "branch:" | awk '{print $2}')
              
# Or from phase pattern
BASE_BRANCH="phase1/integration"  # NOT "main"!
```

### Step 4: Run Tool with Mandatory Parameters
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

# Run with BRANCH NAMES from THIS repository
SIZE=$($TOOL -b $BASE_BRANCH -c $CURRENT_BRANCH | \
       grep "Total" | awk '{print $NF}')
echo "Measured size: $SIZE lines"
```

### ❌❌❌ COMMON FATAL ERRORS ❌❌❌

```bash
# ❌ WRONG: Using directory name as branch
./tools/line-counter.sh -b main -c go-containerregistry-image-builder
# ERROR: "go-containerregistry-image-builder" is a directory, not a branch!

# ❌ WRONG: Running from main repo trying to measure effort
cd /workspaces/project  # Main repo
./tools/line-counter.sh -b main -c efforts/phase1/wave1/effort1
# ERROR: "efforts/phase1/wave1/effort1" is a path, not a branch!

# ❌ WRONG: Measuring uncommitted code
# Making changes...
./tools/line-counter.sh -b phase1/integration -c current-branch
# ERROR: Uncommitted changes won't be measured by git diff!

# ✅ RIGHT: From effort repo with actual branch names
cd efforts/phase1/wave1/my-effort  # CD into effort repo
git add -A && git commit -m "feat: complete" && git push  # Commit & push
CURRENT=$(git branch --show-current)  # Get actual branch name
BASE="phase1/integration"  # Base branch IN THIS REPO
../../tools/line-counter.sh -b "$BASE" -c "$CURRENT"
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

### Real Example of the Problem
```bash
# Code reviewer tried this and got 11,180 lines (WRONG!):
cd /workspaces/project
./tools/line-counter.sh -b main -c go-containerregistry-image-builder
# ERROR: "go-containerregistry-image-builder" is a DIRECTORY, not a branch!
# Tool can't find this "branch" and fails or gives wrong results

# What they SHOULD have done:
cd efforts/phase2/wave1/go-containerregistry-image-builder
git branch --show-current  # Returns: phase2-wave1-gcr-image-builder
../../tools/line-counter.sh -b phase2/integration -c phase2-wave1-gcr-image-builder
# Result: 856 lines (CORRECT!)
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

## Code Reviewer Implementation

```python
def measure_effort_size(effort_dir, project_root, base_branch):
    """Measure effort size using mandatory line counter tool"""
    
    # CRITICAL: Validate inputs first
    if base_branch == "main" or base_branch == "master":
        raise ValueError("❌ FATAL: Cannot use 'main/master' as base for effort measurement!")
    
    if not os.path.exists(f"{project_root}/tools/line-counter.sh"):
        raise FileNotFoundError("❌ FATAL: Line counter tool not found!")
    
    # Get current branch
    current_branch = subprocess.check_output(
        ["git", "branch", "--show-current"],
        cwd=effort_dir
    ).decode().strip()
    
    # Run the tool with MANDATORY parameters
    cmd = [
        f"{project_root}/tools/line-counter.sh",
        "-b", base_branch,
        "-c", current_branch
    ]
    
    result = subprocess.run(
        cmd,
        cwd=effort_dir,
        capture_output=True,
        text=True,
        check=True
    )
    
    # Parse output
    for line in result.stdout.split('\n'):
        if 'Total' in line:
            size = int(line.split()[-1])
            return {
                'size': size,
                'compliant': size <= 800,
                'command': ' '.join(cmd),
                'base_branch': base_branch,
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
  command_executed: "./tools/line-counter.sh -b phase1/integration -c phase1/wave1/effort1"
  base_branch: "phase1/integration"  # NEVER "main"
  current_branch: "phase1/wave1/effort1"
  measured_lines: 687
  limit: 800
  compliant: true
  timestamp: "2025-01-20T10:30:00Z"
  raw_output: |
    Counting lines in phase1/wave1/effort1 against phase1/integration...
    Excluding: *.pb.go, *_generated.go, vendor/*, *.md, *_test.go
    Total non-generated lines: 687
```

## Grading Impact

| Violation | Penalty | Recovery |
|-----------|---------|----------|
| Manual counting (wc -l, etc.) | -100% | None - Automatic failure |
| Using "main" as base | -100% | None - Automatic failure |
| Missing -b parameter | -50% | Must re-measure correctly |
| Wrong base branch (not main) | -40% | Must use phase integration |
| No documentation in review | -30% | Must add measurement details |
| Incorrect tool path | -20% | Must find project root first |

## Summary

**THERE IS ONLY ONE CORRECT WAY TO MEASURE SIZE:**
1. Use the line-counter.sh tool
2. With -b parameter specifying phase integration branch
3. With -c parameter specifying current branch
4. Document the exact command and output

**ANY OTHER METHOD = AUTOMATIC FAILURE**

This is not negotiable. This is not optional. This is MANDATORY.