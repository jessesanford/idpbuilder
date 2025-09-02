package builder

import "github.com/google/go-containerregistry/pkg/v1"

// BuildOptions configures the image build process.
// It provides platform specification, base image settings, and container configuration.
type BuildOptions struct {
	// Platform specifies the target platform (OS/architecture) for the image.
	// Defaults to linux/amd64 if not specified.
	Platform v1.Platform `json:"platform,omitempty"`
	
	// BaseImage specifies an optional base image reference.
	// If empty, builds from scratch (empty base image).
	BaseImage string `json:"baseImage,omitempty"`
	
	// Labels contains OCI labels to apply to the image.
	// Common labels include org.opencontainers.image.* metadata.
	Labels map[string]string `json:"labels,omitempty"`
	
	// Env contains environment variables to set in the container.
	// Format: ["KEY=value", "ANOTHER_KEY=another_value"]
	Env []string `json:"env,omitempty"`
	
	// WorkingDir sets the working directory inside the container.
	// Must be an absolute path if specified.
	WorkingDir string `json:"workingDir,omitempty"`
	
	// Entrypoint defines the container entrypoint.
	// When specified, overrides any base image entrypoint.
	Entrypoint []string `json:"entrypoint,omitempty"`
	
	// Cmd defines the default command and arguments.
	// Used when container is run without explicit command.
	Cmd []string `json:"cmd,omitempty"`
	
	// FeatureFlags enables/disables incomplete functionality per R307.
	// This ensures features under development don't break builds.
	FeatureFlags map[string]bool `json:"featureFlags,omitempty"`
	
	// User sets the user ID or name for container execution.
	// Can be numeric UID or username that exists in the image.
	User string `json:"user,omitempty"`
	
	// ExposedPorts declares ports that the container will listen on.
	// Format: map of "port/protocol" -> struct{} (e.g., "80/tcp")
	ExposedPorts map[string]struct{} `json:"exposedPorts,omitempty"`
	
	// Volumes declares mount points for the container.
	// Format: map of path -> struct{} (e.g., "/data")
	Volumes map[string]struct{} `json:"volumes,omitempty"`
}

// DefaultBuildOptions returns a BuildOptions struct with sensible defaults.
// This provides a starting point for common build scenarios.
func DefaultBuildOptions() BuildOptions {
	return BuildOptions{
		Platform: v1.Platform{
			OS:           "linux",
			Architecture: "amd64",
		},
		Labels: map[string]string{
			"org.opencontainers.image.created": "",  // Will be set to build time
			"org.opencontainers.image.source":  "idpbuilder",
		},
		FeatureFlags: make(map[string]bool),
	}
}

// WithPlatform sets the target platform for the build.
// Common platforms: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
func (opts BuildOptions) WithPlatform(os, arch string) BuildOptions {
	opts.Platform = v1.Platform{
		OS:           os,
		Architecture: arch,
	}
	return opts
}

// WithBaseImage sets the base image for the build.
// The base image will be pulled and used as the foundation.
func (opts BuildOptions) WithBaseImage(ref string) BuildOptions {
	opts.BaseImage = ref
	return opts
}

// WithLabels merges the provided labels with existing ones.
// Existing labels with the same key will be overwritten.
func (opts BuildOptions) WithLabels(labels map[string]string) BuildOptions {
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	for k, v := range labels {
		opts.Labels[k] = v
	}
	return opts
}

// WithWorkingDir sets the working directory for the container.
func (opts BuildOptions) WithWorkingDir(dir string) BuildOptions {
	opts.WorkingDir = dir
	return opts
}

// WithEntrypoint sets the container entrypoint command.
func (opts BuildOptions) WithEntrypoint(entrypoint ...string) BuildOptions {
	opts.Entrypoint = entrypoint
	return opts
}

// WithCmd sets the default container command.
func (opts BuildOptions) WithCmd(cmd ...string) BuildOptions {
	opts.Cmd = cmd
	return opts
}

// WithEnv adds environment variables to the container.
// Variables should be in "KEY=value" format.
func (opts BuildOptions) WithEnv(env ...string) BuildOptions {
	opts.Env = append(opts.Env, env...)
	return opts
}

// WithFeatureFlags enables specific feature flags.
// Used for R307 compliance to gate incomplete functionality.
func (opts BuildOptions) WithFeatureFlags(flags map[string]bool) BuildOptions {
	if opts.FeatureFlags == nil {
		opts.FeatureFlags = make(map[string]bool)
	}
	for k, v := range flags {
		opts.FeatureFlags[k] = v
	}
	return opts
}