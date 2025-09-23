// Package push contains configuration structures for the push command.
package push

import (
	"github.com/spf13/cobra"
)

// PushConfig represents the configuration for the push command
type PushConfig struct {
	RegistryURL string
	Username    string
	Password    string
	Namespace   string
	Dir         string
	Insecure    bool
	PlainHTTP   bool
}

// NewPushConfig creates a new PushConfig with default values
func NewPushConfig() *PushConfig {
	return &PushConfig{
		Namespace: "idpbuilder",
		Dir:       ".",
		Insecure:  false,
		PlainHTTP: false,
	}
}

// parseFlags extracts flag values and returns a populated PushConfig
func parseFlags(cmd *cobra.Command) (*PushConfig, error) {
	config := NewPushConfig()

	// Extract flag values
	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return nil, err
	}
	config.Username = username

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return nil, err
	}
	config.Password = password

	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return nil, err
	}
	config.Namespace = namespace

	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return nil, err
	}
	config.Dir = dir

	insecure, err := cmd.Flags().GetBool("insecure")
	if err != nil {
		return nil, err
	}
	config.Insecure = insecure

	plainHTTP, err := cmd.Flags().GetBool("plain-http")
	if err != nil {
		return nil, err
	}
	config.PlainHTTP = plainHTTP

	return config, nil
}
