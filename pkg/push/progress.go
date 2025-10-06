package push

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// ProgressReporter defines the interface for reporting push progress
type ProgressReporter interface {
	// StartImage indicates the start of pushing an image
	StartImage(digest string, totalSize int64)

	// UpdateLayer reports progress on a specific layer
	UpdateLayer(digest string, written int64)

	// FinishLayer indicates a layer push has completed
	FinishLayer(digest string)

	// FinishImage indicates an image push has completed
	FinishImage(digest string)

	// SetError reports an error during push
	SetError(digest string, err error)
}

// ConsoleProgressReporter implements progress reporting to the console
type ConsoleProgressReporter struct {
	writer     io.Writer
	mu         sync.Mutex
	images     map[string]*imageProgress
	showBytes  bool
	showLayers bool
}

// imageProgress tracks progress for a single image
type imageProgress struct {
	digest    string
	totalSize int64
	written   int64
	layers    map[string]*layerProgress
	startTime time.Time
	err       error
}

// layerProgress tracks progress for a single layer
type layerProgress struct {
	digest  string
	written int64
	done    bool
}

// NewConsoleProgressReporter creates a new console-based progress reporter
func NewConsoleProgressReporter(writer io.Writer) *ConsoleProgressReporter {
	return &ConsoleProgressReporter{
		writer:     writer,
		images:     make(map[string]*imageProgress),
		showBytes:  true,
		showLayers: false, // Keep simple by default
	}
}

// NewDetailedConsoleProgressReporter creates a progress reporter with layer details
func NewDetailedConsoleProgressReporter(writer io.Writer) *ConsoleProgressReporter {
	return &ConsoleProgressReporter{
		writer:     writer,
		images:     make(map[string]*imageProgress),
		showBytes:  true,
		showLayers: true,
	}
}

// StartImage begins tracking progress for an image
func (c *ConsoleProgressReporter) StartImage(digest string, totalSize int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.images[digest] = &imageProgress{
		digest:    digest,
		totalSize: totalSize,
		written:   0,
		layers:    make(map[string]*layerProgress),
		startTime: time.Now(),
	}

	c.printf("📤 Pushing image %s", c.shortDigest(digest))
	if totalSize > 0 {
		c.printf(" (size: %s)", c.formatBytes(totalSize))
	}
	c.printf("\n")
}

// UpdateLayer updates the progress for a specific layer
func (c *ConsoleProgressReporter) UpdateLayer(digest string, written int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find the image this layer belongs to
	// For simplicity, we'll track layer progress per image
	for _, img := range c.images {
		if layer, exists := img.layers[digest]; exists {
			layer.written = written
			c.updateImageProgress(img)
			return
		} else {
			// Create new layer if it doesn't exist
			img.layers[digest] = &layerProgress{
				digest:  digest,
				written: written,
				done:    false,
			}
			c.updateImageProgress(img)
			return
		}
	}
}

// FinishLayer marks a layer as completed
func (c *ConsoleProgressReporter) FinishLayer(digest string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, img := range c.images {
		if layer, exists := img.layers[digest]; exists {
			layer.done = true
			if c.showLayers {
				c.printf("  ✓ Layer %s completed\n", c.shortDigest(digest))
			}
			return
		}
	}
}

// FinishImage marks an image push as completed
func (c *ConsoleProgressReporter) FinishImage(digest string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if img, exists := c.images[digest]; exists {
		duration := time.Since(img.startTime)

		if img.err != nil {
			c.printf("❌ Image %s failed: %v\n", c.shortDigest(digest), img.err)
		} else {
			c.printf("✅ Image %s pushed successfully", c.shortDigest(digest))
			if c.showBytes && img.totalSize > 0 {
				rate := float64(img.totalSize) / duration.Seconds()
				c.printf(" (%s in %s, %s/s)",
					c.formatBytes(img.totalSize),
					c.formatDuration(duration),
					c.formatBytes(int64(rate)))
			}
			c.printf("\n")
		}

		// Clean up completed image
		delete(c.images, digest)
	}
}

