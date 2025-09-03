package builder

import (
	"io"
	"strings"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

func TestLayerTypeString(t *testing.T) {
	tests := []struct {
		layerType LayerType
		expected  string
	}{
		{LayerTypeFile, "file"},
		{LayerTypeTar, "tar"},
		{LayerTypeEmpty, "empty"},
		{LayerTypeMetadata, "metadata"},
		{LayerType(999), "unknown"},
	}

	for _, test := range tests {
		result := test.layerType.String()
		if result != test.expected {
			t.Errorf("LayerType(%d).String() = %s, expected %s", 
				int(test.layerType), result, test.expected)
		}
	}
}

func TestNewBaseLayer(t *testing.T) {
	reader := newMockReadCloser("test content")
	baseLayer := NewBaseLayer(LayerTypeFile, "test layer", reader, 100)

	if baseLayer == nil {
		t.Fatal("base layer should not be nil")
	}

	if baseLayer.GetType() != LayerTypeFile {
		t.Errorf("expected layer type %v, got %v", LayerTypeFile, baseLayer.GetType())
	}

	if baseLayer.GetDescription() != "test layer" {
		t.Errorf("expected description 'test layer', got '%s'", baseLayer.GetDescription())
	}

	size, err := baseLayer.Size()
	if err != nil {
		t.Errorf("unexpected error getting size: %v", err)
	}
	if size != 100 {
		t.Errorf("expected size 100, got %d", size)
	}

	mediaType, err := baseLayer.MediaType()
	if err != nil {
		t.Errorf("unexpected error getting media type: %v", err)
	}
	if mediaType != types.DockerLayer {
		t.Errorf("expected media type %v, got %v", types.DockerLayer, mediaType)
	}
}

func TestEmptyLayer(t *testing.T) {
	layer := NewEmptyLayer()

	if layer == nil {
		t.Fatal("empty layer should not be nil")
	}

	if layer.GetType() != LayerTypeEmpty {
		t.Errorf("expected layer type %v, got %v", LayerTypeEmpty, layer.GetType())
	}

	// Test compressed reader
	compressed, err := layer.Compressed()
	if err != nil {
		t.Errorf("unexpected error getting compressed reader: %v", err)
	}
	if compressed == nil {
		t.Fatal("compressed reader should not be nil")
	}

	// Read from compressed reader - should be empty
	buf := make([]byte, 10)
	n, err := compressed.Read(buf)
	if n != 0 {
		t.Errorf("expected 0 bytes read from empty layer, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("expected EOF from empty layer, got %v", err)
	}

	compressed.Close()

	// Test uncompressed reader
	uncompressed, err := layer.Uncompressed()
	if err != nil {
		t.Errorf("unexpected error getting uncompressed reader: %v", err)
	}
	if uncompressed == nil {
		t.Fatal("uncompressed reader should not be nil")
	}

	n, err = uncompressed.Read(buf)
	if n != 0 {
		t.Errorf("expected 0 bytes read from empty layer, got %d", n)
	}
	if err != io.EOF {
		t.Errorf("expected EOF from empty layer, got %v", err)
	}

	uncompressed.Close()
}

func TestFileLayer(t *testing.T) {
	files := map[string][]byte{
		"file1.txt": []byte("content1"),
		"file2.txt": []byte("content2"),
	}

	layer, err := NewFileLayer(files, "test file layer")
	if err != nil {
		t.Fatalf("unexpected error creating file layer: %v", err)
	}

	if layer == nil {
		t.Fatal("file layer should not be nil")
	}

	if layer.GetType() != LayerTypeFile {
		t.Errorf("expected layer type %v, got %v", LayerTypeFile, layer.GetType())
	}

	if layer.GetDescription() != "test file layer" {
		t.Errorf("expected description 'test file layer', got '%s'", layer.GetDescription())
	}

	// Test size calculation
	size, err := layer.Size()
	if err != nil {
		t.Errorf("unexpected error getting size: %v", err)
	}
	expectedSize := int64(len("content1") + len("content2"))
	if size != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, size)
	}

	// Test getting files
	layerFiles := layer.GetFiles()
	if len(layerFiles) != 2 {
		t.Errorf("expected 2 files, got %d", len(layerFiles))
	}

	if string(layerFiles["file1.txt"]) != "content1" {
		t.Errorf("expected file1.txt content 'content1', got '%s'", string(layerFiles["file1.txt"]))
	}

	if string(layerFiles["file2.txt"]) != "content2" {
		t.Errorf("expected file2.txt content 'content2', got '%s'", string(layerFiles["file2.txt"]))
	}
}

