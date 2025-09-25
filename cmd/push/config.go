package main

import (
	"errors"
	"os"

	"github.com/cnoe-io/idpbuilder-push/client-interface-tests-split-004/pkg/oci"
	"github.com/spf13/viper"
)

// PushConfig holds configuration for the push command
type PushConfig struct {
	// Registry configuration
	RegistryURL      string
	RegistryUsername string
	RegistryPassword string
	RegistryInsecure bool

	// Build configuration
	BuildPath       string
	Dockerfile      string
	BuildArgs       map[string]string
	Target          string
	Platform        string
	Context         string

	// Push configuration
	ImageTag       string
	ImageName      string
	PushToKind     bool
	KindCluster    string
	Force          bool
	DryRun         bool

	// Output configuration
	Verbose bool
	Quiet   bool
}

// NewPushConfig creates a new PushConfig with default values
func NewPushConfig() *PushConfig {
	return &PushConfig{
		RegistryURL:     viper.GetString("registry.url"),
		RegistryInsecure: viper.GetBool("registry.insecure"),
		Dockerfile:      "Dockerfile",
		Context:        ".",
		BuildArgs:      make(map[string]string),
		Platform:       "linux/amd64",
		ImageTag:       "latest",
		KindCluster:    "kind",
		Force:         false,
		DryRun:        false,
		Verbose:       false,
		Quiet:         false,
	}
}

// Validate validates the push configuration
func (c *PushConfig) Validate() error {
	if c.BuildPath == "" {
		return errors.New("build path is required")
	}

	if c.ImageName == "" {
		return errors.New("image name is required")
	}

	// Validate registry URL if pushing to registry
	if !c.PushToKind {
		if c.RegistryURL == "" {
			// Try environment variable
			c.RegistryURL = os.Getenv("REGISTRY_URL")
			if c.RegistryURL == "" {
				return errors.New("registry URL is required when pushing to registry")
			}
		}
	}

	// Validate Kind cluster if pushing to Kind
	if c.PushToKind && c.KindCluster == "" {
		return errors.New("Kind cluster name is required when pushing to Kind")
	}

	// Validate build path exists
	if _, err := os.Stat(c.BuildPath); os.IsNotExist(err) {
		return errors.New("build path does not exist: " + c.BuildPath)
	}

	return nil
}

// ToOptions converts the push configuration to OCI push options
func (c *PushConfig) ToOptions() *oci.PushOptions {
	// Create authentication if needed
	var auth *oci.RegistryAuth
	if c.RegistryUsername != "" || c.RegistryPassword != "" {
		auth = &oci.RegistryAuth{
			Username:      c.RegistryUsername,
			Password:      c.RegistryPassword,
			ServerAddress: c.RegistryURL,
		}
	}

	// Construct image reference
	imageRef := c.ImageName
	if c.ImageTag != "" {
		imageRef += ":" + c.ImageTag
	}

	return &oci.PushOptions{
		ImageRef: imageRef,
		Auth:     auth,
		Insecure: c.RegistryInsecure,
	}
}