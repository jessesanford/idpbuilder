# 🔴🔴🔴 RULE R510 - STATE EXECUTION CHECKLIST COMPLIANCE (SUPREME LAW)

**Criticality:** SUPREME LAW - ABSOLUTE REQUIREMENT
**Grading Impact:** -50% to -100% for violations
**Enforcement:** CONTINUOUS - Every state must have and follow checklist

---

## SUPREME LAW STATEMENT

**EVERY STATE MUST HAVE AN EXECUTION CHECKLIST. AGENTS MUST ACKNOWLEDGE COMPLETION OF EVERY CHECKLIST ITEM BEFORE TRANSITIONING STATES. SKIPPING CHECKLIST ITEMS IS AUTOMATIC FAILURE.**

---

## 🚨🚨🚨 THE ABSOLUTE MANDATE 🚨🚨🚨

### STATE EXECUTION CHECKLISTS ARE:
```markdown
✅ MANDATORY for every state
✅ EXPLICIT requirements that MUST be completed
✅ BLOCKING mechanisms that prevent premature transitions
✅ ACCOUNTABILITY tools that create audit trails
✅ INTEGRATE_WAVE_EFFORTS points with R232 (TodoWrite pending override)
```

### STATE EXECUTION CHECKLISTS ARE NOT:
```markdown
❌ Suggestions or recommendations
❌ Optional guidelines
❌ Things that can be skipped
❌ Decorative documentation
❌ Flexible or negotiable
```

---

## 🔴🔴🔴 CRITICAL: BLOCKING ITEMS ARE EXECUTABLE ACTIONS 🔴🔴🔴

### The Execute-Don't-Defer Principle

**BLOCKING requirements are ACTIONS to perform NOW, not notes to remember later.**

❌ **WRONG Pattern** (causes stalls):
```bash
# Read checklist: "5. Spawn integration agent"
echo "✅ I identified that integration agent should be spawned"
# Create TODO: "⏹️ Spawn integration agent"
# Stop with: CONTINUE-SOFTWARE-FACTORY=FALSE
```

✅ **CORRECT Pattern** (seamless flow):
```bash
# Read checklist: "5. Spawn integration agent"
echo "🚀 Spawning integration agent NOW..."
Task: integration-agent
State: EXECUTE_WAVE_INTEGRATION
Workspace: ${integration_workspace}
# Wait for completion
echo "✅ CHECKLIST[5]: Spawned integration agent [agent-id-timestamp]"
# Continue with: CONTINUE-SOFTWARE-FACTORY=TRUE
```

### Common Misinterpretation

**Checklist says**: "Spawn integration agent to perform integration work"

**WRONG interpretation**: "I should note that an integration agent needs to be spawned"
- Creates pending TODO
- Stops execution
- Waits for human intervention
- **Result**: System stall

**CORRECT interpretation**: "I will use the Task tool right now to spawn the integration agent"
- Invokes Task tool immediately
- Provides agent parameters
- Waits for agent completion
- Acknowledges completion
- **Result**: Seamless flow

### Rule: If You Can Do It Now, Do It Now

Before stopping with CONTINUE-SOFTWARE-FACTORY=FALSE, ask:
1. Can I execute this checklist item with available tools? → YES → **DO IT NOW**
2. Is this blocked by external dependency? → NO → **DO IT NOW**
3. Does this require human decision? → NO → **DO IT NOW**

**Only stop with FALSE if work is truly IMPOSSIBLE**, not just incomplete.

### Examples Across States

**SPAWN_SW_ENGINEERS**: "Spawn 3 SW Engineers for implementation"
- ❌ WRONG: Note "need to spawn engineers" → Stop
- ✅ RIGHT: Spawn 3 agents with Task tool → Acknowledge → Continue

**CREATE_WAVE_FIX_PLAN**: "Spawn Code Reviewer to create fix plan"
- ❌ WRONG: Create TODO "spawn reviewer" → Stop
- ✅ RIGHT: Task: code-reviewer → Wait → Continue

**INTEGRATE_WAVE_EFFORTS**: "Spawn integration agent to perform merges"
- ❌ WRONG: Identify "agent not spawned" → Stop with FALSE
- ✅ RIGHT: Task: integration-agent → Monitor → TRUE

---

## 🔴🔴🔴 CHECKLIST STRUCTURE REQUIREMENTS 🔴🔴🔴

