package buildah

import (
	"fmt"
	"strings"
)

// BuildOptions configures container image build behavior
type BuildOptions struct {
	// Build arguments (--build-arg)
	BuildArgs map[string]string

	// Environment variables
	Env []string

	// Target platform (--platform)
	Platform string

	// Target architecture
	Arch string

	// Target OS
	OS string

	// Build-time network mode
	Network string

	// Additional labels
	Labels map[string]string

	// Cache options
	NoCache bool

	// Squash layers
	Squash bool
}

// NewBuildOptions creates default build options
func NewBuildOptions() *BuildOptions {
	return &BuildOptions{
		BuildArgs: make(map[string]string),
		Labels:    make(map[string]string),
		Env:       []string{},
	}
}

// WithBuildArg adds a build argument
func (opts *BuildOptions) WithBuildArg(key, value string) *BuildOptions {
	if opts.BuildArgs == nil {
		opts.BuildArgs = make(map[string]string)
	}
	opts.BuildArgs[key] = value
	return opts
}

// WithEnv adds an environment variable
func (opts *BuildOptions) WithEnv(envVar string) *BuildOptions {
	opts.Env = append(opts.Env, envVar)
	return opts
}

// WithPlatform sets the target platform
func (opts *BuildOptions) WithPlatform(platform string) *BuildOptions {
	opts.Platform = platform
	// Parse platform to set OS and Arch
	parts := strings.Split(platform, "/")
	if len(parts) == 2 {
		opts.OS = parts[0]
		opts.Arch = parts[1]
	}
	return opts
}

// WithLabel adds a label
func (opts *BuildOptions) WithLabel(key, value string) *BuildOptions {
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	opts.Labels[key] = value
	return opts
}

// Validate checks if options are valid
func (opts *BuildOptions) Validate() error {
	// Validate platform format
	if opts.Platform != "" {
		parts := strings.Split(opts.Platform, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid platform format: %s (expected os/arch)", opts.Platform)
		}
	}

	// Validate environment variables format
	for _, env := range opts.Env {
		if !strings.Contains(env, "=") {
			return fmt.Errorf("invalid environment variable format: %s (expected KEY=value)", env)
		}
	}

	return nil
}

// ToBuildahArgs converts options to buildah command arguments
func (opts *BuildOptions) ToBuildahArgs() []string {
	var args []string

	// Add build arguments
	for key, value := range opts.BuildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	// Add platform
	if opts.Platform != "" {
		args = append(args, "--platform", opts.Platform)
	}

	// Add labels
	for key, value := range opts.Labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", key, value))
	}

	// Add cache options
	if opts.NoCache {
		args = append(args, "--no-cache")
	}

	// Add squash option
	if opts.Squash {
		args = append(args, "--squash")
	}

	return args
}
