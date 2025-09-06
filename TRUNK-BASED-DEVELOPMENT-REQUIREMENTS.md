# 🔴🔴🔴 TRUNK-BASED DEVELOPMENT REQUIREMENTS 🔴🔴🔴

## PARAMOUNT PRINCIPLE

**EVERY branch must be independently mergeable to main at ANY time, even YEARS later.**

This is not optional. This is not negotiable. This is the foundation of true Continuous Deployment.

## The Golden Standard

Your code must pass this test:
```
Can this branch merge to main RIGHT NOW, compile, pass tests, and not break ANYTHING?
Can it STILL merge successfully if the previous PR was 6 months ago?
Can it STILL work if the next PR doesn't come for another 6 months?
```

If the answer to ANY of these is "no", your branch is NOT ready.

## Core Requirements

### 1. Independent Mergeability (R307)
- **Every PR must compile when merged alone**
- **Every PR must pass all tests when merged alone**
- **No PR can break existing functionality**
- **Every incomplete feature must be behind a flag**
- **Every PR must be revertable without breaking others**

### 2. Atomic Design (R220)
- **One effort = One PR**
- **Each PR is self-contained**
- **No multi-part PRs that depend on each other**
- **Complete vertical slices of functionality**

### 3. Feature Flag Strategy
```yaml
principles:
  - Everything new starts behind a flag
  - Flags default to false until fully ready
  - Multiple efforts can contribute to one feature
  - Feature activates only when ALL parts merged
  - Old code paths remain until new is proven
```

## Implementation Patterns

### Pattern 1: Progressive Enhancement
```go
// Phase 1, Wave 1, Effort 1 (Merged January)
type UserService interface {
    GetUser(id string) (*User, error)
}

// Basic implementation - works alone
type BasicUserService struct{}
func (s *BasicUserService) GetUser(id string) (*User, error) {
    // Simple database lookup
    return db.GetUser(id)
}

// Phase 1, Wave 2, Effort 3 (Merged April - 3 months later!)
// Enhanced implementation - builds on basic
type CachedUserService struct {
    BasicUserService
    cache Cache
}
func (s *CachedUserService) GetUser(id string) (*User, error) {
    if Features["user-caching"] {
        if user, ok := s.cache.Get(id); ok {
            return user, nil
        }
    }
    // Falls back to basic if flag off or cache miss
    return s.BasicUserService.GetUser(id)
}
```

### Pattern 2: Feature Flags for Long-Running Features
```go
// Feature spans 6 months and 15 efforts
var Features = map[string]bool{
    "new-payment-system": false,  // Takes 6 efforts
    "payment-notifications": false, // Takes 3 efforts
    "payment-analytics": false,    // Takes 4 efforts
}

// Each effort adds its piece
func ProcessPayment(order Order) error {
    if Features["new-payment-system"] {
        // New system - only if ALL parts ready
        return newPaymentFlow(order)
    }
    // Old system - always works
    return legacyPaymentFlow(order)
}

// After ALL efforts merged (could be months!)
Features["new-payment-system"] = true // Finally activate
```

### Pattern 3: Handling Split Dependencies
```go
// Split 1 (Merged Monday)
type AuthInterface interface {
    Login(creds Credentials) (*Session, error)
    Logout(session *Session) error
}

// Stub for Split 1 - works alone
type StubAuth struct{}
func (s *StubAuth) Login(c Credentials) (*Session, error) {
    if c.Username == "test" && c.Password == "test" {
        return &Session{ID: "test-session"}, nil
    }
    return nil, ErrInvalidCredentials
}

// Split 2 (Merged 2 weeks later)
// Real implementation - replaces stub
type RealAuth struct {
    db Database
}
func (r *RealAuth) Login(c Credentials) (*Session, error) {
    // Full implementation
}

// Factory chooses implementation
func NewAuth() AuthInterface {
    if Features["real-auth"] && dbAvailable() {
        return &RealAuth{db: getDB()}
    }
    return &StubAuth{} // Always have fallback
}
```

## Anti-Patterns to AVOID

### ❌ Anti-Pattern 1: Sequential Dependencies
```go
// WRONG: E1.1.2 requires E1.1.1 to be merged first
import "github.com/project/e1-1-1/client" // FAILS if E1.1.1 not merged

// CORRECT: Check availability
var client ClientInterface
if clientAvailable() {
    client = realClient()
} else {
    client = mockClient()
}
```

### ❌ Anti-Pattern 2: Breaking Changes
```go
// WRONG: Changing existing API
func GetUser(id int) User { } // Was string, now int - BREAKS!

// CORRECT: Add new, deprecate old
func GetUser(id string) User { } // Keep working
func GetUserByID(id int) User { } // New method
```

### ❌ Anti-Pattern 3: Assuming Merge Order
```go
// WRONG: Assumes efforts merge in sequence
if effort1Done() && effort2Done() {
    // Assumes 1 before 2
}

// CORRECT: Handle any order
if effort1Done() {
    enableFeature1()
}
if effort2Done() {
    enableFeature2()
}
// Features work independently
```

