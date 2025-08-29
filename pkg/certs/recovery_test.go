package certs

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewRecoveryManager(t *testing.T) {
	manager := NewRecoveryManager()
	if manager == nil {
		t.Fatal("NewRecoveryManager returned nil")
	}
	
	if manager.maxRetries != 3 {
		t.Errorf("Expected maxRetries 3, got %d", manager.maxRetries)
	}
	
	if manager.baseDelay != 1*time.Second {
		t.Errorf("Expected baseDelay 1s, got %v", manager.baseDelay)
	}
}

func TestRecoverFromError(t *testing.T) {
	manager := NewRecoveryManager()
	ctx := context.Background()
	config := &RecoveryConfig{
		EnableCertRefresh: true,
		EnableTrustUpdate: true,
		EnableChainRepair: true,
		MaxAttempts:       3,
		Timeout:           1 * time.Second,
	}

	tests := []struct {
		name           string
		err            error
		expectedMethod string
		expectSuccess  bool
	}{
		{
			name:           "no error",
			err:            nil,
			expectedMethod: "no-error",
			expectSuccess:  true,
		},
		{
			name:           "expired certificate",
			err:            errors.New("certificate has expired"),
			expectedMethod: "cert-refresh",
			expectSuccess:  false,
		},
		{
			name:           "self signed certificate",
			err:            errors.New("self signed certificate"),
			expectedMethod: "trust-update",
			expectSuccess:  false,
		},
		{
			name:           "network error",
			err:            errors.New("connection timeout"),
			expectedMethod: "retry-with-backoff",
			expectSuccess:  false,
		},
		{
			name:           "chain error",
			err:            errors.New("incomplete chain"),
			expectedMethod: "chain-repair",
			expectSuccess:  false,
		},
		{
			name:           "generic error",
			err:            errors.New("unknown error"),
			expectedMethod: "retry-with-backoff",
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.RecoverFromError(ctx, tt.err, config)
			if err != nil && tt.expectSuccess {
				t.Fatalf("RecoverFromError failed: %v", err)
			}
			
			if result.Method != tt.expectedMethod {
				t.Errorf("Expected method %s, got %s", tt.expectedMethod, result.Method)
			}
			
			if result.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, result.Success)
			}
			
			if len(result.Actions) == 0 && tt.err != nil {
				t.Error("Expected at least one action to be recorded")
			}
		})
	}
}

