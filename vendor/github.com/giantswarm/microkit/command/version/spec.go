package version

import (
	"github.com/spf13/cobra"
)

// Command represents the version command for any microservice.
type Command interface {
	// CobraCommand returns the actual cobra command for the version command.
	CobraCommand() *cobra.Command
	// Execute represents the cobra run method.
	Execute(cmd *cobra.Command, args []string)
}
