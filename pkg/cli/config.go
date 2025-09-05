package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Registry RegistryConfig `mapstructure:"registry"`
	Build    BuildConfig    `mapstructure:"build"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// RegistryConfig holds registry-related configuration
type RegistryConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Insecure bool   `mapstructure:"insecure"`
	Timeout  int    `mapstructure:"timeout"`
}

// BuildConfig holds build-related configuration
type BuildConfig struct {
	Context     string   `mapstructure:"context"`
	Exclude     []string `mapstructure:"exclude"`
	CacheDir    string   `mapstructure:"cache_dir"`
	Platform    string   `mapstructure:"platform"`
	NoCache     bool     `mapstructure:"no_cache"`
	Parallelism int      `mapstructure:"parallelism"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	File   string `mapstructure:"file"`
}

// LoadConfig loads configuration from file and environment
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("IDPBUILDER")
	viper.AutomaticEnv()

	// Set defaults
	setConfigDefaults()

	// Load config file if specified
	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	} else {
		// Try to find config in standard locations
		viper.SetConfigName("idpbuilder")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.idpbuilder")
		viper.AddConfigPath("/etc/idpbuilder")

		// Read config file if found (but don't error if not found)
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read config: %w", err)
			}
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand environment variables in paths
	if err := expandConfigPaths(&config); err != nil {
		return nil, fmt.Errorf("failed to expand config paths: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		return fmt.Errorf("config path cannot be empty")
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.Set("registry", config.Registry)
	viper.Set("build", config.Build)
	viper.Set("logging", config.Logging)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// setConfigDefaults sets default configuration values
func setConfigDefaults() {
	// Registry defaults
	viper.SetDefault("registry.url", "https://gitea.cnoe.localtest.me:443")
	viper.SetDefault("registry.username", "gitea_admin")
	viper.SetDefault("registry.insecure", false)
	viper.SetDefault("registry.timeout", 30)

	// Build defaults
	viper.SetDefault("build.context", ".")
	viper.SetDefault("build.cache_dir", filepath.Join(os.Getenv("HOME"), ".idpbuilder/cache"))
	viper.SetDefault("build.platform", "linux/amd64")
	viper.SetDefault("build.no_cache", false)
	viper.SetDefault("build.parallelism", 1)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")
}

// expandConfigPaths expands environment variables in configuration paths
func expandConfigPaths(config *Config) error {
	var err error

	// Expand cache directory path
	config.Build.CacheDir, err = expandPath(config.Build.CacheDir)
	if err != nil {
		return fmt.Errorf("failed to expand cache directory path: %w", err)
	}

	// Expand log file path if specified
	if config.Logging.File != "" {
		config.Logging.File, err = expandPath(config.Logging.File)
		if err != nil {
			return fmt.Errorf("failed to expand log file path: %w", err)
		}
	}

	return nil
}

// expandPath expands environment variables in a path
func expandPath(path string) (string, error) {
	expanded := os.ExpandEnv(path)
	
	// Handle tilde expansion
	if len(expanded) > 0 && expanded[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		expanded = filepath.Join(home, expanded[1:])
	}

	return expanded, nil
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "idpbuilder.yaml"
	}
	return filepath.Join(home, ".idpbuilder", "config.yaml")
}