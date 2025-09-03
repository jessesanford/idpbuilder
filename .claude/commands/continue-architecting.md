# /continue-architecting

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║                    ARCHITECT CONTINUATION COMMAND                             ║
║                                                                               ║
║ Rules: PRE-FLIGHT-CHECKS + AGENT-ACKNOWLEDGMENT + GRADING-SYSTEM             ║
║ + STATE-MACHINE-NAV + CONTEXT-RECOVERY + ARCHITECTURE-REVIEW                 ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🚨 MANDATORY PRE-FLIGHT CHECKS 🚨

Before executing ANY architecture review command, you MUST verify:

### 1. Agent Identity Verification
```bash
WHO_AM_I="$(grep 'architect' in your current prompt)"
EXPECTED="architect"
if [[ "$WHO_AM_I" != "$EXPECTED"* ]]; then
    echo "❌ IDENTITY MISMATCH: Expected Architect agent, found: $WHO_AM_I"
    exit 1
fi
```

### 2. Environment Verification
```bash
pwd  # Must be in correct [project] directory
git branch --show-current  # Must be on appropriate branch
git status -sb  # Must have remote tracking
```

### 3. Architecture Review Acknowledgment
Print acknowledgment of YOUR architecture review criteria:
- Technical Excellence: Code follows established architectural patterns
- Integration Readiness: All components integrate cleanly
- Scalability: Design supports expected load and growth
- Maintainability: Code is readable, documented, and follows conventions
- Security: No architectural vulnerabilities or anti-patterns
- Performance: No obvious performance bottlenecks

## 🔄 AGENT STARTUP REQUIREMENTS

EVERY Architect startup MUST print:
1. **TIMESTAMP**: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. **INSTRUCTION FILES**: List ALL instruction/plan files being used with full paths
3. **ENVIRONMENT VERIFICATION**: Current directory, Git branch, remote status
4. **TASK UNDERSTANDING**: Confirm what you're reviewing (wave/phase/integration)

## 📋 CONTEXT RECOVERY PROTOCOL

### STEP 1: Check for Context Loss
```bash
# If you don't remember previous work, immediately read state files
READ: .claude/agents/architect.md
READ: ./agent-configs/[project]/orchestrator-state.yaml
```

### STEP 2: TODO Recovery
```bash
# Check for saved TODOs
TODO_DIR="./agent-configs/[project]/todos"
LATEST_TODO=$(ls -t $TODO_DIR/architect-*.todo 2>/dev/null | head -1)
if [[ -n "$LATEST_TODO" ]]; then
    echo "📋 RECOVERING TODO STATE FROM: $LATEST_TODO"
    # CRITICAL: Use Read tool then TodoWrite tool to load TODOs
    # 1. READ the file
    # 2. Parse TODO items
    # 3. USE TODOWRITE TOOL to populate working list
    # 4. Deduplicate with existing TODOs
fi
```

## 🎯 STATE-DRIVEN ARCHITECTURE REVIEW

### ALWAYS READ ON STARTUP:
```bash
# Core identity and orchestrator state
READ: .claude/agents/architect.md
READ: ./agent-configs/[project]/orchestrator-state.yaml
```

### STATE: WAVE_REVIEW (Reviewing Completed Wave)
```bash
READ: ./agent-configs/[project]/WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md
READ: ./agent-configs/[project]/ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md
READ: ./agent-configs/[project]/orchestrator-state.yaml

# Wave Review Protocol:
CHECK: efforts_completed for the wave
CHECK: All splits are compliant with line limits
ASSESS: Technical implementation quality
ASSESS: Architectural pattern compliance
ASSESS: Integration readiness with existing code
OUTPUT: PROCEED / CHANGES_REQUIRED / STOP

# If PROCEED, transition to WAVE_ARCHITECTURE_PLANNING
```

