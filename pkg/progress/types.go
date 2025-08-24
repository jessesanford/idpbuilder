package progress

import (
	"sync"
	"time"
)

// ProgressStatus represents the current state of a progress operation.
type ProgressStatus string

const (
	// StatusPending indicates the operation has not started
	StatusPending ProgressStatus = "pending"

	// StatusInProgress indicates the operation is currently running
	StatusInProgress ProgressStatus = "in_progress"

	// StatusCompleted indicates the operation finished successfully
	StatusCompleted ProgressStatus = "completed"

	// StatusFailed indicates the operation failed with an error
	StatusFailed ProgressStatus = "failed"

	// StatusCancelled indicates the operation was cancelled
	StatusCancelled ProgressStatus = "cancelled"
)

// String returns the string representation of the progress status.
func (ps ProgressStatus) String() string {
	return string(ps)
}

// IsTerminal returns true if this is a final status (completed, failed, cancelled).
func (ps ProgressStatus) IsTerminal() bool {
	return ps == StatusCompleted || ps == StatusFailed || ps == StatusCancelled
}

// ProgressTracker defines the interface for tracking operation progress.
// Implementations should be thread-safe for concurrent access.
type ProgressTracker interface {
	// Start initializes progress tracking with the total number of work units
	Start(total int64)

	// Update sets the current progress value
	Update(current int64)

	// UpdateWithMessage sets progress and updates the status message
	UpdateWithMessage(current int64, message string)

	// Complete marks the operation as successfully completed
	Complete()

	// Failed marks the operation as failed with the given error
	Failed(error)

	// Cancel marks the operation as cancelled
	Cancel()

	// GetProgress returns the current progress snapshot
	GetProgress() *Progress

	// Subscribe registers a callback for progress events
	Subscribe(callback ProgressCallback) func()
}

