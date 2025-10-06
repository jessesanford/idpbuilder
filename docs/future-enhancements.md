# Future Enhancements for idpbuilder push

This document outlines potential future enhancements and features for the `idpbuilder push` command. Each section includes implementation notes and design considerations.

## Rate Limiting

**Priority**: Medium
**Estimated Effort**: 2-3 days

### Description
Implement rate limiting for registry operations to prevent overwhelming registries and respect API rate limits.

### Implementation Notes
```go
// TODO: Implement rate limiting for registry operations
// - Token bucket algorithm for smooth rate limiting
// - Per-registry configuration support
// - Backpressure handling with exponential backoff
// - Dynamic rate adjustment based on registry responses

type RateLimiter struct {
    registry      string
    requestsPerSec float64
    bucket        *TokenBucket
    backoff       *AdaptiveBackoff
}

func NewRateLimiter(registry string, rps float64) *RateLimiter {
    // Implementation here
}
```

### Configuration
```yaml
rate_limiting:
  enabled: true
  default_rps: 100
  per_registry:
    "docker.io": 50
    "gcr.io": 200
```

## Advanced Multi-Architecture Support

**Priority**: High
**Estimated Effort**: 5-7 days

### Description
Enhanced support for multi-architecture images including automatic manifest list creation and platform-specific optimizations.

### Implementation Notes
```go
// TODO: Multi-architecture image support
// - Automatic manifest list generation
// - Platform-specific layer deduplication
// - Cross-platform build coordination
// - ARM, x86_64, and other architecture support

type MultiArchPusher struct {
    platforms []v1.Platform
    strategy  ManifestStrategy
}

// ManifestStrategy determines how to handle multi-arch manifests
type ManifestStrategy int

const (
    MergeManifests ManifestStrategy = iota
    SeparateManifests
    AutoDetect
)
```

## Parallel Layer Uploads

**Priority**: High
**Estimated Effort**: 3-4 days

### Description
Optimize push performance by uploading image layers in parallel where registry supports it.

### Implementation Notes
```go
// TODO: Parallel layer uploads for improved performance
// - Concurrent layer upload with configurable workers
// - Dependency graph analysis for layer ordering
// - Intelligent chunking for large layers
// - Registry capability detection

type ParallelUploader struct {
    maxWorkers    int
    chunkSize     int64
    dependencyMap map[string][]string
}

func (pu *ParallelUploader) UploadLayers(ctx context.Context, layers []v1.Layer) error {
    // Worker pool implementation
    // Dependency-aware scheduling
    // Progress aggregation
}
```

### Performance Targets
- 3-5x faster for images with many independent layers
- Minimal memory overhead (chunked streaming)
- Graceful degradation for registries without parallel support

## Resume Capability

**Priority**: Medium
**Estimated Effort**: 4-5 days

### Description
Allow interrupted pushes to resume from where they left off instead of restarting.

### Implementation Notes
```go
// TODO: Resume capability for interrupted uploads
// - Track upload progress in state file
// - Detect partial uploads on restart
// - Resume from last successful chunk
// - Checksum verification for resumed uploads

type ResumeState struct {
    ImageDigest    string
    UploadedLayers map[string]LayerState
    LastCheckpoint time.Time
}

type LayerState struct {
    Digest       string
    BytesWritten int64
    TotalBytes   int64
    Checksum     string
}

func LoadResumeState(imageName string) (*ResumeState, error) {
    // Load state from ~/.idpbuilder/resume/<image>.json
}
```

## Signature Verification (Cosign Integration)

**Priority**: High (Security)
**Estimated Effort**: 7-10 days

### Description
Integrate with Sigstore's cosign for image signing and verification.

### Implementation Notes
```go
// TODO: Signature verification (cosign integration)
// - Sign images during push with cosign
// - Verify signatures before allowing push
// - Support for keyless signing (OIDC)
// - Integration with transparency logs

import "github.com/sigstore/cosign/v2/pkg/cosign"

type SigningConfig struct {
    Enabled    bool
    KeyPath    string
    UseKeyless bool
    OIDCIssuer string
}

func SignAndPush(ctx context.Context, img v1.Image, opts SigningConfig) error {
    // Sign image with cosign
    // Attach signature to registry
    // Upload to transparency log
}
```

## OpenTelemetry Integration

**Priority**: Medium
**Estimated Effort**: 3-4 days

### Description
Add OpenTelemetry tracing and metrics for observability in production environments.

### Implementation Notes
```go
// TODO: OpenTelemetry integration for distributed tracing
// - Trace push operations end-to-end
// - Track registry API calls with spans
// - Export metrics to OTLP endpoint
// - Baggage propagation for context

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
    "go.opentelemetry.io/otel/metric"
)

type OTelInstrumentation struct {
    tracer trace.Tracer
    meter  metric.Meter
}

func (o *OTelInstrumentation) TracePush(ctx context.Context, imageName string) (context.Context, trace.Span) {
    return o.tracer.Start(ctx, "push_image",
        trace.WithAttributes(
            attribute.String("image.name", imageName),
        ))
}
```

## Prometheus Metrics Export

