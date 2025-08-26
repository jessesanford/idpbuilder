// Package build provides extended configuration and integration features for buildah operations.
// This is split 003 of the buildah-integration implementation, focusing on advanced
// configuration management, inheritance, and integration helpers.
package build

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/containers/buildah"
	"github.com/containers/buildah/define"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
)

// ExtendedBuildConfig represents an advanced configuration for buildah operations
// with support for inheritance, profiles, and environment-specific overrides.
type ExtendedBuildConfig struct {
	// Base configuration
	Name        string            `json:"name" yaml:"name"`
	Version     string            `json:"version" yaml:"version"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Profile     string            `json:"profile,omitempty" yaml:"profile,omitempty"`
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`

	// Build configuration
	Build BuildOptions `json:"build" yaml:"build"`

	// Runtime configuration
	Runtime RuntimeOptions `json:"runtime,omitempty" yaml:"runtime,omitempty"`

	// Storage configuration
	Storage StorageOptions `json:"storage,omitempty" yaml:"storage,omitempty"`

	// Network configuration
	Network NetworkOptions `json:"network,omitempty" yaml:"network,omitempty"`

	// Security configuration
	Security SecurityOptions `json:"security,omitempty" yaml:"security,omitempty"`

	// Cache configuration
	Cache CacheOptions `json:"cache,omitempty" yaml:"cache,omitempty"`

	// Feature flags
	Features FeatureFlags `json:"features,omitempty" yaml:"features,omitempty"`

	// Inheritance chain
	Inherits []string `json:"inherits,omitempty" yaml:"inherits,omitempty"`

	// Profiles to apply
	Profiles map[string]Profile `json:"profiles,omitempty" yaml:"profiles,omitempty"`
}

// BuildOptions defines build-specific configuration
type BuildOptions struct {
	Dockerfile   string            `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty"`
	Context      string            `json:"context,omitempty" yaml:"context,omitempty"`
	Args         map[string]string `json:"args,omitempty" yaml:"args,omitempty"`
	Labels       map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Target       string            `json:"target,omitempty" yaml:"target,omitempty"`
	Platform     string            `json:"platform,omitempty" yaml:"platform,omitempty"`
	Output       []string          `json:"output,omitempty" yaml:"output,omitempty"`
	Pull         bool              `json:"pull,omitempty" yaml:"pull,omitempty"`
	NoCache      bool              `json:"no_cache,omitempty" yaml:"no_cache,omitempty"`
	ForceRm      bool              `json:"force_rm,omitempty" yaml:"force_rm,omitempty"`
	RemoveImages bool              `json:"remove_images,omitempty" yaml:"remove_images,omitempty"`
}

// RuntimeOptions defines runtime configuration
type RuntimeOptions struct {
	Runtime     string            `json:"runtime,omitempty" yaml:"runtime,omitempty"`
	Memory      string            `json:"memory,omitempty" yaml:"memory,omitempty"`
	CPUs        string            `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	CPUShares   uint64            `json:"cpu_shares,omitempty" yaml:"cpu_shares,omitempty"`
	Ulimits     []string          `json:"ulimits,omitempty" yaml:"ulimits,omitempty"`
	ShmSize     string            `json:"shm_size,omitempty" yaml:"shm_size,omitempty"`
	User        string            `json:"user,omitempty" yaml:"user,omitempty"`
	WorkingDir  string            `json:"working_dir,omitempty" yaml:"working_dir,omitempty"`
	Env         map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Volumes     []string          `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Devices     []string          `json:"devices,omitempty" yaml:"devices,omitempty"`
}

// StorageOptions defines storage configuration
type StorageOptions struct {
	Driver      string            `json:"driver,omitempty" yaml:"driver,omitempty"`
	Root        string            `json:"root,omitempty" yaml:"root,omitempty"`
	RunRoot     string            `json:"run_root,omitempty" yaml:"run_root,omitempty"`
	GraphRoot   string            `json:"graph_root,omitempty" yaml:"graph_root,omitempty"`
	Options     map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
	ImageStore  string            `json:"image_store,omitempty" yaml:"image_store,omitempty"`
}

