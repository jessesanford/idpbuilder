# SOFTWARE FACTORY 2.0 - DUAL REPOSITORY ARCHITECTURE

## Overview

Software Factory 2.0 uses a **dual-repository architecture** to separate planning/orchestration from actual code implementation. This design provides clean separation of concerns and allows the factory system to work with any target project.

## The Two Repositories

### 1. Planning Repository (Factory System)
**Example:** `software-factory-template`, `software-factory-2.0-template`

This repository contains:
- **Orchestrator State** (`orchestrator-state.json`) - Tracks all efforts, phases, waves
- **Agent Configurations** (`.claude/agents/`) - Defines agent behaviors
- **Rule Library** (`rule-library/`) - System rules and constraints
- **State Machine** (`SOFTWARE-FACTORY-STATE-MACHINE.md`) - Valid states and transitions
- **Templates** (`templates/`) - Planning and documentation templates
- **Utilities** (`utilities/`) - Helper scripts for management
- **TODO Files** (`todos/`) - Agent task tracking

### 2. Target Repository (Project Code)
**Example:** `idpbuilder-oci-build-push`, `your-actual-project`

This repository contains:
- **Main Branch** - Production code
- **Effort Branches** - Individual implementation branches
  - Pattern: `{project-name}/phase{N}/wave{N}/{effort-name}`
  - Example: `idpbuilder-oci-build-push/phase1/wave1/kind-certificate-extraction`
- **Split Branches** - When efforts exceed size limits
  - Pattern: `{project-name}/phase{N}/wave{N}/{effort-name}-split-{NNN}`
  - Example: `idpbuilder-oci-build-push/phase1/wave1/registry-auth-split-001`
- **Integration Branches** - For merging waves and phases
  - Wave: `phase{N}-wave{N}-integration`
  - Phase: `phase{N}-integration`

## Branch Naming Convention

All effort branches follow this pattern:
```
{project-name}/phase{phase-number}/wave{wave-number}/{effort-name}
```

Where:
- `{project-name}` - The target repository name (e.g., `idpbuilder-oci-build-push`)
- `{phase-number}` - Phase number (1, 2, 3, etc.)
- `{wave-number}` - Wave number within the phase (1, 2, 3, etc.)
- `{effort-name}` - Descriptive effort name (e.g., `kind-certificate-extraction`)

## Working with the Dual Repository System

### Initial Setup

1. **Clone the Planning Repository:**
   ```bash
   git clone https://github.com/org/software-factory-template.git
   cd software-factory-template
   ```

2. **Clone the Target Repository:**
   ```bash
   git clone https://github.com/org/your-project.git
   cd your-project
   ```

3. **Create Orchestrator State:**
   ```bash
   cp orchestrator-state.json.example orchestrator-state.json
   # Edit to match your project
   ```

### During Development

Software Factory agents work in the target repository:

1. **Orchestrator** reads state from planning repo
2. **SW Engineers** create branches in target repo
3. **Code Reviewers** review code in target repo
4. **Architects** assess integration in target repo

### Restoring Efforts

After cloning a fresh planning repository, use the restoration utility:

```bash
# From the planning repository
cd utilities

# Restore all effort branches from the target repository
./restore-all-efforts.sh ../orchestrator-state.json https://github.com/org/your-project.git

# This creates the directory structure:
# efforts/
#   phase1/
#     wave1/
#       effort-name-1/  (full clone)
#       effort-name-2/  (full clone)
#     wave2/
#       effort-name-3/  (full clone)
```

### Pushing Local Efforts

If you have local working copies that haven't been pushed to remote:

```bash
# From the planning repository
cd utilities

# Push all local effort branches to remote
./push-local-efforts.sh ../orchestrator-state.json /path/to/target-repo

# This ensures all branches referenced in orchestrator-state.json
# are available in the remote repository
```

## Common Scenarios

### Scenario 1: Starting Fresh Development

1. Clone planning repository
2. Clone target repository
3. Create orchestrator-state.json
4. Spawn agents to work on efforts
5. Agents create branches in target repository

### Scenario 2: Continuing Development

1. Clone planning repository
2. Run `restore-all-efforts.sh` to get all effort branches
3. Continue with agent work

### Scenario 3: Sharing Work with Team

1. Ensure all local branches are pushed:
   ```bash
   ./push-local-efforts.sh
   ```
2. Team member clones planning repo
3. Team member runs:
   ```bash
   ./restore-all-efforts.sh orchestrator-state.json https://github.com/org/project.git
   ```

### Scenario 4: Recovery After System Reset

1. Clone planning repository
2. Locate orchestrator-state.json (or recover from backup)
3. Run restoration:
   ```bash
   ./restore-all-efforts.sh orchestrator-state.json https://github.com/org/project.git
   ```

## Important Notes

### Branch Persistence
- Effort branches should be pushed to the target repository regularly
- The planning repository tracks branch names but not the code
- Lost local branches can't be recovered without remote backups

### State Synchronization
- The orchestrator-state.json is the source of truth for effort status
- Always commit and push state changes to the planning repository
- State file should match actual branches in target repository

### Working Directory Structure
```
/home/vscode/
├── software-factory-template/        # Planning repository
│   ├── orchestrator-state.json      # Current state
│   ├── utilities/                   # Management scripts
│   ├── rule-library/                # System rules
│   └── efforts/                     # Restored effort working copies
│       ├── phase1/
│       │   ├── wave1/
│       │   │   ├── effort-1/       # Full clone
│       │   │   └── effort-2/       # Full clone
│       │   └── wave2/
│       └── phase2/
└── workspaces/
    └── your-project/                # Target repository clone
        └── .git/                    # Contains all branches
```

## Troubleshooting

### Issue: "Branch not found in remote"
**Cause:** Local branches haven't been pushed to remote
**Solution:** Run `push-local-efforts.sh` from the working copy containing the branches

### Issue: "Could not detect target repository"
**Cause:** Script can't auto-detect the target repository URL
**Solution:** Explicitly provide the target repository URL:
```bash
./restore-all-efforts.sh orchestrator-state.json https://github.com/org/project.git
```

### Issue: "Uncommitted changes detected"
**Cause:** Local working copy has uncommitted changes
**Solution:** Commit or stash changes before restoration

### Issue: Branches exist locally but not in orchestrator-state.json
**Cause:** State file is out of sync with actual development
**Solution:** Update orchestrator-state.json to reflect actual branches

## Best Practices

1. **Regular Commits:** Commit and push both planning and target repositories frequently
2. **State Backups:** Keep backups of orchestrator-state.json
3. **Branch Cleanup:** Delete merged branches to keep repository clean
4. **Documentation:** Document any deviations from standard patterns
5. **Verification:** Regularly verify state file matches actual branches

## Migration from Single Repository

If migrating from a single-repository setup:

1. Create new planning repository from template
2. Move orchestrator-state.json to planning repo
3. Update branch references in state file
4. Push all effort branches to target repository
5. Update agent configurations for dual-repo setup

## Security Considerations

- Planning repository can be private (contains no code)
- Target repository follows your normal security policies
- Credentials needed for both repositories
- Consider using SSH keys for automation

## Conclusion

The dual-repository architecture provides:
- **Separation of Concerns:** Planning vs implementation
- **Reusability:** Factory system works with any project
- **Flexibility:** Different access controls for each repository
- **Clarity:** Clear distinction between system and project