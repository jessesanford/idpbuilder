# Architect - INIT_SYNTHESIZE_PLAN State Rules

## Purpose
Create comprehensive IMPLEMENTATION-PLAN.md from gathered requirements.

## Entry Criteria
- All requirements gathered and stored
- Configuration files generated
- Repository setup complete
- Ready to synthesize into plan

## Required Actions

### 1. Load All Inputs
Read from initialization files:
- `init-state-temp.json`: Initial project idea and description
- `init-decomposition.json`: **CRITICAL** - Phase/wave/effort structure from INIT_DECOMPOSE_PRD
- `init-sizing-justifications.md`: LOC estimates and justifications
- Requirements gathered during INIT_REQUIREMENTS_GATHERING

**IMPORTANT**: The decomposition has already been done! Your job is to SYNTHESIZE it into narrative form, not redo the decomposition.

### 2. Expand Project Overview
Transform initial idea into comprehensive overview:
- Problem statement (2-3 paragraphs)
- Solution approach
- Expected impact
- Key innovations

### 3. Define Goals and Objectives
From requirements, create:
- Primary objectives (3-5)
- Secondary goals
- Success metrics
- Non-goals/out-of-scope

### 4. Document Technical Architecture
Structure as:
```markdown
## Technical Architecture

### Technology Stack
- Language: [primary language]
- Framework: [primary framework]
- Build System: [build tool]
- Testing: [test framework]

### Architecture Pattern
[Description of chosen pattern]

### Key Components
1. [Component 1]: [Description]
2. [Component 2]: [Description]

### Integration Points
- [External system/API]
```

### 5. Transform Decomposition into Implementation Plan

**DO NOT REINVENT**: Use the structure from `init-decomposition.json` directly!

#### Your Job
Transform the JSON decomposition into rich, narrative markdown:

1. **For each phase**: Add narrative description of phase goal
2. **For each wave**: Add context on why these efforts are grouped
3. **For each effort**: Expand description with implementation details
4. **Add sizing section**: Include LOC estimates with justifications
5. **Add dependency section**: Show dependency relationships explicitly

#### Example Transformation

**From init-decomposition.json**:
```json
{
  "effort_id": "E1.1.1",
  "name": "Push command skeleton",
  "description": "Add push subcommand to idpbuilder CLI",
  "estimated_loc": 150
}
```

**To IMPLEMENTATION-PLAN.md**:
```markdown
#### Effort 1.1.1: Push Command Skeleton
**Estimated Size**: 150 lines
**Dependencies**: None (foundation effort)

**Description**:
Implement the `push` subcommand for idpbuilder using the Cobra CLI framework.
This establishes the command structure that all subsequent push functionality
will build upon.

**Deliverables**:
- Cobra command definition for `idpbuilder push`
- Flag definitions: `-username`, `-password`, `-registry-url`, `-insecure`
- Help text and usage examples
- Basic flag validation (required flags present, URL format)

**Sizing Justification**:
- Cobra command struct: 30 lines
- Flag setup and binding: 50 lines
- RunE function skeleton: 40 lines
- Help text and examples: 30 lines
**Total**: 150 lines

**Testing Notes**: Unit tests will be in separate effort (E1.1.5)
```

#### Structure Template
Use exact structure from decomposition:
- Same phase numbers, names, goals
- Same wave numbers, names, efforts
- Same effort IDs and names
- Add narrative and details around the structure

### 6. Define Success Criteria
Measurable criteria for each phase:
```markdown
## Success Criteria

### Phase 1 Completion
- [ ] All tests passing with >70% coverage
- [ ] Core functionality demonstrable
- [ ] Documentation for basic usage

### Phase 2 Completion
- [ ] All planned features implemented
- [ ] Integration tests passing
- [ ] Performance benchmarks met

### Phase 3 Completion
- [ ] Production deployment successful
- [ ] Load testing passed
- [ ] Security audit complete
```

### 7. Risk Mitigation
Identify and address:
```markdown
## Risk Mitigation

### Technical Risks
1. **Risk**: [Description]
   **Mitigation**: [Strategy]

### Schedule Risks
1. **Risk**: [Description]
   **Mitigation**: [Strategy]
```

### 8. Generate Complete Plan
Write IMPLEMENTATION-PLAN.md with:
- All sections properly formatted
- Consistent effort naming (E1.1.1 style)
- Clear deliverables for each effort
- Realistic scope for 700-line limit
- Dependencies noted where applicable

## Plan Validation Checklist
- [ ] Has project overview (expanded from initial description)
- [ ] Has goals and objectives
- [ ] Has technical architecture
- [ ] **CRITICAL**: Uses exact phase/wave/effort structure from init-decomposition.json
- [ ] Every effort includes estimated LOC with justification
- [ ] All efforts are ≤800 lines (verify against decomposition)
- [ ] Dependencies explicitly documented for each effort
- [ ] Has success criteria per phase
- [ ] Has risk mitigation
- [ ] Total scope matches decomposition
- [ ] Narrative adds value beyond raw JSON structure

## Exit Criteria
- Complete IMPLEMENTATION-PLAN.md written
- All sections populated
- Plan validated against template
- Ready for agent customization

## Transition
**MANDATORY**: → INIT_CUSTOMIZE_AGENTS

## Quality Standards

### Good Effort Definition
```
✅ GOOD: "Implement user authentication with JWT tokens and role-based access control"
- Specific deliverable
- Clear scope
- Measurable completion

❌ BAD: "Work on authentication"
- Vague scope
- Unclear deliverable
- Not measurable
```

### Effort Sizing
- Target: 400-600 lines per effort
- Maximum: 700 lines (hard limit)
- If larger: Split into multiple efforts

## Time Guidance
- Plan synthesis: 10-15 minutes
- Should be thorough but not exhaustive
- Focus on clarity and actionability

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

