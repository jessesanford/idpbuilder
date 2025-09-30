package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// BackoffStrategy defines an interface for retry backoff strategies.
type BackoffStrategy interface {
	// NextDelay calculates the delay for the next retry attempt.
	// Returns the duration to wait before the next attempt.
	NextDelay(attempt int) time.Duration

	// Reset resets the backoff strategy to its initial state.
	Reset()
}

// ExponentialBackoff implements exponential backoff with jitter.
// This prevents thundering herd problems in distributed systems.
type ExponentialBackoff struct {
	// BaseDelay is the initial delay for the first retry.
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries.
	MaxDelay time.Duration

	// Multiplier is the factor by which delay increases each attempt.
	// Typically 2.0 for exponential backoff.
	Multiplier float64

	// JitterFraction is the fraction of delay to add as jitter (0.0 to 1.0).
	// For example, 0.1 means up to 10% jitter.
	JitterFraction float64

	// random source for jitter calculation
	rng *rand.Rand
}

// NewExponentialBackoff creates a new exponential backoff strategy with defaults.
func NewExponentialBackoff() *ExponentialBackoff {
	return &ExponentialBackoff{
		BaseDelay:      100 * time.Millisecond,
		MaxDelay:       30 * time.Second,
		Multiplier:     2.0,
		JitterFraction: 0.1,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewExponentialBackoffWithConfig creates a backoff strategy with custom configuration.
func NewExponentialBackoffWithConfig(baseDelay, maxDelay time.Duration, multiplier, jitter float64) *ExponentialBackoff {
	return &ExponentialBackoff{
		BaseDelay:      baseDelay,
		MaxDelay:       maxDelay,
		Multiplier:     multiplier,
		JitterFraction: jitter,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NextDelay calculates the delay for the next retry attempt using exponential backoff with jitter.
func (b *ExponentialBackoff) NextDelay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	// Calculate exponential delay: baseDelay * multiplier^attempt
	delay := float64(b.BaseDelay) * math.Pow(b.Multiplier, float64(attempt))

	// Cap at max delay
	if delay > float64(b.MaxDelay) {
		delay = float64(b.MaxDelay)
	}

	// Add jitter to prevent thundering herd
	jitter := b.calculateJitter(delay)
	finalDelay := time.Duration(delay + jitter)

	// Ensure we don't exceed max delay after jitter
	if finalDelay > b.MaxDelay {
		finalDelay = b.MaxDelay
	}

	return finalDelay
}

// calculateJitter adds randomness to the delay to distribute retry load.
// The jitter is uniformly distributed between 0 and (delay * jitterFraction).
func (b *ExponentialBackoff) calculateJitter(delay float64) float64 {
	if b.JitterFraction <= 0 {
		return 0
	}

	maxJitter := delay * b.JitterFraction
	return b.rng.Float64() * maxJitter
}

// Reset resets the backoff strategy. For exponential backoff, this is a no-op
// as each attempt's delay is calculated independently.
func (b *ExponentialBackoff) Reset() {
	// Exponential backoff is stateless - each NextDelay call uses the attempt number
	// No state to reset
}

// Wait waits for the backoff duration with context cancellation support.
// Returns nil if the wait completed, or ctx.Err() if context was cancelled.
func Wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// ConstantBackoff implements a constant delay backoff strategy.
// Useful for testing or specific retry scenarios.
type ConstantBackoff struct {
	Delay time.Duration
}

// NewConstantBackoff creates a constant backoff strategy.
func NewConstantBackoff(delay time.Duration) *ConstantBackoff {
	return &ConstantBackoff{
		Delay: delay,
	}
}

// NextDelay returns a constant delay regardless of attempt number.
func (c *ConstantBackoff) NextDelay(attempt int) time.Duration {
	return c.Delay
}

// Reset is a no-op for constant backoff.
func (c *ConstantBackoff) Reset() {
	// Constant backoff has no state to reset
}