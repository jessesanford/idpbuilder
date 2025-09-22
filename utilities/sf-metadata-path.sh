#!/usr/bin/env bash
# 🔴🔴🔴 R383: Software Factory Metadata File Organization Helper 🔴🔴🔴
#
# This script provides the MANDATORY helper function for creating
# timestamped metadata files in the correct .software-factory structure.
#
# ALL agents MUST source this file and use the sf_metadata_path function
# to create metadata files. Violations result in -100% grading penalty.

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Main helper function - MANDATORY USE
sf_metadata_path() {
    local phase="$1"
    local wave="$2"
    local effort="$3"
    local filename="$4"
    local ext="$5"

    # Validate inputs
    if [[ -z "$phase" || -z "$wave" || -z "$effort" || -z "$filename" || -z "$ext" ]]; then
        echo -e "${RED}❌ R383 VIOLATION: Missing parameters to sf_metadata_path${NC}" >&2
        echo "Usage: sf_metadata_path PHASE WAVE EFFORT FILENAME EXTENSION" >&2
        echo "Example: sf_metadata_path 1 2 \"registry-client\" \"IMPLEMENTATION-PLAN\" \"md\"" >&2
        exit 1
    fi

    # Create directory structure
    local dir=".software-factory/phase${phase}/wave${wave}/${effort}"
    mkdir -p "$dir"

    # Generate unique timestamped filename
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local full_path="${dir}/${filename}--${timestamp}.${ext}"

    echo "$full_path"
}

# Validation function - checks if a filename is R383 compliant
validate_metadata_filename() {
    local filepath="$1"
    local filename=$(basename "$filepath")

    # Check for timestamp pattern
    if [[ ! "$filename" =~ --[0-9]{8}-[0-9]{6}\. ]]; then
        echo -e "${RED}❌ R383 VIOLATION: Missing timestamp in $filename${NC}" >&2
        echo "   Required format: name--YYYYMMDD-HHMMSS.ext" >&2
        return 1
    fi

    # Check directory structure
    if [[ ! "$filepath" =~ \.software-factory/ ]]; then
        echo -e "${RED}❌ R383 VIOLATION: File not in .software-factory directory${NC}" >&2
        echo "   File: $filepath" >&2
        return 1
    fi

    echo -e "${GREEN}✅ R383 COMPLIANT: $filename${NC}"
    return 0
}

# Migration function - converts old files to new format
migrate_to_timestamped() {
    local old_file="$1"

    if [[ -f "$old_file" && ! "$old_file" =~ --[0-9]{8}-[0-9]{6}\. ]]; then
        # Extract base name and extension
        local base="${old_file%.*}"
        local ext="${old_file##*.}"

        # Create new timestamped name
        local timestamp=$(date +%Y%m%d-%H%M%S)
        local new_file="${base}--${timestamp}.${ext}"

        # Move file
        mv "$old_file" "$new_file"
        echo -e "${GREEN}✅ Migrated: $old_file -> $new_file${NC}"
        return 0
    fi

    echo -e "${YELLOW}⚠️ File already has timestamp or doesn't exist: $old_file${NC}"
    return 1
}

# Scan for non-compliant files
scan_for_violations() {
    local search_dir="${1:-.}"
    local violations=0

    echo "🔍 Scanning for R383 violations in: $search_dir"

    # Find all markdown, log, json files that might be metadata
    while IFS= read -r file; do
        local filename=$(basename "$file")

        # Skip files that are clearly not SF metadata
        if [[ "$filename" == "README.md" || "$filename" == "LICENSE.md" ]]; then
            continue
        fi

        # Check if it's in .software-factory and has timestamp
        if [[ "$file" =~ \.software-factory/ ]]; then
            if [[ ! "$filename" =~ --[0-9]{8}-[0-9]{6}\. ]]; then
                echo -e "${RED}❌ VIOLATION: $file (no timestamp)${NC}"
                ((violations++))
            fi
        else
            # Check if it looks like SF metadata outside .software-factory
            if [[ "$filename" =~ IMPLEMENTATION-PLAN|CODE-REVIEW-REPORT|SPLIT-PLAN|FIX-PLAN|work-log|INTEGRATION-REPORT|MERGE-PLAN|ARCHITECTURE|TEST-PLAN|DEMO-PLAN ]]; then
                echo -e "${RED}❌ VIOLATION: $file (not in .software-factory)${NC}"
                ((violations++))
            fi
        fi
    done < <(find "$search_dir" -type f \( -name "*.md" -o -name "*.log" -o -name "*.json" -o -name "*.yaml" -o -name "*.marker" -o -name "*.status" \) 2>/dev/null)

    if [[ $violations -eq 0 ]]; then
        echo -e "${GREEN}✅ No R383 violations found!${NC}"
    else
        echo -e "${RED}❌ Found $violations R383 violations${NC}"
    fi

    return $violations
}

