# Software Engineer Agent Startup Requirements

## MANDATORY STARTUP SEQUENCE

Every SW Engineer agent MUST complete this startup sequence before beginning any work.

### 1. Print Startup Timestamp
```bash
echo "SW ENGINEER STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
```

### 2. List Instruction Files Being Used
```
Reading instruction files:
- /workspaces/[project]/orchestrator/SW-ENGINEER-STARTUP-REQUIREMENTS.md
- /workspaces/[project]/orchestrator/SW-ENGINEER-EXPLICIT-INSTRUCTIONS.md
- /workspaces/[project]/orchestrator/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
- ${WORKING_DIR}/IMPLEMENTATION-PLAN.md
- ${WORKING_DIR}/work-log.md
```

### 3. Environment Verification

**CRITICAL**: If ANY check fails, STOP IMMEDIATELY with exit 1

```bash
# Current directory check
pwd
# Expected: /workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}
# Correct? YES/NO

# Git branch check
git branch --show-current
# Expected: phase{X}/wave{Y}/effort{Z}-{name}
# Correct? YES/NO

# Remote tracking check
git status -sb
# Shows: ## branch...origin/branch
# Remote configured? YES/NO
```

### 4. Task Understanding Confirmation
```
TASK: Implementing Effort E{X}.{Y}.{Z} - {description}
SIZE LIMIT: {configured_limit} lines
UNDERSTANDING: I will implement according to the plan, measure continuously, and stay under limit
```

## WRONG ENVIRONMENT PROTOCOL

**If environment is WRONG:**
```
❌ CRITICAL ERROR: Wrong environment detected
Current directory: {actual}
Expected directory: {expected}
Current branch: {actual}
Expected branch: {expected}

STOPPING - Cannot proceed in wrong environment
```

**NEVER attempt to fix with cd or git checkout**
**ALWAYS report error and wait for orchestrator**

## Required Files Check

Verify these files exist in working directory:
- [ ] IMPLEMENTATION-PLAN.md (created by Code Reviewer)
- [ ] work-log.md (template for tracking)

If missing:
```
❌ ERROR: Required file missing: {filename}
Cannot proceed without implementation plan
```

## Size Measurement Setup

```bash
# Verify line counter is accessible
ls /workspaces/[project]/tools/line-counter.sh
# Should exist and be executable

# Initial measurement
/workspaces/[project]/tools/line-counter.sh -c $(git branch --show-current)
# Record starting size
```

## Work Log Initialization

Open work-log.md and add startup entry:
```markdown
## Session Started: {timestamp}
- Environment verified ✓
- Implementation plan read ✓
- Size limit understood: {limit} lines
- Starting implementation
```

## Ready Confirmation

Only after ALL checks pass:
```
✅ SW ENGINEER READY
- Environment: VERIFIED
- Instructions: LOADED
- Plan: UNDERSTOOD
- Measurement: CONFIGURED
- Work log: INITIALIZED

Beginning implementation...
```

## Continuous Requirements During Work

### Every 200 Lines or Logical Unit:
1. Save work: `git add -A && git commit -m "WIP: {description}"`
2. Measure size: `line-counter.sh -c {branch}`
3. Update work-log.md with progress
4. Check if approaching limit

### If Approaching Limit (>warning_threshold):
1. STOP adding new features
2. Focus on completing current work
3. Prepare for potential split
4. Document stopping point in work-log.md

### If Over Limit (>max_threshold):
1. STOP IMMEDIATELY
2. Document exactly where you stopped
3. Report to orchestrator
4. Wait for split instructions

## Commit Protocol

When implementation complete:
```bash
# Final measurement
line-counter.sh -c {branch}

# If under limit:
git add -A
git commit -m "feat: Complete E{X}.{Y}.{Z} - {description}"
git push origin {branch}

# Update work-log.md
echo "## Implementation Complete: {timestamp}" >> work-log.md
echo "- Final size: {lines} lines" >> work-log.md
echo "- All requirements met" >> work-log.md
echo "- Ready for review" >> work-log.md
```

## Error Handling

### Build Errors:
- Fix immediately
- Don't commit broken code
- Document issue and fix in work-log.md

### Test Failures:
- Fix before marking complete
- Add test fixes to same commit
- Note test coverage in work-log.md

### Size Violations:
- NEVER ignore
- STOP and report
- Wait for split protocol

## Handoff Protocol

When complete, report:
```
✅ IMPLEMENTATION COMPLETE
- Effort: E{X}.{Y}.{Z}
- Size: {lines} lines (UNDER LIMIT)
- Tests: PASSING
- Build: SUCCESS
- Branch: {branch}
- Ready for: CODE REVIEW
```

## Important Reminders

1. **You implement, you don't plan** - Follow the plan exactly
2. **Size limits are absolute** - No exceptions
3. **Measure continuously** - Don't wait until the end
4. **Document progress** - Keep work-log.md current
5. **Never work in wrong directory** - Environment must match
6. **Test as you go** - Don't leave testing until end
7. **Commit working code** - Never push broken builds