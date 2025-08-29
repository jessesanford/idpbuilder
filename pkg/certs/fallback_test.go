package certs

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Mock implementations for testing
type mockTrustManager struct{}
func (m *mockTrustManager) AddCertificate(ctx context.Context, registry string, cert interface{}) error { return nil }
func (m *mockTrustManager) SetInsecureRegistry(ctx context.Context, registry string, insecure bool) error { return nil }

type mockCertStore struct{}
func (m *mockCertStore) Store(registry string, cert interface{}) error { return nil }
func (m *mockCertStore) Load(registry string, fingerprint string) (interface{}, error) { return nil, nil }

type mockConfigMgr struct{}
func (m *mockConfigMgr) UpdateInsecureRegistry(registry string, insecure bool) error { return nil }
func (m *mockConfigMgr) GetInsecureRegistries() ([]string, error) { return []string{}, nil }

func TestFallbackHandler(t *testing.T) {
	handler := NewFallbackHandler(&mockTrustManager{}, &mockCertStore{}, &mockConfigMgr{})
	if handler == nil {
		t.Fatal("NewFallbackHandler returned nil")
	}
	var _ FallbackHandler = handler

	ctx := context.Background()
	config := &FallbackConfig{AllowInsecure: true, Registry: "test"}

	// Test error handling
	tests := []struct {
		err          error
		expectedType FallbackType
	}{
		{nil, FallbackNone},
		{errors.New("self signed certificate"), FallbackInsecure},
		{errors.New("certificate has expired"), FallbackRetry},
		{errors.New("no such host"), FallbackRetry},
	}

	for _, tt := range tests {
		strategy, err := handler.HandleCertError(ctx, tt.err, config)
		if err != nil || strategy.Type != tt.expectedType {
			t.Errorf("HandleCertError failed for %v", tt.err)
		}
	}

	// Test insecure mode
	insecureConfig := &InsecureConfig{
		Registry: "test", Operation: "push", Duration: 5 * time.Minute, Reason: "testing",
	}
	if err := handler.ApplyInsecureMode(ctx, insecureConfig); err != nil {
		t.Error("ApplyInsecureMode failed")
	}

	// Test recommendations
	recs := handler.GetRecommendations(errors.New("self signed certificate"))
	if len(recs) < 1 {
		t.Error("Expected recommendations")
	}

	// Test auto-recovery
	recoveryConfig := &RecoveryConfig{MaxAttempts: 3, Timeout: 1 * time.Second}
	result, err := handler.AttemptAutoRecovery(ctx, errors.New("no such host"), recoveryConfig)
	if err != nil || result.Method != "dns-retry" {
		t.Error("Auto-recovery failed")
	}
}

func TestEnumValues(t *testing.T) {
	// Verify enum constants are sequential
	fallbackTypes := []FallbackType{FallbackNone, FallbackInsecure, FallbackAlternateTrust, FallbackManualTrust, FallbackRetry}
	impactLevels := []ImpactLevel{ImpactMinimal, ImpactModerate, ImpactHigh, ImpactCritical}
	priorities := []RecommendationPriority{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical}
	
	for i, val := range fallbackTypes {
		if int(val) != i {
			t.Errorf("FallbackType enum mismatch: expected %d, got %d", i, int(val))
		}
	}
	
	for i, val := range impactLevels {
		if int(val) != i {
			t.Errorf("ImpactLevel enum mismatch: expected %d, got %d", i, int(val))
		}
	}
	
	for i, val := range priorities {
		if int(val) != i {
			t.Errorf("RecommendationPriority enum mismatch: expected %d, got %d", i, int(val))
		}
	}
}

// Additional test coverage to reach 85%
func TestWave1Integration(t *testing.T) {
	handler := NewFallbackHandler(&mockTrustManager{}, &mockCertStore{}, &mockConfigMgr{})
	ctx := context.Background()
	
	// Test actual Wave 1 integration
	config := &InsecureConfig{
		Registry: "test.registry.com",
		Operation: "pull",
		Duration: 1*time.Hour,
		Reason: "development testing",
		RequireExplicit: true,
	}
	
	err := handler.ApplyInsecureMode(ctx, config)
	if err != nil {
		t.Errorf("ApplyInsecureMode failed: %v", err)
	}
}

