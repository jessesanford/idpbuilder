# Software Factory Planning Directory Structure

## Overview

This directory contains all planning artifacts for the Software Factory 2.0 system. Planning files are organized hierarchically to match the Phase → Wave → Effort structure.

## Directory Structure

```
planning/
├── project/                    # Project-wide plans and architecture
│   ├── PROJECT-IMPLEMENTATION-PLAN.md      # Master implementation plan (R383)
│   ├── PROJECT-ARCHITECTURE-PLAN.md        # Overall architecture vision
│   ├── PROJECT-TEST-PLAN.md                # Testing strategy
│   └── PROJECT-INTEGRATE_WAVE_EFFORTS-PLAN.md         # Integration approach
│
├── phase{N}/                   # Phase-specific plans
│   ├── PHASE-ARCHITECTURE-PLAN.md          # Phase architecture details
│   ├── PHASE-TEST-PLAN.md                  # Phase testing requirements
│   ├── PHASE-INTEGRATE_WAVE_EFFORTS-PLAN.md           # Phase integration strategy
│   │
│   └── wave{N}/                # Wave-specific plans
│       ├── WAVE-IMPLEMENTATION-PLAN.md     # Wave implementation details
│       ├── WAVE-TEST-PLAN.md               # Wave test requirements
│       ├── EFFORT-PLANS/                   # Individual effort plans
│       │   ├── effort-001-[name].md
│       │   ├── effort-002-[name].md
│       │   └── ...
│       └── reviews/             # Review artifacts
│           ├── code-reviews/
│           └── architecture-reviews/
```

## File Naming Conventions

### Project Level
- `PROJECT-[TYPE]-PLAN.md` - Master plans that span entire project
- All caps for visibility and importance

### Phase Level
- `PHASE-[TYPE]-PLAN.md` - Plans specific to a phase
- Phase directories: `phase1/`, `phase2/`, etc.

### Wave Level
- `WAVE-[TYPE]-PLAN.md` - Plans specific to a wave
- Wave directories: `wave1/`, `wave2/`, etc.

### Effort Level
- `effort-###-[descriptive-name].md` - Individual effort plans
- Use 3-digit numbering (001, 002, etc.)
- Lowercase with hyphens for descriptive names

## Required Plans by Level

### Project (Mandatory)
- **PROJECT-IMPLEMENTATION-PLAN.md** - Master roadmap (see R383)
- **PROJECT-ARCHITECTURE-PLAN.md** - System architecture

### Phase (Mandatory)
- **PHASE-ARCHITECTURE-PLAN.md** - Phase-specific architecture
- **PHASE-TEST-PLAN.md** - Testing requirements for phase

### Wave (Mandatory)
- **WAVE-IMPLEMENTATION-PLAN.md** - Wave breakdown
- **EFFORT-PLANS/** - Directory containing all effort plans

### Effort (Mandatory)
- Individual effort plan following template in `templates/EFFORT-PLAN-TEMPLATE-WITH-SCOPE.md`

## Key Rules

### R383 - Timestamp Patterns
All plans must include:
- Creation timestamp
- Last modified timestamp
- Author/Agent identification

### R287 - TODO Persistence
Planning TODOs must be saved:
- After plan creation
- After major updates
- Before state transitions

### Size Limits
- Efforts: Max 800 lines (hard limit)
- Waves: Collection of sized efforts
- Phases: Logical feature groupings

## Creating New Plans

1. **Start with Templates**
   - Copy from `templates/` directory
   - Use `.example` files as reference
   - Maintain consistent formatting

2. **Required Metadata**
   ```markdown
   ---
   created: YYYY-MM-DD HH:MM:SS Z
   modified: YYYY-MM-DD HH:MM:SS Z
   agent: [agent-name]
   state: [current-state]
   ---
   ```

3. **Version Control**
   - Commit after creation
   - Tag major milestones
   - Never delete, archive instead

## Integration with State Machine

Planning files directly map to state machine transitions:

- `INIT` → Creates PROJECT-IMPLEMENTATION-PLAN
- `PLANNING` → Creates PHASE plans
- `CREATE_NEXT_INFRASTRUCTURE` → Creates WAVE plans
- `SPAWN_SW_ENGINEERS` → Uses EFFORT plans

## Quality Checklist

Before finalizing any plan:
- [ ] Follows naming convention
- [ ] Includes required metadata
- [ ] References appropriate rules
- [ ] Sized within limits
- [ ] Reviewed by architect agent
- [ ] Committed to version control

## Support

For questions about planning structure:
1. Check `rule-library/` for specific rules
2. Review `templates/` for examples
3. Consult state machine documentation