### Mandatory Three-Section Format

Every state MUST have this structure:

```markdown
## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

### BLOCKING REQUIREMENTS (Cannot proceed without)
- [ ] 1. [Critical action description]
  - Context/parameters
  - Validation criteria
  - **BLOCKING**: Cannot transition without this

### STANDARD EXECUTION TASKS (Required)
- [ ] 2. [Normal state operation]
  - Context/details
  - Expected outcome

### EXIT REQUIREMENTS (Must complete before transition)
- [ ] 3. Update state file to NEXT_STATE per R288
- [ ] 4. Save TODOs per R287 (within 30s of last TodoWrite)
- [ ] 5. Commit and push all changes
- [ ] 6. Display checkpoint message per R322 (if applicable)
- [ ] 7. Set CONTINUE-SOFTWARE-FACTORY flag per R405
- [ ] 8. Stop execution (exit 0)
```

### Section-Specific Requirements

#### BLOCKING REQUIREMENTS
- **MUST** contain critical state-specific actions
- **MUST** clearly mark items as "**BLOCKING**"
- **MUST** define validation criteria
- **CANNOT** be empty (if state has no blocking items, state design is wrong)

#### STANDARD EXECUTION TASKS
- **MUST** contain normal state operations
- **SHOULD** be completable without blocking transition
- **CAN** be empty if state is pure checkpoint/waiting

#### EXIT REQUIREMENTS
- **MUST** be identical across all states (standardized)
- **MUST** include R288 state update
- **MUST** include R287 TODO persistence
- **MUST** include R405 continuation flag
- **MUST** include exit 0 stop

---

## 📋 ACKNOWLEDGMENT PROTOCOL (MANDATORY)

### Format Requirements

When completing EACH checklist item, agent MUST output:
```
✅ CHECKLIST[n]: [Description] [proof]
```

### Acknowledgment Examples

```bash
# BLOCKING REQUIREMENT acknowledgment
✅ CHECKLIST[1]: Spawned Code Reviewer [agent-code-reviewer-phase1-20251007-020300]

# STANDARD TASK acknowledgment
✅ CHECKLIST[2]: Updated state field 'phase_integration_review' with reviewer details

# EXIT REQUIREMENT acknowledgments
✅ CHECKLIST[3]: Updated state to WAITING_FOR_PHASE_REVIEW_WAVE_INTEGRATION per R288
✅ CHECKLIST[4]: Saved TODOs to todos/orchestrator-PHASE_REVIEW_WAVE_INTEGRATION-20251007-020315.todo per R287
✅ CHECKLIST[5]: Committed state changes [commit: a9916de]
✅ CHECKLIST[6]: Pushed changes to origin/integration-agent-fixes
✅ CHECKLIST[7]: CONTINUE-SOFTWARE-FACTORY=TRUE (successful completion, factory can proceed)
✅ CHECKLIST[8]: Stopping execution with exit 0 per R322
```

### Proof Requirements

| Item Type | Required Proof Format |
|-----------|----------------------|
| Spawn agent | Agent ID with timestamp |
| Update state | New state name or field updated |
| Save file | File path with timestamp |
| Commit changes | Commit hash (short form) |
| Push changes | Remote branch name |
| Set flag | Flag value (TRUE/FALSE with reason) |
| Stop execution | exit code |

### ENFORCEMENT

```bash
# Before ANY state transition, verify acknowledgments
verify_checklist_acknowledgments() {
    local state="$1"
    local expected_items=$(count_checklist_items "$state")
    local acknowledged_items=$(grep -c "✅ CHECKLIST\[" conversation.log)

    if [ $acknowledged_items -lt $expected_items ]; then
        echo "❌❌❌ R510 VIOLATION ❌❌❌"
        echo "State: $state"
        echo "Expected: $expected_items checklist acknowledgments"
        echo "Found: $acknowledged_items acknowledgments"
        echo "MISSING: $((expected_items - acknowledged_items)) items"
        echo ""
        echo "AUTOMATIC FAILURE - CANNOT TRANSITION"
        return 1
    fi

    echo "✅ All checklist items acknowledged ($acknowledged_items/$expected_items)"
    return 0
}
```

---

## 🚨🚨🚨 INTEGRATE_WAVE_EFFORTS WITH R232 (CRITICAL) 🚨🚨🚨

### How Checklists Enforce R232

