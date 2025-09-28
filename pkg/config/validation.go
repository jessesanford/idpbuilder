package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// ValidationResult represents the result of configuration validation
type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// ValidateRegistryConfig validates a registry configuration and returns detailed results
func ValidateRegistryConfig(config *RegistryConfig) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	if config == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "config cannot be nil")
		return result
	}

	// Validate URL
	if config.URL == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "URL is required")
	} else {
		if _, err := url.Parse(config.URL); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("invalid URL format: %v", err))
		}
	}

	// Validate registry type
	validTypes := []string{"dockerhub", "harbor", "gitea", "generic"}
	if config.Type != "" {
		typeValid := false
		for _, validType := range validTypes {
			if config.Type == validType {
				typeValid = true
				break
			}
		}
		if !typeValid {
			result.Warnings = append(result.Warnings, fmt.Sprintf("unknown registry type '%s', supported types: %s", config.Type, strings.Join(validTypes, ", ")))
		}
	}

	// Validate insecure connections
	if config.Insecure {
		result.Warnings = append(result.Warnings, "insecure connections enabled - not recommended for production")
	}

	// Validate auth configuration
	if config.Auth.Type != "" {
		switch config.Auth.Type {
		case "basic":
			if config.Auth.Username == "" {
				result.Valid = false
				result.Errors = append(result.Errors, "basic auth requires username")
			}
			if config.Auth.Password == "" {
				result.Valid = false
				result.Errors = append(result.Errors, "basic auth requires password")
			}
		case "token", "oauth2":
			if config.Auth.Token == "" {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("%s auth requires token", config.Auth.Type))
			}
		default:
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("unsupported auth type: %s", config.Auth.Type))
		}
	}

	return result
}

// SetDefaultValues applies default values to registry configuration
func SetDefaultValues(config *RegistryConfig) {
	if config == nil {
		return
	}

	// Set default registry type
	if config.Type == "" {
		config.Type = "generic"
	}

	// Set default URL if not provided
	if config.URL == "" {
		if defaultURL := os.Getenv("REGISTRY_URL"); defaultURL != "" {
			config.URL = defaultURL
		} else {
			config.URL = "docker.io"
		}
	}

	// Initialize options map if nil
	if config.Options == nil {
		config.Options = make(map[string]string)
	}

	// Set default timeout if not specified in options
	if _, exists := config.Options["timeout"]; !exists {
		if timeoutStr := os.Getenv("REGISTRY_TIMEOUT"); timeoutStr != "" {
			config.Options["timeout"] = timeoutStr
		} else {
			config.Options["timeout"] = "30s"
		}
	}

	// Set default connection limits if not specified
	if _, exists := config.Options["max_connections"]; !exists {
		if maxConn := os.Getenv("REGISTRY_MAX_CONNECTIONS"); maxConn != "" {
			config.Options["max_connections"] = maxConn
		} else {
			config.Options["max_connections"] = "10"
		}
	}
}

// MergeConfigs merges two registry configurations, with override taking precedence
func MergeConfigs(base, override *RegistryConfig) *RegistryConfig {
	if base == nil && override == nil {
		return nil
	}
	if base == nil {
		result := *override
		return &result
	}
	if override == nil {
		result := *base
		return &result
	}

	// Start with a copy of base
	result := *base

	// Override with non-empty values from override
	if override.URL != "" {
		result.URL = override.URL
	}
	if override.Type != "" {
		result.Type = override.Type
	}
	// Always use override's insecure setting if override config exists
	result.Insecure = override.Insecure

	// Merge auth configuration
	if override.Auth.Type != "" {
		result.Auth = override.Auth
	} else {
		// Merge individual auth fields if override auth type is empty
		if override.Auth.Username != "" {
			result.Auth.Username = override.Auth.Username
		}
		if override.Auth.Password != "" {
			result.Auth.Password = override.Auth.Password
		}
		if override.Auth.Token != "" {
			result.Auth.Token = override.Auth.Token
		}
		if override.Auth.ConfigFile != "" {
			result.Auth.ConfigFile = override.Auth.ConfigFile
		}
	}

	// Merge options
	if result.Options == nil {
		result.Options = make(map[string]string)
	}
	for key, value := range override.Options {
		result.Options[key] = value
	}

	return &result
}

// getIntEnv gets an integer value from environment variable with default fallback
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv gets a duration value from environment variable with default fallback
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}