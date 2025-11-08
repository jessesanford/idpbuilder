# 🔴🔴🔴 RULE R516 - STATE CREATION AND DESIGN PROTOCOL (SUPREME LAW)

**Criticality:** SUPREME LAW - PROTOCOL FOR STATE DESIGN
**Grading Impact:** -100% for states that don't follow this protocol
**Enforcement:** CONTINUOUS - Every new state must follow this protocol

---

## SUPREME LAW STATEMENT

**WHEN CREATING ANY NEW STATE OR MODIFYING EXISTING STATES, THIS PROTOCOL MUST BE FOLLOWED. STATES THAT DON'T FOLLOW THE STANDARD STRUCTURE, MISS REQUIRED SECTIONS, OR FAIL VALIDATION ARE AUTOMATIC FAILURE.**

---

## 🚨🚨🚨 THE STATE CREATION MANDATE 🚨🚨🚨

### States Are Not Arbitrary

Creating or modifying a state is a structured process:
```markdown
✅ REQUIRED: Follow R516 state creation protocol
✅ REQUIRED: Use official state template
✅ REQUIRED: Include ALL required sections
✅ REQUIRED: Design checklist per R511
✅ REQUIRED: Integrate with state machine
✅ REQUIRED: Validate before deployment
❌ FORBIDDEN: Ad-hoc state creation
❌ FORBIDDEN: Missing required sections
❌ FORBIDDEN: States not in state machine
```

---

## 🔴🔴🔴 PART 1: WHEN TO CREATE A NEW STATE 🔴🔴🔴

### Decision Matrix: New State vs Modify Existing

#### CREATE NEW STATE when:

```yaml
distinct_purpose:
  - State represents fundamentally different phase of work
  - Different agent should handle this work
  - Entry/exit criteria are significantly different
  - Requires different validation/error handling

workflow_necessity:
  - State machine needs explicit breakpoint here
  - User review/checkpoint required at this point
  - Different transition paths needed
  - Parallel work streams diverge/converge here

examples:
  - SPAWN_CODE_REVIEWER vs SPAWN_ARCHITECT (different agents)
  - CODE_REVIEW vs INTEGRATE_WAVE_EFFORTS_REVIEW (different validation)
  - CREATE_EFFORT_INFRASTRUCTURE vs CREATE_PHASE_INFRASTRUCTURE (different scope)
```

#### MODIFY EXISTING STATE when:

```yaml
same_purpose_enhanced:
  - Adding more details to existing state
  - Improving validation within same state
  - Adding checklist items for same work
  - Clarifying existing requirements

cosmetic_changes:
  - Updating documentation
  - Fixing typos or formatting
  - Adding examples
  - Clarifying instructions

examples:
  - Adding R510 checklist to existing state
  - Updating spawn parameters in SPAWN_* state
  - Adding validation criteria to existing action
```

### The 5-Question Test

Before creating new state, answer:

1. **Does an existing state already do this?** → If YES, modify existing
2. **Is this a distinct phase with different actor?** → If NO, don't create
3. **Does state machine support this transition?** → If NO, update state machine first
4. **Can this be a checklist item instead?** → If YES, add to existing state
5. **Is there a natural checkpoint/breakpoint here?** → If NO, reconsider

---

## 🔴🔴🔴 PART 2: STATE NAMING CONVENTIONS 🔴🔴🔴

### Naming Pattern Standards

#### Pattern Categories

**SPAWN States:**
```
Format: SPAWN_[AGENT]_[PURPOSE]
Examples:
  - SPAWN_CODE_REVIEWER
  - SPAWN_SW_ENGINEER
  - SPAWN_ARCHITECT_WAVE_PLANNING
```

**WAITING States:**
```
Format: WAITING_FOR_[WHAT]_[CONTEXT]
Examples:
  - WAITING_FOR_CODE_REVIEW
  - WAITING_FOR_INTEGRATE_WAVE_EFFORTS_REVIEW
  - WAITING_FOR_ARCHITECTURE_APPROVAL
```

