# Branch Validation Documentation

## Overview

The Software Factory 2.0 git hooks include branch name validation to ensure consistent naming conventions across different repository types.

## Special Branch: software-factory-2.0

The branch name `software-factory-2.0` has special significance:

### Planning Repositories
- **ALLOWED**: The `software-factory-2.0` branch is explicitly allowed in planning repositories
- **Purpose**: Helps identify when you're working in the planning context vs implementation
- **Examples of planning repos**:
  - `/home/vscode/idpbuilder`
  - `/home/vscode/idpbuilder-push`
  - Any repo with `PROJECT-IMPLEMENTATION-PLAN.md` but no `target-repo-config.yaml`
  - Repos with "software-factory-template" in their remote URL

### Effort Repositories
- **BLOCKED**: The `software-factory-2.0` branch is NOT allowed in effort repositories
- **Reason**: Effort repos must follow strict Software Factory branch naming patterns
- **Examples of effort repos**:
  - Any repo with `target-repo-config.yaml`
  - Repos in `efforts/` directory structure
  - Working copies for implementation

## Repository Type Detection

The master pre-commit hook (`master-pre-commit.sh`) automatically detects repository type:

1. **Planning Repository** detected if:
   - Remote URL contains "software-factory-template" or "software-factory-2.0-template"
   - Repository basename is "idpbuilder" or "idpbuilder-push"
   - Contains `PROJECT-IMPLEMENTATION-PLAN.md` without `target-repo-config.yaml`

2. **Effort Repository** detected if:
   - Contains `target-repo-config.yaml`
   - Located in `efforts/` directory structure
   - Has `../software-factory-template` directory

3. **General Repository**: All other repositories (minimal validation)

## Branch Naming Rules by Repository Type

### Planning Repositories
Allowed branch patterns:
- `main`, `master`, `develop`, `staging`, `production`
- **`software-factory-2.0`** (special planning branch)
- `feature/*`, `bugfix/*`, `hotfix/*`, `release/*`
- `planning/*`, `experiment/*`, `test/*`
- `effort/*`, `phase*/*`, `wave*/*` (for planning work)

### Effort Repositories
Allowed branch patterns (with optional project prefix):
- `[prefix/]phase{X}/wave{Y}/{effort-name}`
- `[prefix/]phase{X}/wave{Y}/{effort-name}--split-{NNN}`
- `[prefix/]phase{X}/wave{Y}/{effort-name}-fix`
- `[prefix/]phase{X}/wave{Y}/integration`
- `[prefix/]phase{X}/integration`
- `[prefix/]integration`
- Traditional branches: `main`, `master`, `feature/*`, `bugfix/*`

Note: `software-factory-2.0` is explicitly NOT allowed in effort repositories.

## Testing Branch Validation

To test if a branch name is valid:

1. **Without committing**: Check branch name manually against the patterns above
2. **With a test commit**: The pre-commit hook will validate and provide feedback
3. **Skip validation**: Set `SKIP_BRANCH_VALIDATION=true` environment variable (not recommended)

## Hook Installation

The hooks are automatically installed when using the Software Factory template. To manually install:

```bash
# Copy the master hook to your repository
cp /home/vscode/software-factory-template/tools/git-commit-hooks/master-pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

## Troubleshooting

If branch validation fails unexpectedly:

1. Check repository type detection: The hook output shows "Repository type: [type]"
2. Verify branch name follows the patterns for that repository type
3. For planning repos, ensure `software-factory-2.0` branch is recognized
4. For effort repos, ensure proper Software Factory naming is used

## Implementation Details

The branch validation is implemented in:
- **Planning repos**: `tools/git-commit-hooks/planning-hooks/branch-name-validation.hook`
- **Effort repos**: `tools/git-commit-hooks/effort-hooks/branch-name-validation.hook`
- **Master hook**: `tools/git-commit-hooks/master-pre-commit.sh`

The master hook coordinates which validation rules to apply based on repository type.