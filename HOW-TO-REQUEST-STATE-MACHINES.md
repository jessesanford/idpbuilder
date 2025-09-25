# HOW TO REQUEST STATE MACHINES FROM SOFTWARE FACTORY MANAGER

## Section 1: Overview

### Purpose of This Guide
This guide provides a structured approach for requesting new state machines from the Software Factory Manager agent. Following this process ensures you receive comprehensive state machine designs with all necessary supporting files, validation scripts, and agent-specific rules.

### When to Create New State Machines
Consider creating a new state machine when:
- You have a complex multi-phase process that needs automation
- A manual process has been repeated successfully 3+ times
- Multiple agents need to coordinate on a structured workflow
- You need consistent error recovery and validation gates
- The process has clear phases with defined success criteria

### Role of the Software-Factory-Manager
The Software Factory Manager is responsible for:
- Translating your process documentation into formal state machines
- Assigning appropriate agents to each state
- Creating validation gates and error recovery states
- Generating all supporting files (rules, scripts, commands)
- Ensuring consistency with existing Software Factory patterns

## Section 2: Pre-Request Preparation

### 1. Document Your Manual Process First
**CRITICAL:** Before requesting a state machine, you MUST have experience with the manual process:
- Run through the process manually at least once
- Document each step and decision point
- Note common errors and how to resolve them
- Track time spent in each phase
- Identify which parts can be parallelized

### 2. Identify the Phases/Stages
Break down your process into logical phases:
- Each phase should have a clear purpose
- Phases should have defined entry/exit criteria
- Consider natural breakpoints for validation
- Group related activities together
- Think about rollback boundaries

### 3. Note Which Agents Should Do What
Map activities to appropriate agents:
- **Orchestrator:** Coordination, state management, agent spawning
- **Integration:** Branch operations, merging, PR creation
- **SW-Engineer:** Code modifications, file operations
- **Code-Reviewer:** Validation, compliance checking
- **Architect:** Design review, impact assessment

### 4. Collect Error Cases and Recovery Needs
Document failure scenarios:
- What errors have you encountered?
- How did you recover from them?
- Which errors are recoverable vs fatal?
- What validation would prevent errors?
- What rollback capabilities are needed?

### 5. Run Through the Process Manually
**Essential validation steps:**
- Execute the complete process end-to-end
- Document the exact commands used
- Note decision points and criteria
- Measure actual vs expected outcomes
- Identify automation opportunities

## Section 3: Required Information Checklist

Before creating your request, ensure you have:

- [ ] **Process name and purpose** - Clear, descriptive name and goal
- [ ] **Detailed phase-by-phase breakdown** - Each phase with steps
- [ ] **Agent assignments for each phase** - Who does what
- [ ] **Success criteria for each phase** - How to know it worked
- [ ] **Error conditions to handle** - What can go wrong
- [ ] **Validation gates needed** - Quality checkpoints
- [ ] **Any special tools/scripts required** - External dependencies
- [ ] **Example of successful manual execution** - Proof it works

## Section 4: Creating the Request Document

### Template Structure
Use the proven structure from successful requests:

```markdown
# Request for Software Factory Manager
## Design State Machine for [PROCESS NAME]

### Background
[Context: Why this is needed, how many times done manually, impact]

### Request
[Clear statement of what state machine should accomplish]

### Process Document
[Reference to detailed process documentation if separate]

### Phase Breakdown
1. **Phase Name** - Description
   - Key activities
   - Success criteria
   - Common issues

### Key Requirements for State Machine Design

#### Critical Capabilities Needed:
[List essential functions the state machine must support]

#### Quality Gates Required:
[List validation points and criteria]

#### Error Recovery Needs:
[List failure scenarios and recovery strategies]

### Specific Challenges to Address
[Unique complexities or edge cases]

### Expected Deliverables
[What files and artifacts you expect]

### Success Criteria
[How to measure if state machine works correctly]
```

