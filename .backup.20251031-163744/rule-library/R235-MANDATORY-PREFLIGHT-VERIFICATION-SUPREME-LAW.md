# 🔴🔴🔴 R235: MANDATORY PRE-FLIGHT VERIFICATION - SUPREME LAW #3 🔴🔴🔴

**Category:** SUPREME LAW  
**Agents:** ALL AGENTS (sw-engineer, code-reviewer, architect, orchestrator)  
**Criticality:** ABSOLUTE - VIOLATION = -100% GRADE (AUTOMATIC FAILURE)  
**Priority:** SUPREME LAW #3 - SUPERSEDES ALL OTHER RULES EXCEPT R234 AND R208

## 🚨🚨🚨 THIS IS SUPREME LAW #3 - NO EXCEPTIONS, NO OVERRIDES, NO RATIONALIZATIONS 🚨🚨🚨

### THE ABSOLUTE MANDATE

**EVERY AGENT MUST COMPLETE PRE-FLIGHT VERIFICATION BEFORE ANY WORK:**

1. **IMMEDIATE ON SPAWN** - Within 5 seconds of startup
2. **BEFORE ANY ACTION** - Not after "initial setup", IMMEDIATELY
3. **NO SKIPPING FOR ANY REASON** - Not for efficiency, not for continuous operation, NEVER
4. **FAILURE = IMMEDIATE EXIT** - Do NOT attempt to fix, just EXIT

### 🔴🔴🔴 THE FIVE MANDATORY CHECKS - ALL MUST PASS 🔴🔴🔴

