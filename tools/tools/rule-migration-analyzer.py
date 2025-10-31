#!/usr/bin/env python3
"""
Rule Migration Analyzer for Software Factory 3.0

Analyzes all SF 2.0 rules to determine compatibility with SF 3.0 architecture.
Identifies rules that need updating, deprecation, or are obsolete.

Usage:
    python tools/rule-migration-analyzer.py [--output analysis.json]
"""

import json
import os
import re
from pathlib import Path
from typing import Dict, List, Set, Tuple
from dataclasses import dataclass, asdict
from collections import defaultdict

@dataclass
class RuleAnalysis:
    """Analysis result for a single rule"""
    rule_id: str
    file_path: str
    category: str  # KEEP_AS_IS, UPDATE_NAMES, UPDATE_STRUCTURE, DEPRECATE, ANALYZE_MANUAL
    confidence: str  # HIGH, MEDIUM, LOW
    issues_found: List[str]
    old_state_references: List[str]
    old_file_references: List[str]
    recommendations: List[str]
    criticality: str  # SUPREME_LAW, BLOCKING, STANDARD
    agent: str  # orchestrator, sw-engineer, code-reviewer, architect, universal

class RuleMigrationAnalyzer:
    def __init__(self, project_root: str = "/home/vscode/software-factory-template"):
        self.project_root = Path(project_root)
        self.rule_library = self.project_root / "rule-library"
        self.agent_states = self.project_root / "agent-states"

        # SF 2.0 → SF 3.0 state name mappings (R516 compliance)
        self.state_mappings = {
            # Wave states
            "WAVE_1_1_SETUP": "SETUP_WAVE_INFRASTRUCTURE",
            "WAVE_1_1_ITERATION_START": "START_WAVE_ITERATION",
            "WAVE_1_1_INTEGRATE": "INTEGRATE_WAVE_EFFORTS",
            "WAVE_1_1_REVIEW": "REVIEW_WAVE_INTEGRATION",
            "WAVE_1_1_CREATE_FIX_PLAN": "CREATE_WAVE_FIX_PLAN",
            "WAVE_1_1_FIX_UPSTREAM": "FIX_WAVE_UPSTREAM_BUGS",
            "WAVE_1_1_ARCHITECT_REVIEW": "REVIEW_WAVE_ARCHITECTURE",
            "WAVE_1_1_COMPLETE": "COMPLETE_WAVE",

            # Phase states
            "PHASE_1_SETUP": "SETUP_PHASE_INFRASTRUCTURE",
            "PHASE_1_INTEGRATE": "INTEGRATE_PHASE_WAVES",
            "PHASE_1_REVIEW": "REVIEW_PHASE_INTEGRATION",

            # Project states
            "PROJECT_SETUP": "SETUP_PROJECT_INFRASTRUCTURE",
            "PROJECT_INTEGRATE": "INTEGRATE_PROJECT_PHASES",

            # Other common renames
            "CREATE_NEXT_INFRASTRUCTURE": "CREATE_NEXT_INFRASTRUCTURE",
            "CREATE_NEXT_SPLIT_INFRASTRUCTURE": "CREATE_NEXT_INFRASTRUCTURE",
            "SPAWN_SW_ENGINEERS": "SPAWN_SW_ENGINEERS",
            "MONITORING": "MONITORING_SWE_PROGRESS",
        }

        # SF 2.0 file references that changed
        self.file_mappings = {
            "orchestrator-state-v3.json": "orchestrator-state-v3.json (+ 3 other files)",
            "current_state": "state_machine.current_state",
            "previous_state": "state_machine.previous_state",
            "bugs_discovered": "bug-tracking.json:bugs[]",
            "integration_attempts": "integration-containers.json:active_integrations[]",
        }

        # Rules we know are deprecated in SF 3.0
        self.known_deprecated = set()

        # Rules we know are SF 3.0 specific
        self.sf3_only_rules = {
            "R288",  # Enhanced for 4-file atomicity
            "R510",  # State execution checklists
            "R511",  # Checklist creation protocol
            "R516",  # State naming protocol
        }

        self.results: List[RuleAnalysis] = []

    def analyze_all_rules(self) -> Dict[str, any]:
        """Analyze all rules in rule-library and agent-states"""
        print("=" * 70)
        print("SOFTWARE FACTORY 3.0 - RULE MIGRATION ANALYZER")
        print("=" * 70)
        print()

        # Analyze rule library
        print(f"Analyzing rule library: {self.rule_library}")
        rule_files = list(self.rule_library.glob("R*.md"))
        print(f"Found {len(rule_files)} rule files")

        for rule_file in sorted(rule_files):
            analysis = self.analyze_rule_file(rule_file)
            if analysis:
                self.results.append(analysis)

        # Analyze state-specific rules
        print(f"\nAnalyzing state-specific rules: {self.agent_states}")
        state_rule_files = list(self.agent_states.glob("**/rules.md"))
        print(f"Found {len(state_rule_files)} state-specific rule files")

        for rule_file in sorted(state_rule_files):
            analysis = self.analyze_state_rule_file(rule_file)
            if analysis:
                self.results.append(analysis)

        # Generate summary
        summary = self.generate_summary()
        return summary

    def analyze_rule_file(self, rule_file: Path) -> RuleAnalysis:
        """Analyze a single rule file from rule-library"""
        with open(rule_file, 'r') as f:
            content = f.read()

        # Extract rule ID from filename
        rule_id = rule_file.stem

        # Determine criticality
        criticality = "STANDARD"
        if "SUPREME LAW" in content or "SUPREME_LAW" in content:
            criticality = "SUPREME_LAW"
        elif "BLOCKING" in content:
            criticality = "BLOCKING"

        # Check for old state name references
        old_states = []
        for old_name in self.state_mappings.keys():
            if old_name in content:
                old_states.append(old_name)

        # Check for old file references
        old_files = []
        for old_ref in self.file_mappings.keys():
            if old_ref in content and not old_ref.startswith("R"):  # Avoid false positives
                old_files.append(old_ref)

        # Categorize
        category, confidence, issues, recommendations = self.categorize_rule(
            rule_id, content, old_states, old_files, criticality
        )

        # Determine agent scope
        agent = self.determine_agent_scope(content)

        return RuleAnalysis(
            rule_id=rule_id,
            file_path=str(rule_file.relative_to(self.project_root)),
            category=category,
            confidence=confidence,
            issues_found=issues,
            old_state_references=old_states,
            old_file_references=old_files,
            recommendations=recommendations,
            criticality=criticality,
            agent=agent
        )

    def analyze_state_rule_file(self, rule_file: Path) -> RuleAnalysis:
        """Analyze a state-specific rules.md file"""
        with open(rule_file, 'r') as f:
            content = f.read()

        # Extract agent and state from path
        # Path format: agent-states/{context...}/{agent}/{STATE}/rules.md (variable depth)
        # State is always parent directory of rules.md
        # Agent is the directory before STATE (looking for known agent names)
        parts = rule_file.parts
        state = parts[-2]

        # Find agent by looking for known agent directory names
        agent = "unknown"
        known_agents = ["orchestrator", "sw-engineer", "code-reviewer", "architect", "state-manager"]
        for i in range(len(parts) - 3, -1, -1):
            if parts[i] in known_agents:
                agent = parts[i]
                break

        # If no known agent found, use parent directory of STATE
        if agent == "unknown":
            agent = parts[-3] if len(parts) >= 3 else "unknown"

        rule_id = f"{agent}:{state}:rules"

        # Check if state name needs updating
        needs_rename = state in self.state_mappings
        new_state = self.state_mappings.get(state, state)

        # Check for old references in content
        old_states = []
        for old_name in self.state_mappings.keys():
            if old_name in content:
                old_states.append(old_name)

        old_files = []
        for old_ref in self.file_mappings.keys():
            if old_ref in content and not old_ref.startswith("R"):
                old_files.append(old_ref)

        issues = []
        recommendations = []

        if needs_rename:
            issues.append(f"State directory name needs rename: {state} → {new_state}")
            recommendations.append(f"Create agent-states/{agent}/{new_state}/rules.md")
            recommendations.append(f"Copy content from old location")
            recommendations.append(f"Update all state references in content")

        if old_states:
            issues.append(f"References old state names: {', '.join(old_states)}")
            recommendations.append("Update state name references per R516")

        if old_files:
            issues.append(f"References old file structure: {', '.join(old_files)}")
            recommendations.append("Update file references for SF 3.0 (4 files)")

        # Categorize
        if needs_rename or old_states or old_files:
            category = "UPDATE_NAMES" if needs_rename else "UPDATE_STRUCTURE"
            confidence = "HIGH"
        else:
            category = "KEEP_AS_IS"
            confidence = "HIGH"

        return RuleAnalysis(
            rule_id=rule_id,
            file_path=str(rule_file.relative_to(self.project_root)),
            category=category,
            confidence="HIGH",
            issues_found=issues,
            old_state_references=old_states,
            old_file_references=old_files,
            recommendations=recommendations,
            criticality="STANDARD",
            agent=agent
        )

    def categorize_rule(self, rule_id: str, content: str, old_states: List[str],
                       old_files: List[str], criticality: str) -> Tuple[str, str, List[str], List[str]]:
        """Categorize a rule and provide recommendations"""
        issues = []
        recommendations = []

        # Check if SF 3.0 specific rule
        if any(sf3_rule in rule_id for sf3_rule in self.sf3_only_rules):
            return ("KEEP_AS_IS", "HIGH", ["SF 3.0 specific rule"], ["Already compliant with SF 3.0"])

        # Check for old state references
        if old_states:
            issues.append(f"References old state names: {', '.join(old_states)}")
            recommendations.append("Update state names per R516 compliance")
            for old_state in old_states:
                new_state = self.state_mappings.get(old_state, "UNKNOWN")
                recommendations.append(f"  {old_state} → {new_state}")

        # Check for old file references
        if old_files:
            issues.append(f"References old file structure: {', '.join(old_files)}")
            recommendations.append("Update file references for SF 3.0 multi-file architecture")
            for old_ref in old_files:
                new_ref = self.file_mappings.get(old_ref, old_ref)
                recommendations.append(f"  {old_ref} → {new_ref}")

        # Check for iteration-related content
        if "iteration" in content.lower() and "container" not in content.lower():
            issues.append("Mentions iterations but not iteration containers")
            recommendations.append("Review for SF 3.0 iteration container concept")

        # Check for bug tracking
        if "bugs_discovered" in content or "bug tracking" in content.lower():
            if "bug-tracking.json" not in content:
                issues.append("References bugs but not bug-tracking.json")
                recommendations.append("Update to use bug-tracking.json (SF 3.0)")

        # Check for State Manager
        if "orchestrator" in content.lower() and "state transition" in content.lower():
            if "state manager" not in content.lower() and "state-manager" not in content.lower():
                issues.append("Orchestrator transition rule without State Manager mention")
                recommendations.append("Consider State Manager role in SF 3.0 (bookend pattern)")

        # Determine category
        if not issues:
            category = "KEEP_AS_IS"
            confidence = "HIGH"
        elif old_states and not old_files:
            category = "UPDATE_NAMES"
            confidence = "HIGH"
        elif old_files:
            category = "UPDATE_STRUCTURE"
            confidence = "MEDIUM"
        elif len(issues) > 3:
            category = "ANALYZE_MANUAL"
            confidence = "LOW"
        else:
            category = "UPDATE_STRUCTURE"
            confidence = "MEDIUM"

        # Special handling for critical rules
        if criticality == "SUPREME_LAW" and issues:
            confidence = "LOW"  # Require manual review for supreme laws
            recommendations.insert(0, "⚠️  SUPREME LAW - Require manual expert review")

        return (category, confidence, issues, recommendations)

    def determine_agent_scope(self, content: str) -> str:
        """Determine which agent(s) this rule applies to"""
        content_lower = content.lower()

        agents = []
        if "orchestrator" in content_lower:
            agents.append("orchestrator")
        if "sw-engineer" in content_lower or "software engineer" in content_lower:
            agents.append("sw-engineer")
        if "code-reviewer" in content_lower or "code reviewer" in content_lower:
            agents.append("code-reviewer")
        if "architect" in content_lower:
            agents.append("architect")
        if "state-manager" in content_lower or "state manager" in content_lower:
            agents.append("state-manager")

        if not agents:
            return "universal"
        elif len(agents) == 1:
            return agents[0]
        else:
            return ",".join(agents)

    def generate_summary(self) -> Dict[str, any]:
        """Generate summary statistics"""
        summary = {
            "total_rules": len(self.results),
            "by_category": defaultdict(int),
            "by_confidence": defaultdict(int),
            "by_criticality": defaultdict(int),
            "by_agent": defaultdict(int),
            "rules_with_issues": 0,
            "supreme_laws_needing_update": [],
            "high_priority_updates": [],
            "deprecated_candidates": [],
            "rules": []
        }

        for result in self.results:
            summary["by_category"][result.category] += 1
            summary["by_confidence"][result.confidence] += 1
            summary["by_criticality"][result.criticality] += 1
            summary["by_agent"][result.agent] += 1

            if result.issues_found:
                summary["rules_with_issues"] += 1

            if result.criticality == "SUPREME_LAW" and result.category != "KEEP_AS_IS":
                summary["supreme_laws_needing_update"].append(result.rule_id)

            if result.criticality in ["SUPREME_LAW", "BLOCKING"] and result.issues_found:
                summary["high_priority_updates"].append({
                    "rule_id": result.rule_id,
                    "category": result.category,
                    "issues": result.issues_found
                })

            if result.category == "DEPRECATE":
                summary["deprecated_candidates"].append(result.rule_id)

            summary["rules"].append(asdict(result))

        return summary

    def print_summary(self, summary: Dict[str, any]):
        """Print human-readable summary"""
        print("\n" + "=" * 70)
        print("ANALYSIS SUMMARY")
        print("=" * 70)
        print(f"\nTotal Rules Analyzed: {summary['total_rules']}")
        print(f"Rules with Issues: {summary['rules_with_issues']}")

        print("\nBy Category:")
        for category, count in sorted(summary['by_category'].items()):
            pct = (count / summary['total_rules'] * 100)
            print(f"  {category:20s}: {count:4d} ({pct:5.1f}%)")

        print("\nBy Confidence:")
        for confidence, count in sorted(summary['by_confidence'].items()):
            pct = (count / summary['total_rules'] * 100)
            print(f"  {confidence:20s}: {count:4d} ({pct:5.1f}%)")

        print("\nBy Criticality:")
        for criticality, count in sorted(summary['by_criticality'].items()):
            pct = (count / summary['total_rules'] * 100)
            print(f"  {criticality:20s}: {count:4d} ({pct:5.1f}%)")

        print("\nBy Agent:")
        for agent, count in sorted(summary['by_agent'].items(), key=lambda x: -x[1]):
            pct = (count / summary['total_rules'] * 100)
            print(f"  {agent:20s}: {count:4d} ({pct:5.1f}%)")

        if summary['supreme_laws_needing_update']:
            print(f"\n🔴 CRITICAL: {len(summary['supreme_laws_needing_update'])} Supreme Laws Need Updates:")
            for rule_id in summary['supreme_laws_needing_update']:
                print(f"  - {rule_id}")

        if summary['high_priority_updates']:
            print(f"\n⚠️  HIGH PRIORITY: {len(summary['high_priority_updates'])} Critical/Blocking Rules Need Updates")

        print("\n" + "=" * 70)

    def save_results(self, output_file: str):
        """Save analysis results to JSON file"""
        summary = self.generate_summary()

        with open(output_file, 'w') as f:
            json.dump(summary, f, indent=2)

        print(f"\n✅ Results saved to: {output_file}")

def main():
    import argparse
    parser = argparse.ArgumentParser(description='Analyze rules for SF 3.0 migration')
    parser.add_argument('--output', '-o', default='rule-migration-analysis.json',
                       help='Output JSON file')
    parser.add_argument('--project-root', default='/home/vscode/software-factory-template',
                       help='Project root directory')
    args = parser.parse_args()

    analyzer = RuleMigrationAnalyzer(project_root=args.project_root)
    summary = analyzer.analyze_all_rules()
    analyzer.print_summary(summary)
    analyzer.save_results(args.output)

    print(f"\n📊 Analysis complete!")
    print(f"Next steps:")
    print(f"  1. Review {args.output}")
    print(f"  2. Address Supreme Law updates FIRST")
    print(f"  3. Update high-priority rules")
    print(f"  4. Batch update rules by category")
    print(f"  5. Run analyzer again to verify 100% KEEP_AS_IS")

if __name__ == "__main__":
    main()
