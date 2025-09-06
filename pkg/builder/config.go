package builder

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

// ConfigFactory generates OCI image configurations.
// It handles platform settings, environment variables, labels, and container runtime configuration.
type ConfigFactory struct {
	platform v1.Platform
}

// NewConfigFactory creates a new configuration factory for the specified platform.
func NewConfigFactory(platform v1.Platform) *ConfigFactory {
	return &ConfigFactory{
		platform: platform,
	}
}

// GenerateConfig creates an OCI configuration for the image based on build options.
// It sets up the container environment, command, working directory, and metadata.
func (f *ConfigFactory) GenerateConfig(opts BuildOptions) (*v1.ConfigFile, error) {
	// Start with basic configuration
	config := &v1.ConfigFile{
		Architecture: opts.Platform.Architecture,
		OS:           opts.Platform.OS,
		Created:      v1.Time{Time: time.Now()},
		Config: v1.Config{
			Env:        opts.Env,
			Cmd:        opts.Cmd,
			WorkingDir: opts.WorkingDir,
			Entrypoint: opts.Entrypoint,
			User:       opts.User,
		},
		RootFS: v1.RootFS{
			Type: "layers",
		},
	}

	// Set platform variant if specified
	if opts.Platform.Variant != "" {
		config.Variant = opts.Platform.Variant
	}

	// Configure exposed ports if specified
	if len(opts.ExposedPorts) > 0 {
		config.Config.ExposedPorts = make(map[string]struct{})
		for port := range opts.ExposedPorts {
			config.Config.ExposedPorts[port] = struct{}{}
		}
	}

	// Configure volumes if specified
	if len(opts.Volumes) > 0 {
		config.Config.Volumes = make(map[string]struct{})
		for volume := range opts.Volumes {
			config.Config.Volumes[volume] = struct{}{}
		}
	}

	// Set labels
	if len(opts.Labels) > 0 {
		config.Config.Labels = make(map[string]string)
		for k, v := range opts.Labels {
			// Handle special label that should be set to build time
			if k == "org.opencontainers.image.created" && v == "" {
				config.Config.Labels[k] = time.Now().UTC().Format(time.RFC3339)
			} else {
				config.Config.Labels[k] = v
			}
		}
	}

	// Validate configuration
	if err := f.validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// ApplyConfig applies the configuration to an existing image.
// This updates the image's configuration while preserving layers.
func (f *ConfigFactory) ApplyConfig(img v1.Image, config *v1.ConfigFile) (v1.Image, error) {
	// Apply the configuration using mutate
	newImage, err := mutate.ConfigFile(img, config)
	if err != nil {
		return nil, fmt.Errorf("failed to apply configuration: %w", err)
	}

	return newImage, nil
}

// validateConfig validates the OCI configuration for compliance and best practices.
func (f *ConfigFactory) validateConfig(config *v1.ConfigFile) error {
	// Validate architecture
	if config.Architecture == "" {
		return fmt.Errorf("architecture cannot be empty")
	}

	// Validate OS
	if config.OS == "" {
		return fmt.Errorf("operating system cannot be empty")
	}

	// Validate working directory format (must be absolute if specified)
	if config.Config.WorkingDir != "" && !strings.HasPrefix(config.Config.WorkingDir, "/") {
		return fmt.Errorf("working directory must be an absolute path: %s", config.Config.WorkingDir)
	}

	// Basic port validation
	for port := range config.Config.ExposedPorts {
		if !strings.Contains(port, "/") {
			return fmt.Errorf("invalid exposed port format: %s", port)
		}
	}

	// Validate environment variables format
	for _, env := range config.Config.Env {
		if !strings.Contains(env, "=") {
			return fmt.Errorf("environment variable must be in KEY=value format: %s", env)
		}
	}

	// Basic user format validation
	if config.Config.User != "" && strings.Count(config.Config.User, ":") > 1 {
		return fmt.Errorf("invalid user format: %s", config.Config.User)
	}

	return nil
}

// DefaultLabels returns a set of recommended OCI labels.
// These provide useful metadata for image management.
func DefaultLabels(source string) map[string]string {
	return map[string]string{
		"org.opencontainers.image.created":     "", // Will be set to build time
		"org.opencontainers.image.source":      source,
		"org.opencontainers.image.title":       "Built with idpbuilder",
		"org.opencontainers.image.description": "OCI image built using go-containerregistry",
		"org.opencontainers.image.vendor":      "CNOE",
		"org.opencontainers.image.version":     "latest",
	}
}

// MergeConfigs merges multiple configuration options.
// Later configs override earlier ones for conflicting settings.
func MergeConfigs(configs ...*v1.ConfigFile) *v1.ConfigFile {
	if len(configs) == 0 {
		return &v1.ConfigFile{}
	}

	result := &v1.ConfigFile{}
	*result = *configs[0] // Start with first config

	for i := 1; i < len(configs); i++ {
		cfg := configs[i]

		// Override basic fields
		if cfg.Architecture != "" {
			result.Architecture = cfg.Architecture
		}
		if cfg.OS != "" {
			result.OS = cfg.OS
		}
		if cfg.Variant != "" {
			result.Variant = cfg.Variant
		}

		// Merge environment variables (later ones override)
		if len(cfg.Config.Env) > 0 {
			result.Config.Env = append(result.Config.Env, cfg.Config.Env...)
		}

		// Override command settings
		if len(cfg.Config.Cmd) > 0 {
			result.Config.Cmd = cfg.Config.Cmd
		}
		if len(cfg.Config.Entrypoint) > 0 {
			result.Config.Entrypoint = cfg.Config.Entrypoint
		}
		if cfg.Config.WorkingDir != "" {
			result.Config.WorkingDir = cfg.Config.WorkingDir
		}
		if cfg.Config.User != "" {
			result.Config.User = cfg.Config.User
		}

		// Merge labels
		if result.Config.Labels == nil {
			result.Config.Labels = make(map[string]string)
		}
		for k, v := range cfg.Config.Labels {
			result.Config.Labels[k] = v
		}

		// Merge exposed ports
		if result.Config.ExposedPorts == nil {
			result.Config.ExposedPorts = make(map[string]struct{})
		}
		for port := range cfg.Config.ExposedPorts {
			result.Config.ExposedPorts[port] = struct{}{}
		}

		// Merge volumes
		if result.Config.Volumes == nil {
			result.Config.Volumes = make(map[string]struct{})
		}
		for volume := range cfg.Config.Volumes {
			result.Config.Volumes[volume] = struct{}{}
		}
	}

	return result
}
