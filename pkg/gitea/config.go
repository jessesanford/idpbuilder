package gitea

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the configuration file structure
type Config struct {
	Registries map[string]RegistryCredentials `json:"registries"`
}

// RegistryCredentials holds credentials for a specific registry
type RegistryCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

// ConfigFileProvider reads from ~/.idpbuilder/config
type ConfigFileProvider struct {
	configPath string
	config     *Config
}

func NewConfigFileProvider() *ConfigFileProvider {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".idpbuilder", "config")
	return &ConfigFileProvider{
		configPath: configPath,
	}
}

func (c *ConfigFileProvider) loadConfig() error {
	if c.config != nil {
		return nil
	}

	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	c.config = &config
	return nil
}

func (c *ConfigFileProvider) GetUsername() (string, error) {
	if err := c.loadConfig(); err != nil {
		return "", err
	}

	// Look for gitea registry config
	for name, creds := range c.config.Registries {
		if name == "gitea" || name == "default" {
			return creds.Username, nil
		}
	}

	return "", fmt.Errorf("no gitea credentials in config")
}

func (c *ConfigFileProvider) GetPassword() (string, error) {
	if err := c.loadConfig(); err != nil {
		return "", err
	}

	for name, creds := range c.config.Registries {
		if name == "gitea" || name == "default" {
			return creds.Password, nil
		}
	}

	return "", fmt.Errorf("no gitea credentials in config")
}

func (c *ConfigFileProvider) IsAvailable() bool {
	if _, err := os.Stat(c.configPath); err != nil {
		return false
	}
	return true
}

func (c *ConfigFileProvider) Priority() int {
	return 3 // Third priority
}