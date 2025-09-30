# SPLIT-PLAN-002: Retry Mechanism with Backoff

## Split 002 of 2: Retry and Backoff Logic
**Planner**: Code Reviewer code-reviewer (same for ALL splits)
**Parent Effort**: E1.2.2-registry-authentication
**Estimated Size**: ~346 lines (well within 800 line limit)

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: Split 001 of phase1/wave2/E1.2.2-registry-authentication
  - Path: efforts/phase1/wave2/E1.2.2-registry-authentication/split-001/
  - Branch: phase1/wave2/registry-authentication-split-001
  - Summary: Implemented core authentication with Authenticator, credentials, and error types

- **This Split**: Split 002 of phase1/wave2/E1.2.2-registry-authentication
  - Path: efforts/phase1/wave2/E1.2.2-registry-authentication/split-002/
  - Branch: phase1/wave2/registry-authentication-split-002

- **Next Split**: None (final split of this effort)

- **File Boundaries**:
  - Previous Split End: pkg/push/errors/auth_errors.go
  - This Split Start: pkg/push/retry/retry.go
  - This Split End: pkg/push/retry/errors.go (last file)

## Files in This Split (EXCLUSIVE - no overlap with other splits)

1. **pkg/push/retry/retry.go** (~211 lines)
   - Core retry logic implementation
   - RetryableFunc type definition
   - IsRetryable function for error classification
   - WithRetry wrapper for operations
   - Network error detection patterns
   - HTTP status code retry logic

2. **pkg/push/retry/backoff.go** (~125 lines)
   - Backoff strategy interface
   - Exponential backoff implementation
   - Jitter for distributed systems
   - Max attempts configuration
   - Delay calculation algorithms
   - Context-aware waiting

3. **pkg/push/retry/errors.go** (~10 lines)
   - Retry-specific error types
   - MaxRetriesExceededError
   - Error wrapping utilities

## Functionality Scope

This split provides complete retry mechanism:
- Intelligent retry decision logic
- Configurable backoff strategies
- Network transient error detection
- HTTP status code based retries
- Context cancellation support
- Jittered exponential backoff

## Dependencies

### External Dependencies:
- Standard library: context, fmt, net, net/http, time, math/rand

### Internal Dependencies from Split 001:
- NONE directly (retry logic is independent)
- May import pkg/push/errors in integration scenarios

**Note**: While retry logic can work with auth errors from Split 001, it's designed to be generic and work with any error type.

## Implementation Instructions

1. **Setup Split Directory**:
   ```bash
   cd efforts/phase1/wave2/E1.2.2-registry-authentication
   # Ensure split-001 is complete and pushed
   git checkout phase1/wave2/registry-authentication-split-001

   mkdir -p split-002
   cd split-002
   git checkout -b phase1/wave2/registry-authentication-split-002
   ```

2. **Create Package Structure**:
   ```bash
   mkdir -p pkg/push/retry
   ```

3. **Copy Split 001 Dependencies** (if needed for testing):
   ```bash
   # Copy error types from split-001 for integration testing
   cp -r ../split-001/pkg/push/errors pkg/push/
   ```

4. **Implement Files in Order**:
   - Start with pkg/push/retry/errors.go (simplest, no dependencies)
   - Then pkg/push/retry/backoff.go (mathematical logic)
   - Finally pkg/push/retry/retry.go (uses backoff and errors)

5. **Testing Requirements**:
   - Unit tests for backoff calculations
   - Tests for retry decision logic
   - Context cancellation tests
   - Maximum retry limit tests
   - Jitter distribution tests
   - Mock failing operations with retries

6. **Verification**:
   ```bash
   go build ./pkg/push/retry/...
   go test ./pkg/push/retry/... -v
   ```

7. **Line Count Verification**:
   ```bash
   $PROJECT_ROOT/tools/line-counter.sh
   # Must show <800 lines (target ~346)
   ```

## Quality Requirements

- Comprehensive godoc comments
- Thread-safe implementations
- No goroutine leaks
- Proper context handling
- Deterministic testing with fixed seeds
- Clear retry reason logging

## Testing Strategy

### Unit Tests Required:
1. **Backoff Tests**:
   - Exponential growth validation
   - Jitter within bounds
   - Max delay enforcement
   - Reset functionality

2. **Retry Logic Tests**:
   - Transient vs permanent error classification
   - HTTP status code decisions
   - Max attempts enforcement
   - Context cancellation mid-retry

3. **Integration Tests**:
   - Retry with mock HTTP client
   - Retry with simulated network errors
   - Backoff timing validation

## Integration Notes

- This split provides generic retry capability
- Can be used with auth operations from Split 001
- Will be used by image push operations in future efforts
- Designed to handle both network and application errors

## Performance Considerations

- Backoff prevents thundering herd
- Jitter spreads retry load
- Context support enables timeout control
- Efficient error type checking

## Success Criteria

- [ ] All 3 files implemented completely
- [ ] Unit tests achieve >85% coverage
- [ ] No race conditions (go test -race)
- [ ] Backoff calculations are correct
- [ ] go build succeeds without errors
- [ ] go test passes all tests
- [ ] Line count <800 (target ~346)
- [ ] Proper error wrapping maintained
- [ ] Documentation complete

## Merge Strategy

1. Complete Split 001 and ensure it's merged
2. Branch Split 002 from Split 001
3. Implement retry logic independently
4. Test with mock scenarios
5. Integrate with Split 001 error types if needed
6. Merge Split 002 back to parent effort branch