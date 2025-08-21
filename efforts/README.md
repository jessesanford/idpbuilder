# Efforts Directory

## Purpose
This directory contains working copies for each effort implementation. Each effort gets its own isolated workspace to prevent conflicts and maintain clean separation of work.

## Structure
```
efforts/
└── phase{X}/
    └── wave{Y}/
        └── effort{Z}-{descriptive-name}/
            ├── .git/                           # Git repository
            ├── [source code directories]       # Project code
            ├── IMPLEMENTATION-PLAN.md          # Created by Code Reviewer
            ├── work-log.md                     # Progress tracking
            └── REVIEW-FEEDBACK.md              # If review required fixes
```

## Usage

### Created By
The orchestrator creates these directories when starting a new effort.

### Working In
SW Engineers work exclusively in their assigned effort directory.

### Isolation
Each effort is a separate working copy (sparse checkout or worktree) to ensure:
- No conflicts between parallel efforts
- Clean git history
- Easy review of changes
- Independent testing

## Related Protocols
- See `SOFTWARE-FACTORY-STATE-MACHINE.md` for workflow
- See `ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md` for setup details
- See `EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md` for splits

## Important Notes
1. Never work directly in main repository when implementing efforts
2. Each effort must have its own branch
3. Efforts are measured independently for size compliance
4. Completed efforts are integrated back through review process