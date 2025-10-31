#!/bin/bash
# Fix orchestrator states by removing R304 references
# R304 (mandatory line counter) conflicts with R319 (orchestrator never measures code)

set -euo pipefail

echo "🔧 Fixing orchestrator state rules - removing R304 references"
echo "R304 conflicts with R319: Orchestrators must NEVER measure code"
echo ""

# Find all orchestrator state files with R304
FILES=$(find /home/vscode/software-factory-template/agent-states/orchestrator -name "rules.md" -exec grep -l "R304" {} \; 2>/dev/null)

COUNT=$(echo "$FILES" | wc -l)
echo "Found $COUNT files with R304 references to fix"
echo ""

FIXED=0
for file in $FILES; do
    echo "Fixing: $(basename $(dirname "$file"))"
    
    # Remove the entire R304 section (typically 7-8 lines)
    # Pattern: From "### R304:" to the next "###" or blank line after requirements
    sed -i '/^### R304: Mandatory Line Counter Tool Enforcement/,/^### \|^$/{ 
        /^### R304:/d
        /^\*\*File\*\*:.*R304/d
        /^\*\*Criticality\*\*:.*Manual counting/d
        /^$/d
        /^\*\*ABSOLUTE REQUIREMENTS/d
        /^- ✅ MUST use.*line-counter\.sh/d
        /^- ❌ NEVER use `wc -l`/d
        /^- ❌ NEVER count lines any other way/d
        /^$/d
    }' "$file"
    
    # Also remove any standalone R304 references in lists
    sed -i '/^[0-9]\. \*\*.*R304\*\* - MANDATORY LINE COUNTER/,+3d' "$file"
    
    # Remove R304 from inline references  
    sed -i 's/, R304//g' "$file"
    sed -i 's/R304, //g' "$file"
    
    FIXED=$((FIXED + 1))
done

echo ""
echo "✅ Fixed $FIXED orchestrator state files"
echo ""
echo "Summary of changes:"
echo "- Removed R304 sections from all orchestrator states"
echo "- R319 (orchestrator never measures) remains enforced"
echo "- Code Reviewers handle ALL size measurements"