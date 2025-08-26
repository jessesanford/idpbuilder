package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// layerMetadata extends Layer with cache-specific metadata
type layerMetadata struct {
	*Layer
	RefCount    int       `json:"ref_count"`
	LastAccess  time.Time `json:"last_access"`
	AccessCount int       `json:"access_count"`
	StoredAt    time.Time `json:"stored_at"`
}

// layerDB manages layer storage and metadata
type layerDB struct {
	mu       sync.RWMutex
	layers   map[string]*layerMetadata
	basePath string
	dataPath string
	metaPath string
}

// newLayerDB creates a new layer database
func newLayerDB(basePath string) (*layerDB, error) {
	dataPath := filepath.Join(basePath, "layers")
	metaPath := filepath.Join(basePath, "metadata")

	// Create directories
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	if err := os.MkdirAll(metaPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create metadata directory: %w", err)
	}

	db := &layerDB{
		layers:   make(map[string]*layerMetadata),
		basePath: basePath,
		dataPath: dataPath,
		metaPath: metaPath,
	}

	// Load existing metadata
	if err := db.loadMetadata(); err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	return db, nil
}

// HasLayer checks if layer exists in database
func (db *layerDB) HasLayer(digest string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	_, exists := db.layers[digest]
	return exists
}

// GetLayer retrieves layer from database
func (db *layerDB) GetLayer(digest string) (*Layer, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	meta, exists := db.layers[digest]
	if !exists {
		return nil, fmt.Errorf("layer not found: %s", digest)
	}

	// Load layer data if not in memory
	if meta.Layer.Data == nil {
		data, err := db.loadLayerData(digest)
		if err != nil {
			return nil, fmt.Errorf("failed to load layer data: %w", err)
		}
		meta.Layer.Data = data
	}

	return meta.Layer, nil
}

// StoreLayer stores layer in database
func (db *layerDB) StoreLayer(layer *Layer) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	digest := layer.Digest.String()

	// Check if already exists
	if existing, exists := db.layers[digest]; exists {
		existing.RefCount++
		existing.LastAccess = time.Now()
		return db.saveMetadata()
	}

	// Store layer data to disk
	if err := db.saveLayerData(digest, layer.Data); err != nil {
		return fmt.Errorf("failed to save layer data: %w", err)
	}

	// Create metadata entry
	meta := &layerMetadata{
		Layer:       layer,
		RefCount:    1,
		LastAccess:  time.Now(),
		AccessCount: 1,
		StoredAt:    time.Now(),
	}

	// Don't keep large data in memory
	meta.Layer.Data = nil

	db.layers[digest] = meta
	return db.saveMetadata()
}

// UpdateAccess updates layer access statistics
func (db *layerDB) UpdateAccess(digest string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if meta, exists := db.layers[digest]; exists {
		meta.LastAccess = time.Now()
		meta.AccessCount++
	}
}

// RemoveLayer removes layer from database
func (db *layerDB) RemoveLayer(digest string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	meta, exists := db.layers[digest]
	if !exists {
		return nil // Already removed
	}

	// Decrement reference count
	meta.RefCount--
	if meta.RefCount > 0 {
		return db.saveMetadata()
	}

	// Remove from memory
	delete(db.layers, digest)

	// Remove from disk
	dataFile := filepath.Join(db.dataPath, digest)
	if err := os.Remove(dataFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove layer data: %w", err)
	}

	return db.saveMetadata()
}

// GetAllLayers returns all layer metadata
func (db *layerDB) GetAllLayers() []*layerMetadata {
	db.mu.RLock()
	defer db.mu.RUnlock()

	layers := make([]*layerMetadata, 0, len(db.layers))
	for _, meta := range db.layers {
		layers = append(layers, meta)
	}
	return layers
}

// PruneLayers removes layers older than specified time
func (db *layerDB) PruneLayers(before time.Time) (removedSize int64, removedCount int, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var toRemove []string
	for digest, meta := range db.layers {
		if meta.StoredAt.Before(before) {
			toRemove = append(toRemove, digest)
			removedSize += meta.Layer.Size
		}
	}

	for _, digest := range toRemove {
		delete(db.layers, digest)
		dataFile := filepath.Join(db.dataPath, digest)
		os.Remove(dataFile) // Ignore errors for cleanup
	}

	removedCount = len(toRemove)
	if removedCount > 0 {
		err = db.saveMetadata()
	}

	return removedSize, removedCount, err
}

// Close closes the database and saves metadata
func (db *layerDB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	return db.saveMetadata()
}

// loadLayerData loads layer data from disk
func (db *layerDB) loadLayerData(digest string) ([]byte, error) {
	dataFile := filepath.Join(db.dataPath, digest)
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read layer data: %w", err)
	}
	return data, nil
}

// saveLayerData saves layer data to disk
func (db *layerDB) saveLayerData(digest string, data []byte) error {
	dataFile := filepath.Join(db.dataPath, digest)
	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write layer data: %w", err)
	}
	return nil
}

// loadMetadata loads metadata from disk
func (db *layerDB) loadMetadata() error {
	metaFile := filepath.Join(db.metaPath, "layers.json")
	
	data, err := os.ReadFile(metaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No metadata yet
		}
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	var layers map[string]*layerMetadata
	if err := json.Unmarshal(data, &layers); err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	db.layers = layers
	return nil
}

// saveMetadata saves metadata to disk
func (db *layerDB) saveMetadata() error {
	metaFile := filepath.Join(db.metaPath, "layers.json")
	
	data, err := json.MarshalIndent(db.layers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metaFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}