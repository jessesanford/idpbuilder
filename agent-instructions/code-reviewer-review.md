# Code Reviewer Review Task

## Task: Review Implementation of [EFFORT_NAME]

### Context
- **Phase**: [PHASE]
- **Wave**: [WAVE]  
- **Effort**: [EFFORT_NUMBER]
- **Working Directory**: [WORKING_DIR]
- **Branch**: [BRANCH_NAME]
- **Developer**: [DEVELOPER_AGENT]

### Your Mission

Review the implementation of [EFFORT_NAME] for code quality, compliance, and completeness.

### Required Reading

1. **Review standards and protocols:**
   - `../../../protocols/CODE-REVIEWER-COMPREHENSIVE-GUIDE.md` - Review guidelines
   - `../../../protocols/IMPERATIVE-LINE-COUNT-RULE.md` - Line count limits
   - `../../../protocols/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md` - Test requirements
   - `./IMPLEMENTATION-PLAN.md` - What was supposed to be built
   - `./work-log.md` - Developer's progress log

### Review Process

1. **Measure Line Count FIRST**
   ```bash
   [LINE_COUNTER_PATH] -c [BRANCH_NAME]
   ```
   - If >800 lines, STOP - requires split
   - If 700-800 lines, note as warning

2. **Code Quality Review**
   - [ ] Code follows project conventions
   - [ ] Proper error handling implemented
   - [ ] No code duplication
   - [ ] Clear variable/function names
   - [ ] Appropriate comments where needed
   - [ ] No sensitive data exposed

3. **Completeness Review**
   - [ ] All items from IMPLEMENTATION-PLAN.md completed
   - [ ] All required files created/modified
   - [ ] Integration points properly connected
   - [ ] No TODOs or incomplete sections

4. **Testing Review**
   - [ ] Unit tests present and comprehensive
   - [ ] Test coverage meets requirements ([COVERAGE]%)
   - [ ] All tests passing
   - [ ] Edge cases covered
   - [ ] Integration tests where appropriate

5. **Functional Review**
   - [ ] Code accomplishes intended purpose
   - [ ] Performance considerations addressed
   - [ ] Security best practices followed
   - [ ] Logging implemented appropriately

### Review Outcomes

Based on your review, determine one of:

#### 1. ACCEPTED ✅
Use when:
- Line count <800
- All requirements met
- Code quality acceptable
- Tests comprehensive and passing

#### 2. CHANGES_REQUIRED 🔄
Use when:
- Minor issues found
- Missing tests
- Code quality issues
- Documentation needed

Create `REVIEW-FEEDBACK.md` with specific issues.

#### 3. NEEDS_SPLIT ⚠️
Use when:
- Line count >800
- Must create split plan

Create `SPLIT-IMPLEMENTATION-PLAN.md` with strategy.

### Output Format

Create `CODE-REVIEW-RESULTS.md`:

```markdown
# Code Review Results: [EFFORT_NAME]

## Summary
- **Verdict**: [ACCEPTED/CHANGES_REQUIRED/NEEDS_SPLIT]
- **Line Count**: [number] lines
- **Test Coverage**: [percentage]%
- **Review Date**: [date]

## Line Count Analysis
[Output from line counter tool]

## Review Findings

### ✅ Strengths
- [What was done well]

### ❌ Issues Found
- [Specific problems]

### 📋 Required Changes
[If CHANGES_REQUIRED, list specific fixes needed]

### 🔀 Split Strategy
[If NEEDS_SPLIT, outline split plan]

## Recommendations
[Specific actionable feedback]
```

### If Split Required

Create `SPLIT-IMPLEMENTATION-PLAN.md`:
```markdown
# Split Plan: [EFFORT_NAME]

## Current State
- Total Lines: [number]
- Must split into parts <800 lines each

## Split Strategy

### Part 1: [Name] (~[X] lines)
- Files: [list]
- Components: [list]

### Part 2: [Name] (~[Y] lines)
- Files: [list]
- Components: [list]

## Implementation Order
1. Complete Part 1 first
2. Review Part 1
3. Then implement Part 2
```

### Success Criteria

- [ ] Line count measured and documented
- [ ] All code reviewed thoroughly
- [ ] Clear verdict provided
- [ ] Specific feedback documented
- [ ] Split plan created if needed

### Remember

- Line count is MANDATORY - check first
- Be specific in feedback
- Consider maintainability
- Verify test coverage
- Document all findings clearly