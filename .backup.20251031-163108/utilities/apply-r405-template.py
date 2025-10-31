#!/usr/bin/env python3
"""
Apply R405 error handling template to all orchestrator states.

This script:
1. Adds error handling to git commit operations (Step 5)
2. Adds error handling to TODO commit operations (Step 6)
3. Normalizes R405 continuation flag language with REASON fields
4. Standardizes exit checklist formatting
"""

import re
import os
import sys
from pathlib import Path

def apply_commit_error_handling(content, state_name):
    """Apply error handling to git commit operations."""

    # Pattern 1: Simple git commit without error handling
    # Look for patterns like: git commit -m "state: ... [R288]"
    pattern1 = re.compile(
        r'(git add[^\n]+\n+)'  # git add line
        r'(git commit -m "[^"]+\[R288\]")',  # git commit line
        re.MULTILINE
    )

    def replacement1(match):
        git_add = match.group(1)
        commit_msg = match.group(2)

        # Extract the commit message
        msg_match = re.search(r'git commit -m "([^"]+)"', commit_msg)
        if not msg_match:
            return match.group(0)

        original_msg = msg_match.group(1)

        # Create error-handled version
        return f'''{git_add}
if ! {commit_msg}; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: {state_name}"
    echo "Attempted transition from: {state_name}"
    echo ""
    echo "Common causes:"
    echo "  - Schema validation failure (check pre-commit hook output above)"
    echo "  - Missing required fields in JSON files"
    echo "  - Invalid JSON syntax"
    echo ""
    echo "🛑 Cannot proceed - manual intervention required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
    exit 1
fi

git push || echo "⚠️ WARNING: Push failed - committed locally"
echo "✅ State file committed and pushed"'''

    # Apply pattern 1
    content = pattern1.sub(replacement1, content)

    # Pattern 2: TODO commits without error handling
    # Look for: git commit -m "todo: ... [R287]"
    pattern2 = re.compile(
        r'(git add todos/\*\.todo\n+)'
        r'(git commit -m "todo:[^"]+\[R287\]")',
        re.MULTILINE
    )

    def replacement2(match):
        git_add = match.group(1)
        commit_msg = match.group(2)

        return f'''{git_add}
if ! {commit_msg}; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"'''

    # Apply pattern 2
    content = pattern2.sub(replacement2, content)

    return content


def normalize_r405_language(content, state_name):
    """Normalize R405 continuation flag language."""

    # Pattern 1: CONTINUE-SOFTWARE-FACTORY=TRUE without REASON
    pattern1 = re.compile(
        r'echo "CONTINUE-SOFTWARE-FACTORY=TRUE"(?!\s+REASON=)',
        re.MULTILINE
    )

    content = pattern1.sub('echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"', content)

    # Pattern 2: CONTINUE-SOFTWARE-FACTORY=FALSE without REASON
    pattern2 = re.compile(
        r'echo "CONTINUE-SOFTWARE-FACTORY=FALSE"(?!\s+REASON=)',
        re.MULTILINE
    )

    # Don't replace these - they were added by error handling template
    # Only replace standalone FALSE without context
    def check_context(match):
        # Get surrounding context
        start = max(0, match.start() - 200)
        context = content[start:match.start()]

        # If this is part of error handling block, leave it
        if 'SCHEMA_VALIDATION' in context or 'COMMIT_FAILURE' in context:
            return match.group(0)

        return 'echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"'

    # Apply with context check
    content = pattern2.sub(check_context, content)

    return content


def standardize_exit_checklist(content):
    """Standardize exit checklist formatting."""

    # Ensure consistent heading format
    patterns = [
        (r'##\s*EXIT CHECKLIST.*R405', '## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol)'),
        (r'##\s*R405.*CONTINUATION.*FLAG', '## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol)'),
        (r'##\s*🔴.*R405.*FLAG', '## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol)'),
    ]

    for pattern, replacement in patterns:
        content = re.sub(pattern, replacement, content, flags=re.IGNORECASE)

    return content


def process_state_file(filepath):
    """Process a single state file."""

    # Get state name from path
    state_name = filepath.parent.name

    # Skip template directory
    if state_name == '_TEMPLATES':
        return False

    # Read file
    with open(filepath, 'r') as f:
        content = f.read()

    # Skip if already has error handling
    if 'if ! git commit' in content and 'REASON=SCHEMA_VALIDATION' in content:
        print(f"  ⏭️  {state_name}: Already has error handling")
        return False

    original_content = content

    # Apply transformations
    content = apply_commit_error_handling(content, state_name)
    content = normalize_r405_language(content, state_name)
    content = standardize_exit_checklist(content)

    # Check if anything changed
    if content == original_content:
        print(f"  ⏭️  {state_name}: No changes needed")
        return False

    # Write back
    with open(filepath, 'w') as f:
        f.write(content)

    print(f"  ✅ {state_name}: Updated successfully")
    return True


def main():
    """Main entry point."""

    # Find all orchestrator state rules files
    base_path = Path('/home/vscode/software-factory-template/agent-states/software-factory/orchestrator')

    if not base_path.exists():
        print(f"❌ Error: Path not found: {base_path}")
        return 1

    # Get all rules.md files
    rules_files = sorted(base_path.glob('*/rules.md'))

    print(f"📋 Found {len(rules_files)} orchestrator state files")
    print()

    updated_count = 0
    skipped_count = 0

    for filepath in rules_files:
        if process_state_file(filepath):
            updated_count += 1
        else:
            skipped_count += 1

    print()
    print("="* 60)
    print(f"✅ Summary:")
    print(f"   Total states: {len(rules_files)}")
    print(f"   Updated: {updated_count}")
    print(f"   Skipped: {skipped_count}")
    print("="* 60)

    return 0


if __name__ == '__main__':
    sys.exit(main())
