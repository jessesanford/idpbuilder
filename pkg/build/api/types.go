// Package api defines the core types for container build operations
package api

import "fmt"

// BuildRequest represents a container build request
type BuildRequest struct {
	DockerfilePath string `json:"dockerfilePath"` // Path to Dockerfile relative to ContextDir
	ContextDir     string `json:"contextDir"`     // Build context directory (absolute path)
	ImageName      string `json:"imageName"`      // Target image name (without registry)
	ImageTag       string `json:"imageTag"`       // Target image tag
}

// BuildResponse represents the result of a build operation
type BuildResponse struct {
	ImageID string `json:"imageID"`           // Built image ID
	FullTag string `json:"fullTag"`           // Complete image reference
	Success bool   `json:"success"`           // Build completion status
	Error   string `json:"error,omitempty"`   // Error details if Success is false
}

// Validate performs basic validation on BuildRequest
func (br *BuildRequest) Validate() error {
	if br.DockerfilePath == "" {
		return fmt.Errorf("DockerfilePath is required")
	}
	if br.ContextDir == "" {
		return fmt.Errorf("ContextDir is required")
	}
	if br.ImageName == "" {
		return fmt.Errorf("ImageName is required")
	}
	if br.ImageTag == "" {
		br.ImageTag = "latest"
	}
	return nil
}