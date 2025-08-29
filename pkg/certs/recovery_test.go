package certs

import (
	"context"
	"crypto/x509"
	"errors"
	"testing"
	"time"
)

// TestRecoveryManager tests the basic recovery manager functionality
func TestRecoveryManager(t *testing.T) {
	config := &RecoveryConfig{
		EnableCertRefresh:   true,
		EnableTrustUpdate:   true,
		EnableChainRepair:   true,
		MaxAttempts:         3,
		Timeout:             100 * time.Millisecond,
		CircuitBreakerThreshold: 5,
		CircuitBreakerTimeout:   60 * time.Second,
	}

	manager := NewRecoveryManager(config)
	if manager == nil {
		t.Fatal("NewRecoveryManager returned nil")
	}

	// Test initial state
	state := manager.GetState()
	if !state.Healthy {
		t.Error("Manager should be healthy initially")
	}
	if state.CircuitOpen {
		t.Error("Circuit should be closed initially")
	}
}

// TestCircuitBreaker tests the circuit breaker pattern
func TestCircuitBreaker(t *testing.T) {
	breaker := &CircuitBreaker{
		failureThreshold: 3,
		resetTimeout:     100 * time.Millisecond,
		state:           CircuitClosed,
	}

	// Test initial state - should allow execution
	if !breaker.CanExecute() {
		t.Error("Circuit should allow execution initially")
	}

	// Record failures to trigger circuit opening
	for i := 0; i < 3; i++ {
		breaker.RecordFailure()
	}

	// Circuit should be open now
	if breaker.CanExecute() {
		t.Error("Circuit should be open after threshold failures")
	}

	// Wait for reset timeout and test half-open
	time.Sleep(150 * time.Millisecond)
	if !breaker.CanExecute() {
		t.Error("Circuit should allow execution after timeout (half-open)")
	}

	// Record success to close circuit
	breaker.RecordSuccess()
	if !breaker.CanExecute() {
		t.Error("Circuit should be closed after success")
	}
}

// TestAutoRecoveryStrategy tests connection recovery with exponential backoff
func TestAutoRecoveryStrategy(t *testing.T) {
	strategy := NewAutoRecoveryStrategy()
	ctx := context.Background()

	// Test CanRecover
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errors.New("connection timeout"), true},
		{errors.New("network unreachable"), true},
		{errors.New("no such host"), true},
		{errors.New("certificate expired"), false},
		{errors.New("connection refused"), true},
	}

	for _, test := range tests {
		if strategy.CanRecover(test.err) != test.expected {
			t.Errorf("CanRecover(%v) = %t, expected %t", test.err, strategy.CanRecover(test.err), test.expected)
		}
	}

	// Test Attempt with connection error
	input := &ValidationInput{
		Registry:  "test-registry",
		Operation: "pull",
		Error:     errors.New("connection timeout"),
	}

	result, err := strategy.Attempt(ctx, input)
	if err != nil && result.Success {
		t.Error("Strategy should handle connection errors")
	}

	// Test metrics are updated
	metrics := strategy.GetMetrics()
	if metrics.AttemptCount == 0 {
		t.Error("Attempt count should be updated")
	}
}

// TestCertificateRefreshStrategy tests certificate refresh functionality
func TestCertificateRefreshStrategy(t *testing.T) {
	strategy := NewCertificateRefreshStrategy()
	ctx := context.Background()

	// Test CanRecover
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errors.New("certificate has expired"), true},
		{errors.New("certificate is not valid"), true},
		{errors.New("certificate not yet valid"), true},
		{errors.New("connection timeout"), false},
	}

	for _, test := range tests {
		if strategy.CanRecover(test.err) != test.expected {
			t.Errorf("CanRecover(%v) = %t, expected %t", test.err, strategy.CanRecover(test.err), test.expected)
		}
	}

	// Test Attempt with expired certificate
	input := &ValidationInput{
		Registry:  "test-registry",
		Operation: "pull",
		Error:     errors.New("certificate has expired"),
		Certificates: []*x509.Certificate{},
	}

	result, err := strategy.Attempt(ctx, input)
	if err != nil {
		t.Errorf("Certificate refresh should not return error: %v", err)
	}
	if !result.Success {
		t.Error("Certificate refresh should succeed")
	}
	if result.Strategy != "refresh-certificates" {
		t.Errorf("Expected strategy 'refresh-certificates', got %s", result.Strategy)
	}

	// Check that NewConfig is set
	if result.NewConfig == nil {
		t.Error("NewConfig should be set on successful refresh")
	}
}