## Testing Requirements

### 1. Test Independence
```go
func TestNewFeature(t *testing.T) {
    // Test with feature ON
    Features["new-feature"] = true
    result := RunFeature()
    assert.Equal(t, expected, result)
    
    // Test with feature OFF (must still pass!)
    Features["new-feature"] = false
    result = RunFeature()
    assert.NotPanics(t, func() { RunFeature() })
}
```

### 2. Test Degradation
```go
func TestGracefulDegradation(t *testing.T) {
    // Simulate missing dependency
    SetDependencyAvailable(false)
    
    // Should fall back gracefully
    result, err := ProcessRequest()
    assert.NoError(t, err)
    assert.Equal(t, "degraded-response", result)
}
```

### 3. Test Flag Combinations
```go
func TestFeatureFlagCombinations(t *testing.T) {
    // Test all possible flag states
    flags := []string{"feature-a", "feature-b", "feature-c"}
    
    for i := 0; i < (1 << len(flags)); i++ {
        // Set flags based on bit pattern
        for j, flag := range flags {
            Features[flag] = (i & (1 << j)) != 0
        }
        
        // App must work in ANY combination
        assert.NotPanics(t, func() {
            app.Run()
        })
    }
}
```

## CI/CD Pipeline Requirements

### Pre-Merge Checks (MANDATORY)
```yaml
pre_merge_checks:
  - name: "Build in Isolation"
    steps:
      - checkout main
      - merge PR branch
      - run build
      - must: PASS
      
  - name: "Test Existing Features"
    steps:
      - run existing test suite
      - must: 100% PASS
      
  - name: "Test with Flags Off"
    steps:
      - disable all new flags
      - run full test suite
      - must: PASS
      
  - name: "Test Revertability"
    steps:
      - merge PR
      - revert PR
      - run tests
      - must: PASS
```

### Post-Merge Monitoring
```yaml
monitoring:
  - metric: "Build Status"
    requirement: "Always Green"
    
  - metric: "Test Pass Rate"
    requirement: "100%"
    
  - metric: "Feature Flag Status"
    requirement: "Document which flags active"
    
  - metric: "Rollback Ready"
    requirement: "Can revert in <5 minutes"
```

## Workflow Example: 6-Month Feature Development

### Month 1: Foundation
```
E1.1.1: Define interfaces (merged Week 1)
E1.1.2: Basic implementation (merged Week 3)
E1.1.3: Unit tests (merged Week 4)

Status: Basic feature works behind flag
Flag: feature-v1 = false (not ready)
```

### Month 2-3: Enhancement
```
E1.2.1: Add caching (merged Month 2, Week 2)
E1.2.2: Add metrics (merged Month 2, Week 4)
E1.2.3: Performance optimization (merged Month 3, Week 2)

Status: Enhanced but still behind flag
Flag: feature-v1 = false (still testing)
```

### Month 4-5: Integration
```
E1.3.1: API integration (merged Month 4, Week 1)
E1.3.2: UI components (merged Month 4, Week 3)
E1.3.3: Documentation (merged Month 5, Week 1)

Status: Fully integrated, testing in staging
Flag: feature-v1 = false (staging only)
```

### Month 6: Activation
```
E1.4.1: Final testing (merged Month 6, Week 1)
E1.4.2: Monitoring setup (merged Month 6, Week 2)
CONFIG: Set feature-v1 = true (Week 3)

Status: Feature fully active
Flag: feature-v1 = true (LIVE!)
```

**Key Point**: At ANY point during these 6 months, main branch is stable, deployable, and working!

## Grading Criteria

### Automatic Failures (-100%)
- PR breaks the build when merged alone
- PR breaks existing functionality
- PR requires specific merge order
- PR cannot be reverted safely

### Major Deductions
- Missing feature flags for incomplete work: -40%
- No graceful degradation: -30%
- Insufficient test coverage: -20%
- Poor flag documentation: -10%

### Excellence Indicators (+bonus)
- PRs mergeable years apart: +10%
- Comprehensive flag testing: +10%
- Perfect rollback capability: +10%
- Zero production incidents: +20%

## The Trunk-Based Development Oath

As a Software Factory agent, I pledge to:

1. **Never break the build** - Every PR leaves main deployable
2. **Always use feature flags** - Incomplete work stays hidden
3. **Design for independence** - My PR needs no other
4. **Test all states** - Flags on, flags off, all must work
5. **Enable rollback** - My PR can be reverted anytime
6. **Think long-term** - My code works regardless of merge timing
7. **Protect production** - Users never see broken features

## Remember

**The goal is not just to merge code. The goal is to maintain a perpetually deployable main branch that can receive contributions at any time, in any order, without coordination, without breakage, without fear.**

This is the path to true Continuous Deployment. This is the Software Factory way.

---

*Rule R307 enforces these requirements. Violation means project failure.*