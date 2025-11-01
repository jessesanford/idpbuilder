// Package progress provides progress tracking and reporting for image push operations
package progress

import "github.com/cnoe-io/idpbuilder/pkg/registry"

// ProgressReporter defines operations for tracking and displaying push progress
type ProgressReporter interface {
	// HandleProgress processes a progress update from registry push
	HandleProgress(update registry.ProgressUpdate)

	// DisplaySummary shows final statistics after push completes
	DisplaySummary()

	// GetCallback returns a callback function for registry.Push()
	GetCallback() registry.ProgressCallback
}
