package retry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"
)

// RetryableFunc is a function that can be retried.
// It receives a context and attempt number, and returns an error if the operation failed.
type RetryableFunc func(ctx context.Context, attempt int) error

// Config holds configuration for retry behavior.
type Config struct {
	// MaxAttempts is the maximum number of attempts (including the initial attempt).
	// Must be at least 1.
	MaxAttempts int

	// BackoffStrategy determines the delay between retry attempts.
	BackoffStrategy BackoffStrategy

	// ShouldRetry is an optional function to determine if an error is retryable.
	// If nil, IsRetryable is used as the default.
	ShouldRetry func(error) bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		MaxAttempts:     3,
		BackoffStrategy: NewExponentialBackoff(),
		ShouldRetry:     nil, // Will use IsRetryable
	}
}

// WithRetry executes the given function with retry logic.
// It will retry transient errors according to the provided configuration.
// Returns the error from the last attempt if all retries are exhausted.
func WithRetry(ctx context.Context, config *Config, fn RetryableFunc) error {
	if config == nil {
		config = DefaultConfig()
	}

	if config.MaxAttempts < 1 {
		return fmt.Errorf("max attempts must be at least 1, got %d", config.MaxAttempts)
	}

	shouldRetry := config.ShouldRetry
	if shouldRetry == nil {
		shouldRetry = IsRetryable
	}

	var lastErr error

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// Check context before attempting
		if err := ctx.Err(); err != nil {
			return err
		}

		// Execute the function
		lastErr = fn(ctx, attempt)

		// Success case
		if lastErr == nil {
			return nil
		}

		// Check if we should retry this error
		if !shouldRetry(lastErr) {
			return lastErr
		}

		// Don't sleep after the last attempt
		if attempt < config.MaxAttempts-1 {
			delay := config.BackoffStrategy.NextDelay(attempt)

			// Wait with context cancellation support
			if err := Wait(ctx, delay); err != nil {
				// Context was cancelled during backoff
				return err
			}
		}
	}

	// All retries exhausted
	return &MaxRetriesExceededError{
		Attempts: config.MaxAttempts,
		LastErr:  lastErr,
	}
}

// IsRetryable determines if an error should trigger a retry.
// It checks for transient network errors and retryable HTTP status codes.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for context errors (never retry these)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Check for network errors
	if isNetworkError(err) {
		return true
	}

	// Check for HTTP status codes
	if isRetryableHTTPStatus(err) {
		return true
	}

	return false
}

// isNetworkError checks if an error is a transient network error.
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific network error types first
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		// Check for connection refused, connection reset, etc.
		return isTransientOpError(opErr)
	}

	// Check for net.Error interface (timeout, temporary)
	var netErr net.Error
	if errors.As(err, &netErr) {
		// Retry on temporary errors or timeouts
		if netErr.Timeout() {
			return true
		}
		// Check if Temporary is implemented and true
		if temp, ok := err.(interface{ Temporary() bool }); ok && temp.Temporary() {
			return true
		}
	}

	// Check error message for common network error patterns
	errMsg := err.Error()
	networkPatterns := []string{
		"connection refused",
		"connection reset",
		"broken pipe",
		"no route to host",
		"network is unreachable",
		"i/o timeout",
		"TLS handshake timeout",
		"EOF",
	}

	for _, pattern := range networkPatterns {
		if strings.Contains(strings.ToLower(errMsg), strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

// isTransientOpError checks if a net.OpError is transient and retryable.
func isTransientOpError(opErr *net.OpError) bool {
	if opErr == nil {
		return false
	}

	// Check for specific syscall errors
	if opErr.Err != nil {
		if errors.Is(opErr.Err, syscall.ECONNREFUSED) ||
			errors.Is(opErr.Err, syscall.ECONNRESET) ||
			errors.Is(opErr.Err, syscall.EPIPE) ||
			errors.Is(opErr.Err, syscall.ETIMEDOUT) {
			return true
		}
	}

	// Check if OpError has timeout
	if opErr.Timeout() {
		return true
	}

	// Check if Temporary method exists and returns true
	if temp, ok := interface{}(opErr).(interface{ Temporary() bool }); ok && temp.Temporary() {
		return true
	}

	return false
}

// isRetryableHTTPStatus checks if an error contains an HTTP status code that should be retried.
func isRetryableHTTPStatus(err error) bool {
	if err == nil {
		return false
	}

	// This is a simplified check. In a real implementation, you might want to
	// extract status codes from custom error types or http.Response errors.
	errMsg := err.Error()

	// Check for common retryable status codes in error messages
	retryableStatuses := []string{
		"408", // Request Timeout
		"429", // Too Many Requests
		"500", // Internal Server Error
		"502", // Bad Gateway
		"503", // Service Unavailable
		"504", // Gateway Timeout
	}

	for _, status := range retryableStatuses {
		if strings.Contains(errMsg, status) {
			return true
		}
	}

	return false
}

// WithRetrySimple is a simplified retry wrapper with default configuration.
// Retries up to 3 times with exponential backoff.
func WithRetrySimple(ctx context.Context, fn func() error) error {
	return WithRetry(ctx, DefaultConfig(), func(ctx context.Context, attempt int) error {
		return fn()
	})
}

// WithRetryN is a convenience wrapper that retries up to N times.
func WithRetryN(ctx context.Context, maxAttempts int, fn RetryableFunc) error {
	config := DefaultConfig()
	config.MaxAttempts = maxAttempts
	return WithRetry(ctx, config, fn)
}

// RetryIfTransient wraps a function and retries only on transient network errors.
func RetryIfTransient(ctx context.Context, maxAttempts int, fn func() error) error {
	config := &Config{
		MaxAttempts:     maxAttempts,
		BackoffStrategy: NewExponentialBackoff(),
		ShouldRetry:     isNetworkError, // Only retry network errors
	}

	return WithRetry(ctx, config, func(ctx context.Context, attempt int) error {
		return fn()
	})
}
