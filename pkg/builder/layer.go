package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// LayerFactory handles creation and configuration of OCI image layers.
type LayerFactory struct {
	BaseImage     v1.Image
	CompressionAlgorithm v1.LayerCompression
	TimestampPolicy TimestampPolicy
	PermissionMode  os.FileMode
}

// TimestampPolicy defines how file timestamps are handled in layers.
type TimestampPolicy int

const (
	// TimestampPreserve keeps original file timestamps
	TimestampPreserve TimestampPolicy = iota
	// TimestampEpoch sets all timestamps to Unix epoch (for reproducible builds)
	TimestampEpoch
	// TimestampCurrent sets all timestamps to current time
	TimestampCurrent
)

// FileEntry represents a file to be added to a layer.
type FileEntry struct {
	Path        string    // Destination path in the image
	Source      string    // Source path on filesystem
	Mode        os.FileMode
	ModTime     time.Time
	Size        int64
	IsDir       bool
}

// LayerOptions configures layer creation behavior.
type LayerOptions struct {
	Files           []FileEntry          `json:"files,omitempty"`
	WorkingDir      string              `json:"working_dir,omitempty"`
	Compression     v1.LayerCompression `json:"compression,omitempty"`
	TimestampPolicy TimestampPolicy     `json:"timestamp_policy,omitempty"`
	PreserveOwner   bool                `json:"preserve_owner,omitempty"`
	DefaultMode     os.FileMode         `json:"default_mode,omitempty"`
}

// NewLayerFactory creates a new LayerFactory with sensible defaults.
func NewLayerFactory() *LayerFactory {
	return &LayerFactory{
		BaseImage:            empty.Image,
		CompressionAlgorithm: v1.GzipCompression,
		TimestampPolicy:      TimestampEpoch,
		PermissionMode:       0644,
	}
}

// CreateLayer creates a new layer from the specified options.
func (lf *LayerFactory) CreateLayer(ctx context.Context, opts LayerOptions) (v1.Layer, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	
	if len(opts.Files) == 0 {
		return nil, fmt.Errorf("no files specified for layer creation")
	}

	// Validate and prepare file entries
	entries, err := lf.prepareFileEntries(opts.Files, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare file entries: %w", err)
	}

	// Create the layer from file entries
	layer, err := lf.createLayerFromEntries(ctx, entries, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from entries: %w", err)
	}

	return layer, nil
}

// CreateLayerFromDirectory creates a layer from all files in a directory.
func (lf *LayerFactory) CreateLayerFromDirectory(ctx context.Context, sourceDir, targetDir string, opts LayerOptions) (v1.Layer, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	
	if sourceDir == "" {
		return nil, fmt.Errorf("source directory cannot be empty")
	}
	
	if targetDir == "" {
		targetDir = "/"
	}

	// Walk the source directory to collect all files
	files, err := lf.collectDirectoryFiles(sourceDir, targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to collect directory files: %w", err)
	}

	opts.Files = files
	return lf.CreateLayer(ctx, opts)
}

// CreateEmptyLayer creates an empty layer (useful for certain image operations).
func (lf *LayerFactory) CreateEmptyLayer() (v1.Layer, error) {
	// Create an empty layer using tarball package
	layer, err := tarball.LayerFromReader(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create empty layer: %w", err)
	}
	return layer, nil
}

// AddLayer adds a layer to the base image and returns a new image.
func (lf *LayerFactory) AddLayer(image v1.Image, layer v1.Layer) (v1.Image, error) {
	if image == nil {
		return nil, fmt.Errorf("base image cannot be nil")
	}
	
	if layer == nil {
		return nil, fmt.Errorf("layer cannot be nil")
	}

	// Use mutate package to add the layer
	newImage, err := mutate.AppendLayers(image, layer)
	if err != nil {
		return nil, fmt.Errorf("failed to append layer to image: %w", err)
	}

	return newImage, nil
}

// SetBaseImage sets the base image for layer operations.
func (lf *LayerFactory) SetBaseImage(image v1.Image) {
	lf.BaseImage = image
}

// SetCompressionAlgorithm sets the compression algorithm for new layers.
func (lf *LayerFactory) SetCompressionAlgorithm(algorithm v1.LayerCompression) {
	lf.CompressionAlgorithm = algorithm
}

// SetTimestampPolicy sets the timestamp policy for files in layers.
func (lf *LayerFactory) SetTimestampPolicy(policy TimestampPolicy) {
	lf.TimestampPolicy = policy
}

