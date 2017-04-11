package flag

import (
	"github.com/giantswarm/microkit/command/daemon/flag/config"
	"github.com/giantswarm/microkit/command/daemon/flag/server"
	"github.com/giantswarm/microkit/flag"
)

type Flag struct {
	Config config.Config
	Server server.Server
}

func New() Flag {
	f := Flag{}
	flag.Init(&f)
	return f
}
