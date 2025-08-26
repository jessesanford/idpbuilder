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
	"reflect"
	"strings"
	"time"

	"github.com/containers/buildah"
	"github.com/containers/buildah/define"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/sirupsen/logrus"
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

	// Inheritance configuration
	Inherits []string          `json:"inherits,omitempty" yaml:"inherits,omitempty"`
	Profiles map[string]Config `json:"profiles,omitempty" yaml:"profiles,omitempty"`

	// Advanced features
	Features FeatureFlags `json:"features,omitempty" yaml:"features,omitempty"`
	
	// Metadata
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	
	// Validation and processing
	lastValidated time.Time
	validated     bool
	merged        bool
}

// BuildOptions defines build-specific configuration options
type BuildOptions struct {
	Dockerfile   string            `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty"`
	Context      string            `json:"context,omitempty" yaml:"context,omitempty"`
	Args         map[string]string `json:"args,omitempty" yaml:"args,omitempty"`
	Target       string            `json:"target,omitempty" yaml:"target,omitempty"`
	Platform     []string          `json:"platform,omitempty" yaml:"platform,omitempty"`
	Output       string            `json:"output,omitempty" yaml:"output,omitempty"`
	Tags         []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Pull         bool              `json:"pull" yaml:"pull"`
	NoCache      bool              `json:"no_cache" yaml:"no_cache"`
	RmIntermediate bool            `json:"rm_intermediate" yaml:"rm_intermediate"`
	Squash       bool              `json:"squash" yaml:"squash"`
}

// RuntimeOptions defines runtime-specific configuration
type RuntimeOptions struct {
	Runtime     string            `json:"runtime,omitempty" yaml:"runtime,omitempty"`
	Isolation   define.Isolation  `json:"isolation,omitempty" yaml:"isolation,omitempty"`
	User        string            `json:"user,omitempty" yaml:"user,omitempty"`
	WorkingDir  string            `json:"working_dir,omitempty" yaml:"working_dir,omitempty"`
	Env         map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Volumes     []string          `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Ulimits     []string          `json:"ulimits,omitempty" yaml:"ulimits,omitempty"`
	Memory      int64             `json:"memory,omitempty" yaml:"memory,omitempty"`
	MemorySwap  int64             `json:"memory_swap,omitempty" yaml:"memory_swap,omitempty"`
	CPUShares   uint64            `json:"cpu_shares,omitempty" yaml:"cpu_shares,omitempty"`
	CPUQuota    int64             `json:"cpu_quota,omitempty" yaml:"cpu_quota,omitempty"`
	CPUPeriod   uint64            `json:"cpu_period,omitempty" yaml:"cpu_period,omitempty"`
}

// StorageOptions defines storage configuration
type StorageOptions struct {
	Driver        string            `json:"driver,omitempty" yaml:"driver,omitempty"`
	Root          string            `json:"root,omitempty" yaml:"root,omitempty"`
	RunRoot       string            `json:"run_root,omitempty" yaml:"run_root,omitempty"`
	Options       map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
	UIDMap        []string          `json:"uid_map,omitempty" yaml:"uid_map,omitempty"`
	GIDMap        []string          `json:"gid_map,omitempty" yaml:"gid_map,omitempty"`
	RemapUser     string            `json:"remap_user,omitempty" yaml:"remap_user,omitempty"`
	RemapGroup    string            `json:"remap_group,omitempty" yaml:"remap_group,omitempty"`
}

// NetworkOptions defines network configuration
type NetworkOptions struct {
	Mode        string            `json:"mode,omitempty" yaml:"mode,omitempty"`
	DNS         []string          `json:"dns,omitempty" yaml:"dns,omitempty"`
	DNSSearch   []string          `json:"dns_search,omitempty" yaml:"dns_search,omitempty"`
	DNSOptions  []string          `json:"dns_options,omitempty" yaml:"dns_options,omitempty"`
	Hosts       map[string]string `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Ports       []string          `json:"ports,omitempty" yaml:"ports,omitempty"`
}

