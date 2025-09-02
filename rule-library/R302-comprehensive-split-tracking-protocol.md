# 🚨🚨🚨 BLOCKING RULE R302: Comprehensive Split Tracking Protocol

## Rule Definition
The Software Factory MUST maintain meticulous tracking of all split operations, including original branches, split branches, relationships, and integration status. This tracking enables proper merge planning, prevents duplicate work, and ensures correct branch selection.

## Criticality: 🚨🚨🚨 BLOCKING
Without proper split tracking, the system cannot determine which branches to integrate, leading to merge failures, duplicate PRs, and corrupted integrations.

## Requirements

### 1. Split Tracking Data Structure
Every split operation MUST be tracked in orchestrator-state.yaml with this structure:
```yaml
split_tracking:
  # Organized by effort ID
  effort-001-authentication:
    original_branch: "project/phase1/wave1/effort-001-authentication"
    status: "SPLIT_DEPRECATED"  # SPLIT_DEPRECATED or SPLIT_ARCHIVED
    split_count: 3
    total_original_lines: 1250
    split_reason: "exceeded 800 line limit"
    split_date: "2025-01-20T14:30:00Z"
    split_strategy: "SEQUENTIAL"  # MANDATORY: Each split based on previous
    split_branches:
      - branch: "project/phase1/wave1/effort-001-authentication--split-001"
        base_branch: "phase1-integration"  # Same base as original
        status: "COMPLETED"  # ACTIVE, COMPLETED, INTEGRATED
        lines: 450  # Measured against base_branch
        description: "Core authentication logic"
        completed_at: "2025-01-20T16:00:00Z"
        reviewed: true
        review_status: "APPROVED"
      - branch: "project/phase1/wave1/effort-001-authentication--split-002"
        base_branch: "project/phase1/wave1/effort-001-authentication--split-001"  # Based on split-001!
        status: "COMPLETED"
        lines: 425  # Measured against split-001, NOT phase1-integration
        description: "Token management and validation"
        completed_at: "2025-01-20T16:30:00Z"
        reviewed: true
        review_status: "APPROVED"
      - branch: "project/phase1/wave1/effort-001-authentication--split-003"
        base_branch: "project/phase1/wave1/effort-001-authentication--split-002"  # Based on split-002!
        status: "COMPLETED"
        lines: 375  # Measured against split-002, NOT phase1-integration
        description: "Session handling and cleanup"
        completed_at: "2025-01-20T17:00:00Z"
        reviewed: true
        review_status: "APPROVED"
    integration_status:
      integrated: true
      integration_branch: "project/phase1/wave1-integration"
      integration_date: "2025-01-20T18:00:00Z"
      integration_lines: 1250  # Sum of all splits
```

### 🔴🔴🔴 CRITICAL: Sequential Split Branching Strategy 🔴🔴🔴

**SPLITS MUST BE CREATED SEQUENTIALLY - EACH BASED ON THE PREVIOUS SPLIT**

#### Why Sequential Branching is MANDATORY:
1. **Correct Line Counting**: Each split measures ONLY its own additions
2. **Clean Integration**: No merge conflicts between splits
3. **Dependency Management**: Later splits can use earlier split code
4. **Progressive Development**: Each split builds on the previous

#### ✅ CORRECT Sequential Pattern:
```bash
# Split 1: Based on same base as original (e.g., phase-integration)
git checkout phase-integration
git checkout -b project/phase1/wave1/effort-001--split-001
# Implement 450 lines
./tools/line-counter.sh -b phase-integration -c $(git branch --show-current)
# Result: 450 lines ✓

# Split 2: Based on split-001 (NOT phase-integration!)
git checkout project/phase1/wave1/effort-001--split-001
git checkout -b project/phase1/wave1/effort-001--split-002
# Implement 425 more lines
./tools/line-counter.sh -b project/phase1/wave1/effort-001--split-001 -c $(git branch --show-current)
# Result: 425 lines ✓ (NOT 875 lines!)

# Split 3: Based on split-002 (NOT phase-integration!)
git checkout project/phase1/wave1/effort-001--split-002
git checkout -b project/phase1/wave1/effort-001--split-003
# Implement 375 more lines
./tools/line-counter.sh -b project/phase1/wave1/effort-001--split-002 -c $(git branch --show-current)
# Result: 375 lines ✓ (NOT 1250 lines!)
```

#### ❌ WRONG Parallel Pattern (FORBIDDEN):
```bash
# WRONG: All splits from same base
git checkout phase-integration
git checkout -b split-001  # 450 lines ✓
git checkout phase-integration  # WRONG!
git checkout -b split-002  # Will show 875 lines ✗
git checkout phase-integration  # WRONG!
git checkout -b split-003  # Will show 1250 lines ✗
```

