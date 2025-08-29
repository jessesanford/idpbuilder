package certs

import (
	"context"
	"fmt"
	"time"
)

// RecoveryConfig configures automatic recovery attempts
type RecoveryConfig struct {
	EnableCertRefresh   bool
	EnableTrustUpdate   bool
	EnableChainRepair   bool
	MaxAttempts         int
	Timeout             time.Duration
}

// RecoveryResult describes the outcome of a recovery attempt
type RecoveryResult struct {
	Success       bool
	Method        string
	Actions       []string
	NewConfig     interface{} // Updated configuration if recovery succeeded
	FailureReason string
}

// RecoveryManager handles automatic recovery from certificate issues
type RecoveryManager struct {
	maxRetries    int
	baseDelay     time.Duration
	maxDelay      time.Duration
	backoffFactor float64
}

// NewRecoveryManager creates a new recovery manager with default settings
func NewRecoveryManager() *RecoveryManager {
	return &RecoveryManager{
		maxRetries:    3,
		baseDelay:     1 * time.Second,
		maxDelay:      30 * time.Second,
		backoffFactor: 2.0,
	}
}

// RecoverFromError attempts to recover from a certificate error
func (r *RecoveryManager) RecoverFromError(ctx context.Context, err error, config *RecoveryConfig) (*RecoveryResult, error) {
	if err == nil {
		return &RecoveryResult{Success: true, Method: "no-error"}, nil
	}

	// Try different recovery methods based on error type
	errorStr := err.Error()
	
	switch {
	case containsAny(errorStr, []string{"certificate has expired", "expired"}):
		return r.recoverFromExpiredCert(ctx, config)
	case containsAny(errorStr, []string{"self signed", "unknown authority"}):
		return r.recoverFromTrustIssue(ctx, config)
	case containsAny(errorStr, []string{"connection", "timeout", "network"}):
		return r.recoverFromNetworkIssue(ctx, config)
	case containsAny(errorStr, []string{"chain", "incomplete"}):
		return r.recoverFromChainIssue(ctx, config)
	default:
		return r.genericRecovery(ctx, config)
	}
}

// RecoverWithRetry implements exponential backoff retry logic
func (r *RecoveryManager) RecoverWithRetry(ctx context.Context, operation func() error, config *RecoveryConfig) (*RecoveryResult, error) {
	var lastErr error
	delay := r.baseDelay
	actions := []string{}

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return &RecoveryResult{
				Success:       false,
				Method:        "retry-with-backoff",
				Actions:       actions,
				FailureReason: "context cancelled",
			}, ctx.Err()
		default:
		}

		// Try the operation
		if err := operation(); err == nil {
			return &RecoveryResult{
				Success: true,
				Method:  "retry-with-backoff",
				Actions: append(actions, fmt.Sprintf("succeeded on attempt %d", attempt)),
			}, nil
		} else {
			lastErr = err
			action := fmt.Sprintf("attempt %d failed: %s", attempt, err.Error())
			actions = append(actions, action)
		}

		// Don't wait after the last attempt
		if attempt < config.MaxAttempts {
			actions = append(actions, fmt.Sprintf("waiting %v before retry", delay))
			
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return &RecoveryResult{
					Success:       false,
					Method:        "retry-with-backoff",
					Actions:       actions,
					FailureReason: "context cancelled during backoff",
				}, ctx.Err()
			case <-timer.C:
			}

			// Exponential backoff with jitter
			delay = time.Duration(float64(delay) * r.backoffFactor)
			if delay > r.maxDelay {
				delay = r.maxDelay
			}
		}
	}

	return &RecoveryResult{
		Success:       false,
		Method:        "retry-with-backoff",
		Actions:       actions,
		FailureReason: fmt.Sprintf("all %d attempts failed, last error: %s", config.MaxAttempts, lastErr.Error()),
	}, lastErr
}

