package errors

// Build Error Codes (1xxx)
const (
	// ErrCodeBuildFailed indicates a general build failure
	ErrCodeBuildFailed = "1000"
	// ErrCodeBuildTimeout indicates a build timeout
	ErrCodeBuildTimeout = "1001"
	// ErrCodeDockerfileMissing indicates missing Dockerfile
	ErrCodeDockerfileMissing = "1002"
	// ErrCodeDockerfileInvalid indicates invalid Dockerfile syntax
	ErrCodeDockerfileInvalid = "1003"
	// ErrCodeBuildContext indicates build context issues
	ErrCodeBuildContext = "1004"
	// ErrCodeBuildTempFailure indicates temporary build failure
	ErrCodeBuildTempFailure = "1005"
)

// Registry Error Codes (2xxx)
const (
	// ErrCodeRegistryUnreachable indicates registry is unreachable
	ErrCodeRegistryUnreachable = "2000"
	// ErrCodeRegistryTimeout indicates registry connection timeout
	ErrCodeRegistryTimeout = "2001"
	// ErrCodeRegistryAuthFailed indicates registry authentication failure
	ErrCodeRegistryAuthFailed = "2002"
	// ErrCodeRegistryPushFailed indicates push operation failure
	ErrCodeRegistryPushFailed = "2003"
	// ErrCodeRegistryRateLimit indicates rate limiting
	ErrCodeRegistryRateLimit = "2004"
	// ErrCodeRegistryTempFailure indicates temporary registry failure
	ErrCodeRegistryTempFailure = "2005"
)

// Configuration Error Codes (3xxx)
const (
	// ErrCodeConfigMissing indicates missing configuration
	ErrCodeConfigMissing = "3000"
	// ErrCodeConfigInvalid indicates invalid configuration
	ErrCodeConfigInvalid = "3001"
	// ErrCodeConfigVersionMismatch indicates version mismatch
	ErrCodeConfigVersionMismatch = "3002"
	// ErrCodeConfigValidationFailed indicates validation failure
	ErrCodeConfigValidationFailed = "3003"
	// ErrCodeConfigPermissionDenied indicates permission denied
	ErrCodeConfigPermissionDenied = "3004"
)

// Stack Error Codes (4xxx)
const (
	// ErrCodeStackNotFound indicates stack not found
	ErrCodeStackNotFound = "4000"
	// ErrCodeStackDeployFailed indicates stack deployment failure
	ErrCodeStackDeployFailed = "4001"
	// ErrCodeStackValidationFailed indicates stack validation failure
	ErrCodeStackValidationFailed = "4002"
	// ErrCodeStackDependencyFailed indicates dependency failure
	ErrCodeStackDependencyFailed = "4003"
	// ErrCodeStackRollbackFailed indicates rollback failure
	ErrCodeStackRollbackFailed = "4004"
)

// Authentication Error Codes (5xxx)
const (
	// ErrCodeAuthTokenInvalid indicates invalid auth token
	ErrCodeAuthTokenInvalid = "5000"
	// ErrCodeAuthTokenExpired indicates expired auth token
	ErrCodeAuthTokenExpired = "5001"
	// ErrCodeAuthCertInvalid indicates invalid certificate
	ErrCodeAuthCertInvalid = "5002"
	// ErrCodeAuthCertExpired indicates expired certificate
	ErrCodeAuthCertExpired = "5003"
	// ErrCodeAuthPermissionDenied indicates permission denied
	ErrCodeAuthPermissionDenied = "5004"
)

// System Error Codes (6xxx)
const (
	// ErrCodeSystemResourceExhausted indicates resource exhaustion
	ErrCodeSystemResourceExhausted = "6000"
	// ErrCodeSystemDiskFull indicates disk full
	ErrCodeSystemDiskFull = "6001"
	// ErrCodeSystemMemoryExhausted indicates memory exhaustion
	ErrCodeSystemMemoryExhausted = "6002"
	// ErrCodeSystemNetworkError indicates network error
	ErrCodeSystemNetworkError = "6003"
	// ErrCodeSystemTempFailure indicates temporary system failure
	ErrCodeSystemTempFailure = "6004"
)

// ErrorMessages provides human-readable messages for error codes
var ErrorMessages = map[string]string{
	// Build errors
	ErrCodeBuildFailed:        "Build operation failed",
	ErrCodeBuildTimeout:       "Build operation timed out",
	ErrCodeDockerfileMissing:  "Dockerfile not found in build context",
	ErrCodeDockerfileInvalid:  "Invalid Dockerfile syntax or structure",
	ErrCodeBuildContext:       "Invalid or inaccessible build context",
	ErrCodeBuildTempFailure:   "Temporary build failure, retry may succeed",

	// Registry errors
	ErrCodeRegistryUnreachable: "Container registry is unreachable",
	ErrCodeRegistryTimeout:     "Connection to registry timed out",
	ErrCodeRegistryAuthFailed:  "Authentication with registry failed",
	ErrCodeRegistryPushFailed:  "Failed to push image to registry",
	ErrCodeRegistryRateLimit:   "Registry rate limit exceeded",
	ErrCodeRegistryTempFailure: "Temporary registry failure, retry may succeed",

	// Configuration errors
	ErrCodeConfigMissing:           "Required configuration is missing",
	ErrCodeConfigInvalid:           "Configuration contains invalid values",
	ErrCodeConfigVersionMismatch:   "Configuration version is incompatible",
	ErrCodeConfigValidationFailed:  "Configuration failed validation checks",
	ErrCodeConfigPermissionDenied:  "Permission denied accessing configuration",

	// Stack errors
	ErrCodeStackNotFound:           "Stack definition not found",
	ErrCodeStackDeployFailed:       "Stack deployment failed",
	ErrCodeStackValidationFailed:   "Stack validation failed",
	ErrCodeStackDependencyFailed:   "Stack dependency resolution failed",
	ErrCodeStackRollbackFailed:     "Stack rollback operation failed",

	// Authentication errors
	ErrCodeAuthTokenInvalid:      "Authentication token is invalid",
	ErrCodeAuthTokenExpired:      "Authentication token has expired",
	ErrCodeAuthCertInvalid:       "Certificate is invalid or corrupted",
	ErrCodeAuthCertExpired:       "Certificate has expired",
	ErrCodeAuthPermissionDenied:  "Permission denied for the requested operation",

	// System errors
	ErrCodeSystemResourceExhausted: "System resources exhausted",
	ErrCodeSystemDiskFull:          "Disk space exhausted",
	ErrCodeSystemMemoryExhausted:   "System memory exhausted",
	ErrCodeSystemNetworkError:      "Network connectivity error",
	ErrCodeSystemTempFailure:       "Temporary system failure, retry may succeed",
}

// GetErrorMessage returns the human-readable message for an error code
func GetErrorMessage(code string) string {
	if message, exists := ErrorMessages[code]; exists {
		return message
	}
	return "Unknown error occurred"
}