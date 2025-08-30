package delete

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
)

// ConfirmationLevel represents different levels of confirmation required
type ConfirmationLevel int

const (
	// BasicConfirmation requires simple y/N confirmation
	BasicConfirmation ConfirmationLevel = iota
	// TypedConfirmation requires typing the resource name
	TypedConfirmation
	// DangerousConfirmation requires typing a confirmation phrase
	DangerousConfirmation
)

// ConfirmationRequest represents a confirmation request
type ConfirmationRequest struct {
	Level         ConfirmationLevel
	Message       string
	ResourceName  string
	WarningCount  int
	RequiredPhrase string
}

// PromptForConfirmation prompts the user for confirmation based on the request
func PromptForConfirmation(request ConfirmationRequest) (bool, error) {
	// Show warnings if any
	if request.WarningCount > 0 {
		helpers.PrintWarning("This operation will affect %d resources", request.WarningCount)
	}

	switch request.Level {
	case BasicConfirmation:
		return promptBasicConfirmation(request.Message)
	case TypedConfirmation:
		return promptTypedConfirmation(request.Message, request.ResourceName)
	case DangerousConfirmation:
		return promptDangerousConfirmation(request.Message, request.RequiredPhrase)
	default:
		return false, fmt.Errorf("unknown confirmation level: %d", request.Level)
	}
}

// promptBasicConfirmation prompts for a basic y/N confirmation
func promptBasicConfirmation(message string) (bool, error) {
	fmt.Printf("%s [y/N]: ", message)
	
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

// promptTypedConfirmation requires typing the resource name for confirmation
func promptTypedConfirmation(message, resourceName string) (bool, error) {
	fmt.Printf("%s\nType the resource name '%s' to confirm: ", message, resourceName)
	
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(response)
	if response != resourceName {
		helpers.PrintError("Input '%s' does not match expected resource name '%s'", response, resourceName)
		return false, nil
	}

	return true, nil
}

// promptDangerousConfirmation requires typing a specific phrase for dangerous operations
func promptDangerousConfirmation(message, requiredPhrase string) (bool, error) {
	helpers.PrintWarning("DANGER: This is a destructive operation!")
	fmt.Printf("%s\nType '%s' to confirm: ", message, requiredPhrase)
	
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(response)
	if response != requiredPhrase {
		helpers.PrintError("Input does not match required confirmation phrase")
		return false, nil
	}

	return true, nil
}

// DetermineConfirmationLevel determines the appropriate confirmation level
func DetermineConfirmationLevel(request *DeleteRequest) ConfirmationLevel {
	// Dangerous operations require dangerous confirmation
	if request.ResourceType == "all" || 
	   request.AllNamespaces || 
	   (request.Selector != "" && len(request.ResourceNames) == 0) {
		return DangerousConfirmation
	}

	// Multiple specific resources require typed confirmation
	if len(request.ResourceNames) > 1 {
		return TypedConfirmation
	}

	// Single resource requires basic confirmation
	return BasicConfirmation
}

// BuildConfirmationRequest builds a confirmation request from a delete request
func BuildConfirmationRequest(deleteReq *DeleteRequest) ConfirmationRequest {
	level := DetermineConfirmationLevel(deleteReq)
	
	var message strings.Builder
	message.WriteString(fmt.Sprintf("Delete %s", deleteReq.ResourceType))
	
	if len(deleteReq.ResourceNames) > 0 {
		message.WriteString(fmt.Sprintf(" '%s'", strings.Join(deleteReq.ResourceNames, "', '")))
	}
	
	if deleteReq.Namespace != "" {
		message.WriteString(fmt.Sprintf(" in namespace '%s'", deleteReq.Namespace))
	} else if deleteReq.AllNamespaces {
		message.WriteString(" in all namespaces")
	}
	
	if deleteReq.Selector != "" {
		message.WriteString(fmt.Sprintf(" matching selector '%s'", deleteReq.Selector))
	}

	warningCount := len(deleteReq.ResourceNames)
	if deleteReq.ResourceType == "all" {
		warningCount = 10 // Estimated for "all" resources
	}

	return ConfirmationRequest{
		Level:         level,
		Message:       message.String(),
		ResourceName:  strings.Join(deleteReq.ResourceNames, " "),
		WarningCount:  warningCount,
		RequiredPhrase: "delete all resources",
	}
}