// Package progress provides comprehensive progress tracking for long-running OCI operations.
// This package defines progress tracking types, interfaces, and implementations for
// monitoring and reporting operation progress across the idpbuilder system.
package progress

import "time"

// Progress status constants define the state of an operation
const (
	// StatusNotStarted indicates the operation has not begun
	StatusNotStarted = "not_started"
	
	// StatusInProgress indicates the operation is currently running
	StatusInProgress = "in_progress"
	
	// StatusPaused indicates the operation is temporarily paused
	StatusPaused = "paused"
	
	// StatusCompleted indicates the operation finished successfully
	StatusCompleted = "completed"
	
	// StatusFailed indicates the operation failed
	StatusFailed = "failed"
	
	// StatusCancelled indicates the operation was cancelled
	StatusCancelled = "cancelled"
	
	// StatusTimeout indicates the operation timed out
	StatusTimeout = "timeout"
)

// Progress stages for OCI operations
const (
	// StageInitialization represents the setup phase
	StageInitialization = "initialization"
	
	// StageAuthentication represents the authentication phase
	StageAuthentication = "authentication"
	
	// StageDownload represents the download phase
	StageDownload = "download"
	
	// StageExtraction represents the extraction phase
	StageExtraction = "extraction"
	
	// StageValidation represents the validation phase
	StageValidation = "validation"
	
	// StageInstallation represents the installation phase
	StageInstallation = "installation"
	
	// StageConfiguration represents the configuration phase
	StageConfiguration = "configuration"
	
	// StageVerification represents the verification phase
	StageVerification = "verification"
	
	// StageCleanup represents the cleanup phase
	StageCleanup = "cleanup"
	
	// StageFinalization represents the final completion phase
	StageFinalization = "finalization"
)

// Event types for progress notifications
const (
	// EventTypeStarted indicates an operation started
	EventTypeStarted = "started"
	
	// EventTypeProgress indicates progress update
	EventTypeProgress = "progress"
	
	// EventTypeStageChanged indicates stage transition
	EventTypeStageChanged = "stage_changed"
	
	// EventTypePaused indicates operation paused
	EventTypePaused = "paused"
	
	// EventTypeResumed indicates operation resumed
	EventTypeResumed = "resumed"
	
	// EventTypeCompleted indicates operation completed
	EventTypeCompleted = "completed"
	
	// EventTypeFailed indicates operation failed
	EventTypeFailed = "failed"
	
	// EventTypeCancelled indicates operation cancelled
	EventTypeCancelled = "cancelled"
	
	// EventTypeTimeout indicates operation timed out
	EventTypeTimeout = "timeout"
)

// Default timeout and retry values
const (
	// DefaultProgressTimeout is the default timeout for progress operations
	DefaultProgressTimeout = 30 * time.Minute
	
	// DefaultRetryInterval is the default interval between retry attempts
	DefaultRetryInterval = 5 * time.Second
	
	// DefaultMaxRetries is the default maximum number of retry attempts
	DefaultMaxRetries = 3
	
	// DefaultUpdateInterval is the default interval for progress updates
	DefaultUpdateInterval = 1 * time.Second
	
	// DefaultStaleThreshold is the default time after which progress is considered stale
	DefaultStaleThreshold = 5 * time.Minute
)

// Progress reporting thresholds
const (
	// MinProgressPercent is the minimum valid progress percentage
	MinProgressPercent = 0
	
	// MaxProgressPercent is the maximum valid progress percentage
	MaxProgressPercent = 100
	
	// ProgressIncrementThreshold is the minimum increment to report
	ProgressIncrementThreshold = 1
)