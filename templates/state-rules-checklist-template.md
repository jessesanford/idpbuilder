# [STATE_NAME] State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

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

[Add other core rules...]

### State-Specific Rules:

[Add state-specific rules...]

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. [Primary critical action that MUST complete]
  - Parameter: [value]
  - Parameter: [value]
  - Validation: [How to verify completion]
  - **BLOCKING**: Cannot transition without this

- [ ] 2. [Secondary critical action if applicable]
  - Context: [details]
  - Validation: [criteria]
  - **BLOCKING**: Required for state success

[Add more BLOCKING items as needed - be specific to this state]

### STANDARD EXECUTION TASKS (Required)

- [ ] 3. [Normal state operation]
  - Details: [what to do]
  - Expected outcome: [what should result]

- [ ] 4. [Additional standard task if applicable]
  - Details: [what to do]

[Add more STANDARD items as needed]

### EXIT REQUIREMENTS (Must complete before transition)

**NOTE**: These are STANDARD across ALL states - copy exactly

- [ ] 5. Update state file to [NEXT_STATE] per R288
  - Field: `current_state`
  - Value: `"[NEXT_STATE]"`
  - Also update: `previous_state`, `transition_time`, `transition_reason`
  - Validation: `jq '.current_state' orchestrator-state-v3.json` shows new state

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
  - Include: Next steps for /continue-orchestrating
  - Refer to: R322 requirements

- [ ] 11. Stop execution (exit 0)
  - Command: `exit 0`
  - Timing: After ALL above items complete
  - Per: R313 (if spawn state), R322 (if checkpoint)

---

## State Purpose

[Brief description of what this state does and why it exists]

---

## Entry Criteria

- **From**: [PREVIOUS_STATE(S)]
- **Condition**: [What must be true to enter this state]
- **Required**:
  - [Required condition 1]
  - [Required condition 2]

---

## State Actions

### 1. [Primary Action]

[Description of the main thing this state does]

```bash
# Example code/commands if helpful
[command examples]
```

### 2. [Secondary Action if applicable]

[Description]

---

## Exit Criteria

### Success Path → [NEXT_STATE]

- [What must be true for success]
- [Required outcomes]
- [Validation checks]

### Failure Scenarios

- **[Failure Type 1]** → ERROR_RECOVERY
  - Condition: [when this happens]
  - Action: [what to do]

- **[Failure Type 2]** → ERROR_RECOVERY
  - Condition: [when this happens]
  - Action: [what to do]

---

## Rules Enforced

- R[XXX]: [Rule summary]
- R[XXX]: [Rule summary]
- R510: State Execution Checklist Compliance (this file)

---

## Transition Rules

- **ALWAYS** → [NEXT_STATE] (after checklist complete)
- **NEVER** skip to [INVALID_STATE] (explain why)
- **CONDITIONAL** → [CONDITIONAL_STATE] (if [condition])

---

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-[agent] to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

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
- ✅ State work completed successfully
- ✅ Ready for /continue-[agent]
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs

**See**: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`

---

## Additional Context

[Any additional notes, warnings, or context specific to this state]

---

**Template Version**: 1.0
**Created**: 2025-10-07
**Purpose**: Standardized state rules format with R510 checklist compliance
