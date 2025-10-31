# 🔴🔴🔴 RULE R511 - CHECKLIST CREATION PROTOCOL (SUPREME LAW)

**Criticality:** SUPREME LAW - PROTOCOL FOR STATE DESIGN
**Grading Impact:** -100% for states with improper checklists
**Enforcement:** CONTINUOUS - Every new state must follow this protocol

---

## SUPREME LAW STATEMENT

**WHEN CREATING OR MODIFYING ANY STATE, THE CHECKLIST MUST BE DESIGNED FOLLOWING THIS PROTOCOL. CHECKLIST ITEMS MUST BE OUTCOME-FOCUSED, PROPERLY CATEGORIZED, AND VERIFIABLE. IMPROPER GRANULARITY OR CATEGORIZATION IS AUTOMATIC FAILURE.**

---

## 🚨🚨🚨 THE CREATION MANDATE 🚨🚨🚨

### Checklists Are Not Optional

When designing a state's rules.md file:
```markdown
✅ REQUIRED: Follow R511 checklist creation protocol
✅ REQUIRED: Use R510 checklist structure (BLOCKING/STANDARD/EXIT)
✅ REQUIRED: Apply proper granularity guidelines
✅ REQUIRED: Categorize items correctly
✅ REQUIRED: Define validation criteria
❌ FORBIDDEN: Skip checklist creation
❌ FORBIDDEN: Use vague or command-level items
❌ FORBIDDEN: Miscategorize BLOCKING vs STANDARD
```

---

## 🔴🔴🔴 PART 1: GRANULARITY DECISION FRAMEWORK 🔴🔴🔴

### The Three Granularity Levels

#### ❌ TOO GRANULAR (Command-Level Micromanagement)

**Anti-Pattern: Individual bash commands**

```markdown
❌ BAD EXAMPLE:
- [ ] 1. Run: cd /efforts/phase1
- [ ] 2. Run: /spawn agent-code-reviewer
- [ ] 3. Run: echo "spawned" >> log.txt
- [ ] 4. Run: jq '.reviewer = "id"' state.json
```

**Why This Fails:**
- Focuses on HOW (commands) not WHAT (outcome)
- No accountability for actual result
- Brittle to implementation changes
- Impossible to measure success
- **Grading:** -50% for command-level checklists

#### ❌ TOO VAGUE (High-Level Hand-Waving)

**Anti-Pattern: Meaningless abstractions**

```markdown
❌ BAD EXAMPLE:
- [ ] 1. Complete the state
- [ ] 2. Do everything required
- [ ] 3. Finish the work
- [ ] 4. Make it good
```

**Why This Fails:**
- No specific actionable outcome
- No validation criteria
- No accountability
- Cannot determine completion
- **Grading:** -50% for vague checklists

#### ✅ JUST RIGHT (Outcome-Focused with Validation)

**Correct Pattern: Verifiable outcomes**

```markdown
✅ GOOD EXAMPLE:
- [ ] 1. Spawn Code Reviewer for phase integration quality assessment
  - Agent: code-reviewer
  - Focus: phase-integration-quality
  - Workspace: ${phase_integration_workspace}
  - Parameters: --phase P1 --branch phase1-integration
  - Validation: Agent ID returned with timestamp
  - **BLOCKING**: Cannot transition without reviewer spawned

- [ ] 2. Update state file with phase integration review metadata
  - Fields: reviewer_id, review_branch, review_focus, started_at
  - Source: Code Reviewer spawn response
  - Per: R288 state update requirements
  - Validation: jq '.phase_integration_review' shows all fields
```

**Why This Works:**
- Specifies WHAT outcome must be achieved
- Provides enough context to execute
- Defines clear validation criteria
- Outcome-focused, not command-focused
- **Grading:** Full credit

### The Granularity Decision Tree

```
Is this a checklist item?
│
├─ Does it describe a single command?
│  └─ YES → ❌ TOO GRANULAR - combine into outcome
│
├─ Does it have measurable completion criteria?
│  └─ NO → ❌ TOO VAGUE - add specifics
│
├─ Can I verify it succeeded without ambiguity?
│  └─ NO → ❌ INSUFFICIENT - add validation
│
├─ Is it outcome-focused (WHAT) not process-focused (HOW)?
│  └─ NO → ❌ WRONG FOCUS - reframe as outcome
│
└─ YES to verification, outcome-focused, not command-level
   └─ ✅ CORRECT GRANULARITY
```

