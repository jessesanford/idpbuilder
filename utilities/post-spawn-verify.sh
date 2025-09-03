#!/bin/bash

# 🔍 POST-SPAWN WORKSPACE VERIFICATION
# Part of Software Factory 2.0 - RULES R181-R185, R176-R180
#
# Purpose: Verify agents are working in correct effort directories after spawn
# Usage: ./post-spawn-verify.sh [phase] [wave]
#
# This script MUST be run by orchestrator AFTER spawning agents
# to ensure workspace isolation compliance

set -e

echo "═══════════════════════════════════════════════════════════════"
echo "🔍 POST-SPAWN WORKSPACE VERIFICATION"
echo "═══════════════════════════════════════════════════════════════"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "═══════════════════════════════════════════════════════════════"

# Parse arguments
PHASE="${1}"
WAVE="${2}"

# If no arguments, scan all efforts
if [ -z "$PHASE" ] || [ -z "$WAVE" ]; then
    echo "Scanning all effort directories..."
    EFFORT_PATTERN="efforts/phase*/wave*/*"
else
    echo "Scanning Phase ${PHASE}, Wave ${WAVE} efforts..."
    EFFORT_PATTERN="efforts/phase${PHASE}/wave${WAVE}/*"
fi

# Track verification results
TOTAL_EFFORTS=0
PASSED_EFFORTS=0
FAILED_EFFORTS=0
CRITICAL_FAILURES=0

# Function to verify individual effort workspace
verify_effort_workspace() {
    local effort_dir="$1"
    local effort_name=$(basename "$effort_dir")
    local checks_passed=true
    local critical_fail=false
    
    echo ""
    echo "─────────────────────────────────────────────────────────────"
    echo "Verifying: $effort_dir"
    echo "─────────────────────────────────────────────────────────────"
    
    # Check 1: Directory exists
    if [ ! -d "$effort_dir" ]; then
        echo "❌ FAIL: Directory does not exist"
        return 1
    fi
    echo "✅ Directory exists"
    
    # Check 2: Git repository present (R182)
    if [ ! -d "$effort_dir/.git" ]; then
        echo "❌ CRITICAL: No git repository found"
        echo "   Orchestrator failed to create sparse clone"
        echo "   This violates R182 - Sparse Clone Requirement"
        critical_fail=true
        checks_passed=false
    else
        echo "✅ Git repository present"
    fi
    
    # Check 3: Correct branch naming (R184)
    if [ -d "$effort_dir/.git" ]; then
        cd "$effort_dir"
        local current_branch=$(git branch --show-current 2>/dev/null || echo "NONE")
        cd - > /dev/null
        
        # Extract phase/wave from path
        local path_phase=$(echo "$effort_dir" | sed -n 's/.*phase\([0-9]*\).*/\1/p')
        local path_wave=$(echo "$effort_dir" | sed -n 's/.*wave\([0-9]*\).*/\1/p')
        local expected_pattern="phase${path_phase}/wave${path_wave}/effort-${effort_name}"
        
        if [[ "$current_branch" == "$expected_pattern" ]]; then
            echo "✅ Branch naming correct: $current_branch"
        else
            echo "❌ CRITICAL: Branch naming violation"
            echo "   Expected: $expected_pattern"
            echo "   Actual: $current_branch"
            echo "   This violates R184 - Effort Branch Naming Scheme"
            critical_fail=true
            checks_passed=false
        fi
    fi
    
    # Check 4: Required files exist (R185)
    if [ ! -f "$effort_dir/IMPLEMENTATION-PLAN.md" ]; then
        echo "⚠️ WARNING: No IMPLEMENTATION-PLAN.md"
        echo "   Code Reviewer should create this"
        checks_passed=false
    else
        echo "✅ IMPLEMENTATION-PLAN.md exists"
    fi
    
    if [ ! -f "$effort_dir/work-log.md" ]; then
        echo "⚠️ WARNING: No work-log.md"
        checks_passed=false
    else
        echo "✅ work-log.md exists"
    fi
    
    # Check 5: pkg directory exists (R176)
    if [ ! -d "$effort_dir/pkg" ]; then
        echo "⚠️ WARNING: No pkg directory"
        echo "   SW Engineer must create code in effort/pkg/"
    else
        echo "✅ pkg directory exists"
        
        # Check if code is being created in the right place
        if [ -n "$(ls -A $effort_dir/pkg 2>/dev/null)" ]; then
            echo "✅ Code is being created in isolated pkg/"
        else
            echo "ℹ️ pkg directory empty (implementation not started)"
        fi
    fi
    
    # Check 6: Sparse checkout configuration (R182)
    if [ -f "$effort_dir/.git/info/sparse-checkout" ]; then
        echo "✅ Sparse checkout configured"
        echo "   Sparse paths:"
        head -5 "$effort_dir/.git/info/sparse-checkout" | sed 's/^/     /'
    elif [ -d "$effort_dir/.git" ]; then
        echo "⚠️ WARNING: No sparse checkout (full clone)"
        echo "   This may violate R182 - Sparse Clone Requirement"
    fi
    
    # Check 7: Remote configured properly
    if [ -d "$effort_dir/.git" ]; then
        cd "$effort_dir"
        if git remote -v 2>/dev/null | grep -q origin; then
            echo "✅ Remote origin configured"
        else
            echo "❌ CRITICAL: No remote configured"
            echo "   Agents cannot push their work"
            critical_fail=true
            checks_passed=false
        fi
        cd - > /dev/null
    fi
    
    # Check 8: No code in main /pkg (R176, R177)
    # This is checked globally, not per effort
    
    # Summary for this effort
    if [ "$critical_fail" = true ]; then
        echo ""
        echo "🚨 CRITICAL FAILURE: Workspace setup violations detected"
        echo "   Orchestrator must fix before agents can work"
        return 2
    elif [ "$checks_passed" = false ]; then
        echo ""
        echo "⚠️ WARNINGS: Some issues detected but not critical"
        return 1
    else
        echo ""
        echo "✅ ALL CHECKS PASSED: Workspace properly configured"
        return 0
    fi
}

