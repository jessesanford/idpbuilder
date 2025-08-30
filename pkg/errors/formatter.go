package errors

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// OutputFormat specifies the format for error output
type OutputFormat string

const (
	FormatConsole OutputFormat = "console" // Human-readable console output
	FormatJSON    OutputFormat = "json"    // JSON format for automation
	FormatYAML    OutputFormat = "yaml"    // YAML format for configuration
)

// FormatOptions controls how errors are formatted
type FormatOptions struct {
	UseColor     bool         `json:"use_color"`     // Enable terminal colors
	Verbose      bool         `json:"verbose"`       // Include full context
	Format       OutputFormat `json:"format"`        // Output format
	ShowGuidance bool         `json:"show_guidance"` // Include resolution steps
	ShowContext  bool         `json:"show_context"`  // Include system context
}

// DefaultFormatOptions returns sensible default formatting options
func DefaultFormatOptions() *FormatOptions {
	return &FormatOptions{
		UseColor:     true,
		Verbose:      false,
		Format:       FormatConsole,
		ShowGuidance: true,
		ShowContext:  false,
	}
}

// ErrorFormatter provides methods for formatting certificate errors
type ErrorFormatter interface {
	Format(err *CertificateError, opts *FormatOptions) (string, error)
}

// ConsoleFormatter formats errors for console output
type ConsoleFormatter struct{}

// Format formats a CertificateError for console display
func (f *ConsoleFormatter) Format(err *CertificateError, opts *FormatOptions) (string, error) {
	var sb strings.Builder
	
	// Error header with severity symbol
	symbol := f.getSeveritySymbol(err.Severity, opts.UseColor)
	title := f.getColoredTitle(err.Type, err.Severity, opts.UseColor)
	
	sb.WriteString(fmt.Sprintf("%s %s [%s]\n", symbol, title, err.Type))
	sb.WriteString(strings.Repeat("━", 60) + "\n")
	
	// Main error message
	sb.WriteString(fmt.Sprintf("%s\n", err.Message))
	
	// Error details
	if len(err.Details) > 0 {
		sb.WriteString(fmt.Sprintf("\n%s Details:\n", f.getIcon("info", opts.UseColor)))
		for key, value := range err.Details {
			if value != "" && !strings.Contains(key, "correlation_id") && !strings.Contains(key, "context_timestamp") {
				sb.WriteString(fmt.Sprintf("  • %s: %s\n", strings.Title(strings.ReplaceAll(key, "_", " ")), value))
			}
		}
	}
	
	// Resolution guidance
	if opts.ShowGuidance && err.Resolution != "" {
		sb.WriteString(fmt.Sprintf("\n%s Resolution:\n", f.getIcon("lightbulb", opts.UseColor)))
		sb.WriteString(f.indentText(err.Resolution, "  "))
	} else if opts.ShowGuidance {
		if guide, exists := GetResolutionGuide(err.Type); exists {
			sb.WriteString(fmt.Sprintf("\n%s Resolution Steps:\n", f.getIcon("lightbulb", opts.UseColor)))
			for i, step := range guide.Steps {
				sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, step))
			}
			
			if len(guide.Examples) > 0 {
				sb.WriteString(fmt.Sprintf("\n%s Example Commands:\n", f.getIcon("terminal", opts.UseColor)))
				for _, example := range guide.Examples {
					sb.WriteString(fmt.Sprintf("  $ %s\n", example))
				}
			}
			
			if guide.Workaround != "" {
				sb.WriteString(fmt.Sprintf("\n%s Workaround: %s\n", f.getIcon("warning", opts.UseColor), guide.Workaround))
			}
			
			if len(guide.DocLinks) > 0 {
				sb.WriteString(fmt.Sprintf("\n%s Documentation:\n", f.getIcon("book", opts.UseColor)))
				for _, link := range guide.DocLinks {
					sb.WriteString(fmt.Sprintf("  %s\n", link))
				}
			}
		}
	}
	
	// Context information (if verbose)
	if opts.ShowContext && opts.Verbose {
		sb.WriteString(fmt.Sprintf("\n%s Context:\n", f.getIcon("info", opts.UseColor)))
		sb.WriteString(fmt.Sprintf("  • Component: %s\n", err.Component))
		sb.WriteString(fmt.Sprintf("  • Operation: %s\n", err.Operation))
		sb.WriteString(fmt.Sprintf("  • Timestamp: %s\n", err.Timestamp.Format(time.RFC3339)))
		if corrID := err.Details["correlation_id"]; corrID != "" {
			sb.WriteString(fmt.Sprintf("  • Correlation ID: %s\n", corrID))
		}
	}
	
	sb.WriteString(strings.Repeat("━", 60))
	return sb.String(), nil
}

// getSeveritySymbol returns the appropriate symbol for error severity
func (f *ConsoleFormatter) getSeveritySymbol(severity Severity, useColor bool) string {
	symbols := map[Severity]string{
		SeverityInfo:     "ℹ️",
		SeverityWarning:  "⚠️",
		SeverityError:    "❌",
		SeverityCritical: "🚨",
	}
	
	if symbol, exists := symbols[severity]; exists {
		return symbol
	}
	return "❌"
}

