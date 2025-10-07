package push

import (
	"context"
	"time"

	"github.com/go-logr/logr"
)

// LogLevel defines the logging levels for push operations
type LogLevel int

const (
	// LogLevelDebug provides detailed debugging information
	LogLevelDebug LogLevel = iota
	// LogLevelInfo provides general information
	LogLevelInfo
	// LogLevelWarn provides warning messages
	LogLevelWarn
	// LogLevelError provides error messages
	LogLevelError
)

// PushLogger wraps logr.Logger with push-specific logging methods
type PushLogger struct {
	logger logr.Logger
}

// NewPushLogger creates a new push-specific logger wrapper
func NewPushLogger(logger logr.Logger) *PushLogger {
	return &PushLogger{
		logger: logger,
	}
}

// LogPushStart logs the beginning of an image push operation
func (p *PushLogger) LogPushStart(imageName, registry string, totalSize int64) {
	p.logger.Info("Starting image push operation",
		"image", imageName,
		"registry", registry,
		"size_bytes", totalSize,
		"timestamp", time.Now().Format(time.RFC3339))
}

// LogPushComplete logs the successful completion of an image push
func (p *PushLogger) LogPushComplete(imageName, digest string, duration time.Duration, totalSize int64) {
	throughput := float64(totalSize) / duration.Seconds()

	p.logger.Info("Image push completed successfully",
		"image", imageName,
		"digest", digest,
		"duration", duration.String(),
		"size_bytes", totalSize,
		"throughput_bytes_per_sec", int64(throughput))
}

// LogPushError logs push operation errors
func (p *PushLogger) LogPushError(imageName string, err error, attempt int, willRetry bool) {
	level := p.logger
	if willRetry {
		level = p.logger.V(1) // Debug level for retryable errors
	}

	level.Error(err, "Image push failed",
		"image", imageName,
		"attempt", attempt,
		"will_retry", willRetry)
}

// LogDiscoveryStart logs the start of image discovery
func (p *PushLogger) LogDiscoveryStart(buildPath string) {
	p.logger.Info("Starting image discovery",
		"build_path", buildPath,
		"timestamp", time.Now().Format(time.RFC3339))
}

// LogDiscoveryComplete logs the completion of image discovery
func (p *PushLogger) LogDiscoveryComplete(imageCount int, duration time.Duration) {
	p.logger.Info("Image discovery completed",
		"images_found", imageCount,
		"duration", duration.String())
}

// LogImageDiscovered logs when an individual image is discovered
func (p *PushLogger) LogImageDiscovered(imageName, format, path string, size int64) {
	p.logger.V(1).Info("Discovered image",
		"name", imageName,
		"format", format,
		"path", path,
		"size_bytes", size)
}

// LogLayerProgress logs layer upload progress (debug level)
func (p *PushLogger) LogLayerProgress(layerDigest string, written, total int64) {
	if total > 0 {
		percent := float64(written) / float64(total) * 100
		p.logger.V(2).Info("Layer upload progress",
			"layer", layerDigest[:12], // Shortened digest
			"written_bytes", written,
			"total_bytes", total,
			"progress_percent", int(percent))
	}
}

// LogRegistryRequest logs HTTP requests to the registry (debug level)
func (p *PushLogger) LogRegistryRequest(method, url string) {
	p.logger.V(2).Info("Registry request",
		"method", method,
		"url", url,
		"timestamp", time.Now().Format(time.RFC3339))
}

// LogRegistryResponse logs HTTP responses from the registry (debug level)
func (p *PushLogger) LogRegistryResponse(method, url string, statusCode int, duration time.Duration) {
	p.logger.V(2).Info("Registry response",
		"method", method,
		"url", url,
		"status_code", statusCode,
		"duration", duration.String())
}

// LogRetryAttempt logs when a retry attempt is being made
func (p *PushLogger) LogRetryAttempt(imageName string, attempt, maxRetries int, backoff time.Duration) {
	p.logger.V(1).Info("Retrying push operation",
		"image", imageName,
		"attempt", attempt,
		"max_retries", maxRetries,
		"backoff", backoff.String())
}

