package builder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// BuildConfig represents the core configuration for building container images.
// It contains all the necessary information to build an image from a given context.
type BuildConfig struct {
	// ContextPath is the path to the build context directory
	ContextPath string `json:"context_path" yaml:"context_path"`

	// Dockerfile is the path to the Dockerfile, relative to ContextPath
	Dockerfile string `json:"dockerfile" yaml:"dockerfile"`

	// Tags is a list of tags to apply to the built image
	Tags []string `json:"tags" yaml:"tags"`

	// Platform specifies the target platform for the image
	Platform PlatformConfig `json:"platform" yaml:"platform"`

	// Registry contains registry-specific configuration
	Registry RegistryConfig `json:"registry" yaml:"registry"`

	// BuildArgs contains build-time arguments
	BuildArgs map[string]string `json:"build_args" yaml:"build_args"`

	// Labels contains metadata labels to apply to the image
	Labels map[string]string `json:"labels" yaml:"labels"`

	// Target specifies the target stage in a multi-stage Dockerfile
	Target string `json:"target,omitempty" yaml:"target,omitempty"`

	// NoCache disables layer caching during build
	NoCache bool `json:"no_cache" yaml:"no_cache"`

	// Pull forces pulling of base images even if present locally
	Pull bool `json:"pull" yaml:"pull"`

	// Remove removes intermediate containers after successful build
	Remove bool `json:"remove" yaml:"remove"`

	// Squash squashes all layers into a single layer (experimental)
	Squash bool `json:"squash" yaml:"squash"`

	// BuildTimeout specifies the maximum time allowed for the build
	BuildTimeout time.Duration `json:"build_timeout" yaml:"build_timeout"`

	// Memory limit for build containers (in bytes)
	MemoryLimit int64 `json:"memory_limit" yaml:"memory_limit"`

	// CPULimit specifies CPU limit for build containers
	CPULimit float64 `json:"cpu_limit" yaml:"cpu_limit"`
}

// PlatformConfig specifies target platform information for multi-platform builds.
type PlatformConfig struct {
	// OS specifies the target operating system (e.g., "linux", "windows")
	OS string `json:"os" yaml:"os"`

	// Architecture specifies the target architecture (e.g., "amd64", "arm64")
	Architecture string `json:"architecture" yaml:"architecture"`

	// Variant specifies the architecture variant (e.g., "v7" for arm)
	Variant string `json:"variant,omitempty" yaml:"variant,omitempty"`

	// OSVersion specifies the OS version for Windows images
	OSVersion string `json:"os_version,omitempty" yaml:"os_version,omitempty"`

	// OSFeatures specifies OS features required by the image
	OSFeatures []string `json:"os_features,omitempty" yaml:"os_features,omitempty"`
}

// RegistryConfig contains configuration for registry authentication and access.
type RegistryConfig struct {
	// Hostname is the registry hostname (e.g., "registry.example.com")
	Hostname string `json:"hostname" yaml:"hostname"`

	// Username for registry authentication
	Username string `json:"username,omitempty" yaml:"username,omitempty"`

	// Password for registry authentication
	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	// Token for token-based authentication
	Token string `json:"token,omitempty" yaml:"token,omitempty"`

	// RegistryToken for registry-specific token authentication
	RegistryToken string `json:"registry_token,omitempty" yaml:"registry_token,omitempty"`

	// Insecure allows insecure registry connections
	Insecure bool `json:"insecure" yaml:"insecure"`

	// PlainHTTP forces the use of HTTP instead of HTTPS
	PlainHTTP bool `json:"plain_http" yaml:"plain_http"`

	// SkipTLSVerify skips TLS certificate verification
	SkipTLSVerify bool `json:"skip_tls_verify" yaml:"skip_tls_verify"`

	// Timeout for registry operations
	Timeout time.Duration `json:"timeout" yaml:"timeout"`

	// RetryCount specifies number of retries for failed operations
	RetryCount int `json:"retry_count" yaml:"retry_count"`

	// RetryDelay specifies delay between retries
	RetryDelay time.Duration `json:"retry_delay" yaml:"retry_delay"`
}

