# 🔴🔴🔴 RULE R307: Independent Branch Mergeability (PARAMOUNT)

## Classification
- **Category**: CD/Integration/Core
- **Criticality**: PARAMOUNT (HIGHEST - Supersedes all other priorities)
- **Penalty**: -100% for breaking build or existing functionality
- **Introduced**: v2.0.0
- **Modified**: v2.0.0

## The PARAMOUNT Requirement

**EVERY branch (efforts AND splits) MUST:**
1. **Compile without errors when merged alone** - No dependency on other unmerged work
2. **Not break ANY existing functionality** - Zero regression tolerance
3. **Be mergeable with MONTHS delay after previous PR** - True independence
4. **Use feature flags for ALL incomplete features** - Progressive enhancement
5. **Maintain working build ALWAYS** - Never leave main broken

## Real-World Scenario That MUST Work

```
January 1: E1.1.1 merged → Build ✅ App works ✅
February 15: E1.1.2 merged → Build ✅ App still works ✅
April 1: E1.1.3-split-001 merged → Build ✅ App still works ✅
June 15: E1.1.3-split-002 merged → Build ✅ App still works ✅
August 30: E1.2.1 merged → Build ✅ App still works ✅
December 25: E1.2.2 merged → Build ✅ App still works ✅
```

**Throughout the ENTIRE YEAR, the app NEVER breaks!**

## Design Patterns for Independence

### 1. Feature Flag Pattern (MANDATORY)
```go
// Define flags at start of feature development
var Features = map[string]bool{
    "new-auth-system": false,      // Will span E1.1.1 through E1.2.3
    "oauth-integration": false,     // Requires all E1.1.3 splits
    "advanced-logging": false,      // E1.2.1 onwards
}

// Check flags before using new code
func HandleAuth(req *Request) (*Response, error) {
    if Features["new-auth-system"] {
        return newAuthHandler(req)  // New implementation
    }
    return legacyAuthHandler(req)   // Existing working code
}

// Enable only when ALL parts are merged
// (After E1.2.3 completes)
Features["new-auth-system"] = true
```

### 2. Interface Before Implementation Pattern
```go
// E1.1.1: Define interface (merged January)
type AuthProvider interface {
    Authenticate(credentials Credentials) (*User, error)
    Refresh(token string) (*Token, error)
}

// E1.1.2: Add stub implementation (merged March)
type StubAuthProvider struct{}
func (s *StubAuthProvider) Authenticate(c Credentials) (*User, error) {
    return nil, ErrNotImplemented
}

// E1.1.3-split-001: Partial implementation (merged May)
func (s *StubAuthProvider) Authenticate(c Credentials) (*User, error) {
    // Basic auth works, OAuth pending
    if c.Type == "basic" {
        return basicAuth(c)
    }
    return nil, ErrNotImplemented
}

// E1.1.3-split-002: Complete implementation (merged July)
// Now OAuth works too
```

### 3. Graceful Degradation Pattern
```go
// New feature tries advanced path, falls back if not ready
func ProcessPayment(order *Order) error {
    if Features["advanced-payment"] && advancedPaymentReady() {
        return advancedPaymentFlow(order)
    }
    
    // Always have working fallback
    return standardPaymentFlow(order)
}
```

### 4. Progressive Enhancement Pattern
```go
// Start with basic functionality
type UserProfile struct {
    ID       string
    Name     string
    Email    string
    // Advanced fields added progressively
    OAuth    *OAuthData    `json:",omitempty"` // E1.1.3
    Settings *UserSettings `json:",omitempty"` // E1.2.1
    Analytics *UserMetrics `json:",omitempty"` // E1.2.2
}
```

## Testing Strategy for Partial Implementations

### 1. Feature Flag Testing
```go
func TestNewAuthSystem(t *testing.T) {
    // Skip if feature not complete
    if !Features["new-auth-system"] {
        t.Skip("New auth system not yet enabled")
    }
    // Actual tests
}

func TestLegacyAuthStillWorks(t *testing.T) {
    // MUST ALWAYS PASS
    Features["new-auth-system"] = false
    defer func() { Features["new-auth-system"] = true }()
    
    // Test legacy path
}
```