# Get latest metadata file of a specific type
get_latest_metadata() {
    local phase="$1"
    local wave="$2"
    local effort="$3"
    local file_prefix="$4"

    local search_dir=".software-factory/phase${phase}/wave${wave}/${effort}"

    if [[ ! -d "$search_dir" ]]; then
        echo -e "${YELLOW}⚠️ Directory not found: $search_dir${NC}" >&2
        return 1
    fi

    local latest=$(ls -t "$search_dir/${file_prefix}"--*.* 2>/dev/null | head -1)

    if [[ -z "$latest" ]]; then
        echo -e "${YELLOW}⚠️ No ${file_prefix} files found in $search_dir${NC}" >&2
        return 1
    fi

    echo "$latest"
}

# Example usage functions for common metadata types
create_implementation_plan() {
    local phase="$1"
    local wave="$2"
    local effort="$3"

    local plan_path=$(sf_metadata_path "$phase" "$wave" "$effort" "IMPLEMENTATION-PLAN" "md")
    echo -e "${GREEN}✅ Creating implementation plan: $plan_path${NC}"
    echo "$plan_path"
}

create_code_review_report() {
    local phase="$1"
    local wave="$2"
    local effort="$3"

    local report_path=$(sf_metadata_path "$phase" "$wave" "$effort" "CODE-REVIEW-REPORT" "md")
    echo -e "${GREEN}✅ Creating review report: $report_path${NC}"
    echo "$report_path"
}

create_work_log() {
    local phase="$1"
    local wave="$2"
    local effort="$3"

    local log_path=$(sf_metadata_path "$phase" "$wave" "$effort" "work-log" "log")
    echo -e "${GREEN}✅ Creating work log: $log_path${NC}"
    echo "$log_path"
}

# If script is run directly (not sourced), show usage
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "🔴🔴🔴 R383: Software Factory Metadata File Organization Helper 🔴🔴🔴"
    echo ""
    echo "This script provides helper functions for R383 compliance."
    echo ""
    echo "Usage:"
    echo "  Source this file:  source $0"
    echo "  Then use:         sf_metadata_path PHASE WAVE EFFORT FILENAME EXT"
    echo ""
    echo "Example:"
    echo '  source utilities/sf-metadata-path.sh'
    echo '  PLAN_PATH=$(sf_metadata_path 1 2 "registry-client" "IMPLEMENTATION-PLAN" "md")'
    echo '  echo "# Implementation Plan" > "$PLAN_PATH"'
    echo ""
    echo "Available functions:"
    echo "  - sf_metadata_path: Generate compliant file path"
    echo "  - validate_metadata_filename: Check if file is compliant"
    echo "  - migrate_to_timestamped: Convert old files to new format"
    echo "  - scan_for_violations: Find non-compliant files"
    echo "  - get_latest_metadata: Find most recent metadata file"
    echo ""
    echo "Shortcut functions:"
    echo "  - create_implementation_plan PHASE WAVE EFFORT"
    echo "  - create_code_review_report PHASE WAVE EFFORT"
    echo "  - create_work_log PHASE WAVE EFFORT"
    echo ""

    # If given arguments, run scan
    if [[ $# -eq 1 ]]; then
        case "$1" in
            scan)
                scan_for_violations "."
                ;;
            *)
                echo "Unknown command: $1"
                echo "Try: $0 scan"
                ;;
        esac
    fi
fi