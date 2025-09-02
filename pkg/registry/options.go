package registry

import (
	"fmt"
	"time"
)

// ClientOptions contains configuration for the Gitea client
type ClientOptions struct {
	BaseURL         string
	Username        string
	Password        string
	Insecure        bool
	Timeout         time.Duration
	MaxRetries      int
	RetryDelay      time.Duration
	UserAgent       string
	
	// Transport configuration
	TransportConfig *TransportConfig
}

// DefaultClientOptions returns default client options for Gitea
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		BaseURL:         "https://gitea.cnoe.localtest.me:443",
		Username:        "gitea_admin",
		Password:        "", // Should be provided or loaded from env
		Insecure:        false,
		Timeout:         30 * time.Second,
		MaxRetries:      3,
		RetryDelay:      1 * time.Second,
		UserAgent:       "idpbuilder-oci/1.0",
		TransportConfig: DefaultTransportConfig(),
	}
}

// Validate checks if the client options are valid
func (o ClientOptions) Validate() error {
	if o.BaseURL == "" {
		return fmt.Errorf("base URL cannot be empty")
	}
	
	if o.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	
	if o.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}
	
	if o.RetryDelay < 0 {
		return fmt.Errorf("retry delay cannot be negative")
	}
	
	if o.UserAgent == "" {
		return fmt.Errorf("user agent cannot be empty")
	}
	
	// Validate transport config if provided
	if o.TransportConfig != nil {
		if err := ValidateTransportConfig(o.TransportConfig); err != nil {
			return fmt.Errorf("invalid transport config: %w", err)
		}
	}
	
	return nil
}

// ApplyDefaults applies default values to missing fields
func (o *ClientOptions) ApplyDefaults() {
	if o.BaseURL == "" {
		o.BaseURL = "https://gitea.cnoe.localtest.me:443"
	}
	
	if o.Timeout <= 0 {
		o.Timeout = 30 * time.Second
	}
	
	if o.MaxRetries < 0 {
		o.MaxRetries = 3
	}
	
	if o.RetryDelay < 0 {
		o.RetryDelay = 1 * time.Second
	}
	
	if o.UserAgent == "" {
		o.UserAgent = "idpbuilder-oci/1.0"
	}
	
	if o.TransportConfig == nil {
		o.TransportConfig = DefaultTransportConfig()
	}
}

// Clone creates a copy of the client options
func (o ClientOptions) Clone() ClientOptions {
	clone := o
	
	// Deep copy transport config
	if o.TransportConfig != nil {
		clone.TransportConfig = CloneTransportConfig(o.TransportConfig)
	}
	
	return clone
}