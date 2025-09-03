# Effort Implementation Plan: [EFFORT_NAME]

<!-- NOTE: [PROJECT_PREFIX/] should be replaced with the actual project prefix from target-repo-config.yaml 
     If project_prefix is "tmc-workspace", branches become: tmc-workspace/phase1/wave1/effort-name
     If project_prefix is empty, branches become: phase1/wave1/effort-name -->

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort**: [EFFORT_NUMBER] - [EFFORT_NAME]  
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-[name]`  
**Can Parallelize**: [Yes/No] (COPIED FROM WAVE PLAN)  
**Parallel With**: [List effort numbers or "None"] (COPIED FROM WAVE PLAN)  
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: [List dependent efforts] (COPIED FROM WAVE PLAN)  

## 📋 Source Information
**Wave Plan**: PHASE-[PHASE]-WAVE-[WAVE]-IMPLEMENTATION-PLAN.md  
**Effort Section**: Effort [NUMBER]  
**Created By**: Code Reviewer Agent  
**Date**: [DATE]  
**Extracted**: [TIMESTAMP]  

## 🚀 Parallelization Context
**Can Parallelize**: [Yes/No]  
**Parallel With**: [Efforts that can run simultaneously]  
**Blocking Status**: [If No - this effort blocks: X, Y, Z]  
**Parallel Group**: [If Yes - member of parallel group: A, B, C]  
**Orchestrator Guidance**: [When orchestrator should spawn this effort]  

## 📁 Files to Create

### Primary Implementation Files
```yaml
new_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/[feature]/[file1].go
    lines: ~[NUMBER]
    purpose: [PURPOSE]
    exports:
      - FunctionName
      - TypeName
      
  - path: pkg/phase[PHASE]/wave[WAVE]/[feature]/[file2].go
    lines: ~[NUMBER]
    purpose: [PURPOSE]
    imports_from_effort:
      - ../api/interfaces.go (from Effort 1)
      - ../lib/client.go (from Effort 2)
```

### Test Files
```yaml
test_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/[feature]/[file]_test.go
    coverage_target: 80%
    test_types:
      - unit
      - integration
```

## 📦 Files to Import/Reuse

### From Previous Efforts (This Wave)
```yaml
this_wave_imports:
  - source: pkg/phase[PHASE]/wave[WAVE]/api/interfaces.go
    from_effort: 1
    usage: Implement Service interface
    
  - source: pkg/phase[PHASE]/wave[WAVE]/lib/client.go
    from_effort: 2
    usage: Use client for API calls
```

### From Previous Waves/Phases
```yaml
previous_work_imports:
  - source: pkg/phase1/common/logger/logger.go
    usage: Logging infrastructure
    
  - source: pkg/phase[PHASE]/wave[PREV]/lib/base.go
    usage: Extend base functionality
```

## 🔗 Dependencies

### Effort Dependencies
- **Must Complete First**: Effort [N] - [Reason]
- **Can Run in Parallel With**: Efforts [X, Y, Z]
- **Blocks**: Efforts [A, B] - [They need this effort's outputs]

### Technical Dependencies
- APIs from Effort 1 (contracts)
- Libraries from Effort 2 (shared code)
- [Other technical dependencies]

## 📝 Implementation Instructions

### Step-by-Step Guide
1. **Setup**
   - Verify you're in correct directory: `efforts/phase[PHASE]/wave[WAVE]/[effort]`
   - Confirm on correct branch: `phase[PHASE]/wave[WAVE]/effort-[name]`
   - Read IMPLEMENTATION-PLAN.md (this file)

2. **Implementation Order**
   - Start with [file1] - defines core types
   - Implement [file2] - main logic
   - Add [file3] - support functions
   - Write tests continuously

3. **Key Implementation Details**
   ```go
   // Example structure from wave plan
   type ServiceImpl struct {
       client *Client // From Effort 2
       logger *Logger // From Phase 1 common
   }
   
   // Must implement interface from Effort 1
   func (s *ServiceImpl) ProcessRequest(ctx context.Context, req *Request) (*Response, error) {
       // Implementation here
   }
   ```

4. **Integration Points**
   - Import contracts from `../api/`
   - Use shared libraries from `../lib/`
   - Follow patterns established in previous waves

## ✅ Test Requirements

### Coverage Requirements
- **Minimum Coverage**: 80%
- **Critical Paths**: 100% coverage required
- **Error Handling**: All error cases must be tested

### Test Categories
```yaml
required_tests:
  unit_tests:
    - All public functions
    - Error conditions
    - Edge cases
    
  integration_tests:
    - API endpoint testing
    - Database interactions
    - External service mocking
    
  performance_tests:
    - Load testing for high-traffic endpoints
    - Memory usage validation
