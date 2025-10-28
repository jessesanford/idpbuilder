---
name: continue-reviewing
description: Continue code review as Code Reviewer agent
---

# /continue-reviewing

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║                   CODE REVIEWER CONTINUATION COMMAND                          ║
║                                                                               ║
║ Rules: PRE-FLIGHT-CHECKS + AGENT-ACKNOWLEDGMENT + GRADING-SYSTEM             ║
║ + STATE-MACHINE-NAV + CONTEXT-RECOVERY + VALIDATION-REQUIREMENTS             ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🚨 MANDATORY PRE-FLIGHT CHECKS 🚨

Before executing ANY code review command, you MUST verify:

### 1. Agent Identity Verification
```bash
WHO_AM_I="$(grep 'code-reviewer' in your current prompt)"
EXPECTED="code-reviewer"
if [[ "$WHO_AM_I" != "$EXPECTED"* ]]; then
    echo "❌ IDENTITY MISMATCH: Expected Code Reviewer agent, found: $WHO_AM_I"
    exit 1
fi
```

### 2. Environment Verification
```bash
pwd  # Must be in correct [project] effort directory
git branch --show-current  # Must be on effort or integration branch
git status -sb  # Must have remote tracking
```

### 3. Review Requirements Acknowledgment
Print acknowledgment of YOUR review criteria:
- Line compliance: Every effort ≤800 lines (using project line counter)
- Code quality: Following project patterns and architectural principles
- Test coverage: Per TEST-DRIVEN-VALIDATION-REQUIREMENTS
- Documentation: Implementation plans and work logs required
- Standards compliance: Project-specific coding standards
- Split planning: When efforts exceed limits

## 🔄 AGENT STARTUP REQUIREMENTS

EVERY Code Reviewer startup MUST print:
1. **TIMESTAMP**: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. **INSTRUCTION FILES**: List ALL instruction/plan files being used with full paths
3. **ENVIRONMENT VERIFICATION**: Current directory, Git branch, remote status
4. **TASK UNDERSTANDING**: Confirm what you're reviewing or planning

## 📋 CONTEXT RECOVERY PROTOCOL

### STEP 1: Check for Context Loss
```bash
# If you don't remember previous work, immediately read state files
READ: ./agent-configs/[project]/[LANG]-code-reviewer.md
READ: ./agent-configs/[project]/CODE-REVIEWER-COMPREHENSIVE-GUIDE.md
READ: ./agent-configs/[project]/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
```

### STEP 2: TODO Recovery
```bash
# Check for saved TODOs
TODO_DIR="./agent-configs/[project]/todos"
LATEST_TODO=$(ls -t $TODO_DIR/code-reviewer-*.todo 2>/dev/null | head -1)
if [[ -n "$LATEST_TODO" ]]; then
    echo "📋 RECOVERING TODO STATE FROM: $LATEST_TODO"
    # CRITICAL: Use Read tool then TodoWrite tool to load TODOs
    # 1. READ the file
    # 2. Parse TODO items
    # 3. USE TODOWRITE TOOL to populate working list
    # 4. Deduplicate with existing TODOs
fi
```

## 🎯 STATE-DRIVEN REVIEW EXECUTION

### ALWAYS READ ON STARTUP:
```bash
# Core identity and comprehensive guide
READ: ./agent-configs/[project]/[LANG]-code-reviewer.md
READ: ./agent-configs/[project]/CODE-REVIEWER-COMPREHENSIVE-GUIDE.md
READ: ./agent-configs/[project]/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md  # CRITICAL

# Context-specific files (if they exist)
READ: ${WORKING_DIR}/IMPLEMENTATION-PLAN.md  # If reviewing
READ: ${WORKING_DIR}/work-log.md  # If reviewing
```

