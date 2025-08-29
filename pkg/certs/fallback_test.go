package certs

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewFallbackHandler(t *testing.T) {
	handler := NewFallbackHandler()
	if handler == nil {
		t.Fatal("NewFallbackHandler returned nil")
	}
	
	// Verify it implements the interface
	var _ FallbackHandler = handler
}

func TestHandleCertError(t *testing.T) {
	handler := NewFallbackHandler()
	ctx := context.Background()
	config := &FallbackConfig{
		AllowInsecure: true,
		Registry:      "test-registry",
	}

	tests := []struct {
		name           string
		err            error
		expectedType   FallbackType
		expectedConsent bool
	}{
		{
			name:           "no error",
			err:            nil,
			expectedType:   FallbackNone,
			expectedConsent: false,
		},
		{
			name:           "self signed certificate",
			err:            errors.New("self signed certificate"),
			expectedType:   FallbackInsecure,
			expectedConsent: true,
		},
		{
			name:           "certificate signed by unknown authority",
			err:            errors.New("certificate signed by unknown authority"),
			expectedType:   FallbackInsecure,
			expectedConsent: true,
		},
		{
			name:           "certificate has expired",
			err:            errors.New("certificate has expired"),
			expectedType:   FallbackRetry,
			expectedConsent: false,
		},
		{
			name:           "certificate name does not match",
			err:            errors.New("certificate name does not match"),
			expectedType:   FallbackRetry,
			expectedConsent: true,
		},
		{
			name:           "no such host",
			err:            errors.New("no such host"),
			expectedType:   FallbackRetry,
			expectedConsent: false,
		},
		{
			name:           "generic error",
			err:            errors.New("unknown error"),
			expectedType:   FallbackRetry,
			expectedConsent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := handler.HandleCertError(ctx, tt.err, config)
			if err != nil {
				t.Fatalf("HandleCertError failed: %v", err)
			}
			
			if strategy.Type != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, strategy.Type)
			}
			
			if strategy.RequiresConsent != tt.expectedConsent {
				t.Errorf("Expected consent %v, got %v", tt.expectedConsent, strategy.RequiresConsent)
			}
		})
	}
}

func TestHandleCertErrorWithoutInsecure(t *testing.T) {
	handler := NewFallbackHandler()
	ctx := context.Background()
	config := &FallbackConfig{
		AllowInsecure: false,
		Registry:      "test-registry",
	}

	err := errors.New("self signed certificate")
	strategy, strategyErr := handler.HandleCertError(ctx, err, config)
	
	if strategyErr != nil {
		t.Fatalf("HandleCertError failed: %v", strategyErr)
	}
	
	// Should suggest manual trust instead of insecure mode
	if strategy.Type != FallbackManualTrust {
		t.Errorf("Expected FallbackManualTrust, got %v", strategy.Type)
	}
}

func TestApplyInsecureMode(t *testing.T) {
	handler := NewFallbackHandler()
	ctx := context.Background()
	
	config := &InsecureConfig{
		Registry:        "test-registry",
		Operation:       "push",
		Duration:        5 * time.Minute,
		Reason:          "testing",
		RequireExplicit: false,
	}
	
	err := handler.ApplyInsecureMode(ctx, config)
	if err != nil {
		t.Fatalf("ApplyInsecureMode failed: %v", err)
	}
}

