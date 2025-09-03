package builder

import (
	"compress/gzip"
	"fmt"
	"io"
	"sync"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/stream"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// Layer represents a single layer in a container image.
type Layer interface {
	v1.Layer
	
	// GetType returns the type of layer
	GetType() LayerType
	
	// GetDescription returns a human-readable description of the layer
	GetDescription() string
}

// LayerType represents different types of layers.
type LayerType int

const (
	// LayerTypeFile represents a layer containing files
	LayerTypeFile LayerType = iota
	
	// LayerTypeTar represents a layer from a tar archive
	LayerTypeTar
	
	// LayerTypeEmpty represents an empty layer (for directory structure)
	LayerTypeEmpty
	
	// LayerTypeMetadata represents a metadata-only layer
	LayerTypeMetadata
)

// String returns the string representation of the layer type.
func (lt LayerType) String() string {
	switch lt {
	case LayerTypeFile:
		return "file"
	case LayerTypeTar:
		return "tar"
	case LayerTypeEmpty:
		return "empty"
	case LayerTypeMetadata:
		return "metadata"
	default:
		return "unknown"
	}
}

// BaseLayer provides common functionality for all layer types.
type BaseLayer struct {
	digest      v1.Hash
	diffID      v1.Hash
	mediaType   types.MediaType
	size        int64
	layerType   LayerType
	description string
	compressed  io.ReadCloser
}

// NewBaseLayer creates a new base layer with the given parameters.
func NewBaseLayer(layerType LayerType, description string, compressed io.ReadCloser, size int64) *BaseLayer {
	return &BaseLayer{
		layerType:   layerType,
		description: description,
		compressed:  compressed,
		size:        size,
		mediaType:   types.DockerLayer,
	}
}

// Digest returns the SHA256 hash of the compressed layer.
func (bl *BaseLayer) Digest() (v1.Hash, error) {
	if bl.digest.Hex == "" {
		return v1.Hash{}, fmt.Errorf("digest not calculated")
	}
	return bl.digest, nil
}

// DiffID returns the SHA256 hash of the uncompressed layer.
func (bl *BaseLayer) DiffID() (v1.Hash, error) {
	if bl.diffID.Hex == "" {
		return v1.Hash{}, fmt.Errorf("diffID not calculated")
	}
	return bl.diffID, nil
}

// Size returns the compressed size of the layer.
func (bl *BaseLayer) Size() (int64, error) {
	return bl.size, nil
}

// MediaType returns the media type of the layer.
func (bl *BaseLayer) MediaType() (types.MediaType, error) {
	return bl.mediaType, nil
}

// Compressed returns a reader for the compressed layer content.
func (bl *BaseLayer) Compressed() (io.ReadCloser, error) {
	if bl.compressed == nil {
		return nil, fmt.Errorf("no compressed content available")
	}
	return bl.compressed, nil
}

// Uncompressed returns a reader for the uncompressed layer content.
func (bl *BaseLayer) Uncompressed() (io.ReadCloser, error) {
	compressed, err := bl.Compressed()
	if err != nil {
		return nil, err
	}
	
	return gzip.NewReader(compressed)
}

// GetType returns the type of layer.
func (bl *BaseLayer) GetType() LayerType {
	return bl.layerType
}

// GetDescription returns a human-readable description of the layer.
func (bl *BaseLayer) GetDescription() string {
	return bl.description
}

// EmptyLayer represents an empty layer, typically used for directory creation.
type EmptyLayer struct {
	*BaseLayer
}

// NewEmptyLayer creates a new empty layer.
func NewEmptyLayer() *EmptyLayer {
	// Empty layers have standard digests for empty tar archives
	emptyTarDigest := v1.Hash{
		Algorithm: "sha256", 
		Hex:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", // sha256 of empty string
	}
	
	layer := &EmptyLayer{
		BaseLayer: NewBaseLayer(LayerTypeEmpty, "empty layer", nil, 0),
	}
	
	// Set the calculated digests
	layer.BaseLayer.digest = emptyTarDigest
	layer.BaseLayer.diffID = emptyTarDigest
	
	return layer
}

// Compressed returns an empty reader for empty layers.
func (el *EmptyLayer) Compressed() (io.ReadCloser, error) {
	return io.NopCloser(&emptyReader{}), nil
}

// Uncompressed returns an empty reader for empty layers.
func (el *EmptyLayer) Uncompressed() (io.ReadCloser, error) {
	return io.NopCloser(&emptyReader{}), nil
}

// emptyReader provides an empty reader implementation.
type emptyReader struct{}

func (er *emptyReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

// FileLayer represents a layer created from files in the build context.
type FileLayer struct {
	*BaseLayer
	files map[string][]byte
}

// NewFileLayer creates a new file layer with the given files.
func NewFileLayer(files map[string][]byte, description string) (*FileLayer, error) {
	if files == nil {
		files = make(map[string][]byte)
	}

	fl := &FileLayer{
		BaseLayer: NewBaseLayer(LayerTypeFile, description, nil, 0),
		files:     files,
	}

	// Calculate the total size
	var totalSize int64
	for _, content := range files {
		totalSize += int64(len(content))
	}
	fl.size = totalSize

	return fl, nil
}

// AddFile adds a file to the layer.
func (fl *FileLayer) AddFile(path string, content []byte) {
	if fl.files == nil {
		fl.files = make(map[string][]byte)
	}
	fl.files[path] = content
	
	// Recalculate size
	var totalSize int64
	for _, content := range fl.files {
		totalSize += int64(len(content))
	}
	fl.size = totalSize
}

// GetFiles returns a copy of the files in this layer.
func (fl *FileLayer) GetFiles() map[string][]byte {
	files := make(map[string][]byte)
	for path, content := range fl.files {
		contentCopy := make([]byte, len(content))
		copy(contentCopy, content)
		files[path] = contentCopy
	}
	return files
}

// TarLayer represents a layer created from a tar archive.
type TarLayer struct {
	*BaseLayer
	tarReader io.ReadCloser
}

// NewTarLayer creates a new tar layer from a tar archive reader.
func NewTarLayer(tarReader io.ReadCloser, size int64, description string) *TarLayer {
	return &TarLayer{
		BaseLayer: NewBaseLayer(LayerTypeTar, description, nil, size),
		tarReader: tarReader,
	}
}

// Compressed returns the tar content as compressed data.
func (tl *TarLayer) Compressed() (io.ReadCloser, error) {
	if tl.tarReader == nil {
		return nil, fmt.Errorf("no tar reader available")
	}
	
	// In a real implementation, this would compress the tar data
	return tl.tarReader, nil
}

// StreamLayer represents a layer that can be streamed.
type StreamLayer struct {
	*BaseLayer
	streamFunc func() (io.ReadCloser, error)
}

// NewStreamLayer creates a new streaming layer.
func NewStreamLayer(streamFunc func() (io.ReadCloser, error), mediaType types.MediaType, description string) *StreamLayer {
	return &StreamLayer{
		BaseLayer:  NewBaseLayer(LayerTypeFile, description, nil, -1), // Size unknown for streaming
		streamFunc: streamFunc,
	}
}

// Compressed returns a stream of the compressed layer content.
func (sl *StreamLayer) Compressed() (io.ReadCloser, error) {
	if sl.streamFunc == nil {
		return nil, fmt.Errorf("no stream function available")
	}
	return sl.streamFunc()
}

// MemoryCache provides an in-memory implementation of LayerCache.
type MemoryCache struct {
	mu     sync.RWMutex
	layers map[string]v1.Layer
}

// NewMemoryCache creates a new memory-based layer cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		layers: make(map[string]v1.Layer),
	}
}

