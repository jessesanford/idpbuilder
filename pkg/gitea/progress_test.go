package gitea

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProgressTracker(t *testing.T) {
	totalSize := int64(1000)
	tracker := NewProgressTracker(totalSize)

	assert.NotNil(t, tracker)
	assert.Equal(t, totalSize, tracker.totalBytes)
	assert.Equal(t, int64(0), tracker.uploadedBytes)
	assert.False(t, tracker.completed)
	assert.NotNil(t, tracker.layers)
	assert.NotZero(t, tracker.startTime)
}

func TestProgressTracker_UpdateProgress(t *testing.T) {
	tracker := NewProgressTracker(1000)

	// Update progress for a layer
	tracker.UpdateProgress("layer-1", 100)

	progress := tracker.GetProgress()
	assert.Equal(t, int64(100), progress.UploadedBytes)
	assert.Equal(t, 10, progress.Percentage) // 100/1000 * 100 = 10%
	assert.Equal(t, 1, len(progress.LayerProgresses))

	// Update progress for same layer
	tracker.UpdateProgress("layer-1", 200)
	progress = tracker.GetProgress()
	assert.Equal(t, int64(200), progress.UploadedBytes)
	assert.Equal(t, 20, progress.Percentage)

	// Add another layer
	tracker.UpdateProgress("layer-2", 300)
	progress = tracker.GetProgress()
	assert.Equal(t, int64(500), progress.UploadedBytes)
	assert.Equal(t, 50, progress.Percentage)
	assert.Equal(t, 2, len(progress.LayerProgresses))
}

func TestProgressTracker_SetLayerSize(t *testing.T) {
	tracker := NewProgressTracker(1000)

	tracker.SetLayerSize("layer-1", 500)

	progress := tracker.GetProgress()
	assert.Equal(t, 1, len(progress.LayerProgresses))

	layer := progress.LayerProgresses[0]
	assert.Equal(t, "layer-1", layer.ID)
	assert.Equal(t, int64(500), layer.Size)
	assert.Equal(t, int64(0), layer.Uploaded)
	assert.Equal(t, "pending", layer.Status)
}

func TestProgressTracker_SetLayerStatus(t *testing.T) {
	tracker := NewProgressTracker(1000)

	tracker.SetLayerSize("layer-1", 500)
	tracker.SetLayerStatus("layer-1", "uploading")

	progress := tracker.GetProgress()
	layer := progress.LayerProgresses[0]
	assert.Equal(t, "uploading", layer.Status)
}

func TestProgressTracker_MarkComplete(t *testing.T) {
	tracker := NewProgressTracker(1000)
	tracker.SetLayerSize("layer-1", 500)
	tracker.SetLayerSize("layer-2", 500)
	tracker.UpdateProgress("layer-1", 300)

	assert.False(t, tracker.IsComplete())

	tracker.MarkComplete()

	assert.True(t, tracker.IsComplete())
	progress := tracker.GetProgress()
	assert.Equal(t, int64(1000), progress.UploadedBytes)
	assert.Equal(t, 100, progress.Percentage)
	assert.True(t, progress.Completed)

	// Check that all layers are marked complete
	for _, layer := range progress.LayerProgresses {
		assert.Equal(t, "complete", layer.Status)
		assert.Equal(t, layer.Size, layer.Uploaded)
	}
}

func TestProgressTracker_IncrementBytes(t *testing.T) {
	tracker := NewProgressTracker(1000)

	tracker.IncrementBytes(100)
	progress := tracker.GetProgress()
	assert.Equal(t, int64(100), progress.UploadedBytes)
	assert.Equal(t, 10, progress.Percentage)

	tracker.IncrementBytes(900)
	progress = tracker.GetProgress()
	assert.Equal(t, int64(1000), progress.UploadedBytes)
	assert.Equal(t, 100, progress.Percentage)
	assert.True(t, progress.Completed)
}

func TestProgressTracker_GetEstimatedTimeRemaining(t *testing.T) {
	tracker := NewProgressTracker(1000)

	// No progress yet, should return 0
	eta := tracker.GetEstimatedTimeRemaining()
	assert.Equal(t, time.Duration(0), eta)

	// Simulate some progress
	time.Sleep(10 * time.Millisecond) // Let some time pass
	tracker.IncrementBytes(100)

	eta = tracker.GetEstimatedTimeRemaining()
	assert.Greater(t, eta, time.Duration(0))

	// Complete should return 0
	tracker.MarkComplete()
	eta = tracker.GetEstimatedTimeRemaining()
	assert.Equal(t, time.Duration(0), eta)
}

func TestLayerProgress_GetPercentage(t *testing.T) {
	layer := LayerProgress{
		ID:       "test-layer",
		Size:     1000,
		Uploaded: 250,
		Status:   "uploading",
	}

	assert.Equal(t, 25, layer.GetPercentage())

	// Test edge cases
	layer.Size = 0
	assert.Equal(t, 0, layer.GetPercentage())

	layer.Size = 100
	layer.Uploaded = 150 // More than size
	assert.Equal(t, 100, layer.GetPercentage())
}

func TestProgressTracker_GetProgress(t *testing.T) {
	tracker := NewProgressTracker(1000)
	tracker.SetLayerSize("layer-1", 400)
	tracker.SetLayerSize("layer-2", 600)

	// Initial progress
	progress := tracker.GetProgress()
	assert.Equal(t, 0, progress.Percentage)
	assert.Equal(t, 0, progress.CurrentLayer)
	assert.Equal(t, 2, progress.TotalLayers)
	assert.False(t, progress.Completed)

	// Mark one layer complete
	tracker.UpdateProgress("layer-1", 400)
	tracker.SetLayerStatus("layer-1", "complete")

	progress = tracker.GetProgress()
	assert.Equal(t, 40, progress.Percentage) // 400/1000
	assert.Equal(t, 1, progress.CurrentLayer) // 1 complete layer
	assert.Equal(t, 2, progress.TotalLayers)

	// Progress on second layer
	tracker.UpdateProgress("layer-2", 300)
	tracker.SetLayerStatus("layer-2", "uploading")

	progress = tracker.GetProgress()
	assert.Equal(t, 70, progress.Percentage) // (400+300)/1000
	assert.Equal(t, 1, progress.CurrentLayer) // Still 1 complete layer

	// Complete second layer
	tracker.UpdateProgress("layer-2", 600)
	tracker.SetLayerStatus("layer-2", "complete")

	progress = tracker.GetProgress()
	assert.Equal(t, 100, progress.Percentage)
	assert.Equal(t, 2, progress.CurrentLayer) // 2 complete layers
	assert.True(t, progress.Completed)
}