# Rule R205: SW Engineer Split Directory Navigation Protocol

## Rule Statement
SW Engineers working on splits MUST read the split plan metadata BEFORE performing any preflight checks. The engineer must navigate to the correct split directory specified in the split plan and verify the branch matches before proceeding with ANY other validation or implementation work.

## Criticality Level
**BLOCKING** - Working in wrong directory causes integration failures and grading penalties

## Enforcement Mechanism
- **Technical**: Read split metadata → Navigate to directory → Verify branch → THEN preflight
- **Behavioral**: Exit immediately if directory/branch mismatch
- **Grading**: -40% for working in wrong location (Critical failure)

## Core Principle

```
Split Navigation = Read Metadata → Change Directory → Verify → THEN Preflight
NEVER start work without being in the correct split directory
The split plan contains the truth about where to work
```

## Detailed Requirements

### SW ENGINEER: Pre-Preflight Split Navigation

```bash
# ✅✅✅ CORRECT - Navigate BEFORE preflight checks
handle_split_work() {
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔍 SPLIT WORK DETECTION AND NAVIGATION"
    echo "═══════════════════════════════════════════════════════════════"
    
    # STEP 1: Detect if this is split work
    if ls SPLIT-PLAN-*.md 2>/dev/null | head -1 > /dev/null; then 
        echo "📋 Split work detected. Reading split plan..."; 
        
        # Find the specific split plan for this agent
        SPLIT_PLAN=$(ls SPLIT-PLAN-*.md 2>/dev/null | head -1); 
        
        if [ -z "$SPLIT_PLAN" ]; then 
            echo "❌ FATAL: Split work indicated but no SPLIT-PLAN found!"; 
            exit 1; 
        fi
        
        echo "Reading: $SPLIT_PLAN"; 
        
        # STEP 2: Extract metadata from split plan
        WORKING_DIR=$(grep "\*\*WORKING_DIRECTORY\*\*:" "$SPLIT_PLAN" | cut -d: -f2- | xargs); 
        EXPECTED_BRANCH=$(grep "\*\*BRANCH\*\*:" "$SPLIT_PLAN" | head -1 | cut -d: -f2- | xargs); 
        SPLIT_NUMBER=$(grep "\*\*SPLIT_NUMBER\*\*:" "$SPLIT_PLAN" | cut -d: -f2- | xargs); 
        
        if [ -z "$WORKING_DIR" ] || [ -z "$EXPECTED_BRANCH" ]; then 
            echo "❌ FATAL: Split plan missing required metadata!"; 
            echo "Orchestrator must update split plan with infrastructure details"; 
            exit 1; 
        fi
        
        echo "═══════════════════════════════════════════════════════════════"; 
        echo "📁 NAVIGATING TO SPLIT DIRECTORY"; 
        echo "Target directory: $WORKING_DIR"; 
        echo "Expected branch: $EXPECTED_BRANCH"; 
        echo "Split number: $SPLIT_NUMBER"; 
        echo "═══════════════════════════════════════════════════════════════"; 
        
        # STEP 3: Navigate to the correct directory
        if [ ! -d "$WORKING_DIR" ]; then 
            echo "❌ FATAL: Split directory does not exist: $WORKING_DIR"; 
            echo "Orchestrator must create infrastructure first!"; 
            exit 1; 
        fi; 
        
        cd "$WORKING_DIR"; 
        echo "✅ Changed to split directory: $(pwd)"; 
        
        # STEP 4: Verify we're on the correct branch
        CURRENT_BRANCH=$(git branch --show-current); 
        if [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then 
            echo "❌ FATAL: Branch mismatch!"; 
            echo "Current: $CURRENT_BRANCH"; 
            echo "Expected: $EXPECTED_BRANCH"; 
            echo "Cannot proceed - infrastructure setup issue"; 
            exit 1; 
        fi; 
        
        echo "✅ Verified branch: $CURRENT_BRANCH"; 
        
        # Note: No need to check for SPLIT-MARKER.txt
        # The presence of SPLIT-PLAN-XXX.md is sufficient to indicate this is a split
        
        echo "═══════════════════════════════════════════════════════════════"; 
        echo "✅ SPLIT NAVIGATION COMPLETE - Ready for preflight checks"; 
        echo "═══════════════════════════════════════════════════════════════"; 
        echo "⚠️ CRITICAL: All subsequent work MUST happen in this directory: $(pwd)"; 
    else 
        echo "ℹ️ Not split work - proceeding with normal preflight checks"; 
    fi
}

# THIS MUST RUN BEFORE ANY PREFLIGHT CHECKS
handle_split_work

# ONLY NOW run preflight checks
run_preflight_checks
```

