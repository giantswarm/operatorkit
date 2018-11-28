package microloggertest

import (
	"github.com/giantswarm/micrologger"
)

// New returns a Logger intance or panics if the creation fails.
func New() micrologger.Logger {
	c := micrologger.Config{}

	logger, err := micrologger.New(c)
	if err != nil {
		panic(err)
	}

	return logger
}