**MONITORING_SWE_PROGRESS States:**
```
Format: MONITORING_[WHAT]_[CONTEXT]
Examples:
  - MONITORING_SWE_PROGRESS
  - MONITORING_INTEGRATE_WAVE_EFFORTS_ATTEMPT
```

**ACTION States:**
```
Format: [VERB]_[NOUN]_[CONTEXT]
Examples:
  - CREATE_EFFORT_INFRASTRUCTURE
  - DISTRIBUTE_EFFORT_ASSIGNMENTS
  - ANALYZE_IMPLEMENTATION_PLAN
  - VALIDATE_WAVE_COMPLETION
```

**PLANNING States:**
```
Format: [PLAN/ANALYZE/DESIGN]_[WHAT]
Examples:
  - PLANNING
  - ANALYZE_EFFORT_SPLITS
  - DESIGN_PHASE_STRUCTURE
```

**INTEGRATE_WAVE_EFFORTS States:**
```
Format: [SCOPE]_INTEGRATE_WAVE_EFFORTS_[STAGE]
Examples:
  - WAVE_INTEGRATE_WAVE_EFFORTS
  - INTEGRATE_PHASE_WAVES
  - EFFORT_INTEGRATE_WAVE_EFFORTS_VALIDATION
```

**ITERATION_CONTAINER States:**
```
Format: [VERB]_[LEVEL]_[NOUN]
Examples:
  - SETUP_WAVE_INFRASTRUCTURE
  - INTEGRATE_WAVE_EFFORTS
  - REVIEW_PHASE_ARCHITECTURE
  - COMPLETE_PROJECT
  - CREATE_WAVE_FIX_PLAN
Description:
  - Used in SF 3.0 for iteration containers
  - [VERB]: Action being performed (SETUP, INTEGRATE, REVIEW, COMPLETE, etc.)
  - [LEVEL]: WAVE, PHASE, or PROJECT
  - [NOUN]: Object of the action (INFRASTRUCTURE, EFFORTS, ARCHITECTURE, etc.)
  - See SF 3.0 Architecture Part 8 and state machine JSON for complete list
```

### Naming Requirements

```yaml
must_haves:
  - ALL_CAPS format
  - Underscore_separated
  - Action verb (except WAITING/MONITORING_SWE_PROGRESS)
  - Clear scope/context
  - No abbreviations (unless standard like SWE)

must_nots:
  - Lowercase or mixed case
  - Hyphens or spaces
  - Vague verbs (DO, HANDLE, PROCESS)
  - Generic names (STATE1, TEMP_STATE)
```

---

## 🔴🔴🔴 PART 3: REQUIRED SECTIONS IN RULES.MD 🔴🔴🔴

### Mandatory File Structure

Every state's `rules.md` MUST have these sections IN THIS ORDER:

#### 1. Header Section

```markdown
# [STATE_NAME] State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES.
**I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS**
*YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE!!!

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE
THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**
```

#### 2. PRIMARY DIRECTIVES Section

```markdown
## 📋 PRIMARY DIRECTIVES FOR [STATE_NAME] STATE

### Core Mandatory Rules (ALL [agent-type] states must have these):

1. **🚨🚨🚨 R[XXX]** - [RULE NAME] (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R[XXX]-[rule-file].md`
   - Criticality: BLOCKING
   - Summary: [Brief summary]

2. **🔴🔴🔴 R510** - STATE EXECUTION CHECKLIST COMPLIANCE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R510-state-execution-checklist-compliance.md`
   - Criticality: SUPREME LAW
   - Summary: MUST complete and acknowledge every checklist item

[... other universal rules ...]

### State-Specific Rules:

[... rules specific to THIS state ...]
```

**Required Universal Rules** (ALL states):
- R510: State Execution Checklist Compliance
- R288: State File Update and Commit Protocol
- R287: TODO Persistence
- R322: Mandatory Checkpoints (if checkpoint state)
- R405: Automation Continuation Flag

**Agent-Specific Rules** (vary by agent):
- Orchestrator states: R232, R313, R234
- SW-Engineer states: R220, R221
- Code-Reviewer states: R355, R322
- etc.

#### 3. MANDATORY EXECUTION CHECKLIST Section