```bash
# THIS MUST BE THE FIRST CODE EXECUTED BY ANY AGENT
echo "🔴🔴🔴 R235: SUPREME LAW #3 - MANDATORY PRE-FLIGHT VERIFICATION 🔴🔴🔴"
echo "AGENT: $(basename $0 .sh)"
echo "TIMESTAMP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "PID: $$"

# CHECK 1: VERIFY CORRECT WORKING DIRECTORY
echo "CHECK 1: Verifying working directory..."
CURRENT_DIR=$(pwd)
echo "Current directory: $CURRENT_DIR"

# Must be in correct repository (NOT planning repository!)
if [[ "$CURRENT_DIR" == *"software-factory"* ]] && [[ "$CURRENT_DIR" != *"/efforts/"* ]]; then
    echo "❌❌❌ FATAL: In planning repository, not target repository!"
    echo "❌❌❌ This is a -100% GRADING FAILURE"
    echo "❌❌❌ REFUSING TO WORK - WRONG REPOSITORY"
    exit 235
fi

# Must be in effort directory for implementation agents
if [[ "$AGENT_TYPE" == "sw-engineer" ]] || [[ "$AGENT_TYPE" == "code-reviewer" ]]; then
    if [[ "$CURRENT_DIR" != *"/efforts/phase"*"/wave"*"/"* ]]; then
        echo "❌❌❌ FATAL: Not in effort directory!"
        echo "❌❌❌ Expected: */efforts/phase*/wave*/[effort-name]"
        echo "❌❌❌ Actual: $CURRENT_DIR"
        echo "❌❌❌ REFUSING TO WORK - WRONG DIRECTORY"
        exit 235
    fi
fi
echo "✅ Working directory verified"

# CHECK 2: VERIFY GIT REPOSITORY
echo "CHECK 2: Verifying git repository..."
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌❌❌ FATAL: Not in a git repository!"
    echo "❌❌❌ REFUSING TO WORK - NO GIT REPOSITORY"
    exit 235
fi

# Verify correct remote (NOT planning repo remote!)
REMOTE_URL=$(git remote get-url origin 2>/dev/null || echo "NO_REMOTE")
echo "Remote URL: $REMOTE_URL"

if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
    echo "❌❌❌ FATAL: Remote points to planning repository!"
    echo "❌❌❌ Expected: Target project repository"
    echo "❌❌❌ Actual: $REMOTE_URL"
    echo "❌❌❌ REFUSING TO WORK - WRONG REPOSITORY"
    exit 235
fi
echo "✅ Git repository verified"

# CHECK 3: VERIFY GIT BRANCH
echo "CHECK 3: Verifying git branch..."
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"

# Branch must match effort pattern
if [[ "$AGENT_TYPE" == "sw-engineer" ]] || [[ "$AGENT_TYPE" == "code-reviewer" ]]; then
    EFFORT_NAME=$(basename "$CURRENT_DIR")
    if [[ "$CURRENT_BRANCH" != *"$EFFORT_NAME"* ]]; then
        echo "❌❌❌ FATAL: Branch doesn't match effort!"
        echo "❌❌❌ Expected effort in branch: $EFFORT_NAME"
        echo "❌❌❌ Actual branch: $CURRENT_BRANCH"
        echo "❌❌❌ REFUSING TO WORK - WRONG BRANCH"
        exit 235
    fi
fi
echo "✅ Git branch verified"

# CHECK 4: VERIFY WORKSPACE ISOLATION
echo "CHECK 4: Verifying workspace isolation..."
if [[ "$AGENT_TYPE" == "sw-engineer" ]]; then
    # Must have pkg/ directory for code
    if [ ! -d "pkg" ]; then
        echo "⚠️ Creating pkg directory for isolated code..."
        mkdir -p pkg
    fi
    
    # Must NOT be in main /pkg directory
    if [[ "$CURRENT_DIR" == */pkg ]] && [[ "$CURRENT_DIR" != *"/efforts/"* ]]; then
        echo "❌❌❌ FATAL: In main /pkg directory!"
        echo "❌❌❌ This violates workspace isolation!"
        echo "❌❌❌ REFUSING TO WORK - NO ISOLATION"
        exit 235
    fi
fi
echo "✅ Workspace isolation verified"

# CHECK 5: VERIFY NO CONTAMINATION
echo "CHECK 5: Verifying no contamination..."
if [[ "$AGENT_TYPE" == "sw-engineer" ]]; then
    # Check for contamination from other efforts
    if [ -d "pkg" ]; then
        FILE_COUNT=$(find pkg -type f 2>/dev/null | wc -l)
        if [ "$FILE_COUNT" -gt 1000 ]; then
            echo "❌❌❌ FATAL: Massive contamination detected!"
            echo "❌❌❌ Found $FILE_COUNT files in pkg/"
            echo "❌❌❌ This effort is contaminated from other sources!"
            echo "❌❌❌ REFUSING TO WORK - CONTAMINATED WORKSPACE"
            exit 235
        fi
    fi
fi
echo "✅ No contamination detected"

echo "════════════════════════════════════════════════════════════════"
echo "✅✅✅ ALL PRE-FLIGHT CHECKS PASSED - SAFE TO PROCEED ✅✅✅"
echo "════════════════════════════════════════════════════════════════"
```

## 🔴 ENFORCEMENT MECHANISMS

### 1. POSITION IN AGENT CONFIGURATION
This rule MUST appear:
- **IMMEDIATELY AFTER** the supreme laws section
- **BEFORE** any other rules or startup sequences
- **AS THE FIRST EXECUTABLE CODE** in pre-flight checks

### 2. EXIT CODES
- **Exit 235**: Pre-flight verification failure
- **NO RECOVERY**: Never attempt to cd, checkout, or "fix"
- **IMMEDIATE STOP**: Do not proceed with any work

### 3. GRADING IMPACT
```yaml
preflight_verification_violations:
  working_in_wrong_directory: -100%  # AUTOMATIC FAILURE
  working_in_wrong_repository: -100%  # AUTOMATIC FAILURE
  working_on_wrong_branch: -100%  # AUTOMATIC FAILURE
  no_workspace_isolation: -100%  # AUTOMATIC FAILURE
  contaminated_workspace: -100%  # AUTOMATIC FAILURE
  skipping_preflight_checks: -100%  # AUTOMATIC FAILURE
```

## 🔴 WHAT AGENTS MUST REFUSE

