package buildah

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// BuildContext represents a Docker build context
type BuildContext interface {
	// GetRoot returns the root directory of the build context
	GetRoot() string

	// AddFile adds a file to the build context
	AddFile(path string, content []byte) error

	// AddDirectory adds a directory recursively to the build context
	AddDirectory(src, dest string) error

	// Validate ensures the build context is valid
	Validate() error

	// Clean removes temporary context files
	Clean() error

	// GetSize returns the total size of the context in bytes
	GetSize() (int64, error)
}

// BuildContextManager creates and manages build contexts
type BuildContextManager interface {
	// CreateContext creates a new build context
	CreateContext(ctx context.Context, dockerfilePath string) (BuildContext, error)

	// CreateFromDirectory creates a context from an existing directory
	CreateFromDirectory(ctx context.Context, dir string) (BuildContext, error)

	// CreateTarball creates a tarball from the context
	CreateTarball(ctx BuildContext) (string, error)
}

// contextImpl implements BuildContext
type contextImpl struct {
	root    string
	files   map[string][]byte
	tempDir bool
}

// NewBuildContext creates a new build context
func NewBuildContext(root string) BuildContext {
	return &contextImpl{
		root:  root,
		files: make(map[string][]byte),
	}
}

// GetRoot returns the root directory
func (c *contextImpl) GetRoot() string {
	return c.root
}

// AddFile adds a file to the context
func (c *contextImpl) AddFile(path string, content []byte) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Clean the path and ensure it's relative
	cleanPath := filepath.Clean(path)
	if filepath.IsAbs(cleanPath) {
		return fmt.Errorf("file path must be relative: %s", path)
	}

	// Store in memory map for tracking
	c.files[cleanPath] = content

	// Write to actual file system
	fullPath := filepath.Join(c.root, cleanPath)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return nil
}

// AddDirectory adds a directory recursively to the build context
func (c *contextImpl) AddDirectory(src, dest string) error {
	if src == "" || dest == "" {
		return fmt.Errorf("source and destination paths cannot be empty")
	}

	// Clean paths
	cleanSrc := filepath.Clean(src)
	cleanDest := filepath.Clean(dest)

	if filepath.IsAbs(cleanDest) {
		return fmt.Errorf("destination path must be relative: %s", dest)
	}

	return filepath.WalkDir(cleanSrc, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(cleanSrc, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Join with destination
		destPath := filepath.Join(cleanDest, relPath)

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Add to context
		return c.AddFile(destPath, content)
	})
}

// Validate ensures the build context is valid
func (c *contextImpl) Validate() error {
	// Check if Dockerfile exists
	dockerfilePaths := []string{"Dockerfile", "dockerfile", "Dockerfile.build"}
	var dockerfileFound bool

	for _, dockerfilePath := range dockerfilePaths {
		fullPath := filepath.Join(c.root, dockerfilePath)
		if _, err := os.Stat(fullPath); err == nil {
			dockerfileFound = true
			break
		}
		// Also check in memory files
		if _, exists := c.files[dockerfilePath]; exists {
			dockerfileFound = true
			break
		}
	}

	if !dockerfileFound {
		return fmt.Errorf("no Dockerfile found in build context")
	}

	// Check if root directory exists
	if _, err := os.Stat(c.root); os.IsNotExist(err) {
		return fmt.Errorf("build context root directory does not exist: %s", c.root)
	}

	return nil
}

// Clean removes temporary context files
func (c *contextImpl) Clean() error {
	if c.tempDir {
		return os.RemoveAll(c.root)
	}
	return nil
}

// GetSize returns the total size of the context in bytes
func (c *contextImpl) GetSize() (int64, error) {
	var totalSize int64

	err := filepath.WalkDir(c.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to calculate context size: %w", err)
	}

	return totalSize, nil
}

// contextManagerImpl implements BuildContextManager
type contextManagerImpl struct {
	workDir string
}

// NewBuildContextManager creates a new context manager
func NewBuildContextManager(workDir string) BuildContextManager {
	return &contextManagerImpl{
		workDir: workDir,
	}
}

// CreateContext creates a new build context
func (m *contextManagerImpl) CreateContext(ctx context.Context, dockerfilePath string) (BuildContext, error) {
	if dockerfilePath == "" {
		return nil, fmt.Errorf("dockerfile path cannot be empty")
	}

	// Check if Dockerfile exists
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dockerfile not found: %s", dockerfilePath)
	}

	// Get directory containing the Dockerfile
	contextRoot := filepath.Dir(dockerfilePath)

	// Create build context
	buildCtx := NewBuildContext(contextRoot)

	return buildCtx, nil
}

// CreateFromDirectory creates a context from an existing directory
func (m *contextManagerImpl) CreateFromDirectory(ctx context.Context, dir string) (BuildContext, error) {
	if dir == "" {
		return nil, fmt.Errorf("directory path cannot be empty")
	}

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory not found: %s", dir)
	}

	// Create build context
	buildCtx := NewBuildContext(dir)

	return buildCtx, nil
}

// CreateTarball creates a tarball from the context
func (m *contextManagerImpl) CreateTarball(ctx BuildContext) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("build context cannot be nil")
	}

	// For now, return the path to the context directory
	// In a full implementation, this would create an actual tarball
	return ctx.GetRoot(), nil
}
