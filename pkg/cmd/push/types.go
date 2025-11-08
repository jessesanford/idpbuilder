package push

import "fmt"

// PushOptions configures the push command execution
type PushOptions struct {
	// ImageName is the local Docker image to push (e.g., "alpine:latest")
	ImageName string

	// Registry is the target registry URL (default: Gitea registry)
	Registry string

	// Username is the registry authentication username
	Username string

	// Password is the registry authentication password
	Password string

	// Insecure controls TLS certificate verification
	Insecure bool

	// Verbose enables detailed progress output
	Verbose bool
}

// Validate checks if PushOptions are valid and complete
func (o *PushOptions) Validate() error {
	if o.ImageName == "" {
		return fmt.Errorf("image name is required")
	}
	if o.Username == "" {
		return fmt.Errorf("username is required")
	}
	if o.Password == "" {
		return fmt.Errorf("password is required")
	}
	if o.Registry == "" {
		return fmt.Errorf("registry is required")
	}
	return nil
}
