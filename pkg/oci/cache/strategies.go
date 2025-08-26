package cache

import (
	"time"
)

// EvictionStrategy defines interface for cache eviction strategies
type EvictionStrategy interface {
	ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool
	Priority(layer *layerMetadata) int
}

// EvictionConfig holds configuration for eviction strategies
type EvictionConfig struct {
	MaxCacheSize  int64         `json:"max_cache_size"`
	MaxLayerAge   time.Duration `json:"max_layer_age"`
	MinRefCount   int           `json:"min_ref_count"`
	LowWaterMark  float64       `json:"low_water_mark"`
	HighWaterMark float64       `json:"high_water_mark"`
}

// LRUStrategy implements Least Recently Used eviction
type LRUStrategy struct{}

// ShouldEvict checks if layer should be evicted based on LRU
func (s *LRUStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	// Always consider for LRU-based eviction if cache is full
	return true
}

// Priority returns priority for LRU (lower = higher priority for eviction)
func (s *LRUStrategy) Priority(layer *layerMetadata) int {
	// Older access = higher priority for eviction (lower number)
	return int(time.Since(layer.LastAccess).Hours())
}

// TTLStrategy implements Time To Live eviction
type TTLStrategy struct {
	MaxAge time.Duration `json:"max_age"`
}

// ShouldEvict checks if layer has exceeded its TTL
func (s *TTLStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	age := time.Since(layer.StoredAt)
	maxAge := s.MaxAge
	if maxAge == 0 {
		maxAge = config.MaxLayerAge
	}
	return age > maxAge
}

// Priority returns priority for TTL (older = higher priority)
func (s *TTLStrategy) Priority(layer *layerMetadata) int {
	return int(time.Since(layer.StoredAt).Hours())
}

// SizeStrategy implements size-based eviction
type SizeStrategy struct {
	MaxSize int64 `json:"max_size"`
}

// ShouldEvict checks if cache size exceeds limits
func (s *SizeStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	// Size strategy doesn't evict individual layers directly
	// It's used to trigger general eviction
	return false
}

// Priority returns priority for size-based eviction (larger = higher priority)
func (s *SizeStrategy) Priority(layer *layerMetadata) int {
	// Larger layers get higher priority for eviction
	return int(layer.Layer.Size / 1024 / 1024) // MB
}

// ReferenceStrategy implements reference count-based eviction
type ReferenceStrategy struct {
	MinRefCount int `json:"min_ref_count"`
}

// ShouldEvict checks if layer has low reference count
func (s *ReferenceStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	minRef := s.MinRefCount
	if minRef == 0 {
		minRef = config.MinRefCount
	}
	return layer.RefCount <= minRef
}

// Priority returns priority for reference-based eviction
func (s *ReferenceStrategy) Priority(layer *layerMetadata) int {
	// Lower ref count = higher priority for eviction
	return -layer.RefCount
}

// CompositeStrategy combines multiple strategies
type CompositeStrategy struct {
	Strategies []EvictionStrategy `json:"strategies"`
	Mode       CompositeMode      `json:"mode"`
}

// CompositeMode defines how strategies are combined
type CompositeMode string

const (
	// CompositeAND requires all strategies to agree on eviction
	CompositeAND CompositeMode = "AND"
	// CompositeOR requires any strategy to agree on eviction
	CompositeOR CompositeMode = "OR"
)

// ShouldEvict applies composite strategy logic
func (s *CompositeStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	if len(s.Strategies) == 0 {
		return false
	}

	switch s.Mode {
	case CompositeAND:
		for _, strategy := range s.Strategies {
			if !strategy.ShouldEvict(layer, config) {
				return false
			}
		}
		return true

	case CompositeOR:
		for _, strategy := range s.Strategies {
			if strategy.ShouldEvict(layer, config) {
				return true
			}
		}
		return false

	default:
		return false
	}
}

// Priority returns combined priority from all strategies
func (s *CompositeStrategy) Priority(layer *layerMetadata) int {
	if len(s.Strategies) == 0 {
		return 0
	}

	totalPriority := 0
	for _, strategy := range s.Strategies {
		totalPriority += strategy.Priority(layer)
	}
	return totalPriority / len(s.Strategies)
}

// NewLRUStrategy creates a new LRU strategy
func NewLRUStrategy() EvictionStrategy {
	return &LRUStrategy{}
}

// NewTTLStrategy creates a new TTL strategy
func NewTTLStrategy(maxAge time.Duration) EvictionStrategy {
	return &TTLStrategy{MaxAge: maxAge}
}

// NewSizeStrategy creates a new size-based strategy
func NewSizeStrategy(maxSize int64) EvictionStrategy {
	return &SizeStrategy{MaxSize: maxSize}
}

// NewReferenceStrategy creates a new reference count strategy
func NewReferenceStrategy(minRefCount int) EvictionStrategy {
	return &ReferenceStrategy{MinRefCount: minRefCount}
}

// NewCompositeStrategy creates a new composite strategy
func NewCompositeStrategy(mode CompositeMode, strategies ...EvictionStrategy) EvictionStrategy {
	return &CompositeStrategy{
		Strategies: strategies,
		Mode:       mode,
	}
}