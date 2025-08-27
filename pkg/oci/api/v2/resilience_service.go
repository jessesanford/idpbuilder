package v2

import (
	"context"
	"time"
)

// ResilienceService defines the contract for resilience patterns
type ResilienceService interface {
	// CreateCircuitBreaker creates a new circuit breaker
	CreateCircuitBreaker(name string, config CircuitBreakerConfig) (CircuitBreaker, error)

	// GetCircuitBreaker retrieves an existing circuit breaker
	GetCircuitBreaker(name string) (CircuitBreaker, error)

	// CreateRetryPolicy creates a new retry policy
	CreateRetryPolicy(name string, config RetryConfig) (RetryPolicy, error)

	// ExecuteWithResilience wraps an operation with resilience patterns
	ExecuteWithResilience(ctx context.Context, op Operation) error

	// RegisterHealthCheck registers a health check
	RegisterHealthCheck(name string, check HealthCheck) error

	// GetHealthStatus returns current health status
	GetHealthStatus() HealthStatus

	// ResetCircuitBreaker manually resets a circuit breaker
	ResetCircuitBreaker(name string) error

	// GetStatistics returns resilience statistics
	GetStatistics() ResilienceStats
}

// CircuitBreaker interface for circuit breaker pattern
type CircuitBreaker interface {
	Execute(ctx context.Context, op Operation) error
	GetState() CircuitState
	Reset() error
	GetStatistics() CircuitStats
}

// CircuitBreakerConfig defines circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold  int
	SuccessThreshold  int
	Timeout           time.Duration
	ResetTimeout      time.Duration
	HalfOpenRequests  int
}

// CircuitState represents circuit breaker states
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "closed"
	CircuitStateOpen     CircuitState = "open"
	CircuitStateHalfOpen CircuitState = "half-open"
)

// RetryPolicy interface for retry strategies
type RetryPolicy interface {
	Execute(ctx context.Context, op Operation) error
	GetConfig() RetryConfig
}

// RetryConfig defines retry policy configuration
type RetryConfig struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	Multiplier      float64
	Jitter          float64
	RetryableErrors []string
}

// Operation represents a resilient operation
type Operation func(ctx context.Context) error

// HealthCheck represents a health check function
type HealthCheck func(ctx context.Context) error

// HealthStatus represents overall health status
type HealthStatus struct {
	Status     string
	Checks     map[string]CheckResult
	LastUpdate time.Time
}

// CheckResult represents individual check result
type CheckResult struct {
	Status  string
	Message string
	Error   error
}

// ResilienceStats provides resilience statistics
type ResilienceStats struct {
	CircuitBreakers map[string]CircuitStats
	RetryPolicies   map[string]RetryStats
}

// CircuitStats provides circuit breaker statistics
type CircuitStats struct {
	State           CircuitState
	FailureCount    int64
	SuccessCount    int64
	LastStateChange time.Time
}

// RetryStats provides retry policy statistics
type RetryStats struct {
	TotalAttempts     int64
	SuccessfulRetries int64
	FailedRetries     int64
}