// SecurityOptions defines security configuration
type SecurityOptions struct {
	Privileged     bool     `json:"privileged" yaml:"privileged"`
	ReadOnly       bool     `json:"read_only" yaml:"read_only"`
	NoNewPrivs     bool     `json:"no_new_privs" yaml:"no_new_privs"`
	SecurityOpts   []string `json:"security_opts,omitempty" yaml:"security_opts,omitempty"`
	AppArmorProfile string  `json:"apparmor_profile,omitempty" yaml:"apparmor_profile,omitempty"`
	SeccompProfile  string  `json:"seccomp_profile,omitempty" yaml:"seccomp_profile,omitempty"`
	SELinuxLabel    string  `json:"selinux_label,omitempty" yaml:"selinux_label,omitempty"`
	CapAdd          []string `json:"cap_add,omitempty" yaml:"cap_add,omitempty"`
	CapDrop         []string `json:"cap_drop,omitempty" yaml:"cap_drop,omitempty"`
}

// CacheOptions defines cache configuration
type CacheOptions struct {
	Enabled       bool              `json:"enabled" yaml:"enabled"`
	Dir           string            `json:"dir,omitempty" yaml:"dir,omitempty"`
	TTL           time.Duration     `json:"ttl,omitempty" yaml:"ttl,omitempty"`
	MaxSize       int64             `json:"max_size,omitempty" yaml:"max_size,omitempty"`
	Compression   string            `json:"compression,omitempty" yaml:"compression,omitempty"`
	Layers        map[string]string `json:"layers,omitempty" yaml:"layers,omitempty"`
}

// FeatureFlags defines experimental and optional features
type FeatureFlags struct {
	Experimental    bool `json:"experimental" yaml:"experimental"`
	LayerCaching    bool `json:"layer_caching" yaml:"layer_caching"`
	ParallelBuilds  bool `json:"parallel_builds" yaml:"parallel_builds"`
	AdvancedCaching bool `json:"advanced_caching" yaml:"advanced_caching"`
	Metrics         bool `json:"metrics" yaml:"metrics"`
	Tracing         bool `json:"tracing" yaml:"tracing"`
}