// getColoredTitle returns a colored title based on severity
func (f *ConsoleFormatter) getColoredTitle(errorType ErrorType, severity Severity, useColor bool) string {
	title := "CERTIFICATE ERROR"
	if !useColor {
		return title
	}
	
	// ANSI color codes
	colors := map[Severity]string{
		SeverityInfo:     "\033[36m", // Cyan
		SeverityWarning:  "\033[33m", // Yellow
		SeverityError:    "\033[31m", // Red
		SeverityCritical: "\033[35m", // Magenta
	}
	reset := "\033[0m"
	
	if color, exists := colors[severity]; exists {
		return fmt.Sprintf("%s%s%s", color, title, reset)
	}
	return title
}

// getIcon returns the appropriate icon for different sections
func (f *ConsoleFormatter) getIcon(iconType string, useColor bool) string {
	icons := map[string]string{
		"info":      "📍",
		"lightbulb": "💡",
		"terminal":  "📋",
		"warning":   "⚠️",
		"book":      "📚",
	}
	
	if icon, exists := icons[iconType]; exists {
		return icon
	}
	return "•"
}

// indentText indents each line of text with the given prefix
func (f *ConsoleFormatter) indentText(text, prefix string) string {
	lines := strings.Split(text, "\n")
	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(prefix + line + "\n")
	}
	return sb.String()
}

// JSONFormatter formats errors as JSON
type JSONFormatter struct{}

// Format formats a CertificateError as JSON
func (f *JSONFormatter) Format(err *CertificateError, opts *FormatOptions) (string, error) {
	// Create a serializable representation
	errorData := map[string]interface{}{
		"type":        err.Type,
		"message":     err.Message,
		"details":     err.Details,
		"severity":    err.Severity,
		"timestamp":   err.Timestamp,
		"component":   err.Component,
		"operation":   err.Operation,
	}
	
	if opts.ShowGuidance {
		if guide, exists := GetResolutionGuide(err.Type); exists {
			errorData["resolution"] = map[string]interface{}{
				"steps":      guide.Steps,
				"examples":   guide.Examples,
				"workaround": guide.Workaround,
				"doc_links":  guide.DocLinks,
				"severity":   guide.Severity,
			}
		}
	}
	
	if err.OriginalErr != nil {
		errorData["original_error"] = err.OriginalErr.Error()
	}
	
	data, err := json.MarshalIndent(errorData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal error to JSON: %w", err)
	}
	
	return string(data), nil
}

// YAMLFormatter formats errors as YAML
type YAMLFormatter struct{}

// Format formats a CertificateError as YAML
func (f *YAMLFormatter) Format(err *CertificateError, opts *FormatOptions) (string, error) {
	// Create a serializable representation
	errorData := map[string]interface{}{
		"type":        string(err.Type),
		"message":     err.Message,
		"details":     err.Details,
		"severity":    string(err.Severity),
		"timestamp":   err.Timestamp.Format(time.RFC3339),
		"component":   err.Component,
		"operation":   err.Operation,
	}
	
	if opts.ShowGuidance {
		if guide, exists := GetResolutionGuide(err.Type); exists {
			errorData["resolution"] = map[string]interface{}{
				"steps":      guide.Steps,
				"examples":   guide.Examples,
				"workaround": guide.Workaround,
				"doc_links":  guide.DocLinks,
				"severity":   guide.Severity,
			}
		}
	}
	
	if err.OriginalErr != nil {
		errorData["original_error"] = err.OriginalErr.Error()
	}
	
	data, yamlErr := yaml.Marshal(errorData)
	if yamlErr != nil {
		return "", fmt.Errorf("failed to marshal error to YAML: %w", yamlErr)
	}
	
	return string(data), nil
}

// Format formats a CertificateError according to the specified options
func Format(err *CertificateError, opts *FormatOptions) (string, error) {
	if opts == nil {
		opts = DefaultFormatOptions()
	}
	
	var formatter ErrorFormatter
	
	switch opts.Format {
	case FormatJSON:
		formatter = &JSONFormatter{}
	case FormatYAML:
		formatter = &YAMLFormatter{}
	default:
		formatter = &ConsoleFormatter{}
	}
	
	return formatter.Format(err, opts)
}

// FormatMultiple formats multiple certificate errors
func FormatMultiple(errors []*CertificateError, opts *FormatOptions) (string, error) {
	if len(errors) == 0 {
		return "", nil
	}
	
	if opts == nil {
		opts = DefaultFormatOptions()
	}
	
	switch opts.Format {
	case FormatJSON:
		var errorList []string
		for _, err := range errors {
			formatted, formatErr := Format(err, opts)
			if formatErr != nil {
				return "", formatErr
			}
			errorList = append(errorList, formatted)
		}
		return fmt.Sprintf("[%s]", strings.Join(errorList, ",")), nil
		
	case FormatYAML:
		var errorList []string
		for i, err := range errors {
			formatted, formatErr := Format(err, opts)
			if formatErr != nil {
				return "", formatErr
			}
			errorList = append(errorList, fmt.Sprintf("error_%d:\n%s", i+1, f.indentYAML(formatted)))
		}
		return strings.Join(errorList, "\n"), nil
		
	default:
		var formattedErrors []string
		for _, err := range errors {
			formatted, formatErr := Format(err, opts)
			if formatErr != nil {
				return "", formatErr
			}
			formattedErrors = append(formattedErrors, formatted)
		}
		return strings.Join(formattedErrors, "\n\n"), nil
	}
}

// indentYAML helper function for YAML formatting
func (f *YAMLFormatter) indentYAML(text string) string {
	lines := strings.Split(text, "\n")
	var sb strings.Builder
	for _, line := range lines {
		if line != "" {
			sb.WriteString("  " + line + "\n")
		}
	}
	return sb.String()
}