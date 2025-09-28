package errors

import (
	"math"
	"os"
	"strconv"
	"time"
)

// RetryStrategy defines when and how to retry failed operations
type RetryStrategy interface {
	// ShouldRetry determines if an operation should be retried
	ShouldRetry(err error, attempt int) bool

	// NextDelay calculates the delay before next retry
	NextDelay(attempt int) time.Duration

	// MaxAttempts returns the maximum number of retry attempts
	MaxAttempts() int
}

// ExponentialBackoff implements exponential backoff retry strategy
type ExponentialBackoff struct {
	baseDelay      time.Duration
	maxDelay       time.Duration
	maxAttempts    int
	retryableCodes map[ErrorCode]bool
}

// NewExponentialBackoff creates a new exponential backoff strategy
func NewExponentialBackoff(base, max time.Duration, maxAttempts int) *ExponentialBackoff {
	return &ExponentialBackoff{
		baseDelay:   base,
		maxDelay:    max,
		maxAttempts: maxAttempts,
		retryableCodes: map[ErrorCode]bool{
			RegistryUnreachable: true,
			BuildTimeout:        true,
		},
	}
}

// NewExponentialBackoffFromEnv creates a backoff strategy from environment variables
func NewExponentialBackoffFromEnv() *ExponentialBackoff {
	// Default values
	baseDelay := 100 * time.Millisecond
	maxDelay := 10 * time.Second
	maxAttempts := 3

	// Override from environment
	if val := os.Getenv("RETRY_BASE_DELAY_MS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			baseDelay = time.Duration(parsed) * time.Millisecond
		}
	}

	if val := os.Getenv("RETRY_MAX_DELAY_MS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			maxDelay = time.Duration(parsed) * time.Millisecond
		}
	}

	if val := os.Getenv("RETRY_MAX_ATTEMPTS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			maxAttempts = parsed
		}
	}

	return NewExponentialBackoff(baseDelay, maxDelay, maxAttempts)
}

// ShouldRetry determines if an operation should be retried
func (e *ExponentialBackoff) ShouldRetry(err error, attempt int) bool {
	if attempt >= e.maxAttempts {
		return false
	}

	// Check if error is retryable
	if structured, ok := err.(*StructuredError); ok {
		return e.retryableCodes[structured.Code]
	}

	return false
}

// NextDelay calculates the delay before next retry
func (e *ExponentialBackoff) NextDelay(attempt int) time.Duration {
	delay := time.Duration(math.Pow(2, float64(attempt))) * e.baseDelay
	if delay > e.maxDelay {
		return e.maxDelay
	}
	return delay
}

// MaxAttempts returns the maximum number of retry attempts
func (e *ExponentialBackoff) MaxAttempts() int {
	return e.maxAttempts
}

// SetRetryableCode adds or removes a code from the retryable list
func (e *ExponentialBackoff) SetRetryableCode(code ErrorCode, retryable bool) {
	e.retryableCodes[code] = retryable
}