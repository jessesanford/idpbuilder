// Copyright 2024 The IDP Builder Authors
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

/*
Package registry-auth-types provides authentication types and utilities for container registry access.

This package defines core interfaces and types for authenticating with container registries,
including Docker Hub, Harbor, Quay, and other OCI-compliant registries. It supports multiple
authentication mechanisms including HTTP Basic authentication, Bearer tokens, and OAuth2 flows.

# Authentication Types

The package supports the following authentication types:

  - Basic: HTTP Basic authentication using username/password
  - Bearer: Bearer token authentication
  - OAuth2: OAuth2 token-based authentication  
  - Anonymous: Unauthenticated access for public registries

# Core Interfaces

The RegistryAuth interface defines the contract for authentication mechanisms:

	type RegistryAuth interface {
		GetCredentials(registry string) (*Credentials, error)
		Validate() error
		Type() AuthType
		Refresh() error
		IsExpired() bool
	}

The AuthStore interface defines credential storage operations:

	type AuthStore interface {
		Get(registry string) (*Credentials, error)
		Set(registry string, creds *Credentials) error
		Delete(registry string) error
		List() ([]string, error)
		Clear() error
	}

# Usage Examples

## Basic Authentication

	// Create basic auth for a registry
	auth := auth.NewBasicAuth("registry.example.com", "username", "password")
	
	// Get credentials
	creds, err := auth.GetCredentials("registry.example.com")
	if err != nil {
		log.Fatal(err)
	}
	
	// Use credentials in HTTP requests
	authHeader := creds.ToAuthHeader()
	req.Header.Set("Authorization", authHeader)

## Bearer Token Authentication

	// Create bearer token auth
	expiresAt := time.Now().Add(1 * time.Hour)
	auth := auth.NewBearerAuth("registry.example.com", "eyJhbGciOiJIUzI1NiIs...", &expiresAt)
	
	// Check if token is expired
	if auth.IsExpired() {
		log.Println("Token has expired")
	}

## Credential Store

	// Create an in-memory credential store
	store := auth.NewCredentialStore()
	
	// Store credentials
	creds := &auth.Credentials{
		Username: "user",
		Password: "pass",
	}
	store.Set("registry.example.com", creds)
	
	// Retrieve credentials
	stored, err := store.Get("registry.example.com")
	if err != nil {
		log.Fatal(err)
	}

## Docker Config Compatibility

The package provides Docker config.json compatibility through the DockerConfig type:

	config := &auth.DockerConfig{
		Auths: map[string]auth.DockerAuthEntry{
			"registry.example.com": {
				Username: "user",
				Password: "pass",
				Auth:     "dXNlcjpwYXNz", // base64 encoded user:pass
			},
		},
		CredStore: "desktop",
		CredHelpers: map[string]string{
			"gcr.io": "gcr",
		},
	}

# Security Considerations

This package implements several security best practices:

1. **Credential Redaction**: Sensitive data is automatically redacted in logs and string representations
2. **Token Expiration**: Built-in support for token expiration checking and refresh
3. **Validation**: Comprehensive validation of credential formats and requirements
4. **Memory Safety**: Secure handling of credentials in memory with explicit clearing

# Error Handling

The package uses structured error types for better error handling:

	if err != nil {
		var authErr *auth.AuthError
		if errors.As(err, &authErr) {
			switch authErr.Type {
			case auth.ErrTokenExpired:
				// Handle token expiration
			case auth.ErrInvalidCredentials:
				// Handle invalid credentials
			default:
				// Handle other auth errors
			}
		}
	}

# Thread Safety

The CredentialStore implementation is not thread-safe. For concurrent access,
external synchronization is required or use a thread-safe implementation.

# Integration

This package is designed to integrate with:

  - Docker daemon authentication
  - Kubernetes image pull secrets
  - OCI registry clients
  - CI/CD pipeline authentication
  - Harbor registry authentication
  - Quay.io authentication

For advanced features like certificate management and TLS configuration,
see the accompanying cert package.
*/
package main