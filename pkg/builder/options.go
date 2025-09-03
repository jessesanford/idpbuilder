package builder

import (
	"fmt"
	"path/filepath"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// BuildOptions contains all configuration options for building an OCI image.
type BuildOptions struct {
	Platform     *v1.Platform          `json:"platform,omitempty"`
	Tags         []string              `json:"tags,omitempty"`
	Labels       map[string]string     `json:"labels,omitempty"`
	Environment  map[string]string     `json:"environment,omitempty"`
	WorkingDir   string                `json:"working_dir,omitempty"`
	Entrypoint   []string              `json:"entrypoint,omitempty"`
	Cmd          []string              `json:"cmd,omitempty"`
	User         string                `json:"user,omitempty"`
	ExposedPorts []string              `json:"exposed_ports,omitempty"`
	FeatureFlags map[string]bool       `json:"feature_flags,omitempty"`
	BuildArgs    map[string]string     `json:"build_args,omitempty"`
	ContextPath  string                `json:"context_path,omitempty"`
}

// NewBuildOptions creates a new BuildOptions with sensible defaults.
func NewBuildOptions() *BuildOptions {
	return &BuildOptions{
		Platform: &v1.Platform{
			Architecture: "amd64",
			OS:           "linux",
		},
		Labels:       make(map[string]string),
		Environment:  make(map[string]string),
		FeatureFlags: make(map[string]bool),
		BuildArgs:    make(map[string]string),
		Tags:         []string{},
		ExposedPorts: []string{},
	}
}

// Validate checks that the BuildOptions are valid and complete.
func (bo *BuildOptions) Validate() error {
	if bo.Platform == nil {
		return fmt.Errorf("platform is required")
	}

	if bo.Platform.OS == "" {
		return fmt.Errorf("platform OS is required")
	}

	if bo.Platform.Architecture == "" {
		return fmt.Errorf("platform architecture is required")
	}

	// Validate tags
	for _, tag := range bo.Tags {
		if err := validateTag(tag); err != nil {
			return fmt.Errorf("invalid tag %q: %w", tag, err)
		}
	}

	// Validate ports
	for _, port := range bo.ExposedPorts {
		if err := validatePort(port); err != nil {
			return fmt.Errorf("invalid exposed port %q: %w", port, err)
		}
	}

	// Validate context path
	if bo.ContextPath != "" {
		if !filepath.IsAbs(bo.ContextPath) {
			return fmt.Errorf("context path must be absolute: %q", bo.ContextPath)
		}
	}

	return nil
}

// SetFeatureFlag enables or disables a specific feature flag.
func (bo *BuildOptions) SetFeatureFlag(flag string, enabled bool) {
	if bo.FeatureFlags == nil {
		bo.FeatureFlags = make(map[string]bool)
	}
	bo.FeatureFlags[flag] = enabled
}

// IsFeatureEnabled checks if a feature flag is enabled.
func (bo *BuildOptions) IsFeatureEnabled(flag string) bool {
	if bo.FeatureFlags == nil {
		return false
	}
	return bo.FeatureFlags[flag]
}

// AddLabel adds a metadata label to the build options.
func (bo *BuildOptions) AddLabel(key, value string) {
	if bo.Labels == nil {
		bo.Labels = make(map[string]string)
	}
	bo.Labels[key] = value
}

// AddEnvironment adds an environment variable to the build options.
func (bo *BuildOptions) AddEnvironment(key, value string) {
	if bo.Environment == nil {
		bo.Environment = make(map[string]string)
	}
	bo.Environment[key] = value
}

// AddBuildArg adds a build argument to the build options.
func (bo *BuildOptions) AddBuildArg(key, value string) {
	if bo.BuildArgs == nil {
		bo.BuildArgs = make(map[string]string)
	}
	bo.BuildArgs[key] = value
}

// AddTag adds a tag to the build options.
func (bo *BuildOptions) AddTag(tag string) error {
	if err := validateTag(tag); err != nil {
		return fmt.Errorf("invalid tag %q: %w", tag, err)
	}
	if bo.Tags == nil {
		bo.Tags = make([]string, 0)
	}
	bo.Tags = append(bo.Tags, tag)
	return nil
}

// AddExposedPort adds an exposed port to the build options.
func (bo *BuildOptions) AddExposedPort(port string) error {
	if err := validatePort(port); err != nil {
		return fmt.Errorf("invalid port %q: %w", port, err)
	}
	if bo.ExposedPorts == nil {
		bo.ExposedPorts = make([]string, 0)
	}
	bo.ExposedPorts = append(bo.ExposedPorts, port)
	return nil
}

// SetPlatform sets the target platform for the build.
func (bo *BuildOptions) SetPlatform(os, arch string) {
	if bo.Platform == nil {
		bo.Platform = &v1.Platform{}
	}
	bo.Platform.OS = os
	bo.Platform.Architecture = arch
}

// Clone creates a deep copy of the BuildOptions.
func (bo *BuildOptions) Clone() *BuildOptions {
	clone := &BuildOptions{
		WorkingDir:  bo.WorkingDir,
		User:        bo.User,
		ContextPath: bo.ContextPath,
	}
	
	// Deep copy platform
	if bo.Platform != nil {
		clone.Platform = &v1.Platform{
			Architecture: bo.Platform.Architecture,
			OS:           bo.Platform.OS,
			Variant:      bo.Platform.Variant,
			OSFeatures:   bo.Platform.OSFeatures,
			OSVersion:    bo.Platform.OSVersion,
		}
	}
	
	// Deep copy maps
	if bo.Labels != nil {
		clone.Labels = make(map[string]string)
		for k, v := range bo.Labels {
			clone.Labels[k] = v
		}
	}
	if bo.Environment != nil {
		clone.Environment = make(map[string]string)
		for k, v := range bo.Environment {
			clone.Environment[k] = v
		}
	}
	if bo.BuildArgs != nil {
		clone.BuildArgs = make(map[string]string)
		for k, v := range bo.BuildArgs {
			clone.BuildArgs[k] = v
		}
	}
	if bo.FeatureFlags != nil {
		clone.FeatureFlags = make(map[string]bool)
		for k, v := range bo.FeatureFlags {
			clone.FeatureFlags[k] = v
		}
	}
	
	// Deep copy slices
	if bo.Tags != nil {
		clone.Tags = make([]string, len(bo.Tags))
		copy(clone.Tags, bo.Tags)
	}
	if bo.Entrypoint != nil {
		clone.Entrypoint = make([]string, len(bo.Entrypoint))
		copy(clone.Entrypoint, bo.Entrypoint)
	}
	if bo.Cmd != nil {
		clone.Cmd = make([]string, len(bo.Cmd))
		copy(clone.Cmd, bo.Cmd)
	}
	if bo.ExposedPorts != nil {
		clone.ExposedPorts = make([]string, len(bo.ExposedPorts))
		copy(clone.ExposedPorts, bo.ExposedPorts)
	}
	
	return clone
}

// validateTag validates that a tag follows Docker tag naming conventions.
func validateTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if strings.Contains(tag, " ") {
		return fmt.Errorf("tag cannot contain spaces")
	}
	if strings.HasPrefix(tag, "-") {
		return fmt.Errorf("tag cannot start with dash")
	}
	return nil
}

// validatePort validates that a port specification is valid.
func validatePort(port string) error {
	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if !strings.Contains(port, "/") {
		return nil // Assume tcp if protocol not specified
	}
	parts := strings.Split(port, "/")
	if len(parts) != 2 {
		return fmt.Errorf("port must be in format 'number/protocol'")
	}
	protocol := parts[1]
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("protocol must be tcp or udp, got: %s", protocol)
	}
	return nil
}

// GetAllLabels returns a copy of all labels.
func (bo *BuildOptions) GetAllLabels() map[string]string {
	if bo.Labels == nil {
		return make(map[string]string)
	}
	result := make(map[string]string)
	for k, v := range bo.Labels {
		result[k] = v
	}
	return result
}

// GetAllEnvironment returns a copy of all environment variables.
func (bo *BuildOptions) GetAllEnvironment() map[string]string {
	if bo.Environment == nil {
		return make(map[string]string)
	}
	result := make(map[string]string)
	for k, v := range bo.Environment {
		result[k] = v
	}
	return result
}