func TestErrorRecovery(t *testing.T) {
	handler := NewFallbackHandler(&mockTrustManager{}, &mockCertStore{}, &mockConfigMgr{})
	ctx := context.Background()
	config := &RecoveryConfig{MaxAttempts: 3}
	
	// Test network error recovery
	networkErr := errors.New("no such host")
	result, err := handler.AttemptAutoRecovery(ctx, networkErr, config)
	if err != nil {
		t.Errorf("Network error recovery failed: %v", err)
	}
	if result == nil {
		t.Error("Expected recovery result")
	}
}

func TestComplexErrorScenarios(t *testing.T) {
	handler := NewFallbackHandler(&mockTrustManager{}, &mockCertStore{}, &mockConfigMgr{})
	ctx := context.Background()
	config := &FallbackConfig{AllowInsecure: false, Registry: "secure.registry.com"}
	
	// Test various error types
	errorMessages := []string{
		"certificate has expired",
		"certificate name does not match", 
		"unknown error type",
	}
	
	for _, errStr := range errorMessages {
		strategy, err := handler.HandleCertError(ctx, errors.New(errStr), config)
		if err != nil {
			t.Errorf("HandleCertError failed for %s: %v", errStr, err)
		}
		if strategy == nil {
			t.Errorf("Expected strategy for error: %s", errStr)
		}
	}
}

func TestFallbackStrategyHelpers(t *testing.T) {
	handler := NewFallbackHandler(&mockTrustManager{}, &mockCertStore{}, &mockConfigMgr{}).(*DefaultFallbackHandler)
	config := &FallbackConfig{AllowInsecure: true, Registry: "test.com"}
	
	// Test all strategy creation methods
	selfSignedStrategy := handler.createSelfSignedStrategy(config)
	if selfSignedStrategy.Type != FallbackInsecure {
		t.Error("Expected insecure fallback for self-signed when allowed")
	}
	
	// Test with insecure disabled
	config.AllowInsecure = false
	trustStrategy := handler.createSelfSignedStrategy(config)
	if trustStrategy.Type != FallbackManualTrust {
		t.Error("Expected manual trust when insecure not allowed")
	}
	
	expiredStrategy := handler.createExpiredStrategy(config)
	if expiredStrategy.Type != FallbackRetry {
		t.Error("Expected retry strategy for expired certificates")
	}
	
	hostnameStrategy := handler.createHostnameStrategy(config)
	if hostnameStrategy.Type != FallbackRetry {
		t.Error("Expected retry strategy for hostname errors")
	}
	
	networkStrategy := handler.createNetworkStrategy(config)
	if networkStrategy.Type != FallbackRetry {
		t.Error("Expected retry strategy for network errors")
	}
	
	genericStrategy := handler.createGenericStrategy(config)
	if genericStrategy.Type != FallbackRetry {
		t.Error("Expected retry strategy for generic errors")
	}
}

func TestGetRecommendationsComprehensive(t *testing.T) {
	handler := NewFallbackHandler(&mockTrustManager{}, &mockCertStore{}, &mockConfigMgr{})
	
	// Test all recommendation scenarios
	testCases := []struct {
		errorMsg string
		expectedRecs int
	}{
		{"self signed certificate", 2}, // Trust store + insecure
		{"certificate has expired", 1}, // Contact admin
		{"certificate name does not match", 1}, // Verify URL
		{"unknown network error", 1}, // Check connectivity
		{"", 0}, // No error, no recommendations
	}
	
	for _, tc := range testCases {
		var err error
		if tc.errorMsg != "" {
			err = errors.New(tc.errorMsg)
		}
		
		recs := handler.GetRecommendations(err)
		if len(recs) != tc.expectedRecs {
			t.Errorf("For error '%s', expected %d recommendations, got %d", tc.errorMsg, tc.expectedRecs, len(recs))
		}
		
		// Verify recommendations have proper structure
		for _, rec := range recs {
			if rec.Title == "" || rec.Description == "" {
				t.Errorf("Recommendation missing title or description for error: %s", tc.errorMsg)
			}
		}
	}
}