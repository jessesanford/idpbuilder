package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// RegistryConfig represents the main registry configuration
type RegistryConfig struct {
	URL     string            `yaml:"url" json:"url"`
	Type    string            `yaml:"type" json:"type"`
	Insecure bool             `yaml:"insecure" json:"insecure"`
	Auth    AuthConfig        `yaml:"auth" json:"auth"`
	Options map[string]string `yaml:"options" json:"options"`
}

// ConnectionConfig represents registry connection pool configuration
type ConnectionConfig struct {
	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
	Timeout         time.Duration `yaml:"timeout" json:"timeout"`
}

// LoadRegistryConfig loads registry configuration from a file with environment variable overrides
func LoadRegistryConfig(path string) (*RegistryConfig, error) {
	if path == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	config := &RegistryConfig{}
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	// Apply environment variable overrides
	if url := os.Getenv("REGISTRY_URL"); url != "" {
		config.URL = url
	}
	if regType := os.Getenv("REGISTRY_TYPE"); regType != "" {
		config.Type = regType
	}
	if insecure := os.Getenv("REGISTRY_INSECURE"); insecure != "" {
		if val, err := strconv.ParseBool(insecure); err == nil {
			config.Insecure = val
		}
	}

	return config, nil
}

// ToConnectionString builds a connection string from registry configuration
func ToConnectionString(config *RegistryConfig) string {
	if config == nil {
		return ""
	}

	scheme := "https"
	if config.Insecure {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%s", scheme, config.URL)
}