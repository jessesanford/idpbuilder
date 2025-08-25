// Package api provides core interfaces and types for OCI operations in IDPBuilder.
// This package defines service contracts for OCI build operations, registry management,
// and foundational OCI configuration.
package api

import (
	"context"
)

// OCIBuildService defines the interface for OCI image build operations.
// This service abstracts buildah functionality and provides a clean API
// for building OCI-compliant images.
type OCIBuildService interface {
	// Initialize prepares the build service with the given configuration.
	// It sets up storage backends, runtime configurations, and validates prerequisites.
	Initialize(ctx context.Context, config *BuildConfig) error

	// BuildImage builds an OCI image from the provided build request.
	// This is the primary build method that handles containerfile processing,
	// layer creation, and image finalization.
	BuildImage(ctx context.Context, req *BuildRequest) (*BuildResult, error)

	// BuildFromDockerfile builds an image specifically from a Dockerfile path.
	// This method provides a convenient wrapper for Dockerfile-based builds.
	BuildFromDockerfile(ctx context.Context, dockerfilePath string, opts *BuildOptions) (*BuildResult, error)

	// GetBuildStatus returns the current status of an active build operation.
	// Use the build ID returned from BuildImage to query status.
	GetBuildStatus(ctx context.Context, buildID string) (*BuildStatus, error)

	// ListBuilds returns a list of recent build operations and their status.
	// Useful for monitoring and debugging build operations.
	ListBuilds(ctx context.Context) ([]*BuildStatus, error)

	// CleanupBuild removes temporary resources associated with a build.
	// Should be called after a build completes or fails to clean up resources.
	CleanupBuild(ctx context.Context, buildID string) error

	// ValidateConfig validates the provided build configuration without initializing.
	// Use this to check configuration validity before calling Initialize.
	ValidateConfig(config *BuildConfig) error

	// Close shuts down the build service and releases all resources.
	// This should be called when the service is no longer needed.
	Close() error
}

// OCIRegistryService defines the interface for OCI registry operations.
// This service handles pushing, pulling, and managing images in OCI registries.
type OCIRegistryService interface {
	// Connect establishes a connection to the registry with the given configuration.
	// This method handles authentication and validates connectivity.
	Connect(ctx context.Context, config *RegistryConfig) error

	// PushImage pushes a built image to the configured registry.
	// The image must be available in local storage before pushing.
	PushImage(ctx context.Context, imageID string, tags []string, opts *PushOptions) error

	// PullImage pulls an image from the registry to local storage.
	// Returns the local image ID after successful pull.
	PullImage(ctx context.Context, imageRef string, opts *PullOptions) (string, error)

	// ListImages returns a list of images available in the registry.
	// Optionally filters by repository name or tag pattern.
	ListImages(ctx context.Context, repository string) ([]*ImageInfo, error)

	// GetImageInfo retrieves detailed information about a specific image.
	// Includes metadata, layers, and manifest information.
	GetImageInfo(ctx context.Context, imageRef string) (*ImageInfo, error)

	// DeleteImage removes an image from the registry.
	// This operation may not be supported by all registry implementations.
	DeleteImage(ctx context.Context, imageRef string) error

	// ValidateConnection tests the registry connection without performing operations.
	// Useful for health checks and configuration validation.
	ValidateConnection(ctx context.Context) error

	// GetRegistry returns the current registry configuration.
	GetRegistry() *RegistryConfig

	// Close terminates the connection and releases resources.
	Close() error
}

// LayerProcessor defines the interface for custom layer processing during builds.
// This allows for custom layer manipulation, optimization, or analysis.
type LayerProcessor interface {
	// ProcessLayer processes a single layer during the build process.
	// Can modify layer contents, add metadata, or perform validation.
	ProcessLayer(ctx context.Context, layer *LayerInfo) error

	// ValidateLayer validates a layer meets specific requirements.
	// Can check size limits, content policies, or security constraints.
	ValidateLayer(layer *LayerInfo) error

	// OptimizeLayer performs optimizations on a layer.
	// May include compression, deduplication, or cleanup operations.
	OptimizeLayer(ctx context.Context, layer *LayerInfo) (*LayerInfo, error)
}