// GetLayer retrieves a cached layer by its digest.
func (mc *MemoryCache) GetLayer(digest v1.Hash) (v1.Layer, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	layer, exists := mc.layers[digest.String()]
	if !exists {
		return nil, fmt.Errorf("layer not found in cache: %s", digest.String())
	}
	
	return layer, nil
}

// PutLayer stores a layer in the cache.
func (mc *MemoryCache) PutLayer(digest v1.Hash, layer v1.Layer) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	mc.layers[digest.String()] = layer
	return nil
}

// HasLayer checks if a layer exists in the cache.
func (mc *MemoryCache) HasLayer(digest v1.Hash) bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	_, exists := mc.layers[digest.String()]
	return exists
}

// Clear clears all cached layers.
func (mc *MemoryCache) Clear() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	mc.layers = make(map[string]v1.Layer)
	return nil
}

// Size returns the number of cached layers.
func (mc *MemoryCache) Size() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	return len(mc.layers)
}

// LayerBuilder helps construct layers with validation.
type LayerBuilder struct {
	layerType   LayerType
	description string
	mediaType   types.MediaType
	files       map[string][]byte
	size        int64
}

// NewLayerBuilder creates a new layer builder.
func NewLayerBuilder(layerType LayerType) *LayerBuilder {
	return &LayerBuilder{
		layerType: layerType,
		mediaType: types.DockerLayer,
		files:     make(map[string][]byte),
	}
}