```markdown
## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. [Primary critical action]
  - Context/parameters
  - Validation criteria
  - **BLOCKING**: Reason this blocks transition

[... more BLOCKING items as needed ...]

### STANDARD EXECUTION TASKS (Required)

- [ ] 2. [Normal state operation]
  - Details
  - Expected outcome

[... more STANDARD items as needed ...]

### EXIT REQUIREMENTS (Must complete before transition)

**NOTE**: These are STANDARD across ALL states - copy exactly

- [ ] 5. Update state file to [NEXT_STATE] per R288
  [... canonical EXIT REQUIREMENTS from template ...]
- [ ] 6. Save TODOs per R287
- [ ] 7. Commit all changes
- [ ] 8. Push to remote
- [ ] 9. Set CONTINUE-SOFTWARE-FACTORY flag per R405
- [ ] 10. Display checkpoint message (if R322 checkpoint)
- [ ] 11. Stop execution (exit 0)
```

**Design per R511:** Follow checklist creation protocol

#### 4. State Purpose Section

```markdown
## State Purpose

[2-4 sentence description of what this state does and why it exists]

**Primary Goal:** [What must be accomplished]
**Key Actions:** [What happens in this state]
**Success Outcome:** [What state produces]
```

#### 5. Entry Criteria Section

```markdown
## Entry Criteria

- **From**: [PREVIOUS_STATE(S)]
- **Condition**: [What must be true to enter this state]
- **Required**:
  - [Required condition 1]
  - [Required condition 2]
  - [etc.]
```

#### 6. State Actions Section

```markdown
## State Actions

### 1. [Primary Action Name]

[Description of the main action this state performs]

**Implementation:**
- [Step or requirement 1]
- [Step or requirement 2]

**Validation:**
- [How to verify this action succeeded]

[Code examples if helpful - optional]

### 2. [Secondary Action if applicable]

[Description]
```

#### 7. Exit Criteria Section

```markdown
## Exit Criteria

### Success Path → [NEXT_STATE]

- [What must be true for successful transition]
- [Required outcomes]
- [Validation checks that must pass]

### Failure Scenarios

- **[Failure Type 1]** → ERROR_RECOVERY
  - Condition: [When this happens]
  - Action: [What to do]
  - Recovery: [How to recover]

- **[Failure Type 2]** → [ALTERNATE_STATE]
  - Condition: [When this happens]
  - Action: [What to do]
```

#### 8. Rules Enforced Section

```markdown
## Rules Enforced

- R[XXX]: [Rule summary - what this state enforces]
- R[YYY]: [Rule summary]
- R510: State Execution Checklist Compliance (this file)
- R511: Checklist Creation Protocol (checklist design)
```

#### 9. Transition Rules Section

```markdown
## Transition Rules

- **ALWAYS** → [NEXT_STATE] (after successful checklist completion)
- **NEVER** skip to [INVALID_STATE] (explain why invalid)
- **CONDITIONAL** → [ALTERNATE_STATE] (if [condition])
- **ERROR** → ERROR_RECOVERY (if [failure scenario])
```

#### 10. R405 Automation Flag Section

```markdown
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

[... canonical R405 section from template ...]
```

#### 11. Additional Context Section (Optional)

```markdown
## Additional Context

[Any state-specific notes, warnings, or context]
[Common pitfalls]
[Related documentation]
```

### Section Ordering Requirements

```yaml
MUST be in this order:
  1. Header + Rule Reading Warning
  2. PRIMARY DIRECTIVES
  3. MANDATORY EXECUTION CHECKLIST
  4. State Purpose
  5. Entry Criteria
  6. State Actions
  7. Exit Criteria
  8. Rules Enforced
  9. Transition Rules
  10. R405 Section
  11. Additional Context (optional)

violation_penalty: -30% for wrong order or missing sections
```

---

## 🔴🔴🔴 PART 4: STATE MACHINE INTEGRATE_WAVE_EFFORTS 🔴🔴🔴

### State Must Exist in State Machine JSON

Before creating state directory and files:

1. **Check state machine:**
```bash
STATE_NAME="YOUR_NEW_STATE"
STATE_MACHINE="$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json"

# Verify state exists
if ! jq -e ".states.\"$STATE_NAME\"" "$STATE_MACHINE" > /dev/null; then
    echo "❌ R516 VIOLATION: State $STATE_NAME not in state machine"
    exit 516
fi
```

2. **Verify transitions defined:**
```bash
# Check allowed transitions
TRANSITIONS=$(jq -r ".states.\"$STATE_NAME\".allowed_transitions[]?" "$STATE_MACHINE")

if [ -z "$TRANSITIONS" ]; then
    echo "❌ R516 VIOLATION: No transitions defined for $STATE_NAME"
    exit 516
fi
```

3. **Validate agent assignment:**
```bash
# Check agent is correct
ASSIGNED_AGENT=$(jq -r ".states.\"$STATE_NAME\".agent" "$STATE_MACHINE")

if [ "$ASSIGNED_AGENT" != "orchestrator" ] && \
   [ "$ASSIGNED_AGENT" != "sw-engineer" ] && \
   [ "$ASSIGNED_AGENT" != "code-reviewer" ] && \
   [ "$ASSIGNED_AGENT" != "architect" ] && \
   [ "$ASSIGNED_AGENT" != "integration-agent" ]; then
    echo "❌ R516 VIOLATION: Invalid agent $ASSIGNED_AGENT for $STATE_NAME"
    exit 516
fi
```

### State Machine Entry Requirements

When adding state to state machine JSON:

```json
{
  "states": {
    "YOUR_NEW_STATE": {
      "description": "Clear description of state purpose",
      "agent": "orchestrator",
      "checkpoint": true,
      "allowed_transitions": [
        "NEXT_STATE",
        "ERROR_RECOVERY"
      ],
      "requires": {
        "conditions": [
          "Condition that must be true to enter"
        ],
        "from_states": [
          "PREVIOUS_STATE"
        ]
      }
    }
  }
}
```

---

## 🔴🔴🔴 PART 5: TEMPLATE USAGE REQUIREMENTS 🔴🔴🔴

### Use Official Template

**Template File:** `/home/vscode/software-factory-template/templates/state-rules-checklist-template.md`

**MANDATORY:** All new states MUST start from this template

### Template Customization Process

```bash
# Step 1: Copy template to new state directory
AGENT="orchestrator"
STATE_NAME="YOUR_NEW_STATE"
STATE_DIR="$CLAUDE_PROJECT_DIR/agent-states/software-factory/$AGENT/$STATE_NAME"

mkdir -p "$STATE_DIR"
cp "$CLAUDE_PROJECT_DIR/templates/state-rules-checklist-template.md" \
   "$STATE_DIR/rules.md"

# Step 2: Replace placeholders
sed -i "s/\[STATE_NAME\]/$STATE_NAME/g" "$STATE_DIR/rules.md"
sed -i "s/\[agent-type\]/$AGENT/g" "$STATE_DIR/rules.md"

# Step 3: Customize for specific state
# - Add state-specific rules to PRIMARY DIRECTIVES
# - Design BLOCKING/STANDARD checklist items per R511
# - Fill in State Purpose, Entry/Exit Criteria
# - Define State Actions
# - Specify transition rules
```

### Template Placeholder Reference

```markdown
Placeholders to replace:
  [STATE_NAME]     → Actual state name (e.g., SPAWN_CODE_REVIEWER)
  [agent-type]     → Agent name (orchestrator, sw-engineer, etc.)
  [NEXT_STATE]     → Primary success transition state
  [PREVIOUS_STATE] → State(s) that transition into this one
  [XXX]            → Specific rule numbers
  [Primary action] → What this state does
```

---

## 🔴🔴🔴 PART 6: VALIDATION REQUIREMENTS 🔴🔴🔴

### Pre-Deployment Validation Checklist

Before deploying new state, MUST validate:

#### Structure Validation

