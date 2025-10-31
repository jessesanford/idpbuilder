#!/usr/bin/env python3
"""
Fix duplicate checklist enforcement sections in orchestrator state rules.

This script removes the second/duplicate enforcement section that conflicts
with the first one, causing R405 compliance failures.
"""

import re
import sys
from pathlib import Path

# States that need fixing (excluding already fixed ones)
BROKEN_STATES = [
    "CREATE_NEXT_INFRASTRUCTURE",
    "ERROR_RECOVERY",
    "IMMEDIATE_BACKPORT_REQUIRED",
    "INJECT_WAVE_METADATA",
    "REVIEW_WAVE_INTEGRATION",
    "MONITORING_BACKPORT_PROGRESS",
    "MONITORING_EFFORT_FIXES",
    "PROJECT_REVIEW_WAVE_INTEGRATION",
    "PR_PLAN_CREATION",
    "PROJECT_DONE",
    "WAITING_FOR_BACKPORT_PLAN",
    "WAITING_FOR_PHASE_FIX_PLANS",
    "WAITING_FOR_PHASE_REVIEW_WAVE_INTEGRATION",
    "WAITING_FOR_PHASE_PLANS",
    "WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION",
    "WAITING_FOR_PROJECT_TEST_PLAN",
    "WAITING_FOR_PROJECT_VALIDATION",
]

def fix_duplicate_checklist(file_path: Path) -> bool:
    """
    Remove duplicate checklist enforcement section from a state rules file.

    Returns True if file was modified, False otherwise.
    """
    content = file_path.read_text()

    # Pattern to match TWO consecutive enforcement sections
    # This is safer than trying to find the exact text since there may be variations
    pattern = r'(## 🚨 CHECKLIST ENFORCEMENT 🚨\s*\n\s*\n.*?ALL \d+ STEPS ARE MANDATORY[^\n]*\n)(## 🚨 CHECKLIST ENFORCEMENT 🚨\s*\n\s*\n.*?ALL \d+ STEPS ARE MANDATORY[^\n]*\n)'

    matches = list(re.finditer(pattern, content, re.DOTALL))

    if not matches:
        print(f"  ⚠️  No duplicate enforcement section found in {file_path.name}")
        return False

    if len(matches) > 1:
        print(f"  ⚠️  Multiple duplicate patterns found in {file_path.name} - manual review needed")
        return False

    # Replace with just the first enforcement section (updated with SF 3.0 language)
    replacement = r'''\1**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

'''

    new_content = re.sub(pattern, replacement, content, count=1, flags=re.DOTALL)

    if new_content == content:
        print(f"  ⚠️  Pattern matched but no changes made in {file_path.name}")
        return False

    # Write back
    file_path.write_text(new_content)
    print(f"  ✅ Fixed {file_path.name}")
    return True

def main():
    base_dir = Path("/home/vscode/software-factory-template/agent-states/software-factory/orchestrator")

    fixed_count = 0
    failed_count = 0

    print("🔧 Fixing duplicate checklist enforcement sections...")
    print()

    for state_name in BROKEN_STATES:
        rules_file = base_dir / state_name / "rules.md"

        if not rules_file.exists():
            print(f"  ❌ File not found: {rules_file}")
            failed_count += 1
            continue

        try:
            if fix_duplicate_checklist(rules_file):
                fixed_count += 1
            else:
                failed_count += 1
        except Exception as e:
            print(f"  ❌ Error fixing {state_name}: {e}")
            failed_count += 1

    print()
    print(f"📊 Results: {fixed_count} fixed, {failed_count} failed/skipped")

    return 0 if failed_count == 0 else 1

if __name__ == "__main__":
    sys.exit(main())
