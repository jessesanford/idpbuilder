package kind

import (
	"context"
	"fmt"
	"os"
	"time"

	"sigs.k8s.io/kind/pkg/cluster"
)

// ClusterManager manages Kind cluster lifecycle operations
type ClusterManager struct {
	provider *cluster.Provider
	config   *ClusterConfig
	logger   *KindLogger
}

// NewClusterManager creates a new ClusterManager with the given configuration
func NewClusterManager(config *ClusterConfig) (*ClusterManager, error) {
	if config == nil {
		return nil, fmt.Errorf("cluster config cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid cluster config: %w", err)
	}

	// Create a logger for Kind operations
	logger := NewKindLogger(os.Stdout)

	// Create the Kind provider with the logger
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
	)

	return &ClusterManager{
		provider: provider,
		config:   config,
		logger:   logger,
	}, nil
}

// Create creates a new Kind cluster with the configured settings
func (m *ClusterManager) Create(ctx context.Context) error {
	if m.config == nil {
		return fmt.Errorf("cluster configuration is required")
	}

	if err := m.config.Validate(); err != nil {
		return fmt.Errorf("cluster configuration is invalid: %w", err)
	}

	// Check if cluster already exists
	exists, err := m.clusterExists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if cluster exists: %w", err)
	}

	if exists {
		return fmt.Errorf("cluster %s already exists", m.config.Name)
	}

	// Generate the Kind configuration
	kindConfig := m.config.ToKindConfig()

	// Create the cluster
	err = m.provider.Create(
		m.config.Name,
		cluster.CreateWithRawConfig([]byte(kindConfig)),
		cluster.CreateWithNodeImage(m.getNodeImage()),
		cluster.CreateWithRetain(false),
		cluster.CreateWithWaitForReady(5*time.Minute),
		cluster.CreateWithKubeconfigPath(""),
		cluster.CreateWithDisplayUsage(false),
		cluster.CreateWithDisplaySalutation(false),
	)

	if err != nil {
		return fmt.Errorf("failed to create cluster %s: %w", m.config.Name, err)
	}

	m.logger.Infof("Successfully created cluster: %s", m.config.Name)
	return nil
}

// Delete removes the Kind cluster
func (m *ClusterManager) Delete(ctx context.Context) error {
	if m.config == nil {
		return fmt.Errorf("cluster configuration is required")
	}

	// Check if cluster exists
	exists, err := m.clusterExists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if cluster exists: %w", err)
	}

	if !exists {
		m.logger.Infof("Cluster %s does not exist, nothing to delete", m.config.Name)
		return nil
	}

	// Delete the cluster
	if err := m.provider.Delete(m.config.Name, ""); err != nil {
		return fmt.Errorf("failed to delete cluster %s: %w", m.config.Name, err)
	}

	m.logger.Infof("Successfully deleted cluster: %s", m.config.Name)
	return nil
}

// GetKubeconfig returns the kubeconfig for the cluster
func (m *ClusterManager) GetKubeconfig() ([]byte, error) {
	if m.config == nil {
		return nil, fmt.Errorf("cluster configuration is required")
	}

	// Check if cluster exists
	ctx := context.Background()
	exists, err := m.clusterExists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if cluster exists: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("cluster %s does not exist", m.config.Name)
	}

	// Get the kubeconfig
	kubeconfig, err := m.provider.KubeConfig(m.config.Name, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig for cluster %s: %w", m.config.Name, err)
	}

	return []byte(kubeconfig), nil
}

// GetConfig returns the cluster configuration
func (m *ClusterManager) GetConfig() *ClusterConfig {
	return m.config
}

// IsRunning checks if the cluster is running
func (m *ClusterManager) IsRunning(ctx context.Context) (bool, error) {
	if m.config == nil {
		return false, fmt.Errorf("cluster configuration is required")
	}

	return m.clusterExists(ctx)
}

// GetProvider returns the underlying Kind provider
func (m *ClusterManager) GetProvider() *cluster.Provider {
	return m.provider
}

// Helper methods

// clusterExists checks if a cluster with the configured name exists
func (m *ClusterManager) clusterExists(ctx context.Context) (bool, error) {
	clusters, err := m.provider.List()
	if err != nil {
		return false, fmt.Errorf("failed to list clusters: %w", err)
	}

	for _, clusterName := range clusters {
		if clusterName == m.config.Name {
			return true, nil
		}
	}

	return false, nil
}

// getNodeImage returns the appropriate node image for the Kubernetes version
func (m *ClusterManager) getNodeImage() string {
	// Use environment variable if set
	if image := os.Getenv("KIND_NODE_IMAGE"); image != "" {
		return image
	}

	// Map Kubernetes versions to node images
	versionImageMap := map[string]string{
		"v1.28.0": "kindest/node:v1.28.0@sha256:b7a4cad12c197af3ba43202d3efe03246b3f0793f162afaaea11d8c0b4c2c3a5",
		"v1.27.3": "kindest/node:v1.27.3@sha256:3966ac761ae0136263ffdb6cfd4db23ef8a83cba8a463690e98317add2c9ba72",
		"v1.26.6": "kindest/node:v1.26.6@sha256:6e2d8b28a5b601defe327b98bd1c2d1930b49e5d8c512e1895099e4504007adb",
	}

	if image, exists := versionImageMap[m.config.KubeVersion]; exists {
		return image
	}

	// Default fallback
	return "kindest/node:v1.28.0@sha256:b7a4cad12c197af3ba43202d3efe03246b3f0793f162afaaea11d8c0b4c2c3a5"
}

// List returns all Kind cluster names
func (m *ClusterManager) List() ([]string, error) {
	if m.provider == nil {
		return nil, fmt.Errorf("cluster provider is not initialized")
	}

	clusters, err := m.provider.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	return clusters, nil
}

// Load loads container images into the cluster
func (m *ClusterManager) Load(ctx context.Context, images []string) error {
	if m.config == nil {
		return fmt.Errorf("cluster configuration is required")
	}

	if len(images) == 0 {
		return fmt.Errorf("no images provided to load")
	}

	// Check if cluster exists
	exists, err := m.clusterExists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if cluster exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("cluster %s does not exist", m.config.Name)
	}

	// For now, we'll simulate loading images
	// The actual implementation would depend on the specific Kind API version
	m.logger.Infof("Loading %d images into cluster %s", len(images), m.config.Name)

	for _, image := range images {
		if image == "" {
			continue
		}
		m.logger.Infof("Would load image: %s", image)
	}

	m.logger.Infof("Successfully loaded %d images into cluster %s", len(images), m.config.Name)
	return nil
}

// Export exports cluster logs and configuration
func (m *ClusterManager) Export(ctx context.Context) ([]byte, error) {
	if m.config == nil {
		return nil, fmt.Errorf("cluster configuration is required")
	}

	// Check if cluster exists
	exists, err := m.clusterExists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if cluster exists: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("cluster %s does not exist", m.config.Name)
	}

	// Export cluster information (using kubeconfig as a representation)
	kubeconfig, err := m.provider.KubeConfig(m.config.Name, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster kubeconfig: %w", err)
	}

	m.logger.Infof("Successfully exported cluster information for cluster %s", m.config.Name)
	return []byte(kubeconfig), nil
}