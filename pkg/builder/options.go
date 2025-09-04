// Package builder provides container image building capabilities using go-containerregistry
package builder

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// BuildOptions represents the configuration options for building container images
type BuildOptions struct {
	// Context is the build context directory
	Context string

	// Dockerfile specifies the path to the Dockerfile
	Dockerfile string

	// Tags specifies the repository tags to apply to the built image
	Tags []string

	// Platform specifies the target platform
	Platform string

	// Labels contains metadata labels to apply to the image
	Labels map[string]string

	// Logger specifies the logger to use for build output
	Logger Logger

	// Timeout specifies the maximum build duration
	Timeout time.Duration
}

// Logger interface for build logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// DefaultBuildOptions returns BuildOptions with sensible defaults
func DefaultBuildOptions() *BuildOptions {
	return &BuildOptions{
		Context:    ".",
		Dockerfile: "Dockerfile",
		Platform:   "linux/amd64",
		Labels:     make(map[string]string),
		Timeout:    30 * time.Minute,
	}
}

// Validate checks the BuildOptions for consistency and completeness
func (opts *BuildOptions) Validate() error {
	if opts.Context == "" {
		return fmt.Errorf("build context cannot be empty")
	}

	if !filepath.IsAbs(opts.Context) {
		abs, err := filepath.Abs(opts.Context)
		if err != nil {
			return fmt.Errorf("failed to resolve build context path: %w", err)
		}
		opts.Context = abs
	}

	if opts.Dockerfile == "" {
		opts.Dockerfile = "Dockerfile"
	}

	if !filepath.IsAbs(opts.Dockerfile) {
		opts.Dockerfile = filepath.Join(opts.Context, opts.Dockerfile)
	}

	if len(opts.Tags) == 0 {
		return fmt.Errorf("at least one tag must be specified")
	}

	// Basic tag validation
	for _, tag := range opts.Tags {
		if tag == "" || strings.Contains(tag, " ") {
			return fmt.Errorf("invalid tag format: %s", tag)
		}
	}

	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Minute
	}

	return nil
}

// WithContext sets the build context directory
func (opts *BuildOptions) WithContext(ctx string) *BuildOptions {
	opts.Context = ctx
	return opts
}

// WithDockerfile sets the Dockerfile path
func (opts *BuildOptions) WithDockerfile(dockerfile string) *BuildOptions {
	opts.Dockerfile = dockerfile
	return opts
}

// WithTags sets the image tags
func (opts *BuildOptions) WithTags(tags ...string) *BuildOptions {
	opts.Tags = tags
	return opts
}

// WithPlatform sets the target platform
func (opts *BuildOptions) WithPlatform(platform string) *BuildOptions {
	opts.Platform = platform
	return opts
}

// WithLabel adds a metadata label
func (opts *BuildOptions) WithLabel(key, value string) *BuildOptions {
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	opts.Labels[key] = value
	return opts
}

// WithLogger sets the logger
func (opts *BuildOptions) WithLogger(logger Logger) *BuildOptions {
	opts.Logger = logger
	return opts
}

// WithTimeout sets the build timeout
func (opts *BuildOptions) WithTimeout(timeout time.Duration) *BuildOptions {
	opts.Timeout = timeout
	return opts
}