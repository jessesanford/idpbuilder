package buildah

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// CacheManager handles build cache storage and retrieval
type CacheManager struct {
	cacheDir   string
	maxSize    int64 // Maximum cache size in bytes
	ttl        time.Duration
	layerCache map[string]*LayerCache
}

// LayerCache represents a cached build layer
type LayerCache struct {
	ID           string            `json:"id"`
	Digest       string            `json:"digest"`
	CreatedAt    time.Time         `json:"created_at"`
	LastUsed     time.Time         `json:"last_used"`
	Size         int64             `json:"size"`
	Dependencies []string          `json:"dependencies"`
	Metadata     map[string]string `json:"metadata"`
}

// CacheStats provides cache usage statistics
type CacheStats struct {
	TotalSize   int64     `json:"total_size"`
	EntryCount  int       `json:"entry_count"`
	HitRate     float64   `json:"hit_rate"`
	OldestEntry time.Time `json:"oldest_entry"`
}

// CacheOption is a functional option for cache configuration
type CacheOption func(*CacheManager)

// WithMaxSize sets the maximum cache size
func WithMaxSize(size int64) CacheOption {
	return func(cm *CacheManager) {
		cm.maxSize = size
	}
}

// WithTTL sets the cache time-to-live
func WithTTL(ttl time.Duration) CacheOption {
	return func(cm *CacheManager) {
		cm.ttl = ttl
	}
}