### Granularity Examples by State Type

#### SPAWN States (SPAWN_CODE_REVIEWER, SPAWN_ARCHITECT, etc.)

```markdown
✅ BLOCKING REQUIREMENT:
- [ ] 1. Spawn [Agent Type] for [specific purpose]
  - Agent: [agent-name]
  - Purpose: [why spawning]
  - Workspace: [where agent works]
  - Parameters: [critical spawn params]
  - Validation: Agent ID with timestamp returned
  - **BLOCKING**: Cannot transition without successful spawn

❌ TOO GRANULAR:
- [ ] 1. Run /spawn command
- [ ] 2. Check exit code
- [ ] 3. Save agent ID

❌ TOO VAGUE:
- [ ] 1. Spawn the reviewer
```

#### MONITORING_SWE_PROGRESS/WAITING States (WAITING_FOR_REVIEW, MONITOR_*, etc.)

```markdown
✅ STANDARD TASK:
- [ ] 1. Display current review status and next steps
  - Show: Reviewer ID, branch under review, expected completion
  - Show: How to check review progress
  - Show: What /continue command to use
  - Expected: Clear user guidance message

❌ TOO GRANULAR:
- [ ] 1. Echo "Waiting for review"
- [ ] 2. Print reviewer ID
- [ ] 3. Print branch name

❌ TOO VAGUE:
- [ ] 1. Show status
```

#### ACTION States (CREATE_*, ANALYZE_*, DISTRIBUTE_*, etc.)

```markdown
✅ BLOCKING REQUIREMENT:
- [ ] 1. Create effort infrastructure for next effort in queue
  - Source: pre_planned_infrastructure.efforts (next with created=false)
  - Action: Clone base branch, create effort branch, push to remote
  - Metadata: Update created=true, created_at timestamp
  - Validation: Branch exists on remote, workspace accessible
  - **BLOCKING**: Cannot spawn SWE without infrastructure

❌ TOO GRANULAR:
- [ ] 1. Get effort ID from queue
- [ ] 2. Run git clone
- [ ] 3. Run git checkout -b
- [ ] 4. Run git push

❌ TOO VAGUE:
- [ ] 1. Create infrastructure
```

---

## 🔴🔴🔴 PART 2: BLOCKING vs STANDARD vs EXIT CATEGORIZATION 🔴🔴🔴

### BLOCKING REQUIREMENTS (Cannot proceed without)

**Definition:** Actions that MUST complete before the state can transition. If these fail, the state CANNOT move forward.

#### BLOCKING Criteria Decision Matrix

Use BLOCKING when:
```yaml
primary_state_purpose: true  # This IS what the state exists to do
blocks_next_state: true      # Next state CANNOT start without this
creates_dependency: true     # Other agents/states depend on this
failure_is_critical: true    # Failure means state must ERROR_RECOVERY
```

#### BLOCKING Examples by Pattern

**Pattern 1: Spawning Agents**
```markdown
✅ ALWAYS BLOCKING:
- Spawn Code Reviewer (SPAWN_CODE_REVIEWER state)
- Spawn Software Engineer (SPAWN_SW_ENGINEER state)
- Spawn Architect (SPAWN_ARCHITECT state)

REASON: State exists SOLELY to spawn agent
FAILURE: Cannot proceed to WAITING_FOR_* without agent
```

**Pattern 2: Creating Infrastructure**
```markdown
✅ ALWAYS BLOCKING:
- Create effort workspace and branch (CREATE_NEXT_INFRASTRUCTURE)
- Create phase integration branch (INTEGRATE_PHASE_WAVES)
- Create split infrastructure (CREATE_SPLIT_INFRASTRUCTURE)

REASON: Next state requires this infrastructure to exist
FAILURE: Software Engineer cannot work without workspace
```

**Pattern 3: Critical Validations**
```markdown
✅ ALWAYS BLOCKING:
- Validate all waves integrated (PHASE_REVIEW_WAVE_INTEGRATION state)
- Verify effort branch merged (WAVE_INTEGRATE_WAVE_EFFORTS state)
- Confirm tests passing (VALIDATION state)

REASON: Cannot proceed with integration if validation fails
FAILURE: Integration would corrupt the cascade
```