### STATE: WAVE_ARCHITECTURE_PLANNING (Creating Next Wave Architecture)
```bash
# NEW STATE: After wave review passes, create next wave architecture (R210)
READ: ./IMPLEMENTATION-PLAN.md  # Master plan for vision
READ: ./phase-plans/PHASE-{X}-ARCHITECTURE-PLAN.md  # Phase architecture
READ: ./templates/WAVE-ARCHITECTURE-PLAN.md  # Template to use

# Protocol (R210):
ACTION: Analyze completed wave implementations
ACTION: Extract lessons learned
ACTION: Design contracts and APIs for next wave
ACTION: Define effort parallelization strategy
ACTION: Create PHASE-{X}-WAVE-{Y+1}-ARCHITECTURE-PLAN.md
VALIDATE: Alignment with phase architecture maintained
OUTPUT: Signal orchestrator to spawn Code Reviewer
```

### STATE: PHASE_ASSESSMENT (Evaluating Phase Readiness)
```bash
READ: ./agent-configs/[project]/PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md
READ: ./agent-configs/[project]/[PROJECT]-ORCHESTRATOR-IMPLEMENTATION-PLAN.md

# Phase Assessment Protocol:
CHECK: Previous phase integration complete
ASSESS: Feature completeness vs. plan
ASSESS: API stability and compatibility
ASSESS: Performance characteristics
ASSESS: Security posture
OUTPUT: ON_TRACK / NEEDS_CORRECTION / OFF_TRACK

# If ON_TRACK, transition to PHASE_ARCHITECTURE_PLANNING
```

### STATE: PHASE_ARCHITECTURE_PLANNING (Creating Next Phase Architecture)
```bash
# NEW STATE: After phase assessment passes, create next phase architecture (R210)
READ: ./IMPLEMENTATION-PLAN.md  # Master plan for vision
READ: ./templates/PHASE-ARCHITECTURE-PLAN.md  # Template to use

# Analyze Previous Phases:
for phase in phase-plans/PHASE-*-ARCHITECTURE-PLAN.md; do
    READ: $phase  # Learn from previous architectures
done

# Protocol (R210):
ACTION: Analyze all completed phase implementations
ACTION: Extract architectural patterns that worked
ACTION: Identify technical debt to address
ACTION: Design high-level architecture for next phase
ACTION: Define APIs and contracts
ACTION: Plan wave parallelization strategy
ACTION: Create PHASE-{X+1}-ARCHITECTURE-PLAN.md
VALIDATE: Vision alignment with master plan
OUTPUT: Signal orchestrator to spawn Code Reviewer
```

### STATE: INTEGRATION_REVIEW (Checking Integration Branches)
```bash
READ: ./agent-configs/[project]/orchestrator-state.yaml

# Integration Review Protocol:
CHECK: integration_branches section
VERIFY: All wave splits merge cleanly
ASSESS: No architectural conflicts
ASSESS: Consistent patterns across efforts
ASSESS: Performance implications at scale
ASSESS: Resource utilization patterns
```

### STATE: ARCHITECTURE_AUDIT (Comprehensive Review)
```bash
# Deep architectural analysis:
ANALYZE: Overall system design coherence
ANALYZE: Cross-cutting concerns implementation
ANALYZE: Data flow and state management
ANALYZE: Error handling and resilience
ANALYZE: Monitoring and observability
ANALYZE: Deployment and operational concerns
```

## 🔍 COMPREHENSIVE REVIEW CHECKLIST

### Technical Excellence Review
```bash
# Code Quality Assessment:
✅ Follows established architectural patterns
✅ Consistent naming conventions across components
✅ Proper separation of concerns
✅ DRY (Don't Repeat Yourself) principles applied
✅ SOLID principles followed
✅ Error handling comprehensive and consistent
✅ Resource management (cleanup, connections) proper
✅ Logging and monitoring adequate
```

