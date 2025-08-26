package cache

import (
	"context"
	"time"
)

// Layer represents a cached OCI image layer with metadata
// Compatible with effort1-contracts/pkg/oci/api/Layer
type Layer struct {
	Digest           string            `json:"digest" yaml:"digest"`
	Size             int64             `json:"size" yaml:"size"`
	MediaType        string            `json:"media_type" yaml:"media_type"`
	Created          time.Time         `json:"created" yaml:"created"`
	LastUsed         time.Time         `json:"last_used" yaml:"last_used"`
	RefCount         int               `json:"ref_count" yaml:"ref_count"`
	Annotations      map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	StoragePath      string            `json:"storage_path" yaml:"storage_path"`
	Compressed       bool              `json:"compressed" yaml:"compressed"`
}

// CacheManager handles layer caching operations
// Compatible with effort1-contracts/pkg/oci/api/CacheManager
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
// Compatible with effort1-contracts/pkg/oci/api/CacheStats
type CacheStats struct {
	TotalLayers int64   `json:"total_layers" yaml:"total_layers"`
	TotalSize   int64   `json:"total_size" yaml:"total_size"`
	HitCount    int64   `json:"hit_count" yaml:"hit_count"`
	MissCount   int64   `json:"miss_count" yaml:"miss_count"`
	HitRate     float64 `json:"hit_rate" yaml:"hit_rate"`
	PruneCount  int64   `json:"prune_count" yaml:"prune_count"`
}