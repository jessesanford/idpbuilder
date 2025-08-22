// Package buildah provides container building capabilities using Buildah
package buildah

import (
	"context"
	"io"
	"time"

	"github.com/cnoe-io/idpbuilder/api/v1alpha1"
)

// Client provides buildah-based container building
type Client interface {
	// Build builds a container image from a Dockerfile
	Build(ctx context.Context, opts BuildOptions) (*BuildResult, error)

	// Push pushes a container image to a registry
	Push(ctx context.Context, opts PushOptions) error

	// Close cleans up client resources
	Close() error
}

// BuildOptions configures a container build
type BuildOptions struct {
	// ContextDir is the path to the build context directory
	ContextDir string

	// Dockerfile is the path to the Dockerfile (relative to ContextDir)
	Dockerfile string

	// Tag is the target image tag
	Tag string

	// BuildArgs are build-time variables passed to the build
	BuildArgs map[string]string

	// Progress is the writer for build progress output
	Progress io.Writer

	// Config contains build customization settings
	Config v1alpha1.BuildCustomizationSpec

	// Labels are metadata labels to apply to the image
	Labels map[string]string

	// NoCache disables build cache if true
	NoCache bool

	// Pull always pulls base images if true
	Pull bool
}

// PushOptions configures a container image push
type PushOptions struct {
	// Image is the image name/tag to push
	Image string

	// Registry is the target registry
	Registry string

	// Username for registry authentication
	Username string

	// Password for registry authentication
	Password string

	// Progress is the writer for push progress output
	Progress io.Writer
}

// BuildResult contains the result of a build operation
type BuildResult struct {
	// ImageID is the ID of the built image
	ImageID string

	// Tag is the tag applied to the image
	Tag string

	// Size is the size of the built image in bytes
	Size int64

	// Duration is the time taken for the build
	Duration time.Duration

	// BuildArgs contains the build arguments that were used
	BuildArgs map[string]string

	// Labels contains the labels applied to the image
	Labels map[string]string
}

// ClientOptions configures a buildah client
type ClientOptions struct {
	// BuildahPath is the path to the buildah binary (optional, defaults to "buildah")
	BuildahPath string

	// WorkDir is the working directory for builds (optional, uses temp dir if empty)
	WorkDir string

	// DefaultRegistry is the default registry for image operations
	DefaultRegistry string

	// Timeout is the default timeout for build operations
	Timeout time.Duration

	// LogLevel controls the verbosity of buildah output
	LogLevel string
}

// BuildError represents an error during a build operation
type BuildError struct {
	// Op is the operation that failed
	Op string

	// Err is the underlying error
	Err error

	// Output is any relevant output from the buildah command
	Output string

	// ExitCode is the exit code from the buildah command
	ExitCode int
}

func (e *BuildError) Error() string {
	if e.Output != "" {
		return e.Op + ": " + e.Err.Error() + "\n" + e.Output
	}
	return e.Op + ": " + e.Err.Error()
}

func (e *BuildError) Unwrap() error {
	return e.Err
}

// ProgressEvent represents a build progress event
type ProgressEvent struct {
	// Type is the type of progress event (step, output, error, etc.)
	Type string

	// Message is the progress message
	Message string

	// Step is the current build step (if applicable)
	Step int

	// TotalSteps is the total number of build steps (if known)
	TotalSteps int

	// Timestamp is when the event occurred
	Timestamp time.Time
}