// Package build provides layer management for OCI image builds.
package build

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
	"time"

	// Import Phase 1 types
	api "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// LayerManager handles layer creation and caching.
type LayerManager struct {
	cacheDir    string
	layerStore  map[string]*api.LayerInfo
	layersMutex sync.RWMutex
}

const (
	// DefaultMediaType for OCI layers
	DefaultMediaType = "application/vnd.oci.image.layer.v1.tar+gzip"
)

// NewLayerManager creates a new layer manager.
func NewLayerManager(cacheDir string) *LayerManager {
	manager := &LayerManager{
		cacheDir:   cacheDir,
		layerStore: make(map[string]*LayerInfo),
	}

	// Ensure cache directory exists
	os.MkdirAll(cacheDir, 0755)
	return manager
}

// CreateLayer creates a new layer from content and instruction.
func (l *LayerManager) CreateLayer(content []byte, instruction *Instruction) (*api.LayerInfo, error) {
	if instruction == nil {
		return nil, fmt.Errorf("instruction cannot be nil")
	}

	// Calculate digest
	digest := l.CalculateDigest(content)

	// Check cache first
	if cached, exists := l.GetCachedLayer(digest); exists {
		return cached, nil
	}

	// Create layer info
	layer := &api.LayerInfo{
		Digest:     digest,
		Size:       int64(len(content)),
		MediaType:  DefaultMediaType,
		Created:    time.Now(),
		CreatedBy:  l.formatCreatedBy(instruction),
		Comment:    l.generateComment(instruction),
		EmptyLayer: l.isEmptyInstruction(instruction),
	}

	// Store layer
	if err := l.StoreLayer(layer); err != nil {
		return nil, fmt.Errorf("failed to store layer: %w", err)
	}

	return layer, nil
}

// GetCachedLayer retrieves a cached layer by digest.
func (l *LayerManager) GetCachedLayer(digest string) (*api.LayerInfo, bool) {
	l.layersMutex.RLock()
	defer l.layersMutex.RUnlock()

	layer, exists := l.layerStore[digest]
	return layer, exists
}

// StoreLayer stores a layer in the cache.
func (l *LayerManager) StoreLayer(layer *api.LayerInfo) error {
	if layer == nil || layer.Digest == "" {
		return fmt.Errorf("invalid layer")
	}

	l.layersMutex.Lock()
	defer l.layersMutex.Unlock()

	l.layerStore[layer.Digest] = layer
	return nil
}

// CalculateDigest calculates SHA256 digest of content.
func (l *LayerManager) CalculateDigest(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("sha256:%x", hash)
}

// Helper methods

// formatCreatedBy formats the created_by field.
func (l *LayerManager) formatCreatedBy(instruction *Instruction) string {
	if len(instruction.Args) > 0 {
		return fmt.Sprintf("%s %s", instruction.Command, instruction.Args[0])
	}
	return instruction.Command
}

// generateComment generates a comment for the layer.
func (l *LayerManager) generateComment(instruction *Instruction) string {
	switch instruction.Command {
	case "RUN":
		return "Execute build command"
	case "COPY", "ADD":
		return "Copy files into image"
	case "FROM":
		return "Base image layer"
	default:
		return fmt.Sprintf("%s instruction", instruction.Command)
	}
}

// isEmptyInstruction determines if instruction creates empty layer.
func (l *LayerManager) isEmptyInstruction(instruction *Instruction) bool {
	switch instruction.Command {
	case "ENV", "LABEL", "EXPOSE", "USER", "WORKDIR", "ENTRYPOINT", "CMD", "VOLUME", "ARG":
		return true
	default:
		return false
	}
}
