package cli

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgressBar(t *testing.T) {
	t.Run("create progress bar", func(t *testing.T) {
		pb := NewProgressBar("Testing")
		require.NotNil(t, pb)

		// Progress bar should implement the interface
		var _ ProgressReporter = pb
	})

	t.Run("update message", func(t *testing.T) {
		pb := NewProgressBar("Initial")

		// Update message should not panic
		pb.UpdateMessage("Updated message")
		pb.UpdateMessage("Another update")
	})

	t.Run("update progress", func(t *testing.T) {
		pb := NewProgressBar("Progress test")

		// Update progress should not panic
		pb.UpdateProgress(0, 100)
		pb.UpdateProgress(50, 100)
		pb.UpdateProgress(100, 100)
	})

	t.Run("finish progress", func(t *testing.T) {
		pb := NewProgressBar("Finish test")

		// Finish should not panic
		pb.Finish()

		// Multiple calls to finish should be safe
		pb.Finish()
		pb.Finish()
	})
}

func TestQuietProgress(t *testing.T) {
	t.Run("create quiet progress", func(t *testing.T) {
		qp := NewQuietProgress("Quiet test")
		require.NotNil(t, qp)

		// Should implement the interface
		var _ ProgressReporter = qp
	})

	t.Run("quiet progress operations", func(t *testing.T) {
		qp := NewQuietProgress("Initial")

		// All operations should be no-op or safe
		qp.UpdateMessage("Updated")
		qp.UpdateProgress(50, 100)
		qp.Finish()
	})
}

func TestMultiProgress(t *testing.T) {
	t.Run("create multi progress", func(t *testing.T) {
		mp := NewMultiProgress()
		require.NotNil(t, mp)
	})

	t.Run("add and use progress bars", func(t *testing.T) {
		mp := NewMultiProgress()

		// Add progress bars
		bar1 := mp.AddBar("task1", "Task 1")
		bar2 := mp.AddBar("task2", "Task 2")

		require.NotNil(t, bar1)
		require.NotNil(t, bar2)

		// Use the bars
		bar1.UpdateMessage("Task 1 working")
		bar1.UpdateProgress(25, 100)

		bar2.UpdateMessage("Task 2 working")
		bar2.UpdateProgress(75, 100)

		// Finish them
		bar1.Finish()
		bar2.Finish()
	})
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"small bytes", 512, "512 B"},
		{"kilobytes", 2048, "2.0 KB"},
		{"megabytes", 1024 * 1024 * 5, "5.0 MB"},
		{"gigabytes", 1024 * 1024 * 1024 * 2, "2.0 GB"},
		{"zero bytes", 0, "0 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProgressBarWithBuffer(t *testing.T) {
	t.Run("progress bar output", func(t *testing.T) {
		var buf bytes.Buffer

		// Create a progress bar that writes to our buffer
		pb := &ProgressBar{
			writer:       &buf,
			message:      "Test",
			startTime:    time.Now(),
			spinnerChars: []string{"⠋", "⠙"},
		}

		pb.render()
		output := buf.String()

		// Should contain the message
		assert.Contains(t, output, "Test")
	})

	t.Run("progress bar with percentage", func(t *testing.T) {
		var buf bytes.Buffer

		pb := &ProgressBar{
			writer:       &buf,
			message:      "Progress",
			current:      50,
			total:        100,
			startTime:    time.Now(),
			spinnerChars: []string{"⠋"},
		}

		pb.render()
		output := buf.String()

		// Should contain percentage
		assert.Contains(t, output, "50%")
		assert.Contains(t, output, "Progress")
	})
}

func TestProgressReporterInterface(t *testing.T) {
	t.Run("all implementations satisfy interface", func(t *testing.T) {
		var reporters []ProgressReporter

		// All these should implement ProgressReporter
		reporters = append(reporters, NewProgressBar("test"))
		reporters = append(reporters, NewQuietProgress("test"))

		mp := NewMultiProgress()
		reporters = append(reporters, mp.AddBar("test", "test"))

		// Test that all implement the interface
		for i, reporter := range reporters {
			assert.NotNil(t, reporter, "Reporter %d should not be nil", i)

			// Should be able to call all interface methods without panic
			reporter.UpdateMessage("test")
			reporter.UpdateProgress(1, 2)
			reporter.Finish()
		}
	})
}
