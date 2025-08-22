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
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client is the concrete implementation of the RegistryClient interface.
type Client struct {
	config        RegistryConfig
	authenticator *Authenticator
	httpClient    *http.Client
	authenticated bool
}

// NewClient creates a new registry client with the given configuration.
func NewClient(config RegistryConfig) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	client := &Client{
		config:        config,
		authenticator: NewAuthenticator(config),
		httpClient:    HTTPClient(config),
		authenticated: false,
	}

	return client, nil
}

// Push pushes a container image to the registry with the given options.
func (c *Client) Push(ctx context.Context, imageRef string, opts PushOptions) error {
	if !c.authenticated && c.config.Auth.AuthMethod != AuthMethodNone {
		return NewAuthError("client is not authenticated", nil)
	}

	// Merge options with defaults
	if opts.Timeout == 0 {
		opts.Timeout = c.config.Timeout
	}
	if opts.Retries == 0 {
		opts.Retries = c.config.RetryAttempts
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	return c.pushWithRetries(ctx, imageRef, opts)
}

// pushWithRetries executes the push operation with retry logic.
func (c *Client) pushWithRetries(ctx context.Context, imageRef string, opts PushOptions) error {
	var lastErr error

	for attempt := 0; attempt <= opts.Retries; attempt++ {
		if attempt > 0 {
			// Wait before retry with exponential backoff
			backoff := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return NewTimeoutError("push operation timeout", ctx.Err())
			case <-time.After(backoff):
			}
		}

		err := c.pushImage(ctx, imageRef, opts)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if regErr, ok := err.(*RegistryError); ok && !regErr.IsRetryable() {
			break // Don't retry non-retryable errors
		}

		// Check if context was cancelled
		if ctx.Err() != nil {
			return NewTimeoutError("push operation timeout", ctx.Err())
		}
	}

	return lastErr
}

// pushImage performs the actual image push operation.
// This is a simplified implementation that focuses on the authentication and error handling patterns.
func (c *Client) pushImage(ctx context.Context, imageRef string, opts PushOptions) error {
	// For this MVP implementation, we simulate a push operation
	// In a real implementation, this would use a proper container registry client library
	
	url := c.buildPushURL(imageRef)
	req, err := c.authenticator.CreateAuthenticatedRequest(ctx, http.MethodPut, url)
	if err != nil {
		return err
	}

	// Set appropriate headers for registry operations
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("User-Agent", "idpbuilder/registry-client")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ClassifyError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return NewRegistryError(fmt.Sprintf("push failed with status %d", resp.StatusCode), nil)
	}

	return nil
}

// buildPushURL constructs the registry URL for pushing the given image reference.
func (c *Client) buildPushURL(imageRef string) string {
	// Extract repository and tag from image reference
	parts := strings.Split(imageRef, ":")
	repository := parts[0]
	tag := "latest"
	if len(parts) > 1 {
		tag = parts[1]
	}

	// Remove any registry prefix if it matches our configured registry
	if strings.Contains(repository, "/") {
		parts := strings.Split(repository, "/")
		if len(parts) > 1 {
			// Use the repository path without the registry prefix
			repository = strings.Join(parts[1:], "/")
		}
	}

	baseURL := strings.TrimSuffix(c.config.URL, "/")
	return fmt.Sprintf("%s/v2/%s/manifests/%s", baseURL, repository, tag)
}

// Authenticate authenticates with the registry using the configured authentication method.
func (c *Client) Authenticate(ctx context.Context, config AuthConfig) error {
	// Update the client's auth config
	c.config.Auth = config
	c.authenticator = NewAuthenticator(c.config)

	err := c.authenticator.Authenticate(ctx)
	if err != nil {
		return err
	}

	c.authenticated = true
	return nil
}

// Health checks registry connectivity and authentication status.
func (c *Client) Health(ctx context.Context) error {
	url := strings.TrimSuffix(c.config.URL, "/") + "/v2/"
	req, err := c.authenticator.CreateAuthenticatedRequest(ctx, http.MethodGet, url)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NewNetworkError("registry health check failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return NewRegistryError(fmt.Sprintf("registry health check failed with status %d", resp.StatusCode), nil)
	}

	return nil
}