// Package push provides metrics collection hooks for monitoring push operations
package push

import (
	"time"
)

// Metrics provides hooks for monitoring and observability of push operations
type Metrics interface {
	// RecordPushStart is called when a push operation begins
	RecordPushStart(image, registry string)

	// RecordPushComplete is called when a push operation completes
	RecordPushComplete(image, registry string, duration time.Duration, err error)

	// RecordRetry is called when a retry is attempted
	RecordRetry(image, registry string, attempt int, reason string)

	// RecordProgress is called periodically during push to report bytes transferred
	RecordProgress(image string, bytes int64, total int64)

	// RecordLayerUpload is called for each layer upload
	RecordLayerUpload(image string, layerDigest string, size int64, duration time.Duration)
}

// NoOpMetrics is a no-op implementation of Metrics interface
type NoOpMetrics struct{}

// RecordPushStart implements Metrics interface
func (n *NoOpMetrics) RecordPushStart(image, registry string) {}

// RecordPushComplete implements Metrics interface
func (n *NoOpMetrics) RecordPushComplete(image, registry string, duration time.Duration, err error) {}

// RecordRetry implements Metrics interface
func (n *NoOpMetrics) RecordRetry(image, registry string, attempt int, reason string) {}

// RecordProgress implements Metrics interface
func (n *NoOpMetrics) RecordProgress(image string, bytes int64, total int64) {}

// RecordLayerUpload implements Metrics interface
func (n *NoOpMetrics) RecordLayerUpload(image string, layerDigest string, size int64, duration time.Duration) {}
