# 🚨🚨🚨 RULE R052: Agent Spawning Protocol

## Classification
- **Category**: Agent Management
- **Criticality Level**: 🚨🚨🚨 CRITICAL
- **Enforcement**: MANDATORY for all agent spawning
- **Penalty**: -50% for incomplete context, -100% for missing deliverables

## The Rule

**When spawning ANY agent, you MUST provide complete context, startup requirements, deliverables, and size limits.**

## Requirements

### 1. Complete Context Package

Every spawned agent MUST receive:

```markdown
PURPOSE: [Clear description of what agent should accomplish]
TARGET_DIRECTORY: [Exact path where agent should work]
EXPECTED_BRANCH: [Git branch agent should be on]
REQUIREMENTS: [Specific requirements for this task]
STARTUP_REQUIREMENTS: [What agent must do first]
DELIVERABLES: [What must be complete when done]
```

### 2. Mandatory Startup Requirements

**EVERY agent spawn message MUST include:**

```markdown
STARTUP REQUIREMENTS:
1. Print: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. Verify pwd matches TARGET_DIRECTORY
3. Verify branch matches EXPECTED_BRANCH
4. Read and acknowledge state-specific rules
5. Load TODO state if recovering from compaction
```

### 3. Directory Navigation Instructions

**For SW Engineers and Code Reviewers:**

```markdown
🔴🔴🔴 CRITICAL: YOU WILL NOT BE IN THE RIGHT DIRECTORY! 🔴🔴🔴
YOU MUST NAVIGATE TO YOUR EFFORT DIRECTORY IMMEDIATELY!

YOUR MANDATORY FIRST ACTIONS:
1. Echo your current directory: pwd
2. Navigate to effort directory: cd [TARGET_DIRECTORY]
3. Verify you're now in correct directory: pwd
4. Verify branch: git branch --show-current
5. If directory doesn't exist or branch is wrong:
   - STOP IMMEDIATELY
   - Report: "❌ ENVIRONMENT ERROR: Directory or branch incorrect"
   - Request orchestrator correction
```

### 4. Size Limits and Constraints

**ALWAYS specify:**
- Maximum lines of code allowed
- Test coverage requirements
- Documentation requirements
- Performance constraints
- Security requirements

### 5. Deliverables Specification

**Clear, measurable deliverables:**

```markdown
DELIVERABLES:
- [ ] Implementation complete per IMPLEMENTATION-PLAN.md
- [ ] Tests passing at {X}% coverage minimum
- [ ] Size under {Y} lines limit
- [ ] Work log updated with progress
- [ ] Code committed and pushed to remote
- [ ] No uncommitted changes remaining
```

## Spawn Message Template

### For SW Engineers:

```markdown
# Task: SW Engineer Implementation
PURPOSE: Implement {effort_id} - {effort_name}

TARGET_DIRECTORY: /efforts/phase{X}/wave{Y}/{effort-name}
EXPECTED_BRANCH: {project-prefix}/phase{X}/wave{Y}/{effort-name}

[Include mandatory directory navigation instructions]

REQUIREMENTS:
- Follow IMPLEMENTATION-PLAN.md exactly
- Size limit: {limit} lines
- Test coverage: {X}% minimum
- Update work-log.md every checkpoint

STARTUP REQUIREMENTS:
[Include mandatory startup requirements]

DELIVERABLES:
[Include clear deliverables list]
```

### For Code Reviewers:

```markdown
# Task: Code Review
PURPOSE: Review {effort_id} implementation

TARGET_DIRECTORY: /efforts/phase{X}/wave{Y}/{effort-name}
EXPECTED_BRANCH: {project-prefix}/phase{X}/wave{Y}/{effort-name}

REQUIREMENTS:
- Verify implementation matches plan
- Check size compliance
- Validate test coverage
- Create review report

STARTUP REQUIREMENTS:
[Include mandatory startup requirements]

DELIVERABLES:
- [ ] CODE-REVIEW-REPORT.md created
- [ ] Size violations identified if any
- [ ] Quality issues documented
- [ ] Recommendations provided
```

## Parallelization Context

When spawning multiple agents in parallel:

```markdown
PARALLELIZATION CONTEXT:
- You are being spawned as part of a parallel group
- Other agents working on: [list other efforts]
- Dependencies: [none/list any dependencies]
- Coordination: Work independently unless specified
```

## Error Handling Instructions

**ALWAYS include:**

```markdown
ERROR HANDLING:
- If environment incorrect: STOP and report immediately
- If plan missing: Request from orchestrator
- If dependencies blocked: Report and wait for resolution
- If size exceeded: Stop and request split plan
```

## Common Violations

### ❌ Incomplete Context
```markdown
# WRONG - Missing critical information
Task: Implement feature
Go implement the user authentication feature.
```

### ❌ No Startup Requirements
```markdown
# WRONG - No startup verification
Task: SW Engineer
Just start coding in the effort directory.
```

### ❌ Vague Deliverables
```markdown
# WRONG - Unmeasurable deliverables
DELIVERABLES:
- Make it work
- Write some tests
- Push when done
```

### ✅ CORRECT Complete Spawn
```markdown
# Task: SW Engineer Implementation
PURPOSE: Implement effort-001-user-auth - User Authentication Module

TARGET_DIRECTORY: /efforts/phase1/wave1/effort-001-user-auth
EXPECTED_BRANCH: myproject/phase1/wave1/effort-001-user-auth

🔴🔴🔴 CRITICAL: YOU WILL NOT BE IN THE RIGHT DIRECTORY! 🔴🔴🔴
[Full directory navigation instructions]

REQUIREMENTS:
- Follow IMPLEMENTATION-PLAN.md exactly
- Size limit: 700 lines maximum
- Test coverage: 80% minimum
- Update work-log.md every 2 hours

STARTUP REQUIREMENTS:
1. Print: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. Verify pwd matches TARGET_DIRECTORY
3. Verify branch matches EXPECTED_BRANCH
4. Read and acknowledge rules R007, R013, R060
5. Check for TODO recovery file

DELIVERABLES:
- [ ] All features from IMPLEMENTATION-PLAN.md implemented
- [ ] Tests passing with 80% coverage
- [ ] Size under 700 lines (verified with line-counter.sh)
- [ ] work-log.md updated with all progress
- [ ] All code committed and pushed
- [ ] No uncommitted changes

ERROR HANDLING:
- If directory missing: STOP and report
- If plan missing: Request from orchestrator
- If size exceeded: Stop at 650 lines and request split
```

## Verification

After spawning, orchestrator should verify:
1. Agent acknowledged startup
2. Agent printed timestamp
3. Agent in correct directory
4. Agent on correct branch
5. Agent began work on correct task

## Grading Criteria

- ✅ **+25%**: Complete context provided
- ✅ **+25%**: Startup requirements included
- ✅ **+25%**: Clear deliverables specified
- ✅ **+25%**: Error handling instructions included

## Related Rules

- R151: Parallel spawning timestamp requirements
- R197: One agent per effort
- R208: CD before spawn requirement
- R255: Post-agent work verification

## Remember

**"A well-spawned agent is a successful agent"**
**"Complete context prevents confusion"**
**"Clear deliverables ensure completion"**

Incomplete spawning leads to failed implementations, wasted time, and grading penalties!