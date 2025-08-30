package build

import (
	"context"
	"time"
)

// Builder handles container image building operations using Buildah
type Builder interface {
	// BuildImage builds a container image from a Dockerfile
	BuildImage(ctx context.Context, opts BuildOptions) (*BuildResult, error)
	
	// ListImages lists available images
	ListImages(ctx context.Context) ([]ImageInfo, error)
	
	// RemoveImage removes an image by ID
	RemoveImage(ctx context.Context, imageID string) error
	
	// TagImage tags an existing image
	TagImage(ctx context.Context, source, target string) error
}

// BuildOptions contains options for building container images
type BuildOptions struct {
	DockerfilePath string
	ContextDir     string
	Tag            string
	Args           map[string]string
	NoCache        bool
	Insecure       bool
}

// BuildResult contains the result of a successful build
type BuildResult struct {
	ImageID      string
	Repository   string
	Tag          string
	Digest       string
	Size         int64
	BuildTime    time.Duration
}

// ImageInfo contains information about an available image
type ImageInfo struct {
	ID         string
	Repository string
	Tag        string
	Digest     string
	Size       int64
	Created    time.Time
}