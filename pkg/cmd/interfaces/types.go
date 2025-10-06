// Package interfaces defines common types and interfaces for CLI commands.
// This package provides the contract layer that enables consistent CLI command
// development across all OCI operations.
package interfaces

// CommandOptions represents common options available to all CLI commands.
// These options provide consistent behavior across different command types.
type CommandOptions struct {
	// Registry URL or configured registry name
	Registry string

	// Authentication options - only one should be used
	Username string // Basic auth username
	Password string // Basic auth password
	Token    string // Bearer token for token-based auth

	// Common behavioral flags
	Insecure bool // Allow insecure connections (skip TLS verification)
	Verbose  bool // Enable verbose output
	DryRun   bool // Perform validation without executing operations
}

// PushOptions contains options specific to push operations.
type PushOptions struct {
	CommandOptions
	Source      string // Local path to push from
	Destination string // Remote destination reference
	Platform    string // Target platform (e.g., "linux/amd64")
	Force       bool   // Force overwrite existing references
}

// PullOptions contains options specific to pull operations.
type PullOptions struct {
	CommandOptions
	Source      string // Remote source reference to pull
	Destination string // Local destination path
	Platform    string // Platform to pull (e.g., "linux/amd64")
	Extract     bool   // Extract contents after pulling
}

// ListOptions contains options specific to list operations.
type ListOptions struct {
	CommandOptions
	Repository string // Repository to list contents from
	Tags       bool   // List tags instead of repositories
	Limit      int    // Maximum number of items to return
}

// InspectOptions contains options specific to inspect operations.
type InspectOptions struct {
	CommandOptions
	Reference string // Reference to inspect
	Raw       bool   // Return raw manifest/config data
	Format    string // Output format (json, yaml, table)
}

// ProgressReporter defines the interface for reporting operation progress.
// Implementations can provide console output, log files, or other progress
// indication mechanisms.
type ProgressReporter interface {
	// ReportProgress reports current progress of an operation
	ReportProgress(current, total int64, message string)

	// Start begins a new progress reporting session
	Start(message string)

	// Complete marks the operation as finished
	Complete(message string)

	// Error reports an error during the operation
	Error(err error)
}

// Extractor defines the interface for extracting pulled content.
// Different extractors can handle various content types and formats.
type Extractor interface {
	// Extract extracts content from source to destination
	Extract(source, destination string) error

	// SupportedFormats returns the formats this extractor can handle
	SupportedFormats() []string
}

// RegistryConfig represents configuration for a container registry.
// This allows managing multiple registries with different settings.
type RegistryConfig struct {
	Name     string // Human-readable name for the registry
	URL      string // Registry base URL
	Type     string // Registry type: "dockerhub", "harbor", "gitea", "generic"
	AuthType string // Authentication type: "basic", "token", "oauth"
	Insecure bool   // Allow insecure connections
}

// RegistryInfo contains metadata about a registry.
type RegistryInfo struct {
	Name         string            // Registry name
	URL          string            // Registry URL
	Version      string            // Registry version
	Capabilities []string          // Supported features/APIs
	Metadata     map[string]string // Additional registry-specific metadata
}

// ErrorType represents different categories of CLI errors for consistent handling.
type ErrorType string

const (
	// ErrorTypeValidation indicates input validation errors
	ErrorTypeValidation ErrorType = "validation"

	// ErrorTypeAuthentication indicates authentication/authorization errors
	ErrorTypeAuthentication ErrorType = "authentication"

	// ErrorTypeNetwork indicates network connectivity errors
	ErrorTypeNetwork ErrorType = "network"

	// ErrorTypeRegistry indicates registry-specific errors
	ErrorTypeRegistry ErrorType = "registry"

	// ErrorTypeFileSystem indicates local file system errors
	ErrorTypeFileSystem ErrorType = "filesystem"
)

// CLIError represents a structured error with additional context for CLI operations.
type CLIError struct {
	Type    ErrorType // Category of error
	Message string    // Human-readable error message
	Cause   error     // Underlying error (if any)
	Context string    // Additional context about when/where the error occurred
}

// Error implements the error interface.
func (e *CLIError) Error() string {
	if e.Context != "" {
		return e.Context + ": " + e.Message
	}
	return e.Message
}

// Unwrap returns the underlying error for error unwrapping.
func (e *CLIError) Unwrap() error {
	return e.Cause
}