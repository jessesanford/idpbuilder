# Rule R204: ORCHESTRATOR (Not SW Engineer!) Must Create Complete Split Infrastructure

## Rule Statement
**THE ORCHESTRATOR** (not SW Engineers, not Code Reviewers) MUST create ALL split directories, working copies, branches, and remote tracking branches BEFORE spawning SW engineering agents for split implementation. Each split MUST have its own directory with `-SPLIT-00Z` suffix and a SEPARATE clone of the target repository. Split branches MUST include the project prefix from target-repo-config.yaml using the branch-naming-helpers.sh functions. Split directories must be in the /efforts folder alongside the too-large branch. Splits MUST follow SEQUENTIAL branching strategy: split-001 based on the same base as original, split-002 based on split-001, split-003 based on split-002, etc. This ensures correct line counting and clean integration.

## 🔴🔴🔴 CRITICAL CLARIFICATION: WHO DOES WHAT 🔴🔴🔴
- **CODE REVIEWER**: Creates split PLANS (SPLIT-INVENTORY.md, SPLIT-PLAN-XXX.md)
- **ORCHESTRATOR**: Creates split INFRASTRUCTURE (directories, clones, branches)
- **SW ENGINEER**: IMPLEMENTS in pre-created infrastructure

## 🔴🔴🔴 PREREQUISITE: Code Reviewer Must Create Split Plans First 🔴🔴🔴
Before the orchestrator can create split infrastructure:
1. Code Reviewer must have created SPLIT-INVENTORY.md in too-large branch
2. Code Reviewer must have created all SPLIT-PLAN-XXX.md files
3. These files must be committed and pushed to the too-large branch remote
4. Orchestrator will fetch these plans from the too-large branch

## 🔴🔴🔴 CRITICAL: Project Prefix Required 🔴🔴🔴
**ALL split branches MUST include the project prefix!**
- Use `get_split_branch_name()` from branch-naming-helpers.sh
- Example: `tmc-workspace/phase1/wave1/api-types--split-001`
- NOT: `phase1/wave1/api-types--split-001` (missing prefix)

## Criticality Level
**BLOCKING** - Improper split setup causes merge conflicts and integration failures

## Enforcement Mechanism
- **Technical**: Orchestrator creates all infrastructure before agent spawn
- **Behavioral**: SW engineers refuse work if infrastructure missing
- **Grading**: -45% for improper split setup (Major architectural failure)

## Core Principle

```
Split Infrastructure = Complete Setup BEFORE Implementation
Orchestrator creates ALL directories, branches, and tracking FIRST
SW Engineers receive ready-to-use split workspaces
Merge strategy planned from the beginning
```

## 🔴🔴🔴 MANDATORY: Sequential Split Branching 🔴🔴🔴

**SPLITS MUST BE CREATED SEQUENTIALLY - THIS IS NOT OPTIONAL!**

### Correct Sequential Strategy:
- **Split-001**: Based on same base as original (e.g., phase-integration)
- **Split-002**: Based on split-001 (NOT phase-integration!)
- **Split-003**: Based on split-002 (NOT phase-integration!)
- **Each split measures ONLY its own additions**

### Why This is MANDATORY:
1. **Line Counting Accuracy**: Each split shows correct line count
2. **No Cumulative Measurement**: Split-002 won't include split-001's lines
3. **Clean Integration**: Splits merge without conflicts
4. **Dependency Chain**: Later splits can use earlier split code

### FORBIDDEN Parallel Strategy:
```bash
# ❌❌❌ NEVER DO THIS - ALL SPLITS FROM SAME BASE
git checkout phase-integration
git checkout -b split-001  # OK: 400 lines
git checkout phase-integration  # WRONG!
git checkout -b split-002  # ERROR: Shows 800 lines (includes split-001!)
git checkout phase-integration  # WRONG!
git checkout -b split-003  # ERROR: Shows 1200 lines (includes all!)
```

