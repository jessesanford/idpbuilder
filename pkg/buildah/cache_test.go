package buildah

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheManager_GenerateCacheKey(t *testing.T) {
	cm, err := NewCacheManager(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	dockerfile := "FROM alpine\nRUN echo hello"
	buildContext := "/tmp/context"
	args := map[string]string{"ENV": "prod", "VERSION": "1.0"}

	// Test deterministic key generation
	key1 := cm.GenerateCacheKey(dockerfile, buildContext, args)
	key2 := cm.GenerateCacheKey(dockerfile, buildContext, args)
	if key1 != key2 {
		t.Error("Cache key generation should be deterministic")
	}

	// Test different inputs produce different keys
	args2 := map[string]string{"ENV": "dev", "VERSION": "1.0"}
	key3 := cm.GenerateCacheKey(dockerfile, buildContext, args2)
	if key1 == key3 {
		t.Error("Different inputs should produce different cache keys")
	}

	// Test empty args
	key4 := cm.GenerateCacheKey(dockerfile, buildContext, nil)
	key5 := cm.GenerateCacheKey(dockerfile, buildContext, map[string]string{})
	if key4 != key5 {
		t.Error("Nil and empty args should produce same key")
	}
}

func TestCacheManager_StoreAndRetrieve(t *testing.T) {
	cm, err := NewCacheManager(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	key := "test-key"
	layerData := []byte("test layer data")

	// Test storing layer in cache
	err = cm.StoreLayer(ctx, key, layerData)
	if err != nil {
		t.Fatalf("Failed to store layer: %v", err)
	}

	// Test retrieving cached layer
	retrieved, err := cm.GetLayer(ctx, key)
	if err != nil {
		t.Fatalf("Failed to retrieve layer: %v", err)
	}
	if retrieved == nil {
		t.Fatal("Retrieved layer should not be nil")
	}

	// Verify cache entry properties
	if retrieved.ID != key {
		t.Errorf("Expected ID %s, got %s", key, retrieved.ID)
	}
	if retrieved.Size != int64(len(layerData)) {
		t.Errorf("Expected size %d, got %d", len(layerData), retrieved.Size)
	}

	// Test cache hit validation
	if !cm.HasValidCache(key) {
		t.Error("Cache should be valid for stored key")
	}
}

func TestCacheManager_TTLExpiration(t *testing.T) {
	shortTTL := 100 * time.Millisecond
	cm, err := NewCacheManager(t.TempDir(), WithTTL(shortTTL))
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	key := "expire-test"
	layerData := []byte("will expire soon")

	// Store layer
	err = cm.StoreLayer(ctx, key, layerData)
	if err != nil {
		t.Fatalf("Failed to store layer: %v", err)
	}

	// Should be valid immediately
	if !cm.HasValidCache(key) {
		t.Error("Cache should be valid immediately after storage")
	}

	// Wait for expiration
	time.Sleep(shortTTL + 50*time.Millisecond)

	// Should be expired now
	if cm.HasValidCache(key) {
		t.Error("Cache should be expired after TTL")
	}

	// Verify expired entries are not returned
	retrieved, err := cm.GetLayer(ctx, key)
	if err != nil {
		t.Fatalf("GetLayer should not error: %v", err)
	}
	if retrieved != nil {
		t.Error("Should not retrieve expired cache entry")
	}
}

func TestCacheManager_SizeLimit(t *testing.T) {
	smallMaxSize := int64(100) // Very small size to trigger eviction
	cm, err := NewCacheManager(t.TempDir(), WithMaxSize(smallMaxSize))
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()

	// Store multiple layers that exceed size limit
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("large-layer-%d", i)
		largeData := make([]byte, 50) // Each layer is 50 bytes
		err = cm.StoreLayer(ctx, key, largeData)
		if err != nil {
			t.Fatalf("Failed to store layer %d: %v", i, err)
		}
	}

	// Check that cache size management worked
	stats := cm.GetCacheStats()
	if stats.TotalSize > smallMaxSize {
		t.Errorf("Cache size %d exceeds limit %d", stats.TotalSize, smallMaxSize)
	}

	// Some entries should have been evicted
	if stats.EntryCount >= 3 {
		t.Error("Expected some cache entries to be evicted due to size limit")
	}
}

func TestCacheManager_Invalidation(t *testing.T) {
	cm, err := NewCacheManager(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()

	// Store multiple layers with different patterns
	keys := []string{"app-layer-1", "app-layer-2", "base-layer-1", "temp-layer-1"}
	for _, key := range keys {
		err = cm.StoreLayer(ctx, key, []byte("data"))
		if err != nil {
			t.Fatalf("Failed to store layer %s: %v", key, err)
		}
	}

	// Test pattern-based invalidation
	err = cm.InvalidateCache("app-*")
	if err != nil {
		t.Fatalf("Failed to invalidate cache: %v", err)
	}

	// Check that app-* entries are invalidated
	for _, key := range []string{"app-layer-1", "app-layer-2"} {
		if cm.HasValidCache(key) {
			t.Errorf("Cache key %s should be invalidated", key)
		}
	}

	// Check that other entries still exist
	for _, key := range []string{"base-layer-1", "temp-layer-1"} {
		if !cm.HasValidCache(key) {
			t.Errorf("Cache key %s should still be valid", key)
		}
	}

	// Test clearing all cache
	err = cm.ClearAllCache()
	if err != nil {
		t.Fatalf("Failed to clear all cache: %v", err)
	}

	stats := cm.GetCacheStats()
	if stats.EntryCount != 0 {
		t.Error("All cache entries should be cleared")
	}
}

func TestCacheManager_ConcurrentAccess(t *testing.T) {
	cm, err := NewCacheManager(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	numGoroutines := 10

	// Test concurrent writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("concurrent-key-%d", id)
			data := []byte(fmt.Sprintf("data-%d", id))
			err := cm.StoreLayer(ctx, key, data)
			if err != nil {
				t.Errorf("Concurrent store failed for key %s: %v", key, err)
			}
		}(i)
	}
	wg.Wait()

	// Test concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("concurrent-key-%d", id)
			layer, err := cm.GetLayer(ctx, key)
			if err != nil {
				t.Errorf("Concurrent get failed for key %s: %v", key, err)
			}
			if layer == nil {
				t.Errorf("Layer should exist for key %s", key)
			}
		}(i)
	}
	wg.Wait()

	// Verify final state
	stats := cm.GetCacheStats()
	if stats.EntryCount != numGoroutines {
		t.Errorf("Expected %d entries, got %d", numGoroutines, stats.EntryCount)
	}
}

func TestCacheManager_ContextCancellation(t *testing.T) {
	cm, err := NewCacheManager(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Test that cancelled context is handled
	err = cm.StoreLayer(ctx, "test", []byte("data"))
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}

	layer, err := cm.GetLayer(ctx, "test")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
	if layer != nil {
		t.Error("Layer should be nil with cancelled context")
	}
}

func TestCacheManager_ErrorConditions(t *testing.T) {
	// Test with invalid cache directory
	_, err := NewCacheManager("/invalid/path/that/cannot/be/created")
	if err == nil {
		t.Error("Should fail to create cache manager with invalid path")
	}

	// Test cache miss scenario
	cm, err := NewCacheManager(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	layer, err := cm.GetLayer(ctx, "nonexistent-key")
	if err != nil {
		t.Errorf("GetLayer should not error for cache miss: %v", err)
	}
	if layer != nil {
		t.Error("Layer should be nil for cache miss")
	}

	if cm.HasValidCache("nonexistent-key") {
		t.Error("Should not have valid cache for nonexistent key")
	}
}