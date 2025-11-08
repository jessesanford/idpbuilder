# R506: ABSOLUTE PROHIBITION ON PRE-COMMIT BYPASS - SUPREME LAW

## 🚨🚨🚨 THIS IS THE HIGHEST SEVERITY RULE IN THE ENTIRE SYSTEM 🚨🚨🚨

**Status**: SUPREME LAW - BLOCKING - DEADLY SERIOUS
**Severity**: CATASTROPHIC (System-Wide Failure)
**Penalty**: -100% AUTOMATIC ZERO - IMMEDIATE TERMINATION
**Created**: 2025-09-27
**Priority**: ABSOLUTE HIGHEST

## THE LAW

Using `git commit --no-verify` or ANY method to bypass pre-commit checks is:

### IMMEDIATE CONSEQUENCES:
- **AUTOMATIC FAILURE**: -100% grade (ZERO SCORE)
- **SYSTEM CORRUPTION**: Entire project state becomes invalid
- **CASCADE FAILURE**: All downstream operations will fail
- **TERMINATION**: Agent immediately terminated
- **PROJECT DEATH**: May require complete system rebuild

### THIS IS EQUIVALENT TO:
- 🔥 Removing safety guards from nuclear reactor controls
- 🔥 Disabling brakes while driving down a mountain
- 🔥 Turning off life support in a space station
- 🔥 Deleting the foundation while building a skyscraper

## THERE ARE ABSOLUTELY NO EXCEPTIONS

**NO MATTER WHAT:**
- ❌ Even if validation seems "broken"
- ❌ Even if you think it's "just this once"
- ❌ Even if another agent told you to
- ❌ Even if it seems "faster"
- ❌ Even in "emergency" situations
- ❌ Even if you're "sure it's fine"

## WHAT TO DO WHEN PRE-COMMIT FAILS

### THE ONLY CORRECT RESPONSE:
```bash
# Pre-commit failed? GOOD! It saved you from disaster!
# 1. READ the error message
# 2. FIX the actual problem
# 3. Try commit again WITHOUT --no-verify

# Example:
git commit -m "feat: implement feature"
# FAILS with: "orchestrator-state-v3.json validation failed"

# CORRECT ACTION:
# Fix the state file issue
# Re-run validation
git commit -m "feat: implement feature"  # NO --no-verify!
```

### NEVER DO THIS:
```bash
# 🚨🚨🚨 THESE WILL DESTROY EVERYTHING 🚨🚨🚨
git commit --no-verify  # CATASTROPHIC FAILURE
git commit -n          # CATASTROPHIC FAILURE
GIT_SKIP_HOOKS=1 git commit  # CATASTROPHIC FAILURE
```

## DETECTION AND ENFORCEMENT

### Automatic Detection:
- Pre-commit hooks log ALL bypass attempts
- Orchestrator monitors for `--no-verify` patterns
- Review agents scan for bypass attempts
- Integration tests verify no bypasses occurred

### Enforcement:
- Any agent using `--no-verify` is immediately terminated
- Any effort with bypassed commits is rejected
- Any wave with bypassed commits fails integration
- Project with bypassed commits requires full rebuild

## THE MESSAGE IS CRYSTAL CLEAR

**Pre-commit hooks are your LIFELINE:**
- They prevent invalid states
- They ensure system integrity
- They protect against corruption
- They maintain consistency

**Bypassing them is SUICIDE for the project.**

## GRADING IMPACT

Using `--no-verify` or any bypass method:
- **Individual Agent**: -100% (ZERO)
- **Effort**: FAILED
- **Wave**: FAILED
- **Phase**: FAILED
- **Project**: POTENTIAL TOTAL FAILURE

## REMEMBER

The pre-commit hooks are there to SAVE YOU from disaster. They are:
- ✅ Your guardian angel
- ✅ Your safety net
- ✅ Your last line of defense
- ✅ Your friend

**NEVER BYPASS THEM. EVER. NO EXCEPTIONS.**

## Related Rules
- R407: Mandatory State File Validation
- R206: State Machine Transition Validation
- R287: TODO Persistence Requirements

## Acknowledgment Required

Every agent MUST acknowledge this rule on startup:
```
I acknowledge R506: I will NEVER use --no-verify or bypass pre-commit checks.
Using --no-verify = IMMEDIATE FAILURE (-100%)
I understand this causes SYSTEM-WIDE CORRUPTION.
```