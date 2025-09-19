package gitea

import (
	"sync"
	"time"
)

// ProgressTracker tracks real image push progress
type ProgressTracker struct {
	mu            sync.RWMutex
	totalBytes    int64
	uploadedBytes int64
	layers        map[string]*LayerProgress
	startTime     time.Time
	completed     bool
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(totalSize int64) *ProgressTracker {
	return &ProgressTracker{
		totalBytes: totalSize,
		layers:     make(map[string]*LayerProgress),
		startTime:  time.Now(),
	}
}

// UpdateProgress updates the progress for a specific layer
func (pt *ProgressTracker) UpdateProgress(layerID string, bytes int64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if layer, exists := pt.layers[layerID]; exists {
		// Update existing layer
		oldUploaded := layer.Uploaded
		layer.Uploaded = bytes

		// Update total uploaded bytes
		pt.uploadedBytes += (bytes - oldUploaded)
	} else {
		// New layer
		pt.layers[layerID] = &LayerProgress{
			ID:       layerID,
			Size:     bytes, // Will be updated with actual size
			Uploaded: bytes,
			Status:   "uploading",
		}
		pt.uploadedBytes += bytes
	}

	// Mark as complete if we've uploaded everything
	if pt.uploadedBytes >= pt.totalBytes {
		pt.completed = true
		for _, layer := range pt.layers {
			layer.Status = "complete"
		}
	}
}

// SetLayerSize sets the total size for a layer
func (pt *ProgressTracker) SetLayerSize(layerID string, size int64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if layer, exists := pt.layers[layerID]; exists {
		layer.Size = size
	} else {
		pt.layers[layerID] = &LayerProgress{
			ID:       layerID,
			Size:     size,
			Uploaded: 0,
			Status:   "pending",
		}
	}
}

// SetLayerStatus sets the status for a layer
func (pt *ProgressTracker) SetLayerStatus(layerID string, status string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if layer, exists := pt.layers[layerID]; exists {
		layer.Status = status
	}
}

// GetProgress returns current progress information
func (pt *ProgressTracker) GetProgress() PushProgress {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	var percentage int
	if pt.totalBytes > 0 {
		percentage = int((pt.uploadedBytes * 100) / pt.totalBytes)
		if percentage > 100 {
			percentage = 100
		}
	}

	// Count layers
	totalLayers := len(pt.layers)
	currentLayer := 0

	for _, layer := range pt.layers {
		if layer.Status == "complete" || layer.Status == "uploaded" {
			currentLayer++
		}
	}

	return PushProgress{
		CurrentLayer:    currentLayer,
		TotalLayers:     totalLayers,
		Percentage:      percentage,
		UploadedBytes:   pt.uploadedBytes,
		TotalBytes:      pt.totalBytes,
		ElapsedTime:     time.Since(pt.startTime),
		Completed:       pt.completed,
		LayerProgresses: pt.getLayerProgresses(),
	}
}

// getLayerProgresses returns a copy of layer progress data
func (pt *ProgressTracker) getLayerProgresses() []LayerProgress {
	var progresses []LayerProgress
	for _, layer := range pt.layers {
		progresses = append(progresses, LayerProgress{
			ID:       layer.ID,
			Size:     layer.Size,
			Uploaded: layer.Uploaded,
			Status:   layer.Status,
		})
	}
	return progresses
}

// IncrementBytes adds bytes to the total uploaded count
func (pt *ProgressTracker) IncrementBytes(bytes int64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.uploadedBytes += bytes

	// Check if completed
	if pt.uploadedBytes >= pt.totalBytes {
		pt.completed = true
	}
}

// MarkComplete marks the entire operation as complete
func (pt *ProgressTracker) MarkComplete() {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.completed = true
	pt.uploadedBytes = pt.totalBytes

	for _, layer := range pt.layers {
		layer.Status = "complete"
		layer.Uploaded = layer.Size
	}
}

// IsComplete returns whether the operation is complete
func (pt *ProgressTracker) IsComplete() bool {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	return pt.completed
}

// GetEstimatedTimeRemaining calculates estimated time remaining
func (pt *ProgressTracker) GetEstimatedTimeRemaining() time.Duration {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	if pt.uploadedBytes == 0 || pt.completed {
		return 0
	}

	elapsed := time.Since(pt.startTime)
	rate := float64(pt.uploadedBytes) / elapsed.Seconds() // bytes per second

	remaining := pt.totalBytes - pt.uploadedBytes
	if rate > 0 {
		return time.Duration(float64(remaining)/rate) * time.Second
	}

	return 0
}

// LayerProgress tracks progress for individual layers
type LayerProgress struct {
	ID       string `json:"id"`
	Size     int64  `json:"size"`
	Uploaded int64  `json:"uploaded"`
	Status   string `json:"status"` // pending, uploading, complete, error
}

// GetPercentage returns the percentage complete for this layer
func (lp *LayerProgress) GetPercentage() int {
	if lp.Size == 0 {
		return 0
	}
	percentage := int((lp.Uploaded * 100) / lp.Size)
	if percentage > 100 {
		percentage = 100
	}
	return percentage
}
