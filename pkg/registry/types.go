// Copyright 2024 idpbuilder Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"time"
)

// RegistryConfig contains registry connection settings.
// Follows idpbuilder configuration patterns with JSON tags for serialization.
type RegistryConfig struct {
	// URL is the registry base URL (e.g., "https://registry.example.com")
	URL string `json:"url,omitempty"`

	// Insecure allows insecure TLS connections for development environments
	Insecure bool `json:"insecure,omitempty"`

	// SkipTLSVerify skips TLS certificate verification
	SkipTLSVerify bool `json:"skipTLSVerify,omitempty"`

	// Auth contains authentication configuration
	Auth AuthConfig `json:"auth,omitempty"`

	// Timeout specifies the default timeout for registry operations
	Timeout time.Duration `json:"timeout,omitempty"`

	// RetryAttempts specifies the number of retry attempts for failed operations
	RetryAttempts int `json:"retryAttempts,omitempty"`
}

// AuthConfig contains authentication settings for registry access.
// Supports both token-based and basic authentication methods.
type AuthConfig struct {
	// Token for token-based authentication (preferred for Gitea)
	Token string `json:"token,omitempty"`

	// Username for basic authentication
	Username string `json:"username,omitempty"`

	// Password for basic authentication
	Password string `json:"password,omitempty"`

	// AuthMethod specifies the authentication method to use
	AuthMethod AuthMethod `json:"authMethod,omitempty"`
}

// AuthMethod defines the type of authentication to use.
type AuthMethod string

const (
	// AuthMethodToken uses token-based authentication
	AuthMethodToken AuthMethod = "token"

	// AuthMethodBasic uses basic username/password authentication
	AuthMethodBasic AuthMethod = "basic"

	// AuthMethodNone indicates no authentication is required
	AuthMethodNone AuthMethod = "none"
)

// DefaultRegistryConfig returns a RegistryConfig with sensible defaults.
func DefaultRegistryConfig() RegistryConfig {
	return RegistryConfig{
		Insecure:      false,
		SkipTLSVerify: false,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		Auth: AuthConfig{
			AuthMethod: AuthMethodNone,
		},
	}
}

// Validate validates the registry configuration and returns any errors.
func (c *RegistryConfig) Validate() error {
	if c.URL == "" {
		return NewConfigError("registry URL is required")
	}

	if c.Auth.AuthMethod == AuthMethodToken && c.Auth.Token == "" {
		return NewConfigError("token is required when using token authentication")
	}

	if c.Auth.AuthMethod == AuthMethodBasic {
		if c.Auth.Username == "" || c.Auth.Password == "" {
			return NewConfigError("username and password are required when using basic authentication")
		}
	}

	return nil
}

// IsSecure returns true if the registry is configured for secure connections.
func (c *RegistryConfig) IsSecure() bool {
	return !c.Insecure && !c.SkipTLSVerify
}