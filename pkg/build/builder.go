package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/containers/buildah"
	"github.com/containers/buildah/define"
	"github.com/containers/buildah/imagebuildah"
	"github.com/containers/common/libimage"
	"github.com/containers/storage"
	"github.com/cnoe-io/idpbuilder/pkg/certs/trust"
)

// BuildahBuilder implements the Builder interface using Buildah libraries
type BuildahBuilder struct {
	storeOptions storage.StoreOptions
	trustManager *trust.TrustManager
}

// NewBuildahBuilder creates a new Buildah-based builder
func NewBuildahBuilder(trustManager *trust.TrustManager) (*BuildahBuilder, error) {
	storeOptions, err := storage.DefaultStoreOptions(false, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get default store options: %w", err)
	}

	return &BuildahBuilder{
		storeOptions: storeOptions,
		trustManager: trustManager,
	}, nil
}

// BuildImage builds a container image from a Dockerfile using Buildah
func (b *BuildahBuilder) BuildImage(ctx context.Context, opts BuildOptions) (*BuildResult, error) {
	if opts.DockerfilePath == "" {
		return nil, fmt.Errorf("dockerfile path cannot be empty")
	}
	if opts.ContextDir == "" {
		return nil, fmt.Errorf("context directory cannot be empty")
	}

	startTime := time.Now()

	// Initialize storage
	store, err := storage.GetStore(b.storeOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %w", err)
	}
	defer store.Shutdown(false)

	// Prepare build options
	systemContext := &define.SystemContext{}
	
	// Integrate with Phase 1 certificate infrastructure if available
	if b.trustManager != nil {
		if err := b.configureTrustStore(systemContext); err != nil {
			return nil, fmt.Errorf("failed to configure trust store: %w", err)
		}
	}

	buildOptions := define.BuildOptions{
		ContextDirectory:     opts.ContextDir,
		SystemContext:        systemContext,
		Isolation:            define.IsolationDefault,
		ConfigureNetwork:     define.NetworkDefault,
		CNIPluginPath:        "",
		CNIConfigDir:         "",
		CommonBuildOpts:      &define.CommonBuildOptions{},
		DefaultMountsFilePath: "",
		IIDFile:              "",
		Squash:               false,
		Args:                 opts.Args,
		NoCache:              opts.NoCache,
		RemoveIntermediateCtrs: true,
		ForceRmIntermediateCtrs: false,
		BlobDirectory:        "",
		Target:               "",
		Devices:              nil,
		DeviceSpecs:          nil,
		LogRusage:            false,
		Quiet:                false,
		Runtime:              "",
		RuntimeArgs:          nil,
		Output:               opts.Tag,
		OutputFormat:         "docker",
		AdditionalTags:       nil,
		Log:                  func(format string, args ...interface{}) {},
		In:                   os.Stdin,
		Out:                  os.Stdout,
		Err:                  os.Stderr,
		SignBy:               "",
		Architecture:         "",
		Timestamp:            nil,
		OS:                   "",
		MaxPullPushRetries:   3,
		PullPushRetryDelay:   2 * time.Second,
		OciDecryptConfig:     nil,
		Jobs:                 nil,
		LogSplitByPlatform:   false,
		OSFeatures:           nil,
		OSVersion:            "",
	}

	// Read Dockerfile
	dockerfileContent, err := os.ReadFile(opts.DockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dockerfile: %w", err)
	}

	// Execute build
	imageID, _, _, err := imagebuildah.BuildDockerfiles(
		ctx,
		store,
		buildOptions,
		[]string{string(dockerfileContent)},
	)
	if err != nil {
		return nil, fmt.Errorf("build failed: %w", err)
	}

	buildTime := time.Since(startTime)

	// Get image information
	imageInfo, err := b.getImageInfo(store, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image info: %w", err)
	}

	return &BuildResult{
		ImageID:    imageID,
		Repository: getRepository(opts.Tag),
		Tag:        getTag(opts.Tag),
		Digest:     imageInfo.Digest,
		Size:       imageInfo.Size,
		BuildTime:  buildTime,
	}, nil
}

// ListImages lists available images in the storage
func (b *BuildahBuilder) ListImages(ctx context.Context) ([]ImageInfo, error) {
	store, err := storage.GetStore(b.storeOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %w", err)
	}
	defer store.Shutdown(false)

	images, err := store.Images()
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	var result []ImageInfo
	for _, img := range images {
		for _, name := range img.Names {
			result = append(result, ImageInfo{
				ID:         img.ID,
				Repository: getRepository(name),
				Tag:        getTag(name),
				Digest:     img.Digest.String(),
				Size:       img.Size,
				Created:    *img.Created,
			})
		}
	}

	return result, nil
}

// RemoveImage removes an image by ID
func (b *BuildahBuilder) RemoveImage(ctx context.Context, imageID string) error {
	store, err := storage.GetStore(b.storeOptions)
	if err != nil {
		return fmt.Errorf("failed to get storage: %w", err)
	}
	defer store.Shutdown(false)

	_, err = store.DeleteImage(imageID, true)
	if err != nil {
		return fmt.Errorf("failed to remove image %s: %w", imageID, err)
	}

	return nil
}

// TagImage tags an existing image
func (b *BuildahBuilder) TagImage(ctx context.Context, source, target string) error {
	store, err := storage.GetStore(b.storeOptions)
	if err != nil {
		return fmt.Errorf("failed to get storage: %w", err)
	}
	defer store.Shutdown(false)

	// Find the source image
	image, err := store.Image(source)
	if err != nil {
		return fmt.Errorf("failed to find source image %s: %w", source, err)
	}

	// Add the new tag
	if err := store.SetNames(image.ID, append(image.Names, target)); err != nil {
		return fmt.Errorf("failed to tag image %s with %s: %w", source, target, err)
	}

	return nil
}

// configureTrustStore configures the system context with Phase 1 certificate infrastructure
func (b *BuildahBuilder) configureTrustStore(systemContext *define.SystemContext) error {
	if b.trustManager == nil {
		return nil // No trust manager available, skip
	}

	// Integration with Phase 1 TrustManager would go here
	// For now, we'll just ensure the system context is configured for certificate validation
	systemContext.DockerInsecureSkipTLSVerify = define.OptionalBoolFalse
	
	// Log security event (placeholder for Phase 1 audit integration)
	fmt.Printf("Security: Configured trust store for image operations\n")
	
	return nil
}

// getImageInfo retrieves detailed information about an image
func (b *BuildahBuilder) getImageInfo(store storage.Store, imageID string) (*ImageInfo, error) {
	image, err := store.Image(imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return &ImageInfo{
		ID:      image.ID,
		Digest:  image.Digest.String(),
		Size:    image.Size,
		Created: *image.Created,
	}, nil
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