# Project Implementation Plan Template

## Project Overview
**Project Name**: [Your Project Name]
**Goal**: [Clear description of what you're building]
**Total Phases**: [Number]
**Total Waves**: [Number]
**Total Efforts**: [Number]
**Estimated Completion**: [Date]

## Configuration
```yaml
size_limits:
  warning_threshold: 700
  error_threshold: 800
  
test_coverage:
  phase_1: 90%
  phase_2: 85%
  phase_3: 80%
  phase_4: 80%
  phase_5: 85%
  
parallelization:
  max_parallel_efforts: 3
  allow_parallel_waves: false
  
review_requirements:
  code_review: mandatory
  architect_review: mandatory
  security_review: optional
```

## Phase 1: Foundation
**Goal**: Establish core infrastructure and basic functionality
**Duration**: [Estimated weeks]
**Success Criteria**: 
- Core types defined
- Basic APIs functional
- Database schema complete
- CI/CD pipeline operational

### Wave 1: Core Types and Models
**Dependencies**: None
**Can Parallelize**: Yes

#### E1.1.1: Basic Data Models
- **Description**: Define fundamental data structures
- **Requirements**:
  - User model with authentication fields
  - Organization/workspace models
  - Permission models
- **Test Requirements**: Unit tests for all models
- **Estimated Size**: 400 lines

#### E1.1.2: API Interfaces
- **Description**: Define API contracts
- **Requirements**:
  - REST/GraphQL/gRPC interfaces
  - Request/response types
  - Error types
- **Test Requirements**: Schema validation tests
- **Estimated Size**: 300 lines

#### E1.1.3: Database Schema
- **Description**: Create database structure
- **Requirements**:
  - Tables/collections definition
  - Indexes and constraints
  - Migration scripts
- **Test Requirements**: Migration tests
- **Estimated Size**: 500 lines

### Wave 2: Basic Services
**Dependencies**: Wave 1 complete
**Can Parallelize**: Yes

#### E1.2.1: Authentication Service
- **Description**: User authentication implementation
- **Requirements**:
  - Login/logout functionality
  - Token management
  - Password handling
- **Test Requirements**: Security tests, unit tests
- **Estimated Size**: 600 lines

#### E1.2.2: Authorization Service
- **Description**: Permission management
- **Requirements**:
  - Role-based access control
  - Permission checking
  - Audit logging
- **Test Requirements**: Permission matrix tests
- **Estimated Size**: 500 lines

### Wave 3: Infrastructure
**Dependencies**: Waves 1-2 complete
**Can Parallelize**: No

#### E1.3.1: CI/CD Pipeline
- **Description**: Automated build and deploy
- **Requirements**:
  - Build automation
  - Test automation
  - Deployment scripts
- **Test Requirements**: Pipeline validation
- **Estimated Size**: 400 lines

## Phase 2: Core Features
**Goal**: Implement primary business functionality
**Duration**: [Estimated weeks]
**Success Criteria**:
- All core features operational
- Integration tests passing
- Performance benchmarks met

### Wave 1: Feature Set A
**Dependencies**: Phase 1 complete
**Can Parallelize**: Yes

#### E2.1.1: [Feature Name]
- **Description**: [What it does]
- **Requirements**:
  - [Requirement 1]
  - [Requirement 2]
- **Test Requirements**: [Test criteria]
- **Estimated Size**: XXX lines

[Continue pattern for all efforts...]

## Phase 3: Advanced Features
**Goal**: Add sophisticated functionality and integrations
**Duration**: [Estimated weeks]
**Success Criteria**:
- Advanced features complete
- Third-party integrations working
- Scalability proven

[Continue with waves and efforts...]

## Phase 4: Optimization
**Goal**: Performance, security, and reliability improvements
**Duration**: [Estimated weeks]
**Success Criteria**:
- Performance targets met
- Security audit passed
- Monitoring complete

[Continue with waves and efforts...]

## Phase 5: Polish and Documentation
**Goal**: Final refinements and comprehensive documentation
**Duration**: [Estimated weeks]
**Success Criteria**:
- All documentation complete
- UI/UX polished
- Ready for production

[Continue with waves and efforts...]

## Dependency Matrix

### Critical Paths
```
Phase 1 Wave 1 → Phase 1 Wave 2 → Phase 1 Wave 3
                                 ↓
                           Phase 2 Wave 1
```

### Parallel Opportunities
- Phase 1, Wave 1: All efforts can run in parallel
- Phase 2, Wave 1: E2.1.1, E2.1.2, E2.1.3 can run in parallel
- Phase 3, Wave 2: E3.2.1, E3.2.2 can run in parallel

### Blocking Dependencies
- Database schema (E1.1.3) blocks all data services
- Authentication (E1.2.1) blocks all secured endpoints
- CI/CD (E1.3.1) blocks automated deployments

## Risk Mitigation

### Technical Risks
| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Performance issues | Medium | High | Early benchmarking, profiling |
| Integration failures | Low | High | Comprehensive integration tests |
| Security vulnerabilities | Low | Critical | Security reviews, penetration testing |

### Process Risks
| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Effort size violations | Medium | Medium | Continuous measurement, early splits |
| Review bottlenecks | Medium | Medium | Parallel reviews where possible |
| Context loss | High | Low | TODO state management, documentation |

## Success Metrics

### Phase Metrics
- Phase completion on schedule: Target 90%
- First-pass review rate: Target 80%
- Size compliance: Target 100%
- Test coverage achieved: Target 100% of goals

### Overall Metrics
- Feature completeness: 100%
- Performance benchmarks: All met
- Security standards: All passed
- Documentation coverage: 100%

## Notes for Orchestrator

### Phase 1 Considerations
- Focus on getting foundations right
- Don't rush through data models
- Ensure CI/CD works early

### Phase 2 Considerations
- This is where most complexity lives
- Watch for effort size carefully
- Consider feature flags for gradual rollout

### Phase 3 Considerations
- Integration points need extra testing
- Performance testing becomes critical
- Consider canary deployments

### Phase 4 Considerations
- Don't break existing functionality
- Benchmark before and after
- Document all optimizations

### Phase 5 Considerations
- User documentation is as important as code
- Polish affects user perception significantly
- Final security audit is mandatory

## Customization Guide

### Adding Efforts
1. Assign to appropriate phase and wave
2. Estimate size conservatively
3. Define clear requirements
4. Specify test criteria
5. Note dependencies

### Adjusting Phases
1. Keep phases focused on specific goals
2. Ensure clear success criteria
3. Don't overload any single phase
4. Consider team capacity

### Setting Dependencies
1. Identify true blockers only
2. Look for parallelization opportunities
3. Document why dependencies exist
4. Review regularly for changes

---

*This plan is a living document. Update as the project evolves.*