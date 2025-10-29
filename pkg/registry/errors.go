package registry

import "fmt"

// AuthenticationError indicates registry authentication failed (401/403).
type AuthenticationError struct {
	Registry string
	Cause    error
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("registry authentication failed for %s: %v", e.Registry, e.Cause)
}

func (e *AuthenticationError) Unwrap() error {
	return e.Cause
}

// NetworkError indicates network connectivity issues with the registry.
type NetworkError struct {
	Registry string
	Cause    error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error connecting to registry %s: %v", e.Registry, e.Cause)
}

func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// RegistryUnavailableError indicates the registry endpoint is not responding correctly.
type RegistryUnavailableError struct {
	Registry   string
	StatusCode int
}

func (e *RegistryUnavailableError) Error() string {
	return fmt.Sprintf("registry %s unavailable (status: %d)", e.Registry, e.StatusCode)
}

// PushFailedError indicates image push operation failed.
type PushFailedError struct {
	TargetRef string
	Cause     error
}

func (e *PushFailedError) Error() string {
	return fmt.Sprintf("push to %s failed: %v", e.TargetRef, e.Cause)
}

func (e *PushFailedError) Unwrap() error {
	return e.Cause
}