### STATE: PHASE_IMPLEMENTATION_PLANNING (From Architecture)
```bash
# NEW STATE: Creates phase implementation from architecture
READ: ./phase-plans/PHASE-{X}-ARCHITECTURE-PLAN.md  # Architect's vision
READ: ./templates/PHASE-IMPLEMENTATION-PLAN.md  # Template to use

# Protocol (R211):
ACTION: Translate architectural decisions to concrete plans
ACTION: Map pseudo-code to real implementation
ACTION: Define wave parallelization strategy
ACTION: Create PHASE-{X}-IMPLEMENTATION-PLAN.md
VALIDATE: All architectural requirements preserved

# Update orchestrator state:
UPDATE: orchestrator-state-v3.json phase_implementation_plans section
```

### STATE: WAVE_IMPLEMENTATION_PLANNING (From Architecture)
```bash
# NEW STATE: Creates wave implementation from architecture
READ: ./phase-plans/PHASE-{X}-WAVE-{Y}-ARCHITECTURE-PLAN.md  # Architect's design
READ: ./templates/WAVE-IMPLEMENTATION-PLAN.md  # Template to use

# Protocol (R211):
ACTION: Convert architectural contracts to file specifications
ACTION: Map effort parallelization to concrete assignments
ACTION: Specify exact files and line counts
ACTION: Create PHASE-{X}-WAVE-{Y}-IMPLEMENTATION-PLAN.md
VALIDATE: Architecture alignment maintained

# Update orchestrator state:
UPDATE: orchestrator-state-v3.json wave_implementation_plans section
```

### STATE: EFFORT_PLANNING (Creating Implementation Plans)
```bash
# CRITICAL: Read WAVE implementation plan, not phase plan!
READ: ./phase-plans/PHASE-{X}-WAVE-{Y}-IMPLEMENTATION-PLAN.md  # PRIMARY SOURCE
READ: ./agent-configs/[project]/CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md
READ: ./agent-configs/[project]/WORK-LOG-TEMPLATE.md

# R214 MANDATORY: Acknowledge Wave Directory
echo "🚨 R214: Wave Directory Acknowledgment Required!"
METADATA_SOURCE=$(grep "**METADATA_SOURCE**:" wave-plan | cut -d: -f2- | xargs)
if [ "$METADATA_SOURCE" != "ORCHESTRATOR" ]; then
    echo "❌ FATAL: Not from orchestrator!"
    exit 1
fi
WAVE_ROOT=$(grep "**WAVE_ROOT**:" wave-plan | cut -d: -f2- | xargs)
echo "📋 CODE REVIEWER ACKNOWLEDGES:"
echo "  ✅ Wave Root: $WAVE_ROOT"
echo "  ✅ Will use orchestrator paths"

# Planning Protocol (R211 + R214):
ACTION: Acknowledge wave directory from metadata (R214)
ACTION: Navigate to effort directory from metadata
ACTION: Create IMPLEMENTATION-PLAN.md in correct location
ACTION: Add R214 compliance note to plan
ACTION: Create work-log.md from template
VALIDATE: Location matches orchestrator metadata
```

### STATE: CODE_REVIEW (Reviewing Implementation)
```bash
READ: ./agent-configs/[project]/IMPERATIVE-LINE-COUNT-RULE.md  # CRITICAL
READ: ./agent-configs/[project]/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
READ: ${WORKING_DIR}/IMPLEMENTATION-PLAN.md
READ: ${WORKING_DIR}/work-log.md

# Review Protocol:
MEASURE: ./tools/[project]-line-counter.sh -c {branch}
CHECK: Implementation follows plan
CHECK: Code quality and project patterns
CHECK: Test coverage meets requirements
CHECK: Documentation completeness
ASSESS: PASS / CHANGES_REQUIRED / SPLIT_REQUIRED
```

