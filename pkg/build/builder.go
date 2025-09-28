package build

import (
	"context"
	"os"
)

// Builder defines the interface for container build operations
type Builder interface {
	// Build creates a new container build context
	Build(ctx context.Context, config *BuildConfig) (*BuildContext, error)

	// AddLayer adds a layer to the build
	AddLayer(ctx context.Context, buildCtx *BuildContext, layer *LayerSpec) error

	// Finalize completes the build and returns the result
	Finalize(ctx context.Context, buildCtx *BuildContext) (*BuildResult, error)
}

// NewBuilder creates a new Builder instance
func NewBuilder(opts *BuildOptions) (Builder, error) {
	if opts == nil {
		opts = &BuildOptions{}
	}

	// Validate and set default storage path
	if opts.StoragePath == "" {
		opts.StoragePath = os.Getenv("BUILDAH_STORAGE_PATH")
		if opts.StoragePath == "" {
			opts.StoragePath = os.Getenv("XDG_DATA_HOME")
			if opts.StoragePath == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return nil, WrapBuildError("newbuilder", err, "failed to get user home directory")
				}
				opts.StoragePath = homeDir + "/.local/share/containers/storage"
			} else {
				opts.StoragePath = opts.StoragePath + "/containers/storage"
			}
		}
	}

	// Validate and set default run root
	if opts.RunRoot == "" {
		opts.RunRoot = os.Getenv("BUILDAH_RUN_ROOT")
		if opts.RunRoot == "" {
			opts.RunRoot = "/tmp/buildah-run-root"
		}
	}

	return newBuildahAdapter(opts)
}