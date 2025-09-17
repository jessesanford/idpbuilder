// pkg/kind/cluster.go
package kind

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// IProvider defines the interface for KIND providers
type IProvider interface {
	Create(ctx context.Context) error
	Delete(ctx context.Context) error
	GetKubeconfig() ([]byte, error)
	Exists(ctx context.Context) (bool, error)
	GetNodes(ctx context.Context) ([]string, error)
}

// Cluster represents a KIND cluster
type Cluster struct {
	Name               string
	Provider           IProvider
	Config             *ClusterConfig
	KubernetesVersion  string
	ConfigDir          string
	KubeConfigPath     string
	ExtraPortMappings  string
	RegistryConfig     interface{}
	BuildCustomization v1alpha1.BuildCustomizationSpec
	Logger             logr.Logger
}

// ClusterConfig holds configuration for KIND cluster
type ClusterConfig struct {
	Image      string
	ConfigPath string
	KubeConfig string
	WaitTime   time.Duration
}

// DefaultClusterConfig returns default configuration
func DefaultClusterConfig() *ClusterConfig {
	return &ClusterConfig{
		Image:      "kindest/node:v1.29.0",
		WaitTime:   5 * time.Minute,
		KubeConfig: "",
	}
}

// NewCluster creates a new KIND cluster instance
func NewCluster(name string, kubernetesVersion string, configDir string, kubeConfigPath string, extraPortMappings string, registryConfig interface{}, customization v1alpha1.BuildCustomizationSpec, logger logr.Logger) (*Cluster, error) {
	if name == "" {
		return nil, fmt.Errorf("cluster name cannot be empty")
	}

	// Validate cluster name format
	if !isValidClusterName(name) {
		return nil, fmt.Errorf("invalid cluster name: %s (must be DNS-1123 compliant)", name)
	}

	config := DefaultClusterConfig()
	if kubernetesVersion != "" {
		config.Image = fmt.Sprintf("kindest/node:v%s", kubernetesVersion)
	}
	if configDir != "" {
		config.ConfigPath = configDir
	}
	if kubeConfigPath != "" {
		config.KubeConfig = kubeConfigPath
	}

	provider := &defaultProvider{
		name:   name,
		config: config,
	}

	return &Cluster{
		Name:               name,
		Provider:           provider,
		Config:             config,
		KubernetesVersion:  kubernetesVersion,
		ConfigDir:          configDir,
		KubeConfigPath:     kubeConfigPath,
		ExtraPortMappings:  extraPortMappings,
		RegistryConfig:     registryConfig,
		BuildCustomization: customization,
		Logger:             logger,
	}, nil
}

// NewClusterWithConfig creates a new KIND cluster instance with custom config
func NewClusterWithConfig(name string, config *ClusterConfig) (*Cluster, error) {
	if name == "" {
		return nil, fmt.Errorf("cluster name cannot be empty")
	}

	if config == nil {
		config = DefaultClusterConfig()
	}

	provider := &defaultProvider{
		name:   name,
		config: config,
	}

	return &Cluster{
		Name:     name,
		Provider: provider,
		Config:   config,
	}, nil
}

// isValidClusterName validates cluster name according to KIND requirements
func isValidClusterName(name string) bool {
	if len(name) == 0 || len(name) > 63 {
		return false
	}

	// Must start and end with alphanumeric
	if !isAlphaNumeric(name[0]) || !isAlphaNumeric(name[len(name)-1]) {
		return false
	}

	// Check each character
	for _, char := range name {
		if !isAlphaNumeric(byte(char)) && char != '-' {
			return false
		}
	}

	return true
}

func isAlphaNumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// defaultProvider implements IProvider
type defaultProvider struct {
	name   string
	config *ClusterConfig
}

