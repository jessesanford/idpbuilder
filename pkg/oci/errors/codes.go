package errors

const (
	CodeBuildFailed                 = "1001"
	CodeBuildTimeout                = "1002"
	CodeDockerfileInvalid           = "1003"
	CodeBuildContextInvalid         = "1004"
	CodeRegistryConnectionFailed    = "2001"
	CodeRegistryAuthFailed          = "2002"
	CodeRegistryPushFailed          = "2003"
	CodeRegistryPullFailed          = "2004"
	CodeRegistryNotFound            = "2005"
	CodeConfigInvalid               = "3001"
	CodeConfigMissing               = "3002"
	CodeConfigParseError            = "3003"
	CodeStackDeployFailed           = "4001"
	CodeStackDeleteFailed           = "4002"
	CodeStackNotFound               = "4003"
	CodeStackResourceConflict       = "4004"
	CodeAuthTokenInvalid            = "5001"
	CodeAuthTokenExpired            = "5002"
	CodeAuthCertInvalid             = "5003"
	CodeAuthPermissionDenied        = "5004"
	CodeSystemStorageError          = "6001"
	CodeSystemPermissionError       = "6002"
	CodeSystemResourceLimitExceeded = "6003"
	CodeSystemNetworkError          = "6004"
)

var errorMessages = map[string]string{
	CodeBuildFailed:                 "Build operation failed",
	CodeBuildTimeout:                "Build operation timed out",
	CodeDockerfileInvalid:           "Invalid or missing Dockerfile",
	CodeBuildContextInvalid:         "Invalid build context",
	CodeRegistryConnectionFailed:    "Failed to connect to registry",
	CodeRegistryAuthFailed:          "Registry authentication failed",
	CodeRegistryPushFailed:          "Failed to push to registry",
	CodeRegistryPullFailed:          "Failed to pull from registry",
	CodeRegistryNotFound:            "Registry or repository not found",
	CodeConfigInvalid:               "Invalid configuration",
	CodeConfigMissing:               "Missing required configuration",
	CodeConfigParseError:            "Configuration parsing error",
	CodeStackDeployFailed:           "Stack deployment failed",
	CodeStackDeleteFailed:           "Stack deletion failed",
	CodeStackNotFound:               "Stack not found",
	CodeStackResourceConflict:       "Resource conflict in stack",
	CodeAuthTokenInvalid:            "Invalid authentication token",
	CodeAuthTokenExpired:            "Authentication token expired",
	CodeAuthCertInvalid:             "Invalid certificate",
	CodeAuthPermissionDenied:        "Permission denied",
	CodeSystemStorageError:          "Storage error",
	CodeSystemPermissionError:       "Permission error",
	CodeSystemResourceLimitExceeded: "Resource limit exceeded",
	CodeSystemNetworkError:          "Network error",
}

func GetErrorMessage(code string) string {
	if message, ok := errorMessages[code]; ok {
		return message
	}
	return "Unknown error"
}
