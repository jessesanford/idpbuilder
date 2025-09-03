package builder

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
)

func TestNewBuilder(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid builder creation
	builder, err := NewBuilder(config)
	if err != nil {
		t.Fatalf("unexpected error creating builder: %v", err)
	}
	
	if builder == nil {
		t.Fatal("builder should not be nil")
	}
	
	if builder.config == nil {
		t.Fatal("builder config should not be nil")
	}
	
	if builder.cache == nil {
		t.Fatal("builder cache should not be nil")
	}
	
	if builder.registry == nil {
		t.Fatal("builder registry should not be nil")
	}
}

func TestNewBuilderWithNilConfig(t *testing.T) {
	_, err := NewBuilder(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}
	
	expectedMsg := "build config cannot be nil"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestNewBuilderWithInvalidConfig(t *testing.T) {
	config := &BuildConfig{
		ContextPath:  "", // Invalid - empty
		Dockerfile:   "Dockerfile",
		Tags:         []string{"test:latest"},
		Platform:     DefaultPlatformConfig(),
		Registry:     DefaultRegistryConfig(),
		BuildTimeout: 30 * time.Minute,
	}
	
	_, err := NewBuilder(config)
	if err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestBuilderOptions(t *testing.T) {
	config := DefaultBuildConfig()
	cache := NewMemoryCache()
	registry := NewMockRegistryClient()
	baseImage := empty.Image
	
	builder, err := NewBuilder(config, 
		WithCache(cache),
		WithRegistry(registry),
		WithBaseImage(baseImage),
	)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if builder.cache != cache {
		t.Error("cache option was not applied")
	}
	
	if builder.registry != registry {
		t.Error("registry option was not applied")
	}
	
	if builder.baseImage != baseImage {
		t.Error("base image option was not applied")
	}
}

func TestBuilderOptionsWithNilValues(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test nil cache
	_, err := NewBuilder(config, WithCache(nil))
	if err == nil {
		t.Error("expected error for nil cache")
	}
	
	// Test nil registry
	_, err = NewBuilder(config, WithRegistry(nil))
	if err == nil {
		t.Error("expected error for nil registry")
	}
	
	// Test nil base image
	_, err = NewBuilder(config, WithBaseImage(nil))
	if err == nil {
		t.Error("expected error for nil base image")
	}
}

func TestBuilderGetters(t *testing.T) {
	config := DefaultBuildConfig()
	config.ContextPath = "/test/context"
	
	builder, err := NewBuilder(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Test GetBuildContext
	context := builder.GetBuildContext()
	if context != "/test/context" {
		t.Errorf("expected build context '/test/context', got '%s'", context)
	}
	
	// Test GetConfig
	configCopy := builder.GetConfig()
	if configCopy == nil {
		t.Fatal("config copy should not be nil")
	}
	
	if configCopy == builder.config {
		t.Error("config should be a copy, not the same instance")
	}
	
	if configCopy.ContextPath != config.ContextPath {
		t.Error("config copy should have same values")
	}
	
	// Test GetLayers
	layers := builder.GetLayers()
	if layers == nil {
		t.Fatal("layers should not be nil")
	}
	
	if len(layers) != 0 {
		t.Errorf("expected 0 layers, got %d", len(layers))
	}
}

func TestBuilderAddLayer(t *testing.T) {
	config := DefaultBuildConfig()
	builder, err := NewBuilder(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Test adding a valid layer
	layer := NewEmptyLayer()
	err = builder.AddLayer(layer)
	if err != nil {
		t.Errorf("unexpected error adding layer: %v", err)
	}
	
	layers := builder.GetLayers()
	if len(layers) != 1 {
		t.Errorf("expected 1 layer, got %d", len(layers))
	}
	
	// Test adding nil layer
	err = builder.AddLayer(nil)
	if err == nil {
		t.Error("expected error for nil layer")
	}
}

func TestBuilderSetImageConfig(t *testing.T) {
	config := DefaultBuildConfig()
	builder, err := NewBuilder(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	imageConfig := v1.Config{
		User: "testuser",
		Env:  []string{"TEST=value"},
	}
	
	builder.SetImageConfig(imageConfig)
	
	if builder.imageConfig.User != "testuser" {
		t.Error("image config was not set correctly")
	}
}

func TestBuilderClose(t *testing.T) {
	config := DefaultBuildConfig()
	builder, err := NewBuilder(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	err = builder.Close()
	if err != nil {
		t.Errorf("unexpected error closing builder: %v", err)
	}
}

func TestBuilderBuildWithProgress(t *testing.T) {
	config := DefaultBuildConfig()
	config.ContextPath = "."  // Use current directory as context
	
	// Use mock registry to avoid network calls
	mockRegistry := NewMockRegistryClient()
	
	builder, err := NewBuilder(config, WithRegistry(mockRegistry))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Create a progress writer to capture output
	var progressOutput strings.Builder
	
	ctx := context.Background()
	result, err := builder.BuildWithProgress(ctx, &progressOutput)
	
	// The build might fail due to missing Dockerfile, but we're testing progress output
	if result != nil {
		if result.BuildTime <= 0 {
			t.Error("build time should be positive")
		}
		
		if len(result.Tags) == 0 {
			t.Error("result should include tags")
		}
	}
	
	// Check that progress was written
	progressText := progressOutput.String()
	if !strings.Contains(progressText, "Starting build") {
		t.Error("progress output should contain 'Starting build'")
	}
}

func TestBuildResult(t *testing.T) {
	result := &BuildResult{
		Image:      empty.Image,
		Digest:     v1.Hash{},
		Size:       1024,
		LayerCount: 3,
		BuildTime:  5 * time.Second,
		Tags:       []string{"test:latest", "test:v1.0"},
	}
	
	if result.Size != 1024 {
		t.Errorf("expected size 1024, got %d", result.Size)
	}
	
	if result.LayerCount != 3 {
		t.Errorf("expected 3 layers, got %d", result.LayerCount)
	}
	
	if result.BuildTime != 5*time.Second {
		t.Errorf("expected 5s build time, got %v", result.BuildTime)
	}
	
	if len(result.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result.Tags))
	}
}

func TestBuilderPush(t *testing.T) {
	config := DefaultBuildConfig()
	config.Tags = []string{"test:latest", "test:v1.0"}
	
	mockRegistry := NewMockRegistryClient()
	
	builder, err := NewBuilder(config, WithRegistry(mockRegistry))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	ctx := context.Background()
	err = builder.Push(ctx, empty.Image)
	if err != nil {
		t.Errorf("unexpected error pushing: %v", err)
	}
	
	pushes := mockRegistry.GetPushes()
	if len(pushes) != 2 {
		t.Errorf("expected 2 pushes, got %d", len(pushes))
	}
	
	expectedTags := []string{"test:latest", "test:v1.0"}
	for _, expectedTag := range expectedTags {
		found := false
		for _, pushed := range pushes {
			if pushed == expectedTag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected tag %s to be pushed", expectedTag)
		}
	}
}

func TestBuilderPushWithNoRegistry(t *testing.T) {
	config := DefaultBuildConfig()
	
	builder, err := NewBuilder(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Clear the registry
	builder.registry = nil
	
	ctx := context.Background()
	err = builder.Push(ctx, empty.Image)
	if err == nil {
		t.Error("expected error when pushing with no registry")
	}
	
	if !strings.Contains(err.Error(), "no registry client configured") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// Mock implementations for testing

type mockReadCloser struct {
	reader io.Reader
	closed bool
}

func newMockReadCloser(content string) *mockReadCloser {
	return &mockReadCloser{
		reader: strings.NewReader(content),
		closed: false,
	}
}

func (mrc *mockReadCloser) Read(p []byte) (n int, err error) {
	if mrc.closed {
		return 0, io.ErrClosedPipe
	}
	return mrc.reader.Read(p)
}

func (mrc *mockReadCloser) Close() error {
	mrc.closed = true
	return nil
}

func TestMockReadCloser(t *testing.T) {
	mock := newMockReadCloser("test content")
	
	// Test reading
	buf := make([]byte, 12)
	n, err := mock.Read(buf)
	if err != nil {
		t.Errorf("unexpected error reading: %v", err)
	}
	
	if n != 12 {
		t.Errorf("expected 12 bytes read, got %d", n)
	}
	
	if string(buf) != "test content" {
		t.Errorf("expected 'test content', got '%s'", string(buf))
	}
	
	// Test closing
	err = mock.Close()
	if err != nil {
		t.Errorf("unexpected error closing: %v", err)
	}
	
	// Test reading after close
	_, err = mock.Read(buf)
	if err != io.ErrClosedPipe {
		t.Errorf("expected ErrClosedPipe after close, got %v", err)
	}
}