# Main verification loop
for effort_dir in $EFFORT_PATTERN; do
    # Skip if not a directory
    [ ! -d "$effort_dir" ] && continue
    
    # Skip the pattern itself if no matches
    [[ "$effort_dir" == "$EFFORT_PATTERN" ]] && continue
    
    TOTAL_EFFORTS=$((TOTAL_EFFORTS + 1))
    
    if verify_effort_workspace "$effort_dir"; then
        PASSED_EFFORTS=$((PASSED_EFFORTS + 1))
    else
        FAILED_EFFORTS=$((FAILED_EFFORTS + 1))
        if [ $? -eq 2 ]; then
            CRITICAL_FAILURES=$((CRITICAL_FAILURES + 1))
        fi
    fi
done

# Global check: Ensure main /pkg is not being used
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "GLOBAL WORKSPACE ISOLATION CHECK"
echo "═══════════════════════════════════════════════════════════════"

if [ -d "./pkg" ] && [ -n "$(ls -A ./pkg 2>/dev/null)" ]; then
    echo "🚨🚨🚨 CRITICAL VIOLATION DETECTED 🚨🚨🚨"
    echo "Code found in main /pkg directory!"
    echo "This violates R176 - Workspace Isolation Requirement"
    echo ""
    echo "Files in main /pkg:"
    ls -la ./pkg | head -10
    echo ""
    echo "IMMEDIATE ACTION REQUIRED:"
    echo "1. Stop all agents immediately"
    echo "2. Move code to appropriate effort directories"
    echo "3. Re-spawn agents with correct working directories"
    echo ""
    echo "This is an AUTOMATIC GRADING FAILURE (20% lost)"
    CRITICAL_FAILURES=$((CRITICAL_FAILURES + 1))
else
    echo "✅ Main /pkg directory clean (no isolation violation)"
fi

# Final summary
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "VERIFICATION SUMMARY"
echo "═══════════════════════════════════════════════════════════════"
echo "Total Efforts Checked: $TOTAL_EFFORTS"
echo "Passed: $PASSED_EFFORTS"
echo "Failed: $FAILED_EFFORTS"
echo "Critical Failures: $CRITICAL_FAILURES"
echo ""

if [ $CRITICAL_FAILURES -gt 0 ]; then
    echo "🚨 CRITICAL FAILURES DETECTED"
    echo "Orchestrator MUST fix these issues before proceeding:"
    echo "1. Ensure all efforts have sparse git clones (R182)"
    echo "2. Verify branch naming follows pattern (R184)"
    echo "3. Check no code in main /pkg (R176)"
    echo "4. Re-spawn agents with correct working directories"
    echo ""
    echo "GRADING IMPACT: -20% for workspace isolation failure"
    exit 2
elif [ $FAILED_EFFORTS -gt 0 ]; then
    echo "⚠️ WARNINGS DETECTED"
    echo "Some non-critical issues found."
    echo "Orchestrator should address these soon."
    exit 1
else
    echo "✅ ALL VERIFICATIONS PASSED"
    echo "Workspace isolation properly maintained."
    echo "Agents are working in correct directories."
    echo "Orchestrator compliance verified."
    exit 0
fi