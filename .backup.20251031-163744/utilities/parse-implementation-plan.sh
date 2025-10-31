#!/bin/bash

# 🚨🚨🚨 R504 Pre-Infrastructure Planning - IMPLEMENTATION-PLAN Parser
# Parses PROJECT-IMPLEMENTATION-PLAN.md and extracts all phases, waves, and efforts
# for pre-infrastructure planning

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
PLAN_FILE="${1:-$CLAUDE_PROJECT_DIR/PROJECT-IMPLEMENTATION-PLAN.md}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to log messages
log() {
    echo -e "${GREEN}[PARSER]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" >&2
}

# Check if plan file exists
if [[ ! -f "$PLAN_FILE" ]]; then
    error "Implementation plan not found: $PLAN_FILE"
fi

# Extract project prefix from plan
extract_project_prefix() {
    local prefix=""

    # Try to find project name/prefix from plan
    prefix=$(grep -i "^# \|^## Project:\|^Project Name:" "$PLAN_FILE" | head -1 | sed -E 's/.*:(.*)/\1/' | tr -d ' ' | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9-]/-/g')

    if [[ -z "$prefix" ]]; then
        # Fall back to directory name or default
        prefix=$(basename "$CLAUDE_PROJECT_DIR" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9-]/-/g')
    fi

    echo "$prefix"
}

# Parse phases from plan
parse_phases() {
    local phase_count=0
    local current_phase=""
    local current_wave=""
    local in_effort=false
    local effort_name=""
    local effort_desc=""

    # JSON output structure
    echo "{"
    echo "  \"project_prefix\": \"$(extract_project_prefix)\","
    echo "  \"phases\": {"

    while IFS= read -r line; do
        # Detect Phase headers (## Phase X: or ### Phase X:)
        if [[ "$line" =~ ^#{2,3}[[:space:]]+Phase[[:space:]]+([0-9]+):?[[:space:]]*(.*) ]]; then
            # Close previous phase if exists
            if [[ -n "$current_phase" ]]; then
                echo "      }"
                echo "    },"
            fi

            current_phase="phase${BASH_REMATCH[1]}"
            local phase_title="${BASH_REMATCH[2]}"
            ((phase_count++))

            echo "    \"$current_phase\": {"
            echo "      \"title\": \"$phase_title\","
            echo "      \"waves\": {"

            current_wave=""

        # Detect Wave headers (### Wave X: or #### Wave X:)
        elif [[ "$line" =~ ^#{3,4}[[:space:]]+Wave[[:space:]]+([0-9]+):?[[:space:]]*(.*) ]]; then
            # Close previous wave if exists
            if [[ -n "$current_wave" ]]; then
                echo "        },"
            fi

            current_wave="wave${BASH_REMATCH[1]}"
            local wave_title="${BASH_REMATCH[2]}"

            echo "        \"$current_wave\": {"
            echo "          \"title\": \"$wave_title\","
            echo "          \"efforts\": ["

        # Detect Effort patterns
        # Common patterns:
        # - **Effort Name**: Description
        # - Effort X: Name - Description
        # - #### Effort: Name
        # - • Effort Name: Description
        elif [[ "$line" =~ ^[-•*]+[[:space:]]*(.*Effort[^:]*):?[[:space:]]*(.*) ]] || \
             [[ "$line" =~ ^#{4,}[[:space:]]+Effort:?[[:space:]]*(.*) ]] || \
             [[ "$line" =~ ^[[:space:]]*Effort[[:space:]]+([0-9]+):?[[:space:]]*(.*) ]]; then

            if [[ -n "$current_wave" ]]; then
                local raw_effort="${BASH_REMATCH[1]}"
                local description="${BASH_REMATCH[2]:-${BASH_REMATCH[1]}}"

                # Clean up effort name
                effort_name=$(echo "$raw_effort" | sed 's/[*:]//g' | tr -d '•' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//' | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9-]/-/g' | sed 's/-\+/-/g' | sed 's/^-//;s/-$//')

                # Skip if effort name is empty or just "effort"
                if [[ -n "$effort_name" ]] && [[ "$effort_name" != "effort" ]]; then
                    echo "            {"
                    echo "              \"name\": \"$effort_name\","
                    echo "              \"description\": \"$description\","
                    echo "              \"phase\": \"$current_phase\","
                    echo "              \"wave\": \"$current_wave\""
                    echo "            },"
                fi
            fi
        fi
    done < "$PLAN_FILE"

    # Close all open structures
    if [[ -n "$current_wave" ]]; then
        # Remove trailing comma from last effort and close
        echo "          ]"
        echo "        }"
    fi

    if [[ -n "$current_phase" ]]; then
        echo "      }"
        echo "    }"
    fi

    echo "  },"
    echo "  \"metadata\": {"
    echo "    \"parsed_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\","
    echo "    \"source_file\": \"$PLAN_FILE\","
    echo "    \"phase_count\": $phase_count"
    echo "  }"
    echo "}"
}

# Main execution
main() {
    log "Parsing implementation plan: $PLAN_FILE"

    # Parse and output JSON structure
    local output=$(parse_phases)

    # Clean up JSON (remove trailing commas)
    echo "$output" | sed -E 's/,([[:space:]]*[}\]])/\1/g'

    log "Parsing complete"
}

# Execute main function
main "$@"