func TestGetRecommendations(t *testing.T) {
	handler := NewFallbackHandler()

	tests := []struct {
		name        string
		err         error
		minRecs     int
		hasHighPrio bool
	}{
		{
			name:        "no error",
			err:         nil,
			minRecs:     0,
			hasHighPrio: false,
		},
		{
			name:        "self signed certificate",
			err:         errors.New("self signed certificate"),
			minRecs:     2,
			hasHighPrio: true,
		},
		{
			name:        "certificate has expired",
			err:         errors.New("certificate has expired"),
			minRecs:     1,
			hasHighPrio: false,
		},
		{
			name:        "certificate name does not match",
			err:         errors.New("certificate name does not match"),
			minRecs:     1,
			hasHighPrio: true,
		},
		{
			name:        "generic error",
			err:         errors.New("unknown error"),
			minRecs:     1,
			hasHighPrio: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recs := handler.GetRecommendations(tt.err)
			
			if len(recs) < tt.minRecs {
				t.Errorf("Expected at least %d recommendations, got %d", tt.minRecs, len(recs))
			}
			
			if tt.hasHighPrio {
				hasHigh := false
				for _, rec := range recs {
					if rec.Priority == PriorityHigh || rec.Priority == PriorityCritical {
						hasHigh = true
						break
					}
				}
				if !hasHigh {
					t.Error("Expected at least one high priority recommendation")
				}
			}
		})
	}
}

func TestAttemptAutoRecovery(t *testing.T) {
	handler := NewFallbackHandler()
	ctx := context.Background()
	config := &RecoveryConfig{
		MaxAttempts: 3,
		Timeout:     1 * time.Second,
	}

	tests := []struct {
		name         string
		err          error
		expectMethod string
		expectSuccess bool
	}{
		{
			name:          "no error",
			err:           nil,
			expectMethod:  "no-op",
			expectSuccess: true,
		},
		{
			name:          "no such host",
			err:           errors.New("no such host"),
			expectMethod:  "dns-retry",
			expectSuccess: false,
		},
		{
			name:          "connection timeout",
			err:           errors.New("connection timeout"),
			expectMethod:  "connection-retry",
			expectSuccess: false,
		},
		{
			name:          "unknown error",
			err:           errors.New("unknown error"),
			expectMethod:  "none",
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handler.AttemptAutoRecovery(ctx, tt.err, config)
			if err != nil {
				t.Fatalf("AttemptAutoRecovery failed: %v", err)
			}
			
			if result.Method != tt.expectMethod {
				t.Errorf("Expected method %s, got %s", tt.expectMethod, result.Method)
			}
			
			if result.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, result.Success)
			}
		})
	}
}

func TestSecurityDecisionLogging(t *testing.T) {
	handler := NewFallbackHandler()
	
	decision := SecurityDecision{
		Timestamp: time.Now(),
		Type:      DecisionUseInsecure,
		Registry:  "test-registry",
		Operation: "push",
		Reason:    "testing",
		Approved:  true,
		Impact: SecurityImpact{
			Level:       ImpactHigh,
			Description: "Test decision",
		},
	}
	
	err := handler.LogSecurityDecision(decision)
	if err != nil {
		t.Fatalf("LogSecurityDecision failed: %v", err)
	}
}

func TestFallbackStrategyTypes(t *testing.T) {
	// Test that all fallback types are properly defined
	types := []FallbackType{
		FallbackNone,
		FallbackInsecure,
		FallbackAlternateTrust,
		FallbackManualTrust,
		FallbackRetry,
	}
	
	for i, typ := range types {
		if int(typ) != i {
			t.Errorf("FallbackType %v should have value %d, got %d", typ, i, int(typ))
		}
	}
}

func TestImpactLevels(t *testing.T) {
	// Test that impact levels are properly ordered
	levels := []ImpactLevel{
		ImpactMinimal,
		ImpactModerate,
		ImpactHigh,
		ImpactCritical,
	}
	
	for i, level := range levels {
		if int(level) != i {
			t.Errorf("ImpactLevel %v should have value %d, got %d", level, i, int(level))
		}
	}
}

func TestRecommendationPriorities(t *testing.T) {
	// Test that recommendation priorities are properly ordered
	priorities := []RecommendationPriority{
		PriorityLow,
		PriorityMedium,
		PriorityHigh,
		PriorityCritical,
	}
	
	for i, priority := range priorities {
		if int(priority) != i {
			t.Errorf("RecommendationPriority %v should have value %d, got %d", priority, i, int(priority))
		}
	}
}