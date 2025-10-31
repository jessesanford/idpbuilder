# IMPLEMENTATION PLAN: [Project Name]

## 🚨🚨🚨 MANDATORY TEST-DRIVEN DEVELOPMENT (TDD) 🚨🚨🚨

**THIS PROJECT FOLLOWS STRICT TDD METHODOLOGY:**
1. **WRITE TESTS FIRST** - Every feature starts with test pseudo-code
2. **RED-GREEN-REFACTOR** - Tests fail first, then pass, then optimize
3. **NO CODE WITHOUT TESTS** - Implementation forbidden until tests defined
4. **VALIDATION AT EVERY LEVEL** - Phase, wave, effort, and master tests required

## 🔍 PRIOR WORK ANALYSIS (R420 MANDATORY)

**BEFORE creating this plan, Code Reviewer MUST execute R420 Cross-Effort Planning Awareness Protocol**

This section documents comprehensive research of ALL previous implementations, plans, and sibling splits to prevent duplicate declarations, API mismatches, and method visibility errors.

### Discovery Phase Results
- **Previous Efforts Reviewed**: [List effort IDs: E1.1.0, E1.1.1, E1.2.0, etc.]
- **Sibling Splits Reviewed**: [For split plans only: split-001, split-002, etc.]
- **Previous Plans Reviewed**: [List plan files examined]
- **Research Timestamp**: [ISO 8601 timestamp when research was conducted]
- **Research Status**: [COMPLETE / INCOMPLETE]

### File Structure Findings
| File Path | Source Effort | Status | Action Required |
|-----------|---------------|--------|-----------------|
| pkg/cmd/push/root.go | E1.2.0 | EXISTS | MUST NOT create duplicate |
| pkg/cmd/push/push.go | E1.2.1 (this effort) | NEW | Safe to create |
| pkg/retry/retry.go | split-001 | EXISTS | MUST import and use |

**Instructions:**
- List ALL files that already exist in previous efforts
- Mark each as EXISTS or NEW
- For existing files: specify action (MUST NOT duplicate, MUST extend, MUST import)

### Interface/API Findings
| Interface/API | Source Effort | Signature | Action Required |
|---------------|---------------|-----------|-----------------|
| Registry | E1.1.0 | `Push(ctx context.Context, image string, content io.Reader) error` | MUST implement exactly (R373) |
| retry.Config | split-001 | `struct{MaxRetries int, BackoffDuration time.Duration}` | MUST use this type |
| Authenticator | E1.1.0 | `Authenticate(ctx context.Context, creds Credentials) error` | MUST implement this interface |

**Instructions:**
- List ALL interfaces and APIs from previous efforts
- Document EXACT signatures (critical for compatibility!)
- Specify whether to implement, use, or extend
- VERIFY these exist - do NOT assume!

### Type/Struct Findings
| Type | Source Effort | Exported | Action Required |
|------|---------------|----------|-----------------|
| GiteaClient | E1.2.0 | YES (public) | Can import and extend |
| MockRegistry | E1.1.0 | YES (public) | Can use in tests |
| configLoader | E1.1.0 | NO (private) | CANNOT import, recreate if needed |

**Instructions:**
- List major types/structs from previous efforts
- Check if exported (uppercase = YES, lowercase = NO)
- Only exported types can be imported from other packages

### Method Visibility Findings
| Method | Type | Visibility | Can Access? | Action Required |
|--------|------|------------|-------------|-----------------|
| Push() | MockRegistry | EXPORTED | YES | Safe to use in test assertions |
| Pull() | MockRegistry | EXPORTED | YES | Safe to use in test assertions |
| validate() | MockRegistry | UNEXPORTED | NO | CANNOT use from external package |
| Reset() | GiteaClient | EXPORTED | YES | Can call from tests |

**Instructions:**
- For types that will be used in tests or other packages
- Check method visibility (uppercase = EXPORTED, lowercase = UNEXPORTED)
- Only exported methods can be called from other packages
- This prevents "cannot refer to unexported field or method" errors

### Package Organization Findings
| Package Path | Source Effort | Purpose | Action Required |
|--------------|---------------|---------|-----------------|
| pkg/cmd/push | E1.2.0 | Push command hierarchy | MUST add to existing structure |
| pkg/retry | split-001 | Retry logic with backoff | MUST import, NOT recreate |
| pkg/auth | E1.1.0 | Authentication interfaces | MUST import and implement |

