package progress

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestJSONReporter(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewJSONReporter(&buf)

	// Test full flow
	reporter.Start("Test operation")
	reporter.ReportProgress(25, 100, "Processing")
	reporter.Complete("Finished")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Fatalf("Expected 3 JSON lines, got %d", len(lines))
	}

	// Parse start event
	var startEvent map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &startEvent); err != nil {
		t.Errorf("Failed to parse start event: %v", err)
	}
	if startEvent["event"] != "start" || startEvent["message"] != "Test operation" {
		t.Errorf("Invalid start event: %v", startEvent)
	}

	// Parse progress event
	var progressEvent map[string]interface{}
	if err := json.Unmarshal([]byte(lines[1]), &progressEvent); err != nil {
		t.Errorf("Failed to parse progress event: %v", err)
	}
	if progressEvent["event"] != "progress" || progressEvent["current"].(float64) != 25 {
		t.Errorf("Invalid progress event: %v", progressEvent)
	}

	// Parse complete event
	var completeEvent map[string]interface{}
	if err := json.Unmarshal([]byte(lines[2]), &completeEvent); err != nil {
		t.Errorf("Failed to parse complete event: %v", err)
	}
	if completeEvent["event"] != "complete" || completeEvent["message"] != "Finished" {
		t.Errorf("Invalid complete event: %v", completeEvent)
	}
}

func TestJSONReporterError(t *testing.T) {
	var buf bytes.Buffer
	reporter := NewJSONReporter(&buf)

	testErr := errors.New("test error")
	reporter.Error(testErr)

	output := strings.TrimSpace(buf.String())
	var errorEvent map[string]interface{}
	if err := json.Unmarshal([]byte(output), &errorEvent); err != nil {
		t.Errorf("Failed to parse error event: %v", err)
	}

	if errorEvent["event"] != "error" || errorEvent["message"] != "test error" {
		t.Errorf("Invalid error event: %v", errorEvent)
	}
}

func TestJSONReporterNilWriter(t *testing.T) {
	// Should not panic with nil writer
	reporter := NewJSONReporter(nil)
	reporter.Start("test")
	// Test passes if no panic occurs
}