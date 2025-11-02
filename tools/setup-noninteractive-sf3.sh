#!/bin/bash
# setup-noninteractive-sf3.sh
# Seeds a target directory with the necessary Software Factory 3.0 machinery.

set -euo pipefail

TARGET_DIR="$1"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Assume this script is in tools/, so source is one level up.
SOURCE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

if [ -z "$TARGET_DIR" ]; then
    echo "Error: Target directory argument is required." >&2
    echo "Usage: $0 <target-directory>" >&2
    exit 1
fi

if [ ! -d "$TARGET_DIR" ]; then
    echo "Error: Target directory does not exist: $TARGET_DIR" >&2
    exit 1
fi

echo "Seeding SF 3.0 machinery from $SOURCE_DIR to $TARGET_DIR..."

# List of essential directories to copy
essential_dirs=(
    ".claude"
    "agent-states"
    "rule-library"
    "schemas"
    "state-machines"
    "templates"
    "tools"
    "utilities"
)

for dir in "${essential_dirs[@]}"; do
    if [ -d "$SOURCE_DIR/$dir" ]; then
        # CRITICAL: For agent-states, use rsync to exclude ARCHIVED directories
        # ARCHIVED states should NEVER exist in active projects
        if [ "$dir" = "agent-states" ]; then
            echo "  Copying $dir/ (excluding ARCHIVED)..."
            rsync -a \
                --exclude='ARCHIVED' \
                --exclude='*/ARCHIVED' \
                --exclude='*/*/ARCHIVED' \
                "$SOURCE_DIR/$dir/" "$TARGET_DIR/$dir/"
        else
            cp -r "$SOURCE_DIR/$dir" "$TARGET_DIR/"
            echo "  Copied $dir/"
        fi
    else
        echo "  Warning: Essential directory $dir/ not found in source. Continuing without it." >&2
    fi
done

echo "Seeding complete."

# Make scripts executable in the target directory
chmod +x "$TARGET_DIR/utilities/"*.sh || true
chmod +x "$TARGET_DIR/tools/"*.sh || true

echo "Set permissions for shell scripts."