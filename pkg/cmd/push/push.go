package push

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/cnoe-io/idpbuilder/pkg/cli"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/spf13/cobra"
)

// PushCmd is the push command
var PushCmd = &cobra.Command{
	Use:   "push IMAGE[:TAG]",
	Short: "Push image to Gitea registry",
	Long: `Push a container image to the builtin Gitea registry with certificate support.
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
		trustStore, err = certs.NewTrustStoreManager("") // Use default directory
		if err != nil {
			return fmt.Errorf("failed to create trust store: %w", err)
		}

		// Auto-configure for Gitea if using default registry
		if strings.Contains(registryURL, "gitea.cnoe.localtest.me") {
			// Detect cluster name dynamically
			clusterName, err := detectKindClusterName()
			if err != nil {
				// Fallback to idpbuilder if detection fails
				clusterName = "idpbuilder"
			}
			extractor := certs.NewDefaultExtractor(clusterName)
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
	client, err := registry.NewGiteaClient(registryURL, username, password, trustStore, clientOpts...)
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	// Push the image
	progress.UpdateMessage("Loading image from tarball")

	// For CLI usage, we expect the image to be provided as a tarball path
	// This follows the pattern: idpbuilder build --output image.tar && idpbuilder push image.tar
	tarballPath := image
	if !strings.Contains(tarballPath, ".tar") {
		return fmt.Errorf("image must be provided as a tarball path (*.tar). Use 'idpbuilder build --output %s.tar' first", image)
	}

	// Load the OCI image from the tarball
	img, err := tarball.ImageFromPath(tarballPath, nil)
	if err != nil {
		return fmt.Errorf("failed to load image from tarball: %w", err)
	}

	progress.UpdateMessage("Pushing to registry")

	// Prepare push options
	pushOpts := registry.PushOptions{
		Options: registry.Options{
			Insecure: insecure,
			Timeout:  30 * time.Second,
		},
	}

	// Push the image
	if err := client.Push(context.Background(), img, image, pushOpts); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	fmt.Printf("Successfully pushed %s\n", image)
	return nil
}

// detectKindClusterName dynamically detects the available Kind cluster name
func detectKindClusterName() (string, error) {
	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get kind clusters: %w", err)
	}

	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(clusters) == 0 {
		return "", fmt.Errorf("no kind clusters found")
	}

	// Use first cluster or look for idpbuilder/localdev
	for _, cluster := range clusters {
		if cluster == "idpbuilder" || cluster == "localdev" {
			return cluster, nil
		}
	}

	// Default to first available cluster
	return clusters[0], nil
}
