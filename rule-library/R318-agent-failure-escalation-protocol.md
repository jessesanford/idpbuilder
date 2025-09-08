# 🚨🚨🚨 RULE R318 - Agent Failure Escalation Protocol

**Criticality:** BLOCKING  
**Grading Impact:** -40% for attempting forbidden fixes  
**Enforcement:** Strict escalation requirements
**Applies To:** Orchestrator when agents encounter errors

## Rule Statement

When an agent fails or encounters an error, the Orchestrator MUST NOT attempt to fix the issue directly. The Orchestrator MUST either respawn the agent with better instructions OR escalate to human intervention. DIY fixes are FORBIDDEN.

## 🔴🔴🔴 ABSOLUTE PROHIBITIONS 🔴🔴🔴

When an agent fails, the Orchestrator MUST NEVER:
- ❌ Attempt to fix code issues directly
- ❌ Complete the agent's task personally  
- ❌ Modify files the agent was working on
- ❌ "Quick fix" problems to save time
- ❌ Debug code on agent's behalf
- ❌ Run tests to diagnose issues
- ❌ Edit implementation to resolve errors

## Required Escalation Protocol

### STEP 1: Detect Agent Failure
```bash
# Monitor agent status
if [[ "$AGENT_EXIT_CODE" -ne 0 ]]; then
    echo "🚨 Agent failure detected: Exit code $AGENT_EXIT_CODE"
    initiate_escalation_protocol
fi
```

### STEP 2: Analyze Failure Type
```bash
analyze_agent_failure() {
    local agent_name="$1"
    local failure_output="$2"
    
    echo "📊 Analyzing failure for: $agent_name"
    
    # Categorize failure
    if [[ "$failure_output" =~ "test.*fail" ]]; then
        echo "Type: Test failure"
        echo "Action: Respawn with fix instructions"
    elif [[ "$failure_output" =~ "compilation.*error" ]]; then
        echo "Type: Compilation error"
        echo "Action: Respawn with error details"
    elif [[ "$failure_output" =~ "permission.*denied" ]]; then
        echo "Type: Infrastructure issue"
        echo "Action: Fix infrastructure, respawn"
    else
        echo "Type: Unknown"
        echo "Action: Escalate to human"
    fi
}
```

### STEP 3: Choose Escalation Path

## ✅ ALLOWED RESPONSES:

### Option A: Respawn with Better Instructions
```bash
# Orchestrator identifies issue
echo "🔍 Agent failed due to missing dependencies"

# Respawn with specific guidance
echo "🚀 Respawning SW Engineer with fix instructions..."

Task: Fix implementation with dependency installation
Agent: sw-engineer
Additional Context:
- Error: "Import error: numpy not found"
- Action needed: Add numpy to requirements.txt
- Then retry the implementation
```

### Option B: Spawn Different Agent
```bash
# If wrong agent was used
echo "🔍 Task requires different expertise"

# Spawn appropriate agent
echo "🚀 Spawning Code Reviewer instead..."

Task: Review and identify issues
Agent: code-reviewer
Context: SW Engineer unable to resolve test failures
```

### Option C: Escalate to Human
```bash
# When truly stuck
echo "🚨 ESCALATION REQUIRED"
echo "Issue: Agent repeatedly failing on authentication module"
echo "Attempts: 3 respawns with different instructions"
echo "Recommendation: Human intervention needed"
echo "Status: BLOCKED"

# Update state
yq -i '.blocked_reason = "Agent failure - human intervention required"' \
    orchestrator-state.json
```

## ❌ FORBIDDEN RESPONSES:

### VIOLATION: Direct Fix Attempt
```bash
# ❌ ORCHESTRATOR MUST NOT DO THIS
echo "Agent can't fix test, I'll do it"
edit test_auth.go  # VIOLATION: Orchestrator editing code
fix_line "assert.Equal" "assert.True"  # VIOLATION: Direct fix
```

### VIOLATION: Completing Agent Task
```bash
# ❌ ORCHESTRATOR MUST NOT DO THIS
echo "SW Engineer failed, I'll finish"
cp implementation.go efforts/phase1/  # VIOLATION: Doing agent's job
git add implementation.go  # VIOLATION: Committing code
```

### VIOLATION: Debug Operations
```bash
# ❌ ORCHESTRATOR MUST NOT DO THIS
cd efforts/phase1/wave1  # VIOLATION: Entering workspace
go test ./...  # VIOLATION: Running tests
cat error.log  # VIOLATION: Debugging directly
```

## Correct Escalation Examples

### ✅ GOOD: Respawn for Test Failure
```bash
# Agent reports test failure
echo "📋 SW Engineer reports 3 test failures"

# Analyze without entering directory
echo "🔍 Reviewing failure report..."

# Respawn with specific instructions
echo "🚀 Respawning SW Engineer..."

Task: Fix failing tests in auth module
Agent: sw-engineer
Specific Issues:
1. TestLogin - expects JWT token format
2. TestLogout - needs session cleanup
3. TestRefresh - missing expiry validation
```

### ✅ GOOD: Infrastructure-Only Fix
```bash
# Agent can't find directory
echo "📋 Agent reports: Directory not found"

# Fix ONLY infrastructure
mkdir -p efforts/phase1/wave1/split2  # OK: Empty directory

# Respawn agent
echo "🚀 Infrastructure ready, respawning agent..."
```

