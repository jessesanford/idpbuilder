package progress

// ProgressReporter defines the interface for reporting progress of long-running operations
// This is based on the interface specification from P1W1-E4
type ProgressReporter interface {
	// Start begins a new progress session with the given message
	Start(message string)

	// ReportProgress updates the current progress
	ReportProgress(current, total int64, message string)

	// Complete marks the operation as completed
	Complete(message string)

	// Error reports an error that occurred during the operation
	Error(err error)
}