package errors

import (
	"fmt"
)

// errorMessages defines simple message templates
var errorMessages = map[ErrorType]string{
	ErrCertNotFound:   "Certificate not found at expected location: %s",
	ErrCertExpired:    "Certificate expired %s days ago",
	ErrCertUntrusted:  "Certificate issuer '%s' not in trust store",
	ErrCertMismatch:   "Certificate CN '%s' doesn't match registry '%s'",
	ErrCertPermission: "Permission denied accessing certificate at %s",
}

// BuildMessage creates a formatted error message
func BuildMessage(errorType ErrorType, args ...interface{}) string {
	template, exists := errorMessages[errorType]
	if !exists {
		return fmt.Sprintf("Unknown error type: %s", errorType)
	}
	
	if len(args) == 0 {
		return template
	}
	
	return fmt.Sprintf(template, args...)
}

// UpdateErrorMessage updates the message in a CertificateError
func UpdateErrorMessage(err *CertificateError) {
	switch err.Type {
	case ErrCertNotFound:
		if path := err.Details["path"]; path != "" {
			err.Message = BuildMessage(err.Type, path)
		}
	case ErrCertExpired:
		if daysAgo := err.Details["days_ago"]; daysAgo != "" {
			err.Message = BuildMessage(err.Type, daysAgo)
		}
	case ErrCertUntrusted:
		if issuer := err.Details["issuer"]; issuer != "" {
			err.Message = BuildMessage(err.Type, issuer)
		}
	case ErrCertMismatch:
		if cn, registry := err.Details["cn"], err.Details["registry"]; cn != "" && registry != "" {
			err.Message = BuildMessage(err.Type, cn, registry)
		}
	case ErrCertPermission:
		if path := err.Details["path"]; path != "" {
			err.Message = BuildMessage(err.Type, path)
		}
	}
}