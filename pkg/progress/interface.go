// Package progress provides progress reporting interfaces and implementations
// for OCI registry operations.
//
// This package defines comprehensive progress tracking for operations like
// artifact push, pull, and layer transfers. It supports different reporting
// mechanisms including console output, logging, and silent operation.
package progress

import (
	"context"
	"time"
)

// OperationType represents the type of operation being tracked.
type OperationType string

const (
	// OperationPush represents an OCI artifact push operation
	OperationPush OperationType = "push"

	// OperationPull represents an OCI artifact pull operation
	OperationPull OperationType = "pull"

	// OperationList represents a repository listing operation
	OperationList OperationType = "list"

	// OperationDelete represents an artifact deletion operation
	OperationDelete OperationType = "delete"
)

// OperationState represents the current state of an operation.
type OperationState string

const (
	// StateStarted indicates the operation has begun
	StateStarted OperationState = "started"

	// StateInProgress indicates the operation is actively progressing
	StateInProgress OperationState = "in_progress"

	// StateCompleted indicates the operation completed successfully
	StateCompleted OperationState = "completed"

	// StateError indicates the operation failed with an error
	StateError OperationState = "error"

	// StateCancelled indicates the operation was cancelled
	StateCancelled OperationState = "cancelled"
)

// OperationProgress contains detailed information about an operation's progress.
type OperationProgress struct {
	// OperationType specifies what kind of operation this is
	OperationType OperationType

	// Reference is the OCI reference being operated on
	Reference string

	// State is the current state of the operation
	State OperationState

	// Current is the current progress amount (bytes, layers, etc.)
	Current int64

	// Total is the total expected amount for completion
	Total int64

	// Message provides human-readable status information
	Message string

	// StartTime records when the operation began
	StartTime time.Time

	// LastUpdate records when progress was last updated
	LastUpdate time.Time

	// Error contains any error that occurred (when State is StateError)
	Error error
}

// ProgressReporter defines the interface for reporting OCI operation progress.
// Implementations can provide console output, structured logging, or other
// progress indication mechanisms optimized for different use cases.
type ProgressReporter interface {
	// StartOperation begins tracking a new operation.
	// Returns a context that can be used to cancel the operation.
	StartOperation(ctx context.Context, opType OperationType, reference string) context.Context

	// UpdateProgress reports current progress for an ongoing operation.
	// The progress information is used to update displays and logs.
	UpdateProgress(ctx context.Context, progress OperationProgress)

	// CompleteOperation marks an operation as successfully completed.
	// Final statistics and timing information should be reported.
	CompleteOperation(ctx context.Context, reference string, message string)

	// ErrorOperation reports that an operation failed with an error.
	// The error should be logged/displayed appropriately.
	ErrorOperation(ctx context.Context, reference string, err error)

	// CancelOperation handles cancellation of an operation.
	// Clean up any ongoing displays and report cancellation.
	CancelOperation(ctx context.Context, reference string, reason string)
}

// LayerProgress contains progress information for individual layer operations.
// This is used when operations involve multiple layers with separate progress.
type LayerProgress struct {
	// LayerDigest identifies the specific layer
	LayerDigest string

	// LayerSize is the total size of this layer in bytes
	LayerSize int64

	// BytesTransferred is how many bytes have been processed
	BytesTransferred int64

	// State indicates if the layer is pending, in progress, or complete
	State OperationState

	// StartTime records when this layer operation began
	StartTime time.Time
}

// MultiLayerReporter extends ProgressReporter with layer-specific progress tracking.
// This interface is useful for operations that transfer multiple layers with
// independent progress tracking for each layer.
type MultiLayerReporter interface {
	ProgressReporter

	// StartLayer begins tracking progress for a specific layer.
	StartLayer(ctx context.Context, reference string, layer LayerProgress)

	// UpdateLayer reports progress for a specific layer transfer.
	UpdateLayer(ctx context.Context, reference string, layer LayerProgress)

	// CompleteLayer marks a specific layer as completed.
	CompleteLayer(ctx context.Context, reference string, layerDigest string)
}

// ProgressCallback is a function type for receiving progress updates.
// This can be used for custom progress handling or testing.
type ProgressCallback func(progress OperationProgress)