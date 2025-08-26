package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/containers/buildah"
	"github.com/containers/buildah/define"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/sirupsen/logrus"
)

// BuildConfig represents the main configuration for build operations
type BuildConfig struct {
	// Build configuration
	From                string            `json:"from,omitempty"`
	Tag                 string            `json:"tag,omitempty"`
	File                string            `json:"file,omitempty"`
	Context             string            `json:"context,omitempty"`
	Platform            string            `json:"platform,omitempty"`
	
	// Build arguments and labels
	Args                map[string]string `json:"args,omitempty"`
	Labels              map[string]string `json:"labels,omitempty"`
	
	// Storage configuration
	Storage             *StoreConfig      `json:"storage,omitempty"`
	
	// Runtime configuration
	Runtime             string            `json:"runtime,omitempty"`
	Isolation           define.Isolation  `json:"isolation,omitempty"`
	
	// Security configuration
	User                string            `json:"user,omitempty"`
	Capabilities        []string          `json:"capabilities,omitempty"`
	
	// Network configuration
	NetworkMode         string            `json:"network_mode,omitempty"`
	HTTPProxy           string            `json:"http_proxy,omitempty"`
	HTTPSProxy          string            `json:"https_proxy,omitempty"`
	
	// Output configuration
	Format              string            `json:"format,omitempty"`
	Push                bool              `json:"push,omitempty"`
	
	// System context
	SystemContext       *types.SystemContext `json:"-"`
}

// ConfigManager handles configuration validation and management
type ConfigManager struct {
	config        *BuildConfig
	systemContext *types.SystemContext
	validated     bool
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(buildConfig *BuildConfig) *ConfigManager {
	if buildConfig == nil {
		buildConfig = DefaultBuildConfig()
	}
	
	return &ConfigManager{
		config:        buildConfig,
		systemContext: buildConfig.SystemContext,
	}
}

// DefaultBuildConfig returns a default build configuration
func DefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		From:        "scratch",
		Context:     ".",
		Platform:    runtime.GOARCH,
		Args:        make(map[string]string),
		Labels:      make(map[string]string),
		Storage:     DefaultStoreConfig(),
		Runtime:     "runc",
		Isolation:   define.IsolationDefault,
		NetworkMode: "bridge",
		Format:      "oci",
		Push:        false,
		SystemContext: &types.SystemContext{},
	}
}

// Validate performs comprehensive validation of the build configuration
func (cm *ConfigManager) Validate(ctx context.Context) error {
	if cm.validated {
		return nil
	}
	
	// Validate basic build parameters
	if err := cm.validateBuildParams(); err != nil {
		return fmt.Errorf("build parameters validation failed: %w", err)
	}
	
	// Validate storage configuration
	if err := cm.validateStorageConfig(); err != nil {
		return fmt.Errorf("storage configuration validation failed: %w", err)
	}
	
	cm.validated = true
	logrus.Info("Build configuration validation completed successfully")
	
	return nil
}

// validateBuildParams validates basic build parameters
func (cm *ConfigManager) validateBuildParams() error {
	config := cm.config
	
	// Validate context directory
	if config.Context != "" {
		if !filepath.IsAbs(config.Context) {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}
			config.Context = filepath.Join(wd, config.Context)
		}
		
		if _, err := os.Stat(config.Context); err != nil {
			return fmt.Errorf("context directory %s does not exist: %w", config.Context, err)
		}
	}
	
	// Validate Containerfile/Dockerfile
	if config.File != "" {
		if !filepath.IsAbs(config.File) {
			config.File = filepath.Join(config.Context, config.File)
		}
		
		if _, err := os.Stat(config.File); err != nil {
			return fmt.Errorf("containerfile %s does not exist: %w", config.File, err)
		}
	}
	
	// Validate platform
	if config.Platform != "" {
		validPlatforms := []string{"amd64", "arm64", "arm", "ppc64le", "s390x", "386"}
		valid := false
		for _, p := range validPlatforms {
			if strings.Contains(config.Platform, p) {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unsupported platform: %s", config.Platform)
		}
	}
	
	return nil
}

