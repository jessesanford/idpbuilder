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
	"time"
)

// RegistryClient defines the interface for container registry operations.
// This interface enables mocking and testing while providing a clean abstraction
// for registry interactions.
type RegistryClient interface {
	// Push pushes a container image to the registry with the given options.
	// Returns an error if the push fails or if authentication is required.
	Push(ctx context.Context, imageRef string, opts PushOptions) error

	// Authenticate authenticates with the registry using the provided configuration.
	// Must be called before attempting push operations that require authentication.
	Authenticate(ctx context.Context, config AuthConfig) error

	// Health checks registry connectivity and authentication status.
	// Useful for validating configuration before attempting operations.
	Health(ctx context.Context) error
}

// PushOptions contains options for push operations.
type PushOptions struct {
	// Tags specifies additional tags to apply to the image during push
	Tags []string `json:"tags,omitempty"`

	// Insecure allows insecure (non-TLS) connections to the registry
	Insecure bool `json:"insecure,omitempty"`

	// Timeout specifies the maximum time to wait for the push operation
	Timeout time.Duration `json:"timeout,omitempty"`

	// Retries specifies the number of retry attempts on transient failures
	Retries int `json:"retries,omitempty"`
}

// DefaultPushOptions returns PushOptions with sensible defaults.
func DefaultPushOptions() PushOptions {
	return PushOptions{
		Tags:     []string{},
		Insecure: false,
		Timeout:  5 * time.Minute,
		Retries:  3,
	}
}