// SetError reports an error for an image push
func (c *ConsoleProgressReporter) SetError(digest string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if img, exists := c.images[digest]; exists {
		img.err = err
	}
}

// updateImageProgress recalculates and potentially displays image progress
func (c *ConsoleProgressReporter) updateImageProgress(img *imageProgress) {
	if !c.showBytes || img.totalSize <= 0 {
		return
	}

	// Calculate total written across all layers
	var totalWritten int64
	for _, layer := range img.layers {
		totalWritten += layer.written
	}

	img.written = totalWritten
	percentage := float64(img.written) / float64(img.totalSize) * 100

	// Only show progress for significant changes (every 10%)
	if int(percentage)/10 > int(float64(img.written-1)/float64(img.totalSize)*100)/10 {
		c.printf("  📊 Progress: %.1f%% (%s/%s)\n",
			percentage,
			c.formatBytes(img.written),
			c.formatBytes(img.totalSize))
	}
}

// printf is a thread-safe wrapper around fmt.Fprintf
func (c *ConsoleProgressReporter) printf(format string, args ...interface{}) {
	fmt.Fprintf(c.writer, format, args...)
}

// shortDigest returns a shortened version of a digest for display
func (c *ConsoleProgressReporter) shortDigest(digest string) string {
	if len(digest) > 12 {
		return digest[:12]
	}
	return digest
}

// formatBytes formats a byte count as a human-readable string
func (c *ConsoleProgressReporter) formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration formats a duration as a human-readable string
func (c *ConsoleProgressReporter) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	} else {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
}

// NoOpProgressReporter is a progress reporter that does nothing
type NoOpProgressReporter struct{}

// NewNoOpProgressReporter creates a no-operation progress reporter
func NewNoOpProgressReporter() *NoOpProgressReporter {
	return &NoOpProgressReporter{}
}

// StartImage does nothing
func (n *NoOpProgressReporter) StartImage(digest string, totalSize int64) {}

// UpdateLayer does nothing
func (n *NoOpProgressReporter) UpdateLayer(digest string, written int64) {}

// FinishLayer does nothing
func (n *NoOpProgressReporter) FinishLayer(digest string) {}

// FinishImage does nothing
func (n *NoOpProgressReporter) FinishImage(digest string) {}

// SetError does nothing
func (n *NoOpProgressReporter) SetError(digest string, err error) {}

// MultiProgressReporter combines multiple progress reporters
type MultiProgressReporter struct {
	reporters []ProgressReporter
}

// NewMultiProgressReporter creates a progress reporter that forwards to multiple reporters
func NewMultiProgressReporter(reporters ...ProgressReporter) *MultiProgressReporter {
	return &MultiProgressReporter{
		reporters: reporters,
	}
}

// StartImage forwards to all reporters
func (m *MultiProgressReporter) StartImage(digest string, totalSize int64) {
	for _, reporter := range m.reporters {
		reporter.StartImage(digest, totalSize)
	}
}

// UpdateLayer forwards to all reporters
func (m *MultiProgressReporter) UpdateLayer(digest string, written int64) {
	for _, reporter := range m.reporters {
		reporter.UpdateLayer(digest, written)
	}
}

// FinishLayer forwards to all reporters
func (m *MultiProgressReporter) FinishLayer(digest string) {
	for _, reporter := range m.reporters {
		reporter.FinishLayer(digest)
	}
}

// FinishImage forwards to all reporters
func (m *MultiProgressReporter) FinishImage(digest string) {
	for _, reporter := range m.reporters {
		reporter.FinishImage(digest)
	}
}

// SetError forwards to all reporters
func (m *MultiProgressReporter) SetError(digest string, err error) {
	for _, reporter := range m.reporters {
		reporter.SetError(digest, err)
	}
}