## Detailed Requirements

### 🔴🔴🔴 ORCHESTRATOR RESPONSIBILITY: Complete Split Infrastructure Creation 🔴🔴🔴

**THIS IS THE ORCHESTRATOR'S JOB - NOT THE SW ENGINEER'S!**

```bash
# ✅✅✅ CORRECT - ORCHESTRATOR Creates ALL split infrastructure FIRST
# This function is called BY THE ORCHESTRATOR, not SW Engineers!
create_split_infrastructure() {
    local phase="1"
    local wave="1"
    local effort_name="api-types"
    local too_large_dir="/efforts/phase${phase}/wave${wave}/${effort_name}"
    
    # STEP 0: Verify split plans exist in too-large branch
    echo "═══════════════════════════════════════════════════════════════"
    echo "📋 CHECKING FOR SPLIT PLANS IN TOO-LARGE BRANCH"
    echo "═══════════════════════════════════════════════════════════════"
    
    cd "$too_large_dir"
    git pull  # Ensure we have latest from remote
    
    if [ ! -f "SPLIT-INVENTORY.md" ]; then
        echo "❌ ERROR: No SPLIT-INVENTORY.md found in too-large branch"
        echo "   Code Reviewer must create split plans first!"
        exit 1
    fi
    
    # Read number of splits from inventory
    local total_splits=$(grep -c "^| [0-9]" SPLIT-INVENTORY.md)
    echo "✅ Found SPLIT-INVENTORY.md with $total_splits splits planned"
    
    # Verify all split plan files exist
    for split_num in $(seq 1 $total_splits); do
        SPLIT_NAME=$(printf "%03d" $split_num)
        if [ ! -f "SPLIT-PLAN-${SPLIT_NAME}.md" ]; then
            echo "❌ ERROR: Missing SPLIT-PLAN-${SPLIT_NAME}.md"
            exit 1
        fi
    done
    echo "✅ All $total_splits split plan files found"
    
    # CRITICAL: Use branch naming helpers for consistency
    SF_ROOT="/workspaces/software-factory-2.0-template"  # Or wherever SF instance is
    source "$SF_ROOT/utilities/branch-naming-helpers.sh"
    
    # Get properly formatted branch name with prefix
    local too_large_branch=$(get_effort_branch_name "$phase" "$wave" "$effort_name")
    local too_large_dir="/efforts/phase${phase}/wave${wave}/${effort_name}"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔧 CREATING COMPLETE SPLIT INFRASTRUCTURE"
    echo "Too-large branch: $too_large_branch"
    echo "Project prefix: ${PROJECT_PREFIX:-none}"
    echo "Total splits needed: $total_splits"
    echo "═══════════════════════════════════════════════════════════════"
    
    # STEP 1: Determine base branch (what too-large branch was based on)
    cd "$too_large_dir"; 
    BASE_BRANCH=$(git log --format=%B -n 1 --grep="from branch" | grep -o "from branch: .*" | cut -d: -f2 | xargs); 
    if [ -z "$BASE_BRANCH" ]; then 
        # Fallback: find merge-base with main/master
        BASE_BRANCH=$(git merge-base HEAD main 2>/dev/null || git merge-base HEAD master 2>/dev/null || echo "main"); 
    fi; 
    echo "Base branch identified: $BASE_BRANCH"
    
    # STEP 2: Think about merge strategy
    echo "═══════════════════════════════════════════════════════════════"
    echo "🤔 MERGE STRATEGY PLANNING"
    echo "═══════════════════════════════════════════════════════════════"
    echo "Original branch ($too_large_branch) based on: $BASE_BRANCH"
    echo "Split strategy: SEQUENTIAL (each based on previous)"
    echo "Merge order will be: split-001 → split-002 → split-003 → $BASE_BRANCH"
    echo "This ensures clean integration without conflicts"
    
    # STEP 3: Create all split directories and branches
    for split_num in $(seq 1 $total_splits); do
        SPLIT_NAME=$(printf "%03d" $split_num)
        
        # Use branch naming helper for split branch
        SPLIT_BRANCH=$(get_split_branch_name "$too_large_branch" "$SPLIT_NAME")
        
        # Directory naming - CRITICAL: Use -SPLIT- suffix (uppercase)
        # This makes splits immediately visible in filesystem
        SPLIT_DIR="/efforts/phase${phase}/wave${wave}/${effort_name}-SPLIT-${SPLIT_NAME}"
        
        echo "═══════════════════════════════════════════════════════════════"
        echo "Creating Split $SPLIT_NAME Infrastructure"
        echo "Directory: $SPLIT_DIR"
        echo "Branch: $SPLIT_BRANCH"
        echo "═══════════════════════════════════════════════════════════════"
        
        # Create directory alongside too-large branch
        mkdir -p "$SPLIT_DIR"
        
        # Clone with proper base
        if [ $split_num -eq 1 ]; then 
            # First split based on original base branch
            echo "Basing split-001 on: $BASE_BRANCH"; 
            git clone --branch "$BASE_BRANCH" --sparse "$TARGET_REPO_URL" "$SPLIT_DIR"; 
        else 
            # Subsequent splits based on previous split (sequential)
            PREV_SPLIT=$(printf "%03d" $((split_num - 1))); 
            PREV_BRANCH=$(get_split_branch_name "$too_large_branch" "$PREV_SPLIT")
            echo "Basing split-${SPLIT_NAME} on: $PREV_BRANCH"; 
            git clone --branch "$PREV_BRANCH" --sparse "$TARGET_REPO_URL" "$SPLIT_DIR"; 
        fi
        
        cd "$SPLIT_DIR"
        
        # Set up sparse checkout
        git sparse-checkout init --cone
        git sparse-checkout set pkg/
        
        # Create and push the split branch
        git checkout -b "$SPLIT_BRANCH"
        git push -u origin "$SPLIT_BRANCH"
        
        # Verify remote tracking
        if git branch -vv | grep -q "$SPLIT_BRANCH.*origin/$SPLIT_BRANCH"; then 
            echo "✅ Remote tracking configured for $SPLIT_BRANCH"; 
        else 
            echo "❌ FATAL: Remote tracking failed for $SPLIT_BRANCH"; 
            exit 1; 
        fi
        
        # Fetch split plans from too-large branch remote
        echo "Fetching split plans from too-large branch..."
        
        # First, ensure too-large directory is up to date
        (cd "$too_large_dir" && git pull)
        
        # Check if split plans exist in too-large branch
        if [ ! -f "$too_large_dir/SPLIT-PLAN-${SPLIT_NAME}.md" ]; then
            echo "❌ ERROR: SPLIT-PLAN-${SPLIT_NAME}.md not found in too-large branch"
            echo "   Expected at: $too_large_dir/SPLIT-PLAN-${SPLIT_NAME}.md"
            echo "   Code Reviewer must create split plans first!"
            exit 1
        fi
        
        if [ ! -f "$too_large_dir/SPLIT-INVENTORY.md" ]; then
            echo "❌ ERROR: SPLIT-INVENTORY.md not found in too-large branch"
            echo "   Expected at: $too_large_dir/SPLIT-INVENTORY.md"
            echo "   Code Reviewer must create split inventory first!"
            exit 1
        fi
        
        # Copy ONLY the specific split plan from too-large branch
        echo "Copying split plan from too-large branch..."
        cp "$too_large_dir/SPLIT-PLAN-${SPLIT_NAME}.md" .
        echo "✅ Split plan SPLIT-PLAN-${SPLIT_NAME}.md copied successfully"
        
        # Note: NOT copying SPLIT-INVENTORY.md to avoid merge conflicts
        # Each split only needs its specific plan
        
        # CRITICAL: Update split plan with directory and branch metadata
        echo "Updating SPLIT-PLAN-${SPLIT_NAME}.md with infrastructure metadata..."
        cat >> "SPLIT-PLAN-${SPLIT_NAME}.md" << EOF

## 🚨 SPLIT INFRASTRUCTURE METADATA (Added by Orchestrator)
**WORKING_DIRECTORY**: $(pwd)
**BRANCH**: $SPLIT_BRANCH
**REMOTE**: origin/$SPLIT_BRANCH
**BASE_BRANCH**: $([ $split_num -eq 1 ] && echo $BASE_BRANCH || echo $PREV_BRANCH)
**SPLIT_NUMBER**: $SPLIT_NAME
**TOTAL_SPLITS**: $total_splits

### SW Engineer Instructions
1. READ this metadata FIRST
2. cd to WORKING_DIRECTORY above
3. Verify branch matches BRANCH above
4. ONLY THEN proceed with preflight checks
EOF
        
        # Also update the original split plan in too-large directory
        cat >> "$too_large_dir/SPLIT-PLAN-${SPLIT_NAME}.md" << EOF

## 🚨 SPLIT INFRASTRUCTURE CREATED (Updated by Orchestrator)
**SPLIT_DIRECTORY**: $SPLIT_DIR
**SPLIT_BRANCH**: $SPLIT_BRANCH
**STATUS**: Infrastructure Ready
**CREATED_AT**: $(date '+%Y-%m-%d %H:%M:%S')
EOF
        
        # Note: Not creating SPLIT-MARKER.txt as it's redundant
        # The presence of SPLIT-PLAN-XXX.md already indicates this is a split
        # All necessary metadata is in the split plan file itself
        
        # Commit initial split setup
        git add -A
        git commit -m "chore: initialize split-${SPLIT_NAME} from $too_large_branch"
        git push
        
        echo "✅ Split $SPLIT_NAME infrastructure complete"
        cd - > /dev/null
    done
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "✅ ALL SPLIT INFRASTRUCTURE CREATED"
    echo "Ready to spawn SW engineer for sequential implementation"
    echo "═══════════════════════════════════════════════════════════════"
}

# R296: After all splits are complete and verified
mark_original_branch_deprecated() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    local too_large_dir="/efforts/phase${phase}/wave${wave}/${effort_name}"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔄 R296: MARKING ORIGINAL BRANCH AS DEPRECATED"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Get properly formatted branch names
    source "$SF_ROOT/utilities/branch-naming-helpers.sh"
    local original_branch=$(get_effort_branch_name "$phase" "$wave" "$effort_name")
    local deprecated_branch="${original_branch}-deprecated-split"
    
    cd "$too_large_dir"
    
    # Rename local branch
    echo "Renaming local branch to mark as deprecated..."
    git branch -m "$original_branch" "$deprecated_branch"
    
    # Delete old remote branch and push renamed one
    echo "Updating remote branches..."
    git push origin ":$original_branch"  # Delete old remote
    git push -u origin "$deprecated_branch"  # Push renamed branch
    
    echo "✅ Branch renamed: $original_branch → $deprecated_branch"
    
    # Update state file per R296
    echo "Updating orchestrator-state.yaml with SPLIT_DEPRECATED status..."
    yq eval ".efforts_completed.\"$effort_name\".status = \"SPLIT_DEPRECATED\"" -i orchestrator-state.yaml
    yq eval ".efforts_completed.\"$effort_name\".deprecated_branch = \"$deprecated_branch\"" -i orchestrator-state.yaml
    yq eval ".efforts_completed.\"$effort_name\".do_not_integrate = true" -i orchestrator-state.yaml
    yq eval ".efforts_completed.\"$effort_name\".split_completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.yaml
    
    # Add replacement splits to state file
    local splits_array="["
    for split_num in $(seq 1 $total_splits); do
        SPLIT_NAME=$(printf "%03d" $split_num)
        SPLIT_BRANCH=$(get_split_branch_name "$original_branch" "$SPLIT_NAME")
        splits_array="${splits_array}\"$SPLIT_BRANCH\""
        [ $split_num -lt $total_splits ] && splits_array="${splits_array}, "
    done
    splits_array="${splits_array}]"
    
    yq eval ".efforts_completed.\"$effort_name\".replacement_splits = $splits_array" -i orchestrator-state.yaml
    
    # Commit state file update
    git add orchestrator-state.yaml
    git commit -m "state: mark $effort_name as SPLIT_DEPRECATED per R296"
    git push
    
    echo "✅ State file updated with SPLIT_DEPRECATED status"
    echo "✅ Original branch deprecated successfully"
}
```

