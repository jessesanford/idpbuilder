package docker

import "fmt"

// DaemonConnectionError indicates the Docker daemon is unreachable or not running.
type DaemonConnectionError struct {
	Cause error
}

func (e *DaemonConnectionError) Error() string {
	return fmt.Sprintf("Docker daemon connection error: %v", e.Cause)
}

func (e *DaemonConnectionError) Unwrap() error {
	return e.Cause
}

// ImageNotFoundError indicates the requested image does not exist in the Docker daemon.
type ImageNotFoundError struct {
	ImageName string
}

func (e *ImageNotFoundError) Error() string {
	return fmt.Sprintf("image '%s' not found in Docker daemon", e.ImageName)
}

// ImageConversionError indicates failure to convert Docker image to OCI format.
type ImageConversionError struct {
	ImageName string
	Cause     error
}

func (e *ImageConversionError) Error() string {
	return fmt.Sprintf("failed to convert image '%s' to OCI format: %v", e.ImageName, e.Cause)
}

func (e *ImageConversionError) Unwrap() error {
	return e.Cause
}

// ValidationError indicates image name validation failed.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error (%s): %s", e.Field, e.Message)
}