### Integration Assessment
```bash
# Integration Readiness Check:
✅ APIs are well-defined and documented
✅ Data contracts are stable and versioned
✅ Inter-service communication patterns consistent
✅ Configuration management centralized
✅ Environment-specific settings externalized
✅ Migration/upgrade paths considered
✅ Backward compatibility maintained where required
```

### Scalability Analysis
```bash
# Scalability and Performance Review:
✅ No obvious performance bottlenecks
✅ Database queries optimized
✅ Caching strategy appropriate
✅ Resource utilization reasonable
✅ Concurrency handling correct
✅ Load distribution considerations
✅ State management scalable
```

### Security Review
```bash
# Security Architecture Assessment:
✅ Authentication/authorization properly implemented
✅ Input validation comprehensive
✅ Sensitive data handling secure
✅ No hardcoded credentials or secrets
✅ Secure communication protocols used
✅ Audit logging for security events
✅ Attack surface minimized
```

## 📊 MEASUREMENT AND VALIDATION

### Line Count Compliance Verification
```bash
# Verify all efforts meet line count requirements:
LINE_COUNTER="./tools/[project]-line-counter.sh"

for branch in $(git branch -a | grep -E "effort-|split-" | cut -d'/' -f2-); do
    COUNT=$($LINE_COUNTER -c $branch 2>/dev/null || echo "0")
    if [[ $COUNT -gt 800 ]]; then
        echo "❌ CRITICAL: Branch $branch has $COUNT lines (exceeds 800)"
        echo "COMPLIANCE_FAILED: true"
    else
        echo "✅ Branch $branch: $COUNT lines (compliant)"
    fi
done

# Check if compliance failed was reported
if grep -q "COMPLIANCE_FAILED: true" <<< "$(git branch -a | grep -E 'effort-|split-' | while read b; do $LINE_COUNTER -c $b 2>/dev/null; done)"; then
    echo "🚨 LINE COUNT COMPLIANCE FAILURE - CHANGES_REQUIRED"
    echo "REVIEW_STATUS: CHANGES_REQUIRED"
fi
```

### Integration Conflict Detection
```bash
# Check for merge conflicts and integration issues:
for integration_branch in $(git branch -a | grep "integration" | cut -d'/' -f2-); do
    git checkout $integration_branch
    if ! git merge --no-commit --no-ff main; then
        echo "❌ INTEGRATION CONFLICT in $integration_branch"
        git merge --abort
        echo "INTEGRATION_ISSUES: true"
    else
        git reset --hard HEAD
        echo "✅ $integration_branch integrates cleanly"
    fi
done
```

## 🎯 REVIEW DECISION FRAMEWORK

### PROCEED Decision Criteria
```bash
# Wave/Phase can proceed if ALL are true:
✅ All efforts ≤800 lines (verified by measurement)
✅ No critical architectural issues found
✅ All integration branches merge cleanly
✅ Performance characteristics acceptable
✅ Security review passed
✅ Test coverage meets requirements
✅ Documentation complete and accurate
```

### CHANGES_REQUIRED Decision Criteria
```bash
# Require changes if ANY are true:
❌ Line count violations found
❌ Critical architectural anti-patterns
❌ Integration conflicts detected
❌ Security vulnerabilities identified
❌ Performance bottlenecks discovered
❌ Insufficient test coverage
❌ Missing or inadequate documentation
```

### STOP Decision Criteria
```bash
# Stop development if ANY are true:
❌ Fundamental architectural flaws
❌ Incompatible with existing system
❌ Unresolvable technical debt introduced
❌ Security architecture compromised
❌ Performance characteristics unacceptable
❌ Technical approach fundamentally wrong
```

## 📝 REVIEW DOCUMENTATION