### Directory Structure After Split Setup

```
efforts/
├── phase1/
│   └── wave1/
│       ├── api-types/                    # Original too-large branch (WILL BE DEPRECATED)
│       │   ├── SPLIT-INVENTORY.md        # Created by Code Reviewer
│       │   ├── SPLIT-PLAN-001.md         # Created by Code Reviewer
│       │   ├── SPLIT-PLAN-002.md         # Created by Code Reviewer
│       │   ├── SPLIT-PLAN-003.md         # Created by Code Reviewer
│       │   └── pkg/                      # Original oversized implementation
│       ├── api-types-SPLIT-001/          # Split 1 (separate clone, based on main)
│       │   ├── .git/                     # Own git repository
│       │   ├── SPLIT-PLAN-001.md         # Copied from too-large branch (with metadata)
│       │   └── pkg/                      # Ready for implementation
│       ├── api-types-SPLIT-002/          # Split 2 (separate clone, based on split-001)
│       │   ├── .git/                     # Own git repository
│       │   ├── SPLIT-PLAN-002.md         # Copied from too-large branch (with metadata)
│       │   └── pkg/                      # Ready for implementation
│       └── api-types-SPLIT-003/          # Split 3 (separate clone, based on split-002)
│           ├── .git/                     # Own git repository
│           ├── SPLIT-PLAN-003.md         # Copied from too-large branch (with metadata)
│           └── pkg/                      # Ready for implementation
```

