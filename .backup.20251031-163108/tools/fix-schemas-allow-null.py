#!/usr/bin/env python3
"""
Fix JSON schemas to allow null values for optional fields.

This script walks through JSON schema files and converts type declarations
from "type": "string" to "type": ["string", "null"] for fields that are
not in the parent's "required" array.
"""

import json
import sys
from pathlib import Path


def is_field_required(field_name, parent_schema):
    """Check if a field is in the parent's required array."""
    if not isinstance(parent_schema, dict):
        return False
    required_fields = parent_schema.get('required', [])
    return field_name in required_fields


def make_nullable(type_value):
    """Convert a type declaration to allow null."""
    if isinstance(type_value, str):
        return [type_value, "null"]
    elif isinstance(type_value, list):
        if "null" not in type_value:
            return type_value + ["null"]
        return type_value
    return type_value


def process_schema(schema, parent=None, field_name=None, path=""):
    """Recursively process schema to make optional fields nullable."""
    if not isinstance(schema, dict):
        return schema

    # Check if this field has a type and is not required
    if 'type' in schema and field_name and parent:
        if not is_field_required(field_name, parent):
            # Make it nullable
            schema['type'] = make_nullable(schema['type'])

            # If field has an enum, add null to the enum as well
            if 'enum' in schema and isinstance(schema['enum'], list):
                if None not in schema['enum']:
                    schema['enum'].append(None)
                    print(f"  Made nullable (with enum): {path}.{field_name}")
                else:
                    print(f"  Made nullable: {path}.{field_name}")
            else:
                print(f"  Made nullable: {path}.{field_name}")

    # Process properties recursively
    if 'properties' in schema:
        for prop_name, prop_schema in schema['properties'].items():
            schema['properties'][prop_name] = process_schema(
                prop_schema,
                parent=schema,
                field_name=prop_name,
                path=f"{path}.{prop_name}" if path else prop_name
            )

    # Process items (for arrays)
    if 'items' in schema:
        if isinstance(schema['items'], dict):
            schema['items'] = process_schema(
                schema['items'],
                parent=schema,
                field_name='items',
                path=f"{path}.items"
            )
        elif isinstance(schema['items'], list):
            schema['items'] = [
                process_schema(item, parent=schema, field_name=f'items[{i}]', path=path)
                for i, item in enumerate(schema['items'])
            ]

    # Process nested schemas (anyOf, oneOf, allOf)
    for key in ['anyOf', 'oneOf', 'allOf']:
        if key in schema:
            schema[key] = [
                process_schema(sub, parent=schema, field_name=key, path=f"{path}.{key}")
                for sub in schema[key]
            ]

    return schema


def fix_schema_file(schema_path):
    """Fix a single schema file to allow null for optional fields."""
    print(f"\nProcessing: {schema_path.name}")

    with open(schema_path, 'r') as f:
        schema = json.load(f)

    # Process the schema
    fixed_schema = process_schema(schema, path=schema_path.stem)

    # Write back with pretty formatting
    with open(schema_path, 'w') as f:
        json.dump(fixed_schema, f, indent=2)
        f.write('\n')  # Add trailing newline

    print(f"  ✅ Updated {schema_path.name}")


def main():
    """Main entry point."""
    schemas_dir = Path(__file__).parent.parent / 'schemas'

    if not schemas_dir.exists():
        print(f"❌ Schemas directory not found: {schemas_dir}")
        sys.exit(1)

    schema_files = list(schemas_dir.glob('*.schema.json'))

    if not schema_files:
        print(f"❌ No schema files found in {schemas_dir}")
        sys.exit(1)

    print(f"Found {len(schema_files)} schema files")

    for schema_file in sorted(schema_files):
        try:
            fix_schema_file(schema_file)
        except Exception as e:
            print(f"  ❌ Error processing {schema_file.name}: {e}")
            sys.exit(1)

    print(f"\n✅ Successfully processed {len(schema_files)} schema files")


if __name__ == '__main__':
    main()
