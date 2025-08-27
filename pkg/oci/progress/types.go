package progress

import (
	"time"
)

// ProgressStatus represents the status of a progress operation
type ProgressStatus int

const (
	// ProgressStatusStarted indicates the operation has started
	ProgressStatusStarted ProgressStatus = iota + 1
	// ProgressStatusInProgress indicates the operation is in progress
	ProgressStatusInProgress
	// ProgressStatusCompleted indicates the operation has completed successfully
	ProgressStatusCompleted
	// ProgressStatusFailed indicates the operation has failed
	ProgressStatusFailed
	// ProgressStatusCanceled indicates the operation was canceled
	ProgressStatusCanceled
)

// String returns the string representation of ProgressStatus
func (ps ProgressStatus) String() string {
	switch ps {
	case ProgressStatusStarted:
		return "Started"
	case ProgressStatusInProgress:
		return "InProgress"
	case ProgressStatusCompleted:
		return "Completed"
	case ProgressStatusFailed:
		return "Failed"
	case ProgressStatusCanceled:
		return "Canceled"
	default:
		return "Unknown"
	}
}

// ProgressEvent represents a progress update event
type ProgressEvent struct {
	// ID is the unique identifier for this progress event
	ID string `json:"id"`
	// Operation is the name of the operation being tracked
	Operation string `json:"operation"`
	// Phase is the current phase within the operation
	Phase string `json:"phase"`
	// Status is the current status of the operation
	Status ProgressStatus `json:"status"`
	// Current is the current progress value
	Current int64 `json:"current"`
	// Total is the total expected value (0 if indeterminate)
	Total int64 `json:"total"`
	// Percent is the completion percentage (0-100)
	Percent float64 `json:"percent"`
	// Message is an optional human-readable message
	Message string `json:"message,omitempty"`
	// Timestamp is when this event was generated
	Timestamp time.Time `json:"timestamp"`
	// Duration is how long the operation has been running
	Duration time.Duration `json:"duration"`
	// ETA is the estimated time to completion
	ETA *time.Duration `json:"eta,omitempty"`
}

// ProgressReporter defines the interface for reporting progress
type ProgressReporter interface {
	// ReportProgress reports a progress event
	ReportProgress(event *ProgressEvent)
	// Start marks the beginning of an operation
	Start(operationID, operation string)
	// Update updates the progress of an operation
	Update(operationID string, current, total int64, message string)
	// Complete marks an operation as completed
	Complete(operationID string, message string)
	// Fail marks an operation as failed
	Fail(operationID string, err error, message string)
	// Cancel marks an operation as canceled
	Cancel(operationID string, message string)
}

// BuildProgress represents progress for build operations
type BuildProgress struct {
	// ImageName is the name of the image being built
	ImageName string `json:"image_name"`
	// BuildID is the unique identifier for this build
	BuildID string `json:"build_id"`
	// Step is the current build step
	Step string `json:"step"`
	// StepNumber is the current step number
	StepNumber int `json:"step_number"`
	// TotalSteps is the total number of build steps
	TotalSteps int `json:"total_steps"`
	// LayersBuilt is the number of layers built
	LayersBuilt int `json:"layers_built"`
	// TotalLayers is the total number of layers to build
	TotalLayers int `json:"total_layers"`
	// BytesProcessed is the number of bytes processed
	BytesProcessed int64 `json:"bytes_processed"`
	// TotalBytes is the total bytes to process
	TotalBytes int64 `json:"total_bytes"`
	// StartTime is when the build started
	StartTime time.Time `json:"start_time"`
	// LastUpdate is when progress was last updated
	LastUpdate time.Time `json:"last_update"`
}

// PushProgress represents progress for registry push operations
type PushProgress struct {
	// ImageName is the name of the image being pushed
	ImageName string `json:"image_name"`
	// Registry is the target registry
	Registry string `json:"registry"`
	// PushID is the unique identifier for this push
	PushID string `json:"push_id"`
	// LayersPushed is the number of layers pushed
	LayersPushed int `json:"layers_pushed"`
	// TotalLayers is the total number of layers to push
	TotalLayers int `json:"total_layers"`
	// LayersSkipped is the number of layers skipped (already exist)
	LayersSkipped int `json:"layers_skipped"`
	// BytesUploaded is the number of bytes uploaded
	BytesUploaded int64 `json:"bytes_uploaded"`
	// TotalBytes is the total bytes to upload
	TotalBytes int64 `json:"total_bytes"`
	// CurrentLayer is the layer currently being pushed
	CurrentLayer string `json:"current_layer"`
	// StartTime is when the push started
	StartTime time.Time `json:"start_time"`
	// LastUpdate is when progress was last updated
	LastUpdate time.Time `json:"last_update"`
}

// NewProgressEvent creates a new progress event
func NewProgressEvent(id, operation, phase string, status ProgressStatus) *ProgressEvent {
	return &ProgressEvent{
		ID:        id,
		Operation: operation,
		Phase:     phase,
		Status:    status,
		Timestamp: time.Now(),
	}
}

// WithProgress adds current/total progress information
func (pe *ProgressEvent) WithProgress(current, total int64) *ProgressEvent {
	pe.Current = current
	pe.Total = total
	if total > 0 {
		pe.Percent = float64(current) / float64(total) * 100
	}
	return pe
}

// WithMessage adds a message to the progress event
func (pe *ProgressEvent) WithMessage(message string) *ProgressEvent {
	pe.Message = message
	return pe
}

// WithDuration sets the duration for the progress event
func (pe *ProgressEvent) WithDuration(duration time.Duration) *ProgressEvent {
	pe.Duration = duration
	return pe
}

// WithETA sets the estimated time to completion
func (pe *ProgressEvent) WithETA(eta time.Duration) *ProgressEvent {
	pe.ETA = &eta
	return pe
}

// NewBuildProgress creates a new build progress tracker
func NewBuildProgress(imageName, buildID string, totalSteps int) *BuildProgress {
	return &BuildProgress{
		ImageName:  imageName,
		BuildID:    buildID,
		TotalSteps: totalSteps,
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}
}

// NewPushProgress creates a new push progress tracker
func NewPushProgress(imageName, registry, pushID string, totalLayers int) *PushProgress {
	return &PushProgress{
		ImageName:    imageName,
		Registry:     registry,
		PushID:       pushID,
		TotalLayers:  totalLayers,
		StartTime:    time.Now(),
		LastUpdate:   time.Now(),
	}
}