func TestFileLayerAddFile(t *testing.T) {
	layer, err := NewFileLayer(nil, "test layer")
	if err != nil {
		t.Fatalf("unexpected error creating file layer: %v", err)
	}

	// Add a file
	layer.AddFile("test.txt", []byte("test content"))

	files := layer.GetFiles()
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}

	if string(files["test.txt"]) != "test content" {
		t.Errorf("expected file content 'test content', got '%s'", string(files["test.txt"]))
	}

	// Test size update
	size, err := layer.Size()
	if err != nil {
		t.Errorf("unexpected error getting size: %v", err)
	}
	expectedSize := int64(len("test content"))
	if size != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, size)
	}
}

func TestTarLayer(t *testing.T) {
	tarContent := "mock tar content"
	reader := newMockReadCloser(tarContent)
	size := int64(len(tarContent))

	layer := NewTarLayer(reader, size, "test tar layer")

	if layer == nil {
		t.Fatal("tar layer should not be nil")
	}

	if layer.GetType() != LayerTypeTar {
		t.Errorf("expected layer type %v, got %v", LayerTypeTar, layer.GetType())
	}

	if layer.GetDescription() != "test tar layer" {
		t.Errorf("expected description 'test tar layer', got '%s'", layer.GetDescription())
	}

	layerSize, err := layer.Size()
	if err != nil {
		t.Errorf("unexpected error getting size: %v", err)
	}
	if layerSize != size {
		t.Errorf("expected size %d, got %d", size, layerSize)
	}
}

func TestStreamLayer(t *testing.T) {
	content := "stream content"
	streamFunc := func() (io.ReadCloser, error) {
		return newMockReadCloser(content), nil
	}

	layer := NewStreamLayer(streamFunc, types.DockerLayer, "test stream layer")

	if layer == nil {
		t.Fatal("stream layer should not be nil")
	}

	if layer.GetDescription() != "test stream layer" {
		t.Errorf("expected description 'test stream layer', got '%s'", layer.GetDescription())
	}

	// Test compressed reader
	compressed, err := layer.Compressed()
	if err != nil {
		t.Errorf("unexpected error getting compressed reader: %v", err)
	}
	if compressed == nil {
		t.Fatal("compressed reader should not be nil")
	}

	buf := make([]byte, len(content))
	n, err := compressed.Read(buf)
	if err != nil && err != io.EOF {
		t.Errorf("unexpected error reading: %v", err)
	}
	if n != len(content) {
		t.Errorf("expected %d bytes read, got %d", len(content), n)
	}
	if string(buf) != content {
		t.Errorf("expected content '%s', got '%s'", content, string(buf))
	}

	compressed.Close()
}

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache()

	if cache == nil {
		t.Fatal("cache should not be nil")
	}

	// Test initial state
	if cache.Size() != 0 {
		t.Errorf("expected cache size 0, got %d", cache.Size())
	}

	// Create a test layer and hash
	layer := NewEmptyLayer()
	hash := v1.Hash{Algorithm: "sha256", Hex: "test123"}

	// Test HasLayer - should not exist
	if cache.HasLayer(hash) {
		t.Error("layer should not exist in empty cache")
	}

	// Test GetLayer - should return error
	_, err := cache.GetLayer(hash)
	if err == nil {
		t.Error("expected error getting non-existent layer")
	}

	// Test PutLayer
	err = cache.PutLayer(hash, layer)
	if err != nil {
		t.Errorf("unexpected error putting layer: %v", err)
	}

	// Test cache size increased
	if cache.Size() != 1 {
		t.Errorf("expected cache size 1, got %d", cache.Size())
	}

	// Test HasLayer - should exist now
	if !cache.HasLayer(hash) {
		t.Error("layer should exist in cache")
	}

	// Test GetLayer - should return the layer
	cachedLayer, err := cache.GetLayer(hash)
	if err != nil {
		t.Errorf("unexpected error getting cached layer: %v", err)
	}
	if cachedLayer != layer {
		t.Error("cached layer should be the same instance")
	}

	// Test Clear
	err = cache.Clear()
	if err != nil {
		t.Errorf("unexpected error clearing cache: %v", err)
	}

	if cache.Size() != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", cache.Size())
	}

	if cache.HasLayer(hash) {
		t.Error("layer should not exist after clear")
	}
}

