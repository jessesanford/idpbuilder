package certificates

import (
	"context"
	"time"
)

// Config represents the configuration for the certificate management system.
type Config struct {
	StorageBackend  *StorageConfig    `json:"storage_backend"`
	ValidationRules []*ValidationRule `json:"validation_rules"`
	DefaultPools    []string          `json:"default_pools"`
	EventHandlers   *EventConfig      `json:"event_handlers"`
}

// StorageConfig defines storage backend configuration.
type StorageConfig struct {
	Type             string                 `json:"type"`
	ConnectionString string                 `json:"connection_string"`
	Options          map[string]interface{} `json:"options,omitempty"`
}

// EventConfig defines event handling configuration.
type EventConfig struct {
	Enabled    bool     `json:"enabled"`
	BufferSize int      `json:"buffer_size"`
	Handlers   []string `json:"handlers"`
}

// ValidationRule represents a validation rule.
type ValidationRule struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Critical    bool                   `json:"critical"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ValidationResult represents the result of validating a certificate.
type ValidationResult struct {
	Valid       bool              `json:"valid"`
	ValidatedAt time.Time         `json:"validated_at"`
	Errors      []ValidationError `json:"errors,omitempty"`
}

// ValidationError represents errors during validation.
type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ConfigLoader defines the interface for loading configuration data.
type ConfigLoader interface {
	LoadConfig(ctx context.Context, source string) (*Config, error)
	SaveConfig(ctx context.Context, config *Config, destination string) error
	ValidateConfig(ctx context.Context, config *Config) (*ValidationResult, error)
	GetDefaultConfig(ctx context.Context) (*Config, error)
}