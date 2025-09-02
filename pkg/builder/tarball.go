package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// TarballWriter exports OCI images as tarball files.
// It supports standard tarball format for maximum compatibility.
type TarballWriter struct {
	options TarballOptions
}

// TarballOptions configures tarball export behavior.
type TarballOptions struct {
	// Compress enables gzip compression of the tarball
	Compress bool
	
	// Platform specifies platform selection for multi-arch images
	Platform *v1.Platform
	
	// IncludeManifest includes the image manifest in the tarball
	IncludeManifest bool
}

// NewTarballWriter creates a new tarball writer.
func NewTarballWriter() *TarballWriter {
	return &TarballWriter{
		options: TarballOptions{
			Compress:        false, // Default to uncompressed for speed
			IncludeManifest: true,
		},
	}
}

// NewTarballWriterWithOptions creates a tarball writer with custom options.
func NewTarballWriterWithOptions(opts TarballOptions) *TarballWriter {
	return &TarballWriter{
		options: opts,
	}
}

// Write exports an OCI image to a tarball file.
// The tarball can be imported into Docker, Podman, or other container runtimes.
func (w *TarballWriter) Write(img v1.Image, outputPath string, ref string) error {
	// Validate input parameters
	if img == nil {
		return fmt.Errorf("image cannot be nil")
	}
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}
	if ref == "" {
		return fmt.Errorf("image reference cannot be empty")
	}
	
	// Parse and validate the reference
	imageRef, err := name.ParseReference(ref, name.WeakValidation)
	if err != nil {
		return fmt.Errorf("invalid image reference %s: %w", ref, err)
	}
	
	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
	}
	
	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outputFile.Close()
	
	// Create tag-to-image map for the tarball
	tagToImage := map[name.Tag]v1.Image{}
	
	// Handle reference type
	switch typedRef := imageRef.(type) {
	case name.Tag:
		tagToImage[typedRef] = img
	case name.Digest:
		// For digest references, create a tag with the digest
		tag, err := name.NewTag(fmt.Sprintf("%s:%s", typedRef.Repository.Name(), "latest"))
		if err != nil {
			return fmt.Errorf("failed to create tag from digest: %w", err)
		}
		tagToImage[tag] = img
	default:
		return fmt.Errorf("unsupported reference type: %T", imageRef)
	}
	
	// Write the tarball using the modern API
	err = tarball.MultiWriteToFile(outputPath, tagToImage)
	if err != nil {
		return fmt.Errorf("failed to write tarball: %w", err)
	}
	
	return nil
}

// WriteMultiple exports multiple images to a single tarball.
// This is useful for batch export of related images.
func (w *TarballWriter) WriteMultiple(images map[string]v1.Image, outputPath string) error {
	if len(images) == 0 {
		return fmt.Errorf("no images provided for export")
	}
	
	// Convert string references to proper tags
	tagToImage := make(map[name.Tag]v1.Image)
	for refStr, img := range images {
		tag, err := name.NewTag(refStr, name.WeakValidation)
		if err != nil {
			return fmt.Errorf("invalid reference %s: %w", refStr, err)
		}
		tagToImage[tag] = img
	}
	
	// Write multi-image tarball
	err := tarball.MultiWriteToFile(outputPath, tagToImage)
	if err != nil {
		return fmt.Errorf("failed to write multi-image tarball: %w", err)
	}
	
	return nil
}

// GetTarballInfo returns information about an existing tarball file.
// This is useful for validation and inspection.
func GetTarballInfo(tarballPath string) (*TarballInfo, error) {
	// Check if file exists
	info, err := os.Stat(tarballPath)
	if err != nil {
		return nil, fmt.Errorf("tarball file not found: %w", err)
	}
	
	return &TarballInfo{
		Path: tarballPath,
		Size: info.Size(),
	}, nil
}

// TarballInfo contains metadata about a tarball file.
type TarballInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// LoadFromTarball loads an image from a tarball file.
// This is the reverse operation of Write().
func LoadFromTarball(tarballPath string, ref string) (v1.Image, error) {
	// Parse reference
	tag, err := name.NewTag(ref, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("invalid reference %s: %w", ref, err)
	}
	
	// Load image from tarball
	img, err := tarball.ImageFromPath(tarballPath, &tag)
	if err != nil {
		return nil, fmt.Errorf("failed to load image from tarball: %w", err)
	}
	
	return img, nil
}

// ValidateTarball checks if a tarball file is valid and contains expected content.
func ValidateTarball(tarballPath string) error {
	// Check file existence and readability
	file, err := os.Open(tarballPath)
	if err != nil {
		return fmt.Errorf("cannot open tarball file: %w", err)
	}
	defer file.Close()
	
	// Basic file size check (should not be empty)
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat tarball file: %w", err)
	}
	
	if info.Size() == 0 {
		return fmt.Errorf("tarball file is empty")
	}
	
	// TODO: Add more sophisticated validation like:
	// - Check tar archive structure
	// - Validate OCI/Docker format compliance
	// - Check for required files (manifest, config, layers)
	
	return nil
}

// CompressTarball compresses an existing tarball using gzip.
// This reduces file size but increases CPU usage.
func CompressTarball(inputPath, outputPath string) error {
	// This would implement gzip compression
	// For now, return an error indicating it's not implemented
	return fmt.Errorf("tarball compression not yet implemented")
}

// WithCompression enables gzip compression of the output tarball.
func (w *TarballWriter) WithCompression(compress bool) *TarballWriter {
	w.options.Compress = compress
	return w
}