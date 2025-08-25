package progress

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNewProgressTracker(t *testing.T) {
	tracker := NewProgressTracker()
	
	if tracker == nil {
		t.Fatal("Expected tracker to be created")
	}
	
	if tracker.operations == nil {
		t.Error("Expected operations map to be initialized")
	}
	
	if tracker.callbacks == nil {
		t.Error("Expected callbacks slice to be initialized")
	}
	
	if tracker.updateInterval != DefaultUpdateInterval {
		t.Errorf("Expected update interval %v, got %v", DefaultUpdateInterval, tracker.updateInterval)
	}
	
	tracker.Close()
}

func TestProgressTracker_Start(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	progress, err := tracker.GetProgress(operationID)
	if err != nil {
		t.Fatalf("Unexpected error getting progress: %v", err)
	}
	
	if progress.OperationID != operationID {
		t.Errorf("Expected operation ID %q, got %q", operationID, progress.OperationID)
	}
	
	if progress.Total != total {
		t.Errorf("Expected total %d, got %d", total, progress.Total)
	}
	
	if progress.Current != 0 {
		t.Errorf("Expected current 0, got %d", progress.Current)
	}
	
	if progress.Status != StatusInProgress {
		t.Errorf("Expected status %q, got %q", StatusInProgress, progress.Status)
	}
	
	if progress.Stage != StageInitialization {
		t.Errorf("Expected stage %q, got %q", StageInitialization, progress.Stage)
	}
}

func TestProgressTracker_Update(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	err = tracker.Update(operationID, 25)
	if err != nil {
		t.Fatalf("Unexpected error updating progress: %v", err)
	}
	
	progress, err := tracker.GetProgress(operationID)
	if err != nil {
		t.Fatalf("Unexpected error getting progress: %v", err)
	}
	
	if progress.Current != 25 {
		t.Errorf("Expected current 25, got %d", progress.Current)
	}
	
	if progress.Percentage != 25.0 {
		t.Errorf("Expected percentage 25.0, got %f", progress.Percentage)
	}
}

func TestProgressTracker_SetProgress(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	err = tracker.SetProgress(operationID, 50)
	if err != nil {
		t.Fatalf("Unexpected error setting progress: %v", err)
	}
	
	progress, err := tracker.GetProgress(operationID)
	if err != nil {
		t.Fatalf("Unexpected error getting progress: %v", err)
	}
	
	if progress.Current != 50 {
		t.Errorf("Expected current 50, got %d", progress.Current)
	}
	
	if progress.Percentage != 50.0 {
		t.Errorf("Expected percentage 50.0, got %f", progress.Percentage)
	}
}

func TestProgressTracker_SetStage(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	err = tracker.SetStage(operationID, StageDownload)
	if err != nil {
		t.Fatalf("Unexpected error setting stage: %v", err)
	}
	
	progress, err := tracker.GetProgress(operationID)
	if err != nil {
		t.Fatalf("Unexpected error getting progress: %v", err)
	}
	
	if progress.Stage != StageDownload {
		t.Errorf("Expected stage %q, got %q", StageDownload, progress.Stage)
	}
}

func TestProgressTracker_SetStatus(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	err = tracker.SetStatus(operationID, StatusPaused)
	if err != nil {
		t.Fatalf("Unexpected error setting status: %v", err)
	}
	
	progress, err := tracker.GetProgress(operationID)
	if err != nil {
		t.Fatalf("Unexpected error getting progress: %v", err)
	}
	
	if progress.Status != StatusPaused {
		t.Errorf("Expected status %q, got %q", StatusPaused, progress.Status)
	}
}

func TestProgressTracker_Complete(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	err = tracker.Complete(operationID)
	if err != nil {
		t.Fatalf("Unexpected error completing operation: %v", err)
	}
	
	// Operation should be removed after completion
	_, err = tracker.GetProgress(operationID)
	if err == nil {
		t.Error("Expected error getting progress for completed operation")
	}
}

func TestProgressTracker_Fail(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	testError := errors.New("test failure")
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	err = tracker.Fail(operationID, testError)
	if err != nil {
		t.Fatalf("Unexpected error failing operation: %v", err)
	}
	
	// Operation should be removed after failure
	_, err = tracker.GetProgress(operationID)
	if err == nil {
		t.Error("Expected error getting progress for failed operation")
	}
}

