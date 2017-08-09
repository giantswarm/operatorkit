package microloggertest

import (
	"io/ioutil"

	"github.com/giantswarm/micrologger"
)

// New returns a Logger intance configured to discard its output.
func New() micrologger.Logger {
	c := micrologger.DefaultConfig()
	{
		c.IOWriter = ioutil.Discard
	}

	logger, err := micrologger.New(c)
	if err != nil {
		panic(err)
	}

	return logger
}