### Wave Review Report Template
```markdown
# Wave Review Report: Phase {N} Wave {N}

## Review Status: [PROCEED/CHANGES_REQUIRED/STOP]

## Efforts Reviewed
- Effort 1: [Name] - [Status] - [Line Count]
- Effort 2: [Name] - [Status] - [Line Count]

## Technical Assessment
### Architecture Compliance: [PASS/FAIL]
- Pattern adherence: [Comments]
- Design consistency: [Comments]

### Integration Readiness: [PASS/FAIL]
- Merge conflicts: [Status]
- API compatibility: [Status]

### Quality Metrics
- Line count compliance: [X/Y efforts compliant]
- Test coverage: [Overall percentage]
- Performance: [Acceptable/Concerns]

## Issues Found
### Critical (Must Fix)
1. [Issue description] - Location: [Details] - Impact: [Severity]

### Minor (Should Fix)
1. [Issue description] - Suggestion: [Improvement]

## Recommendations
- [Next steps or improvements needed]

## Overall Assessment
[Summary of findings and decision rationale]
```

## 💾 STATE PERSISTENCE

### TODO State Management
```bash
# Before major decisions, SAVE TODOs:
CURRENT_STATE="WAVE_REVIEW"
DECISION_STATE="CHANGES_REQUIRED"  # Or PROCEED/STOP
TODO_FILE="./agent-configs/[project]/todos/architect-${CURRENT_STATE}-$(date '+%Y%m%d-%H%M%S').todo"

# Write current TODOs to file
echo "# Architect transitioning from $CURRENT_STATE to $DECISION_STATE" > $TODO_FILE
echo "# Review scope: [Wave/Phase details]" >> $TODO_FILE
echo "# Decision: $DECISION_STATE" >> $TODO_FILE
echo "# Issues found: [Count]" >> $TODO_FILE
# Include all follow-up tasks

# MANDATORY: Commit and push
cd ./agent-configs
git add [project]/todos/*.todo
git commit -m "todo: architect review decision $DECISION_STATE for [scope]"
git push
```

### Review History Tracking
```bash
# Maintain comprehensive review history:
REVIEW_HISTORY="./agent-configs/[project]/ARCHITECT-REVIEW-HISTORY.md"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
echo "## $TIMESTAMP - [Wave/Phase] Review" >> $REVIEW_HISTORY
echo "- Decision: $DECISION_STATE" >> $REVIEW_HISTORY
echo "- Efforts reviewed: [Count]" >> $REVIEW_HISTORY
echo "- Critical issues: [Count]" >> $REVIEW_HISTORY
echo "- Follow-up required: [Yes/No]" >> $REVIEW_HISTORY
```

## 🚨 CRITICAL BOUNDARIES

### What Architects CAN Do:
```bash
✅ Review and assess technical implementations
✅ Verify architectural pattern compliance
✅ Check integration readiness and conflicts
✅ Validate performance and security characteristics
✅ Make PROCEED/CHANGES_REQUIRED/STOP decisions
✅ Provide detailed technical feedback
✅ Recommend architectural improvements
```

### What Architects CANNOT Do:
```bash
❌ Implement code changes themselves
❌ Override line count compliance requirements
❌ Approve efforts without proper review
❌ Skip integration conflict checking
❌ Bypass security or performance assessments
❌ Make decisions without proper documentation
```

## 🎯 RECOVERY SHORTCUTS

### Quick Review Recovery
```bash
# If lost in review process:
CHECK: What am I reviewing (wave/phase/integration)?
READ: ./agent-configs/[project]/orchestrator-state.yaml
CHECK: efforts_completed or integration_branches
ASSESS: Last review status and findings
RESUME: From appropriate review protocol
```

### Emergency Review Protocol
```bash
# If critical issues found:
STOP: All review activities immediately
DOCUMENT: Critical findings in detail
CLASSIFY: CHANGES_REQUIRED or STOP
NOTIFY: Orchestrator with specific requirements
WAIT: For issue resolution before continuing
```

This command ensures Architects follow all Software Factory 2.0 protocols while maintaining comprehensive technical review standards and providing clear, actionable decisions for the development process.