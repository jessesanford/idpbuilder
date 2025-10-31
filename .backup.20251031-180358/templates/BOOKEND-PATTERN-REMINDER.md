# STATE MANAGER CONSULTATION BOOKEND PATTERN (MANDATORY)

## CRITICAL REMINDER FOR ALL ORCHESTRATOR STATES

**This pattern is MANDATORY for ALL orchestrator states. -100% failure for bypass.**

### Before This State Begins

**STARTUP_CONSULTATION (MANDATORY):**
- State Manager was consulted at state entry
- State Manager validated current state
- State Manager provided directive_report
- You are here because State Manager directed you here

### After This State Completes

**SHUTDOWN_CONSULTATION (MANDATORY):**
- Prepare state update payload with work results
- Prepare transition proposal with:
  - Completed work summary
  - Proposed next state
  - Reasoning for proposal
- Spawn State Manager SHUTDOWN_CONSULTATION
- Wait for State Manager's required_next_state
- Transition to State Manager's directed state (NOT your proposal)

### Critical Distinctions

**YOU (Orchestrator) Role:**
- Execute state-specific work
- Prepare state updates (but don't commit them)
- Propose next state based on work results
- **Provide opinion, NOT decision**

**STATE MANAGER Role:**
- Read state machine allowed_transitions
- Validate proposal against state machine rules
- Enforce mandatory_sequences
- Make FINAL DECISION on next state
- Update all 4 state files atomically
- **Make decision, NOT just validate**

### Example: State Manager Overrides Your Proposal

```markdown
You: "I completed master architecture. Proposing: SPAWN_ARCHITECT_PHASE_PLANNING"

State Manager: "Proposal received. Checking state machine...
               R341 TDD requirements detected.
               Required next state: SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
               Reason: Tests MUST be planned before implementation architecture."

You: "Acknowledged. Transitioning to SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING per State Manager directive."
```

### What You MUST NEVER Do

❌ Update orchestrator-state-v3.json yourself
❌ Update bug-tracking.json yourself
❌ Update integration-containers.json yourself
❌ Choose next state without consulting State Manager
❌ Bypass SHUTDOWN_CONSULTATION
❌ Ignore State Manager's required_next_state
❌ Set validated_by: "orchestrator" in state history

### What State Manager Returns

**NOT** `recommended_next_state` - that implies you can ignore it
**YES** `required_next_state` - this is a command, not a suggestion

State Manager's decision is FINAL and binding. You MUST follow it.

---

**See Also:**
- /home/vscode/software-factory-template/.claude/agents/orchestrator.md (bookend pattern section)
- /home/vscode/software-factory-template/.claude/agents/state-manager.md
- /home/vscode/software-factory-template/state-machines/software-factory-3.0-state-machine.json (mandatory_sequences)