### Key Elements to Include

#### Background Section
- Number of times process executed manually
- Time/effort currently required
- Problems the automation will solve
- Business value of automation

#### Process Documentation
- Step-by-step procedures
- Decision trees
- Command sequences
- File locations and formats

#### Requirements Section
- Must-have capabilities
- Nice-to-have features
- Performance expectations
- Integration needs

#### Success Criteria
- Measurable outcomes
- Quality metrics
- Time savings
- Error reduction

### Be Specific About:
- **Agent Responsibilities:** Exactly what each agent should do
- **Validation Points:** Where and what to check
- **Error Handling:** How to detect and recover
- **Parallelization:** What can run simultaneously
- **Dependencies:** What must happen sequentially

### Clearly State What Should NOT Happen
- Forbidden operations
- Protected resources
- Unacceptable side effects
- Breaking changes to avoid

## Section 5: Submitting to Software-Factory-Manager

### Create Request Document
1. Save your request in `/tmp/` directory:
   ```bash
   /tmp/[PROCESS-NAME]-STATE-MACHINE-REQUEST.md
   ```

2. Include all supporting documentation:
   ```bash
   /tmp/[PROCESS-NAME]-MANUAL-PROCESS.md
   /tmp/[PROCESS-NAME]-EXAMPLES.md
   ```

### Use Specific Prompt Structure

```markdown
@software-factory-manager

I need you to design a state machine for [PROCESS NAME].

**Request Document:** `/tmp/[PROCESS-NAME]-STATE-MACHINE-REQUEST.md`
**Process Documentation:** `/tmp/[PROCESS-NAME]-MANUAL-PROCESS.md`

Please create:
1. Complete state machine definition
2. State-specific rules for each agent
3. Validation scripts and tools
4. Example state file
5. Continuation commands

The state machine should handle [KEY REQUIREMENTS].

Priority: [Safety/Speed/Completeness]
```

### Mention Expected Deliverables Explicitly
Always specify what you expect to receive:
- State machine markdown file
- Agent state rules
- Validation scripts
- Example configurations
- Documentation updates

### Reference the Request Document Path
Always provide absolute paths to:
- Your request document
- Supporting documentation
- Example files
- Current manual process docs

## Section 6: Expected Deliverables

### Primary Deliverables

#### 1. State Machine File
- Location: Root directory or `/state-machines/`
- Format: `[PROCESS-NAME]-STATE-MACHINE.md`
- Contents: Complete state definitions, transitions, validations

#### 2. State-Specific Rules
- Location: `/agent-states/[agent-name]/[STATE]/rules.md`
- Format: Structured rules with criticality levels
- Contents: Agent behavior for each state

#### 3. Validation Scripts/Tools
- Location: `/tools/` or `/utilities/`
- Format: Bash scripts with clear documentation
- Contents: Automated validation and checks

#### 4. Continuation Commands
- Location: `/.claude/commands/`
- Format: Markdown with command structure
- Contents: Resume/continue instructions

#### 5. Example State Files
- Location: Root or `/examples/`
- Format: `.json.example` files
- Contents: Sample configurations

### Supporting Deliverables
- Updated documentation indexes
- Integration with existing systems
- Test cases and scenarios
- Recovery procedures

## Section 7: Post-Creation Verification

### Review the State Machine
Verify completeness:
- [ ] All phases from your process are represented
- [ ] Every state has clear entry/exit conditions
- [ ] Error states and recovery paths exist
- [ ] Validation gates are properly placed
- [ ] Agent assignments match capabilities

### Check All States Have Rules
For each state, verify:
- [ ] Agent-specific rules exist
- [ ] Rules have proper criticality levels
- [ ] Success criteria are defined
- [ ] Error conditions are handled
- [ ] Transitions are documented