### 2. Partial Implementation Testing
```go
func TestOAuthIntegration(t *testing.T) {
    if !allSplitsMerged() {
        // Test degraded functionality
        _, err := oauth.Authenticate()
        assert.Equal(t, ErrNotImplemented, err)
        return
    }
    // Full OAuth tests
}
```

## Validation Requirements

### Before EVERY Merge (Automated CI/CD)
```bash
# 1. Build in isolation
git checkout main
git pull origin main
git checkout -b test-branch
git merge --no-ff feature-branch

# 2. Compile check
make build
# MUST SUCCEED

# 3. Existing tests
make test
# ALL EXISTING TESTS MUST PASS

# 4. Deploy test
make deploy-test
# APP MUST START AND WORK

# 5. Feature verification
curl http://localhost:8080/health
# EXISTING FEATURES MUST WORK
```

### Manual Review Checklist
- [ ] Branch builds successfully when merged alone?
- [ ] All existing tests pass?
- [ ] New features are flag-protected?
- [ ] Could this merge 1 year from now and still work?
- [ ] No breaking changes to existing APIs?
- [ ] Graceful degradation for missing dependencies?
- [ ] Documentation explains partial state?

## Common Violations and Fixes

### ❌ Violation: Depending on Unmerged Work
```go
// BAD: Assumes E1.1.1 is merged
import "github.com/project/newauth"

func Handler() {
    user := newauth.GetUser() // FAILS if E1.1.1 not merged
}
```

### ✅ Fix: Check and Fallback
```go
// GOOD: Works regardless of merge state
func Handler() {
    var user *User
    if authAvailable() {
        user = newauth.GetUser()
    } else {
        user = getCurrentUser() // Existing method
    }
}
```

### ❌ Violation: Breaking Existing Features
```go
// BAD: Replaces working code
func GetUserProfile(id string) *Profile {
    return newProfileSystem.Get(id) // Breaks if new system incomplete
}
```

### ✅ Fix: Progressive Enhancement
```go
// GOOD: Enhances without breaking
func GetUserProfile(id string) *Profile {
    profile := legacyProfileSystem.Get(id) // Always works
    
    if Features["enhanced-profiles"] {
        enhance(profile) // Add new fields if available
    }
    
    return profile
}
```

## Enforcement in Agent Workflows

### SW Engineer Requirements
- Design every effort assuming months between merges
- Implement feature flags from day one
- Test both enabled and disabled states
- Never remove working code until replacement is proven

### Code Reviewer Verification
- Verify branch builds in isolation
- Check feature flag implementation
- Ensure graceful degradation
- Validate no breaking changes

### Architect Assessment
- Review independence across phase
- Verify feature flag strategy
- Ensure progressive enhancement approach
- Validate year-long merge compatibility

### Orchestrator Monitoring
- Track feature flag states
- Monitor build health after every merge
- Ensure sequential handling of dependent work
- Validate independence before parallel execution

## Integration with Other Rules

### Relationship to R220/R221 (Size Limits)
- Smaller PRs are easier to make independent
- Feature flags help split large features

### Relationship to R277 (Continuous Build)
- Every merge must maintain green build
- No "temporary" breakage allowed

### Relationship to Split Planning
- Each split must be independently mergeable
- No split can depend on later splits

## The Golden Rule

**If you can't merge this branch 6 months from now with zero other changes and have everything still work, it's NOT independent enough!**

## Penalties for Violations

1. **Breaking the build**: -100% (IMMEDIATE FAILURE)
2. **Breaking existing features**: -100% (IMMEDIATE FAILURE)
3. **Requiring specific merge order**: -50%
4. **Missing feature flags**: -40%
5. **No graceful degradation**: -30%
6. **Incomplete testing of flags**: -20%

## Success Metrics

- 100% of branches compile when merged alone
- 0 regressions in existing functionality
- 100% of new features behind flags
- Months can pass between merges with no issues
- Any PR can be reverted without breaking others

---

**REMEMBER**: This is PARAMOUNT. No other requirement supersedes the need for every branch to be independently mergeable. This is the foundation of true Continuous Deployment and trunk-based development.