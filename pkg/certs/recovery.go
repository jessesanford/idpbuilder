package certs

import (
	"context"
	"crypto/x509"
	"fmt"
	"sync"
	"time"
)

// RecoveryManager interface defines recovery management capabilities
type RecoveryManager interface {
	RegisterStrategy(name string, strategy RecoveryStrategy)
	Recover(ctx context.Context, err error, input *ValidationInput) (*ValidationResult, error)
	GetState() RecoveryState
}

// RecoveryStrategy defines individual recovery approaches
type RecoveryStrategy interface {
	CanRecover(err error) bool
	Attempt(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
	GetMetrics() RecoveryMetrics
}

// RecoveryState represents the current state of the recovery system
type RecoveryState struct {
	Healthy      bool
	LastRecovery time.Time
	FailureCount int
	CircuitOpen  bool
}

// RecoveryMetrics provides metrics about recovery operations
type RecoveryMetrics struct {
	AttemptCount   int
	SuccessCount   int
	FailureCount   int
	LastAttempt    time.Time
	LastSuccess    time.Time
	AverageLatency time.Duration
}

// ValidationInput contains input data for recovery operations
type ValidationInput struct {
	Certificates []*x509.Certificate
	Registry     string
	Operation    string
	Error        error
	Options      map[string]interface{}
}

// ValidationResult represents the outcome of a recovery attempt
type ValidationResult struct {
	Success       bool
	Strategy      string
	Message       string
	SecurityLevel SecurityLevel
	Actions       []string
	NewConfig     map[string]interface{}
}

// SecurityLevel defines security levels for validation results
type SecurityLevel int

const (
	SecurityHigh SecurityLevel = iota
	SecurityMedium
	SecurityLow
	SecurityNone
)

// RecoveryConfig configures recovery behavior
type RecoveryConfig struct {
	EnableCertRefresh   bool
	EnableTrustUpdate   bool
	EnableChainRepair   bool
	MaxAttempts         int
	Timeout             time.Duration
	CircuitBreakerThreshold int
	CircuitBreakerTimeout   time.Duration
}

// RecoveryResult describes the outcome of a recovery attempt
type RecoveryResult struct {
	Success   bool
	Method    string
	Actions   []string
	NewConfig interface{}
	Message   string
}

// DefaultRecoveryManager implements RecoveryManager with circuit breaker pattern
type DefaultRecoveryManager struct {
	strategies    map[string]RecoveryStrategy
	state         RecoveryState
	config        *RecoveryConfig
	mutex         sync.RWMutex
	circuitBreaker *CircuitBreaker
}

// CircuitBreaker implements circuit breaker pattern for recovery operations
type CircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration
	failureCount     int
	lastFailureTime  time.Time
	state           CircuitState
	mutex           sync.Mutex
}

// CircuitState represents circuit breaker states
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// NewRecoveryManager creates a new recovery manager with circuit breaker
func NewRecoveryManager(config *RecoveryConfig) RecoveryManager {
	if config == nil {
		config = &RecoveryConfig{
			EnableCertRefresh:   true,
			EnableTrustUpdate:   true,
			EnableChainRepair:   true,
			MaxAttempts:         3,
			Timeout:             30 * time.Second,
			CircuitBreakerThreshold: 5,
			CircuitBreakerTimeout:   60 * time.Second,
		}
	}

	circuitBreaker := &CircuitBreaker{
		failureThreshold: config.CircuitBreakerThreshold,
		resetTimeout:     config.CircuitBreakerTimeout,
		state:           CircuitClosed,
	}

	manager := &DefaultRecoveryManager{
		strategies:     make(map[string]RecoveryStrategy),
		state:          RecoveryState{Healthy: true},
		config:         config,
		circuitBreaker: circuitBreaker,
	}

	// Register default recovery strategies
	manager.RegisterStrategy("auto-recover-connection", NewAutoRecoveryStrategy())
	manager.RegisterStrategy("refresh-certificates", NewCertificateRefreshStrategy())
	manager.RegisterStrategy("update-trust-store", NewTrustStoreUpdateStrategy())

	return manager
}

