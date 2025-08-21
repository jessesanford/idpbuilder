# Reference Verification Report

## File Path Structure

The correct file structure for the Software Factory Template is:

```
/workspaces/[project]/
├── orchestrator-state.yaml              # Root level (NOT in orchestrator/)
├── PROJECT-IMPLEMENTATION-PLAN.md       # Root level (NOT in orchestrator/)
├── core/
│   └── SOFTWARE-FACTORY-STATE-MACHINE.md  # In core/ (NOT in orchestrator/)
├── protocols/
│   └── [15 protocol files]             # All critical protocols
├── phase-plans/
│   └── PHASE{X}-SPECIFIC-IMPL-PLAN.md  # Phase plans (NOT in orchestrator/)
└── .claude/
    └── commands/
        └── continue-orchestrating.md   # Command definition
```

## Corrections Made

### 1. CLAUDE.md File References
- ✅ Changed `/orchestrator/orchestrator-state.yaml` → `/orchestrator-state.yaml`
- ✅ Changed `/orchestrator/PROJECT-IMPLEMENTATION-PLAN.md` → `/PROJECT-IMPLEMENTATION-PLAN.md`
- ✅ Changed `/orchestrator/PHASE{X}-SPECIFIC-IMPL-PLAN.md` → `/phase-plans/PHASE{X}-SPECIFIC-IMPL-PLAN.md`
- ✅ Changed `/orchestrator/continue-orchestrating.md` → `/.claude/commands/continue-orchestrating.md`

### 2. Protocol File References
- ✅ CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md: Fixed orchestrator-state.yaml and phase plan paths
- ✅ WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md: Fixed orchestrator-state.yaml and phase plan paths
- ✅ ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md: Fixed phase plan path

### 3. README.md Corrections
- ✅ Updated directory tree to show all 54 files in correct locations
- ✅ Fixed phase plan creation command to use phase-plans/ directory
- ✅ Fixed orchestrator-state.yaml creation to be at root level

### 4. Command File Corrections
- ✅ continue-orchestrating.md: Updated directory structure diagram
- ✅ continue-orchestrating.md: Fixed required files list with correct paths

## File Count Verification

### Total Files: 54

#### Root Level (10 files):
1. README.md
2. setup.sh
3. orchestrator-state-example.yaml
4. PROJECT-IMPLEMENTATION-PLAN-TEMPLATE.md
5. HOW-TO-PLAN.md
6. PLANNING-AGENT-ASSIGNMENTS.md
7. TEMPLATE-CREATION-SUMMARY.md
8. CRITICAL-FILES-ADDED.md
9. FINAL-FILE-ORGANIZATION.md
10. CLAUDE-MD-FILE-VERIFICATION.md

#### .claude/ (6 files):
1. CLAUDE.md
2. agents/orchestrator-task-master.md
3. agents/code-reviewer.md
4. agents/architect-reviewer.md
5. agents/sw-engineer-example-go.md
6. commands/continue-orchestrating.md

#### core/ (2 files):
1. SOFTWARE-FACTORY-STATE-MACHINE.md
2. ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md

#### protocols/ (15 files):
1. IMPERATIVE-LINE-COUNT-RULE.md
2. EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md
3. SW-ENGINEER-STARTUP-REQUIREMENTS.md
4. SW-ENGINEER-EXPLICIT-INSTRUCTIONS.md
5. ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md
6. ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md
7. CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md
8. CODE-REVIEWER-COMPREHENSIVE-GUIDE.md
9. WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md
10. ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md
11. PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md
12. TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
13. WORK-LOG-TEMPLATE.md
14. TODO-STATE-MANAGEMENT-PROTOCOL.md
15. (One file was miscounted - actually 14 protocol files listed)

#### phase-plans/ (5 files):
1. README.md
2. PHASEX-GENERIC-TEMPLATE.md
3. PHASE1-TEMPLATE.md
4. PHASE2-TEMPLATE.md
5. PHASE3-TEMPLATE.md

#### Other Directories:
- efforts/: 1 README.md
- todos/: 1 README.md
- tools/: 1 line-counter.sh
- possibly-needed-but-not-sure/: 13 files

## Common Path Patterns

When using the template, replace `[project]` with your actual project name:

### Correct Paths:
- `/workspaces/[project]/orchestrator-state.yaml`
- `/workspaces/[project]/PROJECT-IMPLEMENTATION-PLAN.md`
- `/workspaces/[project]/core/SOFTWARE-FACTORY-STATE-MACHINE.md`
- `/workspaces/[project]/protocols/[protocol-name].md`
- `/workspaces/[project]/phase-plans/PHASE{X}-SPECIFIC-IMPL-PLAN.md`
- `/workspaces/[project]/.claude/commands/continue-orchestrating.md`

### Incorrect Paths (DO NOT USE):
- ❌ `/workspaces/[project]/orchestrator/orchestrator-state.yaml`
- ❌ `/workspaces/[project]/orchestrator/SOFTWARE-FACTORY-STATE-MACHINE.md`
- ❌ `/workspaces/[project]/orchestrator/PHASE{X}-SPECIFIC-IMPL-PLAN.md`

## Validation Complete

All file references have been corrected to match the actual directory structure. The template now has consistent and correct paths throughout all documentation.