**Pattern 4: Distributing Work**
```markdown
✅ ALWAYS BLOCKING:
- Distribute effort assignments to metadata files (DISTRIBUTE_EFFORT_ASSIGNMENTS)
- Load next effort from queue (planning states)
- Update final_merge_plan (planning states)

REASON: Agents depend on this data to work
FAILURE: Agents cannot find their assignments
```

#### NOT BLOCKING Examples

**Pattern 1: Informational Display**
```markdown
❌ NEVER BLOCKING:
- Display current status
- Show next steps
- Print summary information

REASON: Displaying info doesn't block transitions
CATEGORY: STANDARD TASK
```

**Pattern 2: Optional Enhancements**
```markdown
❌ NEVER BLOCKING:
- Log activity to audit trail
- Send notification
- Update dashboard

REASON: Nice to have, but state can succeed without them
CATEGORY: STANDARD TASK
```

**Pattern 3: State Metadata Updates**
```markdown
❌ NOT BLOCKING (Usually):
- Update state file with non-critical metadata
- Record timestamps
- Track metrics

REASON: These document what happened, but aren't prerequisites
CATEGORY: STANDARD TASK
EXCEPTION: R288 state transition update IS in EXIT REQUIREMENTS
```

### STANDARD EXECUTION TASKS (Required but not blocking)

**Definition:** Actions that should complete during the state, but whose failure doesn't prevent transition.

#### STANDARD Criteria Decision Matrix

Use STANDARD when:
```yaml
supports_primary_purpose: true   # Helps the main action succeed
failure_tolerable: true           # State can still succeed if this fails
informational: true               # Provides context/guidance
documentation: true               # Records what happened
```

#### STANDARD Examples

```markdown
✅ STANDARD TASKS:

- [ ] 2. Update state file with reviewer metadata
  - Fields: reviewer_id, review_started_at, review_branch
  - Purpose: Track review details for monitoring
  - Note: Not blocking - EXIT REQUIREMENTS handle critical state updates

- [ ] 3. Display progress summary to user
  - Show: Current phase, wave, effort count
  - Show: Efforts completed vs remaining
  - Purpose: User awareness

- [ ] 4. Log state transition to audit trail
  - File: logs/orchestrator-transitions.log
  - Content: Timestamp, from_state, to_state, reason
  - Purpose: Debugging and compliance
```

### EXIT REQUIREMENTS (Standardized across ALL states)

**Definition:** The EXACT SAME exit requirements for EVERY state. These are never customized.

#### The Canonical EXIT REQUIREMENTS

**COPY THIS EXACTLY into every state:**

```markdown
### EXIT REQUIREMENTS (Must complete before transition)

**NOTE**: These are STANDARD across ALL states - copy exactly

- [ ] 5. Update state file to [NEXT_STATE] per R288
  - Field: `current_state`
  - Value: `"[NEXT_STATE]"`
  - Also update: `previous_state`, `transition_time`, `transition_reason`
  - Validation: `jq '.state_machine.current_state' orchestrator-state-v3.json` shows new state

- [ ] 6. Save TODOs per R287 (within 30s of last TodoWrite)
  - Trigger: "R510_CHECKLIST_COMPLETE"
  - Format: `todos/[agent]-[STATE]-[YYYYMMDD-HHMMSS].todo`
  - Validation: TODO file exists and contains current state

- [ ] 7. Commit all changes with descriptive message
  - Include: State transition details
  - Include: Checklist completion confirmation
  - Include: Rule compliance references (R288, R287, R510)
  - Format: Multi-line commit message with context

- [ ] 8. Push changes to remote
  - Remote: `origin`
  - Branch: Current branch
  - Validation: `git status` shows "up to date with origin"

- [ ] 9. Set CONTINUE-SOFTWARE-FACTORY flag per R405
  - Value: `TRUE` (if successful completion, factory can continue)
  - Value: `FALSE` (if catastrophic error, must halt everything)
  - Context: Explain why TRUE or FALSE
  - **NOTE**: R322 checkpoints = TRUE (agent stops but factory continues)

- [ ] 10. Display checkpoint message (if this is R322 checkpoint state)
  - Format: Clear message about what was done
  - Include: Next steps for /continue-[agent]
  - Refer to: R322 requirements

- [ ] 11. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Per: R313 (if spawn state), R322 (if checkpoint)
```