**Instructions:**
- Map out existing package structure
- Identify what each package provides
- Determine if new code adds to existing package or creates new one

### Conflicts Detected
- ✅ NO duplicate file paths detected
- ✅ NO API mismatches detected
- ✅ NO method visibility violations detected
- ✅ NO package conflicts detected

OR (if conflicts found):

- ❌ CONFLICT: Plan assumes retry.DefaultBackoff exists, but split-001 only exports retry.Config
- ❌ CONFLICT: Plan tries to create pkg/cmd/push/root.go which already exists in E1.2.0
- ❌ CONFLICT: Plan tries to access MockRegistry.validate() which is unexported (private)

**Instructions:**
- List ALL conflicts discovered during analysis
- Be specific about what the conflict is
- These MUST be resolved before implementation begins

### API Assumptions Verified
- ✅ VERIFIED: retry.Config exists in split-001
- ✅ VERIFIED: retry.NewConfig() function available
- ❌ INCORRECT: retry.DefaultBackoff does NOT exist (use Config instead)
- ✅ VERIFIED: MockRegistry.Push() is exported (can use in tests)
- ❌ INCORRECT: MockRegistry.validate() is unexported (cannot access from test package)

**Instructions:**
- For every API/method/type you plan to use, verify it actually exists
- Document what exists and what doesn't
- Correct any wrong assumptions BEFORE implementation