func TestLayerBuilder(t *testing.T) {
	// Test file layer builder
	builder := NewLayerBuilder(LayerTypeFile)

	if builder == nil {
		t.Fatal("layer builder should not be nil")
	}

	layer, err := builder.
		WithDescription("test layer").
		WithFile("file1.txt", []byte("content1")).
		WithFile("file2.txt", []byte("content2")).
		Build()

	if err != nil {
		t.Fatalf("unexpected error building layer: %v", err)
	}

	if layer == nil {
		t.Fatal("built layer should not be nil")
	}

	if layer.GetType() != LayerTypeFile {
		t.Errorf("expected layer type %v, got %v", LayerTypeFile, layer.GetType())
	}

	if layer.GetDescription() != "test layer" {
		t.Errorf("expected description 'test layer', got '%s'", layer.GetDescription())
	}

	// Test file layer specific functionality
	if fileLayer, ok := layer.(*FileLayer); ok {
		files := fileLayer.GetFiles()
		if len(files) != 2 {
			t.Errorf("expected 2 files, got %d", len(files))
		}
	} else {
		t.Error("layer should be a FileLayer")
	}
}

func TestLayerBuilderEmpty(t *testing.T) {
	builder := NewLayerBuilder(LayerTypeEmpty)
	layer, err := builder.Build()

	if err != nil {
		t.Errorf("unexpected error building empty layer: %v", err)
	}

	if layer.GetType() != LayerTypeEmpty {
		t.Errorf("expected layer type %v, got %v", LayerTypeEmpty, layer.GetType())
	}
}

func TestLayerBuilderTar(t *testing.T) {
	builder := NewLayerBuilder(LayerTypeTar)
	_, err := builder.Build()

	if err == nil {
		t.Error("expected error building tar layer without tar reader")
	}

	expectedMsg := "tar layers require a tar reader"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestLayerBuilderUnsupportedType(t *testing.T) {
	builder := NewLayerBuilder(LayerType(999))
	_, err := builder.Build()

	if err == nil {
		t.Error("expected error building layer with unsupported type")
	}

	expectedMsg := "unsupported layer type"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestGetLayerInfo(t *testing.T) {
	layer := NewEmptyLayer()

	info, err := GetLayerInfo(layer)
	if err != nil {
		t.Errorf("unexpected error getting layer info: %v", err)
	}

	if info == nil {
		t.Fatal("layer info should not be nil")
	}

	if info.LayerType != LayerTypeEmpty {
		t.Errorf("expected layer type %v, got %v", LayerTypeEmpty, info.LayerType)
	}

	if info.Description != "empty layer" {
		t.Errorf("expected description 'empty layer', got '%s'", info.Description)
	}
}

func TestLayerFromTar(t *testing.T) {
	content := "mock tar content"
	reader := newMockReadCloser(content)
	size := int64(len(content))

	// Test compressed tar
	layer, err := LayerFromTar(reader, size, true)
	if err != nil {
		t.Errorf("unexpected error creating layer from tar: %v", err)
	}
	if layer == nil {
		t.Error("layer should not be nil")
	}

	// Test uncompressed tar
	reader2 := newMockReadCloser(content)
	layer2, err := LayerFromTar(reader2, size, false)
	if err != nil {
		t.Errorf("unexpected error creating layer from uncompressed tar: %v", err)
	}
	if layer2 == nil {
		t.Error("layer should not be nil")
	}
}