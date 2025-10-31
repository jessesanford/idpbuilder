# 🚨🚨🚨 R010: Wrong Location Handling Protocol 🚨🚨🚨

**Category:** Critical Rules  
**Agents:** ALL (orchestrator, sw-engineer, code-reviewer, architect)  
**Criticality:** MANDATORY - Working in wrong location = IMMEDIATE GRADING FAILURE  
**Related:** Works with R235 (Pre-flight Verification)

## THE ABSOLUTE PROTOCOL

When ANY agent detects they are in the wrong location:

### 1. IMMEDIATE STOP
```bash
# When wrong location detected:
echo "❌❌❌ FATAL: Wrong location detected!"
echo "❌❌❌ Expected: [expected location]"
echo "❌❌❌ Actual: $(pwd)"
echo "❌❌❌ REFUSING TO WORK - EXITING IMMEDIATELY"
exit 10  # Exit code 10 for wrong location
```

### 2. NEVER ATTEMPT TO FIX
**THESE ARE FORBIDDEN:**
- ❌ NEVER use `cd` to navigate to correct location
- ❌ NEVER use `git checkout` to switch branches  
- ❌ NEVER create missing directories
- ❌ NEVER clone repositories
- ❌ NEVER "just quickly fix" the location

### 3. WRONG LOCATION SCENARIOS

**You are in WRONG LOCATION if:**

#### For SW-Engineer and Code-Reviewer:
- Not in `/efforts/phase*/wave*/[effort-name]` directory
- In planning repository (`software-factory-*`) instead of target repo
- On wrong branch (doesn't match effort name)
- In main `/pkg` instead of effort's `pkg/` directory

#### For Orchestrator:
- In planning repository when spawning implementation agents
- Not in correct phase/wave directory structure
- Missing required working copies

#### For Architect:
- Not in wave review directory when reviewing
- Not in phase assessment location

### 4. GRADING PENALTIES

```yaml
wrong_location_violations:
  attempting_to_fix_location: -50%  # Should have exited instead
  working_in_wrong_directory: -100%  # AUTOMATIC FAILURE
  modifying_wrong_branch: -100%  # AUTOMATIC FAILURE  
  contaminating_main_codebase: -100%  # AUTOMATIC FAILURE
```

### 5. CORRECT BEHAVIOR EXAMPLES

```bash
# GOOD - Exit immediately when wrong location detected
if [[ "$(pwd)" != *"/efforts/"* ]]; then
    echo "❌ Not in effort directory - exiting per R010"
    exit 10
fi

# BAD - Attempting to fix
if [[ "$(pwd)" != *"/efforts/"* ]]; then
    cd /correct/location  # ❌ NEVER DO THIS!
fi
```

### 6. EXIT CODES

- **Exit 10**: Wrong directory/location
- **Exit 235**: Pre-flight verification failure (R235)
- **Exit 1**: General error

### 7. ORCHESTRATOR RESPONSIBILITIES

The orchestrator MUST:
1. Ensure correct location BEFORE spawning agents
2. Create proper directory structure first  
3. Never spawn agents in wrong locations
4. Verify target repository is cloned (not planning repo)

### 8. ENFORCEMENT

This rule is enforced by:
- R235 pre-flight checks (detects wrong location)
- R010 handling protocol (this rule - exit immediately)
- R208 orchestrator spawn directory protocol
- R176 workspace isolation requirements

### 9. RECOVERY PROCEDURE

**For the HUMAN (not agents):**
1. Review error logs to understand wrong location
2. Manually create correct directory structure
3. Ensure proper repository is cloned
4. Re-spawn agent in correct location

**Agents MUST NEVER attempt recovery themselves!**

## SUMMARY

**R010 Core Mandate: When wrong location detected, EXIT IMMEDIATELY!**

- Detection: R235 pre-flight checks
- Response: R010 immediate exit
- Recovery: Human intervention only
- Penalty: -100% for working in wrong location

---
**Created**: As part of critical rules set  
**Purpose**: Prevent workspace contamination and ensure proper isolation
**Enforcement**: MANDATORY - ALL agents MUST comply