// NewCacheManager creates a new cache manager with configuration
func NewCacheManager(cacheDir string, options ...CacheOption) (*CacheManager, error) {
	cm := &CacheManager{
		cacheDir:   cacheDir,
		maxSize:    1024 * 1024 * 1024, // 1GB default
		ttl:        24 * time.Hour,      // 24 hours default
		layerCache: make(map[string]*LayerCache),
	}

	// Apply options
	for _, opt := range options {
		opt(cm)
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Load existing cache metadata
	if err := cm.loadCacheMetadata(); err != nil {
		return nil, fmt.Errorf("failed to load cache metadata: %w", err)
	}

	return cm, nil
}

// GenerateCacheKey creates a unique key for cache lookup
func (cm *CacheManager) GenerateCacheKey(dockerfile string, buildContext string, args map[string]string) string {
	hasher := sha256.New()
	hasher.Write([]byte(dockerfile))
	hasher.Write([]byte(buildContext))

	// Sort args for consistent ordering
	var keys []string
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		hasher.Write([]byte(k + "=" + args[k]))
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

// StoreLayer stores a build layer in cache
func (cm *CacheManager) StoreLayer(ctx context.Context, key string, layerData []byte) error {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	// Create cache entry first
	entry := &LayerCache{
		ID:        key,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
		Size:      int64(len(layerData)),
		Metadata:  make(map[string]string),
	}

	// Store in memory cache
	cm.layerCache[key] = entry

	// Check cache size limits and evict if necessary
	for cm.calculateTotalSize() > cm.maxSize {
		if err := cm.EvictOldEntries(); err != nil {
			return fmt.Errorf("failed to evict old entries: %w", err)
		}
		// If we couldn't evict anything, break to avoid infinite loop
		if len(cm.layerCache) <= 1 {
			break
		}
	}

	// Write layer data to cache directory
	layerPath := filepath.Join(cm.cacheDir, key+".layer")
	if err := os.WriteFile(layerPath, layerData, 0600); err != nil {
		return fmt.Errorf("failed to write layer data: %w", err)
	}

	// Save cache metadata
	return cm.saveCacheMetadata()
}

// GetLayer retrieves a cached layer if available
func (cm *CacheManager) GetLayer(ctx context.Context, key string) (*LayerCache, error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Check if key exists in cache
	entry, exists := cm.layerCache[key]
	if !exists {
		return nil, nil // Cache miss
	}

	// Validate cache entry is not expired
	if time.Since(entry.CreatedAt) > cm.ttl {
		// Remove expired entry
		delete(cm.layerCache, key)
		layerPath := filepath.Join(cm.cacheDir, key+".layer")
		os.Remove(layerPath)
		return nil, nil
	}

	// Update last used timestamp
	entry.LastUsed = time.Now()
	cm.layerCache[key] = entry

	// Save updated metadata
	cm.saveCacheMetadata()

	return entry, nil
}

// HasValidCache checks if valid cache exists for the key
func (cm *CacheManager) HasValidCache(key string) bool {
	entry, exists := cm.layerCache[key]
	if !exists {
		return false
	}

	// Check TTL
	if time.Since(entry.CreatedAt) > cm.ttl {
		return false
	}

	// Check if layer file exists
	layerPath := filepath.Join(cm.cacheDir, key+".layer")
	if _, err := os.Stat(layerPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// EvictOldEntries removes expired cache entries
func (cm *CacheManager) EvictOldEntries() error {
	var toRemove []string

	// Find expired entries
	now := time.Now()
	for key, entry := range cm.layerCache {
		if now.Sub(entry.CreatedAt) > cm.ttl {
			toRemove = append(toRemove, key)
		}
	}

	// Remove expired entries
	for _, key := range toRemove {
		delete(cm.layerCache, key)
		layerPath := filepath.Join(cm.cacheDir, key+".layer")
		os.Remove(layerPath)
	}

	// If still over size limit, remove oldest entries
	for cm.calculateTotalSize() > cm.maxSize && len(cm.layerCache) > 0 {
		oldest := cm.findOldestEntry()
		if oldest != "" {
			delete(cm.layerCache, oldest)
			layerPath := filepath.Join(cm.cacheDir, oldest+".layer")
			os.Remove(layerPath)
		}
	}

	return cm.saveCacheMetadata()
}

// InvalidateCache removes specific cache entries
func (cm *CacheManager) InvalidateCache(patterns ...string) error {
	var toRemove []string

	for key := range cm.layerCache {
		for _, pattern := range patterns {
			if matched, _ := filepath.Match(pattern, key); matched {
				toRemove = append(toRemove, key)
				break
			}
		}
	}

	// Remove matching entries
	for _, key := range toRemove {
		delete(cm.layerCache, key)
		layerPath := filepath.Join(cm.cacheDir, key+".layer")
		os.Remove(layerPath)
	}

	return cm.saveCacheMetadata()
}

// ClearAllCache removes all cached data
func (cm *CacheManager) ClearAllCache() error {
	// Remove all cache files
	entries, err := os.ReadDir(cm.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".layer") {
			os.Remove(filepath.Join(cm.cacheDir, entry.Name()))
		}
	}

	// Reset cache metadata
	cm.layerCache = make(map[string]*LayerCache)

	return cm.saveCacheMetadata()
}

// GetCacheStats returns cache usage statistics
func (cm *CacheManager) GetCacheStats() CacheStats {
	stats := CacheStats{
		TotalSize:  cm.calculateTotalSize(),
		EntryCount: len(cm.layerCache),
		HitRate:    0.0, // Would need to track hits/misses separately
	}

	// Find oldest entry
	var oldest time.Time
	for _, entry := range cm.layerCache {
		if oldest.IsZero() || entry.CreatedAt.Before(oldest) {
			oldest = entry.CreatedAt
		}
	}
	stats.OldestEntry = oldest

	return stats
}

// Helper methods

func (cm *CacheManager) loadCacheMetadata() error {
	metadataPath := filepath.Join(cm.cacheDir, "metadata.json")
	data, err := os.ReadFile(metadataPath)
	if os.IsNotExist(err) {
		return nil // No metadata yet, start fresh
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cm.layerCache)
}

func (cm *CacheManager) saveCacheMetadata() error {
	metadataPath := filepath.Join(cm.cacheDir, "metadata.json")
	data, err := json.MarshalIndent(cm.layerCache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(metadataPath, data, 0600)
}

func (cm *CacheManager) calculateTotalSize() int64 {
	var total int64
	for _, entry := range cm.layerCache {
		total += entry.Size
	}
	return total
}

func (cm *CacheManager) findOldestEntry() string {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range cm.layerCache {
		if oldestKey == "" || entry.LastUsed.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.LastUsed
		}
	}

	return oldestKey
}