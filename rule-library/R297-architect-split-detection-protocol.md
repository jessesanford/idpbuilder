# 🚨🚨🚨 RULE R297 - Architect Split Detection Protocol

**Criticality:** BLOCKING - Must check splits BEFORE measuring integration  
**Grading Impact:** -50% for demanding re-split of already-split efforts  
**Enforcement:** EVERY wave/phase review MUST check split_count first

## Rule Statement

Architects MUST check the orchestrator-state.yaml `split_count` field BEFORE measuring any effort size. If an effort has `split_count > 0`, it was already split and is compliant regardless of integration branch size. PRs come from original effort branches, NOT integration branches.

## 🔴🔴🔴 CRITICAL: Integration Size Is Irrelevant! 🔴🔴🔴

### The Fundamental Truth
- **PRs come from**: Original effort branches in `/efforts/phase*/wave*/[effort-name]/`
- **PRs do NOT come from**: Integration branches (these merge all splits together)
- **Integration branches**: Are ONLY for testing merge compatibility
- **Split efforts**: Will ALWAYS exceed limits when integrated (by design!)

## Mandatory Split Detection Workflow

### 🚨 STEP 1: Check Split Count FIRST (MANDATORY)
```bash
check_effort_splits() {
    local effort_name="$1"
    local phase="$2"
    local wave="$3"
    
    echo "🔍 Checking if $effort_name was already split..."
    
    # Check orchestrator-state.yaml for split_count
    SPLIT_COUNT=$(yq ".efforts_completed.\"${effort_name}\".split_count" orchestrator-state.yaml 2>/dev/null || echo "0")
    
    if [ "$SPLIT_COUNT" -gt 0 ]; then
        echo "✅ $effort_name was ALREADY SPLIT into $SPLIT_COUNT parts"
        echo "📦 Original effort branches are compliant"
        echo "🎯 Integration size is IRRELEVANT (PRs come from effort branches)"
        echo "✅ MARKING AS COMPLIANT - No further size check needed"
        return 0  # Compliant, skip size measurement
    fi
    
    echo "📊 $effort_name was not split, proceeding to measure original branch..."
    return 1  # Continue to size measurement
}
```

### 🚨 STEP 2: Measure ORIGINAL Effort Branch (NOT Integration)
```bash
measure_original_effort_size() {
    local effort_name="$1"
    local phase="$2"
    local wave="$3"
    
    # Navigate to ORIGINAL effort directory
    EFFORT_DIR="/efforts/phase${phase}/wave${wave}/${effort_name}"
    
    if [ ! -d "$EFFORT_DIR" ]; then
        echo "❌ ERROR: Effort directory not found: $EFFORT_DIR"
        return 1
    fi
    
    cd "$EFFORT_DIR"
    echo "📍 Measuring ORIGINAL effort branch at: $(pwd)"
    
    # Find project root
    PROJECT_ROOT=$(pwd)
    while [ "$PROJECT_ROOT" != "/" ]; do
        [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break
        PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
    done
    
    # Measure using line counter
    SIZE=$($PROJECT_ROOT/tools/line-counter.sh | grep "Total" | awk '{print $NF}')
    
    echo "📏 Original effort size: $SIZE lines"
    
    if [ "$SIZE" -le 800 ]; then
        echo "✅ COMPLIANT: $SIZE lines (under 800 limit)"
        return 0
    else
        echo "❌ VIOLATION: $SIZE lines (exceeds 800 limit)"
        echo "🚨 This effort needs splitting!"
        return 1
    fi
}
```

## Complete Architect Review Protocol

### The RIGHT Way to Review Efforts
```bash
review_wave_efforts() {
    local phase="$1"
    local wave="$2"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🏗️ ARCHITECT WAVE REVIEW - Phase $phase Wave $wave"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Get list of efforts from state file
    EFFORTS=$(yq ".waves.wave${wave}.efforts[]" orchestrator-state.yaml)
    
    for effort in $EFFORTS; do
        echo ""
        echo "──────────────────────────────────────────────────────────"
        echo "📋 Reviewing effort: $effort"
        
        # STEP 1: Check if already split (R297)
        if check_effort_splits "$effort" "$phase" "$wave"; then
            echo "✅ $effort is COMPLIANT (already split)"
            continue  # Skip to next effort
        fi
        
        # STEP 2: Measure original branch (R022)
        if measure_original_effort_size "$effort" "$phase" "$wave"; then
            echo "✅ $effort is COMPLIANT (size within limit)"
        else
            echo "❌ $effort NEEDS SPLITTING"
            # Document in review report
        fi
    done
    
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "REVIEW COMPLETE"
    echo "═══════════════════════════════════════════════════════════════"
}
```

