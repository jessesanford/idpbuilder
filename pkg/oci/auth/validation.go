package auth

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"github.com/go-playground/validator/v10"
)

// AuthValidator provides validation functionality for authentication types.
type AuthValidator struct {
	validator *validator.Validate
}

// NewAuthValidator creates a new AuthValidator instance with custom validators registered.
func NewAuthValidator() *AuthValidator {
	v := validator.New()
	v.RegisterValidation("hostname_port", validateHostnamePort)
	return &AuthValidator{validator: v}
}

// validateHostnamePort validates hostname:port format.
func validateHostnamePort(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return false
	}
	hostname, portStr := parts[0], parts[1]
	if hostname == "" || len(hostname) > 253 {
		return false
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return false
	}
	if net.ParseIP(hostname) != nil {
		return true
	}
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*$`)
	return hostnameRegex.MatchString(hostname)
}

// ValidateCredentials validates credentials with expiry check.
func (av *AuthValidator) ValidateCredentials(creds *Credentials) error {
	if creds == nil {
		return fmt.Errorf("credentials cannot be nil")
	}
	if err := av.validator.Struct(creds); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if creds.IsExpired() {
		return fmt.Errorf("credentials have expired")
	}
	hasBasicAuth := creds.Username != "" || creds.Password != ""
	hasToken := creds.Token != nil
	if hasBasicAuth && hasToken {
		return fmt.Errorf("cannot specify both username/password and token")
	}
	if !hasBasicAuth && !hasToken {
		return fmt.Errorf("must specify either username/password or token")
	}
	return nil
}

// ValidateToken validates a token and its expiration.
func (av *AuthValidator) ValidateToken(token *Token) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}
	if err := av.validator.Struct(token); err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}
	if token.IsExpired() {
		return fmt.Errorf("token has expired")
	}
	return nil
}

// ValidateRegistryURL validates registry URL format.
func (av *AuthValidator) ValidateRegistryURL(registryURL string) error {
	if registryURL == "" {
		return fmt.Errorf("registry URL cannot be empty")
	}
	u, err := url.Parse(registryURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("registry URL must use http or https scheme")
	}
	if u.Host == "" {
		return fmt.Errorf("registry URL must include host")
	}
	return nil
}