// WithDescription sets the layer description.
func (lb *LayerBuilder) WithDescription(description string) *LayerBuilder {
	lb.description = description
	return lb
}

// WithMediaType sets the layer media type.
func (lb *LayerBuilder) WithMediaType(mediaType types.MediaType) *LayerBuilder {
	lb.mediaType = mediaType
	return lb
}

// WithFile adds a file to the layer.
func (lb *LayerBuilder) WithFile(path string, content []byte) *LayerBuilder {
	if lb.files == nil {
		lb.files = make(map[string][]byte)
	}
	lb.files[path] = content
	lb.size += int64(len(content))
	return lb
}

// Build creates the layer based on the configured parameters.
func (lb *LayerBuilder) Build() (Layer, error) {
	switch lb.layerType {
	case LayerTypeFile:
		return NewFileLayer(lb.files, lb.description)
	case LayerTypeEmpty:
		return NewEmptyLayer(), nil
	case LayerTypeTar:
		return nil, fmt.Errorf("tar layers require a tar reader, use NewTarLayer directly")
	default:
		return nil, fmt.Errorf("unsupported layer type: %v", lb.layerType)
	}
}

// StreamableLayer creates a layer that can be streamed efficiently.
func StreamableLayer(reader func() (io.ReadCloser, error)) v1.Layer {
	// Call the reader function to get the io.ReadCloser
	rc, err := reader()
	if err != nil {
		// In case of error, return an empty layer or handle appropriately
		return NewEmptyLayer()
	}
	return stream.NewLayer(rc)
}

// LayerFromTar creates a layer from a tar archive.
func LayerFromTar(tarReader io.ReadCloser, size int64, compressed bool) (v1.Layer, error) {
	if compressed {
		return stream.NewLayer(tarReader), nil
	}
	
	// For uncompressed tar, we need to wrap it
	// In a real implementation, this would compress the tar
	return stream.NewLayer(tarReader), nil
}

// LayerInfo provides information about a layer without loading its content.
type LayerInfo struct {
	Digest      v1.Hash
	DiffID      v1.Hash
	Size        int64
	MediaType   types.MediaType
	LayerType   LayerType
	Description string
}

// GetLayerInfo extracts information from a layer without loading its content.
func GetLayerInfo(layer v1.Layer) (*LayerInfo, error) {
	digest, err := layer.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to get digest: %w", err)
	}
	
	diffID, err := layer.DiffID()
	if err != nil {
		return nil, fmt.Errorf("failed to get diffID: %w", err)
	}
	
	size, err := layer.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}
	
	mediaType, err := layer.MediaType()
	if err != nil {
		return nil, fmt.Errorf("failed to get media type: %w", err)
	}
	
	info := &LayerInfo{
		Digest:    digest,
		DiffID:    diffID,
		Size:      size,
		MediaType: mediaType,
	}
	
	// If this is our custom layer type, extract additional info
	if customLayer, ok := layer.(Layer); ok {
		info.LayerType = customLayer.GetType()
		info.Description = customLayer.GetDescription()
	}
	
	return info, nil
}