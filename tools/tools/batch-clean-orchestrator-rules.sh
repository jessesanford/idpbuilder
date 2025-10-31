#!/bin/bash

# Batch cleanup script for orchestrator state rule files
# This script removes inline rule duplications from remaining files

ORCHESTRATOR_DIR="/home/vscode/software-factory-template/agent-states/orchestrator"

echo "🔄 Starting batch cleanup of remaining orchestrator state rule files..."

# Function to clean a state rules file
clean_state_file() {
    local STATE="$1"
    local FILE="$ORCHESTRATOR_DIR/$STATE/rules.md"
    
    if [ ! -f "$FILE" ]; then
        echo "⚠️ File not found: $FILE"
        return 1
    fi
    
    # Create backup
    cp "$FILE" "${FILE}.backup.$(date +%Y%m%d-%H%M%S)"
    
    # Count inline blocks before
    BEFORE_COUNT=$(grep -c "^---$" "$FILE" 2>/dev/null || echo 0)
    
    if [ "$BEFORE_COUNT" -eq 0 ]; then
        echo "✅ $STATE: No inline blocks found, skipping"
        return 0
    fi
    
    echo "🔧 Processing $STATE: $BEFORE_COUNT inline blocks to remove"
    
    # Remove all content between --- delimiters
    # This removes inline rule duplications
    awk '
        /^---$/ {
            if (in_block) {
                in_block = 0
                next
            } else {
                in_block = 1
                next
            }
        }
        !in_block {print}
    ' "$FILE" > "${FILE}.tmp"
    
    # Move cleaned file back
    mv "${FILE}.tmp" "$FILE"
    
    # Count lines after
    AFTER_COUNT=$(grep -c "^---$" "$FILE" 2>/dev/null || echo 0)
    
    if [ "$AFTER_COUNT" -eq 0 ]; then
        echo "✅ $STATE: Cleaned successfully (removed all $BEFORE_COUNT blocks)"
    else
        echo "⚠️ $STATE: Still has $AFTER_COUNT blocks remaining"
    fi
}

# Process remaining files identified with inline duplications
STATES_TO_CLEAN=(
    "WAVE_COMPLETE"
    "ERROR_RECOVERY"
    "SPAWN_CODE_REVIEWER_MERGE_PLAN"
    "SPAWN_INTEGRATION_AGENT"
    "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    "ANALYZE_CODE_REVIEWER_PARALLELIZATION"
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION"
    "CREATE_NEXT_INFRASTRUCTURE"
    "SPAWN_SW_ENGINEERS"
    "WAITING_FOR_EFFORT_PLANS"
)

TOTAL_CLEANED=0

for STATE in "${STATES_TO_CLEAN[@]}"; do
    clean_state_file "$STATE"
    if [ $? -eq 0 ]; then
        TOTAL_CLEANED=$((TOTAL_CLEANED + 1))
    fi
done

echo ""
echo "📊 Batch cleanup complete:"
echo "  - States processed: ${#STATES_TO_CLEAN[@]}"
echo "  - Successfully cleaned: $TOTAL_CLEANED"
echo ""
echo "🔍 Verification check - remaining inline blocks:"

# Check for any remaining --- delimiters
for STATE in "${STATES_TO_CLEAN[@]}"; do
    FILE="$ORCHESTRATOR_DIR/$STATE/rules.md"
    if [ -f "$FILE" ]; then
        COUNT=$(grep -c "^---$" "$FILE" 2>/dev/null || echo 0)
        if [ "$COUNT" -gt 0 ]; then
            echo "  ⚠️ $STATE: Still has $COUNT blocks"
        fi
    fi
done

echo ""
echo "✅ Batch cleanup script complete!"