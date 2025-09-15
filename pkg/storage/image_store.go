package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/layout"
)

// ImageStore provides local storage for OCI images
type ImageStore interface {
	Save(tag string, image v1.Image) error
	Load(tag string) (v1.Image, error)
	List() ([]string, error)
	Delete(tag string) error
	GetStoragePath() string
}

// LocalImageStore implements ImageStore using OCI layout on disk
type LocalImageStore struct {
	storageDir string
}

// NewLocalImageStore creates a new local image store
func NewLocalImageStore() (*LocalImageStore, error) {
	// Use a persistent directory instead of temp
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to /tmp if home dir not available
		return NewLocalImageStoreWithPath("/tmp/idpbuilder-images")
	}

	storageDir := filepath.Join(homeDir, ".idpbuilder", "images")
	return NewLocalImageStoreWithPath(storageDir)
}

// NewLocalImageStoreWithPath creates a new local image store at the specified path
func NewLocalImageStoreWithPath(storageDir string) (*LocalImageStore, error) {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory %s: %w", storageDir, err)
	}

	return &LocalImageStore{
		storageDir: storageDir,
	}, nil
}

// Save stores an image in the local store with the given tag
func (s *LocalImageStore) Save(tag string, image v1.Image) error {
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if image == nil {
		return fmt.Errorf("image cannot be nil")
	}

	// Create sanitized directory name from tag
	sanitizedTag := sanitizeTag(tag)
	imagePath := filepath.Join(s.storageDir, sanitizedTag)

	// Create layout path for this image
	layoutPath, err := layout.Write(imagePath, layout.NewIndex())
	if err != nil {
		return fmt.Errorf("failed to create layout path: %w", err)
	}

	// Write the image to the layout
	if err := layoutPath.AppendImage(image); err != nil {
		return fmt.Errorf("failed to write image to layout: %w", err)
	}

	// Write tag information
	tagFile := filepath.Join(imagePath, "tag.txt")
	if err := os.WriteFile(tagFile, []byte(tag), 0644); err != nil {
		return fmt.Errorf("failed to write tag file: %w", err)
	}

	return nil
}

// Load retrieves an image from the local store by tag
func (s *LocalImageStore) Load(tag string) (v1.Image, error) {
	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}

	sanitizedTag := sanitizeTag(tag)
	imagePath := filepath.Join(s.storageDir, sanitizedTag)

	// Check if the image directory exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image %s not found in local store", tag)
	}

	// Load the layout index
	layoutPath, err := layout.FromPath(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load layout from %s: %w", imagePath, err)
	}

	// Get the index
	index, err := layoutPath.ImageIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to get image index: %w", err)
	}

	// Get the index manifest
	indexManifest, err := index.IndexManifest()
	if err != nil {
		return nil, fmt.Errorf("failed to get index manifest: %w", err)
	}

	// We expect exactly one image in the index
	if len(indexManifest.Manifests) == 0 {
		return nil, fmt.Errorf("no images found in layout for tag %s", tag)
	}
	if len(indexManifest.Manifests) > 1 {
		return nil, fmt.Errorf("multiple images found in layout for tag %s, expected one", tag)
	}

	// Get the image by digest
	manifest := indexManifest.Manifests[0]
	image, err := index.Image(manifest.Digest)
	if err != nil {
		return nil, fmt.Errorf("failed to get image by digest: %w", err)
	}

	return image, nil
}

// List returns all tags stored in the local store
func (s *LocalImageStore) List() ([]string, error) {
	entries, err := os.ReadDir(s.storageDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read storage directory: %w", err)
	}

	var tags []string
	for _, entry := range entries {
		if entry.IsDir() {
			tagFile := filepath.Join(s.storageDir, entry.Name(), "tag.txt")
			if tagBytes, err := os.ReadFile(tagFile); err == nil {
				tags = append(tags, string(tagBytes))
			}
		}
	}

	return tags, nil
}

// Delete removes an image from the local store
func (s *LocalImageStore) Delete(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}

	sanitizedTag := sanitizeTag(tag)
	imagePath := filepath.Join(s.storageDir, sanitizedTag)

	if err := os.RemoveAll(imagePath); err != nil {
		return fmt.Errorf("failed to delete image %s: %w", tag, err)
	}

	return nil
}

// GetStoragePath returns the base storage directory path
func (s *LocalImageStore) GetStoragePath() string {
	return s.storageDir
}

// sanitizeTag converts a tag to a filesystem-safe directory name
func sanitizeTag(tag string) string {
	// Replace characters that are problematic for filesystems
	sanitized := strings.ReplaceAll(tag, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, ":", "_")
	sanitized = strings.ReplaceAll(sanitized, "@", "_")
	return sanitized
}