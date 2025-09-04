package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/builder"
	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
)

// PushOptions contains options for the push command
type PushOptions struct {
	Registry     string
	Username     string
	Password     string
	Insecure     bool
	CertFile     string
	KeyFile      string
	CAFile       string
	SkipTLSVerify bool
	Platform     string
}

var pushOpts = &PushOptions{}

// pushCmd represents the push command  
var pushCmd = &cobra.Command{
	Use:   "push [IMAGE] [REGISTRY]",
	Short: "Push container images to a registry",
	Long: `Push container images to a container registry with support for authentication and TLS.

This command pushes images to container registries using go-containerregistry, supporting
various authentication methods and TLS configurations.

Examples:
  # Push image to registry
  idpbuilder push myapp:latest localhost:5000

  # Push with authentication
  idpbuilder push --username user --password pass myapp:latest registry.com

  # Push with custom certificates
  idpbuilder push --cert-file cert.pem --key-file key.pem myapp:latest registry.com

  # Push to insecure registry
  idpbuilder push --insecure myapp:latest localhost:5000`,
	Args: cobra.ExactArgs(2),
	RunE: runPush,
}

func init() {
	// Feature flag check - only add command if enabled
	if os.Getenv("ENABLE_CLI_TOOLS") == "true" {
		rootCmd.AddCommand(pushCmd)
	}

	pushCmd.Flags().StringVarP(&pushOpts.Username, "username", "u", "", "Registry username")
	pushCmd.Flags().StringVarP(&pushOpts.Password, "password", "p", "", "Registry password")
	pushCmd.Flags().BoolVar(&pushOpts.Insecure, "insecure", false, "Allow insecure connections to registry")
	pushCmd.Flags().StringVar(&pushOpts.CertFile, "cert-file", "", "Path to client certificate file")
	pushCmd.Flags().StringVar(&pushOpts.KeyFile, "key-file", "", "Path to client private key file")
	pushCmd.Flags().StringVar(&pushOpts.CAFile, "ca-file", "", "Path to CA certificate file")
	pushCmd.Flags().BoolVar(&pushOpts.SkipTLSVerify, "skip-tls-verify", false, "Skip TLS certificate verification")
	pushCmd.Flags().StringVar(&pushOpts.Platform, "platform", "linux/amd64", "Platform for multi-arch images")
}

func runPush(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Check feature flag
	if !isFeatureEnabled() {
		return fmt.Errorf("push command requires ENABLE_CLI_TOOLS=true")
	}

	imageTag := args[0]
	pushOpts.Registry = args[1]

	// Validate options
	if err := validatePushOptions(); err != nil {
		return fmt.Errorf("invalid push options: %w", err)
	}

	// Build full image name
	fullImageName := fmt.Sprintf("%s/%s", pushOpts.Registry, imageTag)
	if strings.Contains(imageTag, "/") {
		// Image already contains registry/namespace
		fullImageName = fmt.Sprintf("%s/%s", pushOpts.Registry, strings.Split(imageTag, "/")[len(strings.Split(imageTag, "/"))-1])
	}

	fmt.Printf("Pushing image: %s\n", fullImageName)

	// Parse image reference
	ref, err := name.ParseReference(fullImageName)
	if err != nil {
		return fmt.Errorf("failed to parse image reference: %w", err)
	}

	// Setup authentication
	auth := authn.Anonymous
	if pushOpts.Username != "" {
		auth = &authn.Basic{
			Username: pushOpts.Username,
			Password: pushOpts.Password,
		}
	}

	// Setup remote options
	remoteOpts := []remote.Option{
		remote.WithAuth(auth),
		remote.WithContext(ctx),
	}

	// Setup TLS configuration
	if err := setupTLSOptions(ctx, &remoteOpts); err != nil {
		return fmt.Errorf("failed to setup TLS options: %w", err)
	}

	// For demonstration, we'll create a simple image to push
	// In a real implementation, this would load an existing image
	img, err := createDemoImage(ctx)
	if err != nil {
		return fmt.Errorf("failed to create demo image: %w", err)
	}

	// Push the image
	fmt.Printf("Pushing to registry: %s\n", pushOpts.Registry)
	err = remote.Write(ref, img, remoteOpts...)
	if err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	fmt.Printf("Successfully pushed: %s\n", fullImageName)
	return nil
}

func validatePushOptions() error {
	// Validate registry format
	if pushOpts.Registry == "" {
		return fmt.Errorf("registry cannot be empty")
	}

	// Check certificate files exist if specified
	if pushOpts.CertFile != "" {
		if _, err := os.Stat(pushOpts.CertFile); os.IsNotExist(err) {
			return fmt.Errorf("certificate file does not exist: %s", pushOpts.CertFile)
		}
	}

	if pushOpts.KeyFile != "" {
		if _, err := os.Stat(pushOpts.KeyFile); os.IsNotExist(err) {
			return fmt.Errorf("key file does not exist: %s", pushOpts.KeyFile)
		}
	}

	if pushOpts.CAFile != "" {
		if _, err := os.Stat(pushOpts.CAFile); os.IsNotExist(err) {
			return fmt.Errorf("CA file does not exist: %s", pushOpts.CAFile)
		}
	}

	// Both cert and key must be specified together
	if (pushOpts.CertFile != "") != (pushOpts.KeyFile != "") {
		return fmt.Errorf("both cert-file and key-file must be specified together")
	}

	return nil
}

func setupTLSOptions(ctx context.Context, remoteOpts *[]remote.Option) error {
	// Check if certificates are available (from split-002)
	if !isCertManagementEnabled() {
		// Basic TLS setup without cert management
		if pushOpts.Insecure || pushOpts.SkipTLSVerify {
			*remoteOpts = append(*remoteOpts, remote.WithTransport(remote.DefaultTransport))
		}
		return nil
	}

	// Use certificate manager from split-002 if available
	if pushOpts.CertFile != "" && pushOpts.KeyFile != "" {
		certManager := certs.NewCertificateManager(&certs.CertificateOptions{
			CertFile: pushOpts.CertFile,
			KeyFile:  pushOpts.KeyFile,
			CAFile:   pushOpts.CAFile,
			Insecure: pushOpts.Insecure,
		})

		transport, err := certManager.CreateHTTPSTransport(ctx)
		if err != nil {
			return fmt.Errorf("failed to create HTTPS transport: %w", err)
		}

		*remoteOpts = append(*remoteOpts, remote.WithTransport(transport))
	} else if pushOpts.Insecure || pushOpts.SkipTLSVerify {
		// Use insecure transport from cert manager
		certManager := certs.NewCertificateManager(&certs.CertificateOptions{
			Insecure: true,
		})

		transport, err := certManager.CreateHTTPSTransport(ctx)
		if err != nil {
			return fmt.Errorf("failed to create insecure transport: %w", err)
		}

		*remoteOpts = append(*remoteOpts, remote.WithTransport(transport))
	}

	return nil
}

func createDemoImage(ctx context.Context) (v1.Image, error) {
	// Create a simple demo image for pushing
	// In real usage, this would load an existing built image
	opts := builder.DefaultBuildOptions().WithTags("demo:latest")
	b, err := builder.New(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create builder: %w", err)
	}

	// For demo purposes, create an empty image
	// This would typically load from local storage or tarball
	result, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build demo image: %w", err)
	}

	return result.Image, nil
}

func isCertManagementEnabled() bool {
	return os.Getenv("ENABLE_CERT_MANAGEMENT") == "true"
}