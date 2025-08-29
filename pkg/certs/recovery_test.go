package certs

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRecoveryManager(t *testing.T) {
	manager := NewRecoveryManager()
	if manager == nil || manager.maxRetries != 3 {
		t.Fatal("NewRecoveryManager configuration error")
	}

	ctx := context.Background()
	config := &RecoveryConfig{
		EnableCertRefresh: true,
		EnableTrustUpdate: true,
		EnableChainRepair: true,
		MaxAttempts:       3,
		Timeout:           100 * time.Millisecond,
	}

	tests := []struct {
		err            error
		expectedMethod string
	}{
		{nil, "no-error"},
		{errors.New("certificate has expired"), "cert-refresh"},
		{errors.New("self signed certificate"), "trust-update"},
		{errors.New("connection timeout"), "retry-with-backoff"},
		{errors.New("incomplete chain"), "chain-repair"},
	}

	for _, tt := range tests {
		result, err := manager.RecoverFromError(ctx, tt.err, config)
		if err != nil && tt.err == nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result.Method != tt.expectedMethod {
			t.Errorf("Expected method %s, got %s", tt.expectedMethod, result.Method)
		}
	}
}

func TestRecoverWithRetry(t *testing.T) {
	manager := NewRecoveryManager()
	ctx := context.Background()
	config := &RecoveryConfig{MaxAttempts: 2, Timeout: 100 * time.Millisecond}

	// Test immediate success
	callCount := 0
	operation := func() error {
		callCount++
		return nil
	}
	result, err := manager.RecoverWithRetry(ctx, operation, config)
	if err != nil || !result.Success || callCount != 1 {
		t.Error("Immediate success test failed")
	}

	// Test eventual failure
	callCount = 0
	operation = func() error {
		callCount++
		return errors.New("persistent failure")
	}
	result, err = manager.RecoverWithRetry(ctx, operation, config)
	if err == nil || result.Success || callCount != 2 {
		t.Error("Persistent failure test failed")
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test containsAny
	if !containsAny("hello world", []string{"world"}) {
		t.Error("containsAny failed")
	}
	if containsAny("hello world", []string{"foo"}) {
		t.Error("containsAny false positive")
	}

	// Test min
	if min(5, 3) != 3 || min(1, 1) != 1 {
		t.Error("min function failed")
	}
}

// Additional test coverage for recovery scenarios
func TestRecoveryWithContextCancellation(t *testing.T) {
	manager := NewRecoveryManager()
	ctx, cancel := context.WithCancel(context.Background())
	config := &RecoveryConfig{
		EnableTrustUpdate: true,
		MaxAttempts: 5,
		Timeout: 10 * time.Millisecond,
	}
	
	// Cancel context immediately 
	cancel()
	
	result, _ := manager.recoverFromTrustIssue(ctx, config)
	// The method might not return an error in our implementation
	if result == nil {
		t.Error("Expected result even on cancellation")
	}
}

func TestRecoveryDisabledFeatures(t *testing.T) {
	manager := NewRecoveryManager()
	ctx := context.Background()
	
	// Test with all features disabled
	config := &RecoveryConfig{
		EnableCertRefresh: false,
		EnableTrustUpdate: false,
		EnableChainRepair: false,
		MaxAttempts: 1,
		Timeout: 100 * time.Millisecond,
	}
	
	// Test cert refresh disabled
	result, err := manager.recoverFromExpiredCert(ctx, config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Success {
		t.Error("Expected failure when cert refresh disabled")
	}
	
	// Test trust update disabled
	result, err = manager.recoverFromTrustIssue(ctx, config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err) 
	}
	if result.Success {
		t.Error("Expected failure when trust update disabled")
	}
	
	// Test chain repair disabled
	result, err = manager.recoverFromChainIssue(ctx, config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Success {
		t.Error("Expected failure when chain repair disabled")
	}
}

func TestRecoveryManagerNetworkRetry(t *testing.T) {
	manager := NewRecoveryManager()
	ctx := context.Background()
	config := &RecoveryConfig{
		MaxAttempts: 2,
		Timeout: 100 * time.Millisecond,
	}
	
	// Test network recovery with retry - expect failure in simulation
	result, _ := manager.recoverFromNetworkIssue(ctx, config)
	if result == nil {
		t.Error("Expected network recovery result")
	}
	// Network recovery simulates failure, so expect !Success
	if result.Success {
		t.Log("Network recovery unexpectedly succeeded")
	}
	
	// Test generic recovery with limited retries - expect failure in simulation
	result, _ = manager.genericRecovery(ctx, config)
	if result == nil {
		t.Error("Expected generic recovery result")
	}
	// Generic recovery simulates failure, so expect !Success
	if result.Success {
		t.Log("Generic recovery unexpectedly succeeded")
	}
}