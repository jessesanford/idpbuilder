# Code Review Command

## Purpose
Standardized workflow for code review agents to follow when reviewing implementation.

## Command Flow

### 1. Initial Setup
```bash
# Verify you're in the correct directory
pwd
# Should show: [project]/phase[X]/wave[Y]/effort[Z]

# Verify correct branch
git branch --show-current
# Should show: phase[X]/wave[Y]/effort[Z]-[name]
```

### 2. Measure Line Count (CRITICAL FIRST STEP)
```bash
# Run line counter
[project]/tools/line-counter.sh -c $(git branch --show-current)

# Decision tree:
# - If >800 lines: STOP - Create split plan
# - If 700-800: WARNING - Note in review
# - If <700: PROCEED with review
```

### 3. Read Implementation Context
```bash
# Read the implementation plan
cat IMPLEMENTATION-PLAN.md

# Read the work log
cat work-log.md

# Check git history
git log --oneline -10
```

### 4. Review Code Changes
```bash
# See all changes
git diff [base-branch]

# List modified files
git diff --name-only [base-branch]

# Check for generated files that shouldn't be counted
git diff --name-only [base-branch] | grep -E "(generated|pb\.go|\.pb\.|_gen\.)"
```

### 5. Run Validation Checks
```bash
# Run tests
[test-command]  # e.g., npm test, go test ./..., pytest

# Check test coverage
[coverage-command]  # e.g., npm run coverage, go test -cover

# Run linting
[lint-command]  # e.g., npm run lint, golangci-lint run

# Run build
[build-command]  # e.g., npm run build, go build ./...
```

### 6. Architectural Review Checklist

#### Code Quality
- [ ] Follows project coding standards
- [ ] Proper error handling
- [ ] No code duplication
- [ ] Clear naming conventions
- [ ] Appropriate comments

#### Completeness
- [ ] All IMPLEMENTATION-PLAN.md items done
- [ ] No TODO comments left
- [ ] All edge cases handled
- [ ] Proper logging added

#### Testing
- [ ] Unit tests comprehensive
- [ ] Test coverage >= required %
- [ ] Integration tests where needed
- [ ] All tests passing

#### Security
- [ ] No hardcoded secrets
- [ ] Input validation present
- [ ] No SQL injection vulnerabilities
- [ ] Proper authentication/authorization

### 7. Generate Review Results

Create `CODE-REVIEW-RESULTS.md`:

```markdown
# Code Review Results: [Effort Name]

## Summary
- **Verdict**: [ACCEPTED/CHANGES_REQUIRED/NEEDS_SPLIT]
- **Line Count**: [X] lines
- **Test Coverage**: [Y]%
- **Tests Status**: [PASSING/FAILING]

## Line Count Details
[paste output from line counter]

## Review Findings

### Strengths ✅
- [positive points]

### Issues Found ❌
- [problems identified]

### Required Changes 📋
[if CHANGES_REQUIRED]
1. [specific fix needed]
2. [another fix]

## Recommendations
- [suggestions for improvement]
```

### 8. Handle Different Outcomes

#### If ACCEPTED:
```bash
# Update work log
echo "✅ Code review ACCEPTED - $(date)" >> work-log.md

# Report to orchestrator
echo "Review complete: ACCEPTED"
echo "Line count: [X] lines"
echo "Test coverage: [Y]%"
```

#### If CHANGES_REQUIRED:
```bash
# Create feedback file
cat > REVIEW-FEEDBACK.md << 'EOF'
# Review Feedback

## Required Changes
1. [Issue]: [Description]
   - File: [path]
   - Line: [number]
   - Fix: [what to do]

2. [Next issue]...
EOF

# Update work log
echo "🔄 Code review CHANGES_REQUIRED - $(date)" >> work-log.md
```

#### If NEEDS_SPLIT (>800 lines):
```bash
# Get detailed breakdown
[project]/tools/line-counter.sh -c $(git branch --show-current) -d

# Create split plan
cat > SPLIT-IMPLEMENTATION-PLAN.md << 'EOF'
# Split Plan for [Effort Name]

## Current Status
- Total Lines: [X]
- Must split to <800 lines per part

## Proposed Splits

### Part 1: [Component Group] (~[Y] lines)
Files to include:
- [file1]
- [file2]

### Part 2: [Component Group] (~[Z] lines)
Files to include:
- [file3]
- [file4]

## Implementation Strategy
1. Create part1 branch from current
2. Remove part2 files
3. Review and merge part1
4. Create part2 branch
5. Implement remaining
EOF
```

### 9. Finalize Review

```bash
# Commit review artifacts
git add CODE-REVIEW-RESULTS.md
git add REVIEW-FEEDBACK.md  # if exists
git add SPLIT-IMPLEMENTATION-PLAN.md  # if exists
git commit -m "review: Complete code review for [effort name]"

# Push changes
git push -u origin $(git branch --show-current)
```

## Success Criteria

- [ ] Line count measured and under limit
- [ ] All tests passing
- [ ] Code quality acceptable
- [ ] Review results documented
- [ ] Clear verdict provided

## Common Issues and Solutions

### Issue: Line count exceeds 800
**Solution**: Create split plan, work with orchestrator to implement in parts

### Issue: Test coverage too low
**Solution**: Request developer to add more tests

### Issue: Integration conflicts
**Solution**: Rebase on latest base branch, resolve conflicts

### Issue: Performance concerns
**Solution**: Profile code, suggest optimizations

## Remember

1. **ALWAYS check line count first**
2. **Be specific in feedback**
3. **Test everything locally**
4. **Document all findings**
5. **Provide actionable feedback**