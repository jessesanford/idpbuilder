#!/bin/bash

# validate-split-boundaries.sh
# Validates that split plans properly reference splits from the SAME effort only

set -e

SPLIT_PLAN="${1:-SPLIT-PLAN-*.md}"
SPLIT_NUM="${2:-}"

echo "=========================================="
echo "Split Boundary Validation (R207)"
echo "=========================================="

# Detect current effort from path
CURRENT_PHASE=$(pwd | grep -oP 'phase\d+' || echo "unknown")
CURRENT_WAVE=$(pwd | grep -oP 'wave\d+' || echo "unknown")
CURRENT_EFFORT=$(pwd | grep -oP 'effort-[^/]+' || echo "unknown")

if [ "$CURRENT_EFFORT" = "unknown" ]; then
    echo "⚠️ WARNING: Not in an effort directory, cannot validate context"
else
    echo "Current Context:"
    echo "  Phase: $CURRENT_PHASE"
    echo "  Wave: $CURRENT_WAVE"
    echo "  Effort: $CURRENT_EFFORT"
fi

echo ""
echo "Checking split plans..."

# Find and validate all split plans
for plan in $SPLIT_PLAN; do
    if [ ! -f "$plan" ]; then
        continue
    fi
    
    echo ""
    echo "Validating: $plan"
    echo "----------------------------------------"
    
    # Extract split number from filename
    PLAN_NUM=$(echo "$plan" | grep -oP '\d+' | head -1)
    
    # Check for parent effort declaration
    PARENT_EFFORT=$(grep "\*\*Parent Effort\*\*:" "$plan" | head -1 | cut -d: -f2 | xargs)
    echo "  Parent Effort: $PARENT_EFFORT"
    
    # Check previous split reference
    if [ "$PLAN_NUM" -gt 1 ] || ! grep -q "first split" "$plan"; then
        PREV_NUM=$((PLAN_NUM - 1))
        PREV_REF=$(grep "\*\*Previous Split\*\*:\|^- \*\*Previous Split\*\*:" "$plan" | head -1 || echo "Not found")
        
        echo "  Previous Split Reference: $PREV_REF"
        
        # Validate it references the same effort
        if [ "$CURRENT_EFFORT" != "unknown" ]; then
            if echo "$PREV_REF" | grep -q "$CURRENT_EFFORT"; then
                echo "  ✅ Previous split correctly references $CURRENT_EFFORT"
            elif echo "$PREV_REF" | grep -q "None.*first"; then
                echo "  ⚠️ Claims to be first split but numbered $PLAN_NUM"
            else
                echo "  ❌ FATAL: Previous split references different effort!"
                echo "     Expected reference to: $CURRENT_EFFORT"
                echo "     Found: $PREV_REF"
                
                # Try to detect what effort it's referencing
                WRONG_EFFORT=$(echo "$PREV_REF" | grep -oP '\(.*?\)' | tr -d '()')
                if [ -n "$WRONG_EFFORT" ]; then
                    echo "  ❌ Incorrectly references: $WRONG_EFFORT"
                    echo "  ❌ This is CROSS-POLLINATION between efforts!"
                fi
                exit 1
            fi
        fi
        
        # Check for full path
        if grep -q "Path:.*efforts.*$CURRENT_EFFORT/split-" "$plan"; then
            echo "  ✅ Full path information present"
        else
            echo "  ⚠️ Missing or incomplete path information"
        fi
        
        # Check for branch information
        if grep -q "Branch:.*$CURRENT_EFFORT-split-" "$plan"; then
            echo "  ✅ Branch naming consistent with effort"
        else
            echo "  ⚠️ Branch naming may not match effort"
        fi
    else
        echo "  ℹ️ Split 001 - no previous split expected"
        
        # Verify it says "None" or "first split"
        if grep -q "Previous Split:.*None\|Previous Split:.*first" "$plan"; then
            echo "  ✅ Correctly identifies as first split"
        else
            echo "  ⚠️ Split 001 should indicate no previous split"
        fi
    fi
    
    # Check for boundaries section
    if grep -q "## Boundaries" "$plan"; then
        echo "  ✅ Boundaries section present"
        
        # Check for critical warning about same effort
        if grep -q "CRITICAL.*same effort\|MUST reference SAME effort" "$plan"; then
            echo "  ✅ Contains warning about same-effort requirement"
        else
            echo "  ⚠️ Missing critical warning about same-effort splits"
        fi
    else
        echo "  ❌ Missing Boundaries section!"
    fi
    
    # Summary for this plan
    echo "  ----------------------------------------"
    echo "  Summary: Plan structure "
    if grep -q "$CURRENT_EFFORT" "$plan" || [ "$CURRENT_EFFORT" = "unknown" ]; then
        echo "  ✅ Split plan appears valid"
    else
        echo "  ❌ Split plan has boundary issues"
    fi
done

echo ""
echo "=========================================="
echo "Validation Complete"
echo "=========================================="

# Check for cross-pollination patterns
echo ""
echo "Checking for common cross-pollination patterns..."

# Look for references to other common efforts
COMMON_EFFORTS="oci-types registry-auth-types api-client webhook-framework stack-types"
for effort in $COMMON_EFFORTS; do
    if [ "$effort" != "${CURRENT_EFFORT#effort-}" ]; then
        if grep -q "$effort" $SPLIT_PLAN 2>/dev/null; then
            echo "⚠️ Found reference to different effort: $effort"
            grep -n "$effort" $SPLIT_PLAN
        fi
    fi
done

echo ""
echo "Remember: Each effort's splits must be completely isolated!"
echo "Never reference splits from different efforts, even if related."