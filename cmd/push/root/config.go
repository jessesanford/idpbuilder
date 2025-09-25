/*
Copyright 2024 The idpbuilder Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package root

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Config holds the configuration for the push command
type Config struct {
	// KubeConfig path to the kubeconfig file
	KubeConfig string
	// Context is the kubeconfig context to use
	Context string
	// Namespace is the default namespace for operations
	Namespace string
	// Registry configuration
	Registry RegistryConfig
	// Build configuration
	Build BuildConfig
	// Verbose enables verbose logging
	Verbose bool
	// DryRun enables dry-run mode
	DryRun bool
}

// RegistryConfig holds registry configuration
type RegistryConfig struct {
	// URL is the registry base URL
	URL string
	// Username for registry authentication
	Username string
	// Password for registry authentication
	Password string
	// Insecure allows insecure registry connections
	Insecure bool
	// Namespace is the registry namespace/organization
	Namespace string
}

// BuildConfig holds build configuration
type BuildConfig struct {
	// Strategy specifies the build strategy (dockerfile, buildpacks, etc.)
	Strategy string
	// Context is the build context directory
	Context string
	// Dockerfile path relative to context
	Dockerfile string
	// Args are build-time arguments
	Args map[string]string
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		KubeConfig: getDefaultKubeConfig(),
		Namespace:  "default",
		Registry: RegistryConfig{
			URL:       "localhost:5000",
			Insecure:  true,
			Namespace: "idpbuilder",
		},
		Build: BuildConfig{
			Strategy:   "dockerfile",
			Context:    ".",
			Dockerfile: "Dockerfile",
			Args:       make(map[string]string),
		},
		Verbose: false,
		DryRun:  false,
	}
}

// getDefaultKubeConfig returns the default kubeconfig path
func getDefaultKubeConfig() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}

	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	return ""
}

// GetRestConfig creates a Kubernetes REST config
func (c *Config) GetRestConfig() (*rest.Config, error) {
	// If no kubeconfig path is specified, try in-cluster config
	if c.KubeConfig == "" {
		config, err := rest.InClusterConfig()
		if err == nil {
			return config, nil
		}
		// Fall back to default kubeconfig path
		c.KubeConfig = getDefaultKubeConfig()
	}

	// Load config from kubeconfig file
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: c.KubeConfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: c.Context,
		},
	).ClientConfig()

	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	if c.Registry.URL == "" {
		return fmt.Errorf("registry URL cannot be empty")
	}
	if c.Build.Context == "" {
		return fmt.Errorf("build context cannot be empty")
	}
	return nil
}