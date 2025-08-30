package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	// Config holds the application configuration
	Config *AppConfig
	// ConfigFile is the path to the configuration file
	ConfigFile string
	// DefaultConfigPaths are the default locations to look for config files
	DefaultConfigPaths = []string{
		"./idpbuilder.yaml",
		"./idpbuilder.yml",
		"$HOME/.idpbuilder/config.yaml",
		"$HOME/.idpbuilder/config.yml",
		"/etc/idpbuilder/config.yaml",
	}
)

// AppConfig represents the application configuration
type AppConfig struct {
	// Version is the configuration version
	Version string `yaml:"version" json:"version"`
	
	// LogLevel sets the logging level
	LogLevel string `yaml:"logLevel" json:"logLevel"`
	
	// Output configures output settings
	Output OutputConfig `yaml:"output" json:"output"`
	
	// Defaults holds default values for commands
	Defaults DefaultsConfig `yaml:"defaults" json:"defaults"`
	
	// Clusters holds cluster configuration
	Clusters map[string]ClusterConfig `yaml:"clusters" json:"clusters"`
}

// OutputConfig configures output formatting
type OutputConfig struct {
	// Format sets the default output format
	Format string `yaml:"format" json:"format"`
	
	// Color enables colored output
	Color bool `yaml:"color" json:"color"`
	
	// Verbose enables verbose output
	Verbose bool `yaml:"verbose" json:"verbose"`
}

// DefaultsConfig holds default values for commands
type DefaultsConfig struct {
	// Namespace is the default namespace
	Namespace string `yaml:"namespace" json:"namespace"`
	
	// Timeout is the default timeout for operations
	Timeout string `yaml:"timeout" json:"timeout"`
	
	// DryRun enables dry-run mode by default
	DryRun bool `yaml:"dryRun" json:"dryRun"`
}

// ClusterConfig holds cluster-specific configuration
type ClusterConfig struct {
	// Name is the cluster name
	Name string `yaml:"name" json:"name"`
	
	// Context is the kubeconfig context
	Context string `yaml:"context" json:"context"`
	
	// Namespace is the default namespace for this cluster
	Namespace string `yaml:"namespace" json:"namespace"`
	
	// URL is the cluster API URL
	URL string `yaml:"url" json:"url"`
}

// LoadConfig loads configuration from various sources
func LoadConfig() error {
	// Initialize with defaults
	Config = &AppConfig{
		Version:  "v1",
		LogLevel: "info",
		Output: OutputConfig{
			Format:  "table",
			Color:   false,
			Verbose: false,
		},
		Defaults: DefaultsConfig{
			Namespace: "default",
			Timeout:   "30s",
			DryRun:    false,
		},
		Clusters: make(map[string]ClusterConfig),
	}

	// Load from file if specified
	if ConfigFile != "" {
		return loadConfigFromFile(ConfigFile)
	}

	// Try default locations
	for _, path := range DefaultConfigPaths {
		expandedPath := expandPath(path)
		if _, err := os.Stat(expandedPath); err == nil {
			ConfigFile = expandedPath
			return loadConfigFromFile(expandedPath)
		}
	}

	// No config file found, use defaults
	LogDebug("No configuration file found, using defaults")
	return nil
}

// loadConfigFromFile loads configuration from a specific file
func loadConfigFromFile(path string) error {
	LogDebug("Loading configuration from: %s", path)
	
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Determine file format
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, Config)
	case ".json":
		err = json.Unmarshal(data, Config)
	default:
		// Try YAML first, then JSON
		err = yaml.Unmarshal(data, Config)
		if err != nil {
			err = json.Unmarshal(data, Config)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Validate configuration
	return validateConfig()
}

// SaveConfig saves the current configuration to a file
func SaveConfig(path string) error {
	if Config == nil {
		return fmt.Errorf("no configuration to save")
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	// Determine format from extension
	ext := strings.ToLower(filepath.Ext(path))
	var data []byte
	var err error

	switch ext {
	case ".json":
		data, err = json.MarshalIndent(Config, "", "  ")
	default: // Default to YAML
		data, err = yaml.Marshal(Config)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}

	LogInfo("Configuration saved to: %s", path)
	return nil
}

// GetClusterConfig returns configuration for a specific cluster
func GetClusterConfig(name string) (ClusterConfig, bool) {
	if Config == nil || Config.Clusters == nil {
		return ClusterConfig{}, false
	}
	
	cluster, exists := Config.Clusters[name]
	return cluster, exists
}

// SetClusterConfig sets configuration for a specific cluster
func SetClusterConfig(name string, config ClusterConfig) {
	if Config == nil {
		LoadConfig() // Initialize if not loaded
	}
	
	if Config.Clusters == nil {
		Config.Clusters = make(map[string]ClusterConfig)
	}
	
	Config.Clusters[name] = config
}

// validateConfig validates the loaded configuration
func validateConfig() error {
	if Config == nil {
		return fmt.Errorf("configuration is nil")
	}

	// Validate output format
	if Config.Output.Format != "" {
		if _, err := ValidateOutputFormat(Config.Output.Format); err != nil {
			return fmt.Errorf("invalid output format in config: %w", err)
		}
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if Config.LogLevel != "" {
		found := false
		for _, level := range validLogLevels {
			if Config.LogLevel == level {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid log level in config: %s (must be one of: %v)", Config.LogLevel, validLogLevels)
		}
	}

	return nil
}

// expandPath expands environment variables and home directory in paths
func expandPath(path string) string {
	if strings.HasPrefix(path, "$HOME") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return strings.Replace(path, "$HOME", home, 1)
	}
	return os.ExpandEnv(path)
}

// InitializeFromConfig applies configuration values to global variables
func InitializeFromConfig() {
	if Config == nil {
		return
	}

	// Apply output settings
	if Config.Output.Color {
		ColoredOutput = true
	}
	
	if Config.Output.Verbose {
		Verbose = true
	}

	// Apply log level
	if Config.LogLevel != "" {
		LogLevel = Config.LogLevel
	}
}