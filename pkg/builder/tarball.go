package builder

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// TarballManager handles tarball operations for container images
type TarballManager struct {
	options *BuildOptions
	logger  Logger
}

// TarballOptions contains options for tarball creation
type TarballOptions struct {
	Context       string
	IncludeFiles  []string
	ExcludeFiles  []string
	Compression   bool
	PreserveOwner bool
}

// NewTarballManager creates a new TarballManager instance
func NewTarballManager(opts *BuildOptions) *TarballManager {
	return &TarballManager{
		options: opts,
		logger:  opts.Logger,
	}
}

// CreateTarballFromContext creates a tarball from the build context
func (tm *TarballManager) CreateTarballFromContext(ctx context.Context, opts *TarballOptions) (io.Reader, error) {
	if !tm.isFeatureEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	if opts == nil {
		opts = &TarballOptions{
			Context:     tm.options.Context,
			Compression: true,
		}
	}

	tm.logInfo("Creating tarball from context: %s", opts.Context)

	// Create temporary file for tarball
	tmpFile, err := os.CreateTemp("", "context-*.tar.gz")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	// Set up writers
	var writer io.Writer = tmpFile
	if opts.Compression {
		gzw := gzip.NewWriter(tmpFile)
		defer gzw.Close()
		writer = gzw
	}

	tw := tar.NewWriter(writer)
	defer tw.Close()

	// Walk the context directory
	err = filepath.Walk(opts.Context, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(opts.Context, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Check exclusions
		if tm.shouldExclude(relPath, opts.ExcludeFiles) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check inclusions (if specified)
		if len(opts.IncludeFiles) > 0 && !tm.shouldInclude(relPath, opts.IncludeFiles) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("failed to create header for %s: %w", path, err)
		}

		// Set proper path in tarball
		header.Name = filepath.ToSlash(relPath)

		// Preserve ownership if requested
		if !opts.PreserveOwner {
			header.Uid = 0
			header.Gid = 0
			header.Uname = ""
			header.Gname = ""
		}

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header for %s: %w", path, err)
		}

		// Write file content if it's a regular file
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create tarball: %w", err)
	}

	// Reopen file for reading
	tmpFile.Close()
	reader, err := os.Open(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to reopen tarball: %w", err)
	}

	tm.logInfo("Successfully created tarball")
	return reader, nil
}

// CreateImageFromTarball creates a container image from a tarball
func (tm *TarballManager) CreateImageFromTarball(ctx context.Context, tarballPath string, config *v1.Config) (v1.Image, error) {
	if !tm.isFeatureEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	tm.logInfo("Creating image from tarball: %s", tarballPath)

	// Create layer from tarball
	layer, err := tarball.LayerFromFile(tarballPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from tarball: %w", err)
	}

	// Create base image
	img := mutate.MediaType(mutate.ConfigFile(mutate.Config(empty.Image, *config), &v1.ConfigFile{
		Architecture: config.Architecture,
		OS:           config.OS,
		Config:       *config,
		Created:      v1.Time{Time: time.Now()},
	}), types.DockerManifestSchema2)

	// Add layer to image
	img, err = mutate.AppendLayers(img, layer)
	if err != nil {
		return nil, fmt.Errorf("failed to add layer to image: %w", err)
	}

	tm.logInfo("Successfully created image from tarball")
	return img, nil
}

// SaveImageToTarball saves a container image to a tarball file
func (tm *TarballManager) SaveImageToTarball(ctx context.Context, img v1.Image, path string, tag string) error {
	if !tm.isFeatureEnabled() {
		return fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	tm.logInfo("Saving image to tarball: %s", path)

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create or overwrite the tarball file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create tarball file: %w", err)
	}
	defer file.Close()

	// Create tag map
	tagMap := map[string]v1.Image{}
	if tag != "" {
		tagMap[tag] = img
	} else {
		tagMap["latest"] = img
	}

	// Write image to tarball
	err = tarball.MultiWrite(tagMap, file)
	if err != nil {
		return fmt.Errorf("failed to write image to tarball: %w", err)
	}

	tm.logInfo("Successfully saved image to tarball")
	return nil
}

// LoadImageFromTarball loads a container image from a tarball file
func (tm *TarballManager) LoadImageFromTarball(ctx context.Context, path string, tag string) (v1.Image, error) {
	if !tm.isFeatureEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	tm.logInfo("Loading image from tarball: %s", path)

	// Open tarball file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open tarball file: %w", err)
	}
	defer file.Close()

	// Load image from tarball
	img, err := tarball.ImageFromTar(file)
	if err != nil {
		return nil, fmt.Errorf("failed to load image from tarball: %w", err)
	}

	tm.logInfo("Successfully loaded image from tarball")
	return img, nil
}

// ExtractTarball extracts a tarball to the specified directory
func (tm *TarballManager) ExtractTarball(ctx context.Context, tarballPath string, destDir string) error {
	if !tm.isFeatureEnabled() {
		return fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	tm.logInfo("Extracting tarball %s to %s", tarballPath, destDir)

	// Open tarball file
	file, err := os.Open(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to open tarball: %w", err)
	}
	defer file.Close()

	// Determine if it's compressed
	var reader io.Reader = file
	if strings.HasSuffix(tarballPath, ".gz") || strings.HasSuffix(tarballPath, ".tgz") {
		gzr, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzr.Close()
		reader = gzr
	}

	// Create tar reader
	tr := tar.NewReader(reader)

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Create full path
		path := filepath.Join(destDir, header.Name)

		// Security check - ensure path is within destination
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(destDir)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		// Handle different file types
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}
		case tar.TypeReg:
			// Create parent directory
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", path, err)
			}

			// Create file
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", path, err)
			}

			// Copy content
			_, err = io.Copy(file, tr)
			file.Close()
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", path, err)
			}

			// Set modification time
			if err := os.Chtimes(path, time.Now(), header.ModTime); err != nil {
				tm.logWarn("Failed to set modification time for %s: %v", path, err)
			}
		}
	}

	tm.logInfo("Successfully extracted tarball")
	return nil
}

// shouldExclude checks if a path should be excluded based on patterns
func (tm *TarballManager) shouldExclude(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		matched, err := filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}
		
		// Check if path starts with pattern (for directory exclusions)
		if strings.HasPrefix(path, pattern) {
			return true
		}
	}
	return false
}

// shouldInclude checks if a path should be included based on patterns
func (tm *TarballManager) shouldInclude(path string, includePatterns []string) bool {
	for _, pattern := range includePatterns {
		matched, err := filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}
		
		// Check if path starts with pattern
		if strings.HasPrefix(path, pattern) {
			return true
		}
	}
	return false
}

// isFeatureEnabled checks if the CLI tools feature flag is enabled
func (tm *TarballManager) isFeatureEnabled() bool {
	return os.Getenv("ENABLE_CLI_TOOLS") == "true"
}

// Logging helper methods
func (tm *TarballManager) logInfo(msg string, args ...interface{}) {
	if tm.logger != nil {
		tm.logger.Info(msg, args...)
	}
}

func (tm *TarballManager) logWarn(msg string, args ...interface{}) {
	if tm.logger != nil {
		tm.logger.Warn(msg, args...)
	}
}