**CRITICAL RULES:**

1. **NEVER modify EXIT REQUIREMENTS** - they are standardized
2. **ALWAYS include ALL 7 items** (5-11) - no omissions
3. **ONLY change [NEXT_STATE]** placeholder - nothing else
4. **Item numbers continue** from BLOCKING and STANDARD sections

#### Why EXIT REQUIREMENTS Are Standardized

```yaml
consistency_across_states:
  - Every state follows same cleanup protocol
  - No state can skip critical requirements
  - Predictable behavior for all transitions

enforces_critical_rules:
  - R288: State file updates
  - R287: TODO persistence
  - R405: Automation continuation
  - R322: Checkpoint messaging
  - R313: Spawn state stopping

prevents_common_failures:
  - Uncommitted state changes
  - Lost TODOs
  - Missing continuation flags
  - Corrupted state files
```

---

## 🔴🔴🔴 PART 3: WRITING EFFECTIVE CHECKLIST ITEMS 🔴🔴🔴

### Item Format Template

```markdown
- [ ] [n]. [Outcome-focused action description]
  - [Context field 1]: [value or guidance]
  - [Context field 2]: [value or guidance]
  - [Validation]: [How to verify completion]
  - **[Criticality]**: [Why this matters] (if BLOCKING)
```

### Format Standards

#### Required Components

Every checklist item MUST have:

1. **Number** - Sequential within section
2. **Action Description** - Outcome-focused verb phrase
3. **Context** - Enough detail to execute
4. **Validation** - How to verify completion

#### Optional Components

BLOCKING items SHOULD also have:

5. **Criticality Marker** - `**BLOCKING**: [reason]`
6. **Failure Scenario** - What happens if this fails

#### Format Examples

**Minimal STANDARD Item:**
```markdown
- [ ] 3. Display current orchestrator status
  - Show: Current phase, wave, effort
  - Expected: User sees clear status message
```

**Complete BLOCKING Item:**
```markdown
- [ ] 1. Spawn Software Engineer for effort implementation
  - Agent: sw-engineer
  - Workspace: /efforts/idpbuilder-oci-mgmt/phase1/wave1/effort-001
  - Branch: idpbuilder-oci-mgmt/phase1/wave1/effort-001
  - Assignment: ${effort_metadata}/EFFORT-ASSIGNMENT.md
  - Validation: Agent ID returned, workspace verified accessible
  - **BLOCKING**: Cannot start implementation without SWE agent
  - Failure: Transition to ERROR_RECOVERY if spawn fails
```

### Proof/Validation Requirements

#### Validation Types by Action

| Action Type | Required Validation |
|-------------|-------------------|
| Spawn agent | Agent ID with timestamp |
| Update state file | jq query showing field value |
| Create branch | git branch exists on remote |
| Commit changes | git log shows commit hash |
| Push changes | git status shows "up to date" |
| Save file | ls shows file exists with timestamp |
| Run tests | Exit code 0 or test output |
| Display message | Message content shown |

#### Validation Examples

```markdown
✅ GOOD VALIDATIONS:

Spawn:
  - Validation: Agent ID returned with timestamp
  - Proof format: agent-code-reviewer-phase1-20251007-020300

State Update:
  - Validation: jq '.state_machine.current_state' shows "NEXT_STATE"
  - Proof format: State file query result

Branch Creation:
  - Validation: git ls-remote shows branch
  - Proof format: Branch name in remote list

Commit:
  - Validation: git log -1 shows commit message
  - Proof format: Commit hash (short)
```

### Acknowledgment Format

When agent executes checklist, MUST output:

```bash
✅ CHECKLIST[n]: [Description] [proof]
```

**Examples:**
```bash
✅ CHECKLIST[1]: Spawned Code Reviewer [agent-code-reviewer-20251007-020300]
✅ CHECKLIST[2]: Updated state field 'phase_integration_review'
✅ CHECKLIST[3]: Updated state to WAITING_FOR_REVIEW per R288
✅ CHECKLIST[4]: Saved TODOs to todos/orchestrator-REVIEW-20251007.todo per R287
✅ CHECKLIST[5]: Committed changes [commit: a9916de]
✅ CHECKLIST[6]: Pushed to origin/integration-agent-fixes
✅ CHECKLIST[7]: CONTINUE-SOFTWARE-FACTORY=TRUE (successful spawn)
✅ CHECKLIST[8]: Stopped with exit 0 per R322
```

