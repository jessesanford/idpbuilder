package build

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/containers/buildah"
	"github.com/containers/buildah/define"
	"github.com/containers/common/libimage"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
)

// buildahAdapter implements the Builder interface using Buildah
type buildahAdapter struct {
	store   storage.Store
	options *BuildOptions
}

// newBuildahAdapter creates a new buildah-based builder implementation
func newBuildahAdapter(opts *BuildOptions) (Builder, error) {
	if opts == nil {
		return nil, WrapBuildError("adapter_init", ErrInvalidConfig, "build options cannot be nil")
	}

	// Ensure we're running in the right namespace for rootless containers
	unshare.MaybeReexecUsingUserNamespace(false)

	// Initialize storage
	storeOpts, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())
	if err != nil {
		return nil, WrapBuildError("storage_init", err, "failed to get default storage options")
	}

	// Override with provided options
	storeOpts.GraphRoot = opts.StoragePath
	storeOpts.RunRoot = opts.RunRoot

	// Ensure directories exist
	if err := os.MkdirAll(storeOpts.GraphRoot, 0755); err != nil {
		return nil, WrapBuildError("storage_init", err, fmt.Sprintf("failed to create storage directory: %s", storeOpts.GraphRoot))
	}
	if err := os.MkdirAll(storeOpts.RunRoot, 0755); err != nil {
		return nil, WrapBuildError("storage_init", err, fmt.Sprintf("failed to create run directory: %s", storeOpts.RunRoot))
	}

	store, err := storage.GetStore(storeOpts)
	if err != nil {
		return nil, WrapBuildError("storage_init", err, "failed to initialize container storage")
	}

	return &buildahAdapter{
		store:   store,
		options: opts,
	}, nil
}

// Build creates a new container build context
func (b *buildahAdapter) Build(ctx context.Context, config *BuildConfig) (*BuildContext, error) {
	if config == nil {
		return nil, WrapBuildError("build", ErrInvalidConfig, "build configuration cannot be nil")
	}

	if config.BaseImage == "" {
		return nil, WrapBuildError("build", ErrInvalidConfig, "base image must be specified")
	}

	// Generate unique build ID
	buildID, err := generateBuildID()
	if err != nil {
		return nil, WrapBuildError("build", err, "failed to generate build ID")
	}

	// Create buildah builder options
	buildOpts := buildah.BuilderOptions{
		FromImage:        config.BaseImage,
		Container:        buildID,
		PullPolicy:       define.PullIfMissing,
		IsolationDefault: define.IsolationDefault,
	}

	// Create the builder
	builder, err := buildah.NewBuilder(ctx, b.store, buildOpts)
	if err != nil {
		return nil, WrapBuildError("build", err, fmt.Sprintf("failed to create builder for image: %s", config.BaseImage))
	}

	// Configure the container
	if config.WorkingDir != "" {
		builder.SetWorkDir(config.WorkingDir)
	}

	// Set environment variables
	for key, value := range config.Env {
		builder.SetEnv(key, value)
	}

	// Set labels
	for key, value := range config.Labels {
		builder.SetLabel(key, value)
	}

	// Set entrypoint
	if len(config.Entrypoint) > 0 {
		builder.SetEntrypoint(config.Entrypoint)
	}

	// Set command
	if len(config.Cmd) > 0 {
		builder.SetCmd(config.Cmd)
	}

	// Create build context
	buildCtx := &BuildContext{
		ID:         buildID,
		WorkingDir: config.WorkingDir,
		Metadata: map[string]interface{}{
			"base_image": config.BaseImage,
			"created_at": time.Now(),
		},
		internal: builder,
	}

	return buildCtx, nil
}

// AddLayer adds a layer to the build
func (b *buildahAdapter) AddLayer(ctx context.Context, buildCtx *BuildContext, layer *LayerSpec) error {
	if buildCtx == nil {
		return WrapBuildError("add_layer", ErrInvalidConfig, "build context cannot be nil")
	}

	if layer == nil {
		return WrapBuildError("add_layer", ErrInvalidConfig, "layer specification cannot be nil")
	}

	builder, ok := buildCtx.internal.(*buildah.Builder)
	if !ok {
		return WrapBuildError("add_layer", ErrLayerAddFailed, "invalid build context internal state")
	}

	switch layer.Type {
	case LayerTypeCopy:
		return b.addCopyLayer(ctx, builder, layer)
	case LayerTypeAdd:
		return b.addAddLayer(ctx, builder, layer)
	case LayerTypeRun:
		return b.addRunLayer(ctx, builder, layer)
	default:
		return WrapBuildError("add_layer", ErrLayerAddFailed, fmt.Sprintf("unsupported layer type: %s", layer.Type))
	}
}

