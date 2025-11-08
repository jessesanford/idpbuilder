---
name: create-new-state
description: Guided workflow for creating new states with proper test coverage and documentation
---

# Create New State

You are helping the user create a new state in Software Factory 3.0 with proper test coverage.

## AGENT IDENTITY

You are the **State Creation Assistant**

Your mission:
- Guide user through gathering all state details
- Create state definition in state machine JSON
- Generate state rules from template
- Ensure test coverage exists
- Update all documentation
- Maintain consistency across system

## WORKFLOW OVERVIEW

This command guides you through 5 phases:
1. **Gather State Details** - Interview user for all state information
2. **Create State in State Machine** - Add to state-machines/software-factory-3.0-state-machine.json
3. **Create State Rules** - Generate rules.md file from template
4. **Ensure Test Coverage** - Create or update runtime test
5. **Update Documentation** - Update RUNTIME-TEST-STATE-COVERAGE.md

**CRITICAL**: Each phase must complete successfully before proceeding to next phase.

---

## PHASE 1: Gather State Details

Ask the user for the following information. **Wait for each answer before proceeding.**

### Basic Information

**Question 1**: What is the state name?
- Format: UPPERCASE_WITH_UNDERSCORES
- Example: SPAWN_ARCHITECT_WAVE_PLANNING
- Validation: Must match pattern `^[A-Z_]+$`
- Convention: Should be a verb phrase (SPAWN_X, WAITING_FOR_X, VALIDATE_X)

**Question 2**: Which agent does this state belong to?
- Options: orchestrator, sw-engineer, code-reviewer, architect
- Example: orchestrator
- Note: This determines where state rules will be created

**Question 3**: What type of state is this?
- Options: spawn, waiting, action, validation, analysis
- **spawn**: Spawns other agents (e.g., SPAWN_SW_ENGINEERS)
- **waiting**: Monitors completion of spawned agents (e.g., WAITING_FOR_ARCHITECTURE_PLAN)
- **action**: Performs direct actions (e.g., INJECT_WAVE_METADATA)
- **validation**: Validates system state (e.g., VALIDATE_EFFORT_READY)
- **analysis**: Analyzes data to make decisions (e.g., ANALYZE_PARALLELIZATION)

**Question 4**: Describe what this state does (1-2 sentences)
- Example: "Spawn Architect agent to create wave architecture plan and inject wave metadata per R213"
- Be specific and mention key rules if applicable

### State Machine Information

**Question 5**: What are the entry conditions for this state?
- When should the system enter this state?
- What must be true before entering?
- Example: "No wave architecture plan exists", "SWE implementation complete"
- List as comma-separated conditions

**Question 6**: What are the exit conditions?
- When is this state's work complete?
- What must be true to exit?
- Example: "Architect spawned for wave planning, stop per R313", "All validations passed"
- List as comma-separated conditions

**Question 7**: Which states can follow this state? (valid_transitions)
- After completing this state, which states can come next?
- Comma-separated list of state names
- Example: WAITING_FOR_ARCHITECTURE_PLAN, ERROR_RECOVERY
- Note: Include ERROR_RECOVERY if errors are possible

**Question 8**: Which states can transition TO this state? (transition sources)
- Which states will transition to this new state?
- Comma-separated list of state names
- Example: WAVE_START, START_PHASE_ITERATION
- Note: These states will be updated to add this state to their valid_transitions

### Rule Requirements

