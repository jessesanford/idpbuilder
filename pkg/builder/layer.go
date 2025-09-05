package builder

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// LayerFactory creates OCI layers from directory contents.
// It handles file walking, tar archive creation, and metadata preservation.
type LayerFactory struct {
	preservePermissions bool
	preserveTimestamps  bool
}

// NewLayerFactory creates a new layer factory with default settings.
// It preserves permissions but normalizes timestamps for reproducible builds.
func NewLayerFactory() *LayerFactory {
	return &LayerFactory{
		preservePermissions: true,
		preserveTimestamps:  false, // Normalize for reproducible builds
	}
}

// CreateLayer builds an OCI layer from directory contents.
// It walks the directory tree, creates a tar archive, and returns it as a v1.Layer.
func (f *LayerFactory) CreateLayer(contextDir string) (v1.Layer, error) {
	// Validate context directory
	if contextDir == "" {
		return nil, fmt.Errorf("context directory cannot be empty")
	}
	
	// Clean the path to ensure consistent handling
	contextDir = filepath.Clean(contextDir)
	
	// Create tar archive in memory
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()
	
	// Walk the directory tree and add files to tar
	err := filepath.WalkDir(contextDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}
		
		// Skip the root directory itself
		if path == contextDir {
			return nil
		}
		
		// Get file info for metadata
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("failed to get info for %s: %w", path, err)
		}
		
		// Calculate relative path within the context
		relPath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path for %s: %w", path, err)
		}
		
		// Convert Windows paths to Unix paths for OCI compliance
		tarPath := filepath.ToSlash(relPath)
		
		// Add file to tar archive
		if err := f.addFileToTar(tw, path, tarPath, info); err != nil {
			return fmt.Errorf("failed to add file %s to tar: %w", path, err)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk context directory: %w", err)
	}
	
	// Close tar writer to finalize archive
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}
	
	// Create layer from tar archive
	reader := bytes.NewReader(buf.Bytes())
	layer, err := tarball.LayerFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from tar: %w", err)
	}
	
	return layer, nil
}

// addFileToTar adds a single file to the tar archive with appropriate metadata.
// It handles regular files, directories, and symlinks according to OCI specifications.
func (f *LayerFactory) addFileToTar(tw *tar.Writer, srcPath, tarPath string, info os.FileInfo) error {
	// Create tar header from file info
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("failed to create tar header: %w", err)
	}
	
	// Set the name in the tar archive
	header.Name = tarPath
	
	// Handle timestamps for reproducible builds
	if !f.preserveTimestamps {
		// Use a fixed timestamp for reproducible builds
		fixedTime := time.Unix(0, 0)
		header.ModTime = fixedTime
		header.AccessTime = fixedTime
		header.ChangeTime = fixedTime
	}
	
	// Handle permissions
	if f.preservePermissions {
		// Keep original permissions
		header.Mode = int64(info.Mode())
	} else {
		// Normalize permissions
		if info.IsDir() {
			header.Mode = 0755
		} else {
			header.Mode = 0644
		}
	}
	
	// Handle different file types
	switch info.Mode() & os.ModeType {
	case 0: // Regular file
		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header for regular file: %w", err)
		}
		
		// Open and copy file content
		file, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("failed to open file for reading: %w", err)
		}
		defer file.Close()
		
		_, err = io.Copy(tw, file)
		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}
		
	case os.ModeDir: // Directory
		// Ensure directory path ends with slash for tar compliance
		if !strings.HasSuffix(header.Name, "/") {
			header.Name += "/"
		}
		header.Size = 0
		
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header for directory: %w", err)
		}
		
	case os.ModeSymlink: // Symbolic link
		// Read link target
		linkTarget, err := os.Readlink(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read symlink target: %w", err)
		}
		
		header.Linkname = linkTarget
		header.Size = 0
		
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header for symlink: %w", err)
		}
		
	default:
		// Skip special files (devices, named pipes, etc.)
		// These are not commonly needed in container images
		return nil
	}
	
	return nil
}

// WithPermissionPreservation configures whether to preserve original file permissions.
func (f *LayerFactory) WithPermissionPreservation(preserve bool) *LayerFactory {
	f.preservePermissions = preserve
	return f
}

// WithTimestampPreservation configures whether to preserve original file timestamps.
func (f *LayerFactory) WithTimestampPreservation(preserve bool) *LayerFactory {
	f.preserveTimestamps = preserve
	return f
}

// Simple helper functions for basic file operations
func isExecutable(info os.FileInfo) bool {
	return info.Mode()&0111 != 0
}

func getFileOwnership(info os.FileInfo) (int, int) {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return int(stat.Uid), int(stat.Gid)
	}
	return 0, 0
}