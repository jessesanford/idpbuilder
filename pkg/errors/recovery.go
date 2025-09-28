package errors

import "context"

// RecoveryHandler defines error recovery strategies
type RecoveryHandler interface {
	// CanRecover determines if recovery is possible
	CanRecover(err error) bool

	// Recover attempts to recover from the error
	Recover(err error) error

	// RecoverWithContext attempts recovery with context support
	RecoverWithContext(ctx context.Context, err error) error
}

// RecoveryFunc is a function that can recover from an error
type RecoveryFunc func(error) error

// DefaultRecovery provides basic recovery patterns
type DefaultRecovery struct {
	recoveryStrategies map[ErrorCode]RecoveryFunc
}

// NewDefaultRecovery creates a new recovery handler with default strategies
func NewDefaultRecovery() *DefaultRecovery {
	recovery := &DefaultRecovery{
		recoveryStrategies: make(map[ErrorCode]RecoveryFunc),
	}

	// Set up default recovery strategies
	recovery.recoveryStrategies[ConfigurationError] = recovery.recoverConfigurationError
	recovery.recoveryStrategies[ValidationFailed] = recovery.recoverValidationError

	return recovery
}

// CanRecover determines if recovery is possible for the given error
func (d *DefaultRecovery) CanRecover(err error) bool {
	if structured, ok := err.(*StructuredError); ok {
		_, hasStrategy := d.recoveryStrategies[structured.Code]
		return hasStrategy
	}
	return false
}

// Recover attempts to recover from the error
func (d *DefaultRecovery) Recover(err error) error {
	if structured, ok := err.(*StructuredError); ok {
		if strategy, exists := d.recoveryStrategies[structured.Code]; exists {
			return strategy(err)
		}
	}
	return err
}

// RecoverWithContext attempts recovery with context support
func (d *DefaultRecovery) RecoverWithContext(ctx context.Context, err error) error {
	// Check if context is cancelled before attempting recovery
	select {
	case <-ctx.Done():
		return NewStructuredError(InternalError, "recovery", "context cancelled during recovery", ctx.Err())
	default:
		return d.Recover(err)
	}
}

// AddRecoveryStrategy adds a custom recovery strategy for an error code
func (d *DefaultRecovery) AddRecoveryStrategy(code ErrorCode, strategy RecoveryFunc) {
	d.recoveryStrategies[code] = strategy
}

// RemoveRecoveryStrategy removes a recovery strategy for an error code
func (d *DefaultRecovery) RemoveRecoveryStrategy(code ErrorCode) {
	delete(d.recoveryStrategies, code)
}

// recoverConfigurationError attempts to recover from configuration errors
func (d *DefaultRecovery) recoverConfigurationError(err error) error {
	// For configuration errors, we can provide guidance on fixing the issue
	if structured, ok := err.(*StructuredError); ok {
		return NewStructuredError(
			ConfigurationError,
			structured.Op,
			"Configuration error detected. Please check your configuration settings and try again.",
			structured.Cause,
		)
	}
	return err
}

// recoverValidationError attempts to recover from validation errors
func (d *DefaultRecovery) recoverValidationError(err error) error {
	// For validation errors, we can provide actionable feedback
	if structured, ok := err.(*StructuredError); ok {
		return NewStructuredError(
			ValidationFailed,
			structured.Op,
			"Validation failed. Please verify your input parameters and retry.",
			structured.Cause,
		)
	}
	return err
}