func TestRecoverWithRetry(t *testing.T) {
	manager := NewRecoveryManager()
	ctx := context.Background()
	config := &RecoveryConfig{
		MaxAttempts: 3,
		Timeout:     2 * time.Second,
	}

	t.Run("operation succeeds on first try", func(t *testing.T) {
		callCount := 0
		operation := func() error {
			callCount++
			return nil // Success immediately
		}
		
		result, err := manager.RecoverWithRetry(ctx, operation, config)
		if err != nil {
			t.Fatalf("RecoverWithRetry failed: %v", err)
		}
		
		if !result.Success {
			t.Error("Expected success")
		}
		
		if callCount != 1 {
			t.Errorf("Expected 1 call, got %d", callCount)
		}
	})

	t.Run("operation succeeds on second try", func(t *testing.T) {
		callCount := 0
		operation := func() error {
			callCount++
			if callCount == 1 {
				return errors.New("temporary failure")
			}
			return nil // Success on second try
		}
		
		result, err := manager.RecoverWithRetry(ctx, operation, config)
		if err != nil {
			t.Fatalf("RecoverWithRetry failed: %v", err)
		}
		
		if !result.Success {
			t.Error("Expected success")
		}
		
		if callCount != 2 {
			t.Errorf("Expected 2 calls, got %d", callCount)
		}
	})

	t.Run("operation fails all attempts", func(t *testing.T) {
		callCount := 0
		operation := func() error {
			callCount++
			return errors.New("persistent failure")
		}
		
		result, err := manager.RecoverWithRetry(ctx, operation, config)
		if err == nil {
			t.Error("Expected error for persistent failure")
		}
		
		if result.Success {
			t.Error("Expected failure")
		}
		
		if callCount != config.MaxAttempts {
			t.Errorf("Expected %d calls, got %d", config.MaxAttempts, callCount)
		}
		
		if len(result.Actions) == 0 {
			t.Error("Expected actions to be recorded")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		
		operation := func() error {
			cancel() // Cancel context during operation
			return errors.New("failure")
		}
		
		result, err := manager.RecoverWithRetry(ctx, operation, config)
		if err == nil {
			t.Error("Expected context cancellation error")
		}
		
		if result.Success {
			t.Error("Expected failure due to cancellation")
		}
	})
}

func TestRecoveryConfigDisabling(t *testing.T) {
	manager := NewRecoveryManager()
	ctx := context.Background()
	
	t.Run("cert refresh disabled", func(t *testing.T) {
		config := &RecoveryConfig{
			EnableCertRefresh: false,
			MaxAttempts:       3,
			Timeout:           1 * time.Second,
		}
		
		err := errors.New("certificate has expired")
		result, recErr := manager.RecoverFromError(ctx, err, config)
		if recErr != nil {
			t.Fatalf("RecoverFromError failed: %v", recErr)
		}
		
		if result.Success {
			t.Error("Expected failure when cert refresh disabled")
		}
		
		if result.FailureReason == "" {
			t.Error("Expected failure reason to be set")
		}
	})

	t.Run("trust update disabled", func(t *testing.T) {
		config := &RecoveryConfig{
			EnableTrustUpdate: false,
			MaxAttempts:       3,
			Timeout:           1 * time.Second,
		}
		
		err := errors.New("self signed certificate")
		result, recErr := manager.RecoverFromError(ctx, err, config)
		if recErr != nil {
			t.Fatalf("RecoverFromError failed: %v", recErr)
		}
		
		if result.Success {
			t.Error("Expected failure when trust update disabled")
		}
	})

	t.Run("chain repair disabled", func(t *testing.T) {
		config := &RecoveryConfig{
			EnableChainRepair: false,
			MaxAttempts:       3,
			Timeout:           1 * time.Second,
		}
		
		err := errors.New("incomplete chain")
		result, recErr := manager.RecoverFromError(ctx, err, config)
		if recErr != nil {
			t.Fatalf("RecoverFromError failed: %v", recErr)
		}
		
		if result.Success {
			t.Error("Expected failure when chain repair disabled")
		}
	})
}

func TestRecoveryTimeout(t *testing.T) {
	manager := NewRecoveryManager()
	ctx := context.Background()
	
	// Very short timeout to test timeout handling
	config := &RecoveryConfig{
		EnableCertRefresh: true,
		MaxAttempts:       3,
		Timeout:           1 * time.Nanosecond, // Immediate timeout
	}
	
	err := errors.New("certificate has expired")
	result, recErr := manager.RecoverFromError(ctx, err, config)
	if recErr != nil {
		t.Fatalf("RecoverFromError failed: %v", recErr)
	}
	
	if result.Success {
		t.Error("Expected failure due to timeout")
	}
	
	if result.FailureReason == "" {
		t.Error("Expected timeout failure reason")
	}
}

func TestContainsAnyHelper(t *testing.T) {
	tests := []struct {
		str        string
		substrings []string
		expected   bool
	}{
		{"hello world", []string{"hello"}, true},
		{"hello world", []string{"world"}, true},
		{"hello world", []string{"foo", "world"}, true},
		{"hello world", []string{"foo", "bar"}, false},
		{"", []string{"hello"}, false},
		{"hello", []string{""}, true}, // Empty substring always matches
		{"hello", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			result := containsAny(tt.str, tt.substrings)
			if result != tt.expected {
				t.Errorf("containsAny(%q, %v) = %v, expected %v", 
					tt.str, tt.substrings, result, tt.expected)
			}
		})
	}
}

func TestMinHelper(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 1},
		{2, 1, 1},
		{5, 5, 5},
		{0, 10, 0},
		{-1, 1, -1},
	}

	for _, tt := range tests {
		result := min(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("min(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestRecoveryResultFields(t *testing.T) {
	result := &RecoveryResult{
		Success:       true,
		Method:        "test-method",
		Actions:       []string{"action1", "action2"},
		NewConfig:     "test-config",
		FailureReason: "",
	}
	
	if !result.Success {
		t.Error("Expected Success to be true")
	}
	
	if result.Method != "test-method" {
		t.Errorf("Expected Method 'test-method', got %s", result.Method)
	}
	
	if len(result.Actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(result.Actions))
	}
	
	if result.NewConfig != "test-config" {
		t.Errorf("Expected NewConfig 'test-config', got %v", result.NewConfig)
	}
}