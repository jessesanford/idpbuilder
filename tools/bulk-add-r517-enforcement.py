#!/usr/bin/env python3
"""
Bulk update all state rules files with R517 State Manager enforcement section
Created: 2025-11-01
Purpose: Fulfill user mandate for complete State Manager enforcement
"""

import os
import sys
import re
from pathlib import Path
from datetime import datetime

# Enforcement section template
ENFORCEMENT_SECTION = """
## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
"""


def update_file(file_path):
    """Update a single rules.md file with R517 enforcement section."""

    # Read current content
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # Check if already has R517
    if 'R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW' in content:
        return 'skip', 'Already has R517 enforcement'

    # Create backup
    backup_path = f"{file_path}.backup-{datetime.now().strftime('%Y%m%d-%H%M%S')}"
    with open(backup_path, 'w', encoding='utf-8') as f:
        f.write(content)

    # Find insertion point (after PRIMARY DIRECTIVES or first major section)
    lines = content.split('\n')
    insertion_index = None

    # Look for PRIMARY DIRECTIVES, CORE DIRECTIVES, or STATE-SPECIFIC RULES
    for i, line in enumerate(lines):
        if re.match(r'^#\s+(PRIMARY DIRECTIVES|CORE DIRECTIVES|STATE-SPECIFIC RULES)', line):
            # Found the section, now find the next blank line after it
            for j in range(i + 1, len(lines)):
                if lines[j].strip() == '':
                    insertion_index = j + 1
                    break
            break

    # If no directive section found, insert after first header
    if insertion_index is None:
        for i, line in enumerate(lines):
            if line.startswith('# ') and not line.startswith('## '):
                # Find next blank line
                for j in range(i + 1, len(lines)):
                    if lines[j].strip() == '':
                        insertion_index = j + 1
                        break
                break

    # If still not found, insert at beginning
    if insertion_index is None:
        insertion_index = 0

    # Insert the enforcement section
    lines.insert(insertion_index, ENFORCEMENT_SECTION.strip())

    # Write updated content
    updated_content = '\n'.join(lines)
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(updated_content)

    # Verify insertion
    with open(file_path, 'r', encoding='utf-8') as f:
        verify_content = f.read()

    if 'R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW' in verify_content:
        return 'success', 'Updated successfully'
    else:
        # Restore backup
        with open(backup_path, 'r', encoding='utf-8') as f:
            original = f.read()
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(original)
        return 'fail', 'Failed to insert enforcement section'


def main():
    """Main execution function."""

    print("🏭 R517 ENFORCEMENT BULK UPDATE")
    print("=" * 60)
    print(f"Started: {datetime.utcnow().strftime('%Y-%m-%d %H:%M:%S UTC')}")
    print()

    # Find project root
    script_dir = Path(__file__).parent
    project_root = script_dir.parent
    agent_states_dir = project_root / 'agent-states'

    print(f"Project: {project_root}")
    print(f"Agent states: {agent_states_dir}")
    print()

    # Find all rules.md files
    rules_files = sorted(agent_states_dir.rglob('rules.md'))

    print(f"📋 Processing {len(rules_files)} state rules files...")
    print()

    # Counters
    total = 0
    updated = 0
    skipped = 0
    failed = 0
    failed_list = []

    # Process each file
    for file_path in rules_files:
        total += 1
        rel_path = file_path.relative_to(project_root)

        print(f"[{total}] {rel_path}")

        status, message = update_file(file_path)

        if status == 'success':
            print(f"  ✅ {message}")
            updated += 1
        elif status == 'skip':
            print(f"  ⏭️  {message}")
            skipped += 1
        else:  # fail
            print(f"  ❌ {message}")
            failed += 1
            failed_list.append(str(rel_path))

    # Print summary
    print()
    print("=" * 60)
    print("📊 BULK UPDATE SUMMARY")
    print("=" * 60)
    print(f"Total files processed: {total}")
    print(f"✅ Files updated: {updated}")
    print(f"⏭️  Files skipped (already had R517): {skipped}")
    print(f"❌ Files failed: {failed}")
    print()

    if failed > 0:
        print("⚠️  FAILED FILES:")
        for failed_file in failed_list:
            print(f"  - {failed_file}")
        print()

    # Calculate success rate
    success_rate = ((updated + skipped) * 100) // total if total > 0 else 0
    print(f"Success rate: {success_rate}%")
    print()

    if success_rate == 100:
        print("🎉 PERFECT: 100% compliance achieved!")
        return 0
    else:
        print("⚠️  WARNING: Some files failed to update")
        return 1


if __name__ == '__main__':
    sys.exit(main())
