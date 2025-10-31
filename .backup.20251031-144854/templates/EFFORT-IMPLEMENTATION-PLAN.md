# Effort Implementation Plan: [EFFORT_NAME]

<!-- STORAGE LOCATION: This plan should be saved in:
     .software-factory/phase[X]/wave[Y]/[effort-name]/IMPLEMENTATION-PLAN-YYYYMMDD-HHMMSS.md
     within the effort's working directory. This keeps plans organized and separate from code. -->

<!-- NOTE: [PROJECT_PREFIX/] should be replaced with the actual project prefix from target-repo-config.yaml 
     If project_prefix is "tmc-workspace", branches become: tmc-workspace/phase1/wave1/effort-name
     If project_prefix is empty, branches become: phase1/wave1/effort-name -->

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort**: [EFFORT_NUMBER] - [EFFORT_NAME]  
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-[name]`  
**Base Branch**: `[MANDATORY per R337 - e.g., phase1-wave1-integration]`  
**Base Branch Reason**: [WHY this base per R308 incremental strategy]  
**Can Parallelize**: [Yes/No] (COPIED FROM WAVE PLAN)  
**Parallel With**: [List effort numbers or "None"] (COPIED FROM WAVE PLAN)  
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: [List dependent efforts] (COPIED FROM WAVE PLAN)  
**Dependent Efforts**: [List efforts that depend on THIS one]  
**Atomic PR**: ✅ This effort = ONE PR to main (R220 REQUIREMENT)  

## 📋 Source Information
**Wave Plan**: PHASE-[PHASE]-WAVE-[WAVE]-IMPLEMENTATION-PLAN.md  
**Effort Section**: Effort [NUMBER]  
**Created By**: Code Reviewer Agent  
**Date**: [DATE]  
**Extracted**: [TIMESTAMP]

## 🔴 BASE BRANCH VALIDATION (R337 MANDATORY)
**The orchestrator-state-v3.json is the SOLE SOURCE OF TRUTH for base branches!**
- Base branch MUST be explicitly specified above
- Base branch MUST match what's in orchestrator-state-v3.json
- Reason MUST explain why this base (per R308 incremental)
- Orchestrator MUST record this in state file before creating infrastructure  

## 🚀 Parallelization Context
**Can Parallelize**: [Yes/No]  
**Parallel With**: [Efforts that can run simultaneously]  
**Blocking Status**: [If No - this effort blocks: X, Y, Z]  
**Parallel Group**: [If Yes - member of parallel group: A, B, C]  
**Orchestrator Guidance**: [When orchestrator should spawn this effort]  

## 🚨 EXPLICIT SCOPE DEFINITION (R311 MANDATORY)

### IMPLEMENT EXACTLY (BE SPECIFIC!)

#### Functions to Create (EXACTLY [N] - NO MORE)
```go
1. FunctionName1(params) ReturnType    // ~[X] lines - [specific purpose]
2. FunctionName2(params) ReturnType    // ~[Y] lines - [specific purpose]  
3. FunctionName3(params) ReturnType    // ~[Z] lines - [specific purpose]
// STOP HERE - DO NOT ADD MORE FUNCTIONS
```

#### Types/Structs to Define (EXACTLY [N])
```go
// Type 1: Core data model
type ModelName struct {
    Field1 string  // [purpose]
    Field2 int     // [purpose]
    // EXACTLY these fields, NO methods in this effort
}

