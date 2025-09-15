# Integration Test Plan - idpbuilder OCI Build/Push Project

## Executive Summary

This document outlines the comprehensive integration test plan for the idpbuilder OCI build/push project, which combines Phase 1 (certificate validation) and Phase 2 (OCI build/push) features. The tests focus on verifying that both phases work correctly together, handling real-world scenarios with proper certificate validation during OCI operations.

## Test Scope

### Features to Test
1. **Certificate Validation (Phase 1)**
   - Certificate extraction from Kind clusters
   - Chain validation with multiple modes (Strict, Lenient, Insecure)
   - Trust store management
   - Fallback strategies for certificate retrieval

2. **OCI Build/Push (Phase 2)**
   - Image building from context directories
   - Image pushing to Gitea registry
   - Progress tracking and reporting
   - Error handling and recovery

3. **Integration Points**
   - Certificate validation during push operations
   - Fallback to insecure mode when certificates unavailable
   - Certificate storage and reuse across operations
   - Error propagation between phases

## Test Categories

### 1. End-to-End Integration Tests

#### Test 1.1: Secure Build and Push Workflow
**Objective**: Verify complete workflow with certificate validation
**Steps**:
1. Build an OCI image from a test context
2. Extract certificates from Kind cluster
3. Push image to Gitea registry with certificate validation
4. Verify image is accessible in registry

**Expected Results**:
- Build completes successfully
- Certificates are automatically extracted and validated
- Push uses secure HTTPS connection
- No certificate warnings in output

#### Test 1.2: Build and Push with Certificate Rotation
**Objective**: Test handling of certificate changes during operations
**Steps**:
1. Build image with initial certificates
2. Simulate certificate rotation in Kind cluster
3. Attempt push with old certificates (should fail)
4. Re-extract certificates and retry push
5. Verify successful push with new certificates

**Expected Results**:
- First push attempt fails with certificate error
- Certificate re-extraction succeeds
- Second push attempt succeeds with new certificates

#### Test 1.3: Multi-Image Push with Certificate Caching
**Objective**: Verify certificate reuse across multiple operations
**Steps**:
1. Extract certificates once
2. Build and push multiple images (5+)
3. Verify all pushes use cached certificates
4. Monitor certificate validation performance

**Expected Results**:
- Certificates extracted only once
- All subsequent pushes reuse cached certificates
- Performance improvement after first operation

### 2. Certificate Validation Integration Tests

#### Test 2.1: Strict Mode Validation with OCI Push
**Objective**: Ensure strict certificate validation works with push
**Steps**:
1. Configure strict validation mode
2. Attempt push with valid certificate chain
3. Attempt push with incomplete chain (should fail)
4. Verify appropriate error messages

**Expected Results**:
- Valid chain push succeeds
- Invalid chain push fails with clear error
- No fallback to insecure mode in strict mode

#### Test 2.2: Lenient Mode with Self-Signed Certificates
**Objective**: Test lenient mode allows self-signed certificates
**Steps**:
1. Configure lenient validation mode
2. Generate self-signed certificate for test registry
3. Push image using self-signed certificate
4. Verify warning messages but successful push

**Expected Results**:
- Push succeeds with self-signed certificate
- Warning messages displayed about certificate status
- Image successfully stored in registry

#### Test 2.3: Certificate Chain Depth Validation
**Objective**: Test maximum chain depth enforcement
**Steps**:
1. Create certificate chain with 5+ intermediates
2. Configure max chain depth to 3
3. Attempt push (should fail)
4. Increase max depth and retry
5. Verify successful push

**Expected Results**:
- Push fails when chain exceeds max depth
- Clear error message about chain depth
- Push succeeds after increasing limit

### 3. Fallback Strategy Integration Tests

#### Test 3.1: System Certificate Fallback
**Objective**: Test fallback to system certificates
**Steps**:
1. Remove Kind cluster certificates
2. Install registry certificate in system store
3. Attempt push operation
4. Verify system certificate fallback activated

**Expected Results**:
- Fallback strategy triggered automatically
- System certificates used successfully
- Informational message about fallback usage

