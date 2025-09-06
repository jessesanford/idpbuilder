# 🚨 EMERGENCY PROCEDURES QUICK REFERENCE

## 🔥 IMMEDIATE EMERGENCY RESPONSE

### 🚨 WRONG DIRECTORY/BRANCH DETECTED
```
┌─────────────────────────────────────────────┐
│ 🛑 STOP IMMEDIATELY - DO NOT PROCEED       │
│ NEVER attempt to cd or checkout to "fix"   │
│ Report error and wait for correction       │
│ Working in wrong location = GRADING FAIL   │
└─────────────────────────────────────────────┘

EMERGENCY PROTOCOL:
1. echo "🚨 EMERGENCY: Wrong environment detected"
2. pwd && git branch --show-current  
3. echo "Expected: [CORRECT_PATH] on [CORRECT_BRANCH]"
4. echo "STOPPING - Orchestrator intervention required"
5. exit 1
```

### 🚨 SIZE LIMIT EXCEEDED (>800 lines)
```
SW ENGINEER EMERGENCY PROTOCOL:
1. STOP coding immediately
2. lines=$(line-counter.sh -c ${BRANCH} | grep "Total:" | awk '{print $2}')
3. echo "🚨 SIZE VIOLATION: $lines/800 lines - MANDATORY SPLIT"
4. git add . && git commit -m "emergency: size limit exceeded, requesting split"
5. Report to orchestrator: "NEEDS_SPLIT - $lines lines"
6. DO NOT continue implementation

CODE REVIEWER EMERGENCY PROTOCOL:
1. Document violation immediately
2. Create emergency split plan
3. Each split MUST be <800 lines
4. Block all approvals until splits created
5. Report to orchestrator: "CRITICAL_SPLIT_REQUIRED"
```

### 🚨 CRITICAL ARCHITECTURAL ISSUE
```
ARCHITECT EMERGENCY PROTOCOL:
1. Issue immediate STOP decision
2. Document specific problems clearly:
   - Security vulnerabilities
   - Data corruption risks
   - Multi-tenancy violations
   - Performance showstoppers
3. Provide detailed remediation plan
4. Coordinate with orchestrator for recovery
5. No wave progression until resolved

ORCHESTRATOR RESPONSE:
1. Receive STOP decision → Transition to HARD_STOP
2. Record failure reason in state file
3. Notify all agents: "ARCHITECTURAL_EMERGENCY"
4. Do NOT continue to next wave
5. Wait for architect remediation plan
```

## 🔧 TECHNICAL EMERGENCY PROCEDURES

### 🚨 AGENT NOT RESPONDING
```
ORCHESTRATOR DETECTION:
if last_activity_timestamp > 30_minutes_ago:
    echo "🚨 AGENT TIMEOUT: ${AGENT_ID} not responding"
    
RECOVERY PROTOCOL:
1. Check agent's last known state
2. Attempt to restart agent with same context
3. If restart fails twice:
   - Mark agent as failed
   - Spawn replacement agent
   - Transfer work context
   - Update state file
4. Record incident for performance review
```

### 🚨 INTEGRATION CONFLICTS
```
DETECTION (Orchestrator creating wave integration):
git checkout wave-integration
git merge effort-branch-1  # Conflict detected

CONFLICT RESOLUTION PROTOCOL:
1. Stop integration immediately
2. Document conflict details:
   - Conflicting files
   - Architectural implications
   - Risk assessment
3. Spawn architect for conflict review
4. Do NOT resolve conflicts automatically
5. May require effort rework or wave restructure
```

### 🚨 TEST COVERAGE BELOW MINIMUM
```
SW ENGINEER DETECTION:
coverage=$(go test -cover | grep "coverage:" | awk '{print $5}' | sed 's/%//')
required_coverage=75  # Phase-specific

if [ "$coverage" -lt "$required_coverage" ]; then
    echo "🚨 COVERAGE EMERGENCY: $coverage% < $required_coverage%"
fi

RECOVERY PROTOCOL:
1. Stop all non-test development
2. Identify uncovered critical paths:
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
3. Prioritize critical areas:
   - Controllers: Must be >90%
   - Webhooks: Must be >90% 
   - APIs: Must be >85%
4. Add comprehensive table-driven tests
5. Verify coverage before continuing
```