#### Verification of Sequential Structure:
```bash
verify_split_sequence() {
    local effort="$1"
    
    # Get split branches from tracking
    SPLITS=$(yq ".split_tracking.\"${effort}\".split_branches[].branch" orchestrator-state.yaml)
    
    PREV_SPLIT=""
    for split in $SPLITS; do
        if [ -n "$PREV_SPLIT" ]; then
            # Verify this split has previous as ancestor
            if ! git merge-base --is-ancestor "$PREV_SPLIT" "$split"; then
                echo "❌ ERROR: $split not based on $PREV_SPLIT!"
                echo "   This violates sequential split requirement!"
                return 1
            fi
        fi
        PREV_SPLIT="$split"
    done
    
    echo "✅ All splits correctly sequential"
}
```

### 2. Split Lifecycle Tracking

#### Phase 1: Split Planning
```yaml
split_tracking:
  effort-002-database:
    original_branch: "project/phase1/wave2/effort-002-database"
    status: "SPLIT_PLANNED"
    split_count: 2  # Planned number of splits
    total_original_lines: 950
    split_reason: "will exceed limit with tests"
    planned_splits:
      - description: "Schema and migrations"
        estimated_lines: 475
      - description: "CRUD operations and queries"
        estimated_lines: 475
```

#### Phase 2: Split Execution
```yaml
split_tracking:
  effort-002-database:
    original_branch: "project/phase1/wave2/effort-002-database"
    status: "SPLIT_IN_PROGRESS"
    split_count: 2
    total_original_lines: 950
    split_branches:
      - branch: "project/phase1/wave2/effort-002-database--split-001"
        status: "COMPLETED"
        lines: 480
        description: "Schema and migrations"
      - branch: "project/phase1/wave2/effort-002-database--split-002"
        status: "ACTIVE"  # Currently being worked on
        lines: 0  # Will be updated when complete
        description: "CRUD operations and queries"
```

#### Phase 3: Split Completion
```yaml
split_tracking:
  effort-002-database:
    original_branch: "project/phase1/wave2/effort-002-database"
    status: "SPLIT_DEPRECATED"
    deprecated_branch: "project/phase1/wave2/effort-002-database-deprecated-split"
    split_count: 2
    # ... rest of tracking data
```

### 3. Integration Planning with Splits

#### 🔴🔴🔴 CRITICAL: Merge Ordering with Dependencies 🔴🔴🔴

**SUPREME LAW: Dependencies are at EFFORT level, not SPLIT level!**

When Effort E2 depends on Effort E1:
- E2 depends on **ALL** of E1 (not just part of it)
- ALL E1 splits must merge BEFORE E2 can merge
- Splits merge sequentially: split-001 → split-002 → split-003

**Example with Dependencies:**
```
E1 (3 splits) ← E2 depends on E1
E3 (2 splits) ← E4 depends on E3

Correct Merge Order:
1. E1-split-001
2. E1-split-002  
3. E1-split-003  # E1 now complete
4. E2            # Can merge now that E1 is complete
5. E3-split-001
6. E3-split-002  # E3 now complete
7. E4            # Can merge now that E3 is complete
```