Note: Each split directory only contains its specific SPLIT-PLAN-XXX.md to avoid merge conflicts.
The SPLIT-INVENTORY.md remains only in the too-large branch as documentation.

### Branch Structure for Clean Merging

```
main (or base branch)
  │
  ├── tmc-workspace/phase1/wave1/api-types (too large - will be abandoned)
  │
  ├── tmc-workspace/phase1/wave1/api-types--split-001
  │     │
  │     └── tmc-workspace/phase1/wave1/api-types--split-002
  │           │
  │           └── tmc-workspace/phase1/wave1/api-types--split-003
  
Sequential merge strategy:
1. Merge split-003 into split-002
2. Merge split-002 into split-001
3. Merge split-001 into main

Note: Project prefix (e.g., 'tmc-workspace') included in all branch names
```

### SW ENGINEER: Verify Infrastructure Before Starting

```bash
# SW Engineer MUST verify infrastructure exists
verify_split_infrastructure() {
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔍 VERIFYING SPLIT INFRASTRUCTURE"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Check we're in a split directory
    if [[ $(pwd) != *"--split-"* ]]; then 
        echo "❌ FATAL: Not in a split directory!"; 
        echo "Orchestrator must create split infrastructure first!"; 
        exit 1; 
    fi
    
    # Check for split marker
    if [ ! -f "SPLIT-MARKER.txt" ]; then 
        echo "❌ FATAL: No SPLIT-MARKER.txt found!"; 
        echo "This is not a properly configured split workspace!"; 
        exit 1; 
    fi
    
    # Check for split plan
    SPLIT_NUM=$(grep "Split" SPLIT-MARKER.txt | awk '{print $2}'); 
    if [ ! -f "SPLIT-PLAN-${SPLIT_NUM}.md" ]; then 
        echo "❌ FATAL: Split plan not found!"; 
        exit 1; 
    fi
    
    # Check git branch
    CURRENT_BRANCH=$(git branch --show-current); 
    if [[ "$CURRENT_BRANCH" != *"--split-from--"* ]]; then 
        echo "❌ FATAL: Not on a split branch!"; 
        echo "Current: $CURRENT_BRANCH"; 
        echo "Expected: *--split-from--* pattern"; 
        exit 1; 
    fi
    
    # Check remote tracking
    if ! git branch -vv | grep -q "$CURRENT_BRANCH.*origin/"; then 
        echo "❌ FATAL: No remote tracking for split branch!"; 
        exit 1; 
    fi
    
    echo "✅ Split infrastructure verified"
    echo "Split number: $SPLIT_NUM"
    echo "Branch: $CURRENT_BRANCH"
    echo "Ready to implement"
}
```

