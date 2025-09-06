# COMPREHENSIVE FEATURE FLAG STRATEGY

## 🔴 PARAMOUNT PRINCIPLE

**Feature flags are not optional - they are MANDATORY for independent mergeability (R307)**

Every incomplete feature, every partial implementation, every work-in-progress MUST be behind a feature flag.

## Core Philosophy

Feature flags enable:
1. **Independent mergeability** - PRs merge without coordination
2. **Progressive delivery** - Features roll out gradually
3. **Safe experimentation** - Try without commitment
4. **Instant rollback** - Disable without deployment
5. **Parallel development** - Multiple teams, no conflicts

## Feature Flag Lifecycle

### Phase 1: Definition (Before First PR)
```yaml
# In effort planning
feature_flags:
  - name: "payment-redesign"
    description: "New payment processing system"
    default: false
    estimated_efforts: 12
    estimated_duration: "6 months"
    dependencies: []
```

### Phase 2: Implementation (During Development)
```go
// First PR - Define flag
var Features = map[string]bool{
    "payment-redesign": false,
}

// Subsequent PRs - Check flag
func ProcessPayment(order Order) error {
    if Features["payment-redesign"] {
        return newPaymentSystem(order) // New path
    }
    return legacyPaymentSystem(order) // Old path (keep working!)
}
```

### Phase 3: Testing (After Implementation)
```go
// Test both paths
func TestPaymentProcessing(t *testing.T) {
    // Test with flag ON
    Features["payment-redesign"] = true
    assert.NoError(t, ProcessPayment(testOrder))
    
    // Test with flag OFF
    Features["payment-redesign"] = false
    assert.NoError(t, ProcessPayment(testOrder))
}
```

### Phase 4: Rollout (Gradual Activation)
```yaml
rollout_stages:
  - stage: "Internal Testing"
    percentage: 0
    users: ["qa-team"]
    duration: "1 week"
    
  - stage: "Beta Users"
    percentage: 1
    users: ["beta-group"]
    duration: "2 weeks"
    
  - stage: "Gradual Rollout"
    percentage: [5, 10, 25, 50, 75, 100]
    duration: "4 weeks"
```

### Phase 5: Cleanup (After Full Rollout)
```go
// After 100% rollout and stability
// 1. Remove flag checks
func ProcessPayment(order Order) error {
    return newPaymentSystem(order) // Flag removed
}

// 2. Remove old code
// Delete legacyPaymentSystem function

// 3. Remove flag definition
// Delete from Features map
```

## Flag Implementation Patterns

### Pattern 1: Simple Boolean Flag
```go
// Basic on/off switch
if Features["new-feature"] {
    // New behavior
} else {
    // Old behavior
}
```

### Pattern 2: Percentage Rollout
```go
// Gradual rollout by percentage
func ShouldUseNewFeature(userID string) bool {
    if !Features["new-feature"] {
        return false
    }
    
    percentage := getFeaturePercentage("new-feature")
    hash := hashUserID(userID)
    return (hash % 100) < percentage
}
```

### Pattern 3: User Cohort Targeting
```go
// Target specific user groups
func IsFeatureEnabledForUser(feature string, user User) bool {
    if !Features[feature] {
        return false // Global kill switch
    }
    
    // Check user cohorts
    if user.IsBetaTester && Features[feature+"-beta"] {
        return true
    }
    
    if user.IsInternal && Features[feature+"-internal"] {
        return true
    }
    
    return Features[feature+"-general"]
}
```

### Pattern 4: Complex Feature with Sub-flags
```go
// Major feature with multiple components
var AuthFeatures = map[string]bool{
    "new-auth":          false, // Master switch
    "new-auth-ui":       false, // New login UI
    "new-auth-api":      false, // New API endpoints
    "new-auth-oauth":    false, // OAuth providers
    "new-auth-2fa":      false, // Two-factor auth
    "new-auth-biometric": false, // Biometric support
}

func GetLoginPage() Page {
    if !AuthFeatures["new-auth"] {
        return LegacyLoginPage()
    }
    
    page := BasicLoginPage()
    
    if AuthFeatures["new-auth-ui"] {
        page = ModernLoginPage()
    }
    
    if AuthFeatures["new-auth-oauth"] {
        page.AddOAuthProviders()
    }
    
    if AuthFeatures["new-auth-2fa"] {
        page.Add2FAOption()
    }
    
    return page
}
```