## ❌ FORBIDDEN: What Architects Must NEVER Do

### NEVER Measure Integration Branch for Compliance
```bash
# ❌❌❌ CATASTROPHICALLY WRONG ❌❌❌
cd integration-workspace
git checkout phase2-wave2-integration
line-counter.sh  # Shows 1500 lines
echo "VIOLATION! This needs splitting!"  # WRONG! Integration merges all splits!
```

### NEVER Ignore split_count in State File
```bash
# ❌❌❌ WRONG - Ignoring split_count ❌❌❌
# State file shows: split_count: 2
# But architect still demands another split
echo "This effort is 900 lines, needs splitting!"  # WRONG! Already split!
```

### NEVER Demand Re-splitting of Split Efforts
```bash
# ❌❌❌ WRONG - Demanding re-split ❌❌❌
echo "E1.1.2 shows 904 lines in integration"
echo "This needs to be split again!"  # WRONG! Check split_count first!
```

## ✅ CORRECT: What Architects Must ALWAYS Do

### ALWAYS Check Split Count First
```bash
# ✅ CORRECT - Check split_count before anything else
SPLIT_COUNT=$(yq '.efforts_completed."E1.1.2".split_count' orchestrator-state.yaml)
if [ "$SPLIT_COUNT" -gt 0 ]; then
    echo "✅ E1.1.2 already split into $SPLIT_COUNT parts - COMPLIANT"
fi
```

### ALWAYS Measure Original Effort Branches
```bash
# ✅ CORRECT - Measure where PRs come from
cd /efforts/phase1/wave1/E1.1.2
$PROJECT_ROOT/tools/line-counter.sh  # Measure original branch
```

### ALWAYS Document Split Status in Review
```bash
# ✅ CORRECT - Clear documentation
echo "## Size Compliance Assessment"
echo "- E1.1.1: 650 lines (compliant)"
echo "- E1.1.2: Already split into 2 parts (compliant)"
echo "- E1.1.3: 450 lines (compliant)"
echo "✅ All efforts size-compliant"
```

## Integration Branch Clarification

### Why Integration Branches Exceed Limits
```yaml
# Example: E1.1.2 was split into 2 parts
split_1: 450 lines  # Under limit ✅
split_2: 454 lines  # Under limit ✅
integration: 904 lines  # Over limit but EXPECTED! ✅

# The integration branch merges BOTH splits
# This is NORMAL and EXPECTED behavior
# PRs will come from split_1 and split_2, NOT integration
```

### Integration Branches Are For Testing Only
```bash
Purpose of Integration Branches:
✅ Test that all splits work together
✅ Verify no merge conflicts
✅ Run integration tests
✅ Check system stability

NOT for:
❌ Creating PRs (PRs come from effort branches)
❌ Measuring size compliance (measure effort branches)
❌ Determining if splitting is needed (check split_count)
```

## Grading Impact

### Violations and Penalties
- **Ignoring split_count**: -50% (Major review failure)
- **Measuring integration instead of effort**: -30% (Wrong measurement)
- **Demanding re-split of split efforts**: -50% (Critical error)
- **Blocking compliant efforts**: -40% (Process disruption)

### Perfect Review Requirements
```yaml
perfect_review:
  - Always check split_count first
  - Measure original effort branches only
  - Never measure integration for compliance
  - Document split status clearly
  - Understand PRs come from effort branches
```

## Quick Reference Card

```bash
# For EVERY effort in review:
1. CHECK: split_count in orchestrator-state.yaml
   - If > 0: Mark COMPLIANT, skip size check
   - If 0: Continue to step 2

2. MEASURE: Original effort branch
   - cd /efforts/phase*/wave*/[effort-name]
   - Run line-counter.sh
   - Check if ≤ 800 lines

3. NEVER: Measure integration branch for compliance
   - Integration merges all splits (expected to exceed)
   - PRs come from effort branches, not integration

4. DOCUMENT: Clear status in review report
   - Note which efforts were pre-split
   - Show original branch measurements
   - Mark all compliant efforts
```

## The Split Detection Law

```
Check split_count FIRST - ALWAYS
Measure effort branches - NEVER integration
Already-split = compliant - NO EXCEPTIONS
Integration size is irrelevant - PRs from efforts
Document split status - CLEARLY
```

---
**Remember:** Integration branches merge all splits together. They WILL exceed limits. This is EXPECTED. PRs come from the original effort branches, which must each be under the limit. Check split_count FIRST!