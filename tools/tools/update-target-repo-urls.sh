#!/bin/bash
# Script to add target_repo_url to all efforts in orchestrator-state-v3.json
#
# PARAMETERIZED: No hard-coded repository URLs or paths
# Usage:
#   1. Default: Uses orchestrator-state-v3.json in current directory
#      ./update-target-repo-urls.sh
#   2. Custom state file:
#      ./update-target-repo-urls.sh /path/to/orchestrator-state-v3.json
#   3. Custom state file + target URL:
#      ./update-target-repo-urls.sh /path/to/orchestrator-state-v3.json https://github.com/user/repo.git

# Determine state file location
if [ -n "$1" ]; then
    STATE_FILE="$1"
else
    STATE_FILE="orchestrator-state-v3.json"
fi

# Check if state file exists
if [ ! -f "$STATE_FILE" ]; then
    echo "ERROR: State file not found: $STATE_FILE"
    echo "Please provide a valid orchestrator-state-v3.json file"
    exit 1
fi

# Determine target repository URL
if [ -n "$2" ]; then
    # Command-line argument provided
    TARGET_URL="$2"
elif command -v jq &> /dev/null; then
    # Try to read from state file's pre_planned_infrastructure.target_repo_url
    TARGET_URL=$(jq -r '.pre_planned_infrastructure.target_repo_url // empty' "$STATE_FILE")
fi

# If still empty, error
if [ -z "$TARGET_URL" ]; then
    echo "ERROR: TARGET_URL not found!"
    echo "Please either:"
    echo "  1. Add target_repo_url to orchestrator-state-v3.json (pre_planned_infrastructure.target_repo_url)"
    echo "  2. Pass repository URL as argument: $0 <state-file> <repo-url>"
    exit 1
fi

echo "Using state file: $STATE_FILE"
echo "Using target URL: $TARGET_URL"

# Backup the file first
cp "$STATE_FILE" "${STATE_FILE}.backup"

# Get all effort keys
EFFORT_KEYS=$(jq -r '.pre_planned_infrastructure.efforts | keys[]' "$STATE_FILE")

# Update each effort to include target_repo_url
for effort_key in $EFFORT_KEYS; do
    echo "Updating effort: $effort_key"

    # Check if target_repo_url already exists
    has_url=$(jq ".pre_planned_infrastructure.efforts[\"$effort_key\"].target_repo_url // false" "$STATE_FILE")

    if [ "$has_url" = "false" ]; then
        # Add target_repo_url to the effort
        jq ".pre_planned_infrastructure.efforts[\"$effort_key\"] += {\"target_repo_url\": \"$TARGET_URL\"}" "$STATE_FILE" > "${STATE_FILE}.tmp"
        mv "${STATE_FILE}.tmp" "$STATE_FILE"
        echo "  ✅ Added target_repo_url to $effort_key"
    else
        echo "  ⏭️  $effort_key already has target_repo_url"
    fi
done

echo ""
echo "✅ All efforts updated with target_repo_url"
echo "Backup saved to: ${STATE_FILE}.backup"