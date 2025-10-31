## Exit Conditions and Continuation Flag

**⚠️ CRITICAL: Read R405-CONTINUATION-FLAG-MASTER-GUIDE.md before setting this flag!**

**DEFAULT TO TRUE - Only use FALSE for truly unrecoverable errors!**

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (99.9% of cases - ALWAYS DEFAULT)

Use TRUE when:
- ✅ State work completed successfully
- ✅ Transitioning to next state normally
- ✅ Issues found but recoverable through designed protocols
- ✅ Spawning agents (NORMAL workflow)
- ✅ Waiting for results (NORMAL workflow)
- ✅ Monitoring progress (NORMAL workflow)
- ✅ R322 checkpoint (stop ≠ FALSE! This is context preservation!)
- ✅ Tests passed
- ✅ Tests failed but fixable (can enter fix cycle - NORMAL!)
- ✅ Review found issues (can spawn fixes - NORMAL!)
- ✅ Build succeeded
- ✅ Build failed but fixable (can debug and fix - NORMAL!)
- ✅ Integration completed (even with issues if fixable)
- ✅ Ready for next phase of workflow

**THIS IS THE NORMAL CASE!** The Software Factory is designed to handle
issues automatically through review-fix cycles, split protocols, and fix cascades.
If the system CAN recover automatically, use TRUE!

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (0.1% of cases - EXTREMELY RARE)

Use FALSE ONLY when:
- ❌ State file corrupted beyond recovery
- ❌ Infrastructure destroyed with no recovery path
- ❌ Truly unrecoverable error requiring manual fix
- ❌ State machine in invalid state
- ❌ Critical files missing with no way to recreate
- ❌ Unknown error with no defined recovery protocol

**DO NOT set FALSE because:**
- ❌ Work completed successfully (use TRUE!)
- ❌ Spawning agents (NORMAL workflow, use TRUE!)
- ❌ Found issues that can be fixed (NORMAL, use TRUE!)
- ❌ R322 requires stop (checkpoint ≠ FALSE flag!)
- ❌ "User might want to review" (NO! Use TRUE!)
- ❌ Tests/build failed (if fixable, use TRUE!)

**ASK: "Can the system recover automatically?"**
- YES → TRUE
- NO (truly stuck, needs human) → FALSE

**If in doubt, use TRUE. Erring toward FALSE causes -20% penalty per violation.**

## 🔴🔴🔴 R322 vs R405 - CRITICAL DISTINCTION (THE COMMON CONFUSION!) 🔴🔴🔴

### AGENT STOPS ≠ FACTORY STOPS (TWO INDEPENDENT DECISIONS!)

**Decision 1: Should Agent Stop? (R322 Mandatory Stop)**
- Agent completes state work
- Agent exits process to preserve context
- Agent waits for user to run /continue-[agent]
- **This is a TECHNICAL requirement for context preservation**

**Decision 2: Should Factory Continue? (R405 Continuation Flag)**
- Should automation proceed when user continues?
- Is the system healthy and ready for next state?
- Can the factory make forward progress?
- **This is an OPERATIONAL status indicator**

### THE CONFUSION TO ELIMINATE:

❌ **WRONG THINKING:**
"R322 says I must stop → Manual intervention needed → CONTINUE-SOFTWARE-FACTORY=FALSE"

✅ **CORRECT THINKING:**
"R322 says I must stop (technical) → Work completed successfully (operational) → CONTINUE-SOFTWARE-FACTORY=TRUE"

### THE TRUTH ABOUT R322 CHECKPOINTS:

**R322 "mandatory stop" means:**
- End conversation turn (`exit 0`) ✅
- Preserve context (prevent memory overflow) ✅
- Emit continuation flag (usually TRUE!) ✅
- Wait for user to run /continue-[agent] ✅

**R322 does NOT mean:**
- Set FALSE for normal operations ❌
- Require human review/decision ❌
- Halt automation ❌
- Block forward progress ❌

**NORMAL PATTERN AT R322 CHECKPOINTS:**
- **R322 + TRUE = Normal operation** (99.9% of checkpoints)
  - Agent exits (technical requirement)
  - Factory continues (operational status: healthy)
  - User runs /continue-[agent] later
  - Next state executes normally

**RARE PATTERN:**
- **R322 + FALSE = Unrecoverable error** (0.1% of checkpoints)
  - Agent exits (still preserves context)
  - Factory CANNOT continue (operational status: broken)
  - User must fix the catastrophic issue manually
  - System truly cannot make progress

### Implementation Pattern

```bash
# Correct pattern for normal completion (DEFAULT)
echo "✅ [State-specific success message]"
echo "✅ Ready for [NEXT_STATE]"
exit 0  # R322 stop if required
# Last line emitted: CONTINUE-SOFTWARE-FACTORY=TRUE

# Correct pattern for recoverable issues (STILL TRUE!)
echo "⚠️ [Issues found]"
echo "✅ Issues are fixable through [recovery protocol]"
echo "✅ Transitioning to [NEXT_STATE] for handling"
exit 0  # R322 stop if required
# Last line emitted: CONTINUE-SOFTWARE-FACTORY=TRUE (recoverable!)

# ONLY for truly unrecoverable errors (RARE!)
echo "🔴 UNRECOVERABLE: [specific unrecoverable error]"
echo "❌ Cannot proceed: [reason]"
echo "❌ Manual intervention required: [what's needed]"
exit 0
# Last line emitted: CONTINUE-SOFTWARE-FACTORY=FALSE (truly broken!)
```

### Grading Impact

- Wrong FALSE for normal operations: -20% per violation
- Pattern of excessive FALSE (3+): -50% cumulative
- Missing flag entirely: -100% AUTOMATIC FAILURE
- Text after flag: -50% penalty

## References

- R405: Automation Flag Continuation Principle (base rule)
- R405-CONTINUATION-FLAG-MASTER-GUIDE: Decision tree and examples
- R322: Mandatory Stop Before State Transitions (context preservation)
- R287: TODO Persistence (state management)
- R288: State File Updates (state management)
