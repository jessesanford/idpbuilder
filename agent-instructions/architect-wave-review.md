# Architect Wave Review Task

## Task: Review Wave [WAVE] Completion for Phase [PHASE]

### Context
- **Phase**: [PHASE] - [PHASE_NAME]
- **Wave**: [WAVE]
- **Total Efforts in Wave**: [EFFORT_COUNT]
- **Integration Branch**: [INTEGRATION_BRANCH]

### Your Mission

Review the completed wave for architectural consistency, integration readiness, and compliance with system design principles.

### Required Reading

1. **Review protocols:**
   - `../../../protocols/WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md`
   - `../../../protocols/ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md`
   - `../../../core/SOFTWARE-FACTORY-STATE-MACHINE.md`
   - `../../../orchestrator-state.yaml` - Current state

2. **Wave implementation details:**
   - Review each effort's implementation
   - Check integration branch status
   - Verify architectural patterns

### Review Checklist

#### 1. Compliance Verification
- [ ] All efforts have line count <800
- [ ] All efforts have been code reviewed
- [ ] All efforts have passing tests
- [ ] Integration branch created successfully

#### 2. Architectural Review

**Design Patterns**
- [ ] Consistent patterns across efforts
- [ ] Proper separation of concerns
- [ ] Clean interfaces between components
- [ ] No architectural anti-patterns

**Integration Quality**
- [ ] Components integrate cleanly
- [ ] No conflicting implementations
- [ ] Shared resources properly managed
- [ ] Dependencies correctly handled

**Scalability & Performance**
- [ ] Design supports expected scale
- [ ] No obvious performance bottlenecks
- [ ] Resource usage appropriate
- [ ] Caching strategy consistent

**Security & Best Practices**
- [ ] Security patterns properly implemented
- [ ] No exposed sensitive data
- [ ] Input validation present
- [ ] Error handling comprehensive

#### 3. Cross-Effort Consistency
- [ ] Naming conventions consistent
- [ ] Logging approach unified
- [ ] Error handling patterns match
- [ ] Configuration management aligned

### Review Outcomes

Determine one of three verdicts:

#### 1. PROCEED ✅
Use when:
- All architectural standards met
- Integration successful
- Ready for next wave

#### 2. CHANGES_REQUIRED 🔄
Use when:
- Architectural issues found
- Integration problems exist
- Fixes needed before proceeding

List specific issues and required fixes.

#### 3. STOP 🛑
Use when:
- Major architectural problems
- Fundamental design issues
- Requires significant rework

Provide detailed explanation and remediation plan.

### Output Format

Create `WAVE-[WAVE]-ARCHITECT-REVIEW.md`:

```markdown
# Architect Review: Phase [PHASE] Wave [WAVE]

## Summary
- **Verdict**: [PROCEED/CHANGES_REQUIRED/STOP]
- **Review Date**: [date]
- **Efforts Reviewed**: [count]
- **Integration Status**: [status]

## Architectural Assessment

### ✅ Strengths
- [Positive architectural aspects]
- [Good patterns observed]

### ⚠️ Concerns
- [Architectural issues]
- [Integration problems]
- [Pattern violations]

## Effort-by-Effort Review

### Effort 1: [Name]
- **Compliance**: ✅/❌
- **Architecture**: [comments]
- **Integration**: [status]

### Effort 2: [Name]
...

## Integration Analysis
- **Branch**: [branch name]
- **Merge Conflicts**: [none/resolved/pending]
- **Integration Tests**: [status]

## [If CHANGES_REQUIRED] Required Changes

### High Priority
1. [Critical issue to fix]
2. [Another critical issue]

### Medium Priority
1. [Important but not blocking]

### Low Priority
1. [Nice to have improvements]

## [If STOP] Remediation Plan
[Detailed plan to address fundamental issues]

## Recommendations for Next Wave
- [Architectural guidance]
- [Patterns to follow]
- [Things to avoid]

## Sign-off
- Architecture Review: [COMPLETE/PENDING]
- Integration Review: [COMPLETE/PENDING]
- Ready for Next Wave: [YES/NO]
```

### Decision Criteria

**PROCEED when:**
- All efforts properly architected
- Integration successful
- No blocking issues
- Minor issues can be addressed later

**CHANGES_REQUIRED when:**
- Fixable architectural issues
- Integration needs adjustments
- Pattern inconsistencies
- Can be resolved without major rework

**STOP when:**
- Fundamental architecture flaws
- Major integration failures
- Security vulnerabilities
- Requires redesign

### Success Criteria

- [ ] All efforts reviewed for architecture
- [ ] Integration branch verified
- [ ] Clear verdict provided
- [ ] Specific feedback documented
- [ ] Path forward defined

### Remember

- Focus on architecture, not implementation details
- Consider long-term maintainability
- Verify integration actually works
- Provide actionable feedback
- Think about system-wide impact