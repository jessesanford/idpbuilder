// Package push - test harness for running integration tests
package push

import "github.com/spf13/cobra"

// GetCommand returns the push command for testing
func GetCommand() *cobra.Command {
	return pushCmd
}

// AddToParent adds the push command to a parent command
func AddToParent(parent *cobra.Command) {
	parent.AddCommand(pushCmd)
}