// ImageOptions contains options for image creation and metadata.
type ImageOptions struct {
	// Created timestamp for the image
	Created time.Time `json:"created" yaml:"created"`

	// Author of the image
	Author string `json:"author,omitempty" yaml:"author,omitempty"`

	// Architecture of the image (deprecated, use Platform instead)
	Architecture string `json:"architecture,omitempty" yaml:"architecture,omitempty"`

	// OS of the image (deprecated, use Platform instead)
	OS string `json:"os,omitempty" yaml:"os,omitempty"`

	// Config contains the image configuration
	Config ImageConfig `json:"config" yaml:"config"`

	// RootFS describes the image's root filesystem
	RootFS RootFS `json:"rootfs" yaml:"rootfs"`

	// History contains the history of the image layers
	History []History `json:"history,omitempty" yaml:"history,omitempty"`
}

// ImageConfig contains configuration for the container that will run the image.
type ImageConfig struct {
	// User specifies the user that containers should run as
	User string `json:"user,omitempty" yaml:"user,omitempty"`

	// ExposedPorts lists the ports exposed by the container
	ExposedPorts map[string]struct{} `json:"exposed_ports,omitempty" yaml:"exposed_ports,omitempty"`

	// Env contains environment variables for the container
	Env []string `json:"env,omitempty" yaml:"env,omitempty"`

	// Entrypoint specifies the entry point for the container
	Entrypoint []string `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`

	// Cmd specifies the default command for the container
	Cmd []string `json:"cmd,omitempty" yaml:"cmd,omitempty"`

	// Volumes contains volumes used by the container
	Volumes map[string]struct{} `json:"volumes,omitempty" yaml:"volumes,omitempty"`

	// WorkingDir specifies the working directory for the container
	WorkingDir string `json:"working_dir,omitempty" yaml:"working_dir,omitempty"`

	// Labels contains metadata labels for the container
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// StopSignal specifies the signal to stop the container
	StopSignal string `json:"stop_signal,omitempty" yaml:"stop_signal,omitempty"`

	// StopTimeout specifies timeout for stopping the container
	StopTimeout *int `json:"stop_timeout,omitempty" yaml:"stop_timeout,omitempty"`

	// Shell specifies the shell to use for shell-form commands
	Shell []string `json:"shell,omitempty" yaml:"shell,omitempty"`
}

// RootFS describes the root filesystem of the image.
type RootFS struct {
	// Type is typically "layers" for layered filesystems
	Type string `json:"type" yaml:"type"`

	// DiffIDs contains the diff IDs of the filesystem layers
	DiffIDs []string `json:"diff_ids" yaml:"diff_ids"`
}

