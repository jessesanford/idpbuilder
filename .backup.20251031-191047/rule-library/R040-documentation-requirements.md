# Rule R040: Documentation Requirements Protocol

## Rule Statement
All significant work MUST be documented with clear, actionable reports. Documentation must include what was done, why decisions were made, issues encountered, and next steps. Reports must be created for reviews, integrations, phase completions, and error recoveries.

## Criticality Level
**WARNING** - Poor documentation leads to repeated work and lost context

## Enforcement Mechanism
- **Technical**: Mandatory report generation at checkpoints
- **Behavioral**: Reports before state transitions
- **Grading**: -15% for missing documentation, -30% for inadequate documentation

## Core Principle

```
Documentation = Context → Actions → Results → Next Steps
ALWAYS document decisions and rationale
ALWAYS include actionable next steps
NEVER leave work undocumented
Reports enable continuity and recovery
```

## Detailed Requirements

### Required Documentation Points

1. **Code Review Reports**
   ```markdown
   # CODE REVIEW REPORT - [Branch/Effort Name]
   Date: [ISO timestamp]
   Reviewer: [Agent name]
   
   ## Summary
   - Files reviewed: X
   - Issues found: Y
   - Status: PASS/FAIL/NEEDS_FIXES
   
   ## Critical Issues
   1. [Issue description]
      - File: [path]
      - Line: [numbers]
      - Fix: [required action]
   
   ## Recommendations
   - [Actionable recommendation]
   
   ## Next Steps
   - [Specific next action]
   ```

2. **Integration Reports**
   ```markdown
   # INTEGRATE_WAVE_EFFORTS REPORT - [Wave/Phase]
   Date: [ISO timestamp]
   Integrator: [Agent name]
   
   ## Branches Integrated
   - [branch-name]: [status]
   
   ## Conflicts Resolved
   - [File]: [resolution strategy]
   
   ## Test Results
   - [Test suite]: [PASS/FAIL]
   
   ## Issues Requiring Fixes
   - [Issue]: [Owner] [Priority]
   ```

3. **Phase Completion Documentation**
   ```markdown
   # PHASE COMPLETION - [Phase Name]
   
   ## Deliverables Completed
   - [Feature/component]: [Status]
   
   ## Known Issues
   - [Issue]: [Impact] [Mitigation]
   
   ## Lessons Learned
   - [Learning]: [Application]
   
   ## Ready for Next Phase
   - Prerequisites met: YES/NO
   - Blocking issues: [List or NONE]
   ```

### Documentation Standards

1. **Clarity**: Use clear, concise language
2. **Actionability**: Include specific next steps
3. **Completeness**: Cover all relevant aspects
4. **Timeliness**: Create within 5 minutes of completion
5. **Accessibility**: Store in designated locations

### File Locations

```bash
# Standard documentation paths
reports/
├── reviews/
│   └── [effort-name]-review-[timestamp].md
├── integration/
│   └── [wave/phase]-integration-[timestamp].md
├── phases/
│   └── phase-[number]-completion.md
└── errors/
    └── error-[timestamp]-recovery.md
```

### Minimum Content Requirements

Every report MUST include:
- **Timestamp**: ISO format date/time
- **Author**: Agent type and identifier
- **Context**: What triggered this documentation
- **Actions**: What was done
- **Results**: Outcome of actions
- **Next Steps**: Clear, actionable items

## Relationship to Other Rules
- **R263**: Integration documentation requirements
- **R257**: Phase assessment report
- **R258**: Wave review report
- **R303**: Document location protocol

## Implementation Notes
- Reports must be committed within 2 minutes of creation
- Use markdown format for all documentation
- Include code snippets where relevant
- Link to relevant files and commits