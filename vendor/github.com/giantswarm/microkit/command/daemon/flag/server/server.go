package server

import (
	"github.com/giantswarm/microkit/command/daemon/flag/server/listen"
	"github.com/giantswarm/microkit/command/daemon/flag/server/tls"
)

type Server struct {
	Listen listen.Listen
	TLS    tls.TLS
}
