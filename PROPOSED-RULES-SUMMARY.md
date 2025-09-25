# PROPOSED NEW RULES TO PREVENT ORCHESTRATOR CODE VIOLATIONS

## Summary of Proposed Rules

### R314: Orchestrator File Operation Prohibition (SUPREME LAW)
- **Status**: 🔴🔴🔴 SUPREME LAW
- **Purpose**: Absolutely forbid orchestrators from ANY operation on source code files
- **Key Points**:
  - No creating, reading, copying, moving, or modifying source files
  - Explicit list of forbidden file extensions (.go, .py, .js, etc.)
  - Automatic -100% failure for violation

### R315: Infrastructure vs Implementation Boundary (SUPREME LAW)
- **Status**: 🔴🔴🔴 SUPREME LAW
- **Purpose**: Crystal clear definition of what is infrastructure vs implementation
- **Infrastructure** (Allowed):
  - Creating empty directories
  - Git clone/branch operations
  - Creating .md/.yaml/.json files
- **Implementation** (Forbidden):
  - ANY operation on source code
  - Copying files between branches
  - Making implementation commits

### R316: Orchestrator Commit Restrictions (BLOCKING)
- **Status**: 🚨🚨🚨 BLOCKING
- **Purpose**: Prevent orchestrators from committing code files
- **Enforcement**: Git pre-commit hooks that block code file commits

### R317: Orchestrator Working Directory Restrictions (BLOCKING)
- **Status**: 🚨🚨🚨 BLOCKING
- **Purpose**: Keep orchestrators out of implementation directories
- **Key Point**: After creating infrastructure, orchestrator must never enter effort directories

### R318: Agent Failure Escalation Protocol (MANDATORY)
- **Status**: ⚠️⚠️⚠️ WARNING
- **Purpose**: Proper handling when agent spawns fail
- **Required**: Transition to ERROR_RECOVERY, never "do it yourself"

## States Requiring Updates

1. **SPAWN_AGENTS**
   - Add R314-R318 to mandatory rules
   - Explicit warnings about code operations

2. **CREATE_NEXT_SPLIT_INFRASTRUCTURE**
   - Clarify ONLY directory/branch creation allowed
   - Consider renaming to CREATE_NEXT_SPLIT_WORKSPACE

3. **ERROR_RECOVERY**
   - Add R318 escalation protocol
   - Clear recovery paths

## Technical Enforcement

1. **Git Pre-Commit Hooks**: Block orchestrator code commits
2. **Directory ACLs**: Restrict orchestrator access after creation
3. **Command Wrappers**: Validate operations before execution
4. **Real-Time Monitoring**: Detect violations immediately

## The Bright Line Rule

> "If it involves source code in ANY way, orchestrator CANNOT do it."

## Next Steps

**DO NOT IMPLEMENT YET** - Waiting for review and approval before:
1. Creating the new rule files in rule-library
2. Updating agent state rules
3. Implementing enforcement mechanisms
4. Testing with actual orchestrator runs