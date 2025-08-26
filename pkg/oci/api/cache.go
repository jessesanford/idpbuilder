package api

import (
	"context"
	"time"
)

// CacheManager handles layer caching operations for OCI builds.
type CacheManager interface {
	// HasLayer checks if a layer with the given digest exists in cache.
	HasLayer(digest string) bool

	// GetLayer retrieves a cached layer by digest.
	GetLayer(digest string) (*Layer, error)

	// StoreLayer stores a new layer in the cache.
	StoreLayer(layer *Layer) error

	// CalculateCacheKey generates a cache key for a build instruction.
	CalculateCacheKey(instruction string, context []byte) string

	// PruneCache removes old cache entries based on age and usage.
	PruneCache(before time.Time) error

	// GetStats returns comprehensive cache statistics.
	GetStats() *CacheStats

	// ValidateCache performs integrity checks on cached layers.
	ValidateCache(ctx context.Context) (*ValidationReport, error)

	// WarmCache pre-loads commonly used base images and layers.
	WarmCache(ctx context.Context, images []string) error

	// SetCachePolicy configures cache behavior and retention policies.
	SetCachePolicy(policy *CachePolicy) error

	// GetCachePolicy returns the current cache configuration.
	GetCachePolicy() *CachePolicy
}

// Layer represents a cached OCI image layer with metadata.
type Layer struct {
	Digest           string            `json:"digest" yaml:"digest"`
	Size             int64             `json:"size" yaml:"size"`
	MediaType        string            `json:"media_type" yaml:"media_type"`
	Created          time.Time         `json:"created" yaml:"created"`
	LastUsed         time.Time         `json:"last_used" yaml:"last_used"`
	RefCount         int               `json:"ref_count" yaml:"ref_count"`
	Annotations      map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Instruction      string            `json:"instruction,omitempty" yaml:"instruction,omitempty"`
	BaseImage        string            `json:"base_image,omitempty" yaml:"base_image,omitempty"`
	StoragePath      string            `json:"storage_path" yaml:"storage_path"`
	Compressed       bool              `json:"compressed" yaml:"compressed"`
	VerificationHash string            `json:"verification_hash,omitempty" yaml:"verification_hash,omitempty"`
}

// CacheStats contains comprehensive cache performance statistics.
type CacheStats struct {
	TotalLayers      int64            `json:"total_layers" yaml:"total_layers"`
	TotalSize        int64            `json:"total_size" yaml:"total_size"`
	HitCount         int64            `json:"hit_count" yaml:"hit_count"`
	MissCount        int64            `json:"miss_count" yaml:"miss_count"`
	HitRate          float64          `json:"hit_rate" yaml:"hit_rate"`
	AverageLayerSize int64            `json:"average_layer_size" yaml:"average_layer_size"`
	OldestLayer      time.Time        `json:"oldest_layer" yaml:"oldest_layer"`
	NewestLayer      time.Time        `json:"newest_layer" yaml:"newest_layer"`
	PruneCount       int64            `json:"prune_count" yaml:"prune_count"`
	LastPruned       time.Time        `json:"last_pruned" yaml:"last_pruned"`
	DiskUsage        *DiskUsageStats  `json:"disk_usage,omitempty" yaml:"disk_usage,omitempty"`
}

// DiskUsageStats provides detailed storage utilization metrics.
type DiskUsageStats struct {
	Available       int64   `json:"available" yaml:"available"`
	Used            int64   `json:"used" yaml:"used"`
	Total           int64   `json:"total" yaml:"total"`
	UsagePercentage float64 `json:"usage_percentage" yaml:"usage_percentage"`
}

// ValidationReport contains the results of cache integrity validation.
type ValidationReport struct {
	Timestamp       time.Time `json:"timestamp" yaml:"timestamp"`
	TotalLayers     int64     `json:"total_layers" yaml:"total_layers"`
	ValidLayers     int64     `json:"valid_layers" yaml:"valid_layers"`
	CorruptedLayers []string  `json:"corrupted_layers,omitempty" yaml:"corrupted_layers,omitempty"`
	MissingLayers   []string  `json:"missing_layers,omitempty" yaml:"missing_layers,omitempty"`
	Errors          []string  `json:"errors,omitempty" yaml:"errors,omitempty"`
	IsHealthy       bool      `json:"is_healthy" yaml:"is_healthy"`
}

// CachePolicy defines cache behavior and retention policies.
type CachePolicy struct {
	MaxSize              int64         `json:"max_size" yaml:"max_size"`
	MaxAge               time.Duration `json:"max_age" yaml:"max_age"`
	MaxLayers            int64         `json:"max_layers" yaml:"max_layers"`
	PruneThreshold       float64       `json:"prune_threshold" yaml:"prune_threshold"`
	PruneInterval        time.Duration `json:"prune_interval" yaml:"prune_interval"`
	PreferredCompression string        `json:"preferred_compression" yaml:"preferred_compression"`
	EnableIntegrityChecks bool         `json:"enable_integrity_checks" yaml:"enable_integrity_checks"`
}