// TestTrustStoreUpdateStrategy tests trust store update functionality
func TestTrustStoreUpdateStrategy(t *testing.T) {
	strategy := NewTrustStoreUpdateStrategy()
	ctx := context.Background()

	// Test CanRecover
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errors.New("self signed certificate"), true},
		{errors.New("unknown authority"), true},
		{errors.New("certificate signed by unknown authority"), true},
		{errors.New("x509: certificate signed by unknown authority"), true},
		{errors.New("certificate has expired"), false},
	}

	for _, test := range tests {
		if strategy.CanRecover(test.err) != test.expected {
			t.Errorf("CanRecover(%v) = %t, expected %t", test.err, strategy.CanRecover(test.err), test.expected)
		}
	}

	// Test Attempt with self-signed certificate
	input := &ValidationInput{
		Registry:  "test-registry",
		Operation: "pull",
		Error:     errors.New("self signed certificate"),
		Certificates: []*x509.Certificate{},
	}

	result, err := strategy.Attempt(ctx, input)
	if err != nil {
		t.Errorf("Trust store update should not return error: %v", err)
	}
	if !result.Success {
		t.Error("Trust store update should succeed")
	}
	if result.Strategy != "update-trust-store" {
		t.Errorf("Expected strategy 'update-trust-store', got %s", result.Strategy)
	}

	// Check that NewConfig is set
	if result.NewConfig == nil {
		t.Error("NewConfig should be set on successful update")
	}
}

// TestRecoveryWithTimeout tests recovery operations with context timeout
func TestRecoveryWithTimeout(t *testing.T) {
	config := &RecoveryConfig{
		EnableCertRefresh: true,
		MaxAttempts:       3,
		Timeout:           50 * time.Millisecond,
	}

	manager := NewRecoveryManager(config)
	
	// Create a context that times out quickly
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	input := &ValidationInput{
		Registry:  "test-registry",
		Operation: "pull",
		Error:     errors.New("certificate has expired"),
	}

	result, err := manager.Recover(ctx, input.Error, input)
	
	// Should handle timeout gracefully
	if err != nil && err != context.DeadlineExceeded {
		t.Logf("Recovery with timeout returned error: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil even on timeout")
	}
}

// TestRecoveryStrategiesRegistration tests strategy registration and management
func TestRecoveryStrategiesRegistration(t *testing.T) {
	manager := NewRecoveryManager(nil).(*DefaultRecoveryManager)

	// Test that default strategies are registered
	expectedStrategies := []string{
		"auto-recover-connection",
		"refresh-certificates", 
		"update-trust-store",
	}

	for _, strategyName := range expectedStrategies {
		if _, exists := manager.strategies[strategyName]; !exists {
			t.Errorf("Expected strategy %s to be registered", strategyName)
		}
	}

	// Test registering a custom strategy
	customStrategy := &testRecoveryStrategy{name: "custom-test"}
	manager.RegisterStrategy("custom-test", customStrategy)

	if _, exists := manager.strategies["custom-test"]; !exists {
		t.Error("Custom strategy should be registered")
	}
}

// TestRecoveryWithMultipleStrategies tests recovery with multiple applicable strategies
func TestRecoveryWithMultipleStrategies(t *testing.T) {
	manager := NewRecoveryManager(nil)
	ctx := context.Background()

	// Test different error types that trigger different strategies
	testCases := []struct {
		name         string
		err          error
		expectSuccess bool
	}{
		{
			name:         "ConnectionError",
			err:          errors.New("connection timeout"),
			expectSuccess: false, // May succeed or fail due to simulation
		},
		{
			name:         "ExpiredCert",
			err:          errors.New("certificate has expired"),
			expectSuccess: true,
		},
		{
			name:         "SelfSigned",
			err:          errors.New("self signed certificate"),
			expectSuccess: true,
		},
		{
			name:         "UnknownAuthority",
			err:          errors.New("certificate signed by unknown authority"),
			expectSuccess: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &ValidationInput{
				Registry:  "test-registry",
				Operation: "pull",
				Error:     tc.err,
			}

			result, _ := manager.Recover(ctx, tc.err, input)
			
			if result == nil {
				t.Fatal("Result should not be nil")
			}

			// Log results for debugging
			t.Logf("Test: %s", tc.name)
			t.Logf("Success: %t", result.Success)
			t.Logf("Actions: %v", result.Actions)
			if !result.Success {
				t.Logf("Failure message: %s", result.Message)
			}
		})
	}
}

// testRecoveryStrategy is a simple test implementation of RecoveryStrategy
type testRecoveryStrategy struct {
	name    string
	metrics RecoveryMetrics
}

func (t *testRecoveryStrategy) CanRecover(err error) bool {
	return err != nil
}

func (t *testRecoveryStrategy) Attempt(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	return &ValidationResult{
		Success:  true,
		Strategy: t.name,
		Message:  "Test recovery successful",
	}, nil
}

func (t *testRecoveryStrategy) GetMetrics() RecoveryMetrics {
	return t.metrics
}

// TestHelperFunctions tests utility functions
func TestHelperFunctions(t *testing.T) {
	// Test containsAny
	if !containsAny("hello world", []string{"world"}) {
		t.Error("containsAny should find 'world' in 'hello world'")
	}
	if containsAny("hello world", []string{"foo"}) {
		t.Error("containsAny should not find 'foo' in 'hello world'")
	}
	if !containsAny("certificate has expired", []string{"expired", "invalid"}) {
		t.Error("containsAny should find 'expired' in certificate error")
	}

	// Test min
	if min(5, 3) != 3 {
		t.Error("min(5, 3) should return 3")
	}
	if min(1, 1) != 1 {
		t.Error("min(1, 1) should return 1")
	}
}