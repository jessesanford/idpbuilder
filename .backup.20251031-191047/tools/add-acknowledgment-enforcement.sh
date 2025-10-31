#!/bin/bash

# Script to add mandatory acknowledgment enforcement to all orchestrator state rule files

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_DIR="/home/vscode/software-factory-template/agent-states/orchestrator"
UPDATED_FILES=0
SKIPPED_FILES=0
ERRORS=0

echo -e "${GREEN}Starting acknowledgment enforcement update...${NC}"
echo "Base directory: $BASE_DIR"
echo "----------------------------------------"

# Function to add acknowledgment section to a rule file
add_acknowledgment_section() {
    local file_path="$1"
    local state_name="$2"
    
    # Check if file exists
    if [ ! -f "$file_path" ]; then
        echo -e "${YELLOW}  SKIP: No rules.md file found${NC}"
        ((SKIPPED_FILES++))
        return 1
    fi
    
    # Check if acknowledgment section already exists
    if grep -q "MANDATORY STATE RULE READING AND ACKNOWLEDGMENT" "$file_path"; then
        echo -e "${YELLOW}  SKIP: Acknowledgment section already exists${NC}"
        ((SKIPPED_FILES++))
        return 1
    fi
    
    # Create the acknowledgment section with proper state name
    local acknowledgment_section="## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   \`\`\`
   ❌ WRONG: \"I acknowledge R151, R208, R053...\"
   (No Read tool calls detected = AUTOMATIC FAILURE)
   \`\`\`

2. **Bulk Acknowledgment**:
   \`\`\`
   ❌ WRONG: \"I acknowledge all ${state_name} rules\"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   \`\`\`

3. **Silent Reading**:
   \`\`\`
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   \"Now I've read the rules, let me start work...\"
   (MUST explicitly acknowledge EACH rule)
   \`\`\`

4. **Reading From Memory**:
   \`\`\`
   ❌ WRONG: \"I know R208 requires CD before spawn...\"
   (Must READ from file, not recall from memory)
   \`\`\`

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   \`\`\`
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   \`\`\`

### ✅ CORRECT PATTERN FOR ${state_name}:
\`\`\`
1. READ: \$CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. \"I acknowledge [Rule#] - [Rule Name]: [Brief description]\"
3. READ: \$CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. \"I acknowledge [Rule#] - [Rule Name]: [Brief description]\"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. \"Ready to execute ${state_name} work\"
\`\`\`

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY ${state_name} work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute ${state_name} work

"
    
    # Create a temporary file with the new content
    local temp_file=$(mktemp)
    
    # Process the file to insert the acknowledgment section
    local section_inserted=false
    local in_header=false
    
    while IFS= read -r line; do
        # Look for various forms of PRIMARY DIRECTIVES to insert before it
        if [ "$section_inserted" = false ]; then
            if [[ "$line" =~ "PRIMARY DIRECTIVES" ]] || [[ "$line" =~ "Primary Directives" ]] || [[ "$line" =~ "## 📋 PRIMARY DIRECTIVES" ]]; then
                # Add our acknowledgment section before PRIMARY DIRECTIVES
                echo "" >> "$temp_file"
                echo "$acknowledgment_section" >> "$temp_file"
                section_inserted=true
            fi
        fi
        
        echo "$line" >> "$temp_file"
    done < "$file_path"
    
    if [ "$section_inserted" = true ]; then
        # Move the temp file to the original
        mv "$temp_file" "$file_path"
        echo -e "${GREEN}  ✅ Updated successfully${NC}"
        ((UPDATED_FILES++))
        return 0
    else
        rm "$temp_file"
        echo -e "${RED}  ❌ ERROR: Could not find insertion point${NC}"
        ((ERRORS++))
        return 1
    fi
}

# Process all state directories
for state_dir in "$BASE_DIR"/*; do
    if [ -d "$state_dir" ]; then
        state_name=$(basename "$state_dir")
        rules_file="$state_dir/rules.md"
        
        echo -e "\n${YELLOW}Processing: ${state_name}${NC}"
        echo "  File: $rules_file"
        
        add_acknowledgment_section "$rules_file" "$state_name"
    fi
done

echo -e "\n========================================="
echo -e "${GREEN}SUMMARY:${NC}"
echo -e "  Updated files: ${GREEN}$UPDATED_FILES${NC}"
echo -e "  Skipped files: ${YELLOW}$SKIPPED_FILES${NC}"
echo -e "  Errors: ${RED}$ERRORS${NC}"
echo -e "========================================="

if [ $ERRORS -gt 0 ]; then
    echo -e "${RED}⚠️  Some files could not be updated. Please check manually.${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Acknowledgment enforcement update complete!${NC}"
    exit 0
fi