package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// OCI image name validation regex patterns
var (
	// imageNamePattern validates full OCI image references
	// Format: [registry/][namespace/]repository[:tag][@digest]
	imageNamePattern = regexp.MustCompile(`^([a-zA-Z0-9._-]+([.:][0-9]+)?/)?([a-zA-Z0-9._/-]+)(:[a-zA-Z0-9._-]+)?(@sha256:[a-fA-F0-9]{64})?$`)

	// tagPattern validates image tags
	tagPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

	// digestPattern validates SHA256 digests
	digestPattern = regexp.MustCompile(`^sha256:[a-fA-F0-9]{64}$`)
)

// Dangerous characters that could be used for command injection
const dangerousChars = ";|&$`<>(){}[]'\"\\"

// ValidateImageName validates an OCI image name format
//
// Validation Rules:
// 1. Non-empty string required
// 2. No shell metacharacters (command injection prevention)
// 3. Must match OCI image name pattern
// 4. Tag must be valid if present
// 5. Digest must be valid if present
//
// Returns:
//   - nil if valid
//   - ValidationError with actionable suggestion if invalid
func ValidateImageName(imageName string) error {
	if imageName == "" {
		return &ValidationError{
			Field:      "image-name",
			Message:    "image name cannot be empty",
			Suggestion: "provide an image name like 'alpine:latest' or 'docker.io/library/ubuntu:22.04'",
			ExitCode:   1,
		}
	}

	// Check for dangerous characters (command injection prevention)
	if containsAnyChar(imageName, dangerousChars) {
		return &ValidationError{
			Field:      "image-name",
			Message:    fmt.Sprintf("image name contains shell metacharacters: %s", imageName),
			Suggestion: "use only alphanumeric characters, dots, hyphens, underscores, colons, and slashes",
			ExitCode:   1,
		}
	}

	// Validate digest if present (before regex to catch invalid digests)
	if strings.Contains(imageName, "@") {
		digestPart := strings.Split(imageName, "@")[1]
		if !digestPattern.MatchString(digestPart) {
			return &ValidationError{
				Field:      "image-digest",
				Message:    fmt.Sprintf("invalid digest format: %s", digestPart),
				Suggestion: "digest must be in format: sha256:[64 hex characters]",
				ExitCode:   1,
			}
		}
	}

	// Validate against OCI image name pattern
	if !imageNamePattern.MatchString(imageName) {
		return &ValidationError{
			Field:      "image-name",
			Message:    fmt.Sprintf("invalid OCI image name format: %s", imageName),
			Suggestion: "use format: [registry/][namespace/]repository[:tag][@digest]",
			ExitCode:   1,
		}
	}

	// Validate tag if present (must be after last slash and before @ if present)
	// Need to handle registry:port properly - only check tag after last /
	if strings.Contains(imageName, ":") && !strings.Contains(imageName, "@") {
		// Find the last colon after the last slash
		lastSlash := strings.LastIndex(imageName, "/")
		colonPos := strings.LastIndex(imageName, ":")

		// Only validate tag if colon is after the last slash
		if colonPos > lastSlash {
			tag := imageName[colonPos+1:]
			if !tagPattern.MatchString(tag) {
				return &ValidationError{
					Field:      "image-tag",
					Message:    fmt.Sprintf("invalid tag format: %s", tag),
					Suggestion: "tags must contain only alphanumeric characters, dots, hyphens, and underscores",
					ExitCode:   1,
				}
			}
		}
	}

	return nil
}

// containsAnyChar checks if a string contains any character from a set
func containsAnyChar(s, chars string) bool {
	for _, c := range chars {
		if strings.ContainsRune(s, c) {
			return true
		}
	}
	return false
}
