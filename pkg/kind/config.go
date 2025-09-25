package kind

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ClusterConfig represents configuration for a Kind cluster
type ClusterConfig struct {
	Name           string
	KubeVersion    string
	ControlPlanes  int
	Workers        int
	RegistryMirror string
	PortMappings   []PortMapping
}

// PortMapping represents a port mapping for the cluster
type PortMapping struct {
	ContainerPort int
	HostPort      int
	Protocol      string
}

// NewClusterConfig creates a new ClusterConfig with default values
func NewClusterConfig(name string) *ClusterConfig {
	if name == "" {
		name = getDefaultClusterName()
	}

	return &ClusterConfig{
		Name:          name,
		KubeVersion:   getDefaultKubeVersion(),
		ControlPlanes: getDefaultControlPlanes(),
		Workers:       getDefaultWorkers(),
		RegistryMirror: getDefaultRegistryMirror(),
		PortMappings:  getDefaultPortMappings(),
	}
}

// Validate performs validation on the cluster configuration
func (c *ClusterConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("cluster name cannot be empty")
	}

	if !isValidName(c.Name) {
		return fmt.Errorf("cluster name must contain only alphanumeric characters and hyphens")
	}

	if c.KubeVersion == "" {
		return fmt.Errorf("Kubernetes version cannot be empty")
	}

	if c.ControlPlanes < 1 {
		return fmt.Errorf("must have at least 1 control plane node")
	}

	if c.Workers < 0 {
		return fmt.Errorf("worker count cannot be negative")
	}

	return nil
}

// ToKindConfig generates the Kind cluster configuration YAML
func (c *ClusterConfig) ToKindConfig() string {
	var config strings.Builder

	config.WriteString("kind: Cluster\n")
	config.WriteString("apiVersion: kind.x-k8s.io/v1alpha4\n")

	if c.KubeVersion != "" {
		config.WriteString(fmt.Sprintf("kubeadmConfigPatches:\n"))
		config.WriteString(fmt.Sprintf("- |\n"))
		config.WriteString(fmt.Sprintf("  kind: ClusterConfiguration\n"))
		config.WriteString(fmt.Sprintf("  kubernetesVersion: %s\n", c.KubeVersion))
	}

	config.WriteString("nodes:\n")

	// Add control plane nodes
	for i := 0; i < c.ControlPlanes; i++ {
		config.WriteString("- role: control-plane\n")

		if i == 0 && len(c.PortMappings) > 0 {
			config.WriteString("  extraPortMappings:\n")
			for _, pm := range c.PortMappings {
				config.WriteString(fmt.Sprintf("  - containerPort: %d\n", pm.ContainerPort))
				config.WriteString(fmt.Sprintf("    hostPort: %d\n", pm.HostPort))
				config.WriteString(fmt.Sprintf("    protocol: %s\n", pm.Protocol))
			}
		}
	}

	// Add worker nodes
	for i := 0; i < c.Workers; i++ {
		config.WriteString("- role: worker\n")
	}

	return config.String()
}

// Helper functions for getting configuration values
func getDefaultClusterName() string {
	if name := os.Getenv("KIND_CLUSTER_NAME"); name != "" {
		return name
	}
	return "idpbuilder-test"
}

func getDefaultKubeVersion() string {
	if version := os.Getenv("KIND_KUBE_VERSION"); version != "" {
		return version
	}
	return "v1.28.0"
}

func getDefaultControlPlanes() int {
	if cp := os.Getenv("KIND_CONTROL_PLANES"); cp != "" {
		if count, err := strconv.Atoi(cp); err == nil && count > 0 {
			return count
		}
	}
	return 1
}

func getDefaultWorkers() int {
	if workers := os.Getenv("KIND_WORKERS"); workers != "" {
		if count, err := strconv.Atoi(workers); err == nil && count >= 0 {
			return count
		}
	}
	return 0
}

func getDefaultRegistryMirror() string {
	return os.Getenv("KIND_REGISTRY_MIRROR")
}

func getDefaultPortMappings() []PortMapping {
	// Default port mappings for common services
	return []PortMapping{
		{ContainerPort: 30080, HostPort: 30080, Protocol: "TCP"},
		{ContainerPort: 30443, HostPort: 30443, Protocol: "TCP"},
	}
}

func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}

	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
		     (char >= 'A' && char <= 'Z') ||
		     (char >= '0' && char <= '9') ||
		     char == '-') {
			return false
		}
	}

	return true
}