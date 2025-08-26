package api

import (
	"context"
	"testing"
	"time"
)

// MockCacheManager implements CacheManager for testing
type MockCacheManager struct {
	HasLayerFunc         func(string) bool
	GetLayerFunc         func(string) (*Layer, error)
	StoreLayerFunc       func(*Layer) error
	CalculateCacheKeyFunc func(string, []byte) string
	PruneCacheFunc       func(time.Time) error
	GetStatsFunc         func() *CacheStats
	ValidateCacheFunc    func(context.Context) (*ValidationReport, error)
	WarmCacheFunc        func(context.Context, []string) error
	SetCachePolicyFunc   func(*CachePolicy) error
	GetCachePolicyFunc   func() *CachePolicy
}

func (m *MockCacheManager) HasLayer(digest string) bool {
	if m.HasLayerFunc != nil {
		return m.HasLayerFunc(digest)
	}
	return false
}

func (m *MockCacheManager) GetLayer(digest string) (*Layer, error) {
	if m.GetLayerFunc != nil {
		return m.GetLayerFunc(digest)
	}
	return &Layer{}, nil
}

func (m *MockCacheManager) StoreLayer(layer *Layer) error {
	if m.StoreLayerFunc != nil {
		return m.StoreLayerFunc(layer)
	}
	return nil
}

func (m *MockCacheManager) CalculateCacheKey(instruction string, context []byte) string {
	if m.CalculateCacheKeyFunc != nil {
		return m.CalculateCacheKeyFunc(instruction, context)
	}
	return "mock-cache-key"
}

func (m *MockCacheManager) PruneCache(before time.Time) error {
	if m.PruneCacheFunc != nil {
		return m.PruneCacheFunc(before)
	}
	return nil
}

func (m *MockCacheManager) GetStats() *CacheStats {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc()
	}
	return &CacheStats{}
}

func (m *MockCacheManager) ValidateCache(ctx context.Context) (*ValidationReport, error) {
	if m.ValidateCacheFunc != nil {
		return m.ValidateCacheFunc(ctx)
	}
	return &ValidationReport{}, nil
}

func (m *MockCacheManager) WarmCache(ctx context.Context, images []string) error {
	if m.WarmCacheFunc != nil {
		return m.WarmCacheFunc(ctx, images)
	}
	return nil
}

func (m *MockCacheManager) SetCachePolicy(policy *CachePolicy) error {
	if m.SetCachePolicyFunc != nil {
		return m.SetCachePolicyFunc(policy)
	}
	return nil
}

func (m *MockCacheManager) GetCachePolicy() *CachePolicy {
	if m.GetCachePolicyFunc != nil {
		return m.GetCachePolicyFunc()
	}
	return &CachePolicy{}
}

// Test interface compliance
func TestCacheManagerInterface(t *testing.T) {
	var _ CacheManager = &MockCacheManager{}
}

func TestLayer(t *testing.T) {
	now := time.Now()
	layer := &Layer{
		Digest:           "sha256:abc123",
		Size:             1024,
		MediaType:        "application/vnd.oci.image.layer.v1.tar+gzip",
		Created:          now,
		LastUsed:         now,
		RefCount:         1,
		Annotations:      map[string]string{"org.example.label": "value"},
		Instruction:      "RUN apt-get update",
		BaseImage:        "ubuntu:20.04",
		StoragePath:      "/cache/layers/abc123",
		Compressed:       true,
		VerificationHash: "md5:def456",
	}

	if layer.Digest != "sha256:abc123" {
		t.Errorf("Expected digest 'sha256:abc123', got '%s'", layer.Digest)
	}
	if layer.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", layer.Size)
	}
	if !layer.Compressed {
		t.Error("Expected layer to be compressed")
	}
}