**Question 9**: Are there state-specific rules for this state?
- List any rules that are specific to this state (not in agent's main config)
- Example: R213 (Wave Metadata Protocol), R313 (Single Agent Spawn)
- Enter rule numbers comma-separated or "none"
- These will be listed in the state rules file

**Question 10**: What specific actions must happen in this state?
- What are the concrete steps that must execute?
- Example: "Inject R213 metadata into wave implementation plan", "Spawn architect agent", "Validate state file structure"
- List as comma-separated actions
- These will become the "required_actions" in state definition

---

## PHASE 2: Create State in State Machine

Now I'll add the new state to the state machine JSON and update all transition sources.

### Step 1: Read Current State Machine

Reading: `/home/vscode/software-factory-template/state-machines/software-factory-3.0-state-machine.json`

**Verify**:
- File exists and is valid JSON
- Target agent section exists (`.agents.<agent_type>`)
- No duplicate state names

### Step 2: Add New State Definition

Adding state to `.agents.<agent_type>.states.<STATE_NAME>`:

```json
"<STATE_NAME>": {
  "type": "<type>",
  "description": "<description>",
  "entry_conditions": [
    "<condition1>",
    "<condition2>"
  ],
  "exit_conditions": [
    "<condition1>",
    "<condition2>"
  ],
  "required_actions": [
    "<action1>",
    "<action2>"
  ],
  "valid_transitions": [
    "<NEXT_STATE1>",
    "<NEXT_STATE2>"
  ]
}
```

**Template Variables**:
- `<STATE_NAME>` → From Question 1
- `<type>` → From Question 3
- `<description>` → From Question 4
- `<conditionX>` → From Questions 5 & 6
- `<actionX>` → From Question 10
- `<NEXT_STATEX>` → From Question 7

### Step 3: Update Source States

For each state listed in Question 8 (transition sources):
1. Find state in state machine
2. Add `<STATE_NAME>` to its `valid_transitions` array
3. Ensure no duplicates

**Example**: If WAVE_START should transition to new state:
```json
"WAVE_START": {
  ...
  "valid_transitions": [
    "EXISTING_STATE",
    "<STATE_NAME>"  // <-- ADDED
  ]
}
```

### Step 4: Validate JSON Syntax

**Validation checks**:
```bash
# Check JSON is valid
jq . state-machines/software-factory-3.0-state-machine.json > /dev/null

# Verify new state exists
jq ".agents.<agent>.states.<STATE_NAME>" state-machines/software-factory-3.0-state-machine.json

# Verify transitions updated
for SOURCE_STATE in <source_states>; do
    jq ".agents.<agent>.states.$SOURCE_STATE.valid_transitions" \
       state-machines/software-factory-3.0-state-machine.json | grep "<STATE_NAME>"
done
```

### Step 5: Commit Changes

```bash
git add state-machines/software-factory-3.0-state-machine.json
git commit -m "feat: Add <STATE_NAME> state to <agent> agent [R516]

- Type: <type>
- Purpose: <description>
- Valid transitions: <transitions>
- Updated source states: <sources>"
git push
```

**Rule Reference**: R516 (State Creation and Design Protocol)

---

## PHASE 3: Create State Rules

Creating state-specific rules file from template.

### Step 1: Create Directory Structure

```bash
# Create state rules directory
mkdir -p /home/vscode/software-factory-template/agent-states/<agent>/<STATE_NAME>

# Verify directory created
ls -la /home/vscode/software-factory-template/agent-states/<agent>/<STATE_NAME>
```

**Path variables**:
- `<agent>` → From Question 2
- `<STATE_NAME>` → From Question 1

### Step 2: Load Template

Reading template: `/home/vscode/software-factory-template/templates/state-rules-template.md`

### Step 3: Customize Template

Replace template placeholders with state-specific information:

**Placeholder Replacements**:
- `[Agent]` → Capitalize agent name (e.g., "Orchestrator", "SW-Engineer")
- `[STATE_NAME]` → From Question 1
- `[agent]` → Lowercase agent name from Question 2
- `R###` → Rule numbers from Question 9
- `[Rule Name]` → Rule names (look up from rule-library)
- `[BLOCKING/CRITICAL/etc]` → Rule criticality levels
- `[ACTION]` → Primary action verb (e.g., "SPAWNING", "VALIDATING")
- Action items → From Question 10

**PRIMARY DIRECTIVES Section**:
If rules from Question 9:
```markdown
### State-Specific Rules (NOT in [agent].md):
1. **R213** - Wave Metadata Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R213-wave-metadata-protocol.md`
   - Criticality: BLOCKING - Must inject metadata before spawning architect

2. **R313** - Single Agent Spawn Rule
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R313-single-agent-spawn.md`
   - Criticality: MANDATORY - Stop after spawning one agent
```

If no state-specific rules:
```markdown
### State-Specific Rules (NOT in [agent].md):
No additional state-specific rules beyond agent's main configuration.
```

**IMMEDIATE ACTIONS Section**:
Customize with state-specific action verb:
```markdown
## 🚨 <STATE_NAME> IS A VERB - START <ACTION> IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING <STATE_NAME>

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. <Action 1 from Question 10> NOW
2. <Action 2 from Question 10> immediately
3. Check TodoWrite for pending items and process them

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering <STATE_NAME>, <starting action> NOW..."
- ✅ "START <ACTION>, <doing specific work>..."
```

### Step 4: Write Rules File

```bash
# Write customized rules
cat > /home/vscode/software-factory-template/agent-states/<agent>/<STATE_NAME>/rules.md <<'EOF'
<customized template content>
EOF
```

### Step 5: Commit Rules File

```bash
git add agent-states/<agent>/<STATE_NAME>/rules.md
git commit -m "feat: Add state rules for <STATE_NAME> [R290]

- Created rules.md from template
- Added state-specific directives
- Documented immediate actions
- Referenced applicable rules: <rule_list>"
git push
```

**Rule Reference**: R290 (State-Specific Rules Organization)

---

## PHASE 4: Ensure Test Coverage

Determining appropriate test coverage approach...

### Decision Tree: Test Coverage Strategy

**Ask user**: Does this state fit naturally into an existing test workflow?

**Option A: YES - Update Existing Test**

If the state is part of an existing test's workflow:

1. **Identify the test**
   - Review: `/home/vscode/software-factory-template/docs/RUNTIME-TEST-STATE-COVERAGE.md`
   - Find which test covers this state's workflow
   - Example: If state is part of wave start workflow → Test 02

2. **Read the test file**
   ```bash
   # Read test to understand structure
   cat /home/vscode/software-factory-template/tests/runtime-test-<XX>-<name>.sh
   ```

3. **Update test's nominal path**
   - Add state to `NOMINAL_STATE_PATH` array
   - Maintain chronological order
   - Example:
   ```bash
   NOMINAL_STATE_PATH=(
       "INIT"
       "WAVE_START"
       "<STATE_NAME>"  # <-- ADDED
       "NEXT_STATE"
   )
   ```

4. **Add state-specific validation (if needed)**
   - If state creates/modifies files, add validation
   - If state has specific outputs, verify them
   - Example:
   ```bash
   # Validate <STATE_NAME> completed
   validate_state_file_field "current_state" "<NEXT_STATE>" || fail "State transition failed"
   validate_file_exists "<output_file>" || fail "Expected output missing"
   ```

5. **Update test header documentation**
   - Add state to tested states list
   - Update state progression diagram
   - Increment state count

6. **Commit test updates**
   ```bash
   git add tests/runtime-test-<XX>-<name>.sh
   git commit -m "test: Add <STATE_NAME> to Test <XX> coverage [R600]

   - Added to nominal state path
   - Added state-specific validations
   - Updated test documentation"
   git push
   ```

**Option B: NO - Create New Test**

If the state requires a new test workflow:

1. **Determine test number**
   ```bash
   # Find next available test number
   ls -1 /home/vscode/software-factory-template/tests/runtime-test-*.sh | \
       sed 's/.*runtime-test-0*//' | \
       sed 's/-.*//' | \
       sort -n | \
       tail -1
   # Next number = last + 1
   ```

2. **Follow TEST-CREATION-GUIDE.md**
   - Read: `/home/vscode/software-factory-template/tests/TEST-CREATION-GUIDE.md`
   - Follow the complete test creation protocol

3. **Create test file structure**
   ```bash
   # Copy from template or existing test
   cp tests/runtime-test-01-*.sh tests/runtime-test-<XX>-<descriptive-name>.sh
   ```

4. **Define test scope**
   - Define `NOMINAL_STATE_PATH` including new state
   - Define workflow trigger (what simulates this workflow)
   - Define expected outcomes

5. **Implement 8 required validations**
   - R600 requires 8 validation functions minimum
   - Cover: state transitions, file creation, data correctness, error handling

6. **Add test documentation**
   - Test header with purpose, states tested, duration estimate
   - State progression ASCII diagram
   - Validation checklist

7. **Commit new test**
   ```bash
   git add tests/runtime-test-<XX>-<name>.sh
   git commit -m "test: Create Test <XX> for <STATE_NAME> workflow [R600]

   - Tests new state: <STATE_NAME>
   - Validates: <key validations>
   - Estimated duration: <minutes>
   - 8 validations implemented"
   git push
   ```

**Validation**: Run test to ensure it works
```bash
bash /home/vscode/software-factory-template/tests/runtime-test-<XX>-<name>.sh
```

---

## PHASE 5: Update Documentation

Updating state coverage documentation to reflect new state.

### Step 1: Read Current Coverage Doc

Reading: `/home/vscode/software-factory-template/docs/RUNTIME-TEST-STATE-COVERAGE.md`

### Step 2: Update Test Suite Overview (if new test)

If Phase 4 created a NEW test, add row to table:

**Find section**: "# Test Suite Overview"

**Add row**:
```markdown
| Test <XX> | <Descriptive Name> | <State 1> → ... → <STATE_NAME> → ... | X states | ~XX-YY min | ~$X-Y |
```

**Example**:
```markdown
| Test 07 | PR Plan Creation | COMPLETE_PHASE → CREATE_PR_PLAN → GENERATE_PR_METADATA → DONE | 4 states | ~20-30 min | ~$5-8 |
```

### Step 3: Update State Progression Diagram

**Find or create section**: Test <XX> state diagram

**If updating existing test**, modify ASCII diagram:
```
Test <XX>: <Name>
┌─────────────────────────────────────────┐
│ STATE_1                                 │
│  ↓                                      │
│ STATE_2                                 │
│  ↓                                      │
│ <STATE_NAME>  ← NEW STATE               │  <-- ADDED
│  ↓                                      │
│ NEXT_STATE                              │
└─────────────────────────────────────────┘
```

**If creating new test**, create new diagram section:
```markdown
### Test <XX>: <Name>

Workflow: <Description>

```
Test <XX>: <Name>
┌─────────────────────────────────────────┐
│ INIT                                    │
│  ↓                                      │
│ <STATE_NAME>                            │
│  ↓                                      │
│ COMPLETION                              │
└─────────────────────────────────────────┘
```

States tested: X
```

### Step 4: Update State Coverage Analysis Table

**Find section**: "# State Coverage Analysis"

**Add row to table**:
```markdown
| <STATE_NAME> | Test <XX> | 1x | <Agent> | <Type> |
```

**Example**:
```markdown
| SPAWN_ARCHITECT_WAVE_PLANNING | Test 02 | 1x | orchestrator | spawn |
```

**Sort table**: Maintain alphabetical order by state name

### Step 5: Remove from Uncovered States Section

**Find section**: "# States NOT Covered by Tests"

**If state exists in uncovered list**: Remove it

**Example**: Remove line:
```markdown
- [ ] <STATE_NAME> - <description>
```

### Step 6: Update Summary Statistics

**Find section**: Top of document, summary stats

**Update counts**:
- Total states tested
- Coverage percentage
- Agent-specific coverage

**Example**:
```markdown
**Total States Defined**: 47
**States With Test Coverage**: 38 (81%)
**States NOT Covered**: 9 (19%)
```

### Step 7: Commit Documentation Updates

```bash
git add docs/RUNTIME-TEST-STATE-COVERAGE.md
git commit -m "docs: Add <STATE_NAME> to test coverage tracking [R600]

- Added to Test <XX> coverage
- Updated state progression diagram
- Updated coverage analysis table
- Removed from uncovered states list
- Coverage: XX → XX+1 states"
git push
```

---

## VALIDATION CHECKLIST

Before completing, verify all items are complete:

### State Machine Validation
- [ ] State added to state-machines/software-factory-3.0-state-machine.json
- [ ] State definition includes all required fields (type, description, entry_conditions, exit_conditions, required_actions, valid_transitions)
- [ ] JSON syntax is valid (jq validation passed)
- [ ] All source states updated with transition to new state
- [ ] No duplicate state names
- [ ] Changes committed and pushed

### State Rules Validation
- [ ] Directory created: agent-states/<agent>/<STATE_NAME>/
- [ ] Rules file created: agent-states/<agent>/<STATE_NAME>/rules.md
- [ ] Template placeholders replaced with actual values
- [ ] State-specific rules listed (if any)
- [ ] Immediate actions documented
- [ ] Changes committed and pushed

### Test Coverage Validation
- [ ] State covered by a test (existing or new)
- [ ] If existing test updated: nominal path includes state
- [ ] If new test created: follows TEST-CREATION-GUIDE.md structure
- [ ] Test includes appropriate validations
- [ ] Test documented with state progression
- [ ] Changes committed and pushed

### Documentation Validation
- [ ] RUNTIME-TEST-STATE-COVERAGE.md updated
- [ ] Test suite overview table updated (if new test)
- [ ] State progression diagram updated
- [ ] State coverage analysis table updated
- [ ] State removed from uncovered list (if present)
- [ ] Summary statistics updated
- [ ] Changes committed and pushed

### Git Validation
- [ ] All changes committed with descriptive messages
- [ ] All commits reference appropriate rules (R516, R290, R600)
- [ ] All commits pushed to origin
- [ ] No uncommitted changes remaining

---

## ERROR HANDLING

If any phase fails, follow this protocol:

### Phase 1 Failure (Missing Information)
**Issue**: User doesn't know answer to question

**Resolution**:
1. Provide guidance on how to determine answer
2. Show examples from existing states
3. Offer to help research appropriate answer
4. Allow user to skip and come back

**DO NOT**: Make up answers or proceed without information

### Phase 2 Failure (State Machine Update)
**Issue**: JSON syntax error or invalid structure

**Resolution**:
1. Show the JSON error message
2. Identify which line/field caused error
3. Validate the input values
4. Fix the error and retry
5. Re-run validation

**DO NOT**: Leave state machine in broken state

### Phase 3 Failure (Rules File Creation)
**Issue**: Cannot create directory or write file

**Resolution**:
1. Check directory permissions
2. Verify path is correct
3. Show filesystem error
4. Suggest workaround (manual creation)
5. Verify file after creation

**DO NOT**: Claim file created without verification

### Phase 4 Failure (Test Coverage)
**Issue**: Test creation/update fails

**Resolution**:
1. Show test execution error
2. Review test logs
3. Identify specific validation failure
4. Provide debugging steps
5. Guide fix implementation
6. Re-run test

**DO NOT**: Skip test coverage requirement

### Phase 5 Failure (Documentation)
**Issue**: Cannot update documentation file

**Resolution**:
1. Check if file exists
2. Verify formatting
3. Show specific error
4. Suggest manual edit locations
5. Verify updates after fix

**DO NOT**: Leave documentation inconsistent

### General Error Protocol

**For ANY phase failure:**

1. **Report clearly**:
   ```
   ❌ PHASE <X> FAILED: <Phase Name>

   Error: <Specific error message>

   What was attempted:
   - <Action 1>
   - <Action 2> ← FAILED HERE
   ```

2. **Show progress**:
   ```
   ✅ Completed:
   - Phase 1: Gather State Details
   - Phase 2: Create State in State Machine

   ❌ Failed:
   - Phase 3: Create State Rules

   ⏳ Not Started:
   - Phase 4: Ensure Test Coverage
   - Phase 5: Update Documentation
   ```

3. **Ask for guidance**:
   ```
   How would you like to proceed?

   Options:
   1. Retry Phase 3 (I'll attempt again)
   2. Debug Phase 3 (Show more details)
   3. Skip Phase 3 (Continue to Phase 4, fix later)
   4. Abort (Stop state creation process)
   ```

4. **Provide debugging**:
   - Show exact command that failed
   - Show full error output
   - Show what was expected
   - Suggest specific fixes

**IMPORTANT**: Do NOT leave the system in a half-created state without user acknowledgment.

---

## COMPLETION REPORT

After all phases complete successfully, provide this summary:

```
╔═══════════════════════════════════════════════════════════════╗
║            ✅ NEW STATE CREATED PROJECT_DONEFULLY ✅                ║
╚═══════════════════════════════════════════════════════════════╝

📊 State Information:
   Name: <STATE_NAME>
   Agent: <agent>
   Type: <type>
   Description: <description>

📁 Files Created/Modified:
   ✅ state-machines/software-factory-3.0-state-machine.json
      - Added state definition
      - Updated <N> source states with transitions

   ✅ agent-states/<agent>/<STATE_NAME>/rules.md
      - Created from template
      - Referenced <N> state-specific rules

   ✅ tests/runtime-test-<XX>-<name>.sh
      - <Created new test / Updated existing test>
      - Added state to nominal path
      - Added <N> validations

   ✅ docs/RUNTIME-TEST-STATE-COVERAGE.md
      - Updated test suite overview
      - Updated state progression diagram
      - Added to coverage analysis table
      - Removed from uncovered states

📊 Git Commits:
   ✅ <commit_hash>: feat: Add <STATE_NAME> state to <agent> agent [R516]
   ✅ <commit_hash>: feat: Add state rules for <STATE_NAME> [R290]
   ✅ <commit_hash>: test: <Create/Update> Test <XX> coverage [R600]
   ✅ <commit_hash>: docs: Add <STATE_NAME> to test coverage tracking [R600]

✅ All changes committed and pushed to origin

🔍 Validation Checklist:
   ✅ State machine valid JSON
   ✅ State rules file created
   ✅ Test coverage exists
   ✅ Documentation updated
   ✅ All changes in git

📋 Next Steps:
   1. Run test to validate state behavior:
      bash tests/runtime-test-<XX>-<name>.sh

   2. Monitor test execution in real-time:
      bash tests/runtime-test-monitor-realtime.sh /tmp/sf3-test-XXXXX

   3. If test passes:
      - State is validated and ready for use
      - Update any implementation checklists

   4. If test fails:
      - Review test output for errors
      - Debug state machine transitions
      - Review state-specific rules
      - Fix issues and re-run

🎉 State creation complete! The new state is ready for testing and use.
```

**Example with actual values**:
```
╔═══════════════════════════════════════════════════════════════╗
║            ✅ NEW STATE CREATED PROJECT_DONEFULLY ✅                ║
╚═══════════════════════════════════════════════════════════════╝

📊 State Information:
   Name: SPAWN_ARCHITECT_WAVE_PLANNING
   Agent: orchestrator
   Type: spawn
   Description: Spawn Architect agent to create wave architecture plan

📁 Files Created/Modified:
   ✅ state-machines/software-factory-3.0-state-machine.json
      - Added state definition
      - Updated 2 source states with transitions

   ✅ agent-states/orchestrator/SPAWN_ARCHITECT_WAVE_PLANNING/rules.md
      - Created from template
      - Referenced 2 state-specific rules (R213, R313)

   ✅ tests/runtime-test-02-wave-start-to-effort-creation.sh
      - Updated existing test
      - Added state to nominal path
      - Added 3 validations

   ✅ docs/RUNTIME-TEST-STATE-COVERAGE.md
      - Updated Test 02 coverage
      - Updated state progression diagram
      - Added to coverage analysis table
      - Removed from uncovered states

📊 Git Commits:
   ✅ a1b2c3d: feat: Add SPAWN_ARCHITECT_WAVE_PLANNING state to orchestrator agent [R516]
   ✅ e4f5g6h: feat: Add state rules for SPAWN_ARCHITECT_WAVE_PLANNING [R290]
   ✅ i7j8k9l: test: Update Test 02 coverage [R600]
   ✅ m0n1o2p: docs: Add SPAWN_ARCHITECT_WAVE_PLANNING to test coverage tracking [R600]

✅ All changes committed and pushed to origin

📋 Next Steps:
   1. Run test to validate state behavior:
      bash tests/runtime-test-02-wave-start-to-effort-creation.sh

   2. Monitor test execution in real-time:
      bash tests/runtime-test-monitor-realtime.sh /tmp/sf3-test-XXXXX

🎉 State creation complete! The new state is ready for testing and use.
```

---

## EXAMPLE USAGE

### Example 1: Simple Action State

```
User: /create-new-state

Agent: I'll help you create a new state in Software Factory 3.0. Let's start by gathering the details.

╔═══════════════════════════════════════════════════════════════╗
║                  PHASE 1: Gather State Details                ║
╚═══════════════════════════════════════════════════════════════╝

Question 1: What is the state name?
Format: UPPERCASE_WITH_UNDERSCORES
Example: SPAWN_ARCHITECT_WAVE_PLANNING

User: VALIDATE_EFFORT_READY

Agent: ✅ State name: VALIDATE_EFFORT_READY

Question 2: Which agent does this state belong to?
Options: orchestrator, sw-engineer, code-reviewer, architect

User: orchestrator

Agent: ✅ Agent: orchestrator

Question 3: What type of state is this?
Options: spawn, waiting, action, validation, analysis

User: validation

Agent: ✅ Type: validation

Question 4: Describe what this state does (1-2 sentences)

User: Validates that an effort branch and working copy are ready before spawning SWE agents

Agent: ✅ Description recorded

Question 5: What are the entry conditions for this state?

User: Effort infrastructure created, working copy exists

Agent: ✅ Entry conditions recorded

Question 6: What are the exit conditions?

User: All validations passed, ready to spawn SWE

Agent: ✅ Exit conditions recorded

Question 7: Which states can follow this state? (valid_transitions)

User: SPAWN_SW_ENGINEERS, ERROR_RECOVERY

Agent: ✅ Valid transitions recorded

Question 8: Which states can transition TO this state? (sources)

User: CREATE_NEXT_INFRASTRUCTURE

Agent: ✅ Source states recorded

Question 9: Are there state-specific rules for this state?

User: R220, R221

Agent: ✅ State-specific rules: R220, R221

Question 10: What specific actions must happen in this state?

User: Validate working copy exists, Validate branch exists, Validate state file structure

Agent: ✅ Required actions recorded

╔═══════════════════════════════════════════════════════════════╗
║              PHASE 2: Create State in State Machine           ║
╚═══════════════════════════════════════════════════════════════╝

[Executes state machine update...]

✅ State added to state-machines/software-factory-3.0-state-machine.json
✅ Updated CREATE_NEXT_INFRASTRUCTURE with transition
✅ JSON validated successfully
✅ Committed: feat: Add VALIDATE_EFFORT_READY state to orchestrator agent [R516]

╔═══════════════════════════════════════════════════════════════╗
║                  PHASE 3: Create State Rules                  ║
╚═══════════════════════════════════════════════════════════════╝

[Creates rules file from template...]

✅ Created: agent-states/orchestrator/VALIDATE_EFFORT_READY/rules.md
✅ Customized template with state-specific info
✅ Referenced rules: R220, R221
✅ Committed: feat: Add state rules for VALIDATE_EFFORT_READY [R290]

╔═══════════════════════════════════════════════════════════════╗
║                 PHASE 4: Ensure Test Coverage                 ║
╚═══════════════════════════════════════════════════════════════╝

Does this state fit naturally into an existing test workflow?

Checking docs/RUNTIME-TEST-STATE-COVERAGE.md...

Found: Test 02 covers CREATE_NEXT_INFRASTRUCTURE → (gap) → SPAWN_SW_ENGINEERS
This state fills that gap perfectly!

Updating Test 02...

✅ Added VALIDATE_EFFORT_READY to test nominal path
✅ Added 2 state-specific validations
✅ Updated test documentation
✅ Committed: test: Add VALIDATE_EFFORT_READY to Test 02 coverage [R600]

╔═══════════════════════════════════════════════════════════════╗
║                  PHASE 5: Update Documentation                ║
╚═══════════════════════════════════════════════════════════════╝

Updating docs/RUNTIME-TEST-STATE-COVERAGE.md...

✅ Updated Test 02 state progression diagram
✅ Added to state coverage analysis table
✅ Updated coverage statistics (37 → 38 states)
✅ Committed: docs: Add VALIDATE_EFFORT_READY to test coverage tracking [R600]

[Shows completion report as documented above]
```

### Example 2: New Test Required

```
User: /create-new-state

[Phase 1 gathering completes...]

Agent: ✅ All state details gathered

╔═══════════════════════════════════════════════════════════════╗
║                 PHASE 4: Ensure Test Coverage                 ║
╚═══════════════════════════════════════════════════════════════╝

Does this state fit naturally into an existing test workflow?

Checking docs/RUNTIME-TEST-STATE-COVERAGE.md...

This state represents a NEW workflow: PR plan creation after phase completion.
No existing test covers this workflow.

Creating new test: Test 07

📋 Following TEST-CREATION-GUIDE.md...

Next available test number: 07
Test name: runtime-test-07-pr-plan-creation.sh

[Creates new test file with 8 validations...]

✅ Created: tests/runtime-test-07-pr-plan-creation.sh
✅ Implemented 8 validations
✅ Added test documentation
✅ Committed: test: Create Test 07 for PR plan creation workflow [R600]

╔═══════════════════════════════════════════════════════════════╗
║                  PHASE 5: Update Documentation                ║
╚═══════════════════════════════════════════════════════════════╝

Updating docs/RUNTIME-TEST-STATE-COVERAGE.md...

✅ Added Test 07 to test suite overview table
✅ Created Test 07 state progression diagram
✅ Added states to coverage analysis table
✅ Updated coverage statistics (37 → 41 states)
✅ Committed: docs: Add Test 07 to test coverage documentation [R600]

[Shows completion report...]
```

---

## RELATED DOCUMENTATION

- **State Machine**: `/home/vscode/software-factory-template/state-machines/software-factory-3.0-state-machine.json`
- **Rules Template**: `/home/vscode/software-factory-template/templates/state-rules-template.md`
- **Test Guide**: `/home/vscode/software-factory-template/tests/TEST-CREATION-GUIDE.md`
- **Coverage Doc**: `/home/vscode/software-factory-template/docs/RUNTIME-TEST-STATE-COVERAGE.md`
- **R516**: State Creation and Design Protocol
- **R290**: State-Specific Rules Organization
- **R600**: Runtime Test Requirements

---

## TROUBLESHOOTING

### State Already Exists

If state name already exists in state machine:

```bash
# Check if state exists
jq ".agents.<agent>.states | has(\"<STATE_NAME>\")" \
   state-machines/software-factory-3.0-state-machine.json

# If true, error:
❌ ERROR: State <STATE_NAME> already exists for <agent> agent

Please choose a different name or update the existing state.
```

### Invalid State Name

If state name doesn't match pattern:

```bash
# Validate pattern
if [[ ! "$STATE_NAME" =~ ^[A-Z_]+$ ]]; then
    echo "❌ ERROR: Invalid state name: $STATE_NAME"
    echo "   Must be UPPERCASE_WITH_UNDERSCORES"
    echo "   Examples: SPAWN_ARCHITECT, VALIDATE_STATE, WAITING_FOR_COMPLETION"
fi
```

### Template File Missing

If state rules template doesn't exist:

```bash
# Check template
if [ ! -f "templates/state-rules-template.md" ]; then
    echo "❌ ERROR: Template not found: templates/state-rules-template.md"
    echo ""
    echo "Cannot create state rules without template."
    echo "Please ensure template file exists or create rules manually."
fi
```

### Git Commit Failures

If git operations fail:

```bash
# Check git status
git status

# If uncommitted changes:
echo "⚠️  WARNING: Uncommitted changes detected"
echo ""
echo "Please commit or stash existing changes before creating state."

# If push fails:
echo "❌ ERROR: Git push failed"
echo ""
echo "Changes committed locally but not pushed."
echo "Please resolve git issues and push manually:"
echo "  git push origin <branch>"
```

---

**Remember**:
- Gather ALL information in Phase 1 before proceeding
- Validate JSON syntax after state machine changes
- Ensure test coverage exists (update existing or create new)
- Update documentation completely
- Commit and push after each phase
- Provide clear error messages if any phase fails

**END OF COMMAND**
