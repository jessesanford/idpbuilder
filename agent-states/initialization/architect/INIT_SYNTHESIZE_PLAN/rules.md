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
Read from init-state-${PROJECT_PREFIX}.json:
- Initial project idea
- All requirements
- Technology decisions
- Architecture choices

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

### 5. Create Phased Implementation

#### Phase Structure Rules
- **Phase 1**: Foundation/MVP (30-40% of work)
- **Phase 2**: Core Features (40-50% of work)
- **Phase 3**: Polish/Enhancement (20-30% of work)

#### Wave Structure Rules
- Each phase: 2-4 waves
- Each wave: 3-6 efforts
- Each effort: Specific, measurable deliverable

#### Example Structure
```markdown
## Phase 1: Foundation
Goal: Establish core infrastructure and basic functionality

### Wave 1.1: Project Setup
- Effort 1.1.1: Repository structure and build configuration
- Effort 1.1.2: Core domain models and interfaces
- Effort 1.1.3: Basic CLI skeleton with command structure

### Wave 1.2: Core Implementation
- Effort 1.2.1: Primary business logic implementation
- Effort 1.2.2: Data persistence layer
- Effort 1.2.3: Unit test coverage for core

## Phase 2: Feature Development
Goal: Implement all primary features

### Wave 2.1: User-Facing Features
[Similar structure]

## Phase 3: Production Readiness
Goal: Polish, optimize, and prepare for deployment

### Wave 3.1: Quality & Performance
[Similar structure]
```

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
- [ ] Has project overview
- [ ] Has goals and objectives
- [ ] Has technical architecture
- [ ] Has 3 phases defined
- [ ] Each phase has 2-4 waves
- [ ] Each wave has 3-6 efforts
- [ ] Has success criteria
- [ ] Has risk mitigation
- [ ] Total scope realistic

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