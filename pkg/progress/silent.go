package progress

import (
	"context"
)

// SilentReporter implements ProgressReporter with no output.
// This is useful for programmatic usage where progress output is not desired.
type SilentReporter struct{}

// NewSilentReporter creates a new silent progress reporter.
func NewSilentReporter() *SilentReporter {
	return &SilentReporter{}
}

// StartOperation silently begins tracking an operation.
func (s *SilentReporter) StartOperation(ctx context.Context, opType OperationType, reference string) context.Context {
	return ctx
}

// UpdateProgress silently receives progress updates.
func (s *SilentReporter) UpdateProgress(ctx context.Context, progress OperationProgress) {
	// Silent - no output
}

// CompleteOperation silently marks an operation as completed.
func (s *SilentReporter) CompleteOperation(ctx context.Context, reference string, message string) {
	// Silent - no output
}

// ErrorOperation silently receives error notifications.
func (s *SilentReporter) ErrorOperation(ctx context.Context, reference string, err error) {
	// Silent - no output
}

// CancelOperation silently handles cancellation.
func (s *SilentReporter) CancelOperation(ctx context.Context, reference string, reason string) {
	// Silent - no output
}

// CallbackReporter implements ProgressReporter by calling a provided function.
// This is useful for testing and custom progress handling.
type CallbackReporter struct {
	callback ProgressCallback
}

// NewCallbackReporter creates a progress reporter that calls the provided function.
func NewCallbackReporter(callback ProgressCallback) *CallbackReporter {
	return &CallbackReporter{
		callback: callback,
	}
}

// StartOperation calls the callback with start information.
func (c *CallbackReporter) StartOperation(ctx context.Context, opType OperationType, reference string) context.Context {
	if c.callback != nil {
		c.callback(OperationProgress{
			OperationType: opType,
			Reference:     reference,
			State:         StateStarted,
			Message:       "Operation started",
		})
	}
	return ctx
}

// UpdateProgress calls the callback with progress information.
func (c *CallbackReporter) UpdateProgress(ctx context.Context, progress OperationProgress) {
	if c.callback != nil {
		c.callback(progress)
	}
}

// CompleteOperation calls the callback with completion information.
func (c *CallbackReporter) CompleteOperation(ctx context.Context, reference string, message string) {
	if c.callback != nil {
		c.callback(OperationProgress{
			Reference: reference,
			State:     StateCompleted,
			Message:   message,
		})
	}
}

// ErrorOperation calls the callback with error information.
func (c *CallbackReporter) ErrorOperation(ctx context.Context, reference string, err error) {
	if c.callback != nil {
		c.callback(OperationProgress{
			Reference: reference,
			State:     StateError,
			Error:     err,
			Message:   err.Error(),
		})
	}
}

// CancelOperation calls the callback with cancellation information.
func (c *CallbackReporter) CancelOperation(ctx context.Context, reference string, reason string) {
	if c.callback != nil {
		c.callback(OperationProgress{
			Reference: reference,
			State:     StateCancelled,
			Message:   reason,
		})
	}
}