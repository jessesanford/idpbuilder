package progress

import (
	"context"
	"sync"
	"time"
)

// ProgressTracker defines the interface for tracking operation progress
type ProgressTracker interface {
	// Start begins progress tracking for an operation
	Start(operationID string, total int64) error
	
	// Update increments progress by the specified amount
	Update(operationID string, increment int64) error
	
	// SetProgress sets the absolute progress value
	SetProgress(operationID string, current int64) error
	
	// SetStage updates the current stage of the operation
	SetStage(operationID string, stage string) error
	
	// SetStatus updates the status of the operation
	SetStatus(operationID string, status string) error
	
	// Complete marks the operation as completed
	Complete(operationID string) error
	
	// Fail marks the operation as failed with an error
	Fail(operationID string, err error) error
	
	// Cancel marks the operation as cancelled
	Cancel(operationID string, reason string) error
	
	// GetProgress retrieves current progress information
	GetProgress(operationID string) (Progress, error)
	
	// ListActive returns all active operations
	ListActive() []string
	
	// Subscribe adds a callback for progress events
	Subscribe(callback ProgressCallback) error
	
	// Unsubscribe removes a progress callback
	Unsubscribe(callback ProgressCallback) error
	
	// Close shuts down the progress tracker
	Close() error
}

// Progress represents a snapshot of operation progress
type Progress struct {
	// OperationID uniquely identifies the operation
	OperationID string `json:"operation_id"`
	
	// Status indicates the current status
	Status string `json:"status"`
	
	// Stage indicates the current stage
	Stage string `json:"stage"`
	
	// Current progress value
	Current int64 `json:"current"`
	
	// Total expected progress value
	Total int64 `json:"total"`
	
	// Percentage completion (0-100)
	Percentage float64 `json:"percentage"`
	
	// StartTime when the operation began
	StartTime time.Time `json:"start_time"`
	
	// LastUpdate when progress was last updated
	LastUpdate time.Time `json:"last_update"`
	
	// ElapsedTime since operation started
	ElapsedTime time.Duration `json:"elapsed_time"`
	
	// EstimatedTimeRemaining based on current progress rate
	EstimatedTimeRemaining time.Duration `json:"estimated_time_remaining"`
	
	// Rate of progress (units per second)
	Rate float64 `json:"rate"`
	
	// Message provides additional context
	Message string `json:"message,omitempty"`
	
	// Error information if the operation failed
	Error string `json:"error,omitempty"`
	
	// CancelReason if the operation was cancelled
	CancelReason string `json:"cancel_reason,omitempty"`
	
	// Metadata stores additional operation-specific data
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ProgressEvent represents a progress update event
type ProgressEvent struct {
	// Type indicates the event type
	Type string `json:"type"`
	
	// OperationID identifies the operation
	OperationID string `json:"operation_id"`
	
	// Progress snapshot at the time of the event
	Progress Progress `json:"progress"`
	
	// Timestamp when the event occurred
	Timestamp time.Time `json:"timestamp"`
	
	// Message provides additional context
	Message string `json:"message,omitempty"`
}

// ProgressCallback is a function type for progress event notifications
type ProgressCallback func(event ProgressEvent)

// ProgressReporter defines interface for progress reporting
type ProgressReporter interface {
	// ReportProgress sends progress update to external systems
	ReportProgress(progress Progress) error
	
	// ReportEvent sends progress event to external systems
	ReportEvent(event ProgressEvent) error
	
	// Close shuts down the reporter
	Close() error
}

// BaseProgressTracker provides a concrete implementation of ProgressTracker
type BaseProgressTracker struct {
	// mu protects concurrent access
	mu sync.RWMutex
	
	// operations stores progress for active operations
	operations map[string]*progressState
	
	// callbacks stores registered event callbacks
	callbacks []ProgressCallback
	
	// callbacksMu protects callback list
	callbacksMu sync.RWMutex
	
	// reporter sends progress to external systems
	reporter ProgressReporter
	
	// updateInterval for automatic progress updates
	updateInterval time.Duration
	
	// staleThreshold for considering progress stale
	staleThreshold time.Duration
	
	// ctx for cancellation
	ctx context.Context
	
	// cancel function
	cancel context.CancelFunc
	
	// done channel signals completion
	done chan struct{}
}

// progressState tracks internal state for an operation
type progressState struct {
	progress     Progress
	lastRate     float64
	rateHistory  []float64
	updateTicker *time.Ticker
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker() *BaseProgressTracker {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &BaseProgressTracker{
		operations:     make(map[string]*progressState),
		callbacks:      make([]ProgressCallback, 0),
		updateInterval: DefaultUpdateInterval,
		staleThreshold: DefaultStaleThreshold,
		ctx:            ctx,
		cancel:         cancel,
		done:           make(chan struct{}),
	}
}

// NewProgressTrackerWithReporter creates a progress tracker with a reporter
func NewProgressTrackerWithReporter(reporter ProgressReporter) *BaseProgressTracker {
	tracker := NewProgressTracker()
	tracker.reporter = reporter
	return tracker
}

// Start begins progress tracking for an operation
func (t *BaseProgressTracker) Start(operationID string, total int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	now := time.Now()
	progress := Progress{
		OperationID:            operationID,
		Status:                 StatusInProgress,
		Stage:                  StageInitialization,
		Current:                0,
		Total:                  total,
		Percentage:             0.0,
		StartTime:              now,
		LastUpdate:             now,
		ElapsedTime:            0,
		EstimatedTimeRemaining: 0,
		Rate:                   0.0,
		Metadata:               make(map[string]interface{}),
	}
	
	state := &progressState{
		progress:    progress,
		rateHistory: make([]float64, 0, 10),
	}
	
	t.operations[operationID] = state
	
	// Send start event
	event := ProgressEvent{
		Type:        EventTypeStarted,
		OperationID: operationID,
		Progress:    progress,
		Timestamp:   now,
		Message:     "Operation started",
	}
	
	t.notifyCallbacks(event)
	
	if t.reporter != nil {
		t.reporter.ReportEvent(event)
	}
	
	return nil
}

// Update increments progress by the specified amount
func (t *BaseProgressTracker) Update(operationID string, increment int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	state, exists := t.operations[operationID]
	if !exists {
		return newProgressError("operation not found", operationID)
	}
	
	now := time.Now()
	state.progress.Current += increment
	
	return t.updateProgressState(state, now)
}

// SetProgress sets the absolute progress value
func (t *BaseProgressTracker) SetProgress(operationID string, current int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	state, exists := t.operations[operationID]
	if !exists {
		return newProgressError("operation not found", operationID)
	}
	
	now := time.Now()
	state.progress.Current = current
	
	return t.updateProgressState(state, now)
}

// SetStage updates the current stage of the operation
func (t *BaseProgressTracker) SetStage(operationID string, stage string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	state, exists := t.operations[operationID]
	if !exists {
		return newProgressError("operation not found", operationID)
	}
	
	oldStage := state.progress.Stage
	state.progress.Stage = stage
	state.progress.LastUpdate = time.Now()
	
	// Send stage change event
	event := ProgressEvent{
		Type:        EventTypeStageChanged,
		OperationID: operationID,
		Progress:    state.progress,
		Timestamp:   state.progress.LastUpdate,
		Message:     "Stage changed from " + oldStage + " to " + stage,
	}
	
	t.notifyCallbacks(event)
	
	if t.reporter != nil {
		t.reporter.ReportEvent(event)
	}
	
	return nil
}

// SetStatus updates the status of the operation
func (t *BaseProgressTracker) SetStatus(operationID string, status string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	state, exists := t.operations[operationID]
	if !exists {
		return newProgressError("operation not found", operationID)
	}
	
	state.progress.Status = status
	state.progress.LastUpdate = time.Now()
	
	// Send appropriate event based on status
	var eventType string
	var message string
	
	switch status {
	case StatusPaused:
		eventType = EventTypePaused
		message = "Operation paused"
	case StatusInProgress:
		eventType = EventTypeResumed
		message = "Operation resumed"
	case StatusCompleted:
		eventType = EventTypeCompleted
		message = "Operation completed"
	case StatusFailed:
		eventType = EventTypeFailed
		message = "Operation failed"
	case StatusCancelled:
		eventType = EventTypeCancelled
		message = "Operation cancelled"
	case StatusTimeout:
		eventType = EventTypeTimeout
		message = "Operation timed out"
	default:
		eventType = EventTypeProgress
		message = "Status updated to " + status
	}
	
	event := ProgressEvent{
		Type:        eventType,
		OperationID: operationID,
		Progress:    state.progress,
		Timestamp:   state.progress.LastUpdate,
		Message:     message,
	}
	
	t.notifyCallbacks(event)
	
	if t.reporter != nil {
		t.reporter.ReportEvent(event)
	}
	
	return nil
}

// Complete marks the operation as completed
func (t *BaseProgressTracker) Complete(operationID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	state, exists := t.operations[operationID]
	if !exists {
		return newProgressError("operation not found", operationID)
	}
	
	now := time.Now()
	state.progress.Status = StatusCompleted
	state.progress.Current = state.progress.Total
	state.progress.Percentage = 100.0
	state.progress.LastUpdate = now
	state.progress.ElapsedTime = now.Sub(state.progress.StartTime)
	state.progress.EstimatedTimeRemaining = 0
	
	// Send completion event
	event := ProgressEvent{
		Type:        EventTypeCompleted,
		OperationID: operationID,
		Progress:    state.progress,
		Timestamp:   now,
		Message:     "Operation completed successfully",
	}
	
	t.notifyCallbacks(event)
	
	if t.reporter != nil {
		t.reporter.ReportEvent(event)
	}
	
	// Clean up completed operation
	delete(t.operations, operationID)
	
	return nil
}

// Fail marks the operation as failed with an error
func (t *BaseProgressTracker) Fail(operationID string, err error) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	state, exists := t.operations[operationID]
	if !exists {
		return newProgressError("operation not found", operationID)
	}
	
	now := time.Now()
	state.progress.Status = StatusFailed
	state.progress.LastUpdate = now
	state.progress.ElapsedTime = now.Sub(state.progress.StartTime)
	state.progress.Error = err.Error()
	
	// Send failure event
	event := ProgressEvent{
		Type:        EventTypeFailed,
		OperationID: operationID,
		Progress:    state.progress,
		Timestamp:   now,
		Message:     "Operation failed: " + err.Error(),
	}
	
	t.notifyCallbacks(event)
	
	if t.reporter != nil {
		t.reporter.ReportEvent(event)
	}
	
	// Clean up failed operation
	delete(t.operations, operationID)
	
	return nil
}

// Cancel marks the operation as cancelled
func (t *BaseProgressTracker) Cancel(operationID string, reason string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	state, exists := t.operations[operationID]
	if !exists {
		return newProgressError("operation not found", operationID)
	}
	
	now := time.Now()
	state.progress.Status = StatusCancelled
	state.progress.LastUpdate = now
	state.progress.ElapsedTime = now.Sub(state.progress.StartTime)
	state.progress.CancelReason = reason
	
	// Send cancellation event
	event := ProgressEvent{
		Type:        EventTypeCancelled,
		OperationID: operationID,
		Progress:    state.progress,
		Timestamp:   now,
		Message:     "Operation cancelled: " + reason,
	}
	
	t.notifyCallbacks(event)
	
	if t.reporter != nil {
		t.reporter.ReportEvent(event)
	}
	
	// Clean up cancelled operation
	delete(t.operations, operationID)
	
	return nil
}

// GetProgress retrieves current progress information
func (t *BaseProgressTracker) GetProgress(operationID string) (Progress, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	state, exists := t.operations[operationID]
	if !exists {
		return Progress{}, newProgressError("operation not found", operationID)
	}
	
	return state.progress, nil
}

// ListActive returns all active operations
func (t *BaseProgressTracker) ListActive() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	operations := make([]string, 0, len(t.operations))
	for operationID := range t.operations {
		operations = append(operations, operationID)
	}
	
	return operations
}

