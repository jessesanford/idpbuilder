# Phase X: [Phase Name] - Detailed Implementation Plan

## Phase Overview
**Duration:** [X] days  
**Critical Path:** [YES/NO] - [Explanation why this phase is/isn't critical]  
**Base Branch:** `[previous-phase-integration or main]`  
**Target Integration Branch:** `phase[X]-integration`

---

## Wave X.1: [Wave Name/Description]

### E[X.1.1]: [Effort Name]
**Branch:** `/phase[X]/wave[Y]/effort[Z]-[descriptive-name]`  
**Duration:** [X] hours  
**Agent:** [Single agent / Can parallelize with X agents]  
**Dependencies:** [List specific efforts this depends on]

#### Source Material:
```markdown
# Option 1: Reusing existing code
- Primary: `origin/feature/[branch-with-implementation]`
- Secondary: `origin/feature/[fallback-branch]`
- Commits: [List specific commits if known]

# Option 2: Porting from external source
- Source: `[external-repo]/[path]`
- Reference: [Documentation/Design doc URL]

# Option 3: New development
- Design Doc: [Link to design document]
- Reference Implementation: [If any]
```

#### Specific Commits to Cherry-Pick:
```bash
# List exact commits or indicate new development
git cherry-pick [commit-hash]  # Description of what this commit does
git cherry-pick [commit-hash]  # Description of what this commit does

# Or for new development:
# NEW DEVELOPMENT - No existing commits to cherry-pick
```

#### Requirements:
1. **MUST** implement:
   - [Specific requirement 1]
   - [Specific requirement 2]
   - [Specific requirement 3]

2. **MUST** support:
   - [Feature/capability 1]
   - [Feature/capability 2]

3. **MUST NOT**:
   - [Thing to avoid 1]
   - [Thing to avoid 2]

4. **SHOULD** (nice to have):
   - [Optional feature 1]
   - [Optional feature 2]

#### Test Requirements (TDD):
```[language]
// test/[component]/[name]_test.[ext]
// Test Suite 1: Basic Functionality
func/describe Test[ComponentName](t *testing.T) {
    // Test Case 1: Happy path
    test("should successfully [do something]", () => {
        // Given: Initial state
        [setup code]
        
        // When: Action is performed
        [action code]
        
        // Then: Expected outcome
        [assertion code]
    })
    
    // Test Case 2: Error handling
    test("should handle [error condition]", () => {
        // Given: Error condition setup
        [setup code]
        
        // When: Action that triggers error
        [action code]
        
        // Then: Graceful error handling
        [assertion code]
    })
    
    // Test Case 3: Edge cases
    test("should handle [edge case]", () => {
        [test implementation]
    })
}

// Test Suite 2: Integration
func/describe TestIntegration[ComponentName](t *testing.T) {
    // Test integration with other components
}

// Test Suite 3: Performance (if applicable)
func/describe Benchmark[ComponentName](b *testing.B) {
    // Benchmark critical paths
}
```

#### Pseudo-Code Implementation:
```
FUNCTION implement_[effort_name]():
    // Step 1: [Setup/Initialization]
    [PSEUDO CODE describing setup steps]
    
    // Step 2: [Core Implementation]
    IF reusing_existing_code:
        CHERRY_PICK specified_commits
        RESOLVE_CONFLICTS with_strategy
        ADAPT_CODE to_current_structure
    ELSE:
        IMPLEMENT core_logic:
            [DETAILED PSEUDO CODE]
    
    // Step 3: [Integration]
    WIRE_UP with_existing_components
    ADD error_handling
    ADD logging
    ADD metrics
    
    // Step 4: [Validation]
    RUN tests
    CHECK coverage
    VERIFY performance
```

#### Size Estimation and Split Strategy:
```yaml
estimated_lines: [X]
split_threshold: 800

if_exceeds_threshold:
  suggested_splits:
    - Part 1: [Logical component 1] (~X lines)
    - Part 2: [Logical component 2] (~X lines)
    - Part 3: [Tests and documentation] (~X lines)
  
  split_criteria:
    - Each split must build independently
    - Maintain logical cohesion
    - Tests stay with implementation
```

#### Validation Commands:
```bash
# Build validation
[build command] || exit 1

# Test execution
[test command] || exit 1

# Coverage check
[coverage command]
# Requirement: >[X]% coverage for this phase

# Lint/Format check
[lint command] || exit 1

# Line count verification
/workspaces/[project]/tools/line-counter.sh -c $(git branch --show-current)
# MUST be < 800 lines (unless exception granted)

# Performance validation (if applicable)
[benchmark command]
# Must meet: [specific performance criteria]

# Integration test (if applicable)
[integration test command]
```

#### Success Criteria:
- [ ] All tests pass (>[X]% coverage)
- [ ] Build succeeds without warnings
- [ ] Line count within limit (<800 or exception documented)
- [ ] Lint checks pass
- [ ] Performance benchmarks meet targets (if applicable)
- [ ] Documentation updated
- [ ] Integration tests pass (if applicable)
- [ ] Code review completed and approved

#### Rollback Plan:
```bash
# If effort fails validation
git checkout [base-branch]
git branch -D [effort-branch]
# Document failure reason in orchestrator-state.yaml
# Retry with fixes or escalate to architect
```

---

### E[X.1.2]: [Next Effort Name]
**Branch:** `/phase[X]/wave[Y]/effort[Z]-[descriptive-name]`  
**Duration:** [X] hours  
**Dependencies:** [Dependencies]

[Continue pattern for all efforts in wave...]

---

## Wave X.2: [Next Wave Name]

[Continue pattern for all waves in phase...]

---

## Wave Dependency Graph

```mermaid
graph TD
    subgraph "Wave X.1"
        E[X.1.1][Effort Name] 
        E[X.1.2][Effort Name]
    end
    
    subgraph "Wave X.2"
        E[X.2.1][Effort Name]
        E[X.2.2][Effort Name]
    end
    
    E[X.1.1] --> E[X.2.1]
    E[X.1.2] --> E[X.2.1]
    E[X.1.1] -.-> E[X.1.2]
    
    classDef critical fill:#f96
    classDef parallel fill:#9f6
    classDef sequential fill:#69f
    
    class E[X.1.1] critical
    class E[X.1.2],E[X.2.2] parallel
    class E[X.2.1] sequential
```

Legend:
- Red: Critical path
- Green: Can parallelize  
- Blue: Must be sequential
- Solid arrow: Hard dependency
- Dotted arrow: Soft dependency

## Dependency Table

| Effort | Depends On | Can Parallelize With | Blocks | Critical Path |
|--------|------------|---------------------|--------|---------------|
| E[X.1.1] | None/[Previous phase] | E[X.1.2] | E[X.2.1] | Yes |
| E[X.1.2] | None | E[X.1.1] | None | No |
| E[X.2.1] | E[X.1.1] | None | E[X.2.2] | Yes |
| E[X.2.2] | E[X.2.1] | E[X.2.3] | [Next phase] | No |

## Integration Strategy

### Wave Integration:
```bash
# After each wave completes
cd /workspaces/[project]
git checkout [base-branch]
git checkout -b phase[X]/wave[Y]-integration

# Merge efforts in dependency order
for effort in $(cat wave[Y]-efforts.txt); do
    echo "Merging effort: $effort"
    git merge --no-ff $effort -m "Integrate $effort into wave [Y]"
    
    # Validate after each merge
    [build command] || exit 1
    [test command] || exit 1
done

# Wave validation
[integration test command]
```

### Phase Integration:
```bash
# After all waves complete
git checkout [base-branch]
git checkout -b phase[X]-integration

# Merge waves
for wave in wave1-integration wave2-integration wave3-integration; do
    git merge --no-ff phase[X]/$wave
done

# Final phase validation
[full test suite]
[performance benchmarks]
[integration tests]
```

## Risk Assessment

| Risk | Probability | Impact | Mitigation | Contingency |
|------|------------|--------|------------|-------------|
| [Dependency unavailable] | Low/Med/High | Low/Med/High | [Mitigation strategy] | [Fallback plan] |
| [Performance regression] | Medium | High | Benchmark each effort | Rollback and optimize |
| [Merge conflicts] | Low | Medium | Isolated namespaces | Manual resolution |
| [Effort exceeds size limit] | High | Low | Plan splits upfront | Implement split protocol |
| [Test coverage drops] | Low | Medium | TDD approach | Add tests before merge |

## Performance Requirements

| Metric | Target | Measurement | Enforcement |
|--------|--------|-------------|-------------|
| [Operation] latency | <[X]ms p99 | Benchmark test | Gate in CI |
| Memory usage | <[X]MB | pprof analysis | Review check |
| CPU usage | <[X]% | Profile | Benchmark |
| Throughput | >[X] ops/sec | Load test | Integration test |

## Testing Strategy

### Coverage Requirements:
- Phase [X] Overall: [X]%
- Critical paths: [X]%
- Error handling: [X]%

### Test Levels:
1. **Unit Tests** (per effort): Fast, isolated, high coverage
2. **Integration Tests** (per wave): Component interaction
3. **System Tests** (per phase): End-to-end validation
4. **Performance Tests** (critical paths): Benchmarks and profiling

## Documentation Requirements

Each effort must update:
- [ ] API documentation (if APIs changed)
- [ ] README files (if new components)
- [ ] Architecture diagrams (if structure changed)
- [ ] Configuration examples (if new configs)
- [ ] Migration guides (if breaking changes)

## Notes for Orchestrator

1. **Parallelization Opportunities**: [List which efforts can run simultaneously]
2. **Critical Path Items**: [List efforts that block progress]
3. **Resource Requirements**: [Special needs like test clusters, external services]
4. **Review Requirements**: [Which efforts need architect review]
5. **Integration Checkpoints**: [When to run integration tests]
6. **Performance Gates**: [Which efforts have performance requirements]
7. **Split Candidates**: [Which efforts likely need splitting]

## Phase Completion Checklist

- [ ] All efforts completed and reviewed
- [ ] All waves integrated successfully
- [ ] Phase integration branch created
- [ ] All tests passing (coverage >[X]%)
- [ ] Performance benchmarks meet targets
- [ ] Documentation updated
- [ ] Architecture review completed (if required)
- [ ] No outstanding TODOs or FIXMEs
- [ ] Ready for next phase dependencies