### Required Integrations
1. MUST import Registry interface from E1.1.0 and implement ALL methods with exact signatures
2. MUST use retry.Config from split-001 (NOT retry.DefaultBackoff which doesn't exist)
3. MUST extend pkg/cmd/push command structure from E1.2.0 (NOT create new root.go)
4. MUST use only exported MockRegistry methods (Push, Pull) in tests

**Instructions:**
- List everything that MUST be integrated from previous work
- Be specific about what to import, implement, or extend
- Reference the rule requiring it (R373, R374, R420)

### Forbidden Actions
- ❌ DO NOT create pkg/cmd/push/root.go (already exists in E1.2.0 - duplicate declaration)
- ❌ DO NOT use retry.DefaultBackoff (doesn't exist in split-001)
- ❌ DO NOT access MockRegistry.validate() (unexported private method)
- ❌ DO NOT create alternative Registry interface (R373 violation - must use existing)
- ❌ DO NOT recreate retry logic (exists in split-001 - must import)

**Instructions:**
- List everything that is FORBIDDEN based on prior work analysis
- Include the reason why it's forbidden
- This prevents the three failure types: duplicates, API mismatches, visibility errors

---

**R420 VALIDATION REQUIREMENT:**
This plan MUST pass `tools/validate-R420-compliance.sh` validation before approval.

If any required sections are missing or incomplete, plan CANNOT be approved.

---

## 1. Project Overview

[Provide a comprehensive 2-3 paragraph description of the project, expanding on the initial idea. Include the problem being solved, the approach being taken, and the expected impact.]

### Key Features
- [Feature 1]: [Brief description]
- [Feature 2]: [Brief description]
- [Feature 3]: [Brief description]

### Target Users
[Describe who will use this project and how]

## 2. Goals and Objectives

### Primary Objectives
1. [Objective 1]: [Measurable goal]
2. [Objective 2]: [Measurable goal]
3. [Objective 3]: [Measurable goal]

### Secondary Goals
- [Goal 1]
- [Goal 2]

### Non-Goals (Out of Scope)
- [What this project will NOT do]
- [Boundaries and limitations]

## 3. Technical Architecture

### Technology Stack
- **Primary Language**: [Language]
- **Framework**: [Main framework]
- **Build System**: [Build tool]
- **Testing**: [Test framework]
- **Database**: [If applicable]
- **Deployment**: [Target environment]

### Architecture Pattern
[Describe the chosen architecture pattern and why it's appropriate]

### System Components
1. **[Component Name]**: [Description and responsibility]
2. **[Component Name]**: [Description and responsibility]
3. **[Component Name]**: [Description and responsibility]

### External Dependencies
- [Dependency 1]: [Purpose]
- [Dependency 2]: [Purpose]

## 4. MASTER END-TO-END TEST SUITE

### 🎯 PROJECT PROJECT_DONE DEFINITION TEST
**This test PROVES the entire project works as designed:**

```pseudo
test_project_complete_end_to_end():
    # Setup
    [Initialize the complete system]
    [Setup test data and environment]

    # Core functionality validation
    [Test primary feature 1 works correctly]
    [Test primary feature 2 works correctly]
    [Test all features integrate properly]

    # Performance validation
    [Assert response times < threshold]
    [Assert resource usage within limits]

    # Error handling validation
    [Test graceful degradation]
    [Test error recovery]

    # Cleanup
    [Verify clean shutdown]
    [Assert no resource leaks]

    # EXPECTED: All assertions pass, proving project success
```

## 5. Implementation Phases (TDD REQUIRED)

### Phase 1: Foundation and Core Infrastructure
**Goal**: Establish the basic project structure and core functionality
**Duration Estimate**: [Estimate]

#### 📝 PHASE 1 VALIDATION TESTS (WRITE FIRST!)
```pseudo
# Test 1: Core infrastructure operational
test_phase1_infrastructure_complete():
    [Verify build system works]
    [Verify dependency management]
    [Assert all core components initialize]
    [Verify basic CLI/API responds]
    EXPECT: Infrastructure fully operational

# Test 2: Domain models functional
test_phase1_domain_models():
    [Create test instances of each model]
    [Verify serialization/deserialization]
    [Test validation rules]
    [Assert relationships work correctly]
    EXPECT: All models work as designed

# Test 3: End-to-end basic flow
test_phase1_basic_flow():
    [Initialize system]
    [Perform basic operation]
    [Verify output]
    EXPECT: Basic functionality works end-to-end
```

#### Wave 1.1: Project Setup and Basic Structure

##### 🧪 WAVE 1.1 TESTS (MANDATORY BEFORE IMPLEMENTATION)
```pseudo
# Unit tests for build system
test_build_configuration():
    [Run build command]
    [Assert build succeeds]
    [Verify output artifacts]

# Integration test for structure
test_project_structure():
    [Verify all directories exist]
    [Check configuration files present]
    [Assert dependencies resolved]
```
- **Effort 1.1.1**: Repository structure and build configuration
  - Set up directory structure
  - Configure build system
  - Initialize dependency management

- **Effort 1.1.2**: Core domain models and interfaces
  - Define primary data structures
  - Create interface definitions
  - Establish naming conventions

- **Effort 1.1.3**: Basic CLI/API skeleton
  - Implement entry point
  - Set up command/route structure
  - Add configuration loading

#### Wave 1.2: Core Business Logic

##### 🧪 WAVE 1.2 TESTS (MANDATORY BEFORE IMPLEMENTATION)
```pseudo
# Unit tests for business logic
test_core_business_rules():
    [Test each business rule in isolation]
    [Verify edge cases handled]
    [Assert expected outputs]

# Integration tests for logic flow
test_business_logic_integration():
    [Test complete business workflows]
    [Verify data transformations]
    [Assert state changes correct]
```
- **Effort 1.2.1**: [Specific feature implementation]
- **Effort 1.2.2**: [Specific feature implementation]
- **Effort 1.2.3**: Unit test coverage for core

#### Wave 1.3: [Additional wave if needed]
[Similar structure]

### Phase 2: Feature Development
**Goal**: Implement all primary features and capabilities
**Duration Estimate**: [Estimate]

#### 📝 PHASE 2 VALIDATION TESTS (WRITE FIRST!)
```pseudo
# Test 1: All features functional
test_phase2_all_features_work():
    [Test each feature independently]
    [Verify feature interactions]
    [Assert no feature conflicts]
    EXPECT: All features operational

# Test 2: Integration complete
test_phase2_integration():
    [Test feature combinations]
    [Verify data flow between features]
    [Assert system coherence]
    EXPECT: Features work together seamlessly

# Test 3: Performance acceptable
test_phase2_performance():
    [Measure feature response times]
    [Check resource consumption]
    [Verify concurrent operation]
    EXPECT: Performance within requirements
```

#### Wave 2.1: [Feature Category]
- **Effort 2.1.1**: [Specific implementation]
- **Effort 2.1.2**: [Specific implementation]
- **Effort 2.1.3**: [Specific implementation]

#### Wave 2.2: [Feature Category]
[Similar structure]

### Phase 3: Production Readiness
**Goal**: Polish, optimize, and prepare for deployment
**Duration Estimate**: [Estimate]

#### 📝 PHASE 3 VALIDATION TESTS (WRITE FIRST!)
```pseudo
# Test 1: Production deployment ready
test_phase3_deployment_ready():
    [Verify deployment scripts work]
    [Test rollback procedures]
    [Assert monitoring active]
    EXPECT: Deployment successful

# Test 2: Security hardened
test_phase3_security():
    [Run security scans]
    [Test authentication/authorization]
    [Verify data encryption]
    EXPECT: No security vulnerabilities

# Test 3: Production performance
test_phase3_production_performance():
    [Load test with expected traffic]
    [Stress test beyond limits]
    [Verify graceful degradation]
    EXPECT: Production-ready performance
```

#### Wave 3.1: Quality and Performance
- **Effort 3.1.1**: Performance optimization
- **Effort 3.1.2**: Security hardening
- **Effort 3.1.3**: Comprehensive testing

#### Wave 3.2: Documentation and Deployment
- **Effort 3.2.1**: User documentation
- **Effort 3.2.2**: Deployment automation
- **Effort 3.2.3**: Monitoring and logging setup

## 6. Success Criteria (VALIDATED BY TESTS)

### Phase 1 Completion Criteria
- [ ] Core functionality working end-to-end
- [ ] Build system fully configured
- [ ] Unit tests passing with >70% coverage
- [ ] Basic documentation in place

### Phase 2 Completion Criteria
- [ ] All planned features implemented
- [ ] Integration tests passing
- [ ] Performance benchmarks met
- [ ] API/CLI fully functional

### Phase 3 Completion Criteria
- [ ] Production deployment successful
- [ ] Security audit passed
- [ ] Documentation complete
- [ ] Monitoring operational
- [ ] Load testing passed

### Overall Project Success Metrics
- [Metric 1]: [Target value]
- [Metric 2]: [Target value]
- [Metric 3]: [Target value]

## 7. TDD Workflow Requirements

### MANDATORY TDD CYCLE FOR EVERY EFFORT:
1. **RED PHASE**: Write failing tests FIRST
   - Define what success looks like
   - Write comprehensive test cases
   - Run tests - they MUST fail initially

2. **GREEN PHASE**: Write minimal code to pass
   - Implement ONLY enough to pass tests
   - No extra features or optimizations
   - All tests must pass

3. **REFACTOR PHASE**: Improve code quality
   - Clean up implementation
   - Optimize performance
   - Tests must still pass

### TEST COVERAGE REQUIREMENTS:
- **Unit Tests**: Minimum 80% code coverage
- **Integration Tests**: All component interactions
- **End-to-End Tests**: Critical user journeys
- **Performance Tests**: Load and stress testing
- **Security Tests**: Vulnerability scanning

## 8. Risk Mitigation

### Technical Risks
1. **Risk**: [Description of risk]
   - **Impact**: [High/Medium/Low]
   - **Mitigation**: [Strategy to address]

2. **Risk**: [Description of risk]
   - **Impact**: [High/Medium/Low]
   - **Mitigation**: [Strategy to address]

### Schedule Risks
1. **Risk**: [Description of risk]
   - **Impact**: [High/Medium/Low]
   - **Mitigation**: [Strategy to address]

### External Dependencies
1. **Risk**: [Description of risk]
   - **Impact**: [High/Medium/Low]
   - **Mitigation**: [Strategy to address]

## 9. Appendices

### A. Glossary
- [Term]: [Definition]
- [Term]: [Definition]

### B. References
- [Reference document or link]
- [Reference document or link]

### C. Assumptions
- [Assumption about environment/requirements]
- [Assumption about resources/timeline]

---

## 10. TDD Compliance Checklist

### Before Starting ANY Implementation:
- [ ] Master end-to-end test pseudo-code written
- [ ] All phase validation tests defined
- [ ] Wave-level tests specified where applicable
- [ ] Test data and fixtures planned
- [ ] Test environment requirements documented
- [ ] CI/CD pipeline includes all tests

### During Implementation:
- [ ] Following RED-GREEN-REFACTOR cycle
- [ ] Tests written BEFORE code
- [ ] All tests passing before merge
- [ ] Code coverage meets requirements
- [ ] Performance tests validating

### Project Completion:
- [ ] Master end-to-end test PASSES
- [ ] All phase validation tests PASS
- [ ] Test documentation complete
- [ ] Test maintenance plan defined

---

*This TDD-focused plan is generated by Software Factory 2.0 initialization process. Tests MUST be written FIRST before any implementation begins.*