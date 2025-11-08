#!/bin/bash

# Simple script to update all references from markdown to JSON state machine

echo "Updating all references to JSON state machine..."

# Update all markdown references to JSON
find . -type f -name "*.md" -o -name "*.sh" | while read file; do
    # Skip the actual state machine files and backups
    if [[ "$file" == *"software-factory-3.0-state-machine.json"* ]] || [[ "$file" == *".bak"* ]] || [[ "$file" == *".git"* ]]; then
        continue
    fi

    # Check if file contains the old reference
    if grep -q "SOFTWARE-FACTORY-STATE-MACHINE\.md" "$file" 2>/dev/null; then
        echo "Updating: $file"
        sed -i.bak 's/SOFTWARE-FACTORY-STATE-MACHINE\.md/software-factory-3.0-state-machine.json/g' "$file"
        rm -f "${file}.bak"
    fi
done

echo "Creating JSON reference guide..."
cat > rule-library/JSON-STATE-MACHINE-USAGE.md << 'EOF'
# JSON State Machine Usage Guide

## File Locations
- Main: `$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json`
- Sub-machines: `$CLAUDE_PROJECT_DIR/state-machines/[agent].json`

## Common Operations

### Check if state exists
```bash
jq -e '.agents.orchestrator.states.INIT' software-factory-3.0-state-machine.json
```

### Get valid transitions
```bash
jq -r '.transition_matrix.orchestrator.INIT[]' software-factory-3.0-state-machine.json
```

### List all states
```bash
jq -r '.agents.orchestrator.states | keys[]' software-factory-3.0-state-machine.json
```

### Validate transition
```bash
STATE="INIT"
TARGET="WAVE_START"
jq -e ".transition_matrix.orchestrator.$STATE | index(\"$TARGET\")" software-factory-3.0-state-machine.json
```

## Migration Notes
Replace grep patterns with jq queries when working with the state machine.
EOF

echo "Done!"