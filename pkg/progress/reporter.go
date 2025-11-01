package progress

import (
	"fmt"
	"sync"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/registry"
)

// Reporter tracks and displays progress for image push operations
type Reporter struct {
	verbose   bool                        // Enable detailed output
	startTime time.Time                   // Start time for duration calculation
	layers    map[string]*LayerProgress   // Layer digest -> progress
	mu        sync.Mutex                  // Protects layers map
}

// LayerProgress tracks individual layer upload progress
type LayerProgress struct {
	Digest       string     // Layer digest (e.g., sha256:abc123...)
	Size         int64      // Total layer size in bytes
	Pushed       int64      // Bytes pushed so far
	Status       string     // Current status (uploading, complete, exists)
	StartTime    time.Time  // When layer push started
	CompleteTime *time.Time // When layer push completed (nil if in progress)
}

// NewReporter creates a progress reporter
func NewReporter(verbose bool) *Reporter {
	return &Reporter{
		verbose:   verbose,
		startTime: time.Now(),
		layers:    make(map[string]*LayerProgress),
	}
}

// HandleProgress processes a progress update from registry.Push()
// This matches the registry.ProgressCallback signature from Phase 1
func (r *Reporter) HandleProgress(update registry.ProgressUpdate) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get or create layer progress
	layer, exists := r.layers[update.LayerDigest]
	if !exists {
		layer = &LayerProgress{
			Digest:    update.LayerDigest,
			Size:      update.LayerSize,
			StartTime: time.Now(),
		}
		r.layers[update.LayerDigest] = layer
	}

	// Update layer state
	layer.Pushed = update.BytesPushed
	layer.Status = update.Status
	if update.Status == "complete" || update.Status == "exists" {
		now := time.Now()
		layer.CompleteTime = &now
	}

	// Display progress
	if r.verbose {
		r.displayVerbose(layer)
	} else {
		r.displayNormal(layer)
	}
}

// displayNormal shows compact progress (for normal mode)
func (r *Reporter) displayNormal(layer *LayerProgress) {
	shortDigest := layer.Digest
	if len(shortDigest) > 12 {
		shortDigest = shortDigest[:12]
	}

	switch layer.Status {
	case "uploading":
		percent := float64(layer.Pushed) / float64(layer.Size) * 100
		fmt.Printf("  %s: Pushing [%.1f%%]\n", shortDigest, percent)
	case "complete":
		fmt.Printf("  %s: Pushed ✓\n", shortDigest)
	case "exists":
		fmt.Printf("  %s: Already exists ✓\n", shortDigest)
	}
}

// displayVerbose shows detailed progress (for --verbose)
func (r *Reporter) displayVerbose(layer *LayerProgress) {
	elapsed := time.Since(layer.StartTime).Seconds()

	switch layer.Status {
	case "uploading":
		percent := float64(layer.Pushed) / float64(layer.Size) * 100
		rate := 0.0
		if elapsed > 0.001 { // Avoid division by zero
			rate = float64(layer.Pushed) / elapsed / 1024 / 1024 // MB/s
		}
		fmt.Printf("  Layer %s:\n", layer.Digest)
		fmt.Printf("    Status: Uploading [%.1f%%] (%d/%d bytes)\n",
			percent, layer.Pushed, layer.Size)
		fmt.Printf("    Rate: %.2f MB/s\n", rate)
		fmt.Printf("    Elapsed: %.1fs\n", elapsed)
	case "complete":
		duration := elapsed
		if layer.CompleteTime != nil {
			duration = layer.CompleteTime.Sub(layer.StartTime).Seconds()
		}
		fmt.Printf("  Layer %s: Complete (%.1fs)\n", layer.Digest, duration)
	case "exists":
		fmt.Printf("  Layer %s: Already exists (skipped)\n", layer.Digest)
	}
}

// DisplaySummary shows final statistics after push completes
func (r *Reporter) DisplaySummary() {
	r.mu.Lock()
	defer r.mu.Unlock()

	totalLayers := len(r.layers)
	totalBytes := int64(0)
	pushedLayers := 0
	skippedLayers := 0

	for _, layer := range r.layers {
		totalBytes += layer.Size
		if layer.Status == "complete" {
			pushedLayers++
		} else if layer.Status == "exists" {
			skippedLayers++
		}
	}

	duration := time.Since(r.startTime)
	avgRate := 0.0
	if duration.Seconds() > 0.001 { // Avoid division by zero
		avgRate = float64(totalBytes) / duration.Seconds() / 1024 / 1024 // MB/s
	}

	fmt.Println("\nPush Summary:")
	fmt.Printf("  Total layers: %d\n", totalLayers)
	fmt.Printf("  Pushed: %d, Skipped: %d\n", pushedLayers, skippedLayers)
	fmt.Printf("  Total size: %.2f MB\n", float64(totalBytes)/1024/1024)
	fmt.Printf("  Duration: %.1fs\n", duration.Seconds())
	fmt.Printf("  Average rate: %.2f MB/s\n", avgRate)
}

// GetCallback returns a registry.ProgressCallback that uses this reporter
func (r *Reporter) GetCallback() registry.ProgressCallback {
	return func(update registry.ProgressUpdate) {
		r.HandleProgress(update)
	}
}
