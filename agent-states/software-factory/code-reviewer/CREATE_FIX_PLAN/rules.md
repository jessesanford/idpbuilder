# CREATE_FIX_PLAN State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 📋 PRIMARY DIRECTIVES FOR CREATE_FIX_PLAN STATE

### Core Mandatory Rules (ALL code-reviewer states must have these):

1. **🚨🚨🚨 R006** - CODE REVIEWER NEVER WRITES CODE (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-code-reviewer-never-writes-code.md`
   - Criticality: BLOCKING
   - Summary: Code Reviewer ONLY creates plans and instructions, NEVER implements code

2. **🔴🔴🔴 R510** - STATE EXECUTION CHECKLIST COMPLIANCE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R510-state-execution-checklist-compliance.md`
   - Criticality: SUPREME LAW
   - Summary: MUST complete and acknowledge every checklist item

3. **🔴🔴🔴 R321** - IMMEDIATE BACKPORT DURING INTEGRATE_WAVE_EFFORTS PROTOCOL (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
   - Criticality: SUPREME LAW
   - Summary: ALL fixes discovered during integration MUST be documented in source effort branches IMMEDIATELY

4. **🔴🔴🔴 R266** - UPSTREAM BUG DOCUMENTATION AND FIX PROTOCOL (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R266-upstream-bug-documentation.md`
   - Criticality: SUPREME LAW
   - Summary: Document bugs formally, then coordinate fixes through orchestrator

5. **🔴🔴🔴 R287** - TODO PERSISTENCE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME LAW
   - Summary: Save TODOs within 30s of TodoWrite and before all transitions

6. **🔴🔴🔴 R288** - STATE FILE UPDATE AND COMMIT PROTOCOL (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME LAW
   - Summary: Update state file correctly and commit immediately

7. **🔴🔴🔴 R322** - MANDATORY STOP BEFORE STATE TRANSITIONS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW
   - Summary: STOP after completing state work, do NOT auto-continue

8. **🔴🔴🔴 R355** - CODE REVIEWER VALIDATION REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R355-code-reviewer-validation-requirements.md`
   - Criticality: SUPREME LAW
   - Summary: All review outputs must be validated and complete

### State-Specific Rules:

9. **🚨🚨🚨 R533** - ARTIFACT LOCATION REPORTING PROTOCOL (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R533-artifact-location-reporting-protocol.md`
   - Criticality: BLOCKING
   - Summary: ALL artifacts MUST be tracked in orchestrator-state-v3.json with complete metadata

10. **⚠️⚠️⚠️ R383** - METADATA FILE STANDARDS (WARNING)
    - File: `$CLAUDE_PROJECT_DIR/rule-library/R383-metadata-file-standards.md`
    - Criticality: WARNING
    - Summary: Fix plans must follow proper metadata format and naming

11. **⚠️⚠️⚠️ R343** - STRUCTURED OUTPUT REQUIREMENTS (WARNING)
    - File: `$CLAUDE_PROJECT_DIR/rule-library/R343-structured-output-requirements.md`
    - Criticality: WARNING
    - Summary: All outputs must be properly structured and parseable

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Analyze integration review findings and identify all bugs/issues
  - Source: Read integration review report (INTEGRATE_PHASE_WAVES_REPORT.md, WAVE_INTEGRATE_WAVE_EFFORTS_REPORT.md, or CODE-REVIEW-REPORT.md)
  - Action: Extract all issues, categorize by severity (CRITICAL/HIGH/MEDIUM/LOW)
  - Validation: All issues from review are documented in fix plan
  - **BLOCKING**: Cannot create fix plan without understanding all issues

- [ ] 2. Identify source effort branches for each bug per R321
  - Requirement: Map each bug to the originating effort branch
  - Context: Integration bugs must be fixed in source branches, not integration
  - Validation: Each bug has identified source_effort_branch field
  - **BLOCKING**: Cannot backport without knowing source branches

- [ ] 3. Create comprehensive fix plan document
  - Location: `.software-factory/phase{N}/FIX-PLAN-{SCOPE}--{TIMESTAMP}.md`
  - Format: Follow R383 metadata standards
  - Content: Root cause analysis, fix instructions, execution order, validation steps
  - Validation: Fix plan file exists and is complete
  - **BLOCKING**: Fix plan is the deliverable for this state

- [ ] 3.5. Record fix plan location in orchestrator-state-v3.json per R533
  - Action: Update `.artifacts.fix_plans` with complete metadata
  - Required Fields: file_path, created_at, created_by, artifact_type, scope, status, related_bugs
  - Validation: `jq '.artifacts.fix_plans.[artifact_id]' orchestrator-state-v3.json` returns metadata
  - **BLOCKING**: Orchestrator cannot discover fix plan without this (R340/R533 compliance)

### STANDARD EXECUTION TASKS (Required)

- [ ] 4. Create upstream bug files in source effort branches per R321
  - Requirement: For EACH bug, create bug file in source effort branch root
  - Location: `efforts/phase{N}/wave{W}/{effort}/BUG-{NNN}-{DESCRIPTION}.md`
  - Format: Detailed bug documentation with reproduction steps
  - Validation: Bug file created in effort branch for each identified bug
  - **R321 ENFORCEMENT**: Integration branches are READ-ONLY, fixes go to source

- [ ] 5. Document bugs in bug-tracking.json tracking system per R266
  - Action: Add entries to `bugs_discovered` array
  - Fields: bug_id, severity, discovered_in, source_efforts, fix_plan, backport_required, status
  - Example:
    ```json
    {
      "bug_id": "BUG-001-DUPLICATE-PUSHCMD",
      "severity": "CRITICAL",
      "discovered_in": "phase1-integration",
      "source_efforts": ["E1.2.1-command-structure"],
      "fix_plan": ".software-factory/phase1/FIX-PLAN-PHASE-1--20251007-120000.md",
      "backport_required": true,
      "status": "DOCUMENTED"
    }
    ```
  - Validation: All bugs tracked in state file
  - **R266 ENFORCEMENT**: Formal bug tracking before fixes

- [ ] 6. Create backport execution plan if integration issues affect upstream
  - Context: If bugs were found during integration, document backport sequence
  - Content: List affected effort branches, sequence for applying fixes, re-integration plan
  - Validation: Backport plan section exists in fix plan document
  - **R321 REQUIREMENT**: Immediate backport protocol

### EXIT REQUIREMENTS (Must complete before transition)

**NOTE**: These are STANDARD across ALL states - copy exactly

- [ ] 7. Update state file to COMPLETED per R288
  - Field: `current_state`
  - Value: `"COMPLETED"`
  - Also update: `previous_state`, `transition_time`, `transition_reason`
  - Validation: `jq '.state_machine.current_state' orchestrator-state-v3.json` shows COMPLETED

- [ ] 8. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "R510_CHECKLIST_COMPLETE"
  - Format: `todos/code-reviewer-CREATE_FIX_PLAN-{YYYYMMDD-HHMMSS}.todo`
  - Validation: TODO file exists and contains current state

- [ ] 9. Commit all changes with descriptive message
  - Include: Fix plan creation details
  - Include: Bug tracking updates
  - Include: Upstream bug file creation
  - Include: Rule compliance references (R321, R266, R288, R287, R510)
  - Format: Multi-line commit message with context

- [ ] 10. Push changes to remote
  - Remote: `origin`
  - Branch: Current branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 11. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (if fix plan created successfully, factory can continue)
  - Value: `FALSE` (if catastrophic error prevents fix planning)
  - Context: Fix plan creation is normal workflow, use TRUE unless system failure
  - **NOTE**: R322 checkpoints = TRUE (agent stops but factory continues)

- [ ] 12. Display checkpoint message per R322
  - Format: Clear message about fix plan creation
  - Include: Location of fix plan document
  - Include: Number of bugs documented
  - Include: Next steps for orchestrator (spawn SW Engineers to apply fixes)

- [ ] 13. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Per: R322 (checkpoint state)

---

## State Purpose

The CREATE_FIX_PLAN state is where Code Reviewer analyzes integration review findings (from phase, wave, or effort integration) and creates comprehensive fix plans. This state enforces R321 (Immediate Backport) by documenting ALL bugs in upstream source effort branches, and R266 (Bug Documentation) by creating formal bug tracking in bug-tracking.json.

**Primary Goal:** Create actionable fix plans that enable SW Engineers to fix bugs in source branches
**Key Actions:** Analyze issues, create fix plans, document bugs in source branches, track in state file
**Success Outcome:** Complete fix plan with upstream bug files and formal bug tracking

---

## Entry Criteria

- **From**: INIT (after orchestrator spawns Code Reviewer for fix planning)
- **Condition**: Integration review has identified bugs/issues requiring fixes
- **Required**:
  - Integration review report exists (INTEGRATE_PHASE_WAVES_REPORT.md, WAVE_INTEGRATE_WAVE_EFFORTS_REPORT.md, or CODE-REVIEW-REPORT.md)
  - Bugs/issues documented in review report
  - Orchestrator spawned Code Reviewer with CREATE_FIX_PLAN state
  - Source effort branches are identifiable

---

## State Actions

### 1. Analyze Integration Review Findings

Read the integration review report and extract all bugs/issues:

```bash
# Determine scope (phase, wave, or effort)
SCOPE=$(jq -r '.state_machine.fix_plan_scope // "effort"' orchestrator-state-v3.json)
PHASE=$(jq -r '.state_machine.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.state_machine.current_wave // empty' orchestrator-state-v3.json)
EFFORT=$(jq -r '.state_machine.current_effort // empty' orchestrator-state-v3.json)

# Find integration review report
if [ "$SCOPE" = "phase" ]; then
  REVIEW_REPORT=".software-factory/phase${PHASE}/INTEGRATE_PHASE_WAVES_REPORT.md"
elif [ "$SCOPE" = "wave" ]; then
  REVIEW_REPORT=".software-factory/phase${PHASE}/wave${WAVE}/WAVE_INTEGRATE_WAVE_EFFORTS_REPORT.md"
else
  REVIEW_REPORT="efforts/phase${PHASE}/wave${WAVE}/${EFFORT}/.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT}/CODE-REVIEW-REPORT--*.md"
fi

# Read and analyze
echo "📖 Reading review report: $REVIEW_REPORT"
# Extract issues, categorize by severity, identify root causes
```

### 2. Create Fix Plan Document

Generate comprehensive fix plan following R383 metadata standards:

```bash
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
FIX_PLAN_FILE=".software-factory/phase${PHASE}/FIX-PLAN-${SCOPE^^}-${PHASE}--${TIMESTAMP}.md"

# Create fix plan with:
# - Executive summary
# - Bug inventory (all issues with severity)
# - Root cause analysis for each bug
# - Fix instructions (detailed, actionable)
# - Execution order (sequence for applying fixes)
# - Validation steps (how to verify fixes)
# - Backport plan (if integration issues)
```

### 3. Create Upstream Bug Files per R321

For EACH bug identified, create bug documentation in source effort branch:

```bash
for bug in $BUG_LIST; do
  BUG_ID=$(extract_bug_id $bug)
  SOURCE_EFFORT=$(identify_source_effort $bug)
  BUG_DESC=$(extract_description $bug)

  # Create bug file in source effort
  BUG_FILE="efforts/phase${PHASE}/wave${WAVE}/${SOURCE_EFFORT}/BUG-${BUG_ID}-${BUG_DESC}.md"

  cat > "$BUG_FILE" << 'EOF'
# Bug: {BUG_ID} - {DESCRIPTION}

## Severity: {SEVERITY}
## Type: {TYPE}
## Discovered During: {SCOPE} integration

## Location
- File: {FILE}
- Line: {LINE}
- Function: {FUNCTION}

## Description
{DETAILED_DESCRIPTION}

## Error Output
```
{ERROR_TRACE}
```

## Root Cause Analysis
{ROOT_CAUSE}

## Recommended Fix
{FIX_INSTRUCTIONS}

## Validation Steps
{VALIDATION_STEPS}

## Impact
- Affects: {AFFECTED_FUNCTIONALITY}
- Integration Impact: {INTEGRATE_WAVE_EFFORTS_IMPACT}
EOF

  echo "✅ Created bug file: $BUG_FILE"
done
```

### 4. Update bug-tracking.json Bug Tracking per R266

Add formal bug tracking entries:

```bash
# For each bug, add to bugs_discovered array
jq --arg bug_id "$BUG_ID" \
   --arg severity "$SEVERITY" \
   --arg discovered_in "${SCOPE}-integration" \
   --argjson source_efforts '["'$SOURCE_EFFORT'"]' \
   --arg fix_plan "$FIX_PLAN_FILE" \
   '.bugs_discovered += [{
     "bug_id": $bug_id,
     "severity": $severity,
     "discovered_in": $discovered_in,
     "source_efforts": $source_efforts,
     "fix_plan": $fix_plan,
     "backport_required": true,
     "status": "DOCUMENTED"
   }]' bug-tracking.json > tmp && mv tmp bug-tracking.json
```

---

## Exit Criteria

### Success Path → COMPLETED

- Fix plan document created with all required sections
- All bugs documented in upstream source effort branches per R321
- All bugs tracked in bug-tracking.json per R266
- Backport plan created if integration issues exist
- All checklist items completed and acknowledged
- State file updated to COMPLETED
- TODOs saved per R287
- Changes committed and pushed

### Failure Scenarios

- **Integration review report missing or unreadable** → ERROR_RECOVERY
  - Condition: Cannot find or parse review report
  - Action: Log error, set CONTINUE-SOFTWARE-FACTORY=FALSE
  - Recovery: Orchestrator must verify integration review completed

- **Cannot identify source effort branches** → ERROR_RECOVERY
  - Condition: Bug source unclear or effort branches missing
  - Action: Document issue, escalate to orchestrator
  - Recovery: Manual investigation required

- **State file corruption** → ERROR_RECOVERY
  - Condition: Cannot update bug-tracking.json
  - Action: Log error, preserve data, stop
  - Recovery: Restore from backup or manual intervention

---

## Rules Enforced

- R321: Immediate Backport During Integration - ALL bugs documented in source effort branches
- R266: Upstream Bug Documentation - Formal bug tracking before fixes
- R006: Code Reviewer Never Writes Code - Only creates plans, never implements
- R383: Metadata File Standards - Fix plans follow proper format
- R355: Code Reviewer Validation - All outputs validated and complete
- R510: State Execution Checklist Compliance (this file)
- R511: Checklist Creation Protocol (checklist design)

---

## Transition Rules

- **ALWAYS** → COMPLETED (after successful fix plan creation)
- **NEVER** skip to orchestrator states (Code Reviewer completes and stops)
- **ERROR** → ERROR_RECOVERY (if critical failure prevents fix planning)

**Note**: Code Reviewer does NOT transition to fix implementation. That's orchestrator's job to spawn SW Engineers.

---

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-code-reviewing to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ Fix plan creation complete"

# 2. Update state file
jq '.state_machine.current_state = "COMPLETED"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ Fix plan created successfully
- ✅ All bugs documented in source branches
- ✅ Bug tracking updated in state file
- ✅ Ready for orchestrator to spawn SW Engineers
- ✅ Waiting for user to continue (NORMAL)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Integration review data corrupted/missing
- ❌ Cannot access source effort branches
- ❌ State file corruption prevents bug tracking
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs

**See**: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`

---

## Additional Context

### R321 Enforcement Critical Points

**Integration branches are READ-ONLY.** This state exists to document bugs in SOURCE effort branches, NOT to fix them in integration branches. Every bug identified during integration MUST be:

1. Documented in the source effort branch with a BUG-{NNN}-{DESC}.md file
2. Tracked in bug-tracking.json bugs array
3. Fixed by SW Engineers in the source branch
4. Re-integrated after fixes complete

**Violation of R321 (fixing in integration instead of source) = -100% FAILURE**

### R266 Two-Step Protocol

1. **Code Reviewer (this state)**: DOCUMENT bugs only
2. **Orchestrator (next)**: COORDINATE fixes through SW Engineer spawns

**Code Reviewer NEVER fixes bugs. Code Reviewer DOCUMENTS bugs for others to fix.**

### Bug File Naming Convention

```
BUG-{NNN}-{SHORT-DESCRIPTION}.md
```

Examples:
- `BUG-001-DUPLICATE-PUSHCMD.md`
- `BUG-002-UNDEFINED-VARIABLE.md`
- `BUG-003-TEST-TIMEOUT.md`

### bug-tracking.json Bug Schema

```json
{
  "bugs_discovered": [
    {
      "bug_id": "BUG-{NNN}-{DESC}",
      "severity": "CRITICAL|HIGH|MEDIUM|LOW",
      "discovered_in": "phase1-integration|wave1-integration|effort-review",
      "source_efforts": ["E1.2.1-effort-name"],
      "fix_plan": ".software-factory/phase1/FIX-PLAN-*.md",
      "backport_required": true|false,
      "status": "DOCUMENTED|FIXING|FIXED|REINTEGRATED"
    }
  ]
}
```

---

**Template Version**: 2.0
**Created**: 2025-10-07
**Purpose**: Code Reviewer fix planning with R321/R266 enforcement
**Compliance**: R516 State Creation Protocol, R510 Checklist Compliance, R511 Checklist Design