## 🔴 MANDATORY 3-FAILURE ESCALATION THRESHOLD 🔴

**CRITICAL REQUIREMENT**: The orchestrator MUST track failure attempts and escalate after 3 failures:

### Failure Attempt Tracking Protocol
```yaml
# In orchestrator-state.json
failure_tracking:
  effort_name:
    agent_type: "sw-engineer"
    failure_count: 3  # Current count
    max_attempts: 3   # HARD LIMIT
    failures:
      - attempt: 1
        timestamp: "2024-01-20T10:00:00Z"
        reason: "test failures in auth module"
        action: "respawned with specific fix instructions"
      - attempt: 2
        timestamp: "2024-01-20T10:15:00Z"
        reason: "test failures persist"
        action: "spawned code reviewer for assistance"
      - attempt: 3
        timestamp: "2024-01-20T10:30:00Z"
        reason: "still failing after review"
        action: "ESCALATED TO HUMAN - BLOCKED"
```

### The 3-Strike Rule
```bash
handle_agent_failure() {
    local effort="$1"
    local agent_type="$2"
    local failure_reason="$3"
    
    # Get current failure count
    local failure_count=$(yq ".failure_tracking.\"$effort\".failure_count // 0" orchestrator-state.json)
    failure_count=$((failure_count + 1))
    
    echo "🚨 Failure #$failure_count for $effort ($agent_type)"
    
    case $failure_count in
        1)
            echo "📋 First failure - Respawning with detailed instructions"
            # Update tracking
            yq -i ".failure_tracking.\"$effort\".failure_count = $failure_count" orchestrator-state.json
            # Respawn with specific guidance
            spawn_with_fix_instructions "$agent_type" "$effort" "$failure_reason"
            ;;
        2)
            echo "⚠️ Second failure - Trying different approach"
            # Update tracking
            yq -i ".failure_tracking.\"$effort\".failure_count = $failure_count" orchestrator-state.json
            # Try different agent or approach
            spawn_alternate_approach "$effort" "$failure_reason"
            ;;
        3)
            echo "🔴 THIRD FAILURE - MANDATORY ESCALATION"
            echo "❌ 3-STRIKE LIMIT REACHED"
            # Update tracking
            yq -i ".failure_tracking.\"$effort\".failure_count = $failure_count" orchestrator-state.json
            yq -i ".failure_tracking.\"$effort\".status = \"BLOCKED_HUMAN_REQUIRED\"" orchestrator-state.json
            # Create escalation report
            create_escalation_report "$effort" "$agent_type" "$failure_reason"
            echo "🛑 STOPPED - Human intervention required"
            ;;
        *)
            echo "❌ CRITICAL: Exceeded 3-failure limit!"
            echo "This should never happen - check escalation logic"
            exit 318
            ;;
    esac
}
```

## Escalation Decision Matrix

| Failure Type | Attempt # | Required Action | Notes |
|-------------|-----------|-----------------|-------|
| Test failure | 1 | Respawn with specific fixes | Include exact test names |
| Test failure | 2 | Spawn Code Reviewer | Get different perspective |
| Test failure | 3 | **ESCALATE TO HUMAN** | **MANDATORY STOP** |
| Build error | 1 | Respawn with fix hints | Include error messages |
| Build error | 2 | Check dependencies/environment | May spawn different agent |
| Build error | 3 | **ESCALATE TO HUMAN** | **MANDATORY STOP** |
| Missing deps | 1 | Respawn with dep list | Provide exact packages |
| Missing deps | 2 | Verify infrastructure setup | Check environment |
| Missing deps | 3 | **ESCALATE TO HUMAN** | **MANDATORY STOP** |
| Any error | 3 | **ALWAYS ESCALATE** | **NO EXCEPTIONS** |

## Tracking Escalations

```yaml
# In orchestrator-state.json
escalations:
  - timestamp: "2024-01-20T10:30:00Z"
    agent: "sw-engineer-1"
    failure: "test failures"
    action: "respawned with instructions"
    attempt: 1
  - timestamp: "2024-01-20T10:45:00Z"  
    agent: "sw-engineer-1"
    failure: "test failures persist"
    action: "escalated to human"
    attempt: 2
```

## Enforcement Mechanism

```bash
# Detect forbidden fix attempts
detect_forbidden_fix() {
    local orchestrator_action="$1"
    
    # Check for direct intervention
    if [[ "$orchestrator_action" =~ (edit|fix|debug|test|compile|run) ]]; then
        echo "🚨 R318 VIOLATION: Orchestrator attempting direct fix!"
        echo "❌ FORBIDDEN: $orchestrator_action"
        echo "✅ REQUIRED: Respawn agent or escalate"
        return 1
    fi
}
```

## Grading Impact

- Attempting direct code fix: -40%
- Running tests/debugging: -30%
- Completing agent's task: -50%
- Proper escalation: +5% bonus

## Self-Check Protocol

Before responding to any agent failure:
1. ❓ Am I about to fix this myself? → STOP
2. ❓ Am I entering their workspace? → STOP
3. ❓ Am I touching code files? → STOP
4. ✅ Am I respawning or escalating? → PROCEED

---
**REMEMBER:** Orchestrator is a manager, not a fixer. When agents fail, get them help, don't become the help.