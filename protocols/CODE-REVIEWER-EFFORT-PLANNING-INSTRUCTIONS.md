# Code Reviewer Effort Planning Instructions

## Purpose
This document provides detailed instructions for Code Reviewers when creating effort implementation plans before SW Engineers begin work.

## When You Are Tasked for Planning

You will receive a task like:
```
Task @agent-code-reviewer:
PURPOSE: Create implementation plan for Effort E{X}.{Y}.{Z}
```

## Your Planning Process

### 1. Context Gathering
```bash
# Read these files IN ORDER:
1. /workspaces/[project]/phase-plans/PHASE{X}-SPECIFIC-IMPL-PLAN.md
2. /workspaces/[project]/orchestrator-state.yaml
3. Any existing addendums or corrections for this phase
```

### 2. Analyze Current State
- What efforts have been completed?
- What patterns have emerged?
- What dependencies exist?
- What integration points need consideration?

### 3. Create Implementation Plan

#### File Location
```
/workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}/IMPLEMENTATION-PLAN.md
```

#### Required Sections

##### Header
```markdown
# Implementation Plan for E{X}.{Y}.{Z}: {Name}
Created: {timestamp}
Created by: @agent-code-reviewer
Phase: {X} - {Phase Name}
Wave: {Y} - {Wave Description}
```

##### Context Analysis
```markdown
## Context Analysis

### Completed Related Efforts
- E{X}.{Y}.1: {summary of what it did}
- E{X}.{Y}.2: {summary of what it did}

### Established Patterns
- Pattern 1: {description}
- Pattern 2: {description}

### Integration Points
- API: {description}
- Database: {description}
- Services: {description}
```

##### Requirements Extraction
```markdown
## Requirements (from Phase Plan)

### Primary Requirements
1. {Exact requirement from phase plan}
2. {Exact requirement from phase plan}

### Derived Requirements
1. {What's implied but not stated}
2. {What's needed for integration}

### Non-Functional Requirements
- Performance: {specifics}
- Security: {specifics}
- Scalability: {specifics}
```

##### Implementation Strategy
```markdown
## Implementation Strategy

### Approach
{Describe the overall approach}

### Design Decisions
1. {Decision}: {Rationale}
2. {Decision}: {Rationale}

### Patterns to Follow
- {Pattern from completed efforts}
- {Pattern from project standards}
```

##### Step-by-Step Instructions
```markdown
## Implementation Steps

### Step 1: {Title}
**Action**: {Exact action to take}
**Files**: {Files to create/modify}
**Validation**: {How to verify this step}

### Step 2: {Title}
**Action**: {Exact action to take}
**Files**: {Files to create/modify}
**Validation**: {How to verify this step}

[Continue for all steps...]
```

##### File Structure
```markdown
## Files to Create/Modify

### New Files
```
{directory}/
├── {file1.ext}       # {purpose}
├── {file2.ext}       # {purpose}
└── {directory2}/
    └── {file3.ext}   # {purpose}
```

### Modified Files
- `{existing/file.ext}`: {what to add/change}
- `{existing/file2.ext}`: {what to add/change}
```

##### Code Templates
```markdown
## Code Templates

### {Component Name}
```{language}
// Starter template for {component}
{starter code structure}
```

### {Component Name 2}
```{language}
// Starter template for {component}
{starter code structure}
```
```

##### Testing Requirements
```markdown
## Testing Requirements

### Unit Tests
- [ ] Test for {functionality}
- [ ] Test for {edge case}
- [ ] Test for {error handling}

### Integration Tests
- [ ] Test {integration point}
- [ ] Test {workflow}

### Coverage Target
- Minimum: {X}%
- Target: {Y}%

### Test File Structure
```
tests/
├── unit/
│   └── {test_file}.{ext}
└── integration/
    └── {test_file}.{ext}
```
```

##### Size Management
```markdown
## Size Management

### Estimated Size
- Core implementation: ~{X} lines
- Tests: ~{Y} lines
- Total: ~{Z} lines

### Size Limit
- Maximum: {limit} lines
- Measurement: line-counter.sh

### Split Strategy (if needed)
If approaching limit:
1. Complete {core functionality} first
2. Split {additional features} to separate effort
3. Prioritize {critical path}
```

##### Success Criteria
```markdown
## Success Criteria

### Functional
- [ ] All requirements implemented
- [ ] Integration with {existing} works
- [ ] No hardcoded values

### Quality
- [ ] Tests pass
- [ ] Coverage >= {X}%
- [ ] Lint clean
- [ ] Build successful

### Size
- [ ] Under {limit} lines per line-counter.sh
- [ ] Properly organized if split needed

### Documentation
- [ ] Code comments for complex logic
- [ ] API documentation if applicable
- [ ] Updated work-log.md
```

##### Integration Notes
```markdown
## Integration Notes

### Dependencies
- Depends on: E{X}.{Y}.{Z}
- Required by: E{X}.{Y}.{Z}

### API Contracts
- Endpoint: {description}
- Data format: {description}

### Breaking Changes
- None expected / {list if any}
```

### 4. Create Work Log Template

Also create `work-log.md`:
```markdown
# Work Log for E{X}.{Y}.{Z}

## Effort: {Name}
## Plan: IMPLEMENTATION-PLAN.md
## Size Limit: {limit} lines

## Sessions
[To be filled by SW Engineer]

## Size Tracking
| Checkpoint | Lines | Status |
|------------|-------|--------|
| Start | 0 | OK |

## Issues
[To be documented as encountered]
```

## Review Considerations

When creating plans, consider:

### Technical Factors
- Current architecture
- Existing patterns
- Performance requirements
- Security requirements
- Scalability needs

### Process Factors
- Size limits
- Testing requirements
- Review criteria
- Integration needs

### Risk Factors
- Complexity areas
- Unknown dependencies
- Potential blockers
- Split likelihood

## Common Patterns

### For API Endpoints
1. Define data models
2. Create validation
3. Implement handler
4. Add tests
5. Document API

### For Data Models
1. Define structure
2. Add validation
3. Create migrations
4. Add CRUD operations
5. Test thoroughly

### For Services
1. Define interface
2. Implement core logic
3. Add error handling
4. Create tests
5. Add monitoring

## Quality Checklist

Before submitting plan:
- [ ] Requirements fully understood
- [ ] Steps are clear and actionable
- [ ] Size estimate is realistic
- [ ] Test requirements defined
- [ ] Success criteria measurable
- [ ] Integration points identified
- [ ] Split strategy defined if needed

## Important Notes

1. **Be specific** - Vague plans lead to confusion
2. **Consider context** - Build on what exists
3. **Plan for testing** - Tests are part of implementation
4. **Watch size** - Better to plan splits upfront
5. **Document decisions** - Explain the why, not just what

This planning phase is critical for successful implementation. A good plan makes implementation straightforward and reviewable.