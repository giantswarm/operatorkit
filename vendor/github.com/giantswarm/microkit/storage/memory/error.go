package memory

import (
	"github.com/juju/errgo"
)

var notFoundError = errgo.New("not found")

func IsNotFound(err error) bool {
	return errgo.Cause(err) == notFoundError
}
