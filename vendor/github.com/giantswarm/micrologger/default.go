package micrologger

import (
	"fmt"
	"os"
	"time"

	"github.com/go-stack/stack"
)

var DefaultCaller = func() interface{} {
	return fmt.Sprintf("%+v", stack.Caller(4))
}

var DefaultIOWriter = os.Stdout

var DefaultTimestampFormatter = func() interface{} {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999999-07:00")
}
