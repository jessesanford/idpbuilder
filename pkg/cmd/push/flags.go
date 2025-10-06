package push

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Flag variables
	username     string
	password     string
	insecure     bool
	dryRun       bool
	verbose      bool
	repository   string
	tag          string
)

const (
	// Flag usage strings
	usernameUsage    = "Registry username (env: REGISTRY_USERNAME)"
	passwordUsage    = "Registry password (env: REGISTRY_PASSWORD)"
	insecureUsage    = "Skip TLS certificate verification (use for self-signed certificates)"
	dryRunUsage      = "Perform a dry run without actually pushing"
	verboseUsage     = "Enable verbose output"
	repositoryUsage  = "Override target repository name"
	tagUsage         = "Override image tag"

	// Environment variable names
	EnvRegistryUsername = "REGISTRY_USERNAME"
	EnvRegistryPassword = "REGISTRY_PASSWORD"
	EnvRegistryInsecure = "REGISTRY_INSECURE"
)

func init() {
	// Authentication flags
	PushCmd.Flags().StringVarP(&username, "username", "u", "", usernameUsage)
	PushCmd.Flags().StringVarP(&password, "password", "p", "", passwordUsage)

	// TLS configuration
	PushCmd.Flags().BoolVar(&insecure, "insecure", false, insecureUsage)

	// Behavior flags
	PushCmd.Flags().BoolVar(&dryRun, "dry-run", false, dryRunUsage)
	PushCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, verboseUsage)

	// Repository overrides
	PushCmd.Flags().StringVar(&repository, "repository", "", repositoryUsage)
	PushCmd.Flags().StringVar(&tag, "tag", "", tagUsage)

	// Mark password as hidden in help
	PushCmd.Flags().MarkHidden("password")
}

// buildPushOptions constructs PushOptions from flags and environment
func buildPushOptions(cmd *cobra.Command, args []string) (*PushOptions, error) {
	opts := &PushOptions{
		DryRun:   dryRun,
		Verbose:  verbose,
		Insecure: insecure,
	}

	// Parse positional arguments
	if len(args) > 0 {
		opts.ImageRef = args[0]
	}
	if len(args) > 1 {
		opts.RegistryURL = args[1]
	}

	// Authentication: flags override environment variables
	opts.Username = getStringValue(username, os.Getenv(EnvRegistryUsername))
	opts.Password = getStringValue(password, os.Getenv(EnvRegistryPassword))

	// Check insecure from environment if not set via flag
	if !opts.Insecure && os.Getenv(EnvRegistryInsecure) == "true" {
		opts.Insecure = true
	}

	// Repository overrides
	if repository != "" {
		opts.Repository = repository
	}
	if tag != "" {
		opts.Tag = tag
	}

	return opts, nil
}

// getStringValue returns flag value if set, otherwise env value
func getStringValue(flagValue, envValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return envValue
}