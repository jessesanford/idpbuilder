package retry

import (
	"math"
	"math/rand"
	"time"
)

// BackoffStrategy defines the exponential backoff configuration for retries
type BackoffStrategy struct {
	// InitialInterval is the initial retry interval
	InitialInterval time.Duration
	// MaxInterval is the maximum retry interval
	MaxInterval time.Duration
	// Multiplier is the backoff multiplier for exponential growth
	Multiplier float64
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int
	// Jitter adds randomness to retry intervals to avoid thundering herd
	Jitter bool
}

// DefaultBackoff returns a default backoff strategy suitable for registry operations
func DefaultBackoff() *BackoffStrategy {
	return &BackoffStrategy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		MaxRetries:      10, // 5-10 retries as specified in requirements
		Jitter:          true,
	}
}

// ConservativeBackoff returns a more conservative backoff strategy
func ConservativeBackoff() *BackoffStrategy {
	return &BackoffStrategy{
		InitialInterval: 2 * time.Second,
		MaxInterval:     60 * time.Second,
		Multiplier:      1.5,
		MaxRetries:      5,
		Jitter:          true,
	}
}

// AggressiveBackoff returns an aggressive backoff strategy for quick retries
func AggressiveBackoff() *BackoffStrategy {
	return &BackoffStrategy{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.5,
		MaxRetries:      15,
		Jitter:          true,
	}
}

// NextInterval calculates the next retry interval based on the attempt number
func (b *BackoffStrategy) NextInterval(attempt int) time.Duration {
	if attempt >= b.MaxRetries {
		return 0 // No more retries
	}

	// Calculate exponential backoff interval
	interval := float64(b.InitialInterval) * math.Pow(b.Multiplier, float64(attempt))

	// Cap at maximum interval
	if interval > float64(b.MaxInterval) {
		interval = float64(b.MaxInterval)
	}

	// Add jitter to prevent thundering herd
	if b.Jitter {
		// Apply jitter: multiply by random factor between 0.5 and 1.5
		jitterFactor := 0.5 + rand.Float64()
		interval = interval * jitterFactor
	}

	return time.Duration(interval)
}

// TotalMaxDuration calculates the maximum total duration for all retries
func (b *BackoffStrategy) TotalMaxDuration() time.Duration {
	total := time.Duration(0)
	for i := 0; i < b.MaxRetries; i++ {
		// Use worst-case scenario without jitter for max calculation
		interval := float64(b.InitialInterval) * math.Pow(b.Multiplier, float64(i))
		if interval > float64(b.MaxInterval) {
			interval = float64(b.MaxInterval)
		}

		// Add jitter worst case (1.5x multiplier)
		if b.Jitter {
			interval = interval * 1.5
		}

		total += time.Duration(interval)
	}
	return total
}

// Validate checks if the backoff strategy configuration is valid
func (b *BackoffStrategy) Validate() error {
	if b.InitialInterval <= 0 {
		return ErrInvalidInterval
	}
	if b.MaxInterval <= 0 || b.MaxInterval < b.InitialInterval {
		return ErrInvalidMaxInterval
	}
	if b.Multiplier <= 1.0 {
		return ErrInvalidMultiplier
	}
	if b.MaxRetries < 0 {
		return ErrInvalidMaxRetries
	}
	return nil
}

// Clone creates a copy of the backoff strategy
func (b *BackoffStrategy) Clone() *BackoffStrategy {
	return &BackoffStrategy{
		InitialInterval: b.InitialInterval,
		MaxInterval:     b.MaxInterval,
		Multiplier:      b.Multiplier,
		MaxRetries:      b.MaxRetries,
		Jitter:          b.Jitter,
	}
}