#### Test 3.2: Manual Certificate Fallback
**Objective**: Test manual certificate specification
**Steps**:
1. Provide certificate via environment variable
2. Attempt push without Kind cluster
3. Verify manual certificate used
4. Test with invalid manual certificate

**Expected Results**:
- Manual certificate takes precedence
- Clear error for invalid certificate
- Success with valid manual certificate

#### Test 3.3: Insecure Mode Fallback
**Objective**: Test final fallback to insecure mode
**Steps**:
1. Ensure no certificates available
2. Attempt push (should fail)
3. Enable --insecure flag
4. Retry push operation
5. Verify warning messages

**Expected Results**:
- Initial push fails with certificate error
- Insecure mode push succeeds
- Prominent warning about insecure operation

### 4. Error Handling Integration Tests

#### Test 4.1: Certificate Extraction Failure Recovery
**Objective**: Test recovery from certificate extraction failures
**Steps**:
1. Simulate Kind API unavailable
2. Attempt certificate extraction (fails)
3. Trigger fallback strategies
4. Complete push with fallback

**Expected Results**:
- Clear error about extraction failure
- Automatic fallback strategy activation
- Successful completion via fallback

#### Test 4.2: Registry Connection Failure with Retry
**Objective**: Test connection failure handling
**Steps**:
1. Start push to unavailable registry
2. Verify connection timeout
3. Start registry during retry window
4. Verify successful completion

**Expected Results**:
- Initial connection fails gracefully
- Retry mechanism activates
- Push completes after registry available

#### Test 4.3: Partial Push Failure Recovery
**Objective**: Test recovery from interrupted push
**Steps**:
1. Start pushing large image
2. Interrupt connection mid-push
3. Retry push operation
4. Verify resume from interruption point

**Expected Results**:
- Push interruption handled gracefully
- Resume capability if supported
- Complete push on retry

### 5. Performance Integration Tests

#### Test 5.1: Large Image Push with Certificates
**Objective**: Test performance with large images
**Steps**:
1. Build 1GB+ test image
2. Push with certificate validation
3. Measure validation overhead
4. Compare with insecure push

**Expected Results**:
- Certificate validation < 5% overhead
- Progress reporting remains accurate
- No timeout issues

#### Test 5.2: Concurrent Push Operations
**Objective**: Test concurrent pushes with shared certificates
**Steps**:
1. Extract certificates once
2. Start 5 concurrent push operations
3. Monitor certificate access patterns
4. Verify all pushes succeed

**Expected Results**:
- No certificate access conflicts
- All pushes complete successfully
- Reasonable performance scaling

### 6. Configuration Integration Tests

#### Test 6.1: Environment Variable Configuration
**Objective**: Test environment-based configuration
**Steps**:
1. Set IDPBUILDER_CERT_PATH variable
2. Set IDPBUILDER_REGISTRY variable
3. Run build and push without flags
4. Verify configuration applied

**Expected Results**:
- Environment variables respected
- No command-line flags needed
- Correct registry and certificates used

#### Test 6.2: Configuration File Integration
**Objective**: Test config file with certificate settings
**Steps**:
1. Create .idpbuilder.yaml with cert config
2. Run operations without explicit settings
3. Verify config file settings applied
4. Test config override with flags

**Expected Results**:
- Config file loaded automatically
- Certificate settings applied correctly
- Command flags override config file

### 7. Kubernetes Integration Tests

#### Test 7.1: Kind Cluster Certificate Rotation
**Objective**: Test with Kind cluster lifecycle
**Steps**:
1. Create fresh Kind cluster
2. Extract certificates
3. Delete and recreate cluster
4. Verify certificate invalidation
5. Re-extract and verify new certificates

**Expected Results**:
- Old certificates properly invalidated
- New certificates extracted successfully
- Clear indication of certificate change

#### Test 7.2: Multi-Cluster Certificate Management
**Objective**: Test with multiple Kind clusters
**Steps**:
1. Create two Kind clusters with different names
2. Extract certificates from both
3. Push to registries in each cluster
4. Verify correct certificate selection