```bash
#!/bin/bash
# File: validate-new-state.sh

validate_state_structure() {
    local state_file="$1"
    local errors=0

    echo "🔍 R516: Validating state structure..."

    # Check required sections exist
    required_sections=(
        "MANDATORY STATE RULE READING"
        "PRIMARY DIRECTIVES"
        "MANDATORY EXECUTION CHECKLIST"
        "State Purpose"
        "Entry Criteria"
        "State Actions"
        "Exit Criteria"
        "Rules Enforced"
        "Transition Rules"
        "R405.*AUTOMATION CONTINUATION FLAG"
    )

    for section in "${required_sections[@]}"; do
        if ! grep -q "$section" "$state_file"; then
            echo "❌ Missing required section: $section"
            ((errors++))
        fi
    done

    # Check checklist structure (R510)
    if ! grep -q "### BLOCKING REQUIREMENTS" "$state_file"; then
        echo "❌ Missing BLOCKING REQUIREMENTS section"
        ((errors++))
    fi

    if ! grep -q "### EXIT REQUIREMENTS" "$state_file"; then
        echo "❌ Missing EXIT REQUIREMENTS section"
        ((errors++))
    fi

    # Check R510 reference
    if ! grep -q "R510.*State Execution Checklist" "$state_file"; then
        echo "❌ Missing R510 reference"
        ((errors++))
    fi

    if [ $errors -eq 0 ]; then
        echo "✅ R516: State structure valid"
        return 0
    else
        echo "❌ R516: Found $errors structural issues"
        return 1
    fi
}
```

#### State Machine Integration Validation

```bash
validate_state_machine_integration() {
    local state_name="$1"
    local state_machine="$2"
    local errors=0

    echo "🔍 R516: Validating state machine integration..."

    # State exists in state machine
    if ! jq -e ".states.\"$state_name\"" "$state_machine" > /dev/null 2>&1; then
        echo "❌ State $state_name not found in state machine"
        ((errors++))
    fi

    # Has transitions defined
    local transition_count=$(jq ".states.\"$state_name\".allowed_transitions | length" "$state_machine")
    if [ "$transition_count" -eq 0 ]; then
        echo "❌ No transitions defined for $state_name"
        ((errors++))
    fi

    # Has valid agent assignment
    local agent=$(jq -r ".states.\"$state_name\".agent" "$state_machine")
    if [[ ! "$agent" =~ ^(orchestrator|sw-engineer|code-reviewer|architect|integration-agent)$ ]]; then
        echo "❌ Invalid agent assignment: $agent"
        ((errors++))
    fi

    if [ $errors -eq 0 ]; then
        echo "✅ R516: State machine integration valid"
        return 0
    else
        echo "❌ R516: Found $errors integration issues"
        return 1
    fi
}
```

#### Checklist Quality Validation

```bash
validate_checklist_quality() {
    local state_file="$1"
    local errors=0

    echo "🔍 R516: Validating checklist quality (R511 compliance)..."

    # Check for command-level granularity
    if grep -q "^- \[ \] [0-9]*\. Run:" "$state_file"; then
        echo "❌ R511 VIOLATION: Command-level granularity detected"
        ((errors++))
    fi

    # Check for validation criteria
    local checklist_items=$(grep -c "^- \[ \] [0-9]*\." "$state_file")
    local validation_count=$(grep -c "Validation:" "$state_file")

    if [ $validation_count -lt $checklist_items ]; then
        echo "⚠️  WARNING: Some checklist items missing validation criteria"
    fi

    # Check EXIT REQUIREMENTS match canonical
    if ! grep -q "Update state file to.*R288" "$state_file"; then
        echo "❌ EXIT REQUIREMENTS missing R288 state update"
        ((errors++))
    fi

    if ! grep -q "Save TODOs.*R287" "$state_file"; then
        echo "❌ EXIT REQUIREMENTS missing R287 TODO save"
        ((errors++))
    fi

    if [ $errors -eq 0 ]; then
        echo "✅ R516: Checklist quality valid"
        return 0
    else
        echo "❌ R516: Found $errors checklist quality issues"
        return 1
    fi
}
```

### Complete Validation Script

