# Create Wave Implementation Plan Command

## Purpose
Guide the orchestrator in creating implementation plans for all efforts in a wave.

## Prerequisites
- Phase plan exists in `/phase-plans/PHASE[X]-TEMPLATE.md`
- Orchestrator state file is current
- Previous wave (if any) is complete and integrated

## Command Flow

### 1. Load Current State
```bash
# Check current position
cat orchestrator-state.yaml | grep -E "current_phase|current_wave"

# Verify we're in WAVE_START state
cat orchestrator-state.yaml | grep current_state
# Should show: current_state: "WAVE_START"
```

### 2. Read Phase Plan
```bash
# Load phase plan for context
cat phase-plans/PHASE[X]-TEMPLATE.md

# Identify efforts for current wave
grep -A 20 "Wave [Y]" phase-plans/PHASE[X]-TEMPLATE.md
```

### 3. Create Wave Directory Structure
```bash
# Create wave directory
mkdir -p phase[X]/wave[Y]

# Create effort directories
for i in {1..N}; do
  mkdir -p phase[X]/wave[Y]/effort${i}
done
```

### 4. Plan Each Effort

For each effort in the wave:

#### A. Create Effort Working Directory
```bash
EFFORT_DIR="phase[X]/wave[Y]/effort[N]"
cd $EFFORT_DIR
```

#### B. Spawn Code Reviewer for Planning
```yaml
Task: @agent-code-reviewer
Instruction: Create implementation plan for [effort-name]
Context:
  - Phase: [X]
  - Wave: [Y]
  - Effort: [N]
  - Description: [from phase plan]
  - Working Dir: $EFFORT_DIR
Template: agent-instructions/code-reviewer-planning.md
```

#### C. Verify Plan Created
```bash
# Check plan exists
test -f IMPLEMENTATION-PLAN.md && echo "✅ Plan created" || echo "❌ Plan missing"

# Check work log initialized
test -f work-log.md && echo "✅ Work log created" || echo "❌ Work log missing"
```

### 5. Create Wave Summary

Create `phase[X]/wave[Y]/WAVE-PLAN-SUMMARY.md`:

```markdown
# Wave [Y] Implementation Plan

## Overview
- **Phase**: [X] - [Phase Name]
- **Wave**: [Y]
- **Total Efforts**: [N]
- **Estimated Total Lines**: [sum of estimates]

## Efforts Planned

### Effort 1: [Name]
- **Directory**: phase[X]/wave[Y]/effort1
- **Estimated Lines**: [X]
- **Dependencies**: [list]
- **Plan Status**: ✅ Created

### Effort 2: [Name]
- **Directory**: phase[X]/wave[Y]/effort2
- **Estimated Lines**: [Y]
- **Dependencies**: [list]
- **Plan Status**: ✅ Created

[Continue for all efforts...]

## Implementation Order
Based on dependencies:
1. Effort [X] - [reason]
2. Effort [Y] - [reason]
...

## Risk Assessment
- **High Risk**: [efforts that might exceed 800 lines]
- **Medium Risk**: [complex integrations]
- **Low Risk**: [straightforward implementations]

## Success Criteria
- [ ] All efforts implemented
- [ ] All efforts <800 lines
- [ ] All efforts reviewed
- [ ] Wave integration successful
```

### 6. Update Orchestrator State

```yaml
# Update orchestrator-state.yaml
current_state: "EFFORT_SELECTION"
current_wave_plan:
  phase: [X]
  wave: [Y]
  total_efforts: [N]
  efforts_planned:
    - name: [effort1-name]
      dir: phase[X]/wave[Y]/effort1
      status: planned
      estimated_lines: [X]
    - name: [effort2-name]
      dir: phase[X]/wave[Y]/effort2
      status: planned
      estimated_lines: [Y]
```

### 7. Validate Wave Plan

```bash
# Check all efforts have plans
for dir in phase[X]/wave[Y]/effort*/; do
  echo "Checking $dir..."
  test -f "$dir/IMPLEMENTATION-PLAN.md" || echo "Missing plan in $dir"
  test -f "$dir/work-log.md" || echo "Missing work-log in $dir"
done

# Verify total estimated lines
echo "Total estimated lines for wave:"
grep -h "Estimated Total Lines:" phase[X]/wave[Y]/effort*/IMPLEMENTATION-PLAN.md | \
  awk '{sum += $NF} END {print sum}'
```

### 8. Prepare for Implementation

```bash
# Create implementation tracking
cat > phase[X]/wave[Y]/IMPLEMENTATION-STATUS.md << 'EOF'
# Wave [Y] Implementation Status

## Progress
- [ ] Effort 1: [name] - NOT_STARTED
- [ ] Effort 2: [name] - NOT_STARTED
...

## Notes
- Created: [date]
- Ready for implementation
EOF
```

## Decision Points

### If any effort estimates >600 lines:
- Flag as high risk
- Plan for potential split
- Consider breaking into smaller efforts

### If dependencies complex:
- Create dependency graph
- Define implementation order
- Plan integration points

### If similar efforts exist:
- Plan for code reuse
- Create shared components first
- Avoid duplication

## Success Criteria

- [ ] All efforts have implementation plans
- [ ] All efforts have work logs
- [ ] Wave summary created
- [ ] State file updated
- [ ] Risk assessment complete
- [ ] Implementation order defined

## Common Patterns

### API/Model Efforts
```
1. Define data models/types
2. Create validators
3. Add serialization
4. Write tests
```

### Service/Logic Efforts
```
1. Define interfaces
2. Implement core logic
3. Add error handling
4. Create unit tests
```

### Integration Efforts
```
1. Define connection points
2. Implement adapters
3. Add retry logic
4. Test integration
```

## Troubleshooting

### Problem: Effort too complex to estimate
**Solution**: Break into sub-efforts or request architect input

### Problem: Dependencies unclear
**Solution**: Review phase plan, consult existing code

### Problem: Similar to existing code
**Solution**: Plan refactoring to share implementation

## Remember

1. **Each effort needs clear boundaries**
2. **Consider line count from planning stage**
3. **Plan for testing upfront**
4. **Document dependencies clearly**
5. **Keep efforts focused and cohesive**