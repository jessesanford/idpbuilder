# Standard Exit Conditions Template for Orchestrator States

## Purpose
This template provides standard guidance for when to use CONTINUE-SOFTWARE-FACTORY=TRUE vs FALSE in orchestrator states. Use this template for states that spawn agents, monitor work, or wait for completion.

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ All required agents spawned successfully
- ✅ Infrastructure and metadata verified
- ✅ State transitions as expected
- ✅ Following designed workflow
- ✅ Everything working normally
- ✅ Monitoring/waiting proceeding as planned
- ✅ Ready to transition to next state

**THIS IS NORMAL WORKFLOW.** Spawning agents, waiting for results, and
transitioning between states is the DESIGNED PROCESS of the Software Factory.
Automation should continue unless something is genuinely broken.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Required files/infrastructure missing or corrupt
- ❌ Cannot spawn agents due to errors
- ❌ State machine corruption detected
- ❌ Unrecoverable error prevents proceeding
- ❌ Manual intervention truly required
- ❌ System in inconsistent state that automation cannot fix
- ❌ Critical validation failures

**DO NOT set FALSE because:**
- ❌ Normal state transitions
- ❌ Spawning agents (designed workflow)
- ❌ Waiting for work to complete
- ❌ R322 checkpoint (stop ≠ FALSE!)
- ❌ "User might want to review" (unless exceptional)
- ❌ Monitoring loops (NORMAL operation)
- ❌ Review-fix cycles (EXPECTED process)

### Critical Distinction: R322 Stop vs Continuation Flag

**R322 "stop"** = End conversation turn (`exit 0`)
- Purpose: Context preservation
- Always required at state checkpoints
- Prevents context window overflow

**Continuation flag** = Can system auto-restart?
- TRUE = Normal operations (default)
- FALSE = Exceptional/error conditions (rare)

**THESE ARE INDEPENDENT!**

**Correct pattern for normal operations:**
```bash
# Do state work
# Update state file
# Save todos
exit 0  # R322 stop (context preservation)
```

**Last line before exit:**
```bash
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Normal operation
```

**NOT:**
```bash
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # Only for errors!
```

### State-Specific Customization

When using this template:
1. Copy the entire "Exit Conditions and Continuation Flag" section
2. Customize the TRUE conditions list for your specific state
3. Customize the FALSE conditions list for your specific state
4. Add specific examples from your state if helpful
5. Keep the "Critical Distinction" section unchanged

### Grading Impact

- Using FALSE for normal operations: -20% per violation
- Pattern of incorrect FALSE usage: -50%
- Breaking automation unnecessarily: -30%
- Missing flag guidance entirely: -15%

### References

- R405: Automation Flag Continuation Principle
- R322: Mandatory Stop Before State Transitions
- R287: TODO Persistence Requirements
