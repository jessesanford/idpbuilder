# R175.0.0 - Manual Utility Usage

## Rule Statement
Utility scripts MUST be explicitly documented as manual tools requiring user execution, not automatic hooks.

## Rationale
Clear expectations prevent confusion and ensure users understand when manual intervention is required.

## Implementation

### Documentation Requirements
Every utility script MUST include:
```bash
#!/bin/bash
# This is a MANUAL utility script - it does NOT run automatically
# Usage: ./utilities/script-name.sh [args]
```

### When to Run Utilities
| Utility | When to Run | Purpose |
|---------|------------|---------|
| pre-compact.sh | Before expected compaction | Comprehensive state save |
| todo-preservation.sh | State transitions | Save/load TODO states |
| state-snapshot.sh | Major milestones | Create recovery points |
| post-compact.sh | After resuming | Check for compaction |
| recovery-assistant.sh | When confused | Interactive recovery |

### Integration in Workflow
```markdown
# In agent instructions
Before transitioning from WAVE_COMPLETE to INTEGRATE_WAVE_EFFORTS:
1. Run: ./utilities/todo-preservation.sh save orchestrator WAVE_COMPLETE
2. Proceed with integration
```

## Enforcement
- README MUST clarify manual execution
- Setup.sh MUST inform users these are manual
- Agent instructions SHOULD include utility calls

## Validation
```bash
# Check for proper documentation
grep -l "MANUAL\|manual" utilities/*.sh

# Verify executable permissions
ls -la utilities/*.sh | grep -c "^-rwx"

# Check for incorrect hooks directory
[ ! -d hooks ] && echo "✅ No hooks directory (correct)"
```

## Common Mistakes
- ❌ Expecting automatic execution
- ❌ Not including in workflow instructions
- ❌ Not making scripts executable
- ✅ Clear manual execution documentation