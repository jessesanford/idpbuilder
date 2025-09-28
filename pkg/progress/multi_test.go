package progress

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestMultiReporter(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	reporter1 := NewConsoleReporterWithWriter(&buf1)
	reporter2 := NewConsoleReporterWithWriter(&buf2)

	multiReporter := NewMultiReporter(reporter1, reporter2)

	// Test that all operations are forwarded to all reporters
	multiReporter.Start("Multi test")
	multiReporter.ReportProgress(30, 100, "Working")
	multiReporter.Complete("All done")

	output1 := buf1.String()
	output2 := buf2.String()

	// Both outputs should be identical
	if output1 != output2 {
		t.Errorf("Multi-reporter outputs don't match:\n1: %s\n2: %s", output1, output2)
	}

	// Verify content is present
	if !strings.Contains(output1, "Starting: Multi test") {
		t.Errorf("Expected start message in output1: %s", output1)
	}
	if !strings.Contains(output1, "[30.0%]") {
		t.Errorf("Expected progress in output1: %s", output1)
	}
	if !strings.Contains(output1, "Completed: All done") {
		t.Errorf("Expected completion in output1: %s", output1)
	}
}

func TestMultiReporterError(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	reporter1 := NewConsoleReporterWithWriter(&buf1)
	reporter2 := NewConsoleReporterWithWriter(&buf2)

	multiReporter := NewMultiReporter(reporter1, reporter2)

	testErr := errors.New("multi error")
	multiReporter.Error(testErr)

	output1 := buf1.String()
	output2 := buf2.String()

	if !strings.Contains(output1, "Error: multi error") {
		t.Errorf("Expected error in output1: %s", output1)
	}
	if !strings.Contains(output2, "Error: multi error") {
		t.Errorf("Expected error in output2: %s", output2)
	}
}

func TestMultiReporterWithNilReporters(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewConsoleReporterWithWriter(&buf)

	// Include nil reporters - they should be filtered out
	multiReporter := NewMultiReporter(reporter, nil, nil)

	multiReporter.Start("Nil test")
	multiReporter.Complete("Done")

	output := buf.String()
	if !strings.Contains(output, "Starting: Nil test") {
		t.Errorf("Expected start message despite nil reporters: %s", output)
	}
}

func TestMultiReporterEmpty(t *testing.T) {
	// Should not panic with no reporters
	multiReporter := NewMultiReporter()
	multiReporter.Start("Empty test")
	multiReporter.ReportProgress(10, 20, "test")
	multiReporter.Complete("test")
	multiReporter.Error(errors.New("test"))
	// Test passes if no panic occurs
}