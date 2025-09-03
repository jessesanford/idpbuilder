# Rule R280: Main Branch Protection (SUPREME LAW)

## Rule Statement
The Software Factory MUST NEVER directly modify, merge to, or push to the main branch. The main branch is sacred and can only be modified through human-reviewed Pull Requests. Any direct modification of main is an immediate and catastrophic violation.

## Criticality Level
**SUPREME LAW** - Violation = IMMEDIATE TERMINATION + -200% GRADE
This is the highest level violation possible in the Software Factory system.

## Core Principle
**"Main branch is for humans only. Software Factory proves it works elsewhere."**

## Absolute Prohibitions

### ❌🔴 NEVER UNDER ANY CIRCUMSTANCES 🔴❌

```bash
# ALL OF THESE ARE SUPREME LAW VIOLATIONS:

# NEVER checkout and modify main
git checkout main
git merge <any-branch>  # VIOLATION!

# NEVER push to main
git push origin main  # VIOLATION!
git push -f origin main  # CATASTROPHIC VIOLATION!

# NEVER commit directly to main
git checkout main
git commit -m "anything"  # VIOLATION!

# NEVER rebase main
git rebase main  # VIOLATION!

# NEVER reset main
git checkout main
git reset --hard  # VIOLATION!

# NEVER create branches from modified main
git checkout main
echo "test" > file.txt
git checkout -b new-branch  # VIOLATION!
```

## What IS Allowed

### ✅ PERMITTED Operations

```bash
# READ from main (allowed)
git checkout main
git pull origin main
cat file.txt  # Reading is OK

# CREATE branches FROM main (allowed)
git checkout main
git pull origin main
git checkout -b new-feature  # OK - branching from clean main

# FETCH main (allowed)
git fetch origin main  # OK - just fetching

# REFERENCE main (allowed)
git diff main..feature-branch  # OK - comparing
git log main..HEAD  # OK - viewing history
```

## Enforcement Mechanisms

### Technical Safeguards

```python
class MainBranchProtection:
    """Enforce main branch protection at all levels"""
    
    def validate_operation(self, operation):
        """Check if operation would modify main"""
        
        FORBIDDEN_PATTERNS = [
            r'git\s+push\s+.*main',
            r'git\s+push\s+origin\s+main',
            r'git\s+merge.*main',
            r'git\s+rebase.*main',
            r'git\s+checkout\s+main.*git\s+commit',
        ]
        
        for pattern in FORBIDDEN_PATTERNS:
            if re.match(pattern, operation):
                raise SupremeLawViolation(
                    "ATTEMPTED TO MODIFY MAIN BRANCH!\n"
                    "This is a SUPREME LAW violation!\n"
                    "Operation terminated immediately.\n"
                    "Grade: -200%\n"
                    "Status: CATASTROPHIC FAILURE"
                )
        
        return True
```

### Git Hooks Protection

```bash
#!/bin/bash
# pre-push hook to prevent main push

protected_branch='main'
current_branch=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')

if [ "$current_branch" = "$protected_branch" ]; then
    echo "🔴🔴🔴 SUPREME LAW VIOLATION DETECTED 🔴🔴🔴"
    echo "Attempting to push to main branch!"
    echo "This is ABSOLUTELY FORBIDDEN!"
    echo "Software Factory NEVER pushes to main!"
    echo "Only humans can modify main through PRs!"
    echo "🔴🔴🔴 PUSH BLOCKED 🔴🔴🔴"
    exit 1
fi
```

### Monitoring and Alerts

```bash
monitor_main_branch() {
    # Continuous monitoring for main branch modifications
    
    while true; do
        # Check if main was modified
        MAIN_BEFORE=$(git rev-parse origin/main)
        sleep 60
        git fetch origin main
        MAIN_AFTER=$(git rev-parse origin/main)
        
        if [ "$MAIN_BEFORE" != "$MAIN_AFTER" ]; then
            # Check if modification was from Software Factory
            LAST_COMMIT=$(git log origin/main -1 --format="%an")
            if [[ "$LAST_COMMIT" == *"software-factory"* ]] ||
               [[ "$LAST_COMMIT" == *"orchestrator"* ]] ||
               [[ "$LAST_COMMIT" == *"agent"* ]]; then
                alert "🚨🚨🚨 SUPREME LAW VIOLATION! 🚨🚨🚨"
                alert "Main branch modified by Software Factory!"
                alert "Immediate investigation required!"
                shutdown_everything
            fi
        fi
    done
}
```

