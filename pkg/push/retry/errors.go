package retry

import "errors"

// Backoff strategy validation errors
var (
	ErrInvalidInterval    = errors.New("initial interval must be positive")
	ErrInvalidMaxInterval = errors.New("max interval must be positive and >= initial interval")
	ErrInvalidMultiplier  = errors.New("multiplier must be > 1.0")
	ErrInvalidMaxRetries  = errors.New("max retries must be >= 0")
)