### Pattern 5: Feature Flag with Configuration
```go
type FeatureConfig struct {
    Enabled    bool
    Percentage int
    Whitelist  []string
    Blacklist  []string
    Config     map[string]interface{}
}

var FeatureConfigs = map[string]FeatureConfig{
    "advanced-search": {
        Enabled:    true,
        Percentage: 50,
        Whitelist:  []string{"power-users"},
        Config: map[string]interface{}{
            "max-results":     100,
            "enable-fuzzy":    true,
            "timeout-seconds": 30,
        },
    },
}
```

## Flag Management Infrastructure

### 1. Flag Definition File
```yaml
# features.yaml - Source of truth for all flags
features:
  - name: payment-redesign
    description: Complete payment system overhaul
    owner: payments-team
    created: 2024-01-15
    type: release
    default: false
    dependencies: []
    cleanup_date: 2024-12-31
    
  - name: experimental-ai-suggestions
    description: AI-powered product suggestions
    owner: ml-team
    created: 2024-03-01
    type: experiment
    default: false
    percentage: 0
    success_metrics:
      - conversion_rate > 1.05x
      - user_satisfaction > 4.0
```

### 2. Flag Service Implementation
```go
package features

type FlagService struct {
    store      FlagStore
    cache      Cache
    metrics    MetricsCollector
    evaluator  RuleEvaluator
}

func (s *FlagService) IsEnabled(flag string, context Context) bool {
    // Check cache first
    if cached, ok := s.cache.Get(flag, context); ok {
        return cached
    }
    
    // Get flag configuration
    config, err := s.store.GetFlag(flag)
    if err != nil {
        s.metrics.RecordError(flag, err)
        return false // Default to off on error
    }
    
    // Evaluate rules
    result := s.evaluator.Evaluate(config, context)
    
    // Cache result
    s.cache.Set(flag, context, result)
    
    // Record metrics
    s.metrics.RecordEvaluation(flag, result)
    
    return result
}
```

### 3. Flag Monitoring and Metrics
```go
package monitoring

func MonitorFeatureFlag(flag string) {
    metrics := CollectMetrics(flag)
    
    // Alert on anomalies
    if metrics.ErrorRate > 0.05 {
        Alert("High error rate for feature: " + flag)
        AutoDisableFlag(flag)
    }
    
    if metrics.PerformanceDegradation > 0.20 {
        Alert("Performance degradation for feature: " + flag)
        ReducePercentage(flag, 50)
    }
    
    // Track adoption
    LogMetric("feature.adoption", flag, metrics.UsageCount)
    LogMetric("feature.success", flag, metrics.SuccessRate)
}
```

## Testing Strategy for Feature Flags

### 1. Unit Tests - Both States
```go
func TestFeatureWithFlagStates(t *testing.T) {
    tests := []struct {
        name     string
        flagOn   bool
        expected string
    }{
        {"Flag Off", false, "legacy-behavior"},
        {"Flag On", true, "new-behavior"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            Features["test-feature"] = tt.flagOn
            result := FunctionUnderTest()
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 2. Integration Tests - Transition Scenarios
```go
func TestFeatureFlagTransition(t *testing.T) {
    // Start with flag off
    Features["migration-feature"] = false
    setupLegacyData()
    
    // Enable flag mid-test
    Features["migration-feature"] = true
    
    // Should handle transition gracefully
    assert.NoError(t, PerformOperation())
    assert.True(t, DataConsistent())
}
```

### 3. Chaos Testing - Random Flag States
```go
func TestChaosFlags(t *testing.T) {
    flags := []string{"feature-a", "feature-b", "feature-c"}
    
    for i := 0; i < 100; i++ {
        // Randomize all flags
        for _, flag := range flags {
            Features[flag] = rand.Float32() < 0.5
        }
        
        // System should never crash
        assert.NotPanics(t, func() {
            RunApplication()
        })
    }
}
```

## Flag Anti-Patterns to Avoid

### ❌ Anti-Pattern 1: Nested Flag Dependencies
```go
// BAD: Complex nested dependencies
if Features["feature-a"] {
    if Features["feature-b"] {
        if Features["feature-c"] {
            // Too complex!
        }
    }
}

