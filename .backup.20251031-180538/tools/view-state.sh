#!/bin/bash
# View minified state file in readable format

STATE_FILE="${1:-orchestrator-state-v3.json}"

if [[ ! -f "$STATE_FILE" ]]; then
    echo "ERROR: State file not found: $STATE_FILE"
    exit 1
fi

# Pretty-print to pager
jq '.' "$STATE_FILE" | less -R