**R232 Problem**: Agents create PENDING TODOs and stop without executing
**R510 Solution**: Checklists make execution explicit and trackable

### Integration Pattern

```markdown
1. Agent enters STATE
2. Agent reads state rules including checklist
3. Checklist items loaded as TodoWrite PENDING items
4. R232 prevents stopping with pending todos
5. Agent MUST complete each item
6. Agent acknowledges each completion
7. All items completed → TodoWrite empty
8. R232 allows stopping
9. Agent transitions
```

### Example Integration

```bash
# State: PHASE_REVIEW_WAVE_INTEGRATION
# Checklist item 1: Spawn Code Reviewer

# WRONG (R232 violation):
TodoWrite:
  - ⏳ Spawn Code Reviewer for phase integration
[Agent stops] ❌❌❌ VIOLATION

# RIGHT (R510 + R232 compliance):
TodoWrite:
  - ⏳ Spawn Code Reviewer for phase integration
[Agent spawns reviewer]
✅ CHECKLIST[1]: Spawned Code Reviewer [agent-code-reviewer-20251007]
[TodoWrite item marked complete]
[Agent continues with next checklist item]
```

---

## 🔴 GRANULARITY GUIDELINES

### Too High-Level ❌

```markdown
- [ ] 1. Complete the state
- [ ] 2. Do everything required
```

**Problem**: No accountability, no validation criteria

### Too Low-Level ❌

```markdown
- [ ] 1. Run: cd /efforts/phase1
- [ ] 2. Run: /spawn agent-code-reviewer --state REVIEW
- [ ] 3. Run: echo "spawned" >> log.txt
- [ ] 4. Run: jq '.state = "NEXT"' state.json
```

**Problem**: Command-level micromanagement, no focus on outcomes

### Right Level ✅

```markdown
- [ ] 1. Spawn Code Reviewer for phase integration quality review
  - Agent: code-reviewer
  - Focus: phase-integration-quality
  - Workspace: ${phase_integration_workspace}
  - Validation: Agent ID returned with timestamp
  - **BLOCKING**: Must complete before transition

- [ ] 2. Update state file with phase integration review metadata
  - Fields: reviewer, branch, focus, waves_integrated, started_at
  - Per: R288 state update requirements

- [ ] 3. Verify wave integration branches present
  - Check: W1, W2, W3 all merged into phase branch
  - Validation: git log shows wave merge commits
```

