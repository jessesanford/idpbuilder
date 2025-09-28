package errors

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// RetryHandler implements retry logic with exponential backoff
type RetryHandler struct {
	config RetryInfo
	random *rand.Rand
}

// NewRetryHandler creates a new retry handler with the given configuration
func NewRetryHandler(config RetryInfo) *RetryHandler {
	return &RetryHandler{
		config: config,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewDefaultRetryHandler creates a retry handler with default configuration
func NewDefaultRetryHandler() *RetryHandler {
	return NewRetryHandler(DefaultRetryInfo())
}

// Retry executes a function with retry logic based on error classification
func (r *RetryHandler) Retry(ctx context.Context, operation string, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on the last attempt
		if attempt == r.config.MaxAttempts {
			break
		}

		// Check if error should be retried
		if !IsRetryable(err) {
			break
		}

		// Calculate delay with exponential backoff
		delay := r.calculateDelay(attempt)

		// Wait for the delay, respecting context cancellation
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return lastErr
}

// RetryWithResource executes a function with retry logic and resource context
func (r *RetryHandler) RetryWithResource(ctx context.Context, operation, resource string, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}

		// Wrap error with operation context
		wrappedErr := &OperationError{
			Operation: operation,
			Resource:  resource,
			Category:  ClassifyError(err),
			Retryable: IsRetryable(err),
			Attempt:   attempt,
			Timestamp: time.Now(),
			Cause:     err,
			Context:   make(map[string]interface{}),
		}

		lastErr = wrappedErr

		// Don't retry on the last attempt
		if attempt == r.config.MaxAttempts {
			break
		}

		// Check if error should be retried
		if !wrappedErr.Retryable {
			break
		}

		// Calculate delay with exponential backoff
		delay := r.calculateDelay(attempt)

		// Add delay information to error context
		wrappedErr.Context["retry_delay"] = delay
		wrappedErr.Context["next_attempt"] = attempt + 1

		// Wait for the delay, respecting context cancellation
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return lastErr
}

// calculateDelay computes the delay for a given attempt using exponential backoff
func (r *RetryHandler) calculateDelay(attempt int) time.Duration {
	// Calculate exponential backoff: baseDelay * multiplier^(attempt-1)
	delay := float64(r.config.BaseDelay) * math.Pow(r.config.Multiplier, float64(attempt-1))

	// Apply maximum delay limit
	if delay > float64(r.config.MaxDelay) {
		delay = float64(r.config.MaxDelay)
	}

	// Add jitter if enabled (±25% randomness)
	if r.config.Jitter && r.random != nil {
		jitter := delay * 0.25 * (r.random.Float64()*2 - 1) // Random value between -0.25 and 0.25
		delay += jitter
	}

	// Ensure delay is not negative
	if delay < 0 {
		delay = float64(r.config.BaseDelay)
	}

	return time.Duration(delay)
}

// ShouldRetry determines if an error should trigger a retry attempt
func (r *RetryHandler) ShouldRetry(err error, attempt int) bool {
	if err == nil {
		return false
	}

	if attempt >= r.config.MaxAttempts {
		return false
	}

	return IsRetryable(err)
}

// GetConfig returns the current retry configuration
func (r *RetryHandler) GetConfig() RetryInfo {
	return r.config
}

// SetConfig updates the retry configuration
func (r *RetryHandler) SetConfig(config RetryInfo) {
	r.config = config
}

// RetryFunc is a helper type for functions that can be retried
type RetryFunc func() error

// RetryWithBackoff is a convenience function for simple retry scenarios
func RetryWithBackoff(ctx context.Context, operation string, maxAttempts int, baseDelay time.Duration, fn RetryFunc) error {
	handler := NewRetryHandler(RetryInfo{
		MaxAttempts: maxAttempts,
		BaseDelay:   baseDelay,
		MaxDelay:    10 * time.Second,
		Multiplier:  2.0,
		Jitter:      true,
	})

	return handler.Retry(ctx, operation, fn)
}

// RetryConfig provides a builder pattern for retry configuration
type RetryConfig struct {
	info RetryInfo
}

// NewRetryConfig creates a new retry configuration builder
func NewRetryConfig() *RetryConfig {
	return &RetryConfig{
		info: DefaultRetryInfo(),
	}
}

// WithMaxAttempts sets the maximum number of retry attempts
func (c *RetryConfig) WithMaxAttempts(attempts int) *RetryConfig {
	c.info.MaxAttempts = attempts
	return c
}

// WithBaseDelay sets the base delay between retries
func (c *RetryConfig) WithBaseDelay(delay time.Duration) *RetryConfig {
	c.info.BaseDelay = delay
	return c
}

// WithMaxDelay sets the maximum delay between retries
func (c *RetryConfig) WithMaxDelay(delay time.Duration) *RetryConfig {
	c.info.MaxDelay = delay
	return c
}

// WithMultiplier sets the backoff multiplier
func (c *RetryConfig) WithMultiplier(multiplier float64) *RetryConfig {
	c.info.Multiplier = multiplier
	return c
}

// WithJitter enables or disables jitter
func (c *RetryConfig) WithJitter(enabled bool) *RetryConfig {
	c.info.Jitter = enabled
	return c
}

// Build returns the configured RetryInfo
func (c *RetryConfig) Build() RetryInfo {
	return c.info
}