// validateStorageConfig validates storage configuration
func (cm *ConfigManager) validateStorageConfig() error {
	if cm.config.Storage == nil {
		return nil
	}
	
	storage := cm.config.Storage
	
	// Validate root directory
	if storage.RootDir != "" {
		if err := os.MkdirAll(storage.RootDir, 0755); err != nil {
			return fmt.Errorf("failed to create storage root directory: %w", err)
		}
	}
	
	// Validate run root
	if storage.RunRoot != "" {
		if err := os.MkdirAll(storage.RunRoot, 0755); err != nil {
			return fmt.Errorf("failed to create storage run directory: %w", err)
		}
	}
	
	// Validate graph driver
	validDrivers := []string{"overlay", "vfs", "btrfs", "zfs"}
	if storage.GraphDriver != "" {
		valid := false
		for _, driver := range validDrivers {
			if storage.GraphDriver == driver {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unsupported storage driver: %s", storage.GraphDriver)
		}
	}
	
	return nil
}

// GetConfig returns the validated build configuration
func (cm *ConfigManager) GetConfig() *BuildConfig {
	return cm.config
}

// GetSystemContext returns the system context for image operations
func (cm *ConfigManager) GetSystemContext() *types.SystemContext {
	if cm.systemContext == nil {
		cm.systemContext = &types.SystemContext{}
	}
	return cm.systemContext
}

// CreateBuildahOptions converts the configuration to Buildah builder options
func (cm *ConfigManager) CreateBuildahOptions(ctx context.Context, store storage.Store) (buildah.BuilderOptions, error) {
	if !cm.validated {
		return buildah.BuilderOptions{}, fmt.Errorf("configuration not validated")
	}
	
	config := cm.config
	
	options := buildah.BuilderOptions{
		FromImage:         config.From,
		Container:         "",
		PullPolicy:        define.PullIfMissing,
		ReportWriter:      os.Stderr,
		Store:             store,
		SystemContext:     cm.GetSystemContext(),
		Isolation:         config.Isolation,
		NamespaceOptions:  []define.NamespaceOption{},
		ConfigureNetwork:  define.NetworkDefault,
		Capabilities:      config.Capabilities,
		Args:              config.Args,
		Format:            config.Format,
		MaxPullPushRetries: 3,
		PullPushRetryDelay: 2 * time.Second,
	}
	
	// Configure network settings
	switch config.NetworkMode {
	case "host":
		options.ConfigureNetwork = define.NetworkDisabled
	case "none":
		options.ConfigureNetwork = define.NetworkDisabled
	}
	
	return options, nil
}

// ApplyEnvironmentVariables applies environment configuration
func (cm *ConfigManager) ApplyEnvironmentVariables() {
	config := cm.config
	
	// Apply proxy settings
	if config.HTTPProxy != "" {
		os.Setenv("HTTP_PROXY", config.HTTPProxy)
		os.Setenv("http_proxy", config.HTTPProxy)
	}
	if config.HTTPSProxy != "" {
		os.Setenv("HTTPS_PROXY", config.HTTPSProxy)
		os.Setenv("https_proxy", config.HTTPSProxy)
	}
}

// UpdateConfig updates the configuration with new values
func (cm *ConfigManager) UpdateConfig(updates map[string]interface{}) error {
	cm.validated = false // Mark as needing re-validation
	
	for key, value := range updates {
		switch key {
		case "from":
			if v, ok := value.(string); ok {
				cm.config.From = v
			}
		case "tag":
			if v, ok := value.(string); ok {
				cm.config.Tag = v
			}
		case "platform":
			if v, ok := value.(string); ok {
				cm.config.Platform = v
			}
		default:
			logrus.Warnf("Unknown configuration key: %s", key)
		}
	}
	
	return nil
}