### Merge Strategy Documentation

```yaml
# Created by orchestrator in each split directory
merge_strategy:
  original_branch: "tmc-workspace/phase1/wave1/api-types"
  base_branch: "main"
  split_count: 3
  project_prefix: "tmc-workspace"
  
  splits:
    - number: "001"
      branch: "tmc-workspace/phase1/wave1/api-types--split-001"
      based_on: "main"
      merge_to: "main"
      
    - number: "002"
      branch: "tmc-workspace/phase1/wave1/api-types--split-002"
      based_on: "tmc-workspace/phase1/wave1/api-types--split-001"
      merge_to: "tmc-workspace/phase1/wave1/api-types--split-001"
      
    - number: "003"
      branch: "tmc-workspace/phase1/wave1/api-types--split-003"
      based_on: "tmc-workspace/phase1/wave1/api-types--split-002"
      merge_to: "tmc-workspace/phase1/wave1/api-types--split-002"
  
  merge_order:
    1: "Complete all splits first"
    2: "Merge split-003 → split-002"
    3: "Merge split-002 → split-001"
    4: "Merge split-001 → main"
    5: "Abandon original too-large branch"
```

### Common Violations to Avoid

### ❌ Spawning Before Infrastructure
```bash
# WRONG - No infrastructure created
Task: sw-engineer
Implement split-001
# Agent has nowhere to work!
```

