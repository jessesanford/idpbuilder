# Request for Software Factory Manager
## Design State Machine for [PROCESS NAME]

### Background
[Why this state machine is needed]
[What manual process it replaces]
[How many times process has been done manually]
[Current pain points and time investment]
[Expected benefits from automation]

### Request
Design a state machine with appropriate states and agent assignments to [PURPOSE/GOAL].

The state machine should [HIGH-LEVEL OBJECTIVE].

### Process Document
[Reference to detailed process documentation, e.g., `/tmp/PROCESS-NAME-MANUAL.md`]
[Any additional supporting documents]

### Phase Breakdown

#### Phase 1: [PHASE NAME]
**Purpose:** [What this phase accomplishes]
**Duration:** [Typical time in manual process]
**Key Activities:**
- [Activity 1]
- [Activity 2]
- [Activity 3]

**Success Criteria:**
- [ ] [Criterion 1]
- [ ] [Criterion 2]

**Common Issues:**
- [Issue and resolution]
- [Issue and resolution]

#### Phase 2: [PHASE NAME]
**Purpose:** [What this phase accomplishes]
**Duration:** [Typical time in manual process]
**Key Activities:**
- [Activity 1]
- [Activity 2]
- [Activity 3]

**Success Criteria:**
- [ ] [Criterion 1]
- [ ] [Criterion 2]

**Common Issues:**
- [Issue and resolution]
- [Issue and resolution]

#### Phase 3: [PHASE NAME]
[Continue pattern for all phases...]

### Key Requirements for State Machine Design

#### Critical Capabilities Needed:
1. **[CAPABILITY NAME]** - [Description of what's needed and why]
2. **[CAPABILITY NAME]** - [Description of what's needed and why]
3. **[CAPABILITY NAME]** - [Description of what's needed and why]
4. **[CAPABILITY NAME]** - [Description of what's needed and why]
5. **[CAPABILITY NAME]** - [Description of what's needed and why]

#### Quality Gates Required:
- **[GATE NAME]:** [What to validate and acceptance criteria]
- **[GATE NAME]:** [What to validate and acceptance criteria]
- **[GATE NAME]:** [What to validate and acceptance criteria]
- **[GATE NAME]:** [What to validate and acceptance criteria]

#### Error Recovery Needs:
- **[ERROR SCENARIO]:** [How to detect and recover]
- **[ERROR SCENARIO]:** [How to detect and recover]
- **[ERROR SCENARIO]:** [How to detect and recover]
- **[ERROR SCENARIO]:** [How to detect and recover]

### Agent Assignments

#### Orchestrator Role:
- [Responsibility 1]
- [Responsibility 2]
- [Responsibility 3]

#### Integration Agent Role:
- [Responsibility 1]
- [Responsibility 2]
- [Responsibility 3]

#### SW-Engineer Role:
- [Responsibility 1]
- [Responsibility 2]
- [Responsibility 3]

#### Code-Reviewer Role:
- [Responsibility 1]
- [Responsibility 2]
- [Responsibility 3]

#### Architect Role (if needed):
- [Responsibility 1]
- [Responsibility 2]
- [Responsibility 3]

### Specific Challenges to Address

#### Challenge 1: [NAME]
**Description:** [Detailed explanation of the challenge]
**Current Manual Solution:** [How you handle it manually]
**Required Automation:** [What the state machine must do]

#### Challenge 2: [NAME]
**Description:** [Detailed explanation of the challenge]
**Current Manual Solution:** [How you handle it manually]
**Required Automation:** [What the state machine must do]

#### Challenge 3: [NAME]
[Continue pattern for all challenges...]

### Success Criteria

The state machine implementation will be considered successful when it can:

- [ ] [Success criterion 1 - measurable outcome]
- [ ] [Success criterion 2 - measurable outcome]
- [ ] [Success criterion 3 - measurable outcome]
- [ ] [Success criterion 4 - measurable outcome]
- [ ] [Success criterion 5 - measurable outcome]

### Performance Requirements

- **Execution Time:** [Expected vs current manual time]
- **Parallelization:** [What can run in parallel]
- **Resource Usage:** [Any constraints]
- **Scalability:** [How many instances/branches/etc]

### Integration Requirements

- **Existing Tools:** [Tools that must be integrated]
- **Data Formats:** [Input/output formats needed]
- **APIs/Interfaces:** [External systems to interact with]
- **Reporting:** [Progress and status reporting needs]

### Example Execution

#### Successful Manual Execution Example:
[Provide a real example of successful manual execution, including:]
- Initial state/setup
- Commands executed in each phase
- Validation performed
- Issues encountered and resolved
- Final outcome

```bash
# Example commands from manual process
[Include actual commands used]
```

### Questions for Design Consideration

1. [Specific question about design choice]
2. [Specific question about error handling]
3. [Specific question about validation]
4. [Specific question about parallelization]
5. [Specific question about user interaction]

### Priority Order

Please prioritize design decisions in this order:
1. **[PRIORITY 1]** - [Why this is most important]
2. **[PRIORITY 2]** - [Why this is second]
3. **[PRIORITY 3]** - [Why this is third]
4. **[PRIORITY 4]** - [Why this is fourth]

### Expected Deliverables

Please provide:
1. **State Machine Definition** - Complete state machine with all states and transitions
2. **State Transition Rules** - Clear entry/exit conditions for each state
3. **Agent State Rules** - State-specific rules for each agent type
4. **Error States and Recovery** - Comprehensive error handling
5. **Validation Scripts** - Automated validation tools
6. **Example State File** - Sample JSON configuration
7. **Continuation Commands** - How to resume from any state
8. **Documentation Updates** - Integration with existing docs

### Additional Context

[Any additional information that might be helpful:]
- Related systems or processes
- Historical context
- Future plans that might affect design
- Constraints or limitations
- Dependencies on other work

### Timeline and Usage

**Expected Frequency:** [How often will this be used]
**First Use Target:** [When you plan to first use it]
**Critical Dependencies:** [What must be in place first]

---

**Note to Factory Manager:** Please feel free to ask clarifying questions before proceeding with the design. The goal is to create a robust, reliable state machine that can handle our [PROCESS NAME] process with minimal human intervention while maintaining safety and quality standards.