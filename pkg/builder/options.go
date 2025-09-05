package builder

import "github.com/google/go-containerregistry/pkg/v1"

// BuildOptions configures the image build process.
type BuildOptions struct {
	Platform     v1.Platform         `json:"platform,omitempty"`
	BaseImage    string              `json:"baseImage,omitempty"`
	Labels       map[string]string   `json:"labels,omitempty"`
	Env          []string            `json:"env,omitempty"`
	WorkingDir   string              `json:"workingDir,omitempty"`
	Entrypoint   []string            `json:"entrypoint,omitempty"`
	Cmd          []string            `json:"cmd,omitempty"`
	FeatureFlags map[string]bool     `json:"featureFlags,omitempty"`
	User         string              `json:"user,omitempty"`
	ExposedPorts map[string]struct{} `json:"exposedPorts,omitempty"`
	Volumes      map[string]struct{} `json:"volumes,omitempty"`
}

// DefaultBuildOptions returns BuildOptions with defaults.
func DefaultBuildOptions() BuildOptions {
	return BuildOptions{
		Platform:     v1.Platform{OS: "linux", Architecture: "amd64"},
		Labels:       map[string]string{"org.opencontainers.image.source": "idpbuilder"},
		FeatureFlags: make(map[string]bool),
	}
}

// WithPlatform sets the target platform.
func (opts BuildOptions) WithPlatform(os, arch string) BuildOptions {
	opts.Platform = v1.Platform{OS: os, Architecture: arch}
	return opts
}

// WithBaseImage sets the base image.
func (opts BuildOptions) WithBaseImage(ref string) BuildOptions {
	opts.BaseImage = ref
	return opts
}

// WithLabels merges labels.
func (opts BuildOptions) WithLabels(labels map[string]string) BuildOptions {
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	for k, v := range labels {
		opts.Labels[k] = v
	}
	return opts
}

// WithWorkingDir sets the working directory.
func (opts BuildOptions) WithWorkingDir(dir string) BuildOptions {
	opts.WorkingDir = dir
	return opts
}

// WithEntrypoint sets the container entrypoint.
func (opts BuildOptions) WithEntrypoint(entrypoint ...string) BuildOptions {
	opts.Entrypoint = entrypoint
	return opts
}

// WithCmd sets the default command.
func (opts BuildOptions) WithCmd(cmd ...string) BuildOptions {
	opts.Cmd = cmd
	return opts
}

// WithEnv adds environment variables.
func (opts BuildOptions) WithEnv(env ...string) BuildOptions {
	opts.Env = append(opts.Env, env...)
	return opts
}

// WithFeatureFlags enables feature flags.
func (opts BuildOptions) WithFeatureFlags(flags map[string]bool) BuildOptions {
	if opts.FeatureFlags == nil {
		opts.FeatureFlags = make(map[string]bool)
	}
	for k, v := range flags {
		opts.FeatureFlags[k] = v
	}
	return opts
}