---

## 🔴🔴🔴 PART 4: STATE-SPECIFIC GUIDANCE 🔴🔴🔴

### Spawn States (SPAWN_*)

**Pattern:** Spawn agent, update metadata, stop

```markdown
BLOCKING REQUIREMENTS:
- [ ] 1. Spawn [specific agent] for [specific purpose]
  - Agent: [agent-name]
  - Purpose: [why spawning]
  - Workspace: [effort or review workspace]
  - Parameters: [critical params]
  - Validation: Agent ID returned
  - **BLOCKING**: Cannot transition without agent

STANDARD TASKS:
- [ ] 2. Update state file with [agent type] metadata
  - Fields: [agent]_id, [agent]_workspace, started_at
  - Purpose: Track spawned agent details

EXIT REQUIREMENTS:
- [ ] 3-9. [Standard EXIT REQUIREMENTS]
```

**Customization:** Only change agent type, purpose, workspace

### Waiting States (WAITING_FOR_*, MONITOR_*)

**Pattern:** Display status, guide user, checkpoint

```markdown
BLOCKING REQUIREMENTS:
- Usually NONE for pure waiting states
- If checking condition, that becomes BLOCKING

STANDARD TASKS:
- [ ] 1. Display current [process] status
  - Show: What's being waited for
  - Show: How to check progress
  - Show: What command resumes (/continue-[agent])
  - Expected: Clear user guidance

- [ ] 2. Verify [condition being waited for] (if applicable)
  - Check: [what to verify]
  - Validation: [how to confirm]

EXIT REQUIREMENTS:
- [ ] 3-9. [Standard EXIT REQUIREMENTS including R322 checkpoint message]
```

**Customization:** Specify what's being waited for

### Action States (CREATE_*, ANALYZE_*, DISTRIBUTE_*, INTEGRATE_*)

**Pattern:** Perform action, validate, update state

```markdown
BLOCKING REQUIREMENTS:
- [ ] 1. [Primary action the state exists to perform]
  - Input: [what data/state needed]
  - Action: [what to do]
  - Output: [what's created/updated]
  - Validation: [how to verify success]
  - **BLOCKING**: [Why this blocks transition]

- [ ] 2. Validate [action result]
  - Check: [critical validation]
  - Criteria: [success criteria]
  - **BLOCKING**: Cannot proceed with invalid result

STANDARD TASKS:
- [ ] 3. Update state file with [action] results
  - Fields: [what to update]
  - Purpose: [why tracking this]

EXIT REQUIREMENTS:
- [ ] 4-10. [Standard EXIT REQUIREMENTS]
```

**Customization:** Specify action, validation, results

### Planning States (PLANNING_*, ANALYSIS_*)

**Pattern:** Analyze, make decisions, record plan

```markdown
BLOCKING REQUIREMENTS:
- [ ] 1. Analyze [what to analyze]
  - Source: [input data]
  - Analysis: [what to determine]
  - Decision: [what to decide]
  - Output: [plan/document created]
  - Validation: [plan completeness check]
  - **BLOCKING**: Next state requires this plan

- [ ] 2. Validate plan against [criteria]
  - Check: [validation criteria]
  - Requirement: [what must be true]
  - **BLOCKING**: Invalid plan blocks execution

STANDARD TASKS:
- [ ] 3. Record plan in state file
  - Field: [where to store]
  - Content: [what to record]

EXIT REQUIREMENTS:
- [ ] 4-10. [Standard EXIT REQUIREMENTS]
```

**Customization:** Specify analysis type, plan format

---

## 🔴🔴🔴 PART 5: THE CHECKLIST CREATION WORKFLOW 🔴🔴🔴

### Step-by-Step Checklist Creation Process

When creating a new state or updating an existing one:

#### STEP 1: Understand State Purpose

```yaml
questions_to_answer:
  - What is this state's PRIMARY purpose?
  - What MUST happen before transitioning?
  - What's nice-to-have but not critical?
  - What's the success criteria?
  - What's the failure scenario?
```

