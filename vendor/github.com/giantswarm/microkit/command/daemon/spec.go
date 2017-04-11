package daemon

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/giantswarm/microkit/server"
)

type ServerFactory func(v *viper.Viper) server.Server

// Command represents the daemon command for any microservice.
type Command interface {
	// CobraCommand returns the actual cobra command for the daemon command.
	CobraCommand() *cobra.Command
	// Execute represents the cobra run method.
	Execute(cmd *cobra.Command, args []string)
}