```bash
#!/bin/bash
# File: utilities/validate-state-creation.sh

validate_new_state() {
    local state_name="$1"
    local agent="$2"

    local state_file="$CLAUDE_PROJECT_DIR/agent-states/software-factory/$agent/$state_name/rules.md"
    local state_machine="$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json"

    echo "🏭 R516: Validating new state: $state_name"
    echo "=================================="

    local total_errors=0

    # Structure validation
    validate_state_structure "$state_file"
    ((total_errors += $?))

    # State machine integration
    validate_state_machine_integration "$state_name" "$state_machine"
    ((total_errors += $?))

    # Checklist quality
    validate_checklist_quality "$state_file"
    ((total_errors += $?))

    echo "=================================="
    if [ $total_errors -eq 0 ]; then
        echo "✅✅✅ R516: State $state_name FULLY VALIDATED ✅✅✅"
        return 0
    else
        echo "❌❌❌ R516: State $state_name FAILED VALIDATION ❌❌❌"
        echo "Total validation failures: $total_errors"
        return 1
    fi
}

# Run validation
validate_new_state "$1" "$2"
```

---

## 🔴🔴🔴 PART 7: THE STATE CREATION WORKFLOW 🔴🔴🔴

### Step-by-Step Creation Process

When creating a new state from scratch:

#### STEP 1: Justify the New State

```yaml
questions_to_answer:
  - Why does this state need to exist?
  - What existing state could NOT handle this?
  - What is the distinct purpose?
  - Who is the responsible agent?
  - Where in workflow does this fit?

document:
  - Create justification document
  - Review with team/architect
  - Get approval before proceeding
```

#### STEP 2: Update State Machine First

```bash
# CRITICAL: State machine BEFORE state directory

# 1. Edit state machine JSON
vim $CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json

# 2. Add state entry
{
  "YOUR_NEW_STATE": {
    "description": "...",
    "agent": "orchestrator",
    "checkpoint": true/false,
    "allowed_transitions": ["NEXT_STATE", "ERROR_RECOVERY"],
    "requires": {...}
  }
}

# 3. Update transitions in PREVIOUS states
# Add YOUR_NEW_STATE to allowed_transitions of states that lead here

# 4. Validate JSON
jq '.' $CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json > /dev/null
```

#### STEP 3: Create State Directory and Copy Template

```bash
# Determine paths
AGENT="orchestrator"  # or sw-engineer, code-reviewer, architect
STATE_NAME="YOUR_NEW_STATE"
STATE_DIR="$CLAUDE_PROJECT_DIR/agent-states/software-factory/$AGENT/$STATE_NAME"

# Create directory
mkdir -p "$STATE_DIR"

# Copy template
cp "$CLAUDE_PROJECT_DIR/templates/state-rules-checklist-template.md" \
   "$STATE_DIR/rules.md"
```

#### STEP 4: Replace Template Placeholders

```bash
# Automated replacement
sed -i "s/\[STATE_NAME\]/$STATE_NAME/g" "$STATE_DIR/rules.md"
sed -i "s/\[agent-type\]/$AGENT/g" "$STATE_DIR/rules.md"

# Manual replacements needed:
# - [NEXT_STATE] → Actual next state name
# - [PREVIOUS_STATE] → Actual previous state(s)
# - [XXX] → Specific rule numbers
# - [Primary action] → State's main action
```

#### STEP 5: Design State-Specific Rules

```yaml
primary_directives:
  - Add state-specific rules
  - Ensure universal rules present (R510, R288, R287, R405)
  - Add agent-specific rules
  - Reference rules by file path

example:
  - "3. **🔴🔴🔴 R514** - Infrastructure Creation Protocol
       - File: $CLAUDE_PROJECT_DIR/rule-library/R514-...md
       - Criticality: SUPREME LAW
       - Summary: Must create branches following cascade pattern"
```

#### STEP 6: Design Execution Checklist (per R511)

```yaml
follow_r511_protocol:
  blocking_requirements:
    - What MUST complete for state to succeed?
    - Use outcome-focused granularity
    - Add validation criteria
    - Mark with **BLOCKING**

  standard_tasks:
    - What supports the main purpose?
    - What's nice-to-have but not critical?
    - Informational items

  exit_requirements:
    - Copy canonical EXIT REQUIREMENTS exactly
    - Replace [NEXT_STATE] placeholder only
    - DO NOT modify other EXIT text
```

