// Package logger implements a logging interface used to log messages.
package micrologger

import (
	"context"
	"io"

	kitlog "github.com/go-kit/kit/log"

	"github.com/giantswarm/micrologger/loggermeta"
)

type Config struct {
	Caller             kitlog.Valuer
	IOWriter           io.Writer
	TimestampFormatter kitlog.Valuer
}

type MicroLogger struct {
	logger kitlog.Logger
}

func New(config Config) (*MicroLogger, error) {
	if config.Caller == nil {
		config.Caller = DefaultCaller
	}
	if config.TimestampFormatter == nil {
		config.TimestampFormatter = DefaultTimestampFormatter
	}
	if config.IOWriter == nil {
		config.IOWriter = DefaultIOWriter
	}

	kitLogger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(config.IOWriter))
	kitLogger = kitlog.With(
		kitLogger,
		"caller", config.Caller,
		"time", config.TimestampFormatter,
	)

	l := &MicroLogger{
		logger: kitLogger,
	}

	return l, nil
}

func (l *MicroLogger) Log(keyVals ...interface{}) error {
	return l.logger.Log(keyVals...)
}

func (l *MicroLogger) LogCtx(ctx context.Context, keyVals ...interface{}) error {
	meta, ok := loggermeta.FromContext(ctx)
	if !ok {
		return l.logger.Log(keyVals...)
	}

	var newKeyVals []interface{}
	{
		newKeyVals = append(newKeyVals, keyVals...)

		for k, v := range meta.KeyVals {
			newKeyVals = append(newKeyVals, k)
			newKeyVals = append(newKeyVals, v)
		}
	}

	return l.logger.Log(newKeyVals...)
}

func (l *MicroLogger) With(keyVals ...interface{}) Logger {
	return &MicroLogger{
		logger: kitlog.With(l.logger, keyVals...),
	}
}
