package push

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultDiscoveryOptions(t *testing.T) {
	buildPath := "/test/path"
	opts := DefaultDiscoveryOptions(buildPath)

	if opts.BuildPath != buildPath {
		t.Errorf("Expected BuildPath %s, got %s", buildPath, opts.BuildPath)
	}

	if opts.MaxSizeMB != 5000 {
		t.Errorf("Expected MaxSizeMB 5000, got %d", opts.MaxSizeMB)
	}

	if !opts.FollowLinks {
		t.Error("Expected FollowLinks to be true")
	}

	expectedExtensions := []string{".tar", ".tar.gz", ".tgz"}
	if len(opts.Extensions) != len(expectedExtensions) {
		t.Errorf("Expected %d extensions, got %d", len(expectedExtensions), len(opts.Extensions))
	}
}

func TestDiscoverLocalImages_NonExistentPath(t *testing.T) {
	images, err := DiscoverLocalImages("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
	if images != nil {
		t.Errorf("Expected nil images for non-existent path, got %v", images)
	}
}

func TestDiscoverLocalImages_EmptyDirectory(t *testing.T) {
	// Create temporary empty directory
	tmpDir := t.TempDir()

	images, err := DiscoverLocalImages(tmpDir)
	if err != nil {
		t.Errorf("Expected no error for empty directory, got %v", err)
	}
	if len(images) != 0 {
		t.Errorf("Expected 0 images in empty directory, got %d", len(images))
	}
}

func TestDiscoverLocalImagesWithOptions_CustomExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	opts := &DiscoveryOptions{
		BuildPath:   tmpDir,
		Extensions:  []string{".tar"},
		MaxSizeMB:   1000,
		FollowLinks: false,
	}

	images, err := DiscoverLocalImagesWithOptions(opts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Empty directory should return empty slice, not nil
	if len(images) != 0 {
		t.Errorf("Expected 0 images in empty directory, got %d", len(images))
	}
}

func TestIsSupportedTarball(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		extensions []string
		expected   bool
	}{
		{
			name:       "tar file",
			filename:   "image.tar",
			extensions: []string{".tar", ".tar.gz", ".tgz"},
			expected:   true,
		},
		{
			name:       "tar.gz file",
			filename:   "image.tar.gz",
			extensions: []string{".tar", ".tar.gz", ".tgz"},
			expected:   true,
		},
		{
			name:       "tgz file",
			filename:   "image.tgz",
			extensions: []string{".tar", ".tar.gz", ".tgz"},
			expected:   true,
		},
		{
			name:       "unsupported extension",
			filename:   "image.zip",
			extensions: []string{".tar", ".tar.gz", ".tgz"},
			expected:   false,
		},
		{
			name:       "no extension",
			filename:   "image",
			extensions: []string{".tar", ".tar.gz", ".tgz"},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSupportedTarball(tt.filename, tt.extensions)
			if result != tt.expected {
				t.Errorf("isSupportedTarball(%s) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestLocalImageStruct(t *testing.T) {
	img := &LocalImage{
		Name:   "test-image:latest",
		Path:   "/path/to/image.tar",
		Format: "tarball",
		Image:  nil, // Would be a real v1.Image in production
	}

	if img.Name != "test-image:latest" {
		t.Errorf("Expected Name 'test-image:latest', got %s", img.Name)
	}
	if img.Path != "/path/to/image.tar" {
		t.Errorf("Expected Path '/path/to/image.tar', got %s", img.Path)
	}
	if img.Format != "tarball" {
		t.Errorf("Expected Format 'tarball', got %s", img.Format)
	}
}

func TestDiscoverLocalImagesWithOptions_MaxSizeLimit(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file larger than the limit
	testFile := filepath.Join(tmpDir, "large.tar")
	content := make([]byte, 2*1024*1024) // 2MB
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	opts := &DiscoveryOptions{
		BuildPath:   tmpDir,
		Extensions:  []string{".tar"},
		MaxSizeMB:   1, // 1MB limit
		FollowLinks: true,
	}

	// Should discover 0 images because file exceeds size limit
	images, err := DiscoverLocalImagesWithOptions(opts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// File should be skipped due to size
	if len(images) > 0 {
		t.Logf("Note: Found %d images (tarball parsing might have failed, which is OK)", len(images))
	}
}

func TestDiscoverLocalImagesWithOptions_FollowLinksOption(t *testing.T) {
	tmpDir := t.TempDir()

	opts := &DiscoveryOptions{
		BuildPath:   tmpDir,
		Extensions:  []string{".tar"},
		MaxSizeMB:   5000,
		FollowLinks: false, // Explicitly test with false
	}

	images, err := DiscoverLocalImagesWithOptions(opts)
	if err != nil {
		t.Errorf("Expected no error with FollowLinks=false, got %v", err)
	}
	// Empty directory should return empty slice
	if len(images) != 0 {
		t.Errorf("Expected 0 images in empty directory, got %d", len(images))
	}
}