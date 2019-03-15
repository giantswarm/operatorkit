package micrologger

import (
	"context"

	"github.com/giantswarm/microerror"
)

const (
	KeyLevel     = "level"
	KeyVerbosity = "verbosity"
)

const (
	levelDebug levelID = 1 << iota
	levelInfo
	levelWarning
	levelError
)

var (
	levelMapping = map[string]levelID{
		"debug":   levelDebug,
		"info":    levelInfo,
		"warning": levelWarning,
		"error":   levelError,
	}
)

type levelID byte

type ActivationLoggerConfig struct {
	Underlying Logger

	Activations map[string]interface{}
}

type activationLogger struct {
	underlying Logger

	activations map[string]interface{}
}

// NewActivation creates a new activation key logger. This logger kind can be
// used on command line tools to improve situations in which log filtering using
// other command line tools like grep is not sufficient. Due to certain filter
// mechanisms this Logger implementation should not be used in performance
// critical applications. The idea of the activation key logger is to have a
// multi dimensional log filter mechanism. This logger here provides three
// different features which can be combined and used simultaneously at will.
//
//     Filtering arbitrary key-value pairs. The structured nature of the Logger
//     interface expects key-value pairs to be logged. The activation key logger
//     can be configured with any kind of activation key-pairs which, when
//     configured, all have to match against an emitted logging call, in order
//     to be dispatched. In case none, or not all activation keys match, the
//     emitted logging call is going to be ignored.
//
//     Filtering log levels works using the special log levels debug, info,
//     warning and error. The level based nature of this activation mechanism is
//     that lower log levels match just like exact log levels match. When the
//     Logger is configured to activate on info log levels, the Logger will
//     activate on debug related logs, as well as info related logs, but not on
//     warning or error related logs.
//
//     Filtering log verbosity works similar like the log level mechanism, but
//     on arbitrary verbosity levels, which are represented as numbers. As long
//     as the configured verbosity is higher or equal to the perceived verbosity
//     obtained by the emitted logging call, the log will be dispatched.
//
func NewActivation(config ActivationLoggerConfig) (Logger, error) {
	if config.Underlying == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Underlying must not be empty", config)
	}

	l := &activationLogger{
		underlying: config.Underlying,

		activations: config.Activations,
	}

	return l, nil
}

func (l *activationLogger) Log(keyVals ...interface{}) error {
	activated, err := shouldActivate(l.activations, keyVals)
	if err != nil {
		return microerror.Mask(err)
	}

	if activated {
		return l.underlying.Log(keyVals...)
	}

	return nil
}

func (l *activationLogger) LogCtx(ctx context.Context, keyVals ...interface{}) error {
	activated, err := shouldActivate(l.activations, keyVals)
	if err != nil {
		return microerror.Mask(err)
	}

	if activated {
		return l.underlying.LogCtx(ctx, keyVals...)
	}

	return nil
}

func (l *activationLogger) With(keyVals ...interface{}) Logger {
	return l.underlying.With(keyVals...)
}

func valueFor(keyVals []interface{}, key string) (interface{}, bool) {
	for i := 1; i < len(keyVals); i += 2 {
		if key == keyVals[i-1] {
			return keyVals[i], true
		}
	}

	return nil, false
}

func isLevelAllowed(keyVals []interface{}, aVal interface{}) bool {
	s, ok := aVal.(string)
	if !ok {
		return false
	}
	activationLevel, ok := levelMapping[s]
	if !ok {
		return false
	}

	for i := 0; i < len(keyVals); i += 2 {
		k, ok := keyVals[i].(string)
		if !ok {
			continue
		}
		if k != KeyLevel {
			continue
		}
		v, ok := keyVals[i+1].(string)
		if !ok {
			continue
		}
		keyValsLevel, ok := levelMapping[v]
		if !ok {
			continue
		}

		return activationLevel >= keyValsLevel
	}

	return false
}

func isVerbosityAllowed(keyVals []interface{}, aVal interface{}) bool {
	activationVerbosity, ok := aVal.(int)
	if !ok {
		return false
	}

	for i := 0; i < len(keyVals); i += 2 {
		k, ok := keyVals[i].(string)
		if !ok {
			continue
		}
		if k != KeyVerbosity {
			continue
		}
		keyValsVerbosity, ok := keyVals[i+1].(int)
		if !ok {
			continue
		}

		return activationVerbosity >= keyValsVerbosity
	}

	return false
}

func shouldActivate(activations map[string]interface{}, keyVals []interface{}) (bool, error) {
	var activationCount int

	for aKey, aVal := range activations {
		v, ok := valueFor(keyVals, aKey)
		if ok && v == aVal {
			activationCount++
			continue
		}
		if aKey == KeyLevel && isLevelAllowed(keyVals, aVal) {
			activationCount++
			continue
		}
		if aKey == KeyVerbosity && isVerbosityAllowed(keyVals, aVal) {
			activationCount++
			continue
		}
	}

	if len(activations) != 0 && len(activations) == activationCount {
		return true, nil
	}

	return false, nil
}
