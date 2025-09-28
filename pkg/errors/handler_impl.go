package errors

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"
)

// ErrorHandlerImpl provides full error handling implementation
type ErrorHandlerImpl struct {
	retry    RetryStrategy
	recovery RecoveryHandler
}

// NewErrorHandler creates a new error handler with default configuration from environment
func NewErrorHandler() *ErrorHandlerImpl {
	return &ErrorHandlerImpl{
		retry:    NewExponentialBackoffFromEnv(),
		recovery: NewDefaultRecovery(),
	}
}

// NewErrorHandlerWithConfig creates a new error handler with custom configuration
func NewErrorHandlerWithConfig(retry RetryStrategy, recovery RecoveryHandler) *ErrorHandlerImpl {
	return &ErrorHandlerImpl{
		retry:    retry,
		recovery: recovery,
	}
}

// Handle processes an error and returns a structured error
func (h *ErrorHandlerImpl) Handle(err error) error {
	if err == nil {
		return nil
	}

	// Check if already structured
	if structured, ok := err.(*StructuredError); ok {
		return h.processStructuredError(structured)
	}

	// Convert to structured error
	structured := h.classify(err)
	return h.processStructuredError(structured)
}

// HandleWithContext processes an error with context for cancellation
func (h *ErrorHandlerImpl) HandleWithContext(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return NewStructuredError(InternalError, "error_handler", "context cancelled during error handling", ctx.Err())
	default:
		return h.Handle(err)
	}
}

// WithRetry adds retry capability to the handler
func (h *ErrorHandlerImpl) WithRetry(strategy RetryStrategy) ErrorHandler {
	h.retry = strategy
	return h
}

// WithRecovery adds recovery capability to the handler
func (h *ErrorHandlerImpl) WithRecovery(recovery RecoveryHandler) ErrorHandler {
	h.recovery = recovery
	return h
}

// classify converts generic errors to structured errors based on error content
func (h *ErrorHandlerImpl) classify(err error) *StructuredError {
	errMsg := strings.ToLower(err.Error())

	// Classification based on error message patterns
	switch {
	case strings.Contains(errMsg, "build") && strings.Contains(errMsg, "failed"):
		return NewStructuredError(BuildFailed, "build", "Build operation failed", err)
	case strings.Contains(errMsg, "build") && strings.Contains(errMsg, "timeout"):
		return NewStructuredError(BuildTimeout, "build", "Build operation timed out", err)
	case strings.Contains(errMsg, "registry") && strings.Contains(errMsg, "unreachable"):
		return NewStructuredError(RegistryUnreachable, "registry", "Registry is unreachable", err)
	case strings.Contains(errMsg, "registry") && strings.Contains(errMsg, "auth"):
		return NewStructuredError(RegistryAuthFailed, "registry", "Registry authentication failed", err)
	case strings.Contains(errMsg, "registry") && strings.Contains(errMsg, "push"):
		return NewStructuredError(RegistryPushFailed, "registry", "Registry push operation failed", err)
	case strings.Contains(errMsg, "certificate") && strings.Contains(errMsg, "invalid"):
		return NewStructuredError(CertificateInvalid, "certificate", "Certificate is invalid", err)
	case strings.Contains(errMsg, "certificate") && strings.Contains(errMsg, "expired"):
		return NewStructuredError(CertificateExpired, "certificate", "Certificate has expired", err)
	case strings.Contains(errMsg, "validation"):
		return NewStructuredError(ValidationFailed, "validation", "Validation failed", err)
	case strings.Contains(errMsg, "config"):
		return NewStructuredError(ConfigurationError, "configuration", "Configuration error", err)
	default:
		return NewStructuredError(InternalError, "unknown", "Unknown error occurred", err)
	}
}

// processStructuredError handles a structured error with retry/recovery logic
func (h *ErrorHandlerImpl) processStructuredError(err *StructuredError) error {
	// First, check if recovery is possible
	if h.recovery != nil && h.recovery.CanRecover(err) {
		recovered := h.recovery.Recover(err)
		// If recovery returns a different error (not the original), use it
		if recovered != err {
			return recovered
		}
	}

	// If no recovery or recovery didn't help, return the structured error
	return err
}

// ExecuteWithRetry executes a function with retry logic
func (h *ErrorHandlerImpl) ExecuteWithRetry(ctx context.Context, operation func() error) error {
	var lastError error

	for attempt := 0; attempt < h.retry.MaxAttempts(); attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return NewStructuredError(InternalError, "retry_executor", "context cancelled during retry", ctx.Err())
		default:
		}

		// Execute the operation
		err := operation()
		if err == nil {
			return nil // Success
		}

		// Process the error
		processedErr := h.Handle(err)
		lastError = processedErr

		// Check if we should retry
		if !h.retry.ShouldRetry(processedErr, attempt) {
			break
		}

		// Calculate delay and wait
		delay := h.retry.NextDelay(attempt)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return NewStructuredError(InternalError, "retry_executor", "context cancelled during retry delay", ctx.Err())
		case <-timer.C:
			// Continue to next attempt
		}
	}

	return lastError
}

// GetRetryableErrorCodes returns the list of error codes that can be retried
func (h *ErrorHandlerImpl) GetRetryableErrorCodes() []ErrorCode {
	codes := make([]ErrorCode, 0)

	// Test each known error code to see if it's retryable
	testCodes := []ErrorCode{
		BuildFailed, BuildTimeout, RegistryUnreachable, RegistryAuthFailed,
		RegistryPushFailed, CertificateInvalid, CertificateExpired,
		ValidationFailed, ConfigurationError, InternalError,
	}

	for _, code := range testCodes {
		testErr := NewStructuredError(code, "test", "test", nil)
		if h.retry.ShouldRetry(testErr, 0) {
			codes = append(codes, code)
		}
	}

	return codes
}

// GetRecoverableErrorCodes returns the list of error codes that can be recovered
func (h *ErrorHandlerImpl) GetRecoverableErrorCodes() []ErrorCode {
	codes := make([]ErrorCode, 0)

	// Test each known error code to see if it's recoverable
	testCodes := []ErrorCode{
		BuildFailed, BuildTimeout, RegistryUnreachable, RegistryAuthFailed,
		RegistryPushFailed, CertificateInvalid, CertificateExpired,
		ValidationFailed, ConfigurationError, InternalError,
	}

	for _, code := range testCodes {
		testErr := NewStructuredError(code, "test", "test", nil)
		if h.recovery.CanRecover(testErr) {
			codes = append(codes, code)
		}
	}

	return codes
}

// loadConfigFromEnv loads configuration values from environment variables
func loadConfigFromEnv() (maxRetries int, baseDelay, maxDelay time.Duration) {
	// Default values
	maxRetries = 3
	baseDelay = 100 * time.Millisecond
	maxDelay = 10 * time.Second

	// Override from environment
	if val := os.Getenv("ERROR_MAX_RETRIES"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			maxRetries = parsed
		}
	}

	if val := os.Getenv("ERROR_BACKOFF_BASE_MS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			baseDelay = time.Duration(parsed) * time.Millisecond
		}
	}

	if val := os.Getenv("ERROR_BACKOFF_MAX_MS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			maxDelay = time.Duration(parsed) * time.Millisecond
		}
	}

	return maxRetries, baseDelay, maxDelay
}