#### Wave Integration Must Check Split Tracking
```bash
prepare_wave_integration() {
    local phase="$1"
    local wave="$2"
    
    echo "Collecting branches for wave integration..."
    
    # Get all efforts for this wave
    EFFORTS=$(yq ".waves.wave${wave}.efforts[]" orchestrator-state.yaml)
    
    BRANCHES_TO_MERGE=()
    
    for effort in $EFFORTS; do
        # Check split tracking first
        SPLIT_COUNT=$(yq ".split_tracking.\"${effort}\".split_count // 0" orchestrator-state.yaml)
        
        if [ "$SPLIT_COUNT" -gt 0 ]; then
            echo "✅ $effort was split into $SPLIT_COUNT branches"
            # Get all split branches IN ORDER
            for i in $(seq 1 $SPLIT_COUNT); do
                SPLIT_NUM=$(printf "%03d" $i)
                SPLIT_BRANCH=$(yq ".split_tracking.\"${effort}\".split_branches[$((i-1))].branch" orchestrator-state.yaml)
                BRANCHES_TO_MERGE+=("$SPLIT_BRANCH")
                echo "  Added split $i: $SPLIT_BRANCH"
            done
        else
            # Not split, use original branch
            BRANCH=$(yq ".efforts_completed.\"${effort}\".branch" orchestrator-state.yaml)
            BRANCHES_TO_MERGE+=("$BRANCH")
            echo "  Added original: $BRANCH"
        fi
    done
    
    echo "Total branches to merge: ${#BRANCHES_TO_MERGE[@]}"
    
    # CRITICAL: Respect dependencies when ordering merges!
    order_branches_by_dependencies "${BRANCHES_TO_MERGE[@]}"
}

# Ensure dependent efforts wait for ALL splits of dependencies
validate_dependency_completion() {
    local effort="$1"
    local dependencies=$(yq ".efforts.\"${effort}\".dependencies[]" orchestrator-state.yaml)
    
    for dep in $dependencies; do
        echo "Checking dependency $dep for $effort..."
        
        # Check if dependency has splits
        SPLIT_COUNT=$(yq ".split_tracking.\"${dep}\".split_count // 0" orchestrator-state.yaml)
        
        if [ "$SPLIT_COUNT" -gt 0 ]; then
            # ALL splits must be integrated
            for i in $(seq 1 $SPLIT_COUNT); do
                STATUS=$(yq ".split_tracking.\"${dep}\".split_branches[$((i-1))].status" orchestrator-state.yaml)
                if [ "$STATUS" != "INTEGRATED" ]; then
                    echo "❌ BLOCKED: $effort cannot merge!"
                    echo "   Dependency $dep split $i not integrated (status: $STATUS)"
                    return 1
                fi
            done
            echo "✅ All $SPLIT_COUNT splits of $dep are integrated"
        else
            # Check single effort integration
            STATUS=$(yq ".efforts_completed.\"${dep}\".integration_status" orchestrator-state.yaml)
            if [ "$STATUS" != "INTEGRATED" ]; then
                echo "❌ BLOCKED: $effort cannot merge!"
                echo "   Dependency $dep not integrated"
                return 1
            fi
        fi
    done
    
    echo "✅ All dependencies complete for $effort"
}
```

### 4. Split Branch Status Management

#### Valid Split Branch Statuses
```yaml
split_branch_statuses:
  ACTIVE:      # Currently being implemented
  COMPLETED:   # Implementation complete, pending review
  REVIEWED:    # Code review complete
  INTEGRATED:  # Merged into integration branch
  FAILED:      # Split failed or abandoned
```

#### Status Transition Rules
```
ACTIVE → COMPLETED → REVIEWED → INTEGRATED
   ↓         ↓          ↓
FAILED    FAILED     FAILED
```

### 5. Architect Split Verification

When Architect reviews efforts, MUST check split tracking:
```bash
verify_effort_compliance() {
    local effort="$1"
    
    # R297: Check split_count first
    SPLIT_DATA=$(yq ".split_tracking.\"${effort}\"" orchestrator-state.yaml)
    
    if [ "$SPLIT_DATA" != "null" ]; then
        SPLIT_COUNT=$(echo "$SPLIT_DATA" | yq ".split_count")
        STATUS=$(echo "$SPLIT_DATA" | yq ".status")
        
        echo "📊 Split Tracking for $effort:"
        echo "  Status: $STATUS"
        echo "  Split Count: $SPLIT_COUNT"
        
        if [ "$SPLIT_COUNT" -gt 0 ]; then
            echo "  ✅ COMPLIANT: Effort was properly split"
            
            # List all splits for documentation
            echo "$SPLIT_DATA" | yq ".split_branches[]" -o json | while read -r split; do
                BRANCH=$(echo "$split" | jq -r '.branch')
                LINES=$(echo "$split" | jq -r '.lines')
                STATUS=$(echo "$split" | jq -r '.status')
                echo "    - $BRANCH: $LINES lines ($STATUS)"
            done
            
            return 0
        fi
    fi
    
    # Not split, measure original
    return 1
}
```

### 6. Split Tracking Queries

#### Find All Current Splits for Integration
```bash
get_current_splits_for_integration() {
    yq '.split_tracking | to_entries | .[] | 
        select(.value.status == "SPLIT_DEPRECATED") | 
        .value.split_branches[] | 
        select(.status == "COMPLETED" or .status == "REVIEWED") | 
        .branch' orchestrator-state.yaml
}
```

#### Find Efforts Needing Split
```bash
find_efforts_needing_split() {
    yq '.split_tracking | to_entries | .[] | 
        select(.value.status == "SPLIT_PLANNED") | 
        .key' orchestrator-state.yaml
}
```

#### Get Split Summary for Wave
```bash
get_wave_split_summary() {
    local wave="$1"
    
    echo "Wave $wave Split Summary:"
    yq ".waves.wave${wave}.efforts[]" orchestrator-state.yaml | while read effort; do
        SPLIT_COUNT=$(yq ".split_tracking.\"${effort}\".split_count // 0" orchestrator-state.yaml)
        if [ "$SPLIT_COUNT" -gt 0 ]; then
            TOTAL_LINES=$(yq ".split_tracking.\"${effort}\".total_original_lines" orchestrator-state.yaml)
            echo "  $effort: Split into $SPLIT_COUNT branches (was $TOTAL_LINES lines)"
        fi
    done
}
```