#### STEP 2: Identify BLOCKING Items

```yaml
blocking_item_criteria:
  - Is this THE reason the state exists?
  - Does next state depend on this completion?
  - Would failure require ERROR_RECOVERY?
  - Is this action irreplaceable?

decision:
  - YES to any → BLOCKING
  - NO to all → STANDARD or skip
```

#### STEP 3: Draft Items at Correct Granularity

```yaml
for_each_item:
  1. Write outcome-focused description
  2. Add context (parameters, source, target)
  3. Define validation criteria
  4. Add BLOCKING marker if applicable
  5. Specify proof format
```

#### STEP 4: Add STANDARD Tasks

```yaml
standard_task_types:
  - State metadata updates (non-critical)
  - Status displays
  - Logging/audit trail
  - Informational messages
  - Optional validations
```

#### STEP 5: Copy EXIT REQUIREMENTS

```yaml
action:
  - Copy canonical EXIT REQUIREMENTS exactly
  - Replace [NEXT_STATE] placeholder
  - Continue item numbering from STANDARD
  - DO NOT modify any other text
```

#### STEP 6: Review Checklist

```yaml
review_checklist:
  granularity:
    - No command-level items?
    - No vague items?
    - Outcome-focused?

  categorization:
    - BLOCKING items truly block transition?
    - STANDARD items truly optional?
    - EXIT REQUIREMENTS unmodified?

  completeness:
    - All validation criteria defined?
    - All context provided?
    - Proof formats specified?

  consistency:
    - Follows state pattern?
    - Matches similar states?
    - Aligns with state machine?
```

---

## 🔴🔴🔴 PART 6: COMMON MISTAKES AND FIXES 🔴🔴🔴

### Mistake 1: Command-Level Granularity

```markdown
❌ WRONG:
- [ ] 1. Run: jq '.state = "NEXT"' state.json
- [ ] 2. Run: git add state.json
- [ ] 3. Run: git commit -m "update"

✅ RIGHT:
- [ ] 5. Update state file to NEXT_STATE per R288
  - Field: current_state
  - Value: "NEXT_STATE"
  - Also update: previous_state, transition_time
  - Validation: jq '.state_machine.current_state' shows new state
```

**Fix:** Combine commands into outcome, add validation

### Mistake 2: Vague Descriptions

```markdown
❌ WRONG:
- [ ] 1. Do the review
- [ ] 2. Make it work
- [ ] 3. Check things

✅ RIGHT:
- [ ] 1. Spawn Code Reviewer for effort quality assessment
  - Agent: code-reviewer
  - Focus: effort-quality-review
  - Workspace: ${effort_workspace}
  - Validation: Agent ID returned with timestamp
```

**Fix:** Add specifics, context, validation criteria

### Mistake 3: Wrong BLOCKING Categorization

```markdown
❌ WRONG (marked BLOCKING but isn't):
- [ ] 1. Log state transition to audit file
  - **BLOCKING**: Must track all transitions

✅ RIGHT (STANDARD):
- [ ] 3. Log state transition to audit file
  - File: logs/transitions.log
  - Purpose: Debugging and compliance
  (No BLOCKING marker - nice to have, not critical)
```

**Fix:** Apply BLOCKING criteria decision matrix

### Mistake 4: Modified EXIT REQUIREMENTS

```markdown
❌ WRONG (customized EXIT):
- [ ] 5. Update state file if needed
- [ ] 6. Maybe save TODOs
- [ ] 7. Commit changes (optional)

✅ RIGHT (canonical EXIT):
- [ ] 5. Update state file to NEXT_STATE per R288
  [... exact canonical text ...]
- [ ] 6. Save TODOs per R287 (within 30s of last TodoWrite)
  [... exact canonical text ...]
- [ ] 7. Commit all changes with descriptive message
  [... exact canonical text ...]
```

**Fix:** Copy canonical EXIT REQUIREMENTS exactly

### Mistake 5: Missing Validation

```markdown
❌ WRONG (no way to verify):
- [ ] 1. Spawn the code reviewer

✅ RIGHT (clear validation):
- [ ] 1. Spawn Code Reviewer for review
  - Agent: code-reviewer
  - Validation: Agent ID returned with timestamp
  - Proof format: agent-code-reviewer-20251007-HHMMSS
```