### 🚨 GRADING SYSTEM FAILURE
```
ANY AGENT DETECTION:
echo "🚨 GRADING EMERGENCY: Cannot track performance metrics"

IMMEDIATE ACTIONS:
1. Document current work state
2. Save all progress to state files
3. Report grading system status
4. Continue work but flag for manual review
5. Escalate to system administrator

MANUAL GRADING FALLBACK:
- Orchestrator: Count spawn timestamps manually
- SW Engineer: Track lines/hour and coverage
- Code Reviewer: Document review accuracy
- Architect: Track decision consistency
```

## 🛡️ STATE MACHINE EMERGENCY RECOVERY

### 🚨 LOST STATE INFORMATION
```
CONTEXT RECOVERY PROTOCOL:
1. Check for compaction marker:
   if [ -f /tmp/compaction_marker.txt ]; then
       cat /tmp/compaction_marker.txt
   fi

2. Load most recent state:
   READ: orchestrator-state.yaml
   READ: latest TODO files in todos directory

3. Determine current position:
   - What phase/wave are we in?
   - What efforts are in progress?
   - What is blocking progress?

4. Resume from safe checkpoint:
   - Load appropriate state rules
   - Verify environment correctness
   - Continue from last known good state
```

### 🚨 STATE MACHINE DEADLOCK
```
DEADLOCK DETECTION:
- Agents waiting for each other
- No state transitions occurring
- All agents blocked simultaneously

DEADLOCK RESOLUTION:
1. Map current agent states and dependencies
2. Identify circular wait conditions
3. Break deadlock by:
   - Failing least critical blocked operation
   - Reverting to previous wave if necessary  
   - Manual intervention to resolve blocking condition
4. Update state machine to prevent future deadlocks
```

### 🚨 TERMINAL STATE REACHED UNEXPECTEDLY
```
HARD_STOP ANALYSIS:
1. Verify legitimacy of HARD_STOP:
   - Architect issued STOP decision?
   - Critical system failure occurred?
   - Unrecoverable error detected?

2. If legitimate HARD_STOP:
   - Document failure cause
   - Preserve all work artifacts
   - Notify stakeholders
   - Begin failure analysis

3. If illegitimate HARD_STOP:
   - Override with manual state reset
   - Address agent malfunction
   - Resume from last valid checkpoint
   - Review state machine logic
```

## 🔍 DIAGNOSTIC PROCEDURES

### 🚨 PERFORMANCE DEGRADATION
```
ORCHESTRATOR DIAGNOSTICS:
1. Check spawn timing trends:
   grep "spawn_delta" work-log.md | tail -10
   
2. Monitor agent response times:
   for agent in sw-eng code-reviewer architect; do
       echo "Last response from $agent: $(find . -name "*$agent*" -exec stat -c %Y {} \; | sort -nr | head -1)"
   done

3. Analyze resource usage:
   ps aux | grep -E "(claude|agent)"
   df -h /tmp /workspaces
```

### 🚨 QUALITY DEGRADATION  
```
CODE REVIEWER DIAGNOSTICS:
1. Track first-try success rate trend:
   success_rate=$(echo "scale=2; $successful_first_attempts / $total_attempts" | bc)
   echo "Current success rate: $success_rate (target: >0.8)"

2. Analyze common failure patterns:
   grep "NEEDS_FIXES" review-reports/*.md | sort | uniq -c

3. Check measurement tool usage:
   grep -c "line-counter.sh" work-log.md
   grep -c "manual.*count" work-log.md  # Should be 0
```