### Reading Split Metadata

```bash
# Extract and validate split metadata
extract_split_metadata() {
    local split_plan=$1
    
    # Create associative array for metadata
    declare -A METADATA
    
    # Extract all metadata fields
    METADATA[WORKING_DIRECTORY]=$(grep "\*\*WORKING_DIRECTORY\*\*:" "$split_plan" | cut -d: -f2- | xargs)
    METADATA[BRANCH]=$(grep "\*\*BRANCH\*\*:" "$split_plan" | head -1 | cut -d: -f2- | xargs)
    METADATA[REMOTE]=$(grep "\*\*REMOTE\*\*:" "$split_plan" | cut -d: -f2- | xargs)
    METADATA[BASE_BRANCH]=$(grep "\*\*BASE_BRANCH\*\*:" "$split_plan" | cut -d: -f2- | xargs)
    METADATA[SPLIT_NUMBER]=$(grep "\*\*SPLIT_NUMBER\*\*:" "$split_plan" | cut -d: -f2- | xargs)
    METADATA[TOTAL_SPLITS]=$(grep "\*\*TOTAL_SPLITS\*\*:" "$split_plan" | cut -d: -f2- | xargs)
    
    # Validate required fields
    REQUIRED_FIELDS=("WORKING_DIRECTORY" "BRANCH" "SPLIT_NUMBER")
    for field in "${REQUIRED_FIELDS[@]}"; do
        if [ -z "${METADATA[$field]}" ]; then
            echo "❌ Missing required metadata: $field"
            return 1
        fi
    done
    
    # Display extracted metadata
    echo "📋 Split Metadata Extracted:"
    for key in "${!METADATA[@]}"; do
        echo "  $key: ${METADATA[$key]}"
    done
    
    # Export for use
    export SPLIT_WORKING_DIR="${METADATA[WORKING_DIRECTORY]}"
    export SPLIT_BRANCH="${METADATA[BRANCH]}"
    export SPLIT_NUMBER="${METADATA[SPLIT_NUMBER]}"
    
    return 0
}
```

### Sequence Validation

```python
def validate_split_sequence():
    """Validate the complete split sequence is correct"""
    
    sequence_checks = {
        'split_plan_exists': False,
        'metadata_present': False,
        'directory_exists': False,
        'branch_exists': False,
        'remote_tracking': False,
        'in_correct_directory': False,
        'on_correct_branch': False
    }
    
    # 1. Check split plan exists
    split_plans = glob.glob("SPLIT-PLAN-*.md")
    if split_plans:
        sequence_checks['split_plan_exists'] = True
        split_plan = split_plans[0]
        
        # 2. Check metadata is present
        with open(split_plan, 'r') as f:
            content = f.read()
            if 'WORKING_DIRECTORY:' in content and 'BRANCH:' in content:
                sequence_checks['metadata_present'] = True
                
                # Extract metadata
                working_dir = None
                expected_branch = None
                for line in content.split('\n'):
                    if 'WORKING_DIRECTORY:' in line:
                        working_dir = line.split(':', 1)[1].strip()
                    if line.startswith('BRANCH:'):
                        expected_branch = line.split(':', 1)[1].strip()
                
                # 3. Check directory exists
                if working_dir and os.path.exists(working_dir):
                    sequence_checks['directory_exists'] = True
                    
                    # 4. Check branch exists
                    os.chdir(working_dir)
                    current_branch = subprocess.check_output(
                        ['git', 'branch', '--show-current'],
                        text=True
                    ).strip()
                    
                    if current_branch:
                        sequence_checks['branch_exists'] = True
                        
                        # 5. Check remote tracking
                        tracking = subprocess.check_output(
                            ['git', 'branch', '-vv'],
                            text=True
                        )
                        if f'origin/{current_branch}' in tracking:
                            sequence_checks['remote_tracking'] = True
                        
                        # 6. Check we're in correct directory
                        if os.getcwd() == working_dir:
                            sequence_checks['in_correct_directory'] = True
                        
                        # 7. Check we're on correct branch
                        if current_branch == expected_branch:
                            sequence_checks['on_correct_branch'] = True
    
    # Determine if sequence is valid
    all_valid = all(sequence_checks.values())
    
    if not all_valid:
        print("❌ SPLIT SEQUENCE VALIDATION FAILED:")
        for check, passed in sequence_checks.items():
            status = "✅" if passed else "❌"
            print(f"  {status} {check}")
        return False
    
    print("✅ All split sequence checks passed")
    return True
```

