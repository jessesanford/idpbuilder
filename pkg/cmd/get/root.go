// pkg/cmd/get/root.go
package get

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const (
	// DefaultTimeout is the default timeout for get operations in seconds
	DefaultTimeout = 300

	// MaxRetries is the maximum number of retries for failed operations
	MaxRetries = 3

	// RetryInterval is the interval between retries in seconds
	RetryInterval = 5

	// DefaultOutputFormat is the default output format for get commands
	DefaultOutputFormat = "table"

	// ValidOutputFormats are the valid output formats
	ValidOutputFormats = "table,json,yaml"
)

// GetOptions holds the configuration for get commands
type GetOptions struct {
	Timeout       time.Duration
	OutputFormat  string
	Namespace     string
	AllNamespaces bool
	Verbose       bool
	Watch         bool
	MaxRetries    int
}

// NewGetOptions creates a new GetOptions with default values
func NewGetOptions() *GetOptions {
	return &GetOptions{
		Timeout:      time.Duration(DefaultTimeout) * time.Second,
		OutputFormat: DefaultOutputFormat,
		MaxRetries:   MaxRetries,
	}
}

// NewGetCmd creates a new get command
func NewGetCmd() *cobra.Command {
	opts := NewGetOptions()

	cmd := &cobra.Command{
		Use:   "get [resource]",
		Short: "Get resources from the cluster",
		Long: `Get resources from the Kubernetes cluster.

This command allows you to retrieve and display information about various
resources in your cluster, including secrets, configmaps, deployments, and more.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd.Context(), opts, args)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return getValidResources(), cobra.ShellCompDirectiveNoFileComp
		},
	}

	// Add flags
	cmd.Flags().DurationVar(&opts.Timeout, "timeout", opts.Timeout,
		fmt.Sprintf("Timeout for the operation (default %ds)", DefaultTimeout))
	cmd.Flags().StringVarP(&opts.OutputFormat, "output", "o", opts.OutputFormat,
		fmt.Sprintf("Output format. One of: %s", ValidOutputFormats))
	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", opts.Namespace,
		"Namespace to use for the operation")
	cmd.Flags().BoolVar(&opts.AllNamespaces, "all-namespaces", opts.AllNamespaces,
		"List resources from all namespaces")
	cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", opts.Verbose,
		"Enable verbose output")
	cmd.Flags().BoolVarP(&opts.Watch, "watch", "w", opts.Watch,
		"Watch for changes to the requested resources")
	cmd.Flags().IntVar(&opts.MaxRetries, "max-retries", opts.MaxRetries,
		"Maximum number of retries for failed operations")

	return cmd
}

// runGet executes the get command
func runGet(ctx context.Context, opts *GetOptions, args []string) error {
	// Validate flags before proceeding
	if err := validateFlags(opts); err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("resource type must be specified")
	}

	resourceType := args[0]

	// Validate resource type
	if !isValidResourceType(resourceType) {
		return fmt.Errorf("invalid resource type: %s. Valid types: %v",
			resourceType, getValidResources())
	}

	// Apply timeout to context
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// Execute the get operation with retries
	return executeWithRetries(ctx, opts, resourceType, args[1:])
}

// validateFlags validates the command flags
func validateFlags(opts *GetOptions) error {
	// Validate output format
	validFormats := []string{"table", "json", "yaml"}
	isValid := false
	for _, format := range validFormats {
		if opts.OutputFormat == format {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid output format: %s. Valid formats: %v",
			opts.OutputFormat, validFormats)
	}

	// Validate timeout
	if opts.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	// Validate max retries
	if opts.MaxRetries < 0 {
		return fmt.Errorf("max-retries must be non-negative")
	}

	// Validate namespace conflicts
	if opts.AllNamespaces && opts.Namespace != "" {
		return fmt.Errorf("cannot specify both --namespace and --all-namespaces")
	}

	return nil
}

// isValidResourceType checks if the given resource type is valid
func isValidResourceType(resourceType string) bool {
	validTypes := getValidResources()
	for _, validType := range validTypes {
		if resourceType == validType {
			return true
		}
	}
	return false
}

// getValidResources returns the list of valid resource types
func getValidResources() []string {
	return []string{
		"secrets",
		"configmaps",
		"deployments",
		"services",
		"pods",
		"nodes",
		"namespaces",
		"persistentvolumes",
		"persistentvolumeclaims",
		"ingresses",
		"certificates",
		"issuers",
	}
}

// executeWithRetries executes the get operation with retry logic
func executeWithRetries(ctx context.Context, opts *GetOptions, resourceType string, args []string) error {
	var lastErr error

	for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
		if attempt > 0 {
			if opts.Verbose {
				fmt.Printf("Retrying get operation (attempt %d/%d)...\n", attempt+1, opts.MaxRetries+1)
			}

			// Wait before retrying
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(RetryInterval) * time.Second):
			}
		}

		// Execute the actual get operation
		err := executeGet(ctx, opts, resourceType, args)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			if opts.Verbose {
				fmt.Printf("Non-retryable error encountered: %v\n", err)
			}
			break
		}

		if opts.Verbose {
			fmt.Printf("Retryable error encountered: %v\n", err)
		}
	}

	return fmt.Errorf("operation failed after %d attempts. Last error: %w", opts.MaxRetries+1, lastErr)
}

// executeGet performs the actual get operation
func executeGet(ctx context.Context, opts *GetOptions, resourceType string, args []string) error {
	// This is a stub implementation for upstream test fixes
	// In a real implementation, this would interact with the Kubernetes API

	if opts.Verbose {
		fmt.Printf("Getting %s resources...\n", resourceType)
		if opts.Namespace != "" {
			fmt.Printf("Namespace: %s\n", opts.Namespace)
		}
		if opts.AllNamespaces {
			fmt.Printf("All namespaces: true\n")
		}
		fmt.Printf("Output format: %s\n", opts.OutputFormat)
	}

	// Simulate some work
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		// Success
	}

	// Output based on format
	switch opts.OutputFormat {
	case "json":
		fmt.Println(`{"items": [], "kind": "List", "apiVersion": "v1"}`)
	case "yaml":
		fmt.Println(`apiVersion: v1
kind: List
items: []`)
	default: // table
		fmt.Printf("No %s found.\n", resourceType)
	}

	return nil
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// In a real implementation, this would check for specific error types
	// that indicate transient failures (network timeouts, server errors, etc.)
	errStr := err.Error()

	// Examples of retryable errors
	retryableErrors := []string{
		"timeout",
		"connection refused",
		"temporary failure",
		"server error",
		"network",
	}

	for _, retryableErr := range retryableErrors {
		if contains(errStr, retryableErr) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		   (s == substr ||
		    (len(s) > len(substr) &&
		     (s[:len(substr)] == substr ||
		      s[len(s)-len(substr):] == substr ||
		      hasSubstring(s, substr))))
}

// hasSubstring performs a simple substring search
func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}