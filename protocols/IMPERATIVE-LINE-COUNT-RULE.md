# IMPERATIVE LINE COUNT RULE

## 🚨 ABSOLUTE REQUIREMENT - NO EXCEPTIONS 🚨

**EVERY EFFORT MUST BE UNDER THE CONFIGURED SIZE LIMIT**

This is not a guideline. This is not a suggestion. This is an IMPERATIVE RULE.

## The Rule

```
IF effort_size > MAX_LIMIT:
    STOP IMMEDIATELY
    IMPLEMENT SPLIT PROTOCOL
    NO EXCEPTIONS
```

## Measurement Protocol

### ONLY Valid Method
```bash
/home/vscode/workspaces/idpbuilder/tools/line-counter.sh -c {branch}
```

### INVALID Methods ❌
- Manual counting
- wc -l
- git diff --stat
- IDE line counters
- Any other method

## What Counts

### INCLUDED in Count ✅
- Hand-written source code
- Hand-written tests
- Hand-written documentation in code files
- Configuration that is hand-maintained

### EXCLUDED from Count ❌
- Generated code (any file with 'generated' in name)
- Vendor/dependencies
- Build artifacts
- Test fixtures/data
- Documentation files (.md, .rst, .txt)
- Binary files
- Images

## Thresholds

```yaml
warning_threshold: 700 lines  # Plan to complete
error_threshold: 800 lines    # MUST STOP
```

## Continuous Measurement

### Required Checkpoints
1. **Before Starting**: Baseline measurement
2. **Every 200 lines**: Progress check
3. **Every logical unit**: Feature complete check
4. **Before commit**: Final measurement
5. **After any major addition**: Sanity check

### SW Engineer Responsibility
```bash
# After each logical change
line-counter.sh -c $(git branch --show-current)

# If approaching warning
if [ $lines -gt 700 ]; then
    echo "WARNING: Approaching limit"
    echo "Focus on completing current work"
fi

# If over limit
if [ $lines -gt 800 ]; then
    echo "ERROR: OVER LIMIT - STOPPING"
    exit 1
fi
```

## Split Protocol Trigger

When limit is exceeded:

1. **IMMEDIATE STOP** - No more code
2. **Document stopping point** - In work-log.md
3. **Report to orchestrator** - With exact count
4. **Wait for split plan** - From Code Reviewer
5. **Execute splits sequentially** - Never parallel

## Enforcement Hierarchy

### Level 1: SW Engineer
- Self-monitor continuously
- Stop before hitting limit
- Report approaching limit

### Level 2: Code Reviewer
- Verify measurement
- Reject if over limit
- Design split if needed

### Level 3: Orchestrator
- Block progress if over limit
- Enforce split protocol
- Update state file

### Level 4: System
- line-counter.sh returns exit 1
- Commits blocked
- CI/CD fails

## Common Violations and Consequences

| Violation | Consequence | Recovery |
|-----------|-------------|----------|
| "Just 50 more lines" | Review rejection | Split required |
| "Tests don't count" | They do count | Split required |
| "Generated doesn't count" | Correct, but verify | Re-measure |
| "I'll split later" | Work lost | Split now |
| "Review anyway" | Review refused | Split first |

## The Psychology

### Why This Rule Exists
1. **Reviewability**: Large PRs hide bugs
2. **Cognitive Load**: Humans can't review 2000 lines effectively
3. **Quality**: Smaller changes = better review = fewer bugs
4. **Velocity**: Small PRs merge faster
5. **Rollback**: Smaller changes easier to revert

### The Temptation
"I'm almost done, just a bit more..."

### The Reality
- Over-limit efforts get rejected
- Time is wasted
- Splits are harder after the fact
- Quality suffers

## Mathematical Certainty

Given:
- Human review capacity ≈ 400-800 lines effectively
- Bug detection rate drops >50% after 800 lines
- Review time increases exponentially with size

Therefore:
- Limit MUST be enforced
- No exceptions can be made
- System integrity depends on this

## Implementation in State Machine

```python
def measure_size(effort):
    size = run_line_counter(effort.branch)
    
    if size > ERROR_THRESHOLD:
        # MANDATORY STATE TRANSITION
        return "CREATE_SPLIT_PLAN"
    elif size > WARNING_THRESHOLD:
        log.warning("Approaching limit")
        return "CONTINUE_CAUTIOUSLY"
    else:
        return "CONTINUE"
```

## Audit Trail

Every measurement is logged:
```yaml
effort_metrics:
  - timestamp: "2025-01-21T10:00:00Z"
    effort: "E1.2.3"
    measurement: 450
    status: "OK"
  - timestamp: "2025-01-21T11:00:00Z"
    effort: "E1.2.3"
    measurement: 720
    status: "WARNING"
  - timestamp: "2025-01-21T12:00:00Z"
    effort: "E1.2.3"
    measurement: 850
    status: "ERROR - SPLIT REQUIRED"
```

## Final Word

**This rule is not negotiable.**

Every effort. Every time. No exceptions.

The line counter is the arbiter of truth.
The limit is absolute.
The protocol is mandatory.

Trust the process. The system works when the rules are followed.