### 7. Split Tracking Validation

#### Pre-Integration Validation
```bash
validate_splits_before_integration() {
    local wave="$1"
    
    echo "Validating all splits for wave $wave..."
    
    VALIDATION_ERRORS=0
    
    yq ".waves.wave${wave}.efforts[]" orchestrator-state.yaml | while read effort; do
        SPLIT_DATA=$(yq ".split_tracking.\"${effort}\"" orchestrator-state.yaml)
        
        if [ "$SPLIT_DATA" != "null" ]; then
            # Verify all splits are completed
            INCOMPLETE=$(echo "$SPLIT_DATA" | yq '.split_branches[] | select(.status != "COMPLETED" and .status != "REVIEWED" and .status != "INTEGRATED") | .branch')
            
            if [ -n "$INCOMPLETE" ]; then
                echo "❌ ERROR: $effort has incomplete splits:"
                echo "$INCOMPLETE"
                ((VALIDATION_ERRORS++))
            fi
            
            # Verify deprecated branch is marked
            STATUS=$(echo "$SPLIT_DATA" | yq '.status')
            if [ "$STATUS" != "SPLIT_DEPRECATED" ]; then
                echo "❌ ERROR: $effort splits complete but status not SPLIT_DEPRECATED"
                ((VALIDATION_ERRORS++))
            fi
        fi
    done
    
    if [ $VALIDATION_ERRORS -gt 0 ]; then
        echo "❌ BLOCKED: Fix $VALIDATION_ERRORS split tracking issues before integration"
        return 1
    fi
    
    echo "✅ All splits properly tracked and ready for integration"
    return 0
}
```

## Implementation Requirements

### For Orchestrator:
1. **Initialize split_tracking section** in state file
2. **Update tracking** when splits are planned
3. **Monitor split progress** during MONITORING state
4. **Update statuses** as splits complete
5. **Validate tracking** before integration
6. **Use split branches** (not originals) for integration

### For SW Engineer:
1. **Report split creation** to orchestrator
2. **Update branch status** when complete
3. **Provide line counts** for each split
4. **Mark original** as deprecated when done

### For Code Reviewer:
1. **Create split plan** with tracking structure
2. **Update review status** for each split
3. **Verify all splits** reviewed before integration

### For Architect:
1. **Check split_tracking** before measuring
2. **Verify split compliance** in reviews
3. **Document split status** in reports
4. **Validate deprecated marking**

### For Integration Agent:
1. **Read split_tracking** for branch selection
2. **Use split branches** not originals
3. **Verify all splits** complete
4. **Update integration_status** after merge

## Error Messages

### Missing Split Tracking
```
❌ CRITICAL: No split tracking found for effort-003
This effort was split but tracking is missing!
Action: Update split_tracking section in orchestrator-state.yaml
```

### Incomplete Split Tracking
```
❌ ERROR: Split tracking incomplete for effort-004
Missing: deprecated_branch, split_branches[1].lines
Action: Complete all required fields before integration
```

### Using Original Instead of Splits
```
❌ BLOCKED: Attempting to integrate original branch
Branch: project/phase1/wave1/effort-005
Status: SPLIT_DEPRECATED
Use these splits instead:
  - project/phase1/wave1/effort-005--split-001
  - project/phase1/wave1/effort-005--split-002
```

## Verification Steps

### Check Split Tracking Completeness
```bash
verify_split_tracking_complete() {
    ERRORS=0
    
    # Check all split_tracking entries have required fields
    yq '.split_tracking | to_entries | .[]' orchestrator-state.yaml -o json | while read -r entry; do
        EFFORT=$(echo "$entry" | jq -r '.key')
        
        # Required fields
        for field in "original_branch" "status" "split_count" "split_branches"; do
            VALUE=$(echo "$entry" | jq -r ".value.$field // \"missing\"")
            if [ "$VALUE" = "missing" ]; then
                echo "❌ $EFFORT missing required field: $field"
                ((ERRORS++))
            fi
        done
    done
    
    return $ERRORS
}
```

## Related Rules
- R296: Deprecated Branch Marking Protocol
- R297: Architect Split Detection Protocol
- R204: Orchestrator Split Infrastructure
- R014: Branch Naming Convention
- R260: Integration Agent Core Requirements

## Penalties
- Missing split tracking: -30%
- Incomplete tracking data: -20%
- Using wrong branches for integration: -40%
- Not updating split statuses: -15%
- Missing deprecated marking: -25%

---
*Rule Type*: Protocol
*Agents*: All Agents
*Enforcement*: State file validation and integration checks