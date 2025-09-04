package registry

import (
	"fmt"
	"time"
)

// ClientOptions contains basic configuration for the Gitea client
type ClientOptions struct {
	BaseURL    string
	Username   string
	Password   string
	Insecure   bool
	Timeout    time.Duration
	UserAgent  string
}

// NewDefaultOptions returns default client options for Gitea
func NewDefaultOptions() ClientOptions {
	return ClientOptions{
		BaseURL:   "https://gitea.cnoe.localtest.me:443",
		Username:  "gitea_admin",
		Insecure:  false,
		Timeout:   30 * time.Second,
		UserAgent: "idpbuilder-oci/1.0",
	}
}

// Validate performs basic validation of client options
func (o ClientOptions) Validate() error {
	if o.BaseURL == "" {
		return fmt.Errorf("base URL cannot be empty")
	}
	if o.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	return nil
}