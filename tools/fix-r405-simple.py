#!/usr/bin/env python3
"""
Fix R405 example sections in orchestrator state rule files.
Replace old jq pattern with SF 3.0 State Manager pattern.
"""

import os
import sys
from pathlib import Path

# Old pattern (what we're replacing)
OLD_PATTERN = """### THE PATTERN AT R322 CHECKPOINTS

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

# New pattern (SF 3.0)
NEW_PATTERN = """### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

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

def main():
    project_root = Path(__file__).parent.parent
    orchestrator_dir = project_root / "agent-states" / "software-factory" / "orchestrator"

    print("🔧 Fixing R405 example sections in orchestrator state rule files...")
    print("=" * 68)

    fixed_count = 0
    skipped_count = 0

    # Find all rules.md files
    for rules_file in sorted(orchestrator_dir.glob("*/rules.md")):
        state_name = rules_file.parent.name

        # Read file content
        content = rules_file.read_text()

        # Check if old pattern exists
        if OLD_PATTERN in content:
            # Replace old pattern with new pattern
            new_content = content.replace(OLD_PATTERN, NEW_PATTERN)

            # Write back
            rules_file.write_text(new_content)

            print(f"✅ Fixed: {state_name}")
            fixed_count += 1
        else:
            # Check if file has any jq violations (might be different format)
            if "jq '.state_machine.current_state" in content:
                print(f"⚠️  Skipped: {state_name} (pattern mismatch, needs manual fix)")
                skipped_count += 1

    print("=" * 68)
    print(f"✅ Fixed: {fixed_count} files")
    if skipped_count > 0:
        print(f"⚠️  Skipped: {skipped_count} files (manual fixes needed)")
    print("")
    print("Next steps:")
    print("1. Review changes: git diff agent-states/software-factory/orchestrator/*/rules.md")
    print("2. Commit changes")
    print("3. Push to remote")

if __name__ == "__main__":
    main()
