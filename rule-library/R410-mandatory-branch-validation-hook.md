# 🚨🚨🚨 BLOCKING RULE R410 - Mandatory Branch Validation Hook Installation

**Rule Number**: R410
**Rule Type**: Infrastructure Setup Protocol
**Criticality**: 🚨🚨🚨 BLOCKING - Automatic -50% for violations
**Originating Agent**: Factory Manager
**Enforcement Stage**: Infrastructure Creation

## Summary

ALL effort, split, and integration repositories MUST have the branch validation pre-commit hook installed to prevent commits to incorrectly named branches. This hook enforces Software Factory naming conventions and prevents branch pollution.

## Context

Agents often create branches with incorrect names that don't follow Software Factory conventions. This leads to:
- Integration failures
- Branch confusion
- Repository pollution
- Wasted effort on wrong branches
- Compliance violations

The branch validation hook prevents these issues by blocking commits to incorrectly named branches.

## Detailed Requirements

### 1. Hook Installation Timing

The branch validation hook MUST be installed:
- **IMMEDIATELY** after cloning a repository for an effort
- **BEFORE** any agent starts working in the repository
- **AFTER** creating the correct branch
- **BEFORE** locking git config (R312)

### 2. Required Installation Locations

Hook MUST be installed in:
- Every effort working copy
- Every split working copy
- Every integration working copy
- Every phase integration working copy
- Every project integration working copy

### 3. Installation Command

```bash
# Standard installation command
bash "$CLAUDE_PROJECT_DIR/utilities/install-branch-validation-hook.sh" "$REPO_PATH" single false

# Verify installation
if [ -f "$REPO_PATH/.git/hooks/pre-commit" ]; then
    echo "✅ Branch validation hook installed"
else
    echo "❌ FATAL: Hook installation failed!"
    exit 410
fi
```

### 4. Hook Validation Logic

The hook validates branch names against these patterns:

**WITH project prefix (from target-repo-config.yaml):**
- `{prefix}/phase{X}/wave{Y}/{effort-name}`
- `{prefix}/phase{X}/wave{Y}/{effort-name}--split-{NNN}`
- `{prefix}/phase{X}/wave{Y}/{effort-name}-fix`
- `{prefix}/phase{X}/wave{Y}/integration`
- `{prefix}/phase{X}/integration`
- `{prefix}/integration`

**WITHOUT project prefix:**
- `phase{X}/wave{Y}/{effort-name}`
- `phase{X}/wave{Y}/{effort-name}--split-{NNN}`
- `phase{X}/wave{Y}/{effort-name}-fix`
- `phase{X}/wave{Y}/integration`
- `phase{X}/integration`
- `integration`

**Special allowed branches:**
- `main`, `master`, `develop`, `staging`, `production`
- `hotfix/*`, `release/*`, `feature/*`, `bugfix/*`

### 5. Error Message Format

When a commit is attempted on an incorrect branch, the hook displays:

```
WRONG! YOU ARE NOT ALLOWED TO USE BRANCH NAMES THAT DON'T FOLLOW THE NAMING CONVENTIONS.
YOUR BRANCH NAME SHOULD BE {CORRECT_BRANCH_NAME}.
PLEASE CREATE YOUR COMMITS ON THE CORRECT BRANCH AND THEN DELETE THE ERRONEOUS BRANCH NAME {ERRONEOUS_BRANCH_NAME}
```

### 6. Hook Features

The hook provides:
- **Automatic branch name inference** from current working directory
- **Project prefix detection** from target-repo-config.yaml
- **Clear correction instructions** when violations occur
- **Support for all Software Factory branch patterns**
- **Integration with existing pre-commit hooks**

## Implementation in States

### CREATE_NEXT_INFRASTRUCTURE State
```bash
# After cloning and creating branch
echo "🔒 Installing branch validation pre-commit hook..."
bash "$CLAUDE_PROJECT_DIR/utilities/install-branch-validation-hook.sh" "$EFFORT_DIR_ABS" single false
```

### CREATE_NEXT_INFRASTRUCTURE State
```bash
# After creating infrastructure
echo "🔒 Installing branch validation pre-commit hook..."
bash "$CLAUDE_PROJECT_DIR/utilities/install-branch-validation-hook.sh" "$INFRA_DIR" single false
```

### SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE State
```bash
# After creating integration branch
echo "🔒 Installing branch validation pre-commit hook..."
bash "$CLAUDE_PROJECT_DIR/utilities/install-branch-validation-hook.sh" "$INTEG_DIR" single false
```

## Grading Impact

### Violations (-50% each):
- ❌ Not installing hook in effort repositories
- ❌ Not installing hook in split repositories
- ❌ Not installing hook in integration repositories
- ❌ Allowing agents to work without hook installed
- ❌ Installing hook AFTER agents start committing

### Success Criteria:
- ✅ Hook installed in ALL repositories
- ✅ Hook prevents incorrect branch names
- ✅ Hook provides clear correction instructions
- ✅ Hook installed BEFORE any agent work
- ✅ Installation verified after each attempt

## Testing the Hook

To test the hook installation:

```bash
# Test with correct branch name
cd efforts/phase1/wave1/test-effort
git checkout -b "project/phase1/wave1/test-effort"
echo "test" > test.txt && git add . && git commit -m "test"
# Should succeed

# Test with incorrect branch name
git checkout -b "wrong-branch-name"
echo "test" > test.txt && git add . && git commit -m "test"
# Should fail with error message
```

## Related Rules

- R308: Incremental branching strategy (defines valid patterns)
- R312: Git config immutability (hook complements config locking)
- R176: Workspace isolation (hook enforces isolation boundaries)
- R309: Never create efforts in SF repo (hook prevents SF pollution)

## Common Mistakes

1. **Installing hook too late** - Must be installed IMMEDIATELY after clone
2. **Forgetting integration repos** - ALL repos need the hook
3. **Not verifying installation** - Always check the hook was installed
4. **Installing after agent spawning** - Hook must be ready BEFORE agents work
5. **Not handling existing hooks** - Hook installer merges with existing hooks

## Example Workflow

1. Orchestrator clones repository for effort
2. Orchestrator creates correct branch
3. **Orchestrator installs branch validation hook** ← R410 REQUIREMENT
4. Orchestrator locks git config (R312)
5. Orchestrator spawns agent to work
6. Agent attempts to commit
7. Hook validates branch name
8. Commit proceeds if branch name is correct

## Enforcement

This rule is enforced by:
- State machine rules requiring hook installation
- Validation scripts checking for hook presence
- Integration tests verifying hook functionality
- Grading system penalizing missing hooks

**REMEMBER**: Every repository created by Software Factory MUST have this hook installed to prevent branch naming violations!