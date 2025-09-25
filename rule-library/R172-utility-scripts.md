# R172.0.0 - Utility Script Execution

## Rule Statement
Utility scripts in the `utilities/` directory MUST be executed manually and MUST NOT be referred to as "hooks".

## Rationale
These are helper scripts, not Claude Code hooks. Calling them hooks creates false expectations of automatic execution.

## Implementation

### Directory Structure
```
project/
└── utilities/              # Manual helper scripts
    ├── pre-compact.sh     # Manual state preservation
    ├── post-compact.sh    # Manual recovery helper
    ├── todo-preservation.sh # Manual TODO management
    ├── state-snapshot.sh  # Manual checkpoint creation
    └── recovery-assistant.sh # Manual recovery wizard
```

### Execution Protocol
```bash
# All utilities must be run manually
./utilities/pre-compact.sh          # Before expected compaction
./utilities/todo-preservation.sh save # During transitions
./utilities/post-compact.sh         # After context loss
```

## Enforcement
- Documentation MUST clarify manual execution
- Setup.sh MUST NOT configure these as hooks
- Scripts MUST be in `utilities/` not `hooks/`

## Validation
```bash
# Verify correct directory name
ls -la utilities/

# Verify scripts are executable
ls -la utilities/*.sh

# Verify no hooks directory exists
ls -la hooks/ 2>&1 | grep "No such file"
```

## Common Mistakes
- ❌ Expecting automatic execution
- ❌ Calling them "hooks"
- ❌ Putting them in hooks/ directory
- ✅ Manual execution from utilities/ directory