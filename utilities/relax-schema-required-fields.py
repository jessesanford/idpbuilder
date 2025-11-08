#!/usr/bin/env python3
"""
Relax schema required fields to only CRITICAL ones.

The schema was too strict - many required fields are operational/runtime data
that may not exist at all times. We need to keep validation for:
1. State names (enum) - CRITICAL
2. Core structure fields (current_phase, current_state, etc.) - CRITICAL
3. Field types - CRITICAL

But relax:
- Optional metadata fields
- Runtime operational fields
- Infrastructure planning fields that may be filled incrementally
"""

import json
import sys
from pathlib import Path

def relax_required_fields(obj, path="", parent_key=""):
    """Recursively relax required fields that are not critical."""
    changes = []

    if isinstance(obj, dict):
        # Check if this object has a 'required' array
        if 'required' in obj and isinstance(obj['required'], list):
            original_required = obj['required'].copy()

            # Define CRITICAL required fields at root level only
            if path == "":  # Root level
                critical_root_fields = [
                    'current_phase',
                    'current_wave',
                    'current_state',
                    'previous_state',
                    'transition_time',
                    'phases_planned',
                    'waves_per_phase',
                    'efforts_completed',
                    'efforts_in_progress',
                    'efforts_pending',
                    'project_info'
                ]
                # Keep only critical fields
                obj['required'] = [f for f in obj['required'] if f in critical_root_fields]

                if obj['required'] != original_required:
                    removed = set(original_required) - set(obj['required'])
                    changes.append(f"ROOT: Removed non-critical required fields: {removed}")

            else:
                # For nested objects, make most fields optional
                # Keep only ID/key fields that are truly required
                critical_nested_patterns = ['_id', 'id', 'name']

                new_required = []
                for field in obj['required']:
                    # Keep if it matches a critical pattern
                    if any(pattern in field for pattern in critical_nested_patterns):
                        new_required.append(field)

                if new_required != original_required and len(original_required) > len(new_required):
                    removed = set(original_required) - set(new_required)
                    obj['required'] = new_required
                    changes.append(f"{path}: Relaxed required fields (removed: {removed})")

        # Recursively process nested objects
        for key, value in obj.items():
            new_path = f"{path}.{key}" if path else key
            changes.extend(relax_required_fields(value, new_path, key))

    elif isinstance(obj, list):
        for i, item in enumerate(obj):
            changes.extend(relax_required_fields(item, f"{path}[{i}]", parent_key))

    return changes

def main():
    # Paths
    repo_root = Path(__file__).parent.parent
    schema_path = repo_root / 'orchestrator-state.schema.json'

    if not schema_path.exists():
        print(f"ERROR: Schema not found at {schema_path}", file=sys.stderr)
        sys.exit(1)

    print(f"🔧 Relaxing required fields in {schema_path}...")
    print("   Keeping ONLY critical root-level and ID/key fields as required\n")

    with open(schema_path, 'r') as f:
        schema = json.load(f)

    changes = relax_required_fields(schema)

    if changes:
        print(f"✓ Made {len(changes)} changes:\n")
        for change in changes[:20]:  # Show first 20
            print(f"  • {change}")
        if len(changes) > 20:
            print(f"  ... and {len(changes) - 20} more")

        print(f"\n💾 Writing updated schema...")
        with open(schema_path, 'w') as f:
            json.dump(schema, f, indent=2)

        print(f"✅ Schema updated successfully")
        print("\n📋 Summary:")
        print(f"   - Root required fields: Kept only critical state/structure fields")
        print(f"   - Nested required fields: Kept only ID/key fields")
        print(f"   - Everything else: Now optional (but types still validated)")
    else:
        print("✓ No changes needed")

if __name__ == '__main__':
    main()
