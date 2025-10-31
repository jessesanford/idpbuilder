# RULE R254 - Agent Error Reporting Protocol

## Purpose
Establish a clear protocol for agents to report critical errors back to the orchestrator, especially environment mismatches that prevent work from proceeding.

## The Problem
- Agents are spawned via Task tool and start in unpredictable directories
- Agents cannot modify their spawning environment
- Orchestrator needs to know when agents cannot proceed
- Silent failures waste time and violate grading requirements

## R254 Requirements

### For All Agents (SW Engineer, Code Reviewer, Architect)

When detecting a critical environment error, agents MUST:

1. **CAPTURE THE ORCHESTRATOR'S PROMPT** - Store it for echo-back
2. **STOP IMMEDIATELY** - Do not attempt to proceed or fix
3. **Report clearly** using this format:
```
❌ ENVIRONMENT ERROR: [Brief description]

🔴 ORCHESTRATOR, YOU GAVE ME THE WRONG PROMPT!

THIS IS THE PROMPT YOU GAVE:
════════════════════════════════════════
[COMPLETE ORCHESTRATOR PROMPT HERE]
════════════════════════════════════════

I FAILED TO FIND MY WORKING DIRECTORY BASED ON THIS PROMPT.

EXPECTED:
- Directory: /efforts/phase{X}/wave{Y}/{effort-name}
- Branch: {project-prefix}/phase{X}/wave{Y}/{effort-name}

ACTUAL:
- Directory: {current pwd}
- Branch: {current branch}
- Directory from prompt: {what was extracted or NOT SPECIFIED}

PLEASE TRY AGAIN WITH:
1. Correct directory path for my effort
2. Verification that infrastructure exists
3. Clear TARGET_DIRECTORY specification

GRADING VIOLATION: Proceeding would cause automatic failure per R208/R209
```

3. **Exit with error code** matching the issue:
   - `exit 208` - Wrong directory (R208 violation)
   - `exit 209` - Directory isolation violation (R209)
   - `exit 221` - Cannot navigate to directory (R221)
   - `exit 254` - General environment error

### For the Orchestrator

When receiving agent error reports, the orchestrator MUST:

1. **Read the error message** completely
2. **Update orchestrator-state-v3.json** with error details:
```yaml
agent_errors:
  - agent_type: "sw-engineer"
    effort: "effort-name"
    error_time: "timestamp"
    error_code: 208
    error_message: "Wrong directory"
    action_taken: "Will re-spawn with clearer instructions"
```
3. **Take corrective action**:
   - Verify infrastructure exists
   - Check directory permissions
   - Re-spawn with explicit navigation instructions

## Error Codes

| Code | Meaning | Agent Action | Orchestrator Action |
|------|---------|--------------|-------------------|
| 208 | Wrong directory | Stop and report | Verify infrastructure, re-spawn |
| 209 | Isolation violation | Stop and report | Check boundaries, re-spawn |
| 221 | Cannot CD to directory | Stop and report | Check permissions, verify path |
| 254 | General environment error | Stop and report | Investigate and correct |

## Implementation Examples

### SW Engineer Startup Check:
```bash
# First, capture the orchestrator's prompt
ORCHESTRATOR_PROMPT="[The complete prompt I received from orchestrator]"

# Mandatory environment verification
EXPECTED_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
CURRENT_DIR=$(pwd)

if [ "$CURRENT_DIR" != "$EXPECTED_DIR" ]; then
    # Try to navigate first
    if [ -d "$EXPECTED_DIR" ]; then
        cd "$EXPECTED_DIR"
        echo "✅ Navigated to correct directory"
    else
        # Directory doesn't exist - critical error with prompt echo-back
        echo "❌ ENVIRONMENT ERROR: Expected directory doesn't exist"
        echo ""
        echo "🔴 ORCHESTRATOR, YOU GAVE ME THE WRONG PROMPT!"
        echo ""
        echo "THIS IS THE PROMPT YOU GAVE:"
        echo "════════════════════════════════════════"
        echo "$ORCHESTRATOR_PROMPT"
        echo "════════════════════════════════════════"
        echo ""
        echo "EXPECTED:"
        echo "- Directory: $EXPECTED_DIR"
        echo ""
        echo "ACTUAL:"
        echo "- Directory: $CURRENT_DIR"
        echo "- Directory exists: NO"
        echo ""
        echo "ORCHESTRATOR, PLEASE:"
        echo "1. Run CREATE_NEXT_INFRASTRUCTURE for this effort"
        echo "2. Verify directory was created successfully"
        echo "3. Re-spawn me with the correct path"
        echo ""
        echo "GRADING VIOLATION: Cannot proceed without proper infrastructure"
        exit 208
    fi
fi
```

### Code Reviewer Validation:
```bash
# Verify we can access implementation plan
if [ ! -f "IMPLEMENTATION-PLAN.md" ]; then
    echo "❌ ENVIRONMENT ERROR: Cannot find IMPLEMENTATION-PLAN.md"
    echo ""
    echo "EXPECTED:"
    echo "- File: $(pwd)/IMPLEMENTATION-PLAN.md"
    echo ""
    echo "ACTUAL:"
    echo "- Directory: $(pwd)"
    echo "- Files present: $(ls -la)"
    echo ""
    echo "ORCHESTRATOR ACTION REQUIRED:"
    echo "1. Verify this effort had planning completed"
    echo "2. Check if plan was created in correct location"
    echo "3. Re-spawn with correct directory"
    exit 254
fi
```

## Success Confirmation

When environment is correct, agents should confirm:
```bash
echo "✅ ENVIRONMENT VERIFIED:"
echo "   Directory: $(pwd)"
echo "   Branch: $(git branch --show-current)"
echo "   State: READY TO PROCEED"
```

## Orchestrator Acknowledgment

The orchestrator should acknowledge when it sees an error:
```bash
echo "📍 ORCHESTRATOR: Received error report from $AGENT_TYPE"
echo "   Error code: $ERROR_CODE"
echo "   Taking corrective action..."
# Update state file
jq ".agent_errors += [{\"agent\": \"$AGENT_TYPE\", \"code\": $ERROR_CODE, \"time\": \"$(date -Iseconds)\"}]" orchestrator-state-v3.json
git add orchestrator-state-v3.json && git commit -m "state: agent error $ERROR_CODE recorded [R288]" && git push
```

## Grading Impact

- **Agents proceeding in wrong directory**: AUTOMATIC FAIL
- **Orchestrator ignoring error reports**: -25% orchestration score
- **Repeated same errors**: -10% per recurrence
- **Proper error handling**: +5% bonus for robustness

## The Goal

Create a robust feedback loop where:
1. Agents detect problems immediately
2. Report clearly to orchestrator
3. Orchestrator takes corrective action
4. Work proceeds in correct environment
5. No grading violations occur