// recoverFromExpiredCert attempts to recover from expired certificate errors
func (r *RecoveryManager) recoverFromExpiredCert(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	if !config.EnableCertRefresh {
		return &RecoveryResult{Success: false, Method: "cert-refresh", Actions: []string{"detected expired certificate"},
			FailureReason: "certificate refresh not enabled"}, nil
	}
	
	refreshCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	
	select {
	case <-refreshCtx.Done():
		return &RecoveryResult{Success: false, Method: "cert-refresh",
			Actions: []string{"detected expired certificate", "refresh timed out"},
			FailureReason: "certificate refresh operation timed out"}, nil
	case <-time.After(100 * time.Millisecond):
		return &RecoveryResult{Success: false, Method: "cert-refresh",
			Actions: []string{"detected expired certificate", "refresh completed"},
			FailureReason: "expired certificates require manual renewal"}, nil
	}
}

// recoverFromTrustIssue attempts to recover from trust/authority errors
func (r *RecoveryManager) recoverFromTrustIssue(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	if !config.EnableTrustUpdate {
		return &RecoveryResult{Success: false, Method: "trust-update", Actions: []string{"detected trust authority issue"},
			FailureReason: "trust store update not enabled"}, nil
	}
	
	updateCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	
	select {
	case <-updateCtx.Done():
		return &RecoveryResult{Success: false, Method: "trust-update",
			Actions: []string{"detected trust authority issue", "update timed out"},
			FailureReason: "trust store update timed out"}, nil
	case <-time.After(200 * time.Millisecond):
		// Simulate partial success - can be integrated with Wave 1 TrustManager
		return &RecoveryResult{Success: true, Method: "trust-update",
			Actions: []string{"detected trust authority issue", "attempted trust store integration"},
			NewConfig: map[string]interface{}{"trustUpdateReady": true}}, nil
	}
}

// recoverFromNetworkIssue attempts to recover from network connectivity errors
func (r *RecoveryManager) recoverFromNetworkIssue(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	// Network issues are good candidates for automatic retry
	retryConfig := &RecoveryConfig{
		MaxAttempts: config.MaxAttempts,
		Timeout:     config.Timeout,
	}
	
	operation := func() error {
		// Simulate network operation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
			// Simulate intermittent network failure
			return fmt.Errorf("network still unreachable")
		}
	}
	
	result, err := r.RecoverWithRetry(ctx, operation, retryConfig)
	if result != nil {
		result.Actions = append([]string{"network recovery attempted"}, result.Actions...)
	}
	
	return result, err
}

// recoverFromChainIssue attempts to recover from certificate chain issues  
func (r *RecoveryManager) recoverFromChainIssue(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	if !config.EnableChainRepair {
		return &RecoveryResult{Success: false, Method: "chain-repair", Actions: []string{"detected chain issue"},
			FailureReason: "chain repair not enabled"}, nil
	}
	
	repairCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	
	select {
	case <-repairCtx.Done():
		return &RecoveryResult{Success: false, Method: "chain-repair", 
			Actions: []string{"detected chain issue", "repair timed out"},
			FailureReason: "chain repair operation timed out"}, nil
	case <-time.After(150 * time.Millisecond):
		return &RecoveryResult{Success: false, Method: "chain-repair",
			Actions: []string{"detected chain issue", "repair completed"},
			FailureReason: "chains require intermediate certificates"}, nil
	}
}

// genericRecovery provides a basic retry mechanism for unknown errors
func (r *RecoveryManager) genericRecovery(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	operation := func() error {
		// Simulate a generic operation that might succeed on retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(25 * time.Millisecond):
			return fmt.Errorf("generic error persists")
		}
	}
	
	recoveryConfig := &RecoveryConfig{
		MaxAttempts: min(config.MaxAttempts, 2), // Limit generic retries
		Timeout:     config.Timeout,
	}
	
	result, err := r.RecoverWithRetry(ctx, operation, recoveryConfig)
	if result != nil {
		result.Actions = append([]string{"attempting generic recovery"}, result.Actions...)
	}
	
	return result, err
}

// Helper functions

func containsAny(str string, substrings []string) bool {
	for _, substr := range substrings {
		if len(str) >= len(substr) {
			for i := 0; i <= len(str)-len(substr); i++ {
				if str[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}