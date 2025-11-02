#!/bin/bash

# Software Factory 3.0 Upgrade Script Wrapper
# This wrapper runs the actual upgrade script with minimal CPU and IO priority
# to prevent system slowdowns during upgrades.
#
# CPU Priority: nice 19 (lowest - only runs when CPU is idle)
# IO Priority:  ionice class 3 (idle - only performs IO when system is idle)
#               Falls back to class 2 priority 7 (lowest best-effort) if idle class unavailable
#
# All keyboard commands and signals are passed through to the underlying script.

set -e

# Get the directory where this wrapper script lives
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UNDERLYING_SCRIPT="$SCRIPT_DIR/_upgrade-no-nice.sh"

# Verify the underlying script exists
if [ ! -f "$UNDERLYING_SCRIPT" ]; then
    echo "ERROR: Underlying upgrade script not found: $UNDERLYING_SCRIPT" >&2
    exit 1
fi

# Check if we're already running with nice/ionice (to avoid double-wrapping)
if [ -n "$SF_UPGRADE_NICED" ]; then
    # Already wrapped, just run the underlying script
    exec bash "$UNDERLYING_SCRIPT" "$@"
fi

# Mark that we're now running in the niced wrapper
export SF_UPGRADE_NICED=1

# Determine ionice settings
# Try idle class (3) first, fall back to best-effort lowest priority (class 2, priority 7)
IONICE_CMD=""
if command -v ionice &> /dev/null; then
    # Test if idle class is available
    if ionice -c 3 -p $$ &> /dev/null; then
        # Idle class works
        IONICE_CMD="ionice -c 3"
    else
        # Fall back to best-effort with lowest priority
        IONICE_CMD="ionice -c 2 -n 7"
    fi
fi

# Execute the underlying script with maximum nice and ionice
# Using exec ensures:
# 1. This wrapper is replaced by the child process (signals pass through)
# 2. The PID remains the same (important for process management)
# 3. All arguments are passed through exactly as received
# 4. Keyboard interrupts (Ctrl+C) work correctly

if [ -n "$IONICE_CMD" ]; then
    # Run with both nice and ionice
    exec $IONICE_CMD nice -n 19 bash "$UNDERLYING_SCRIPT" "$@"
else
    # ionice not available, just use nice
    exec nice -n 19 bash "$UNDERLYING_SCRIPT" "$@"
fi
