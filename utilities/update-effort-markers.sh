#!/bin/bash

# Update effort-specific marker files (plans, reports) to new .software-factory/ structure

set -e

echo "📦 Updating effort marker files to new .software-factory/ structure..."
echo ""

# Function to update plan/report file references
update_effort_markers() {
    local file="$1"
    
    # Skip backup files
    if [[ "$file" == *".backup"* ]]; then
        return
    fi
    
    echo "  Processing: $file"
    
    # Create backup
    cp "$file" "${file}.backup.effort-markers" 2>/dev/null || true
    
    local changes=false
    
    # Update IMPLEMENTATION-PLAN.md creation/references
    if grep -q "IMPLEMENTATION-PLAN\.md" "$file"; then
        # For creation commands
        sed -i 's|echo.*> IMPLEMENTATION-PLAN\.md|PLAN_DIR=".software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/plans" \&\& mkdir -p "$PLAN_DIR" \&\& echo "..." > "$PLAN_DIR/IMPLEMENTATION-PLAN-$(date +%Y%m%d-%H%M%S).md"|g' "$file"
        
        # For checking existence
        sed -i 's|\[ -f "IMPLEMENTATION-PLAN\.md" \]|[ -n "$(ls .software-factory/phase*/wave*/${EFFORT_NAME}/plans/IMPLEMENTATION-PLAN-*.md 2>/dev/null)" ]|g' "$file"
        
        # For references in text
        sed -i 's|IMPLEMENTATION-PLAN\.md|.software-factory/phase\${PHASE}/wave\${WAVE}/\${EFFORT_NAME}/plans/IMPLEMENTATION-PLAN-*.md|g' "$file"
        
        echo "    ✅ Updated IMPLEMENTATION-PLAN.md references"
        changes=true
    fi
    
    # Update CODE-REVIEW-REPORT.md creation/references
    if grep -q "CODE-REVIEW-REPORT\.md" "$file"; then
        # For creation commands
        sed -i 's|echo.*> CODE-REVIEW-REPORT\.md|REPORT_DIR=".software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/reports" \&\& mkdir -p "$REPORT_DIR" \&\& echo "..." > "$REPORT_DIR/CODE-REVIEW-REPORT-$(date +%Y%m%d-%H%M%S).md"|g' "$file"
        
        # For checking existence
        sed -i 's|\[ -f "CODE-REVIEW-REPORT\.md" \]|[ -n "$(ls .software-factory/phase*/wave*/${EFFORT_NAME}/reports/CODE-REVIEW-REPORT-*.md 2>/dev/null)" ]|g' "$file"
        
        # For references in text
        sed -i 's|CODE-REVIEW-REPORT\.md|.software-factory/phase\${PHASE}/wave\${WAVE}/\${EFFORT_NAME}/reports/CODE-REVIEW-REPORT-*.md|g' "$file"
        
        echo "    ✅ Updated CODE-REVIEW-REPORT.md references"
        changes=true
    fi
    
    # Update SPLIT-PLAN.md creation/references
    if grep -q "SPLIT-PLAN\.md" "$file"; then
        # For creation commands
        sed -i 's|echo.*> SPLIT-PLAN\.md|PLAN_DIR=".software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/plans" \&\& mkdir -p "$PLAN_DIR" \&\& echo "..." > "$PLAN_DIR/SPLIT-PLAN-$(date +%Y%m%d-%H%M%S).md"|g' "$file"
        
        # For checking existence
        sed -i 's|\[ -f "SPLIT-PLAN\.md" \]|[ -n "$(ls .software-factory/phase*/wave*/${EFFORT_NAME}/plans/SPLIT-PLAN-*.md 2>/dev/null)" ]|g' "$file"
        
        # For references in text
        sed -i 's|SPLIT-PLAN\.md|.software-factory/phase\${PHASE}/wave\${WAVE}/\${EFFORT_NAME}/plans/SPLIT-PLAN-*.md|g' "$file"
        
        echo "    ✅ Updated SPLIT-PLAN.md references"
        changes=true
    fi
    
    # Update FIX-PLAN.md creation/references
    if grep -q "FIX-PLAN\.md" "$file"; then
        # For creation commands
        sed -i 's|echo.*> FIX-PLAN\.md|PLAN_DIR=".software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/plans" \&\& mkdir -p "$PLAN_DIR" \&\& echo "..." > "$PLAN_DIR/FIX-PLAN-$(date +%Y%m%d-%H%M%S).md"|g' "$file"
        
        # For checking existence
        sed -i 's|\[ -f "FIX-PLAN\.md" \]|[ -n "$(ls .software-factory/phase*/wave*/${EFFORT_NAME}/plans/FIX-PLAN-*.md 2>/dev/null)" ]|g' "$file"
        
        # For references in text
        sed -i 's|FIX-PLAN\.md|.software-factory/phase\${PHASE}/wave\${WAVE}/\${EFFORT_NAME}/plans/FIX-PLAN-*.md|g' "$file"
        
        echo "    ✅ Updated FIX-PLAN.md references"
        changes=true
    fi
    
    # Clean up backup if changes were made
    if [ "$changes" = true ]; then
        rm "${file}.backup.effort-markers"
    else
        mv "${file}.backup.effort-markers" "$file" 2>/dev/null || true
    fi
}

# Process all state rule files
echo "🔍 Processing state rule files..."
find /home/vscode/software-factory-template/agent-states -name "rules.md" -type f | while read -r file; do
    update_effort_markers "$file"
done

echo ""
echo "🔍 Processing agent configuration files..."
for file in /home/vscode/software-factory-template/.claude/agents/*.md; do
    if [ -f "$file" ]; then
        update_effort_markers "$file"
    fi
done

echo ""
echo "🔍 Processing rule library files..."
find /home/vscode/software-factory-template/rule-library -name "*.md" -type f | while read -r file; do
    update_effort_markers "$file"
done

echo ""
echo "✅ Effort marker updates complete!"
echo ""
echo "📋 New structure for effort markers:"
echo "  .software-factory/"
echo "  └── phase\${PHASE}/"
echo "      └── wave\${WAVE}/"
echo "          └── \${EFFORT_NAME}/"
echo "              ├── plans/"
echo "              │   ├── IMPLEMENTATION-PLAN-TIMESTAMP.md"
echo "              │   ├── SPLIT-PLAN-TIMESTAMP.md"
echo "              │   └── FIX-PLAN-TIMESTAMP.md"
echo "              └── reports/"
echo "                  └── CODE-REVIEW-REPORT-TIMESTAMP.md"