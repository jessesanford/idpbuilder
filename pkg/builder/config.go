package builder

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// ConfigFactory handles creation and configuration of OCI image configurations.
type ConfigFactory struct {
	BaseConfig         *v1.Config
	DefaultLabels      map[string]string
	DefaultEnvironment map[string]string
}

// NewConfigFactory creates a new ConfigFactory with sensible defaults.
func NewConfigFactory() *ConfigFactory {
	return &ConfigFactory{
		BaseConfig: &v1.Config{
			User:       "root",
			WorkingDir: "/",
			Env: []string{
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			},
		},
		DefaultLabels: map[string]string{
			"org.opencontainers.image.created":     time.Now().Format(time.RFC3339),
			"org.opencontainers.image.title":       "OCI Image",
			"org.opencontainers.image.description": "Built with go-containerregistry",
		},
		DefaultEnvironment: make(map[string]string),
	}
}

// CreateConfig generates an OCI image configuration based on the provided BuildOptions.
func (cf *ConfigFactory) CreateConfig(opts *BuildOptions) (*v1.Config, error) {
	if opts == nil {
		return nil, fmt.Errorf("build options cannot be nil")
	}

	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid build options: %w", err)
	}

	config := &v1.Config{}
	if cf.BaseConfig != nil {
		*config = *cf.BaseConfig
		if cf.BaseConfig.Env != nil {
			config.Env = append([]string(nil), cf.BaseConfig.Env...)
		}
		if cf.BaseConfig.Entrypoint != nil {
			config.Entrypoint = append([]string(nil), cf.BaseConfig.Entrypoint...)
		}
		if cf.BaseConfig.Cmd != nil {
			config.Cmd = append([]string(nil), cf.BaseConfig.Cmd...)
		}
	}

	// Apply build options
	if opts.WorkingDir != "" {
		config.WorkingDir = opts.WorkingDir
	}

	if opts.User != "" {
		config.User = opts.User
	}

	if len(opts.Entrypoint) > 0 {
		config.Entrypoint = append([]string(nil), opts.Entrypoint...)
	}

	if len(opts.Cmd) > 0 {
		config.Cmd = append([]string(nil), opts.Cmd...)
	}

	// Merge labels
	config.Labels = cf.mergeLabels(opts.Labels)

	// Merge environment variables
	config.Env = cf.mergeEnvironment(config.Env, opts.Environment)

	// Handle exposed ports
	if len(opts.ExposedPorts) > 0 {
		if config.ExposedPorts == nil {
			config.ExposedPorts = make(map[string]struct{})
		}
		for _, port := range opts.ExposedPorts {
			config.ExposedPorts[port] = struct{}{}
		}
	}

	return config, nil
}

// CreatePlatformConfig creates a platform-specific configuration.
func (cf *ConfigFactory) CreatePlatformConfig(opts *BuildOptions) (*v1.Platform, error) {
	if opts == nil || opts.Platform == nil {
		return &v1.Platform{
			Architecture: "amd64",
			OS:           "linux",
		}, nil
	}

	platform := &v1.Platform{
		Architecture: opts.Platform.Architecture,
		OS:           opts.Platform.OS,
		Variant:      opts.Platform.Variant,
		OSVersion:    opts.Platform.OSVersion,
		OSFeatures:   append([]string(nil), opts.Platform.OSFeatures...),
	}

	// Validate platform
	if err := validatePlatform(platform); err != nil {
		return nil, fmt.Errorf("invalid platform: %w", err)
	}

	return platform, nil
}

// AddDefaultLabel adds a label that will be applied to all configurations.
func (cf *ConfigFactory) AddDefaultLabel(key, value string) {
	if cf.DefaultLabels == nil {
		cf.DefaultLabels = make(map[string]string)
	}
	cf.DefaultLabels[key] = value
}

// AddDefaultEnvironment adds an environment variable that will be applied to all configurations.
func (cf *ConfigFactory) AddDefaultEnvironment(key, value string) {
	if cf.DefaultEnvironment == nil {
		cf.DefaultEnvironment = make(map[string]string)
	}
	cf.DefaultEnvironment[key] = value
}

// SetBaseConfig sets the base configuration that will be used as a template.
func (cf *ConfigFactory) SetBaseConfig(config *v1.Config) {
	cf.BaseConfig = config
}

// mergeLabels combines default labels with build option labels.
func (cf *ConfigFactory) mergeLabels(optLabels map[string]string) map[string]string {
	result := make(map[string]string, len(cf.DefaultLabels)+len(optLabels))
	// Add default labels first
	for k, v := range cf.DefaultLabels {
		result[k] = v
	}
	// Override with option labels
	for k, v := range optLabels {
		result[k] = v
	}
	return result
}

// mergeEnvironment combines existing environment with build option environment.
func (cf *ConfigFactory) mergeEnvironment(existingEnv []string, optEnv map[string]string) []string {
	envMap := make(map[string]string)
	
	for _, env := range existingEnv {
		if parts := strings.SplitN(env, "=", 2); len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	for k, v := range cf.DefaultEnvironment {
		envMap[k] = v
	}
	for k, v := range optEnv {
		envMap[k] = v
	}

	result := make([]string, 0, len(envMap))
	for k, v := range envMap {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

// validatePlatform checks that a platform specification is valid.
func validatePlatform(platform *v1.Platform) error {
	if platform.OS == "" {
		return fmt.Errorf("platform OS is required")
	}

	if platform.Architecture == "" {
		return fmt.Errorf("platform architecture is required")
	}

	// Validate supported OS values
	supportedOS := map[string]bool{
		"linux":   true,
		"windows": true,
		"darwin":  true,
	}

	if !supportedOS[platform.OS] {
		return fmt.Errorf("unsupported OS: %s", platform.OS)
	}

	// Validate supported architecture values
	supportedArch := map[string]bool{
		"amd64": true,
		"arm64": true,
		"arm":   true,
		"386":   true,
	}

	if !supportedArch[platform.Architecture] {
		return fmt.Errorf("unsupported architecture: %s", platform.Architecture)
	}

	return nil
}

// ParsePlatform parses a platform string in the format "os/arch[/variant]".
func ParsePlatform(platformStr string) (*v1.Platform, error) {
	if platformStr == "" {
		return &v1.Platform{
			OS:           "linux",
			Architecture: "amd64",
		}, nil
	}

	parts := strings.Split(platformStr, "/")
	if len(parts) < 2 || len(parts) > 3 {
		return nil, fmt.Errorf("platform must be in format 'os/arch[/variant]', got: %s", platformStr)
	}

	platform := &v1.Platform{
		OS:           parts[0],
		Architecture: parts[1],
	}

	if len(parts) == 3 {
		platform.Variant = parts[2]
	}

	if err := validatePlatform(platform); err != nil {
		return nil, err
	}

	return platform, nil
}

// FormatPlatform formats a platform as a string in the format "os/arch[/variant]".
func FormatPlatform(platform *v1.Platform) string {
	if platform == nil {
		return "linux/amd64"
	}

	result := fmt.Sprintf("%s/%s", platform.OS, platform.Architecture)
	if platform.Variant != "" {
		result += "/" + platform.Variant
	}

	return result
}

// GenerateConfigDigest creates a simple digest for a configuration.
func GenerateConfigDigest(config *v1.Config) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config cannot be nil")
	}
	
	configStr := fmt.Sprintf("%s:%s:%v:%v", config.User, config.WorkingDir, 
		config.Entrypoint, config.Cmd)
	hash := 0
	for _, c := range configStr {
		hash = hash*31 + int(c)
	}
	return "sha256:" + strconv.FormatInt(int64(hash), 16), nil
}