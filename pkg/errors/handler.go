package errors

import "context"

// ErrorHandler provides error handling with recovery strategies
type ErrorHandler interface {
	// Handle processes an error and returns a structured error
	Handle(err error) error

	// HandleWithContext processes an error with context for cancellation
	HandleWithContext(ctx context.Context, err error) error

	// WithRetry adds retry capability to the handler
	WithRetry(strategy RetryStrategy) ErrorHandler

	// WithRecovery adds recovery capability to the handler
	WithRecovery(recovery RecoveryHandler) ErrorHandler
}