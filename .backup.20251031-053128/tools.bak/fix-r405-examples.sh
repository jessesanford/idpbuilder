#!/bin/bash
# Fix R405 example sections in all orchestrator state rule files
# Replace old jq pattern with SF 3.0 State Manager pattern

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "🔧 Fixing R405 example sections in orchestrator state rule files..."
echo "===================================================================="

# Old pattern (what we're replacing)
OLD_PATTERN='### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '\''.state_machine.current_state = "NEXT_STATE"'\'' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```'

# New pattern (SF 3.0)
NEW_PATTERN='### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 6. Agent stops (technical requirement)
exit 0
```'

# Counter
FIXED_COUNT=0

# Find all orchestrator state rule files
for file in agent-states/software-factory/orchestrator/*/rules.md; do
    # Skip if file doesn't contain the old pattern
    if ! grep -q 'jq.*current_state.*orchestrator-state' "$file"; then
        continue
    fi

    # Check if it's in the R405 example section (around line 291-308)
    LINE_NUM=$(grep -n "jq.*current_state.*orchestrator-state" "$file" | head -1 | cut -d: -f1)

    # If the violation is in the R405 example section (typically lines 290-310)
    if [ "$LINE_NUM" -ge 285 ] && [ "$LINE_NUM" -le 315 ]; then
        STATE_NAME=$(basename $(dirname "$file"))
        echo "Fixing: $STATE_NAME (line $LINE_NUM)"

        # Create a Python script to do the replacement (easier for multiline)
        python3 << 'EOF'
import sys
import re

file_path = sys.argv[1]

with open(file_path, 'r') as f:
    content = f.read()

# Old pattern
old = """### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```"""

# New pattern
new = """### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \\
  --current-state "CURRENT_STATE" \\
  --proposed-next-state "$PROPOSED_NEXT_STATE" \\
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 6. Agent stops (technical requirement)
exit 0
```"""

if old in content:
    content = content.replace(old, new)
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"✅ Fixed {file_path}")
    sys.exit(0)
else:
    print(f"⚠️  Pattern not found in {file_path}")
    sys.exit(1)
EOF

        if python3 -c "import sys; sys.argv = ['', '$file']" "$file"; then
            ((FIXED_COUNT++))
        fi
    fi
done

echo "===================================================================="
echo "✅ Fixed $FIXED_COUNT R405 example sections"
echo ""
echo "Next steps:"
echo "1. Review changes: git diff"
echo "2. Commit: git add . && git commit -m 'fix: Update R405 examples to SF 3.0 pattern [R288]'"
echo "3. Push: git push"
