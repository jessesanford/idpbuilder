package cache

import (
	"time"
	"github.com/opencontainers/go-digest"
)

// Layer represents a cached OCI image layer
type Layer struct {
	Digest      digest.Digest     `json:"digest"`
	Size        int64             `json:"size"`
	MediaType   string            `json:"media_type"`
	Data        []byte            `json:"-"`
	Created     time.Time         `json:"created"`
	LastUsed    time.Time         `json:"last_used"`
	RefCount    int               `json:"ref_count"`
	Annotations map[string]string `json:"annotations,omitempty"`
	StoragePath string            `json:"storage_path"`
	Compressed  bool              `json:"compressed"`
}

// CacheManager handles layer caching operations
type CacheManager interface {
	HasLayer(digest string) bool
	GetLayer(digest string) (*Layer, error)
	StoreLayer(layer *Layer) error
	CalculateCacheKey(instruction string, context []byte) string
	PruneCache(before time.Time) error
	GetStats() *CacheStats
	Close() error
}

// CacheStats contains cache performance statistics
type CacheStats struct {
	TotalLayers int64   `json:"total_layers"`
	TotalSize   int64   `json:"total_size"`
	HitCount    int64   `json:"hit_count"`
	MissCount   int64   `json:"miss_count"`
	HitRate     float64 `json:"hit_rate"`
	PruneCount  int64   `json:"prune_count"`
}