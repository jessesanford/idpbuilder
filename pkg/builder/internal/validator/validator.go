package validator

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidatePortFormat validates that a port specification follows the correct format.
// Valid formats: "port/protocol" (e.g., "80/tcp", "443/tcp", "53/udp")
func ValidatePortFormat(port string) error {
	parts := strings.Split(port, "/")
	if len(parts) != 2 {
		return fmt.Errorf("port must be in format 'port/protocol'")
	}
	
	// Validate port number
	portNum, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid port number: %s", parts[0])
	}
	
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port number must be between 1 and 65535: %d", portNum)
	}
	
	// Validate protocol
	protocol := strings.ToLower(parts[1])
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("protocol must be 'tcp' or 'udp': %s", protocol)
	}
	
	return nil
}

// ValidateUserFormat validates the user specification format.
// Valid formats: numeric UID, username, UID:GID, username:group
func ValidateUserFormat(user string) error {
	if user == "" {
		return nil // Empty is valid (defaults to root)
	}
	
	// Handle user:group format
	parts := strings.Split(user, ":")
	if len(parts) > 2 {
		return fmt.Errorf("user specification can have at most one colon")
	}
	
	// Validate user part
	userPart := parts[0]
	if userPart == "" {
		return fmt.Errorf("user part cannot be empty")
	}
	
	// Check if it's a numeric UID
	if _, err := strconv.Atoi(userPart); err != nil {
		// Not numeric, validate as username
		if !IsValidUsername(userPart) {
			return fmt.Errorf("invalid username format")
		}
	}
	
	// Validate group part if present
	if len(parts) == 2 {
		groupPart := parts[1]
		if groupPart == "" {
			return fmt.Errorf("group part cannot be empty")
		}
		
		// Check if it's a numeric GID
		if _, err := strconv.Atoi(groupPart); err != nil {
			// Not numeric, validate as group name
			if !IsValidUsername(groupPart) { // Same rules as username
				return fmt.Errorf("invalid group name format")
			}
		}
	}
	
	return nil
}

// IsValidUsername validates that a string is a valid Unix username.
// This is a basic validation - actual validation depends on the target system.
func IsValidUsername(name string) bool {
	if len(name) == 0 || len(name) > 32 {
		return false
	}
	
	// Must start with letter or underscore
	if !(name[0] >= 'a' && name[0] <= 'z') && !(name[0] >= 'A' && name[0] <= 'Z') && name[0] != '_' {
		return false
	}
	
	// Can contain letters, numbers, underscores, hyphens
	for _, char := range name[1:] {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || char == '_' || char == '-') {
			return false
		}
	}
	
	return true
}