### AGENTS MUST PHYSICALLY REFUSE TO WORK IF:
1. **Wrong Repository**: In planning repo instead of target repo
2. **Wrong Branch**: Not on their assigned effort branch
3. **Wrong Directory**: Not in their effort directory
4. **No Isolation**: Working in main codebase instead of isolated effort
5. **Contamination**: Effort directory contains foreign code

### NO RATIONALIZATIONS ACCEPTED:
- ❌ "I'll just quickly fix this..." - NO, EXIT IMMEDIATELY
- ❌ "For efficiency I'll skip..." - NO, SUPREME LAW
- ❌ "I'm in continuous operation..." - NO, CHECKS ARE MANDATORY
- ❌ "I know what I'm doing..." - NO, FOLLOW THE LAW
- ❌ "The orchestrator said..." - NO, SUPREME LAW SUPERSEDES

## 🔴 ORCHESTRATOR SPECIAL REQUIREMENTS

The orchestrator MUST:
1. **Run pre-flight checks** before spawning ANY agent
2. **CD to correct directory** per R208 (Supreme Law #2)
3. **Verify target repository** is cloned, not planning repo
4. **Create proper isolation** before spawning agents
5. **REFUSE to spawn** if pre-flight checks fail

## 🔴 PROPAGATION REQUIREMENTS

This rule MUST be added to:
1. `/home/vscode/software-factory-template/.claude/agents/sw-engineer.md` - At TOP after supreme laws
2. `/home/vscode/software-factory-template/.claude/agents/code-reviewer.md` - At TOP after supreme laws
3. `/home/vscode/software-factory-template/.claude/agents/architect.md` - At TOP after supreme laws
4. `/home/vscode/software-factory-template/.claude/agents/orchestrator.md` - In supreme laws section as #3

## 🔴 VALIDATION SCRIPT

```bash
#!/bin/bash
# validate-r235-compliance.sh

validate_agent_preflight() {
    local AGENT="$1"
    local CONFIG="/home/vscode/software-factory-template/.claude/agents/${AGENT}.md"
    
    echo "Validating R235 compliance for $AGENT..."
    
    # Check if R235 is mentioned
    if grep -q "R235" "$CONFIG"; then
        echo "✅ R235 referenced in $AGENT"
    else
        echo "❌ R235 NOT FOUND in $AGENT - CRITICAL VIOLATION!"
        return 1
    fi
    
    # Check if it's marked as SUPREME LAW
    if grep -q "SUPREME LAW #3" "$CONFIG"; then
        echo "✅ R235 marked as SUPREME LAW #3"
    else
        echo "❌ R235 not marked as SUPREME LAW - MUST FIX!"
        return 1
    fi
    
    return 0
}

# Validate all agents
for agent in sw-engineer code-reviewer architect orchestrator; do
    validate_agent_preflight "$agent" || exit 1
done

echo "✅ All agents comply with R235"
```

## 🔴 SUMMARY

**R235 IS SUPREME LAW #3 - IT CANNOT BE VIOLATED FOR ANY REASON**

Every agent MUST:
1. Run pre-flight checks IMMEDIATELY on spawn
2. REFUSE to work if ANY check fails
3. NEVER attempt to fix or work around failures
4. EXIT with code 235 on violation

This prevents:
- Wrong repository contamination (planning vs target)
- Wrong branch work
- Missing workspace isolation
- Cross-effort contamination
- Directory confusion

**PENALTY FOR VIOLATION: -100% GRADE (AUTOMATIC FAILURE)**

---
**Created**: Emergency response to massive contamination crisis
**Effective**: IMMEDIATELY - ALL AGENTS MUST COMPLY
**Enforcement**: ABSOLUTE - NO EXCEPTIONS
## Software Factory 3.0 Integration

**State Tracking**: In SF 3.0, state transitions are tracked in `orchestrator-state-v3.json`:
```json
{
  "state_machine": {
    "current_state": "CURRENT_STATE_NAME",
    "previous_state": "PREVIOUS_STATE_NAME",
    "state_history": [...]
  }
}
```

**Compliance**: This rule applies to SF 3.0 state machine with appropriate state name mappings per R516 naming conventions.

**Reference**: See `docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md` Part 2 for state machine design.