// Type 2: Service interface
type ServiceInterface interface {
    Method1(ctx context.Context) error  // Required
    Method2(id string) (*Model, error)  // Required
    // NO additional methods
}
```

#### Endpoints/Handlers (if applicable)
```go
// EXACTLY these endpoints:
POST   /api/v1/resource   // CreateResource handler - ~100 lines
GET    /api/v1/resource   // GetResource handler - ~80 lines
// NO additional endpoints (no DELETE, PUT, etc.)
```

### 🛑 DO NOT IMPLEMENT (SCOPE BOUNDARIES)

**EXPLICITLY FORBIDDEN IN THIS EFFORT:**
- ❌ DO NOT add validation beyond basic nil/empty checks
- ❌ DO NOT implement Update/Delete operations (future effort)
- ❌ DO NOT add caching or optimization (future effort)
- ❌ DO NOT implement authentication/authorization (separate effort)
- ❌ DO NOT add comprehensive error handling (basic only)
- ❌ DO NOT create helper/utility functions not listed above
- ❌ DO NOT refactor or "improve" existing code
- ❌ DO NOT add logging beyond critical errors
- ❌ DO NOT write performance tests or benchmarks
- ❌ DO NOT implement nice-to-have features

### 📊 REALISTIC SIZE CALCULATION

```
Component Breakdown:
- Function 1:                    [X] lines
- Function 2:                    [Y] lines
- Function 3:                    [Z] lines
- Type definitions:             [A] lines
- Basic tests (N × 30):        [B] lines
- Integration setup:           [C] lines

TOTAL ESTIMATE: [SUM] lines (must be <800)
BUFFER: [800-SUM] lines for unforeseen needs
```

## 🔴🔴🔴 PRE-PLANNING RESEARCH RESULTS (R374 MANDATORY) 🔴🔴🔴

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| [MUST LIST ALL FOUND] | [branch/path] | [exact signature] | YES/NO |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| [MUST LIST ALL REUSABLE CODE] | [branch/path] | [what it does] | [import/use directly] |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| [MUST LIST ALL EXISTING APIs] | [method name] | [exact signature] | [notes] |

### FORBIDDEN DUPLICATIONS (R373)
- ❌ DO NOT create new [interface] (exists in [location])
- ❌ DO NOT reimplement [functionality] (exists in [location])
- ❌ DO NOT create alternative [method] signatures

### REQUIRED INTEGRATE_WAVE_EFFORTSS (R373)
- ✅ MUST implement [interface] from [location] with EXACT signature
- ✅ MUST reuse [component] from [location]
- ✅ MUST import and use [package] for [functionality]

## 📁 Files to Create

### Primary Implementation Files
```yaml
new_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/[feature]/[file1].go
    lines: ~[NUMBER] MAX
    purpose: [SPECIFIC PURPOSE]
    contains:
      - FunctionName1 (NO additional functions)
      - TypeName1 (NO additional types)
      
  - path: pkg/phase[PHASE]/wave[WAVE]/[feature]/[file2].go
    lines: ~[NUMBER] MAX  
    purpose: [SPECIFIC PURPOSE]
    contains:
      - FunctionName2 ONLY
      - NO helper functions
```

### Test Files
```yaml
test_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/[feature]/[file]_test.go
    lines: ~[NUMBER] MAX
    coverage_target: 80%
    test_functions:
      - TestFunction1Basic  # ~30 lines
      - TestFunction2Basic  # ~30 lines
      # NO edge cases, NO benchmarks
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

## 🔴 ATOMIC PR REQUIREMENTS (R220 - SUPREME LAW)

### 🔴🔴🔴 PARAMOUNT: Independent Mergeability (R307) 🔴🔴🔴
**This effort MUST be mergeable at ANY time, even YEARS later:**
- ✅ Must compile when merged alone to main
- ✅ Must NOT break any existing functionality
- ✅ Must use feature flags for incomplete features
- ✅ Must work even if next effort merges 6 months later
- ✅ Must gracefully degrade if dependencies missing

### Feature Flags for This Effort
```yaml
feature_flags:
  - flag: "EFFORT_[NAME]_ENABLED"
    location: "config/features.yaml"
    default: false
    purpose: "Control activation of [feature]"
    activation: "Set true when [condition]"
```

### 🚨🚨🚨 R355 PRODUCTION READY CODE (SUPREME LAW #5) 🚨🚨🚨

**ALL CODE MUST BE PRODUCTION READY - NO EXCEPTIONS**

