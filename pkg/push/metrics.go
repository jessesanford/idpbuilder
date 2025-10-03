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

// TODO: Add OpenTelemetry integration for distributed tracing
// type OTelMetrics struct {
//     tracer trace.Tracer
//     meter  metric.Meter
// }
//
// func NewOTelMetrics(provider trace.TracerProvider, meterProvider metric.MeterProvider) *OTelMetrics {
//     return &OTelMetrics{
//         tracer: provider.Tracer("idpbuilder/push"),
//         meter:  meterProvider.Meter("idpbuilder/push"),
//     }
// }

// TODO: Add Prometheus metrics exporter
// type PrometheusMetrics struct {
//     pushDuration    prometheus.Histogram
//     retryCount      prometheus.Counter
//     bytesTransferred prometheus.Counter
//     layerUploads    prometheus.Histogram
// }
//
// func NewPrometheusMetrics(registry *prometheus.Registry) *PrometheusMetrics {
//     pm := &PrometheusMetrics{
//         pushDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
//             Name: "idpbuilder_push_duration_seconds",
//             Help: "Duration of push operations in seconds",
//         }),
//         retryCount: prometheus.NewCounter(prometheus.CounterOpts{
//             Name: "idpbuilder_push_retry_total",
//             Help: "Total number of retry attempts",
//         }),
//         bytesTransferred: prometheus.NewCounter(prometheus.CounterOpts{
//             Name: "idpbuilder_push_bytes_total",
//             Help: "Total bytes transferred during push",
//         }),
//         layerUploads: prometheus.NewHistogram(prometheus.HistogramOpts{
//             Name: "idpbuilder_push_layer_duration_seconds",
//             Help: "Duration of layer uploads in seconds",
//         }),
//     }
//     registry.MustRegister(pm.pushDuration, pm.retryCount, pm.bytesTransferred, pm.layerUploads)
//     return pm
// }