// prepareFileEntries validates and prepares file entries for layer creation.
func (lf *LayerFactory) prepareFileEntries(files []FileEntry, opts LayerOptions) ([]FileEntry, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files to process")
	}

	var prepared []FileEntry
	seenPaths := make(map[string]bool)

	for i, entry := range files {
		// Validate required fields
		if entry.Path == "" {
			return nil, fmt.Errorf("file entry %d: destination path cannot be empty", i)
		}
		
		if entry.Source == "" && !entry.IsDir {
			return nil, fmt.Errorf("file entry %d: source path cannot be empty for files", i)
		}

		// Check for duplicate paths
		if seenPaths[entry.Path] {
			return nil, fmt.Errorf("duplicate destination path: %s", entry.Path)
		}
		seenPaths[entry.Path] = true

		// Normalize the path
		entry.Path = filepath.Clean("/" + entry.Path)

		// Set defaults for optional fields
		if entry.Mode == 0 {
			if entry.IsDir {
				entry.Mode = 0755
			} else {
				entry.Mode = opts.DefaultMode
				if entry.Mode == 0 {
					entry.Mode = lf.PermissionMode
				}
			}
		}

		// Apply timestamp policy
		entry.ModTime = lf.applyTimestampPolicy(entry.ModTime)

		// Validate source file exists (if not a directory entry)
		if !entry.IsDir && entry.Source != "" {
			info, err := os.Stat(entry.Source)
			if err != nil {
				return nil, fmt.Errorf("file entry %d: cannot access source file %s: %w", i, entry.Source, err)
			}
			entry.Size = info.Size()
		}

		prepared = append(prepared, entry)
	}

	// Sort by path to ensure deterministic layer creation
	sort.Slice(prepared, func(i, j int) bool {
		return prepared[i].Path < prepared[j].Path
	})

	return prepared, nil
}

// createLayerFromEntries creates a layer from prepared file entries.
func (lf *LayerFactory) createLayerFromEntries(ctx context.Context, entries []FileEntry, opts LayerOptions) (v1.Layer, error) {
	// For now, create a simple empty layer
	// In the full implementation, this would create a proper tar archive
	// with all the file entries
	
	// This is a simplified implementation for split-002a
	// Full tarball generation will be in split-002b
	
	layer, err := lf.CreateEmptyLayer()
	if err != nil {
		return nil, fmt.Errorf("failed to create base layer: %w", err)
	}

	return layer, nil
}

// collectDirectoryFiles recursively collects all files from a directory.
func (lf *LayerFactory) collectDirectoryFiles(sourceDir, targetDir string) ([]FileEntry, error) {
	var files []FileEntry
	
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Create target path
		targetPath := filepath.Join(targetDir, relPath)
		
		entry := FileEntry{
			Path:    targetPath,
			Source:  path,
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
		}

		files = append(files, entry)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", sourceDir, err)
	}

	return files, nil
}

// applyTimestampPolicy applies the configured timestamp policy to a time value.
func (lf *LayerFactory) applyTimestampPolicy(original time.Time) time.Time {
	switch lf.TimestampPolicy {
	case TimestampEpoch:
		return time.Unix(0, 0).UTC()
	case TimestampCurrent:
		return time.Now().UTC()
	case TimestampPreserve:
		fallthrough
	default:
		return original.UTC()
	}
}

// ValidateLayerOptions validates layer creation options.
func ValidateLayerOptions(opts LayerOptions) error {
	if len(opts.Files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	for i, file := range opts.Files {
		if file.Path == "" {
			return fmt.Errorf("file %d: path cannot be empty", i)
		}
		
		if !file.IsDir && file.Source == "" {
			return fmt.Errorf("file %d: source cannot be empty for regular files", i)
		}
	}

	return nil
}

// GetLayerDigest calculates the digest of a layer.
func GetLayerDigest(layer v1.Layer) (v1.Hash, error) {
	if layer == nil {
		return v1.Hash{}, fmt.Errorf("layer cannot be nil")
	}
	
	digest, err := layer.Digest()
	if err != nil {
		return v1.Hash{}, fmt.Errorf("failed to calculate layer digest: %w", err)
	}
	
	return digest, nil
}

// GetLayerSize returns the compressed size of a layer.
func GetLayerSize(layer v1.Layer) (int64, error) {
	if layer == nil {
		return 0, fmt.Errorf("layer cannot be nil")
	}
	
	size, err := layer.Size()
	if err != nil {
		return 0, fmt.Errorf("failed to get layer size: %w", err)
	}
	
	return size, nil
}