#### ❌ ABSOLUTELY FORBIDDEN:
- NO stubs or placeholder implementations
- NO mocks except in test directories
- NO hardcoded credentials or secrets
- NO static configuration values
- NO TODO/FIXME markers in code
- NO returning nil for "later implementation"
- NO panic("not implemented") patterns

#### ✅ REQUIRED PATTERNS:
```go
// ❌ WRONG - Hardcoded (AUTOMATIC FAILURE)
dbURL := "postgres://localhost/mydb"
apiKey := "sk-12345"
timeout := 30

// ✅ CORRECT - Configuration-driven
dbURL := os.Getenv("DATABASE_URL")
apiKey := os.Getenv("API_KEY")
timeout := config.GetInt("timeout", 30)
```

### Interface Implementations (Instead of Stubs)
```yaml
interfaces:
  - name: "[ServiceName]"
    implements: "I[ServiceName]"
    type: "Minimal Working Implementation"
    notes: "Use in-memory storage if DB not ready"
    production_ready: true
```

### PR Mergeability Checklist
- [ ] PR can merge to main independently
- [ ] Build passes with just this PR
- [ ] All tests pass in isolation
- [ ] Feature flags hide incomplete features
- [ ] Stubs replace missing dependencies
- [ ] No breaking changes to existing code
- [ ] Backward compatible with main

## 🔴 MANDATORY ADHERENCE CHECKPOINTS (R311)

### Before Starting:
```bash
echo "EFFORT SCOPE LOCKED:"
echo "✓ Functions: EXACTLY [N] (list them)"
echo "✓ Types: EXACTLY [N] (list them)"
echo "✓ Endpoints: EXACTLY [N] (list them)"
echo "✓ Tests: EXACTLY [N] basic tests"
echo "✗ Validation: MINIMAL ONLY"
echo "✗ Extra features: NONE"
echo "✗ Optimizations: NONE"
```

### During Implementation:
```bash
# Check scope adherence after each component
FUNC_COUNT=$(grep -c "^func [A-Z]" *.go 2>/dev/null || echo 0)
if [ "$FUNC_COUNT" -gt [EXPECTED_COUNT] ]; then
    echo "⚠️ WARNING: Exceeding function count! Stop adding!"
fi
```

### If Addition Seems Necessary:
```markdown
## Justification for Addition
- **What needs adding**: [describe]
- **Why it's ESSENTIAL**: [not just nice to have]
- **What BREAKS without it**: [specific failure]
- **Line count of addition**: [estimate]
- **Decision**: [Add with justification / Skip for now]
```

## 📝 Implementation Instructions

### Step-by-Step Guide
1. **Scope Acknowledgment**
   - Read and acknowledge DO NOT IMPLEMENT section
   - Extract exact function/type counts
   - Create .scope-acknowledgment file

2. **Implementation Order**
   - Start with [file1] - defines core types ONLY
   - Implement [file2] - EXACTLY the listed functions
   - Add [file3] - ONLY specified components
   - Write MINIMAL tests (basic cases only)

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

### Common Pitfalls to Avoid (R311 ENFORCEMENT)
1. **SCOPE CREEP**: Adding "helpful" features = AUTOMATIC FAILURE
2. **OVER-ENGINEERING**: Making it "production-ready" = 3-5X overrun
3. **ASSUMPTIONS**: Implementing what "seems needed" = VIOLATION
4. **Size Limit**: Monitor continuously with line-counter.sh
5. **Dependencies**: Import from correct efforts, not reimplement
6. **Test Coverage**: Basic tests only - no edge cases unless specified
7. **Isolation**: Stay in your effort directory
8. **Parallelization**: Don't create dependencies on parallel efforts

### Success Criteria Checklist
- [ ] Read and acknowledged DO NOT IMPLEMENT section
- [ ] Implemented EXACTLY the specified functions (no more)
- [ ] Created EXACTLY the specified types (no more)
- [ ] Wrote EXACTLY the specified tests (no more)
- [ ] Total lines under 800
- [ ] NO unauthorized features added
- [ ] Any additions have written justification
- [ ] Followed all scope boundaries

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