// LogAuthentication logs authentication setup (without sensitive data)
func (p *PushLogger) LogAuthentication(registry string, authMethod string) {
	p.logger.V(1).Info("Setting up registry authentication",
		"registry", registry,
		"auth_method", authMethod)
}

// LogImageValidation logs image validation results
func (p *PushLogger) LogImageValidation(imageName string, isValid bool, issues []string) {
	if isValid {
		p.logger.V(1).Info("Image validation passed", "image", imageName)
	} else {
		p.logger.Info("Image validation failed",
			"image", imageName,
			"issues", issues)
	}
}

// LogPerformanceMetrics logs detailed performance metrics
func (p *PushLogger) LogPerformanceMetrics(metrics *PushMetrics) {
	p.logger.Info("Push operation metrics",
		"total_duration", metrics.TotalDuration.String(),
		"discovery_duration", metrics.DiscoveryDuration.String(),
		"push_duration", metrics.PushDuration.String(),
		"total_bytes", metrics.TotalBytes,
		"images_pushed", metrics.ImagesPushed,
		"layers_uploaded", metrics.LayersUploaded,
		"retry_count", metrics.RetryCount,
		"average_throughput_mbps", metrics.AverageThroughputMBps())
}

// PushMetrics contains performance metrics for push operations
type PushMetrics struct {
	TotalDuration     time.Duration
	DiscoveryDuration time.Duration
	PushDuration      time.Duration
	TotalBytes        int64
	ImagesPushed      int
	LayersUploaded    int
	RetryCount        int
}

// AverageThroughputMBps calculates the average throughput in MB/s
func (m *PushMetrics) AverageThroughputMBps() float64 {
	if m.PushDuration.Seconds() == 0 {
		return 0
	}
	bytesPerSecond := float64(m.TotalBytes) / m.PushDuration.Seconds()
	return bytesPerSecond / (1024 * 1024) // Convert to MB/s
}

// WithContext adds context information to the logger
func (p *PushLogger) WithContext(ctx context.Context) *PushLogger {
	// Extract context values if available
	logger := p.logger

	// Add trace ID if available
	if traceID, ok := ctx.Value("trace-id").(string); ok {
		logger = logger.WithValues("trace_id", traceID)
	}

	// Add user context if available
	if userID, ok := ctx.Value("user-id").(string); ok {
		logger = logger.WithValues("user_id", userID)
	}

	return &PushLogger{logger: logger}
}

// WithValues adds key-value pairs to the logger
func (p *PushLogger) WithValues(keysAndValues ...interface{}) *PushLogger {
	return &PushLogger{
		logger: p.logger.WithValues(keysAndValues...),
	}
}

// Structured logging helpers

// LogStructuredError logs an error with structured context
func (p *PushLogger) LogStructuredError(err error, operation string, context map[string]interface{}) {
	values := make([]interface{}, 0, len(context)*2+2)
	values = append(values, "operation", operation)

	for k, v := range context {
		values = append(values, k, v)
	}

	p.logger.WithValues(values...).Error(err, "Operation failed")
}

// LogStructuredInfo logs an info message with structured context
func (p *PushLogger) LogStructuredInfo(message string, context map[string]interface{}) {
	values := make([]interface{}, 0, len(context)*2)

	for k, v := range context {
		values = append(values, k, v)
	}

	p.logger.WithValues(values...).Info(message)
}

// LogDebug logs a debug message (V(2) level)
func (p *PushLogger) LogDebug(message string, keysAndValues ...interface{}) {
	p.logger.V(2).Info(message, keysAndValues...)
}

// LogInfo logs an info message
func (p *PushLogger) LogInfo(message string, keysAndValues ...interface{}) {
	p.logger.Info(message, keysAndValues...)
}

// LogWarn logs a warning message
func (p *PushLogger) LogWarn(message string, keysAndValues ...interface{}) {
	p.logger.V(1).Info("[WARNING] "+message, keysAndValues...)
}

// LogError logs an error message
func (p *PushLogger) LogError(err error, message string, keysAndValues ...interface{}) {
	p.logger.Error(err, message, keysAndValues...)
}