**Priority**: Medium
**Estimated Effort**: 2-3 days

### Description
Export detailed metrics in Prometheus format for monitoring and alerting.

### Implementation Notes
```go
// TODO: Prometheus metrics endpoint
// - HTTP endpoint exposing /metrics
// - Push duration histogram
// - Retry counter by registry
// - Bytes transferred gauge
// - Error rate by type

import "github.com/prometheus/client_golang/prometheus"

var (
    pushDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "idpbuilder_push_duration_seconds",
            Help: "Duration of push operations",
            Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
        },
        []string{"registry", "status"},
    )

    retryCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "idpbuilder_push_retry_total",
            Help: "Total number of retry attempts",
        },
        []string{"registry", "reason"},
    )
)
```

## Structured Logging with Trace IDs

**Priority**: Low
**Estimated Effort**: 2 days

### Description
Enhanced logging with structured fields and correlation trace IDs.

### Implementation Notes
```go
// TODO: Structured logging with trace IDs
// - JSON log output option
// - Trace ID propagation across operations
// - Log aggregation friendly format
// - Sampling for high-volume operations

import "go.uber.org/zap"

type StructuredLogger struct {
    logger  *zap.Logger
    traceID string
}

func (sl *StructuredLogger) LogPush(imageName string, fields ...zap.Field) {
    fields = append(fields,
        zap.String("trace_id", sl.traceID),
        zap.String("image", imageName),
    )
    sl.logger.Info("pushing image", fields...)
}
```

## Image Vulnerability Scanning

**Priority**: High (Security)
**Estimated Effort**: 5-7 days

### Description
Integrate vulnerability scanning before push to prevent pushing vulnerable images.

### Implementation Notes
```go
// TODO: Image vulnerability scanning integration
// - Scan with Trivy/Grype before push
// - Configurable severity thresholds
// - Policy-based blocking (fail on HIGH/CRITICAL)
// - Scan result caching

type VulnScanner interface {
    Scan(ctx context.Context, img v1.Image) (*ScanResult, error)
}

type ScanResult struct {
    Vulnerabilities []Vulnerability
    Summary         VulnSummary
}

type VulnSummary struct {
    Critical int
    High     int
    Medium   int
    Low      int
}
```

## Registry-Specific Optimizations

**Priority**: Low
**Estimated Effort**: 3-4 days

### Description
Add registry-specific optimizations for popular registries.

### Implementation Notes
```go
// TODO: Registry-specific optimizations
// - AWS ECR: Use ECR-specific APIs for faster uploads
// - GCR: Leverage Google Cloud Storage acceleration
// - Azure ACR: Use ACR import for cross-region copies
// - Harbor: Use Harbor replication APIs

type RegistryOptimizer interface {
    CanOptimize(registry string) bool
    OptimizedPush(ctx context.Context, img v1.Image) error
}

type ECROptimizer struct {
    client *ecr.Client
}

type GCROptimizer struct {
    client *storage.Client
}
```

## Content Addressable Storage (CAS) Cache

**Priority**: Medium
**Estimated Effort**: 4-5 days

### Description
Local CAS cache to avoid re-uploading layers that exist in registry.

### Implementation Notes
```go
// TODO: Local CAS cache for layer deduplication
// - Track uploaded layers by digest
// - Skip upload if layer exists remotely
// - Periodic cache cleanup (LRU)
// - Cross-registry deduplication

type CASCache struct {
    storage map[string]CacheEntry
    maxSize int64
}

type CacheEntry struct {
    Digest     string
    Registries []string
    LastUsed   time.Time
}

func (c *CASCache) HasLayer(digest string, registry string) bool {
    // Check if layer exists in registry
}
```

## Webhook Notifications

**Priority**: Low
**Estimated Effort**: 2-3 days

### Description
Send webhook notifications on push completion or failure.

### Implementation Notes
```go
// TODO: Webhook notifications for push events
// - Configurable webhook endpoints
// - Payload templates
// - Retry logic for failed webhooks
// - Support for Slack, Discord, generic HTTP

type WebhookNotifier struct {
    endpoints []WebhookEndpoint
    template  *template.Template
}

type WebhookEndpoint struct {
    URL     string
    Type    string // "slack", "discord", "generic"
    Headers map[string]string
}

func (wn *WebhookNotifier) NotifyPushComplete(result *PushResult) error {
    // Send notification to all configured endpoints
}
```

## Implementation Priority

1. **High Priority** (Next Sprint):
   - Multi-architecture support
   - Parallel layer uploads
   - Signature verification (cosign)
   - Vulnerability scanning

2. **Medium Priority** (Future Sprints):
   - Rate limiting
   - Resume capability
   - OpenTelemetry integration
   - Prometheus metrics
   - CAS cache

3. **Low Priority** (Backlog):
   - Structured logging enhancements
   - Registry-specific optimizations
   - Webhook notifications

## Contributing

When implementing these enhancements:
1. Follow the established code patterns in `pkg/push/`
2. Maintain backward compatibility
3. Add comprehensive tests
4. Update documentation
5. Consider performance impact
6. Ensure graceful degradation if features are disabled