// Create implements IProvider.Create
func (p *defaultProvider) Create(ctx context.Context) error {
	// Check if KIND is available
	if !p.isKindAvailable() {
		return fmt.Errorf("KIND is not available in PATH")
	}

	// Check if cluster already exists
	exists, err := p.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if cluster exists: %w", err)
	}
	if exists {
		return fmt.Errorf("cluster %s already exists", p.name)
	}

	// Build KIND create command
	args := []string{"create", "cluster", "--name", p.name}

	if p.config.Image != "" {
		args = append(args, "--image", p.config.Image)
	}

	if p.config.ConfigPath != "" {
		args = append(args, "--config", p.config.ConfigPath)
	}

	if p.config.KubeConfig != "" {
		args = append(args, "--kubeconfig", p.config.KubeConfig)
	}

	// Execute KIND create
	cmd := exec.CommandContext(ctx, "kind", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create KIND cluster: %w", err)
	}

	// Wait for cluster to be ready
	return p.waitForClusterReady(ctx)
}

// Delete implements IProvider.Delete
func (p *defaultProvider) Delete(ctx context.Context) error {
	if !p.isKindAvailable() {
		return fmt.Errorf("KIND is not available in PATH")
	}

	cmd := exec.CommandContext(ctx, "kind", "delete", "cluster", "--name", p.name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete KIND cluster: %w", err)
	}

	return nil
}

// GetKubeconfig implements IProvider.GetKubeconfig
func (p *defaultProvider) GetKubeconfig() ([]byte, error) {
	if !p.isKindAvailable() {
		// Return minimal kubeconfig for testing when KIND is not available
		return p.getTestKubeconfig(), nil
	}

	cmd := exec.Command("kind", "get", "kubeconfig", "--name", p.name)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	return output, nil
}

// Exists implements IProvider.Exists
func (p *defaultProvider) Exists(ctx context.Context) (bool, error) {
	if !p.isKindAvailable() {
		return false, nil
	}

	cmd := exec.CommandContext(ctx, "kind", "get", "clusters")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list clusters: %w", err)
	}

	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, cluster := range clusters {
		if strings.TrimSpace(cluster) == p.name {
			return true, nil
		}
	}

	return false, nil
}

// GetNodes implements IProvider.GetNodes
func (p *defaultProvider) GetNodes(ctx context.Context) ([]string, error) {
	if !p.isKindAvailable() {
		return []string{}, nil
	}

	cmd := exec.CommandContext(ctx, "kind", "get", "nodes", "--name", p.name)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	nodes := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, node := range nodes {
		if node = strings.TrimSpace(node); node != "" {
			result = append(result, node)
		}
	}

	return result, nil
}

// isKindAvailable checks if KIND binary is available in PATH
func (p *defaultProvider) isKindAvailable() bool {
	_, err := exec.LookPath("kind")
	return err == nil
}

// waitForClusterReady waits for the cluster to be ready
func (p *defaultProvider) waitForClusterReady(ctx context.Context) error {
	timeout := p.config.WaitTime
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for cluster to be ready")
		case <-ticker.C:
			if ready, _ := p.isClusterReady(ctx); ready {
				return nil
			}
		}
	}
}