### Verify Agent Assignments
Ensure correct agent usage:
- [ ] Orchestrator only coordinates (no implementation)
- [ ] Integration handles branch/PR operations
- [ ] SW-Engineers do code modifications
- [ ] Code-Reviewers validate and check
- [ ] No agent exceeds its capabilities

### Test with Small Example
Before full deployment:
1. Create minimal test case
2. Run through state machine manually
3. Verify each state transition
4. Test error recovery paths
5. Validate all outputs

### Validation Checklist
- [ ] State machine loads without errors
- [ ] All referenced files exist
- [ ] Scripts have proper permissions
- [ ] Example state file is valid JSON
- [ ] Documentation is complete

## Section 8: Example Workflow

### Successful PR-Ready State Machine Request

This example shows the exact workflow used to create the PR-Ready State Machine:

#### 1. Orchestrator Documented Manual Process
The orchestrator agent documented 16 manual branch transformations in detail:
- Tracked each step and command
- Noted common conflicts and resolutions
- Identified patterns and automation opportunities
- Created comprehensive process document

#### 2. User Created Structured Request
Following the template, created `/tmp/SOFTWARE-FACTORY-MANAGER-REQUEST.md`:
- Clear background and motivation
- Detailed 7-phase process breakdown
- Specific requirements and challenges
- Expected deliverables list

#### 3. User Invoked Software-Factory-Manager
Used specific prompt:
```
@software-factory-manager

I've documented the PR-ready branch transformation process that we've
successfully used 16 times. Please review the comprehensive request at:
`/tmp/SOFTWARE-FACTORY-MANAGER-REQUEST.md`

Create a complete state machine with all supporting files.
```

#### 4. Manager Created Complete Implementation
Delivered:
- `SOFTWARE-FACTORY-PR-READY-STATE-MACHINE.md` (main state machine)
- Agent state rules for integration-agent
- `pr-ready-state.json.example` (example configuration)
- Validation scripts in `/tools/`
- Updated documentation

#### 5. Verification and Testing
- Reviewed state machine for completeness
- Checked all states had proper rules
- Verified agent assignments were appropriate
- Tested with sample branch transformation
- Confirmed all validations worked

### Key Success Factors
1. **Detailed manual process documentation** - Critical for accuracy
2. **Clear requirements and challenges** - Helped design decisions
3. **Specific deliverable requests** - Got complete implementation
4. **Structured request format** - Easy for manager to parse
5. **Reference to examples** - Grounded in reality

## Appendix A: Common Pitfalls to Avoid

### 1. Requesting Without Manual Experience
Never request a state machine for a process you haven't done manually. The manager needs real-world examples and edge cases.

### 2. Vague Requirements
Avoid statements like "handle errors appropriately." Be specific: "detect file deletion >1000 lines and abort with rollback."

### 3. Missing Agent Assignments
Don't leave agent assignments to the manager's discretion. Specify who should handle what based on agent capabilities.

### 4. No Success Criteria
Every phase needs measurable success criteria. "Complete" is not enough; specify what complete looks like.

### 5. Ignoring Dependencies
Document sequential requirements clearly. State machines can parallelize, but only when safe.

## Appendix B: Quick Reference

### File Locations
```
/tmp/                           # Request documents
/state-machines/                # State machine definitions
/agent-states/                  # Agent-specific rules
/tools/                         # Validation scripts
/.claude/commands/              # Continuation commands
/examples/                      # Example requests and configs
```

### Agent Capabilities Quick Reference
```
Orchestrator: Coordinate, spawn agents, manage state
Integration: Branches, merges, PRs, git operations
SW-Engineer: Code changes, file operations, testing
Code-Reviewer: Validation, compliance, quality checks
Architect: Design review, impact analysis, approval
```

### Request Template Location
See: `/templates/STATE-MACHINE-REQUEST-TEMPLATE.md`

### Example Requests
See: `/examples/state-machine-requests/`

---

**Remember:** The quality of your state machine depends on the quality of your request. Invest time in preparation and documentation for best results.