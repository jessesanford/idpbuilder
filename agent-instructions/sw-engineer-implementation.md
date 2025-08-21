# Software Engineer Implementation Task

## Task: Implement [EFFORT_NAME]

### Context
- **Phase**: [PHASE]
- **Wave**: [WAVE]
- **Effort**: [EFFORT_NUMBER]
- **Working Directory**: [WORKING_DIR]
- **Branch**: [BRANCH_NAME]

### Your Mission

You are tasked with implementing [EFFORT_NAME] as part of the software factory process.

### Required Reading

1. **First, verify your environment:**
   ```bash
   pwd  # Should be [WORKING_DIR]
   git branch --show-current  # Should be [BRANCH_NAME]
   ```

2. **Read these files in order:**
   - `./IMPLEMENTATION-PLAN.md` - Your implementation roadmap
   - `./work-log.md` - Track your progress here
   - `../../../protocols/SW-ENGINEER-EXPLICIT-INSTRUCTIONS.md` - Your operating rules
   - `../../../protocols/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md` - Testing requirements
   - `../../../protocols/IMPERATIVE-LINE-COUNT-RULE.md` - Line count limits

### Implementation Requirements

1. **Follow the implementation plan exactly**
   - Implement each component listed in IMPLEMENTATION-PLAN.md
   - Update work-log.md as you complete each item

2. **Code Quality Standards**
   - Write clean, maintainable code
   - Follow project conventions and patterns
   - Include appropriate error handling
   - Add logging where appropriate

3. **Testing Requirements**
   - Write unit tests for all new code
   - Ensure [TEST_COVERAGE]% test coverage minimum
   - All tests must pass before marking complete

4. **Line Count Compliance**
   - Monitor your line count every ~200 lines
   - Run: `[LINE_COUNTER_PATH] -c [BRANCH_NAME]`
   - MUST stay under 800 lines total
   - If approaching 700 lines, notify immediately

### Deliverables

1. **Implementation** matching the plan
2. **Tests** with required coverage
3. **Updated work-log.md** showing progress
4. **All checks passing** (lint, tests, build)

### Success Criteria

- [ ] All items in IMPLEMENTATION-PLAN.md completed
- [ ] Tests written and passing
- [ ] Line count under 800 (measured with line counter)
- [ ] work-log.md updated with completion status
- [ ] Code follows project conventions

### When Complete

1. Run final validation:
   ```bash
   # Check line count
   [LINE_COUNTER_PATH] -c [BRANCH_NAME]
   
   # Run tests
   [TEST_COMMAND]
   
   # Run lint
   [LINT_COMMAND]
   ```

2. Update work-log.md with:
   - Completion timestamp
   - Final line count
   - Test coverage percentage
   - Any notes for reviewer

3. Report completion with:
   - Final line count
   - Test results
   - Any issues encountered

### If You Get Blocked

- Document the blocker in work-log.md
- Provide specific error messages
- Suggest potential solutions
- Request orchestrator assistance

### Remember

- You are implementing, not designing
- Follow the plan exactly as written
- Keep line count under control
- Update work-log.md regularly
- Test everything you write