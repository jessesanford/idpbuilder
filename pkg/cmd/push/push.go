package push

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/cnoe-io/idpbuilder/pkg/cli"
	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// PushCmd is the push command
var PushCmd = &cobra.Command{
	Use:   "push IMAGE[:TAG]",
	Short: "Push image to Gitea registry",
	Long:  `Push a container image to the builtin Gitea registry with certificate support.
Automatically handles certificate trust configuration for secure connections.`,
	Example: `  idpbuilder push myapp:v1
  idpbuilder push --insecure myapp:latest
  idpbuilder push --registry https://gitea.example.com myapp:v2
  idpbuilder push --username myuser --password mypass myapp:latest`,
	Args: cobra.ExactArgs(1),
	RunE: runPush,
}

func init() {
	PushCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
	PushCmd.Flags().String("registry", "", "Registry URL (default: auto-detect)")
	PushCmd.Flags().String("username", "", "Registry username (default: gitea_admin)")
	PushCmd.Flags().String("password", "", "Registry password")
	PushCmd.Flags().Int("retry", 3, "Number of retry attempts")
}

func runPush(cmd *cobra.Command, args []string) error {
	image := args[0]
	insecure, _ := cmd.Flags().GetBool("insecure")
	registryURL, _ := cmd.Flags().GetString("registry")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	retryCount, _ := cmd.Flags().GetInt("retry")

	// Validate image name
	if !strings.Contains(image, ":") {
		image = image + ":latest"
	}

	// Create progress reporter
	progress := cli.NewProgressBar("Pushing " + image)
	defer progress.Finish()

	// Set defaults for registry access
	if registryURL == "" {
		registryURL = "https://gitea.cnoe.localtest.me:443"
	}
	if username == "" {
		username = "gitea_admin"
	}

	// Get password from environment if not provided
	if password == "" {
		if envPassword := os.Getenv("GITEA_PASSWORD"); envPassword != "" {
			password = envPassword
		} else if envPassword := os.Getenv("IDPBUILDER_REGISTRY_PASSWORD"); envPassword != "" {
			password = envPassword
		}
	}

	var trustStore certs.TrustStoreManager
	
	// Setup certificate trust (unless --insecure)
	if !insecure {
		progress.UpdateMessage("Configuring certificate trust")
		
		// Create trust store from Phase 1 certificate infrastructure
		var err error
		trustStore, err = certs.NewTrustStoreManager("")  // Use default directory
		if err != nil {
			return fmt.Errorf("failed to create trust store: %w", err)
		}
		
		// Auto-configure for Gitea if using default registry
		if strings.Contains(registryURL, "gitea.cnoe.localtest.me") {
			extractor := certs.NewDefaultExtractor("idpbuilder")
			ctx := context.Background()
			cert, err := extractor.ExtractGiteaCert(ctx)
			if err != nil {
				return fmt.Errorf("certificate extraction failed: %w", err)
			}
			
			if err := trustStore.AddCertificate("gitea.cnoe.localtest.me", cert); err != nil {
				return fmt.Errorf("certificate setup failed: %w", err)
			}
		}
	}

	// Create registry client options
	clientOpts := []registry.ClientOption{
		registry.WithTimeout(30 * time.Second),
		registry.WithRetryConfig(retryCount, 2*time.Second),
	}

	if insecure {
		clientOpts = append(clientOpts, registry.WithInsecure(true))
	}

	// Create registry client
	_, err := registry.NewGiteaClient(registryURL, username, password, trustStore, clientOpts...)
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	// Push the image
	progress.UpdateMessage("Pushing to registry")
	
	// Note: This is a placeholder - in a real implementation we'd need to:
	// 1. Load the image from local storage/daemon  
	// 2. Parse the image reference properly
	// 3. Handle the actual image data with client.Push()
	
	// For now, we'll return an informative error until image loading is implemented
	return fmt.Errorf("image loading not yet implemented - this command needs integration with local image storage")
}