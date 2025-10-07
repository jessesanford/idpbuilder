// Package push provides OCI image push operations for idpbuilder
package push

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// LocalImage represents a discovered local OCI image
type LocalImage struct {
	Name   string   // Name/tag of the image
	Path   string   // Local path to the image
	Format string   // Format: "tarball", "oci-layout"
	Image  v1.Image // The actual image object
}

// DiscoveryOptions configures image discovery behavior
type DiscoveryOptions struct {
	BuildPath   string   // Path to search for images
	Extensions  []string // File extensions to consider
	MaxSizeMB   int64    // Maximum image size to consider
	FollowLinks bool     // Whether to follow symlinks
}

// DefaultDiscoveryOptions returns sensible defaults for image discovery
func DefaultDiscoveryOptions(buildPath string) *DiscoveryOptions {
	return &DiscoveryOptions{
		BuildPath:   buildPath,
		Extensions:  []string{".tar", ".tar.gz", ".tgz"},
		MaxSizeMB:   5000, // 5GB max
		FollowLinks: true,
	}
}

// DiscoverLocalImages scans the specified path for OCI images and returns a list of discoverable images.
// It supports both tarball formats (docker save) and OCI layout directories.
func DiscoverLocalImages(buildPath string) ([]*LocalImage, error) {
	opts := DefaultDiscoveryOptions(buildPath)
	return DiscoverLocalImagesWithOptions(opts)
}

// DiscoverLocalImagesWithOptions scans for images with custom options
func DiscoverLocalImagesWithOptions(opts *DiscoveryOptions) ([]*LocalImage, error) {
	var images []*LocalImage

	// Check if build path exists
	if _, err := os.Stat(opts.BuildPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("build path does not exist: %s", opts.BuildPath)
	}

	// Walk the directory tree
	err := filepath.Walk(opts.BuildPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories unless they're OCI layout
		if info.IsDir() {
			if isOCILayout(path) {
				image, err := loadOCILayoutImage(path)
				if err != nil {
					return fmt.Errorf("failed to load OCI layout from %s: %w", path, err)
				}
				if image != nil {
					images = append(images, image)
				}
			}
			return nil
		}

		// Check file size limits
		if opts.MaxSizeMB > 0 && info.Size() > opts.MaxSizeMB*1024*1024 {
			return nil // Skip large files
		}

		// Check for tarball images
		if isSupportedTarball(path, opts.Extensions) {
			image, err := loadTarballImage(path)
			if err != nil {
				// Log but don't fail on individual image errors
				return nil
			}
			if image != nil {
				images = append(images, image)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover images in %s: %w", opts.BuildPath, err)
	}

	return images, nil
}

// isOCILayout checks if a directory contains an OCI image layout
func isOCILayout(path string) bool {
	// Check for oci-layout file and index.json
	layoutFile := filepath.Join(path, "oci-layout")
	indexFile := filepath.Join(path, "index.json")

	layoutExists := fileExists(layoutFile)
	indexExists := fileExists(indexFile)

	return layoutExists && indexExists
}

// loadOCILayoutImage loads an image from an OCI layout directory
func loadOCILayoutImage(path string) (*LocalImage, error) {
	// Load the OCI layout
	layoutPath, err := layout.FromPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load OCI layout: %w", err)
	}

	// Get the image index
	index, err := layoutPath.ImageIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to get image index: %w", err)
	}

	// For simplicity, take the first image from the index
	manifest, err := index.IndexManifest()
	if err != nil {
		return nil, fmt.Errorf("failed to get index manifest: %w", err)
	}

	if len(manifest.Manifests) == 0 {
		return nil, fmt.Errorf("no images found in OCI layout")
	}

	// Get the first image
	image, err := layoutPath.Image(manifest.Manifests[0].Digest)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}

	// Extract name from annotations or use directory name as fallback
	name := filepath.Base(path)
	if manifest.Manifests[0].Annotations != nil {
		if refName, ok := manifest.Manifests[0].Annotations["org.opencontainers.image.ref.name"]; ok {
			name = refName
		}
	}

	return &LocalImage{
		Name:   name,
		Path:   path,
		Format: "oci-layout",
		Image:  image,
	}, nil
}

// isSupportedTarball checks if a file is a supported tarball format
func isSupportedTarball(path string, extensions []string) bool {
	lower := strings.ToLower(path)
	for _, ext := range extensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// loadTarballImage loads an image from a tarball file
func loadTarballImage(path string) (*LocalImage, error) {
	// First, check if it's a valid Docker tarball by examining its structure
	if !isValidDockerTarball(path) {
		return nil, nil // Not a Docker tarball, skip
	}

	// Load the tarball image
	image, err := tarball.ImageFromPath(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load tarball image from %s: %w", path, err)
	}

	// Extract name from file path (remove extension)
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	if strings.HasSuffix(name, ".tar") {
		name = strings.TrimSuffix(name, ".tar")
	}

	return &LocalImage{
		Name:   name,
		Path:   path,
		Format: "tarball",
		Image:  image,
	}, nil
}

// isValidDockerTarball checks if a tarball contains Docker image data
func isValidDockerTarball(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	var reader io.Reader = file

	// Handle gzipped files
	if strings.HasSuffix(strings.ToLower(path), ".gz") || strings.HasSuffix(strings.ToLower(path), ".tgz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return false
		}
		defer gzReader.Close()
		reader = gzReader
	}

	tarReader := tar.NewReader(reader)

	// Look for Docker-specific files
	foundManifest := false
	foundRepositories := false

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false
		}

		switch header.Name {
		case "manifest.json":
			foundManifest = true
		case "repositories":
			foundRepositories = true
		}

		// If we found both required files, it's likely a Docker tarball
		if foundManifest && foundRepositories {
			return true
		}

		// Stop after checking a reasonable number of entries
		// to avoid reading the entire large file
		if foundManifest {
			return true // manifest.json is sufficient for modern Docker exports
		}
	}

	return foundManifest || foundRepositories
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// FilterPushTargets filters images based on push criteria
func FilterPushTargets(images []*LocalImage, criteria *FilterCriteria) []*LocalImage {
	if criteria == nil {
		return images
	}

	var filtered []*LocalImage
	for _, image := range images {
		if shouldIncludeImage(image, criteria) {
			filtered = append(filtered, image)
		}
	}
	return filtered
}

// FilterCriteria defines criteria for filtering images before push
type FilterCriteria struct {
	IncludePatterns []string // Glob patterns for image names to include
	ExcludePatterns []string // Glob patterns for image names to exclude
	MinSizeBytes    int64    // Minimum image size
	MaxSizeBytes    int64    // Maximum image size
}

// shouldIncludeImage checks if an image matches the filter criteria
func shouldIncludeImage(image *LocalImage, criteria *FilterCriteria) bool {
	// Check size constraints
	if criteria.MinSizeBytes > 0 || criteria.MaxSizeBytes > 0 {
		if info, err := os.Stat(image.Path); err == nil {
			size := info.Size()
			if criteria.MinSizeBytes > 0 && size < criteria.MinSizeBytes {
				return false
			}
			if criteria.MaxSizeBytes > 0 && size > criteria.MaxSizeBytes {
				return false
			}
		}
	}

	// Check include patterns
	if len(criteria.IncludePatterns) > 0 {
		matched := false
		for _, pattern := range criteria.IncludePatterns {
			if matched, _ := filepath.Match(pattern, image.Name); matched {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check exclude patterns
	for _, pattern := range criteria.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, image.Name); matched {
			return false
		}
	}

	return true
}
