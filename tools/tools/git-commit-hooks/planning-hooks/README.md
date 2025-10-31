# Planning Repository Hooks - STRICT ENFORCEMENT

## Purpose
These hooks enforce STRICT branch discipline in planning repositories to prevent code contamination.

## The Problem
The orchestrator discovered 361 inappropriate branches with implementation code in the idpbuilder-push planning repository. This contamination makes it impossible to maintain clean separation between planning and implementation.

## The Solution: STRICT Branch Enforcement

### Allowed Branches
Only TWO branches are allowed in planning repositories:

1. **software-factory-2.0** - For all Software Factory planning work
2. **main** - ONLY for initial repository setup (before PROJECT-IMPLEMENTATION-PLAN.md exists)

### ALL Other Branches Are BLOCKED
Any attempt to commit on other branches will receive this error:
```
WRONG! YOU ARE ON THE PLANNING REPO, YOU MUST NOT SUBMIT ANY CODE
OR WORK OTHER THAN PLANS TO THIS REPO!

MOVE YOUR CODE OR OTHER CHANGES TO THE APPROPRIATE EFFORT DIRECTORY
(located under $CLAUDE_PROJECT_DIR/efforts/phaseX/waveY/**/*)
OR THE APPROPRIATE INTEGRATE_WAVE_EFFORTS DIRECTORY
(under $CLAUDE_PROJECT_DIR/efforts/**/integration)
AND THE APPROPRIATE BRANCH ON THE TARGET REPO.
```

## Installation

### Quick Installation
```bash
# From the template directory
./tools/install-planning-hooks.sh /path/to/planning/repo

# Or from within a planning repo
/path/to/template/tools/install-planning-hooks.sh
```

### Manual Installation
1. Copy `master-pre-commit.sh` to `.git/hooks/pre-commit`
2. Copy all hook files to `tools/git-commit-hooks/` in the repo
3. Create `.planning-repo` marker file in repo root
4. Make pre-commit hook executable

## How It Works

### Detection
The master hook detects planning repositories by:
- Repository name patterns (idpbuilder, *-push, *-planning)
- Presence of PROJECT-IMPLEMENTATION-PLAN.md without target-repo-config.yaml
- Presence of `.planning-repo` marker file
- Remote URL containing "software-factory-template"

### Enforcement
When detected as a planning repo:
1. **branch-name-validation.hook** runs FIRST
2. Only allows `software-factory-2.0` or `main` (conditionally)
3. Blocks ALL other branches with clear error message
4. Directs users to proper effort directories

### Override (Emergency Only)
```bash
# ONLY use in true emergencies
SKIP_BRANCH_VALIDATION=true git commit -m "Emergency commit"
```

## Files

- **branch-name-validation.hook** - STRICT enforcement (current)
- **branch-name-validation-old.hook** - Flexible version (deprecated)
- **efforts-protection.hook** - Prevents effort directory commits
- **README.md** - This file

## Key Rules
- Planning repos are for PLANS ONLY
- NO implementation code
- NO feature branches
- NO effort branches
- Use effort directories for ALL implementation work

## Support
If you need to work on implementation:
1. Check `$CLAUDE_PROJECT_DIR/target-repo-config.yaml` for target repo
2. Navigate to appropriate effort directory
3. Work on proper branch there
4. NEVER bring code back to planning repo