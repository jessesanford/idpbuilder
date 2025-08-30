package errors

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
)

// errorTemplates defines message templates for each error type
var errorTemplates = map[ErrorType]string{
	ErrCertNotFound:    "Certificate not found at expected location: {{.Path}}",
	ErrCertInvalid:     "Certificate validation failed: {{.Reason}}",
	ErrCertExpired:     "Certificate expired on {{.ExpiryDate}} ({{.DaysAgo}} days ago)",
	ErrCertUntrusted:   "Certificate issuer '{{.Issuer}}' not in trust store",
	ErrCertMismatch:    "Certificate CN '{{.CN}}' doesn't match registry '{{.Registry}}'",
	ErrCertPermission:  "Permission denied accessing certificate at {{.Path}}: {{.Error}}",
	ErrCertChainBroken: "Certificate chain incomplete: missing {{.MissingLink}}",
	ErrCertFormat:      "Unsupported certificate format: expected {{.Expected}}, got {{.Actual}}",
}

// MessageBuilder handles error message template processing
type MessageBuilder struct {
	templates map[ErrorType]*template.Template
	mu        sync.RWMutex
}

// NewMessageBuilder creates a new message builder with pre-parsed templates
func NewMessageBuilder() *MessageBuilder {
	mb := &MessageBuilder{
		templates: make(map[ErrorType]*template.Template),
	}
	
	// Parse all templates at initialization
	for errorType, templateStr := range errorTemplates {
		tmpl, err := template.New(string(errorType)).Parse(templateStr)
		if err != nil {
			// If template parsing fails, create a simple fallback
			tmpl, _ = template.New(string(errorType)).Parse("Template error: {{.Error}}")
		}
		mb.templates[errorType] = tmpl
	}
	
	return mb
}

// defaultMessageBuilder is the global message builder instance
var defaultMessageBuilder = NewMessageBuilder()

// BuildMessage creates a formatted error message using template data
func (mb *MessageBuilder) BuildMessage(errorType ErrorType, data interface{}) (string, error) {
	mb.mu.RLock()
	tmpl, exists := mb.templates[errorType]
	mb.mu.RUnlock()
	
	if !exists {
		return fmt.Sprintf("Unknown error type: %s", errorType), fmt.Errorf("no template found for error type: %s", errorType)
	}
	
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		// Fallback to basic error message if template execution fails
		return fmt.Sprintf("Error executing template for %s: %v", errorType, err), err
	}
	
	return buf.String(), nil
}

// AddTemplate adds or updates a template for a specific error type
func (mb *MessageBuilder) AddTemplate(errorType ErrorType, templateStr string) error {
	tmpl, err := template.New(string(errorType)).Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template for %s: %w", errorType, err)
	}
	
	mb.mu.Lock()
	mb.templates[errorType] = tmpl
	mb.mu.Unlock()
	
	return nil
}

// GetTemplate returns the raw template string for an error type
func (mb *MessageBuilder) GetTemplate(errorType ErrorType) (string, bool) {
	template, exists := errorTemplates[errorType]
	return template, exists
}

// MessageData provides common data structure for template rendering
type MessageData struct {
	Path        string `json:"path,omitempty"`
	Reason      string `json:"reason,omitempty"`
	ExpiryDate  string `json:"expiry_date,omitempty"`
	DaysAgo     string `json:"days_ago,omitempty"`
	Issuer      string `json:"issuer,omitempty"`
	CN          string `json:"cn,omitempty"`
	Registry    string `json:"registry,omitempty"`
	Error       string `json:"error,omitempty"`
	MissingLink string `json:"missing_link,omitempty"`
	Expected    string `json:"expected,omitempty"`
	Actual      string `json:"actual,omitempty"`
}

// NewMessageDataFromError extracts template data from a CertificateError
func NewMessageDataFromError(err *CertificateError) *MessageData {
	data := &MessageData{}
	
	if err.Details != nil {
		data.Path = err.Details["path"]
		data.Reason = err.Details["reason"]
		data.ExpiryDate = err.Details["expiry_date"]
		data.DaysAgo = err.Details["days_ago"]
		data.Issuer = err.Details["issuer"]
		data.CN = err.Details["cn"]
		data.Registry = err.Details["registry"]
		data.MissingLink = err.Details["missing_link"]
		data.Expected = err.Details["expected"]
		data.Actual = err.Details["actual"]
	}
	
	if err.OriginalErr != nil {
		data.Error = err.OriginalErr.Error()
	}
	
	return data
}

// BuildMessageFromError creates a formatted message from a CertificateError
func (mb *MessageBuilder) BuildMessageFromError(err *CertificateError) (string, error) {
	data := NewMessageDataFromError(err)
	return mb.BuildMessage(err.Type, data)
}

// Global convenience functions using the default message builder

// BuildMessage builds a message using the default message builder
func BuildMessage(errorType ErrorType, data interface{}) (string, error) {
	return defaultMessageBuilder.BuildMessage(errorType, data)
}

// BuildMessageFromError builds a message from a CertificateError using the default builder
func BuildMessageFromError(err *CertificateError) (string, error) {
	return defaultMessageBuilder.BuildMessageFromError(err)
}

// AddTemplate adds a template to the default message builder
func AddTemplate(errorType ErrorType, templateStr string) error {
	return defaultMessageBuilder.AddTemplate(errorType, templateStr)
}

// GetTemplate gets a template from the default message builder
func GetTemplate(errorType ErrorType) (string, bool) {
	return defaultMessageBuilder.GetTemplate(errorType)
}

// UpdateErrorMessage updates the message in a CertificateError using the template system
func UpdateErrorMessage(err *CertificateError) error {
	message, buildErr := BuildMessageFromError(err)
	if buildErr != nil {
		return fmt.Errorf("failed to build message for error %s: %w", err.Type, buildErr)
	}
	
	err.Message = message
	return nil
}