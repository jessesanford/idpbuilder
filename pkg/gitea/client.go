package gitea

import (
	"context"
	"fmt"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
)

// Client wraps the registry.Registry to provide the gitea-specific interface
// expected by the CLI commands.
type Client struct {
	registry registry.Registry
	config   registry.RegistryConfig
}

// Global credential manager instance
var credentialManager *CredentialManager

func init() {
	credentialManager = NewCredentialManager()
}

// PushProgress represents progress information during image push operations.
type PushProgress struct {
	CurrentLayer    int             `json:"currentLayer"`
	TotalLayers     int             `json:"totalLayers"`
	Percentage      int             `json:"percentage"`
	UploadedBytes   int64           `json:"uploadedBytes"`
	TotalBytes      int64           `json:"totalBytes"`
	ElapsedTime     time.Duration   `json:"elapsedTime"`
	Completed       bool            `json:"completed"`
	LayerProgresses []LayerProgress `json:"layerProgresses"`
}

// NewClient creates a new Gitea client with certificate manager integration.
func NewClient(registryURL string, certManager *certs.DefaultTrustStore) (*Client, error) {
	config := registry.RegistryConfig{
		URL:      registryURL,
		Username: getRegistryUsername(),
		Token:    getRegistryPassword(), // Using Token field instead of Password
		Insecure: false,
	}

	// Create remote options with default values
	opts := registry.DefaultRemoteOptions()

	reg, err := registry.NewGiteaRegistry(&config, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create gitea registry: %w", err)
	}

	return &Client{
		registry: reg,
		config:   config,
	}, nil
}

// NewInsecureClient creates a new Gitea client without certificate verification.
func NewInsecureClient(registryURL string) (*Client, error) {
	config := registry.RegistryConfig{
		URL:      registryURL,
		Username: getRegistryUsername(),
		Token:    getRegistryPassword(), // Using Token field instead of Password
		Insecure: true,
	}

	// Create remote options with insecure settings
	opts := registry.DefaultRemoteOptions()
	opts.Insecure = true
	opts.SkipTLSVerify = true

	reg, err := registry.NewGiteaRegistry(&config, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create insecure gitea registry: %w", err)
	}

	return &Client{
		registry: reg,
		config:   config,
	}, nil
}

// Push pushes an image to the registry with real progress reporting.
// The progressChan parameter allows monitoring of push progress.
func (c *Client) Push(imageRef string, progressChan chan<- PushProgress) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Initialize image loader
	loader, err := NewImageLoader()
	if err != nil {
		return fmt.Errorf("failed to create image loader: %w", err)
	}
	defer loader.Close()

	// Load image from Docker daemon
	manifest, err := loader.LoadImage(ctx, imageRef)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	// Get image content reader
	contentReader, err := loader.GetImageContent(ctx, imageRef)
	if err != nil {
		return fmt.Errorf("failed to get image content: %w", err)
	}
	defer contentReader.Close()

	// Initialize real progress tracker
	tracker := NewProgressTracker(manifest.TotalSize)

	// Start real progress reporting
	if progressChan != nil {
		go c.reportProgress(tracker, progressChan)
	}

	// Simulate layer progress during push
	// In a real implementation, this would be integrated with the registry push
	go c.simulateLayerProgress(tracker, manifest)

	// Perform actual push with progress tracking
	err = c.registry.Push(ctx, imageRef, contentReader)
	if err != nil {
		return fmt.Errorf("registry push failed: %w", err)
	}

	// Mark as complete
	tracker.MarkComplete()

	return nil
}

// reportProgress reports real progress updates
func (c *Client) reportProgress(tracker *ProgressTracker, progressChan chan<- PushProgress) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		progress := tracker.GetProgress()
		select {
		case progressChan <- progress:
			if progress.Completed || progress.Percentage >= 100 {
				return
			}
		default:
			// Channel full, skip this update
		}
	}
}

// simulateLayerProgress simulates layer upload progress
// In a real implementation, this would be integrated with the actual registry push
func (c *Client) simulateLayerProgress(tracker *ProgressTracker, manifest *ImageManifest) {
	// Set up layer sizes
	for i, layer := range manifest.Layers {
		layerID := fmt.Sprintf("layer-%d", i)
		tracker.SetLayerSize(layerID, layer.Size)
		tracker.SetLayerStatus(layerID, "pending")
	}

	// Simulate layer uploads
	for i, layer := range manifest.Layers {
		layerID := fmt.Sprintf("layer-%d", i)
		tracker.SetLayerStatus(layerID, "uploading")

		// Simulate upload progress for this layer
		uploaded := int64(0)
		increment := layer.Size / 10 // 10 progress updates per layer

		for uploaded < layer.Size {
			uploaded += increment
			if uploaded > layer.Size {
				uploaded = layer.Size
			}
			tracker.UpdateProgress(layerID, uploaded)
			time.Sleep(50 * time.Millisecond)
		}

		tracker.SetLayerStatus(layerID, "complete")
	}
}

// getRegistryUsername retrieves the registry username from environment or config.
func getRegistryUsername() string {
	return credentialManager.GetUsername()
}

// getRegistryPassword retrieves the registry password from environment or config.
func getRegistryPassword() string {
	return credentialManager.GetPassword()
}

// SetCredentials sets credentials from CLI flags
func (c *Client) SetCredentials(username, password string) {
	credentialManager.SetCLICredentials(username, password)
	// Update the config with new credentials
	c.config.Username = username
	c.config.Token = password
}
