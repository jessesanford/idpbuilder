package validator

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

var (
	// registryPattern validates registry hostnames (domain or IP with optional port)
	registryPattern = regexp.MustCompile(`^([a-zA-Z0-9.-]+|\[[0-9a-fA-F:]+\])(:[0-9]{1,5})?$`)

	// Private IP ranges (SSRF protection)
	privateIPRanges = []string{
		"10.0.0.0/8",     // Private class A
		"172.16.0.0/12",  // Private class B
		"192.168.0.0/16", // Private class C
		"127.0.0.0/8",    // Loopback
		"169.254.0.0/16", // Link-local
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link-local
	}
)

// ValidateRegistryURL validates a registry URL or hostname
//
// Validation Rules:
// 1. Non-empty string required
// 2. No shell metacharacters (command injection prevention)
// 3. Must be valid hostname, IP, or full URL
// 4. SSRF warning for private IP ranges (non-blocking)
//
// Returns:
//   - nil if valid
//   - ValidationError if invalid format or contains dangerous characters
//   - SSRFWarning if private IP detected (allows continuation)
func ValidateRegistryURL(registry string) error {
	if registry == "" {
		return &ValidationError{
			Field:      "registry",
			Message:    "registry URL cannot be empty",
			Suggestion: "provide a registry URL like 'docker.io' or 'localhost:5000'",
			ExitCode:   1,
		}
	}

	// Check for dangerous characters
	if containsAnyChar(registry, dangerousChars) {
		return &ValidationError{
			Field:      "registry",
			Message:    fmt.Sprintf("registry URL contains shell metacharacters: %s", registry),
			Suggestion: "use only alphanumeric characters, dots, hyphens, colons, and brackets",
			ExitCode:   1,
		}
	}

	// Parse URL or treat as hostname
	var hostname string

	// Try parsing as full URL first
	if strings.Contains(registry, "://") {
		parsedURL, err := url.Parse(registry)
		if err != nil {
			return &ValidationError{
				Field:      "registry",
				Message:    fmt.Sprintf("invalid registry URL: %v", err),
				Suggestion: "use format: hostname[:port] or https://hostname[:port]",
				ExitCode:   1,
			}
		}
		hostname = parsedURL.Hostname()
	} else {
		// Check if it looks like an IPv6 address (contains multiple colons)
		colonCount := strings.Count(registry, ":")
		if colonCount > 1 {
			// Likely IPv6 - use as-is
			hostname = registry
		} else {
			// Treat as hostname:port
			parts := strings.Split(registry, ":")
			hostname = parts[0]
		}
	}

	// SSRF protection: check for private IP ranges BEFORE regex validation
	// (some private IPs like ::1 may not match the registry pattern)
	if isPrivateIP(hostname) {
		// This is a warning, not an error (some users intentionally use private registries)
		return &SSRFWarning{
			Target:     registry,
			Message:    fmt.Sprintf("registry appears to be in a private IP range: %s", registry),
			Suggestion: "ensure this is intentional and you trust the target registry",
		}
	}

	// Validate hostname pattern
	if !registryPattern.MatchString(hostname) && !registryPattern.MatchString(registry) {
		return &ValidationError{
			Field:      "registry",
			Message:    fmt.Sprintf("invalid registry hostname format: %s", hostname),
			Suggestion: "use a valid domain name, IP address, or 'localhost'",
			ExitCode:   1,
		}
	}

	return nil
}

// isPrivateIP checks if a hostname resolves to a private IP address
func isPrivateIP(hostname string) bool {
	// Check localhost variants
	if hostname == "localhost" || hostname == "127.0.0.1" || hostname == "::1" {
		return true
	}

	// Try to resolve as IP
	ip := net.ParseIP(hostname)
	if ip == nil {
		// Try DNS resolution
		ips, err := net.LookupIP(hostname)
		if err != nil || len(ips) == 0 {
			return false
		}
		ip = ips[0]
	}

	// Check against private ranges
	for _, cidr := range privateIPRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}