func TestProgressTracker_Cancel(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	reason := "user cancelled"
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	err = tracker.Cancel(operationID, reason)
	if err != nil {
		t.Fatalf("Unexpected error cancelling operation: %v", err)
	}
	
	// Operation should be removed after cancellation
	_, err = tracker.GetProgress(operationID)
	if err == nil {
		t.Error("Expected error getting progress for cancelled operation")
	}
}

func TestProgressTracker_ListActive(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operations := []string{"op1", "op2", "op3"}
	
	for _, op := range operations {
		err := tracker.Start(op, 100)
		if err != nil {
			t.Fatalf("Unexpected error starting operation %s: %v", op, err)
		}
	}
	
	active := tracker.ListActive()
	if len(active) != len(operations) {
		t.Errorf("Expected %d active operations, got %d", len(operations), len(active))
	}
	
	// Verify all operations are in the list
	for _, op := range operations {
		found := false
		for _, activeOp := range active {
			if activeOp == op {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Operation %s not found in active list", op)
		}
	}
}

func TestProgressTracker_Callbacks(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	eventCh := make(chan ProgressEvent, 10)
	callback := func(event ProgressEvent) {
		eventCh <- event
	}
	
	err := tracker.Subscribe(callback)
	if err != nil {
		t.Fatalf("Unexpected error subscribing: %v", err)
	}
	
	operationID := "test-operation"
	total := int64(100)
	
	err = tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	// Wait for start event
	select {
	case event := <-eventCh:
		if event.Type != EventTypeStarted {
			t.Errorf("Expected event type %q, got %q", EventTypeStarted, event.Type)
		}
		if event.OperationID != operationID {
			t.Errorf("Expected operation ID %q, got %q", operationID, event.OperationID)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for start event")
	}
	
	err = tracker.Update(operationID, 50)
	if err != nil {
		t.Fatalf("Unexpected error updating progress: %v", err)
	}
	
	// Wait for progress event
	select {
	case event := <-eventCh:
		if event.Type != EventTypeProgress {
			t.Errorf("Expected event type %q, got %q", EventTypeProgress, event.Type)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for progress event")
	}
}

func TestProgressTracker_ConcurrentAccess(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	numOperations := 10
	numUpdates := 100
	
	var wg sync.WaitGroup
	
	// Start multiple operations concurrently
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(opIndex int) {
			defer wg.Done()
			
			operationID := "operation-" + string(rune('0'+opIndex))
			err := tracker.Start(operationID, int64(numUpdates))
			if err != nil {
				t.Errorf("Error starting operation %s: %v", operationID, err)
				return
			}
			
			// Update progress multiple times
			for j := 0; j < numUpdates; j++ {
				err := tracker.Update(operationID, 1)
				if err != nil {
					t.Errorf("Error updating operation %s: %v", operationID, err)
					return
				}
			}
			
			err = tracker.Complete(operationID)
			if err != nil {
				t.Errorf("Error completing operation %s: %v", operationID, err)
			}
		}(i)
	}
	
	wg.Wait()
	
	// All operations should be completed and removed
	active := tracker.ListActive()
	if len(active) != 0 {
		t.Errorf("Expected 0 active operations, got %d", len(active))
	}
}

func TestProgressTracker_RateCalculation(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	// Sleep briefly to ensure time passes
	time.Sleep(10 * time.Millisecond)
	
	err = tracker.Update(operationID, 50)
	if err != nil {
		t.Fatalf("Unexpected error updating progress: %v", err)
	}
	
	progress, err := tracker.GetProgress(operationID)
	if err != nil {
		t.Fatalf("Unexpected error getting progress: %v", err)
	}
	
	// Rate should be positive since we've made progress
	if progress.Rate <= 0 {
		t.Errorf("Expected positive rate, got %f", progress.Rate)
	}
	
	// ETA should be calculated (may be 0 if calculation conditions aren't met)
	if progress.EstimatedTimeRemaining < 0 {
		t.Errorf("Expected non-negative ETA, got %v", progress.EstimatedTimeRemaining)
	}
}

func TestProgressTracker_EdgeCases(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	t.Run("update non-existent operation", func(t *testing.T) {
		err := tracker.Update("non-existent", 10)
		if err == nil {
			t.Error("Expected error updating non-existent operation")
		}
	})
	
	t.Run("get progress for non-existent operation", func(t *testing.T) {
		_, err := tracker.GetProgress("non-existent")
		if err == nil {
			t.Error("Expected error getting progress for non-existent operation")
		}
	})
	
	t.Run("complete non-existent operation", func(t *testing.T) {
		err := tracker.Complete("non-existent")
		if err == nil {
			t.Error("Expected error completing non-existent operation")
		}
	})
	
	t.Run("progress over 100%", func(t *testing.T) {
		operationID := "over-100"
		err := tracker.Start(operationID, 100)
		if err != nil {
			t.Fatalf("Unexpected error starting operation: %v", err)
		}
		
		err = tracker.SetProgress(operationID, 150) // Over total
		if err != nil {
			t.Fatalf("Unexpected error setting progress: %v", err)
		}
		
		progress, err := tracker.GetProgress(operationID)
		if err != nil {
			t.Fatalf("Unexpected error getting progress: %v", err)
		}
		
		// Percentage should be capped at 100
		if progress.Percentage > 100.0 {
			t.Errorf("Expected percentage <= 100, got %f", progress.Percentage)
		}
	})
	
	t.Run("zero total progress", func(t *testing.T) {
		operationID := "zero-total"
		err := tracker.Start(operationID, 0)
		if err != nil {
			t.Fatalf("Unexpected error starting operation: %v", err)
		}
		
		progress, err := tracker.GetProgress(operationID)
		if err != nil {
			t.Fatalf("Unexpected error getting progress: %v", err)
		}
		
		// Should handle zero total gracefully
		if progress.Percentage != 0.0 {
			t.Errorf("Expected percentage 0.0 for zero total, got %f", progress.Percentage)
		}
	})
}

func TestProgress_Metadata(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	progress, err := tracker.GetProgress(operationID)
	if err != nil {
		t.Fatalf("Unexpected error getting progress: %v", err)
	}
	
	// Metadata map should be initialized
	if progress.Metadata == nil {
		t.Error("Expected metadata map to be initialized")
	}
	
	// Should start with empty metadata
	if len(progress.Metadata) != 0 {
		t.Errorf("Expected empty metadata, got %d items", len(progress.Metadata))
	}
}

func TestProgressEvent_Fields(t *testing.T) {
	tracker := NewProgressTracker()
	defer tracker.Close()
	
	eventCh := make(chan ProgressEvent, 1)
	callback := func(event ProgressEvent) {
		eventCh <- event
	}
	
	err := tracker.Subscribe(callback)
	if err != nil {
		t.Fatalf("Unexpected error subscribing: %v", err)
	}
	
	operationID := "test-operation"
	total := int64(100)
	
	before := time.Now()
	err = tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	after := time.Now()
	
	select {
	case event := <-eventCh:
		if event.Type == "" {
			t.Error("Event type should not be empty")
		}
		
		if event.OperationID != operationID {
			t.Errorf("Expected operation ID %q, got %q", operationID, event.OperationID)
		}
		
		if event.Timestamp.Before(before) || event.Timestamp.After(after) {
			t.Error("Event timestamp should be within test execution window")
		}
		
		if event.Message == "" {
			t.Error("Event message should not be empty")
		}
		
	case <-time.After(time.Second):
		t.Error("Timeout waiting for event")
	}
}

func TestProgressError(t *testing.T) {
	err := &ProgressError{
		Message:     "test error",
		OperationID: "test-operation",
	}
	
	expectedError := "progress error for operation test-operation: test error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, err.Error())
	}
}

func TestNewProgressTrackerWithReporter(t *testing.T) {
	mockReporter := &MockProgressReporter{}
	tracker := NewProgressTrackerWithReporter(mockReporter)
	defer tracker.Close()
	
	if tracker.reporter != mockReporter {
		t.Error("Expected reporter to be set")
	}
	
	// Test that reporter is called
	operationID := "test-operation"
	total := int64(100)
	
	err := tracker.Start(operationID, total)
	if err != nil {
		t.Fatalf("Unexpected error starting operation: %v", err)
	}
	
	// Check that reporter was called
	if len(mockReporter.events) == 0 {
		t.Error("Expected reporter to receive events")
	}
	
	err = tracker.Update(operationID, 50)
	if err != nil {
		t.Fatalf("Unexpected error updating progress: %v", err)
	}
	
	if len(mockReporter.progress) == 0 {
		t.Error("Expected reporter to receive progress updates")
	}
}

// MockProgressReporter for testing
type MockProgressReporter struct {
	progress []Progress
	events   []ProgressEvent
	mu       sync.Mutex
}

func (m *MockProgressReporter) ReportProgress(progress Progress) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.progress = append(m.progress, progress)
	return nil
}

func (m *MockProgressReporter) ReportEvent(event ProgressEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	return nil
}

func (m *MockProgressReporter) Close() error {
	return nil
}