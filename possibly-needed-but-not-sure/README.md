# Possibly Needed But Not Sure

This directory contains additional protocol and instruction files that may be useful depending on your project's specific needs. Review these files and move them to the appropriate directories if needed for your implementation.

## Files That Might Be Useful

### Planning Protocols
- `CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md` - Detailed instructions for creating effort plans
- `ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md` - How orchestrator manages effort planning

### Review Protocols  
- `CODE-REVIEWER-COMPREHENSIVE-GUIDE.md` - Complete code review process
- `ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md` - Architectural review procedures
- `WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md` - Wave completion reviews
- `PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md` - Phase assessment protocols

### Testing Protocols
- `TEST-DRIVEN-VALIDATION-REQUIREMENTS.md` - Testing coverage requirements
- `PHASE-COMPLETION-FUNCTIONAL-TESTING.md` - End-of-phase testing

### State Management
- `TODO-STATE-MANAGEMENT-PROTOCOL.md` - Detailed TODO persistence rules
- `ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md` - Complete orchestration guide

### Templates
- `WORK-LOG-TEMPLATE.md` - Template for effort work logs
- `IMPLEMENTATION-PLAN-TEMPLATE.md` - Template for effort plans

## When You Might Need These

### Starting a Complex Project
Move these to active use:
- All planning protocols
- All review protocols
- Templates

### Projects with Strict Quality Requirements
Move these to active use:
- Testing protocols
- Comprehensive review guides
- Phase completion testing

### Projects with Multiple Teams
Move these to active use:
- Detailed planning protocols
- Architect review protocols
- State management protocols

### Simple Projects
You might not need:
- Phase-specific protocols (if single phase)
- Architect reviews (if small scope)
- Complex testing protocols (if simple validation)

## How to Activate

1. Review the file content
2. If needed, move to appropriate directory:
   ```bash
   # Move to protocols directory
   mv CODE-REVIEWER-COMPREHENSIVE-GUIDE.md ../protocols/
   
   # Move to agent instructions
   mv ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md ../agent-instructions/
   ```
3. Update references in your agent configurations
4. Add to agent reading lists in `.claude/CLAUDE.md`

## Note
These files are from a production system and contain battle-tested protocols. Even if not immediately needed, they provide valuable reference material for handling edge cases and scaling your software factory.