package progress

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewConsoleReporter(t *testing.T) {
	// Test with custom output
	buf := &bytes.Buffer{}
	reporter := NewConsoleReporter(buf, false)
	if reporter.output != buf {
		t.Error("NewConsoleReporter should use provided output")
	}
	if reporter.quiet {
		t.Error("NewConsoleReporter should not be quiet when quiet=false")
	}

	// Test with nil output (should default to os.Stdout)
	reporter = NewConsoleReporter(nil, true)
	if reporter.quiet != true {
		t.Error("NewConsoleReporter should be quiet when quiet=true")
	}
}

func TestConsoleReporter_StartOperation(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewConsoleReporter(buf, false)

	ctx := context.Background()
	ctx = reporter.StartOperation(ctx, OperationPush, "example.com/test:latest")

	output := buf.String()
	if !strings.Contains(output, "Starting push: example.com/test:latest") {
		t.Errorf("Expected start message, got: %s", output)
	}

	// Check that operation is tracked
	if len(reporter.active) != 1 {
		t.Error("Operation should be tracked in active map")
	}
}

func TestConsoleReporter_UpdateProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewConsoleReporter(buf, false)

	ctx := context.Background()
	reference := "example.com/test:latest"
	reporter.StartOperation(ctx, OperationPush, reference)

	progress := OperationProgress{
		OperationType: OperationPush,
		Reference:     reference,
		State:         StateInProgress,
		Current:       50,
		Total:         100,
		Message:       "Uploading layer",
		LastUpdate:    time.Now(),
	}

	buf.Reset() // Clear start message
	reporter.UpdateProgress(ctx, progress)

	output := buf.String()
	if !strings.Contains(output, "push") {
		t.Errorf("Expected operation type in progress, got: %s", output)
	}
	if !strings.Contains(output, "50.0%") {
		t.Errorf("Expected percentage in progress, got: %s", output)
	}
	if !strings.Contains(output, "50/100") {
		t.Errorf("Expected current/total in progress, got: %s", output)
	}
}

func TestConsoleReporter_CompleteOperation(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewConsoleReporter(buf, false)

	ctx := context.Background()
	reference := "example.com/test:latest"
	reporter.StartOperation(ctx, OperationPush, reference)

	buf.Reset() // Clear start message
	reporter.CompleteOperation(ctx, reference, "Upload successful")

	output := buf.String()
	if !strings.Contains(output, "Completed push") {
		t.Errorf("Expected completion message, got: %s", output)
	}
	if !strings.Contains(output, reference) {
		t.Errorf("Expected reference in completion, got: %s", output)
	}

	// Check that operation is no longer tracked
	if len(reporter.active) != 0 {
		t.Error("Completed operation should be removed from active map")
	}
}

func TestConsoleReporter_ErrorOperation(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewConsoleReporter(buf, false)

	ctx := context.Background()
	reference := "example.com/test:latest"
	testError := errors.New("network timeout")

	reporter.ErrorOperation(ctx, reference, testError)

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
	if !strings.Contains(output, reference) {
		t.Errorf("Expected reference in error, got: %s", output)
	}
	if !strings.Contains(output, "network timeout") {
		t.Errorf("Expected error details, got: %s", output)
	}
}

func TestConsoleReporter_CancelOperation(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewConsoleReporter(buf, false)

	ctx := context.Background()
	reference := "example.com/test:latest"
	reason := "user requested cancellation"

	reporter.CancelOperation(ctx, reference, reason)

	output := buf.String()
	if !strings.Contains(output, "Cancelled:") {
		t.Errorf("Expected cancellation message, got: %s", output)
	}
	if !strings.Contains(output, reference) {
		t.Errorf("Expected reference in cancellation, got: %s", output)
	}
	if !strings.Contains(output, reason) {
		t.Errorf("Expected reason in cancellation, got: %s", output)
	}
}

func TestConsoleReporter_QuietMode(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewConsoleReporter(buf, true) // Quiet mode

	ctx := context.Background()
	reference := "example.com/test:latest"

	// These operations should produce no output in quiet mode
	reporter.StartOperation(ctx, OperationPush, reference)
	reporter.UpdateProgress(ctx, OperationProgress{
		Reference: reference,
		Current:   50,
		Total:     100,
		Message:   "test",
	})
	reporter.CompleteOperation(ctx, reference, "done")
	reporter.CancelOperation(ctx, reference, "cancelled")

	output := buf.String()
	if output != "" {
		t.Errorf("Quiet mode should produce no output, got: %s", output)
	}
}

func TestConsoleReporter_IndeterminateProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewConsoleReporter(buf, false)

	ctx := context.Background()
	reference := "example.com/test:latest"

	progress := OperationProgress{
		OperationType: OperationList,
		Reference:     reference,
		State:         StateInProgress,
		Current:       0,
		Total:         0, // Unknown total
		Message:       "Fetching repository list",
	}

	reporter.UpdateProgress(ctx, progress)

	output := buf.String()
	if !strings.Contains(output, "[???]") {
		t.Errorf("Expected indeterminate progress indicator, got: %s", output)
	}
}