**Benefits**:
- Outcome-focused (what must be achieved)
- Clear validation criteria (how to verify)
- Appropriate detail (enough to execute, not command-level)
- Explicit blocking markers (what's critical)

---

## ❌ VIOLATION EXAMPLES

### VIOLATION 1: Skipping BLOCKING Item

```bash
# Checklist has BLOCKING item:
# - [ ] 1. Spawn Code Reviewer [BLOCKING]

# Agent behavior:
echo "Moving to next state..."
jq '.state_machine.current_state = "WAITING_FOR_REVIEW"' state.json
git add . && git commit -m "state update" && git push

# ❌❌❌ CATASTROPHIC VIOLATION ❌❌❌
# - No spawn occurred
# - BLOCKING item ignored
# - No acknowledgment
# GRADE: -100% IMMEDIATE FAILURE
```

### VIOLATION 2: Missing Acknowledgment

```bash
# Checklist completed but no acknowledgment:
[Agent spawns code reviewer]
[Agent updates state]
[Agent transitions]

# ❌❌❌ R510 VIOLATION ❌❌❌
# - Work done but no proof
# - No audit trail
# - Cannot verify compliance
# GRADE: -20% PER MISSING ACKNOWLEDGMENT
```

### VIOLATION 3: Incomplete EXIT Requirements

```bash
# Agent completes BLOCKING items
✅ CHECKLIST[1]: Spawned reviewer
# But skips EXIT requirements:
[No state update]
[No TODO save]
[No commit/push]

# ❌❌❌ R510 + R288 + R287 VIOLATIONS ❌❌❌
# GRADE: -100% SYSTEMATIC FAILURE
```

---

## ✅ CORRECT BEHAVIOR EXAMPLES

### CORRECT 1: Full Checklist Compliance

```bash
# State: PHASE_REVIEW_WAVE_INTEGRATION
echo "📋 Executing MANDATORY EXECUTION CHECKLIST..."

# BLOCKING REQUIREMENT
echo "🚀 CHECKLIST[1]: Spawning Code Reviewer for phase integration review..."
/spawn agent-code-reviewer PHASE_REVIEW_WAVE_INTEGRATION \
  --phase "P1" \
  --branch "idpbuilder-oci-mgmt/phase1-integration" \
  --focus "phase-integration-quality"
echo "✅ CHECKLIST[1]: Spawned Code Reviewer [agent-code-reviewer-phase1-20251007-020300]"

# STANDARD TASKS
echo "📝 CHECKLIST[2]: Updating state file with review metadata..."
jq '.phase_integration_review = {
  reviewer: "agent-code-reviewer-phase1-20251007-020300",
  branch: "idpbuilder-oci-mgmt/phase1-integration",
  focus: "phase-integration-quality",
  started_at: "'$(date -Iseconds)'"
}' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json
echo "✅ CHECKLIST[2]: Updated state field 'phase_integration_review'"

# EXIT REQUIREMENTS
echo "📝 CHECKLIST[3]: Updating state to WAITING_FOR_PHASE_REVIEW_WAVE_INTEGRATION..."
jq '.state_machine.current_state = "WAITING_FOR_PHASE_REVIEW_WAVE_INTEGRATION"' \
  orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json
echo "✅ CHECKLIST[3]: Updated state to WAITING_FOR_PHASE_REVIEW_WAVE_INTEGRATION per R288"

echo "💾 CHECKLIST[4]: Saving TODOs per R287..."
save_todos "R510_CHECKLIST_COMPLETE"
echo "✅ CHECKLIST[4]: Saved TODOs to todos/orchestrator-PHASE_REVIEW_WAVE_INTEGRATION-20251007-020315.todo per R287"

echo "📦 CHECKLIST[5]: Committing changes..."
git add orchestrator-state-v3.json todos/
git commit -m "state: PHASE_REVIEW_WAVE_INTEGRATION → WAITING_FOR_PHASE_REVIEW_WAVE_INTEGRATION

- Spawned Code Reviewer for phase integration quality review
- Updated state file with review metadata per R288
- Saved TODOs per R287
- All checklist items completed per R510"
echo "✅ CHECKLIST[5]: Committed changes [commit: $(git rev-parse --short HEAD)]"

echo "🚀 CHECKLIST[6]: Pushing to remote..."
git push
echo "✅ CHECKLIST[6]: Pushed changes to origin/$(git branch --show-current)"

echo "🏭 CHECKLIST[7]: Setting automation continuation flag..."
# 🚨 CRITICAL: Use TRUE here! R322 checkpoint is NORMAL workflow, not a failure
# State work completed successfully → TRUE (even though we're stopping for user)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
echo "✅ CHECKLIST[7]: CONTINUE-SOFTWARE-FACTORY=TRUE (successful spawn, factory can proceed)"

echo "🛑 CHECKLIST[8]: Stopping execution per R322..."
# Note: Stopping at checkpoint (R322) does NOT mean FALSE flag!
# Checkpoint = designed workflow = TRUE flag above
exit 0
echo "✅ CHECKLIST[8]: Stopped with exit 0"
```

### CORRECT 2: Acknowledgment Protocol

```bash
# Clear, verifiable acknowledgments with proof:
✅ CHECKLIST[1]: Spawned Code Reviewer [agent-code-reviewer-phase1-20251007-020300]
✅ CHECKLIST[2]: Updated state field 'phase_integration_review'
✅ CHECKLIST[3]: Updated state to WAITING_FOR_PHASE_REVIEW_WAVE_INTEGRATION per R288
✅ CHECKLIST[4]: Saved TODOs to todos/orchestrator-20251007-020315.todo per R287
✅ CHECKLIST[5]: Committed changes [commit: a9916de]
✅ CHECKLIST[6]: Pushed to origin/integration-agent-fixes
✅ CHECKLIST[7]: CONTINUE-SOFTWARE-FACTORY=TRUE (successful completion)
✅ CHECKLIST[8]: Stopped with exit 0 per R322
```

---

## 🔴 GRADING IMPACT

### Compliance Grading

```yaml
checklist_compliance_grading:
  blocking_items_complete: 40%      # All BLOCKING items done
  standard_items_complete: 30%      # All STANDARD items done
  exit_requirements_complete: 20%   # All EXIT items done
  acknowledgment_protocol: 10%      # Proper format used

total_possible: 100%
```

### Violation Penalties

```yaml
violation_penalties:
  skipped_blocking_item: -50% to -100%     # Per item
  skipped_standard_item: -20%              # Per item
  skipped_exit_requirement: -30% to -100%  # Depends on which
  missing_acknowledgment: -10%             # Per missing ack
  incorrect_ack_format: -5%                # Per incorrect format
  no_proof_in_ack: -10%                    # Per missing proof

catastrophic_violations:
  no_checklist_in_state: -100%             # State missing checklist
  agent_ignores_checklist: -100%           # Agent didn't use it
  transition_without_completion: -100%     # Premature transition
```

---

## 🎯 INTEGRATE_WAVE_EFFORTS WITH OTHER RULES

### R510 + R232 (TodoWrite Pending Override)
- **Relationship**: Checklist items → TodoWrite pending items
- **Enforcement**: R232 prevents stopping with pending, R510 defines what's pending
- **Synergy**: Natural integration, dual enforcement

### R510 + R288 (State File Updates)
- **Relationship**: EXIT REQUIREMENTS always include R288 compliance
- **Enforcement**: Cannot exit without state update
- **Synergy**: Consistent state management

### R510 + R287 (TODO Persistence)
- **Relationship**: EXIT REQUIREMENTS always include R287 save
- **Enforcement**: Cannot exit without saving TODOs
- **Synergy**: Prevents TODO loss at transitions

### R510 + R322 (Mandatory Checkpoints)
- **Relationship**: EXIT REQUIREMENTS include R322 messaging
- **Enforcement**: Checkpoint states must display message
- **Synergy**: User review opportunities preserved

### R510 + R405 (Automation Continuation)
- **Relationship**: EXIT REQUIREMENTS always include R405 flag
- **Enforcement**: Must set operational status flag
- **Synergy**: Clear factory continuation decision

---

## 📢 THE CHECKLIST MANTRA

### Repeat Before Every State:

1. **"Read the checklist for this state"**
2. **"BLOCKING items MUST complete before anything else"**
3. **"Acknowledge EVERY item with ✅ CHECKLIST[n]: proof"**
4. **"EXIT REQUIREMENTS are MANDATORY for every state"**
5. **"No transition without complete checklist"**

---

## 🛡️ STATE-SPECIFIC CUSTOMIZATION

### Spawn States (SPAWN_*)

**BLOCKING REQUIREMENTS:**
- Spawn specified agent(s) with correct parameters
- Verify spawn successful (agent ID returned)
- Record agent metadata in state file

**EXIT REQUIREMENTS:**
- Standard EXIT REQUIREMENTS +
- **MUST** stop after spawn (R313)

### Monitoring/Waiting States (MONITOR_*, WAITING_FOR_*)

**BLOCKING REQUIREMENTS:**
- Often none (waiting states don't have critical actions)
- Or: Verify condition being waited for

**STANDARD TASKS:**
- Display current status
- Show what's being waited for
- Provide next steps guidance

**EXIT REQUIREMENTS:**
- Standard EXIT REQUIREMENTS +
- Display R322 checkpoint message

### Action States (CREATE_*, ANALYZE_*, DISTRIBUTE_*)

**BLOCKING REQUIREMENTS:**
- Complete the named action
- Validate action succeeded
- Create required artifacts

**STANDARD TASKS:**
- Update state file with results
- Create documentation
- Log activity

**EXIT REQUIREMENTS:**
- Standard EXIT REQUIREMENTS

---

## 💪 VALIDATION TOOLS

### Pre-Transition Validation

```bash
# File: utilities/validate-checklist-completion.sh

validate_checklist_completion() {
    local current_state="$1"
    local conversation_log="$2"

    # Load checklist
    local checklist_file="$CLAUDE_PROJECT_DIR/agent-states/software-factory/orchestrator/${current_state}/rules.md"

    if [ ! -f "$checklist_file" ]; then
        echo "❌ No checklist found for state: $current_state"
        return 1
    fi

    # Extract BLOCKING items count
    local blocking_count=$(grep -c "**BLOCKING**" "$checklist_file")

    # Extract acknowledgments count
    local ack_count=$(grep -c "✅ CHECKLIST\[" "$conversation_log")

    # Extract total items count
    local total_items=$(grep -c "^- \[ \] " "$checklist_file")

    echo "📊 Checklist Validation for $current_state:"
    echo "  Total items: $total_items"
    echo "  BLOCKING items: $blocking_count"
    echo "  Acknowledged: $ack_count"

    if [ $ack_count -lt $total_items ]; then
        echo "❌ INCOMPLETE: Missing $((total_items - ack_count)) acknowledgments"
        return 1
    fi

    echo "✅ All checklist items acknowledged"
    return 0
}
```

### State Checklist Validator

```bash
# File: utilities/validate-state-checklists.sh

validate_all_state_checklists() {
    local errors=0

    for state_dir in $CLAUDE_PROJECT_DIR/agent-states/software-factory/orchestrator/*/; do
        state_name=$(basename "$state_dir")
        rules_file="${state_dir}rules.md"

        # Check file exists
        if [ ! -f "$rules_file" ]; then
            echo "❌ $state_name: Missing rules.md"
            ((errors++))
            continue
        fi

        # Check checklist section exists
        if ! grep -q "## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST" "$rules_file"; then
            echo "❌ $state_name: Missing MANDATORY EXECUTION CHECKLIST section"
            ((errors++))
            continue
        fi

        # Check required sections
        if ! grep -q "### BLOCKING REQUIREMENTS" "$rules_file"; then
            echo "❌ $state_name: Missing BLOCKING REQUIREMENTS section"
            ((errors++))
        fi

        if ! grep -q "### EXIT REQUIREMENTS" "$rules_file"; then
            echo "❌ $state_name: Missing EXIT REQUIREMENTS section"
            ((errors++))
        fi

        # Check EXIT REQUIREMENTS content
        if ! grep -q "Update state file.*R288" "$rules_file"; then
            echo "❌ $state_name: EXIT REQUIREMENTS missing R288 state update"
            ((errors++))
        fi

        if ! grep -q "Save TODOs.*R287" "$rules_file"; then
            echo "❌ $state_name: EXIT REQUIREMENTS missing R287 TODO save"
            ((errors++))
        fi

        if ! grep -q "CONTINUE-SOFTWARE-FACTORY.*R405" "$rules_file"; then
            echo "❌ $state_name: EXIT REQUIREMENTS missing R405 flag"
            ((errors++))
        fi

        if [ $errors -eq 0 ]; then
            echo "✅ $state_name: Checklist valid"
        fi
    done

    if [ $errors -eq 0 ]; then
        echo ""
        echo "✅✅✅ ALL STATE CHECKLISTS VALIDATED PROJECT_DONEFULLY ✅✅✅"
        return 0
    else
        echo ""
        echo "❌❌❌ FOUND $errors CHECKLIST ISSUES ❌❌❌"
        return 1
    fi
}

validate_all_state_checklists
```

---

## 📜 THE CHECKLIST OATH

```
I, the Agent, swear by R510:

State checklists are MANDATORY execution requirements.
I WILL NOT skip any checklist item.
I WILL acknowledge EVERY item with proof.
I WILL complete ALL BLOCKING items before proceeding.
I WILL complete ALL EXIT REQUIREMENTS before transitioning.
I WILL NOT transition without a complete checklist.

My checklist is my execution contract.
Every item MUST be completed.
Every completion MUST be acknowledged.
Every acknowledgment MUST have proof.

This is SUPREME LAW.
Violation means FAILURE.
I WILL EXECUTE COMPLETELY.
```

---

**Remember:** A checklist without acknowledgments is worthless. An acknowledgment without proof is suspect. A transition without complete checklist is FAILURE. R510 makes execution explicit, trackable, and enforceable.

**See Also:**
- R232: TodoWrite Pending Items Override (integration point)
- R288: State File Update Requirements (EXIT REQUIREMENTS)
- R287: TODO Persistence (EXIT REQUIREMENTS)
- R322: Mandatory Checkpoints (EXIT REQUIREMENTS)
- R405: Automation Continuation Flag (EXIT REQUIREMENTS)

## State Manager Coordination (SF 3.0)

State Manager supports checklist execution through state transitions:
- **Each checklist item** may trigger state transitions
- **State validation** ensures checklist compliance (can't skip required states)
- **Atomic updates** track checklist progress in orchestrator-state-v3.json
- **Rollback protection** if checklist item fails

Checklist items that modify state use State Manager bookend pattern for safe updates.

See: R600 (checklist execution protocol), R601 (worklog maintenance)
