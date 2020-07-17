package sentry

import (
	"context"
	"encoding/json"

	"github.com/getsentry/sentry-go"
	"github.com/giantswarm/microerror"
)

type Default struct {
}

type sentryFriendlyError struct {
	err        error
	stackTrace []stackTraceEntry
}

type stackTraceEntry struct {
	ProgramCounter uintptr
}

func (e sentryFriendlyError) StackTrace() []stackTraceEntry {
	return e.stackTrace
}

func (e sentryFriendlyError) Error() string {
	return e.err.Error()
}

func fromErr(err error) sentryFriendlyError {
	return sentryFriendlyError{
		err:        err,
		stackTrace: []stackTraceEntry{},
	}
}

func (s *Default) Capture(ctx context.Context, err error) {
	var jsonErr *microerror.JSONError
	myerr := json.Unmarshal([]byte(microerror.JSON(err)), &jsonErr)
	if myerr != nil {

		return
	}

	e := fromErr(err)
	stackTrace := e.stackTrace
	for _, entry := range jsonErr.Stack {
		stackTrace = append(stackTrace, stackTraceEntry{entry.ProgramCounter})
	}
	e.stackTrace = stackTrace

	sentry.CaptureException(e)
}
