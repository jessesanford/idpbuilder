# Code Reviewer Planning Task

## Task: Create Implementation Plan for [EFFORT_NAME]

### Context
- **Phase**: [PHASE]
- **Wave**: [WAVE]
- **Effort**: [EFFORT_NUMBER]
- **Working Directory**: [WORKING_DIR]
- **Description**: [EFFORT_DESCRIPTION]

### Your Mission

Create a detailed implementation plan for [EFFORT_NAME] that a software engineer can follow to implement this effort.

### Required Reading

1. **Read these planning resources:**
   - `../../../protocols/CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md` - How to create plans
   - `../../../phase-plans/PHASE[PHASE]-TEMPLATE.md` - Phase context
   - `../../../protocols/WORK-LOG-TEMPLATE.md` - Work log format
   - `../../../protocols/IMPERATIVE-LINE-COUNT-RULE.md` - Line limits to consider

2. **Understand the project context:**
   - Review existing codebase structure
   - Identify dependencies and integration points
   - Consider architectural patterns in use

### Planning Requirements

1. **Create IMPLEMENTATION-PLAN.md with:**
   - Clear scope definition
   - Step-by-step implementation tasks
   - Specific files to create/modify
   - Required components and their relationships
   - Integration points with existing code
   - Testing requirements
   - Estimated line count per component

2. **Consider Line Count:**
   - Estimate total lines for implementation
   - If likely >600 lines, design for modularity
   - Plan logical split points if needed

3. **Create work-log.md from template:**
   - Copy from WORK-LOG-TEMPLATE.md
   - Initialize with effort details
   - Add planned milestones

### Deliverables

1. **IMPLEMENTATION-PLAN.md** containing:
   ```markdown
   # Implementation Plan: [EFFORT_NAME]
   
   ## Overview
   [Brief description of what this effort implements]
   
   ## Components to Implement
   
   ### 1. [Component Name]
   - **File**: [path/to/file]
   - **Purpose**: [what this does]
   - **Estimated Lines**: [number]
   - **Key Functions/Classes**:
     - [Function/Class name]: [purpose]
   
   ### 2. [Next Component]
   ...
   
   ## Integration Points
   - [Where this connects to existing code]
   
   ## Testing Requirements
   - Unit tests for [components]
   - Integration tests for [workflows]
   - Coverage target: [percentage]%
   
   ## Implementation Order
   1. [First thing to implement]
   2. [Second thing]
   ...
   
   ## Success Criteria
   - [ ] All components implemented
   - [ ] Tests passing
   - [ ] Line count <800
   - [ ] Integrates with [existing systems]
   
   ## Estimated Total Lines: [number]
   ```

2. **work-log.md** initialized from template

### Planning Best Practices

1. **Be Specific**
   - Name exact files and functions
   - Provide clear acceptance criteria
   - Include error handling requirements

2. **Consider Maintainability**
   - Plan for clean interfaces
   - Design for testability
   - Follow project patterns

3. **Think About Splits**
   - If approaching 800 lines, plan modular components
   - Identify natural boundaries for splitting
   - Keep related functionality together

### Success Criteria

- [ ] IMPLEMENTATION-PLAN.md is detailed and actionable
- [ ] Every component has clear requirements
- [ ] Line count estimates are realistic
- [ ] Testing requirements are comprehensive
- [ ] work-log.md is properly initialized
- [ ] Plan follows project architecture patterns

### Output Format

When complete, confirm:
```
✅ Created IMPLEMENTATION-PLAN.md
✅ Created work-log.md from template
✅ Estimated total lines: [number]
✅ Plan is ready for implementation
```

### Remember

- You're planning, not implementing
- Be specific enough that any engineer can follow
- Consider line count limits from the start
- Include all necessary components
- Plan for comprehensive testing