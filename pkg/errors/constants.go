// Package errors provides comprehensive error handling for OCI management operations.
// This package defines error categories, codes, and types for consistent error reporting
// across the idpbuilder OCI management system.
package errors

// Error categories define the general classification of errors
const (
	// CategoryTransient represents errors that may resolve with retry
	CategoryTransient = "transient"
	
	// CategoryPermanent represents errors that will not resolve with retry
	CategoryPermanent = "permanent"
	
	// CategoryConfiguration represents errors related to configuration issues
	CategoryConfiguration = "configuration"
	
	// CategoryValidation represents errors related to input validation
	CategoryValidation = "validation"
	
	// CategoryNetwork represents errors related to network operations
	CategoryNetwork = "network"
	
	// CategoryAuthentication represents errors related to authentication
	CategoryAuthentication = "authentication"
	
	// CategoryAuthorization represents errors related to authorization
	CategoryAuthorization = "authorization"
	
	// CategoryResource represents errors related to resource management
	CategoryResource = "resource"
)

// Error severity levels
const (
	// SeverityLow represents minor errors that don't affect core functionality
	SeverityLow = "low"
	
	// SeverityMedium represents errors that affect some functionality
	SeverityMedium = "medium"
	
	// SeverityHigh represents errors that affect core functionality
	SeverityHigh = "high"
	
	// SeverityCritical represents errors that prevent system operation
	SeverityCritical = "critical"
)

// Error context keys for structured error information
const (
	// ContextKeyOperation identifies the operation that failed
	ContextKeyOperation = "operation"
	
	// ContextKeyResource identifies the resource involved in the error
	ContextKeyResource = "resource"
	
	// ContextKeyNamespace identifies the Kubernetes namespace
	ContextKeyNamespace = "namespace"
	
	// ContextKeyRegistry identifies the OCI registry
	ContextKeyRegistry = "registry"
	
	// ContextKeyImage identifies the OCI image
	ContextKeyImage = "image"
	
	// ContextKeyTag identifies the image tag
	ContextKeyTag = "tag"
	
	// ContextKeyRepository identifies the repository
	ContextKeyRepository = "repository"
	
	// ContextKeyURL identifies a URL involved in the operation
	ContextKeyURL = "url"
	
	// ContextKeyPath identifies a file or directory path
	ContextKeyPath = "path"
	
	// ContextKeyTimeout identifies timeout duration
	ContextKeyTimeout = "timeout"
)