#!/usr/bin/env python3
"""
Comprehensive schema validation for orchestrator-state-v3.json
Tests that all fields in the schema are properly defined and validated.
"""

import json
import sys
from pathlib import Path

def load_json_file(filepath):
    """Load a JSON file and return its contents."""
    with open(filepath, 'r') as f:
        return json.load(f)

def validate_basic_structure(schema_file, state_file):
    """Perform basic structural validation without jsonschema module."""
    schema = load_json_file(schema_file)
    state = load_json_file(state_file)

    errors = []
    warnings = []

    # Check required fields
    required_fields = schema.get('required', [])
    for field in required_fields:
        if field not in state:
            errors.append(f"Missing required field: {field}")

    # Check for unknown fields not in schema
    schema_properties = set(schema.get('properties', {}).keys())
    state_fields = set(state.keys())

    unknown_fields = state_fields - schema_properties
    if unknown_fields:
        warnings.append(f"Unknown fields not in schema: {', '.join(unknown_fields)}")

    # Check that all schema properties are documented
    for prop_name, prop_def in schema.get('properties', {}).items():
        if 'description' not in prop_def and prop_name not in ['_comment']:
            warnings.append(f"Schema property '{prop_name}' lacks description")

    return errors, warnings

def check_schema_coverage(schema_file, example_files):
    """Check that schema covers all fields in example files."""
    schema = load_json_file(schema_file)
    schema_properties = set(schema.get('properties', {}).keys())

    all_example_fields = set()
    for example_file in example_files:
        if Path(example_file).exists():
            example = load_json_file(example_file)
            all_example_fields.update(example.keys())

    missing_in_schema = all_example_fields - schema_properties
    unused_in_examples = schema_properties - all_example_fields

    return missing_in_schema, unused_in_examples

def main():
    # File paths
    schema_file = "orchestrator-state.schema.json"
    example_files = [
        "orchestrator-state-demo.json",
        "integration-state.json.example",
        "orchestrator-fix-cascade-state.json.example",
        "pr-ready-state.json.example"
    ]

    print("=" * 60)
    print("ORCHESTRATOR STATE SCHEMA VALIDATION REPORT")
    print("=" * 60)
    print()

    # Basic validation of example file
    print("1. VALIDATING EXAMPLE FILE STRUCTURE")
    print("-" * 40)
    errors, warnings = validate_basic_structure(schema_file, example_files[0])

    if errors:
        print("❌ ERRORS FOUND:")
        for error in errors:
            print(f"   - {error}")
    else:
        print("✅ No structural errors found")

    if warnings:
        print("\n⚠️  WARNINGS:")
        for warning in warnings:
            print(f"   - {warning}")
    print()

    # Check schema coverage
    print("2. SCHEMA COVERAGE ANALYSIS")
    print("-" * 40)
    missing_in_schema, unused_in_examples = check_schema_coverage(schema_file, example_files)

    if missing_in_schema:
        print("❌ Fields in examples but not in schema:")
        for field in sorted(missing_in_schema):
            print(f"   - {field}")
    else:
        print("✅ All example fields are covered by schema")

    if unused_in_examples:
        print("\n⚠️  Fields in schema but not used in examples:")
        for field in sorted(unused_in_examples):
            print(f"   - {field}")
    print()

    # Schema statistics
    print("3. SCHEMA STATISTICS")
    print("-" * 40)
    schema = load_json_file(schema_file)
    properties = schema.get('properties', {})

    print(f"Total properties defined: {len(properties)}")
    print(f"Required properties: {len(schema.get('required', []))}")

    # Count properties with descriptions
    with_desc = sum(1 for p in properties.values() if 'description' in p)
    print(f"Properties with descriptions: {with_desc}/{len(properties)}")

    # Count complex types
    arrays = sum(1 for p in properties.values() if p.get('type') == 'array')
    objects = sum(1 for p in properties.values() if p.get('type') == 'object')
    print(f"Array properties: {arrays}")
    print(f"Object properties: {objects}")
    print()

    # Check for sub-orchestrator fields (new additions)
    print("4. SUB-ORCHESTRATOR FIELDS VALIDATION")
    print("-" * 40)
    sub_orch_fields = [
        'sub_state_machine',
        'sub_state_history',
        'active_sub_orchestrators',
        'sub_orchestrator_history',
        'pending_sub_orchestrations',
        'sub_orchestrator_monitoring',
        'sub_orchestrator_recovery'
    ]

    found_sub_orch = []
    missing_sub_orch = []
    for field in sub_orch_fields:
        if field in properties:
            found_sub_orch.append(field)
        else:
            missing_sub_orch.append(field)

    if found_sub_orch:
        print(f"✅ Sub-orchestrator fields found: {len(found_sub_orch)}/{len(sub_orch_fields)}")
        for field in found_sub_orch:
            desc = properties[field].get('description', 'No description')[:50] + "..."
            print(f"   - {field}: {desc}")

    if missing_sub_orch:
        print(f"\n❌ Missing sub-orchestrator fields:")
        for field in missing_sub_orch:
            print(f"   - {field}")
    print()

    # Overall assessment
    print("5. OVERALL ASSESSMENT")
    print("-" * 40)

    total_issues = len(errors) + len(missing_in_schema) + len(missing_sub_orch)
    if total_issues == 0:
        print("✅ SCHEMA VALIDATION PROJECT_DONEFUL!")
        print("   The schema is comprehensive and covers all known fields.")
    else:
        print(f"❌ ISSUES FOUND: {total_issues}")
        print("   Please address the errors listed above.")

    return 0 if total_issues == 0 else 1

if __name__ == "__main__":
    sys.exit(main())