// addCopyLayer adds a COPY layer to the container
func (b *buildahAdapter) addCopyLayer(ctx context.Context, builder *buildah.Builder, layer *LayerSpec) error {
	if layer.Source == "" || layer.Destination == "" {
		return WrapBuildError("copy_layer", ErrLayerAddFailed, "source and destination must be specified for COPY operation")
	}

	// Check if source exists
	if _, err := os.Stat(layer.Source); err != nil {
		return WrapBuildError("copy_layer", err, fmt.Sprintf("source file does not exist: %s", layer.Source))
	}

	// Add the file to the container
	err := builder.Add(layer.Destination, false, buildah.AddAndCopyOptions{}, layer.Source)
	if err != nil {
		return WrapBuildError("copy_layer", err, fmt.Sprintf("failed to copy %s to %s", layer.Source, layer.Destination))
	}

	return nil
}

// addAddLayer adds an ADD layer to the container
func (b *buildahAdapter) addAddLayer(ctx context.Context, builder *buildah.Builder, layer *LayerSpec) error {
	if layer.Source == "" || layer.Destination == "" {
		return WrapBuildError("add_layer", ErrLayerAddFailed, "source and destination must be specified for ADD operation")
	}

	// ADD operation with extraction capabilities
	err := builder.Add(layer.Destination, true, buildah.AddAndCopyOptions{}, layer.Source)
	if err != nil {
		return WrapBuildError("add_layer", err, fmt.Sprintf("failed to add %s to %s", layer.Source, layer.Destination))
	}

	return nil
}

// addRunLayer executes a RUN command in the container
func (b *buildahAdapter) addRunLayer(ctx context.Context, builder *buildah.Builder, layer *LayerSpec) error {
	if layer.Source == "" {
		return WrapBuildError("run_layer", ErrLayerAddFailed, "command must be specified for RUN operation")
	}

	// Parse command
	var command []string
	if strings.Contains(layer.Source, " ") {
		// Shell form
		command = []string{"/bin/sh", "-c", layer.Source}
	} else {
		// Exec form (single command)
		command = []string{layer.Source}
	}

	// Run the command
	err := builder.Run(command, buildah.RunOptions{})
	if err != nil {
		return WrapBuildError("run_layer", err, fmt.Sprintf("failed to execute command: %s", layer.Source))
	}

	return nil
}

// Finalize completes the build and returns the result
func (b *buildahAdapter) Finalize(ctx context.Context, buildCtx *BuildContext) (*BuildResult, error) {
	if buildCtx == nil {
		return nil, WrapBuildError("finalize", ErrInvalidConfig, "build context cannot be nil")
	}

	builder, ok := buildCtx.internal.(*buildah.Builder)
	if !ok {
		return nil, WrapBuildError("finalize", ErrFinalizeFailed, "invalid build context internal state")
	}

	// Commit the container to create an image
	imageRef, err := builder.Commit(ctx, "", buildah.CommitOptions{})
	if err != nil {
		return nil, WrapBuildError("finalize", err, "failed to commit container")
	}

	// Get image information
	runtime, err := libimage.RuntimeFromStore(b.store, nil)
	if err != nil {
		return nil, WrapBuildError("finalize", err, "failed to create runtime from store")
	}

	image, _, err := runtime.LookupImage(imageRef, nil)
	if err != nil {
		return nil, WrapBuildError("finalize", err, "failed to lookup committed image")
	}

	// Get image size
	size, err := image.Size()
	if err != nil {
		// Log error but don't fail the build
		if b.options.Debug {
			fmt.Fprintf(os.Stderr, "Warning: failed to get image size: %v\n", err)
		}
		size = -1
	}

	// Get image digest
	digest := ""
	if image.Digest() != nil {
		digest = image.Digest().String()
	}

	result := &BuildResult{
		ImageID:   imageRef,
		Digest:    digest,
		Size:      size,
		CreatedAt: time.Now(),
	}

	return result, nil
}

// generateBuildID creates a unique identifier for the build
func generateBuildID() (string, error) {
	// Generate 8 random bytes
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert to hex string with timestamp prefix
	timestamp := time.Now().Unix()
	return fmt.Sprintf("build-%d-%x", timestamp, bytes), nil
}