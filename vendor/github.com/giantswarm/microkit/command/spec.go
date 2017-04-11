package command

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/microkit/command/daemon"
	"github.com/giantswarm/microkit/command/version"
)

// Command represents the root command for any microservice.
type Command interface {
	// CobraCommand returns the actual cobra command for the root command.
	CobraCommand() *cobra.Command
	// DaemonCommand returns the daemon sub command.
	DaemonCommand() daemon.Command
	// Execute represents the cobra run method.
	Execute(cmd *cobra.Command, args []string)
	// VersionCommand returns the version sub command.
	VersionCommand() version.Command
}