## Common Violations to Avoid

### ❌ Starting Preflight in Wrong Directory
```bash
# WRONG - Running preflight before navigation
if [[ $(pwd) != */efforts/phase*/wave*/* ]]; then
    echo "Not in effort directory!"  # Already too late!
    exit 1
fi
```

### ❌ Ignoring Split Metadata
```bash
# WRONG - Not reading the split plan metadata
if [ -f "SPLIT-PLAN-001.md" ]; then
    # Just starts working without checking directory
    implement_split
fi
```

### ❌ Manual Directory Navigation
```bash
# WRONG - Guessing the directory
cd ../api-types--split-001  # Don't guess! Read the metadata!
```

### ✅ Correct Split Navigation
```bash
# RIGHT - Read metadata, navigate, verify, then proceed
SPLIT_PLAN=$(ls SPLIT-PLAN-*.md | head -1)
WORKING_DIR=$(grep "\*\*WORKING_DIRECTORY\*\*:" "$SPLIT_PLAN" | cut -d: -f2- | xargs)
cd "$WORKING_DIR"
verify_branch_matches_metadata
# NOW safe to run preflight checks
```

## Integration with Other Rules

- **R204**: Orchestrator creates infrastructure and updates metadata
- **R202**: Single SW engineer handles all splits sequentially
- **R199**: Single Code Reviewer creates all split plans
- **R001**: Preflight checks run AFTER navigation
- **R177**: Working directory enforcement validates AFTER navigation

## Grading Impact

- **No metadata in split plan**: -30% (Planning failure)
- **Wrong directory navigation**: -40% (Critical failure)
- **Branch mismatch**: -35% (Infrastructure failure)
- **Skipping navigation**: -45% (Protocol violation)
- **Starting work before verification**: -40% (Safety violation)

## Critical Working Directory Rule

### 🚨 ALL WORK HAPPENS IN THE SPLIT DIRECTORY

Once you've navigated to the split directory, **EVERY SINGLE COMMAND** must be executed there:

```bash
# After navigation to split directory
echo "Current directory: $(pwd)"  # e.g., /efforts/phase1/wave1/api-types--split-001

# ✅ CORRECT - All work in split directory
git add .                    # Adds files in split directory
git commit -m "feat: split"  # Commits in split branch
go test ./pkg/...            # Tests split code
make build                   # Builds split code

# ❌ WRONG - Trying to work elsewhere
cd ../api-types             # NO! Stay in split directory
cd /workspace/main-repo     # NO! Stay in split directory
```

### Tool Usage in Split Directory

When using Claude Code tools in split work:
1. **Read**: Read files relative to split directory
2. **Write**: Write files in split directory pkg/ structure
3. **Edit**: Edit files in split directory
4. **Bash**: Run commands IN the split directory
5. **Grep**: Search within split directory

```python
# Example: Verify you're in the right place before ANY work
import os
current_dir = os.getcwd()
if "--split-" not in current_dir:
    raise Exception("NOT IN SPLIT DIRECTORY! Navigate first!")
```

## Summary

**Remember**:
- Split plan contains the truth about where to work
- Navigate to directory BEFORE any checks
- Verify branch matches metadata
- ONLY THEN proceed with preflight
- **ALL WORK HAPPENS IN THE SPLIT DIRECTORY**
- Exit immediately on any mismatch
- This prevents ALL downstream issues