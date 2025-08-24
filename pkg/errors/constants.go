package errors

import "time"

// Severity levels for error reporting and logging
type Severity string

const (
	// SeverityCritical indicates system-wide failure
	SeverityCritical Severity = "critical"

	// SeverityError indicates operation failure
	SeverityError Severity = "error"

	// SeverityWarning indicates potential issues
	SeverityWarning Severity = "warning"

	// SeverityInfo indicates informational messages
	SeverityInfo Severity = "info"
)

// Common error message templates
const (
	MsgRegistryConnectionFailed = "failed to connect to registry %s"
	MsgAuthenticationRequired   = "authentication required for registry %s"
	MsgManifestNotFound        = "manifest not found for %s:%s"
	MsgInvalidManifest         = "invalid manifest format for %s:%s"
	MsgNetworkTimeout          = "network timeout after %s"
	MsgInvalidConfiguration    = "invalid configuration: %s"
)

// Retry policy constants
const (
	// DefaultMaxRetries is the default number of retry attempts
	DefaultMaxRetries = 3

	// DefaultRetryDelay is the default delay between retries
	DefaultRetryDelay = time.Second * 2

	// DefaultRetryBackoff is the default backoff multiplier
	DefaultRetryBackoff = 2.0

	// MaxRetryDelay is the maximum delay between retries
	MaxRetryDelay = time.Minute * 5
)

// Default error formats for consistent error reporting
const (
	// DefaultErrorFormat is the standard error message format
	DefaultErrorFormat = "[%s] %s: %s"

	// DetailedErrorFormat includes timestamp and context
	DetailedErrorFormat = "[%s] %s at %s: %s (context: %v)"

	// ShortErrorFormat is for brief error messages
	ShortErrorFormat = "%s: %s"
)

// Error stack limits
const (
	// DefaultMaxStackDepth limits error stack depth
	DefaultMaxStackDepth = 10

	// MaxErrorContextSize limits context data size
	MaxErrorContextSize = 1024
)