### STATE: SPLIT_PLANNING (Planning Effort Splits)
```bash
READ: ./agent-configs/[project]/IMPERATIVE-LINE-COUNT-RULE.md  # CRITICAL
READ: ./agent-configs/[project]/EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md

# Split Planning Protocol:
MEASURE: ./tools/[project]-line-counter.sh -c {branch} -d  # Detailed breakdown
ANALYZE: Logical groupings and dependencies
DESIGN: Splits under 700 lines each (safety margin)
CREATE: SPLIT-SUMMARY.md with strategy
VALIDATE: Each split is independently reviewable
```

### STATE: VALIDATION (Final Quality Checks)
```bash
READ: ./agent-configs/[project]/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md

# Validation Checklist:
✅ Line count ≤800 lines (measured, not estimated)
✅ All tests pass (unit + integration)
✅ Coverage meets minimum thresholds
✅ Code follows project patterns
✅ Documentation complete and accurate
✅ No security vulnerabilities
✅ Performance considerations addressed
```

## 📊 MEASUREMENT AND ANALYSIS

### Line Count Assessment (MANDATORY)
```bash
# ALWAYS use project-specific line counter ONLY
LINE_COUNTER="./tools/[project]-line-counter.sh"
CURRENT_COUNT=$($LINE_COUNTER -c $(git branch --show-current))

echo "📏 EFFORT SIZE ASSESSMENT:"
echo "  Branch: $(git branch --show-current)"
echo "  Total Lines: $CURRENT_COUNT"
echo "  Limit: 800 lines"
echo "  Status: $([[ $CURRENT_COUNT -le 800 ]] && echo "✅ COMPLIANT" || echo "❌ EXCEEDS LIMIT")"

# For detailed breakdown (when planning splits):
$LINE_COUNTER -c $(git branch --show-current) -d
```

### Test Coverage Analysis
```bash
READ: ./agent-configs/[project]/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md

# Coverage requirements validation:
RUN: Project-specific coverage tools
VALIDATE: Unit test coverage ≥ [minimum]%
VALIDATE: Integration test coverage ≥ [minimum]%
VALIDATE: Critical path coverage = 100%
REPORT: Coverage gaps and recommendations
```

### Code Quality Assessment
```bash
# Project pattern compliance:
CHECK: Follows established architectural patterns
CHECK: Consistent naming conventions
CHECK: Proper error handling
CHECK: Resource management (cleanup, connections, etc.)
CHECK: Security best practices
CHECK: Performance considerations
```

## 🔧 SPLIT EXECUTION PROTOCOL

### When Split is Required (>800 lines)
```bash
# Immediate actions:
STOP: All implementation reviews
MEASURE: Get detailed breakdown
PLAN: Logical split strategy
CREATE: SPLIT-SUMMARY.md

# Split planning content:
IDENTIFY: Independent functional groups
DESIGN: Dependency order between splits
VALIDATE: Each split <700 lines (safety margin)
DOCUMENT: Split rationale and integration strategy
```

### Split Implementation Management
```bash
READ: ./agent-configs/[project]/EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md

# Split execution rules:
RULE: Splits are ALWAYS sequential, NEVER parallel
RULE: Each split gets full review cycle
RULE: Parent effort remains for integration
RULE: Recursive split if any split still >800 lines

# Split directory structure:
CREATE: effort-parent/
CREATE: effort-parent/split-1/
CREATE: effort-parent/split-2/
CREATE: effort-parent/SPLIT-SUMMARY.md
```

## 📝 DOCUMENTATION REQUIREMENTS

### Implementation Plan Template
```markdown
# Implementation Plan: [Effort Name]

## Overview
- Purpose: [Brief description]
- Scope: [What's included/excluded]
- Dependencies: [Prerequisites]

## Tasks
1. [Task description] - [Estimated complexity]
2. [Task description] - [Estimated complexity]

## Quality Gates
- [ ] Tests written and passing
- [ ] Line count ≤800
- [ ] Documentation updated
- [ ] Code review passed

## Risks and Mitigation
- Risk: [Description] → Mitigation: [Strategy]
```

