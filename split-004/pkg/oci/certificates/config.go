package certificates

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// TrustStoreConfig represents the complete trust store configuration.
type TrustStoreConfig struct {
	StoragePath string                 `yaml:"storagePath" json:"storage_path"`
	Pools       map[string]*PoolConfig `yaml:"pools" json:"pools"`
	Validation  *ValidationConfig      `yaml:"validation" json:"validation"`
	Events      *EventsConfig          `yaml:"events" json:"events"`
}

// PoolConfig defines a certificate pool configuration.
type PoolConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Path    string `yaml:"path" json:"path"`
	Type    string `yaml:"type,omitempty" json:"type,omitempty"`
}

// ValidationConfig defines validation behavior.
type ValidationConfig struct {
	CheckExpiry   bool     `yaml:"checkExpiry" json:"check_expiry"`
	CheckChain    bool     `yaml:"checkChain" json:"check_chain"`
	CustomRules   []string `yaml:"customRules" json:"custom_rules"`
	WarnDays      int      `yaml:"warnDays,omitempty" json:"warn_days,omitempty"`
	FailOnExpired bool     `yaml:"failOnExpired,omitempty" json:"fail_on_expired,omitempty"`
}

// EventsConfig defines event handling configuration.
type EventsConfig struct {
	Enabled  bool     `yaml:"enabled" json:"enabled"`
	Handlers []string `yaml:"handlers" json:"handlers"`
	Buffer   int      `yaml:"buffer,omitempty" json:"buffer,omitempty"`
}

// ConfigLoaderImpl implements the ConfigLoader interface.
type ConfigLoaderImpl struct {
	defaultPath string
	envPrefix   string
}

// NewConfigLoader creates a new ConfigLoaderImpl instance.
func NewConfigLoader() *ConfigLoaderImpl {
	return &ConfigLoaderImpl{
		defaultPath: "/etc/trust-store/config.yaml",
		envPrefix:   "TRUSTSTORE_",
	}
}