```

## 📏 Size Constraints
**Target Size**: [NUMBER] lines (from wave plan)  
**Maximum Size**: 800 lines (HARD LIMIT)  
**Current Size**: [To be updated during implementation]  

### Size Monitoring Protocol
```bash
# Check size every ~200 lines
cd efforts/phase[PHASE]/wave[WAVE]/[effort]
$PROJECT_ROOT/tools/line-counter.sh

# If approaching 700 lines:
# 1. Alert Code Reviewer
# 2. Prepare for potential split
# 3. Focus on completing current functionality
```

## 🏁 Completion Criteria

### Implementation Checklist
- [ ] All files created as specified
- [ ] Size verified under 800 lines
- [ ] All imports properly referenced
- [ ] Interfaces from Effort 1 implemented
- [ ] Libraries from Effort 2 utilized

### Quality Checklist
- [ ] Test coverage ≥80%
- [ ] All tests passing
- [ ] No linting errors
- [ ] Error handling complete
- [ ] Logging added where appropriate

### Documentation Checklist
- [ ] Code comments for complex logic
- [ ] API documentation for exported functions
- [ ] README updated if needed
- [ ] Work log updated with progress

### Review Checklist
- [ ] Self-review completed
- [ ] Code committed and pushed
- [ ] Ready for Code Reviewer assessment
- [ ] No blocking issues

## 📊 Progress Tracking

### Work Log
```markdown
## [DATE] - Session 1
- Created file1.go (150 lines)
- Implemented core logic
- Added unit tests

## [DATE] - Session 2
- Created file2.go (200 lines)
- Integrated with Effort 1 contracts
- Current total: 350 lines

[Continue updating during implementation]
```

## ⚠️ Important Notes

### Parallelization Reminder
[IF Can Parallelize = Yes]
- This effort can run simultaneously with efforts [X, Y, Z]
- No need to wait for those efforts to complete
- Ensure no shared state with parallel efforts

[IF Can Parallelize = No]
- This is a BLOCKING effort
- Other efforts depend on this completion
- Priority: Complete ASAP to unblock team

### Common Pitfalls to Avoid
1. **Size Limit**: Monitor continuously, don't discover at end
2. **Dependencies**: Import from correct efforts, not reimpl> **Test Coverage**: Write tests as you go, not after
4. **Isolation**: Stay in your effort directory
5. **Parallelization**: Don't create dependencies on parallel efforts

## 📚 References

### Source Documents
- [Wave Implementation Plan](../../../phase-plans/PHASE-[PHASE]-WAVE-[WAVE]-IMPLEMENTATION-PLAN.md)
- [Phase Architecture](../../../phase-plans/PHASE-[PHASE]-ARCHITECTURE-PLAN.md)
- [Master Plan](../../../IMPLEMENTATION-PLAN.md)

### Code Examples
- [Previous Wave Implementation](../../wave[PREV]/)
- [Phase 1 Patterns](../../../phase1/)

### Standards
- [Coding Standards](../../../docs/coding-standards.md)
- [Testing Guide](../../../docs/testing-guide.md)

---

**Remember**: This plan was extracted from the wave plan. All headers, especially parallelization info, MUST match the wave plan exactly. The orchestrator depends on this information for spawning decisions.