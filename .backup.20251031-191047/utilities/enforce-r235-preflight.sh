#!/bin/bash
# 🔴🔴🔴 R235: MANDATORY PRE-FLIGHT VERIFICATION ENFORCEMENT SCRIPT 🔴🔴🔴
# 
# This script MUST be sourced by ALL agents at startup
# Violation of pre-flight checks = -100% GRADE = AUTOMATIC FAILURE
#
# Usage: source /path/to/enforce-r235-preflight.sh <agent-type>

set -euo pipefail

# Get agent type from parameter
AGENT_TYPE="${1:-unknown}"

echo "════════════════════════════════════════════════════════════════"
echo "🔴🔴🔴 R235: SUPREME LAW #3 - MANDATORY PRE-FLIGHT VERIFICATION 🔴🔴🔴"
echo "════════════════════════════════════════════════════════════════"
echo "AGENT: $AGENT_TYPE"
echo "TIMESTAMP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "PID: $$"
echo "════════════════════════════════════════════════════════════════"

# Track violations
VIOLATIONS=0
VIOLATION_MESSAGES=()

# Function to record violation
record_violation() {
    local message="$1"
    VIOLATIONS=$((VIOLATIONS + 1))
    VIOLATION_MESSAGES+=("$message")
    echo "❌❌❌ VIOLATION #$VIOLATIONS: $message"
}

# CHECK 1: VERIFY CORRECT WORKING DIRECTORY
echo ""
echo "CHECK 1: Verifying working directory..."
CURRENT_DIR=$(pwd)
echo "Current directory: $CURRENT_DIR"

# Check if in planning repository (FATAL ERROR)
if [[ "$CURRENT_DIR" == *"software-factory"* ]] && [[ "$CURRENT_DIR" != *"/efforts/"* ]]; then
    record_violation "In planning repository, not target repository!"
    echo "   This indicates orchestrator spawned agent in wrong location"
    echo "   Or agent was manually started in planning repo"
fi

# For implementation agents, must be in effort directory
if [[ "$AGENT_TYPE" == "sw-engineer" ]] || [[ "$AGENT_TYPE" == "code-reviewer" ]]; then
    if [[ "$CURRENT_DIR" != *"/efforts/phase"*"/wave"*"/"* ]]; then
        record_violation "Not in effort directory! Expected: */efforts/phase*/wave*/[effort-name]"
        echo "   Actual: $CURRENT_DIR"
    else
        echo "✅ In effort directory structure"
    fi
fi

# For orchestrator, can be in project root or planning repo root
if [[ "$AGENT_TYPE" == "orchestrator" ]]; then
    if [[ "$CURRENT_DIR" == *"/efforts/"* ]]; then
        record_violation "Orchestrator should not be in effort directory!"
    else
        echo "✅ Orchestrator in appropriate directory"
    fi
fi

# CHECK 2: VERIFY GIT REPOSITORY
echo ""
echo "CHECK 2: Verifying git repository..."
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    record_violation "Not in a git repository!"
    echo "   Cannot perform any version control operations"
else
    echo "✅ Git repository exists"
    
    # Verify correct remote (NOT planning repo remote!)
    REMOTE_URL=$(git remote get-url origin 2>/dev/null || echo "NO_REMOTE")
    echo "Remote URL: $REMOTE_URL"
    
    if [[ "$REMOTE_URL" == "NO_REMOTE" ]]; then
        record_violation "No git remote configured!"
    elif [[ "$REMOTE_URL" == *"software-factory"* ]]; then
        record_violation "Remote points to planning repository! Should point to target project!"
        echo "   This means wrong repository was cloned or remote misconfigured"
    else
        echo "✅ Remote appears to be target repository"
    fi
fi

# CHECK 3: VERIFY GIT BRANCH
echo ""
echo "CHECK 3: Verifying git branch..."
if git rev-parse --git-dir > /dev/null 2>&1; then
    CURRENT_BRANCH=$(git branch --show-current)
    echo "Current branch: $CURRENT_BRANCH"
    
    # For implementation agents, branch must match effort pattern
    if [[ "$AGENT_TYPE" == "sw-engineer" ]] || [[ "$AGENT_TYPE" == "code-reviewer" ]]; then
        EFFORT_NAME=$(basename "$CURRENT_DIR")
        if [[ "$CURRENT_BRANCH" != *"$EFFORT_NAME"* ]]; then
            record_violation "Branch doesn't match effort! Expected effort '$EFFORT_NAME' in branch name"
            echo "   Actual branch: $CURRENT_BRANCH"
        else
            echo "✅ Branch matches effort name"
        fi
        
        # Check branch pattern
        if ! [[ "$CURRENT_BRANCH" =~ phase[0-9]+/wave[0-9]+/ ]]; then
            record_violation "Branch doesn't follow phase/wave pattern!"
        else
            echo "✅ Branch follows correct pattern"
        fi
    fi
else
    echo "⚠️ Skipping branch check (no git repository)"
fi

