#!/bin/bash

# Script to convert old box-delimiter rules to new clean delimiter + criticality format
# Usage: ./convert-rules-to-new-format.sh [directory]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Directory to process (default to current)
TARGET_DIR="${1:-.}"

echo -e "${GREEN}Converting rule files to new delimiter + criticality format...${NC}"

# Function to determine criticality based on rule content
determine_criticality() {
    local rule_num="$1"
    local content="$2"
    
    # Check for blocking keywords
    if echo "$content" | grep -qi "exit 1\|immediate stop\|blocking\|must.*before.*any.*work"; then
        echo "BLOCKING"
        return
    fi
    
    # Check for mandatory keywords
    if echo "$content" | grep -qi "mandatory\|required.*approval\|must.*complete\|grading.*failure"; then
        echo "MANDATORY"
        return
    fi
    
    # Check for critical keywords
    if echo "$content" | grep -qi "critical\|size.*limit\|800.*lines\|automatic.*split"; then
        echo "CRITICAL"
        return
    fi
    
    # Check for important keywords
    if echo "$content" | grep -qi "important\|should\|recommend"; then
        echo "IMPORTANT"
        return
    fi
    
    # Default to INFO
    echo "INFO"
}

# Function to get emoji for criticality
get_emoji() {
    case "$1" in
        "BLOCKING") echo "🚨🚨🚨" ;;
        "MANDATORY") echo "🚨🚨" ;;
        "CRITICAL") echo "🚨" ;;
        "IMPORTANT") echo "⚠️" ;;
        *) echo "ℹ️" ;;
    esac
}

# Find all markdown files with box delimiters
echo -e "${YELLOW}Searching for files with old box delimiters...${NC}"
FILES=$(grep -r "┌─────" "$TARGET_DIR" --include="*.md" 2>/dev/null | cut -d: -f1 | sort -u || true)

if [ -z "$FILES" ]; then
    echo -e "${GREEN}No files with old box delimiters found!${NC}"
    exit 0
fi

# Process each file
for file in $FILES; do
    echo -e "${YELLOW}Processing: $file${NC}"
    
    # Create backup
    cp "$file" "${file}.backup"
    
    # Create temporary file for conversion
    temp_file="${file}.tmp"
    > "$temp_file"
    
    # Process file line by line
    in_rule=false
    rule_num=""
    rule_title=""
    rule_source=""
    rule_content=""
    
    while IFS= read -r line; do
        # Check for rule box start
        if [[ "$line" =~ ┌─────.*─────┐ ]]; then
            in_rule=true
            rule_content=""
            continue
        fi
        
        # Check for rule box end
        if [[ "$line" =~ └─────.*─────┘ ]]; then
            if [ "$in_rule" = true ]; then
                # Determine criticality
                criticality=$(determine_criticality "$rule_num" "$rule_content")
                emoji=$(get_emoji "$criticality")
                
                # Write new format
                echo "---" >> "$temp_file"
                echo "### $emoji RULE $rule_num - $rule_title" >> "$temp_file"
                echo "**Source:** $rule_source" >> "$temp_file"
                
                # Add criticality line based on level
                case "$criticality" in
                    "BLOCKING")
                        echo "**Criticality:** BLOCKING - Failure = Immediate Stop (exit 1)" >> "$temp_file"
                        ;;
                    "MANDATORY")
                        echo "**Criticality:** MANDATORY - Required for approval" >> "$temp_file"
                        ;;
                    "CRITICAL")
                        echo "**Criticality:** CRITICAL - Major impact on grading" >> "$temp_file"
                        ;;
                    "IMPORTANT")
                        echo "**Criticality:** IMPORTANT - Affects workflow" >> "$temp_file"
                        ;;
                    *)
                        echo "**Criticality:** INFO - Best practice" >> "$temp_file"
                        ;;
                esac
                
                echo "" >> "$temp_file"
                echo "$rule_content" >> "$temp_file"
                echo "---" >> "$temp_file"
                
                in_rule=false
                rule_content=""
                continue
            fi
        fi
        
        # Process lines inside rule box
        if [ "$in_rule" = true ]; then
            # Extract rule header info
            if [[ "$line" =~ │\ RULE\ R([0-9.]+)\ -\ (.*)\ *│ ]]; then
                rule_num="R${BASH_REMATCH[1]}"
                rule_title="${BASH_REMATCH[2]}"
                # Remove trailing spaces from title
                rule_title=$(echo "$rule_title" | sed 's/[[:space:]]*$//')
                continue
            fi
            
            # Extract source
            if [[ "$line" =~ │\ Source:\ (.*)\ *│ ]]; then
                rule_source="${BASH_REMATCH[1]}"
                # Remove trailing spaces from source
                rule_source=$(echo "$rule_source" | sed 's/[[:space:]]*$//')
                continue
            fi
            
            # Skip box divider lines
            if [[ "$line" =~ ├─────.*─────┤ ]]; then
                continue
            fi
            
            # Extract content (remove box sides)
            if [[ "$line" =~ ^│(.*)│$ ]]; then
                content_line="${BASH_REMATCH[1]}"
                # Remove leading/trailing spaces but preserve internal spacing
                content_line=$(echo "$content_line" | sed 's/^ *//;s/ *$//')
                if [ -n "$rule_content" ]; then
                    rule_content="$rule_content"$'\n'"$content_line"
                else
                    rule_content="$content_line"
                fi
                continue
            fi
        fi
        
        # Write non-rule lines as-is
        echo "$line" >> "$temp_file"
    done < "$file"
    
    # Replace original file with converted version
    mv "$temp_file" "$file"
    
    echo -e "${GREEN}✓ Converted: $file${NC}"
done

echo -e "${GREEN}Conversion complete! Backup files created with .backup extension${NC}"
echo -e "${YELLOW}To review changes: diff file.md file.md.backup${NC}"
echo -e "${YELLOW}To restore: mv file.md.backup file.md${NC}"