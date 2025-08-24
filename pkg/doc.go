// Package errors provides structured error handling types and utilities
// for OCI (Open Container Initiative) management operations.
//
// This package defines standard error interfaces, codes, and categories
// that enable consistent error handling across the idpbuilder OCI
// management system. It follows Go 1.13+ error wrapping patterns
// and provides enhanced context for debugging and monitoring.
//
// Error Handling Example:
//
//	// Creating a structured error
//	err := &errors.BaseError{
//		Code:      errors.CodeRegistryUnavailable,
//		Category:  errors.CategoryTransient,
//		Message:   "registry connection failed",
//		Context:   map[string]interface{}{"registry": "docker.io"},
//		Timestamp: time.Now().UTC(),
//	}
//
//	// Wrapping an existing error
//	wrappedErr := err.Wrap(originalErr)
//
//	// Checking error properties
//	if ociErr, ok := err.(errors.OCIError); ok {
//		if ociErr.Category().IsRetryable() {
//			// Retry the operation
//		}
//	}
//
// Progress Tracking Example:
//
//	// Creating a progress tracker
//	tracker := progress.NewBaseProgressTracker()
//
//	// Starting an operation
//	tracker.Start(100) // 100 total work units
//
//	// Updating progress
//	tracker.UpdateWithMessage(25, "Downloaded 25 layers")
//
//	// Completing the operation
//	tracker.Complete()
//
// The progress package provides thread-safe progress tracking with
// event notifications, percentage calculations, and time estimates.
// It supports both determinate (known total) and indeterminate
// (unknown total) progress tracking patterns.
//
// Best Practices:
//
// 1. Use structured errors with appropriate codes and categories
// 2. Include relevant context information for debugging
// 3. Wrap errors to preserve the error chain
// 4. Check error categories to determine retry strategies
// 5. Use progress tracking for long-running operations
// 6. Subscribe to progress events for user feedback
package pkg