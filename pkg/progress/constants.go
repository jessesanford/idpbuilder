package progress

import "time"

// Default update frequencies for progress reporting
const (
	// DefaultUpdateInterval is the default time between progress updates
	DefaultUpdateInterval = time.Millisecond * 500

	// FastUpdateInterval for high-frequency updates
	FastUpdateInterval = time.Millisecond * 100

	// SlowUpdateInterval for low-frequency updates
	SlowUpdateInterval = time.Second * 2

	// MinUpdateInterval prevents too frequent updates
	MinUpdateInterval = time.Millisecond * 50
)

// Progress percentage thresholds for significant milestones
const (
	// QuarterComplete represents 25% completion
	QuarterComplete = 25.0

	// HalfComplete represents 50% completion
	HalfComplete = 50.0

	// ThreeQuartersComplete represents 75% completion
	ThreeQuartersComplete = 75.0

	// NearlyComplete represents 95% completion
	NearlyComplete = 95.0
)

// Default progress message templates
const (
	MsgStarted           = "Operation started"
	MsgInProgress        = "Processing... %d%% complete"
	MsgCompleted         = "Operation completed successfully"
	MsgFailed            = "Operation failed: %s"
	MsgCancelled         = "Operation was cancelled"
	MsgIndeterminate     = "Processing... (progress unknown)"
	MsgEstimatedTime     = "Estimated time remaining: %s"
	MsgElapsedTime       = "Elapsed time: %s"
)

// Buffer sizes for event channels and queues
const (
	// DefaultEventBufferSize for progress event channels
	DefaultEventBufferSize = 100

	// LargeEventBufferSize for high-throughput scenarios
	LargeEventBufferSize = 1000

	// SmallEventBufferSize for low-memory environments
	SmallEventBufferSize = 10
)

// Limits for progress tracking to prevent resource exhaustion
const (
	// MaxCallbacks limits the number of progress callbacks
	MaxCallbacks = 50

	// MaxMetadataSize limits the size of progress metadata
	MaxMetadataSize = 4096

	// MaxMessageLength limits progress message length
	MaxMessageLength = 512
)