package progress

import (
	"time"
)

// Operation represents a long-running operation being tracked
type Operation struct {
	Name      string    // Operation name
	Total     int64     // Total units of work
	Current   int64     // Current progress
	StartTime time.Time // When operation started
}

// Result represents the completion result of an operation
type Result struct {
	Success  bool          // Whether operation succeeded
	Message  string        // Completion or error message
	Duration time.Duration // How long operation took
}