**Fix:** Always specify validation criteria and proof format

---

## 🔴🔴🔴 PART 7: TEMPLATE USAGE 🔴🔴🔴

### Use the Official Template

**File:** `/home/vscode/software-factory-template/templates/state-rules-checklist-template.md`

**When creating new state:**

1. Copy template to new state directory
2. Rename placeholders: [STATE_NAME], [agent-type], [NEXT_STATE]
3. Follow R511 protocol to design checklist
4. Add state-specific rules to PRIMARY DIRECTIVES
5. Customize BLOCKING and STANDARD sections
6. Keep EXIT REQUIREMENTS exactly as-is
7. Fill in State Purpose, Entry/Exit Criteria

**Template ensures:**
- Consistent structure across all states
- R510 compliance built-in
- Correct EXIT REQUIREMENTS format
- All required sections present

---

## 🔴🔴🔴 PART 8: ENFORCEMENT AND VALIDATION 🔴🔴🔴

### Pre-Commit Validation

States MUST validate checklist structure:

```bash
# File: .pre-commit/validate-state-checklist.sh

validate_state_checklist() {
    local state_file="$1"

    # Check checklist section exists
    if ! grep -q "## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST" "$state_file"; then
        echo "❌ R511 VIOLATION: Missing checklist section"
        return 1
    fi

    # Check BLOCKING section
    if ! grep -q "### BLOCKING REQUIREMENTS" "$state_file"; then
        echo "❌ R511 VIOLATION: Missing BLOCKING REQUIREMENTS section"
        return 1
    fi

    # Check EXIT REQUIREMENTS match canonical
    if ! diff <(extract_exit_requirements "$state_file") \
              <(get_canonical_exit_requirements); then
        echo "❌ R511 VIOLATION: EXIT REQUIREMENTS modified from canonical"
        return 1
    fi

    # Check item granularity (heuristic)
    if grep -q "^- \[ \] [0-9]*\. Run:" "$state_file"; then
        echo "❌ R511 VIOLATION: Command-level granularity detected"
        return 1
    fi

    echo "✅ R511 COMPLIANCE: Checklist structure valid"
    return 0
}
```

### Grading Criteria

```yaml
checklist_quality_grading:
  structure:
    has_all_sections: 20%           # BLOCKING, STANDARD, EXIT
    exit_requirements_canonical: 20% # Unmodified from template

  granularity:
    outcome_focused: 20%             # Not command-level
    proper_detail: 15%               # Not too vague

  categorization:
    blocking_correct: 15%            # True blockers only
    standard_appropriate: 10%        # Non-critical items

total: 100%

violations:
  command_level_items: -50%          # Per item
  vague_items: -30%                  # Per item
  wrong_blocking: -40%               # Per misclassified item
  modified_exit: -100%               # Automatic failure
  missing_validation: -20%           # Per item
```

---

## 🔴 THE CHECKLIST CREATION OATH

```
I, the State Designer, swear by R511:

When I create a state checklist:
- I WILL use outcome-focused granularity
- I WILL categorize items correctly (BLOCKING vs STANDARD)
- I WILL define clear validation criteria
- I WILL use the canonical EXIT REQUIREMENTS unmodified
- I WILL NOT use command-level items
- I WILL NOT create vague items
- I WILL NOT miscategorize criticality

My checklists make execution explicit and verifiable.
R511 + R510 together ensure systematic state completion.

This is SUPREME LAW.
Improper checklists mean FAILURE.
I WILL DESIGN CHECKLISTS PROPERLY.
```

---

**Remember:** A well-designed checklist makes state execution deterministic. Agents know exactly what to do, how to validate, and when they're done. R511 ensures every state's checklist is properly designed from the start.

**See Also:**
- R510: State Execution Checklist Compliance (structure and execution)
- R288: State File Update Requirements (EXIT REQUIREMENTS component)
- R287: TODO Persistence (EXIT REQUIREMENTS component)
- R322: Mandatory Checkpoints (EXIT REQUIREMENTS component)
- R405: Automation Continuation Flag (EXIT REQUIREMENTS component)
- Template: `/home/vscode/software-factory-template/templates/state-rules-checklist-template.md`