// Progress represents a snapshot of progress information.
type Progress struct {
	// Total is the total number of work units (0 means indeterminate)
	Total int64 `json:"total"`

	// Current is the current number of completed work units
	Current int64 `json:"current"`

	// Status is the current progress status
	Status ProgressStatus `json:"status"`

	// Message is an optional descriptive message
	Message string `json:"message,omitempty"`

	// StartTime records when the operation started
	StartTime time.Time `json:"start_time"`

	// EndTime records when the operation ended (for terminal statuses)
	EndTime *time.Time `json:"end_time,omitempty"`

	// LastUpdate records when progress was last updated
	LastUpdate time.Time `json:"last_update"`

	// Error contains error information if status is failed
	Error error `json:"error,omitempty"`

	// Metadata contains additional operation-specific information
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Percentage returns the completion percentage (0-100).
// Returns -1 for indeterminate progress (Total == 0).
func (p *Progress) Percentage() float64 {
	if p.Total == 0 {
		return -1 // indeterminate
	}
	if p.Current >= p.Total {
		return 100.0
	}
	return float64(p.Current) / float64(p.Total) * 100.0
}

// Duration returns how long the operation has been running.
func (p *Progress) Duration() time.Duration {
	if p.EndTime != nil {
		return p.EndTime.Sub(p.StartTime)
	}
	return time.Since(p.StartTime)
}

// EstimatedTimeRemaining returns an estimate of remaining time.
// Returns -1 for indeterminate progress or if insufficient data.
func (p *Progress) EstimatedTimeRemaining() time.Duration {
	if p.Total == 0 || p.Current == 0 {
		return -1 // indeterminate
	}
	
	elapsed := p.Duration()
	remaining := p.Total - p.Current
	if remaining <= 0 {
		return 0
	}
	
	rate := float64(p.Current) / elapsed.Seconds()
	if rate <= 0 {
		return -1
	}
	
	return time.Duration(float64(remaining)/rate) * time.Second
}

// EventType represents different types of progress events.
type EventType string

const (
	// EventStarted indicates progress tracking started
	EventStarted EventType = "started"

	// EventUpdated indicates progress was updated
	EventUpdated EventType = "updated"

	// EventCompleted indicates operation completed successfully
	EventCompleted EventType = "completed"

	// EventFailed indicates operation failed
	EventFailed EventType = "failed"

	// EventCancelled indicates operation was cancelled
	EventCancelled EventType = "cancelled"
)

// ProgressEvent represents a progress update event.
type ProgressEvent struct {
	// Type identifies the event type
	Type EventType `json:"type"`

	// Timestamp records when the event occurred
	Timestamp time.Time `json:"timestamp"`

	// Progress contains the progress snapshot at event time
	Progress *Progress `json:"progress"`

	// Details contains event-specific additional information
	Details map[string]interface{} `json:"details,omitempty"`
}

// ProgressCallback is a function type for progress event notifications.
type ProgressCallback func(event *ProgressEvent)

// ProgressReporter defines the interface for reporting progress events.
// This allows for different reporting strategies (logging, metrics, etc.).
type ProgressReporter interface {
	// ReportEvent sends a progress event to the reporting system
	ReportEvent(event *ProgressEvent)

	// ReportProgress sends a progress snapshot
	ReportProgress(progress *Progress)

	// Close shuts down the reporter and flushes any pending reports
	Close() error
}

// BaseProgressTracker provides a basic implementation of ProgressTracker.
type BaseProgressTracker struct {
	mu        sync.RWMutex
	progress  *Progress
	callbacks []ProgressCallback
	nextID    int
}

// NewBaseProgressTracker creates a new BaseProgressTracker.
func NewBaseProgressTracker() *BaseProgressTracker {
	return &BaseProgressTracker{
		progress: &Progress{
			Status:     StatusPending,
			StartTime:  time.Now().UTC(),
			LastUpdate: time.Now().UTC(),
			Metadata:   make(map[string]interface{}),
		},
		callbacks: make([]ProgressCallback, 0),
	}
}

// Start initializes progress tracking with the total number of work units
func (bpt *BaseProgressTracker) Start(total int64) {
	bpt.mu.Lock()
	defer bpt.mu.Unlock()
	
	bpt.progress.Total = total
	bpt.progress.Status = StatusInProgress
	bpt.progress.StartTime = time.Now().UTC()
	bpt.progress.LastUpdate = bpt.progress.StartTime
	
	event := &ProgressEvent{
		Type:      EventStarted,
		Timestamp: time.Now().UTC(),
		Progress:  bpt.copyProgress(),
	}
	bpt.notifyCallbacks(event)
}

// Update sets the current progress value
func (bpt *BaseProgressTracker) Update(current int64) {
	bpt.UpdateWithMessage(current, "")
}

// UpdateWithMessage sets progress and updates the status message
func (bpt *BaseProgressTracker) UpdateWithMessage(current int64, message string) {
	bpt.mu.Lock()
	defer bpt.mu.Unlock()
	
	bpt.progress.Current = current
	bpt.progress.Message = message
	bpt.progress.LastUpdate = time.Now().UTC()
	
	event := &ProgressEvent{
		Type:      EventUpdated,
		Timestamp: time.Now().UTC(),
		Progress:  bpt.copyProgress(),
	}
	bpt.notifyCallbacks(event)
}

// Complete marks the operation as successfully completed
func (bpt *BaseProgressTracker) Complete() {
	bpt.mu.Lock()
	defer bpt.mu.Unlock()
	
	bpt.progress.Status = StatusCompleted
	endTime := time.Now().UTC()
	bpt.progress.EndTime = &endTime
	bpt.progress.LastUpdate = endTime
	
	event := &ProgressEvent{
		Type:      EventCompleted,
		Timestamp: endTime,
		Progress:  bpt.copyProgress(),
	}
	bpt.notifyCallbacks(event)
}

// Failed marks the operation as failed with the given error
func (bpt *BaseProgressTracker) Failed(err error) {
	bpt.mu.Lock()
	defer bpt.mu.Unlock()
	
	bpt.progress.Status = StatusFailed
	bpt.progress.Error = err
	endTime := time.Now().UTC()
	bpt.progress.EndTime = &endTime
	bpt.progress.LastUpdate = endTime
	
	event := &ProgressEvent{
		Type:      EventFailed,
		Timestamp: endTime,
		Progress:  bpt.copyProgress(),
	}
	bpt.notifyCallbacks(event)
}

// Cancel marks the operation as cancelled
func (bpt *BaseProgressTracker) Cancel() {
	bpt.mu.Lock()
	defer bpt.mu.Unlock()
	
	bpt.progress.Status = StatusCancelled
	endTime := time.Now().UTC()
	bpt.progress.EndTime = &endTime
	bpt.progress.LastUpdate = endTime
	
	event := &ProgressEvent{
		Type:      EventCancelled,
		Timestamp: endTime,
		Progress:  bpt.copyProgress(),
	}
	bpt.notifyCallbacks(event)
}

// GetProgress returns the current progress snapshot
func (bpt *BaseProgressTracker) GetProgress() *Progress {
	bpt.mu.RLock()
	defer bpt.mu.RUnlock()
	return bpt.copyProgress()
}

// Subscribe registers a callback for progress events
func (bpt *BaseProgressTracker) Subscribe(callback ProgressCallback) func() {
	bpt.mu.Lock()
	defer bpt.mu.Unlock()
	
	bpt.callbacks = append(bpt.callbacks, callback)
	callbackIndex := len(bpt.callbacks) - 1
	
	// Return unsubscribe function
	return func() {
		bpt.mu.Lock()
		defer bpt.mu.Unlock()
		
		// Remove callback by setting to nil (avoid slice reallocation)
		if callbackIndex < len(bpt.callbacks) {
			bpt.callbacks[callbackIndex] = nil
		}
	}
}

// copyProgress creates a deep copy of the current progress
func (bpt *BaseProgressTracker) copyProgress() *Progress {
	// Create metadata copy
	metadata := make(map[string]interface{})
	for k, v := range bpt.progress.Metadata {
		metadata[k] = v
	}
	
	return &Progress{
		Total:      bpt.progress.Total,
		Current:    bpt.progress.Current,
		Status:     bpt.progress.Status,
		Message:    bpt.progress.Message,
		StartTime:  bpt.progress.StartTime,
		EndTime:    bpt.progress.EndTime,
		LastUpdate: bpt.progress.LastUpdate,
		Error:      bpt.progress.Error,
		Metadata:   metadata,
	}
}

// notifyCallbacks sends event to all registered callbacks
func (bpt *BaseProgressTracker) notifyCallbacks(event *ProgressEvent) {
	for _, callback := range bpt.callbacks {
		if callback != nil {
			go callback(event) // Call asynchronously to avoid blocking
		}
	}
}