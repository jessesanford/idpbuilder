# IMPLEMENTATION PLAN: Effort 1.2.2 - Fallback Strategies

## EFFORT INFRASTRUCTURE METADATA
**EFFORT_NAME**: fallback-strategies
**PHASE**: 1
**WAVE**: 2
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave2/fallback-strategies
**BRANCH**: idpbuilder-oci-mvp/phase1/wave2/fallback-strategies
**ISOLATION_BOUNDARY**: /home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase1/wave2/fallback-strategies

## Executive Summary

**Effort**: Fallback Strategies (1.2.2)
**Estimated Size**: ~400 lines (HARD LIMIT: 800 lines)
**Timeline**: 2 days
**Dependencies**: Wave 1 TrustManager interface

## Purpose

Implement intelligent fallback mechanisms for certificate issues, including:
- `--insecure` flag handler with proper warnings
- Auto-recovery suggestions and mechanisms
- Clear error messages with user guidance
- Security decision logging and audit trail

## Architecture Components

### Primary Implementation Files

1. **pkg/certs/fallback.go** (~200 lines)
   - FallbackHandler interface and implementation
   - Error analysis and strategy generation
   - Security decision logging

2. **pkg/certs/recovery.go** (~100 lines)
   - Auto-recovery mechanisms
   - Retry logic and timeout handling
   - Recovery result reporting

3. **pkg/certs/insecure.go** (~50 lines)
   - Insecure mode configuration
   - Warning system implementation
   - Explicit user consent handling

4. **Test files** (~50 lines total)
   - fallback_test.go
   - recovery_test.go
   - insecure_test.go

### Key Types to Implement

```go
// Core interfaces
type FallbackHandler interface {
    HandleCertError(ctx context.Context, err error, config *FallbackConfig) (*FallbackStrategy, error)
    ApplyInsecureMode(ctx context.Context, config *InsecureConfig) error
    LogSecurityDecision(decision SecurityDecision) error
    GetRecommendations(err error) []Recommendation
    AttemptAutoRecovery(ctx context.Context, err error, config *RecoveryConfig) (*RecoveryResult, error)
}

// Configuration types
type FallbackConfig struct {
    AllowInsecure       bool
    AutoRecoveryEnabled bool
    MaxRetries          int
    RetryDelay          time.Duration
    Registry            string
}

// Strategy and result types
type FallbackStrategy struct {
    Type            FallbackType
    Description     string
    SecurityImpact  SecurityImpact
    Implementation  string
    RequiresConsent bool
}

type RecoveryResult struct {
    Success       bool
    Method        string
    Actions       []string
    NewConfig     interface{}
    FailureReason string
}
```

## Implementation Strategy

### Phase 1: Core Fallback Infrastructure (~2 hours)
1. Implement FallbackHandler interface
2. Create error categorization system
3. Implement basic strategy generation

### Phase 2: Insecure Mode Handler (~1.5 hours)
1. Implement --insecure flag handling
2. Create warning and consent system
3. Add security decision logging

### Phase 3: Auto-Recovery Mechanisms (~1.5 hours)
1. Implement retry logic with exponential backoff
2. Create certificate refresh mechanisms
3. Add trust store update capabilities

### Phase 4: Testing and Integration (~1 hour)
1. Unit tests for all components
2. Integration with Wave 1 TrustManager
3. Error scenario testing

## Integration Points

### Wave 1 Dependencies
- Uses `TrustManager` interface for trust store updates
- Leverages `CertificateStore` for certificate persistence
- Extends error types from Wave 1 validation
- Integrates with `RegistryConfigManager` for insecure registry configs

### Error Categories to Handle
1. **Certificate Validation Errors**
2. **Chain Building Failures**
3. **Trust Anchor Issues**
4. **Hostname Mismatches**
5. **Expiry Problems**

## Security Requirements

### Explicit Security Model
- NEVER silently bypass certificate checks
- ALL security decisions must be user-initiated
- Security decisions must be logged for audit
- Clear warnings about security implications

### User Consent Requirements
- Explicit --insecure flag required for insecure mode
- Clear warnings about security implications
- Audit trail of all security decisions
- Time-limited insecure mode operations

## Testing Strategy

### Test Coverage Target: 85%

1. **Error Analysis Tests**
   - Test all error categories
   - Verify strategy generation
   - Check recommendation accuracy

2. **Insecure Mode Tests**
   - Test flag handling
   - Verify warning displays
   - Check audit logging

3. **Recovery Tests**
   - Test retry mechanisms
   - Verify timeout handling
   - Check recovery success/failure

4. **Integration Tests**
   - Wave 1 TrustManager integration
   - End-to-end recovery flows
   - Security decision persistence

## File Structure

```
efforts/phase1/wave2/fallback-strategies/pkg/certs/
├── fallback.go           # Main fallback handler
├── fallback_test.go      # Fallback handler tests
├── recovery.go           # Auto-recovery mechanisms
├── recovery_test.go      # Recovery tests
├── insecure.go          # Insecure mode handling
├── insecure_test.go     # Insecure mode tests
└── types_fallback.go    # Fallback-specific types
```

## Success Metrics

- ✅ --insecure flag properly implemented with warnings
- ✅ All security decisions logged for audit
- ✅ Recommendations generated for common certificate errors
- ✅ Recovery mechanisms attempted where safe
- ✅ No silent security bypasses
- ✅ 85% test coverage achieved
- ✅ Integration with Wave 1 components working
- ✅ Code stays under 400 lines (target) / 800 lines (hard limit)

## Risk Mitigation

### Security Risks
- **Bypass Risk**: Require explicit user consent for all bypasses
- **Silent Failure**: All decisions logged and auditable
- **Escalation Risk**: Time-limit all insecure operations

### Technical Risks
- **Recovery Loop**: Implement retry limits and circuit breakers
- **State Corruption**: Atomic operations for trust store updates
- **Integration Issues**: Well-defined interfaces with Wave 1

## Implementation Notes

### Error Message Guidelines
- Clear, actionable error messages
- Specific remediation steps
- Links to documentation where helpful
- Security implications clearly stated

### Logging Requirements
- All security decisions logged with timestamp
- User identification and operation context
- Decision rationale and impact level
- Audit trail compliance

## Dependencies

### External Dependencies
- Go crypto/x509 package for certificate handling
- Go context package for cancellation
- Go time package for timeout handling

### Internal Dependencies
- Wave 1 TrustManager interface
- Wave 1 CertificateStore interface
- Wave 1 error types and patterns
- Wave 1 RegistryConfigManager

## Completion Criteria

1. All interfaces implemented and tested
2. Integration with Wave 1 components verified
3. Security decision audit system working
4. User consent mechanisms functional
5. Auto-recovery mechanisms tested
6. Code quality and test coverage met
7. Documentation complete

This implementation will provide robust, secure fallback strategies for certificate handling while maintaining explicit user control over security decisions.