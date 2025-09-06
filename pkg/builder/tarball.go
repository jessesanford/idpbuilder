package builder

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// TarballWriter handles exporting OCI images as tarball files.
// It supports both single and multiple image exports.
type TarballWriter struct {
	options TarballOptions
}

// TarballOptions configures tarball export behavior.
type TarballOptions struct {
	Compress        bool
	IncludeManifest bool
}

// NewTarballWriter creates a new tarball writer with default settings.
func NewTarballWriter() *TarballWriter {
	return &TarballWriter{
		options: TarballOptions{
			Compress:        false,
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

// Write exports a single image to a tarball file.
func (w *TarballWriter) Write(img v1.Image, outputPath string, ref string) error {
	if img == nil {
		return fmt.Errorf("image cannot be nil")
	}

	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	if ref == "" {
		return fmt.Errorf("image reference cannot be empty")
	}

	// Parse the reference
	tag, err := name.NewTag(ref, name.WeakValidation)
	if err != nil {
		return fmt.Errorf("invalid reference %s: %w", ref, err)
	}

	// Create tag map for the image
	tagMap := map[name.Tag]v1.Image{
		tag: img,
	}

	// Write the tarball
	if err := tarball.MultiWriteToFile(outputPath, tagMap); err != nil {
		return fmt.Errorf("failed to write tarball: %w", err)
	}

	return nil
}

// WriteMultiple exports multiple images to a single tarball file.
func (w *TarballWriter) WriteMultiple(images map[string]v1.Image, outputPath string) error {
	if len(images) == 0 {
		return fmt.Errorf("no images provided for export")
	}

	// Build tag map
	tagMap := make(map[name.Tag]v1.Image)
	for ref, img := range images {
		tag, err := name.NewTag(ref, name.WeakValidation)
		if err != nil {
			return fmt.Errorf("invalid reference %s: %w", ref, err)
		}
		tagMap[tag] = img
	}

	// Write the tarball
	if err := tarball.MultiWriteToFile(outputPath, tagMap); err != nil {
		return fmt.Errorf("failed to write multi-image tarball: %w", err)
	}

	return nil
}

// WithCompression enables gzip compression of the output tarball.
func (w *TarballWriter) WithCompression(compress bool) *TarballWriter {
	w.options.Compress = compress
	return w
}
