// cmd/idpbuilder/main.go
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/cnoe-io/idpbuilder/pkg/cmd"
)

var (
	// Version information, typically set during build
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Set up runtime configuration
	setupRuntime()

	// Create and configure the root command
	rootCmd := cmd.NewRootCmd()

	// Add version information to the root command
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date)

	// Set version template for better formatting
	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}
`)

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// setupRuntime configures the Go runtime for optimal performance
func setupRuntime() {
	// Set GOMAXPROCS to use all available CPU cores
	// This is usually set automatically, but we make it explicit
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// Set garbage collector target percentage
	// Lower values trade CPU for memory usage
	if os.Getenv("GOGC") == "" {
		os.Setenv("GOGC", "100")
	}
}

// init performs any necessary initialization before main
func init() {
	// Configure default behavior for missing environment variables
	ensureRequiredEnvVars()
}

// ensureRequiredEnvVars sets default values for required environment variables
func ensureRequiredEnvVars() {
	// Set default log level if not specified
	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "info")
	}

	// Set default configuration directory
	if os.Getenv("IDPBUILDER_CONFIG_DIR") == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			os.Setenv("IDPBUILDER_CONFIG_DIR", homeDir+"/.idpbuilder")
		}
	}

	// Set default cache directory
	if os.Getenv("IDPBUILDER_CACHE_DIR") == "" {
		cacheDir, err := os.UserCacheDir()
		if err == nil {
			os.Setenv("IDPBUILDER_CACHE_DIR", cacheDir+"/idpbuilder")
		}
	}

	// Set default temporary directory for idpbuilder operations
	if os.Getenv("IDPBUILDER_TEMP_DIR") == "" {
		os.Setenv("IDPBUILDER_TEMP_DIR", os.TempDir()+"/idpbuilder")
	}
}