package cache

import (
	"time"
)

type EvictionStrategy interface {
	ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool
	Priority(layer *layerMetadata) int
}

type EvictionConfig struct {
	MaxCacheSize  int64         `json:"max_cache_size"`
	MaxLayerAge   time.Duration `json:"max_layer_age"`
	MinRefCount   int           `json:"min_ref_count"`
	LowWaterMark  float64       `json:"low_water_mark"`
	HighWaterMark float64       `json:"high_water_mark"`
}

type LRUStrategy struct{}

func (s *LRUStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	return true
}

func (s *LRUStrategy) Priority(layer *layerMetadata) int {
	return int(time.Since(layer.LastAccess).Hours())
}

type TTLStrategy struct {
	MaxAge time.Duration `json:"max_age"`
}

func (s *TTLStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	age := time.Since(layer.StoredAt)
	maxAge := s.MaxAge
	if maxAge == 0 {
		maxAge = config.MaxLayerAge
	}
	return age > maxAge
}

func (s *TTLStrategy) Priority(layer *layerMetadata) int {
	return int(time.Since(layer.StoredAt).Hours())
}

type SizeStrategy struct {
	MaxSize int64 `json:"max_size"`
}

func (s *SizeStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	// Size strategy doesn't evict individual layers directly
	// It's used to trigger general eviction
	return false
}

func (s *SizeStrategy) Priority(layer *layerMetadata) int {
	// Larger layers get higher priority for eviction
	return int(layer.Layer.Size / 1024 / 1024) // MB
}

type ReferenceStrategy struct {
	MinRefCount int `json:"min_ref_count"`
}

func (s *ReferenceStrategy) ShouldEvict(layer *layerMetadata, config *EvictionConfig) bool {
	minRef := s.MinRefCount
	if minRef == 0 {
		minRef = config.MinRefCount
	}
	return layer.RefCount <= minRef
}

func (s *ReferenceStrategy) Priority(layer *layerMetadata) int {
	// Lower ref count = higher priority for eviction
	return -layer.RefCount
}

type CompositeStrategy struct {
	Strategies []EvictionStrategy `json:"strategies"`
	Mode       CompositeMode      `json:"mode"`
}

type CompositeMode string

const (
	// CompositeAND requires all strategies to agree on eviction
	CompositeAND CompositeMode = "AND"
	// CompositeOR requires any strategy to agree on eviction
	CompositeOR CompositeMode = "OR"
)

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

func NewLRUStrategy() EvictionStrategy {
	return &LRUStrategy{}
}

func NewTTLStrategy(maxAge time.Duration) EvictionStrategy {
	return &TTLStrategy{MaxAge: maxAge}
}

func NewSizeStrategy(maxSize int64) EvictionStrategy {
	return &SizeStrategy{MaxSize: maxSize}
}

func NewReferenceStrategy(minRefCount int) EvictionStrategy {
	return &ReferenceStrategy{MinRefCount: minRefCount}
}

func NewCompositeStrategy(mode CompositeMode, strategies ...EvictionStrategy) EvictionStrategy {
	return &CompositeStrategy{
		Strategies: strategies,
		Mode:       mode,
	}
}