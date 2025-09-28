package config

import (
	"fmt"
	"os"
)

// AuthConfig represents authentication configuration for registry access
type AuthConfig struct {
	Type       string `yaml:"type" json:"type"`             // auth type: basic, token, oauth2
	Username   string `yaml:"username" json:"username"`     // username for basic auth
	Password   string `yaml:"password" json:"password"`     // password for basic auth (from env/secret)
	Token      string `yaml:"token" json:"token"`           // bearer token
	ConfigFile string `yaml:"config_file" json:"config_file"` // path to docker config
}

// TLSConfig represents TLS configuration for secure registry connections
type TLSConfig struct {
	CACert     string `yaml:"ca_cert" json:"ca_cert"`         // CA certificate path
	ClientCert string `yaml:"client_cert" json:"client_cert"` // Client certificate path
	ClientKey  string `yaml:"client_key" json:"client_key"`   // Client key path
	SkipVerify bool   `yaml:"skip_verify" json:"skip_verify"` // Skip TLS verification (dev only)
}

// GetAuthConfig extracts and resolves authentication configuration with environment variable support
func GetAuthConfig(config *RegistryConfig) (*AuthConfig, error) {
	if config == nil {
		return nil, fmt.Errorf("registry config cannot be nil")
	}

	auth := config.Auth

	// Resolve environment variables for sensitive fields
	if auth.Username == "" {
		if username := os.Getenv("REGISTRY_USERNAME"); username != "" {
			auth.Username = username
		}
	}

	if auth.Password == "" {
		if password := os.Getenv("REGISTRY_PASSWORD"); password != "" {
			auth.Password = password
		}
	}

	if auth.Token == "" {
		if token := os.Getenv("REGISTRY_TOKEN"); token != "" {
			auth.Token = token
		}
	}

	if auth.ConfigFile == "" {
		if configFile := os.Getenv("DOCKER_CONFIG"); configFile != "" {
			auth.ConfigFile = configFile
		}
	}

	// Validate auth configuration
	if auth.Type == "" {
		return &auth, nil // No auth configured
	}

	switch auth.Type {
	case "basic":
		if auth.Username == "" || auth.Password == "" {
			return nil, fmt.Errorf("basic auth requires both username and password")
		}
	case "token":
		if auth.Token == "" {
			return nil, fmt.Errorf("token auth requires a token")
		}
	case "oauth2":
		if auth.Token == "" {
			return nil, fmt.Errorf("oauth2 auth requires a token")
		}
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", auth.Type)
	}

	return &auth, nil
}