// Subscribe adds a callback for progress events
func (t *BaseProgressTracker) Subscribe(callback ProgressCallback) error {
	t.callbacksMu.Lock()
	defer t.callbacksMu.Unlock()
	
	t.callbacks = append(t.callbacks, callback)
	return nil
}

// Unsubscribe removes a progress callback
func (t *BaseProgressTracker) Unsubscribe(callback ProgressCallback) error {
	t.callbacksMu.Lock()
	defer t.callbacksMu.Unlock()
	
	// Remove callback by comparing function addresses
	for i, cb := range t.callbacks {
		if &cb == &callback {
			t.callbacks = append(t.callbacks[:i], t.callbacks[i+1:]...)
			break
		}
	}
	
	return nil
}

// Close shuts down the progress tracker
func (t *BaseProgressTracker) Close() error {
	t.cancel()
	close(t.done)
	
	if t.reporter != nil {
		return t.reporter.Close()
	}
	
	return nil
}

// Helper methods

// updateProgressState updates the internal state and calculations
func (t *BaseProgressTracker) updateProgressState(state *progressState, now time.Time) error {
	// Update timing
	elapsed := now.Sub(state.progress.StartTime)
	state.progress.ElapsedTime = elapsed
	state.progress.LastUpdate = now
	
	// Calculate percentage
	if state.progress.Total > 0 {
		state.progress.Percentage = float64(state.progress.Current) / float64(state.progress.Total) * 100.0
		if state.progress.Percentage > 100.0 {
			state.progress.Percentage = 100.0
		}
	}
	
	// Calculate rate and ETA
	if elapsed.Seconds() > 0 {
		rate := float64(state.progress.Current) / elapsed.Seconds()
		state.progress.Rate = rate
		
		// Update rate history for smoothing
		state.rateHistory = append(state.rateHistory, rate)
		if len(state.rateHistory) > 10 {
			state.rateHistory = state.rateHistory[1:]
		}
		
		// Calculate average rate
		var avgRate float64
		for _, r := range state.rateHistory {
			avgRate += r
		}
		avgRate /= float64(len(state.rateHistory))
		
		// Estimate time remaining
		if avgRate > 0 && state.progress.Current < state.progress.Total {
			remaining := state.progress.Total - state.progress.Current
			state.progress.EstimatedTimeRemaining = time.Duration(float64(remaining)/avgRate) * time.Second
		}
	}
	
	// Send progress event
	event := ProgressEvent{
		Type:        EventTypeProgress,
		OperationID: state.progress.OperationID,
		Progress:    state.progress,
		Timestamp:   now,
		Message:     "Progress updated",
	}
	
	t.notifyCallbacks(event)
	
	if t.reporter != nil {
		t.reporter.ReportProgress(state.progress)
	}
	
	return nil
}

// notifyCallbacks sends events to all registered callbacks
func (t *BaseProgressTracker) notifyCallbacks(event ProgressEvent) {
	t.callbacksMu.RLock()
	defer t.callbacksMu.RUnlock()
	
	for _, callback := range t.callbacks {
		go callback(event) // Run callbacks asynchronously
	}
}

// newProgressError creates a progress-related error
func newProgressError(message, operationID string) error {
	return &ProgressError{
		Message:     message,
		OperationID: operationID,
	}
}

// ProgressError represents errors related to progress tracking
type ProgressError struct {
	Message     string
	OperationID string
}

// Error implements the error interface
func (e *ProgressError) Error() string {
	return "progress error for operation " + e.OperationID + ": " + e.Message
}