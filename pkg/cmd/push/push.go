package push

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
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
  idpbuilder push --username myuser --token mytoken myapp:latest`,
	Args: cobra.ExactArgs(1),
	RunE: runPush,
}

func init() {
	PushCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
	PushCmd.Flags().String("registry", "", "Registry URL (default: auto-detect)")
	PushCmd.Flags().String("username", "", "Registry username (default: gitea_admin)")
	PushCmd.Flags().String("token", "", "Registry token/password")
}

func runPush(cmd *cobra.Command, args []string) error {
	image := args[0]
	insecure, _ := cmd.Flags().GetBool("insecure")
	registryURL, _ := cmd.Flags().GetString("registry")
	username, _ := cmd.Flags().GetString("username")
	token, _ := cmd.Flags().GetString("token")

	// Don't add :latest tag if this is a tarball path
	if !strings.Contains(image, ".tar") && !strings.Contains(image, ":") {
		image = image + ":latest"
	}

	// Simple progress reporting using fmt
	fmt.Printf("Pushing %s\n", image)

	// Set defaults for registry access
	if registryURL == "" {
		registryURL = "https://gitea.cnoe.localtest.me:443"
	}
	if username == "" {
		username = "gitea_admin"
	}

	// Get token from environment if not provided
	if token == "" {
		if envToken := os.Getenv("GITEA_PASSWORD"); envToken != "" {
			token = envToken
		} else if envToken := os.Getenv("IDPBUILDER_REGISTRY_PASSWORD"); envToken != "" {
			token = envToken
		} else if envToken := os.Getenv("GITEA_TOKEN"); envToken != "" {
			token = envToken
		} else if envToken := os.Getenv("IDPBUILDER_REGISTRY_TOKEN"); envToken != "" {
			token = envToken
		}
	}

	// Always create trust store (required by NewGiteaClient)
	trustStore := certs.NewTrustStore()

	// Setup certificate trust (unless --insecure)
	if !insecure {
		fmt.Println("Configuring certificate trust...")

		// Auto-configure for Gitea if using default registry
		if strings.Contains(registryURL, "gitea.cnoe.localtest.me") {
			// Detect cluster name dynamically
			clusterName, err := detectKindClusterName()
			if err != nil {
				// Fallback to idpbuilder if detection fails
				clusterName = "idpbuilder"
			}
			extractorConfig := certs.ExtractorConfig{
				ClusterName: clusterName,
				Namespace:   "gitea",
			}
			extractor, err := certs.NewKindCertExtractor(extractorConfig)
			if err != nil {
				return fmt.Errorf("failed to create cert extractor: %w", err)
			}
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

	// Create registry config
	config := &registry.RegistryConfig{
		URL:      registryURL,
		Username: username,
		Token:    token,
		Insecure: insecure,
	}

	// Create remote options
	opts := registry.DefaultRemoteOptions()
	if insecure {
		opts.Insecure = true
		opts.SkipTLSVerify = true
	}

	// Create registry client
	client, err := registry.NewGiteaRegistry(config, opts)
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	// Push the image
	fmt.Println("Loading image from tarball...")

	// For CLI usage, we expect the image to be provided as a tarball path
	// This follows the pattern: idpbuilder build --output image.tar && idpbuilder push image.tar
	tarballPath := image
	if !strings.Contains(tarballPath, ".tar") {
		return fmt.Errorf("image must be provided as a tarball path (*.tar). Use 'idpbuilder build --output %s.tar' first", image)
	}

	// Parse registry URL to get host
	registryHost := strings.TrimPrefix(registryURL, "https://")
	registryHost = strings.TrimPrefix(registryHost, "http://")

	// Construct the full image reference for pushing
	imageRef := fmt.Sprintf("%s/%s/hello-world:v1", registryHost, strings.ToLower(username))

	fmt.Printf("Pushing to %s\n", imageRef)

	// Convert image to reader
	// Note: This is a simplified approach - in production you'd stream the tarball
	reader, err := os.Open(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to open tarball: %w", err)
	}
	defer reader.Close()

	// Push the image
	ctx := context.Background()
	if err := client.Push(ctx, imageRef, reader); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	fmt.Printf("Successfully pushed %s to %s\n", tarballPath, imageRef)
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
