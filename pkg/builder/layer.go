package builder

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// LayerManager handles container image layer operations
type LayerManager struct {
	options *BuildOptions
	logger  Logger
}

// LayerInfo contains metadata about a layer
type LayerInfo struct {
	Digest     string
	Size       int64
	MediaType  types.MediaType
	CreatedBy  string
	Comment    string
	EmptyLayer bool
}

// NewLayerManager creates a new LayerManager instance
func NewLayerManager(opts *BuildOptions) *LayerManager {
	return &LayerManager{
		options: opts,
		logger:  opts.Logger,
	}
}

// CreateLayerFromDirectory creates a new layer from a directory
func (lm *LayerManager) CreateLayerFromDirectory(ctx context.Context, dir string) (v1.Layer, error) {
	if !lm.isFeatureEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	lm.logInfo("Creating layer from directory: %s", dir)

	// Create tarball buffer
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	// Walk directory and add files to tar
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from directory
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("failed to create tar header for %s: %w", path, err)
		}

		// Set the name to the relative path with proper separators
		header.Name = filepath.ToSlash(relPath)

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header for %s: %w", path, err)
		}

		// If it's a regular file, write the content
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return fmt.Errorf("failed to write file content for %s: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create tar from directory: %w", err)
	}

	// Close tar and gzip writers
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := gzw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Create layer from tarball
	layer, err := tarball.LayerFromReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from tarball: %w", err)
	}

	lm.logInfo("Successfully created layer from directory")
	return layer, nil
}

// CreateLayerFromFiles creates a layer from a list of files
func (lm *LayerManager) CreateLayerFromFiles(ctx context.Context, files []string, basePath string) (v1.Layer, error) {
	if !lm.isFeatureEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	lm.logInfo("Creating layer from %d files", len(files))

	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	for _, file := range files {
		// Get file info
		info, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %w", file, err)
		}

		// Calculate relative path
		relPath, err := filepath.Rel(basePath, file)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path for %s: %w", file, err)
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create tar header for %s: %w", file, err)
		}

		header.Name = filepath.ToSlash(relPath)

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return nil, fmt.Errorf("failed to write tar header for %s: %w", file, err)
		}

		// Write file content if regular file
		if info.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", file, err)
			}

			if _, err := io.Copy(tw, f); err != nil {
				f.Close()
				return nil, fmt.Errorf("failed to copy file content for %s: %w", file, err)
			}
			f.Close()
		}
	}

	// Close writers
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := gzw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Create layer
	layer, err := tarball.LayerFromReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from files: %w", err)
	}

	lm.logInfo("Successfully created layer from files")
	return layer, nil
}

// CreateEmptyLayer creates an empty layer (useful for certain instructions)
func (lm *LayerManager) CreateEmptyLayer() (v1.Layer, error) {
	if !lm.isFeatureEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	// Create empty tar
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer for empty layer: %w", err)
	}
	if err := gzw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer for empty layer: %w", err)
	}

	layer, err := tarball.LayerFromReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to create empty layer: %w", err)
	}

	lm.logInfo("Created empty layer")
	return layer, nil
}

// GetLayerInfo extracts information about a layer
func (lm *LayerManager) GetLayerInfo(layer v1.Layer) (*LayerInfo, error) {
	digest, err := layer.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to get layer digest: %w", err)
	}

	size, err := layer.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get layer size: %w", err)
	}

	mediaType, err := layer.MediaType()
	if err != nil {
		return nil, fmt.Errorf("failed to get layer media type: %w", err)
	}

	return &LayerInfo{
		Digest:     digest.String(),
		Size:       size,
		MediaType:  mediaType,
		EmptyLayer: size == 0,
	}, nil
}

// ExtractLayer extracts a layer to the specified directory
func (lm *LayerManager) ExtractLayer(ctx context.Context, layer v1.Layer, destDir string) error {
	if !lm.isFeatureEnabled() {
		return fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	lm.logInfo("Extracting layer to directory: %s", destDir)

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Get layer reader
	rc, err := layer.Uncompressed()
	if err != nil {
		return fmt.Errorf("failed to get layer reader: %w", err)
	}
	defer rc.Close()

	// Create tar reader
	tr := tar.NewReader(rc)

	// Extract files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Create the full path
		path := filepath.Join(destDir, header.Name)

		// Ensure we don't extract outside the destination directory
		if !strings.HasPrefix(path, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		// Create directories as needed
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(path, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", path, err)
		}

		// Create and write file
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, header.FileInfo().Mode())
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}

		if _, err := io.Copy(file, tr); err != nil {
			file.Close()
			return fmt.Errorf("failed to write file content to %s: %w", path, err)
		}
		file.Close()

		// Set file modification time
		if err := os.Chtimes(path, time.Now(), header.ModTime); err != nil {
			lm.logWarn("Failed to set modification time for %s: %v", path, err)
		}
	}

	lm.logInfo("Successfully extracted layer")
	return nil
}

// isFeatureEnabled checks if the CLI tools feature flag is enabled
func (lm *LayerManager) isFeatureEnabled() bool {
	return os.Getenv("ENABLE_CLI_TOOLS") == "true"
}

// Logging helper methods
func (lm *LayerManager) logInfo(msg string, args ...interface{}) {
	if lm.logger != nil {
		lm.logger.Info(msg, args...)
	}
}

func (lm *LayerManager) logWarn(msg string, args ...interface{}) {
	if lm.logger != nil {
		lm.logger.Warn(msg, args...)
	}
}