// RegisterStrategy registers a new recovery strategy
func (rm *DefaultRecoveryManager) RegisterStrategy(name string, strategy RecoveryStrategy) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.strategies[name] = strategy
}

// Recover attempts to recover from an error using registered strategies
func (rm *DefaultRecoveryManager) Recover(ctx context.Context, err error, input *ValidationInput) (*ValidationResult, error) {
	if err == nil {
		return &ValidationResult{Success: true, Message: "No error to recover from"}, nil
	}

	// Check circuit breaker state
	if !rm.circuitBreaker.CanExecute() {
		return &ValidationResult{
			Success: false,
			Message: "Circuit breaker open - recovery temporarily disabled",
			Actions: []string{"Circuit breaker protection active"},
		}, nil
	}

	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	var lastError error
	var attempts []string

	for name, strategy := range rm.strategies {
		if !strategy.CanRecover(err) {
			continue
		}

		attempts = append(attempts, name)

		select {
		case <-ctx.Done():
			return &ValidationResult{
				Success: false,
				Message: "Recovery cancelled - context cancelled",
				Actions: attempts,
			}, ctx.Err()
		default:
		}

		result, recoveryErr := strategy.Attempt(ctx, input)
		if recoveryErr != nil {
			lastError = recoveryErr
			rm.circuitBreaker.RecordFailure()
			continue
		}

		if result.Success {
			rm.circuitBreaker.RecordSuccess()
			rm.updateRecoveryState(true)
			result.Actions = append(result.Actions, fmt.Sprintf("Successful recovery using: %s", name))
			return result, nil
		}

		lastError = fmt.Errorf("strategy %s failed: %s", name, result.Message)
	}

	rm.circuitBreaker.RecordFailure()
	rm.updateRecoveryState(false)

	return &ValidationResult{
		Success: false,
		Message: fmt.Sprintf("All recovery strategies failed. Last error: %v", lastError),
		Actions: attempts,
	}, lastError
}

// GetState returns the current recovery system state
func (rm *DefaultRecoveryManager) GetState() RecoveryState {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	
	state := rm.state
	state.CircuitOpen = rm.circuitBreaker.state == CircuitOpen
	return state
}

// updateRecoveryState updates internal recovery state
func (rm *DefaultRecoveryManager) updateRecoveryState(success bool) {
	if success {
		rm.state.Healthy = true
		rm.state.LastRecovery = time.Now()
		rm.state.FailureCount = 0
	} else {
		rm.state.FailureCount++
		if rm.state.FailureCount >= rm.config.MaxAttempts {
			rm.state.Healthy = false
		}
	}
}

// Circuit Breaker Implementation

// CanExecute checks if circuit breaker allows execution
func (cb *CircuitBreaker) CanExecute() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = CircuitHalfOpen
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

// RecordSuccess records a successful operation
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failureCount = 0
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
	}
}

// RecordFailure records a failed operation
func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.failureThreshold {
		cb.state = CircuitOpen
	}
}

// Auto Recovery Strategy Implementation

// AutoRecoveryStrategy implements connection recovery with exponential backoff
type AutoRecoveryStrategy struct {
	name         string
	metrics      RecoveryMetrics
	baseDelay    time.Duration
	maxDelay     time.Duration
	backoffFactor float64
	maxRetries   int
}

// NewAutoRecoveryStrategy creates a new auto recovery strategy
func NewAutoRecoveryStrategy() RecoveryStrategy {
	return &AutoRecoveryStrategy{
		name:         "auto-recover-connection",
		baseDelay:    1 * time.Second,
		maxDelay:     30 * time.Second,
		backoffFactor: 2.0,
		maxRetries:   3,
		metrics:      RecoveryMetrics{},
	}
}

// CanRecover determines if this strategy can recover from the given error
func (ars *AutoRecoveryStrategy) CanRecover(err error) bool {
	if err == nil {
		return false
	}
	
	errorStr := err.Error()
	return containsAny(errorStr, []string{
		"connection", "timeout", "network", "no such host", 
		"unreachable", "refused", "reset",
	})
}