#### STEP 7: Fill State Metadata

```markdown
State Purpose:
  - 2-4 sentence description
  - Primary goal statement
  - Key actions summary

Entry Criteria:
  - From which state(s)?
  - What conditions required?
  - What must be true?

State Actions:
  - List each major action
  - Provide implementation details
  - Add validation steps

Exit Criteria:
  - Success path definition
  - Failure scenarios
  - Recovery procedures

Transition Rules:
  - ALWAYS transitions
  - NEVER transitions
  - CONDITIONAL transitions
```

#### STEP 8: Validate Before Deployment

```bash
# Run comprehensive validation
bash $CLAUDE_PROJECT_DIR/utilities/validate-state-creation.sh \
    "$STATE_NAME" "$AGENT"

# Fix any issues found
# Re-validate until clean
```

#### STEP 9: Test State in Isolation

```bash
# Manual testing:
# 1. Spawn agent in this state manually
# 2. Verify agent reads rules correctly
# 3. Check agent executes checklist
# 4. Verify transitions work
# 5. Test error scenarios

# Document test results
```

#### STEP 10: Commit and Deploy

```bash
# Add to git
git add agent-states/software-factory/$AGENT/$STATE_NAME/
git add state-machines/software-factory-3.0-state-machine.json

# Commit with descriptive message
git commit -m "feat: Add $STATE_NAME state for $AGENT

PURPOSE: [Why this state exists]

STRUCTURE:
- Follows R516 state creation protocol
- R510/R511 compliant checklist
- Template-based with customizations

TRANSITIONS:
- From: [PREVIOUS_STATE]
- To: [NEXT_STATE]

VALIDATION: Passed utilities/validate-state-creation.sh

[Additional context]"

# Push
git push
```

---

## 🔴🔴🔴 PART 8: COMMON MISTAKES AND FIXES 🔴🔴🔴

### Mistake 1: Creating State Before State Machine

```markdown
❌ WRONG ORDER:
1. Create state directory and rules.md
2. Try to reference state in agent code
3. State machine doesn't know about it
4. Transitions fail

✅ RIGHT ORDER:
1. Update state machine JSON first
2. Define transitions
3. Assign agent
4. THEN create state directory
5. Create rules.md from template
```

**Fix:** Always state machine first, state directory second

### Mistake 2: Missing Required Sections

```markdown
❌ WRONG (incomplete):
# STATE_NAME Rules
## Checklist
- Do the work
## Done

✅ RIGHT (complete):
# STATE_NAME Rules
## Rule Reading Warning
## PRIMARY DIRECTIVES
## MANDATORY EXECUTION CHECKLIST
## State Purpose
## Entry Criteria
## State Actions
## Exit Criteria
## Rules Enforced
## Transition Rules
## R405 Section
```

**Fix:** Use template, don't create from scratch

### Mistake 3: Modifying EXIT REQUIREMENTS

```markdown
❌ WRONG (customized EXIT):
- [ ] 5. Update state if needed
- [ ] 6. Maybe commit
- [ ] 7. Set flag

✅ RIGHT (canonical EXIT):
- [ ] 5. Update state file to [NEXT_STATE] per R288
  [exact canonical text from template]
- [ ] 6. Save TODOs per R287
  [exact canonical text]
- [ ] 7. Commit all changes
  [exact canonical text]
```

**Fix:** Copy EXIT REQUIREMENTS exactly from template

### Mistake 4: Ad-Hoc State Names

```markdown
❌ WRONG:
- do-review
- handle_stuff
- TempState
- State1

✅ RIGHT:
- SPAWN_CODE_REVIEWER
- WAITING_FOR_REVIEW
- CREATE_EFFORT_INFRASTRUCTURE
- INTEGRATE_PHASE_WAVES
```

**Fix:** Follow naming conventions in Part 2

### Mistake 5: No Validation Criteria