// History describes a single layer in the image history.
type History struct {
	// Created timestamp for when the layer was created
	Created time.Time `json:"created" yaml:"created"`

	// CreatedBy contains the command that created the layer
	CreatedBy string `json:"created_by,omitempty" yaml:"created_by,omitempty"`

	// Author of the layer
	Author string `json:"author,omitempty" yaml:"author,omitempty"`

	// Comment contains any comment about the layer
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`

	// EmptyLayer indicates if this history entry corresponds to an empty layer
	EmptyLayer bool `json:"empty_layer,omitempty" yaml:"empty_layer,omitempty"`
}

// DefaultBuildConfig returns a BuildConfig with sensible defaults.
func DefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		ContextPath:  ".",
		Dockerfile:   "Dockerfile",
		Tags:         []string{"latest"},
		Platform:     DefaultPlatformConfig(),
		Registry:     DefaultRegistryConfig(),
		BuildArgs:    make(map[string]string),
		Labels:       make(map[string]string),
		Remove:       true,
		BuildTimeout: 30 * time.Minute,
		MemoryLimit:  0, // No limit
		CPULimit:     0, // No limit
	}
}

// DefaultPlatformConfig returns a PlatformConfig with defaults for the current platform.
func DefaultPlatformConfig() PlatformConfig {
	return PlatformConfig{
		OS:           "linux",
		Architecture: "amd64",
	}
}

// DefaultRegistryConfig returns a RegistryConfig with sensible defaults.
func DefaultRegistryConfig() RegistryConfig {
	return RegistryConfig{
		Hostname:    "index.docker.io",
		Insecure:    false,
		PlainHTTP:   false,
		Timeout:     30 * time.Second,
		RetryCount:  3,
		RetryDelay:  1 * time.Second,
	}
}

// Validate validates the BuildConfig and returns an error if invalid.
func (c *BuildConfig) Validate() error {
	if c.ContextPath == "" {
		return fmt.Errorf("context path cannot be empty")
	}

	if c.Dockerfile == "" {
		return fmt.Errorf("dockerfile cannot be empty")
	}

	if len(c.Tags) == 0 {
		return fmt.Errorf("at least one tag must be specified")
	}

	// Validate tags format
	for _, tag := range c.Tags {
		if strings.TrimSpace(tag) == "" {
			return fmt.Errorf("tag cannot be empty or whitespace")
		}
		if strings.Contains(tag, " ") {
			return fmt.Errorf("tag '%s' cannot contain spaces", tag)
		}
	}

	// Validate platform
	if err := c.Platform.Validate(); err != nil {
		return fmt.Errorf("platform validation failed: %w", err)
	}

	// Validate registry if hostname is set
	if c.Registry.Hostname != "" {
		if err := c.Registry.Validate(); err != nil {
			return fmt.Errorf("registry validation failed: %w", err)
		}
	}

	// Validate build timeout
	if c.BuildTimeout <= 0 {
		return fmt.Errorf("build timeout must be greater than zero")
	}

	// Validate memory limit
	if c.MemoryLimit < 0 {
		return fmt.Errorf("memory limit cannot be negative")
	}

	// Validate CPU limit
	if c.CPULimit < 0 {
		return fmt.Errorf("CPU limit cannot be negative")
	}

	return nil
}

// Validate validates the PlatformConfig.
func (p *PlatformConfig) Validate() error {
	if p.OS == "" {
		return fmt.Errorf("OS cannot be empty")
	}

	if p.Architecture == "" {
		return fmt.Errorf("architecture cannot be empty")
	}

	// Validate OS values
	validOS := []string{"linux", "windows", "darwin", "freebsd", "openbsd", "netbsd"}
	validOSMap := make(map[string]bool)
	for _, os := range validOS {
		validOSMap[os] = true
	}
	if !validOSMap[p.OS] {
		return fmt.Errorf("unsupported OS: %s", p.OS)
	}

	// Validate architecture values
	validArch := []string{"amd64", "arm64", "arm", "386", "ppc64le", "s390x", "mips64le", "mips64", "riscv64"}
	validArchMap := make(map[string]bool)
	for _, arch := range validArch {
		validArchMap[arch] = true
	}
	if !validArchMap[p.Architecture] {
		return fmt.Errorf("unsupported architecture: %s", p.Architecture)
	}

	return nil
}

// Validate validates the RegistryConfig.
func (r *RegistryConfig) Validate() error {
	if r.Hostname == "" {
		return fmt.Errorf("registry hostname cannot be empty")
	}

	// Validate authentication - either username/password or token, but not both
	hasUserPass := r.Username != "" || r.Password != ""
	hasToken := r.Token != "" || r.RegistryToken != ""

	if hasUserPass && hasToken {
		return fmt.Errorf("cannot specify both username/password and token authentication")
	}

	// If username is provided, password should also be provided (and vice versa)
	if (r.Username != "") != (r.Password != "") {
		return fmt.Errorf("username and password must both be provided or both be empty")
	}

	// Validate timeout
	if r.Timeout <= 0 {
		return fmt.Errorf("registry timeout must be greater than zero")
	}

	// Validate retry settings
	if r.RetryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}

	if r.RetryDelay < 0 {
		return fmt.Errorf("retry delay cannot be negative")
	}

	return nil
}

// String returns the platform string in the format OS/Architecture[/Variant].
func (p *PlatformConfig) String() string {
	platform := fmt.Sprintf("%s/%s", p.OS, p.Architecture)
	if p.Variant != "" {
		platform += "/" + p.Variant
	}
	return platform
}

// ContextDir returns the absolute path to the build context directory.
func (c *BuildConfig) ContextDir() (string, error) {
	return filepath.Abs(c.ContextPath)
}

// DockerfilePath returns the absolute path to the Dockerfile.
func (c *BuildConfig) DockerfilePath() (string, error) {
	contextDir, err := c.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(contextDir, c.Dockerfile), nil
}

// Clone creates a deep copy of the BuildConfig.
func (c *BuildConfig) Clone() *BuildConfig {
	clone := *c

	// Deep copy slices and maps
	clone.Tags = make([]string, len(c.Tags))
	copy(clone.Tags, c.Tags)

	clone.BuildArgs = make(map[string]string)
	for k, v := range c.BuildArgs {
		clone.BuildArgs[k] = v
	}

	clone.Labels = make(map[string]string)
	for k, v := range c.Labels {
		clone.Labels[k] = v
	}

	return &clone
}

// WithContext returns a new BuildConfig with the given context path.
func (c *BuildConfig) WithContext(ctx context.Context) *BuildConfig {
	// This method can be extended to support context-based cancellation
	// For now, it just returns a clone
	return c.Clone()
}