// NetworkOptions defines network configuration
type NetworkOptions struct {
	Mode      string            `json:"mode,omitempty" yaml:"mode,omitempty"`
	Networks  []string          `json:"networks,omitempty" yaml:"networks,omitempty"`
	Ports     []string          `json:"ports,omitempty" yaml:"ports,omitempty"`
	DNS       []string          `json:"dns,omitempty" yaml:"dns,omitempty"`
	DNSSearch []string          `json:"dns_search,omitempty" yaml:"dns_search,omitempty"`
	ExtraHosts []string         `json:"extra_hosts,omitempty" yaml:"extra_hosts,omitempty"`
	Options   map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
}

// SecurityOptions defines security configuration
type SecurityOptions struct {
	NoNewPrivileges bool              `json:"no_new_privileges,omitempty" yaml:"no_new_privileges,omitempty"`
	Capabilities    CapabilityOptions `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	SELinux         SELinuxOptions    `json:"selinux,omitempty" yaml:"selinux,omitempty"`
	AppArmor        AppArmorOptions   `json:"apparmor,omitempty" yaml:"apparmor,omitempty"`
	Seccomp         SeccompOptions    `json:"seccomp,omitempty" yaml:"seccomp,omitempty"`
	Privileged      bool              `json:"privileged,omitempty" yaml:"privileged,omitempty"`
}

// CapabilityOptions defines capability configuration
type CapabilityOptions struct {
	Add  []string `json:"add,omitempty" yaml:"add,omitempty"`
	Drop []string `json:"drop,omitempty" yaml:"drop,omitempty"`
}

// SELinuxOptions defines SELinux configuration
type SELinuxOptions struct {
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	Level   string `json:"level,omitempty" yaml:"level,omitempty"`
}

// AppArmorOptions defines AppArmor configuration
type AppArmorOptions struct {
	Profile string `json:"profile,omitempty" yaml:"profile,omitempty"`
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

// SeccompOptions defines Seccomp configuration
type SeccompOptions struct {
	Profile string `json:"profile,omitempty" yaml:"profile,omitempty"`
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

// CacheOptions defines cache configuration
type CacheOptions struct {
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Dir     string `json:"dir,omitempty" yaml:"dir,omitempty"`
	MaxSize string `json:"max_size,omitempty" yaml:"max_size,omitempty"`
	TTL     string `json:"ttl,omitempty" yaml:"ttl,omitempty"`
}

// FeatureFlags defines feature flags
type FeatureFlags struct {
	ParallelBuilds bool `json:"parallel_builds,omitempty" yaml:"parallel_builds,omitempty"`
	LayerCaching   bool `json:"layer_caching,omitempty" yaml:"layer_caching,omitempty"`
	BuildKit       bool `json:"buildkit,omitempty" yaml:"buildkit,omitempty"`
	Squash         bool `json:"squash,omitempty" yaml:"squash,omitempty"`
}

// Profile represents a configuration profile
type Profile struct {
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Settings    map[string]interface{} `json:"settings" yaml:"settings"`
}

// ExtendedConfigManager manages extended build configurations
type ExtendedConfigManager struct {
	store  storage.Store
	logger *Logger
	cache  map[string]*ExtendedBuildConfig
	cacheTTL time.Duration
}

// Logger interface for logging operations
type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// NewExtendedConfigManager creates a new extended configuration manager
func NewExtendedConfigManager(store storage.Store, logger Logger) *ExtendedConfigManager {
	return &ExtendedConfigManager{
		store:    store,
		logger:   logger,
		cache:    make(map[string]*ExtendedBuildConfig),
		cacheTTL: 5 * time.Minute,
	}
}

// LoadConfig loads a configuration from file with inheritance resolution
func (m *ExtendedConfigManager) LoadConfig(configPath string) (*ExtendedBuildConfig, error) {
	// Check cache first
	if cached, exists := m.cache[configPath]; exists {
		m.logger.Debugf("Using cached configuration: %s", configPath)
		return cached, nil
	}

	// Load from file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config ExtendedBuildConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Resolve inheritance
	if err := m.resolveInheritance(&config, filepath.Dir(configPath)); err != nil {
		return nil, fmt.Errorf("failed to resolve inheritance: %w", err)
	}

	// Apply profiles
	if err := m.applyProfiles(&config); err != nil {
		return nil, fmt.Errorf("failed to apply profiles: %w", err)
	}

	// Apply environment overrides
	m.applyEnvironmentOverrides(&config)

	// Validate configuration
	if err := m.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Cache the resolved configuration
	m.cache[configPath] = &config

	m.logger.Infof("Loaded extended configuration: %s (profile: %s)", config.Name, config.Profile)
	return &config, nil
}

// SaveConfig saves a configuration to file
func (m *ExtendedConfigManager) SaveConfig(configPath string, config *ExtendedBuildConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Update cache
	m.cache[configPath] = config

	m.logger.Infof("Saved configuration: %s", config.Name)
	return nil
}

// resolveInheritance resolves configuration inheritance
func (m *ExtendedConfigManager) resolveInheritance(config *ExtendedBuildConfig, baseDir string) error {
	if len(config.Inherits) == 0 {
		return nil
	}

	for _, inherit := range config.Inherits {
		inheritPath := filepath.Join(baseDir, inherit)
		parentConfig, err := m.LoadConfig(inheritPath)
		if err != nil {
			return fmt.Errorf("failed to load parent config %s: %w", inherit, err)
		}

		// Merge parent configuration (simple override strategy)
		if config.Description == "" {
			config.Description = parentConfig.Description
		}
		if config.Build.Dockerfile == "" {
			config.Build.Dockerfile = parentConfig.Build.Dockerfile
		}
		if config.Build.Context == "" {
			config.Build.Context = parentConfig.Build.Context
		}

		// Merge environment variables
		if config.Environment == nil {
			config.Environment = make(map[string]string)
		}
		for k, v := range parentConfig.Environment {
			if _, exists := config.Environment[k]; !exists {
				config.Environment[k] = v
			}
		}

		// Merge build args
		if config.Build.Args == nil {
			config.Build.Args = make(map[string]string)
		}
		for k, v := range parentConfig.Build.Args {
			if _, exists := config.Build.Args[k]; !exists {
				config.Build.Args[k] = v
			}
		}
	}

	return nil
}

// applyProfiles applies configuration profiles
func (m *ExtendedConfigManager) applyProfiles(config *ExtendedBuildConfig) error {
	if config.Profile == "" {
		return nil
	}

	profile, exists := config.Profiles[config.Profile]
	if !exists {
		return fmt.Errorf("profile %s not found", config.Profile)
	}

	// Apply profile settings
	for key, value := range profile.Settings {
		switch key {
		case "build.dockerfile":
			if str, ok := value.(string); ok && config.Build.Dockerfile == "" {
				config.Build.Dockerfile = str
			}
		case "build.context":
			if str, ok := value.(string); ok && config.Build.Context == "" {
				config.Build.Context = str
			}
		case "runtime.memory":
			if str, ok := value.(string); ok && config.Runtime.Memory == "" {
				config.Runtime.Memory = str
			}
		case "runtime.cpus":
			if str, ok := value.(string); ok && config.Runtime.CPUs == "" {
				config.Runtime.CPUs = str
			}
		case "storage.driver":
			if str, ok := value.(string); ok && config.Storage.Driver == "" {
				config.Storage.Driver = str
			}
		case "cache.enabled":
			if b, ok := value.(bool); ok {
				config.Cache.Enabled = b
			}
		}
	}

	m.logger.Debugf("Applied profile: %s", profile.Name)
	return nil
}

// applyEnvironmentOverrides applies environment variable overrides
func (m *ExtendedConfigManager) applyEnvironmentOverrides(config *ExtendedBuildConfig) {
	// Override with environment variables
	if val := os.Getenv("BUILDAH_DOCKERFILE"); val != "" {
		config.Build.Dockerfile = val
	}
	if val := os.Getenv("BUILDAH_CONTEXT"); val != "" {
		config.Build.Context = val
	}
	if val := os.Getenv("BUILDAH_RUNTIME"); val != "" {
		config.Runtime.Runtime = val
	}
	if val := os.Getenv("BUILDAH_STORAGE_DRIVER"); val != "" {
		config.Storage.Driver = val
	}

	// Set defaults if not specified
	if config.Build.Dockerfile == "" {
		config.Build.Dockerfile = "Dockerfile"
	}
	if config.Runtime.Runtime == "" {
		config.Runtime.Runtime = "runc"
	}
	if config.Storage.Driver == "" {
		config.Storage.Driver = "overlay"
	}
	if config.Network.Mode == "" {
		config.Network.Mode = "bridge"
	}
}

// validateConfig validates the configuration
func (m *ExtendedConfigManager) validateConfig(config *ExtendedBuildConfig) error {
	if config.Name == "" {
		return fmt.Errorf("configuration name is required")
	}

	if len(config.Name) > 253 {
		return fmt.Errorf("configuration name too long (max 253 characters)")
	}

	if config.Build.Context != "" {
		if !filepath.IsAbs(config.Build.Context) {
			abs, err := filepath.Abs(config.Build.Context)
			if err != nil {
				return fmt.Errorf("failed to resolve context path: %w", err)
			}
			config.Build.Context = abs
		}
	}

	return nil
}

// CreateBuilder creates a buildah builder from the configuration
func (m *ExtendedConfigManager) CreateBuilder(ctx context.Context, config *ExtendedBuildConfig) (*buildah.Builder, error) {
	buildOptions := define.BuildOptions{
		FromImage:        config.Build.Target,
		PullPolicy:       define.PullIfMissing,
		ContextDirectory: config.Build.Context,
		Args:             config.Build.Args,
		Labels:           config.Build.Labels,
	}

	if config.Build.Pull {
		buildOptions.PullPolicy = define.PullAlways
	}

	if config.Build.NoCache {
		buildOptions.NoCache = true
	}

	builder, err := buildah.NewBuilder(ctx, m.store, buildOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create builder: %w", err)
	}

	// Apply runtime options
	if config.Runtime.Memory != "" {
		builder.SetMemory(config.Runtime.Memory)
	}
	if config.Runtime.CPUs != "" {
		builder.SetCPUs(config.Runtime.CPUs)
	}
	if config.Runtime.User != "" {
		builder.SetUser(config.Runtime.User)
	}
	if config.Runtime.WorkingDir != "" {
		builder.SetWorkDir(config.Runtime.WorkingDir)
	}

	// Apply environment variables
	for k, v := range config.Runtime.Env {
		builder.SetEnv(k, v)
	}

	// Apply volumes
	for _, volume := range config.Runtime.Volumes {
		builder.AddVolume(volume)
	}

	m.logger.Infof("Created builder for configuration: %s", config.Name)
	return builder, nil
}

// GetConfigSummary returns a summary of the configuration
func (m *ExtendedConfigManager) GetConfigSummary(config *ExtendedBuildConfig) map[string]interface{} {
	summary := map[string]interface{}{
		"name":    config.Name,
		"version": config.Version,
		"profile": config.Profile,
		"build": map[string]interface{}{
			"dockerfile": config.Build.Dockerfile,
			"context":    config.Build.Context,
			"pull":       config.Build.Pull,
			"no_cache":   config.Build.NoCache,
		},
		"runtime": map[string]interface{}{
			"runtime": config.Runtime.Runtime,
			"memory":  config.Runtime.Memory,
			"cpus":    config.Runtime.CPUs,
		},
		"storage": map[string]interface{}{
			"driver": config.Storage.Driver,
			"root":   config.Storage.Root,
		},
		"features": map[string]interface{}{
			"parallel_builds": config.Features.ParallelBuilds,
			"layer_caching":   config.Features.LayerCaching,
			"buildkit":        config.Features.BuildKit,
		},
	}

	return summary
}

// ClearCache clears the configuration cache
func (m *ExtendedConfigManager) ClearCache() {
	m.cache = make(map[string]*ExtendedBuildConfig)
	m.logger.Debugf("Configuration cache cleared")
}

// SetCacheTTL sets the cache TTL
func (m *ExtendedConfigManager) SetCacheTTL(ttl time.Duration) {
	m.cacheTTL = ttl
	m.logger.Debugf("Cache TTL set to: %s", ttl)
}

// EnableCache enables or disables caching
func (m *ExtendedConfigManager) EnableCache(enabled bool) {
	if !enabled {
		m.ClearCache()
	}
	m.logger.Debugf("Configuration caching: %t", enabled)
}