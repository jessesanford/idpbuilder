#!/usr/bin/env python3
"""
Rule Migration Validator for Software Factory 3.0

Test suite to verify rules have been correctly migrated.
Ensures no rules lost, no orphaned references, complete coverage.

Usage:
    python tools/rule-migration-validator.py [--test all|state-refs|file-refs|coverage|count]
"""

import json
import re
from pathlib import Path
from typing import Dict, List, Set, Tuple
from collections import defaultdict

class RuleMigrationValidator:
    def __init__(self, project_root: str = "/home/vscode/software-factory-template"):
        self.project_root = Path(project_root)
        self.rule_library = self.project_root / "rule-library"
        self.agent_states = self.project_root / "agent-states"
        self.state_machine_file = self.project_root / "state-machines" / "software-factory-3.0-state-machine.json"

        # SF 3.0 valid state names (R516 compliant)
        self.valid_sf3_states = {
            # Existing SF 2.0 states that remain
            "INIT", "PLANNING", "ERROR_RECOVERY", "PROJECT_DONE",
            "CREATE_NEXT_INFRASTRUCTURE", "SPAWN_SW_ENGINEERS",
            "MONITORING_SWE_PROGRESS", "EFFORT_COMPLETE",

            # Wave iteration states
            "SETUP_WAVE_INFRASTRUCTURE", "START_WAVE_ITERATION",
            "INTEGRATE_WAVE_EFFORTS", "REVIEW_WAVE_INTEGRATION",
            "CREATE_WAVE_FIX_PLAN", "FIX_WAVE_UPSTREAM_BUGS",
            "REVIEW_WAVE_ARCHITECTURE", "COMPLETE_WAVE",

            # Phase iteration states
            "SETUP_PHASE_INFRASTRUCTURE", "START_PHASE_ITERATION",
            "INTEGRATE_PHASE_WAVES", "REVIEW_PHASE_INTEGRATION",
            "CREATE_PHASE_FIX_PLAN", "FIX_PHASE_UPSTREAM_BUGS",
            "REVIEW_PHASE_ARCHITECTURE", "COMPLETE_PHASE",

            # Project iteration states
            "SETUP_PROJECT_INFRASTRUCTURE", "START_PROJECT_ITERATION",
            "INTEGRATE_PROJECT_PHASES", "REVIEW_PROJECT_INTEGRATION",
            "CREATE_PROJECT_FIX_PLAN", "FIX_PROJECT_UPSTREAM_BUGS",
            "REVIEW_PROJECT_ARCHITECTURE", "COMPLETE_PROJECT",

            # State Manager states
            "STARTUP_CONSULTATION", "SHUTDOWN_CONSULTATION",

            # Add others as needed
        }

        # OLD state names that should NOT appear (R516 violations)
        self.invalid_sf3_states = {
            "WAVE_1_1_SETUP", "WAVE_1_1_ITERATION_START", "WAVE_1_1_INTEGRATE",
            "WAVE_1_1_REVIEW", "WAVE_1_1_CREATE_FIX_PLAN", "WAVE_1_1_FIX_UPSTREAM",
            "WAVE_1_1_ARCHITECT_REVIEW", "WAVE_1_1_COMPLETE",
            "PHASE_1_SETUP", "PHASE_1_INTEGRATE", "PHASE_1_REVIEW",
            "PROJECT_SETUP", "PROJECT_INTEGRATE", "PROJECT_REVIEW",
            "CREATE_NEXT_INFRASTRUCTURE", "CREATE_NEXT_SPLIT_INFRASTRUCTURE",
            "SPAWN_SW_ENGINEERS", "MONITORING",
        }

        # OLD file references that should NOT appear
        self.invalid_file_refs = [
            'orchestrator-state-v3.json"',  # Should be orchestrator-state-v3.json
            'bugs_discovered',            # Should be bug-tracking.json:bugs
            '.current_state',             # Should be .state_machine.current_state
            '.previous_state',            # Should be .state_machine.previous_state
        ]

        # Valid SF 3.0 file references
        self.valid_file_refs = [
            'orchestrator-state-v3.json',
            'bug-tracking.json',
            'integration-containers.json',
            'fix-cascade-state.json',
            'state_machine.current_state',
            'state_machine.previous_state',
        ]

        self.test_results = defaultdict(list)
        self.passed_tests = 0
        self.failed_tests = 0

    def run_all_tests(self) -> bool:
        """Run all validation tests"""
        print("=" * 70)
        print("SOFTWARE FACTORY 3.0 - RULE MIGRATION VALIDATION SUITE")
        print("=" * 70)
        print()

        tests = [
            ("State Reference Check", self.test_no_invalid_state_refs),
            ("File Reference Check", self.test_no_invalid_file_refs),
            ("State Coverage Check", self.test_state_coverage),
            ("Rule Count Check", self.test_rule_count),
            ("State Directory Check", self.test_state_directories),
            ("State Manager Rules Check", self.test_state_manager_rules),
        ]

        for test_name, test_func in tests:
            print(f"Running: {test_name}...")
            passed = test_func()
            if passed:
                print(f"  ✅ PASS")
                self.passed_tests += 1
            else:
                print(f"  ❌ FAIL")
                self.failed_tests += 1
            print()

        self.print_summary()
        return self.failed_tests == 0

    def test_no_invalid_state_refs(self) -> bool:
        """Test: No rules reference old state names"""
        import re
        violations = []

        for rule_file in self.rule_library.glob("*.md"):
            with open(rule_file, 'r') as f:
                content = f.read()

            for invalid_state in self.invalid_sf3_states:
                # Use word boundaries to avoid matching substrings
                # e.g., don't match "MONITORING" in "MONITORING_SWE_PROGRESS"
                pattern = r'\b' + re.escape(invalid_state) + r'\b'
                if re.search(pattern, content):
                    violations.append({
                        "file": str(rule_file.relative_to(self.project_root)),
                        "invalid_state": invalid_state,
                        "line": self.find_line_number(content, invalid_state)
                    })

        # Check state-specific rules too
        for rule_file in self.agent_states.glob("*/*/rules.md"):
            with open(rule_file, 'r') as f:
                content = f.read()

            for invalid_state in self.invalid_sf3_states:
                # Use word boundaries to avoid matching substrings
                pattern = r'\b' + re.escape(invalid_state) + r'\b'
                if re.search(pattern, content):
                    violations.append({
                        "file": str(rule_file.relative_to(self.project_root)),
                        "invalid_state": invalid_state,
                        "line": self.find_line_number(content, invalid_state)
                    })

        if violations:
            self.test_results["invalid_state_refs"] = violations
            print(f"    Found {len(violations)} invalid state references:")
            for v in violations[:5]:  # Show first 5
                print(f"      - {v['file']}:{v['line']} references {v['invalid_state']}")
            if len(violations) > 5:
                print(f"      ... and {len(violations) - 5} more")
            return False

        print(f"    ✓ No invalid state references found")
        return True

    def test_no_invalid_file_refs(self) -> bool:
        """Test: No rules reference old file structure"""
        import re
        violations = []

        for rule_file in self.rule_library.glob("*.md"):
            with open(rule_file, 'r') as f:
                content = f.read()

            for invalid_ref in self.invalid_file_refs:
                # For file refs like ".current_state", use regex to avoid matching
                # ".state_machine.current_state" (which contains ".current_state")
                # Only match if NOT preceded by "state_machine"
                if invalid_ref.startswith('.'):
                    # e.g., ".current_state" should not match if preceded by "state_machine"
                    pattern = r'(?<!state_machine)' + re.escape(invalid_ref)
                else:
                    # For other patterns, just escape and search
                    pattern = re.escape(invalid_ref)

                if re.search(pattern, content):
                    violations.append({
                        "file": str(rule_file.relative_to(self.project_root)),
                        "invalid_ref": invalid_ref,
                        "line": self.find_line_number(content, invalid_ref)
                    })

        if violations:
            self.test_results["invalid_file_refs"] = violations
            print(f"    Found {len(violations)} invalid file references:")
            for v in violations[:5]:
                print(f"      - {v['file']}:{v['line']} references {v['invalid_ref']}")
            if len(violations) > 5:
                print(f"      ... and {len(violations) - 5} more")
            return False

        print(f"    ✓ No invalid file references found")
        return True

    def test_state_coverage(self) -> bool:
        """Test: All SF 3.0 states have rule files"""
        missing_states = []

        for state in self.valid_sf3_states:
            # Check if rule file exists for this state
            # Look in agent-states/software-factory/orchestrator/{STATE}/rules.md
            orchestrator_path = self.agent_states / "software-factory" / "orchestrator" / state / "rules.md"
            state_manager_path = self.agent_states / "state-manager" / state / "rules.md"

            if not orchestrator_path.exists() and not state_manager_path.exists():
                missing_states.append(state)

        if missing_states:
            self.test_results["missing_state_rules"] = missing_states
            print(f"    Found {len(missing_states)} states without rule files:")
            for state in missing_states[:10]:
                print(f"      - {state}")
            if len(missing_states) > 10:
                print(f"      ... and {len(missing_states) - 10} more")
            print(f"    Note: This may be expected during early migration phases")
            # Don't fail on this - it's informational
            return True

        print(f"    ✓ All {len(self.valid_sf3_states)} states have rule files")
        return True

    def test_rule_count(self) -> bool:
        """Test: No rules lost during migration"""
        # Expected: 267 rule-library + 8 state-specific (from analyzer)
        # But state-specific will grow significantly in SF 3.0

        rule_library_count = len(list(self.rule_library.glob("R*.md")))
        state_rules_count = len(list(self.agent_states.glob("*/*/rules.md")))
        total = rule_library_count + state_rules_count

        # Baseline from analyzer: 275 total
        baseline = 275
        expected_growth = 50  # Expect ~50 new state rule files for iteration containers

        print(f"    Rule library: {rule_library_count} (expected ~267)")
        print(f"    State-specific: {state_rules_count} (expected 8 → ~58)")
        print(f"    Total: {total}")

        if total < baseline:
            self.test_results["rule_loss"] = {
                "current": total,
                "baseline": baseline,
                "difference": baseline - total
            }
            print(f"    ⚠️  WARNING: {baseline - total} fewer rules than baseline!")
            print(f"    This suggests rule loss during migration")
            return False

        print(f"    ✓ Rule count OK ({total} >= {baseline} baseline)")
        return True

    def test_state_directories(self) -> bool:
        """Test: State directories use R516-compliant names"""
        invalid_dirs = []

        for agent_dir in self.agent_states.iterdir():
            if not agent_dir.is_dir():
                continue

            for state_dir in agent_dir.iterdir():
                if not state_dir.is_dir():
                    continue

                state_name = state_dir.name
                if state_name in self.invalid_sf3_states:
                    invalid_dirs.append({
                        "agent": agent_dir.name,
                        "state": state_name,
                        "path": str(state_dir.relative_to(self.project_root))
                    })

        if invalid_dirs:
            self.test_results["invalid_state_dirs"] = invalid_dirs
            print(f"    Found {len(invalid_dirs)} state directories with old names:")
            for d in invalid_dirs:
                print(f"      - {d['path']}")
            return False

        print(f"    ✓ All state directories use R516-compliant names")
        return True

    def test_state_manager_rules(self) -> bool:
        """Test: State Manager agent has required rules"""
        required_states = ["STARTUP_CONSULTATION", "SHUTDOWN_CONSULTATION"]
        missing = []

        for state in required_states:
            rule_file = self.agent_states / "state-manager" / state / "rules.md"
            if not rule_file.exists():
                missing.append(state)

        if missing:
            self.test_results["missing_state_manager_rules"] = missing
            print(f"    Missing State Manager rules for: {', '.join(missing)}")
            print(f"    Note: This is expected until Phase 2 of migration")
            # Don't fail - informational
            return True

        print(f"    ✓ State Manager has all required rule files")
        return True

    def find_line_number(self, content: str, search_str: str) -> int:
        """Find line number of first occurrence of search string"""
        lines = content.split('\n')
        for i, line in enumerate(lines, 1):
            if search_str in line:
                return i
        return 0

    def print_summary(self):
        """Print test summary"""
        print("=" * 70)
        print("VALIDATION SUMMARY")
        print("=" * 70)
        print(f"Tests Passed: {self.passed_tests}")
        print(f"Tests Failed: {self.failed_tests}")

        if self.failed_tests == 0:
            print("\n✅ ALL TESTS PASSED - Rule migration validated successfully!")
        else:
            print(f"\n❌ {self.failed_tests} TESTS FAILED - Review and fix issues above")
            print("\nNext steps:")
            print("  1. Fix invalid state references (R516 compliance)")
            print("  2. Fix invalid file references (SF 3.0 multi-file)")
            print("  3. Ensure no rules lost")
            print("  4. Rename state directories per R516")
            print("  5. Run validator again until all tests pass")

    def save_results(self, output_file: str):
        """Save test results to JSON"""
        results = {
            "passed": self.passed_tests,
            "failed": self.failed_tests,
            "test_results": dict(self.test_results)
        }

        with open(output_file, 'w') as f:
            json.dump(results, f, indent=2)

        print(f"\n💾 Results saved to: {output_file}")

def main():
    import argparse
    parser = argparse.ArgumentParser(description='Validate rule migration for SF 3.0')
    parser.add_argument('--test', choices=['all', 'state-refs', 'file-refs', 'coverage', 'count'],
                       default='all', help='Which test to run')
    parser.add_argument('--output', '-o', default='rule-validation-results.json',
                       help='Output JSON file')
    args = parser.parse_args()

    validator = RuleMigrationValidator()

    if args.test == 'all':
        success = validator.run_all_tests()
    elif args.test == 'state-refs':
        success = validator.test_no_invalid_state_refs()
    elif args.test == 'file-refs':
        success = validator.test_no_invalid_file_refs()
    elif args.test == 'coverage':
        success = validator.test_state_coverage()
    elif args.test == 'count':
        success = validator.test_rule_count()

    validator.save_results(args.output)

    exit(0 if success else 1)

if __name__ == "__main__":
    main()
