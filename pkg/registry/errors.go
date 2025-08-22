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
	"fmt"
	"net"
	"net/url"
)

// ErrorType defines the category of registry error.
type ErrorType string

const (
	// ErrorTypeAuth indicates an authentication or authorization error
	ErrorTypeAuth ErrorType = "authentication"

	// ErrorTypeNetwork indicates a network connectivity error
	ErrorTypeNetwork ErrorType = "network"

	// ErrorTypeConfig indicates a configuration error
	ErrorTypeConfig ErrorType = "configuration"

	// ErrorTypeRegistry indicates a registry-specific error
	ErrorTypeRegistry ErrorType = "registry"

	// ErrorTypeTimeout indicates a timeout error
	ErrorTypeTimeout ErrorType = "timeout"

	// ErrorTypeRetry indicates a retryable error
	ErrorTypeRetry ErrorType = "retry"
)

// RegistryError represents a domain-specific registry error.
type RegistryError struct {
	Type     ErrorType
	Message  string
	Cause    error
	Retryable bool
}

// Error implements the error interface.
func (e *RegistryError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s error: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error for error chain compatibility.
func (e *RegistryError) Unwrap() error {
	return e.Cause
}

// IsRetryable returns true if the error is retryable.
func (e *RegistryError) IsRetryable() bool {
	return e.Retryable
}

// NewAuthError creates a new authentication error.
func NewAuthError(message string, cause error) *RegistryError {
	return &RegistryError{
		Type:      ErrorTypeAuth,
		Message:   message,
		Cause:     cause,
		Retryable: false,
	}
}

// NewNetworkError creates a new network error.
func NewNetworkError(message string, cause error) *RegistryError {
	return &RegistryError{
		Type:      ErrorTypeNetwork,
		Message:   message,
		Cause:     cause,
		Retryable: true,
	}
}

// NewConfigError creates a new configuration error.
func NewConfigError(message string) *RegistryError {
	return &RegistryError{
		Type:      ErrorTypeConfig,
		Message:   message,
		Cause:     nil,
		Retryable: false,
	}
}

// NewRegistryError creates a new registry-specific error.
func NewRegistryError(message string, cause error) *RegistryError {
	return &RegistryError{
		Type:      ErrorTypeRegistry,
		Message:   message,
		Cause:     cause,
		Retryable: false,
	}
}

// NewTimeoutError creates a new timeout error.
func NewTimeoutError(message string, cause error) *RegistryError {
	return &RegistryError{
		Type:      ErrorTypeTimeout,
		Message:   message,
		Cause:     cause,
		Retryable: true,
	}
}

// ClassifyError classifies a generic error into a registry error type.
func ClassifyError(err error) *RegistryError {
	if err == nil {
		return nil
	}

	// Already a registry error
	if regErr, ok := err.(*RegistryError); ok {
		return regErr
	}

	// Network errors
	if _, ok := err.(net.Error); ok {
		return NewNetworkError("network operation failed", err)
	}
	if _, ok := err.(*net.OpError); ok {
		return NewNetworkError("network operation failed", err)
	}
	if _, ok := err.(*url.Error); ok {
		return NewNetworkError("URL operation failed", err)
	}

	// Default to registry error
	return NewRegistryError("registry operation failed", err)
}