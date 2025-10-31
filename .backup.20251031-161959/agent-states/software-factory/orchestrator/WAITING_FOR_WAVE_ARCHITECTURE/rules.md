# WAITING_FOR_WAVE_ARCHITECTURE State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

## 📋 PRIMARY DIRECTIVES

### Core Mandatory Rules:
1. R006 - Orchestrator Never Writes Code
2. R287 - TODO Persistence
3. R288 - State File Updates
4. R510 - Checklist Compliance
5. R405 - Continuation Flag
6. R233 - Active Monitoring Protocol

### State-Specific Rules:
7. **🚨🚨🚨 R340** - Architecture Plan Quality Gates
   - Summary: Validate wave architecture contains REAL CODE examples (not pseudocode)

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

### BLOCKING REQUIREMENTS

- [ ] 1. Wait for Architect to complete wave architecture
  - Expected: `wave-plans/WAVE-N-ARCHITECTURE.md`
  - Polling: Every 30 seconds (R233)
  - Timeout: 30 minutes

- [ ] 2. Validate wave architecture file exists
  - File: `wave-plans/WAVE-{N}-ARCHITECTURE.md`
  - Minimum size: >1000 bytes

- [ ] 3. Validate wave architecture quality (R340 - CONCRETE fidelity)
  - Check: Contains real function signatures with types
  - Check: Contains actual code examples (not pseudocode)
  - Check: Shows concrete interfaces
  - Check: NO abstract patterns only
  - **BLOCKING**: Must have REAL CODE (not "when X do Y")

- [ ] 4. Record adaptation notes
  - Check: Architect reviewed previous wave code
  - Check: Adaptation notes section exists
  - Purpose: Ensure learning from previous waves

### STANDARD EXECUTION TASKS

- [ ] 5. Record validation results in state file
- [ ] 6. Display validation summary

### EXIT REQUIREMENTS

- [ ] 7. Set proposed next state: `SPAWN_CODE_REVIEWER_WAVE_IMPL`
- [ ] 8. Spawn State Manager for SHUTDOWN_CONSULTATION
- [ ] 9. Save TODOs per R287
- [ ] 10. Set CONTINUE-SOFTWARE-FACTORY=TRUE
- [ ] 11. Stop execution (exit 0)

## State Purpose

Monitor Architect creating wave architecture, validate it contains REAL CODE examples from previous waves (not pseudocode), and ensure adaptation notes are present.

**Fidelity Enforced:** CONCRETE (real code, actual interfaces)

## Validation Criteria

**Must Have (CONCRETE fidelity):**
- Real function signatures: `async def process_payment(amount: Decimal) -> PaymentResult:`
- Actual class definitions with methods
- Working code examples
- Concrete type annotations
- Adaptation notes from previous waves

**Must NOT Have:**
- Pseudocode: "when payment received, process it"
- Abstract descriptions without code
- Missing type annotations
- No references to previous wave code

## Entry Criteria

- **From**: SPAWN_ARCHITECT_WAVE_PLANNING
- **Required**: Architect spawned for wave architecture

## Exit Criteria

**Success** → SPAWN_CODE_REVIEWER_WAVE_IMPL
**Failure** → ERROR_RECOVERY (if fidelity wrong or timeout)

## Rules Enforced

- R510, R340, R233, R288, R287, R405, R006

## Additional Context

SF 3.0 Progressive Planning Step 3 Validation:
- Phase used pseudocode ✅
- Phase had wave list ✅  
- **Wave uses REAL CODE** ← THIS STATE validates
- Wave will have exact specs (next)

Common Failures:
1. Architect uses pseudocode instead of real code → Request revision
2. Missing adaptation notes → Flag but may proceed
3. No previous wave references → Warning (first wave OK)