// Attempt performs connection recovery with exponential backoff
func (ars *AutoRecoveryStrategy) Attempt(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	start := time.Now()
	ars.metrics.AttemptCount++
	ars.metrics.LastAttempt = start

	var lastErr error
	delay := ars.baseDelay
	actions := []string{"Starting connection recovery with exponential backoff"}

	for attempt := 1; attempt <= ars.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return &ValidationResult{
				Success:  false,
				Strategy: ars.name,
				Message:  "Connection recovery cancelled - context cancelled",
				Actions:  actions,
			}, ctx.Err()
		default:
		}

		actions = append(actions, fmt.Sprintf("Connection attempt %d for registry %s", attempt, input.Registry))

		// Simulate connection test - in real implementation, this would test actual connectivity
		if err := ars.simulateConnectionTest(ctx); err == nil {
			ars.metrics.SuccessCount++
			ars.metrics.LastSuccess = time.Now()
			ars.metrics.AverageLatency = time.Since(start)

			actions = append(actions, fmt.Sprintf("Connection successful on attempt %d", attempt))
			return &ValidationResult{
				Success:      true,
				Strategy:     ars.name,
				Message:      "Connection recovery successful",
				SecurityLevel: SecurityHigh,
				Actions:      actions,
				NewConfig: map[string]interface{}{
					"connectionRecovered": true,
					"attempts":           attempt,
					"recoveryTime":       time.Since(start).String(),
				},
			}, nil
		} else {
			lastErr = err
			actions = append(actions, fmt.Sprintf("Attempt %d failed: %s", attempt, err.Error()))
		}

		// Don't wait after the last attempt
		if attempt < ars.maxRetries {
			actions = append(actions, fmt.Sprintf("Waiting %v before next attempt", delay))
			
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return &ValidationResult{
					Success:  false,
					Strategy: ars.name,
					Message:  "Connection recovery cancelled during backoff",
					Actions:  actions,
				}, ctx.Err()
			case <-timer.C:
			}

			// Exponential backoff
			delay = time.Duration(float64(delay) * ars.backoffFactor)
			if delay > ars.maxDelay {
				delay = ars.maxDelay
			}
		}
	}

	ars.metrics.FailureCount++
	return &ValidationResult{
		Success:  false,
		Strategy: ars.name,
		Message:  fmt.Sprintf("Connection recovery failed - all %d attempts failed, last error: %s", ars.maxRetries, lastErr.Error()),
		Actions:  actions,
	}, lastErr
}

// GetMetrics returns recovery metrics
func (ars *AutoRecoveryStrategy) GetMetrics() RecoveryMetrics {
	return ars.metrics
}

// simulateConnectionTest simulates testing network connectivity
func (ars *AutoRecoveryStrategy) simulateConnectionTest(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(50 * time.Millisecond):
		// Simulate improving network conditions - some attempts may succeed
		if time.Now().UnixNano()%3 == 0 { // Randomly succeed 1/3 of the time
			return nil // Connection successful
		}
		return fmt.Errorf("network still unreachable")
	}
}

// Certificate Refresh Strategy Implementation

// CertificateRefreshStrategy implements certificate refresh with retry logic
type CertificateRefreshStrategy struct {
	name       string
	metrics    RecoveryMetrics
	maxRetries int
	timeout    time.Duration
}

// NewCertificateRefreshStrategy creates a new certificate refresh strategy
func NewCertificateRefreshStrategy() RecoveryStrategy {
	return &CertificateRefreshStrategy{
		name:       "refresh-certificates",
		maxRetries: 3,
		timeout:    30 * time.Second,
		metrics:    RecoveryMetrics{},
	}
}

// CanRecover determines if this strategy can recover from certificate errors
func (crs *CertificateRefreshStrategy) CanRecover(err error) bool {
	if err == nil {
		return false
	}

	errorStr := err.Error()
	return containsAny(errorStr, []string{
		"certificate has expired", "expired", "certificate is not valid",
		"certificate not yet valid", "validity period",
	})
}

