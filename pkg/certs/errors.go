// pkg/certs/errors.go
package certs

import "fmt"

// ClusterNotFoundError indicates that the Kind cluster could not be found
type ClusterNotFoundError struct {
	ClusterName string
}

func (e ClusterNotFoundError) Error() string {
	return fmt.Sprintf("Kind cluster not found: %s", e.ClusterName)
}

// PodNotFoundError indicates that the specified pod could not be found
type PodNotFoundError struct {
	PodName   string
	Namespace string
}

func (e PodNotFoundError) Error() string {
	return fmt.Sprintf("pod '%s' not found in namespace '%s'", e.PodName, e.Namespace)
}

// CertificateInvalidError indicates that the certificate is invalid or cannot be processed
type CertificateInvalidError struct {
	Reason string
}

func (e CertificateInvalidError) Error() string {
	return fmt.Sprintf("certificate invalid: %s", e.Reason)
}

// PermissionError indicates insufficient permissions for file operations
type PermissionError struct {
	Path   string
	Action string
}

func (e PermissionError) Error() string {
	return fmt.Sprintf("permission denied: cannot %s file at path %s", e.Action, e.Path)
}