package certificates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadYAMLConfig loads configuration from a YAML file.
func LoadYAMLConfig(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	expandedPath := ResolvePath(path)
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML config: %w", err)
	}

	var trustConfig TrustStoreConfig
	if err := yaml.Unmarshal(data, &trustConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	loader := NewConfigLoader()
	return loader.convertToConfig(&trustConfig)
}

// LoadJSONConfig loads configuration from a JSON file.
func LoadJSONConfig(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	expandedPath := ResolvePath(path)
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON config: %w", err)
	}

	var trustConfig TrustStoreConfig
	if err := json.Unmarshal(data, &trustConfig); err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %w", err)
	}

	loader := NewConfigLoader()
	return loader.convertToConfig(&trustConfig)
}

// MergeConfigs merges two configurations with override taking precedence.
func MergeConfigs(base, override *Config) *Config {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	merged := &Config{}

	// Merge storage backend
	if override.StorageBackend != nil {
		merged.StorageBackend = override.StorageBackend
	} else if base.StorageBackend != nil {
		merged.StorageBackend = base.StorageBackend
	}

	// Merge validation rules
	ruleMap := make(map[string]*ValidationRule)
	for _, rule := range base.ValidationRules {
		ruleMap[rule.Name] = rule
	}
	for _, rule := range override.ValidationRules {
		ruleMap[rule.Name] = rule
	}
	for _, rule := range ruleMap {
		merged.ValidationRules = append(merged.ValidationRules, rule)
	}

	// Merge default pools and event handlers
	if len(override.DefaultPools) > 0 {
		merged.DefaultPools = override.DefaultPools
	} else {
		merged.DefaultPools = base.DefaultPools
	}

	if override.EventHandlers != nil {
		merged.EventHandlers = override.EventHandlers
	} else {
		merged.EventHandlers = base.EventHandlers
	}

	return merged
}

// ResolvePaths resolves all path references in the configuration.
func ResolvePaths(config *Config) error {
	if config == nil || config.StorageBackend == nil {
		return nil
	}

	config.StorageBackend.ConnectionString = ResolvePath(config.StorageBackend.ConnectionString)
	return nil
}

// ApplyEnvironmentOverrides applies environment variable overrides.
func ApplyEnvironmentOverrides(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	envPrefix := "TRUSTSTORE_"

	if storagePath := os.Getenv(envPrefix + "STORAGE_PATH"); storagePath != "" {
		if config.StorageBackend == nil {
			config.StorageBackend = &StorageConfig{}
		}
		config.StorageBackend.ConnectionString = ResolvePath(storagePath)
	}

	if eventsEnabled := os.Getenv(envPrefix + "EVENTS_ENABLED"); eventsEnabled != "" {
		enabled := eventsEnabled == "true"
		if config.EventHandlers == nil {
			config.EventHandlers = &EventConfig{}
		}
		config.EventHandlers.Enabled = enabled
	}

	return nil
}

// ResolvePath expands environment variables and resolves relative paths.
func ResolvePath(path string) string {
	if path == "" {
		return ""
	}

	expanded := os.ExpandEnv(path)
	if !filepath.IsAbs(expanded) {
		if abs, err := filepath.Abs(expanded); err == nil {
			expanded = abs
		}
	}

	return filepath.Clean(expanded)
}