```markdown
❌ WRONG:
- [ ] 1. Spawn reviewer

✅ RIGHT:
- [ ] 1. Spawn Code Reviewer for effort quality assessment
  - Agent: code-reviewer
  - Workspace: ${effort_workspace}
  - Validation: Agent ID returned with timestamp
  - **BLOCKING**: Cannot transition without reviewer
```

**Fix:** Add validation to every checklist item

---

## 🔴🔴🔴 PART 9: ENFORCEMENT AND GRADING 🔴🔴🔴

### Validation Points

States are validated at:

1. **Pre-commit hooks** - Structure validation
2. **State machine transitions** - Integration validation
3. **Agent startup** - Rule reading verification (R290)
4. **State execution** - Checklist compliance (R510)
5. **Code review** - Manual quality check

### Grading Criteria

```yaml
state_creation_grading:
  structure_compliance:
    all_required_sections: 20%
    correct_section_order: 10%
    template_usage: 10%

  state_machine_integration:
    exists_in_state_machine: 15%
    valid_transitions: 10%
    correct_agent_assignment: 5%

  checklist_quality:
    r511_compliance: 15%
    proper_granularity: 10%
    validation_criteria: 5%

total: 100%

violations:
  missing_required_section: -20%        # Per section
  wrong_section_order: -10%
  no_template_usage: -30%
  not_in_state_machine: -100%           # FATAL
  invalid_transitions: -50%
  modified_exit_requirements: -100%     # FATAL
  command_level_checklist: -50%
  no_validation_criteria: -20%
```

### Automatic Failure Conditions

```yaml
automatic_failure:
  - State not in state machine JSON
  - Modified EXIT REQUIREMENTS
  - Missing MANDATORY EXECUTION CHECKLIST
  - Missing R510 reference
  - No transitions defined
  - Invalid agent assignment
  - No validation run before deployment
```

---

## 🔴 THE STATE CREATION OATH

```
I, the State Designer, swear by R516:

When I create or modify a state:
- I WILL update state machine JSON first
- I WILL use the official template
- I WILL include ALL required sections in correct order
- I WILL design checklist per R511 protocol
- I WILL define clear validation criteria
- I WILL validate before deploying
- I WILL NOT create ad-hoc states
- I WILL NOT skip required sections
- I WILL NOT modify canonical EXIT REQUIREMENTS

My states follow the protocol.
My states integrate with the system.
My states are validated and tested.

This is SUPREME LAW.
Improper state creation means FAILURE.
I WILL CREATE STATES PROPERLY.
```

---

**Remember:** States are the building blocks of the Software Factory state machine. Creating them properly ensures agents know exactly what to do, the system can validate correctness, and workflows execute reliably. R516 + R510 + R511 together ensure every state is well-designed, properly structured, and fully enforceable.

**See Also:**
- R510: State Execution Checklist Compliance (structure)
- R511: Checklist Creation Protocol (checklist design)
- R288: State File Update Requirements
- R287: TODO Persistence
- R322: Mandatory Checkpoints
- R405: Automation Continuation Flag
- R290: State Rule Reading Verification
- R234: Mandatory State Traversal
- Template: `/home/vscode/software-factory-template/templates/state-rules-checklist-template.md`
- State Machine: `/home/vscode/software-factory-template/state-machines/software-factory-3.0-state-machine.json`
- Validation Script: `/home/vscode/software-factory-template/utilities/validate-state-creation.sh`

## State Manager Coordination (SF 3.0)

State Manager uses R516 naming conventions for state machine operations:
- **State names** in `state-machines/software-factory-3.0-state-machine.json` follow R516 patterns
- **Current state** in `orchestrator-state-v3.json` uses R516-compliant names
- **State directories** use R516 naming (e.g., `SETUP_WAVE_INFRASTRUCTURE`)
- **Validation** ensures state names match R516 format

State Manager enforces R516 compliance through:
1. State machine validation (state names must match patterns)
2. State directory validation (must exist for current_state)
3. Transition validation (next_state must be R516-compliant)

See: `state-machines/software-factory-3.0-state-machine.json`, R206 (state validation)