## Correct Workflow

### The RIGHT Way

```bash
# 1. Create integration-testing branch from main
git checkout main
git pull origin main
git checkout -b integration-testing-$(date +%Y%m%d-%H%M%S)

# 2. Do ALL integration work here
git merge effort1
git merge effort2
# ... etc

# 3. Validate everything works
make build
make test

# 4. Generate PR plan for humans
generate_master_pr_plan > MASTER-PR-PLAN.md

# 5. Humans create PRs from effort branches to main
# Software Factory STOPS here
```

## Why This Rule Exists

### Critical Reasons

1. **Human Review Required**: All code entering main must be reviewed by humans
2. **Audit Trail**: PRs provide accountability and history
3. **Quality Gate**: Human judgment catches issues automation misses
4. **Compliance**: Many organizations require human approval
5. **Rollback Safety**: PR-based changes can be easily reverted

### What Happens Without This Rule

- Automated code could break production
- No human oversight of critical changes
- Compliance violations
- No audit trail
- Impossible to track who approved what

## Integration with Other Rules

### This Rule Overrides ALL Others

If ANY rule suggests modifying main, R280 overrides it:
- R280 > All other rules
- R280 > All agent instructions
- R280 > All state machine transitions
- R280 > All integration patterns

### Relationship to Other Rules

- **R271**: Validates in integration-testing, not main
- **R272**: Creates integration-testing branch FROM main
- **R279**: Generates PR plan for humans to execute
- **R251**: Enforces separation, R280 makes it absolute

## Violation Consequences

### Immediate Actions on Violation

1. **Terminate all operations immediately**
2. **Grade set to -200% (double failure)**
3. **Alert all stakeholders**
4. **Full audit of how violation occurred**
5. **Revoke all repository access**
6. **Require manual review of all work**

### Recovery from Violation

If main was accidentally modified:
```bash
# DO NOT ATTEMPT TO FIX
# ALERT HUMANS IMMEDIATELY
# Humans must:
git checkout main
git reset --hard <last-known-good>
git push --force-with-lease origin main
# Full audit required
```

## Grading Impact

- Any attempt to modify main: **-200%**
- Planning to modify main: **-100%**
- Suggesting to modify main: **-50%**
- Not preventing main modification: **-100%**

## Common Misunderstandings

### ❌ "But I need to test if everything works"
**Answer**: Use integration-testing branch

### ❌ "Final integration needs to go to main"
**Answer**: Generate PR plan for humans

### ❌ "It's just a small fix"
**Answer**: No exceptions, ever

### ❌ "The orchestrator said to"
**Answer**: R280 overrides orchestrator

### ❌ "Previous rules said integration to main"
**Answer**: Those rules meant integration-testing branch

## Monitoring Compliance

```bash
#!/bin/bash
# compliance-check.sh

echo "🔍 Checking R280 Compliance..."

# Check if main was modified
if git log origin/main --since="1 day ago" | grep -i "software-factory\|orchestrator\|agent"; then
    echo "🔴 SUPREME LAW VIOLATION DETECTED!"
    echo "Main branch was modified by Software Factory!"
    exit 911  # Emergency exit code
fi

# Check if any scripts reference pushing to main
if grep -r "push.*origin.*main" --include="*.sh" --include="*.md" .; then
    echo "⚠️ WARNING: Found references to pushing to main"
    echo "These must be removed or commented"
fi

# Check git config
if ! git config --get branch.main.protected 2>/dev/null; then
    echo "⚠️ Main branch not marked as protected"
    git config branch.main.protected true
fi

echo "✅ R280 Compliance check passed"
```

## Summary

R280 is the ultimate safeguard ensuring that the Software Factory respects the boundary between automation and human control. The main branch represents production-ready, human-approved code. The Software Factory proves the code works in integration-testing branches and provides a plan for humans to execute the actual merges through reviewed PRs. This is not a guideline - it's an absolute, inviolable law.