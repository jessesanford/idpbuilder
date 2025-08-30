package errors

import (
	"fmt"
	"strings"
)

// FormatConsole formats a CertificateError for console display
func FormatConsole(err *CertificateError) string {
	var sb strings.Builder
	
	// Error header
	symbol := getSeveritySymbol(err.Severity)
	sb.WriteString(fmt.Sprintf("%s CERTIFICATE ERROR [%s]\n", symbol, err.Type))
	sb.WriteString(strings.Repeat("━", 50) + "\n")
	
	// Main error message
	sb.WriteString(fmt.Sprintf("%s\n", err.Message))
	
	// Error details (only important ones)
	if len(err.Details) > 0 {
		sb.WriteString("\n📍 Details:\n")
		for key, value := range err.Details {
			if value != "" && isImportantDetail(key) {
				displayKey := strings.Title(strings.ReplaceAll(key, "_", " "))
				sb.WriteString(fmt.Sprintf("  • %s: %s\n", displayKey, value))
			}
		}
	}
	
	// Resolution guidance
	if resolution := err.Details["resolution"]; resolution != "" {
		sb.WriteString(fmt.Sprintf("\n💡 %s", resolution))
	} else {
		resolution := FormatResolution(err.Type)
		sb.WriteString(fmt.Sprintf("\n💡 %s", resolution))
	}
	
	sb.WriteString(strings.Repeat("━", 50))
	return sb.String()
}

// getSeveritySymbol returns the appropriate symbol for error severity
func getSeveritySymbol(severity Severity) string {
	switch severity {
	case SeverityWarning:
		return "⚠️"
	case SeverityError:
		return "❌"
	default:
		return "❌"
	}
}

// isImportantDetail filters which details to show
func isImportantDetail(key string) bool {
	importantKeys := []string{"path", "cn", "registry", "issuer", "days_ago", "expiry_date"}
	for _, important := range importantKeys {
		if key == important {
			return true
		}
	}
	return false
}

// Format formats a CertificateError for display (defaults to console)
func Format(err *CertificateError) string {
	return FormatConsole(err)
}

// FormatMultiple formats multiple certificate errors
func FormatMultiple(errors []*CertificateError) string {
	if len(errors) == 0 {
		return ""
	}
	
	var formatted []string
	for _, err := range errors {
		formatted = append(formatted, Format(err))
	}
	return strings.Join(formatted, "\n\n")
}