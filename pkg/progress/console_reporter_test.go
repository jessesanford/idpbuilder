package progress

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestConsoleReporter(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewConsoleReporterWithWriter(&buf)

	// Test full flow
	reporter.Start("Test operation")
	reporter.ReportProgress(50, 100, "Processing files")
	reporter.Complete("Done processing")

	output := buf.String()

	// Verify start message
	if !strings.Contains(output, "Starting: Test operation") {
		t.Errorf("Expected start message, got: %s", output)
	}

	// Verify progress percentage
	if !strings.Contains(output, "[50.0%]") {
		t.Errorf("Expected progress percentage, got: %s", output)
	}

	// Verify completion message
	if !strings.Contains(output, "Completed: Done processing") {
		t.Errorf("Expected completion message, got: %s", output)
	}
}

func TestConsoleReporterProgressWithoutTotal(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewConsoleReporterWithWriter(&buf)

	reporter.ReportProgress(42, 0, "Counting items")

	output := buf.String()
	if !strings.Contains(output, "[42]") {
		t.Errorf("Expected count format without percentage, got: %s", output)
	}
}

func TestConsoleReporterError(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewConsoleReporterWithWriter(&buf)

	testErr := errors.New("test error")
	reporter.Error(testErr)

	output := buf.String()
	if !strings.Contains(output, "Error: test error") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestConsoleReporterNilError(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewConsoleReporterWithWriter(&buf)

	reporter.Error(nil)

	output := buf.String()
	if output != "" {
		t.Errorf("Expected no output for nil error, got: %s", output)
	}
}