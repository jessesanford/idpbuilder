package progress

import (
	"context"
	"errors"
	"testing"
)

func TestSilentReporter(t *testing.T) {
	reporter := NewSilentReporter()
	ctx := context.Background()
	reference := "example.com/test:latest"
	testError := errors.New("test error")

	// All operations should complete without error or output
	ctx = reporter.StartOperation(ctx, OperationPush, reference)
	if ctx == nil {
		t.Error("StartOperation should return a context")
	}

	reporter.UpdateProgress(ctx, OperationProgress{
		Reference: reference,
		Current:   50,
		Total:     100,
	})

	reporter.CompleteOperation(ctx, reference, "success")
	reporter.ErrorOperation(ctx, reference, testError)
	reporter.CancelOperation(ctx, reference, "cancelled")

	// Silent reporter should not panic or error on any operation
}

func TestCallbackReporter(t *testing.T) {
	var receivedProgress []OperationProgress
	callback := func(progress OperationProgress) {
		receivedProgress = append(receivedProgress, progress)
	}

	reporter := NewCallbackReporter(callback)
	ctx := context.Background()
	reference := "example.com/test:latest"

	// Start operation
	reporter.StartOperation(ctx, OperationPush, reference)
	if len(receivedProgress) != 1 {
		t.Errorf("Expected 1 progress update, got %d", len(receivedProgress))
	}
	if receivedProgress[0].State != StateStarted {
		t.Errorf("Expected StateStarted, got %v", receivedProgress[0].State)
	}

	// Update progress
	progress := OperationProgress{
		OperationType: OperationPush,
		Reference:     reference,
		State:         StateInProgress,
		Current:       75,
		Total:         100,
		Message:       "Uploading",
	}
	reporter.UpdateProgress(ctx, progress)
	if len(receivedProgress) != 2 {
		t.Errorf("Expected 2 progress updates, got %d", len(receivedProgress))
	}
	if receivedProgress[1].Current != 75 {
		t.Errorf("Expected Current=75, got %d", receivedProgress[1].Current)
	}

	// Complete operation
	reporter.CompleteOperation(ctx, reference, "Upload complete")
	if len(receivedProgress) != 3 {
		t.Errorf("Expected 3 progress updates, got %d", len(receivedProgress))
	}
	if receivedProgress[2].State != StateCompleted {
		t.Errorf("Expected StateCompleted, got %v", receivedProgress[2].State)
	}

	// Error operation
	testError := errors.New("network error")
	reporter.ErrorOperation(ctx, reference, testError)
	if len(receivedProgress) != 4 {
		t.Errorf("Expected 4 progress updates, got %d", len(receivedProgress))
	}
	if receivedProgress[3].State != StateError {
		t.Errorf("Expected StateError, got %v", receivedProgress[3].State)
	}
	if receivedProgress[3].Error != testError {
		t.Errorf("Expected error to be passed through, got %v", receivedProgress[3].Error)
	}

	// Cancel operation
	reporter.CancelOperation(ctx, reference, "User cancelled")
	if len(receivedProgress) != 5 {
		t.Errorf("Expected 5 progress updates, got %d", len(receivedProgress))
	}
	if receivedProgress[4].State != StateCancelled {
		t.Errorf("Expected StateCancelled, got %v", receivedProgress[4].State)
	}
}

func TestCallbackReporter_NilCallback(t *testing.T) {
	reporter := NewCallbackReporter(nil)
	ctx := context.Background()
	reference := "example.com/test:latest"

	// Operations with nil callback should not panic
	reporter.StartOperation(ctx, OperationPush, reference)
	reporter.UpdateProgress(ctx, OperationProgress{Reference: reference})
	reporter.CompleteOperation(ctx, reference, "done")
	reporter.ErrorOperation(ctx, reference, errors.New("error"))
	reporter.CancelOperation(ctx, reference, "cancelled")
}