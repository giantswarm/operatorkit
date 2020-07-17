package sentry

import (
	"context"
	"encoding/json"
	"fmt"

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

func fromErr(err error) error {
	sfe := sentryFriendlyError{
		err:        err,
		stackTrace: []stackTraceEntry{},
	}

	var jsonErr *microerror.JSONError
	myerr := json.Unmarshal([]byte(microerror.JSON(err)), &jsonErr)
	if myerr != nil {
		// TODO print some output
		return err
	}

	stackTrace := sfe.stackTrace
	for _, entry := range jsonErr.Stack {
		fmt.Printf("%v %s:%d\n", entry.ProgramCounter, entry.File, entry.Line)
		stackTrace = append(stackTrace, stackTraceEntry{entry.ProgramCounter})
	}
	sfe.stackTrace = stackTrace

	return microerror.Mask(sfe)
}

func (s *Default) Capture(ctx context.Context, err error) {
	sentry.CaptureException(fromErr(err))
}
