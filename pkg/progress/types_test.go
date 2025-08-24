package progress

import (
	"testing"
	"time"
)

func TestProgressStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   ProgressStatus
		expected string
	}{
		{"pending", StatusPending, "pending"},
		{"in_progress", StatusInProgress, "in_progress"},
		{"completed", StatusCompleted, "completed"},
		{"failed", StatusFailed, "failed"},
		{"cancelled", StatusCancelled, "cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("ProgressStatus.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProgressStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		name     string
		status   ProgressStatus
		expected bool
	}{
		{"pending is not terminal", StatusPending, false},
		{"in_progress is not terminal", StatusInProgress, false},
		{"completed is terminal", StatusCompleted, true},
		{"failed is terminal", StatusFailed, true},
		{"cancelled is terminal", StatusCancelled, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsTerminal(); got != tt.expected {
				t.Errorf("ProgressStatus.IsTerminal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProgress_Percentage(t *testing.T) {
	tests := []struct {
		name     string
		progress *Progress
		expected float64
	}{
		{
			name: "zero total (indeterminate)",
			progress: &Progress{
				Total:   0,
				Current: 50,
			},
			expected: -1,
		},
		{
			name: "50% complete",
			progress: &Progress{
				Total:   100,
				Current: 50,
			},
			expected: 50.0,
		},
		{
			name: "100% complete",
			progress: &Progress{
				Total:   100,
				Current: 100,
			},
			expected: 100.0,
		},
		{
			name: "over 100% complete",
			progress: &Progress{
				Total:   100,
				Current: 150,
			},
			expected: 100.0,
		},
		{
			name: "zero progress",
			progress: &Progress{
				Total:   100,
				Current: 0,
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.progress.Percentage(); got != tt.expected {
				t.Errorf("Progress.Percentage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProgress_Duration(t *testing.T) {
	start := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 12, 0, 30, 0, time.UTC)

	tests := []struct {
		name     string
		progress *Progress
		expected time.Duration
	}{
		{
			name: "completed operation",
			progress: &Progress{
				StartTime: start,
				EndTime:   &end,
			},
			expected: 30 * time.Second,
		},
		{
			name: "ongoing operation",
			progress: func() *Progress {
				// Create progress with recent start time
				return &Progress{
					StartTime: time.Now().Add(-5 * time.Second),
					EndTime:   nil,
				}
			}(),
			expected: 5 * time.Second, // Approximate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.progress.Duration()
			if tt.name == "completed operation" {
				if got != tt.expected {
					t.Errorf("Progress.Duration() = %v, want %v", got, tt.expected)
				}
			} else {
				// For ongoing operation, check within reasonable range
				if got < 4*time.Second || got > 6*time.Second {
					t.Errorf("Progress.Duration() = %v, expected around %v", got, tt.expected)
				}
			}
		})
	}
}

func TestProgress_EstimatedTimeRemaining(t *testing.T) {
	tests := []struct {
		name     string
		progress *Progress
		expected time.Duration
	}{
		{
			name: "indeterminate progress (zero total)",
			progress: &Progress{
				Total:     0,
				Current:   50,
				StartTime: time.Now().Add(-10 * time.Second),
			},
			expected: -1,
		},
		{
			name: "no progress yet",
			progress: &Progress{
				Total:     100,
				Current:   0,
				StartTime: time.Now().Add(-10 * time.Second),
			},
			expected: -1,
		},
		{
			name: "completed",
			progress: &Progress{
				Total:     100,
				Current:   100,
				StartTime: time.Now().Add(-10 * time.Second),
			},
			expected: 0,
		},
		{
			name: "half complete, estimate remaining",
			progress: &Progress{
				Total:     100,
				Current:   50,
				StartTime: time.Now().Add(-10 * time.Second),
			},
			expected: 10 * time.Second, // Should take same time for remaining half
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.progress.EstimatedTimeRemaining()
			if tt.expected == -1 {
				if got != -1 {
					t.Errorf("Progress.EstimatedTimeRemaining() = %v, want %v", got, tt.expected)
				}
			} else if tt.expected == 0 {
				if got != 0 {
					t.Errorf("Progress.EstimatedTimeRemaining() = %v, want %v", got, tt.expected)
				}
			} else {
				// Allow some tolerance for time-based estimates
				tolerance := 2 * time.Second
				if got < tt.expected-tolerance || got > tt.expected+tolerance {
					t.Errorf("Progress.EstimatedTimeRemaining() = %v, want %v (±%v)", got, tt.expected, tolerance)
				}
			}
		})
	}
}

func TestNewBaseProgressTracker(t *testing.T) {
	tracker := NewBaseProgressTracker()

	if tracker == nil {
		t.Fatal("NewBaseProgressTracker() returned nil")
	}

	progress := tracker.GetProgress()
	if progress == nil {
		t.Fatal("GetProgress() returned nil")
	}

	if progress.Status != StatusPending {
		t.Errorf("initial status = %v, want %v", progress.Status, StatusPending)
	}

	if progress.Total != 0 {
		t.Errorf("initial total = %v, want %v", progress.Total, 0)
	}

	if progress.Current != 0 {
		t.Errorf("initial current = %v, want %v", progress.Current, 0)
	}

	if progress.Metadata == nil {
		t.Error("metadata should be initialized")
	}
}

func TestProgressEvent(t *testing.T) {
	timestamp := time.Now().UTC()
	progress := &Progress{
		Total:   100,
		Current: 50,
		Status:  StatusInProgress,
	}

	event := &ProgressEvent{
		Type:      EventUpdated,
		Timestamp: timestamp,
		Progress:  progress,
		Details: map[string]interface{}{
			"message": "processing",
		},
	}

	if event.Type != EventUpdated {
		t.Errorf("event type = %v, want %v", event.Type, EventUpdated)
	}

	if !event.Timestamp.Equal(timestamp) {
		t.Errorf("event timestamp = %v, want %v", event.Timestamp, timestamp)
	}

	if event.Progress != progress {
		t.Errorf("event progress = %v, want %v", event.Progress, progress)
	}

	if msg, ok := event.Details["message"]; !ok || msg != "processing" {
		t.Errorf("event details = %v, want message=processing", event.Details)
	}
}

func TestProgressCallback(t *testing.T) {
	var receivedEvent *ProgressEvent
	callback := func(event *ProgressEvent) {
		receivedEvent = event
	}

	// Test that callback type matches expected signature
	event := &ProgressEvent{
		Type:      EventStarted,
		Timestamp: time.Now(),
		Progress: &Progress{
			Status: StatusInProgress,
		},
	}

	callback(event)

	if receivedEvent == nil {
		t.Error("callback was not called")
	}

	if receivedEvent.Type != EventStarted {
		t.Errorf("callback received event type = %v, want %v", receivedEvent.Type, EventStarted)
	}
}

func TestEventType_Values(t *testing.T) {
	events := []EventType{
		EventStarted,
		EventUpdated,
		EventCompleted,
		EventFailed,
		EventCancelled,
	}

	expectedStrings := []string{
		"started",
		"updated",
		"completed",
		"failed",
		"cancelled",
	}

	for i, event := range events {
		if string(event) != expectedStrings[i] {
			t.Errorf("event %d = %v, want %v", i, string(event), expectedStrings[i])
		}
	}
}

func TestProgressTrackerInterface(t *testing.T) {
	// Test that BaseProgressTracker implements ProgressTracker interface
	var tracker ProgressTracker = NewBaseProgressTracker()
	
	if tracker == nil {
		t.Fatal("BaseProgressTracker does not implement ProgressTracker interface")
	}

	// Test that we can call interface methods
	progress := tracker.GetProgress()
	if progress == nil {
		t.Error("ProgressTracker.GetProgress() returned nil")
	}
}

// mockReporter is a test implementation of ProgressReporter
type mockReporter struct{}

func (m *mockReporter) ReportEvent(event *ProgressEvent) {}
func (m *mockReporter) ReportProgress(progress *Progress) {}
func (m *mockReporter) Close() error { return nil }

func TestProgressReporterInterface(t *testing.T) {
	// Define a mock reporter to test the interface

	// Test that mock implements ProgressReporter
	var reporter ProgressReporter = &mockReporter{}
	
	if reporter == nil {
		t.Fatal("mockReporter does not implement ProgressReporter interface")
	}

	// Test interface methods
	reporter.ReportEvent(&ProgressEvent{})
	reporter.ReportProgress(&Progress{})
	
	if err := reporter.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}