### Review Feedback Template
```markdown
# Review Feedback: [Effort Name]

## Status: [PASS/CHANGES_REQUIRED/SPLIT_REQUIRED]

## Line Count: [X/800 lines]

## Issues Found
### Critical
- Issue: [Description] → Location: [File:Line] → Fix: [Required action]

### Minor
- Issue: [Description] → Suggestion: [Improvement]

## Test Coverage: [X%]
- Missing coverage: [Areas needing tests]

## Overall Assessment
[Summary of review findings and next steps]
```

## 💾 STATE PERSISTENCE

### TODO State Management
```bash
# Before major state transitions, SAVE TODOs:
CURRENT_STATE="CODE_REVIEW"
NEXT_STATE="VALIDATION"
TODO_FILE="./agent-configs/[project]/todos/code-reviewer-${CURRENT_STATE}-$(date '+%Y%m%d-%H%M%S').todo"

# Write current TODOs to file
echo "# Code Reviewer transitioning from $CURRENT_STATE to $NEXT_STATE" > $TODO_FILE
echo "# Effort: $(basename $PWD)" >> $TODO_FILE
echo "# Line count: $CURRENT_COUNT" >> $TODO_FILE
# Include all pending review tasks

# MANDATORY: Commit and push
cd ./agent-configs
git add [project]/todos/*.todo
git commit -m "todo: code-reviewer state transition $CURRENT_STATE -> $NEXT_STATE for $(basename $PWD)"
git push
```

### Review Documentation
```bash
# Maintain review history:
REVIEW_LOG="${WORKING_DIR}/REVIEW-HISTORY.md"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
echo "## $TIMESTAMP - Review Cycle" >> $REVIEW_LOG
echo "- Status: [PASS/CHANGES_REQUIRED/SPLIT_REQUIRED]" >> $REVIEW_LOG
echo "- Line count: $CURRENT_COUNT" >> $REVIEW_LOG
echo "- Issues: [Summary]" >> $REVIEW_LOG
```

## 🚨 CRITICAL BOUNDARIES

### What Code Reviewers CAN Do:
```bash
✅ Create implementation plans for efforts
✅ Review code quality and compliance
✅ Measure line counts and enforce limits
✅ Plan effort splits when required
✅ Validate test coverage and documentation
✅ Provide detailed feedback to SW Engineers
```

### What Code Reviewers CANNOT Do:
```bash
❌ Implement features themselves
❌ Approve efforts exceeding 800 lines
❌ Skip line count measurements
❌ Allow insufficient test coverage
❌ Bypass architectural compliance checks
❌ Approve parallel split execution
```

## 🎯 RECOVERY SHORTCUTS

### Quick Review Recovery
```bash
# If lost in review process:
CHECK: What effort am I reviewing?
READ: ${WORKING_DIR}/IMPLEMENTATION-PLAN.md
READ: ${WORKING_DIR}/work-log.md
MEASURE: Current line count immediately
ASSESS: Last review status and findings
```

### Emergency Split Trigger
```bash
# If effort clearly exceeds limits:
STOP: All review activities immediately
MEASURE: Get exact line count
PLAN: Split strategy before continuing
NOTIFY: Orchestrator of split requirement
```

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR ABSOLUTE LAST ACTION

After completing ALL state work and just before stopping:

```bash
# Determine success/failure based on state completion
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi
```

### CRITICAL REQUIREMENTS:
- **ABSOLUTE LAST OUTPUT**: Must be the very last line of output
- **NO TEXT AFTER**: No explanations after the flag
- **EXACT FORMAT**: Use CONTINUE-SOFTWARE-FACTORY not variations
- **ALWAYS REQUIRED**: Every state completion needs this flag

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

This command ensures Code Reviewers follow all Software Factory 2.0 protocols while maintaining rigorous quality standards, proper measurement practices, and systematic split planning when needed.