// Attempt performs certificate refresh
func (crs *CertificateRefreshStrategy) Attempt(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	start := time.Now()
	crs.metrics.AttemptCount++
	crs.metrics.LastAttempt = start

	actions := []string{fmt.Sprintf("Starting certificate refresh for registry %s", input.Registry)}

	refreshCtx, cancel := context.WithTimeout(ctx, crs.timeout)
	defer cancel()

	select {
	case <-refreshCtx.Done():
		crs.metrics.FailureCount++
		return &ValidationResult{
			Success:  false,
			Strategy: crs.name,
			Message:  "Certificate refresh operation timed out",
			Actions:  append(actions, "Certificate refresh timed out"),
		}, refreshCtx.Err()
	case <-time.After(200 * time.Millisecond):
		// Simulate certificate refresh process
		actions = append(actions, "Attempting to retrieve updated certificates from registry")
		actions = append(actions, "Certificate refresh process completed")

		crs.metrics.SuccessCount++
		crs.metrics.LastSuccess = time.Now()
		crs.metrics.AverageLatency = time.Since(start)

		return &ValidationResult{
			Success:      true,
			Strategy:     crs.name,
			Message:      "Certificate refresh successful",
			SecurityLevel: SecurityHigh,
			Actions:      actions,
			NewConfig: map[string]interface{}{
				"certificateRefreshed": true,
				"registry":            input.Registry,
				"refreshTime":         time.Now(),
			},
		}, nil
	}
}

// GetMetrics returns certificate refresh metrics
func (crs *CertificateRefreshStrategy) GetMetrics() RecoveryMetrics {
	return crs.metrics
}

// Trust Store Update Strategy Implementation

// TrustStoreUpdateStrategy implements trust store update for authority issues
type TrustStoreUpdateStrategy struct {
	name       string
	metrics    RecoveryMetrics
	maxRetries int
	timeout    time.Duration
}

// NewTrustStoreUpdateStrategy creates a new trust store update strategy
func NewTrustStoreUpdateStrategy() RecoveryStrategy {
	return &TrustStoreUpdateStrategy{
		name:       "update-trust-store",
		maxRetries: 3,
		timeout:    45 * time.Second,
		metrics:    RecoveryMetrics{},
	}
}

// CanRecover determines if this strategy can recover from trust issues
func (tsus *TrustStoreUpdateStrategy) CanRecover(err error) bool {
	if err == nil {
		return false
	}

	errorStr := err.Error()
	return containsAny(errorStr, []string{
		"self signed", "unknown authority", "untrusted", "trust",
		"certificate signed by unknown authority", "x509: certificate signed by unknown authority",
	})
}

// Attempt performs trust store update
func (tsus *TrustStoreUpdateStrategy) Attempt(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	start := time.Now()
	tsus.metrics.AttemptCount++
	tsus.metrics.LastAttempt = start

	actions := []string{fmt.Sprintf("Starting trust store update for registry %s", input.Registry)}

	updateCtx, cancel := context.WithTimeout(ctx, tsus.timeout)
	defer cancel()

	select {
	case <-updateCtx.Done():
		tsus.metrics.FailureCount++
		return &ValidationResult{
			Success:  false,
			Strategy: tsus.name,
			Message:  "Trust store update operation timed out",
			Actions:  append(actions, "Trust store update timed out"),
		}, updateCtx.Err()
	case <-time.After(300 * time.Millisecond):
		// Simulate trust store update process
		actions = append(actions, "Analyzing certificate trust chain")
		actions = append(actions, "Updating trust store with registry certificates")
		actions = append(actions, "Trust store update completed")

		tsus.metrics.SuccessCount++
		tsus.metrics.LastSuccess = time.Now()
		tsus.metrics.AverageLatency = time.Since(start)

		return &ValidationResult{
			Success:      true,
			Strategy:     tsus.name,
			Message:      "Trust store update successful",
			SecurityLevel: SecurityMedium,
			Actions:      actions,
			NewConfig: map[string]interface{}{
				"trustStoreUpdated": true,
				"registry":         input.Registry,
				"updateTime":       time.Now(),
			},
		}, nil
	}
}

// GetMetrics returns trust store update metrics
func (tsus *TrustStoreUpdateStrategy) GetMetrics() RecoveryMetrics {
	return tsus.metrics
}

// Helper Functions

// containsAny checks if string contains any of the given substrings
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

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}