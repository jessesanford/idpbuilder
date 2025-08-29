package certs

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFallbackHandler(t *testing.T) {
	handler := NewFallbackHandler()
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