**Expected Results**:
- Certificates isolated per cluster
- Correct certificate used for each registry
- No certificate confusion

### 8. Diagnostic Integration Tests

#### Test 8.1: Certificate Diagnostics During Push
**Objective**: Test diagnostic information availability
**Steps**:
1. Enable verbose/debug mode
2. Perform push operation
3. Verify certificate details in output
4. Check diagnostic logs

**Expected Results**:
- Certificate chain displayed in debug
- Validation steps shown
- Performance metrics available

#### Test 8.2: Troubleshooting Failed Push
**Objective**: Test troubleshooting information
**Steps**:
1. Attempt push with known bad certificate
2. Review error messages
3. Follow suggested remediation
4. Verify fix resolves issue

**Expected Results**:
- Clear, actionable error messages
- Specific remediation suggestions
- Successful resolution following guidance

## Test Environment Requirements

### Infrastructure
- Kind cluster (v0.20.0+)
- Gitea registry instance
- Test OCI images of various sizes
- Certificate generation tools

### Test Data
- Valid certificate chains
- Self-signed certificates
- Expired certificates
- Malformed certificates
- Test container contexts (small, medium, large)

### Tools
- Go test framework
- Container runtime (Docker/Podman)
- Network simulation tools (for failure testing)
- Performance monitoring tools

## Test Execution Strategy

### Phase 1: Unit Test Verification
- Verify all unit tests pass
- Check test coverage meets requirements (>80%)
- Identify any gaps in unit testing

### Phase 2: Integration Test Implementation
1. Set up test infrastructure
2. Implement test helpers and utilities
3. Create test data generators
4. Implement integration test suites

### Phase 3: Test Execution
- Run tests in isolation first
- Execute full integration suite
- Perform stress testing
- Conduct failure injection testing

### Phase 4: Performance Testing
- Baseline performance measurements
- Load testing with concurrent operations
- Certificate validation overhead analysis
- Resource consumption monitoring

## Success Criteria

### Functional Requirements
- ✅ All certificate validation modes work correctly
- ✅ OCI build and push operations complete successfully
- ✅ Fallback strategies activate appropriately
- ✅ Error messages are clear and actionable

### Performance Requirements
- ✅ Certificate validation adds <5% overhead
- ✅ Push operations complete within expected timeframes
- ✅ Concurrent operations scale linearly
- ✅ Memory usage remains bounded

### Reliability Requirements
- ✅ Graceful handling of all failure scenarios
- ✅ Automatic recovery where possible
- ✅ No data corruption or loss
- ✅ Consistent behavior across environments

## Risk Mitigation

### High Risk Areas
1. **Certificate Expiration**: Test with soon-to-expire certificates
2. **Network Interruptions**: Implement retry mechanisms
3. **Large Image Handling**: Test with realistic image sizes
4. **Concurrent Access**: Verify thread safety

### Mitigation Strategies
- Comprehensive error handling tests
- Timeout and retry configuration tests
- Resource cleanup verification
- State consistency checks

## Test Maintenance

### Regular Updates
- Update tests for new certificate formats
- Add tests for new OCI specifications
- Refresh test certificates quarterly
- Review and update test data

### Documentation
- Maintain test case documentation
- Document known issues and workarounds
- Keep troubleshooting guide current
- Update performance baselines

## Conclusion

This integration test plan ensures comprehensive validation of the idpbuilder OCI build/push project with certificate validation. The tests verify that Phase 1 and Phase 2 features work correctly together, handle edge cases gracefully, and provide a robust solution for secure container image operations.

### Next Steps
1. Review and approve test plan
2. Set up test infrastructure
3. Implement test suites
4. Execute test plan
5. Address any findings
6. Document results

### Timeline
- Test Infrastructure Setup: 2 days
- Test Implementation: 5 days
- Test Execution: 3 days
- Result Analysis: 2 days
- Total: 12 days

### Resources Required
- 2 QA Engineers
- 1 DevOps Engineer (infrastructure)
- Test environment (Kind, Gitea)
- CI/CD pipeline integration