package build

import "time"

// BuildConfig specifies the configuration for a build
type BuildConfig struct {
	BaseImage    string            // Base image reference
	WorkingDir   string            // Working directory for build
	Env          map[string]string // Environment variables
	Labels       map[string]string // Container labels
	Entrypoint   []string          // Container entrypoint
	Cmd          []string          // Container command
}

// LayerSpec defines a layer to add to the container
type LayerSpec struct {
	Source      string    // Source path or URL
	Destination string    // Destination in container
	Type        LayerType // Type of layer operation
	Permissions string    // File permissions
	Owner       string    // File ownership
}

// BuildResult contains the result of a successful build
type BuildResult struct {
	ImageID   string    // Generated image ID
	Digest    string    // Image digest
	Size      int64     // Image size in bytes
	CreatedAt time.Time // Build timestamp
}

// BuildOptions configures the builder behavior
type BuildOptions struct {
	StoragePath string // Path to storage directory
	RunRoot     string // Runtime root directory
	Debug       bool   // Enable debug logging
}

// LayerType defines the type of layer operation
type LayerType string

const (
	LayerTypeCopy LayerType = "copy"
	LayerTypeAdd  LayerType = "add"
	LayerTypeRun  LayerType = "run"
)

// BuildContext holds the state of an in-progress build
type BuildContext struct {
	ID         string                 // Unique build identifier
	WorkingDir string                 // Current working directory
	Metadata   map[string]interface{} // Additional build metadata
	internal   interface{}            // Internal Buildah state (opaque)
}