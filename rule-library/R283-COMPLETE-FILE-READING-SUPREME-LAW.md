# 🔴🔴🔴 RULE R283: COMPLETE FILE READING (SUPREME LAW) 🔴🔴🔴

## Rule Type
**SUPREME LAW - ABSOLUTE ENFORCEMENT**

## Criticality  
**🔴🔴🔴 BLOCKING - Violation = Immediate -100% Failure**

## Description
All agents MUST read EVERY SINGLE LINE of EVERY rule file they are required to read. Partial reading of critical files is an IMMEDIATE FAILURE condition.

## Requirements

### 1. COMPLETE FILE READING
- **NEVER** stop at 100, 200, or any arbitrary line limit
- **ALWAYS** read until the LAST LINE of the file
- If a file is too large, read it in chunks but **READ IT ALL**
- **NEVER** "mark as read" without reading the ENTIRE file
- **MUST** print the last line of each file in acknowledgment

### 2. ACKNOWLEDGMENT FORMAT
For EACH rule file read, acknowledgment MUST include:
```
"I have read ALL [X] lines of [filename], ending with: '[last line of file]'"
```

### 3. CONTINUATION READING
When files exceed initial read limit:
```
Read(filename)                    # First chunk
Read(filename, offset=100)        # Continue from line 100
Read(filename, offset=200)        # Continue from line 200
[Continue until end of file]
```

### 4. VALIDATION REQUIREMENTS
Before proceeding with any work:
- Total lines read for each file must be verified
- Last line of each file must be quoted
- Explicit confirmation of complete reading required

## Anti-Patterns (IMMEDIATE FAILURE)

### ❌ PARTIAL READING
```
Read(SOFTWARE-FACTORY-STATE-MACHINE.md)
⎿ Read 100 lines
Thinking: "This file is large, but I need to mark it as read"
```
**FAILURE - File not completely read!**

### ❌ ASSUMING COMPLETION
```
Read(rule-library/R234.md)
⎿ Read 200 lines
"I've read R234"
```
**FAILURE - Unknown if more lines exist!**

### ❌ SKIPPING CONTENT
```
Read(orchestrator-rules.md)
⎿ Read 100 lines
"I'll continue with other files"
```
**FAILURE - Critical content missed!**

## Correct Patterns

### ✅ COMPLETE READING WITH CONTINUATION
```
Read(SOFTWARE-FACTORY-STATE-MACHINE.md)
⎿ Read 100 lines
Read(SOFTWARE-FACTORY-STATE-MACHINE.md, offset=100)
⎿ Read 100 lines
Read(SOFTWARE-FACTORY-STATE-MACHINE.md, offset=200)
⎿ Read 100 lines
[Continue until reaching end]
"I have read ALL 847 lines of SOFTWARE-FACTORY-STATE-MACHINE.md, ending with: '## END OF STATE MACHINE SPECIFICATION'"
```

### ✅ VERIFICATION OF COMPLETION
```
Read(rule-library/R234.md)
⎿ Read 300 lines
Read(rule-library/R234.md, offset=300)  
⎿ Read 45 lines
"I have read ALL 345 lines of R234.md, ending with: '## Enforcement: Immediate failure for violations'"
```

## Error Recovery

If unable to read entire file:
1. **STOP** immediately
2. **Report**: "CRITICAL ERROR: Cannot read entire [filename]"
3. **DO NOT PROCEED** with partial knowledge
4. Transition to ERROR_RECOVERY state if applicable

## Enforcement

- **Monitoring**: System tracks Read tool calls and verifies complete reading
- **Validation**: Line counts and last lines must match actual file contents
- **Penalty**: Partial reading = IMMEDIATE -100% FAILURE
- **No Exceptions**: This rule applies to ALL agents, ALL rule files, ALL states

## Affected Files
- All agent configurations in `.claude/agents/`
- All rule files in `rule-library/`
- All state rules in `agent-states/*/rules.md`
- SOFTWARE-FACTORY-STATE-MACHINE.md
- Any file marked as "MUST READ" or "MANDATORY"

## Related Rules
- R203: State-aware startup (requires reading state rules)
- R234: Mandatory state traversal (requires reading state machine)
- R290: State rule reading and verification (consolidated from R236/R237)
- R217: Post-transition re-acknowledgment

## Version History
- v1.0.0 (2025-01-29): Initial creation to prevent partial file reading failures