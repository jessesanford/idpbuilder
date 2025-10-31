#!/usr/bin/env python3
"""
Fix schema to allow additional properties in runtime objects.

Many nested objects have additionalProperties: false which breaks real usage.
We need to be permissive for operational/runtime fields while keeping strict
validation for state names and required fields.
"""

import json
import sys
from pathlib import Path

def fix_additional_properties(obj, path=""):
    """Recursively find and fix additionalProperties: false in nested objects."""
    changes = []

    if isinstance(obj, dict):
        # If this object has additionalProperties: false, change to true
        if obj.get('additionalProperties') == False:
            # Skip if it's a very specific/controlled object where we want strict validation
            # For example, enum types or very specific metadata
            skip_paths = []  # Add specific paths to skip if needed

            if path not in skip_paths:
                obj['additionalProperties'] = True
                changes.append(f"Changed to true: {path}")

        # Recursively process all nested objects
        for key, value in obj.items():
            new_path = f"{path}.{key}" if path else key
            changes.extend(fix_additional_properties(value, new_path))

    elif isinstance(obj, list):
        for i, item in enumerate(obj):
            changes.extend(fix_additional_properties(item, f"{path}[{i}]"))

    return changes

def main():
    # Paths
    repo_root = Path(__file__).parent.parent
    schema_path = repo_root / 'orchestrator-state.schema.json'

    if not schema_path.exists():
        print(f"ERROR: Schema not found at {schema_path}", file=sys.stderr)
        sys.exit(1)

    print(f"🔧 Fixing additionalProperties in {schema_path}...")

    with open(schema_path, 'r') as f:
        schema = json.load(f)

    changes = fix_additional_properties(schema)

    if changes:
        print(f"\n✓ Made {len(changes)} changes:")
        for change in changes:
            print(f"  • {change}")

        print(f"\n💾 Writing updated schema...")
        with open(schema_path, 'w') as f:
            json.dump(schema, f, indent=2)

        print(f"✅ Schema updated successfully")
    else:
        print("✓ No changes needed")

if __name__ == '__main__':
    main()
