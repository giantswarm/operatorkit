package transaction

import (
	"bytes"
	"net/http"

	"golang.org/x/net/context"
)

// ExecuteConfig is used to configure the Executer.
type ExecuteConfig struct {
	// Replay is the action being executed to replay a transaction result in case
	// a former call to Trial was successful.
	Replay func(context context.Context, v interface{}) error
	// ReplayDecoder is the decoder used to convert persisted trial outputs so
	// they can be consumed by replay functions. The underlying type of the
	// returned interface value is string.
	ReplayDecoder func(b []byte) (interface{}, error)
	// Trial is the action being executed to fulfil a transaction.
	Trial func(context context.Context) (interface{}, error)
	// TrialEncoder is the encoder used to convert created trial outputs so they
	// can be persisted.
	TrialEncoder func(v interface{}) ([]byte, error)
	// TrialID is an identifier scoped to the transaction ID obtained by the
	// context provided to transaction.Executer.Execute. The trial ID is used to
	// keep track of the state of the current transaction.
	TrialID string
}

// Executer provides a single transactional execution according to the provided
// configuration.
type Executer interface {
	// Execute actually executes the configured transaction. Transactions are
	// identified by the transaction ID obtained by the given context. In case the
	// given context contains a transaction ID, the configured trial associated
	// with the given trial ID will only be executed successfully once for the
	// current transaction ID. In case the trial fails it will be executed again
	// on the next call to Execute. In case the trial succeeded the given deplay
	// function, if any, will be executed.
	Execute(ctx context.Context, config ExecuteConfig) error
	// ExecuteConfig provides a default configuration for calls to Execute by best
	// effort.
	ExecuteConfig() ExecuteConfig
}

// Responder is able to reply to requests for which transactions have been
// tracked.
type Responder interface {
	// Exists checks whether a transaction response is stored under the given
	// transaction ID.
	Exists(ctx context.Context, transactionID string) (bool, error)
	// Reply uses the information onbtained by the given response replier to
	// create a response to reply to the current request.
	Reply(ctx context.Context, transactionID string, rr ResponseReplier) error
	// Track persists information obtained by the given response tracker to create
	// a transaction response. This can be used to reply to upcoming requests.
	Track(ctx context.Context, transactionID string, rt ResponseTracker) error
}

// Response is the transaction response object obtaining response relevant
// information used to track responses and reply to requests.
type Response struct {
	Body   string              `json:"body"`
	Code   int                 `json:"code"`
	Header map[string][]string `json:"header"`
}

// ResponseReplier is used to create response information to reply to requests.
type ResponseReplier interface {
	// Header is only a wrapper around http.ResponseWriter.Header.
	Header() http.Header
	// Write is only a wrapper around http.ResponseWriter.Write.
	Write(b []byte) (int, error)
	// WriteHeader is a wrapper around http.ResponseWriter.Write. In addition to
	// that it is used to track the written status code.
	WriteHeader(c int)
}

// ResponseTracker is used to track information about responses.
type ResponseTracker interface {
	// BodyBuffer returns the buffer which is used to track the bytes being
	// written to the response.
	BodyBuffer() *bytes.Buffer
	// Header is only a wrapper around http.ResponseWriter.Header.
	Header() http.Header
	// StatusCode returns either the default status code of the one that was
	// actually written using WriteHeader.
	StatusCode() int
}
