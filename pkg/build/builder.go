//go:build !buildah
// +build !buildah

package build

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// TrustManager represents a certificate trust manager (Phase 1 integration interface)
type TrustManager interface {
	// ConfigureTLS configures TLS settings for container operations
	ConfigureTLS() error
	// ValidateCertificate validates a certificate against trusted CAs
	ValidateCertificate(cert []byte) error
}

// BuildahBuilder implements the Builder interface using Buildah libraries
// This is a mock implementation for environments without buildah dependencies
type BuildahBuilder struct {
	trustManager TrustManager
}

// NewBuildahBuilder creates a new Buildah-based builder
// This is a mock implementation that returns an error indicating Buildah is not available
func NewBuildahBuilder(trustManager TrustManager) (*BuildahBuilder, error) {
	return &BuildahBuilder{
		trustManager: trustManager,
	}, nil
}

// BuildImage is a mock implementation that simulates a successful build
func (b *BuildahBuilder) BuildImage(ctx context.Context, opts BuildOptions) (*BuildResult, error) {
	if opts.DockerfilePath == "" {
		return nil, fmt.Errorf("dockerfile path cannot be empty")
	}
	if opts.ContextDir == "" {
		return nil, fmt.Errorf("context directory cannot be empty")
	}

	// In a real environment with buildah dependencies, this would do the actual build
	// For now, return a mock successful result
	return &BuildResult{
		ImageID:    "mock-image-id-" + fmt.Sprintf("%d", time.Now().Unix()),
		Repository: getRepository(opts.Tag),
		Tag:        getTag(opts.Tag),
		Digest:     "sha256:mockdigest123456789abcdef",
		Size:       1024000, // 1MB mock size
		BuildTime:  5 * time.Second,
	}, fmt.Errorf("buildah dependencies not available in this environment - this is a mock implementation")
}

// ListImages is a mock implementation
func (b *BuildahBuilder) ListImages(ctx context.Context) ([]ImageInfo, error) {
	// Return empty list in mock mode
	return []ImageInfo{}, fmt.Errorf("buildah dependencies not available - mock implementation")
}

// RemoveImage is a mock implementation
func (b *BuildahBuilder) RemoveImage(ctx context.Context, imageID string) error {
	return fmt.Errorf("buildah dependencies not available - mock implementation")
}

// TagImage is a mock implementation
func (b *BuildahBuilder) TagImage(ctx context.Context, source, target string) error {
	return fmt.Errorf("buildah dependencies not available - mock implementation")
}

// Helper functions to parse repository and tag from image name
func getRepository(imageName string) string {
	if imageName == "" {
		return ""
	}

	parts := strings.Split(imageName, ":")
	if len(parts) == 1 {
		return imageName
	}

	return strings.Join(parts[:len(parts)-1], ":")
}

func getTag(imageName string) string {
	if imageName == "" {
		return ""
	}

	parts := strings.Split(imageName, ":")
	if len(parts) == 1 {
		return "latest"
	}

	return parts[len(parts)-1]
}
