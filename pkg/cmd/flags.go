package cmd

import "github.com/spf13/cobra"

var (
	// Global flags
	configFile string
	verbose    bool
	quiet      bool
)

// AddGlobalFlags adds global flags to the root command
func AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file path")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")
}

// AddBuildFlags adds build-specific flags to a command
func AddBuildFlags(cmd *cobra.Command) {
	cmd.Flags().String("exclude", "", "Exclusion patterns (comma-separated)")
	cmd.Flags().String("platform", "linux/amd64", "Target platform")
	cmd.Flags().String("cache-dir", "", "Build cache directory")
	cmd.Flags().Bool("no-cache", false, "Disable build cache")
}

// AddPushFlags adds push-specific flags to a command
func AddPushFlags(cmd *cobra.Command) {
	cmd.Flags().String("registry", "", "Registry URL (default: auto-detect)")
	cmd.Flags().Int("retry", 3, "Number of retry attempts")
	cmd.Flags().String("auth-config", "", "Path to authentication config file")
	cmd.Flags().Bool("skip-tls-verify", false, "Skip TLS certificate verification")
}

// AddRegistryFlags adds common registry-related flags
func AddRegistryFlags(cmd *cobra.Command) {
	cmd.Flags().String("username", "", "Registry username")
	cmd.Flags().String("password", "", "Registry password")
	cmd.Flags().String("registry-config", "", "Registry configuration file path")
}

// AddOutputFlags adds output-related flags
func AddOutputFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("output", "o", "", "Output format or file path")
	cmd.Flags().Bool("json", false, "Output in JSON format")
	cmd.Flags().Bool("yaml", false, "Output in YAML format")
}

// GetGlobalFlags returns the global flag values
func GetGlobalFlags() (string, bool, bool) {
	return configFile, verbose, quiet
}
