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
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// Authenticator handles authentication with the registry.
type Authenticator struct {
	client *http.Client
	config RegistryConfig
}

// NewAuthenticator creates a new authenticator with the given configuration.
func NewAuthenticator(config RegistryConfig) *Authenticator {
	return &Authenticator{
		client: HTTPClient(config),
		config: config,
	}
}

// Authenticate performs authentication with the registry based on the configured method.
func (a *Authenticator) Authenticate(ctx context.Context) error {
	switch a.config.Auth.AuthMethod {
	case AuthMethodNone:
		return nil
	case AuthMethodToken:
		return a.authenticateWithToken(ctx)
	case AuthMethodBasic:
		return a.authenticateWithBasic(ctx)
	default:
		return NewAuthError("unsupported authentication method", nil)
	}
}

// authenticateWithToken performs token-based authentication.
// This is the preferred method for Gitea registry authentication.
func (a *Authenticator) authenticateWithToken(ctx context.Context) error {
	if a.config.Auth.Token == "" {
		return NewAuthError("token is required for token authentication", nil)
	}

	// Test token authentication by making a request to the registry
	url := strings.TrimSuffix(a.config.URL, "/") + "/v2/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NewAuthError("failed to create authentication request", err)
	}

	// Set authorization header with bearer token
	req.Header.Set("Authorization", "Bearer "+a.config.Auth.Token)

	resp, err := a.client.Do(req)
	if err != nil {
		return ClassifyError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return NewAuthError("token authentication failed: invalid token", nil)
	}

	if resp.StatusCode >= 400 {
		return NewAuthError(fmt.Sprintf("token authentication failed with status %d", resp.StatusCode), nil)
	}

	return nil
}

// authenticateWithBasic performs basic authentication with username and password.
func (a *Authenticator) authenticateWithBasic(ctx context.Context) error {
	if a.config.Auth.Username == "" || a.config.Auth.Password == "" {
		return NewAuthError("username and password are required for basic authentication", nil)
	}

	// Test basic authentication by making a request to the registry
	url := strings.TrimSuffix(a.config.URL, "/") + "/v2/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NewAuthError("failed to create authentication request", err)
	}

	// Set basic authentication header
	auth := base64.StdEncoding.EncodeToString([]byte(a.config.Auth.Username + ":" + a.config.Auth.Password))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := a.client.Do(req)
	if err != nil {
		return ClassifyError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return NewAuthError("basic authentication failed: invalid credentials", nil)
	}

	if resp.StatusCode >= 400 {
		return NewAuthError(fmt.Sprintf("basic authentication failed with status %d", resp.StatusCode), nil)
	}

	return nil
}

// CreateAuthenticatedRequest creates an HTTP request with appropriate authentication headers.
func (a *Authenticator) CreateAuthenticatedRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, NewAuthError("failed to create request", err)
	}

	switch a.config.Auth.AuthMethod {
	case AuthMethodToken:
		req.Header.Set("Authorization", "Bearer "+a.config.Auth.Token)
	case AuthMethodBasic:
		auth := base64.StdEncoding.EncodeToString([]byte(a.config.Auth.Username + ":" + a.config.Auth.Password))
		req.Header.Set("Authorization", "Basic "+auth)
	case AuthMethodNone:
		// No authentication headers needed
	}

	return req, nil
}