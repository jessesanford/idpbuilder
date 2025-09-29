// Package interfaces defines the core command interfaces for OCI operations.
// These interfaces provide the foundation for implementing consistent CLI commands
// across all OCI-related operations such as push, pull, list, and inspect.
package interfaces

import "context"

// OCICommand defines the base interface that all OCI-related commands must implement.
// This interface ensures consistent behavior and validation across all command types.
type OCICommand interface {
	// Execute runs the command with the provided options.
	// The context allows for cancellation and timeout handling.
	Execute(ctx context.Context, opts CommandOptions) error

	// Validate checks if the provided options are valid for this command.
	// This should be called before Execute to catch configuration errors early.
	Validate(opts CommandOptions) error

	// GetDescription returns a human-readable description of the command.
	// This is used for help text and documentation generation.
	GetDescription() string
}

// PushCommand extends OCICommand with push-specific functionality.
// Implementations handle uploading local content to remote registries.
type PushCommand interface {
	OCICommand

	// SetProgressReporter sets a progress reporter for upload tracking.
	// The reporter will receive updates during the push operation.
	SetProgressReporter(reporter ProgressReporter)

	// ExecutePush performs the push operation with push-specific options.
	// This provides typed options for better type safety and validation.
	ExecutePush(ctx context.Context, opts PushOptions) error

	// ValidatePush validates push-specific options and configuration.
	ValidatePush(opts PushOptions) error

	// GetPushCapabilities returns the capabilities supported by this push implementation.
	// This allows callers to determine what features are available.
	GetPushCapabilities() []string
}

// PullCommand extends OCICommand with pull-specific functionality.
// Implementations handle downloading content from remote registries.
type PullCommand interface {
	OCICommand

	// SetProgressReporter sets a progress reporter for download tracking.
	SetProgressReporter(reporter ProgressReporter)

	// SetExtractor sets the extractor to use for processing pulled content.
	// Different extractors can handle various content types and formats.
	SetExtractor(extractor Extractor)

	// ExecutePull performs the pull operation with pull-specific options.
	ExecutePull(ctx context.Context, opts PullOptions) error

	// ValidatePull validates pull-specific options and configuration.
	ValidatePull(opts PullOptions) error

	// GetPullCapabilities returns the capabilities supported by this pull implementation.
	GetPullCapabilities() []string
}

// ListCommand extends OCICommand with list-specific functionality.
// Implementations handle listing repositories, tags, and other registry content.
type ListCommand interface {
	OCICommand

	// ExecuteList performs the list operation with list-specific options.
	ExecuteList(ctx context.Context, opts ListOptions) error

	// ValidateList validates list-specific options and configuration.
	ValidateList(opts ListOptions) error

	// GetListCapabilities returns the capabilities supported by this list implementation.
	GetListCapabilities() []string

	// SupportsPagination returns true if the implementation supports paginated results.
	SupportsPagination() bool
}

// InspectCommand extends OCICommand with inspect-specific functionality.
// Implementations handle examining metadata and configuration of registry content.
type InspectCommand interface {
	OCICommand

	// ExecuteInspect performs the inspect operation with inspect-specific options.
	ExecuteInspect(ctx context.Context, opts InspectOptions) error

	// ValidateInspect validates inspect-specific options and configuration.
	ValidateInspect(opts InspectOptions) error

	// GetInspectCapabilities returns the capabilities supported by this inspect implementation.
	GetInspectCapabilities() []string

	// GetSupportedFormats returns the output formats supported by this implementation.
	GetSupportedFormats() []string
}

// CommandExecutor orchestrates the execution of OCI commands.
// This interface provides a unified way to execute any type of OCI command
// with consistent error handling and progress reporting.
type CommandExecutor interface {
	// ExecuteCommand runs any OCICommand with the provided context and options.
	// This provides a unified execution path for all command types.
	ExecuteCommand(ctx context.Context, cmd OCICommand, opts CommandOptions) error

	// ValidateCommand validates any OCICommand with the provided options.
	ValidateCommand(cmd OCICommand, opts CommandOptions) error

	// SetDefaultProgressReporter sets a default progress reporter for all commands.
	// Individual commands can override this with their own reporters.
	SetDefaultProgressReporter(reporter ProgressReporter)

	// GetSupportedCommands returns a list of command types this executor supports.
	GetSupportedCommands() []string
}

// CommandFactory creates instances of OCI commands based on command type.
// This interface enables dynamic command creation and plugin-style architectures.
type CommandFactory interface {
	// CreatePushCommand creates a new push command instance.
	CreatePushCommand() (PushCommand, error)

	// CreatePullCommand creates a new pull command instance.
	CreatePullCommand() (PullCommand, error)

	// CreateListCommand creates a new list command instance.
	CreateListCommand() (ListCommand, error)

	// CreateInspectCommand creates a new inspect command instance.
	CreateInspectCommand() (InspectCommand, error)

	// GetSupportedCommandTypes returns the command types this factory can create.
	GetSupportedCommandTypes() []string
}