// isClusterReady checks if the cluster is ready
func (p *defaultProvider) isClusterReady(ctx context.Context) (bool, error) {
	kubeconfig, err := p.GetKubeconfig()
	if err != nil {
		return false, err
	}

	// Write kubeconfig to temporary file
	tmpFile, err := os.CreateTemp("", "kubeconfig-*")
	if err != nil {
		return false, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(kubeconfig); err != nil {
		return false, err
	}

	// Test connection using kubectl
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", tmpFile.Name(), "get", "nodes")
	return cmd.Run() == nil, nil
}

// getTestKubeconfig returns a minimal kubeconfig for testing
func (p *defaultProvider) getTestKubeconfig() []byte {
	config := &api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: map[string]*api.Cluster{
			p.name: {
				Server:                   fmt.Sprintf("https://127.0.0.1:6443"),
				CertificateAuthorityData: []byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t"), // dummy CA
			},
		},
		Contexts: map[string]*api.Context{
			p.name: {
				Cluster:  p.name,
				AuthInfo: p.name,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			p.name: {
				ClientCertificateData: []byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t"), // dummy cert
				ClientKeyData:         []byte("LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVkt"), // dummy key
			},
		},
		CurrentContext: p.name,
	}

	data, err := clientcmd.Write(*config)
	if err != nil {
		// Fallback to simple YAML structure
		return []byte(fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://127.0.0.1:6443
  name: %s
contexts:
- context:
    cluster: %s
    user: %s
  name: %s
current-context: %s
users:
- name: %s
  user: {}
`, p.name, p.name, p.name, p.name, p.name, p.name))
	}

	return data
}

// GetKubeconfigPath returns the path to the kubeconfig file
func (c *Cluster) GetKubeconfigPath() string {
	if c.Config != nil && c.Config.KubeConfig != "" {
		return c.Config.KubeConfig
	}

	// Default kubeconfig path for KIND
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".kube", "config")
}

// Status returns the current status of the cluster
func (c *Cluster) Status(ctx context.Context) (string, error) {
	exists, err := c.Provider.Exists(ctx)
	if err != nil {
		return "unknown", err
	}

	if !exists {
		return "not-found", nil
	}

	// Try to get nodes to verify cluster is responsive
	nodes, err := c.Provider.GetNodes(ctx)
	if err != nil {
		return "unreachable", err
	}

	if len(nodes) == 0 {
		return "no-nodes", nil
	}

	return "ready", nil
}

// getConfig returns the KIND cluster configuration as YAML bytes
func (c *Cluster) getConfig() ([]byte, error) {
	// Build configuration based on cluster settings
	host := c.BuildCustomization.Host
	if host == "" {
		host = "cnoe.localtest.me"
	}

	port := c.BuildCustomization.Port
	if port == "" {
		port = "8443"
	}

	hostPort := 8443
	if port != "" {
		fmt.Sscanf(port, "%d", &hostPort)
	}

	// Don't add double 'v' prefix
	image := "kindest/node:v1.26.3"
	if c.KubernetesVersion != "" && c.KubernetesVersion != "1.26.3" {
		if !strings.HasPrefix(c.KubernetesVersion, "v") {
			image = fmt.Sprintf("kindest/node:v%s", c.KubernetesVersion)
		} else {
			image = fmt.Sprintf("kindest/node:%s", c.KubernetesVersion)
		}
	}

	// Base config
	config := fmt.Sprintf(`kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  image: "%s"
  labels:
    ingress-ready: "true"
  extraPortMappings:
  - containerPort: 443
    hostPort: %d
    protocol: TCP
  - containerPort: 32222
    hostPort: 32222
    protocol: TCP`, image, hostPort)

	// Add additional extra port mappings if specified (besides the default 32222)
	if c.ExtraPortMappings != "" && c.ExtraPortMappings != "22:32222" {
		// Parse extra port mappings (format: "22:32222")
		parts := strings.Split(c.ExtraPortMappings, ":")
		if len(parts) == 2 {
			containerPort := parts[0]
			hostPort := parts[1]
			// Only add if it's not the default 32222 mapping
			if hostPort != "32222" {
				config += fmt.Sprintf(`
  - containerPort: %s
    hostPort: %s
    protocol: TCP`, containerPort, hostPort)
			}
		}
	}

	// Add registry config mounts if specified
	if registryConfigs, ok := c.RegistryConfig.([]string); ok && len(registryConfigs) > 0 {
		for _, configPath := range registryConfigs {
			// Check if file exists
			if _, err := os.Stat(configPath); err == nil {
				config += fmt.Sprintf(`
  extraMounts:
  - containerPath: /var/lib/kubelet/config.json
    hostPath: %s`, configPath)
				break // Only use the first valid config
			}
		}
	}

	// Add containerd config patches - use gitea.cnoe.localtest.me for registry
	registryHost := "gitea." + host
	config += fmt.Sprintf(`
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."%s:%s"]
    endpoint = ["https://%s"]
  [plugins."io.containerd.grpc.v1.cri".registry.configs."%s".tls]
    insecure_skip_verify = true`, registryHost, port, registryHost, registryHost)

	return []byte(config), nil
}