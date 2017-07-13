// Package logger implements a logging interface used to log messages.
package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-stack/stack"

	microerror "github.com/giantswarm/microkit/error"
)

// Config represents the configuration used to create a new logger.
type Config struct {
	// Settings.
	Caller             kitlog.Valuer
	IOWriter           io.Writer
	TimestampFormatter kitlog.Valuer
}

// DefaultConfig provides a default configuration to create a new logger by best
// effort.
func DefaultConfig() Config {

	return Config{
		// Settings.
		Caller: func() interface{} {
			return fmt.Sprintf("%+v", stack.Caller(4))
		},
		IOWriter: ioutil.Discard,
		TimestampFormatter: func() interface{} {
			return time.Now().UTC().Format("06-01-02 15:04:05.000")
		},
	}
}

// New creates a new configured logger.
func New(config Config) (Logger, error) {
	// Settings.
	if config.Caller == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "caller must not be empty")
	}
	if config.TimestampFormatter == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "timestamp formatter must not be empty")
	}

	kitLogger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(config.IOWriter))
	kitLogger = kitlog.With(
		kitLogger,
		"caller", config.Caller,
		"time", config.TimestampFormatter,
	)

	newLogger := &logger{
		Logger: kitLogger,
	}

	return newLogger, nil
}

type logger struct {
	Logger kitlog.Logger
}

func (l *logger) Log(keyvals ...interface{}) error {
	return l.Logger.Log(keyvals...)
}
