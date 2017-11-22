// Package logger implements a logging interface used to log messages.
package micrologger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-stack/stack"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/loggercontext"
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
		IOWriter: os.Stdout,
		TimestampFormatter: func() interface{} {
			return time.Now().UTC().Format("2006-01-02 15:04:05.000")
		},
	}
}

// New creates a new configured logger.
func New(config Config) (Logger, error) {
	// Settings.
	if config.Caller == nil {
		return nil, microerror.Maskf(invalidConfigError, "caller must not be empty")
	}
	if config.TimestampFormatter == nil {
		return nil, microerror.Maskf(invalidConfigError, "timestamp formatter must not be empty")
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

func (l *logger) Log(keyVals ...interface{}) error {
	return l.Logger.Log(keyVals...)
}

func (l *logger) LogWithCtx(ctx context.Context, keyVals ...interface{}) error {
	container, ok := loggercontext.FromContext(ctx)
	if !ok {
		return l.Logger.Log(keyVals...)
	}

	var newKeyVals []interface{}
	{
		newKeyVals = append(newKeyVals, keyVals...)

		for k, v := range container.KeyVals {
			newKeyVals = append(newKeyVals, k)
			newKeyVals = append(newKeyVals, v)
		}
	}

	return l.Logger.Log(newKeyVals...)
}

func (l *logger) With(keyVals ...interface{}) Logger {
	return &logger{
		Logger: kitlog.With(l.Logger, keyVals...),
	}
}