### 🚨 INTEGRATION FAILURES
```
ARCHITECT DIAGNOSTICS:
1. Check merge readiness:
   for branch in $(git branch -r | grep effort-); do
       git checkout $branch
       lines=$(line-counter.sh -c $branch | grep "Total:" | awk '{print $2}')
       echo "$branch: $lines lines"
   done

2. Test integration feasibility:
   git checkout -b integration-test
   for branch in effort-*; do
       git merge --no-commit $branch || echo "CONFLICT: $branch"
   done

3. Analyze pattern compliance:
   find . -name "*.go" -exec grep -l "workspace" {} \; | wc -l
   find . -name "*.go" -exec grep -l "logicalcluster" {} \; | wc -l
```

## 📞 ESCALATION PROCEDURES

### Level 1 - Internal Agent Recovery
```
Scope: Single agent issues, recoverable errors
Response Time: Immediate (0-5 minutes)
Actions:
- Restart agent with preserved context
- Apply standard recovery procedures
- Continue with checkpoint recovery
```

### Level 2 - Cross-Agent Coordination Issues
```  
Scope: Multiple agents affected, state synchronization
Response Time: 5-15 minutes
Actions:
- Orchestrator intervention required
- State machine reset to last good checkpoint
- Manual coordination required
```

### Level 3 - System-Wide Failure
```
Scope: Infrastructure, tooling, or fundamental issues
Response Time: 15-30 minutes  
Actions:
- Escalate to system administrator
- Preserve all work artifacts
- Implement manual fallback procedures
- Document incident thoroughly
```

### Level 4 - Critical Architecture Emergency
```
Scope: Security, data integrity, or design fundamentals
Response Time: Immediate STOP
Actions:
- Issue immediate HARD_STOP
- Notify all stakeholders
- Begin emergency architecture review
- No work continues until resolved
```

## 📋 EMERGENCY CONTACT PROCEDURES

### Orchestrator Emergency Contacts:
```
Primary: System state management
Backup: Manual state file recovery
Escalation: Infrastructure team
```

### SW Engineer Emergency Contacts:
```
Primary: Code reviewer (for splits/reviews)
Backup: Orchestrator (for blocking issues)
Escalation: Technical lead (for environment)
```

### Code Reviewer Emergency Contacts:
```
Primary: Architect (for critical issues)
Backup: Orchestrator (for process issues)
Escalation: Quality assurance team
```

### Architect Emergency Contacts:
```
Primary: Orchestrator (for STOP decisions)
Backup: System administrator (for tooling)
Escalation: Executive leadership (for critical decisions)
```

## 🔄 POST-EMERGENCY PROCEDURES

### Incident Documentation:
```markdown
# Emergency Report: [YYYY-MM-DD-HHMMSS]

## Emergency Type
- Classification: [TECHNICAL/PROCESS/GRADING/ARCHITECTURAL]
- Severity: [1-4] 
- Duration: [start] - [end]

## Root Cause Analysis
- Initial trigger: [what happened]
- Contributing factors: [why it happened]
- Detection method: [how we found it]

## Response Actions
- Immediate actions: [what we did first]
- Recovery steps: [how we fixed it]
- Verification: [how we confirmed fix]

## Prevention Measures
- Process improvements: [what changed]
- Monitoring enhancements: [what we watch now]
- Training updates: [what agents learned]

## Lessons Learned
- What worked well: [positive observations]
- What needs improvement: [areas for enhancement]
- Recommendations: [future improvements]
```

### System Health Verification:
```bash
# Verify all systems operational after emergency
echo "=== POST-EMERGENCY SYSTEM CHECK ==="

# 1. Environment verification
echo "Directory: $(pwd)"
echo "Branch: $(git branch --show-current)"
echo "Remote status: $(git status -sb)"

# 2. Tool availability  
echo "Line counter: $(which line-counter.sh)"
echo "Git access: $(git --version)"
echo "Go tools: $(go version)"

# 3. State consistency
echo "State file: $(ls -la *state*.yaml)"
echo "TODO files: $(ls -la todos/ | wc -l) files"

# 4. Agent readiness
echo "All agents ready for normal operations"
```

---
**REMEMBER**: In emergency situations, SAFETY FIRST. Stop work, assess situation, follow procedures, document everything.