// LoadConfig loads configuration from the specified source.
func (c *ConfigLoaderImpl) LoadConfig(ctx context.Context, source string) (*Config, error) {
	if source == "" {
		source = c.defaultPath
	}

	trustConfig, err := c.loadTrustStoreConfig(source)
	if err != nil {
		return nil, fmt.Errorf("failed to load trust store config: %w", err)
	}

	config, err := c.convertToConfig(trustConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to convert config format: %w", err)
	}

	if err := c.applyEnvironmentOverrides(config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to the specified destination.
func (c *ConfigLoaderImpl) SaveConfig(ctx context.Context, config *Config, destination string) error {
	if destination == "" {
		return fmt.Errorf("destination path cannot be empty")
	}

	trustConfig, err := c.convertFromConfig(config)
	if err != nil {
		return fmt.Errorf("failed to convert config format: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(destination))
	var data []byte

	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(trustConfig)
	case ".json":
		data, err = json.MarshalIndent(trustConfig, "", "  ")
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(destination, data, 0644)
}

// ValidateConfig validates the provided configuration.
func (c *ConfigLoaderImpl) ValidateConfig(ctx context.Context, config *Config) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
	}

	var errors []ValidationError

	if config.StorageBackend == nil {
		errors = append(errors, ValidationError{Code: "MISSING_STORAGE", Message: "storage backend required"})
	} else if config.StorageBackend.Type == "" {
		errors = append(errors, ValidationError{Code: "EMPTY_STORAGE_TYPE", Message: "storage type cannot be empty"})
	}

	if config.EventHandlers != nil && config.EventHandlers.BufferSize < 0 {
		errors = append(errors, ValidationError{Code: "NEGATIVE_BUFFER", Message: "buffer size cannot be negative"})
	}

	for _, rule := range config.ValidationRules {
		if rule.Name == "" {
			errors = append(errors, ValidationError{Code: "EMPTY_RULE_NAME", Message: "validation rule name cannot be empty"})
		}
	}

	result.Errors = errors
	result.Valid = len(errors) == 0
	return result, nil
}

// GetDefaultConfig returns a default configuration.
func (c *ConfigLoaderImpl) GetDefaultConfig(ctx context.Context) (*Config, error) {
	trustConfig := &TrustStoreConfig{
		StoragePath: "/var/lib/trust-store",
		Pools: map[string]*PoolConfig{
			"system": {Enabled: true, Path: "/etc/ssl/certs", Type: "system"},
			"custom": {Enabled: true, Path: "/opt/custom-certs", Type: "custom"},
		},
		Validation: &ValidationConfig{
			CheckExpiry: true, CheckChain: true, CustomRules: []string{}, WarnDays: 30,
		},
		Events: &EventsConfig{
			Enabled: true, Handlers: []string{"log", "metrics"}, Buffer: 100,
		},
	}

	return c.convertToConfig(trustConfig)
}

// loadTrustStoreConfig loads trust store configuration from file.
func (c *ConfigLoaderImpl) loadTrustStoreConfig(path string) (*TrustStoreConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config TrustStoreConfig
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	case ".json":
		err = json.Unmarshal(data, &config)
	default:
		return nil, fmt.Errorf("unsupported config format: %s", ext)
	}

	return &config, err
}

// convertToConfig converts TrustStoreConfig to Config format.
func (c *ConfigLoaderImpl) convertToConfig(trustConfig *TrustStoreConfig) (*Config, error) {
	config := &Config{DefaultPools: []string{}}

	if trustConfig.StoragePath != "" {
		config.StorageBackend = &StorageConfig{
			Type:             "filesystem",
			ConnectionString: trustConfig.StoragePath,
			Options:          map[string]interface{}{"pools": trustConfig.Pools},
		}
	}

	if trustConfig.Validation != nil {
		config.ValidationRules = []*ValidationRule{
			{Name: "check_expiry", Enabled: trustConfig.Validation.CheckExpiry, Critical: trustConfig.Validation.FailOnExpired},
			{Name: "check_chain", Enabled: trustConfig.Validation.CheckChain, Critical: true},
		}
	}

	if trustConfig.Events != nil {
		config.EventHandlers = &EventConfig{
			Enabled: trustConfig.Events.Enabled, BufferSize: trustConfig.Events.Buffer, Handlers: trustConfig.Events.Handlers,
		}
	}

	for poolName, poolConfig := range trustConfig.Pools {
		if poolConfig.Enabled {
			config.DefaultPools = append(config.DefaultPools, poolName)
		}
	}

	return config, nil
}

// convertFromConfig converts Config format to TrustStoreConfig.
func (c *ConfigLoaderImpl) convertFromConfig(config *Config) (*TrustStoreConfig, error) {
	trustConfig := &TrustStoreConfig{Pools: make(map[string]*PoolConfig)}

	if config.StorageBackend != nil {
		trustConfig.StoragePath = config.StorageBackend.ConnectionString
		if pools, ok := config.StorageBackend.Options["pools"].(map[string]*PoolConfig); ok {
			trustConfig.Pools = pools
		}
	}

	trustConfig.Validation = &ValidationConfig{CustomRules: []string{}}
	for _, rule := range config.ValidationRules {
		switch rule.Name {
		case "check_expiry":
			trustConfig.Validation.CheckExpiry = rule.Enabled
			trustConfig.Validation.FailOnExpired = rule.Critical
		case "check_chain":
			trustConfig.Validation.CheckChain = rule.Enabled
		}
	}

	if config.EventHandlers != nil {
		trustConfig.Events = &EventsConfig{
			Enabled: config.EventHandlers.Enabled, Buffer: config.EventHandlers.BufferSize, Handlers: config.EventHandlers.Handlers,
		}
	}

	return trustConfig, nil
}

// applyEnvironmentOverrides applies environment variable overrides.
func (c *ConfigLoaderImpl) applyEnvironmentOverrides(config *Config) error {
	if storagePath := os.Getenv(c.envPrefix + "STORAGE_PATH"); storagePath != "" {
		if config.StorageBackend == nil {
			config.StorageBackend = &StorageConfig{}
		}
		config.StorageBackend.ConnectionString = storagePath
	}

	if eventsEnabled := os.Getenv(c.envPrefix + "EVENTS_ENABLED"); eventsEnabled != "" {
		enabled := strings.ToLower(eventsEnabled) == "true"
		if config.EventHandlers == nil {
			config.EventHandlers = &EventConfig{}
		}
		config.EventHandlers.Enabled = enabled
	}

	return nil
}