// GOOD: Single flag or clear combination
if Features["complete-feature"] {
    // Clear and simple
}
```

### ❌ Anti-Pattern 2: Permanent Flags
```go
// BAD: Flag that never gets removed
if Features["temp-fix-2019"] { // Still here in 2024!
    // Technical debt
}

// GOOD: Flags have cleanup dates
if Features["migration-2024-q1"] {
    // Will be removed after Q1 2024
}
```

### ❌ Anti-Pattern 3: Flag Explosion
```go
// BAD: Too many granular flags
Features["button-color-blue"]
Features["button-size-large"]
Features["button-text-bold"]
Features["button-animation-fade"]

// GOOD: Grouped configuration
Features["new-button-design"] // One flag for related changes
```

## Flag Documentation Requirements

### For Each Flag, Document:
```markdown
## Feature: [FLAG_NAME]

**Owner**: [Team/Person]
**Created**: [Date]
**Purpose**: [Why this flag exists]
**Expected Duration**: [How long until removal]

### Activation Criteria
- [ ] All implementation PRs merged
- [ ] Tests passing with flag on
- [ ] Performance benchmarks met
- [ ] Security review complete

### Rollout Plan
1. Internal testing (Week 1)
2. 5% of users (Week 2)
3. 25% of users (Week 3)
4. 50% of users (Week 4)
5. 100% of users (Week 5)

### Rollback Plan
If issues detected:
1. Set flag to false
2. Monitor error rates
3. Fix issues
4. Resume rollout

### Cleanup Plan
After successful rollout:
1. Remove flag checks (PR #XXX)
2. Remove old code (PR #XXX)
3. Remove flag definition (PR #XXX)
```

## Enforcement and Compliance

### Code Review Checklist
- [ ] New features have flags?
- [ ] Flags default to false?
- [ ] Both paths tested?
- [ ] Graceful degradation?
- [ ] Flag documented?
- [ ] Cleanup plan defined?

### Automated Checks
```yaml
# CI/CD pipeline checks
feature_flag_validation:
  - check: "All new features flagged"
    command: "detect-unflagged-features"
    
  - check: "Flags have tests"
    command: "verify-flag-test-coverage"
    
  - check: "No orphaned flags"
    command: "find-unused-flags"
    
  - check: "Flags have owners"
    command: "verify-flag-ownership"
```

## Success Metrics

Track these metrics for feature flag success:
1. **Deployment Frequency**: Increase with confidence
2. **Rollback Rate**: Decrease with better testing
3. **Incident Rate**: Lower with gradual rollouts
4. **Feature Adoption**: Measure actual usage
5. **Development Velocity**: Parallel work enabled

## The Feature Flag Commitment

As Software Factory agents, we commit to:
1. **Always flag incomplete work**
2. **Never break existing features**
3. **Test all flag states**
4. **Document flag lifecycle**
5. **Clean up after rollout**
6. **Monitor flag health**
7. **Enable safe experimentation**

## Remember

Feature flags are not just about hiding incomplete work. They're about:
- **Decoupling deployment from release**
- **Enabling continuous deployment**
- **Reducing risk**
- **Increasing velocity**
- **Improving quality**

This is the path to true trunk-based development where any PR can merge at any time.

---

*This strategy is enforced by R307 (Independent Mergeability) and R220 (Atomic PR Design). Violation results in project failure.*