# CHECK 4: VERIFY WORKSPACE ISOLATION
echo ""
echo "CHECK 4: Verifying workspace isolation..."
if [[ "$AGENT_TYPE" == "sw-engineer" ]]; then
    # Must have pkg/ directory for code
    if [ ! -d "pkg" ]; then
        echo "⚠️ No pkg directory found - will be created on first file"
    else
        echo "✅ pkg directory exists for isolated code"
    fi
    
    # Must NOT be in main /pkg directory
    if [[ "$CURRENT_DIR" == */pkg ]] && [[ "$CURRENT_DIR" != *"/efforts/"* ]]; then
        record_violation "In main /pkg directory! This violates workspace isolation!"
    fi
    
    # Check for IMPLEMENTATION-PLAN.md
    if [ ! -f "IMPLEMENTATION-PLAN.md" ]; then
        record_violation "No IMPLEMENTATION-PLAN.md found! Orchestrator failed to set up infrastructure!"
    else
        echo "✅ Implementation plan exists"
        
        # Check for R209 metadata
        if ! grep -q "EFFORT INFRASTRUCTURE METADATA" IMPLEMENTATION-PLAN.md 2>/dev/null; then
            record_violation "Implementation plan missing R209 metadata! Orchestrator violation!"
        else
            echo "✅ R209 metadata present in plan"
        fi
    fi
fi

# CHECK 5: VERIFY NO CONTAMINATION
echo ""
echo "CHECK 5: Verifying no contamination..."
if [[ "$AGENT_TYPE" == "sw-engineer" ]] || [[ "$AGENT_TYPE" == "code-reviewer" ]]; then
    # Check for massive file counts indicating contamination
    if [ -d "." ]; then
        # Count all files (excluding .git)
        FILE_COUNT=$(find . -type f -not -path "./.git/*" 2>/dev/null | wc -l)
        echo "Total files in workspace: $FILE_COUNT"
        
        if [ "$FILE_COUNT" -gt 10000 ]; then
            record_violation "MASSIVE CONTAMINATION! Found $FILE_COUNT files!"
            echo "   This workspace is contaminated from other sources"
            
            # Try to identify contamination source
            if [ -d "vendor" ] || [ -d "node_modules" ]; then
                echo "   Detected dependency directories (vendor/node_modules)"
            fi
            if find . -name "*.exe" -o -name "*.dll" 2>/dev/null | head -1 | grep -q .; then
                echo "   Detected binary files"
            fi
        elif [ "$FILE_COUNT" -gt 1000 ]; then
            echo "⚠️ WARNING: High file count ($FILE_COUNT) - possible contamination"
        else
            echo "✅ File count reasonable ($FILE_COUNT files)"
        fi
        
        # Specific check for pkg directory
        if [ -d "pkg" ]; then
            PKG_COUNT=$(find pkg -type f 2>/dev/null | wc -l)
            echo "Files in pkg/: $PKG_COUNT"
            if [ "$PKG_COUNT" -gt 500 ]; then
                record_violation "pkg/ directory has $PKG_COUNT files - likely contaminated!"
            fi
        fi
    fi
fi

# FINAL VERDICT
echo ""
echo "════════════════════════════════════════════════════════════════"
if [ "$VIOLATIONS" -eq 0 ]; then
    echo "✅✅✅ ALL PRE-FLIGHT CHECKS PASSED - SAFE TO PROCEED ✅✅✅"
    echo "════════════════════════════════════════════════════════════════"
    echo ""
    echo "Agent $AGENT_TYPE may begin work in:"
    echo "  Directory: $CURRENT_DIR"
    if git rev-parse --git-dir > /dev/null 2>&1; then
        echo "  Branch: $(git branch --show-current)"
        echo "  Remote: $(git remote get-url origin 2>/dev/null || echo 'none')"
    fi
    exit 0
else
    echo "❌❌❌ PRE-FLIGHT VERIFICATION FAILED ❌❌❌"
    echo "════════════════════════════════════════════════════════════════"
    echo ""
    echo "TOTAL VIOLATIONS: $VIOLATIONS"
    echo ""
    echo "VIOLATION SUMMARY:"
    for i in "${!VIOLATION_MESSAGES[@]}"; do
        echo "$((i+1)). ${VIOLATION_MESSAGES[$i]}"
    done
    echo ""
    echo "🔴🔴🔴 R235 SUPREME LAW #3 VIOLATED 🔴🔴🔴"
    echo "REFUSING TO WORK - WRONG LOCATION/CONFIGURATION"
    echo "PENALTY: -100% GRADE (AUTOMATIC FAILURE)"
    echo ""
    echo "REQUIRED ACTIONS:"
    echo "1. Orchestrator must set up proper infrastructure"
    echo "2. Clone target repository (not planning repo)"
    echo "3. Create effort directories with isolation"
    echo "4. Initialize git with correct remote"
    echo "5. Create branches following naming convention"
    echo "6. Inject R209 metadata into implementation plans"
    echo ""
    echo "EXITING WITH CODE 235"
    exit 235
fi