### ❌ Creating Splits in Wrong Location
```bash
# WRONG - Not alongside original
efforts/splits/api-types-001/  # Wrong location!
# Should be: efforts/phase1/wave1/api-types--split-001/
```

### ❌ Wrong Base Branch
```bash
# WRONG - All splits from main
git checkout -b split-001 main
git checkout -b split-002 main  # Will cause conflicts!
git checkout -b split-003 main  # Will cause conflicts!
```

### ✅ Correct Sequential Setup
```bash
# RIGHT - Sequential dependency
git checkout -b split-001 main
git checkout -b split-002 split-001  # Based on previous
git checkout -b split-003 split-002  # Based on previous
```

## Integration with Other Rules

- **R199**: Single reviewer plans all splits
- **R202**: Single SW engineer implements all splits
- **R204**: Orchestrator creates infrastructure (THIS RULE)
- **R196**: Base branch selection from config
- **R193**: Effort clone protocol
- **R296**: Deprecated branch marking after splits complete

## Grading Impact

- **No infrastructure before spawn**: -45% (Major failure)
- **Wrong directory structure**: -30% (Organization failure)
- **No remote tracking**: -25% (Integration failure)
- **Wrong base branches**: -35% (Merge failure)
- **No merge strategy planning**: -20% (Planning failure)

## Summary

**Remember**:
- Orchestrator creates ALL split infrastructure FIRST
- Directories alongside too-large branch with clear naming
- Sequential base branches for clean merging
- Remote tracking for all split branches
- SW engineers verify before starting
- Merge strategy planned from the beginning