// Config represents a reusable configuration profile
type Config struct {
	Build    *BuildOptions    `json:"build,omitempty" yaml:"build,omitempty"`
	Runtime  *RuntimeOptions  `json:"runtime,omitempty" yaml:"runtime,omitempty"`
	Storage  *StorageOptions  `json:"storage,omitempty" yaml:"storage,omitempty"`
	Network  *NetworkOptions  `json:"network,omitempty" yaml:"network,omitempty"`
	Security *SecurityOptions `json:"security,omitempty" yaml:"security,omitempty"`
	Cache    *CacheOptions    `json:"cache,omitempty" yaml:"cache,omitempty"`
	Features *FeatureFlags    `json:"features,omitempty" yaml:"features,omitempty"`
	Labels   map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// ExtendedConfigManager manages advanced configuration operations
type ExtendedConfigManager struct {
	store            storage.Store
	configs          map[string]*ExtendedBuildConfig
	profiles         map[string]*Config
	validationRules  []ValidationRule
	transformers     []ConfigTransformer
	logger           *logrus.Logger
	cacheEnabled     bool
	cacheTTL         time.Duration
}

// ValidationRule defines a configuration validation rule
type ValidationRule func(*ExtendedBuildConfig) error

// ConfigTransformer defines a configuration transformation function
type ConfigTransformer func(*ExtendedBuildConfig) (*ExtendedBuildConfig, error)

// NewExtendedConfigManager creates a new extended configuration manager
func NewExtendedConfigManager(store storage.Store, logger *logrus.Logger) *ExtendedConfigManager {
	if logger == nil {
		logger = logrus.New()
	}

	return &ExtendedConfigManager{
		store:           store,
		configs:         make(map[string]*ExtendedBuildConfig),
		profiles:        make(map[string]*Config),
		validationRules: getDefaultValidationRules(),
		transformers:    getDefaultTransformers(),
		logger:          logger,
		cacheEnabled:    true,
		cacheTTL:        time.Hour,
	}
}

// LoadConfig loads a configuration from file with inheritance and profile resolution
func (ecm *ExtendedConfigManager) LoadConfig(configPath string) (*ExtendedBuildConfig, error) {
	ecm.logger.WithField("path", configPath).Debug("Loading extended configuration")

	// Check cache first
	if ecm.cacheEnabled {
		if config, exists := ecm.configs[configPath]; exists {
			if time.Since(config.lastValidated) < ecm.cacheTTL {
				ecm.logger.WithField("path", configPath).Debug("Returning cached configuration")
				return config, nil
			}
		}
	}

	// Load from file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config ExtendedBuildConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Resolve inheritance
	if err := ecm.resolveInheritance(&config, filepath.Dir(configPath)); err != nil {
		return nil, fmt.Errorf("failed to resolve inheritance: %w", err)
	}

	// Apply profile overrides
	if err := ecm.applyProfile(&config); err != nil {
		return nil, fmt.Errorf("failed to apply profile: %w", err)
	}

	// Apply environment-specific overrides
	if err := ecm.applyEnvironmentOverrides(&config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	// Validate configuration
	if err := ecm.ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Apply transformations
	transformedConfig, err := ecm.applyTransformations(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to apply transformations: %w", err)
	}

	// Mark as validated and cache
	transformedConfig.lastValidated = time.Now()
	transformedConfig.validated = true
	transformedConfig.merged = true

	if ecm.cacheEnabled {
		ecm.configs[configPath] = transformedConfig
	}

	ecm.logger.WithField("name", transformedConfig.Name).Info("Extended configuration loaded successfully")
	return transformedConfig, nil
}

// resolveInheritance resolves configuration inheritance chain
func (ecm *ExtendedConfigManager) resolveInheritance(config *ExtendedBuildConfig, baseDir string) error {
	if len(config.Inherits) == 0 {
		return nil
	}

	ecm.logger.WithField("inherits", config.Inherits).Debug("Resolving inheritance chain")

	// Process inheritance in reverse order (last takes precedence)
	for i := len(config.Inherits) - 1; i >= 0; i-- {
		parentPath := config.Inherits[i]
		if !filepath.IsAbs(parentPath) {
			parentPath = filepath.Join(baseDir, parentPath)
		}

		parentConfig, err := ecm.LoadConfig(parentPath)
		if err != nil {
			return fmt.Errorf("failed to load parent config %s: %w", parentPath, err)
		}

		// Merge parent config into current config
		if err := ecm.mergeConfigs(config, parentConfig); err != nil {
			return fmt.Errorf("failed to merge parent config %s: %w", parentPath, err)
		}
	}

	return nil
}

// mergeConfigs merges parent configuration into child configuration
func (ecm *ExtendedConfigManager) mergeConfigs(child, parent *ExtendedBuildConfig) error {
	// Merge using reflection to handle all fields dynamically
	childValue := reflect.ValueOf(child).Elem()
	parentValue := reflect.ValueOf(parent).Elem()

	return ecm.mergeStructs(childValue, parentValue)
}

// mergeStructs recursively merges struct fields
func (ecm *ExtendedConfigManager) mergeStructs(child, parent reflect.Value) error {
	for i := 0; i < parent.NumField(); i++ {
		parentField := parent.Field(i)
		childField := child.Field(i)

		if !childField.CanSet() {
			continue
		}

		// Skip internal fields
		fieldName := parent.Type().Field(i).Name
		if strings.HasPrefix(fieldName, "last") || fieldName == "validated" || fieldName == "merged" {
			continue
		}

		switch parentField.Kind() {
		case reflect.Map:
			if parentField.Len() > 0 {
				if childField.IsNil() {
					childField.Set(reflect.MakeMap(parentField.Type()))
				}
				for _, key := range parentField.MapKeys() {
					if !childField.MapIndex(key).IsValid() {
						childField.SetMapIndex(key, parentField.MapIndex(key))
					}
				}
			}
		case reflect.Slice:
			if parentField.Len() > 0 && childField.Len() == 0 {
				childField.Set(parentField)
			}
		case reflect.Struct:
			if !parentField.IsZero() {
				if childField.IsZero() {
					childField.Set(parentField)
				} else {
					if err := ecm.mergeStructs(childField, parentField); err != nil {
						return err
					}
				}
			}
		case reflect.Ptr:
			if !parentField.IsNil() && childField.IsNil() {
				childField.Set(parentField)
			}
		case reflect.String:
			if parentField.String() != "" && childField.String() == "" {
				childField.Set(parentField)
			}
		case reflect.Bool:
			if !parentField.Bool() && !childField.Bool() {
				// Keep child value for booleans
			} else if parentField.Bool() && !childField.Bool() {
				childField.Set(parentField)
			}
		default:
			if !parentField.IsZero() && childField.IsZero() {
				childField.Set(parentField)
			}
		}
	}

	return nil
}

// applyProfile applies profile-specific configuration overrides
func (ecm *ExtendedConfigManager) applyProfile(config *ExtendedBuildConfig) error {
	if config.Profile == "" {
		return nil
	}

	profile, exists := config.Profiles[config.Profile]
	if !exists {
		return fmt.Errorf("profile %s not found", config.Profile)
	}

	ecm.logger.WithField("profile", config.Profile).Debug("Applying profile configuration")

	// Apply profile overrides
	if profile.Build != nil {
		if err := ecm.mergeStructs(reflect.ValueOf(&config.Build).Elem(), reflect.ValueOf(profile.Build).Elem()); err != nil {
			return fmt.Errorf("failed to merge build profile: %w", err)
		}
	}

	if profile.Runtime != nil {
		if err := ecm.mergeStructs(reflect.ValueOf(&config.Runtime).Elem(), reflect.ValueOf(profile.Runtime).Elem()); err != nil {
			return fmt.Errorf("failed to merge runtime profile: %w", err)
		}
	}

	if profile.Storage != nil {
		if err := ecm.mergeStructs(reflect.ValueOf(&config.Storage).Elem(), reflect.ValueOf(profile.Storage).Elem()); err != nil {
			return fmt.Errorf("failed to merge storage profile: %w", err)
		}
	}

	if profile.Network != nil {
		if err := ecm.mergeStructs(reflect.ValueOf(&config.Network).Elem(), reflect.ValueOf(profile.Network).Elem()); err != nil {
			return fmt.Errorf("failed to merge network profile: %w", err)
		}
	}

	if profile.Security != nil {
		if err := ecm.mergeStructs(reflect.ValueOf(&config.Security).Elem(), reflect.ValueOf(profile.Security).Elem()); err != nil {
			return fmt.Errorf("failed to merge security profile: %w", err)
		}
	}

	if profile.Cache != nil {
		if err := ecm.mergeStructs(reflect.ValueOf(&config.Cache).Elem(), reflect.ValueOf(profile.Cache).Elem()); err != nil {
			return fmt.Errorf("failed to merge cache profile: %w", err)
		}
	}

	if profile.Features != nil {
		if err := ecm.mergeStructs(reflect.ValueOf(&config.Features).Elem(), reflect.ValueOf(profile.Features).Elem()); err != nil {
			return fmt.Errorf("failed to merge features profile: %w", err)
		}
	}

	// Merge labels and annotations
	if profile.Labels != nil {
		if config.Labels == nil {
			config.Labels = make(map[string]string)
		}
		for k, v := range profile.Labels {
			if _, exists := config.Labels[k]; !exists {
				config.Labels[k] = v
			}
		}
	}

	if profile.Annotations != nil {
		if config.Annotations == nil {
			config.Annotations = make(map[string]string)
		}
		for k, v := range profile.Annotations {
			if _, exists := config.Annotations[k]; !exists {
				config.Annotations[k] = v
			}
		}
	}

	return nil
}

// applyEnvironmentOverrides applies environment-specific configuration overrides
func (ecm *ExtendedConfigManager) applyEnvironmentOverrides(config *ExtendedBuildConfig) error {
	if config.Environment == nil {
		return nil
	}

	ecm.logger.WithField("environment", config.Environment).Debug("Applying environment overrides")

	// Apply environment variable overrides
	for key, value := range config.Environment {
		envValue := os.Getenv(key)
		if envValue != "" {
			config.Environment[key] = envValue
			ecm.logger.WithFields(logrus.Fields{
				"key":   key,
				"value": envValue,
			}).Debug("Applied environment override")
		} else {
			// Set environment variable if not present
			if err := os.Setenv(key, value); err != nil {
				ecm.logger.WithError(err).WithField("key", key).Warn("Failed to set environment variable")
			}
		}
	}

	return nil
}

// ValidateConfig validates the extended configuration
func (ecm *ExtendedConfigManager) ValidateConfig(config *ExtendedBuildConfig) error {
	ecm.logger.WithField("name", config.Name).Debug("Validating extended configuration")

	for _, rule := range ecm.validationRules {
		if err := rule(config); err != nil {
			return err
		}
	}

	return nil
}

// applyTransformations applies configuration transformations
func (ecm *ExtendedConfigManager) applyTransformations(config *ExtendedBuildConfig) (*ExtendedBuildConfig, error) {
	transformedConfig := config

	for _, transformer := range ecm.transformers {
		var err error
		transformedConfig, err = transformer(transformedConfig)
		if err != nil {
			return nil, err
		}
	}

	return transformedConfig, nil
}

// getDefaultValidationRules returns default validation rules
func getDefaultValidationRules() []ValidationRule {
	return []ValidationRule{
		validateName,
		validateBuildContext,
		validateStorageDriver,
		validateNetworkMode,
		validateSecurityOptions,
		validateCacheConfiguration,
		validateResourceLimits,
	}
}

// Validation rule implementations
func validateName(config *ExtendedBuildConfig) error {
	if config.Name == "" {
		return fmt.Errorf("configuration name is required")
	}
	if len(config.Name) > 253 {
		return fmt.Errorf("configuration name too long: %d characters (max: 253)", len(config.Name))
	}
	return nil
}

func validateBuildContext(config *ExtendedBuildConfig) error {
	if config.Build.Context != "" {
		if _, err := os.Stat(config.Build.Context); os.IsNotExist(err) {
			return fmt.Errorf("build context does not exist: %s", config.Build.Context)
		}
	}
	return nil
}

func validateStorageDriver(config *ExtendedBuildConfig) error {
	if config.Storage.Driver != "" {
		validDrivers := []string{"overlay", "devicemapper", "btrfs", "zfs", "vfs"}
		for _, driver := range validDrivers {
			if config.Storage.Driver == driver {
				return nil
			}
		}
		return fmt.Errorf("invalid storage driver: %s", config.Storage.Driver)
	}
	return nil
}

func validateNetworkMode(config *ExtendedBuildConfig) error {
	if config.Network.Mode != "" {
		validModes := []string{"bridge", "host", "none", "container"}
		for _, mode := range validModes {
			if config.Network.Mode == mode {
				return nil
			}
		}
		return fmt.Errorf("invalid network mode: %s", config.Network.Mode)
	}
	return nil
}

func validateSecurityOptions(config *ExtendedBuildConfig) error {
	// Validate security context combinations
	if config.Security.Privileged && config.Security.NoNewPrivs {
		return fmt.Errorf("privileged and no-new-privs are mutually exclusive")
	}
	return nil
}

func validateCacheConfiguration(config *ExtendedBuildConfig) error {
	if config.Cache.Enabled {
		if config.Cache.Dir == "" {
			return fmt.Errorf("cache directory is required when caching is enabled")
		}
		if config.Cache.MaxSize < 0 {
			return fmt.Errorf("cache max size cannot be negative")
		}
	}
	return nil
}

func validateResourceLimits(config *ExtendedBuildConfig) error {
	if config.Runtime.Memory < 0 {
		return fmt.Errorf("memory limit cannot be negative")
	}
	if config.Runtime.CPUQuota < 0 {
		return fmt.Errorf("CPU quota cannot be negative")
	}
	if config.Runtime.CPUPeriod < 0 {
		return fmt.Errorf("CPU period cannot be negative")
	}
	return nil
}

// getDefaultTransformers returns default configuration transformers
func getDefaultTransformers() []ConfigTransformer {
	return []ConfigTransformer{
		expandPaths,
		setDefaults,
		optimizeSettings,
	}
}

// Transformer implementations
func expandPaths(config *ExtendedBuildConfig) (*ExtendedBuildConfig, error) {
	// Expand relative paths to absolute paths
	if config.Build.Context != "" {
		abs, err := filepath.Abs(config.Build.Context)
		if err != nil {
			return nil, fmt.Errorf("failed to expand build context path: %w", err)
		}
		config.Build.Context = abs
	}

	if config.Storage.Root != "" {
		abs, err := filepath.Abs(config.Storage.Root)
		if err != nil {
			return nil, fmt.Errorf("failed to expand storage root path: %w", err)
		}
		config.Storage.Root = abs
	}

	if config.Cache.Dir != "" {
		abs, err := filepath.Abs(config.Cache.Dir)
		if err != nil {
			return nil, fmt.Errorf("failed to expand cache directory path: %w", err)
		}
		config.Cache.Dir = abs
	}

	return config, nil
}

func setDefaults(config *ExtendedBuildConfig) (*ExtendedBuildConfig, error) {
	// Set default values for unspecified options
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

	if config.Cache.Enabled && config.Cache.TTL == 0 {
		config.Cache.TTL = 24 * time.Hour
	}

	return config, nil
}

func optimizeSettings(config *ExtendedBuildConfig) (*ExtendedBuildConfig, error) {
	// Apply optimization based on configuration
	if config.Features.ParallelBuilds {
		// Enable parallel build optimizations
		if config.Runtime.CPUShares == 0 {
			config.Runtime.CPUShares = 1024
		}
	}

	if config.Features.LayerCaching {
		// Enable layer caching optimizations
		if !config.Cache.Enabled {
			config.Cache.Enabled = true
			config.Cache.Dir = filepath.Join(os.TempDir(), "buildah-cache")
		}
	}

	return config, nil
}

// SaveConfig saves configuration to file
func (ecm *ExtendedConfigManager) SaveConfig(configPath string, config *ExtendedBuildConfig) error {
	ecm.logger.WithField("path", configPath).Debug("Saving extended configuration")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	ecm.logger.WithField("name", config.Name).Info("Extended configuration saved successfully")
	return nil
}

// CreateBuildahBuilder creates a buildah builder from extended configuration
func (ecm *ExtendedConfigManager) CreateBuildahBuilder(ctx context.Context, config *ExtendedBuildConfig) (buildah.Builder, error) {
	ecm.logger.WithField("name", config.Name).Debug("Creating buildah builder from extended config")

	if !config.validated {
		if err := ecm.ValidateConfig(config); err != nil {
			return buildah.Builder{}, fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	// Create buildah options from extended configuration
	options := buildah.BuilderOptions{
		FromImage:        "scratch", // Will be overridden by Dockerfile FROM
		Container:        config.Name,
		PullPolicy:       define.PullIfMissing,
		SystemContext:    &types.SystemContext{},
		ConfigureNetwork: define.NetworkDefault,
		Isolation:        config.Runtime.Isolation,
	}

	// Configure pull policy
	if config.Build.Pull {
		options.PullPolicy = define.PullAlways
	}

	// Configure isolation
	if config.Runtime.Isolation != "" {
		options.Isolation = config.Runtime.Isolation
	}

	// Configure system context
	if config.Security.AppArmorProfile != "" {
		options.SystemContext.DockerInsecureSkipTLSVerify = types.OptionalBoolTrue
	}

	// Create the builder
	builder, err := buildah.NewBuilder(ctx, ecm.store, options)
	if err != nil {
		return buildah.Builder{}, fmt.Errorf("failed to create buildah builder: %w", err)
	}

	ecm.logger.WithField("name", config.Name).Info("Buildah builder created successfully")
	return *builder, nil
}

// GetConfigurationSummary returns a summary of the configuration
func (ecm *ExtendedConfigManager) GetConfigurationSummary(config *ExtendedBuildConfig) map[string]interface{} {
	summary := map[string]interface{}{
		"name":        config.Name,
		"version":     config.Version,
		"profile":     config.Profile,
		"validated":   config.validated,
		"merged":      config.merged,
		"inherits":    len(config.Inherits),
		"profiles":    len(config.Profiles),
		"features":    config.Features,
		"cache":       config.Cache.Enabled,
		"last_validated": config.lastValidated,
	}

	if len(config.Labels) > 0 {
		summary["labels"] = len(config.Labels)
	}

	if len(config.Annotations) > 0 {
		summary["annotations"] = len(config.Annotations)
	}

	return summary
}

// ClearCache clears the configuration cache
func (ecm *ExtendedConfigManager) ClearCache() {
	ecm.logger.Debug("Clearing extended configuration cache")
	ecm.configs = make(map[string]*ExtendedBuildConfig)
}

// SetCacheTTL sets the cache time-to-live
func (ecm *ExtendedConfigManager) SetCacheTTL(ttl time.Duration) {
	ecm.cacheTTL = ttl
	ecm.logger.WithField("ttl", ttl).Debug("Updated cache TTL")
}

// EnableCache enables or disables configuration caching
func (ecm *ExtendedConfigManager) EnableCache(enabled bool) {
	ecm.cacheEnabled = enabled
	if !enabled {
		ecm.ClearCache()
	}
	ecm.logger.WithField("enabled", enabled).Debug("Updated cache settings")
}