# Agent Instruction Templates

This directory contains templates for instructing agents when they are spawned by the orchestrator.

## Purpose

When the orchestrator spawns a sub-agent, it must provide clear instructions about:
- What work needs to be done
- Which files to read for context
- What deliverables are expected
- Success criteria

## Templates

### For Software Engineer Agent
- `sw-engineer-implementation.md` - Template for implementation tasks
- `sw-engineer-fix-review.md` - Template for fixing code review issues
- `sw-engineer-split-work.md` - Template for working on split branches

### For Code Reviewer Agent
- `code-reviewer-planning.md` - Template for creating implementation plans
- `code-reviewer-review.md` - Template for reviewing implementation
- `code-reviewer-split-planning.md` - Template for planning splits

### For Architect Reviewer Agent
- `architect-wave-review.md` - Template for wave completion reviews
- `architect-phase-review.md` - Template for phase assessment
- `architect-integration-review.md` - Template for integration reviews

## Usage

The orchestrator should use these templates as a base and customize them with:
1. Specific effort details
2. Current phase/wave/effort numbers
3. Relevant file paths
4. Specific requirements

Example:
```markdown
# Copy template
cp sw-engineer-implementation.md /tmp/task.md

# Customize with specific details
sed -i 's/\[EFFORT_NAME\]/user-authentication/g' /tmp/task.md
sed -i 's/\[PHASE\]/1/g' /tmp/task.md
sed -i 's/\[WAVE\]/2/g' /tmp/task.md

# Spawn agent with customized instructions
Task @agent-sw-engineer with /tmp/task.md
```