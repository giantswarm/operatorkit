package server

import (
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/juju/errgo"
)

func errorDomain(err error) string {
	switch e := err.(type) {
	case kithttp.Error:
		switch e.Domain {
		case kithttp.DomainEncode:
			return "encode"
		case kithttp.DomainDecode:
			return "decode"
		case kithttp.DomainDo:
			return "domain"
		}
	}
	return "server"
}

func errorMessage(err error) string {
	if err == nil {
		return ""
	}
	switch kitErr := err.(type) {
	case kithttp.Error:
		switch errgoErr := kitErr.Err.(type) {
		case *errgo.Err:
			return errgoErr.Error()
		}
	}
	return err.Error()
}

func errorTrace(err error) string {
	switch kitErr := err.(type) {
	case kithttp.Error:
		switch errgoErr := kitErr.Err.(type) {
		case *errgo.Err:
			return errgoErr.GoString()
		}
	}
	return "n/a"
}

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var invalidContextError = errgo.New("invalid context")

// IsInvalidContext asserts invalidContextError.
func IsInvalidContext(err error) bool {
	return errgo.Cause(err) == invalidContextError
}

var invalidTransactionIDError = errgo.New("invalid transaction ID")

// IsInvalidTransactionID asserts invalidTransactionIDError.
func IsInvalidTransactionID(err error) bool {
	return errgo.Cause(err) == invalidTransactionIDError
}
