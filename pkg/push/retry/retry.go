package retry

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// RetryableFunc is a function that can be retried
type RetryableFunc func() error

// RetryableError indicates whether an error should trigger a retry
type RetryableError interface {
	IsRetryable() bool
}

// IsRetryable determines if an error should trigger a retry based on common transient error patterns
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check if error implements RetryableError interface
	if retryableErr, ok := err.(RetryableError); ok {
		return retryableErr.IsRetryable()
	}

	errMsg := strings.ToLower(err.Error())

	// Network-related transient errors
	networkErrors := []string{
		"connection refused",
		"connection reset",
		"connection timeout",
		"timeout",
		"temporary failure",
		"network is unreachable",
		"no such host",
		"i/o timeout",
		"broken pipe",
		"connection aborted",
	}

	for _, netErr := range networkErrors {
		if strings.Contains(errMsg, netErr) {
			return true
		}
	}

	// Check for specific error types
	if netErr, ok := err.(net.Error); ok {
		return netErr.Temporary() || netErr.Timeout()
	}

	// HTTP status code based retry logic
	if httpErr, ok := err.(*http.Response); ok {
		switch httpErr.StatusCode {
		case http.StatusTooManyRequests,      // 429
			 http.StatusInternalServerError,   // 500
			 http.StatusBadGateway,           // 502
			 http.StatusServiceUnavailable,   // 503
			 http.StatusGatewayTimeout:       // 504
			return true
		}
	}

	// Registry-specific transient errors
	registryErrors := []string{
		"manifest unknown",
		"blob unknown",
		"upload unknown",
		"layer does not exist",
		"repository does not exist",
	}

	for _, regErr := range registryErrors {
		if strings.Contains(errMsg, regErr) {
			// These might be retryable in some contexts
			return true
		}
	}

	return false
}

// WithRetry executes a function with retry logic using exponential backoff
func WithRetry(ctx context.Context, fn RetryableFunc, strategy *BackoffStrategy) error {
	if fn == nil {
		return fmt.Errorf("retry function cannot be nil")
	}

	if strategy == nil {
		strategy = DefaultBackoff()
	}

	if err := strategy.Validate(); err != nil {
		return fmt.Errorf("invalid backoff strategy: %w", err)
	}

	var lastErr error
	attempt := 0

	for attempt <= strategy.MaxRetries {
		// Check context cancellation before attempt
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled after %d attempts: %w", attempt, ctx.Err())
		default:
		}

		// Wait before retry (except for first attempt)
		if attempt > 0 {
			interval := strategy.NextInterval(attempt - 1)
			if interval <= 0 {
				break // No more retries
			}

			fmt.Printf("Attempt %d failed, retrying in %v: %v\n", attempt, interval, lastErr)

			select {
			case <-time.After(interval):
				// Continue with retry
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled during backoff after %d attempts: %w", attempt, ctx.Err())
			}
		}

		// Execute the function
		err := fn()
		if err == nil {
			if attempt > 0 {
				fmt.Printf("Success after %d retries\n", attempt)
			}
			return nil // Success
		}

		lastErr = err
		attempt++

		// Check if error is retryable
		if !IsRetryable(err) {
			return fmt.Errorf("non-retryable error after %d attempts: %w", attempt, err)
		}

		// Log the attempt if not the last one
		if attempt <= strategy.MaxRetries {
			fmt.Printf("Attempt %d failed with retryable error: %v\n", attempt, err)
		}
	}

	// All retries exhausted
	return fmt.Errorf("max retries (%d) exceeded: %w", strategy.MaxRetries, lastErr)
}

// WithRetryAndCallback executes a function with retry logic and calls a callback on each attempt
func WithRetryAndCallback(ctx context.Context, fn RetryableFunc, strategy *BackoffStrategy,
	onRetry func(attempt int, err error, nextInterval time.Duration)) error {

	if fn == nil {
		return fmt.Errorf("retry function cannot be nil")
	}

	if strategy == nil {
		strategy = DefaultBackoff()
	}

	var lastErr error
	attempt := 0

	for attempt <= strategy.MaxRetries {
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled after %d attempts: %w", attempt, ctx.Err())
		default:
		}

		if attempt > 0 {
			interval := strategy.NextInterval(attempt - 1)
			if interval <= 0 {
				break
			}

			// Call the callback before waiting
			if onRetry != nil {
				onRetry(attempt, lastErr, interval)
			}

			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled during backoff: %w", ctx.Err())
			}
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err
		attempt++

		if !IsRetryable(err) {
			return err
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}