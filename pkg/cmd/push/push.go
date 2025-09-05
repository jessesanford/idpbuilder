package push

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
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
	PushCmd.Flags().String("output", "", "Tarball path to load image from (default: derive from image name)")
	PushCmd.Flags().Int("retry", 3, "Number of retry attempts")
}

func runPush(cmd *cobra.Command, args []string) error {
	image := args[0]
	insecure, _ := cmd.Flags().GetBool("insecure")
	registryURL, _ := cmd.Flags().GetString("registry")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	output, _ := cmd.Flags().GetString("output")
	retryCount, _ := cmd.Flags().GetInt("retry")

	// Validate image name
	if !strings.Contains(image, ":") {
		image = image + ":latest"
	}

	// Create progress reporter
	// TODO: Integrate cli progress bar
	// progress := cli.NewProgressBar("Pushing " + image)
	// defer progress.Finish()
	fmt.Printf("Pushing %s\n", image)

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
		// progress.UpdateMessage("Configuring certificate trust")

		// Create trust store from Phase 1 certificate infrastructure
		var err error
		trustStore, err = certs.NewTrustStoreManager("") // Use default directory
		if err != nil {
			return fmt.Errorf("failed to create trust store: %w", err)
		}

		// Auto-configure for Gitea if using default registry
		if strings.Contains(registryURL, "gitea.cnoe.localtest.me") {
			// Detect Kind cluster name dynamically
			clusterName, err := detectKindClusterName()
			if err != nil {
				// Fall back to default
				clusterName = "localdev"
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

	// Determine tarball path
	tarballPath := output
	if tarballPath == "" {
		// Derive tarball path from image name
		tarballPath = fmt.Sprintf("%s.tar", strings.ReplaceAll(image, ":", "-"))
	}

	// Check if tarball exists
	fmt.Printf("Looking for image tarball: %s\n", tarballPath)
	if _, err := os.Stat(tarballPath); os.IsNotExist(err) {
		return fmt.Errorf("image tarball not found: %s (run 'idpbuilder build --output %s' first)", 
			tarballPath, tarballPath)
	}

	// Load OCI image from tarball
	fmt.Printf("Loading image from %s\n", tarballPath)
	img, err := tarball.ImageFromPath(tarballPath, nil)
	if err != nil {
		return fmt.Errorf("failed to load image from tarball: %w", err)
	}

	// Get image digest for reporting
	digest, err := img.Digest()
	if err != nil {
		return fmt.Errorf("failed to get image digest: %w", err)
	}

	// Push the image
	fmt.Printf("Pushing %s to %s\n", image, registryURL)
	
	// Create push options
	pushOpts := registry.PushOptions{
		Options: registry.Options{
			Timeout:  30 * time.Second,
			Insecure: insecure,
		},
	}

	if err := client.Push(context.Background(), img, image, pushOpts); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	fmt.Printf("\n✅ Image pushed successfully: %s\n", image)
	fmt.Printf("   Registry: %s\n", registryURL)
	fmt.Printf("   Digest: %s\n", digest)
	return nil
}

// detectKindClusterName detects the currently running Kind cluster name
func detectKindClusterName() (string, error) {
	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(clusters) > 0 && clusters[0] != "" {
		return clusters[0], nil
	}
	return "", fmt.Errorf("no kind clusters found")
}
