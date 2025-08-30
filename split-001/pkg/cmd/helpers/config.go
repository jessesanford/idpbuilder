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
	Config *AppConfig
	ConfigFile string
	DefaultConfigPaths = []string{
		"./idpbuilder.yaml",
		"$HOME/.idpbuilder/config.yaml",
		"/etc/idpbuilder/config.yaml",
	}
)

type AppConfig struct {
	Version  string              `yaml:"version" json:"version"`
	LogLevel string              `yaml:"logLevel" json:"logLevel"`
	Output   OutputConfig        `yaml:"output" json:"output"`
	Defaults DefaultsConfig      `yaml:"defaults" json:"defaults"`
	Clusters map[string]ClusterConfig `yaml:"clusters" json:"clusters"`
}

type OutputConfig struct {
	Format  string `yaml:"format" json:"format"`
	Color   bool   `yaml:"color" json:"color"`
	Verbose bool   `yaml:"verbose" json:"verbose"`
}

type DefaultsConfig struct {
	Namespace string `yaml:"namespace" json:"namespace"`
	Timeout   string `yaml:"timeout" json:"timeout"`
	DryRun    bool   `yaml:"dryRun" json:"dryRun"`
}

type ClusterConfig struct {
	Name      string `yaml:"name" json:"name"`
	Context   string `yaml:"context" json:"context"`
	Namespace string `yaml:"namespace" json:"namespace"`
	URL       string `yaml:"url" json:"url"`
}

func LoadConfig() error {
	Config = &AppConfig{
		Version:  "v1",
		LogLevel: "info",
		Output:   OutputConfig{Format: "table", Color: false, Verbose: false},
		Defaults: DefaultsConfig{Namespace: "default", Timeout: "30s", DryRun: false},
		Clusters: make(map[string]ClusterConfig),
	}

	if ConfigFile != "" {
		return loadConfigFromFile(ConfigFile)
	}

	for _, path := range DefaultConfigPaths {
		expandedPath := expandPath(path)
		if _, err := os.Stat(expandedPath); err == nil {
			ConfigFile = expandedPath
			return loadConfigFromFile(expandedPath)
		}
	}
	return nil
}

func loadConfigFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".json" {
		err = json.Unmarshal(data, Config)
	} else {
		err = yaml.Unmarshal(data, Config)
	}

	if err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", path, err)
	}
	return validateConfig()
}

func SaveConfig(path string) error {
	if Config == nil {
		return fmt.Errorf("no configuration to save")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	var data []byte
	var err error
	if strings.ToLower(filepath.Ext(path)) == ".json" {
		data, err = json.MarshalIndent(Config, "", "  ")
	} else {
		data, err = yaml.Marshal(Config)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func GetClusterConfig(name string) (ClusterConfig, bool) {
	if Config == nil || Config.Clusters == nil {
		return ClusterConfig{}, false
	}
	cluster, exists := Config.Clusters[name]
	return cluster, exists
}

func SetClusterConfig(name string, config ClusterConfig) {
	if Config == nil {
		LoadConfig()
	}
	if Config.Clusters == nil {
		Config.Clusters = make(map[string]ClusterConfig)
	}
	Config.Clusters[name] = config
}

func validateConfig() error {
	if Config == nil {
		return fmt.Errorf("configuration is nil")
	}
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if Config.LogLevel != "" {
		for _, level := range validLogLevels {
			if Config.LogLevel == level {
				return nil
			}
		}
		return fmt.Errorf("invalid log level: %s", Config.LogLevel)
	}
	return nil
}

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

func InitializeFromConfig() {
	if Config == nil {
		return
	}
	if Config.Output.Color {
		ColoredOutput = true
	}
	if Config.Output.Verbose {
		Verbose = true
	}
	if Config.LogLevel != "" {
		LogLevel = Config.LogLevel
	}
}