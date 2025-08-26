package cache

import (
	"fmt"
	"sync"
	"time"
)

// CacheKeyParams defines parameters for cache key calculation
type CacheKeyParams struct {
	Instruction     string            `json:"instruction"`
	Context         []byte            `json:"context,omitempty"`
	BuildArgs       []BuildArg        `json:"build_args,omitempty"`
	BaseImageDigest string            `json:"base_image_digest"`
	Timestamp       time.Time         `json:"timestamp,omitempty"`
}

type BuildArg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// managerImpl implements the CacheManager interface
type managerImpl struct {
	mu         sync.RWMutex
	db         *layerDB
	calculator *keyCalculator
	strategies []EvictionStrategy
	config     *Config
	stats      CacheStats
	closed     bool
}

// Config holds cache manager configuration
type Config struct {
	BasePath       string        `yaml:"base_path"`
	MaxSize        int64         `yaml:"max_size"`
	MaxLayers      int           `yaml:"max_layers"`
	MaxAge         time.Duration `yaml:"max_age"`
	EnableMetrics  bool          `yaml:"enable_metrics"`
	SyncInterval   time.Duration `yaml:"sync_interval"`
	LowWaterMark   float64       `yaml:"low_water_mark"`
	HighWaterMark  float64       `yaml:"high_water_mark"`
}

func DefaultConfig() *Config {
	return &Config{"/tmp/oci-cache", 10*1024*1024*1024, 10000, 7*24*time.Hour, true, 30*time.Second, 0.9, 0.7}
}

// NewManager creates a new cache manager instance
func NewManager(config *Config) (CacheManager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Initialize layer database
	db, err := newLayerDB(config.BasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize layer database: %w", err)
	}

	// Initialize cache key calculator
	calculator := &keyCalculator{
		includeContext:    true,
		includeBuildArgs:  true,
		includeTimestamps: false,
	}

	// Initialize eviction strategies
	strategies := []EvictionStrategy{
		&LRUStrategy{},
		&TTLStrategy{MaxAge: config.MaxAge},
		&SizeStrategy{MaxSize: config.MaxSize},
	}

	manager := &managerImpl{
		db:         db,
		calculator: calculator,
		strategies: strategies,
		config:     config,
		stats:      CacheStats{},
		closed:     false,
	}

	return manager, nil
}

// HasLayer checks if a layer exists in the cache
func (m *managerImpl) HasLayer(digest string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return false
	}

	has := m.db.HasLayer(digest)
	
	// Update access statistics
	if has {
		m.stats.HitCount++
	} else {
		m.stats.MissCount++
	}
	
	// Recalculate hit ratio
	total := m.stats.HitCount + m.stats.MissCount
	if total > 0 {
		m.stats.HitRate = float64(m.stats.HitCount) / float64(total)
	}

	return has
}

// GetLayer retrieves a layer from the cache
func (m *managerImpl) GetLayer(digest string) (*Layer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, fmt.Errorf("cache manager is closed")
	}

	layer, err := m.db.GetLayer(digest)
	if err != nil {
		m.stats.MissCount++
		return nil, fmt.Errorf("layer not found in cache: %w", err)
	}

	// Update access statistics and metadata
	m.stats.HitCount++
	m.db.UpdateAccess(digest)

	// Recalculate hit ratio
	total := m.stats.HitCount + m.stats.MissCount
	if total > 0 {
		m.stats.HitRate = float64(m.stats.HitCount) / float64(total)
	}

	return layer, nil
}

// StoreLayer stores a layer in the cache
func (m *managerImpl) StoreLayer(layer *Layer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("cache manager is closed")
	}

	// Check if we need to evict layers before storing
	if m.shouldEvict() {
		if err := m.evictLayers(); err != nil {
			return fmt.Errorf("failed to evict layers: %w", err)
		}
	}

	// Store the layer
	if err := m.db.StoreLayer(layer); err != nil {
		return fmt.Errorf("failed to store layer: %w", err)
	}

	// Update statistics
	m.stats.TotalSize += layer.Size
	m.stats.TotalLayers++

	return nil
}

// CalculateCacheKey generates a cache key for the given parameters
func (m *managerImpl) CalculateCacheKey(instruction string, context []byte) string {
	params := CacheKeyParams{
		Instruction: instruction,
		Context:     context,
	}
	
	return m.calculator.Calculate(params)
}

// PruneCache removes layers older than the specified time
func (m *managerImpl) PruneCache(before time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("cache manager is closed")
	}

	removedSize, removedCount, err := m.db.PruneLayers(before)
	if err != nil {
		return fmt.Errorf("failed to prune cache: %w", err)
	}

	// Update statistics
	m.stats.TotalSize -= removedSize
	m.stats.TotalLayers -= int64(removedCount)
	m.stats.PruneCount += int64(removedCount)

	return nil
}

// GetStats returns current cache statistics
func (m *managerImpl) GetStats() *CacheStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	return &CacheStats{
		HitCount:      m.stats.HitCount,
		MissCount:     m.stats.MissCount,
		HitRate:       m.stats.HitRate,
		TotalSize:     m.stats.TotalSize,
		TotalLayers:   m.stats.TotalLayers,
	}
}

// Close shuts down the cache manager and releases resources
func (m *managerImpl) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	m.closed = true

	if m.db != nil {
		if err := m.db.Close(); err != nil {
			return fmt.Errorf("failed to close layer database: %w", err)
		}
	}

	return nil
}

// shouldEvict determines if cache eviction is needed
func (m *managerImpl) shouldEvict() bool {
	// Check size-based eviction
	if m.config.MaxSize > 0 && m.stats.TotalSize >= int64(float64(m.config.MaxSize)*m.config.LowWaterMark) {
		return true
	}

	// Check layer count-based eviction
	if m.config.MaxLayers > 0 && m.stats.TotalLayers >= int64(float64(m.config.MaxLayers)*m.config.LowWaterMark) {
		return true
	}

	return false
}

// evictLayers performs cache eviction based on configured strategies
func (m *managerImpl) evictLayers() error {
	layers := m.db.GetAllLayers()
	
	// Apply eviction strategies
	var layersToEvict []*layerMetadata
	
	for _, layer := range layers {
		for _, strategy := range m.strategies {
			if strategy.ShouldEvict(layer, &EvictionConfig{
				MaxCacheSize:  m.config.MaxSize,
				MaxLayerAge:   m.config.MaxAge,
				MinRefCount:   1,
				LowWaterMark:  m.config.LowWaterMark,
				HighWaterMark: m.config.HighWaterMark,
			}) {
				layersToEvict = append(layersToEvict, layer)
				break
			}
		}
	}

	// Remove selected layers
	var removedSize int64
	for _, layer := range layersToEvict {
		if err := m.db.RemoveLayer(layer.Layer.Digest.String()); err != nil {
			return fmt.Errorf("failed to remove layer %s: %w", layer.Layer.Digest.String(), err)
		}
		removedSize += layer.Layer.Size
		m.stats.PruneCount++
		
		// Stop evicting if we've reached the high water mark
		newSize := m.stats.TotalSize - removedSize
		if float64(newSize) <= float64(m.config.MaxSize)*m.config.HighWaterMark {
			break
		}
	}

	// Update statistics
	m.stats.TotalSize -= removedSize
	m.stats.TotalLayers -= int64(len(layersToEvict))

	return nil
}